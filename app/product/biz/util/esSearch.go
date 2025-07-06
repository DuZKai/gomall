package util

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/elastic/go-elasticsearch/v7"
	"github.com/elastic/go-elasticsearch/v7/esapi"
	"io"
	"log"
	"strings"
)

// PageParams holds pagination parameters
type PageParams struct {
	PageNo   int
	PageSize int
}

// SearchCourseParamDto holds search filter and sort parameters
type SearchCourseParamDto struct {
	Keywords string
	Mt       string
	St       string
	Grade    string
	SortType string // "1" for price ASC, "2" for price DESC
}

type _hit struct {
	Source    CourseIndex         `json:"_source"`
	Highlight map[string][]string `json:"highlight"`
}

// CourseIndex represents the course document stored in ES
type CourseIndex struct {
	ID            int64           `json:"id"`
	Name          string          `json:"name"`
	Grade         string          `json:"grade"`
	Mt            string          `json:"mt"`
	St            string          `json:"st"`
	Charge        string          `json:"charge"`
	Pic           string          `json:"pic"`
	Price         float64         `json:"price"`
	OriginalPrice float64         `json:"originalPrice"`
	Teachmode     string          `json:"teachmode"`
	ValidDays     int             `json:"validDays"`
	CreateDate    string          `json:"createDate"`
	CompanyName   string          `json:"companyName"`
	IsAd          string          `json:"isAd"`         // 改成 string
	RawTeachers   json.RawMessage `json:"teacherNames"` // 原始字符串
	RawTags       json.RawMessage `json:"tags"`
	TeacherNames  []string
	Tags          string
}

// SearchPageResultDto is the paginated result
type SearchPageResultDto struct {
	List     []CourseIndex
	Total    int64
	PageNo   int
	PageSize int
}

// QueryCoursePubNewIndex searches the course-publish index with given params
func QueryCoursePubNewIndex(ctx context.Context, client *elasticsearch.Client, pageParams PageParams, courseSearchParam SearchCourseParamDto, index string, sourceFields string) (SearchPageResultDto, error) {
	// Build the search body
	from := (pageParams.PageNo - 1) * pageParams.PageSize

	// Prepare the query map
	query := map[string]interface{}{
		"bool": map[string]interface{}{
			"must":   []interface{}{},
			"filter": []interface{}{},
		},
	}
	log.Printf("query: %+v", query)

	// Keywords multi-match
	if courseSearchParam.Keywords != "" {
		multi := map[string]interface{}{
			"multi_match": map[string]interface{}{
				"query":                courseSearchParam.Keywords,
				"fields":               []string{"name^2", "description"},
				"minimum_should_match": "75%",
			},
		}
		query["bool"].(map[string]interface{})["must"] = append(
			query["bool"].(map[string]interface{})["must"].([]interface{}),
			multi,
		)
	}

	// Term filters
	filters := query["bool"].(map[string]interface{})["filter"].([]interface{})
	if courseSearchParam.Mt != "" {
		filters = append(filters, map[string]interface{}{"term": map[string]interface{}{"mt.keyword": courseSearchParam.Mt}})
	}
	if courseSearchParam.St != "" {
		filters = append(filters, map[string]interface{}{"term": map[string]interface{}{"st.keyword": courseSearchParam.St}})
	}
	if courseSearchParam.Grade != "" {
		filters = append(filters, map[string]interface{}{"term": map[string]interface{}{"grade": courseSearchParam.Grade}})
	}
	query["bool"].(map[string]interface{})["filter"] = filters

	// Function score to boost ads
	functionScore := map[string]interface{}{
		"function_score": map[string]interface{}{
			"query": query,
			"functions": []interface{}{
				map[string]interface{}{
					"filter": map[string]interface{}{"term": map[string]interface{}{"isAd": 800002}},
					"weight": 10,
				},
			},
		},
	}

	strArray := strings.Split(sourceFields, ",")

	// Assemble the final body
	body := map[string]interface{}{
		"from":    from,
		"size":    pageParams.PageSize,
		"_source": strArray,
		"query":   functionScore,
		"highlight": map[string]interface{}{
			"pre_tags":  []string{"<font class='eslight'>"},
			"post_tags": []string{"</font>"},
			"fields":    map[string]interface{}{"name": map[string]interface{}{}},
		},
	}

	// Sort by price if required
	if courseSearchParam.SortType == "1" {
		body["sort"] = []interface{}{map[string]interface{}{"price": map[string]interface{}{"order": "asc"}}}
	} else if courseSearchParam.SortType == "2" {
		body["sort"] = []interface{}{map[string]interface{}{"price": map[string]interface{}{"order": "desc"}}}
	}

	// Encode body
	buf, err := json.Marshal(body)
	if err != nil {
		return SearchPageResultDto{}, err
	}

	// Perform search
	req := esapi.SearchRequest{
		Index: []string{index},
		Body:  io.NopCloser(bytes.NewReader(buf)),
	}

	res, err := req.Do(ctx, client)
	if err != nil {
		return SearchPageResultDto{}, err
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			panic(err)
		}
	}(res.Body)
	// log.Printf("search response: %+v", res)

	// Parse response
	var r struct {
		Hits struct {
			Total struct {
				Value int64 `json:"value"`
			} `json:"total"`
			Hits []_hit `json:"hits"`
		} `json:"hits"`
	}

	// 反序列化
	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		return SearchPageResultDto{}, err
	}

	list := make([]CourseIndex, 0, len(r.Hits.Hits))
	for _, hit := range r.Hits.Hits {
		ci := hit.Source

		// 解析 teacherNames
		var raw string
		if err := json.Unmarshal(hit.Source.RawTeachers, &raw); err != nil {
			log.Fatalf("outer unmarshal failed: %v", err)
		}
		if err := json.Unmarshal([]byte(raw), &ci.TeacherNames); err != nil {
			log.Fatalf("inner unmarshal failed: %v", err)
		}

		// 解析 tags
		if err := json.Unmarshal(hit.Source.RawTags, &ci.Tags); err != nil {
			log.Fatalf("unmarshal failed: %v", err)
		}

		// 高亮覆盖
		if hs, ok := hit.Highlight["name"]; ok && len(hs) > 0 {
			ci.Name = hs[0]
		}
		list = append(list, ci)
	}

	return SearchPageResultDto{
		List:     list,
		Total:    r.Hits.Total.Value,
		PageNo:   pageParams.PageNo,
		PageSize: pageParams.PageSize,
	}, nil
}

// AddCourseIndex 在指定索引 indexName 中写入 ID=id 的文档 doc。
// 返回 true 代表 result 为 "created" 或 "updated"（ES 的幂等语义），否则 false。
func AddCourseIndex(ctx context.Context, client *elasticsearch.Client,
	indexName, id string, doc interface{}) (bool, error) {

	// 序列化文档
	data, err := json.Marshal(doc)
	if err != nil {
		return false, err
	}

	// 构建 IndexRequest
	req := esapi.IndexRequest{
		Index:      indexName,
		DocumentID: id,
		Body:       bytes.NewReader(data),
		// 如果希望写完立即可见可加 Refresh:"true"
	}

	res, err := req.Do(ctx, client)
	if err != nil {
		return false, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			panic(err)
		}
	}(res.Body)

	// 解析响应，只关心 result 字段
	var resp struct {
		Result string `json:"result"`
	}
	if err := json.NewDecoder(res.Body).Decode(&resp); err != nil {
		return false, err
	}
	return resp.Result == "created" || resp.Result == "updated", nil
}

// UpdateCourseIndex 更新指定索引 indexName 中 ID=id 的文档 doc。
// 仅当 ES 返回 "updated" 才视为成功。
func UpdateCourseIndex(ctx context.Context, client *elasticsearch.Client,
	indexName, id string, doc interface{}) (bool, error) {

	// ES Update API 需要 {"doc":{…}}
	body := map[string]interface{}{"doc": doc}
	data, err := json.Marshal(body)
	if err != nil {
		return false, err
	}

	req := esapi.UpdateRequest{
		Index:      indexName,
		DocumentID: id,
		Body:       bytes.NewReader(data),
	}

	res, err := req.Do(ctx, client)
	if err != nil {
		return false, err
	}
	defer res.Body.Close()

	var resp struct {
		Result string `json:"result"`
	}
	if err := json.NewDecoder(res.Body).Decode(&resp); err != nil {
		return false, err
	}
	return resp.Result == "updated", nil
}

// DeleteCourseIndex 删除指定索引 indexName 中 ID=id 的文档。
// 当 result=="deleted" 表示成功。
func DeleteCourseIndex(ctx context.Context, client *elasticsearch.Client,
	indexName, id string) (bool, error) {

	req := esapi.DeleteRequest{
		Index:      indexName,
		DocumentID: id,
	}

	res, err := req.Do(ctx, client)
	if err != nil {
		return false, err
	}
	defer res.Body.Close()

	var resp struct {
		Result string `json:"result"`
	}
	if err := json.NewDecoder(res.Body).Decode(&resp); err != nil {
		return false, err
	}
	return resp.Result == "deleted", nil
}

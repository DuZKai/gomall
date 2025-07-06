package util

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/elastic/go-elasticsearch/v7/esapi"
	"gomall/app/product/biz/dal/es"
	"io"
	"net/http"
)

// GetDocument 查询文档
func GetDocument(ctx context.Context, index, docID string) (map[string]interface{}, error) {
	req := esapi.GetRequest{
		Index:      index,
		DocumentID: docID,
	}
	client := es.ESClient

	res, err := req.Do(ctx, client)
	if err != nil {
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(res.Body)

	if res.StatusCode == 404 {
		return nil, fmt.Errorf("Document %s/%s not found", index, docID)
	} else if res.IsError() {
		body, _ := io.ReadAll(res.Body)
		return nil, fmt.Errorf("ES error [%d %s]: %s",
			res.StatusCode,
			http.StatusText(res.StatusCode),
			string(body),
		)
	}

	var result map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return nil, err
	}
	return result["_source"].(map[string]interface{}), nil
}

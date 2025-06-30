package service

import (
	"context"
	"encoding/json"
	"fmt"
	"gomall/app/product/biz/dal/mysql"
	"gomall/app/product/biz/dal/redis"
	"gomall/app/product/biz/model"
	product "gomall/rpc_gen/kitex_gen/product"
	"time"
)

type ListProductsService struct {
	ctx context.Context
} // NewListProductsService new ListProductsService
func NewListProductsService(ctx context.Context) *ListProductsService {
	return &ListProductsService{ctx: ctx}
}

// Run create note info
func (s *ListProductsService) Run(req *product.ListProductsReq) (resp *product.ListProductsResp, err error) {
	categoryQuery := model.NewCategoryQuery(s.ctx, mysql.DB)
	resp = &product.ListProductsResp{}

	// 构造 Redis key（唯一标识该请求参数）
	cacheKey := fmt.Sprintf("list_products:%s:%d:%d", req.CategoryName, req.Page, req.PageSize)

	// 尝试从 Redis 中读取缓存
	val, err := redis.RedisClient.Get(s.ctx, cacheKey).Result()
	if err == nil && val != "" {
		// 命中缓存，反序列化返回
		if err := json.Unmarshal([]byte(val), resp); err == nil {
			return resp, nil
		}
		// 如果反序列化失败，也继续查数据库
	}

	// ② 未命中缓存或反序列化失败，查询数据库
	categories, err := categoryQuery.GetProductsByCategoryNameAndPage(req.CategoryName, req.Page, req.PageSize)
	if err != nil {
		return nil, err
	}

	for _, v1 := range categories {
		for _, v := range v1.Products {
			resp.Products = append(resp.Products, &product.Product{
				Id:          uint32(v.ID),
				Name:        v.Name,
				Description: v.Description,
				Picture:     v.Picture,
				Price:       v.Price,
			})
		}
	}

	// ③ 写入 Redis 缓存（设置过期时间）
	data, _ := json.Marshal(resp)
	_ = redis.RedisClient.Set(s.ctx, cacheKey, data, time.Minute*10).Err()

	return resp, nil
}

package service

import (
	"context"
	"fmt"
	"gomall/app/product/biz/dal/mysql"
	"gomall/app/product/biz/dal/redis"
	"gomall/app/product/biz/model"
	product "gomall/rpc_gen/kitex_gen/product"
	"log"
)

type UpdateProductService struct {
	ctx context.Context
} // NewUpdateProductService new UpdateProductService
func NewUpdateProductService(ctx context.Context) *UpdateProductService {
	return &UpdateProductService{ctx: ctx}
}

// Run create note info
func (s *UpdateProductService) Run(req *product.UpdateProductReq) (resp *product.UpdateProductResp, err error) {
	// 直接构建更新模型
	updateData := &model.Product{
		Name:        req.Name,
		Description: req.Description,
		Picture:     req.Picture,
		Price:       req.Price,
	}

	// 使用链式调用
	query := model.NewProductQuery(s.ctx, mysql.DB)
	if err := query.UpdateProduct(int(req.Id), updateData); err != nil {
		return nil, err
	}

	// 获取更新后的记录
	finalP, err := query.GetById(int(req.Id))
	if err != nil {
		return nil, err
	}

	// 直接返回结果
	resp = &product.UpdateProductResp{
		Product: &product.Product{
			Id:          uint32(finalP.ID),
			Name:        finalP.Name,
			Description: finalP.Description,
			Picture:     finalP.Picture,
			Price:       finalP.Price, // 确保包含价格字段
		},
	}

	// 更新成功后删除相关缓存
	if err := s.deleteProductListCache(); err != nil {
		log.Printf("警告: 缓存清理失败 - %v", err)
		// 不返回错误，因为主操作已成功
	}
	return resp, nil
}

// 删除所有匹配的缓存键
func (s *UpdateProductService) deleteProductListCache() error {
	// 定义匹配模式（根据您的实际键结构调整）
	pattern := "list_products:*"

	// 使用 SCAN 迭代查找所有匹配的键
	var cursor uint64
	var keys []string
	for {
		var err error
		// 一次扫描 100 个键（可调整）
		keys, cursor, err = redis.RedisClient.Scan(s.ctx, cursor, pattern, 100).Result()
		if err != nil {
			return fmt.Errorf("扫描缓存键失败: %w", err)
		}

		// 删除当前批次的键
		if len(keys) > 0 {
			if err := redis.RedisClient.Del(s.ctx, keys...).Err(); err != nil {
				return fmt.Errorf("删除缓存键失败: %w", err)
			}
			log.Printf("已删除缓存键: %v", keys)
		}

		// 扫描完成
		if cursor == 0 {
			break
		}
	}
	return nil
}

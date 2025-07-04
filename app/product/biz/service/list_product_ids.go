package service

import (
	"context"
	"gomall/app/product/biz/dal/mysql"
	"gomall/app/product/biz/model"
	product "gomall/rpc_gen/kitex_gen/product"
)

type ListProductIdsService struct {
	ctx context.Context
} // NewListProductIdsService new ListProductIdsService
func NewListProductIdsService(ctx context.Context) *ListProductIdsService {
	return &ListProductIdsService{ctx: ctx}
}

// Run create note info
func (s *ListProductIdsService) Run() (resp *product.ListProductIdsResp, err error) {
	productQuery := model.NewProductQuery(s.ctx, mysql.DB)
	p, err := productQuery.GetAllId()
	if err != nil {
		return nil, err
	}
	resp = &product.ListProductIdsResp{
		Id: make([]uint32, 0, len(p)),
	}
	return
}

package service

import (
	"context"
	"github.com/cloudwego/hertz/pkg/common/utils"
	"gomall/app/frontend/infra/rpc"
	product "gomall/rpc_gen/kitex_gen/product"
	rpcproduct "gomall/rpc_gen/kitex_gen/product"
)

type SearchProductsService struct {
	Context context.Context
}

func NewSearchProductsService(Context context.Context) *SearchProductsService {
	return &SearchProductsService{Context: Context}
}

// Run create note info
func (h *SearchProductsService) Run(req *product.SearchProductsReq) (resp map[string]any, err error) {
	products, err := rpc.ProductClient.SearchProducts(h.Context, &rpcproduct.SearchProductsReq{
		Query: req.Query,
	})
	if err != nil {
		return nil, err
	}

	resp = utils.H{
		"items": products.Results,
		"q":     req.Query,
	}
	return
}

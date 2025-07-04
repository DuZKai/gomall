package main

import (
	"context"
	"errors"
	"gomall/app/product/biz/service"
	product "gomall/rpc_gen/kitex_gen/product"
)

// ProductCatalogServiceImpl implements the last service interface defined in the IDL.
type ProductCatalogServiceImpl struct{}

// ListProducts implements the ProductCatalogServiceImpl interface.
func (s *ProductCatalogServiceImpl) ListProducts(ctx context.Context, req *product.ListProductsReq) (resp *product.ListProductsResp, err error) {
	resp, err = service.NewListProductsService(ctx).Run(req)

	return resp, err
}

// GetProduct implements the ProductCatalogServiceImpl interface.
func (s *ProductCatalogServiceImpl) GetProduct(ctx context.Context, req *product.GetProductReq) (resp *product.GetProductResp, err error) {
	resp, err = service.NewGetProductService(ctx).Run(req)

	return resp, err
}

// SearchProducts implements the ProductCatalogServiceImpl interface.
func (s *ProductCatalogServiceImpl) SearchProducts(ctx context.Context, req *product.SearchProductsReq) (resp *product.SearchProductsResp, err error) {
	rawResp, err := service.NewSearchProductsService(ctx).Run(req)
	if err != nil {
		return nil, err
	}

	products, ok := rawResp["items"].([]*product.Product)
	if !ok {
		return nil, errors.New("invalid type for items")
	}

	resp = &product.SearchProductsResp{
		Results: products,
	}
	return resp, nil
}

// UpdateProduct implements the ProductCatalogServiceImpl interface.
func (s *ProductCatalogServiceImpl) UpdateProduct(ctx context.Context, req *product.UpdateProductReq) (resp *product.UpdateProductResp, err error) {
	resp, err = service.NewUpdateProductService(ctx).Run(req)
	return resp, err
}

// ListProductIds implements the ProductCatalogServiceImpl interface.
func (s *ProductCatalogServiceImpl) ListProductIds(ctx context.Context, req *product.Empty) (resp *product.ListProductIdsResp, err error) {
	resp, err = service.NewListProductIdsService(ctx).Run()

	return resp, err
}

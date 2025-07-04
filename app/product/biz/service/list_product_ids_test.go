package service

import (
	"context"
	product "gomall/rpc_gen/kitex_gen/product"
	"testing"
)

func TestListProductIds_Run(t *testing.T) {
	ctx := context.Background()
	s := NewListProductIdsService(ctx)
	// init req and assert value

	req := &product.ListProductIdsReq{}
	resp, err := s.Run()
	t.Logf("err: %v", err)
	t.Logf("resp: %v", resp)

	// todo: edit your unit test

}

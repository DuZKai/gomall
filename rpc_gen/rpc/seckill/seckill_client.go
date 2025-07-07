package seckill

import (
	"gomall/rpc_gen/kitex_gen/seckill/seckillservice"
	"github.com/cloudwego/kitex/client"
	"github.com/cloudwego/kitex/client/callopt"
)

type RPCClient interface {
	KitexClient() seckillservice.Client
	Service() string
}

func NewRPCClient(dstService string, opts ...client.Option) (RPCClient, error) {
	kitexClient, err := seckillservice.NewClient(dstService, opts...)
	if err != nil {
		return nil, err
	}
	cli := &clientImpl{
		service:     dstService,
		kitexClient: kitexClient,
	}

	return cli, nil
}

type clientImpl struct {
	service     string
	kitexClient seckillservice.Client
}

func (c *clientImpl) Service() string {
	return c.service
}

func (c *clientImpl) KitexClient() seckillservice.Client {
	return c.kitexClient
}

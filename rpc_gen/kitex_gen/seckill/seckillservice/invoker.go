// Code generated by Kitex v0.9.1. DO NOT EDIT.

package seckillservice

import (
	server "github.com/cloudwego/kitex/server"
	seckill "gomall/rpc_gen/kitex_gen/seckill"
)

// NewInvoker creates a server.Invoker with the given handler and options.
func NewInvoker(handler seckill.SeckillService, opts ...server.Option) server.Invoker {
	var options []server.Option

	options = append(options, opts...)

	s := server.NewInvoker(options...)
	if err := s.RegisterService(serviceInfo(), handler); err != nil {
		panic(err)
	}
	if err := s.Init(); err != nil {
		panic(err)
	}
	return s
}

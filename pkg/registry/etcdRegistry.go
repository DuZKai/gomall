package registry

import (
	"context"
	"fmt"
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"
)

type EtcdRegistry struct {
	cli        *clientv3.Client
	serviceKey string
	cancel     context.CancelFunc
}

// NewEtcdRegistry 创建一个服务注册器
func NewEtcdRegistry(endpoints []string) (*EtcdRegistry, error) {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   endpoints,
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		return nil, err
	}
	return &EtcdRegistry{cli: cli}, nil
}

// Register 注册服务到 registry
func (r *EtcdRegistry) Register(serviceName, addr string, ttl int64) error {
	r.serviceKey = fmt.Sprintf("/services/%s/%s", serviceName, addr)
	leaseResp, err := r.cli.Grant(context.Background(), ttl)
	if err != nil {
		return err
	}

	// 设置 key 并附加租约
	_, err = r.cli.Put(context.Background(), r.serviceKey, addr, clientv3.WithLease(leaseResp.ID))
	if err != nil {
		return err
	}

	// 启动 keepalive
	ctx, cancel := context.WithCancel(context.Background())
	r.cancel = cancel
	ch, err := r.cli.KeepAlive(ctx, leaseResp.ID)
	if err != nil {
		return err
	}

	// 异步监听 keepalive 响应
	go func() {
		for ka := range ch {
			if ka == nil {
				fmt.Println("KeepAlive channel closed")
				break
			}
		}
	}()
	return nil
}

// Unregister 注销服务
func (r *EtcdRegistry) Unregister() error {
	if r.cancel != nil {
		r.cancel()
	}
	_, err := r.cli.Delete(context.Background(), r.serviceKey)
	return err
}

// Discover 获取指定服务的所有实例
func (r *EtcdRegistry) Discover(serviceName string) ([]string, error) {
	prefix := fmt.Sprintf("/services/%s/", serviceName)
	resp, err := r.cli.Get(context.Background(), prefix, clientv3.WithPrefix())
	if err != nil {
		return nil, err
	}
	var addrs []string
	for _, kv := range resp.Kvs {
		addrs = append(addrs, string(kv.Value))
	}
	return addrs, nil
}

// Watch 监听服务变化
func (r *EtcdRegistry) Watch(serviceName string, onChange func([]string)) {
	prefix := fmt.Sprintf("/services/%s/", serviceName)
	go func() {
		watchChan := r.cli.Watch(context.Background(), prefix, clientv3.WithPrefix())
		for range watchChan {
			addrs, err := r.Discover(serviceName)
			if err == nil {
				onChange(addrs)
			}
		}
	}()
}

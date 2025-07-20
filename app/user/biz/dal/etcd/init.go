// registry/registry.go
package etcd

import (
	"context"
	"fmt"
	"gomall/app/user/conf"
	"sync"
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"
	"go.uber.org/zap"
)

var (
	cli    *clientv3.Client
	once   sync.Once
	logger *zap.Logger
)

// Init 初始化 registry 客户端，单例
func Init() {
	endpoints := conf.GetConf().Etcd.Address
	var err error
	once.Do(func() {
		cli, err = clientv3.New(clientv3.Config{
			Endpoints:   endpoints,
			DialTimeout: 5 * time.Second,
		})
		if err != nil {
			return
		}
		// 可以初始化 zap logger
		logger, _ = zap.NewProduction()
	})
	if err != nil {
		panic(err)
	}
}

// Put 写 key
func Put(key, val string) error {
	if cli == nil {
		return fmt.Errorf("registry client not initialized")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	_, err := cli.Put(ctx, key, val)
	return err
}

// Get 读 key
func Get(key string) (string, error) {
	if cli == nil {
		return "", fmt.Errorf("registry client not initialized")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	resp, err := cli.Get(ctx, key)
	if err != nil {
		return "", err
	}
	if len(resp.Kvs) == 0 {
		return "", fmt.Errorf("key not found")
	}
	return string(resp.Kvs[0].Value), nil
}

// Watch 监听 key 变化，回调触发热更新
func Watch(key string, onChange func(key, val string)) {
	if cli == nil {
		logger.Error("registry client not initialized")
		return
	}
	go func() {
		watchChan := cli.Watch(context.Background(), key)
		for wresp := range watchChan {
			for _, ev := range wresp.Events {
				logger.Info("registry watch event",
					zap.String("type", ev.Type.String()),
					zap.String("key", string(ev.Kv.Key)),
					zap.String("value", string(ev.Kv.Value)),
				)
				onChange(string(ev.Kv.Key), string(ev.Kv.Value))
			}
		}
	}()
}

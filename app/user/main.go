package main

import (
	"context"
	"fmt"
	"github.com/cloudwego/hertz/pkg/app"
	hserver "github.com/cloudwego/hertz/pkg/app/server"
	"github.com/joho/godotenv"
	"gomall/app/user/biz/dal"
	"gomall/app/user/biz/dal/etcd"
	etcdService "gomall/pkg/registry"
	"gomall/rpc_gen/kitex_gen/user"
	"net"
	"time"

	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/cloudwego/kitex/pkg/rpcinfo"
	kserver "github.com/cloudwego/kitex/server"
	kitexlogrus "github.com/kitex-contrib/obs-opentelemetry/logging/logrus"
	"go.uber.org/zap/zapcore"
	"gomall/app/user/conf"
	"gomall/rpc_gen/kitex_gen/user/userservice"
	"gopkg.in/natefinch/lumberjack.v2"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		klog.Error(err.Error())
	}

	dal.Init()

	opts := kitexInit()
	svr := userservice.NewServer(new(UserServiceImpl), opts...)

	// 注册 Watch 监听某配置
	etcd.Watch("/config/app/config1", func(key, val string) {
		fmt.Printf("配置变更: %s = %s\n", key, val)
		// TODO: 这里触发你的配置热更新逻辑
	})

	// 启动 Hertz HTTP 服务（goroutine）
	go func() {
		h := hserver.Default(
			hserver.WithHostPorts(":8080"),
		)

		h.POST("/user/register", func(ctx context.Context, c *app.RequestContext) {
			req := &user.RegisterReq{}
			if err := c.BindAndValidate(req); err != nil {
				c.JSON(400, map[string]string{"error": "Invalid request body"})
				return
			}
			resp, err := new(UserServiceImpl).Register(ctx, req)
			if err != nil {
				c.JSON(500, map[string]string{"error": err.Error()})
				return
			}
			c.JSON(200, resp)
		})

		// 测试是否热更新
		// docker exec -it etcd1 etcdctl --endpoints=http://etcd1:2379 put /config/app/config1 '{"log_level":"INFO","db_host":"127.0.0.1","db_port":3306,"max_connections":50}'
		h.GET("/config/update", func(ctx context.Context, c *app.RequestContext) {
			key := "/config/app/config1"
			val := fmt.Sprintf("val-%d", time.Now().Unix())
			if err := etcd.Put(key, val); err != nil {
				c.JSON(500, map[string]string{"error": err.Error()})
				return
			}
			c.JSON(200, map[string]interface{}{"message": "updated", "key": key, "val": val})
		})

		if err := h.Run(); err != nil {
			klog.Error("Failed to run Hertz server: ", err)
		}
	}()

	// 启动 Kitex RPC 服务（goroutine）
	go func() {
		if err := svr.Run(); err != nil {
			klog.Error(err.Error())
		}
	}()

	// 主协程阻塞，防止程序退出
	select {}
}

func kitexInit() (opts []kserver.Option) {
	// address
	addr, err := net.ResolveTCPAddr("tcp", conf.GetConf().Kitex.Address)
	if err != nil {
		panic(err)
	}
	opts = append(opts, kserver.WithServiceAddr(addr))

	// service info
	opts = append(opts, kserver.WithServerBasicInfo(&rpcinfo.EndpointBasicInfo{
		ServiceName: conf.GetConf().Kitex.Service,
	}))

	etcdReg, err := etcdService.NewEtcdRegistry(conf.GetConf().Etcd.Address)
	if err != nil {
		klog.Fatal("registry init error: ", err)
	}

	// 注册到etcd, TTL 10秒
	err = etcdReg.Register("user-service", "192.168.101.65:8888", 10)
	if err != nil {
		klog.Fatal("registry register error: ", err)
	} else {
		klog.Infof("[registry] 服务已注册: service=user-service addr=192.168.101.65:8888")
	}

	// klog
	logger := kitexlogrus.NewLogger()
	klog.SetLogger(logger)
	klog.SetLevel(conf.LogLevel())
	asyncWriter := &zapcore.BufferedWriteSyncer{
		WS: zapcore.AddSync(&lumberjack.Logger{
			Filename:   conf.GetConf().Kitex.LogFileName,
			MaxSize:    conf.GetConf().Kitex.LogMaxSize,
			MaxBackups: conf.GetConf().Kitex.LogMaxBackups,
			MaxAge:     conf.GetConf().Kitex.LogMaxAge,
		}),
		FlushInterval: time.Minute,
	}
	klog.SetOutput(asyncWriter)
	kserver.RegisterShutdownHook(func() {
		err := asyncWriter.Sync()
		if err != nil {
			return
		}
	})
	return
}

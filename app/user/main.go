package main

import (
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	consul "github.com/kitex-contrib/registry-consul"
	"gomall/app/user/biz/dal"
	"gomall/rpc_gen/kitex_gen/user"
	"net"
	"time"

	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/cloudwego/kitex/pkg/rpcinfo"
	"github.com/cloudwego/kitex/server"
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

	// 启动 Gin HTTP 服务（在 goroutine 中）
	go func() {
		r := gin.Default()
		r.POST("/user/register", func(c *gin.Context) {
			req := &user.RegisterReq{}
			if err := c.ShouldBindJSON(req); err != nil {
				c.JSON(400, gin.H{"error": "Invalid request body"})
				return
			}
			resp, err := new(UserServiceImpl).Register(c.Request.Context(), req)
			if err != nil {
				c.JSON(500, gin.H{"error": err.Error()})
				return
			}
			c.JSON(200, resp)
		})
		if err := r.Run(":8080"); err != nil {
			klog.Error("Failed to run Gin server: ", err)
		}
	}()

	// 启动 Kitex RPC 服务（阻塞）
	if err := svr.Run(); err != nil {
		klog.Error(err.Error())
	}
}

func kitexInit() (opts []server.Option) {
	// address
	addr, err := net.ResolveTCPAddr("tcp", conf.GetConf().Kitex.Address)
	if err != nil {
		panic(err)
	}
	opts = append(opts, server.WithServiceAddr(addr))

	// service info
	opts = append(opts, server.WithServerBasicInfo(&rpcinfo.EndpointBasicInfo{
		ServiceName: conf.GetConf().Kitex.Service,
	}))

	r, err := consul.NewConsulRegister(conf.GetConf().Registry.RegistryAddress[0])
	if err != nil {
		klog.Fatal(err)
	}
	opts = append(opts, server.WithRegistry(r))

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
	server.RegisterShutdownHook(func() {
		asyncWriter.Sync()
	})
	return
}

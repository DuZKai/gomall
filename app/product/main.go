package main

import (
	"context"
	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/cloudwego/kitex/pkg/rpcinfo"
	"github.com/cloudwego/kitex/server"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	kitexlogrus "github.com/kitex-contrib/obs-opentelemetry/logging/logrus"
	consul "github.com/kitex-contrib/registry-consul"
	"go.uber.org/zap/zapcore"
	"gomall/app/product/biz/dal"
	"gomall/app/product/biz/service"
	"gomall/app/product/conf"
	"gomall/rpc_gen/kitex_gen/product"
	"gomall/rpc_gen/kitex_gen/product/productcatalogservice"
	"gopkg.in/natefinch/lumberjack.v2"
	"net"
	"net/http"
	"strconv"
	"time"
)

func main() {
	_ = godotenv.Load()
	dal.Init()
	opts := kitexInit()

	// 启动 Gin 接口服务（单独协程）
	go func() {
		r := setupGinRouter()
		if err := r.Run(":8079"); err != nil {
			panic("failed to start gin: " + err.Error())
		}
	}()

	svr := productcatalogservice.NewServer(new(ProductCatalogServiceImpl), opts...)

	err := svr.Run()
	if err != nil {
		klog.Error(err.Error())
	}
}

func setupGinRouter() *gin.Engine {
	r := gin.Default()

	r.GET("/products", func(c *gin.Context) {
		categoryName := c.Query("category_name")
		if categoryName == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "category_name is required"})
			return
		}
		pageStr := c.DefaultQuery("page", "1")           // 默认第1页
		pageSizeStr := c.DefaultQuery("page_size", "10") // 默认每页10条

		page, err := strconv.Atoi(pageStr)
		if err != nil || page < 1 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid page"})
			return
		}

		pageSize, err := strconv.Atoi(pageSizeStr)
		if err != nil || pageSize < 1 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid page_size"})
			return
		}

		// 构造 proto 请求结构体
		req := &product.ListProductsReq{
			CategoryName: categoryName,
			Page:         int32(page),
			PageSize:     int32(pageSize),
		}
		svc := service.NewListProductsService(context.Background())

		resp, err := svc.Run(req)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, resp)
	})

	r.POST("/products", func(c *gin.Context) {
		var req product.UpdateProductReq
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// 创建服务并执行
		svc := service.NewUpdateProductService(c.Request.Context()) // 使用请求上下文
		resp, err := svc.Run(&req)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, resp)
	})

	return r
}

func kitexInit() (opts []server.Option) {
	// address
	addr, err := net.ResolveTCPAddr("tcp", conf.GetConf().Kitex.Address)
	if err != nil {
		panic(err)
	}
	opts = append(opts, server.WithServiceAddr(addr))

	r, err := consul.NewConsulRegister(conf.GetConf().Registry.RegistryAddress[0])

	// service info
	opts = append(opts, server.WithServerBasicInfo(&rpcinfo.EndpointBasicInfo{
		ServiceName: conf.GetConf().Kitex.Service,
	}), server.WithRegistry(r))

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

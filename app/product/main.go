package main

import (
	"context"
	"fmt"
	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/cloudwego/kitex/pkg/rpcinfo"
	"github.com/cloudwego/kitex/server"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	kitexlogrus "github.com/kitex-contrib/obs-opentelemetry/logging/logrus"
	consul "github.com/kitex-contrib/registry-consul"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap/zapcore"
	"gomall/app/product/biz/dal"
	redisInit "gomall/app/product/biz/dal/redis"
	"gomall/app/product/biz/service"
	"gomall/app/product/conf"
	"gomall/rpc_gen/kitex_gen/product"
	"gomall/rpc_gen/kitex_gen/product/productcatalogservice"
	"gopkg.in/natefinch/lumberjack.v2"
	"log"
	"net"
	"net/http"
	"strconv"
	"time"
)

func main() {
	_ = godotenv.Load()
	dal.Init()
	opts := kitexInit()
	// bloomExample()

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

	r.GET("/productsAll", func(c *gin.Context) {
		svc := service.NewListProductIdsService(context.Background())

		resp, err := svc.Run()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, resp)
	})

	r.GET("/product", func(c *gin.Context) {
		productId := c.Query("product_id")
		if productId == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "product_id is required"})
			return
		}

		pid, err := strconv.ParseUint(productId, 10, 32)
		if err != nil || pid < 1 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid page"})
			return
		}

		req := &product.GetProductReq{
			Id: uint32(pid),
		}
		svc := service.NewGetProductService(context.Background())
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

func bloomExample() {
	// 创建 Redis 客户端
	client := redisInit.RedisBloomClient
	ctx := context.Background()

	// 删除已有的布隆过滤器（如果存在）
	fmt.Println("Deleting existing Bloom Filter...")
	_, err := client.Del(ctx, "my_filter").Result()
	if err != nil {
		log.Fatalf("Error deleting filter: %v", err)
	}

	// 创建布隆过滤器
	fmt.Println("Creating Bloom Filter...")
	_, err = client.BFReserve(ctx, "my_filter", 0.01, 10000).Result()
	if err != nil {
		log.Fatalf("Error creating filter: %v", err)
	}

	// 添加元素
	fmt.Println("Adding elements to Bloom Filter...")
	elements := []string{"element1", "element2", "element3"}
	for _, el := range elements {
		_, err := client.BFAdd(ctx, "my_filter", el).Result()
		if err != nil {
			log.Fatalf("Error adding element %s: %v", el, err)
		}
	}

	// 检查元素是否存在
	fmt.Println("Checking elements existence...")
	for _, el := range elements {
		exists, err := client.BFExists(ctx, "my_filter", el).Result()
		if err != nil {
			log.Fatalf("Error checking element %s: %v", el, err)
		}
		fmt.Printf("Element %s exists: %v\n", el, exists)
	}

	// 获取布隆过滤器的信息
	fmt.Println("Getting Bloom Filter information...")
	info, err := client.BFInfo(ctx, "my_filter").Result()
	if err != nil {
		log.Fatalf("Error getting filter info: %v", err)
	}
	fmt.Printf("Filter Info: %v\n", info)

	// 获取布隆过滤器的特定信息
	fmt.Println("Getting Bloom Filter specific information...")
	capacity, err := client.BFInfoCapacity(ctx, "my_filter").Result()
	if err != nil {
		log.Fatalf("Error getting filter capacity: %v", err)
	}
	size, err := client.BFInfoSize(ctx, "my_filter").Result()
	if err != nil {
		log.Fatalf("Error getting filter size: %v", err)
	}
	filters, err := client.BFInfoFilters(ctx, "my_filter").Result()
	if err != nil {
		log.Fatalf("Error getting filter filters: %v", err)
	}
	items, err := client.BFInfoItems(ctx, "my_filter").Result()
	if err != nil {
		log.Fatalf("Error getting filter items: %v", err)
	}
	fmt.Printf("Filter Capacity: %d\n", capacity)
	fmt.Printf("Filter Size: %d\n", size)
	fmt.Printf("Filter Filters: %v\n", filters)
	fmt.Printf("Filter Items: %d\n", items)

	// 获取布隆过滤器的扩展信息
	fmt.Println("Getting Bloom Filter expansion information...")
	expansion, err := client.BFInfoExpansion(ctx, "my_filter").Result()
	if err != nil {
		log.Fatalf("Error getting filter expansion: %v", err)
	}
	fmt.Printf("Filter Expansion: %d\n", expansion)

	// 批量插入元素
	fmt.Println("Batch inserting elements...")
	insertOptions := &redis.BFInsertOptions{}
	insertResult, err := client.BFInsert(ctx, "my_filter", insertOptions, "element4", "element5").Result()
	if err != nil {
		log.Fatalf("Error batch inserting elements: %v", err)
	}
	fmt.Printf("Insert Results: %v\n", insertResult)

	// 批量检查元素
	fmt.Println("Batch checking elements...")
	existResults, err := client.BFMExists(ctx, "my_filter", "element1", "element4").Result()
	if err != nil {
		log.Fatalf("Error batch checking elements: %v", err)
	}
	fmt.Printf("Batch Exist Results: %v\n", existResults)

	// 返回布隆过滤器的基数，即：布隆过滤器中添加的项目数量
	fmt.Println("Getting Bloom Filter card...")
	card, err := client.BFCard(ctx, "my_filter").Result()
	if err != nil {
		log.Fatalf("Error getting filter card: %v", err)
	}
	fmt.Printf("Filter Card: %d\n", card)

	// 结束
	fmt.Println("All operations completed.")
}

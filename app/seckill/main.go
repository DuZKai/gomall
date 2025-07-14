package main

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/zsais/go-gin-prometheus"
	"gomall/app/seckill/biz/dal"
	"gomall/app/seckill/biz/dal/asynq"
	"gomall/app/seckill/biz/dal/kafka"
	"gomall/app/seckill/biz/util"
	"gomall/app/seckill/config"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {

	if err := godotenv.Load(); err != nil {
		log.Fatalf("failed to load .env: %v", err)
	}

	// 创建上下文用于退出
	_, cancel := context.WithCancel(context.Background())
	defer cancel()

	dal.Init()
	config.LoadConfigFromConsul()
	config.StartConfigWatcher()

	// 启动 Seckill HTTP Server（后台）
	go seckillInit()

	go kafka.InitKafkaConsumer() // 后台启动消费者

	// 等待退出信号
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	<-sig
	log.Println("Shutting down gracefully...")

	cancel() // 通知消费者退出
	if kafka.ConsumerGroup != nil {
		if err := kafka.ConsumerGroup.Close(); err != nil {
			log.Printf("Error closing Kafka ConsumerGroup: %v", err)
		}
	}
	asynq.ShutdownAll()
}

func seckillInit() {
	r := gin.Default()

	// 集成gin-prometheus自动采集请求指标
	p := ginprometheus.NewPrometheus("seckill")
	p.Use(r)

	// 配置
	r.GET("/config", func(c *gin.Context) {
		appConfig := config.AppConfig
		c.JSON(200, gin.H{
			"token_ttl":             appConfig.TokenTTL,
			"blacklist_ttl":         appConfig.BlacklistTTL,
			"freq_limit_expire":     appConfig.FreqLimitExpire,
			"idempotent_key_expire": appConfig.IdempotentKeyExpire,
			"bucket_expire_seconds": appConfig.BucketExpireSeconds,
			"capacity_factor":       appConfig.CapacityFactor,
			"rate_factor":           appConfig.RateFactor,
			"base_token_rate":       appConfig.BaseTokenRate,
			"token_bucket_factor":   appConfig.TokenBucketFactor,
		})
	})
	// 秒杀请求
	r.POST("/seckill/request", util.SeckillRequestHandler)
	// 短轮询状态
	r.GET("/seckill/status", util.SeckillStatusHandler)
	// 支付下单入库
	r.POST("/seckill/checkout", util.SeckillCheckoutHandler)
	// 缓存预热
	r.POST("/seckill/activity/create", util.CreateSeckillActivity)
	// 分布式锁测试
	r.GET("/seckill/redisLock", util.RedisLockHandler)
	r.GET("/seckill/test", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Seckill service is running",
		})
	})

	// 自定义 http.Server
	srv := &http.Server{
		Addr:           ":8080",
		Handler:        r,
		ReadTimeout:    5 * time.Second,   // 可视业务需求调大
		WriteTimeout:   10 * time.Second,  // 可视业务需求调大
		IdleTimeout:    120 * time.Second, // KeepAlive连接保持时间
		MaxHeaderBytes: 1 << 20,           // 1MB header 限制
	}

	// 启动服务
	log.Fatal(srv.ListenAndServe())
}

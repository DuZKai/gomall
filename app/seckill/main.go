package main

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gomall/app/seckill/biz/dal"
	"gomall/app/seckill/biz/dal/asynq"
	"gomall/app/seckill/biz/dal/kafka"
	"gomall/app/seckill/config"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/cloudwego/kitex/pkg/rpcinfo"
	"github.com/cloudwego/kitex/server"
	kitexlogrus "github.com/kitex-contrib/obs-opentelemetry/logging/logrus"
	"go.uber.org/zap/zapcore"
	"gomall/app/seckill/biz/util"
	"gomall/app/seckill/conf"
	"gomall/rpc_gen/kitex_gen/seckill/seckillservice"
	"gopkg.in/natefinch/lumberjack.v2"
)

func main() {

	if err := godotenv.Load(); err != nil {
		log.Fatalf("failed to load .env: %v", err)
	}

	// 创建上下文用于退出
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	dal.Init()
	kafka.InitKafkaConsumerGroup(ctx) // 后台启动消费者
	config.LoadConfigFromConsul()
	config.StartConfigWatcher()

	// 启动 Seckill HTTP Server（后台）
	go seckillInit()

	go kitexRun() // 启动Kitex服务

	// 等待退出信号
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	<-sig
	log.Println("Shutting down gracefully...")

	cancel() // 通知消费者退出
	if err := kafka.ConsumerGroup.Close(); err != nil {
		log.Printf("Error closing Kafka ConsumerGroup: %v", err)
	}
	asynq.ShutdownAll()
}

func seckillInit() {
	r := gin.Default()
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
	err := r.Run(":8080")
	if err != nil {
		return
	}
}

func kitexRun() {
	opts := kitexInit()
	svr := seckillservice.NewServer(new(SeckillServiceImpl), opts...)

	err := svr.Run()
	if err != nil {
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

package dal

import (
	"github.com/bwmarrin/snowflake"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"gomall/app/seckill/biz/dal/asynq"
	"gomall/app/seckill/biz/dal/consul"
	"gomall/app/seckill/biz/dal/kafka"
	"gomall/app/seckill/biz/dal/mysql"
	"gomall/app/seckill/biz/dal/redis"
	"gomall/app/seckill/biz/dal/sentinel"
)

var (
	// 雪花算法节点
	Node *snowflake.Node

	// HTTP 请求指标
	HttpRequestsTotal = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "http_requests_total",
		Help: "Total HTTP requests",
	}, []string{"method", "path", "status_code"})

	HttpRequestDuration = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "http_request_duration_seconds",
		Help:    "HTTP request duration in seconds",
		Buckets: []float64{0.01, 0.05, 0.1, 0.3, 0.5, 1, 2, 5}, // 自定义时间桶
	}, []string{"method", "path"})

	// 业务指标
	SeckillBussinessRequestsTotal = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "seckill_business_requests_total",
		Help: "Total seckill requests",
	}, []string{"activity_id", "status"}) // status: success, rate_limited, invalid, etc.

	KafkaMessagesSent = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "seckill_kafka_messages_sent_total",
		Help: "Total messages sent to Kafka",
	}, []string{"status"}) // status: success, failed

	RateLimitedRequests = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "seckill_rate_limited_requests_total",
		Help: "Total rate-limited requests",
	}, []string{"priority"}) // priority: vip, normal
)

func Init() {
	mysql.Init()
	redis.Init()
	consul.Init()
	sentinel.Init()
	kafka.Init()
	asynq.Init()

	// 雪花算法初始化
	// 设置节点ID（0~1023）
	var err error
	Node, err = snowflake.NewNode(1)
	if err != nil {
		panic(err)
	}
}

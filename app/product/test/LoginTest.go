package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"gomall/app/product/biz/dal"
	"gomall/app/product/rpc"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func main() {
	_ = godotenv.Load()
	dal.Init()
	// 假设 rpc 包的 init() 已经完成了 Sentinel + userClient 的初始化
	r := gin.Default()

	// 测试接口：多次调用 rpc.Login，触发限流/熔断
	// 使用示例：http://192.168.101.65:8080/test/sentinel?n=200
	r.GET("/test/sentinel", func(c *gin.Context) {
		// 从 query 读取调用次数 n，默认 100
		n := 100
		if s := c.Query("n"); s != "" {
			if v, err := strconv.Atoi(s); err == nil && v > 0 {
				n = v
			}
		}

		// 并发触发
		for i := 0; i < n; i++ {
			go func(i int) {
				// 这里每次都调用 rpc.Login，内部已经做了 Entry/Exit/TraceError
				rpc.Login("abc@example.com", "123456")
				fmt.Printf("trigger #%d done\n", i)
			}(i)
		}

		c.JSON(http.StatusOK, gin.H{
			"message":     "started",
			"calls":       n,
			"description": "并发调用 rpc.Login，查看控制台日志或 Sentinel Dashboard",
		})
	})

	// 启动在 8080
	err := r.Run(":8888")
	if err != nil {
		return
	}
}

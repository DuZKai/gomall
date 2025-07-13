package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"net"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"sync/atomic"
	"syscall"
	"time"
)

type SeckillRequest struct {
	UserID     string `json:"user_id"`
	ActivityID string `json:"activity_id"`
	Captcha    string `json:"captcha"`
	Priority   int    `json:"priority"`
}

const (
	url         = "http://192.168.101.65:8080/seckill/request"
	threadCount = 1000 // 并发线程数
	loopCount   = 3    // 每个线程请求次数
	timeout     = 3 * time.Second
)

var (
	successCount int64
	failCount    int64
	timeoutCount int64
)

func isTimeoutErr(err error) bool {
	if err == nil {
		return false
	}
	// 检查是否是网络超时错误
	netErr, ok := err.(net.Error)
	return ok && netErr.Timeout()
}

func main() {
	rand.Seed(time.Now().UnixNano())
	start := time.Now()

	ctx, cancel := context.WithCancel(context.Background())
	var wg sync.WaitGroup

	// 捕获 Ctrl+C 退出
	go func() {
		sigCh := make(chan os.Signal, 1)
		signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
		<-sigCh
		fmt.Println("\n🛑 收到退出信号，正在终止请求...")
		cancel()
	}()

	client := &http.Client{Timeout: timeout}

	for i := 0; i < threadCount; i++ {
		wg.Add(1)
		go func(tid int) {
			defer wg.Done()
			for j := 0; j < loopCount; j++ {
				select {
				case <-ctx.Done():
					return
				default:
					reqData := SeckillRequest{
						UserID:     fmt.Sprintf("%d", rand.Intn(100000)+1000),
						ActivityID: "5001",
						Captcha:    fmt.Sprintf("%03d", rand.Intn(1000)),
						Priority:   rand.Intn(2), // 随机优先级 0 或 1
					}

					body, _ := json.Marshal(reqData)
					req, _ := http.NewRequest("POST", url, bytes.NewBuffer(body))
					req.Header.Set("Content-Type", "application/json")

					startReq := time.Now()
					resp, err := client.Do(req)
					duration := time.Since(startReq)

					if err != nil {
						if isTimeoutErr(err) || duration >= timeout {
							atomic.AddInt64(&timeoutCount, 1)
							fmt.Printf("⏱️ 超时线程 %d 请求 %d: %v\n", tid, j, err)
						} else {
							atomic.AddInt64(&failCount, 1)
							fmt.Printf("❌ 请求失败线程 %d 请求 %d: %v\n", tid, j, err)
						}
						continue
					}

					if resp.StatusCode == http.StatusOK {
						atomic.AddInt64(&successCount, 1)
					} else {
						atomic.AddInt64(&failCount, 1)
						fmt.Printf("❌ 非200响应线程 %d 请求 %d: 状态码 %d\n", tid, j, resp.StatusCode)
					}
					err = resp.Body.Close()
					if err != nil {
						return
					}
				}
			}
		}(i)
	}

	wg.Wait()
	duration := time.Since(start).Seconds()
	totalReq := int64(threadCount * loopCount)
	qps := float64(successCount) / duration

	fmt.Printf("\n📊 压测结果统计：\n")
	fmt.Printf("总请求数: %d\n", totalReq)
	fmt.Printf("成功数 ✅: %d\n", successCount)
	fmt.Printf("失败数 ❌: %d\n", failCount)
	fmt.Printf("超时数 ⏱️ : %d\n", timeoutCount)
	fmt.Printf("总耗时: %.2fs\n", duration)
	fmt.Printf("最大 QPS（仅成功）: %.2f\n", qps)
}

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
	url          = "http://192.168.101.65:8080/seckill/request"
	startQPS     = 100              // 初始 QPS
	maxQPS       = 2000             // 最大 QPS
	increaseStep = 100              // 每轮递增 QPS
	stepDuration = 10 * time.Second // 每个阶段持续时间
	timeout      = 3 * time.Second  // 请求超时
)

var (
	successCount int64
	failCount    int64
	timeoutCount int64
	histogram    = make([]int64, 10)
	histMu       sync.Mutex
)

func isTimeoutErr(err error) bool {
	if err == nil {
		return false
	}
	netErr, ok := err.(net.Error)
	return ok && netErr.Timeout()
}

func recordLatency(duration time.Duration) {
	ms := duration.Milliseconds()
	idx := 0
	switch {
	case ms <= 5:
		idx = 0
	case ms <= 10:
		idx = 1
	case ms <= 20:
		idx = 2
	case ms <= 50:
		idx = 3
	case ms <= 100:
		idx = 4
	case ms <= 200:
		idx = 5
	case ms <= 500:
		idx = 6
	case ms <= 1000:
		idx = 7
	case ms <= 2000:
		idx = 8
	default:
		idx = 9
	}
	histMu.Lock()
	histogram[idx]++
	histMu.Unlock()
}

func sendRequest(client *http.Client) {
	reqData := SeckillRequest{
		UserID:     fmt.Sprintf("%d", rand.Intn(100000)+1000),
		ActivityID: "5001",
		Captcha:    fmt.Sprintf("%03d", rand.Intn(1000)),
		Priority:   rand.Intn(2),
	}
	body, _ := json.Marshal(reqData)
	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	start := time.Now()
	resp, err := client.Do(req)
	duration := time.Since(start)
	recordLatency(duration)

	if err != nil {
		if isTimeoutErr(err) || duration >= timeout {
			atomic.AddInt64(&timeoutCount, 1)
		} else {
			atomic.AddInt64(&failCount, 1)
		}
		return
	}

	if resp.StatusCode == http.StatusOK {
		atomic.AddInt64(&successCount, 1)
	} else {
		atomic.AddInt64(&failCount, 1)
	}
	_ = resp.Body.Close()
}

func runRampUp(ctx context.Context) {
	client := &http.Client{Timeout: timeout}

	for currentQPS := startQPS; currentQPS <= maxQPS; currentQPS += increaseStep {
		fmt.Printf("\n当前 QPS: %d，压测中（持续 %.0fs）...\n", currentQPS, stepDuration.Seconds())
		atomic.StoreInt64(&successCount, 0)
		atomic.StoreInt64(&failCount, 0)
		atomic.StoreInt64(&timeoutCount, 0)

		stepCtx, cancel := context.WithTimeout(ctx, stepDuration)
		ticker := time.NewTicker(time.Second / time.Duration(currentQPS))
		var wg sync.WaitGroup

	Loop:
		for {
			select {
			case <-stepCtx.Done():
				break Loop
			case <-ctx.Done():
				cancel()
				return
			case <-ticker.C:
				wg.Add(1)
				go func() {
					defer wg.Done()
					sendRequest(client)
				}()
			}
		}

		ticker.Stop()
		cancel()
		wg.Wait()

		total := atomic.LoadInt64(&successCount) + atomic.LoadInt64(&failCount) + atomic.LoadInt64(&timeoutCount)
		successRate := float64(successCount) / float64(total) * 100

		fmt.Printf("成功: %d | 失败: %d | 超时: %d | 成功率: %.2f%%\n",
			successCount, failCount, timeoutCount, successRate)

		if successRate < 90.0 {
			fmt.Println("成功率低于 90%，系统可能达到极限，停止压测。")
			break
		}
	}
}

func printHistogram() {
	fmt.Println("\n响应时间分布（单位：ms）:")
	buckets := []string{"<=5", "<=10", "<=20", "<=50", "<=100", "<=200", "<=500", "<=1000", "<=2000", ">2000"}
	for i, v := range histogram {
		fmt.Printf("%-8s : %d\n", buckets[i], v)
	}
}

func main() {
	rand.Seed(time.Now().UnixNano())

	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		sig := make(chan os.Signal, 1)
		signal.Notify(sig, os.Interrupt, syscall.SIGTERM)
		<-sig
		fmt.Println("\n收到退出信号，终止压测...")
		cancel()
	}()

	fmt.Printf("开始 Ramp-Up 压测：起始 %d QPS，每 %ds 增加 %d，最大 QPS %d\n",
		startQPS, int(stepDuration.Seconds()), increaseStep, maxQPS)

	runRampUp(ctx)
	printHistogram()
}

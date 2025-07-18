package sentinel

import (
	"fmt"
	"github.com/alibaba/sentinel-golang/api"
	"github.com/alibaba/sentinel-golang/core/circuitbreaker"
)

func Init() {
	// 初始化 Sentinel
	if err := api.InitDefault(); err != nil {
		panic(fmt.Errorf("sentinel init error: %w", err))
	}

	// 熔断规则：当 user.Login 接口 30 次请求中错误比例超过 50% 时熔断
	_, err := circuitbreaker.LoadRules([]*circuitbreaker.Rule{
		{
			Resource:         "user.Login", // Resource 名称需和 api.Entry 使用的一致
			Strategy:         circuitbreaker.ErrorRatio,
			Threshold:        0.5,   // 错误比例阈值
			MinRequestAmount: 30,    // 最小统计请求数
			StatIntervalMs:   10000, // 统计窗口 10 秒
			RetryTimeoutMs:   3000,  // 熔断后 3 秒尝试重试一次
		},
	})
	if err != nil {
		panic(fmt.Errorf("circuitbreaker.LoadRules error: %w", err))
	}

	fmt.Println("Sentinel client init OK")
}

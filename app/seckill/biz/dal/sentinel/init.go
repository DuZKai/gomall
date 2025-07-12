package sentinel

import (
	"fmt"
	"github.com/alibaba/sentinel-golang/api"
	"github.com/alibaba/sentinel-golang/core/flow"
)

func Init() {
	// 初始化 Sentinel
	err := api.InitDefault()
	if err != nil {
		panic(err)
	}
	// 配置VIP用户流控规则（匀速排队）
	_, err = flow.LoadRules([]*flow.Rule{
		{
			Resource:               "seckill_vip", // 专门为VIP用户设置的资源名
			Threshold:              5000,          // VIP用户高QPS阈值
			TokenCalculateStrategy: flow.Direct,
			ControlBehavior:        flow.Throttling, // 匀速排队模式
			MaxQueueingTimeMs:      2000,            // 最长排队2秒
			StatIntervalInMs:       1000,
		},
	})
	if err != nil {
		panic(err)
	}
	fmt.Println("Sentinel init ok")
}

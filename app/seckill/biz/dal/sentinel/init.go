package sentinel

import (
	"fmt"
	"github.com/alibaba/sentinel-golang/api"
	"github.com/alibaba/sentinel-golang/core/flow"
)

func Init(prodNum int) {
	// 初始化 Sentinel
	err := api.InitDefault()
	if err != nil {
		panic(err)
	}
	// 配置流控规则
	_, err = flow.LoadRules([]*flow.Rule{
		{
			Resource:               "seckill_request",
			Threshold:              1.2 * float64(prodNum), // 每秒最多允许 10 个请求
			TokenCalculateStrategy: flow.Direct,
			ControlBehavior:        flow.Reject,
			StatIntervalInMs:       1000,
		},
	})
	if err != nil {
		panic(err)
	}
	fmt.Println("Sentinel init ok")
}

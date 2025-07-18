package rpc

import (
	"context"
	"fmt"
	"github.com/alibaba/sentinel-golang/api"
	"github.com/alibaba/sentinel-golang/core/base"
	"gomall/rpc_gen/kitex_gen/user"
	"gomall/rpc_gen/kitex_gen/user/userservice"
	"log"
	"time"

	"github.com/cloudwego/kitex/client"
)

var userClient userservice.Client

func InitUserRPC() {
	var err error
	userClient, err = userservice.NewClient(
		"user",
		client.WithHostPorts("192.168.101.65:8881"), // 改为 user 服务真实地址
		client.WithRPCTimeout(2*time.Second),
	)
	if err != nil {
		log.Fatalf("init user client failed: %v", err)
	}
}

func Login(email, password string) {
	// 对同一个 Resource 做限流/熔断
	entry, blockErr := api.Entry(
		"user.Login",
		api.WithResourceType(base.ResTypeRPC), // RPC 调用
		api.WithTrafficType(base.Inbound),     // 入站
	)
	if blockErr != nil {
		// 降级逻辑：被限流或熔断时走这里
		fmt.Println("被限流或熔断，降级处理用户 Login:", blockErr)
		return
	}
	// 记得 Exit，统计成功或失败
	defer entry.Exit()

	// 真正的 RPC 调用
	req := &user.LoginReq{
		Email:    email,
		Password: password,
	}
	resp, err := userClient.Login(context.Background(), req)
	if err != nil {
		// 上报错误给 Sentinel，用于熔断统计
		api.TraceError(entry, err)
		// 这里也可以再做一次降级或 fallback
		fmt.Println("RPC 登录失败，降级返回默认值")
		return
	}
	fmt.Printf("Login success! User ID = %d\n", resp.UserId)
}

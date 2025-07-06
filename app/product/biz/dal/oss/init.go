package oss

import (
	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss"
	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss/credentials"
	"gomall/app/product/conf"
)

var (
	OSSClient  *oss.Client
	BucketName string
)

func Init() {
	cfg := conf.GetConf().OSS
	// 创建静态凭证提供者
	credProvider := credentials.NewStaticCredentialsProvider(
		cfg.AccessKeyId,
		cfg.AccessKeySecret,
		"", // security token, 可选
	)
	// 加载默认配置并指定凭证、endpoint
	cliCfg := oss.LoadDefaultConfig().
		WithCredentialsProvider(credProvider).
		WithEndpoint(cfg.Endpoint).
		WithRegion("oss-cn-shenzhen")

	// 创建客户端
	OSSClient = oss.NewClient(cliCfg)
	BucketName = cfg.BucketName
	if OSSClient == nil {
		panic("failed to create OSS client")
	}
}

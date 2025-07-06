package util

import (
	"fmt"
	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss"
	"github.com/go-co-op/gocron"
	"os/exec"
	"time"
)

// UploadFile 演示如何上传本地文件到 OSS。
func UploadFileToOss(client *oss.Client, bucketName, objectKey, localPath string) {
	// // 打开本地文件
	// f, err := os.Open(localPath)
	// log.Printf("uploading file %s to OSS bucket %s with object key %s\n", localPath, bucketName, objectKey)
	// if err != nil {
	// 	panic(err)
	// }
	// defer func(f *os.File) {
	// 	err := f.Close()
	// 	if err != nil {
	// 		panic(err)
	// 	}
	// }(f)
	//
	// // 发起上传请求
	// res, err := client.PutObject(context.Background(), &oss.PutObjectRequest{
	// 	Bucket: oss.Ptr(bucketName),
	// 	Key:    oss.Ptr(objectKey),
	// 	Body:   io.Reader(f),
	// })
	// if err != nil {
	// 	panic(err)
	// }
	//
	// log.Printf("upload succeeded, ETag: %s\n", oss.ToString(res.ETag))

	s := gocron.NewScheduler(time.Local)
	_, err := s.Cron("* * * * *").Do(func() {
		fmt.Println("每分钟执行：", time.Now())
		cmd := exec.Command("ossutil", "cp", localPath, "oss://dzk-xiaokai-tea-house/"+objectKey)
		output, err := cmd.CombinedOutput()
		if err != nil {
			fmt.Println("执行命令失败:", err)
		} else {
			fmt.Println("输出结果:", string(output))
		}
	})
	if err != nil {
		fmt.Println("添加任务失败:", err)
		return
	}
	s.StartBlocking()
}

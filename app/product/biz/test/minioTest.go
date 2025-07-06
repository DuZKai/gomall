package minioTest

import (
	"context"
	"gomall/app/product/biz/util"
	"log"
	"time"
)

const (
	filePath         = "/root/go/src/gomall/app/product/main.go"
	objectName       = "main.go"
	downloadFilePath = "/root/go/src/gomall/app/product/main.txt"
	bucket           = "testbucket"
)

// 测试上传文件
func TestMinioUpload() {
	ctx := context.Background()

	err := util.UploadFile(ctx, bucket, objectName, filePath)
	if err != nil {
		panic(err)
	}
}

// 测试下载文件
func TestMinioDownloadFile() {
	ctx := context.Background()
	err := util.DownloadFile(ctx, bucket, objectName, downloadFilePath)
	if err != nil {
		panic(err)
	}
}

// 测试列出bucket下所有的对象
func TestListObjects() {
	ctx := context.Background()
	objects, err := util.ListObjects(ctx, bucket, "")
	if err != nil {
		panic(err)
	}
	log.Println("objects: ", objects)
}

// 删除对象
func TestDeleteFile() {
	ctx := context.Background()
	ret, err := util.DeleteFile(ctx, bucket, objectName)
	if err != nil {
		panic(err)
	}
	log.Println("delete object: ", ret)
}

// 获取对象，返回url
func TestGetPresignedGetObject() {
	ctx := context.Background()
	object, err := util.GetPresignedURL(ctx, bucket, objectName, 24*time.Hour)
	if err != nil {
		panic(err)
	}
	log.Println("GetPresignedGetObject: ", object)
}

package minioOperaion

import (
	"context"
	"fmt"
	"github.com/minio/minio-go"
	minioInit "gomall/app/product/biz/dal/minio"
	"os"
	"time"
)

// UploadFile 上传文件到指定 bucket
func UploadFile(ctx context.Context, bucket, objectName, filePath string) error {
	client := minioInit.MinioClient
	f, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("open %s: %w", filePath, err)
	}
	defer f.Close()

	_, err = client.PutObjectWithContext(ctx, bucket, objectName, f, -1, minio.PutObjectOptions{})
	if err != nil {
		return fmt.Errorf("put object: %w", err)
	}
	return nil
}

// DownloadFile 下载对象到本地
func DownloadFile(ctx context.Context, bucket, objectName, filePath string) error {
	client := minioInit.MinioClient
	if err := client.FGetObject(bucket, objectName, filePath, minio.GetObjectOptions{}); err != nil {
		return fmt.Errorf("download object: %w", err)
	}
	return nil
}

// DeleteFile 删除对象；返回 true/false 表示是否真正删除
func DeleteFile(ctx context.Context, bucket, objectName string) (bool, error) {
	client := minioInit.MinioClient
	if err := client.RemoveObject(bucket, objectName); err != nil {
		return false, fmt.Errorf("remove object: %w", err)
	}
	return true, nil
}

// ListObjects 按 prefix 遍历对象名
func ListObjects(ctx context.Context, bucket, prefix string) ([]string, error) {
	var names []string
	client := minioInit.MinioClient
	for obj := range client.ListObjects(bucket, prefix, true, nil) {
		if obj.Err != nil {
			return nil, obj.Err
		}
		names = append(names, obj.Key)
	}
	return names, nil
}

// GetPresignedURL 生成临时访问链接
func GetPresignedURL(ctx context.Context, bucket, objectName string, expires time.Duration) (string, error) {
	client := minioInit.MinioClient
	url, err := client.PresignedGetObject(bucket, objectName, expires, nil)
	if err != nil {
		return "", err
	}
	return url.String(), nil
}

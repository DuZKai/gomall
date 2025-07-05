package minio

import (
	"log"

	"github.com/minio/minio-go"
	"gomall/app/product/conf"
)

var MinioClient *minio.Client

func Init() {
	c := conf.GetConf().Minio
	var err error

	MinioClient, err = minio.New(c.Endpoint, c.AccessKey, c.SecretKey, false)
	if err != nil {
		panic(err)
	}
	log.Printf("MinIO client initialized: %s", c.Endpoint)
}

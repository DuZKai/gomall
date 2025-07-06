package es

import (
	"github.com/elastic/go-elasticsearch/v7"
	"gomall/app/product/conf"
	"log"
)

var (
	ESClient     *elasticsearch.Client
	Index        string
	SourceFields string
)

func Init() {
	cfg := conf.GetConf().ES
	esCfg := elasticsearch.Config{
		Addresses: cfg.Hostlist,
	}
	Index = cfg.Index
	SourceFields = cfg.SourceFields

	var err error

	ESClient, err = elasticsearch.NewClient(esCfg)
	if err != nil {
		panic(err)
	}

	// 测试连接
	res, err := ESClient.Info()
	if err != nil {
		panic(err)
	}

	if res.IsError() {
		panic(res.String())
	}

	log.Printf("Elasticsearch client initialized: %s", cfg.Hostlist)
}

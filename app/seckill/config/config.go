package config

import (
	"encoding/json"
	"fmt"
	"github.com/hashicorp/consul/api"
	"gomall/app/seckill/biz/model"
	"gomall/app/seckill/conf"
	"log"
	"time"
)

var AppConfig *model.SeckillLimitConfig

func LoadConfigFromConsul() {
	client, err := api.NewClient(api.DefaultConfig())
	if err != nil {
		log.Fatalf("Failed to create consul client: %v", err)
	}

	kv := client.KV()
	key := fmt.Sprintf("config/%s/seckill_limits", conf.GetConf().Registry.Env)

	pair, _, err := kv.Get(key, nil)
	if err != nil {
		log.Fatalf("Failed to get config from consul: %v", err)
	}

	if pair == nil {
		log.Fatal("No config found in consul at key: config/seckill_limits")
	}

	var cfg model.SeckillLimitConfig
	if err := json.Unmarshal(pair.Value, &cfg); err != nil {
		log.Fatalf("Failed to parse config json: %v", err)
	}

	AppConfig = &cfg
	log.Printf("Loaded config from Consul: %+v", AppConfig)
}

func StartConfigWatcher() {
	ticker := time.NewTicker(30 * time.Second)
	go func() {
		for range ticker.C {
			LoadConfigFromConsul()
		}
	}()
}

func GetSeckillLimitConfig() *model.SeckillLimitConfig {
	if AppConfig == nil {
		log.Println("AppConfig is nil, loading from Consul")
		LoadConfigFromConsul()
	}
	return AppConfig
}

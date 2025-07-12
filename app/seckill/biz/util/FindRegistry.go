package util

import (
	"fmt"
	"github.com/hashicorp/consul/api"
	"gomall/app/seckill/conf"
	"log"
)

func DiscoverService(serviceName string) []string {
	config := api.DefaultConfig()
	config.Address = conf.GetConf().Registry.RegistryAddress[0]

	client, err := api.NewClient(config)
	if err != nil {
		log.Fatal(err)
	}

	services, _, err := client.Health().Service(serviceName, "", true, nil)
	if err != nil {
		log.Fatal(err)
	}

	var addresses []string
	for _, svc := range services {
		addr := fmt.Sprintf("%s:%d", svc.Service.Address, svc.Service.Port)
		addresses = append(addresses, addr)
	}
	return addresses
}

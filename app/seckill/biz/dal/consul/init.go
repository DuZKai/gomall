package consul

import (
	"fmt"
	"github.com/hashicorp/consul/api"
	"gomall/app/seckill/conf"
)

func Init() {
	config := api.DefaultConfig()
	config.Address = conf.GetConf().Registry.RegistryAddress[0]

	address := conf.GetConf().Consul.Address
	port := conf.GetConf().Consul.Port

	client, err := api.NewClient(config)
	if err != nil {
		panic(err)
	}

	registration := &api.AgentServiceRegistration{
		ID:      "seckill-service-1",
		Name:    "seckill-service",
		Address: address,
		Port:    port,
		Tags:    []string{"go"},
		Check: &api.AgentServiceCheck{
			HTTP:     fmt.Sprintf("http://%s:%d/health", address, port),
			Interval: "10s",
			Timeout:  "1s",
		},
	}

	err = client.Agent().ServiceRegister(registration)
	if err != nil {
		panic(err)
	}
}

package consul

import (
	"github.com/hashicorp/consul/api"
	"gomall/app/seckill/conf"
)

func Init() {
	config := api.DefaultConfig()
	config.Address = conf.GetConf().Registry.RegistryAddress[0]

	client, err := api.NewClient(config)
	if err != nil {
		panic(err)
	}

	registration := &api.AgentServiceRegistration{
		ID:      "seckill-service-1",
		Name:    "seckill-service",
		Address: config.Address,
		Port:    8080,
		Tags:    []string{"go"},
		Check: &api.AgentServiceCheck{
			HTTP:     "http://" + config.Address + ":8080/health",
			Interval: "10s",
			Timeout:  "1s",
		},
	}

	err = client.Agent().ServiceRegister(registration)
	if err != nil {
		panic(err)
	}
}

package main

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
)

const (
	adminAPI   = "http://127.0.0.1:9180/apisix/admin" // APISIX Admin API 地址
	adminToken = "edd1c9f034335f136f87ad84b625c8f1"   // 默认 X-API-KEY
)

func addProductRoute() error {
	url := fmt.Sprintf("%s/routes/1", adminAPI)

	// 定义路由，将 /product/* 转发到 192.168.101.65:8079
	body := []byte(`{
		"uri": "/product*",
		"upstream": {
			"type": "roundrobin",
			"nodes": {
				"192.168.101.65:8079": 1
			}
		},
		"plugins": {
			"ip-restriction": {
				"whitelist": ["127.0.0.1", "192.168.101.0/24"]
			}
		}
	}`)

	req, err := http.NewRequest(http.MethodPut, url, bytes.NewReader(body))
	if err != nil {
		return err
	}
	req.Header.Set("X-API-KEY", adminToken)
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			panic(err)
		}
	}(resp.Body)

	data, _ := io.ReadAll(resp.Body)
	fmt.Println("Add Route Response:", string(data))
	return nil
}

func addWhitelist() error {
	url := fmt.Sprintf("%s/routes/1", adminAPI)

	// 全量更新而不是仅 PATCH
	body := []byte(`{
        "uri": "/product*",
        "upstream": {
            "type": "roundrobin",
            "nodes": { "192.168.101.65:8079": 1 }
        },
        "plugins": {
            "ip-restriction": {
                "whitelist": ["127.0.0.1", "192.168.101.0/24"]
            }
        }
    }`)

	req, err := http.NewRequest(http.MethodPatch, url, bytes.NewReader(body))
	if err != nil {
		return err
	}
	req.Header.Set("X-API-KEY", adminToken)
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			panic(err)
		}
	}(resp.Body)

	data, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("Add Whitelist Response:", string(data))
	return nil
}

func main() {
	if err := addProductRoute(); err != nil {
		fmt.Println("Error:", err)
	}
	if err := addWhitelist(); err != nil {
		fmt.Println("Error adding whitelist:", err)
	}
}

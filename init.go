package es

import (
	"fmt"
	"github.com/elastic/go-elasticsearch/v7"
)

// Es 全局
var (
	Es *elasticsearch.Client //elasticsearch

)

// 初始化
func init() {
	EsInit()
}

// EsInit 初始化es连接
func EsInit() {
	//注册es连接
	cfg := elasticsearch.Config{
		Addresses: []string{"http://127.0.0.1:9200"},
	}
	var err error
	Es, _ = elasticsearch.NewClient(cfg)
	if err != nil {
		fmt.Println("ES连接失败" + err.Error())
	}
}

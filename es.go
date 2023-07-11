package es

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"

	"github.com/elastic/go-elasticsearch/v7/esapi"
	"github.com/elastic/go-elasticsearch/v7/esutil"
)

// Data 最外层
type Data struct {
	Hits Hits `json:"hits,omitempty"`
}

// Hits 第一层
type Hits struct {
	Total map[string]interface{} `json:"total,omitempty"`
	Hits  []HitsTwo              `json:"hits,omitempty"`
}

// HitsTwo 第二层
type HitsTwo struct {
	Source json.RawMessage `json:"_source,omitempty"`
}

// InterfaceEs 获取切片长度
func InterfaceEs(len int) []interface{} {
	return make([]interface{}, len)
}

// AddEs 循环同步添加到es
func AddEs(indexName string, data []interface{}) {
	for _, v := range data {
		//循环数据添加到es
		res, err := Es.Index(
			indexName, //文档索引
			esutil.NewJSONReader(&v),
		)
		fmt.Println(res)
		if err != nil || res.IsError() {
			fmt.Println("Es索引添加失败" + err.Error())
		}
	}
}

// AddEsBuck es批量添加，buck
func AddEsBuck(indexName string, data []interface{}) {
	// 创建批量请求的主体
	var body bytes.Buffer

	for _, doc := range data {
		// 构建索引操作的元数据
		metaData := map[string]interface{}{
			"index": map[string]interface{}{
				"_index": indexName,
			},
		}

		// 将元数据和文档内容拼接为一条操作记录
		metaDataBytes, err := json.Marshal(metaData)
		if err != nil {
			fmt.Println(err.Error())
		}
		body.Write(metaDataBytes)
		body.WriteByte('\n')

		docBytes, err := json.Marshal(doc)
		if err != nil {
			fmt.Println(err.Error())
		}
		body.Write(docBytes)
		body.WriteByte('\n')
	}

	// 创建批量请求
	req := esapi.BulkRequest{
		Body:    bytes.NewReader(body.Bytes()),
		Refresh: "true",
	}

	// 执行批量请求
	res, err := req.Do(context.Background(), Es)
	if err != nil {
		fmt.Println(err.Error())
	}
	defer res.Body.Close()
	// 检查响应中的错误
	if res.IsError() {
		fmt.Println("批量插入过程中发生错误:", res.String())
	}
}

// ListEs 展示出es数据
func ListEs(keyword string, queryIndex string, from int, size int) []HitsTwo {
	//es搜索字段
	query := map[string]interface{}{
		"query": map[string]interface{}{
			"match": map[string]interface{}{
				"username": keyword,
			},
		},
	}
	result := Data{}
	//搜索判断
	if len(keyword) == 0 {
		//发送索引请求到es之中
		res, err := esapi.SearchRequest{
			Index: []string{queryIndex}, //索引内容为
			From:  &from,                //起始页数
			Size:  &size,                //分页
			//Sort:  []string{"num:asc"},  //排序
		}.Do(context.Background(), Es)
		if err != nil {
			fmt.Println(err.Error())
		}
		json.NewDecoder(res.Body).Decode(&result)
	} else {
		res, err := esapi.SearchRequest{
			Index: []string{queryIndex},
			Body:  esutil.NewJSONReader(query),
		}.Do(context.Background(), Es)
		if err != nil {
			fmt.Println(err.Error())
		}
		json.NewDecoder(res.Body).Decode(&result)
	}
	return result.Hits.Hits
}

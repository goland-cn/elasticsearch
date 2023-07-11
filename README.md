# elasticsearch

##### 安装:

```go
go get github.com/goland-cn/elasticsearch
```

##### 初始化elasticsearch:

```go
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
```

##### 循环同步es:

其中需要把参数数据转换为切片泛型

```go
//定义切片长度
interfaceSlice := make([]interface{}, len(data))
for i, item := range data {
	interfaceSlice[i] = item
}
```

##### 示例:

```go
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
```

##### 展示es数据:

```go
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
	//非空判断
	if len(keyword) == 0 {
		//发送索引请求到es之中
		res, err := esapi.SearchRequest{
			Index: []string{queryIndex}, //索引内容为
			From:  &from,                //起始页数
			Size:  &size,                //分页
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
```

*附加*

##### 进一步数据处理:

```go
//示例
//重新获取数据
var datas []data
var data Data
for _, v := range data {
	_ = json.Unmarshal(v.Source, &data)
	datas = append(datas, shop)
}
return datas
```


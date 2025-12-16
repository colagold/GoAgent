package main

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types"
)

// 定义一个结构体，用于序列化文档
type Article struct {
	Title   string    `json:"title"`
	Content string    `json:"content"`
	Created time.Time `json:"created"`
}

func main() {
	// 1. 初始化客户端 (Typed Client)
	cfg := elasticsearch.Config{
		Addresses: []string{"https://localhost:9200"},
		// 如果有账号密码：
		Username: "elastic",
		Password: "Wanwu123456",
		// 核心修改在这里：配置 Transport 忽略证书错误
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true, // ⚠️ 跳过证书验证
			},
		},
	}

	client, err := elasticsearch.NewTypedClient(cfg)
	if err != nil {
		log.Fatalf("Error creating the client: %s", err)
	}

	ctx := context.Background()
	indexName := "my-go-index"

	// 2. 创建索引 (如果不存在)
	// 判断索引是否存在
	exists, err := client.Indices.Exists(indexName).Do(ctx)
	if err != nil {
		log.Fatal(err)
	}

	if !exists {
		_, err := client.Indices.Create(indexName).Do(ctx)
		if err != nil {
			log.Fatalf("无法创建索引: %s", err)
		}
		fmt.Println("索引创建成功！")
	}

	// 3. 写入数据 (Index Document)
	doc := Article{
		Title:   "Golang Elasticsearch 入门",
		Content: "使用官方 Typed API 非常简单",
		Created: time.Now(),
	}

	// 使用 Index API 写入，指定 ID 为 "1"
	_, err = client.Index(indexName).
		Id("1").
		Request(doc).
		Do(ctx)
	if err != nil {
		log.Fatalf("写入失败: %s", err)
	}
	fmt.Println("文档写入成功！")

	// 4. 获取单条文档 (Get)
	getRes, err := client.Get(indexName, "1").Do(ctx)
	if err != nil {
		log.Printf("获取文档失败: %s", err)
	} else {
		// 需要反序列化 Source
		var foundArticle Article
		if err := json.Unmarshal(getRes.Source_, &foundArticle); err == nil {
			fmt.Printf("查找到文档: %+v\n", foundArticle)
		}
	}

	// 5. 搜索数据 (Search)
	// 等待 1 秒让 ES 刷新索引，否则刚写的数据可能搜不到（生产环境不需要这步）
	time.Sleep(1 * time.Second)

	searchRes, err := client.Search().
		Index(indexName).
		Query(&types.Query{
			Match: map[string]types.MatchQuery{
				"content": {Query: "API"}, // 搜索 content 字段包含 "API" 的文档
			},
		}).
		Do(ctx)

	if err != nil {
		log.Fatalf("搜索失败: %s", err)
	}

	fmt.Printf("搜索到 %d 条结果:\n", searchRes.Hits.Total.Value)

	// 遍历搜索结果
	for _, hit := range searchRes.Hits.Hits {
		var a Article
		if err := json.Unmarshal(hit.Source_, &a); err == nil {
			fmt.Printf(" - 标题: %s, 内容: %s\n", a.Title, a.Content)
		}
	}
}

package es

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/http"
	"strings"
	"sync"

	"github.com/elastic/go-elasticsearch/v8"
)

type Config struct {
	Address  string `json:"address" mapstructure:"address"`
	Username string `json:"username" mapstructure:"username"`
	Password string `json:"password" mapstructure:"password"`
}

type client struct {
	ctx context.Context
	cli *elasticsearch.Client

	mutex   sync.Mutex
	stopped bool
	stop    chan struct{}
}

func newClient(ctx context.Context, c Config) (*client, error) {
	// 智能判断协议，如果地址没有协议前缀，则尝试HTTPS，失败后尝试HTTP
	addresses := []string{}

	// 如果地址已经包含协议，直接使用
	if strings.HasPrefix(c.Address, "http://") || strings.HasPrefix(c.Address, "https://") {
		addresses = append(addresses, c.Address)
	} else {
		// 优先尝试HTTPS，然后HTTP
		addresses = append(addresses, "https://"+c.Address, "http://"+c.Address)
	}

	var lastErr error

	// 尝试每个地址
	for _, addr := range addresses {
		cfg := elasticsearch.Config{
			Addresses: []string{addr},
			Username:  c.Username,
			Password:  c.Password,
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{
					InsecureSkipVerify: true,
				},
			},
		}

		esClient, err := elasticsearch.NewClient(cfg)
		if err != nil {
			lastErr = fmt.Errorf("创建ES客户端失败 [%s]: %v", addr, err)
			//log.Warnf("创建ES客户端失败，地址: %s, 错误: %v", addr, err)
			continue
		}

		// 测试连接
		res, err := esClient.Info()
		if err != nil {
			lastErr = fmt.Errorf("ES连接测试失败 [%s]: %v", addr, err)
			//log.Warnf("ES连接测试失败，地址: %s, 错误: %v", addr, err)
			continue
		}

		if res != nil {
			defer res.Body.Close()

			if res.IsError() {
				lastErr = fmt.Errorf("ES连接响应错误 [%s]: %s", addr, res.String())
				//log.Warnf("ES连接响应错误，地址: %s, 响应: %s", addr, res.String())
				continue
			}
		}

		//log.Infof("ES连接成功，地址: %s", addr)
		return &client{
			ctx:  ctx,
			cli:  esClient,
			stop: make(chan struct{}, 1),
		}, nil
	}
	// 所有地址都失败了
	if lastErr != nil {
		return nil, lastErr
	}

	return nil, fmt.Errorf("无法连接到ES，尝试的地址: %v", addresses)
}

package tools

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/components/tool/utils"
)

// TavilyRequest 定义工具的输入参数
type TavilyRequest struct {
	Query      string `json:"query" jsonschema:"description=The search query to verify facts or look up information"`
	MaxResults int    `json:"max_results,omitempty" jsonschema:"description=Maximum number of results to return (default 5)"`
}

// TavilyResponse 定义工具的输出结构
type TavilyResponse struct {
	Query   string         `json:"query"`
	Results []TavilyResult `json:"results"`
	Error   string         `json:"error,omitempty"`
}

// TavilyResult 单条搜索结果
type TavilyResult struct {
	Title   string  `json:"title"`
	URL     string  `json:"url"`
	Content string  `json:"content"`
	Score   float64 `json:"score"`
}

// internalTavilyAPIRequest 用于构建发送给 Tavily API 的 JSON Body
type internalTavilyAPIRequest struct {
	APIKey      string `json:"api_key"`
	Query       string `json:"query"`
	MaxResults  int    `json:"max_results,omitempty"`
	SearchDepth string `json:"search_depth,omitempty"` // basic or advanced
}

// internalTavilyAPIResponse 用于解析 Tavily API 的原始响应
type internalTavilyAPIResponse struct {
	Query   string         `json:"query"`
	Results []TavilyResult `json:"results"`
}

// NewTavilySearchTool 创建一个 Tavily 搜索工具实例
func NewTavilySearchTool(ctx context.Context) (tool.BaseTool, error) {
	// 定义工具的名称和描述
	return utils.InferTool("tavily_search", "Search the internet using Tavily API to get up-to-date information, news, and facts.max_results not exceed 5",
		func(ctx context.Context, req *TavilyRequest) (*TavilyResponse, error) {
			// 1. 参数校验
			if req.Query == "" {
				return &TavilyResponse{Error: "Search query is required"}, nil
			}
			if req.MaxResults <= 0 {
				req.MaxResults = 5 // 设置默认值
			}

			// 2. 构建 HTTP 请求体
			apiReq := internalTavilyAPIRequest{
				APIKey:      os.Getenv("TAVILY_API_KEY"),
				Query:       req.Query,
				MaxResults:  req.MaxResults,
				SearchDepth: "basic", // 可以根据需求改为 "advanced"
			}

			jsonData, err := json.Marshal(apiReq)
			if err != nil {
				return nil, fmt.Errorf("failed to marshal request: %w", err)
			}

			// 3. 发起 HTTP 请求
			// 建议使用带有超时的 http client
			client := &http.Client{Timeout: 10 * time.Second}
			httpReq, err := http.NewRequestWithContext(ctx, "POST", "https://api.tavily.com/search", bytes.NewBuffer(jsonData))
			if err != nil {
				return nil, fmt.Errorf("failed to create http request: %w", err)
			}
			httpReq.Header.Set("Content-Type", "application/json")

			resp, err := client.Do(httpReq)
			if err != nil {
				return nil, fmt.Errorf("tavily api request failed: %w", err)
			}
			defer resp.Body.Close()

			// 4. 处理 HTTP 响应状态
			if resp.StatusCode != http.StatusOK {
				return nil, fmt.Errorf("tavily api returned status: %s", resp.Status)
			}

			// 5. 解析响应数据
			var apiResp internalTavilyAPIResponse
			if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
				return nil, fmt.Errorf("failed to decode response: %w", err)
			}

			// 6. 返回符合 Tool 定义的结构
			return &TavilyResponse{
				Query:   apiResp.Query,
				Results: apiResp.Results,
			}, nil
		})
}

func GetAllSearchTools(ctx context.Context) ([]tool.BaseTool, error) {

	tavilyTool, err := NewTavilySearchTool(ctx)
	if err != nil {
		return nil, err
	}

	return []tool.BaseTool{tavilyTool}, nil
}

package model

import (
	"context"
	"log"
	"os"

	"github.com/cloudwego/eino-ext/components/model/ark"
	"github.com/cloudwego/eino/components/model"
)

func CreateArkToolChatModel(ctx context.Context) model.ToolCallingChatModel {
	chatModel, err := ark.NewChatModel(ctx, &ark.ChatModelConfig{
		BaseURL: os.Getenv("ARK_BASE_URL"), // Ollama 服务地址
		Model:   os.Getenv("ARK_MODEL"),    // 模型名称
		APIKey:  os.Getenv("ARK_API_KEY"),
	})
	if err != nil {
		log.Fatalf("create ollama chat model failed: %v", err)
	}
	return chatModel
}

func CreateArkChatModel(ctx context.Context) *ark.ChatModel {
	chatModel, err := ark.NewChatModel(ctx, &ark.ChatModelConfig{
		BaseURL: os.Getenv("ARK_BASE_URL"), // Ollama 服务地址
		Model:   os.Getenv("ARK_MODEL"),    // 模型名称
		APIKey:  os.Getenv("ARK_API_KEY"),
	})
	if err != nil {
		log.Fatalf("create ollama chat model failed: %v", err)
	}
	return chatModel
}

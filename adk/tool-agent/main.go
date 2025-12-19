package main

import (
	"context"
	"log"

	chat_model "GoAgent/chat-model"
	"GoAgent/pkg/prints"
	"GoAgent/tools"

	"github.com/cloudwego/eino/adk"
	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/compose"
	"github.com/joho/godotenv"
)

func main() {
	// Entry point for the tool-agent application.
	ctx := context.Background()
	godotenv.Load(".env")
	// 初始化模型
	model := chat_model.CreateArkChatModel(ctx)
	weatherTool := tools.CreateWeatherTool()
	// 创建 ChatModelAgent
	agent, err := adk.NewChatModelAgent(ctx, &adk.ChatModelAgentConfig{
		Name:        "tool_agent",
		Description: "This agent can get the current weather for a given city.",
		Instruction: "Your sole purpose is to get the current weather for a given city by using the 'get_weather' tool.After calling the tool, report the result directly to the user.",
		Model:       model,
		ToolsConfig: adk.ToolsConfig{
			ToolsNodeConfig: compose.ToolsNodeConfig{
				Tools: []tool.BaseTool{weatherTool},
			},
		},
	})
	if err != nil {
		log.Fatal(err)
	}
	runner := adk.NewRunner(ctx, adk.RunnerConfig{
		EnableStreaming: true, // you can disable streaming here
		Agent:           agent,
	})

	// query weather
	println("\n\n>>>>>>>>>query weather<<<<<<<<<")
	iter := runner.Query(ctx, "What's the weather in Beijing?")
	for {
		event, ok := iter.Next()
		if !ok {
			break
		}
		if event.Err != nil {
			log.Fatal(event.Err)
		}

		prints.Event(event)
	}

}

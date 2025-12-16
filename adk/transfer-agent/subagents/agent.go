package subagents

import (
	"context"
	"fmt"
	"log"

	"github.com/cloudwego/eino/adk"
	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/components/tool/utils"
	"github.com/cloudwego/eino/compose"

	chat_model "GoAgent/chat-model"
)

type GetWeatherInput struct {
	City string `json:"city"`
}

func NewWeatherAgent() adk.Agent {
	weatherTool, err := utils.InferTool(
		"get_weather",
		"Gets the current weather for a specific city.", // English description
		func(ctx context.Context, input *GetWeatherInput) (string, error) {
			return fmt.Sprintf(`the temperature in %s is 25°C`, input.City), nil
		},
	)
	if err != nil {
		log.Fatal(err)
	}

	a, err := adk.NewChatModelAgent(context.Background(), &adk.ChatModelAgentConfig{
		Name:        "WeatherAgent",
		Description: "This agent can get the current weather for a given city.",
		Instruction: `Your sole purpose is to get the current weather for a given city by using the 'get_weather' tool.
After calling the tool, report the result directly to the user.`,
		Model: chat_model.CreateArkChatModel(context.Background()),
		ToolsConfig: adk.ToolsConfig{
			ToolsNodeConfig: compose.ToolsNodeConfig{
				Tools: []tool.BaseTool{weatherTool},
			},
		},
	})
	if err != nil {
		log.Fatal(err)
	}
	return a
}

func NewChatAgent() adk.Agent {
	a, err := adk.NewChatModelAgent(context.Background(), &adk.ChatModelAgentConfig{
		Name:        "ChatAgent",
		Description: "A general-purpose agent for handling conversational chat.", // English description
		Instruction: `You are a friendly conversational assistant.
Your role is to handle general chit-chat and answer questions that are not related to any specific tool-based tasks.`,
		Model: chat_model.CreateArkChatModel(context.Background()),
	})
	if err != nil {
		log.Fatal(err)
	}
	return a
}

func NewRouterAgent() adk.Agent {
	a, err := adk.NewChatModelAgent(context.Background(), &adk.ChatModelAgentConfig{
		Name:        "RouterAgent",
		Description: "A manual router that transfers tasks to other expert agents.",
		Instruction: `You are an intelligent task router.
Your responsibility is to analyze the user's request and delegate it to the most appropriate expert agent.
If no Agent can handle the task, simply inform the user it cannot be processed.`,
		Model: chat_model.CreateArkChatModel(context.Background()),
	})
	if err != nil {
		log.Fatal(err)
	}
	return a
}

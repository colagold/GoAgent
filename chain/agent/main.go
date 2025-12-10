package main

import (
	"context"
	"fmt"
	"log"

	"github.com/cloudwego/eino/callbacks"
	"github.com/cloudwego/eino/components/model"
	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/compose"
	"github.com/cloudwego/eino/schema"
	callbackHelpers "github.com/cloudwego/eino/utils/callbacks"
	"github.com/joho/godotenv"

	chat_model "GoAgent/chat-model"
	"GoAgent/tools"
)

func main() {
	//加载环境变量
	godotenv.Load(".env")
	//创建模型
	ctx := context.Background()
	chatModel := chat_model.CreateArkChatModel(ctx)
	//大模型回调函数
	modelHandler := &callbackHelpers.ModelCallbackHandler{
		OnEnd: func(ctx context.Context, info *callbacks.RunInfo, output *model.CallbackOutput) context.Context {
			// 1. output.Result 类型是 string
			fmt.Println("模型思考过程为：")
			fmt.Println(output.Message.Content)
			return ctx
		},
	}
	//工具回调函数
	toolHandler := &callbackHelpers.ToolCallbackHandler{
		OnStart: func(ctx context.Context, info *callbacks.RunInfo, input *tool.CallbackInput) context.Context {
			fmt.Printf("开始执行工具，参数: %s\n", input.ArgumentsInJSON)
			return ctx
		},
		OnEnd: func(ctx context.Context, info *callbacks.RunInfo, output *tool.CallbackOutput) context.Context {
			fmt.Printf("工具执行完成，结果: %s\n", output.Response)
			return ctx
		},
	}
	//构建实际回调函数Handler
	handler := callbackHelpers.NewHandlerHelper().
		ChatModel(modelHandler).
		Tool(toolHandler).
		Handler()

	getGameTool := tools.CreateTool()
	//绑定工具
	info, err := getGameTool.Info(ctx)
	if err != nil {
		panic(err)
	}
	infos := []*schema.ToolInfo{
		info,
	}
	err = chatModel.BindTools(infos)
	if err != nil {
		panic(err)
	}
	//创建tools节点
	ToolsNode, err := compose.NewToolNode(context.Background(), &compose.ToolsNodeConfig{
		Tools: []tool.BaseTool{
			getGameTool,
		},
	})
	if err != nil {
		panic(err)
	}
	//创建完整的处理链
	chain := compose.NewChain[[]*schema.Message, []*schema.Message]()
	chain.
		AppendChatModel(chatModel, compose.WithNodeName("chat_model")).
		AppendToolsNode(ToolsNode, compose.WithNodeName("tools"))

	// 编译并运行 chain
	agent, err := chain.Compile(ctx)
	if err != nil {
		log.Fatal(err)
	}
	//运行Agent
	resp, err := agent.Invoke(ctx, []*schema.Message{
		{
			Role:    schema.User,
			Content: "请告诉我原神的URL是什么",
		},
	}, compose.WithCallbacks(handler))
	if err != nil {
		log.Fatal(err)
	}

	// 输出结果
	for _, msg := range resp {
		fmt.Println(msg.Content)
	}
}

package main

import (
	model "GoAgent/chat-model"
	"context"
	"fmt"

	"github.com/cloudwego/eino/compose"
	"github.com/cloudwego/eino/schema"
	"github.com/joho/godotenv"
)

func main() {
	//加载环境变量
	godotenv.Load(".env")
	//示例代码
	ctx := context.Background()
	model := model.CreateArkChatModel(ctx)
	//编写lambda节点,封装[]*schema.Message
	lambda := compose.InvokableLambda(func(ctx context.Context, input string) (output []*schema.Message, err error) {
		desuwa := input + "回答结尾加上：以上内容由AI生成"
		output = []*schema.Message{
			{
				Role:    schema.User,
				Content: desuwa,
			},
		}
		return output, nil
	})
	//注册链条，其中[string, *schema.Message]分别为input类型和output类型，因为这是给模型输入，所以输出固定。
	chain := compose.NewChain[string, *schema.Message]()
	//连接起各个节点
	chain.AppendLambda(lambda).AppendChatModel(model)
	r, err := chain.Compile(ctx)
	if err != nil {
		panic(err)
	}
	answer, err := r.Invoke(ctx, "你好，请告诉我你的名字")
	if err != nil {
		panic(err)
	}
	fmt.Println(answer.Content)
}

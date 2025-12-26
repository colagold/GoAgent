package agents

import (
	"GoAgent/adk/deep-research/model"
	"context"

	"github.com/cloudwego/eino-ext/components/tool/commandline"
	"github.com/cloudwego/eino/adk"
)

func newDefaultAgent(ctx context.Context, op commandline.Operator) (adk.Agent, error) {
	m := model.NewChatModel()
	a, err := adk.NewChatModelAgent(ctx, &adk.ChatModelAgentConfig{
		Name:        "assistant",
		Description: "A friendly greeting assistant",
		Instruction: `
		You are a friendly assistant. Please respond to the user in a warm tone.`,
		Model: m,
		Exit:  &adk.ExitTool{},
	})
	if err != nil {
		return nil, err
	}
	return NewWrite2PlanMDWrapper(a, op), nil
}

package pkg

import (
	"context"
	"fmt"
	"time"

	"github.com/cloudwego/eino/callbacks"
	"github.com/cloudwego/eino/components/model"
	cbutils "github.com/cloudwego/eino/utils/callbacks"
)

func GetInputLoggerCallback() callbacks.Handler {
	return cbutils.NewHandlerHelper().ChatModel(&cbutils.ModelCallbackHandler{
		OnStart: func(ctx context.Context, info *callbacks.RunInfo, input *model.CallbackInput) context.Context {
			time.Sleep(20 * time.Second)
			fmt.Printf("\n========================================\n")
			fmt.Printf("[ChatModel Input] Agent: %s\n", info.Name)
			fmt.Printf("========================================\n")
			for i, msg := range input.Messages {
				fmt.Printf("  Message %d [%s]: %s\n", i+1, msg.Role, msg.Content)
				if len(msg.ToolCalls) > 0 {
					fmt.Printf("    Tool Calls: %d\n", len(msg.ToolCalls))
					for j, tc := range msg.ToolCalls {
						fmt.Printf("      %d. %s: %s\n", j+1, tc.Function.Name, tc.Function.Arguments)
					}
				}
			}
			fmt.Printf("========================================\n\n")
			return ctx
		},
	}).Handler()
}

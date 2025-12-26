package agents

import (
	"context"
	"log"

	"github.com/cloudwego/eino-ext/components/tool/commandline"
	"github.com/cloudwego/eino/adk"
	"github.com/cloudwego/eino/adk/prebuilt/planexecute"
)

func newDeepSearchAgent(ctx context.Context, op commandline.Operator) (adk.Agent, error) {
	planAgent, err := NewPlanner(ctx, op)
	if err != nil {
		log.Fatalf("agent.NewPlanner failed, err: %v", err)
		return nil, err
	}

	executeAgent, err := NewExecutor(ctx)
	if err != nil {
		log.Fatalf("agent.NewExecutor failed, err: %v", err)
		return nil, err
	}

	replanAgent, err := NewPlanner(ctx, op)
	if err != nil {
		log.Fatalf("agent.NewReplanAgent failed, err: %v", err)
		return nil, err
	}

	entryAgent, err := planexecute.New(ctx, &planexecute.Config{
		Planner:       planAgent,
		Executor:      executeAgent,
		Replanner:     replanAgent,
		MaxIterations: 20,
	})
	if err != nil {
		log.Fatalf("NewPlanExecuteAgent failed, err: %v", err)
		return nil, err
	}
	return entryAgent, nil
}

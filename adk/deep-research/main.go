package main

import (
	"GoAgent/adk/deep-research/agents"
	"GoAgent/adk/deep-research/operator"
	"GoAgent/adk/deep-research/params"
	"GoAgent/pkg/prints"
	"context"
	"log"

	"github.com/cloudwego/eino/adk"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load(".env")
	ctx := context.Background()

	operator := &operator.LocalOperator{}

	a, err := agents.BuildSupervisor(ctx, operator)
	if err != nil {
		log.Fatal(err)
	}
	runner := adk.NewRunner(ctx, adk.RunnerConfig{
		EnableStreaming: true, // you can disable streaming here
		Agent:           a,
	})

	ctx = params.InitContextParams(ctx)
	params.AppendContextParams(ctx, map[string]interface{}{
		params.WorkDirSessionKey: "akd/deep-research/logs",
		params.TaskIDKey:         uuid.New().String(),
	})

	// query weather
	println("\n\n>>>>>>>>>query<<<<<<<<<")
	iter := runner.Query(ctx, "调研联通万悟")
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

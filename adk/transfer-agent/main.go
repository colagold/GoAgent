/*
 * Copyright 2025 CloudWeGo Authors
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package main

import (
	"GoAgent/adk/transfer-agent/subagents"
	"GoAgent/pkg/prints"
	"context"
	"log"

	"github.com/cloudwego/eino/adk"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load(".env")
	weatherAgent := subagents.NewWeatherAgent()
	chatAgent := subagents.NewChatAgent()
	routerAgent := subagents.NewRouterAgent()

	ctx := context.Background()
	//设置父子关系：routerAgent为
	a, err := adk.SetSubAgents(ctx, routerAgent, []adk.Agent{chatAgent, weatherAgent})
	if err != nil {
		log.Fatal(err)
	}

	runner := adk.NewRunner(ctx, adk.RunnerConfig{
		EnableStreaming: true, // you can disable streaming here
		Agent:           a,
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

	// failed to route
	println("\n\n>>>>>>>>>failed to route<<<<<<<<<")
	iter = runner.Query(ctx, "Book me a flight from New York to London tomorrow.")
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

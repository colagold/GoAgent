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

package agents

import (
	"GoAgent/adk/deep-research/model"
	"context"

	"github.com/cloudwego/eino-ext/components/tool/commandline"
	"github.com/cloudwego/eino/adk"
	"github.com/cloudwego/eino/adk/prebuilt/supervisor"
)

func BuildSupervisor(ctx context.Context, op commandline.Operator) (adk.Agent, error) {
	m := model.NewChatModel()

	sv, err := adk.NewChatModelAgent(ctx, &adk.ChatModelAgentConfig{
		Name:        "supervisor",
		Description: "the agent responsible to supervise tasks",
		Instruction: `
		You are a supervisor managing two agents:

        - a plan agent. Assign research-related tasks to this agent
        - a chat agent. Assign another tasks to this agent
        Assign work to one agent at a time, do not call agents in parallel.
        Do not do any work yourself.`,
		Model: m,
		Exit:  &adk.ExitTool{},
	})
	if err != nil {
		return nil, err
	}

	deepSearchAgent, err := newDeepSearchAgent(ctx, op)
	if err != nil {
		return nil, err
	}
	defaultAgent, err := newDefaultAgent(ctx, op)
	if err != nil {
		return nil, err
	}

	sv_agent, err := supervisor.New(ctx, &supervisor.Config{
		Supervisor: sv,
		SubAgents:  []adk.Agent{deepSearchAgent, defaultAgent},
	})
	if err != nil {
		return nil, err
	}
	return NewWrite2PlanMDWrapper(sv_agent, op), nil
}

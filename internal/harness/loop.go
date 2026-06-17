package harness

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/ashintv/Zeu/internal/ai"
	"github.com/ashintv/Zeu/internal/logger"
	"github.com/ashintv/Zeu/internal/tools"
	"github.com/ashintv/Zeu/internal/types"
)

type State struct {
	messages []types.Coversation
	currIter int
}

type Agent struct {
	state    State
	Ai       *ai.AI
	MaxIter  int
	System   string
	registry *tools.ToolRegistry
	excuter  tools.ToolExcuter
}

func CreateAgent(Ai *ai.AI, reg *tools.ToolRegistry) *Agent {
	return &Agent{
		state: State{
			messages: []types.Coversation{},
			currIter: 0,
		},

		Ai:       Ai,
		MaxIter:  20,
		System:   ai.DefaultSystemPrompt(),
		registry: reg,
		excuter: tools.ToolExcuter{
			Registry: reg,
		},
	}
}

func (a *Agent) Invoke(ctx context.Context, Prompt string) <-chan string {
	conv := types.Coversation{
		Role:    "user",
		Content: Prompt,
	}

	a.state.messages = append(a.state.messages, conv)
	a.state.currIter = 0

	logger.Info("starting agent loop")
	return a.agentLoop(ctx)
}

func (a *Agent) agentLoop(ctx context.Context) <-chan string {
	streamCh := make(chan string)
	go func() {
		defer close(streamCh)
		for {
			if err := ctx.Err(); err != nil {
				logger.Info("Agent loop cancelled: ", err)
				break
			}

			if a.state.currIter >= a.MaxIter {
				logger.Info("Max Iteration reached")
				break
			}

			a.state.currIter += 1
			resChan := make(chan types.AiResponse)
			a.Ai.Invoke(
				ctx,
				resChan,
				ai.WithSystem(a.System),
				ai.WithMessages(a.state.messages),
				ai.WithTools(a.registry.List()),
			)
			toolCalls := []types.ToolCall{}

			cancelled := false

			var assistantMsg string
			logger.Infof("starting iteration %d", a.state.currIter)
		resLoop:
			for {
				select {
				case <-ctx.Done():
					logger.Info("Agent loop cancelled during AI invocation: ", ctx.Err())
					cancelled = true
					break resLoop
				case res, ok := <-resChan:
					if !ok {
						break resLoop
					}
					if res.Err != nil {
						logger.Error("AI invocation error: ", res.Err)
						cancelled = true
						break resLoop
					}
					if len(res.ToolsCalls) > 0 {
						toolCalls = append(toolCalls, res.ToolsCalls...)
					}

					if res.Messages != "" {
						assistantMsg += res.Messages
						select {
						case streamCh <- res.Messages:
						case <-ctx.Done():
							logger.Info("Agent loop cancelled during stream: ", ctx.Err())
							cancelled = true
							break resLoop
						}
					}
				}
			}

			if cancelled {
				break
			}

			if assistantMsg != "" || len(toolCalls) > 0 {
				a.state.messages = append(a.state.messages, types.Coversation{
					Role:      "assistant",
					Content:   assistantMsg,
					ToolCalls: toolCalls,
				})
			}

			if len(toolCalls) > 0 {
				for _, tc := range toolCalls {
					select {
					case streamCh <- fmt.Sprintf("\n[Calling tool: %s with args: %s]\n", tc.Name, string(tc.Args)):
					case <-ctx.Done():
						cancelled = true
						break
					}
				}
				if cancelled {
					break
				}

				res := a.excuter.ExcuteConcurrent(ctx, toolCalls)
				if err := ctx.Err(); err != nil {
					logger.Info("Agent loop cancelled after tool execution: ", err)
					break
				}
				for i, r := range res {
					var content string
					if errVal, ok := r.(error); ok {
						content = fmt.Sprintf("Tool error: %v", errVal)
					} else {
						jsonData, err := json.Marshal(r)
						if err != nil {
							content = fmt.Sprintf("Tool results: %v", r)
						} else {
							content = string(jsonData)
						}
					}
					a.state.messages = append(a.state.messages, types.Coversation{
						Role:       "tool",
						Content:    content,
						ToolCallID: toolCalls[i].Id,
					})

					select {
					case streamCh <- fmt.Sprintf("\n[Tool result: %s -> %s]\n", toolCalls[i].Name, content):
					case <-ctx.Done():
						cancelled = true
						break
					}
				}
				if cancelled {
					break
				}

				continue
			} else {
				logger.Info("Execution ended")
				break
			}

		}
	}()

	return streamCh
}

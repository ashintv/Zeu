package tools

import (
	"context"
	"fmt"
	"sync"

	"github.com/ashintv/Zeu/internal/logger"
	"github.com/ashintv/Zeu/internal/types"
)

type ToolExcuter struct {
	Registry *ToolRegistry
}

func (ex *ToolExcuter) Excute(ctx context.Context, toolCall types.ToolCall) (error, interface{}) {
	if err := ctx.Err(); err != nil {
		logger.Infof("Tool call [ID: %s, Name: %s] cancelled before starting: %v", toolCall.Id, toolCall.Name, err)
		return err, nil
	}

	err, toolFn := ex.Registry.Get(toolCall.Name)

	if err != nil {
		errRef := fmt.Errorf("Error finding tool %s", err)
		logger.Errorf("Failed to find tool [Name: %s]: %v", toolCall.Name, errRef)
		return errRef, nil
	}

	logger.Infof("Executing tool call [ID: %s, Name: %s] with args: %s", toolCall.Id, toolCall.Name, string(toolCall.Args))
	res, err := toolFn(ctx, toolCall.Args)
	if ctx.Err() != nil {
		logger.Infof("Tool call [ID: %s, Name: %s] cancelled during execution: %v", toolCall.Id, toolCall.Name, ctx.Err())
		return ctx.Err(), nil
	}
	if err != nil {
		errRef := fmt.Errorf(
			"Error excuting toolcall id: %s  name: %s with args %s ",
			toolCall.Id, toolCall.Name, toolCall.Args)
		logger.Errorf("Failed to execute tool call [ID: %s, Name: %s]: %v", toolCall.Id, toolCall.Name, err)
		return errRef, nil
	}

	logger.Infof("Successfully completed tool call [ID: %s, Name: %s] with result: %v", toolCall.Id, toolCall.Name, res)
	return nil, res
}

func (ex *ToolExcuter) ExcuteConcurrent(ctx context.Context, toolCalls []types.ToolCall) []interface{} {
	result := make([]interface{}, len(toolCalls))
	if err := ctx.Err(); err != nil {
		for i := range toolCalls {
			result[i] = err
		}
		return result
	}

	logger.Infof("Starting execution of %d tool calls concurrently", len(toolCalls))

	var wg sync.WaitGroup
	var mu sync.Mutex

	for i, toolCall := range toolCalls {
		wg.Add(1)

		go func(idx int, tc types.ToolCall) {
			defer wg.Done()
			
			if err := ctx.Err(); err != nil {
				mu.Lock()
				result[idx] = err
				mu.Unlock()
				return
			}

			err, res := ex.Excute(ctx, tc)
			
			mu.Lock()
			if ctx.Err() != nil {
				result[idx] = ctx.Err()
			} else if err != nil {
				result[idx] = err
			} else {
				result[idx] = res
			}
			mu.Unlock()
		}(i, toolCall)
	}

	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()

	select {
	case <-ctx.Done():
		logger.Info("Concurrent tool execution cancelled")
		mu.Lock()
		for idx := range result {
			if result[idx] == nil {
				result[idx] = ctx.Err()
			}
		}
		mu.Unlock()
		return result
	case <-done:
		logger.Info("All concurrent tool calls completed")
		return result
	}
}

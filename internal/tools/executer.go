package tools

import (
	"context"
	"fmt"
	"sync"

	"github.com/ashintv/Zeu/internal/types"
)

type ToolExcuter struct {
	registry ToolRegistry
}

func (ex *ToolExcuter) Excute(ctx context.Context, toolCall types.ToolCall) (error, interface{}) {
	err, toolFn := ex.registry.Get(toolCall.Name)

	if err != nil {
		return fmt.Errorf("Error finding tool %s", err), nil

	}

	res, err := toolFn(ctx, toolCall.Args)
	if err != nil {
		err := fmt.Errorf(
			"Error excuting toolcall id: %s  name: %s with args %s ",
			toolCall.Id, toolCall.Name, toolCall.Args)
		return err, nil
	}

	return nil, res
}

func (ex *ToolExcuter) ExcuteConcurrent(ctx context.Context, toolCalls []types.ToolCall) []interface{} {
	result := []interface{}{}

	var wg sync.WaitGroup

	for _, toolCall := range toolCalls {
		wg.Add(1)

		go func() {

			defer wg.Done()
			err, res := ex.Excute(ctx, toolCall)
			if err != nil {
				result = append(result, err)
				return
			}
			result = append(result, res)

		}()
	}

	wg.Wait()

	return result
}

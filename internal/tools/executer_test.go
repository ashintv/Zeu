package tools

import (
	"context"
	"encoding/json"
	"errors"
	"testing"
	"time"

	"github.com/ashintv/Zeu/internal/types"
)

func TestExcuteConcurrent_Cancel(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())

	called := make(chan struct{})

	slowTool := Tool{
		Name:   "slow_tool",
		Source: "built_in",
		Schema: types.Tool{},
		Excutable: func(ctx context.Context, args json.RawMessage) (any, error) {
			close(called)
			select {
			case <-ctx.Done():
				return nil, ctx.Err()
			case <-time.After(2 * time.Second):
				return "done", nil
			}
		},
	}

	registry := NewToolRegistry(WithTools([]Tool{slowTool}))
	ex := &ToolExcuter{Registry: registry}

	toolCalls := []types.ToolCall{
		{Id: "1", Name: "slow_tool", Args: nil},
	}

	done := make(chan struct{})
	var result []interface{}
	go func() {
		result = ex.ExcuteConcurrent(ctx, toolCalls)
		close(done)
	}()

	// Wait for tool to start executing
	<-called

	// Cancel context immediately
	cancel()

	// Ensure ExcuteConcurrent returns quickly without waiting 2 seconds
	select {
	case <-done:
		if len(result) != 1 {
			t.Fatalf("expected 1 result, got %d", len(result))
		}
		if !errors.Is(result[0].(error), context.Canceled) {
			t.Errorf("expected context.Canceled error, got %v", result[0])
		}
	case <-time.After(500 * time.Millisecond):
		t.Fatal("ExcuteConcurrent did not terminate quickly after cancelation")
	}
}

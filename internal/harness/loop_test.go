package harness

import (
	"context"
	"testing"
	"time"

	"github.com/ashintv/Zeu/internal/ai"
	"github.com/ashintv/Zeu/internal/tools"
	"github.com/ashintv/Zeu/internal/types"
)

type mockProvider struct {
	processFn func(ctx context.Context, req *types.AiRequest, streamCh chan<- types.AiResponse)
}

func (m *mockProvider) Info() types.ProviderInfo {
	return types.ProviderInfo{Name: "mock", Model: "mock"}
}

func (m *mockProvider) Default() *types.DefaultOptions {
	return &types.DefaultOptions{}
}

func (m *mockProvider) Process(ctx context.Context, req *types.AiRequest, streamCh chan<- types.AiResponse) {
	defer close(streamCh)
	if m.processFn != nil {
		m.processFn(ctx, req, streamCh)
	}
}

func TestAgentInvoke_Cancel(t *testing.T) {
	// Create a context that will be cancelled
	ctx, cancel := context.WithCancel(context.Background())

	calledProcess := make(chan struct{})

	mp := &mockProvider{
		processFn: func(ctx context.Context, req *types.AiRequest, streamCh chan<- types.AiResponse) {
			close(calledProcess)
			
			// Simulate long-running generation
			select {
			case <-ctx.Done():
				return
			case <-time.After(1 * time.Second):
				select {
				case streamCh <- types.AiResponse{Messages: "hello"}:
				case <-ctx.Done():
				}
			}
		},
	}

	registry := tools.NewToolRegistry()
	Ai := ai.NewAI(ai.Withprovider(mp))
	a := CreateAgent(Ai, registry)
	a.MaxIter = 5

	// Invoke the agent in a goroutine
	done := make(chan struct{})
	go func() {
		resChan := a.Invoke(ctx, "hello")
		for range resChan {
		}
		close(done)
	}()

	// Wait for process function to be entered
	select {
	case <-calledProcess:
	case <-time.After(500 * time.Millisecond):
		t.Fatal("mock provider process was not called")
	}

	// Cancel the context
	cancel()

	// Ensure agent loop exits quickly
	select {
	case <-done:
		// Success!
	case <-time.After(500 * time.Millisecond):
		t.Fatal("Agent loop did not terminate quickly after cancelation")
	}
}

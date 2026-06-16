package ai

import (
	"context"

	"github.com/ashintv/Zeu/internal/ai/provider"
	"github.com/ashintv/Zeu/internal/logger"
	"github.com/ashintv/Zeu/internal/types"
)

type Provider interface {
	Info() types.ProviderInfo
	Process(ctx context.Context, req *types.AiRequest, streamCh chan<- types.AiResponse)
	Default() *types.DefaultOptions
}

type AI struct {
	provider Provider
}

type AiOpts func(*AI)

func Withprovider(prv Provider) AiOpts {
	return func(a *AI) {
		a.provider = prv
	}
}

func NewAI(opts ...AiOpts) *AI {
	a := AI{
		provider: provider.NewOllama(),
	}

	for _, opt := range opts {
		opt(&a)
	}
	return &a
}

type InvokeOpts func(*types.AiRequest)

func WithMessages(Messages []types.Coversation) InvokeOpts {
	return func(ar *types.AiRequest) {
		ar.Messages = Messages
	}
}

func WithSystem(system string) InvokeOpts {
	return func(ar *types.AiRequest) {
		ar.System = system
	}
}

func WithTools(tools []types.Tool) InvokeOpts {
	return func(ar *types.AiRequest) {
		ar.Tools = tools
	}
}

func DefaultSystemPrompt() string {
	return `You are AI Harness Core.
			Your purpose is to understand the user's intent and provide the most accurate, useful, and complete response possible.

			Guidelines:
			- Follow user instructions carefully.
			- Be clear, concise, and practical.
			- Adapt your response style to the task.
			- Ask for clarification only when required to proceed.
			- Make reasonable assumptions when information is missing and state them when relevant.
			- Do not fabricate facts, sources, tool results, or actions.
			- If tools are available, use them when they improve accuracy or efficiency.
			- If a tool is unavailable, continue with the information you have and explain any limitations.
			- Prioritize correctness over confidence.
			- For complex tasks, think through the problem before responding.
			- Return results in the format most useful for the task.

			Your goal is successful task completion while maintaining accuracy, reliability, and helpfulness.`
}

func (a *AI) Invoke(ctx context.Context, opts ...InvokeOpts) (<-chan types.AiResponse) {
	req := types.AiRequest{
		System:   DefaultSystemPrompt(),
		Tools:    []types.Tool{},
		Messages: []types.Coversation{},
	}

	for _, opt := range opts {
		opt(&req)
	}

	logger.Info(
		"Request created",
		req,
		"with provider",
		a.provider.Info(),
	)

	streamChan := make(chan types.AiResponse)
	go a.provider.Process(ctx, &req, streamChan)
	return  streamChan
}

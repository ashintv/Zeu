package provider

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	
	"net/http"

	"github.com/ashintv/Zeu/internal/logger"
	"github.com/ashintv/Zeu/internal/types"
)

type Ollama struct {
	Name     string
	Model    string
	ApiKey   string
	Url      string
	DataType string
	Temp     float32
}

type Option func(*Ollama)

func WithOllamnaModel(model string) Option {
	return func(o *Ollama) {
		o.Model = model
	}
}

func WithOllamUrl(url string) Option {
	return func(o *Ollama) {
		o.Url = url

	}
}

func WithOllamaApiKey(apiKey string) Option {
	return func(o *Ollama) {
		o.ApiKey = apiKey

	}
}

func (O *Ollama) Default() *types.DefaultOptions {
	return &types.DefaultOptions{
		Model:    O.Model,
		ApiKey:   O.ApiKey,
		Temp:     O.Temp,
		Url:      O.Url,
		DataType: O.DataType,
	}
}

func NewOllama(opts ...Option) *Ollama {
	o := Ollama{
		Model:    "qwen3",
		Temp:     0.7,
		Url:      "http://localhost:11434/api/chat",
		ApiKey:   types.DEFAULT_API_KEY,
		DataType: "application/json",
	}

	for _, opt := range opts {
		opt(&o)
	}
	return &o
}

func (O *Ollama) Process(ctx context.Context, req *types.AiRequest, streamCh chan<- types.AiResponse) {
	defer close(streamCh)

	if err := ctx.Err(); err != nil {
		select {
		case streamCh <- types.AiResponse{Err: err}:
		case <-ctx.Done():
		}
		return
	}

	reqBody := O.BuildRequest(req)
	parsed, err := json.Marshal(reqBody)

	if err != nil {
		logger.Error("Unable to parse", reqBody)
		select {
		case streamCh <- types.AiResponse{Err: err}:
		case <-ctx.Done():
		}
		return
	}

	httpReq, err := http.NewRequestWithContext(ctx, "POST", O.Url, bytes.NewBuffer(parsed))
	if err != nil {
		select {
		case streamCh <- types.AiResponse{Err: err}:
		case <-ctx.Done():
		}
		return
	}
	httpReq.Header.Set("Content-Type", O.DataType)

	resp, err := http.DefaultClient.Do(httpReq)
	if err != nil {
		select {
		case streamCh <- types.AiResponse{Err: err}:
		case <-ctx.Done():
		}
		return
	}
	defer resp.Body.Close()

	scanner := bufio.NewScanner(resp.Body)
	for scanner.Scan() {
		if err := ctx.Err(); err != nil {
			return
		}

		line := scanner.Text()

		var chunk types.OllamaChatResponse

		if err := json.Unmarshal(
			[]byte(line),
			&chunk,
		); err != nil {
			continue
		}

		var toolCalls []types.ToolCall
		if len(chunk.Message.ToolCalls) > 0 {
			toolCalls = make([]types.ToolCall, len(chunk.Message.ToolCalls))
			for i, tc := range chunk.Message.ToolCalls {
				toolCalls[i] = types.ToolCall{
					Id:   tc.Id,
					Name: tc.Function.Name,
					Args: tc.Function.Arguments,
				}
			}
		}

		aiResponse := types.AiResponse{
			Messages:   chunk.Message.Content,
			Err:        nil,
			TimeStamp:  chunk.CreatedAt,
			ToolsCalls: toolCalls,
		}
		
		select {
		case streamCh <- aiResponse:
		case <-ctx.Done():
			return
		}
	}

	err = scanner.Err()
	if err != nil && ctx.Err() == nil {
		logger.Error(err)
		select {
		case streamCh <- types.AiResponse{Err: err}:
		case <-ctx.Done():
		}
	}
}

func (O *Ollama) BuildRequest(req *types.AiRequest) *types.OllamaChatRequest {
	messages := make([]types.OllamaRequestMessage, 0, len(req.Messages)+1)
	
	messages = append(messages, types.OllamaRequestMessage{
		Role:    "system",
		Content: req.System,
	})

	for _, msg := range req.Messages {
		var toolCalls []types.OllamaToolCall
		if len(msg.ToolCalls) > 0 {
			toolCalls = make([]types.OllamaToolCall, len(msg.ToolCalls))
			for i, tc := range msg.ToolCalls {
				toolCalls[i] = types.OllamaToolCall{
					Id:   tc.Id,
					Type: "function",
					Function: types.OllamaToolCallFunction{
						Name:      tc.Name,
						Arguments: tc.Args,
					},
				}
			}
		}

		messages = append(messages, types.OllamaRequestMessage{
			Role:       msg.Role,
			Content:    msg.Content,
			ToolCalls:  toolCalls,
			ToolCallID: msg.ToolCallID,
		})
	}

	return &types.OllamaChatRequest{
		Model:    O.Model,
		Messages: messages,
		Stream:   true,
		Tools:    req.Tools,
	}
}

func (O *Ollama) Info() types.ProviderInfo {
	return types.ProviderInfo{
		Model: O.Model,
		Name: O.Name,
	}
}

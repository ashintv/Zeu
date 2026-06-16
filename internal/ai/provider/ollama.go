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

	_ = ctx
	_ = streamCh
	reqBody := O.BuildRequest(req)
	parsed, err := json.Marshal(reqBody)

	if err != nil {
		logger.Error("Unable to parse", reqBody)
		streamCh <- types.AiResponse{
			Err: err,
		}
		return
	}

	resp, err := http.Post(O.Url, O.DataType, bytes.NewBuffer(parsed))
	if err != nil {
		streamCh <- types.AiResponse{
			Err: err,
		}
		return
	}
	defer resp.Body.Close()
	scanner := bufio.NewScanner(resp.Body)
	for scanner.Scan() {

		line := scanner.Text()

		var chunk types.OllamaChatResponse

		if err := json.Unmarshal(
			[]byte(line),
			&chunk,
		); err != nil {
			continue
		}

		aiResponse := types.AiResponse{
			Messages:   chunk.Message.Content,
			Err:        nil,
			TimeStamp:  chunk.CreatedAt,
			ToolsCalls: chunk.Message.ToolCalls,
		}
		
		streamCh <- aiResponse
	}

	err = scanner.Err()
	if err!= nil{
		logger.Error(err)
	}
	
}

func (O *Ollama) BuildRequest(req *types.AiRequest) *types.OllamaChatRequest {

	message := types.Coversation{
		Role: "system",
		Content: req.System,
	}

	messages := append([]types.Coversation{message},req.Messages...)

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

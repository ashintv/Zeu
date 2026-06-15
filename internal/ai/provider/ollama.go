package provider

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
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

func WithModel(model string) Option {
	return func(o *Ollama) {
		o.Model = model
	}
}

func WithUrl(url string) Option {
	return func(o *Ollama) {
		o.Url = url

	}
}

func WithApiKey(apiKey string) Option {
	return func(o *Ollama) {
		o.ApiKey = apiKey

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

func (O *Ollama) Process(ctx context.Context, req *types.AiRequest, streamCh chan<- types.AiResponse) (err error) {
	_ = ctx
	_ = streamCh

	reqBody := O.BuildRequest(req)
	parsed, err := json.Marshal(reqBody)

	if err != nil {
		logger.Error("Unable to parse", reqBody)
		return err
	}

	resp, err := http.Post(O.Url, O.DataType, bytes.NewBuffer(parsed))
	if err != nil {
		logger.Error("Error sending Request", reqBody)
		return err
	}
	
	defer resp.Body.Close()
	defer close(streamCh)
	
	scanner := bufio.NewScanner(resp.Body)

	for scanner.Scan() {
		
		line := scanner.Text()

		var chunk map[string]any

		if err := json.Unmarshal(
			[]byte(line),
			&chunk,
		); err != nil {
			continue
		}

		fmt.Println(chunk)

	}
	return scanner.Err()
}

func (O *Ollama) BuildRequest(req *types.AiRequest) *types.OllamaChatRequest {
	return &types.OllamaChatRequest{
		Model:    O.Model,
		Messages: req.Messages,
		Stream:   true,
		Tools:    req.Tools,
	}
}

func StreamChat() error {

	reqBody := []byte(`{
		"model":"qwen3",
		"messages":[
			{
				"role":"user",
				"content":"Hello"
			}
		],
		"stream":true
	}`)

	resp, err := http.Post(
		"http://localhost:11434/api/chat",
		"application/json",
		bytes.NewBuffer(reqBody),
	)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	scanner := bufio.NewScanner(resp.Body)

	for scanner.Scan() {
		line := scanner.Text()

		log.Println(line)

		var chunk map[string]any

		if err := json.Unmarshal(
			[]byte(line),
			&chunk,
		); err != nil {
			continue
		}

		fmt.Println(chunk)
	}

	return scanner.Err()
}

package types

import "encoding/json"

// Ollama-specific types for Ollama server communication

type OllamaToolFunction struct {
	Name        string          `json:"name"`
	Description string          `json:"description"`
	Parameters  json.RawMessage `json:"parameters"`
}

type OllamaTool struct {
	Type     string             `json:"type"`
	Function OllamaToolFunction `json:"function"`
}

type OllamaChatRequest struct {
	Model    string            `json:"model"`
	Messages []Coversation     `json:"messages"`
	Stream   bool              `json:"stream"`
	Tools    []json.RawMessage `json:"tools,omitempty"`
}

type OllamaChatMessage struct {
	Role      string     `json:"role"`
	Content   string     `json:"content"`
	Thinking  string     `json:"thinking,omitempty"`
	ToolCalls []ToolCall `json:"tool_calls,omitempty"`
}

type OllamaChatResponse struct {
	Model              string            `json:"model"`
	CreatedAt          string            `json:"created_at"`
	Message            OllamaChatMessage `json:"message"`
	Done               bool              `json:"done"`
	DoneReason         string            `json:"done_reason,omitempty"`
	TotalDuration      int64             `json:"total_duration,omitempty"`
	LoadDuration       int64             `json:"load_duration,omitempty"`
	PromptEvalCount    int               `json:"prompt_eval_count,omitempty"`
	PromptEvalDuration int64             `json:"prompt_eval_duration,omitempty"`
	EvalCount          int               `json:"eval_count,omitempty"`
	EvalDuration       int64             `json:"eval_duration,omitempty"`
}



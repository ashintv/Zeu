package types

import "encoding/json"

const DEFAULT_API_KEY = "KEY_NOT_SET_DEFAULT_API_KEY"

// ToolParameterProperty defines a single property in the parameters object
type ToolParameterProperty struct {
	Type        string `json:"type"`
	Description string `json:"description,omitempty"`
}

// ToolParameters defines the parameters schema for a tool function
type ToolParameters struct {
	Type       string                           `json:"type"`
	Required   []string                         `json:"required,omitempty"`
	Properties map[string]ToolParameterProperty `json:"properties"`
}

// ToolFunction defines the function details for a tool
type ToolFunction struct {
	Name        string         `json:"name"`
	Description string         `json:"description"`
	Parameters  ToolParameters `json:"parameters"`
}

// Tool represents a generic tool that can be used across different providers
type Tool struct {
	Type     string       `json:"type"`
	Function ToolFunction `json:"function"`
}

type ToolCall struct {
	Id   string
	Name string
	Args json.RawMessage
}

type Coversation struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type AiResponse struct {
	Messages   string     `json:"messages"`
	Err        error      `json:"error,omitempty"`
	TimeStamp  string     `json:"timestamp"`
	ToolsCalls []ToolCall `json:"tool_calls,omitempty"`
}

type AiRequest struct {
	System   string
	Tools    []Tool
	Messages []Coversation
}



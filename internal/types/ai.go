package types

import "encoding/json"


const DEFAULT_API_KEY = "KEY_NOT_SET_DEFAULT_API_KEY"

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
	messages   string
	err        error
	timeStamp  string
	toolsCalls []ToolCall
}

type AiRequest struct {
	System   string
	ApiKey   string
	Tools    []json.RawMessage
	Messages []Coversation
}



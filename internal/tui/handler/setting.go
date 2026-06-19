package handler

type ModelConfig struct {
	Model       string
	Temperature float64
	MaxTokens   int
}

func DefaultConfig() ModelConfig {
	return ModelConfig{
		Model:       "ollama/llama3",
		Temperature: 0.7,
		MaxTokens:   2048,
	}
}

package main

import (
	"encoding/json"
	"fmt"

	"github.com/ashintv/Zeu/internal/types"
)

func main() {
	// Example: Creating tools for weather API
	tools := []types.Tool{
		{
			Type: "function",
			Function: types.ToolFunction{
				Name:        "get_temperature",
				Description: "Get the current temperature for a city",
				Parameters: types.ToolParameters{
					Type:     "object",
					Required: []string{"city"},
					Properties: map[string]types.ToolParameterProperty{
						"city": {
							Type:        "string",
							Description: "The name of the city",
						},
					},
				},
			},
		},
		{
			Type: "function",
			Function: types.ToolFunction{
				Name:        "get_conditions",
				Description: "Get the current weather conditions for a city",
				Parameters: types.ToolParameters{
					Type:     "object",
					Required: []string{"city"},
					Properties: map[string]types.ToolParameterProperty{
						"city": {
							Type:        "string",
							Description: "The name of the city",
						},
					},
				},
			},
		},
	}

	// Serialize to JSON to see the structure
	jsonData, err := json.MarshalIndent(tools, "", "  ")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Println("Tools JSON structure:")
	fmt.Println(string(jsonData))

	// Example: Creating an AI request with tools
	request := types.AiRequest{
		System: "You are a helpful weather assistant.",
		Tools:  tools,
		Messages: []types.Coversation{
			{
				Role:    "user",
				Content: "What are the current weather conditions and temperature in New York and London?",
			},
		},
	}

	fmt.Println("\nAI Request created with", len(request.Tools), "tools")
}

// Made with Bob

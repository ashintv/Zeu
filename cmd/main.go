package main

import (
	"context"
	"fmt"

	"github.com/ashintv/Zeu/internal/ai"
	"github.com/ashintv/Zeu/internal/logger"
	"github.com/ashintv/Zeu/internal/types"
)

func main() {
	Ai := ai.NewAI()
	logger.Info("CREATED AI", Ai)

	logger.Info("Invoking ... ")

	messages := []types.Coversation{{
		Role:    "user",
		Content: "what is the weather conditon look like in kochi india",
	}}

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


	resChan := Ai.Invoke(context.Background(), ai.WithMessages(messages) , ai.WithTools(tools))
	var messagges string
    var toolCalls []types.ToolCall


	fmt.Println()
	fmt.Println()



	for res := range resChan {
		messag := res.Messages
		if messag != "" {
			fmt.Print(messag)
			messagges += messag
		}

		if len(res.ToolsCalls) > 0 {
			toolCalls = append(toolCalls, res.ToolsCalls...)
		}

	}
	fmt.Println()
	fmt.Println()
	for _, tool := range tools{
		logger.Info("Callinng tools" , tool.Function.Name , tool.Function.Parameters)
	}
	logger.Info()
	logger.Info("Completated query with", messages, tools)

}

package main

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/ashintv/Zeu/internal/ai"
	"github.com/ashintv/Zeu/internal/harness"
	"github.com/ashintv/Zeu/internal/logger"
	"github.com/ashintv/Zeu/internal/tools"
	"github.com/ashintv/Zeu/internal/types"
)

func main() {

	logger.Info("Invoking ... ")

	GetTempSchema := types.Tool{
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
	}

	GetCondtionSchema := types.Tool{
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
	}

	getCondTool := tools.Tool{
		Name:   GetCondtionSchema.Function.Name,
		Source: "built_in",
		Schema: GetCondtionSchema,
		Excutable: func(ctx context.Context, args json.RawMessage) (any, error) {
			var params struct {
				City string `json:"city"`
			}
			_ = json.Unmarshal(args, &params)
			loc := "Bangalore"
			if params.City != "" {
				loc = params.City
			}

			result := map[string]any{
				"location":    loc,
				"condition":   "Partly Cloudy",
				"temperature": 28,
				"feels_like":  31,
				"humidity":    72,
				"wind_speed":  12,
				"wind_unit":   "km/h",
				"success":     true,
			}

			return result, nil
		},
	}

	getTempTool := tools.Tool{
		Name:   GetTempSchema.Function.Name,
		Source: "built_in",
		Schema: GetTempSchema,
		Excutable: func(ctx context.Context, args json.RawMessage) (any, error) {
			var params struct {
				City string `json:"city"`
			}
			_ = json.Unmarshal(args, &params)
			loc := "Bangalore"
			if params.City != "" {
				loc = params.City
			}

			result := map[string]any{
				"location":    loc,
				"temperature": 28,
				"success":     true,
			}

			return result, nil
		},
	}

	registry := tools.NewToolRegistry(
		tools.WithTools(
			[]tools.Tool{
				getCondTool,
				getTempTool,
			},
		),
	)
	Ai := ai.NewAI()

	agent := harness.CreateAgent(Ai, registry)

	resChan := agent.Invoke(context.Background(), "what is the weather conditon look like in kochi india")

	fmt.Println()
	fmt.Println()

	for res := range resChan {
		fmt.Print(res)
	}

}

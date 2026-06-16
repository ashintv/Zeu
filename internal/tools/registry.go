package tools

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/ashintv/Zeu/internal/types"
)

type ToolExcutable func(ctx context.Context, args json.RawMessage) (any, error)

type tool struct {
	name      string
	source    string //mcp // built in
	schema    types.Tool
	excutable ToolExcutable
}

type ToolRegistry struct {
	caltelog map[string]tool
}

func (r *ToolRegistry) Get(name string) (error, ToolExcutable) {
	tool, ok := r.caltelog[name]

	if !ok {
		return fmt.Errorf("Invalid key please check function name"), nil
	}

	return nil, tool.excutable
}

func (r *ToolRegistry) Register(tools []tool) (bool, error) {
	for _, t := range tools {
		r.caltelog[t.name] = t
	}

	return true, nil
}

func (r *ToolRegistry) List() []types.Tool {
	tools := []types.Tool{}

	for _, t := range r.caltelog {
		tools = append(tools, t.schema)
	}

	return tools
}

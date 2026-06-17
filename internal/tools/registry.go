package tools

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/ashintv/Zeu/internal/types"
)

type ToolExcutable func(ctx context.Context, args json.RawMessage) (any, error)

type Tool struct {
	Name      string
	Source    string //mcp // built in
	Schema    types.Tool
	Excutable ToolExcutable
}

type ToolRegistry struct {
	caltelog map[string]Tool
}

type toolRegOpts func(*ToolRegistry)

func WithTools(tools []Tool) toolRegOpts {
	return func(tr *ToolRegistry) {
		for _, tl := range tools {
			tr.caltelog[tl.Name] = tl
		}
	}
}

func NewToolRegistry(opts ...toolRegOpts) *ToolRegistry{
	reg := ToolRegistry{
		caltelog: make(map[string]Tool),
	}

	for _, fn := range opts {
		fn(&reg)
	}

	return &reg
}


func (r *ToolRegistry) Get(name string) (error, ToolExcutable) {
	tool, ok := r.caltelog[name]

	if !ok {
		return fmt.Errorf("Invalid key please check function name"), nil
	}

	return nil, tool.Excutable
}

func (r *ToolRegistry) Register(tools []Tool) (bool, error) {
	for _, t := range tools {
		r.caltelog[t.Name] = t
	}

	return true, nil
}

func (r *ToolRegistry) List() []types.Tool {
	tools := []types.Tool{}

	for _, t := range r.caltelog {
		tools = append(tools, t.Schema)
	}

	return tools
}

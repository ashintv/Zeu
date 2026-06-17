package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/ashintv/Zeu/internal/ai"
	"github.com/ashintv/Zeu/internal/cli"
	"github.com/ashintv/Zeu/internal/harness"
	"github.com/ashintv/Zeu/internal/tools"
	"github.com/ashintv/Zeu/internal/types"
)

const (
	workspaceRoot = "/home/runner/work/Zeu/Zeu"
	maxReadBytes  = 1024 * 1024
)

func main() {
	registry := tools.NewToolRegistry(tools.WithTools([]tools.Tool{
		newReadFileTool(),
		newWriteFileTool(),
		newReadMultipleFilesTool(),
		newApplyDiffTool(),
		newShellTool(),
	}))

	agent := harness.CreateAgent(ai.NewAI(), registry)
	agent.System = codingSystemPrompt()

	terminalCLI := cli.NewCLI(agent)
	terminalCLI.Run()
}

func codingSystemPrompt() string {
	return `You are a minimal coding agent.

You have 5 tools:
1) read_file(path): read one file.
2) write_file(path, content): overwrite or create one file.
3) read_multiple_files(paths): read many files from a comma/newline separated path list.
4) apply_diff(diff): apply a unified diff patch with git apply.
5) shell(command): run a shell command in the workspace.

Tool usage rules:
- Use tools only when needed.
- Prefer read_file/read_multiple_files before write_file/apply_diff.
- Keep edits minimal and targeted.
- Use shell mainly for validation commands and diagnostics.
- Paths must stay inside the workspace.
- If a tool fails, explain the exact failure and next step.
`
}

func newReadFileTool() tools.Tool {
	schema := types.Tool{
		Type: "function",
		Function: types.ToolFunction{
			Name:        "read_file",
			Description: "Read one file from workspace and return content",
			Parameters: types.ToolParameters{
				Type:     "object",
				Required: []string{"path"},
				Properties: map[string]types.ToolParameterProperty{
					"path": {
						Type:        "string",
						Description: "Absolute or workspace-relative file path",
					},
				},
			},
		},
	}

	return tools.Tool{
		Name:   schema.Function.Name,
		Source: "built_in",
		Schema: schema,
		Excutable: func(ctx context.Context, args json.RawMessage) (any, error) {
			var in struct {
				Path string `json:"path"`
			}
			if err := json.Unmarshal(args, &in); err != nil {
				return nil, err
			}
			absPath, err := resolveWorkspacePath(in.Path)
			if err != nil {
				return nil, err
			}
			b, err := os.ReadFile(absPath)
			if err != nil {
				return nil, err
			}
			if len(b) > maxReadBytes {
				b = b[:maxReadBytes]
			}
			select {
			case <-ctx.Done():
				return nil, ctx.Err()
			default:
			}
			return map[string]any{
				"path":    absPath,
				"content": string(b),
			}, nil
		},
	}
}

func newWriteFileTool() tools.Tool {
	schema := types.Tool{
		Type: "function",
		Function: types.ToolFunction{
			Name:        "write_file",
			Description: "Write content to one file in workspace",
			Parameters: types.ToolParameters{
				Type:     "object",
				Required: []string{"path", "content"},
				Properties: map[string]types.ToolParameterProperty{
					"path": {
						Type:        "string",
						Description: "Absolute or workspace-relative file path",
					},
					"content": {
						Type:        "string",
						Description: "File content to write",
					},
				},
			},
		},
	}

	return tools.Tool{
		Name:   schema.Function.Name,
		Source: "built_in",
		Schema: schema,
		Excutable: func(ctx context.Context, args json.RawMessage) (any, error) {
			var in struct {
				Path    string `json:"path"`
				Content string `json:"content"`
			}
			if err := json.Unmarshal(args, &in); err != nil {
				return nil, err
			}
			absPath, err := resolveWorkspacePath(in.Path)
			if err != nil {
				return nil, err
			}
			if err := os.MkdirAll(filepath.Dir(absPath), 0o755); err != nil {
				return nil, err
			}
			select {
			case <-ctx.Done():
				return nil, ctx.Err()
			default:
			}
			if err := os.WriteFile(absPath, []byte(in.Content), 0o644); err != nil {
				return nil, err
			}
			return map[string]any{
				"path":    absPath,
				"written": true,
				"bytes":   len(in.Content),
			}, nil
		},
	}
}

func newReadMultipleFilesTool() tools.Tool {
	schema := types.Tool{
		Type: "function",
		Function: types.ToolFunction{
			Name:        "read_multiple_files",
			Description: "Read multiple files using a comma or newline separated path list",
			Parameters: types.ToolParameters{
				Type:     "object",
				Required: []string{"paths"},
				Properties: map[string]types.ToolParameterProperty{
					"paths": {
						Type:        "string",
						Description: "Comma or newline separated list of file paths",
					},
				},
			},
		},
	}

	return tools.Tool{
		Name:   schema.Function.Name,
		Source: "built_in",
		Schema: schema,
		Excutable: func(ctx context.Context, args json.RawMessage) (any, error) {
			var in struct {
				Paths string `json:"paths"`
			}
			if err := json.Unmarshal(args, &in); err != nil {
				return nil, err
			}
			rawPaths := splitPathList(in.Paths)
			if len(rawPaths) == 0 {
				return nil, errors.New("paths must include at least one file")
			}

			files := make([]map[string]any, 0, len(rawPaths))
			for _, p := range rawPaths {
				select {
				case <-ctx.Done():
					return nil, ctx.Err()
				default:
				}
				absPath, err := resolveWorkspacePath(p)
				if err != nil {
					files = append(files, map[string]any{
						"path":  p,
						"error": err.Error(),
					})
					continue
				}
				b, err := os.ReadFile(absPath)
				if err != nil {
					files = append(files, map[string]any{
						"path":  absPath,
						"error": err.Error(),
					})
					continue
				}
				if len(b) > maxReadBytes {
					b = b[:maxReadBytes]
				}
				files = append(files, map[string]any{
					"path":    absPath,
					"content": string(b),
				})
			}
			return map[string]any{"files": files}, nil
		},
	}
}

func newApplyDiffTool() tools.Tool {
	schema := types.Tool{
		Type: "function",
		Function: types.ToolFunction{
			Name:        "apply_diff",
			Description: "Apply a unified diff patch using git apply in workspace",
			Parameters: types.ToolParameters{
				Type:     "object",
				Required: []string{"diff"},
				Properties: map[string]types.ToolParameterProperty{
					"diff": {
						Type:        "string",
						Description: "Unified diff content",
					},
				},
			},
		},
	}

	return tools.Tool{
		Name:   schema.Function.Name,
		Source: "built_in",
		Schema: schema,
		Excutable: func(ctx context.Context, args json.RawMessage) (any, error) {
			var in struct {
				Diff string `json:"diff"`
			}
			if err := json.Unmarshal(args, &in); err != nil {
				return nil, err
			}
			cmd := exec.CommandContext(ctx, "git", "apply", "--whitespace=nowarn", "-")
			cmd.Dir = workspaceRoot
			cmd.Stdin = strings.NewReader(in.Diff)
			out, err := cmd.CombinedOutput()
			if err != nil {
				return nil, fmt.Errorf("git apply failed: %w: %s", err, string(out))
			}
			return map[string]any{
				"applied": true,
				"output":  strings.TrimSpace(string(out)),
			}, nil
		},
	}
}

func newShellTool() tools.Tool {
	schema := types.Tool{
		Type: "function",
		Function: types.ToolFunction{
			Name:        "shell",
			Description: "Run a shell command in workspace and return output",
			Parameters: types.ToolParameters{
				Type:     "object",
				Required: []string{"command"},
				Properties: map[string]types.ToolParameterProperty{
					"command": {
						Type:        "string",
						Description: "Shell command to execute",
					},
				},
			},
		},
	}

	return tools.Tool{
		Name:   schema.Function.Name,
		Source: "built_in",
		Schema: schema,
		Excutable: func(ctx context.Context, args json.RawMessage) (any, error) {
			var in struct {
				Command string `json:"command"`
			}
			if err := json.Unmarshal(args, &in); err != nil {
				return nil, err
			}
			runCtx, cancel := context.WithTimeout(ctx, 60*time.Second)
			defer cancel()

			cmd := exec.CommandContext(runCtx, "bash", "-lc", in.Command)
			cmd.Dir = workspaceRoot
			out, err := cmd.CombinedOutput()
			if err != nil {
				return map[string]any{
					"ok":       false,
					"output":   string(out),
					"exit_err": err.Error(),
				}, nil
			}
			return map[string]any{
				"ok":     true,
				"output": string(out),
			}, nil
		},
	}
}

func resolveWorkspacePath(path string) (string, error) {
	if strings.TrimSpace(path) == "" {
		return "", errors.New("path is required")
	}
	cleanInput := filepath.Clean(path)
	var absPath string
	if filepath.IsAbs(cleanInput) {
		absPath = cleanInput
	} else {
		absPath = filepath.Join(workspaceRoot, cleanInput)
	}

	absPath = filepath.Clean(absPath)
	rel, err := filepath.Rel(workspaceRoot, absPath)
	if err != nil {
		return "", err
	}
	if rel == ".." || strings.HasPrefix(rel, ".."+string(filepath.Separator)) {
		return "", fmt.Errorf("path is outside workspace: %s", path)
	}
	return absPath, nil
}

func splitPathList(v string) []string {
	v = strings.ReplaceAll(v, "\n", ",")
	parts := strings.Split(v, ",")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p != "" {
			out = append(out, p)
		}
	}
	return out
}

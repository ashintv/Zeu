# Zeu: AI Agent Harness Core

Zeu is a lightweight, layered, and context-cancellable AI agent harness built in Go. It enables developer-friendly LLM integration, concurrent tool invocation, real-time logging, streaming output, and a terminal interface.

---

## 📖 Table of Contents
1. [Architecture & Layers](#1-architecture--layers)
2. [Project Directory Structure](#2-project-directory-structure)
3. [Key Features](#3-key-features)
4. [Package Reference](#4-package-reference)
5. [Getting Started](#5-getting-started)
6. [Writing Custom Tools](#6-writing-custom-tools)
7. [Running Tests](#7-running-tests)

---

## 1. Architecture & Layers

Zeu is designed as a layered pipeline separating presentation, orchestration, model connection, and tool execution.

```
       [User / Main Client]
                │
                ▼
       [CLI Presentation Layer]   (bufio, os/signal)
                │
                ▼
          [Harness Layer]         (agent control loop, history state)
          /             \
         ▼               ▼
     [AI Layer]     [Tools Layer] (ToolRegistry, ToolExcuter)
         │
         ▼
  [Provider Layer] (Ollama API endpoint interface)
```

For a detailed layer-to-layer sequence and context propagation overview, see the [Architecture Flow Specification](flow.md).

---

## 2. Project Directory Structure

```
Zeu/
├── cmd/
│   └── main.go                 # Entrypoint initializing registry, AI, and starting CLI
├── internal/
│   ├── ai/
│   │   ├── provider/
│   │   │   └── ollama.go       # Ollama chat provider client implementation
│   │   └── llm.go              # Abstract model front-face & invocation
│   ├── cli/
│   │   └── cli.go              # Interactive CLI wrapper (SIGINT cancel logic)
│   ├── harness/
│   │   ├── loop.go             # Main Agent loop, streams messages & coordinates toolcalls
│   │   └── state.go            # Simple State structure definitions
│   ├── logger/
│   │   └── logger.go           # Thread-safe tag-based CLI structural logger
│   ├── tools/
│   │   ├── executer.go         # Concurrent & context-aware tool executers
│   │   └── registry.go         # Registry catalog holding executable tools
│   └── types/
│       ├── ai.go               # Shared types for Tool, Coversation, and requests
│       ├── ollama.go           # Ollama API serialization payload structures
│       └── provider.go         # Generic Provider information models
├── examples/
│   └── tools_example.go        # Quick example demonstrating schema generation
├── flow.md                     # Deep-dive design specification document
└── go.mod                      # Go module configurations
```

---

## 3. Key Features

- **Context-Based Cancellation Flow**: Full `context.Context` integration at every level. You can abort incoming LLM streaming calls, concurrent tool executions, and HTTP requests instantly using context cancellation.
- **Ctrl+C Terminal Interrupt Hook**: The interactive CLI listens for interrupt signals (`os.Interrupt`). Pressing Ctrl+C during agent execution stops the current run gracefully without terminating the CLI loop.
- **Concurrent & Safe Tool Executer**: Tools run concurrently in individual goroutines, synchronized via a `sync.WaitGroup` and protected by a `sync.Mutex` to prevent data races during cancellation.
- **JSON-based Tool Schema & Conversational State**: Implements standard OpenAI/Ollama conversation structure (`system`, `user`, `assistant`, `tool`), tracking `tool_call_id` to prevent LLM execution loops.
- **Interactive Console Output Streaming**: Real-time streaming of model responses, active tool invocations (`[Calling tool: ...]`), and execution returns (`[Tool result: ...]`).

---

## 4. Package Reference

### `internal/cli`
Provides the CLI wrapper that reads prompts from standard input and launches agent executions under a cancellable context. Handles OS signals cleanly.

### `internal/harness`
Controls the central agent loop. Resolves iteration bounds (`MaxIter`), delegates prompting, manages state history, and formats tool execution payloads into JSON for the model.

### `internal/ai` & `internal/ai/provider`
Wraps connection logic to LLMs. The `Ollama` provider connects to `http://localhost:11434`, streaming tokens and normalizing Ollama's nested tool call structure into generic `types.ToolCall` formats.

### `internal/tools`
Manages tool storage and concurrent task execution. The registry handles catalog lookups while the executor manages Go routines and context safety.

### `internal/logger`
Thread-safe, colorized terminal log utility. Provides level filters (`DEBUG`, `INFO`, `ERROR`) and caller location traces.

---

## 5. Getting Started

### Prerequisites
- Go 1.20 or later installed.
- Local [Ollama](https://ollama.com) instance running.
- Local model `qwen3` pulled (or customize the model in `internal/ai/provider/ollama.go` under `NewOllama`).

### Running the CLI
Run the main program from the project root:
```bash
go run cmd/main.go
```

This starts the interactive Zeu CLI:
```
Zeu CLI - Interactive Agent Terminal
Type 'exit' or 'quit' to exit. Press Ctrl+C during execution to cancel.
--------------------------------------------------

Zeu > what is the weather like in Kochi, India?

[Calling tool: get_conditions with args: {"city":"Kochi"}]

[Tool result: get_conditions -> {"condition":"Partly Cloudy","feels_like":31,"humidity":72,"location":"Kochi","success":true,"temperature":28,"wind_speed":12,"wind_unit":"km/h"}]

The weather in Kochi, India is currently Partly Cloudy at 28°C (feels like 31°C). The humidity is 72% with wind speeds of 12 km/h.
```

---

## 6. Writing Custom Tools

To define and register custom tools in Zeu, add them to `cmd/main.go` inside the `tools.Tool` registration sequence:

```go
customTool := tools.Tool{
    Name:   "custom_action",
    Source: "built_in",
    Schema: types.Tool{
        Type: "function",
        Function: types.ToolFunction{
            Name:        "custom_action",
            Description: "Execute a custom action",
            Parameters: types.ToolParameters{
                Type:     "object",
                Required: []string{"param1"},
                Properties: map[string]types.ToolParameterProperty{
                    "param1": {
                        Type:        "string",
                        Description: "Input parameter description",
                    },
                },
            },
        },
    },
    Excutable: func(ctx context.Context, args json.RawMessage) (any, error) {
        // Parse arguments
        var params struct {
            Param1 string `json:"param1"`
        }
        if err := json.Unmarshal(args, &params); err != nil {
            return nil, err
        }
        
        // Execute tool logic respecting context
        if err := ctx.Err(); err != nil {
            return nil, err
        }
        
        return map[string]any{"result": "success", "input": params.Param1}, nil
    },
}
```

---

## 7. Running Tests

Run the test suite across the packages using:
```bash
go test ./... -v
```

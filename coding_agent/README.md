# Coding Agent (Minimal Tool Set)

This is a separate minimal coding agent built on top of the existing Zeu harness.

## Run

```bash
go run /home/runner/work/Zeu/Zeu/coding_agent/main.go
```

## Included Tools (5)

1. `read_file`
   - Input: `{"path":"<file-path>"}`
2. `write_file`
   - Input: `{"path":"<file-path>","content":"<text>"}`
3. `read_multiple_files`
   - Input: `{"paths":"file1.go,file2.go"}` (comma/newline separated)
4. `apply_diff`
   - Input: `{"diff":"<unified-diff-content>"}`
5. `shell`
   - Input: `{"command":"go test ./..."}`

## Prompt System

The agent uses a dedicated system prompt in `coding_agent/main.go` (`codingSystemPrompt()`), which defines:

- what tools exist
- when to use each tool
- preference for minimal changes
- workspace path safety expectations

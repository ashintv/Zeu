# Logger Package

A lightweight, thread-safe, colored, and tag-based console logging package for Go.

## Implementation Details

### 1. Structure Definitions

#### Logger
The core logger holds synchronization locks, output writer configurations, active severity level, tag slice context, and formatting flags.

```go
type Logger struct {
	mu           sync.Mutex
	writer       io.Writer
	level        Level
	tags         []string
	enableColors bool
	showCaller   bool
}
```

#### StreamWriter
Handles chunk-by-chunk console streaming. It maintains a state machine to write the prefix headers on the first chunk and terminal delimiter upon completion.

```go
type StreamWriter struct {
	mu           sync.Mutex
	writer       io.Writer
	tags         []string
	enableColors bool
	started      bool
	closed       bool
	lastChar     byte
}
```

### 2. Concurrency Model
All log operations acquire a package-level or instance-level `sync.Mutex` lock prior to formatting and writing outputs. This ensures that concurrent writes from multiple goroutines do not interleave output lines.

### 3. Dynamic Caller Resolution
When caller tracing is enabled, the call stack is resolved dynamically using the standard `runtime.Caller` API. It iterates through frames to bypass internal logger wrappers automatically:

```go
func (l *Logger) getCallerInfo() string {
	for i := 2; i < 15; i++ {
		_, file, line, ok := runtime.Caller(i)
		if !ok {
			break
		}
		isLoggerPkg := strings.Contains(file, "internal/logger/") && !strings.HasSuffix(file, "_test.go")
		isLoggerFile := strings.HasSuffix(file, "logger.go")
		if isLoggerPkg || isLoggerFile {
			continue
		}
		short := file
		if idx := strings.LastIndex(file, "/"); idx >= 0 {
			short = file[idx+1:]
		}
		return fmt.Sprintf("%s:%d", short, line)
	}
	return "???:0"
}
```

### 4. Text Formatting & Terminal Styling
ANSI escape sequences are used for terminal formatting:
- Grey (`\033[90m`): Timestamps and callers
- Cyan (`\033[36m`): DEBUG level and stream tags
- Green (`\033[32m`): INFO level
- Red (`\033[31m`): ERROR level
- Bold Red (`\033[1;31m`): ERROR message bodies
- Magenta (`\033[35m`): User log tags

---

## Usage Examples

### Standard Logging

```go
package main

import "github.com/ashintv/Zeu/internal/logger"

func main() {
	logger.ShowCaller(true)
	logger.SetDefaultLevel(logger.DebugLevel)

	logger.Debug("Fetching database credentials")
	logger.Info("Listening on port 8080")
	logger.Error("Database connection timed out")
}
```

### Tagged Loggers

```go
package main

import "github.com/ashintv/Zeu/internal/logger"

func main() {
	dbLog := logger.Tagged("DATABASE")
	dbLog.Info("Connected")

	queryLog := dbLog.Tagged("QUERY")
	queryLog.Debug("SELECT * FROM users")
}
```

### Streaming Logs

```go
package main

import (
	"time"
	"github.com/ashintv/Zeu/internal/logger"
)

func main() {
	stream := logger.Tagged("AI").Stream()
	defer stream.Close()

	chunks := []string{"Hello", " world", "!", " Ready"}
	for _, chunk := range chunks {
		stream.WriteString(chunk)
		time.Sleep(50 * time.Millisecond)
	}
}
```

package logger

import (
	"fmt"
	"io"
	"os"
	"runtime"
	"strings"
	"sync"
	"time"
)

// Level represents the severity of a log message.
type Level int

const (
	DebugLevel Level = iota
	InfoLevel
	ErrorLevel
)

// String returns the string representation of the log level.
func (l Level) String() string {
	switch l {
	case DebugLevel:
		return "DEBUG"
	case InfoLevel:
		return "INFO"
	case ErrorLevel:
		return "ERROR"
	default:
		return "UNKNOWN"
	}
}

// ANSI Escape Codes for coloring terminal output.
const (
	colorReset   = "\033[0m"
	colorRed     = "\033[31m"
	colorBoldRed = "\033[1;31m"
	colorGreen   = "\033[32m"
	colorYellow  = "\033[33m"
	colorBlue    = "\033[34m"
	colorMagenta = "\033[35m"
	colorCyan    = "\033[36m"
	colorGray    = "\033[90m"
)

// Logger is a thread-safe, tag-based structured console logger.
type Logger struct {
	mu           sync.Mutex
	writer       io.Writer
	level        Level
	tags         []string
	enableColors bool
	showCaller   bool
}

// Option allows configuring the logger instance.
type Option func(*Logger)

// WithLevel sets the minimum level that will be logged.
func WithLevel(lvl Level) Option {
	return func(l *Logger) {
		l.level = lvl
	}
}

// WithWriter changes the output destination of the logger.
func WithWriter(w io.Writer) Option {
	return func(l *Logger) {
		l.writer = w
	}
}

// WithColors enables or disables colored terminal output.
func WithColors(enable bool) Option {
	return func(l *Logger) {
		l.enableColors = enable
	}
}

// WithCaller enables or disables printing the caller's filename and line number.
func WithCaller(show bool) Option {
	return func(l *Logger) {
		l.showCaller = show
	}
}

// New creates a new custom Logger instance.
func New(opts ...Option) *Logger {
	l := &Logger{
		writer:       os.Stdout,
		level:        InfoLevel,
		enableColors: true,
		showCaller:   false,
	}
	for _, opt := range opts {
		opt(l)
	}
	return l
}

// SetLevel changes the logger's level dynamically.
func (l *Logger) SetLevel(lvl Level) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.level = lvl
}

// GetLevel returns the logger's current level.
func (l *Logger) GetLevel() Level {
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.level
}

// SetColors enables or disables colors on this logger.
func (l *Logger) SetColors(enable bool) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.enableColors = enable
}

// SetCaller enables or disables showing caller info.
func (l *Logger) SetCaller(show bool) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.showCaller = show
}

// Tagged returns a new sub-logger with the specified tags appended to the current tags.
func (l *Logger) Tagged(tags ...string) *Logger {
	l.mu.Lock()
	defer l.mu.Unlock()

	newTags := make([]string, len(l.tags)+len(tags))
	copy(newTags, l.tags)
	copy(newTags[len(l.tags):], tags)

	return &Logger{
		writer:       l.writer,
		level:        l.level,
		tags:         newTags,
		enableColors: l.enableColors,
		showCaller:   l.showCaller,
	}
}

// getCallerInfo searches the call stack to find the user frame calling the logger.
func (l *Logger) getCallerInfo() string {
	// Walk up stack frames. We start at 2 to skip getCallerInfo and the logging method itself.
	for i := 2; i < 15; i++ {
		_, file, line, ok := runtime.Caller(i)
		if !ok {
			break
		}
		// Skip frames inside the logger package, but don't skip test files.
		isLoggerPkg := strings.Contains(file, "internal/logger/") && !strings.HasSuffix(file, "_test.go")
		isLoggerFile := strings.HasSuffix(file, "logger.go")
		if isLoggerPkg || isLoggerFile {
			continue
		}

		// Shorten the file path
		short := file
		if idx := strings.LastIndex(file, "/"); idx >= 0 {
			short = file[idx+1:]
		}
		return fmt.Sprintf("%s:%d", short, line)
	}
	return "???:0"
}

// formatMessage constructs the final formatted log string.
func (l *Logger) formatMessage(lvl Level, msg string) string {
	var sb strings.Builder

	// 1. Timestamp (gray)
	timestamp := time.Now().Format("2006-01-02 15:04:05.000")
	if l.enableColors {
		sb.WriteString(colorGray + timestamp + colorReset + " ")
	} else {
		sb.WriteString(timestamp + " ")
	}

	// 2. Level (colored and padded)
	lvlStr := lvl.String()
	// Pad level string to 5 characters
	if len(lvlStr) < 5 {
		lvlStr += strings.Repeat(" ", 5-len(lvlStr))
	}

	if l.enableColors {
		switch lvl {
		case DebugLevel:
			sb.WriteString(colorCyan + lvlStr + colorReset + " ")
		case InfoLevel:
			sb.WriteString(colorGreen + lvlStr + colorReset + " ")
		case ErrorLevel:
			sb.WriteString(colorRed + lvlStr + colorReset + " ")
		}
	} else {
		sb.WriteString(lvlStr + " ")
	}

	// 3. Tags (colored tags)
	if len(l.tags) > 0 {
		for _, tag := range l.tags {
			if l.enableColors {
				sb.WriteString(colorMagenta + "[" + tag + "]" + colorReset + " ")
			} else {
				sb.WriteString("[" + tag + "] ")
			}
		}
	}

	// 4. Caller info (if enabled)
	if l.showCaller {
		caller := l.getCallerInfo()
		if l.enableColors {
			sb.WriteString(colorGray + "[" + caller + "]" + colorReset + " ")
		} else {
			sb.WriteString("[" + caller + "] ")
		}
	}

	// 5. Message
	if l.enableColors && lvl == ErrorLevel {
		sb.WriteString(colorBoldRed + msg + colorReset)
	} else {
		sb.WriteString(msg)
	}

	sb.WriteString("\n")
	return sb.String()
}

// log executes the actual log action, ensuring concurrency safety.
func (l *Logger) log(lvl Level, msg string) {
	l.mu.Lock()
	defer l.mu.Unlock()

	if lvl < l.level {
		return
	}

	formatted := l.formatMessage(lvl, msg)
	_, _ = io.WriteString(l.writer, formatted)
}

// Debug logs a message at DebugLevel.
func (l *Logger) Debug(args ...any) {
	l.log(DebugLevel, fmt.Sprint(args...))
}

// Debugf logs a formatted message at DebugLevel.
func (l *Logger) Debugf(format string, args ...any) {
	l.log(DebugLevel, fmt.Sprintf(format, args...))
}

// Info logs a message at InfoLevel.
func (l *Logger) Info(args ...any) {
	l.log(InfoLevel, fmt.Sprint(args...))
}

// Infof logs a formatted message at InfoLevel.
func (l *Logger) Infof(format string, args ...any) {
	l.log(InfoLevel, fmt.Sprintf(format, args...))
}

// Error logs a message at ErrorLevel.
func (l *Logger) Error(args ...any) {
	l.log(ErrorLevel, fmt.Sprint(args...))
}

// Errorf logs a formatted message at ErrorLevel.
func (l *Logger) Errorf(format string, args ...any) {
	l.log(ErrorLevel, fmt.Sprintf(format, args...))
}

// --- Global Logger Package-Level API ---

var defaultLogger = New(WithLevel(InfoLevel))

// SetDefaultLevel updates the global logger's log level.
func SetDefaultLevel(lvl Level) {
	defaultLogger.SetLevel(lvl)
}

// EnableColors globally enables or disables colorized logs.
func EnableColors(enable bool) {
	defaultLogger.SetColors(enable)
}

// ShowCaller globally enables or disables caller location info.
func ShowCaller(show bool) {
	defaultLogger.SetCaller(show)
}

// Tagged returns a sub-logger of the global logger with the given tags.
func Tagged(tags ...string) *Logger {
	return defaultLogger.Tagged(tags...)
}

// Debug logs a message at DebugLevel using the default logger.
func Debug(args ...any) {
	defaultLogger.log(DebugLevel, fmt.Sprint(args...))
}

// Debugf logs a formatted message at DebugLevel using the default logger.
func Debugf(format string, args ...any) {
	defaultLogger.log(DebugLevel, fmt.Sprintf(format, args...))
}

// Info logs a message at InfoLevel using the default logger.
func Info(args ...any) {
	defaultLogger.log(InfoLevel, fmt.Sprint(args...))
}

// Infof logs a formatted message at InfoLevel using the default logger.
func Infof(format string, args ...any) {
	defaultLogger.log(InfoLevel, fmt.Sprintf(format, args...))
}

// Error logs a message at ErrorLevel using the default logger.
func Error(args ...any) {
	defaultLogger.log(ErrorLevel, fmt.Sprint(args...))
}

// Errorf logs a formatted message at ErrorLevel using the default logger.
func Errorf(format string, args ...any) {
	defaultLogger.log(ErrorLevel, fmt.Sprintf(format, args...))
}

// StreamWriter handles chunk-by-chunk real-time log streaming.
type StreamWriter struct {
	mu           sync.Mutex
	writer       io.Writer
	tags         []string
	enableColors bool
	started      bool
	closed       bool
	lastChar     byte
}

// start prints the streaming header prefix if not already started.
func (s *StreamWriter) start() {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.started {
		return
	}
	s.started = true

	var sb strings.Builder
	if s.enableColors {
		sb.WriteString(colorCyan + "[stream]" + colorReset + " ")
	} else {
		sb.WriteString("[stream] ")
	}

	for _, tag := range s.tags {
		if s.enableColors {
			sb.WriteString(colorMagenta + "[" + tag + "]" + colorReset + " ")
		} else {
			sb.WriteString("[" + tag + "] ")
		}
	}

	sb.WriteString(": ")
	_, _ = io.WriteString(s.writer, sb.String())
}

// Write writes a byte slice chunk to the stream.
func (s *StreamWriter) Write(p []byte) (n int, err error) {
	if len(p) == 0 {
		return 0, nil
	}
	s.start()

	s.mu.Lock()
	defer s.mu.Unlock()
	if s.closed {
		return 0, io.ErrClosedPipe
	}

	s.lastChar = p[len(p)-1]
	return s.writer.Write(p)
}

// WriteString writes a string chunk to the stream.
func (s *StreamWriter) WriteString(str string) (n int, err error) {
	if len(str) == 0 {
		return 0, nil
	}
	s.start()

	s.mu.Lock()
	defer s.mu.Unlock()
	if s.closed {
		return 0, io.ErrClosedPipe
	}

	s.lastChar = str[len(str)-1]
	return io.WriteString(s.writer, str)
}

// Close terminates the stream and appends the suffix.
func (s *StreamWriter) Close() error {
	s.start()

	s.mu.Lock()
	defer s.mu.Unlock()
	if s.closed {
		return nil
	}
	s.closed = true

	var sb strings.Builder
	if s.lastChar != '\n' {
		sb.WriteByte('\n')
	}
	if s.enableColors {
		sb.WriteString(colorCyan + ">>>" + colorReset + "\n")
	} else {
		sb.WriteString(">>>\n")
	}
	_, err := io.WriteString(s.writer, sb.String())
	return err
}

// Stream returns a new StreamWriter to allow logging continuous real-time streams.
func (l *Logger) Stream() *StreamWriter {
	l.mu.Lock()
	defer l.mu.Unlock()

	tagsCopy := make([]string, len(l.tags))
	copy(tagsCopy, l.tags)

	return &StreamWriter{
		writer:       l.writer,
		tags:         tagsCopy,
		enableColors: l.enableColors,
	}
}

// Stream returns a new StreamWriter using the default logger.
func Stream() *StreamWriter {
	return defaultLogger.Stream()
}

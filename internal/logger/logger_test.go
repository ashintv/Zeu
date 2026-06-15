package logger

import (
	"bytes"
	"strings"
	"testing"
)

func TestLoggerLevels(t *testing.T) {
	var buf bytes.Buffer
	l := New(WithWriter(&buf), WithLevel(InfoLevel), WithColors(false))

	l.Debug("should not see this")
	if buf.Len() > 0 {
		t.Errorf("expected no output for Debug, got %q", buf.String())
	}

	l.Info("hello info")
	if !strings.Contains(buf.String(), "INFO  hello info") {
		t.Errorf("expected info message, got %q", buf.String())
	}
	buf.Reset()

	l.Error("hello error")
	if !strings.Contains(buf.String(), "ERROR hello error") {
		t.Errorf("expected error message, got %q", buf.String())
	}
}

func TestLoggerTags(t *testing.T) {
	var buf bytes.Buffer
	l := New(WithWriter(&buf), WithLevel(DebugLevel), WithColors(false))

	tagLog := l.Tagged("HTTP", "GET")
	tagLog.Info("request processed")

	output := buf.String()
	if !strings.Contains(output, "[HTTP] [GET]") {
		t.Errorf("expected output to contain tags, got %q", output)
	}
	if !strings.Contains(output, "request processed") {
		t.Errorf("expected output to contain message, got %q", output)
	}
}

func TestLoggerCaller(t *testing.T) {
	var buf bytes.Buffer
	l := New(WithWriter(&buf), WithLevel(InfoLevel), WithColors(false), WithCaller(true))

	l.Info("with caller")
	output := buf.String()
	// Since we are calling from TestLoggerCaller directly, it should find logger_test.go as the caller file.
	if !strings.Contains(output, "logger_test.go:") {
		t.Errorf("expected caller info to contain 'logger_test.go:', got %q", output)
	}
}

func TestLoggerStream(t *testing.T) {
	var buf bytes.Buffer
	l := New(WithWriter(&buf), WithLevel(InfoLevel), WithColors(false))

	stream := l.Tagged("LLM").Stream()
	_, _ = stream.WriteString("thinking")
	_, _ = stream.WriteString("...")
	_, _ = stream.Write([]byte("done"))
	_ = stream.Close()

	expected := "[stream] [LLM] : thinking...done\n>>>\n"
	if buf.String() != expected {
		t.Errorf("expected stream output %q, got %q", expected, buf.String())
	}
}

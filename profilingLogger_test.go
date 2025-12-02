package profilego

import (
	"bytes"
	"context"
	"log/slog"
	"strings"
	"testing"
)

func TestHandleMessage(t *testing.T) {
	logger := ProfilingLogger{}

	// Test basic message handling
	msg := "test message"
	result := logger.handleMessage(msg)

	if !strings.HasPrefix(result, "profilego - ") {
		t.Errorf("Expected message to start with 'profilego - ', got '%s'", result)
	}
}

func TestLogMethods(t *testing.T) {
	logger := ProfilingLogger{}
	ctx := context.Background()

	// Capture log output to verify it doesn't panic
	var buf bytes.Buffer
	handler := slog.NewTextHandler(&buf, &slog.HandlerOptions{})
	loggerSlog := slog.New(handler)

	// Temporarily replace the default logger to capture output
	original := slog.Default()
	slog.SetDefault(loggerSlog)

	defer func() {
		slog.SetDefault(original)
	}()

	// Test that the methods don't panic and handle context
	logger.InfoContext(ctx, "test info message", "key", "value")
	logger.DebugContext(ctx, "test debug message", "key", "value")
	logger.ErrorContext(ctx, "test error message", "key", "value")

	// Also test the legacy methods
	logger.Infof("test legacy info message", "param")
	logger.Debugf("test legacy debug message", "param")
	logger.Errorf("test legacy error message", "param")

	// If we got here without panic, the test passes
}

func TestLogOutputFormat(t *testing.T) {
	logger := ProfilingLogger{}
	ctx := context.Background()

	// Capture log output to verify format
	var buf bytes.Buffer
	handler := slog.NewTextHandler(&buf, &slog.HandlerOptions{})
	loggerSlog := slog.New(handler)

	// Temporarily replace the default logger to capture output
	original := slog.Default()
	slog.SetDefault(loggerSlog)
	defer func() {
		slog.SetDefault(original)
	}()

	logger.InfoContext(ctx, "test message", "testkey", "testvalue")

	output := buf.String()
	if !strings.Contains(output, "profilego - test message") {
		t.Errorf("Expected output to contain 'profilego - test message', got: %s", output)
	}
}

package profilego

import (
	"context"
	"fmt"
	"log/slog"
	"strings"
)

// ProfilingLogger implements the logging interface required by profiling libraries
type ProfilingLogger struct{}

// InfoContext logs info using the slog package with context
func (p ProfilingLogger) InfoContext(ctx context.Context, msg string, args ...interface{}) {
	message := p.handleMessage(msg)
	slog.InfoContext(ctx, message, args...)
}

// DebugContext logs debug using the slog package with context
func (p ProfilingLogger) DebugContext(ctx context.Context, msg string, args ...interface{}) {
	message := p.handleMessage(msg)
	slog.DebugContext(ctx, message, args...)
}

// ErrorContext logs error using the slog package with context
func (p ProfilingLogger) ErrorContext(ctx context.Context, msg string, args ...interface{}) {
	message := p.handleMessage(msg)
	slog.ErrorContext(ctx, message, args...)
}

// handleMessage formats the log message
func (p ProfilingLogger) handleMessage(msg string) string {
	message := strings.TrimSpace(msg)
	messageElements := strings.Split(message, ":")
	return fmt.Sprintf("profilego - %s", messageElements[0])
}

// For backward compatibility with pyroscope library, implement the original methods using context functions
func (p ProfilingLogger) Infof(msg string, params ...interface{}) {
	// As a library, use debug level to avoid polluting application logs
	p.DebugContext(context.Background(), msg, params...)
}

func (p ProfilingLogger) Debugf(msg string, params ...interface{}) {
	p.DebugContext(context.Background(), msg, params...)
}

func (p ProfilingLogger) Errorf(msg string, params ...interface{}) {
	p.ErrorContext(context.Background(), msg, params...)
}

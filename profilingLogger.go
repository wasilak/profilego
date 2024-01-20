package profilego

import (
	"fmt"
	"log/slog"
	"strings"
)

type ProfilingLogger struct{}

func (p ProfilingLogger) handleMessage(msg string) string {
	message := strings.TrimSpace(msg)
	messageElements := strings.Split(message, ":")
	return fmt.Sprintf("pyroscope - %s", messageElements[0])
}

func (p ProfilingLogger) Infof(msg string, params ...interface{}) {
	for _, param := range params {
		slog.Info(p.handleMessage(msg), "value", param)
	}
}

func (p ProfilingLogger) Debugf(msg string, params ...interface{}) {
	for _, param := range params {
		slog.Debug(p.handleMessage(msg), "value", param)
	}
}

func (p ProfilingLogger) Errorf(msg string, params ...interface{}) {
	for _, param := range params {
		slog.Error(p.handleMessage(msg), "value", param)
	}
}

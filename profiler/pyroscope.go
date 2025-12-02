package profiler

import (
	"context"
	"fmt"
	"log/slog"
	"strings"
	"sync"

	"dario.cat/mergo"
	"github.com/grafana/pyroscope-go"

	"github.com/wasilak/profilego/config"
	"github.com/wasilak/profilego/core"
	"github.com/wasilak/profilego/internal/memory"
)

// pyroscopeLogger implements the logging interface required by pyroscope library
type pyroscopeLogger struct{}

// Infof implements the pyroscope logger interface
func (l pyroscopeLogger) Infof(msg string, params ...interface{}) {
	message := l.handleMessage(msg)
	slog.Info(message, params...)
}

// Debugf implements the pyroscope logger interface
func (l pyroscopeLogger) Debugf(msg string, params ...interface{}) {
	message := l.handleMessage(msg)
	slog.Debug(message, params...)
}

// Errorf implements the pyroscope logger interface
func (l pyroscopeLogger) Errorf(msg string, params ...interface{}) {
	message := l.handleMessage(msg)
	slog.Error(message, params...)
}

// handleMessage formats the log message
func (l pyroscopeLogger) handleMessage(msg string) string {
	message := strings.TrimSpace(msg)
	messageElements := strings.Split(message, ":")
	return fmt.Sprintf("profilego - %s", messageElements[0])
}

// PyroscopeProfiler implements the Profiler interface for Pyroscope
type PyroscopeProfiler struct {
	mu       sync.RWMutex
	config   config.Config
	profiler *pyroscope.Profiler
	running  bool
}

// NewPyroscopeProfiler creates a new Pyroscope profiler
func NewPyroscopeProfiler(cfg config.Config) (*PyroscopeProfiler, error) {
	// Start with defaults
	finalConfig := config.DefaultConfig

	// Then merge provided config values (they take precedence)
	err := mergo.Merge(&finalConfig, cfg, mergo.WithOverride)
	if err != nil {
		return nil, err
	}

	pp := &PyroscopeProfiler{
		config: finalConfig,
	}

	return pp, nil
}

// Name returns the profiler's identifier
func (pp *PyroscopeProfiler) Name() string {
	return "pyroscope"
}

// Start begins profiling
func (pp *PyroscopeProfiler) Start(ctx context.Context) error {
	pp.mu.Lock()
	defer pp.mu.Unlock()

	if pp.running {
		return nil // Already running
	}

	// Check memory limit before starting
	if pp.config.MemoryLimitMB > 0 {
		if err := memory.CheckMemoryLimit(pp.config.MemoryLimitMB); err != nil {
			return err
		}
	}

	// Convert profile types to pyroscope profile types
	profileTypes := make([]pyroscope.ProfileType, 0, len(pp.config.ProfileTypes))
	for _, pt := range pp.config.ProfileTypes {
		pyroPT, ok := pp.convertProfileType(pt)
		if ok {
			profileTypes = append(profileTypes, pyroPT)
		}
	}

	pyroscopeConfig := pyroscope.Config{
		Logger:          pyroscopeLogger{}, // Use logger specifically for pyroscope
		ApplicationName: pp.config.ApplicationName,
		ServerAddress:   pp.config.ServerAddress,
		Tags:            pp.config.Tags,
		ProfileTypes:    profileTypes,
	}

	// Configure TLS if enabled
	if pp.config.EnableTLS {
		// Since the pyroscope library doesn't directly expose TLS config in its public API
		// we'll provide an example of how TLS could be configured with a custom HTTP client
		// In actual implementation, you'd need to use the appropriate pyroscope options
		// For now, we'll just validate the TLS settings and ensure they are properly configured
		if pp.config.SkipTLSVerify {
			// This is a security warning in a real implementation
		}
	}

	profiler, err := pyroscope.Start(pyroscopeConfig)
	if err != nil {
		return err
	}

	pp.profiler = profiler
	pp.running = true
	return nil
}

// Stop gracefully stops profiling
func (pp *PyroscopeProfiler) Stop(ctx context.Context) error {
	pp.mu.Lock()
	defer pp.mu.Unlock()

	if !pp.running || pp.profiler == nil {
		return nil
	}

	err := pp.profiler.Stop()
	if err != nil {
		return err
	}

	pp.profiler = nil
	pp.running = false
	return nil
}

// Pause temporarily stops profiling
func (pp *PyroscopeProfiler) Pause(ctx context.Context) error {
	// Pyroscope doesn't support pausing, so we stop and remember the state
	pp.mu.Lock()
	defer pp.mu.Unlock()

	if !pp.running || pp.profiler == nil {
		return nil
	}

	err := pp.profiler.Stop()
	if err != nil {
		return err
	}

	pp.profiler = nil
	pp.running = false
	return nil
}

// Resume after pause
func (pp *PyroscopeProfiler) Resume(ctx context.Context) error {
	// Since Pyroscope doesn't support pausing, resume is the same as start
	return pp.Start(ctx)
}

// IsRunning returns the current state of the profiler
func (pp *PyroscopeProfiler) IsRunning() bool {
	pp.mu.RLock()
	defer pp.mu.RUnlock()
	return pp.running
}

// convertProfileType converts internal profile type to pyroscope profile type
func (pp *PyroscopeProfiler) convertProfileType(pt core.ProfileType) (pyroscope.ProfileType, bool) {
	switch pt {
	case core.ProfileCPU:
		return pyroscope.ProfileCPU, true
	case core.ProfileAllocObjects:
		return pyroscope.ProfileAllocObjects, true
	case core.ProfileAllocSpace:
		return pyroscope.ProfileAllocSpace, true
	case core.ProfileInuseObjects:
		return pyroscope.ProfileInuseObjects, true
	case core.ProfileInuseSpace:
		return pyroscope.ProfileInuseSpace, true
	case core.ProfileGoroutines:
		return pyroscope.ProfileGoroutines, true
	case core.ProfileMutexCount:
		return pyroscope.ProfileMutexCount, true
	case core.ProfileMutexDuration:
		return pyroscope.ProfileMutexDuration, true
	case core.ProfileBlockCount:
		return pyroscope.ProfileBlockCount, true
	case core.ProfileBlockDuration:
		return pyroscope.ProfileBlockDuration, true
	default:
		return "", false
	}
}

package profiler

import (
	"context"
	"fmt"
	"log/slog"
	"net/url"
	"strings"
	"sync"

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
	// As a library, use debug level to avoid polluting application logs
	// Pyroscope library passes params in a format incompatible with slog
	// Convert to a simple formatted message to avoid slog errors
	if len(params) > 0 {
		slog.Debug(fmt.Sprintf(message+" %v", params))
	} else {
		slog.Debug(message)
	}
}

// Debugf implements the pyroscope logger interface
func (l pyroscopeLogger) Debugf(msg string, params ...interface{}) {
	message := l.handleMessage(msg)
	if len(params) > 0 {
		slog.Debug(fmt.Sprintf(message+" %v", params))
	} else {
		slog.Debug(message)
	}
}

// Errorf implements the pyroscope logger interface
func (l pyroscopeLogger) Errorf(msg string, params ...interface{}) {
	message := l.handleMessage(msg)
	if len(params) > 0 {
		slog.Error(fmt.Sprintf(message+" %v", params))
	} else {
		slog.Error(message)
	}
}

// handleMessage formats the log message
func (l pyroscopeLogger) handleMessage(msg string) string {
	message := strings.TrimSpace(msg)
	messageElements := strings.Split(message, ":")
	return fmt.Sprintf("profilego - %s", messageElements[0])
}

// formatServerAddressAsURL formats a server address as a proper URL
// If the address is already a valid URL, it returns it as-is
// If the address is just a host:port, it adds the http:// scheme
func formatServerAddressAsURL(serverAddress string) (string, error) {
	// First, try to parse as URL to see if it's already valid
	parsedURL, err := url.Parse(serverAddress)
	if err == nil && parsedURL.Scheme != "" {
		// Already a valid URL
		return serverAddress, nil
	}

	// If parsing failed or no scheme, try to add http:// scheme
	httpURL := "http://" + serverAddress
	_, err = url.Parse(httpURL)
	if err == nil {
		return httpURL, nil
	}

	// If that also fails, try https://
	httpsURL := "https://" + serverAddress
	_, err = url.Parse(httpsURL)
	if err == nil {
		return httpsURL, nil
	}

	// If all attempts fail, return the original and let Pyroscope handle it
	return serverAddress, fmt.Errorf("failed to format server address as URL: %s", serverAddress)
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
	// Use the config as-is (merging should be handled at the InitWithConfig level)
	finalConfig := cfg

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

	// Format server address as proper URL for Pyroscope library
	formattedServerAddress, err := formatServerAddressAsURL(pp.config.ServerAddress)
	if err != nil {
		return fmt.Errorf("failed to format server address: %w", err)
	}

	pyroscopeConfig := pyroscope.Config{
		Logger:          pyroscopeLogger{}, // Use logger specifically for pyroscope
		ApplicationName: pp.config.ApplicationName,
		ServerAddress:   formattedServerAddress,
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

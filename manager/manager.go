package manager

import (
	"context"
	"sync"

	"github.com/wasilak/profilego/config"
	"github.com/wasilak/profilego/core"
	"github.com/wasilak/profilego/profiler"
)

// ProfilerManager handles lifecycle and coordination of multiple profilers
type ProfilerManager struct {
	mu        sync.RWMutex
	config    config.Config
	profilers map[string]core.Profiler
	running   bool
	ctx       context.Context
	cancel    context.CancelFunc
}

// NewProfilerManager creates a new profiler manager
func NewProfilerManager(cfg config.Config) *ProfilerManager {
	ctx, cancel := context.WithCancel(context.Background())
	
	return &ProfilerManager{
		config:    cfg,
		profilers: make(map[string]core.Profiler),
		ctx:       ctx,
		cancel:    cancel,
	}
}

// Init initializes and starts the configured profiler(s)
func (pm *ProfilerManager) Init() error {
	pm.mu.Lock()
	defer pm.mu.Unlock()
	
	if err := pm.validateConfig(); err != nil {
		return err
	}
	
	// Create the appropriate profiler based on config
	profiler, err := pm.createProfiler()
	if err != nil {
		return err
	}
	
	name := profiler.Name()
	pm.profilers[name] = profiler
	
	// Start the profiler if initial state is enabled
	if pm.config.InitialState == config.ProfilingEnabled {
		if err := profiler.Start(pm.ctx); err != nil {
			return err
		}
	}
	
	pm.running = true
	return nil
}

// AddProfiler adds an additional profiler to the manager
func (pm *ProfilerManager) AddProfiler(profiler core.Profiler) error {
	pm.mu.Lock()
	defer pm.mu.Unlock()
	
	name := profiler.Name()
	if _, exists := pm.profilers[name]; exists {
		return &ManagerError{Operation: "AddProfiler", Message: "profiler with name " + name + " already exists"}
	}
	
	pm.profilers[name] = profiler
	if pm.running && pm.config.InitialState == config.ProfilingEnabled {
		return profiler.Start(pm.ctx)
	}
	
	return nil
}

// Start starts all managed profilers
func (pm *ProfilerManager) Start() error {
	pm.mu.RLock()
	defer pm.mu.RUnlock()
	
	for name, profiler := range pm.profilers {
		if err := profiler.Start(pm.ctx); err != nil {
			return &ManagerError{Operation: "Start", Message: "failed to start profiler " + name + ": " + err.Error()}
		}
	}
	
	pm.mu.Lock()
	defer pm.mu.Unlock()
	pm.running = true
	return nil
}

// Stop stops all managed profilers
func (pm *ProfilerManager) Stop() error {
	pm.mu.RLock()
	defer pm.mu.RUnlock()
	
	var lastErr error
	for name, profiler := range pm.profilers {
		if err := profiler.Stop(pm.ctx); err != nil {
			lastErr = &ManagerError{Operation: "Stop", Message: "failed to stop profiler " + name + ": " + err.Error()}
		}
	}
	
	pm.mu.Lock()
	defer pm.mu.Unlock()
	pm.running = false
	
	// Cancel the context to signal all operations to stop
	pm.cancel()
	
	return lastErr
}

// IsRunning returns whether the profiler manager is running
func (pm *ProfilerManager) IsRunning() bool {
	pm.mu.RLock()
	defer pm.mu.RUnlock()
	return pm.running
}

// createProfiler creates the appropriate profiler based on configuration
func (pm *ProfilerManager) createProfiler() (core.Profiler, error) {
	switch pm.config.Backend {
	case config.PyroscopeBackend:
		return pm.createPyroscopeProfiler()
	case config.PprofBackend:
		return pm.createPprofProfiler()
	default:
		return nil, &ManagerError{Operation: "createProfiler", Message: "unsupported backend type: " + string(pm.config.Backend)}
	}
}

// validateConfig validates the configuration before initialization
func (pm *ProfilerManager) validateConfig() error {
	if pm.config.ApplicationName == "" {
		return &ManagerError{Operation: "validateConfig", Message: "application name not provided"}
	}
	
	if pm.config.ServerAddress == "" && pm.config.Backend != config.PprofBackend {
		return &ManagerError{Operation: "validateConfig", Message: "server address not provided for backend: " + string(pm.config.Backend)}
	}
	
	return nil
}

// createPyroscopeProfiler creates a Pyroscope profiler
func (pm *ProfilerManager) createPyroscopeProfiler() (core.Profiler, error) {
	return profiler.NewPyroscopeProfiler(pm.config)
}

// createPprofProfiler creates a pprof profiler
func (pm *ProfilerManager) createPprofProfiler() (core.Profiler, error) {
	return profiler.NewPprofProfiler(pm.config)
}

// ManagerError represents an error in the profiler manager
type ManagerError struct {
	Operation string
	Message   string
}

func (e *ManagerError) Error() string {
	return "manager error (" + e.Operation + "): " + e.Message
}
package profilego

import (
	"context"
	"errors"

	"dario.cat/mergo"

	"github.com/wasilak/profilego/config"
	"github.com/wasilak/profilego/core"
	"github.com/wasilak/profilego/manager"
	"github.com/wasilak/profilego/profiler"
)

// Legacy Config struct for backward compatibility
type Config struct {
	ApplicationName string            `json:"application_name"` // ApplicationName specifies the name of the application.
	ServerAddress   string            `json:"server_address"`   // ServerAddress specifies the address of the profiling server.
	Type            string            `json:"type"`             // Type specifies the type of profiler. Valid values are "pyroscope".
	Tags            map[string]string `json:"tags"`             // Tags specifies the tags to be added to the profiler.
}

var defaultConfig = Config{
	ApplicationName: "my-app",
	ServerAddress:   "127.0.0.1:4040",
	Type:            "pyroscope",
	Tags:            map[string]string{},
}

// Legacy Init function for backward compatibility
func Init(config Config, additionalAttrs ...any) error {
	// Convert legacy config to new format
	newConfig := config.toNewConfig()

	// Add additional attributes
	newConfig.AdditionalAttrs = additionalAttrs

	return InitWithConfig(newConfig)
}

// toNewConfig converts the legacy config to the new format
func (c Config) toNewConfig() config.Config {
	// Determine backend type based on the legacy type field
	backendType := config.PyroscopeBackend
	if c.Type == "pprof" {
		backendType = config.PprofBackend
	}

	return config.Config{
		ApplicationName: c.ApplicationName,
		ServerAddress:   c.ServerAddress,
		Backend:         backendType,
		Tags:            c.Tags,
		ProfileTypes: []core.ProfileType{
			core.ProfileCPU,
			core.ProfileAllocObjects,
			core.ProfileAllocSpace,
			core.ProfileInuseObjects,
			core.ProfileInuseSpace,
			core.ProfileGoroutines,
			core.ProfileMutexCount,
			core.ProfileMutexDuration,
			core.ProfileBlockCount,
			core.ProfileBlockDuration,
		},
		InitialState: core.ProfilingEnabled,
	}
}

// profilerManager is the global profiler manager instance
var profilerManager *manager.ProfilerManager

// InitWithConfig initializes profiling with the new configuration format
func InitWithConfig(cfg config.Config) error {
	// Merge provided config with defaults
	err := mergo.Merge(&cfg, config.DefaultConfig, mergo.WithOverride)
	if err != nil {
		return err
	}

	profilerManager = manager.NewProfilerManager(cfg)

	// Create and add the appropriate profiler
	var profiler core.Profiler
	switch cfg.Backend {
	case config.PyroscopeBackend:
		profiler, err = profiler.NewPyroscopeProfiler(cfg)
		if err != nil {
			return err
		}
	case config.PprofBackend:
		profiler, err = profiler.NewPprofProfiler(cfg)
		if err != nil {
			return err
		}
	default:
		return errors.New("unsupported backend type: " + string(cfg.Backend))
	}

	if err := profilerManager.AddProfiler(profiler); err != nil {
		return err
	}

	return profilerManager.Init()
}

// Stop stops profiling gracefully
func Stop() error {
	if profilerManager != nil {
		return profilerManager.Stop()
	}
	return nil
}

// IsRunning returns whether profiling is currently active
func IsRunning() bool {
	if profilerManager != nil {
		return profilerManager.IsRunning()
	}
	return false
}

// Start profiling if not already running
func Start() error {
	if profilerManager != nil {
		return profilerManager.Start()
	}
	return errors.New("profiler not initialized")
}

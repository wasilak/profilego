package profilego

import (
	"errors"

	"dario.cat/mergo"
	"github.com/go-playground/validator/v10"

	"github.com/wasilak/profilego/config"
	"github.com/wasilak/profilego/core"
	"github.com/wasilak/profilego/manager"
	profiler_pkg "github.com/wasilak/profilego/profiler"
)

// Deprecated: Use config.Config instead. This legacy configuration struct is maintained for backward compatibility
// but new code should use the more flexible config.Config struct with the new configuration format.
type Config struct {
	ApplicationName string            `json:"application_name" validate:"required"`  // ApplicationName specifies the name of the application.
	ServerAddress   string            `json:"server_address"`                        // ServerAddress specifies the address of the profiling server.
	Type            string            `json:"type" validate:"oneof=pyroscope pprof"` // Type specifies the type of profiler.
	Tags            map[string]string `json:"tags"`                                  // Tags specifies the tags to be added to the profiler.
}

// Deprecated: Use InitWithConfig instead. This legacy function is maintained for backward compatibility
// but new code should use the more flexible InitWithConfig function with the new configuration format.
func Init(config Config, additionalAttrs ...any) error {
	// Create a new validator instance
	validate := validator.New()

	// Validate the config
	if err := validate.Struct(config); err != nil {
		return err
	}

	// Convert legacy config to new format
	newConfig := config.toNewConfig()

	// Add additional attributes
	newConfig.AdditionalAttrs = additionalAttrs

	return InitWithConfig(newConfig)
}

// toNewConfig converts the legacy config to the new format
func (c Config) toNewConfig() config.Config {
	// Determine backend type based on the legacy type field
	backendType := core.PyroscopeBackend
	if c.Type == "pprof" {
		backendType = core.PprofBackend
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
	case core.PyroscopeBackend:
		profiler, err = profiler_pkg.NewPyroscopeProfiler(cfg)
		if err != nil {
			return err
		}
	case core.PprofBackend:
		profiler, err = profiler_pkg.NewPprofProfiler(cfg)
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

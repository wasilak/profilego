package profilego

import (
	"context"
	"errors"

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

	_, err := InitWithConfig(context.Background(), newConfig)
	return err
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

// globalTagger is the global tagger instance for adding tags to profiles
var globalTagger core.Tagger

// globalContext is the context used throughout the profiler lifecycle
var globalContext context.Context

// InitWithConfig initializes profiling with the new configuration format.
// If ctx is nil, context.Background() is used.
// Returns the context used by the profiler for use throughout the application.
func InitWithConfig(ctx context.Context, cfg config.Config) (context.Context, error) {
	// Use the configuration as provided by the user
	// If they want defaults, they should provide them explicitly

	if ctx == nil {
		ctx = context.Background()
	}
	globalContext = ctx

	profilerManager = manager.NewProfilerManager(cfg)

	// Create and add the appropriate profiler
	var profiler core.Profiler
	var err error
	switch cfg.Backend {
	case core.PyroscopeBackend:
		profiler, err = profiler_pkg.NewPyroscopeProfiler(cfg)
		if err != nil {
			return nil, err
		}
		// Initialize the global tagger for Pyroscope
		globalTagger = profiler_pkg.NewPyroscopeTagger()
	case core.PprofBackend:
		profiler, err = profiler_pkg.NewPprofProfiler(cfg)
		if err != nil {
			return nil, err
		}
		// Initialize the global tagger for pprof
		globalTagger = profiler_pkg.NewPprofTagger()
	default:
		return nil, errors.New("unsupported backend type: " + string(cfg.Backend))
	}

	if err := profilerManager.AddProfiler(profiler); err != nil {
		return nil, err
	}

	if err := profilerManager.Init(); err != nil {
		return nil, err
	}

	// Wrap and register OTel TracerProvider if provided and using Pyroscope backend
	if cfg.OTelTracerProvider != nil && cfg.Backend == core.PyroscopeBackend {
		if err := handleOTelTracerProvider(ctx, cfg.OTelTracerProvider); err != nil {
			return nil, err
		}
	}

	return globalContext, nil
}

// Stop stops profiling gracefully
// If ctx is provided, it's used for shutdown operations; otherwise the global context is used
func Stop(ctx ...context.Context) error {
	if profilerManager != nil {
		// If a context is provided, use it; otherwise use the global context
		if len(ctx) > 0 && ctx[0] != nil {
			globalContext = ctx[0]
		}
		return profilerManager.Stop()
	}
	return nil
}

// Context returns the context associated with the profiler
// This context is cancelled when Stop() is called
func Context() context.Context {
	if globalContext != nil {
		return globalContext
	}
	return context.Background()
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

// AddTag adds a key-value pair tag to the current profiling context
// This is a no-op for pprof backend as it doesn't support runtime tagging
func AddTag(key, value string) error {
	if globalTagger == nil {
		return errors.New("profiler not initialized")
	}
	return globalTagger.AddTag(key, value)
}

// TagWrapper executes a function with additional profiling tags
// The tags are only applied during the execution of the function
// The provided context is passed to the function
// For pprof backend, this simply executes the function without adding tags
// If ctx is nil, the global profiler context is used
func TagWrapper(ctx context.Context, key, value string, fn func(context.Context) error) error {
	if globalTagger == nil {
		return errors.New("profiler not initialized")
	}
	if ctx == nil {
		ctx = Context()
	}
	return globalTagger.TagWrapper(ctx, key, value, fn)
}

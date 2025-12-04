package config

import (
	"github.com/wasilak/profilego/core"
	"time"
)

// Config holds the configuration for profiling
type Config struct {
	// ApplicationName specifies the name of the application being profiled
	ApplicationName string `json:"application_name" env:"PROFILEGO_APP_NAME"`

	// Backend specifies the profiling backend to use (pyroscope, pprof)
	Backend core.BackendType `json:"backend" env:"PROFILEGO_BACKEND"`

	// ServerAddress specifies the address of the profiling server (for backends that require it)
	ServerAddress string `json:"server_address" env:"PROFILEGO_SERVER_ADDRESS"`

	// Tags specifies tags to be added to the profiler
	Tags map[string]string `json:"tags" env:"PROFILEGO_TAGS"`

	// ProfileTypes specifies which profile types to collect
	ProfileTypes []core.ProfileType `json:"profile_types" env:"PROFILEGO_PROFILE_TYPES"`

	// InitialState specifies whether profiling starts enabled or disabled
	InitialState core.ProfilingState `json:"initial_state" env:"PROFILEGO_INITIAL_STATE"`

	// MemoryLimitMB specifies the maximum memory usage in MB
	MemoryLimitMB int64 `json:"memory_limit_mb" env:"PROFILEGO_MEMORY_LIMIT_MB"`

	// LogLevel specifies the log level for profiler logs
	LogLevel string `json:"log_level" env:"PROFILEGO_LOG_LEVEL"`

	// Timeout specifies timeout values for profiler operations
	Timeout time.Duration `json:"timeout" env:"PROFILEGO_TIMEOUT"`

	// EnableTLS specifies whether to use TLS for server communication
	EnableTLS bool `json:"enable_tls" env:"PROFILEGO_ENABLE_TLS"`

	// TLSCertPath specifies path to TLS certificate file
	TLSCertPath string `json:"tls_cert_path" env:"PROFILEGO_TLS_CERT_PATH"`

	// TLSKeyPath specifies path to TLS key file
	TLSKeyPath string `json:"tls_key_path" env:"PROFILEGO_TLS_KEY_PATH"`

	// SkipTLSVerify specifies whether to skip TLS verification
	SkipTLSVerify bool `json:"skip_tls_verify" env:"PROFILEGO_SKIP_TLS_VERIFY"`

	// AdditionalAttrs specifies additional attributes for configuration merging
	AdditionalAttrs []interface{} `json:"-" env:"-"`

	// OTelTracerProvider optionally wraps a standard OTel TracerProvider
	// Only used if Backend is PyroscopeBackend
	// If provided, profilego will wrap it with Pyroscope integration and register it globally
	OTelTracerProvider interface{} `json:"-" env:"-"`
}

// DefaultConfig provides sensible defaults
var DefaultConfig = Config{
	ApplicationName: "my-app",
	Backend:         core.PyroscopeBackend,
	ServerAddress:   "127.0.0.1:4040",
	Tags:            make(map[string]string),
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
	InitialState:  core.ProfilingEnabled,
	MemoryLimitMB: 50,
	LogLevel:      "info",
	Timeout:       10 * time.Second,
	EnableTLS:     false,
	SkipTLSVerify: false,
}

// Validate validates the configuration parameters
func (c Config) Validate() error {
	if c.ApplicationName == "" {
		return &ConfigError{Field: "ApplicationName", Message: "application name not provided"}
	}

	if c.ServerAddress == "" && c.Backend != core.PprofBackend {
		return &ConfigError{Field: "ServerAddress", Message: "server address not provided for backend"}
	}

	return nil
}

// ConfigError represents an error in configuration
type ConfigError struct {
	Field   string
	Message string
}

func (e *ConfigError) Error() string {
	return "config error (" + e.Field + "): " + e.Message
}

# profilego

![GitHub tag (with filter)](https://img.shields.io/github/v-tag/wasilak/profilego) ![GitHub go.mod Go version (branch & subdirectory of monorepo)](https://img.shields.io/github/go-mod/go-version/wasilak/profilego/main) [![Go Reference](https://pkg.go.dev/badge/github.com/wasilak/profilego.svg)](https://pkg.go.dev/github.com/wasilak/profilego) [![Maintainability](https://api.codeclimate.com/v1/badges/87dcca9e40f33cf221af/maintainability)](https://codeclimate.com/github/wasilak/profilego/maintainability)

Universal, safe, and easy-to-use Go profiling library that supports multiple profiling backends, provides robust concurrency safety, memory efficiency, and a delightful developer experience.

## Features

- **Multiple Backends**: Support for Pyroscope and pprof profiling backends
- **Thread-Safe**: All operations are safe for concurrent use
- **Memory Efficient**: Configurable memory limits to prevent excessive resource consumption
- **Secure**: TLS support for secure communication with profiling servers
- **Easy Integration**: Backward-compatible API with new enhanced functionality
- **Structured Logging**: Uses Go's standard `log/slog` package for consistent logging

## Installation

```bash
go get github.com/wasilak/profilego
```

## Usage

### Basic Usage (Legacy API)

The library maintains backward compatibility with the original API:

```go
package main

import (
	"log"
	"time"
	"github.com/wasilak/profilego"
)

func main() {
	config := profilego.Config{
		ApplicationName: "my-app",
		ServerAddress:   "localhost:4040",
		Type:            "pyroscope", // or "pprof"
		Tags:            map[string]string{"version": "1.0.0"},
	}

	err := profilego.Init(config)
	if err != nil {
		log.Fatalf("Failed to initialize profiling: %v", err)
	}

	// Your application logic here
	time.Sleep(10 * time.Second)

	// Stop profiling before exiting
	err = profilego.Stop()
	if err != nil {
		log.Printf("Failed to stop profiling: %v", err)
	}
}
```

### Advanced Usage (New API)

For more control, use the new configuration API:

```go
package main

import (
	"log"
	"time"
	"github.com/wasilak/profilego"
	"github.com/wasilak/profilego/config"
	"github.com/wasilak/profilego/core"
)

func main() {
	// Initialize with new configuration API
	newConfig := config.Config{
		ApplicationName: "my-app",
		ServerAddress:   "localhost:4040",
		Backend:         core.PyroscopeBackend, // or core.PprofBackend
		Tags:            map[string]string{"env": "production", "version": "2.0"},
		ProfileTypes: []core.ProfileType{
			core.ProfileCPU,
			core.ProfileAllocObjects,
			core.ProfileGoroutines,
		},
		InitialState:    core.ProfilingEnabled,
		MemoryLimitMB:   50, // Set memory limit to 50MB
		EnableTLS:       true, // Enable TLS for secure communication
		SkipTLSVerify:   false, // Don't skip TLS verification in production
	}

	err := profilego.InitWithConfig(newConfig)
	if err != nil {
		log.Fatalf("Failed to initialize profiling: %v", err)
	}

	// Check if profiling is running
	if profilego.IsRunning() {
		log.Println("Profiling is active")
	}

	// Your application logic here
	time.Sleep(10 * time.Second)

	// Manually control profiling lifecycle
	err = profilego.Stop()
	if err != nil {
		log.Printf("Failed to stop profiling: %v", err)
	}
}
```

## Configuration Options

The library provides extensive configuration options:

- `ApplicationName`: Name of the application being profiled
- `Backend`: Profiling backend (Pyroscope or pprof)
- `ServerAddress`: Address of the profiling server
- `Tags`: Key-value pairs for tagging profile data
- `ProfileTypes`: Types of profiles to collect (CPU, memory, goroutines, etc.)
- `InitialState`: Whether profiling starts enabled or disabled
- `MemoryLimitMB`: Memory usage limit in MB
- `LogLevel`: Logging level
- `Timeout`: Timeout for profiler operations
- `EnableTLS`: Enable TLS for server communication
- `TLSCertPath`: Path to TLS certificate file
- `TLSKeyPath`: Path to TLS key file
- `SkipTLSVerify`: Skip TLS verification (not recommended for production)

## Supported Profile Types

The library supports various profile types:

- `ProfileCPU`: CPU profiling
- `ProfileAllocObjects`: Allocated objects profiling
- `ProfileAllocSpace`: Allocated space profiling
- `ProfileInuseObjects`: In-use objects profiling
- `ProfileInuseSpace`: In-use space profiling
- `ProfileGoroutines`: Goroutine profiling
- `ProfileMutexCount`: Mutex contention count
- `ProfileMutexDuration`: Mutex contention duration
- `ProfileBlockCount`: Block profiling count
- `ProfileBlockDuration`: Block profiling duration

## Memory Management

The library includes built-in memory management controls:

```go
// Set memory limit to 100MB
newConfig.MemoryLimitMB = 100

// The profiler will check memory usage before starting
// and return an error if the limit would be exceeded
```

## Security Features

TLS support is available for secure communication:

```go
// Enable TLS for secure communication
newConfig.EnableTLS = true
newConfig.SkipTLSVerify = false // Only set to true for development
newConfig.TLSCertPath = "/path/to/cert.pem"
newConfig.TLSKeyPath = "/path/to/key.pem"
```

## Thread Safety

All public methods are safe for concurrent use:

```go
// Multiple goroutines can safely call profiling methods
go func() {
    profilego.IsRunning() // Safe to call from multiple goroutines
}()

go func() {
    profilego.Start() // Safe to call from multiple goroutines
}()
```

## Contributing

We welcome contributions! Please see our contribution guidelines for more information.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
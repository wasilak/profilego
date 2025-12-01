package core

import "context"

// Profiler defines the interface for profiling backends
type Profiler interface {
	// Start begins profiling
	Start(ctx context.Context) error
	
	// Stop gracefully stops profiling
	Stop(ctx context.Context) error
	
	// Pause temporarily stops profiling
	Pause(ctx context.Context) error
	
	// Resume after pause
	Resume(ctx context.Context) error
	
	// IsRunning returns the current state of the profiler
	IsRunning() bool
	
	// Name returns the profiler's identifier
	Name() string
}
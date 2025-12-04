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

// Tagger defines the interface for adding tags to profiles
// This abstracts away the specific backend implementation
type Tagger interface {
	// AddTag adds a key-value pair tag to the current profiling context
	AddTag(key, value string) error

	// TagWrapper executes a function with additional profiling tags
	// The tags are only applied during the execution of the function
	// The provided context is passed to the function
	TagWrapper(ctx context.Context, key, value string, fn func(context.Context) error) error
}

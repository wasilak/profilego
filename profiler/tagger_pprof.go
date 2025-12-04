package profiler

import (
	"context"
)

// PprofTagger implements the Tagger interface for pprof backend
// Note: pprof doesn't support runtime tagging like Pyroscope, so this is a no-op implementation
type PprofTagger struct{}

// NewPprofTagger creates a new pprof tagger
func NewPprofTagger() *PprofTagger {
	return &PprofTagger{}
}

// AddTag is a no-op for pprof as it doesn't support runtime tagging
func (pt *PprofTagger) AddTag(key, value string) error {
	// pprof doesn't support runtime tagging, so this is a no-op
	// Users can use build flags or comments for static profiling labels
	return nil
}

// TagWrapper executes the function without adding tags (pprof limitation)
// For pprof, we just execute the function and return its error
func (pt *PprofTagger) TagWrapper(ctx context.Context, key, value string, fn func(context.Context) error) error {
	// pprof doesn't support runtime tagging, so we just execute the function with context
	return fn(ctx)
}

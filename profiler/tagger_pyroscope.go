package profiler

import (
	"context"
	"runtime/pprof"

	"github.com/grafana/pyroscope-go"
)

// PyroscopeTagger implements the Tagger interface for Pyroscope backend
type PyroscopeTagger struct{}

// NewPyroscopeTagger creates a new Pyroscope tagger
func NewPyroscopeTagger() *PyroscopeTagger {
	return &PyroscopeTagger{}
}

// AddTag adds a key-value pair tag to the current Pyroscope profiling context
func (pt *PyroscopeTagger) AddTag(key, value string) error {
	labels := pprof.Labels(key, value)
	pyroscope.TagWrapper(context.Background(), labels, func(ctx context.Context) {
		// No-op function, just adding the tag
	})
	return nil
}

// TagWrapper executes a function with additional Pyroscope profiling tags
// The tags are only applied during the execution of the function
func (pt *PyroscopeTagger) TagWrapper(ctx context.Context, key, value string, fn func(context.Context) error) error {
	labels := pprof.Labels(key, value)
	var err error
	pyroscope.TagWrapper(ctx, labels, func(c context.Context) {
		err = fn(c)
	})
	return err
}

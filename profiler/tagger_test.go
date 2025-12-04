package profiler

import (
	"context"
	"testing"
)

// TestPyroscopeTaggerAddTag tests adding a tag to Pyroscope profiler
func TestPyroscopeTaggerAddTag(t *testing.T) {
	tagger := NewPyroscopeTagger()
	
	// AddTag should not return an error
	err := tagger.AddTag("test_key", "test_value")
	if err != nil {
		t.Fatalf("AddTag returned error: %v", err)
	}
}

// TestPyroscopeTaggerTagWrapper tests TagWrapper with Pyroscope profiler
func TestPyroscopeTaggerTagWrapper(t *testing.T) {
	tagger := NewPyroscopeTagger()
	ctx := context.Background()
	
	executed := false
	err := tagger.TagWrapper(ctx, "test_key", "test_value", func(c context.Context) error {
		executed = true
		return nil
	})
	
	if err != nil {
		t.Fatalf("TagWrapper returned error: %v", err)
	}
	
	if !executed {
		t.Fatal("function was not executed inside TagWrapper")
	}
}

// TestPprofTaggerAddTag tests adding a tag to pprof profiler (no-op)
func TestPprofTaggerAddTag(t *testing.T) {
	tagger := NewPprofTagger()
	
	// AddTag should not return an error (it's a no-op for pprof)
	err := tagger.AddTag("test_key", "test_value")
	if err != nil {
		t.Fatalf("AddTag returned error: %v", err)
	}
}

// TestPprofTaggerTagWrapper tests TagWrapper with pprof profiler
func TestPprofTaggerTagWrapper(t *testing.T) {
	tagger := NewPprofTagger()
	ctx := context.Background()
	
	executed := false
	err := tagger.TagWrapper(ctx, "test_key", "test_value", func(c context.Context) error {
		executed = true
		return nil
	})
	
	if err != nil {
		t.Fatalf("TagWrapper returned error: %v", err)
	}
	
	if !executed {
		t.Fatal("function was not executed inside TagWrapper")
	}
}

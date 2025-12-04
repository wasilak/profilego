package profilego

import (
	"context"
	"testing"

	"github.com/wasilak/profilego/config"
	"github.com/wasilak/profilego/core"
)

// TestAddTag tests the global AddTag function
func TestAddTag(t *testing.T) {
	// Initialize with pprof backend for testing
	cfg := config.Config{
		ApplicationName: "test-app",
		Backend:         core.PprofBackend,
	}

	_, err := InitWithConfig(context.Background(), cfg)
	if err != nil {
		t.Fatalf("failed to initialize profiler: %v", err)
	}
	defer Stop()

	// AddTag should work without error
	err = AddTag("test_key", "test_value")
	if err != nil {
		t.Fatalf("AddTag returned error: %v", err)
	}
}

// TestTagWrapper tests the global TagWrapper function
func TestTagWrapper(t *testing.T) {
	// Initialize with pprof backend for testing
	cfg := config.Config{
		ApplicationName: "test-app",
		Backend:         core.PprofBackend,
	}

	_, err := InitWithConfig(context.Background(), cfg)
	if err != nil {
		t.Fatalf("failed to initialize profiler: %v", err)
	}
	defer Stop()

	// TagWrapper should execute the function and return no error
	ctx := context.Background()
	executed := false
	var receivedCtx context.Context
	err = TagWrapper(ctx, "test_key", "test_value", func(c context.Context) error {
		executed = true
		receivedCtx = c
		return nil
	})

	if err != nil {
		t.Fatalf("TagWrapper returned error: %v", err)
	}

	if !executed {
		t.Fatal("function was not executed inside TagWrapper")
	}

	if receivedCtx == nil {
		t.Fatal("context was not passed to function")
	}

	if receivedCtx != ctx {
		t.Fatal("wrong context passed to function")
	}
}

// TestAddTagWithoutInit tests AddTag without initialization
func TestAddTagWithoutInit(t *testing.T) {
	// Ensure profiler is not initialized
	profilerManager = nil
	globalTagger = nil

	// AddTag should return an error
	err := AddTag("test_key", "test_value")
	if err == nil {
		t.Fatal("AddTag should return error when profiler not initialized")
	}
}

// TestTagWrapperWithoutInit tests TagWrapper without initialization
func TestTagWrapperWithoutInit(t *testing.T) {
	// Ensure profiler is not initialized
	profilerManager = nil
	globalTagger = nil

	// TagWrapper should return an error
	ctx := context.Background()
	err := TagWrapper(ctx, "test_key", "test_value", func(c context.Context) error {
		return nil
	})

	if err == nil {
		t.Fatal("TagWrapper should return error when profiler not initialized")
	}
}

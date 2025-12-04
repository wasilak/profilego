package main

import (
	"context"
	"log"
	"time"

	"github.com/wasilak/profilego"
	"github.com/wasilak/profilego/config"
	"github.com/wasilak/profilego/core"
)

func main() {
	// Create a context for the profiler
	ctx := context.Background()

	// Initialize profiling with Pyroscope backend
	// The context is used throughout the profiler lifecycle
	cfg := config.Config{
		ApplicationName: "tagging-example",
		ServerAddress:   "localhost:4040",
		Backend:         core.PyroscopeBackend,
		Tags: map[string]string{
			"version": "1.0",
			"env":     "development",
		},
	}

	profilerCtx, err := profilego.InitWithConfig(ctx, cfg)
	if err != nil {
		log.Fatalf("failed to initialize profiler: %v", err)
	}
	defer profilego.Stop()

	// Add tags to the profiling context using the new API
	if err := profilego.AddTag("request_id", "12345"); err != nil {
		log.Fatalf("failed to add tag: %v", err)
	}

	// Use TagWrapper to execute code with additional tags
	// Can pass nil to use the profiler context automatically
	if err := profilego.TagWrapper(nil, "route", "api_converter", func(ctx context.Context) error {
		// Your handler logic here
		// ctx is the profiler context passed automatically
		time.Sleep(100 * time.Millisecond)
		return nil
	}); err != nil {
		log.Fatalf("error in tag wrapper: %v", err)
	}

	log.Println("Example completed successfully")
	log.Printf("Profiler context: %v", profilerCtx)
}

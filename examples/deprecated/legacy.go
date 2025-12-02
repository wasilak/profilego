package main

import (
	"log"
	"time"

	"github.com/wasilak/profilego"
)

func main() {
	// DEPRECATED EXAMPLE: This demonstrates the legacy API usage
	// Use this only for maintaining existing code that uses the deprecated API
	// For new projects, use the examples in basic/current.go and advanced/advanced.go

	// Initialize profiling with deprecated configuration API
	// Note: profilego.Config and profilego.Init are deprecated
	// Use config.Config and profilego.InitWithConfig instead
	config := profilego.Config{
		ApplicationName: "deprecated-legacy-example",
		ServerAddress:   "localhost:4040",
		Type:            "pyroscope",
		Tags:            map[string]string{"version": "1.0.0"},
	}

	err := profilego.Init(config)
	if err != nil {
		log.Fatalf("Failed to initialize profiling: %v", err)
	}

	// Simulate some work
	for i := 0; i < 10; i++ {
		time.Sleep(1 * time.Second)
		log.Printf("Working... %d", i)
	}

	// Stop profiling before exiting
	err = profilego.Stop()
	if err != nil {
		log.Printf("Failed to stop profiling: %v", err)
	}
}

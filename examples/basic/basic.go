package main

import (
	"log"
	"time"

	"github.com/wasilak/profilego"
	"github.com/wasilak/profilego/config"
	"github.com/wasilak/profilego/core"
)

func main() {
	// BASIC EXAMPLE: Simple usage of the current API
	// This demonstrates the minimal recommended way to use the library

	// Initialize profiling with basic configuration
	basicConfig := config.Config{
		ApplicationName: "basic-example",
		ServerAddress:   "localhost:4040",
		Backend:         core.PyroscopeBackend,
		Tags:            map[string]string{"version": "1.0.0"},
		InitialState:    core.ProfilingEnabled,
	}

	// Using the current API
	err := profilego.InitWithConfig(basicConfig)
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

package main

import (
	"log"
	"time"

	"github.com/wasilak/profilego"
	"github.com/wasilak/profilego/config"
	"github.com/wasilak/profilego/core"
)

func main() {
	// Initialize profiling with new configuration API
	newConfig := config.Config{
		ApplicationName: "multiple-example",
		ServerAddress:   "localhost:4040",
		Backend:         core.PyroscopeBackend,
		Tags:            map[string]string{"version": "1.0.0", "env": "dev"},
		InitialState:    core.ProfilingEnabled,
	}

	// Using the new API
	err := profilego.InitWithConfig(newConfig)
	if err != nil {
		log.Fatalf("Failed to initialize profiling: %v", err)
	}

	// Simulate some work
	for i := 0; i < 10; i++ {
		time.Sleep(1 * time.Second)
		log.Printf("Working... %d", i)

		// Check if profiling is running
		if profilego.IsRunning() {
			log.Println("Profiling is active")
		}
	}

	// Stop profiling before exiting
	err = profilego.Stop()
	if err != nil {
		log.Printf("Failed to stop profiling: %v", err)
	}
}
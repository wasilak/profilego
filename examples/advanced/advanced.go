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
	// ADVANCED EXAMPLE: Comprehensive usage of the current API
	// This demonstrates advanced configuration and usage patterns

	// Create advanced configuration with multiple profile types
	advancedConfig := config.Config{
		ApplicationName: "advanced-example",
		ServerAddress:   "localhost:4040",
		Backend:         core.PyroscopeBackend,
		Tags: map[string]string{
			"version":   "1.0.0",
			"env":       "production",
			"component": "backend-service",
			"team":      "performance",
		},
		ProfileTypes: []core.ProfileType{
			core.ProfileCPU,
			core.ProfileAllocObjects,
			core.ProfileAllocSpace,
			core.ProfileInuseObjects,
			core.ProfileInuseSpace,
			core.ProfileGoroutines,
			core.ProfileMutexCount,
			core.ProfileMutexDuration,
			core.ProfileBlockCount,
			core.ProfileBlockDuration,
		},
		InitialState: core.ProfilingEnabled,
		AdditionalAttrs: []any{
			"custom_attribute_1",
			"custom_attribute_2",
		},
	}

	// Initialize with advanced configuration and retrieve the context
	ctx, err := profilego.InitWithConfig(context.Background(), advancedConfig)
	if err != nil {
		log.Fatalf("Failed to initialize profiling: %v", err)
	}
	_ = ctx // Use context throughout your application

	// Demonstrate advanced usage patterns
	for i := 0; i < 15; i++ {
		time.Sleep(500 * time.Millisecond)

		// Check profiling status
		if profilego.IsRunning() {
			log.Printf("Profiling active - iteration %d", i)

			// Demonstrate state management
			if i == 5 {
				// Pause profiling temporarily
				err := profilego.Stop()
				if err != nil {
					log.Printf("Failed to stop profiling: %v", err)
				} else {
					log.Println("Profiling paused for demonstration")
					time.Sleep(1 * time.Second)

					// Resume profiling
					err := profilego.Start()
					if err != nil {
						log.Printf("Failed to restart profiling: %v", err)
					} else {
						log.Println("Profiling resumed")
					}
				}
			}
		} else {
			log.Printf("Profiling inactive - iteration %d", i)
		}
	}

	// Clean shutdown
	err = profilego.Stop()
	if err != nil {
		log.Printf("Failed to stop profiling: %v", err)
	} else {
		log.Println("Profiling stopped successfully")
	}
}

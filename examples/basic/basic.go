package main

import (
	"log"
	"time"

	"github.com/wasilak/profilego"
)

func main() {
	// Initialize profiling with basic configuration (legacy API)
	config := profilego.Config{
		ApplicationName: "basic-example",
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
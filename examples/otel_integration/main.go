package main

import (
	"context"
	"log"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/sdk/trace"

	"github.com/wasilak/profilego"
	"github.com/wasilak/profilego/config"
	"github.com/wasilak/profilego/core"
)

func main() {
	ctx := context.Background()

	// Example 1: Using WrapTracerProvider with explicit control
	// This approach gives users full control over when wrapping happens
	explicitWrappingExample(ctx)
}

// explicitWrappingExample demonstrates explicit TracerProvider wrapping
func explicitWrappingExample(ctx context.Context) {
	log.Println("=== Explicit OTel Integration Example ===")

	// Initialize profilego with Pyroscope backend
	cfg := config.Config{
		ApplicationName: "otel-integration-example",
		Backend:         core.PyroscopeBackend,
		ServerAddress:   "localhost:4040",
		Tags: map[string]string{
			"environment": "development",
			"version":     "1.0",
		},
	}

	profilerCtx, err := profilego.InitWithConfig(ctx, cfg)
	if err != nil {
		log.Fatalf("Failed to initialize profiler: %v", err)
	}
	defer profilego.Stop()

	// Create a standard OpenTelemetry TracerProvider
	// In a real application, you would configure this with processors, exporters, etc.
	standardTP := trace.NewTracerProvider(
		// Add your processors and exporters here
		// trace.WithBatcher(otlptracehttp.NewClient(...)),
	)

	// Let profilego wrap it with Pyroscope integration
	// This wraps the TracerProvider with Pyroscope's OTel integration
	wrappedTP, err := profilego.WrapTracerProvider(ctx, standardTP)
	if err != nil {
		log.Fatalf("Failed to wrap TracerProvider: %v", err)
	}

	// Register the wrapped TracerProvider globally
	otel.SetTracerProvider(wrappedTP)

	log.Println("OTel TracerProvider wrapped and registered successfully")
	log.Println("Profiler context available from InitWithConfig")
	_ = profilerCtx // Use this context in your handlers/middleware
}

// integratedSetupExample demonstrates OTel integration via Config
// This approach provides automatic wrapping during initialization
func integratedSetupExample(ctx context.Context) {
	log.Println("=== Integrated OTel Setup Example ===")

	// Create a standard TracerProvider
	standardTP := trace.NewTracerProvider()

	// Initialize profilego with OTel integration built-in
	// profilego will automatically wrap and register the TracerProvider
	cfg := config.Config{
		ApplicationName:    "otel-integrated-example",
		Backend:            core.PyroscopeBackend,
		ServerAddress:      "localhost:4040",
		OTelTracerProvider: standardTP, // Automatic wrapping
	}

	profilerCtx, err := profilego.InitWithConfig(ctx, cfg)
	if err != nil {
		log.Fatalf("Failed to initialize profiler: %v", err)
	}
	defer profilego.Stop()

	log.Println("Profiler and OTel TracerProvider initialized and wrapped automatically")
	_ = profilerCtx
}

// convenienceSetterExample demonstrates the SetTracerProvider helper
func convenienceSetterExample(ctx context.Context) {
	log.Println("=== Convenience SetTracerProvider Example ===")

	// Initialize profilego first
	cfg := config.Config{
		ApplicationName: "otel-convenience-example",
		Backend:         core.PyroscopeBackend,
		ServerAddress:   "localhost:4040",
	}

	profilerCtx, err := profilego.InitWithConfig(ctx, cfg)
	if err != nil {
		log.Fatalf("Failed to initialize profiler: %v", err)
	}
	defer profilego.Stop()

	// Create and wrap+register in one call
	standardTP := trace.NewTracerProvider()
	if err := profilego.SetTracerProvider(ctx, standardTP); err != nil {
		log.Fatalf("Failed to set TracerProvider: %v", err)
	}

	log.Println("TracerProvider wrapped and registered using convenience function")
	_ = profilerCtx
}

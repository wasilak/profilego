package profilego

import (
	"context"
	"testing"

	"go.opentelemetry.io/otel/sdk/trace"

	"github.com/wasilak/profilego/config"
	"github.com/wasilak/profilego/core"
)

// TestWrapTracerProviderPyroscope tests wrapping with Pyroscope backend
func TestWrapTracerProviderPyroscope(t *testing.T) {
	ctx := context.Background()

	// Initialize with Pyroscope backend
	cfg := config.Config{
		ApplicationName: "test-app",
		Backend:         core.PyroscopeBackend,
		ServerAddress:   "localhost:4040",
	}

	_, err := InitWithConfig(ctx, cfg)
	if err != nil {
		t.Fatalf("Failed to initialize: %v", err)
	}
	defer Stop()

	// Create a standard TracerProvider
	standardTP := trace.NewTracerProvider()

	// Wrap it
	wrapped, err := WrapTracerProvider(ctx, standardTP)
	if err != nil {
		t.Fatalf("WrapTracerProvider failed: %v", err)
	}

	// Should return a wrapped provider (not the same as input)
	if wrapped == nil {
		t.Fatal("wrapped TracerProvider is nil")
	}
}

// TestWrapTracerProviderPprof tests wrapping with pprof backend
func TestWrapTracerProviderPprof(t *testing.T) {
	ctx := context.Background()

	// Initialize with pprof backend
	cfg := config.Config{
		ApplicationName: "test-app",
		Backend:         core.PprofBackend,
	}

	_, err := InitWithConfig(ctx, cfg)
	if err != nil {
		t.Fatalf("Failed to initialize: %v", err)
	}
	defer Stop()

	// Create a standard TracerProvider
	standardTP := trace.NewTracerProvider()

	// Wrap it
	wrapped, err := WrapTracerProvider(ctx, standardTP)
	if err != nil {
		t.Fatalf("WrapTracerProvider failed: %v", err)
	}

	// For pprof backend, should return the same provider unchanged
	if wrapped != standardTP {
		t.Fatal("for pprof backend, wrapped should be the same as input")
	}
}

// TestWrapTracerProviderNotInitialized tests wrapping without initialization
func TestWrapTracerProviderNotInitialized(t *testing.T) {
	ctx := context.Background()

	// Ensure profiler is not initialized
	profilerManager = nil

	// Create a standard TracerProvider
	standardTP := trace.NewTracerProvider()

	// Try to wrap - should fail
	_, err := WrapTracerProvider(ctx, standardTP)
	if err == nil {
		t.Fatal("WrapTracerProvider should fail when profiler not initialized")
	}
}

// TestSetTracerProvider tests the convenience wrapper function
func TestSetTracerProvider(t *testing.T) {
	ctx := context.Background()

	// Initialize with Pyroscope
	cfg := config.Config{
		ApplicationName: "test-app",
		Backend:         core.PyroscopeBackend,
		ServerAddress:   "localhost:4040",
	}

	_, err := InitWithConfig(ctx, cfg)
	if err != nil {
		t.Fatalf("Failed to initialize: %v", err)
	}
	defer Stop()

	// Create a TracerProvider
	tp := trace.NewTracerProvider()

	// Set it using the convenience function
	if err := SetTracerProvider(ctx, tp); err != nil {
		t.Fatalf("SetTracerProvider failed: %v", err)
	}
}

// TestOTelTracerProviderInConfig tests OTel integration via Config
func TestOTelTracerProviderInConfig(t *testing.T) {
	ctx := context.Background()

	// Create a TracerProvider
	tp := trace.NewTracerProvider()

	// Initialize with OTel integration
	cfg := config.Config{
		ApplicationName:    "test-app",
		Backend:            core.PyroscopeBackend,
		ServerAddress:      "localhost:4040",
		OTelTracerProvider: tp,
	}

	_, err := InitWithConfig(ctx, cfg)
	if err != nil {
		t.Fatalf("Failed to initialize with OTel: %v", err)
	}
	defer Stop()

	// Success if no error
}

// TestOTelTracerProviderInvalidType tests Config with invalid OTelTracerProvider type
func TestOTelTracerProviderInvalidType(t *testing.T) {
	ctx := context.Background()

	// Initialize basic profiler first
	cfg := config.Config{
		ApplicationName: "test-app",
		Backend:         core.PyroscopeBackend,
		ServerAddress:   "localhost:4040",
	}

	_, err := InitWithConfig(ctx, cfg)
	if err != nil {
		t.Fatalf("Failed to initialize: %v", err)
	}
	defer Stop()

	// Try to handle invalid OTel type
	if err := handleOTelTracerProvider(ctx, "not-a-tracer-provider"); err == nil {
		t.Fatal("handleOTelTracerProvider should fail with invalid type")
	}
}

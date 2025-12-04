package profilego

import (
	"context"
	"fmt"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
	otelpyroscope "github.com/grafana/otel-profiling-go"

	"github.com/wasilak/profilego/core"
)

// WrapTracerProvider wraps a standard TracerProvider with Pyroscope profiling
// when using the Pyroscope backend. For other backends, returns the original provider.
//
// This function is optional and allows users to integrate with OpenTelemetry
// without directly importing github.com/grafana/otel-profiling-go
//
// If profiler is not initialized, returns an error.
// If backend is not Pyroscope, returns the original provider unchanged.
func WrapTracerProvider(ctx context.Context, tp trace.TracerProvider) (trace.TracerProvider, error) {
	if profilerManager == nil {
		return nil, fmt.Errorf("profiler not initialized - call InitWithConfig first")
	}

	// Only wrap for Pyroscope backend
	if profilerManager.Backend() != core.PyroscopeBackend {
		return tp, nil // Other backends don't need wrapping
	}

	// Wrap with Pyroscope OTel integration
	wrappedTP := otelpyroscope.NewTracerProvider(
		tp,
		otelpyroscope.WithAppName(profilerManager.AppName()),
		otelpyroscope.WithPyroscopeURL(profilerManager.ServerAddress()),
	)

	return wrappedTP, nil
}

// SetTracerProvider wraps and registers a TracerProvider globally.
// This is a convenience function combining WrapTracerProvider and otel.SetTracerProvider.
//
// For Pyroscope backend, wraps the provided TracerProvider with Pyroscope integration.
// For other backends, registers the provider unchanged.
func SetTracerProvider(ctx context.Context, tp trace.TracerProvider) error {
	wrappedTP, err := WrapTracerProvider(ctx, tp)
	if err != nil {
		return err
	}
	otel.SetTracerProvider(wrappedTP)
	return nil
}

// handleOTelTracerProvider handles OTel integration if provided in config.
// The tp parameter should be a trace.TracerProvider instance.
func handleOTelTracerProvider(ctx context.Context, tp interface{}) error {
	// Type assert to trace.TracerProvider
	tracerProvider, ok := tp.(trace.TracerProvider)
	if !ok {
		return fmt.Errorf("OTelTracerProvider is not a valid trace.TracerProvider")
	}

	wrappedTP, err := WrapTracerProvider(ctx, tracerProvider)
	if err != nil {
		return err
	}

	return SetTracerProvider(ctx, wrappedTP)
}

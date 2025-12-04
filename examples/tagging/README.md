# Tagging Example

This example demonstrates how to use the new tagging API with context flow in profilego.

## Overview

The tagging API allows you to add profiling context to your applications without needing to import backend-specific libraries. Context is properly integrated throughout the profiler lifecycle.

## Running the Example

```bash
# From the profilego root directory
go run ./examples/tagging
```

## What It Does

1. Creates a context for the profiler
2. Initializes profilego with the context
3. Adds simple tags using `AddTag()`
4. Executes code with contextual tags using `TagWrapper()`
5. Retrieves and uses the profiler context

## Code Structure

```go
// 1. Create context for the profiler
ctx := context.Background()

// 2. Initialize profilego with context
profilego.InitWithConfig(ctx, cfg)

// 3. Get profiler context for use in handlers
profilerCtx := profilego.Context()

// 4. Add a simple tag
profilego.AddTag("request_id", "12345")

// 5. Execute code with tags (pass nil to use profiler context)
profilego.TagWrapper(nil, "route", "api_converter", func() error {
    // This code runs with the "route" tag set to "api_converter"
    return nil
})
```

## Key Features

- **No backend imports**: Your handler code doesn't import `github.com/grafana/pyroscope-go`
- **Context aware**: Context flows through entire profiler lifecycle
- **Backend agnostic**: Switch from Pyroscope to pprof by changing config only
- **Consistent API**: Same tagging API works with all backends
- **Graceful shutdown**: Control profiler context lifecycle

## Key Benefits

✅ Proper Go idioms with context propagation  
✅ Supports context cancellation for graceful shutdown  
✅ Supports context timeouts  
✅ Clean API - pass nil to use profiler context automatically  

## Next Steps

See `examples/echo_integration` for a complete web framework example with middleware integration and context-aware tagging throughout the request lifecycle.

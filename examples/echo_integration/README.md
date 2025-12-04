# Echo Framework Integration Example

This example demonstrates how to integrate profilego's tagging API with the Echo web framework.

## Overview

This shows a production-ready pattern for adding profiling context to HTTP handlers without tight coupling to profiling backends.

## Installation

First, add Echo as a dependency:

```bash
go get github.com/labstack/echo/v4
```

## Running the Example

```bash
# From the profilego root directory
go run ./examples/echo_integration
```

This will start an HTTP server on `localhost:8080`.

## Example Requests

```bash
# Get user data
curl http://localhost:8080/users/123

# Convert with type parameter
curl http://localhost:8080/convert?type=json

# Health check
curl http://localhost:8080/health
```

## Code Pattern

### 1. Middleware-Based Tagging

```go
func ProfilingMiddleware() echo.MiddlewareFunc {
    return func(next echo.HandlerFunc) echo.HandlerFunc {
        return func(c echo.Context) error {
            route := c.Path()
            method := c.Request().Method

            if err := profilego.AddTag("http_method", method); err != nil {
                log.Printf("failed to add tag: %v", err)
            }

            // Pass nil to use profiler context automatically
            return profilego.TagWrapper(
                nil,
                "route",
                route,
                func() error {
                    return next(c)
                },
            )
        }
    }
}
```

### 2. Handler-Level Tagging

```go
func UserHandler(c echo.Context) error {
    userID := c.Param("user_id")

    // Add contextual tag with error handling
    if err := profilego.AddTag("user_id", userID); err != nil {
        log.Printf("failed to add tag: %v", err)
    }

    return c.JSON(http.StatusOK, map[string]string{
        "user_id": userID,
    })
}
```

### 3. Nested TagWrapper with Profiler Context

```go
func ConvertHandler(c echo.Context) error {
    conversionType := c.QueryParam("type")

    // Pass nil to use profiler context automatically
    return profilego.TagWrapper(
        nil,
        "conversion_type",
        conversionType,
        func() error {
            if err := profilego.AddTag("status", "processing"); err != nil {
                log.Printf("failed to add tag: %v", err)
            }

            // Do work here
            return c.JSON(http.StatusOK, map[string]string{
                "type": conversionType,
            })
        },
    )
}
```

## Tags Generated

For a request to `GET /users/123`:
- `service`: "api-server" (from init config)
- `env`: "development" (from init config)
- `http_method`: "GET" (from middleware)
- `route`: "/users/:user_id" (from middleware)
- `user_id`: "123" (from handler)

## Profiling Backend

To view profiles, you'll need a Pyroscope server running:

```bash
# Using Docker
docker run -p 4040:4040 grafana/pyroscope:latest
```

Then visit `http://localhost:4040` to view profiles.

## Key Insights

1. **No Pyroscope imports in handlers** - Completely decoupled
2. **Middleware encapsulation** - HTTP concerns (route, method) in middleware
3. **Handler-specific tags** - Each handler can add its own context
4. **Composable approach** - Tags stack naturally through the call stack
5. **Framework agnostic** - Same pattern works with Gin, Fiber, net/http, etc.

## Benefits Over Direct API Usage

### Before (Direct Pyroscope)
```go
import "github.com/grafana/pyroscope-go"
import "runtime/pprof"

func handler(c echo.Context) error {
    pyroscope.TagWrapper(c.Request().Context(),
        pprof.Labels("route", c.Path()),
        func(ctx context.Context) {
            // handler logic
        })
}
```

### After (Profilego with Context)
```go
func handler(c echo.Context) error {
    // Pass nil - uses profiler context automatically
    return profilego.TagWrapper(nil,
        "route", c.Path(),
        func() error {
            // handler logic
            return nil
        })
}
```

The latter is:
- ✓ Cleaner API (no `pprof.Labels` boilerplate)
- ✓ Backend agnostic (works with any profiling backend)
- ✓ Easier to test (no Pyroscope imports needed)
- ✓ Smaller cognitive load (one API instead of three)
- ✓ Context-aware (profiler context flows throughout app)
- ✓ Simpler handlers (pass nil instead of c.Request().Context())
- ✓ Proper Go idioms (context propagation built-in)

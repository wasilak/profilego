package main

import (
	"context"
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/wasilak/profilego"
	"github.com/wasilak/profilego/config"
	"github.com/wasilak/profilego/core"
)

// ProfilingMiddleware adds profiling tags to each request
func ProfilingMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Extract request info for tags
			route := c.Path()
			method := c.Request().Method

			// Add tags to the current request context
			if err := profilego.AddTag("http_method", method); err != nil {
				// Log but don't fail
				log.Printf("failed to add tag: %v", err)
			}

			// Use TagWrapper to wrap the entire handler execution
			// Pass nil to use the profiler context automatically
			return profilego.TagWrapper(
				nil,
				"route",
				route,
				func(ctx context.Context) error {
					// Execute the actual handler
					// ctx is the profiler context passed automatically
					return next(c)
				},
			)
		}
	}
}

// UserHandler demonstrates tagging within a handler
func UserHandler(c echo.Context) error {
	userID := c.Param("user_id")

	// Additional tag for this specific handler
	if err := profilego.AddTag("user_id", userID); err != nil {
		log.Printf("failed to add tag: %v", err)
	}

	return c.JSON(http.StatusOK, map[string]string{
		"user_id": userID,
		"message": "User data retrieved",
	})
}

// ConvertHandler demonstrates nested TagWrapper usage
func ConvertHandler(c echo.Context) error {
	conversionType := c.QueryParam("type")

	// Wrap the conversion logic with tags
	// Pass nil to use the profiler context automatically
	return profilego.TagWrapper(
		nil,
		"conversion_type",
		conversionType,
		func(ctx context.Context) error {
			// Simulate some work
			if err := profilego.AddTag("status", "processing"); err != nil {
				log.Printf("failed to add tag: %v", err)
			}

			// In a real app, this would do actual conversion work
			return c.JSON(http.StatusOK, map[string]string{
				"type":   conversionType,
				"status": "completed",
			})
		},
	)
}

func main() {
	// Create a context for the profiler
	ctx := context.Background()

	// Initialize profilego with Pyroscope backend
	cfg := config.Config{
		ApplicationName: "echo-integration-example",
		ServerAddress:   "localhost:4040", // Pyroscope server
		Backend:         core.PyroscopeBackend,
		Tags: map[string]string{
			"service": "api-server",
			"env":     "development",
		},
	}

	profilerCtx, err := profilego.InitWithConfig(ctx, cfg)
	if err != nil {
		log.Fatalf("Failed to initialize profiler: %v", err)
	}
	defer profilego.Stop()

	log.Println("Profiler initialized with context")
	_ = profilerCtx // Available for use in handlers

	// Create Echo server
	e := echo.New()

	// Add logging middleware
	e.Use(middleware.Logger())

	// Add recovery middleware
	e.Use(middleware.Recover())

	// Add profiling middleware - this wraps all routes
	// All requests will be tagged with route and http_method
	e.Use(ProfilingMiddleware())

	// Routes
	e.GET("/users/:user_id", UserHandler)
	e.GET("/convert", ConvertHandler)

	// Health check endpoint
	e.GET("/health", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{
			"status": "healthy",
		})
	})

	log.Println("Starting server on :8080")
	log.Println("Profiler context is active throughout request lifecycle")
	log.Println("Pyroscope profiling enabled - check http://localhost:4040")
	if err := e.Start(":8080"); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}

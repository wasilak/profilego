package manager

import (
	"context"
	"testing"
	"time"

	"github.com/wasilak/profilego/config"
	"github.com/wasilak/profilego/core"
)

// TestProfiler is a test implementation of the Profiler interface
type TestProfiler struct {
	name    string
	started bool
	paused  bool
}

func (t *TestProfiler) Start(ctx context.Context) error {
	t.started = true
	return nil
}

func (t *TestProfiler) Stop(ctx context.Context) error {
	t.started = false
	return nil
}

func (t *TestProfiler) Pause(ctx context.Context) error {
	t.paused = true
	return nil
}

func (t *TestProfiler) Resume(ctx context.Context) error {
	t.paused = false
	return nil
}

func (t *TestProfiler) IsRunning() bool {
	return t.started && !t.paused
}

func (t *TestProfiler) Name() string {
	return t.name
}

func TestNewProfilerManager(t *testing.T) {
	cfg := config.Config{
		ApplicationName: "test-app",
		Backend:         config.PyroscopeBackend,
	}
	
	manager := NewProfilerManager(cfg)
	
	if manager.config.ApplicationName != "test-app" {
		t.Errorf("Expected ApplicationName 'test-app', got '%s'", manager.config.ApplicationName)
	}
	
	if len(manager.profilers) != 0 {
		t.Errorf("Expected empty profilers map, got %d profilers", len(manager.profilers))
	}
	
	if manager.running {
		t.Error("Manager should not be running initially")
	}
}

func TestAddProfiler(t *testing.T) {
	cfg := config.Config{
		ApplicationName: "test-app",
		Backend:         config.PyroscopeBackend,
	}
	
	manager := NewProfilerManager(cfg)
	
	profiler := &TestProfiler{name: "test1"}
	err := manager.AddProfiler(profiler)
	if err != nil {
		t.Errorf("AddProfiler returned error: %v", err)
	}
	
	if len(manager.profilers) != 1 {
		t.Errorf("Expected 1 profiler, got %d", len(manager.profilers))
	}
	
	// Try adding the same profiler again
	err = manager.AddProfiler(profiler)
	if err == nil {
		t.Error("Adding duplicate profiler should return error")
	}
}

func TestIsRunning(t *testing.T) {
	cfg := config.Config{
		ApplicationName: "test-app",
		Backend:         config.PyroscopeBackend,
	}
	
	manager := NewProfilerManager(cfg)
	
	// Should not be running initially
	if manager.IsRunning() {
		t.Error("Manager should not be running initially")
	}
	
	// Set running to true manually for test
	manager.mu.Lock()
	manager.running = true
	manager.mu.Unlock()
	
	// Should be running now
	if !manager.IsRunning() {
		t.Error("Manager should be running after setting to true")
	}
}

func TestManagerError(t *testing.T) {
	err := &ManagerError{
		Operation: "test",
		Message:   "test message",
	}
	
	expected := "manager error (test): test message"
	if err.Error() != expected {
		t.Errorf("Expected error '%s', got '%s'", expected, err.Error())
	}
}

func TestConcurrentAccess(t *testing.T) {
	cfg := config.Config{
		ApplicationName: "test-app",
		Backend:         config.PyroscopeBackend,
	}
	
	manager := NewProfilerManager(cfg)
	
	// Add a profiler
	profiler := &TestProfiler{name: "concurrent-test"}
	manager.AddProfiler(profiler)
	
	// Run multiple goroutines to test concurrent access
	errChan := make(chan error, 10)
	
	for i := 0; i < 10; i++ {
		go func() {
			// Test reading IsRunning safely
			running := manager.IsRunning()
			_ = running // Just to check if it panics
			
			// Test that no panic occurs when accessing the manager
			manager.mu.RLock()
			_ = len(manager.profilers)
			manager.mu.RUnlock()
			
			errChan <- nil
		}()
	}
	
	// Wait for all goroutines to finish
	for i := 0; i < 10; i++ {
		err := <-errChan
		if err != nil {
			t.Errorf("Concurrent access caused error: %v", err)
		}
	}
}
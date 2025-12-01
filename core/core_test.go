package core

import (
	"context"
	"testing"
)

// MockProfiler is a mock implementation of the Profiler interface for testing
type MockProfiler struct {
	started bool
	paused  bool
	name    string
}

func (m *MockProfiler) Start(ctx context.Context) error {
	m.started = true
	return nil
}

func (m *MockProfiler) Stop(ctx context.Context) error {
	m.started = false
	return nil
}

func (m *MockProfiler) Pause(ctx context.Context) error {
	m.paused = true
	return nil
}

func (m *MockProfiler) Resume(ctx context.Context) error {
	m.paused = false
	return nil
}

func (m *MockProfiler) IsRunning() bool {
	return m.started && !m.paused
}

func (m *MockProfiler) Name() string {
	return m.name
}

// Test Profiler interface implementation
func TestProfilerInterface(t *testing.T) {
	profiler := &MockProfiler{name: "test"}

	ctx := context.Background()

	// Test Start
	err := profiler.Start(ctx)
	if err != nil {
		t.Errorf("Start returned error: %v", err)
	}
	if !profiler.started {
		t.Error("Profiler should be started after Start() call")
	}

	// Test IsRunning
	if !profiler.IsRunning() {
		t.Error("Profiler should be running after Start()")
	}

	// Test Pause
	err = profiler.Pause(ctx)
	if err != nil {
		t.Errorf("Pause returned error: %v", err)
	}
	if !profiler.paused {
		t.Error("Profiler should be paused after Pause() call")
	}
	if profiler.IsRunning() {
		t.Error("Profiler should not be running after Pause()")
	}

	// Test Resume
	err = profiler.Resume(ctx)
	if err != nil {
		t.Errorf("Resume returned error: %v", err)
	}
	if profiler.paused {
		t.Error("Profiler should not be paused after Resume() call")
	}
	if !profiler.IsRunning() {
		t.Error("Profiler should be running after Resume()")
	}

	// Test Stop
	err = profiler.Stop(ctx)
	if err != nil {
		t.Errorf("Stop returned error: %v", err)
	}
	if profiler.started {
		t.Error("Profiler should not be started after Stop()")
	}
	if profiler.IsRunning() {
		t.Error("Profiler should not be running after Stop()")
	}

	// Test Name
	if profiler.Name() != "test" {
		t.Errorf("Expected name 'test', got '%s'", profiler.Name())
	}
}
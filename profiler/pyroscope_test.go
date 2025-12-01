package profiler

import (
	"context"
	"testing"

	"github.com/wasilak/profilego/config"
)

func TestNewPyroscopeProfiler(t *testing.T) {
	cfg := config.Config{
		ApplicationName: "test-app",
		ServerAddress:   "localhost:4040",
		Backend:         config.PyroscopeBackend,
	}
	
	profiler, err := NewPyroscopeProfiler(cfg)
	if err != nil {
		t.Errorf("NewPyroscopeProfiler returned error: %v", err)
	}
	
	if profiler.config.ApplicationName != "test-app" {
		t.Errorf("Expected ApplicationName 'test-app', got '%s'", profiler.config.ApplicationName)
	}
	
	if profiler.Name() != "pyroscope" {
		t.Errorf("Expected Name 'pyroscope', got '%s'", profiler.Name())
	}
}

func TestPyroscopeProfilerName(t *testing.T) {
	cfg := config.Config{}
	profiler, _ := NewPyroscopeProfiler(cfg)
	
	if profiler.Name() != "pyroscope" {
		t.Errorf("Expected Name 'pyroscope', got '%s'", profiler.Name())
	}
}

func TestPyroscopeProfilerIsRunning(t *testing.T) {
	cfg := config.Config{}
	profiler, _ := NewPyroscopeProfiler(cfg)
	
	// Initially should not be running
	if profiler.IsRunning() {
		t.Error("PyroscopeProfiler should not be running initially")
	}
}

func TestConvertProfileType(t *testing.T) {
	cfg := config.Config{}
	profiler, _ := NewPyroscopeProfiler(cfg)
	
	// Test each profile type conversion
	testCases := []struct {
		input    config.ProfileType
		expected bool
	}{
		{config.ProfileCPU, true},
		{config.ProfileAllocObjects, true},
		{config.ProfileAllocSpace, true},
		{config.ProfileInuseObjects, true},
		{config.ProfileInuseSpace, true},
		{config.ProfileGoroutines, true},
		{config.ProfileMutexCount, true},
		{config.ProfileMutexDuration, true},
		{config.ProfileBlockCount, true},
		{config.ProfileBlockDuration, true},
		{config.ProfileType("invalid"), false},
	}
	
	for _, tc := range testCases {
		_, ok := profiler.convertProfileType(tc.input)
		if ok != tc.expected {
			t.Errorf("convertProfileType(%s) = %v, want %v", tc.input, ok, tc.expected)
		}
	}
}

// Note: We don't test Start/Stop/Pause/Resume extensively since they interact with external services
// and would require mocking the pyroscope library which is complex.
// The important logic is covered in the conversion and helper methods.
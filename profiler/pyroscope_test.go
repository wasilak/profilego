package profiler

import (
	"testing"

	"github.com/wasilak/profilego/config"
	"github.com/wasilak/profilego/core"
)

func TestNewPyroscopeProfiler(t *testing.T) {
	cfg := config.Config{
		ApplicationName: "test-app",
		ServerAddress:   "localhost:4040",
		Backend:         core.PyroscopeBackend,
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
		input    core.ProfileType
		expected bool
	}{
		{core.ProfileCPU, true},
		{core.ProfileAllocObjects, true},
		{core.ProfileAllocSpace, true},
		{core.ProfileInuseObjects, true},
		{core.ProfileInuseSpace, true},
		{core.ProfileGoroutines, true},
		{core.ProfileMutexCount, true},
		{core.ProfileMutexDuration, true},
		{core.ProfileBlockCount, true},
		{core.ProfileBlockDuration, true},
		{core.ProfileType("invalid"), false},
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

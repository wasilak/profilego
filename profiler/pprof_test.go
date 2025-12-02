package profiler

import (
	"testing"

	"github.com/wasilak/profilego/config"
	"github.com/wasilak/profilego/core"
)

func TestNewPprofProfiler(t *testing.T) {
	cfg := config.Config{
		ApplicationName: "test-app",
		Backend:         core.PprofBackend,
	}

	profiler, err := NewPprofProfiler(cfg)
	if err != nil {
		t.Errorf("NewPprofProfiler returned error: %v", err)
	}

	if profiler.config.ApplicationName != "test-app" {
		t.Errorf("Expected ApplicationName 'test-app', got '%s'", profiler.config.ApplicationName)
	}

	if profiler.Name() != "pprof" {
		t.Errorf("Expected Name 'pprof', got '%s'", profiler.Name())
	}
}

func TestPprofProfilerName(t *testing.T) {
	cfg := config.Config{}
	profiler, _ := NewPprofProfiler(cfg)

	if profiler.Name() != "pprof" {
		t.Errorf("Expected Name 'pprof', got '%s'", profiler.Name())
	}
}

func TestPprofProfilerIsRunning(t *testing.T) {
	cfg := config.Config{}
	profiler, _ := NewPprofProfiler(cfg)

	// Initially should not be running
	if profiler.IsRunning() {
		t.Error("PprofProfiler should not be running initially")
	}
}

// Note: We don't test Start/Stop/Pause/Resume extensively since they interact with runtime profiling
// and would require complex setup to test properly. The structure and basic functionality is tested.

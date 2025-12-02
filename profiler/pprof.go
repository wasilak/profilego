package profiler

import (
	"context"
	"log/slog"
	"os"
	"runtime"
	"runtime/pprof"
	"sync"
	"time"

	"dario.cat/mergo"
	"github.com/wasilak/profilego/config"
	"github.com/wasilak/profilego/core"
	"github.com/wasilak/profilego/internal/memory"
)

// PprofProfiler implements the Profiler interface for pprof
type PprofProfiler struct {
	mu      sync.RWMutex
	config  config.Config
	running bool
	stopCh  chan struct{}
}

// NewPprofProfiler creates a new pprof profiler
func NewPprofProfiler(cfg config.Config) (*PprofProfiler, error) {
	// Start with defaults
	finalConfig := config.DefaultConfig

	// Then merge provided config values (they take precedence)
	err := mergo.Merge(&finalConfig, cfg, mergo.WithOverride)
	if err != nil {
		return nil, err
	}

	pp := &PprofProfiler{
		config: finalConfig,
		stopCh: make(chan struct{}),
	}

	return pp, nil
}

// Name returns the profiler's identifier
func (pp *PprofProfiler) Name() string {
	return "pprof"
}

// Start begins profiling
func (pp *PprofProfiler) Start(ctx context.Context) error {
	pp.mu.Lock()
	defer pp.mu.Unlock()

	if pp.running {
		return nil // Already running
	}

	// Check memory limit before starting
	if pp.config.MemoryLimitMB > 0 {
		if err := memory.CheckMemoryLimit(pp.config.MemoryLimitMB); err != nil {
			return err
		}
	}

	// Start profiling based on configured profile types
	for _, profileType := range pp.config.ProfileTypes {
		switch profileType {
		case core.ProfileCPU:
			f, err := os.Create(pp.config.ApplicationName + "_cpu.pprof")
			if err != nil {
				return err
			}
			if err := pprof.StartCPUProfile(f); err != nil {
				return err
			}
		case core.ProfileGoroutines:
			// Goroutine profiling is ongoing, no need to start
			slog.InfoContext(ctx, "Goroutine profiling is ongoing")
		case core.ProfileMutexCount, core.ProfileMutexDuration:
			runtime.SetMutexProfileFraction(1) // Enable mutex profiling
		case core.ProfileBlockCount, core.ProfileBlockDuration:
			runtime.SetBlockProfileRate(1) // Enable block profiling
		}
	}

	pp.running = true

	// Start a goroutine to periodically profile if needed
	go pp.profileLoop()

	return nil
}

// Stop gracefully stops profiling
func (pp *PprofProfiler) Stop(ctx context.Context) error {
	pp.mu.Lock()
	defer pp.mu.Unlock()

	if !pp.running {
		return nil
	}

	// Stop profiling based on configured profile types
	for _, profileType := range pp.config.ProfileTypes {
		switch profileType {
		case core.ProfileCPU:
			pprof.StopCPUProfile()
		case core.ProfileMutexCount, core.ProfileMutexDuration:
			runtime.SetMutexProfileFraction(0) // Disable mutex profiling
		case core.ProfileBlockCount, core.ProfileBlockDuration:
			runtime.SetBlockProfileRate(0) // Disable block profiling
		}
	}

	// Notify the profile loop to stop
	close(pp.stopCh)

	pp.running = false
	return nil
}

// Pause temporarily stops profiling
func (pp *PprofProfiler) Pause(ctx context.Context) error {
	pp.mu.Lock()
	defer pp.mu.Unlock()

	if !pp.running {
		return nil
	}

	// Stop only CPU profiling during pause
	for _, profileType := range pp.config.ProfileTypes {
		if profileType == core.ProfileCPU {
			pprof.StopCPUProfile()
		}
	}

	pp.running = false
	return nil
}

// Resume after pause
func (pp *PprofProfiler) Resume(ctx context.Context) error {
	return pp.Start(ctx)
}

// IsRunning returns the current state of the profiler
func (pp *PprofProfiler) IsRunning() bool {
	pp.mu.RLock()
	defer pp.mu.RUnlock()
	return pp.running
}

// profileLoop runs profiling operations in a loop
func (pp *PprofProfiler) profileLoop() {
	ticker := time.NewTicker(10 * time.Second) // Profile every 10 seconds
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			// Perform periodic profiling tasks if needed
			// For pprof, most profiling is done through runtime
		case <-pp.stopCh:
			// Stop the profiling loop
			return
		}
	}
}

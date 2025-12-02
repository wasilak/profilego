package memory

import (
	"fmt"
	"runtime"
	"sync"
	"time"
)

// MemoryMonitor provides memory monitoring capabilities
type MemoryMonitor struct {
	mu         sync.RWMutex
	limitMB    int64
	monitoring bool
	stopCh     chan struct{}
	usageFunc  func() uint64 // Function to get current memory usage in bytes
}

// NewMemoryMonitor creates a new memory monitor with the specified limit in MB
func NewMemoryMonitor(limitMB int64) *MemoryMonitor {
	return &MemoryMonitor{
		limitMB:   limitMB,
		stopCh:    make(chan struct{}),
		usageFunc: getCurrentMemoryUsage,
	}
}

// Start begins monitoring memory usage
func (mm *MemoryMonitor) Start() {
	mm.mu.Lock()
	defer mm.mu.Unlock()

	if mm.monitoring {
		return
	}

	mm.monitoring = true
	go mm.monitorLoop()
}

// Stop stops monitoring memory usage
func (mm *MemoryMonitor) Stop() {
	mm.mu.Lock()
	defer mm.mu.Unlock()

	if !mm.monitoring {
		return
	}

	close(mm.stopCh)
	mm.monitoring = false
}

// IsWithinLimit checks if current memory usage is within the configured limit
func (mm *MemoryMonitor) IsWithinLimit() bool {
	mm.mu.RLock()
	defer mm.mu.RUnlock()

	if mm.limitMB <= 0 {
		return true // No limit set
	}

	currentMB := mm.getCurrentMemoryMB()
	return currentMB <= float64(mm.limitMB)
}

// GetMemoryUsage returns current memory usage in MB
func (mm *MemoryMonitor) GetMemoryUsage() float64 {
	mm.mu.RLock()
	defer mm.mu.RUnlock()

	return mm.getCurrentMemoryMB()
}

// GetMemoryLimit returns the configured memory limit in MB
func (mm *MemoryMonitor) GetMemoryLimit() int64 {
	mm.mu.RLock()
	defer mm.mu.RUnlock()

	return mm.limitMB
}

// SetLimit updates the memory limit in MB
func (mm *MemoryMonitor) SetLimit(limitMB int64) {
	mm.mu.Lock()
	defer mm.mu.Unlock()

	mm.limitMB = limitMB
}

// getCurrentMemoryMB gets current memory usage in MB
func (mm *MemoryMonitor) getCurrentMemoryMB() float64 {
	usage := mm.usageFunc()
	return float64(usage) / (1024 * 1024) // Convert bytes to MB
}

// monitorLoop continuously monitors memory usage
func (mm *MemoryMonitor) monitorLoop() {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			if !mm.IsWithinLimit() {
				// Could trigger alerts or other actions when limit is exceeded
				// For now, we just log it (using internal logging)
			}
		case <-mm.stopCh:
			return
		}
	}
}

// getCurrentMemoryUsage returns the current memory usage in bytes
func getCurrentMemoryUsage() uint64 {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	// Using Sys as a representation of the memory reserved from the system
	return m.Sys
}

// CheckMemoryLimit returns an error if memory usage exceeds the limit
func CheckMemoryLimit(limitMB int64) error {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	currentMB := float64(m.Sys) / (1024 * 1024)

	if currentMB > float64(limitMB) {
		return &MemoryError{
			CurrentMB: currentMB,
			LimitMB:   float64(limitMB),
			Message:   fmt.Sprintf("memory usage %.2f MB exceeds limit of %.2f MB", currentMB, float64(limitMB)),
		}
	}

	return nil
}

// MemoryError represents a memory limit error
type MemoryError struct {
	CurrentMB float64
	LimitMB   float64
	Message   string
}

func (e *MemoryError) Error() string {
	return e.Message
}

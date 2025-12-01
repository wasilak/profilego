package core

import "context"

// BackendType represents the type of profiling backend
type BackendType string

const (
	PyroscopeBackend BackendType = "pyroscope"
	PprofBackend     BackendType = "pprof"
)

// ProfileType represents different types of profiles
type ProfileType string

const (
	ProfileCPU            ProfileType = "cpu"
	ProfileAllocObjects   ProfileType = "alloc_objects"
	ProfileAllocSpace     ProfileType = "alloc_space"
	ProfileInuseObjects   ProfileType = "inuse_objects"
	ProfileInuseSpace     ProfileType = "inuse_space"
	ProfileGoroutines     ProfileType = "goroutines"
	ProfileMutexCount     ProfileType = "mutex_count"
	ProfileMutexDuration  ProfileType = "mutex_duration"
	ProfileBlockCount     ProfileType = "block_count"
	ProfileBlockDuration  ProfileType = "block_duration"
)

// ProfilingState represents the initial state of profiling
type ProfilingState string

const (
	ProfilingEnabled  ProfilingState = "enabled"
	ProfilingDisabled ProfilingState = "disabled"
)
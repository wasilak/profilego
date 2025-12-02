package config

import (
	"github.com/wasilak/profilego/core"
	"testing"
	"time"
)

func TestConfigValidate(t *testing.T) {
	// Test with valid config
	validConfig := Config{
		ApplicationName: "test-app",
		ServerAddress:   "localhost:4040",
		Backend:         core.PyroscopeBackend,
	}

	err := validConfig.Validate()
	if err != nil {
		t.Errorf("Valid config should not return error: %v", err)
	}

	// Test with missing ApplicationName
	invalidConfig := Config{
		ApplicationName: "",
		ServerAddress:   "localhost:4040",
		Backend:         core.PyroscopeBackend,
	}

	err = invalidConfig.Validate()
	if err == nil {
		t.Error("Config with missing ApplicationName should return error")
	}

	// Test with missing ServerAddress for Pyroscope backend
	invalidConfig2 := Config{
		ApplicationName: "test-app",
		ServerAddress:   "",
		Backend:         core.PyroscopeBackend,
	}

	err = invalidConfig2.Validate()
	if err == nil {
		t.Error("Config with missing ServerAddress for Pyroscope backend should return error")
	}

	// Test with missing ServerAddress for Pprof backend (should be OK)
	pprofConfig := Config{
		ApplicationName: "test-app",
		ServerAddress:   "",
		Backend:         core.PprofBackend,
	}

	err = pprofConfig.Validate()
	if err != nil {
		t.Errorf("Config with missing ServerAddress for Pprof backend should not return error: %v", err)
	}
}

func TestDefaultConfig(t *testing.T) {
	// Check that DefaultConfig has expected values
	if DefaultConfig.ApplicationName != "my-app" {
		t.Errorf("Expected ApplicationName 'my-app', got '%s'", DefaultConfig.ApplicationName)
	}

	if DefaultConfig.Backend != core.PyroscopeBackend {
		t.Errorf("Expected Backend PyroscopeBackend, got '%s'", DefaultConfig.Backend)
	}

	if DefaultConfig.ServerAddress != "127.0.0.1:4040" {
		t.Errorf("Expected ServerAddress '127.0.0.1:4040', got '%s'", DefaultConfig.ServerAddress)
	}

	if DefaultConfig.InitialState != core.ProfilingEnabled {
		t.Errorf("Expected InitialState ProfilingEnabled, got '%s'", DefaultConfig.InitialState)
	}

	if DefaultConfig.MemoryLimitMB != 50 {
		t.Errorf("Expected MemoryLimitMB 50, got %d", DefaultConfig.MemoryLimitMB)
	}

	if DefaultConfig.LogLevel != "info" {
		t.Errorf("Expected LogLevel 'info', got '%s'", DefaultConfig.LogLevel)
	}

	if DefaultConfig.Timeout != 10*time.Second {
		t.Errorf("Expected Timeout 10s, got %v", DefaultConfig.Timeout)
	}

	if DefaultConfig.EnableTLS != false {
		t.Errorf("Expected EnableTLS false, got %v", DefaultConfig.EnableTLS)
	}

	if DefaultConfig.SkipTLSVerify != false {
		t.Errorf("Expected SkipTLSVerify false, got %v", DefaultConfig.SkipTLSVerify)
	}

	// Check ProfileTypes slice has expected length
	expectedProfileTypes := 10
	if len(DefaultConfig.ProfileTypes) != expectedProfileTypes {
		t.Errorf("Expected %d ProfileTypes, got %d", expectedProfileTypes, len(DefaultConfig.ProfileTypes))
	}
}

func TestConfigError(t *testing.T) {
	err := &ConfigError{
		Field:   "test",
		Message: "test message",
	}

	expected := "config error (test): test message"
	if err.Error() != expected {
		t.Errorf("Expected error '%s', got '%s'", expected, err.Error())
	}
}

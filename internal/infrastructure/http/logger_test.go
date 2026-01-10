package http

import (
	"testing"

	"github.com/go-resty/resty/v2"
	"github.com/zevwings/workflow/internal/logging"
)

// TestLogrusLogger_ImplementsRestyLogger tests whether LogrusLogger implements resty.Logger interface
func TestLogrusLogger_ImplementsRestyLogger(t *testing.T) {
	var _ resty.Logger = (*LogrusLogger)(nil)
}

// TestNewLogrusLogger tests NewLogrusLogger function
func TestNewLogrusLogger(t *testing.T) {
	logger := NewLogrusLogger()
	if logger == nil {
		t.Fatal("NewLogrusLogger() returned nil")
	}
}

// TestLogrusLogger_Methods tests all methods of LogrusLogger
func TestLogrusLogger_Methods(t *testing.T) {
	// Initialize logging system (required for testing)
	logging.Init("debug", "text", nil)

	logger := NewLogrusLogger()

	// Test Errorf
	logger.Errorf("Test error message: %s", "error")

	// Test Warnf
	logger.Warnf("Test warning message: %s", "warning")

	// Test Debugf
	logger.Debugf("Test debug message: %s", "debug")
}

// TestLogrusLogger_WithRestyClient tests LogrusLogger integration with Resty client
func TestLogrusLogger_WithRestyClient(t *testing.T) {
	// Initialize logging system (required for testing)
	logging.Init("debug", "text", nil)

	// Create Resty client and set Logger
	client := resty.New()
	logger := NewLogrusLogger()
	client.SetLogger(logger)

	// Verify Logger is set (calling SetLogger should not return error)
	// Resty will internally use the Logger we set
	if logger == nil {
		t.Error("Logger is nil")
	}
}

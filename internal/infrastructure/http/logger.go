package http

import (
	"github.com/go-resty/resty/v2"
	"github.com/zevwings/workflow/internal/logging"
)

// LogrusLogger implements resty.Logger interface, forwards Resty logs to Logrus
//
// This adapter converts go-resty's logging interface to the project's internal logging package.
// Through this adapter, all Resty log output will be recorded through the logging package.
//
// Usage:
//
//	import adapterhttp "github.com/zevwings/workflow/internal/infrastructure/http"
//
//	client := resty.New()
//	client.SetLogger(adapterhttp.NewLogrusLogger())
type LogrusLogger struct{}

// NewLogrusLogger creates a new LogrusLogger instance
//
// Returns:
//   - *LogrusLogger: LogrusLogger instance that implements resty.Logger interface
func NewLogrusLogger() *LogrusLogger {
	return &LogrusLogger{}
}

// Errorf records error log
//
// Implements Errorf method of resty.Logger interface.
// Forwards Resty error logs to logging and identifies module as "http".
//
// Parameters:
//   - format: Log format string
//   - v: Format parameters
func (l *LogrusLogger) Errorf(format string, v ...interface{}) {
	logging.WithField("module", "http").Errorf(format, v...)
}

// Warnf records warning log
//
// Implements Warnf method of resty.Logger interface.
// Forwards Resty warning logs to logging and identifies module as "http".
//
// Parameters:
//   - format: Log format string
//   - v: Format parameters
func (l *LogrusLogger) Warnf(format string, v ...interface{}) {
	logging.WithField("module", "http").Warnf(format, v...)
}

// Debugf records debug log
//
// Implements Debugf method of resty.Logger interface.
// Forwards Resty debug logs to logging and identifies module as "http".
//
// Parameters:
//   - format: Log format string
//   - v: Format parameters
func (l *LogrusLogger) Debugf(format string, v ...interface{}) {
	logging.WithField("module", "http").Debugf(format, v...)
}

// Verify that LogrusLogger implements resty.Logger interface
var _ resty.Logger = (*LogrusLogger)(nil)

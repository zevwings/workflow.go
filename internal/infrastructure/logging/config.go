package logging

import (
	"io"
	"os"
	"path/filepath"

	"github.com/sirupsen/logrus"
	"github.com/zevwings/workflow/internal/config"
	"github.com/zevwings/workflow/internal/logging"
)

// InitLogging initializes the logging system
func InitLogging() {
	// First get log directory and configuration, then decide temporary logger output target
	var logDir string
	workflowStateDir, err := config.StateDir()
	if err == nil {
		logDir = filepath.Join(workflowStateDir, "logs")
	} else {
		logDir = ""
	}

	// Try to load configuration to get log level (silent load, no log output)
	var logLevel string
	manager, err := config.Global()
	if err == nil {
		// Try to load configuration, ignore errors (config file may not exist)
		_ = manager.Load()
		logLevel = manager.LogConfig.Level
	}

	// Before initializing logging system, clean up old log files first
	if logDir != "" {
		_ = ClearLogFiles(logDir)
	}

	// Decide temporary logger output target based on configuration
	// If level is empty, temporary logger should output to file (if logDir is not empty) or not output
	var tempLoggerOutput io.Writer
	if logLevel == "" {
		// level is empty, only write to file, don't output to console
		if logDir != "" {
			// Ensure log directory exists
			_ = os.MkdirAll(logDir, 0755)
			// Create temporary log file for logs during initialization
			tempLogFile := filepath.Join(logDir, "init.log")
			file, err := os.OpenFile(tempLogFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
			if err == nil {
				tempLoggerOutput = file
			} else {
				// If unable to create file, use io.Discard to discard output
				tempLoggerOutput = io.Discard
			}
		} else {
			tempLoggerOutput = io.Discard
		}
	} else {
		// level is not empty, output to console
		tempLoggerOutput = os.Stdout
	}

	// Create temporary logger
	var logger *logrus.Logger
	if logging.Logger != nil {
		logger = logging.Logger
	} else {
		logger = logrus.New()
		logger.SetOutput(tempLoggerOutput)
		logger.SetLevel(logrus.InfoLevel)
	}

	logger.Info("Initializing logging system")

	// If config manager creation fails, use default settings (level is empty, only write to file)
	if manager == nil {
		logger.WithError(err).Warn("Failed to create config manager, using default logging settings")
		logging.InitWithFiles("", "text", nil, logDir, false)
		return
	}

	// Get log level and format from configuration
	// If level is empty, keep it as empty string, so it will only write to file and only record info/warn/error levels
	// If level is not empty, control console output level based on level
	logFormat := "text" // Default text format

	// If level is not empty, output to console; if empty, don't output to console (InitWithFiles will handle automatically)
	consoleOut := logLevel != ""
	logging.InitWithFiles(logLevel, logFormat, nil, logDir, consoleOut)

	// Record initialization completion (using initialized logger)
	initLogger := logging.GetLogger()
	initLogger.WithFields(logging.Fields{
		"level":   logLevel,
		"format":  logFormat,
		"log_dir": logDir,
	}).Info("Logging system initialized")
}

// ClearLogFiles deletes all log files in the log directory
func ClearLogFiles(logDir string) error {
	// Check if directory exists
	if _, err := os.Stat(logDir); os.IsNotExist(err) {
		// Directory does not exist, no need to clean
		return nil
	}

	// Read directory contents
	entries, err := os.ReadDir(logDir)
	if err != nil {
		return err
	}

	// Delete all .log files
	for _, entry := range entries {
		if !entry.IsDir() && filepath.Ext(entry.Name()) == ".log" {
			logFile := filepath.Join(logDir, entry.Name())
			if err := os.Remove(logFile); err != nil {
				// Record error but continue deleting other files
				_ = err
			}
		}
	}

	return nil
}

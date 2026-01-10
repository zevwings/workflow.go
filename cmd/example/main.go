//go:build example

package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/zevwings/workflow/internal/commands/example"
	"github.com/zevwings/workflow/internal/config"
	"github.com/zevwings/workflow/internal/logging"
)

var (
	version   string = "0.1.0"
	buildDate string = "unknown"
	gitCommit string = "unknown"
)

func main() {
	// Initialize logging system
	initLogging()

	// Create root command
	rootCmd := &cobra.Command{
		Use:   "example",
		Short: "Workflow CLI Example - Demonstration and testing tool",
		Long: `Workflow CLI Example is a demonstration and testing tool,
used to showcase various features and components of Workflow CLI.

This tool includes the following demo commands:
- demo-prompt: Demonstrate interactive features of all Prompt components
- demo-form: Demonstrate interactive form features of Form module`,
		Version: version,
	}

	// Register example commands
	rootCmd.AddCommand(example.NewDemoCmd())
	rootCmd.AddCommand(example.NewDemoFormCmd())

	// Set version template
	rootCmd.SetVersionTemplate(fmt.Sprintf("example version %s\nBuild Date: %s\nGit Commit: %s\n", version, buildDate, gitCommit))

	// Execute command
	if err := rootCmd.Execute(); err != nil {
		logger := logging.GetLogger()
		logger.WithError(err).Error("Command execution failed")
		os.Exit(1)
	}
}

// initLogging initializes the logging system
func initLogging() {
	// If Logger is not yet initialized, initialize a default one first
	if logging.Logger == nil {
		logging.Init("info", "text", nil)
	}
	// Use global logger (module logger is not yet initialized at this point)
	logger := logging.Logger
	logger.Info("Initializing logging system")

	// Get log directory (using XDG_STATE_HOME)
	var logDir string
	workflowStateDir, err := config.StateDir()
	if err == nil {
		logDir = filepath.Join(workflowStateDir, "logs")
	} else {
		logger.WithError(err).Warn("Failed to get XDG state directory, using empty log directory")
		logDir = ""
	}

	// Try to load configuration to get log level
	manager, err := config.Global()
	if err != nil {
		// If config manager creation fails, use default settings
		logger.WithError(err).Warn("Failed to create config manager, using default logging settings")
		logging.InitWithFiles("info", "text", nil, logDir, false)
		return
	}

	// Try to load configuration
	if err := manager.Load(); err != nil {
		// Use default settings when config file does not exist
		logger.Debug("Config file not found, using default logging settings")
		logging.InitWithFiles("info", "text", nil, logDir, false)
		return
	}

	// Get log level and format from configuration
	// Directly access configuration fields
	logLevel := manager.LogConfig.Level
	if logLevel == "" {
		logLevel = "info"
	}

	logFormat := "text" // Default text format
	logging.InitWithFiles(logLevel, logFormat, nil, logDir, false)

	// Record initialization completion
	logger.WithFields(logrus.Fields{
		"level":   logLevel,
		"format":  logFormat,
		"log_dir": logDir,
	}).Info("Logging system initialized")
}

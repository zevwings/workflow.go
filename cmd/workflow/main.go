package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/zevwings/workflow/internal/commands"
	configCmd "github.com/zevwings/workflow/internal/commands/config"
	repoCmd "github.com/zevwings/workflow/internal/commands/repo"
	infrastructurelogging "github.com/zevwings/workflow/internal/infrastructure/logging"
	"github.com/zevwings/workflow/internal/logging"
)

var (
	version   string = "0.0.1"
	buildDate string = "2026-01-10"
	gitCommit string = "0000000000000000000000000000000000000000"
)

func main() {
	// Initialize logging system
	infrastructurelogging.InitLogging()

	// Create root command
	rootCmd := &cobra.Command{
		Use:   "workflow",
		Short: "Workflow CLI - Git workflow automation tool",
		Long: `Workflow CLI is a powerful Git workflow automation tool,
supporting PR management, Jira integration, LLM integration, and more.`,
		Version: version,
	}

	// Register subcommands
	rootCmd.AddCommand(commands.NewSetupCmd())
	rootCmd.AddCommand(configCmd.NewConfigCmd())
	rootCmd.AddCommand(repoCmd.NewRepoCmd())
	rootCmd.AddCommand(commands.NewCheckCmd())
	rootCmd.AddCommand(commands.NewVersionCmd(version, buildDate, gitCommit))

	// Set version template
	rootCmd.SetVersionTemplate(fmt.Sprintf("workflow version %s\nBuild Date: %s\nGit Commit: %s\n", version, buildDate, gitCommit))

	// Execute command
	if err := rootCmd.Execute(); err != nil {
		logger := logging.GetLogger()
		logger.WithError(err).Error("Command execution failed")
		os.Exit(1)
	}
}

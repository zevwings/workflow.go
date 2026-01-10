package commands

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/zevwings/workflow/internal/config"
	"github.com/zevwings/workflow/internal/infrastructure/verify"
	"github.com/zevwings/workflow/internal/prompt"
)

// NewCheckCmd creates the check command
func NewCheckCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "check",
		Short: "View current configuration and run environment checks",
		Long: `Display the contents of the current Workflow CLI configuration file, grouped logically, with sensitive information automatically masked.
Also verify configuration validity and run environment checks, including:
- Git status check
- Network connection check
- Configuration file check`,
		RunE: runCheck,
	}

	return cmd
}

func runCheck(cmd *cobra.Command, args []string) error {
	out := prompt.GetMessage()
	out.Info("Starting check command")
	// Display configuration information
	manager, err := config.Global()
	if err != nil {
		return fmt.Errorf("failed to create config manager: %w", err)
	}

	if err := manager.Load(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			out.Warning("Config file not found")
		} else {
			return fmt.Errorf("failed to load configuration: %w", err)
		}
	} else {
		configPath := manager.GetConfigPath()

		// Display title and configuration path
		out.Break('=', 80, "Current Configuration")
		out.Break()
		out.Info("Workflow config: %q", configPath)
		out.Break()
	}

	// Environment check
	verify.VerifyEnvironment()

	verify.VerifyLogConfig(manager.LogConfig)
	verify.VerifyLLMConfig(manager.LLMConfig)
	verify.VerifyJiraConfig(manager.JiraConfig)
	verify.VerifyGitHubConfig(manager.GitHubConfig)

	return nil
}

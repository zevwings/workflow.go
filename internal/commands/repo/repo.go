package repo

import (
	"github.com/spf13/cobra"
)

// NewRepoCmd creates the repo command
func NewRepoCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "repo",
		Short: "Repository configuration management",
		Long:  `Manage repository-level configuration files.`,
	}

	// Add subcommands
	cmd.AddCommand(NewSetupCmd())
	cmd.AddCommand(NewShowCmd())

	return cmd
}

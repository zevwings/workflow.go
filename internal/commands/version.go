package commands

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/zevwings/workflow/internal/prompt"
)

// NewVersionCmd creates the version command
func NewVersionCmd(version, buildDate, gitCommit string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "version",
		Short: "Display version information",
		Long:  `Display Workflow CLI version information, including version number, build date, and Git commit hash.`,
		Run: func(cmd *cobra.Command, args []string) {
			out := prompt.GetMessage()
			out.Print("%s", "Workflow CLI")
			out.Break()
			out.Print("%s", fmt.Sprintf("Version:    %s", version))
			out.Print("%s", fmt.Sprintf("Build Date: %s", buildDate))
			out.Print("%s", fmt.Sprintf("Git Commit: %s", gitCommit))
		},
	}

	return cmd
}

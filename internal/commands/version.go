package commands

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/zevwings/workflow/internal/prompt"
)

// NewVersionCmd 创建 version 命令
func NewVersionCmd(version, buildDate, gitCommit string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "version",
		Short: "显示版本信息",
		Long:  `显示 Workflow CLI 的版本信息，包括版本号、构建日期和 Git 提交哈希。`,
		Run: func(cmd *cobra.Command, args []string) {
			out := prompt.NewMessage(false)
			out.Println("%s", "Workflow CLI")
			out.Println("%s", "")
			out.Println("%s", fmt.Sprintf("Version:    %s", version))
			out.Println("%s", fmt.Sprintf("Build Date: %s", buildDate))
			out.Println("%s", fmt.Sprintf("Git Commit: %s", gitCommit))
		},
	}

	return cmd
}

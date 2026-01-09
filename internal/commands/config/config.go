package config

import (
	"github.com/spf13/cobra"
)

// NewConfigCmd 创建 config 命令
func NewConfigCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "config",
		Short: "配置管理",
		Long:  `管理 Workflow CLI 配置文件。`,
	}

	cmd.AddCommand(newConfigShowCmd())

	return cmd
}

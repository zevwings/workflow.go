package config

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/zevwings/workflow/internal/config"
	"github.com/zevwings/workflow/internal/logging"
	"github.com/zevwings/workflow/internal/prompt"
)

// config show 子命令
func newConfigShowCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "show",
		Short: "查看当前配置",
		Long:  `显示当前 Workflow CLI 配置文件的内容，按逻辑分组显示，敏感信息会自动掩码。同时验证配置有效性。`,
		RunE: func(cmd *cobra.Command, args []string) error {
			logger := logging.GetLogger()
			logger.Info("Starting config show command")

			out := prompt.NewMessage(false)

			manager, err := config.Global()
			if err != nil {
				return fmt.Errorf("创建配置管理器失败: %w", err)
			}

			if err := manager.Load(); err != nil {
				if _, ok := err.(viper.ConfigFileNotFoundError); ok {
					logger.Warn("Config file not found")
					out.Error("配置文件不存在")
					return fmt.Errorf("配置文件不存在")
				}
				return fmt.Errorf("加载配置失败: %w", err)
			}

			configPath := manager.GetConfigPath()

			// 显示标题和配置路径
			out.Break('=', 80, "Current Configuration")
			out.Break()
			out.Info("Workflow config: %q", configPath)
			out.Break()

			// 格式化显示配置
			showConfig(out, manager.Config)

			return nil
		},
	}

	return cmd
}

// showConfig 格式化显示配置
func showConfig(out *prompt.Message, cfg *config.GlobalConfig) {
	showConfigWithOptions(out, cfg, true)
}

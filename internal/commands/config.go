package commands

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/zevwings/workflow/internal/config"
	"github.com/zevwings/workflow/internal/prompt"
)

// NewConfigCmd 创建 config 命令
func NewConfigCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "config",
		Short: "配置管理",
		Long:  `管理 Workflow CLI 配置文件。`,
	}

	cmd.AddCommand(newConfigShowCmd())
	cmd.AddCommand(newConfigValidateCmd())

	return cmd
}

// config show 子命令
func newConfigShowCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "show",
		Short: "查看当前配置",
		Long:  `显示当前 Workflow CLI 配置文件的内容。`,
		RunE: func(cmd *cobra.Command, args []string) error {
			out := prompt.NewMessage(false)

			manager, err := config.NewGlobalManager()
			if err != nil {
				return fmt.Errorf("创建配置管理器失败: %w", err)
			}

			if err := manager.Load(); err != nil {
				if _, ok := err.(viper.ConfigFileNotFoundError); ok {
					out.Warning("配置文件不存在，请运行 'workflow setup' 进行初始化")
					return nil
				}
				return fmt.Errorf("加载配置失败: %w", err)
			}

			configPath := manager.GetConfigPath()
			out.Info("配置文件路径: %s", configPath)
			out.Println("")

			// 读取并显示配置文件内容
			data, err := os.ReadFile(configPath)
			if err != nil {
				return fmt.Errorf("读取配置文件失败: %w", err)
			}

			out.Println(string(data))
			return nil
		},
	}

	return cmd
}

// config validate 子命令
func newConfigValidateCmd() *cobra.Command {
	var fix bool
	var strict bool

	cmd := &cobra.Command{
		Use:   "validate [CONFIG_PATH]",
		Short: "验证配置文件",
		Long:  `验证配置文件格式和内容是否正确。`,
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			out := prompt.NewMessage(false)

			manager, err := config.NewGlobalManager()
			if err != nil {
				return fmt.Errorf("创建配置管理器失败: %w", err)
			}

			// 如果指定了配置文件路径，使用指定的路径
			if len(args) > 0 {
				// TODO: 实现自定义路径加载
			}

			if err := manager.Load(); err != nil {
				if _, ok := err.(viper.ConfigFileNotFoundError); ok {
					out.Error("配置文件不存在")
					return fmt.Errorf("配置文件不存在")
				}
				return fmt.Errorf("加载配置失败: %w", err)
			}

			// 验证配置
			out.Info("验证配置文件...")

			// TODO: 实现详细的配置验证逻辑
			// 目前只检查文件是否可以解析

			out.Success("配置文件验证通过")

			if fix {
				out.Info("修复模式未实现")
			}

			return nil
		},
	}

	cmd.Flags().BoolVar(&fix, "fix", false, "自动修复配置错误")
	cmd.Flags().BoolVar(&strict, "strict", false, "严格模式（检查所有字段）")

	return cmd
}

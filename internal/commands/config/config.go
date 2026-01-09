package config

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/zevwings/workflow/internal/config"

	// "github.com/zevwings/workflow/internal/jira"
	"github.com/zevwings/workflow/internal/logging"
	// githubAuth "github.com/zevwings/workflow/internal/pr/github"
	"github.com/zevwings/workflow/internal/prompt"
	"github.com/zevwings/workflow/internal/util"
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
			out.Println("%s", "========  Current Configuration ========")
			out.Println("%s", "")
			out.Info("Workflow config: %q", configPath)
			out.Println("%s", "")

			// 格式化显示配置
			showConfig(out, manager.Config)

			return nil
		},
	}

	return cmd
}

// showConfig 格式化显示配置
func showConfig(out *prompt.Message, cfg *config.GlobalConfig) {
	// 日志配置
	if cfg.Log.Level != "" {
		out.Println("%s", "------------------------------------------  Configuration ------------------------------------------")
		out.Println("%s", fmt.Sprintf("Log Output Folder Name: %s", cfg.Log.Level))
		out.Println("%s", "")
	}

	// LLM 配置
	if cfg.LLM.Provider != "" {
		out.Info("LLM Configuration")
		table := prompt.NewTable([]string{"Provider", "Model", "Key", "Output Language"})

		var model, key, language string
		language = cfg.LLM.Language
		if language == "" {
			language = "en"
		}

		switch cfg.LLM.Provider {
		case "openai":
			model = cfg.LLM.OpenAI.Model
			if model == "" {
				model = "gpt-3.5-turbo"
			}
			key = util.MaskSensitiveValue(cfg.LLM.OpenAI.APIKey)
		case "deepseek":
			model = cfg.LLM.DeepSeek.Model
			if model == "" {
				model = "deepseek-chat"
			}
			key = util.MaskSensitiveValue(cfg.LLM.DeepSeek.APIKey)
		case "proxy":
			model = cfg.LLM.Proxy.Model
			if cfg.LLM.Proxy.URL != "" {
				model = fmt.Sprintf("%s(%s)", model, cfg.LLM.Proxy.URL)
			}
			key = util.MaskSensitiveValue(cfg.LLM.Proxy.APIKey)
		default:
			model = "-"
			key = "-"
		}

		table.AddRow([]string{cfg.LLM.Provider, model, key, language})
		table.Render()
		out.Println("%s", "")
	}

	// Jira 配置
	if cfg.Jira.Email != "" || cfg.Jira.APIToken != "" || cfg.Jira.ServiceAddress != "" {
		opts := DefaultVerifyOptions()
		VerifyJiraConfig(out, &cfg.Jira, opts)
	}

	// GitHub 配置
	if len(cfg.GitHub.Accounts) > 0 {
		opts := DefaultVerifyOptions()
		VerifyGitHubConfig(out, &cfg.GitHub, opts)
	}
}

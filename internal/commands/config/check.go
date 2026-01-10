package config

import (
	"context"
	"fmt"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/zevwings/workflow/internal/config"
	"github.com/zevwings/workflow/internal/http"
	"github.com/zevwings/workflow/internal/logging"
	"github.com/zevwings/workflow/internal/prompt"
)

// NewCheckCmd 创建 check 命令
func NewCheckCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "check",
		Short: "查看当前配置并运行环境检查",
		Long: `显示当前 Workflow CLI 配置文件的内容，按逻辑分组显示，敏感信息会自动掩码。
同时验证配置有效性，并运行环境检查，包括：
- Git 状态检查
- 网络连接检查
- 配置文件检查`,
		RunE: runCheck,
	}

	return cmd
}

func runCheck(cmd *cobra.Command, args []string) error {
	logger := logging.GetLogger()
	logger.Info("Starting check command")

	out := prompt.GetMessage()

	// 显示配置信息
	manager, err := config.Global()
	if err != nil {
		return fmt.Errorf("创建配置管理器失败: %w", err)
	}

	if err := manager.Load(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			logger.Warn("Config file not found")
			out.Error("配置文件不存在")
		} else {
			return fmt.Errorf("加载配置失败: %w", err)
		}
	} else {
		configPath := manager.GetConfigPath()

		// 显示标题和配置路径
		out.Break('=', 80, "Current Configuration")
		out.Break()
		out.Info("Workflow config: %q", configPath)
		out.Break()
	}

	// 环境检查
	out.Break('=', 80, "Environment Check")
	out.Break()
	out.Info("开始环境检查...")
	out.Break()

	table := prompt.NewTable([]string{"检查项", "状态", "说明"})

	// 验证配置文件
	configOK := verifyConfig()
	if configOK {
		table.AddRow([]string{"配置文件", "✓", "配置文件存在且有效"})
	} else {
		table.AddRow([]string{"配置文件", "✗", "配置文件不存在或无效"})
	}

	// 验证网络连接
	networkOK := verifyNetwork()
	if networkOK {
		table.AddRow([]string{"网络连接", "✓", "网络连接正常"})
	} else {
		table.AddRow([]string{"网络连接", "✗", "网络连接失败"})
	}

	out.Break()
	table.Render()
	out.Break()

	VerifyLogConfig(manager.LogConfig)
	VerifyLLMConfig(manager.LLMConfig)
	VerifyLLMConnection(manager.LLMConfig)
	VerifyJiraConfig(manager.JiraConfig)
	VerifyGitHubConfig(manager.GitHubConfig)

	return nil
}

func verifyConfig() bool {
	manager, err := config.Global()
	if err != nil {
		return false
	}

	if err := manager.Load(); err != nil {
		return false
	}

	return true
}

func verifyNetwork() bool {
	client := http.Global()

	spinner := prompt.NewSpinner("Verifying network connection...")
	var respStatusCode int
	err := spinner.Do(func() error {
		// 检查 GitHub 连接（使用超时上下文）
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		// 使用 resty 客户端发送请求
		restyClient := client.GetRestyClient()
		resp, err := restyClient.R().SetContext(ctx).Get("https://api.github.com")
		if err != nil {
			return err
		}

		// 检查响应状态码
		respStatusCode = resp.StatusCode()
		return nil
	})
	spinner.Stop()

	if err != nil {
		return false
	}

	return respStatusCode == 200
}

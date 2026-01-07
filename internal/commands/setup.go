package commands

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/zevwings/workflow/internal/config"
	"github.com/zevwings/workflow/internal/prompt"
)

// NewSetupCmd 创建 setup 命令
func NewSetupCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "setup",
		Short: "初始化或更新配置（交互式）",
		Long: `交互式初始化或更新 Workflow CLI 配置。
此命令将引导您完成配置过程，包括 GitHub、Jira 等服务的配置。`,
		RunE: runSetup,
	}

	return cmd
}

func runSetup(cmd *cobra.Command, args []string) error {
	msg := prompt.NewMessage(false)

	msg.Info("欢迎使用 Workflow CLI 配置向导")
	msg.Println("")

	// 创建配置管理器
	manager, err := config.NewGlobalManager()
	if err != nil {
		return fmt.Errorf("创建配置管理器失败: %w", err)
	}

	// 尝试加载现有配置
	var cfg config.GlobalConfig
	if err := manager.Load(); err == nil {
		// 如果配置文件存在，询问是否更新
		update, err := prompt.AskConfirm("检测到现有配置文件，是否要更新？", false)
		if err != nil {
			return err
		}
		if !update {
			msg.Info("已取消配置更新")
			return nil
		}
		// 加载现有配置
		cfg = config.GlobalConfig{
			Log: config.LogConfig{
				Level: manager.GetString("log.level"),
			},
		}
	}

	// 收集用户信息
	msg.Info("配置用户信息")

	name, err := prompt.AskInput("请输入您的姓名", cfg.User.Name)
	if err != nil {
		return err
	}

	email, err := prompt.AskInput("请输入您的邮箱", cfg.User.Email)
	if err != nil {
		return err
	}

	cfg.User.Name = name
	cfg.User.Email = email

	// 询问是否配置 GitHub
	configureGitHub, err := prompt.AskConfirm("是否配置 GitHub？", true)
	if err != nil {
		return err
	}

	if configureGitHub {
		githubToken, err := prompt.AskPassword("请输入 GitHub Personal Access Token")
		if err != nil {
			return err
		}

		if githubToken != "" {
			cfg.GitHub.Accounts = []config.GitHubAccount{
				{
					Name:  "default",
					Token: githubToken,
				},
			}
			cfg.GitHub.Current = "default"
		}
	}

	// 询问是否配置 Jira
	configureJira, err := prompt.AskConfirm("是否配置 Jira？", false)
	if err != nil {
		return err
	}

	if configureJira {
		jiraURL, err := prompt.AskInput("请输入 Jira URL", cfg.Jira.URL)
		if err != nil {
			return err
		}

		jiraUsername, err := prompt.AskInput("请输入 Jira 用户名", cfg.Jira.Username)
		if err != nil {
			return err
		}

		jiraToken, err := prompt.AskPassword("请输入 Jira API Token")
		if err != nil {
			return err
		}

		cfg.Jira.URL = jiraURL
		cfg.Jira.Username = jiraUsername
		cfg.Jira.Token = jiraToken
	}

	// 保存配置
	if err := manager.Save(&cfg); err != nil {
		return fmt.Errorf("保存配置失败: %w", err)
	}

	msg.Success("配置已保存到: %s", manager.GetConfigPath())
	msg.Info("您可以使用 'workflow config show' 查看配置")

	return nil
}

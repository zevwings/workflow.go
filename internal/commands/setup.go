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
	manager, err := config.Global()
	if err != nil {
		return fmt.Errorf("创建配置管理器失败: %w", err)
	}

	// 尝试加载现有配置
	var cfg *config.GlobalConfig
	configExists := false
	if err := manager.Load(); err == nil {
		configExists = true
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
		// 直接访问 Config 字段
		cfg = manager.Config
	} else {
		// 配置文件不存在，创建新配置
		cfg = &config.GlobalConfig{
			Log: config.LogConfig{
				Level: "info",
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

	// 如果已有 GitHub 账号，询问是否更新
	hasGitHub := len(cfg.GitHub.Accounts) > 0
	if hasGitHub && configExists {
		updateGitHub, err := prompt.AskConfirm("是否更新 GitHub 配置？", false)
		if err != nil {
			return err
		}
		if !updateGitHub {
			// 保留现有 GitHub 配置
		} else {
			hasGitHub = false // 标记为需要重新配置
		}
	}

	// 如果已有 Jira 配置，询问是否更新
	hasJira := cfg.Jira.URL != "" || cfg.Jira.Username != "" || cfg.Jira.Token != ""
	if hasJira && configExists {
		updateJira, err := prompt.AskConfirm("是否更新 Jira 配置？", false)
		if err != nil {
			return err
		}
		if !updateJira {
			// 保留现有 Jira 配置
		} else {
			hasJira = false // 标记为需要重新配置
		}
	}

	// 询问是否配置 GitHub
	if !hasGitHub {
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
				// 如果已有账号，添加到列表；否则创建新列表
				if len(cfg.GitHub.Accounts) == 0 {
					cfg.GitHub.Accounts = []config.GitHubAccount{
						{
							Name:  "default",
							Token: githubToken,
						},
					}
					cfg.GitHub.Current = "default"
				} else {
					// 更新第一个账号的 token
					cfg.GitHub.Accounts[0].Token = githubToken
					if cfg.GitHub.Current == "" {
						cfg.GitHub.Current = cfg.GitHub.Accounts[0].Name
					}
				}
			}
		}
	}

	// 询问是否配置 Jira
	if !hasJira {
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
	}

	// 保存配置
	manager.Config = cfg
	if err := manager.Save(); err != nil {
		return fmt.Errorf("保存配置失败: %w", err)
	}

	msg.Success("配置已保存到: %s", manager.GetConfigPath())
	msg.Info("您可以使用 'workflow config show' 查看配置")

	return nil
}

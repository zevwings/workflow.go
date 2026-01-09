package config

import (
	"fmt"

	"github.com/zevwings/workflow/internal/config"
	"github.com/zevwings/workflow/internal/jira"
	"github.com/zevwings/workflow/internal/pr/github"
	"github.com/zevwings/workflow/internal/prompt"
	"github.com/zevwings/workflow/internal/util"
)

// VerifyOptions 验证选项
// 保留此结构体以保持 API 兼容性，但当前不需要任何选项
type VerifyOptions struct{}

// DefaultVerifyOptions 返回默认验证选项
func DefaultVerifyOptions() *VerifyOptions {
	return &VerifyOptions{}
}

// VerifyJiraConfig 验证 Jira 配置
func VerifyJiraConfig(msg *prompt.Message, jiraConfig *config.JiraConfig, opts *VerifyOptions) {
	if jiraConfig.Email == "" && jiraConfig.APIToken == "" && jiraConfig.ServiceAddress == "" {
		return
	}

	if opts == nil {
		opts = DefaultVerifyOptions()
	}

	msg.Info("Jira Configuration")
	table := prompt.NewTable([]string{"Email", "Service Address", "API Token"})

	// 验证 Jira 认证
	jiraConfigForAuth := &jira.Config{
		ServiceAddress: jiraConfig.ServiceAddress,
		Email:          jiraConfig.Email,
		APIToken:       jiraConfig.APIToken,
	}

	var jiraResult *jira.AuthResult
	var err error

	spinner := prompt.NewSpinner("Verifying Jira configuration...")
	err = spinner.Do(func() error {
		jiraResult, err = jira.ValidateAuth(jiraConfigForAuth)
		return err
	})

	spinner.Stop()

	if err != nil {
		msg.Warning("Jira verification error: %v", err)
	}

	table.AddRow([]string{
		jiraConfig.Email,
		jiraConfig.ServiceAddress,
		util.MaskSensitiveValue(jiraConfig.APIToken),
	})
	table.Render()

	if jiraResult != nil && jiraResult.Valid {
		if accountID, ok := jiraResult.Details["account_id"].(string); ok && accountID != "" {
			msg.Success("Jira verified successfully! Email: %s (Account ID: %s)", jiraConfig.Email, accountID)
		} else {
			msg.Success("Jira verified successfully! Email: %s", jiraConfig.Email)
		}
	} else if jiraResult != nil {
		msg.Error("Jira verification failed: %s", jiraResult.Message)
	}
	msg.Println("%s", "")
}

// VerifyGitHubConfig 验证 GitHub 配置
func VerifyGitHubConfig(msg *prompt.Message, githubConfig *config.GitHubConfig, opts *VerifyOptions) bool {
	if len(githubConfig.Accounts) == 0 {
		return true
	}

	if opts == nil {
		opts = DefaultVerifyOptions()
	}

	msg.Info("GitHub Configuration")
	table := prompt.NewTable([]string{"Name", "Email", "API Token", "Status", "Verification"})

	allValid := true
	for _, account := range githubConfig.Accounts {
		status := ""
		if account.Name == githubConfig.Current {
			status = "Current"
		}

		// 验证 go-github
		var githubResult *github.AuthResult
		var githubErr error

		// 检查配置完整性
		if account.APIToken == "" {
			githubResult = &github.AuthResult{
				Valid:   false,
				Message: "GitHub API Token 未配置",
				Details: make(map[string]interface{}),
			}
		} else {
			spinner := prompt.NewSpinner(fmt.Sprintf("Verifying go-github for %s...", account.Name))
			githubErr = spinner.Do(func() error {
				githubResult, githubErr = github.ValidateAuth(account.APIToken)
				return githubErr
			})
			spinner.Stop()

			// 如果配置中有 email 但 API 返回中没有，使用配置中的 email
			if githubResult != nil && account.Email != "" && githubResult.Details["email"] == nil {
				if githubResult.Details == nil {
					githubResult.Details = make(map[string]interface{})
				}
				githubResult.Details["email"] = account.Email
			}
		}

		// 判断验证结果
		githubValid := githubResult != nil && githubResult.Valid && githubErr == nil

		// 显示验证结果（compact 格式：github 或留空）
		var verification string
		if githubValid {
			verification = "github"
		} else {
			verification = ""
			allValid = false
		}

		// 如果 github 验证失败，设置 allValid = false
		if !githubValid {
			allValid = false
		}

		table.AddRow([]string{
			account.Name,
			account.Email,
			util.MaskSensitiveValue(account.APIToken),
			status,
			verification,
		})
	}

	table.Render()
	if allValid {
		msg.Success("All %d GitHub account(s) verified successfully!", len(githubConfig.Accounts))
	} else {
		msg.Warning("Some GitHub account(s) verification failed. Please check the configuration.")
	}
	msg.Println("%s", "")

	return allValid
}

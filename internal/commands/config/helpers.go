package config

import (
	"fmt"
	"strings"

	llmadapter "github.com/zevwings/workflow/internal/adapter/llm"
	"github.com/zevwings/workflow/internal/config"
	"github.com/zevwings/workflow/internal/jira"
	llmclient "github.com/zevwings/workflow/internal/llm/client"
	"github.com/zevwings/workflow/internal/pr/github"
	"github.com/zevwings/workflow/internal/prompt"
	"github.com/zevwings/workflow/internal/util"
)

// VerifyJiraConfig 验证 Jira 配置
func VerifyJiraConfig(jiraConfig *config.JiraConfig) {
	if jiraConfig.Email == "" || jiraConfig.APIToken == "" || jiraConfig.ServiceAddress == "" {
		return
	}

	msg := prompt.GetMessage()
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
	msg.Break()
}

// VerifyGitHubConfig 验证 GitHub 配置
func VerifyGitHubConfig(githubConfig *config.GitHubConfig) bool {
	if len(githubConfig.Accounts) == 0 {
		return true
	}

	msg := prompt.GetMessage()
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
	msg.Break()

	return allValid
}

// VerifyLLMConfig 验证 LLM 配置
func VerifyLLMConfig(llmConfig *config.LLMConfig) {
	if llmConfig.Provider == "" {
		return
	}

	msg := prompt.GetMessage()
	msg.Info("LLM Configuration")
	table := prompt.NewTable([]string{"Provider", "Model", "Key", "Output Language"})

	var model, key, language string
	language = llmConfig.Language
	if language == "" {
		language = "en"
	}

	switch llmConfig.Provider {
	case "openai":
		model = llmConfig.OpenAI.Model
		if model == "" {
			model = "gpt-3.5-turbo"
		}
		key = util.MaskSensitiveValue(llmConfig.OpenAI.APIKey)
	case "deepseek":
		model = llmConfig.DeepSeek.Model
		if model == "" {
			model = "deepseek-chat"
		}
		key = util.MaskSensitiveValue(llmConfig.DeepSeek.APIKey)
	case "proxy":
		model = llmConfig.Proxy.Model
		if llmConfig.Proxy.URL != "" {
			model = fmt.Sprintf("%s(%s)", model, llmConfig.Proxy.URL)
		}
		key = util.MaskSensitiveValue(llmConfig.Proxy.APIKey)
	default:
		model = "-"
		key = "-"
	}

	table.AddRow([]string{llmConfig.Provider, model, key, language})
	table.Render()
	msg.Break()
}

// VerifyLogConfig 验证日志配置
func VerifyLogConfig(logConfig *config.LogConfig) {
	if logConfig.Level == "" {
		return
	}
	msg := prompt.GetMessage()
	msg.Info("Log Configuration")
	msg.Info("%s", fmt.Sprintf("Log Level: %s", logConfig.Level))
	msg.Break()
}

// VerifyLLMConnection 验证 LLM 连接
//
// 通过发送一个简单的 "hello" 测试请求来验证 LLM 配置和连接是否正常。
// 验证包括：
//   - 配置完整性检查
//   - API 连接测试
//   - API Key 有效性验证
//   - LLM 响应验证
func VerifyLLMConnection(llmConfig *config.LLMConfig) {
	if llmConfig.Provider == "" {
		return
	}

	msg := prompt.GetMessage()
	msg.Info("LLM Connection Verification")

	// 1. 配置完整性验证
	apiKey, _, _, err := llmConfig.CurrentProvider()
	if err != nil {
		msg.Error("LLM 配置验证失败: %v", err)
		msg.Break()
		return
	}

	if apiKey == "" {
		msg.Error("LLM API Key 未配置")
		msg.Break()
		return
	}

	// 2. 创建 LLM 客户端并发送测试请求
	var response string

	spinner := prompt.NewSpinner("Verifying LLM connection...")
	err = spinner.Do(func() error {
		// 创建配置提供者
		provider := llmadapter.NewLLMConfigProvider()

		// 获取 ProviderConfig
		providerConfig, err := provider.GetProviderConfig()
		if err != nil {
			return fmt.Errorf("获取 LLM provider 配置失败: %w", err)
		}

		// 创建 LLM 客户端
		llmClient := llmclient.Global(providerConfig)

		// 发送简单的测试请求
		params := &llmclient.LLMRequestParams{
			SystemPrompt: "You are a helpful assistant.",
			UserPrompt:   "Say hello",
			Temperature:  0.7,
			// MaxTokens:    intPtr(10), // 限制 token 数，节省成本
		}

		response, err = llmClient.Call(params)
		if err != nil {
			return err
		}

		// 验证响应不为空
		if response == "" {
			return fmt.Errorf("LLM 返回空响应")
		}

		return nil
	})
	// 注意：spinner.Do() 内部已经通过 defer 自动调用 Stop()，无需手动调用

	// 3. 显示验证结果
	if err != nil {
		// 根据错误类型提供不同的提示
		errorMsg := err.Error()
		if strings.Contains(errorMsg, "API key") || strings.Contains(errorMsg, "API Key") {
			msg.Error("LLM API Key 无效或未配置")
		} else if strings.Contains(errorMsg, "timeout") {
			msg.Error("LLM API 连接超时，请检查网络连接")
		} else if strings.Contains(errorMsg, "network") || strings.Contains(errorMsg, "连接") {
			msg.Error("无法连接到 LLM API，请检查网络连接")
		} else {
			msg.Error("LLM 验证失败: %v", err)
		}
		msg.Break()
		return
	}

	// 验证成功
	msg.Success("LLM 验证成功！")
	msg.Info("测试响应: %s", response)
	msg.Break()
}

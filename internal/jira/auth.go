package jira

import (
	"github.com/zevwings/workflow/internal/logging"
)

// AuthResult Jira 认证验证结果
type AuthResult struct {
	// Valid 验证是否通过
	Valid bool
	// Message 验证消息（成功或失败的原因）
	Message string
	// Error 错误信息（如果有）
	Error error
	// Details 额外信息（如用户名、邮箱等）
	Details map[string]interface{}
}

// ValidateAuth 验证 Jira 认证
//
// 验证 Jira 配置是否有效，通过调用 Jira API 测试认证。
// 使用 GetUserInfo API 验证。
//
// 参数:
//   - config: Jira 配置
//
// 返回:
//   - *AuthResult: 验证结果
//   - error: 如果验证过程出错，返回错误
func ValidateAuth(config *Config) (*AuthResult, error) {
	logger := logging.GetLogger()
	result := &AuthResult{
		Details: make(map[string]interface{}),
	}

	// 1. 检查配置完整性
	if config == nil {
		result.Valid = false
		result.Message = "Jira 配置为空"
		return result, nil
	}

	if config.ServiceAddress == "" {
		result.Valid = false
		result.Message = "Jira Service Address 未配置"
		return result, nil
	}

	if config.Email == "" {
		result.Valid = false
		result.Message = "Jira Email 未配置"
		return result, nil
	}

	if config.APIToken == "" {
		result.Valid = false
		result.Message = "Jira API Token 未配置"
		return result, nil
	}

	// 2. 创建 Jira 客户端
	logger.Debug("Creating Jira client for validation...")
	jiraClient, err := NewJiraClient(config)
	if err != nil {
		logger.WithError(err).Error("Failed to create Jira client for validation")
		result.Valid = false
		result.Message = "创建 Jira 客户端失败"
		result.Error = err
		return result, nil
	}

	// 3. 调用 API 验证认证（使用 GetUserInfo API）
	logger.Debug("Validating Jira authentication...")
	// 注意：GetUserInfo 使用 JiraClient 内部的 context，不需要额外传入
	user, err := jiraClient.GetUserInfo()
	if err != nil {
		logger.WithError(err).Error("Jira authentication validation failed")
		result.Valid = false
		result.Message = "Jira 认证失败"
		result.Error = err
		return result, nil
	}

	// 4. 验证成功，提取用户信息
	result.Valid = true
	result.Message = "Jira 认证成功"
	if user != nil {
		if user.DisplayName != "" {
			result.Details["display_name"] = user.DisplayName
		}
		if user.EmailAddress != "" {
			result.Details["email"] = user.EmailAddress
		}
		if user.AccountID != "" {
			result.Details["account_id"] = user.AccountID
		}
	}

	logger.WithFields(logging.Fields{
		"display_name": result.Details["display_name"],
		"valid":        result.Valid,
	}).Info("Jira authentication validated successfully")

	return result, nil
}

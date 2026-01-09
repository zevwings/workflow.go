package github

import (
	"context"
	"time"

	"github.com/google/go-github/v57/github"
	"github.com/zevwings/workflow/internal/logging"
	"golang.org/x/oauth2"
)

// AuthResult GitHub 认证验证结果
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

// ValidateAuth 验证 GitHub 认证
//
// 验证 GitHub token 是否有效，通过调用 GitHub API 测试 token。
// 使用 GetUser API 验证，不需要 owner/repo 信息。
//
// 参数:
//   - token: GitHub Personal Access Token
//
// 返回:
//   - *AuthResult: 验证结果
//   - error: 如果验证过程出错，返回错误
func ValidateAuth(token string) (*AuthResult, error) {
	logger := logging.GetLogger()
	result := &AuthResult{
		Details: make(map[string]interface{}),
	}

	// 1. 检查配置完整性
	if token == "" {
		result.Valid = false
		result.Message = "GitHub API Token 未配置"
		return result, nil
	}

	// 2. 创建 GitHub 客户端
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)

	// 3. 调用 API 验证 token（使用 GetUser API，获取当前认证用户）
	logger.Debug("Validating GitHub authentication...")
	user, _, err := client.Users.Get(ctx, "")
	if err != nil {
		logger.WithError(err).Error("GitHub authentication validation failed")
		result.Valid = false
		result.Message = "GitHub 认证失败"
		result.Error = err
		return result, nil
	}

	// 4. 验证成功，提取用户信息
	result.Valid = true
	result.Message = "GitHub 认证成功"
	if user.Login != nil {
		result.Details["username"] = *user.Login
	}
	if user.Email != nil {
		result.Details["email"] = *user.Email
	}
	if user.Name != nil {
		result.Details["name"] = *user.Name
	}

	logger.WithFields(logging.Fields{
		"username": result.Details["username"],
		"valid":    result.Valid,
	}).Info("GitHub authentication validated successfully")

	return result, nil
}

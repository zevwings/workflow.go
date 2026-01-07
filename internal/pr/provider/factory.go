package provider

import (
	"fmt"

	"github.com/zevwings/workflow/internal/pr"
	"github.com/zevwings/workflow/internal/pr/github"
)

// NewPlatformProvider 创建平台提供者实例
//
// 根据平台名称自动创建对应的平台提供者。
// 目前支持：
//   - "github": 创建 GitHub 平台提供者
//
// 未来可以扩展支持：
//   - "gitlab": GitLab 平台
//   - "bitbucket": Bitbucket 平台
//
// 参数:
//   - platform: 平台名称（如 "github"）
//   - token: 平台认证 token（如 GitHub Personal Access Token）
//   - owner: 仓库所有者（如 "zevwings"）
//   - repo: 仓库名称（如 "workflow"）
//
// 返回:
//   - PlatformProvider: 平台提供者实例
//   - error: 如果平台不支持或创建失败，返回错误
func NewPlatformProvider(platform, token, owner, repo string) (pr.PlatformProvider, error) {
	switch platform {
	case "github":
		return github.NewGitHub(token, owner, repo)
	default:
		return nil, fmt.Errorf("unsupported platform: %s", platform)
	}
}

// AutoDetectPlatform 自动检测平台
//
// 从当前 Git 仓库的远程 URL 自动检测平台类型。
//
// 返回:
//   - string: 平台名称（如 "github"）
//   - error: 如果检测失败，返回错误
func AutoDetectPlatform() (string, error) {
	// TODO: 实现自动检测逻辑
	// 1. 获取当前仓库的远程 URL
	// 2. 根据 URL 判断平台（github.com -> github, gitlab.com -> gitlab）
	// 3. 返回平台名称

	// 目前默认返回 github
	return "github", nil
}

// NewPlatformProviderAuto 自动检测并创建平台提供者
//
// 自动检测当前仓库的平台类型，并创建对应的平台提供者。
// 注意：此函数需要 token、owner 和 repo 参数，业务逻辑应由调用方（commands 层）完成。
//
// 参数:
//   - token: 平台认证 token（如 GitHub Personal Access Token）
//   - owner: 仓库所有者（如 "zevwings"）
//   - repo: 仓库名称（如 "workflow"）
//
// 返回:
//   - PlatformProvider: 平台提供者实例
//   - error: 如果检测失败或创建失败，返回错误
func NewPlatformProviderAuto(token, owner, repo string) (pr.PlatformProvider, error) {
	platform, err := AutoDetectPlatform()
	if err != nil {
		return nil, fmt.Errorf("failed to detect platform: %w", err)
	}

	return NewPlatformProvider(platform, token, owner, repo)
}


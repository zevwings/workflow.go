package helpers

import (
	"fmt"
	"strings"
)

// ExtractRepoFromURL 从 URL 中提取仓库信息
//
// 支持多种 URL 格式：
//   - https://github.com/owner/repo
//   - https://github.com/owner/repo.git
//   - git@github.com:owner/repo.git
//   - ssh://git@github.com/owner/repo.git
//
// 返回:
//   - owner: 仓库所有者
//   - repo: 仓库名称
//   - error: 如果解析失败，返回错误
func ExtractRepoFromURL(url string) (owner, repo string, err error) {
	// 移除 .git 后缀
	url = strings.TrimSuffix(url, ".git")

	// SSH 格式: git@host:owner/repo
	if strings.Contains(url, "@") && strings.Contains(url, ":") {
		parts := strings.Split(url, ":")
		if len(parts) == 2 {
			repoPath := parts[1]
			repoParts := strings.Split(repoPath, "/")
			if len(repoParts) == 2 {
				return repoParts[0], repoParts[1], nil
			}
		}
	}

	// HTTPS 格式: https://host/owner/repo
	if strings.Contains(url, "://") {
		parts := strings.Split(url, "://")
		if len(parts) == 2 {
			path := strings.TrimPrefix(parts[1], "www.")
			// 移除 host 部分
			pathParts := strings.SplitN(path, "/", 2)
			if len(pathParts) == 2 {
				repoParts := strings.Split(pathParts[1], "/")
				if len(repoParts) >= 2 {
					return repoParts[0], repoParts[1], nil
				}
			}
		}
	}

	return "", "", fmt.Errorf("failed to extract repo from URL: %s", url)
}

// BuildPRURL 构建 PR URL
//
// 参数:
//   - baseURL: 基础 URL（如 https://github.com）
//   - owner: 仓库所有者
//   - repo: 仓库名称
//   - prNumber: PR 编号
//
// 返回:
//   - string: PR URL
func BuildPRURL(baseURL, owner, repo string, prNumber int) string {
	baseURL = strings.TrimSuffix(baseURL, "/")
	return fmt.Sprintf("%s/%s/%s/pull/%d", baseURL, owner, repo, prNumber)
}


package git

import (
	"fmt"
	"strings"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/transport"
)

// AddRemote 添加远程仓库
func (r *Repository) AddRemote(name, url string) error {
	_, err := r.repo.CreateRemote(&config.RemoteConfig{
		Name: name,
		URLs: []string{url},
	})
	if err != nil {
		return fmt.Errorf("failed to add remote %s: %w", name, err)
	}
	return nil
}

// RemoveRemote 删除远程仓库
func (r *Repository) RemoveRemote(name string) error {
	err := r.repo.DeleteRemote(name)
	if err != nil {
		return fmt.Errorf("failed to remove remote %s: %w", name, err)
	}
	return nil
}

// ListRemotes 列出所有远程仓库
func (r *Repository) ListRemotes() ([]RemoteInfo, error) {
	remotes, err := r.repo.Remotes()
	if err != nil {
		return nil, fmt.Errorf("failed to list remotes: %w", err)
	}

	infos := []RemoteInfo{}
	for _, remote := range remotes {
		url := ""
		if len(remote.Config().URLs) > 0 {
			url = remote.Config().URLs[0]
		}
		infos = append(infos, RemoteInfo{
			Name: remote.Config().Name,
			URL:  url,
		})
	}

	return infos, nil
}

// GetRemoteURL 获取远程仓库 URL
func (r *Repository) GetRemoteURL(name string) (string, error) {
	remote, err := r.repo.Remote(name)
	if err != nil {
		return "", fmt.Errorf("failed to get remote %s: %w", name, err)
	}

	urls := remote.Config().URLs
	if len(urls) == 0 {
		return "", fmt.Errorf("remote %s has no URL", name)
	}

	return urls[0], nil
}

// Fetch 从远程获取更新
func (r *Repository) Fetch(remoteName string, auth transport.AuthMethod) error {
	remote, err := r.repo.Remote(remoteName)
	if err != nil {
		return fmt.Errorf("failed to get remote %s: %w", remoteName, err)
	}

	err = remote.Fetch(&git.FetchOptions{
		Auth: auth,
	})
	if err != nil && err != git.NoErrAlreadyUpToDate {
		return fmt.Errorf("failed to fetch from %s: %w", remoteName, err)
	}

	return nil
}

// Push 推送到远程
func (r *Repository) Push(remoteName string, branchName string, auth transport.AuthMethod) error {
	remote, err := r.repo.Remote(remoteName)
	if err != nil {
		return fmt.Errorf("failed to get remote %s: %w", remoteName, err)
	}

	refSpec := fmt.Sprintf("refs/heads/%s:refs/heads/%s", branchName, branchName)
	err = remote.Push(&git.PushOptions{
		RefSpecs: []config.RefSpec{config.RefSpec(refSpec)},
		Auth:     auth,
	})
	if err != nil && err != git.NoErrAlreadyUpToDate {
		return fmt.Errorf("failed to push to %s: %w", remoteName, err)
	}

	return nil
}

// PushWithUpstream 推送并设置上游分支
// 注意：go-git v5 不直接支持设置上游分支，此方法只执行推送
// 如果需要设置上游，可以使用 git 命令：git branch --set-upstream-to=origin/branch branch
func (r *Repository) PushWithUpstream(remoteName string, branchName string, auth transport.AuthMethod) error {
	// 推送分支
	return r.Push(remoteName, branchName, auth)
}

// ListRemoteRefs 列出远程引用
func (r *Repository) ListRemoteRefs(remoteName string) (map[string]plumbing.Hash, error) {
	remote, err := r.repo.Remote(remoteName)
	if err != nil {
		return nil, fmt.Errorf("failed to get remote %s: %w", remoteName, err)
	}

	refs, err := remote.List(&git.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to list remote refs: %w", err)
	}

	result := make(map[string]plumbing.Hash)
	for _, ref := range refs {
		result[ref.Name().String()] = ref.Hash()
	}

	return result, nil
}

// ExtractRepoName 从远程 URL 提取仓库名（owner/repo 格式）
func ExtractRepoName(url string) (string, error) {
	// 支持多种 URL 格式
	// git@github.com:owner/repo.git
	// https://github.com/owner/repo.git
	// ssh://git@github.com/owner/repo.git

	// 移除 .git 后缀
	url = strings.TrimSuffix(url, ".git")

	// SSH 格式: git@host:owner/repo
	if strings.Contains(url, "@") && strings.Contains(url, ":") {
		parts := strings.Split(url, ":")
		if len(parts) == 2 {
			return parts[1], nil
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
				return pathParts[1], nil
			}
		}
	}

	return "", fmt.Errorf("failed to extract repo name from URL: %s", url)
}


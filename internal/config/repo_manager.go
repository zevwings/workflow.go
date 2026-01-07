package config

import (
	"crypto/sha256"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/pelletier/go-toml/v2"
	"github.com/spf13/viper"
)

// RepoManager 仓库配置管理器
//
// 管理仓库级别的配置：
//   - 项目公共配置：.workflow/config.toml（项目根目录，提交到 Git）
//   - 项目私有配置：~/.workflow/config/repository.toml（用户主目录，不提交）
type RepoManager struct {
	// 项目公共配置
	publicViper *viper.Viper
	publicPath  string

	// 项目私有配置
	privatePath string
	repoID      string
	repoPath    string
}

// NewRepoManager 创建仓库配置管理器
//
// 参数:
//   - repoPath: 仓库根目录路径（如果为空，使用当前目录）
//
// 返回:
//   - *RepoManager: 仓库配置管理器实例
//   - error: 如果创建失败，返回错误
func NewRepoManager(repoPath string) (*RepoManager, error) {
	if repoPath == "" {
		var err error
		repoPath, err = os.Getwd()
		if err != nil {
			return nil, fmt.Errorf("获取当前目录失败: %w", err)
		}
	}

	// 检查是否为 Git 仓库
	if !isGitRepo(repoPath) {
		return nil, fmt.Errorf("路径不是 Git 仓库: %s", repoPath)
	}

	// 生成 repo_id
	repoID, err := generateRepoID(repoPath)
	if err != nil {
		return nil, fmt.Errorf("生成仓库 ID 失败: %w", err)
	}

	// 初始化项目公共配置 viper
	publicViper := viper.New()
	publicViper.SetConfigName("config")
	publicViper.SetConfigType("toml")
	publicPath := filepath.Join(repoPath, ".workflow", "config.toml")
	publicViper.AddConfigPath(filepath.Join(repoPath, ".workflow"))
	publicViper.AddConfigPath(repoPath)

	// 获取用户主目录
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("获取用户主目录失败: %w", err)
	}

	// 设置私有配置路径
	privateConfigDir := filepath.Join(homeDir, ".workflow", "config")
	privatePath := filepath.Join(privateConfigDir, "repository.toml")

	return &RepoManager{
		publicViper: publicViper,
		publicPath:  publicPath,
		privatePath: privatePath,
		repoID:      repoID,
		repoPath:    repoPath,
	}, nil
}

// Load 加载仓库配置
//
// 加载顺序：
// 1. 先加载项目公共配置（如果存在）
// 2. 再加载项目私有配置（如果存在）
// 3. 私有配置覆盖公共配置
func (r *RepoManager) Load() error {
	// 加载项目公共配置（如果存在）
	if _, err := os.Stat(r.publicPath); err == nil {
		if err := r.publicViper.ReadInConfig(); err != nil {
			return fmt.Errorf("读取项目公共配置失败: %w", err)
		}
	}

	// 加载项目私有配置（如果存在）
	// 注意：私有配置的加载逻辑会在后续方法中实现
	// 这里先不实现，因为需要解析 TOML 并合并到 publicViper

	return nil
}

// GetTemplateConfig 获取模板配置
//
// 从项目公共配置中读取模板配置。
//
// 返回:
//   - *TemplateConfig: 模板配置结构
func (r *RepoManager) GetTemplateConfig() *TemplateConfig {
	cfg := &TemplateConfig{
		Commit:       make(map[string]interface{}),
		Branch:       make(map[string]interface{}),
		PullRequests: make(map[string]interface{}),
	}

	// 读取 template.commit
	if commitMap := r.publicViper.GetStringMap("template.commit"); commitMap != nil {
		for k, v := range commitMap {
			cfg.Commit[k] = v
		}
	}

	// 读取 template.branch
	if branchMap := r.publicViper.GetStringMap("template.branch"); branchMap != nil {
		for k, v := range branchMap {
			cfg.Branch[k] = v
		}
	}

	// 读取 template.pull_requests
	if prMap := r.publicViper.GetStringMap("template.pull_requests"); prMap != nil {
		for k, v := range prMap {
			cfg.PullRequests[k] = v
		}
	}

	return cfg
}

// GetBranchPrefix 获取分支前缀（个人偏好）
//
// 从项目私有配置中读取分支前缀。
//
// 返回:
//   - string: 分支前缀，如果未配置则返回空字符串
func (r *RepoManager) GetBranchPrefix() string {
	// TODO: 实现从私有配置读取
	// 需要加载 ~/.workflow/config/repository.toml 并解析 [${repo_id}.branch.prefix]
	return ""
}

// GetIgnoreBranches 获取忽略的分支列表（个人偏好）
//
// 从项目私有配置中读取忽略的分支列表。
//
// 返回:
//   - []string: 忽略的分支列表
func (r *RepoManager) GetIgnoreBranches() []string {
	// TODO: 实现从私有配置读取
	return []string{}
}

// GetAutoAcceptChangeType 获取自动接受变更类型（个人偏好）
//
// 从项目私有配置中读取自动接受变更类型设置。
//
// 返回:
//   - bool: 是否自动接受变更类型，默认返回 false
func (r *RepoManager) GetAutoAcceptChangeType() bool {
	// TODO: 实现从私有配置读取
	return false
}

// SaveTemplateConfig 保存模板配置
//
// 保存模板配置到项目公共配置文件。
//
// 参数:
//   - cfg: 模板配置
//
// 返回:
//   - error: 如果保存失败，返回错误
func (r *RepoManager) SaveTemplateConfig(cfg *TemplateConfig) error {
	// 确保目录存在
	if err := os.MkdirAll(filepath.Dir(r.publicPath), 0755); err != nil {
		return fmt.Errorf("创建配置目录失败: %w", err)
	}

	// 读取现有配置（如果存在）
	var existingConfig map[string]interface{}
	if _, err := os.Stat(r.publicPath); err == nil {
		data, err := os.ReadFile(r.publicPath)
		if err == nil {
			if err := toml.Unmarshal(data, &existingConfig); err != nil {
				existingConfig = make(map[string]interface{})
			}
		}
	}
	if existingConfig == nil {
		existingConfig = make(map[string]interface{})
	}

	// 更新 template 部分
	templateSection := make(map[string]interface{})
	if existingTemplate, ok := existingConfig["template"].(map[string]interface{}); ok {
		templateSection = existingTemplate
	}

	if len(cfg.Commit) > 0 {
		templateSection["commit"] = cfg.Commit
	}
	if len(cfg.Branch) > 0 {
		templateSection["branch"] = cfg.Branch
	}
	if len(cfg.PullRequests) > 0 {
		templateSection["pull_requests"] = cfg.PullRequests
	}

	existingConfig["template"] = templateSection

	// 保存配置
	data, err := toml.Marshal(existingConfig)
	if err != nil {
		return fmt.Errorf("序列化配置失败: %w", err)
	}

	if err := os.WriteFile(r.publicPath, data, 0644); err != nil {
		return fmt.Errorf("写入配置文件失败: %w", err)
	}

	return nil
}

// GetPublicConfigPath 获取项目公共配置文件路径
func (r *RepoManager) GetPublicConfigPath() string {
	return r.publicPath
}

// GetPrivateConfigPath 获取项目私有配置文件路径
func (r *RepoManager) GetPrivateConfigPath() string {
	return r.privatePath
}

// GetRepoID 获取仓库 ID
func (r *RepoManager) GetRepoID() string {
	return r.repoID
}

// isGitRepo 检查路径是否为 Git 仓库
func isGitRepo(path string) bool {
	gitDir := filepath.Join(path, ".git")
	if info, err := os.Stat(gitDir); err == nil && info.IsDir() {
		return true
	}

	// 也检查 git rev-parse 命令
	cmd := exec.Command("git", "rev-parse", "--git-dir")
	cmd.Dir = path
	if err := cmd.Run(); err == nil {
		return true
	}

	return false
}

// generateRepoID 生成仓库 ID
//
// 基于 Git remote URL 生成唯一的仓库标识符。
// 格式：{repo_name}_{hash}，其中 hash 是 URL 的 SHA256 前 8 个字符。
func generateRepoID(repoPath string) (string, error) {
	// 获取 remote URL
	cmd := exec.Command("git", "config", "--get", "remote.origin.url")
	cmd.Dir = repoPath
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("获取 remote URL 失败: %w", err)
	}
	url := strings.TrimSpace(string(output))
	if url == "" {
		return "", fmt.Errorf("未找到 remote.origin.url")
	}

	// 提取仓库名称（从 URL 中）
	repoName := extractRepoNameFromURL(url)

	// 计算 SHA256 hash
	hash := sha256.Sum256([]byte(url))
	hashStr := fmt.Sprintf("%x", hash)

	// 取前 8 个字符
	return fmt.Sprintf("%s_%s", repoName, hashStr[:8]), nil
}

// extractRepoNameFromURL 从 URL 中提取仓库名称
func extractRepoNameFromURL(url string) string {
	// 移除 .git 后缀
	url = strings.TrimSuffix(url, ".git")

	// 处理不同的 URL 格式
	// git@github.com:owner/repo.git
	// https://github.com/owner/repo.git
	// https://github.com/owner/repo

	parts := strings.Split(url, "/")
	if len(parts) > 0 {
		name := parts[len(parts)-1]
		// 移除可能的协议前缀
		if strings.Contains(name, ":") {
			name = strings.Split(name, ":")[1]
		}
		return name
	}

	return "unknown"
}

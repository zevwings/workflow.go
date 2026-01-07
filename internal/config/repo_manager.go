package config

import (
	"crypto/sha256"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/pelletier/go-toml/v2"
	"github.com/spf13/viper"
)

// GitRepository 定义 Git 仓库操作接口
//
// 用于解耦 config 模块对 git 模块的依赖，支持依赖注入。
// 接口定义在 config 包中，避免循环依赖。
type GitRepository interface {
	// GetRepoPath 获取仓库路径
	//
	// 返回此 GitRepository 实例管理的仓库路径。
	// 如果路径为空，表示使用当前工作目录。
	GetRepoPath() string
	// IsGitRepo 检查指定路径是否是 Git 仓库
	IsGitRepo(path string) bool
	// Open 打开指定路径的 Git 仓库
	Open(path string) (GitRepo, error)
}

// GitRepo 定义 Git 仓库实例接口
//
// 表示一个已打开的 Git 仓库，提供获取 remote URL 的方法。
type GitRepo interface {
	// GetRemoteURL 获取指定名称的 remote URL
	GetRemoteURL(name string) (string, error)
}

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

	// 缓存的私有配置（延迟加载）
	privateConfig *PrivateRepoConfig
}

// NewRepoManager 创建仓库配置管理器
//
// 通过依赖注入解耦 config 模块对 git 模块的依赖。
// 仓库路径从 GitRepository 接口中获取。
//
// 参数:
//   - gitRepo: Git 仓库操作接口（如果为 nil，则跳过 Git 相关操作，使用当前目录和简单 ID）
//
// 返回:
//   - *RepoManager: 仓库配置管理器实例
//   - error: 如果创建失败，返回错误
func NewRepoManager(gitRepo GitRepository) (*RepoManager, error) {
	var repoPath string
	var repoID string
	var err error

	// 如果提供了 gitRepo，从接口获取路径
	if gitRepo != nil {
		repoPath = gitRepo.GetRepoPath()
		if repoPath == "" {
			// 如果路径为空，使用当前目录
			var err error
			repoPath, err = os.Getwd()
			if err != nil {
				return nil, fmt.Errorf("获取当前目录失败: %w", err)
			}
		}

		// 检查是否为 Git 仓库
		if !gitRepo.IsGitRepo(repoPath) {
			return nil, fmt.Errorf("路径不是 Git 仓库: %s", repoPath)
		}

		// 生成 repo_id
		repoID, err = generateRepoIDWithGit(repoPath, gitRepo)
		if err != nil {
			return nil, fmt.Errorf("生成仓库 ID 失败: %w", err)
		}
	} else {
		// 如果没有提供 gitRepo，使用当前目录和基于路径的简单 ID
		var err error
		repoPath, err = os.Getwd()
		if err != nil {
			return nil, fmt.Errorf("获取当前目录失败: %w", err)
		}
		repoID = generateSimpleRepoID(repoPath)
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
	cfg := r.loadPrivateConfig()
	if cfg == nil {
		return ""
	}

	// 查找当前 repo_id 的配置
	repoSection, ok := cfg.Repositories[r.repoID]
	if !ok {
		return ""
	}

	if repoSection.Branch != nil && repoSection.Branch.Prefix != nil {
		return *repoSection.Branch.Prefix
	}

	return ""
}

// GetIgnoreBranches 获取忽略的分支列表（个人偏好）
//
// 从项目私有配置中读取忽略的分支列表。
//
// 返回:
//   - []string: 忽略的分支列表
func (r *RepoManager) GetIgnoreBranches() []string {
	cfg := r.loadPrivateConfig()
	if cfg == nil {
		return []string{}
	}

	// 查找当前 repo_id 的配置
	repoSection, ok := cfg.Repositories[r.repoID]
	if !ok {
		return []string{}
	}

	if repoSection.Branch != nil && len(repoSection.Branch.Ignore) > 0 {
		return repoSection.Branch.Ignore
	}

	return []string{}
}

// GetAutoAcceptChangeType 获取自动接受变更类型（个人偏好）
//
// 从项目私有配置中读取自动接受变更类型设置。
//
// 返回:
//   - bool: 是否自动接受变更类型，默认返回 false
func (r *RepoManager) GetAutoAcceptChangeType() bool {
	cfg := r.loadPrivateConfig()
	if cfg == nil {
		return false
	}

	// 查找当前 repo_id 的配置
	repoSection, ok := cfg.Repositories[r.repoID]
	if !ok {
		return false
	}

	if repoSection.AutoAcceptChangeType != nil {
		return *repoSection.AutoAcceptChangeType
	}

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

	// 使用辅助函数保存配置
	return SaveConfigToFile(r.publicPath, existingConfig)
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

// generateRepoIDWithGit 使用 Git 接口生成仓库 ID
//
// 基于 Git remote URL 生成唯一的仓库标识符。
// 格式：{repo_name}_{hash}，其中 hash 是 URL 的 SHA256 前 8 个字符。
func generateRepoIDWithGit(repoPath string, gitRepo GitRepository) (string, error) {
	// 打开 Git 仓库
	repo, err := gitRepo.Open(repoPath)
	if err != nil {
		return "", fmt.Errorf("打开 Git 仓库失败: %w", err)
	}

	// 获取 remote URL
	url, err := repo.GetRemoteURL("origin")
	if err != nil {
		return "", fmt.Errorf("获取 remote URL 失败: %w", err)
	}
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

// generateSimpleRepoID 基于路径生成简单的仓库 ID（不依赖 Git）
//
// 当没有提供 Git 接口时，使用此方法生成基于路径的简单 ID。
func generateSimpleRepoID(repoPath string) string {
	// 使用路径的绝对路径生成 hash
	absPath, err := filepath.Abs(repoPath)
	if err != nil {
		absPath = repoPath
	}
	hash := sha256.Sum256([]byte(absPath))
	hashStr := fmt.Sprintf("%x", hash)
	return fmt.Sprintf("repo_%s", hashStr[:8])
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

// PrivateRepoConfig 私有仓库配置结构
//
// 用于解析 ~/.workflow/config/repository.toml 文件。
// 格式：
//
//	[${repo_id}.branch]
//	prefix = "..."
//	ignore = ["branch1", "branch2"]
//
//	[${repo_id}]
//	auto_accept_change_type = true
type PrivateRepoConfig struct {
	// Repositories 按 repo_id 组织的配置
	Repositories map[string]PrivateRepoSection `toml:",inline"`
}

// PrivateRepoSection 单个仓库的私有配置
type PrivateRepoSection struct {
	// Branch 分支相关配置
	Branch *BranchConfig `toml:"branch,omitempty"`
	// AutoAcceptChangeType 自动接受变更类型
	AutoAcceptChangeType *bool `toml:"auto_accept_change_type,omitempty"`
}

// loadPrivateConfig 加载私有配置（延迟加载，带缓存）
//
// 如果配置文件不存在或解析失败，返回 nil（不返回错误，因为私有配置是可选的）。
//
// 返回:
//   - *PrivateRepoConfig: 私有配置，如果不存在或解析失败返回 nil
func (r *RepoManager) loadPrivateConfig() *PrivateRepoConfig {
	// 如果已缓存，直接返回
	if r.privateConfig != nil {
		return r.privateConfig
	}

	// 检查配置文件是否存在
	if _, err := os.Stat(r.privatePath); os.IsNotExist(err) {
		// 配置文件不存在，返回 nil（不缓存，下次可能创建）
		return nil
	}

	// 读取配置文件
	data, err := os.ReadFile(r.privatePath)
	if err != nil {
		// 读取失败，返回 nil（不缓存）
		return nil
	}

	// 解析 TOML
	var config map[string]interface{}
	if err := toml.Unmarshal(data, &config); err != nil {
		// 解析失败，返回 nil（不缓存）
		return nil
	}

	// 转换为结构化配置
	privateConfig := &PrivateRepoConfig{
		Repositories: make(map[string]PrivateRepoSection),
	}

	// 初始化仓库配置段
	repoSection := PrivateRepoSection{}

	// TOML 解析器会将 [repo_id.branch] 解析为嵌套结构：
	// config["repo_id"] = map[string]interface{}{
	//     "branch": map[string]interface{}{...}
	// }
	// 所以需要先检查 repo_id 是否存在
	if repoValue, ok := config[r.repoID]; ok {
		if repoMap, ok := repoValue.(map[string]interface{}); ok {
			// 解析 branch 配置
			if branchValue, ok := repoMap["branch"]; ok {
				if branchMap, ok := branchValue.(map[string]interface{}); ok {
					branchConfig := &BranchConfig{}
					if prefix, ok := branchMap["prefix"].(string); ok {
						branchConfig.Prefix = &prefix
					}
					if ignore, ok := branchMap["ignore"].([]interface{}); ok {
						branchConfig.Ignore = make([]string, 0, len(ignore))
						for _, item := range ignore {
							if str, ok := item.(string); ok {
								branchConfig.Ignore = append(branchConfig.Ignore, str)
							}
						}
					}
					repoSection.Branch = branchConfig
				}
			}

			// 解析 auto_accept_change_type
			if autoAccept, ok := repoMap["auto_accept_change_type"].(bool); ok {
				repoSection.AutoAcceptChangeType = &autoAccept
			}
		}
	}

	// 如果有任何配置，保存到结果中
	if repoSection.Branch != nil || repoSection.AutoAcceptChangeType != nil {
		privateConfig.Repositories[r.repoID] = repoSection
	}

	// 缓存配置
	r.privateConfig = privateConfig

	return privateConfig
}

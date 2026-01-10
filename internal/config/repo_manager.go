package config

import (
	"crypto/sha256"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/pelletier/go-toml/v2"
	"github.com/spf13/viper"
	"github.com/zevwings/workflow/internal/logging"
)

var (
	// repoManager 仓库配置管理器单例
	repoManager *RepoManager
	repoOnce    sync.Once
	repoErr     error
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
//   - 项目私有配置：$XDG_CONFIG_HOME/Workflow/config/repository.toml（遵循 XDG 规范，不提交）
//
// 配置字段可以直接访问，例如：
//   - manager.TemplateConfig.Commit
//   - manager.Config.Template.Branch
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

	// Config 仓库公共配置数据
	// 在 Load() 时自动加载，可以直接访问配置字段
	Config *RepoConfig

	// 便捷字段：直接访问子配置（指向 Config 中的对应字段）
	TemplateConfig *TemplateConfig // 指向 Config.Template
}

// newRepoManager 创建仓库配置管理器（私有函数）
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
func newRepoManager(gitRepo GitRepository) (*RepoManager, error) {
	logger := logging.GetLogger()
	logger.Debug("Creating repository config manager")

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
				logger.WithError(err).Error("Failed to get current directory")
				return nil, fmt.Errorf("获取当前目录失败: %w", err)
			}
		}

		// 检查是否为 Git 仓库
		if !gitRepo.IsGitRepo(repoPath) {
			logger.WithField("path", repoPath).Error("Path is not a Git repository")
			return nil, fmt.Errorf("路径不是 Git 仓库: %s", repoPath)
		}

		// 生成 repo_id
		repoID, err = generateRepoIDWithGit(repoPath, gitRepo)
		if err != nil {
			logger.WithError(err).WithField("path", repoPath).Error("Failed to generate repository ID")
			return nil, fmt.Errorf("生成仓库 ID 失败: %w", err)
		}
	} else {
		// 如果没有提供 gitRepo，使用当前目录和基于路径的简单 ID
		var err error
		repoPath, err = os.Getwd()
		if err != nil {
			logger.WithError(err).Error("Failed to get current directory")
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

	// 使用 XDG 配置目录设置私有配置路径
	workflowConfigDir, err := ConfigDir()
	if err != nil {
		return nil, fmt.Errorf("获取配置目录失败: %w", err)
	}

	privateConfigDir := filepath.Join(workflowConfigDir, "config")
	privatePath := filepath.Join(privateConfigDir, "repository.toml")

	config := &RepoConfig{
		Template: TemplateConfig{
			Commit:       make(map[string]interface{}),
			Branch:       make(map[string]interface{}),
			PullRequests: make(map[string]interface{}),
		},
	}

	manager := &RepoManager{
		publicViper: publicViper,
		publicPath:  publicPath,
		privatePath: privatePath,
		repoID:      repoID,
		repoPath:    repoPath,
		Config:      config,
	}

	// 初始化便捷字段，指向 Config 中的对应字段
	manager.TemplateConfig = &config.Template

	return manager, nil
}

// Global 获取全局 RepoManager 单例
//
// 返回进程级别的 RepoManager 单例。
// 单例会在首次调用时初始化，后续调用会复用同一个实例。
//
// 参数:
//   - gitRepo: Git 仓库操作接口（如果为 nil，则跳过 Git 相关操作，使用当前目录和简单 ID）
//     首次调用时传入的参数会被保存，后续调用会忽略参数
//
// 返回:
//   - *RepoManager: 仓库配置管理器实例
//   - error: 如果创建失败，返回错误
//
// 注意:
//   - 首次调用时如果创建失败，后续调用会返回相同的错误
//   - 首次调用时传入的参数会被保存，后续调用会忽略参数
//   - 线程安全：可以在多线程环境中安全使用
//
// 优势:
//   - 减少资源消耗：避免重复创建管理器实例
//   - 统一管理：所有配置操作使用同一个管理器实例
//   - 配置一致性：确保整个进程使用相同的配置状态
func GlobalRepoManager(gitRepo GitRepository) (*RepoManager, error) {
	repoOnce.Do(func() {
		repoManager, repoErr = newRepoManager(gitRepo)
	})
	return repoManager, repoErr
}

// Load 加载仓库公共配置
//
// 加载项目公共配置（.workflow/config.toml）到内存，并更新 Config 字段。
// 此方法只加载公共配置，私有配置通过 LoadPrivateConfig() 或按需延迟加载。
//
// 注意：
// - 公共配置和私有配置是分离的，不会合并
// - 公共配置：项目标准，可提交到 Git
// - 私有配置：个人偏好，不提交到 Git，通过 GetBranchPrefix() 等方法按需加载
func (r *RepoManager) Load() error {
	return r.LoadPublicConfig()
}

// LoadPublicConfig 加载仓库公共配置
//
// 加载项目公共配置（.workflow/config.toml）到内存，并更新 Config 字段。
// 公共配置包含模板配置（提交、分支、PR 模板），可以提交到 Git 与团队共享。
func (r *RepoManager) LoadPublicConfig() error {
	logger := logging.GetLogger()

	// 加载项目公共配置（如果存在）
	if _, err := os.Stat(r.publicPath); err == nil {
		logger.Debugf("Loading public config from: %s", r.publicPath)
		if err := r.publicViper.ReadInConfig(); err != nil {
			logger.WithError(err).WithField("config_path", r.publicPath).Error("Failed to load config file")
			return fmt.Errorf("读取项目公共配置失败: %w", err)
		}
	}

	// 从 publicViper 加载配置到 Config 字段
	r.Config = r.getRepoConfig()

	// 更新便捷字段的指针（确保指向最新的 Config）
	r.TemplateConfig = &r.Config.Template

	return nil
}

// GetTemplateConfig 获取模板配置
//
// 返回 Config 字段中的模板配置的引用。
// 可以直接使用 manager.Config.Template 或 manager.TemplateConfig 访问，此方法保留以保持向后兼容。
//
// 返回:
//   - *TemplateConfig: 模板配置结构（指向 Config.Template）
func (r *RepoManager) GetTemplateConfig() *TemplateConfig {
	if r.Config == nil {
		return &TemplateConfig{
			Commit:       make(map[string]interface{}),
			Branch:       make(map[string]interface{}),
			PullRequests: make(map[string]interface{}),
		}
	}
	return r.TemplateConfig
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

// Save 保存配置到文件
//
// 保存当前 Config 字段的内容到文件。
// 保存后会自动重新加载以同步 publicViper。
func (r *RepoManager) Save() error {
	logger := logging.GetLogger()
	logger.Infof("Saving config to: %s", r.publicPath)

	err := SaveConfigToFile(r.publicPath, r.Config)
	if err != nil {
		logger.WithError(err).WithField("config_path", r.publicPath).Error("Failed to save config file")
		return err
	}

	// 重新加载以同步 publicViper 和 Config 字段
	if err := r.Load(); err != nil {
		return err
	}

	return nil
}

// SaveTemplateConfig 保存模板配置（已废弃，请使用 Save()）
//
// 保存模板配置到项目公共配置文件。
// 此方法保留以保持向后兼容，但建议直接设置 manager.Config.Template 或 manager.TemplateConfig 后调用 Save()。
//
// 参数:
//   - cfg: 模板配置
//
// 返回:
//   - error: 如果保存失败，返回错误
//
// 已废弃: 请使用 manager.TemplateConfig = cfg; manager.Save() 或直接操作 manager.Config.Template
func (r *RepoManager) SaveTemplateConfig(cfg *TemplateConfig) error {
	// 更新 Config 中的模板配置
	if r.Config == nil {
		r.Config = &RepoConfig{
			Template: TemplateConfig{
				Commit:       make(map[string]interface{}),
				Branch:       make(map[string]interface{}),
				PullRequests: make(map[string]interface{}),
			},
		}
		r.TemplateConfig = &r.Config.Template
	}
	r.Config.Template = *cfg
	r.TemplateConfig = &r.Config.Template
	return r.Save()
}

// getRepoConfig 获取完整的仓库配置（私有方法）
//
// 从 publicViper 中读取完整的仓库公共配置。
// 此方法仅在内部使用（Load 方法中），外部应直接访问 Config 字段。
//
// 返回:
//   - *RepoConfig: 仓库配置结构
func (r *RepoManager) getRepoConfig() *RepoConfig {
	cfg := &RepoConfig{
		Template: TemplateConfig{
			Commit:       make(map[string]interface{}),
			Branch:       make(map[string]interface{}),
			PullRequests: make(map[string]interface{}),
		},
	}

	// 读取 template.commit
	if commitMap := r.publicViper.GetStringMap("template.commit"); commitMap != nil {
		for k, v := range commitMap {
			cfg.Template.Commit[k] = v
		}
	}

	// 读取 template.branch
	if branchMap := r.publicViper.GetStringMap("template.branch"); branchMap != nil {
		for k, v := range branchMap {
			cfg.Template.Branch[k] = v
		}
	}

	// 读取 template.pull_requests
	if prMap := r.publicViper.GetStringMap("template.pull_requests"); prMap != nil {
		for k, v := range prMap {
			cfg.Template.PullRequests[k] = v
		}
	}

	return cfg
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

// LoadPrivateConfig 加载私有配置
//
// 加载项目私有配置（~/.workflow/config/repository.toml）到内存。
// 私有配置包含个人偏好（分支前缀、自动接受变更类型等），不提交到 Git。
//
// 此方法使用延迟加载和缓存机制：
// - 首次调用时从文件加载并缓存
// - 后续调用直接返回缓存结果
// - 如果配置文件不存在或解析失败，返回 nil（不返回错误，因为私有配置是可选的）
//
// 返回:
//   - *PrivateRepoConfig: 私有配置，如果不存在或解析失败返回 nil
func (r *RepoManager) LoadPrivateConfig() *PrivateRepoConfig {
	return r.loadPrivateConfig()
}

// SavePrivateConfig 保存私有配置
//
// 保存私有配置到文件。
func (r *RepoManager) SavePrivateConfig(cfg *PrivateRepoConfig) error {
	// 更新缓存
	r.privateConfig = cfg

	// 保存到文件
	return SaveConfigToFile(r.privatePath, cfg)
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
// 用于解析 $XDG_CONFIG_HOME/Workflow/config/repository.toml 文件。
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

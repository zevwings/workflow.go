package config

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/spf13/viper"
	"github.com/zevwings/workflow/internal/logging"
)

var (
	// globalManager 全局配置管理器单例
	globalManager *GlobalManager
	globalOnce    sync.Once
	globalErr     error
)

// GlobalManager 全局配置管理器
//
// 管理用户级别的全局配置：遵循 XDG Base Directory Specification
// 配置文件位置：$XDG_CONFIG_HOME/workflow/config.toml（默认：~/.config/workflow/config.toml）
// 包含：用户信息、认证配置（GitHub、Jira）、工具配置（LLM、Proxy、Log）
//
// 配置字段可以直接访问，例如：
//   - manager.LLMConfig.Provider
//   - manager.GitHubConfig.Current
//   - manager.Config.Log.Level
type GlobalManager struct {
	viper *viper.Viper
	path  string

	// Config 全局配置数据
	// 在 Load() 时自动加载，可以直接访问配置字段
	Config *GlobalConfig

	// 便捷字段：直接访问子配置（指向 Config 中的对应字段）
	LLMConfig    *LLMConfig    // 指向 Config.LLM
	GitHubConfig *GitHubConfig // 指向 Config.GitHub
	UserConfig   *UserConfig   // 指向 Config.User
	JiraConfig   *JiraConfig   // 指向 Config.Jira
	LogConfig    *LogConfig    // 指向 Config.Log
	ProxyConfig  *ProxyConfig  // 指向 Config.Proxy
}

// newGlobalManager 创建全局配置管理器（私有函数）
//
// 返回:
//   - *GlobalManager: 全局配置管理器实例
//   - error: 如果创建失败，返回错误
func newGlobalManager() (*GlobalManager, error) {
	// 使用 XDG 配置目录
	configDir, err := ConfigDir()
	if err != nil {
		return nil, fmt.Errorf("获取配置目录失败: %w", err)
	}

	configPath := filepath.Join(configDir, "config.toml")

	// 确保配置目录存在
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return nil, fmt.Errorf("创建配置目录失败: %w", err)
	}

	v := viper.New()
	v.SetConfigName("config")
	v.SetConfigType("toml")
	v.AddConfigPath(configDir)
	v.AddConfigPath(".")

	// 设置默认值
	v.SetDefault("log.level", "info")

	config := &GlobalConfig{
		Log: LogConfig{
			Level: "info", // 默认值
		},
	}

	manager := &GlobalManager{
		viper:  v,
		path:   configPath,
		Config: config,
	}

	// 初始化便捷字段，指向 Config 中的对应字段
	manager.LLMConfig = &config.LLM
	manager.GitHubConfig = &config.GitHub
	manager.UserConfig = &config.User
	manager.JiraConfig = &config.Jira
	manager.LogConfig = &config.Log
	manager.ProxyConfig = &config.Proxy

	return manager, nil
}

// Load 加载配置文件
//
// 从配置文件加载配置到内存，并更新 Config 字段。
// 如果配置文件不存在，会创建默认配置文件。
func (m *GlobalManager) Load() error {
	logger := logging.GetLogger()
	logger.Infof("Loading config from: %s", m.path)

	if err := m.viper.ReadInConfig(); err != nil {
		// 配置文件不存在时，创建默认配置
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			logger.Infof("Creating default config file: %s", m.path)
			if err := m.SaveDefault(); err != nil {
				return err
			}
			// 重新读取刚创建的配置文件
			if err := m.viper.ReadInConfig(); err != nil {
				logger.WithError(err).Error("Failed to read newly created config file")
				return fmt.Errorf("读取配置文件失败: %w", err)
			}
		} else {
			logger.WithError(err).Error("Config operation failed")
			return fmt.Errorf("读取配置文件失败: %w", err)
		}
	}

	// 从 viper 加载配置到 Config 字段
	m.Config = m.getGlobalConfig()

	// 更新便捷字段的指针（确保指向最新的 Config）
	m.LLMConfig = &m.Config.LLM
	m.GitHubConfig = &m.Config.GitHub
	m.UserConfig = &m.Config.User
	m.JiraConfig = &m.Config.Jira
	m.LogConfig = &m.Config.Log
	m.ProxyConfig = &m.Config.Proxy

	return nil
}

// Save 保存配置到文件
//
// 保存当前 Config 字段的内容到文件。
// 保存后会自动重新加载以同步 viper。
func (m *GlobalManager) Save() error {
	logger := logging.GetLogger()
	logger.Infof("Saving config to: %s", m.path)

	err := SaveConfigToFile(m.path, m.Config)
	if err != nil {
		logger.WithError(err).Error("Config operation failed")
		return err
	}

	// 重新加载以同步 viper 和 Config 字段
	if err := m.Load(); err != nil {
		return err
	}

	return nil
}

// SaveDefault 保存默认配置
func (m *GlobalManager) SaveDefault() error {
	m.Config = &GlobalConfig{
		Log: LogConfig{
			Level: "info",
		},
	}
	return m.Save()
}

// GetConfigPath 获取配置文件路径
func (m *GlobalManager) GetConfigPath() string {
	return m.path
}

// GetLLMConfig 获取 LLM 配置
//
// 返回 Config 字段中的 LLM 配置的引用。
// 可以直接使用 manager.Config.LLM 访问，此方法保留以保持向后兼容。
//
// 返回:
//   - *LLMConfig: LLM 配置结构（指向 Config.LLM）
func (m *GlobalManager) GetLLMConfig() *LLMConfig {
	if m.Config == nil {
		return &LLMConfig{}
	}
	return &m.Config.LLM
}

// GetGitHubConfig 获取 GitHub 配置
//
// 返回 Config 字段中的 GitHub 配置的引用。
// 可以直接使用 manager.Config.GitHub 或 manager.GitHubConfig 访问，此方法保留以保持向后兼容。
//
// 返回:
//   - *GitHubConfig: GitHub 配置结构（指向 Config.GitHub）
func (m *GlobalManager) GetGitHubConfig() *GitHubConfig {
	if m.Config == nil {
		return &GitHubConfig{}
	}
	return &m.Config.GitHub
}

// GetUserConfig 获取用户配置
//
// 返回 Config 字段中的用户配置的引用。
// 可以直接使用 manager.Config.User 或 manager.UserConfig 访问，此方法保留以保持向后兼容。
//
// 返回:
//   - *UserConfig: 用户配置结构（指向 Config.User）
func (m *GlobalManager) GetUserConfig() *UserConfig {
	if m.Config == nil {
		return &UserConfig{}
	}
	return &m.Config.User
}

// GetJiraConfig 获取 Jira 配置
//
// 返回 Config 字段中的 Jira 配置的引用。
// 可以直接使用 manager.Config.Jira 或 manager.JiraConfig 访问，此方法保留以保持向后兼容。
//
// 返回:
//   - *JiraConfig: Jira 配置结构（指向 Config.Jira）
func (m *GlobalManager) GetJiraConfig() *JiraConfig {
	if m.Config == nil {
		return &JiraConfig{}
	}
	return &m.Config.Jira
}

// GetLogConfig 获取日志配置
//
// 返回 Config 字段中的日志配置的引用。
// 可以直接使用 manager.Config.Log 或 manager.LogConfig 访问，此方法保留以保持向后兼容。
//
// 返回:
//   - *LogConfig: 日志配置结构（指向 Config.Log）
func (m *GlobalManager) GetLogConfig() *LogConfig {
	if m.Config == nil {
		return &LogConfig{}
	}
	return &m.Config.Log
}

// GetProxyConfig 获取代理配置
//
// 返回 Config 字段中的代理配置的引用。
// 可以直接使用 manager.Config.Proxy 或 manager.ProxyConfig 访问，此方法保留以保持向后兼容。
//
// 返回:
//   - *ProxyConfig: 代理配置结构（指向 Config.Proxy）
func (m *GlobalManager) GetProxyConfig() *ProxyConfig {
	if m.Config == nil {
		return &ProxyConfig{}
	}
	return &m.Config.Proxy
}

// Global 获取全局 GlobalManager 单例
//
// 返回进程级别的 GlobalManager 单例。
// 单例会在首次调用时初始化，后续调用会复用同一个实例。
//
// 返回:
//   - *GlobalManager: 全局配置管理器实例
//   - error: 如果创建失败，返回错误
//
// 注意:
//   - 首次调用时如果创建失败，后续调用会返回相同的错误
//   - 线程安全：可以在多线程环境中安全使用
//
// 优势:
//   - 减少资源消耗：避免重复创建管理器实例
//   - 统一管理：所有配置操作使用同一个管理器实例
//   - 配置一致性：确保整个进程使用相同的配置状态
func Global() (*GlobalManager, error) {
	globalOnce.Do(func() {
		globalManager, globalErr = newGlobalManager()
	})
	return globalManager, globalErr
}

// NewGlobalManager 创建全局配置管理器（已废弃，请使用 Global()）
//
// 此函数保留以保持向后兼容，但建议使用 Global() 获取单例。
//
// 返回:
//   - *GlobalManager: 全局配置管理器实例
//   - error: 如果创建失败，返回错误
//
// 已废弃: 请使用 Global() 获取单例
func NewGlobalManager() (*GlobalManager, error) {
	return newGlobalManager()
}

// getGlobalConfig 获取完整的全局配置（私有方法）
//
// 从 viper 中读取完整的全局配置。
// 此方法仅在内部使用（Load 方法中），外部应直接访问 Config 字段。
//
// 返回:
//   - *GlobalConfig: 全局配置结构
func (m *GlobalManager) getGlobalConfig() *GlobalConfig {
	cfg := &GlobalConfig{}

	// 读取用户配置
	cfg.User.Name = m.viper.GetString("user.name")
	cfg.User.Email = m.viper.GetString("user.email")

	// 读取日志配置
	cfg.Log.Level = m.viper.GetString("log.level")
	if cfg.Log.Level == "" {
		cfg.Log.Level = "info"
	}

	// 读取 GitHub 配置
	cfg.GitHub.Current = m.viper.GetString("github.current")
	// 读取账号列表
	if accountsVal := m.viper.Get("github.accounts"); accountsVal != nil {
		if accounts, ok := accountsVal.([]interface{}); ok {
			for _, acc := range accounts {
				if accMap, ok := acc.(map[string]interface{}); ok {
					account := GitHubAccount{}
					if name, ok := accMap["name"].(string); ok {
						account.Name = name
					}
					if token, ok := accMap["token"].(string); ok {
						account.Token = token
					}
					if account.Name != "" || account.Token != "" {
						cfg.GitHub.Accounts = append(cfg.GitHub.Accounts, account)
					}
				}
			}
		}
	}

	// 读取 Jira 配置
	cfg.Jira.URL = m.viper.GetString("jira.url")
	cfg.Jira.Username = m.viper.GetString("jira.username")
	cfg.Jira.Token = m.viper.GetString("jira.token")

	// 读取 LLM 配置
	cfg.LLM.Provider = m.viper.GetString("llm.provider")
	cfg.LLM.Language = m.viper.GetString("llm.language")
	cfg.LLM.OpenAI.APIKey = m.viper.GetString("llm.openai.api_key")
	cfg.LLM.OpenAI.Model = m.viper.GetString("llm.openai.model")
	cfg.LLM.DeepSeek.APIKey = m.viper.GetString("llm.deepseek.api_key")
	cfg.LLM.DeepSeek.Model = m.viper.GetString("llm.deepseek.model")
	cfg.LLM.Proxy.URL = m.viper.GetString("llm.proxy.url")
	cfg.LLM.Proxy.APIKey = m.viper.GetString("llm.proxy.api_key")
	cfg.LLM.Proxy.Model = m.viper.GetString("llm.proxy.model")

	// 读取代理配置
	if enabled := m.viper.Get("proxy.enabled"); enabled != nil {
		if enabledBool, ok := enabled.(bool); ok {
			cfg.Proxy.Enabled = enabledBool
		}
	}
	cfg.Proxy.HTTP = m.viper.GetString("proxy.http")
	cfg.Proxy.HTTPS = m.viper.GetString("proxy.https")

	return cfg
}

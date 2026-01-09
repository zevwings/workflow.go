package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
	"github.com/zevwings/workflow/internal/logging"
)

// GlobalManager 全局配置管理器
//
// 管理用户级别的全局配置：~/.workflow/config.toml
// 包含：用户信息、认证配置（GitHub、Jira）、工具配置（LLM、Proxy、Log）
type GlobalManager struct {
	viper *viper.Viper
	path  string
}

// NewGlobalManager 创建全局配置管理器
//
// 返回:
//   - *GlobalManager: 全局配置管理器实例
//   - error: 如果创建失败，返回错误
func NewGlobalManager() (*GlobalManager, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("获取用户主目录失败: %w", err)
	}

	configDir := filepath.Join(homeDir, ".workflow")
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

	return &GlobalManager{
		viper: v,
		path:  configPath,
	}, nil
}

// Load 加载配置文件
func (m *GlobalManager) Load() error {
	logger := logging.GetLogger()
	logger.Infof("Loading config from: %s", m.path)

	if err := m.viper.ReadInConfig(); err != nil {
		// 配置文件不存在时，创建默认配置
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			logger.Infof("Creating default config file: %s", m.path)
			return m.SaveDefault()
		}
		logger.WithError(err).Error("Config operation failed")
		return fmt.Errorf("读取配置文件失败: %w", err)
	}
	return nil
}

// Save 保存配置到文件
func (m *GlobalManager) Save(config interface{}) error {
	logger := logging.GetLogger()
	logger.Infof("Saving config to: %s", m.path)

	err := SaveConfigToFile(m.path, config)
	if err != nil {
		logger.WithError(err).Error("Config operation failed")
		return err
	}
	return nil
}

// SaveDefault 保存默认配置
func (m *GlobalManager) SaveDefault() error {
	defaultConfig := &GlobalConfig{
		Log: LogConfig{
			Level: "info",
		},
	}
	return m.Save(defaultConfig)
}

// Get 获取配置值
func (m *GlobalManager) Get(key string) interface{} {
	return m.viper.Get(key)
}

// GetString 获取字符串配置值
func (m *GlobalManager) GetString(key string) string {
	return m.viper.GetString(key)
}

// Set 设置配置值
func (m *GlobalManager) Set(key string, value interface{}) {
	m.viper.Set(key, value)
}

// GetConfigPath 获取配置文件路径
func (m *GlobalManager) GetConfigPath() string {
	return m.path
}

// GetLLMConfig 获取 LLM 配置
//
// 从 viper 中读取完整的 LLM 配置。
//
// 返回:
//   - *LLMConfig: LLM 配置结构
func (m *GlobalManager) GetLLMConfig() *LLMConfig {
	cfg := &LLMConfig{
		Provider: m.GetString("llm.provider"),
		Language: m.GetString("llm.language"),
	}

	cfg.OpenAI.APIKey = m.GetString("llm.openai.api_key")
	cfg.OpenAI.Model = m.GetString("llm.openai.model")

	cfg.DeepSeek.APIKey = m.GetString("llm.deepseek.api_key")
	cfg.DeepSeek.Model = m.GetString("llm.deepseek.model")

	cfg.Proxy.URL = m.GetString("llm.proxy.url")
	cfg.Proxy.APIKey = m.GetString("llm.proxy.api_key")
	cfg.Proxy.Model = m.GetString("llm.proxy.model")

	return cfg
}

// GetGitHubConfig 获取 GitHub 配置
//
// 从 viper 中读取完整的 GitHub 配置。
//
// 返回:
//   - *GitHubConfig: GitHub 配置结构
func (m *GlobalManager) GetGitHubConfig() *GitHubConfig {
	cfg := &GitHubConfig{
		Current: m.GetString("github.current"),
	}

	// 读取账号列表
	if accountsVal := m.Get("github.accounts"); accountsVal != nil {
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
						cfg.Accounts = append(cfg.Accounts, account)
					}
				}
			}
		}
	}

	return cfg
}

// GetGlobalConfig 获取完整的全局配置
//
// 从 viper 中读取完整的全局配置。
//
// 返回:
//   - *GlobalConfig: 全局配置结构
func (m *GlobalManager) GetGlobalConfig() *GlobalConfig {
	cfg := &GlobalConfig{}

	// 读取用户配置
	cfg.User.Name = m.GetString("user.name")
	cfg.User.Email = m.GetString("user.email")

	// 读取日志配置
	cfg.Log.Level = m.GetString("log.level")
	if cfg.Log.Level == "" {
		cfg.Log.Level = "info"
	}

	// 读取 GitHub 配置
	cfg.GitHub = *m.GetGitHubConfig()

	// 读取 Jira 配置
	cfg.Jira.URL = m.GetString("jira.url")
	cfg.Jira.Username = m.GetString("jira.username")
	cfg.Jira.Token = m.GetString("jira.token")

	// 读取 LLM 配置
	cfg.LLM = *m.GetLLMConfig()

	// 读取代理配置
	if enabled := m.Get("proxy.enabled"); enabled != nil {
		if enabledBool, ok := enabled.(bool); ok {
			cfg.Proxy.Enabled = enabledBool
		}
	}
	cfg.Proxy.HTTP = m.GetString("proxy.http")
	cfg.Proxy.HTTPS = m.GetString("proxy.https")

	return cfg
}

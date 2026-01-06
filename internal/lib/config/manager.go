package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/pelletier/go-toml/v2"
	"github.com/spf13/viper"
)

// Manager 配置管理器
type Manager struct {
	viper *viper.Viper
	path  string
}

// NewManager 创建新的配置管理器
func NewManager() (*Manager, error) {
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

	return &Manager{
		viper: v,
		path:  configPath,
	}, nil
}

// Load 加载配置文件
func (m *Manager) Load() error {
	if err := m.viper.ReadInConfig(); err != nil {
		// 配置文件不存在时，创建默认配置
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			return m.SaveDefault()
		}
		return fmt.Errorf("读取配置文件失败: %w", err)
	}
	return nil
}

// Save 保存配置到文件
func (m *Manager) Save(config interface{}) error {
	data, err := toml.Marshal(config)
	if err != nil {
		return fmt.Errorf("序列化配置失败: %w", err)
	}

	if err := os.WriteFile(m.path, data, 0644); err != nil {
		return fmt.Errorf("写入配置文件失败: %w", err)
	}

	return nil
}

// SaveDefault 保存默认配置
func (m *Manager) SaveDefault() error {
	defaultConfig := &Config{
		Log: LogConfig{
			Level: "info",
		},
	}
	return m.Save(defaultConfig)
}

// Get 获取配置值
func (m *Manager) Get(key string) interface{} {
	return m.viper.Get(key)
}

// GetString 获取字符串配置值
func (m *Manager) GetString(key string) string {
	return m.viper.GetString(key)
}

// Set 设置配置值
func (m *Manager) Set(key string, value interface{}) {
	m.viper.Set(key, value)
}

// GetConfigPath 获取配置文件路径
func (m *Manager) GetConfigPath() string {
	return m.path
}

// Config 配置结构
type Config struct {
	User   UserConfig   `toml:"user,omitempty"`
	Jira   JiraConfig   `toml:"jira,omitempty"`
	GitHub GitHubConfig `toml:"github,omitempty"`
	Log    LogConfig    `toml:"log,omitempty"`
	LLM    LLMConfig    `toml:"llm,omitempty"`
	Proxy  ProxyConfig  `toml:"proxy,omitempty"`
}

// UserConfig 用户配置
type UserConfig struct {
	Name  string `toml:"name,omitempty"`
	Email string `toml:"email,omitempty"`
}

// JiraConfig Jira 配置
type JiraConfig struct {
	URL      string `toml:"url,omitempty"`
	Username string `toml:"username,omitempty"`
	Token    string `toml:"token,omitempty"`
}

// GitHubConfig GitHub 配置
type GitHubConfig struct {
	Accounts []GitHubAccount `toml:"accounts,omitempty"`
	Current  string          `toml:"current,omitempty"`
}

// GitHubAccount GitHub 账号
type GitHubAccount struct {
	Name  string `toml:"name"`
	Token string `toml:"token"`
}

// LogConfig 日志配置
type LogConfig struct {
	Level string `toml:"level"`
}

// LLMConfig LLM 配置
type LLMConfig struct {
	Provider string `toml:"provider,omitempty"`
	OpenAI   struct {
		APIKey string `toml:"api_key,omitempty"`
	} `toml:"openai,omitempty"`
	DeepSeek struct {
		APIKey string `toml:"api_key,omitempty"`
	} `toml:"deepseek,omitempty"`
}

// ProxyConfig 代理配置
type ProxyConfig struct {
	Enabled bool   `toml:"enabled,omitempty"`
	HTTP    string `toml:"http,omitempty"`
	HTTPS   string `toml:"https,omitempty"`
}


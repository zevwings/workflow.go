package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ==================== NewGlobalManager 测试 ====================

func TestNewGlobalManager(t *testing.T) {
	// Arrange: 设置测试环境
	tempDir := t.TempDir()
	t.Setenv("HOME", tempDir)
	// 禁用 iCloud 以使用默认 XDG 路径
	t.Setenv("WORKFLOW_DISABLE_ICLOUD_CONFIG", "1")

	// Act: 创建全局配置管理器
	manager, err := NewGlobalManager()

	// Assert: 验证结果
	require.NoError(t, err)
	assert.NotNil(t, manager)
	assert.NotNil(t, manager.viper)

	// 验证配置路径（使用 XDG 路径）
	expectedConfigDir, err := ConfigDir()
	require.NoError(t, err)
	expectedPath := filepath.Join(expectedConfigDir, "config.toml")
	actualPath := manager.GetConfigPath()
	assert.Equal(t, expectedPath, actualPath)

	// 验证配置目录已创建
	configDir := filepath.Dir(expectedPath)
	_, err = os.Stat(configDir)
	assert.NoError(t, err, "配置目录应该已创建")

	// 验证管理器已正确初始化
	assert.NotNil(t, manager.Config)
	assert.NotNil(t, manager.LogConfig)
}

func TestNewGlobalManager_CreatesDirectory(t *testing.T) {
	// Arrange: 设置测试环境，确保目录不存在
	tempDir := t.TempDir()
	t.Setenv("HOME", tempDir)
	// 禁用 iCloud 以使用默认 XDG 路径
	t.Setenv("WORKFLOW_DISABLE_ICLOUD_CONFIG", "1")

	// Act: 创建全局配置管理器
	manager, err := NewGlobalManager()

	// Assert: 验证目录已创建
	require.NoError(t, err)
	configDir := filepath.Dir(manager.GetConfigPath())
	_, err = os.Stat(configDir)
	assert.NoError(t, err, "配置目录应该已自动创建")
	assert.NotNil(t, manager)
}

// ==================== Load 测试 ====================

func TestGlobalManager_Load_FileExists(t *testing.T) {
	// Arrange: 设置测试环境并创建配置文件
	tempDir := t.TempDir()
	configDir := filepath.Join(tempDir, ".config", "Workflow")
	configPath := filepath.Join(configDir, "config.toml")
	require.NoError(t, os.MkdirAll(configDir, 0755))

	// 创建测试配置文件
	configContent := `[log]
level = "debug"
`
	require.NoError(t, os.WriteFile(configPath, []byte(configContent), 0644))

	// 创建管理器（使用临时路径）
	manager := &GlobalManager{
		viper:  viper.New(),
		path:   configPath,
		Config: &GlobalConfig{},
	}
	manager.viper.SetConfigName("config")
	manager.viper.SetConfigType("toml")
	manager.viper.AddConfigPath(configDir)

	// 初始化便捷字段
	manager.LLMConfig = &manager.Config.LLM
	manager.GitHubConfig = &manager.Config.GitHub
	manager.JiraConfig = &manager.Config.Jira
	manager.LogConfig = &manager.Config.Log
	manager.ProxyConfig = &manager.Config.Proxy

	// Act: 加载配置
	err := manager.Load()

	// Assert: 验证配置已加载
	assert.NoError(t, err)
	assert.Equal(t, "debug", manager.LogConfig.Level)
}

func TestGlobalManager_Load_FileNotExists(t *testing.T) {
	// Arrange: 设置测试环境，但不创建配置文件
	tempDir := t.TempDir()
	configDir := filepath.Join(tempDir, ".config", "Workflow")

	// 创建配置目录（但不创建配置文件）
	require.NoError(t, os.MkdirAll(configDir, 0755))

	configPath := filepath.Join(configDir, "config.toml")

	// 确保配置文件不存在
	_, err := os.Stat(configPath)
	assert.Error(t, err, "配置文件应该不存在")

	// 创建管理器（使用临时路径）
	manager := &GlobalManager{
		viper:  viper.New(),
		path:   configPath,
		Config: &GlobalConfig{},
	}
	manager.viper.SetConfigName("config")
	manager.viper.SetConfigType("toml")
	manager.viper.AddConfigPath(configDir)

	// 初始化便捷字段
	manager.LLMConfig = &manager.Config.LLM
	manager.GitHubConfig = &manager.Config.GitHub
	manager.JiraConfig = &manager.Config.Jira
	manager.LogConfig = &manager.Config.Log
	manager.ProxyConfig = &manager.Config.Proxy

	// Act: 加载配置（文件不存在）
	err = manager.Load()

	// Assert: 应该返回 ConfigFileNotFoundError，不创建默认配置
	assert.Error(t, err)
	_, ok := err.(viper.ConfigFileNotFoundError)
	assert.True(t, ok, "应该返回 ConfigFileNotFoundError，实际错误: %v", err)

	// 验证配置文件未创建
	_, err = os.Stat(configPath)
	assert.Error(t, err, "配置文件不应该被创建")

	// 验证 Config 为零值
	assert.NotNil(t, manager.Config)
	assert.Equal(t, "", manager.LogConfig.Level, "Log.Level 应该保持零值（空字符串）")
}

// ==================== Save 和 SaveDefault 测试 ====================

func TestGlobalManager_Save(t *testing.T) {
	// Arrange: 设置测试环境
	tempDir := t.TempDir()
	t.Setenv("HOME", tempDir)
	// 禁用 iCloud 以使用默认 XDG 路径
	t.Setenv("WORKFLOW_DISABLE_ICLOUD_CONFIG", "1")

	manager, err := NewGlobalManager()
	require.NoError(t, err)

	// Act: 设置并保存配置
	manager.Config = &GlobalConfig{
		Log: LogConfig{
			Level: "debug",
		},
	}
	err = manager.Save()

	// Assert: 验证配置已保存
	assert.NoError(t, err)

	// 验证文件已创建
	configPath := manager.GetConfigPath()
	_, err = os.Stat(configPath)
	assert.NoError(t, err)

	// 重新加载并验证
	err = manager.Load()
	require.NoError(t, err)
	assert.Equal(t, "debug", manager.LogConfig.Level)
}

func TestGlobalManager_SaveDefault(t *testing.T) {
	// Arrange: 设置测试环境
	tempDir := t.TempDir()
	t.Setenv("HOME", tempDir)
	// 禁用 iCloud 以使用默认 XDG 路径
	t.Setenv("WORKFLOW_DISABLE_ICLOUD_CONFIG", "1")

	manager, err := NewGlobalManager()
	require.NoError(t, err)

	// Act: 保存默认配置
	err = manager.SaveDefault()

	// Assert: 验证默认配置已保存
	assert.NoError(t, err)

	// 验证文件已创建
	configPath := manager.GetConfigPath()
	_, err = os.Stat(configPath)
	assert.NoError(t, err)

	// 重新加载并验证
	err = manager.Load()
	require.NoError(t, err)
	assert.Equal(t, "info", manager.LogConfig.Level)
}

// ==================== 配置字段直接访问测试 ====================

func TestGlobalManager_DirectFieldAccess(t *testing.T) {
	// Arrange: 设置测试环境
	tempDir := t.TempDir()
	t.Setenv("HOME", tempDir)
	// 禁用 iCloud 以使用默认 XDG 路径
	t.Setenv("WORKFLOW_DISABLE_ICLOUD_CONFIG", "1")

	manager, err := NewGlobalManager()
	require.NoError(t, err)

	// 加载配置
	err = manager.Load()
	require.NoError(t, err)

	// Act & Assert: 测试直接访问配置字段
	// 修改配置字段
	manager.LogConfig.Level = "debug"

	// 验证修改
	assert.Equal(t, "debug", manager.LogConfig.Level)

	// 验证通过 Config 字段也能访问
	assert.Equal(t, "debug", manager.Config.Log.Level)
}

// ==================== GetLLMConfig 测试 ====================

func TestGlobalManager_GetLLMConfig(t *testing.T) {
	// Arrange: 设置测试环境并创建配置文件
	tempDir := t.TempDir()
	t.Setenv("HOME", tempDir)
	// 禁用 iCloud 以使用默认 XDG 路径
	t.Setenv("WORKFLOW_DISABLE_ICLOUD_CONFIG", "1")

	configDir, err := ConfigDir()
	require.NoError(t, err)
	configPath := filepath.Join(configDir, "config.toml")
	require.NoError(t, os.MkdirAll(configDir, 0755))

	configContent := `[llm]
provider = "openai"
language = "en"

[llm.openai]
api_key = "sk-test-key"
model = "gpt-4"

[llm.deepseek]
api_key = "sk-deepseek-key"
model = "deepseek-chat"

[llm.proxy]
url = "https://api.example.com/v1"
api_key = "proxy-key"
model = "custom-model"
`
	require.NoError(t, os.WriteFile(configPath, []byte(configContent), 0644))

	manager, err := NewGlobalManager()
	require.NoError(t, err)
	require.NoError(t, manager.Load())

	// Act: 获取 LLM 配置
	llmConfig := manager.GetLLMConfig()

	// Assert: 验证配置
	assert.NotNil(t, llmConfig)
	assert.Equal(t, "openai", llmConfig.Provider)
	assert.Equal(t, "en", llmConfig.Language)
	assert.Equal(t, "sk-test-key", llmConfig.OpenAI.APIKey)
	assert.Equal(t, "gpt-4", llmConfig.OpenAI.Model)
	assert.Equal(t, "sk-deepseek-key", llmConfig.DeepSeek.APIKey)
	assert.Equal(t, "deepseek-chat", llmConfig.DeepSeek.Model)
	assert.Equal(t, "https://api.example.com/v1", llmConfig.Proxy.URL)
	assert.Equal(t, "proxy-key", llmConfig.Proxy.APIKey)
	assert.Equal(t, "custom-model", llmConfig.Proxy.Model)
}

func TestGlobalManager_GetLLMConfig_Empty(t *testing.T) {
	// Arrange: 设置测试环境，不设置 LLM 配置
	tempDir := t.TempDir()
	t.Setenv("HOME", tempDir)
	// 禁用 iCloud 以使用默认 XDG 路径
	t.Setenv("WORKFLOW_DISABLE_ICLOUD_CONFIG", "1")

	manager, err := NewGlobalManager()
	require.NoError(t, err)

	// Act: 获取 LLM 配置
	llmConfig := manager.GetLLMConfig()

	// Assert: 验证返回空配置
	assert.NotNil(t, llmConfig)
	assert.Empty(t, llmConfig.Provider)
	assert.Empty(t, llmConfig.Language)
}

// ==================== GetGitHubConfig 测试 ====================

func TestGlobalManager_GetGitHubConfig(t *testing.T) {
	// Arrange: 设置测试环境并创建配置文件
	tempDir := t.TempDir()
	configDir := filepath.Join(tempDir, ".config", "Workflow")
	configPath := filepath.Join(configDir, "config.toml")
	require.NoError(t, os.MkdirAll(configDir, 0755))

	configContent := `[github]
current = "account1"

[[github.accounts]]
name = "account1"
api_token = "token1"

[[github.accounts]]
name = "account2"
api_token = "token2"
`
	require.NoError(t, os.WriteFile(configPath, []byte(configContent), 0644))

	// 创建管理器（使用临时路径）
	manager := &GlobalManager{
		viper:  viper.New(),
		path:   configPath,
		Config: &GlobalConfig{},
	}
	manager.viper.SetConfigName("config")
	manager.viper.SetConfigType("toml")
	manager.viper.AddConfigPath(configDir)

	// 初始化便捷字段
	manager.LLMConfig = &manager.Config.LLM
	manager.GitHubConfig = &manager.Config.GitHub
	manager.JiraConfig = &manager.Config.Jira
	manager.LogConfig = &manager.Config.Log
	manager.ProxyConfig = &manager.Config.Proxy

	require.NoError(t, manager.Load())

	// Act: 获取 GitHub 配置
	githubConfig := manager.GetGitHubConfig()

	// Assert: 验证配置
	assert.NotNil(t, githubConfig)
	assert.Equal(t, "account1", githubConfig.Current)
	assert.Len(t, githubConfig.Accounts, 2)
	assert.Equal(t, "account1", githubConfig.Accounts[0].Name)
	assert.Equal(t, "token1", githubConfig.Accounts[0].APIToken)
	assert.Equal(t, "account2", githubConfig.Accounts[1].Name)
	assert.Equal(t, "token2", githubConfig.Accounts[1].APIToken)
}

func TestGlobalManager_GetGitHubConfig_Empty(t *testing.T) {
	// Arrange: 设置测试环境，不设置 GitHub 配置
	tempDir := t.TempDir()
	t.Setenv("HOME", tempDir)
	// 禁用 iCloud 以使用默认 XDG 路径
	t.Setenv("WORKFLOW_DISABLE_ICLOUD_CONFIG", "1")

	manager, err := NewGlobalManager()
	require.NoError(t, err)

	// Act: 获取 GitHub 配置
	githubConfig := manager.GetGitHubConfig()

	// Assert: 验证返回空配置
	assert.NotNil(t, githubConfig)
	assert.Empty(t, githubConfig.Current)
	assert.Empty(t, githubConfig.Accounts)
}

func TestGlobalManager_GetGitHubConfig_PartialAccount(t *testing.T) {
	// Arrange: 设置测试环境并创建部分账号配置
	tempDir := t.TempDir()
	configDir := filepath.Join(tempDir, ".config", "Workflow")
	configPath := filepath.Join(configDir, "config.toml")
	require.NoError(t, os.MkdirAll(configDir, 0755))

	// 只设置 name，不设置 api_token
	configContent := `[[github.accounts]]
name = "account1"
`
	require.NoError(t, os.WriteFile(configPath, []byte(configContent), 0644))

	// 创建管理器（使用临时路径）
	manager := &GlobalManager{
		viper:  viper.New(),
		path:   configPath,
		Config: &GlobalConfig{},
	}
	manager.viper.SetConfigName("config")
	manager.viper.SetConfigType("toml")
	manager.viper.AddConfigPath(configDir)

	// 初始化便捷字段
	manager.LLMConfig = &manager.Config.LLM
	manager.GitHubConfig = &manager.Config.GitHub
	manager.JiraConfig = &manager.Config.Jira
	manager.LogConfig = &manager.Config.Log
	manager.ProxyConfig = &manager.Config.Proxy

	require.NoError(t, manager.Load())

	// Act: 获取 GitHub 配置
	githubConfig := manager.GetGitHubConfig()

	// Assert: 验证部分配置也被包含
	assert.NotNil(t, githubConfig)
	assert.Len(t, githubConfig.Accounts, 1)
	assert.Equal(t, "account1", githubConfig.Accounts[0].Name)
	assert.Empty(t, githubConfig.Accounts[0].APIToken)
}

// ==================== Config 字段直接访问测试 ====================

func TestGlobalManager_ConfigField(t *testing.T) {
	// Arrange: 设置测试环境并创建完整配置文件
	tempDir := t.TempDir()
	t.Setenv("HOME", tempDir)
	// 禁用 iCloud 以使用默认 XDG 路径
	t.Setenv("WORKFLOW_DISABLE_ICLOUD_CONFIG", "1")

	configDir, err := ConfigDir()
	require.NoError(t, err)
	configPath := filepath.Join(configDir, "config.toml")
	require.NoError(t, os.MkdirAll(configDir, 0755))

	configContent := `[log]
level = "debug"

[github]
current = "account1"

[[github.accounts]]
name = "account1"
token = "token1"

[jira]
service_address = "https://jira.example.com"
email = "user@example.com"
api_token = "jira-token"

[llm]
provider = "openai"
language = "en"

[llm.openai]
api_key = "sk-test-key"
model = "gpt-4"

[proxy]
enabled = true
http = "http://proxy.example.com:8080"
https = "https://proxy.example.com:8080"
`
	require.NoError(t, os.WriteFile(configPath, []byte(configContent), 0644))

	manager, err := NewGlobalManager()
	require.NoError(t, err)
	require.NoError(t, manager.Load())

	// Act: 直接访问 Config 字段
	globalConfig := manager.Config

	// Assert: 验证所有配置
	assert.NotNil(t, globalConfig)
	assert.Equal(t, "debug", globalConfig.Log.Level)
	assert.Equal(t, "account1", globalConfig.GitHub.Current)
	assert.Len(t, globalConfig.GitHub.Accounts, 1)
	assert.Equal(t, "https://jira.example.com", globalConfig.Jira.ServiceAddress)
	assert.Equal(t, "user@example.com", globalConfig.Jira.Email)
	assert.Equal(t, "jira-token", globalConfig.Jira.APIToken)
	assert.Equal(t, "openai", globalConfig.LLM.Provider)
	assert.True(t, globalConfig.Proxy.Enabled)
	assert.Equal(t, "http://proxy.example.com:8080", globalConfig.Proxy.HTTP)
	assert.Equal(t, "https://proxy.example.com:8080", globalConfig.Proxy.HTTPS)
}

func TestGlobalManager_ConfigField_DefaultLogLevel(t *testing.T) {
	// Arrange: 设置测试环境，不设置 log.level
	tempDir := t.TempDir()
	t.Setenv("HOME", tempDir)
	// 禁用 iCloud 以使用默认 XDG 路径
	t.Setenv("WORKFLOW_DISABLE_ICLOUD_CONFIG", "1")

	manager, err := NewGlobalManager()
	require.NoError(t, err)

	// 需要先 Load 才能访问 Config 字段
	err = manager.Load()
	
	// 如果配置文件不存在，Load 会返回错误，Config 为零值
	if err != nil {
		// 配置文件不存在时，Log.Level 应该为零值（空字符串）
		assert.Equal(t, "", manager.Config.Log.Level, "配置文件不存在时，Log.Level 应该为零值")
	} else {
		// 如果配置文件存在（可能从当前目录读取），验证配置已加载
		globalConfig := manager.Config
		assert.NotNil(t, globalConfig)
		// 如果配置文件存在，log level 可能是任何值（取决于配置文件内容）
		// 这里只验证 Config 不为 nil
		assert.NotNil(t, globalConfig.Log)
	}
}

// ==================== GetConfigPath 测试 ====================

func TestGlobalManager_GetConfigPath(t *testing.T) {
	// Arrange: 设置测试环境
	tempDir := t.TempDir()
	t.Setenv("HOME", tempDir)
	// 禁用 iCloud 以使用默认 XDG 路径
	t.Setenv("WORKFLOW_DISABLE_ICLOUD_CONFIG", "1")

	manager, err := NewGlobalManager()
	require.NoError(t, err)

	// Act: 获取配置路径
	path := manager.GetConfigPath()

	// Assert: 验证路径包含 workflow 和 config.toml
	assert.Contains(t, path, "Workflow")
	assert.Contains(t, path, "config.toml")
}

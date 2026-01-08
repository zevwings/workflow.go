package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ==================== NewGlobalManager 测试 ====================

func TestNewGlobalManager(t *testing.T) {
	// Arrange: 设置测试环境
	tempDir := t.TempDir()
	t.Setenv("HOME", tempDir)

	// Act: 创建全局配置管理器
	manager, err := NewGlobalManager()

	// Assert: 验证结果
	require.NoError(t, err)
	assert.NotNil(t, manager)
	assert.NotNil(t, manager.viper)

	// 验证配置路径
	expectedPath := filepath.Join(tempDir, ".workflow", "config.toml")
	assert.Equal(t, expectedPath, manager.GetConfigPath())

	// 验证配置目录已创建
	configDir := filepath.Dir(expectedPath)
	_, err = os.Stat(configDir)
	assert.NoError(t, err, "配置目录应该已创建")

	// 验证默认值已设置
	assert.Equal(t, "info", manager.GetString("log.level"))
}

func TestNewGlobalManager_CreatesDirectory(t *testing.T) {
	// Arrange: 设置测试环境，确保目录不存在
	tempDir := t.TempDir()
	t.Setenv("HOME", tempDir)

	// 删除可能存在的目录
	configDir := filepath.Join(tempDir, ".workflow")
	os.RemoveAll(configDir)

	// Act: 创建全局配置管理器
	manager, err := NewGlobalManager()

	// Assert: 验证目录已创建
	require.NoError(t, err)
	_, err = os.Stat(configDir)
	assert.NoError(t, err, "配置目录应该已自动创建")
	assert.NotNil(t, manager)
}

// ==================== Load 测试 ====================

func TestGlobalManager_Load_FileExists(t *testing.T) {
	// Arrange: 设置测试环境并创建配置文件
	tempDir := t.TempDir()
	t.Setenv("HOME", tempDir)

	configDir := filepath.Join(tempDir, ".workflow")
	configPath := filepath.Join(configDir, "config.toml")
	require.NoError(t, os.MkdirAll(configDir, 0755))

	// 创建测试配置文件
	configContent := `[log]
level = "debug"

[user]
name = "Test User"
email = "test@example.com"
`
	require.NoError(t, os.WriteFile(configPath, []byte(configContent), 0644))

	// 创建管理器
	manager, err := NewGlobalManager()
	require.NoError(t, err)

	// Act: 加载配置
	err = manager.Load()

	// Assert: 验证配置已加载
	assert.NoError(t, err)
	assert.Equal(t, "debug", manager.GetString("log.level"))
	assert.Equal(t, "Test User", manager.GetString("user.name"))
	assert.Equal(t, "test@example.com", manager.GetString("user.email"))
}

func TestGlobalManager_Load_FileNotExists(t *testing.T) {
	// Arrange: 设置测试环境，但不创建配置文件
	tempDir := t.TempDir()
	t.Setenv("HOME", tempDir)

	// 创建管理器
	manager, err := NewGlobalManager()
	require.NoError(t, err)

	// Act: 加载配置（文件不存在）
	err = manager.Load()

	// Assert: 应该创建默认配置
	assert.NoError(t, err)

	// 验证默认配置文件已创建
	configPath := manager.GetConfigPath()
	_, err = os.Stat(configPath)
	assert.NoError(t, err, "默认配置文件应该已创建")

	// 验证默认值
	assert.Equal(t, "info", manager.GetString("log.level"))
}

// ==================== Save 和 SaveDefault 测试 ====================

func TestGlobalManager_Save(t *testing.T) {
	// Arrange: 设置测试环境
	tempDir := t.TempDir()
	t.Setenv("HOME", tempDir)

	manager, err := NewGlobalManager()
	require.NoError(t, err)

	config := &GlobalConfig{
		Log: LogConfig{
			Level: "debug",
		},
		User: UserConfig{
			Name:  "Test User",
			Email: "test@example.com",
		},
	}

	// Act: 保存配置
	err = manager.Save(config)

	// Assert: 验证配置已保存
	assert.NoError(t, err)

	// 验证文件已创建
	configPath := manager.GetConfigPath()
	_, err = os.Stat(configPath)
	assert.NoError(t, err)

	// 重新加载并验证
	err = manager.Load()
	require.NoError(t, err)
	assert.Equal(t, "debug", manager.GetString("log.level"))
	assert.Equal(t, "Test User", manager.GetString("user.name"))
	assert.Equal(t, "test@example.com", manager.GetString("user.email"))
}

func TestGlobalManager_SaveDefault(t *testing.T) {
	// Arrange: 设置测试环境
	tempDir := t.TempDir()
	t.Setenv("HOME", tempDir)

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
	assert.Equal(t, "info", manager.GetString("log.level"))
}

// ==================== Get, GetString, Set 测试 ====================

func TestGlobalManager_GetAndSet(t *testing.T) {
	// Arrange: 设置测试环境
	tempDir := t.TempDir()
	t.Setenv("HOME", tempDir)

	manager, err := NewGlobalManager()
	require.NoError(t, err)

	// Act & Assert: 测试 Set 和 Get
	manager.Set("test.key", "test-value")
	value := manager.Get("test.key")
	assert.Equal(t, "test-value", value)

	// 测试 GetString
	strValue := manager.GetString("test.key")
	assert.Equal(t, "test-value", strValue)

	// 测试不存在的键
	nonExistent := manager.GetString("non.existent")
	assert.Empty(t, nonExistent)
}

// ==================== GetLLMConfig 测试 ====================

func TestGlobalManager_GetLLMConfig(t *testing.T) {
	// Arrange: 设置测试环境并创建配置文件
	tempDir := t.TempDir()
	t.Setenv("HOME", tempDir)

	configDir := filepath.Join(tempDir, ".workflow")
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
	t.Setenv("HOME", tempDir)

	configDir := filepath.Join(tempDir, ".workflow")
	configPath := filepath.Join(configDir, "config.toml")
	require.NoError(t, os.MkdirAll(configDir, 0755))

	configContent := `[github]
current = "account1"

[[github.accounts]]
name = "account1"
token = "token1"

[[github.accounts]]
name = "account2"
token = "token2"
`
	require.NoError(t, os.WriteFile(configPath, []byte(configContent), 0644))

	manager, err := NewGlobalManager()
	require.NoError(t, err)
	require.NoError(t, manager.Load())

	// Act: 获取 GitHub 配置
	githubConfig := manager.GetGitHubConfig()

	// Assert: 验证配置
	assert.NotNil(t, githubConfig)
	assert.Equal(t, "account1", githubConfig.Current)
	assert.Len(t, githubConfig.Accounts, 2)
	assert.Equal(t, "account1", githubConfig.Accounts[0].Name)
	assert.Equal(t, "token1", githubConfig.Accounts[0].Token)
	assert.Equal(t, "account2", githubConfig.Accounts[1].Name)
	assert.Equal(t, "token2", githubConfig.Accounts[1].Token)
}

func TestGlobalManager_GetGitHubConfig_Empty(t *testing.T) {
	// Arrange: 设置测试环境，不设置 GitHub 配置
	tempDir := t.TempDir()
	t.Setenv("HOME", tempDir)

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
	t.Setenv("HOME", tempDir)

	configDir := filepath.Join(tempDir, ".workflow")
	configPath := filepath.Join(configDir, "config.toml")
	require.NoError(t, os.MkdirAll(configDir, 0755))

	// 只设置 name，不设置 token
	configContent := `[[github.accounts]]
name = "account1"
`
	require.NoError(t, os.WriteFile(configPath, []byte(configContent), 0644))

	manager, err := NewGlobalManager()
	require.NoError(t, err)
	require.NoError(t, manager.Load())

	// Act: 获取 GitHub 配置
	githubConfig := manager.GetGitHubConfig()

	// Assert: 验证部分配置也被包含
	assert.NotNil(t, githubConfig)
	assert.Len(t, githubConfig.Accounts, 1)
	assert.Equal(t, "account1", githubConfig.Accounts[0].Name)
	assert.Empty(t, githubConfig.Accounts[0].Token)
}

// ==================== GetGlobalConfig 测试 ====================

func TestGlobalManager_GetGlobalConfig(t *testing.T) {
	// Arrange: 设置测试环境并创建完整配置文件
	tempDir := t.TempDir()
	t.Setenv("HOME", tempDir)

	configDir := filepath.Join(tempDir, ".workflow")
	configPath := filepath.Join(configDir, "config.toml")
	require.NoError(t, os.MkdirAll(configDir, 0755))

	configContent := `[user]
name = "Test User"
email = "test@example.com"

[log]
level = "debug"

[github]
current = "account1"

[[github.accounts]]
name = "account1"
token = "token1"

[jira]
url = "https://jira.example.com"
username = "user"
token = "jira-token"

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

	// Act: 获取完整全局配置
	globalConfig := manager.GetGlobalConfig()

	// Assert: 验证所有配置
	assert.NotNil(t, globalConfig)
	assert.Equal(t, "Test User", globalConfig.User.Name)
	assert.Equal(t, "test@example.com", globalConfig.User.Email)
	assert.Equal(t, "debug", globalConfig.Log.Level)
	assert.Equal(t, "account1", globalConfig.GitHub.Current)
	assert.Len(t, globalConfig.GitHub.Accounts, 1)
	assert.Equal(t, "https://jira.example.com", globalConfig.Jira.URL)
	assert.Equal(t, "user", globalConfig.Jira.Username)
	assert.Equal(t, "jira-token", globalConfig.Jira.Token)
	assert.Equal(t, "openai", globalConfig.LLM.Provider)
	assert.True(t, globalConfig.Proxy.Enabled)
	assert.Equal(t, "http://proxy.example.com:8080", globalConfig.Proxy.HTTP)
	assert.Equal(t, "https://proxy.example.com:8080", globalConfig.Proxy.HTTPS)
}

func TestGlobalManager_GetGlobalConfig_DefaultLogLevel(t *testing.T) {
	// Arrange: 设置测试环境，不设置 log.level
	tempDir := t.TempDir()
	t.Setenv("HOME", tempDir)

	manager, err := NewGlobalManager()
	require.NoError(t, err)

	// Act: 获取全局配置
	globalConfig := manager.GetGlobalConfig()

	// Assert: 验证默认 log level
	assert.NotNil(t, globalConfig)
	assert.Equal(t, "info", globalConfig.Log.Level)
}

// ==================== GetConfigPath 测试 ====================

func TestGlobalManager_GetConfigPath(t *testing.T) {
	// Arrange: 设置测试环境
	tempDir := t.TempDir()
	t.Setenv("HOME", tempDir)

	manager, err := NewGlobalManager()
	require.NoError(t, err)

	// Act: 获取配置路径
	path := manager.GetConfigPath()

	// Assert: 验证路径
	expectedPath := filepath.Join(tempDir, ".workflow", "config.toml")
	assert.Equal(t, expectedPath, path)
}


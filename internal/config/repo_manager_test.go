package config

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ==================== Mock 实现 ====================

// mockGitRepository 实现 GitRepository 接口用于测试
type mockGitRepository struct {
	repoPath string
	isGitRepo bool
	remoteURL string
	openError error
}

func (m *mockGitRepository) GetRepoPath() string {
	return m.repoPath
}

func (m *mockGitRepository) IsGitRepo(path string) bool {
	return m.isGitRepo
}

func (m *mockGitRepository) Open(path string) (GitRepo, error) {
	if m.openError != nil {
		return nil, m.openError
	}
	return &mockGitRepo{remoteURL: m.remoteURL}, nil
}

// mockGitRepo 实现 GitRepo 接口用于测试
type mockGitRepo struct {
	remoteURL string
	getRemoteURLError error
}

func (m *mockGitRepo) GetRemoteURL(name string) (string, error) {
	if m.getRemoteURLError != nil {
		return "", m.getRemoteURLError
	}
	if name == "origin" {
		return m.remoteURL, nil
	}
	return "", fmt.Errorf("remote %s not found", name)
}

// ==================== NewRepoManager 测试 ====================

func TestNewRepoManager_WithGitRepo(t *testing.T) {
	// Arrange: 设置测试环境
	tempDir := t.TempDir()
	t.Setenv("HOME", tempDir)

	mockGitRepo := &mockGitRepository{
		repoPath:  tempDir,
		isGitRepo: true,
		remoteURL: "https://github.com/owner/repo.git",
	}

	// Act: 创建仓库配置管理器
	manager, err := NewRepoManager(mockGitRepo)

	// Assert: 验证结果
	require.NoError(t, err)
	assert.NotNil(t, manager)
	assert.NotEmpty(t, manager.GetRepoID())
	assert.Equal(t, tempDir, manager.repoPath)

	// 验证配置路径
	expectedPublicPath := filepath.Join(tempDir, ".workflow", "config.toml")
	assert.Equal(t, expectedPublicPath, manager.GetPublicConfigPath())

	expectedPrivatePath := filepath.Join(tempDir, ".workflow", "config", "repository.toml")
	assert.Equal(t, expectedPrivatePath, manager.GetPrivateConfigPath())
}

func TestNewRepoManager_WithGitRepo_EmptyPath(t *testing.T) {
	// Arrange: 设置测试环境
	tempDir := t.TempDir()
	t.Setenv("HOME", tempDir)

	// 设置当前工作目录
	originalWd, err := os.Getwd()
	require.NoError(t, err)
	defer os.Chdir(originalWd)

	require.NoError(t, os.Chdir(tempDir))

	mockGitRepo := &mockGitRepository{
		repoPath:  "", // 空路径，应该使用当前目录
		isGitRepo: true,
		remoteURL: "https://github.com/owner/repo.git",
	}

	// Act: 创建仓库配置管理器
	manager, err := NewRepoManager(mockGitRepo)

	// Assert: 验证结果
	require.NoError(t, err)
	assert.NotNil(t, manager)
}

func TestNewRepoManager_WithGitRepo_NotGitRepo(t *testing.T) {
	// Arrange: 设置测试环境
	tempDir := t.TempDir()
	t.Setenv("HOME", tempDir)

	mockGitRepo := &mockGitRepository{
		repoPath:  tempDir,
		isGitRepo: false, // 不是 Git 仓库
	}

	// Act: 创建仓库配置管理器
	manager, err := NewRepoManager(mockGitRepo)

	// Assert: 应该返回错误
	assert.Error(t, err)
	assert.Nil(t, manager)
	assert.Contains(t, err.Error(), "不是 Git 仓库")
}

func TestNewRepoManager_WithoutGitRepo(t *testing.T) {
	// Arrange: 设置测试环境
	tempDir := t.TempDir()
	t.Setenv("HOME", tempDir)

	// 设置当前工作目录
	originalWd, err := os.Getwd()
	require.NoError(t, err)
	defer os.Chdir(originalWd)

	require.NoError(t, os.Chdir(tempDir))

	// Act: 创建仓库配置管理器（不提供 Git 接口）
	manager, err := NewRepoManager(nil)

	// Assert: 验证结果
	require.NoError(t, err)
	assert.NotNil(t, manager)
	assert.NotEmpty(t, manager.GetRepoID())
	// repo_id 应该以 "repo_" 开头（简单 ID）
	assert.Contains(t, manager.GetRepoID(), "repo_")
}

func TestNewRepoManager_WithGitRepo_OpenError(t *testing.T) {
	// Arrange: 设置测试环境
	tempDir := t.TempDir()
	t.Setenv("HOME", tempDir)

	mockGitRepo := &mockGitRepository{
		repoPath:  tempDir,
		isGitRepo: true,
		openError: fmt.Errorf("failed to open repository"),
	}

	// Act: 创建仓库配置管理器
	manager, err := NewRepoManager(mockGitRepo)

	// Assert: 应该返回错误
	assert.Error(t, err)
	assert.Nil(t, manager)
	assert.Contains(t, err.Error(), "生成仓库 ID 失败")
}

func TestNewRepoManager_WithGitRepo_NoRemoteURL(t *testing.T) {
	// Arrange: 设置测试环境
	tempDir := t.TempDir()
	t.Setenv("HOME", tempDir)

	mockGitRepo := &mockGitRepository{
		repoPath:  tempDir,
		isGitRepo: true,
		remoteURL: "", // 没有 remote URL
	}

	// Act: 创建仓库配置管理器
	manager, err := NewRepoManager(mockGitRepo)

	// Assert: 应该返回错误
	assert.Error(t, err)
	assert.Nil(t, manager)
	assert.Contains(t, err.Error(), "未找到 remote.origin.url")
}

// ==================== Load 测试 ====================

func TestRepoManager_Load_PublicConfigExists(t *testing.T) {
	// Arrange: 设置测试环境并创建公共配置文件
	tempDir := t.TempDir()
	t.Setenv("HOME", tempDir)

	configDir := filepath.Join(tempDir, ".workflow")
	publicConfigPath := filepath.Join(configDir, "config.toml")
	require.NoError(t, os.MkdirAll(configDir, 0755))

	configContent := `[template.commit]
format = "conventional"

[template.branch]
prefix = "feature/"
`
	require.NoError(t, os.WriteFile(publicConfigPath, []byte(configContent), 0644))

	mockGitRepo := &mockGitRepository{
		repoPath:  tempDir,
		isGitRepo: true,
		remoteURL: "https://github.com/owner/repo.git",
	}

	manager, err := NewRepoManager(mockGitRepo)
	require.NoError(t, err)

	// Act: 加载配置
	err = manager.Load()

	// Assert: 验证配置已加载
	assert.NoError(t, err)

	// 验证可以读取配置
	templateConfig := manager.GetTemplateConfig()
	assert.NotNil(t, templateConfig)
}

func TestRepoManager_Load_PublicConfigNotExists(t *testing.T) {
	// Arrange: 设置测试环境，不创建配置文件
	tempDir := t.TempDir()
	t.Setenv("HOME", tempDir)

	mockGitRepo := &mockGitRepository{
		repoPath:  tempDir,
		isGitRepo: true,
		remoteURL: "https://github.com/owner/repo.git",
	}

	manager, err := NewRepoManager(mockGitRepo)
	require.NoError(t, err)

	// Act: 加载配置（文件不存在）
	err = manager.Load()

	// Assert: 应该不返回错误（文件不存在是允许的）
	assert.NoError(t, err)
}

// ==================== GetTemplateConfig 测试 ====================

func TestRepoManager_GetTemplateConfig(t *testing.T) {
	// Arrange: 设置测试环境并创建配置文件
	tempDir := t.TempDir()
	t.Setenv("HOME", tempDir)

	configDir := filepath.Join(tempDir, ".workflow")
	publicConfigPath := filepath.Join(configDir, "config.toml")
	require.NoError(t, os.MkdirAll(configDir, 0755))

	configContent := `[template.commit]
format = "conventional"
type = "feat"

[template.branch]
prefix = "feature/"
pattern = "kebab-case"

[template.pull_requests]
title_format = "{{type}}: {{description}}"
`
	require.NoError(t, os.WriteFile(publicConfigPath, []byte(configContent), 0644))

	mockGitRepo := &mockGitRepository{
		repoPath:  tempDir,
		isGitRepo: true,
		remoteURL: "https://github.com/owner/repo.git",
	}

	manager, err := NewRepoManager(mockGitRepo)
	require.NoError(t, err)
	require.NoError(t, manager.Load())

	// Act: 获取模板配置
	templateConfig := manager.GetTemplateConfig()

	// Assert: 验证配置
	assert.NotNil(t, templateConfig)
	assert.NotEmpty(t, templateConfig.Commit)
	assert.Equal(t, "conventional", templateConfig.Commit["format"])
	assert.Equal(t, "feat", templateConfig.Commit["type"])
	assert.NotEmpty(t, templateConfig.Branch)
	assert.Equal(t, "feature/", templateConfig.Branch["prefix"])
	assert.NotEmpty(t, templateConfig.PullRequests)
}

func TestRepoManager_GetTemplateConfig_Empty(t *testing.T) {
	// Arrange: 设置测试环境，不创建配置文件
	tempDir := t.TempDir()
	t.Setenv("HOME", tempDir)

	mockGitRepo := &mockGitRepository{
		repoPath:  tempDir,
		isGitRepo: true,
		remoteURL: "https://github.com/owner/repo.git",
	}

	manager, err := NewRepoManager(mockGitRepo)
	require.NoError(t, err)

	// Act: 获取模板配置
	templateConfig := manager.GetTemplateConfig()

	// Assert: 验证返回空配置
	assert.NotNil(t, templateConfig)
	assert.Empty(t, templateConfig.Commit)
	assert.Empty(t, templateConfig.Branch)
	assert.Empty(t, templateConfig.PullRequests)
}

// ==================== GetBranchPrefix 测试 ====================

func TestRepoManager_GetBranchPrefix(t *testing.T) {
	// Arrange: 设置测试环境并创建私有配置文件
	tempDir := t.TempDir()
	t.Setenv("HOME", tempDir)

	privateConfigDir := filepath.Join(tempDir, ".workflow", "config")
	privateConfigPath := filepath.Join(privateConfigDir, "repository.toml")
	require.NoError(t, os.MkdirAll(privateConfigDir, 0755))

	mockGitRepo := &mockGitRepository{
		repoPath:  tempDir,
		isGitRepo: true,
		remoteURL: "https://github.com/owner/repo.git",
	}

	manager, err := NewRepoManager(mockGitRepo)
	require.NoError(t, err)
	repoID := manager.GetRepoID()

	// 创建私有配置文件
	configContent := fmt.Sprintf(`[%s.branch]
prefix = "my-prefix/"
`, repoID)
	require.NoError(t, os.WriteFile(privateConfigPath, []byte(configContent), 0644))

	// Act: 获取分支前缀
	prefix := manager.GetBranchPrefix()

	// Assert: 验证结果
	assert.Equal(t, "my-prefix/", prefix)
}

func TestRepoManager_GetBranchPrefix_NotConfigured(t *testing.T) {
	// Arrange: 设置测试环境，不创建私有配置文件
	tempDir := t.TempDir()
	t.Setenv("HOME", tempDir)

	mockGitRepo := &mockGitRepository{
		repoPath:  tempDir,
		isGitRepo: true,
		remoteURL: "https://github.com/owner/repo.git",
	}

	manager, err := NewRepoManager(mockGitRepo)
	require.NoError(t, err)

	// Act: 获取分支前缀
	prefix := manager.GetBranchPrefix()

	// Assert: 应该返回空字符串
	assert.Empty(t, prefix)
}

// ==================== GetIgnoreBranches 测试 ====================

func TestRepoManager_GetIgnoreBranches(t *testing.T) {
	// Arrange: 设置测试环境并创建私有配置文件
	tempDir := t.TempDir()
	t.Setenv("HOME", tempDir)

	privateConfigDir := filepath.Join(tempDir, ".workflow", "config")
	privateConfigPath := filepath.Join(privateConfigDir, "repository.toml")
	require.NoError(t, os.MkdirAll(privateConfigDir, 0755))

	mockGitRepo := &mockGitRepository{
		repoPath:  tempDir,
		isGitRepo: true,
		remoteURL: "https://github.com/owner/repo.git",
	}

	manager, err := NewRepoManager(mockGitRepo)
	require.NoError(t, err)
	repoID := manager.GetRepoID()

	// 创建私有配置文件
	configContent := fmt.Sprintf(`[%s.branch]
ignore = ["main", "develop", "master"]
`, repoID)
	require.NoError(t, os.WriteFile(privateConfigPath, []byte(configContent), 0644))

	// Act: 获取忽略分支列表
	ignoreBranches := manager.GetIgnoreBranches()

	// Assert: 验证结果
	assert.Len(t, ignoreBranches, 3)
	assert.Contains(t, ignoreBranches, "main")
	assert.Contains(t, ignoreBranches, "develop")
	assert.Contains(t, ignoreBranches, "master")
}

func TestRepoManager_GetIgnoreBranches_NotConfigured(t *testing.T) {
	// Arrange: 设置测试环境，不创建私有配置文件
	tempDir := t.TempDir()
	t.Setenv("HOME", tempDir)

	mockGitRepo := &mockGitRepository{
		repoPath:  tempDir,
		isGitRepo: true,
		remoteURL: "https://github.com/owner/repo.git",
	}

	manager, err := NewRepoManager(mockGitRepo)
	require.NoError(t, err)

	// Act: 获取忽略分支列表
	ignoreBranches := manager.GetIgnoreBranches()

	// Assert: 应该返回空切片
	assert.Empty(t, ignoreBranches)
}

// ==================== GetAutoAcceptChangeType 测试 ====================

func TestRepoManager_GetAutoAcceptChangeType(t *testing.T) {
	// Arrange: 设置测试环境并创建私有配置文件
	tempDir := t.TempDir()
	t.Setenv("HOME", tempDir)

	privateConfigDir := filepath.Join(tempDir, ".workflow", "config")
	privateConfigPath := filepath.Join(privateConfigDir, "repository.toml")
	require.NoError(t, os.MkdirAll(privateConfigDir, 0755))

	mockGitRepo := &mockGitRepository{
		repoPath:  tempDir,
		isGitRepo: true,
		remoteURL: "https://github.com/owner/repo.git",
	}

	manager, err := NewRepoManager(mockGitRepo)
	require.NoError(t, err)
	repoID := manager.GetRepoID()

	// 创建私有配置文件
	configContent := fmt.Sprintf(`[%s]
auto_accept_change_type = true
`, repoID)
	require.NoError(t, os.WriteFile(privateConfigPath, []byte(configContent), 0644))

	// Act: 获取自动接受变更类型设置
	autoAccept := manager.GetAutoAcceptChangeType()

	// Assert: 验证结果
	assert.True(t, autoAccept)
}

func TestRepoManager_GetAutoAcceptChangeType_False(t *testing.T) {
	// Arrange: 设置测试环境并创建私有配置文件
	tempDir := t.TempDir()
	t.Setenv("HOME", tempDir)

	privateConfigDir := filepath.Join(tempDir, ".workflow", "config")
	privateConfigPath := filepath.Join(privateConfigDir, "repository.toml")
	require.NoError(t, os.MkdirAll(privateConfigDir, 0755))

	mockGitRepo := &mockGitRepository{
		repoPath:  tempDir,
		isGitRepo: true,
		remoteURL: "https://github.com/owner/repo.git",
	}

	manager, err := NewRepoManager(mockGitRepo)
	require.NoError(t, err)
	repoID := manager.GetRepoID()

	// 创建私有配置文件
	configContent := fmt.Sprintf(`[%s]
auto_accept_change_type = false
`, repoID)
	require.NoError(t, os.WriteFile(privateConfigPath, []byte(configContent), 0644))

	// Act: 获取自动接受变更类型设置
	autoAccept := manager.GetAutoAcceptChangeType()

	// Assert: 验证结果
	assert.False(t, autoAccept)
}

func TestRepoManager_GetAutoAcceptChangeType_NotConfigured(t *testing.T) {
	// Arrange: 设置测试环境，不创建私有配置文件
	tempDir := t.TempDir()
	t.Setenv("HOME", tempDir)

	mockGitRepo := &mockGitRepository{
		repoPath:  tempDir,
		isGitRepo: true,
		remoteURL: "https://github.com/owner/repo.git",
	}

	manager, err := NewRepoManager(mockGitRepo)
	require.NoError(t, err)

	// Act: 获取自动接受变更类型设置
	autoAccept := manager.GetAutoAcceptChangeType()

	// Assert: 应该返回 false（默认值）
	assert.False(t, autoAccept)
}

// ==================== SaveTemplateConfig 测试 ====================

func TestRepoManager_SaveTemplateConfig(t *testing.T) {
	// Arrange: 设置测试环境
	tempDir := t.TempDir()
	t.Setenv("HOME", tempDir)

	mockGitRepo := &mockGitRepository{
		repoPath:  tempDir,
		isGitRepo: true,
		remoteURL: "https://github.com/owner/repo.git",
	}

	manager, err := NewRepoManager(mockGitRepo)
	require.NoError(t, err)

	templateConfig := &TemplateConfig{
		Commit: map[string]interface{}{
			"format": "conventional",
			"type":   "feat",
		},
		Branch: map[string]interface{}{
			"prefix": "feature/",
		},
		PullRequests: map[string]interface{}{
			"title_format": "{{type}}: {{description}}",
		},
	}

	// Act: 保存模板配置
	err = manager.SaveTemplateConfig(templateConfig)

	// Assert: 验证配置已保存
	assert.NoError(t, err)

	// 验证文件已创建
	publicConfigPath := manager.GetPublicConfigPath()
	_, err = os.Stat(publicConfigPath)
	assert.NoError(t, err)

	// 重新加载并验证
	require.NoError(t, manager.Load())
	loadedConfig := manager.GetTemplateConfig()
	assert.Equal(t, "conventional", loadedConfig.Commit["format"])
	assert.Equal(t, "feat", loadedConfig.Commit["type"])
	assert.Equal(t, "feature/", loadedConfig.Branch["prefix"])
}

func TestRepoManager_SaveTemplateConfig_OverwriteExisting(t *testing.T) {
	// Arrange: 设置测试环境并创建现有配置文件
	tempDir := t.TempDir()
	t.Setenv("HOME", tempDir)

	configDir := filepath.Join(tempDir, ".workflow")
	publicConfigPath := filepath.Join(configDir, "config.toml")
	require.NoError(t, os.MkdirAll(configDir, 0755))

	// 创建现有配置
	existingConfig := `[template.commit]
format = "old-format"
`
	require.NoError(t, os.WriteFile(publicConfigPath, []byte(existingConfig), 0644))

	mockGitRepo := &mockGitRepository{
		repoPath:  tempDir,
		isGitRepo: true,
		remoteURL: "https://github.com/owner/repo.git",
	}

	manager, err := NewRepoManager(mockGitRepo)
	require.NoError(t, err)

	templateConfig := &TemplateConfig{
		Commit: map[string]interface{}{
			"format": "new-format",
		},
	}

	// Act: 保存新配置
	err = manager.SaveTemplateConfig(templateConfig)

	// Assert: 验证配置已更新
	assert.NoError(t, err)

	// 重新加载并验证
	require.NoError(t, manager.Load())
	loadedConfig := manager.GetTemplateConfig()
	assert.Equal(t, "new-format", loadedConfig.Commit["format"])
}

// ==================== 辅助函数测试 ====================

func TestExtractRepoNameFromURL(t *testing.T) {
	tests := []struct {
		name     string
		url      string
		expected string
	}{
		{
			name:     "HTTPS URL with .git",
			url:      "https://github.com/owner/repo.git",
			expected: "repo",
		},
		{
			name:     "HTTPS URL without .git",
			url:      "https://github.com/owner/repo",
			expected: "repo",
		},
		{
			name:     "SSH URL with .git",
			url:      "git@github.com:owner/repo.git",
			expected: "repo",
		},
		{
			name:     "SSH URL without .git",
			url:      "git@github.com:owner/repo",
			expected: "repo",
		},
		{
			name:     "URL with nested path",
			url:      "https://github.com/org/group/repo.git",
			expected: "repo",
		},
		{
			name:     "Invalid URL",
			url:      "invalid-url",
			expected: "invalid-url", // 实际实现会返回最后一个部分
		},
		{
			name:     "Empty URL",
			url:      "",
			expected: "", // 实际实现会返回空字符串
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := extractRepoNameFromURL(tt.url)
			assert.Equal(t, tt.expected, result)
		})
	}
}


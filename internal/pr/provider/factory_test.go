package provider

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ==================== NewPlatformProvider 测试 ====================

func TestNewPlatformProvider(t *testing.T) {
	tests := []struct {
		name     string
		platform string
		token    string
		owner    string
		repo     string
		wantErr  bool
	}{
		{
			name:     "GitHub 平台",
			platform: "github",
			token:    "test-token",
			owner:    "owner",
			repo:     "repo",
			wantErr:  false,
		},
		{
			name:     "不支持的平台",
			platform: "gitlab",
			token:    "test-token",
			owner:    "owner",
			repo:     "repo",
			wantErr:  true,
		},
		{
			name:     "空平台名称",
			platform: "",
			token:    "test-token",
			owner:    "owner",
			repo:     "repo",
			wantErr:  true,
		},
		{
			name:     "无效平台名称",
			platform: "invalid",
			token:    "test-token",
			owner:    "owner",
			repo:     "repo",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			provider, err := NewPlatformProvider(tt.platform, tt.token, tt.owner, tt.repo)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, provider)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, provider)
				assert.Equal(t, tt.platform, provider.GetPlatformName())
			}
		})
	}
}

// ==================== AutoDetectPlatform 测试 ====================

func TestAutoDetectPlatform(t *testing.T) {
	// 注意：当前实现默认返回 "github"
	// 未来实现自动检测逻辑后，需要更新此测试
	platform, err := AutoDetectPlatform()

	assert.NoError(t, err)
	assert.Equal(t, "github", platform)
}

// ==================== NewPlatformProviderAuto 测试 ====================

func TestNewPlatformProviderAuto(t *testing.T) {
	tests := []struct {
		name    string
		token   string
		owner   string
		repo    string
		wantErr bool
	}{
		{
			name:    "有效配置",
			token:   "test-token",
			owner:   "owner",
			repo:    "repo",
			wantErr: false,
		},
		{
			name:    "空 token",
			token:   "",
			owner:   "owner",
			repo:    "repo",
			wantErr: true, // GitHub 创建时会验证 token
		},
		{
			name:    "空 owner",
			token:   "test-token",
			owner:   "",
			repo:    "repo",
			wantErr: true, // GitHub 创建时会验证 owner
		},
		{
			name:    "空 repo",
			token:   "test-token",
			owner:   "owner",
			repo:    "",
			wantErr: true, // GitHub 创建时会验证 repo
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			provider, err := NewPlatformProviderAuto(tt.token, tt.owner, tt.repo)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, provider)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, provider)
				// 当前实现默认检测为 GitHub
				assert.Equal(t, "github", provider.GetPlatformName())
			}
		})
	}
}

// ==================== 集成测试 ====================

func TestNewPlatformProvider_ThenGetPlatformName(t *testing.T) {
	// 测试创建提供者后获取平台名称
	provider, err := NewPlatformProvider("github", "test-token", "owner", "repo")
	require.NoError(t, err)
	require.NotNil(t, provider)

	platformName := provider.GetPlatformName()
	assert.Equal(t, "github", platformName)
}

func TestNewPlatformProviderAuto_Integration(t *testing.T) {
	// 测试自动检测和创建提供者的完整流程
	provider, err := NewPlatformProviderAuto("test-token", "owner", "repo")
	require.NoError(t, err)
	require.NotNil(t, provider)

	// 验证提供者类型
	assert.Equal(t, "github", provider.GetPlatformName())
}

func TestNewPlatformProviderAuto_AutoDetectFails(t *testing.T) {
	// 注意：当前 AutoDetectPlatform 不会失败
	// 如果未来实现会失败，这个测试会验证错误处理
	// 目前这个测试主要验证 AutoDetectPlatform 返回错误时的处理
	// 由于当前实现不会失败，这个测试主要作为占位符
	provider, err := NewPlatformProviderAuto("test-token", "owner", "repo")

	// 当前实现不会因为检测失败而返回错误
	// 只有在创建 GitHub 客户端时才会失败
	if err != nil {
		assert.Contains(t, err.Error(), "failed to detect platform")
		assert.Nil(t, provider)
	} else {
		assert.NotNil(t, provider)
	}
}

// ==================== 错误消息测试 ====================

func TestNewPlatformProvider_ErrorMessages(t *testing.T) {
	t.Run("不支持的平台错误消息", func(t *testing.T) {
		_, err := NewPlatformProvider("gitlab", "test-token", "owner", "repo")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "unsupported platform")
		assert.Contains(t, err.Error(), "gitlab")
	})

	t.Run("自动检测失败错误消息", func(t *testing.T) {
		// 注意：当前 AutoDetectPlatform 不会失败
		// 如果未来实现会失败，需要更新此测试
		_, err := NewPlatformProviderAuto("test-token", "owner", "repo")
		// 当前实现不会因为检测失败而返回错误
		// 只有在创建 GitHub 客户端时才会失败
		if err != nil {
			// 如果返回错误，应该是创建 GitHub 客户端的错误
			assert.Contains(t, err.Error(), "failed to detect platform")
		}
	})
}


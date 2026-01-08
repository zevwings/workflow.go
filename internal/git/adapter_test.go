package git

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ==================== NewGitAdapter 测试 ====================

func TestNewGitAdapter(t *testing.T) {
	tests := []struct {
		name     string
		repoPath string
	}{
		{
			name:     "with repo path",
			repoPath: "/path/to/repo",
		},
		{
			name:     "empty repo path",
			repoPath: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			adapter := NewGitAdapter(tt.repoPath)
			assert.NotNil(t, adapter)
			assert.Equal(t, tt.repoPath, adapter.GetRepoPath())
		})
	}
}

// ==================== GetRepoPath 测试 ====================

func TestGitAdapter_GetRepoPath(t *testing.T) {
	repoPath := "/path/to/repo"
	adapter := NewGitAdapter(repoPath)

	result := adapter.GetRepoPath()
	assert.Equal(t, repoPath, result)
}

// ==================== IsGitRepo 测试 ====================

func TestGitAdapter_IsGitRepo(t *testing.T) {
	tests := []struct {
		name  string
		setup func(t *testing.T) string
		want  bool
	}{
		{
			name: "valid git repo",
			setup: func(t *testing.T) string {
				_, tempDir := setupTestRepo(t)
				return tempDir
			},
			want: true,
		},
		{
			name: "not a git repo",
			setup: func(t *testing.T) string {
				return t.TempDir()
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			adapter := NewGitAdapter("")
			path := tt.setup(t)

			result := adapter.IsGitRepo(path)
			assert.Equal(t, tt.want, result)
		})
	}
}

// ==================== Open 测试 ====================

func TestGitAdapter_Open(t *testing.T) {
	tests := []struct {
		name    string
		setup   func(t *testing.T) string
		wantErr bool
	}{
		{
			name: "valid git repo",
			setup: func(t *testing.T) string {
				_, tempDir := setupTestRepo(t)
				return tempDir
			},
			wantErr: false,
		},
		{
			name: "not a git repo",
			setup: func(t *testing.T) string {
				return t.TempDir()
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			adapter := NewGitAdapter("")
			path := tt.setup(t)

			repoAdapter, err := adapter.Open(path)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, repoAdapter)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, repoAdapter)
			}
		})
	}
}

// ==================== GitRepoAdapter 测试 ====================

func TestGitRepoAdapter_GetRemoteURL(t *testing.T) {
	repo, _ := setupTestRepoWithCommit(t)

	// 添加远程
	remoteName := "origin"
	remoteURL := "https://github.com/user/repo.git"
	err := repo.AddRemote(remoteName, remoteURL)
	require.NoError(t, err)

	// 创建适配器
	adapter := NewGitAdapter("")
	repoAdapter, err := adapter.Open(repo.Path())
	require.NoError(t, err)

	// 获取远程 URL
	url, err := repoAdapter.GetRemoteURL(remoteName)
	assert.NoError(t, err)
	assert.Equal(t, remoteURL, url)
}

func TestGitRepoAdapter_GetRemoteURL_NonExistent(t *testing.T) {
	repo, _ := setupTestRepoWithCommit(t)

	// 创建适配器
	adapter := NewGitAdapter("")
	repoAdapter, err := adapter.Open(repo.Path())
	require.NoError(t, err)

	// 获取不存在的远程 URL
	url, err := repoAdapter.GetRemoteURL("non-existent")
	assert.Error(t, err)
	assert.Empty(t, url)
}

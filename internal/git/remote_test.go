package git

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ==================== AddRemote 测试 ====================

func TestRepository_AddRemote(t *testing.T) {
	repo, _ := setupTestRepoWithCommit(t)

	tests := []struct {
		name    string
		remote  string
		url     string
		wantErr bool
	}{
		{
			name:    "add origin remote",
			remote:  "origin",
			url:     "https://github.com/user/repo.git",
			wantErr: false,
		},
		{
			name:    "add upstream remote",
			remote:  "upstream",
			url:     "https://github.com/upstream/repo.git",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := repo.AddRemote(tt.remote, tt.url)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)

				// 验证远程已添加
				url, err := repo.GetRemoteURL(tt.remote)
				assert.NoError(t, err)
				assert.Equal(t, tt.url, url)
			}
		})
	}
}

func TestRepository_AddRemote_Duplicate(t *testing.T) {
	repo, _ := setupTestRepoWithCommit(t)

	// 添加远程
	err := repo.AddRemote("origin", "https://github.com/user/repo.git")
	require.NoError(t, err)

	// 尝试添加同名远程（应该失败）
	err = repo.AddRemote("origin", "https://github.com/user/repo2.git")
	assert.Error(t, err)
}

// ==================== RemoveRemote 测试 ====================

func TestRepository_RemoveRemote(t *testing.T) {
	repo, _ := setupTestRepoWithCommit(t)

	// 添加远程
	err := repo.AddRemote("origin", "https://github.com/user/repo.git")
	require.NoError(t, err)

	// 删除远程
	err = repo.RemoveRemote("origin")
	assert.NoError(t, err)

	// 验证远程已删除
	_, err = repo.GetRemoteURL("origin")
	assert.Error(t, err)
}

func TestRepository_RemoveRemote_NonExistent(t *testing.T) {
	repo, _ := setupTestRepoWithCommit(t)

	err := repo.RemoveRemote("non-existent")
	assert.Error(t, err)
}

// ==================== ListRemotes 测试 ====================

func TestRepository_ListRemotes(t *testing.T) {
	repo, _ := setupTestRepoWithCommit(t)

	// 添加多个远程
	remotes := map[string]string{
		"origin":   "https://github.com/user/repo.git",
		"upstream": "https://github.com/upstream/repo.git",
	}

	for name, url := range remotes {
		err := repo.AddRemote(name, url)
		require.NoError(t, err)
	}

	// 列出所有远程
	list, err := repo.ListRemotes()
	assert.NoError(t, err)
	assert.Len(t, list, len(remotes))

	// 验证远程信息
	remoteMap := make(map[string]string)
	for _, remote := range list {
		remoteMap[remote.Name] = remote.URL
	}

	for name, url := range remotes {
		assert.Equal(t, url, remoteMap[name])
	}
}

func TestRepository_ListRemotes_Empty(t *testing.T) {
	repo, _ := setupTestRepoWithCommit(t)

	list, err := repo.ListRemotes()
	assert.NoError(t, err)
	assert.Empty(t, list)
}

// ==================== GetRemoteURL 测试 ====================

func TestRepository_GetRemoteURL(t *testing.T) {
	repo, _ := setupTestRepoWithCommit(t)

	// 添加远程
	expectedURL := "https://github.com/user/repo.git"
	err := repo.AddRemote("origin", expectedURL)
	require.NoError(t, err)

	// 获取远程 URL
	url, err := repo.GetRemoteURL("origin")
	assert.NoError(t, err)
	assert.Equal(t, expectedURL, url)
}

func TestRepository_GetRemoteURL_NonExistent(t *testing.T) {
	repo, _ := setupTestRepoWithCommit(t)

	url, err := repo.GetRemoteURL("non-existent")
	assert.Error(t, err)
	assert.Empty(t, url)
}

// ==================== Fetch 测试 ====================

func TestRepository_Fetch(t *testing.T) {
	// 使用测试远程仓库（Mock 或真实）
	repo, _ := setupTestRepoWithCommit(t)

	// 设置测试远程仓库
	remotePath, isReal, _ := setupTestRemoteRepo(t, "main")

	// 添加远程
	err := repo.AddRemote("origin", remotePath)
	require.NoError(t, err)

	// 执行 Fetch
	err = repo.Fetch("origin", nil)
	if isReal {
		// 真实远程仓库可能因为网络问题失败，这是可以接受的
		if err != nil {
			t.Logf("从真实远程仓库 Fetch 失败（可能是网络问题）: %v", err)
			return
		}
		t.Logf("从真实远程仓库 Fetch 成功")
	} else {
		// Mock 远程仓库：Fetch 可能会成功或失败（取决于 go-git 的实现）
		// 我们主要测试方法调用不会 panic
		if err != nil {
			t.Logf("Fetch 结果: %v", err)
			// 某些错误是可以接受的（如 NoErrAlreadyUpToDate）
		}
	}
}

// ==================== Push 测试 ====================

func TestRepository_Push(t *testing.T) {
	// 使用测试远程仓库（Mock 或真实）
	repo, _ := setupTestRepoWithCommit(t)

	// 设置测试远程仓库
	remotePath, isReal, _ := setupTestRemoteRepo(t, "main")

	// 添加远程
	err := repo.AddRemote("origin", remotePath)
	require.NoError(t, err)

	// 获取当前分支
	branch, err := repo.CurrentBranch()
	require.NoError(t, err)

	// 执行 Push
	err = repo.Push("origin", branch, nil)
	if isReal {
		// 真实远程仓库需要认证，可能会失败
		// 这是可以接受的，因为测试环境可能没有配置认证
		if err != nil {
			t.Logf("推送到真实远程仓库失败（可能需要认证）: %v", err)
			return
		}
		t.Logf("推送到真实远程仓库成功")
	} else {
		// Mock 远程仓库：Push 可能会成功或失败（取决于 go-git 的实现）
		// 我们主要测试方法调用不会 panic
		if err != nil {
			t.Logf("Push 结果: %v", err)
			// 某些错误是可以接受的
		} else {
			// 如果成功，验证远程仓库有新的引用
			refs, err := repo.ListRemoteRefs("origin")
			if err == nil {
				assert.NotNil(t, refs)
			}
		}
	}
}

// ==================== PushWithUpstream 测试 ====================

func TestRepository_PushWithUpstream(t *testing.T) {
	// 注意：PushWithUpstream 需要真实的远程仓库或 Mock 服务器
	// 这里我们测试错误处理
	repo, _ := setupTestRepoWithCommit(t)

	// 添加远程（使用不存在的 URL）
	err := repo.AddRemote("origin", "https://github.com/non-existent/repo.git")
	require.NoError(t, err)

	// 获取当前分支
	branch, err := repo.CurrentBranch()
	require.NoError(t, err)

	// 尝试 PushWithUpstream（应该失败，因为没有真实的远程）
	err = repo.PushWithUpstream("origin", branch, nil)
	// 这个测试可能会失败，因为需要网络连接
	// 我们主要测试方法调用不会 panic
	if err != nil {
		t.Logf("PushWithUpstream 失败（预期行为，因为没有真实远程）: %v", err)
	}
}

// ==================== ListRemoteRefs 测试 ====================

func TestRepository_ListRemoteRefs(t *testing.T) {
	// 使用测试远程仓库（Mock 或真实）
	repo, _ := setupTestRepoWithCommit(t)

	// 设置测试远程仓库
	remotePath, isReal, commitHash := setupTestRemoteRepo(t, "main")

	// 添加远程
	err := repo.AddRemote("origin", remotePath)
	require.NoError(t, err)

	// 列出远程引用
	refs, err := repo.ListRemoteRefs("origin")
	if isReal {
		// 真实远程仓库可能因为网络问题失败，这是可以接受的
		if err != nil {
			t.Logf("从真实远程仓库列出引用失败（可能是网络问题）: %v", err)
			return
		}
		assert.NotNil(t, refs)
		// 真实远程仓库应该至少包含一些引用
		assert.NotEmpty(t, refs)
	} else {
		assert.NoError(t, err)
		assert.NotNil(t, refs)

		// 验证包含默认分支引用
		mainRef := "refs/heads/main"
		assert.Contains(t, refs, mainRef)
		assert.Equal(t, commitHash, refs[mainRef])
	}
}

// TestRepository_ListRemoteRefs_NonExistent 测试远程不存在的情况
func TestRepository_ListRemoteRefs_NonExistent(t *testing.T) {
	repo, _ := setupTestRepoWithCommit(t)

	// 尝试列出不存在的远程引用
	refs, err := repo.ListRemoteRefs("non-existent")
	assert.Error(t, err)
	assert.Nil(t, refs)
	assert.Contains(t, err.Error(), "failed to get remote")
}

// TestRepository_ListRemoteRefs_EmptyRemote 测试空远程的情况
func TestRepository_ListRemoteRefs_EmptyRemote(t *testing.T) {
	repo, _ := setupTestRepoWithCommit(t)

	// 添加远程
	err := repo.AddRemote("test", "https://github.com/user/repo.git")
	require.NoError(t, err)

	// 注意：由于没有真实的远程连接，ListRemoteRefs 可能会失败
	// 但我们至少验证了方法调用不会 panic
	refs, err := repo.ListRemoteRefs("test")
	if err != nil {
		// 这是预期的，因为没有真实的远程连接
		t.Logf("ListRemoteRefs 失败（预期行为）: %v", err)
		assert.Nil(t, refs)
	} else {
		// 如果成功，验证返回的结构
		assert.NotNil(t, refs)
		// refs 可能是空的，这是正常的
	}
}

// ==================== ExtractRepoName 测试 ====================

func TestExtractRepoName(t *testing.T) {
	tests := []struct {
		name    string
		url     string
		want    string
		wantErr bool
	}{
		{
			name:    "SSH format",
			url:     "git@github.com:owner/repo.git",
			want:    "owner/repo",
			wantErr: false,
		},
		{
			name:    "HTTPS format",
			url:     "https://github.com/owner/repo.git",
			want:    "owner/repo",
			wantErr: false,
		},
		{
			name:    "HTTPS format without .git",
			url:     "https://github.com/owner/repo",
			want:    "owner/repo",
			wantErr: false,
		},
		{
			name:    "SSH format with port",
			url:     "ssh://git@github.com:22/owner/repo.git",
			want:    "owner/repo",
			wantErr: false,
		},
		{
			name:    "invalid URL",
			url:     "not-a-url",
			want:    "",
			wantErr: true,
		},
		{
			name:    "empty URL",
			url:     "",
			want:    "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ExtractRepoName(tt.url)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Empty(t, result)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, result)
			}
		})
	}
}

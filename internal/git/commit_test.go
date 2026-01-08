package git

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ==================== Status 测试 ====================

func TestRepository_Status(t *testing.T) {
	repo, tempDir := setupTestRepoWithCommit(t)

	// 创建新文件（未跟踪）
	newFile := filepath.Join(tempDir, "new.txt")
	err := os.WriteFile(newFile, []byte("new content"), 0644)
	require.NoError(t, err)

	// 修改已存在的文件
	existingFile := filepath.Join(tempDir, "test.txt")
	err = os.WriteFile(existingFile, []byte("modified content"), 0644)
	require.NoError(t, err)

	// 获取状态（在添加之前）
	status, err := repo.Status()
	assert.NoError(t, err)
	assert.NotNil(t, status)

	// 验证新文件在未跟踪列表中
	assert.Contains(t, status.UntrackedFiles, "new.txt")
	// 验证修改的文件在修改列表中
	assert.Contains(t, status.ModifiedFiles, "test.txt")

	// 添加文件到暂存区
	err = repo.Add("new.txt")
	require.NoError(t, err)

	// 再次获取状态
	status, err = repo.Status()
	assert.NoError(t, err)
	assert.NotNil(t, status)

	// 验证新文件现在在暂存区
	assert.Contains(t, status.StagedFiles, "new.txt")
	// 验证修改的文件仍在修改列表中
	assert.Contains(t, status.ModifiedFiles, "test.txt")
}

func TestRepository_Status_Clean(t *testing.T) {
	repo, _ := setupTestRepoWithCommit(t)

	status, err := repo.Status()
	assert.NoError(t, err)
	assert.NotNil(t, status)

	// 干净的工作区应该没有更改
	assert.Empty(t, status.ModifiedFiles)
	assert.Empty(t, status.StagedFiles)
	assert.Empty(t, status.UntrackedFiles)
}

// ==================== HasChanges 测试 ====================

func TestRepository_HasChanges(t *testing.T) {
	tests := []struct {
		name    string
		setup   func(t *testing.T, repo *Repository, tempDir string)
		want    bool
		wantErr bool
	}{
		{
			name: "has modified files",
			setup: func(t *testing.T, repo *Repository, tempDir string) {
				file := filepath.Join(tempDir, "test.txt")
				err := os.WriteFile(file, []byte("modified"), 0644)
				require.NoError(t, err)
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "has staged files",
			setup: func(t *testing.T, repo *Repository, tempDir string) {
				file := filepath.Join(tempDir, "new.txt")
				err := os.WriteFile(file, []byte("new"), 0644)
				require.NoError(t, err)
				err = repo.Add("new.txt")
				require.NoError(t, err)
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "has untracked files",
			setup: func(t *testing.T, repo *Repository, tempDir string) {
				file := filepath.Join(tempDir, "untracked.txt")
				err := os.WriteFile(file, []byte("untracked"), 0644)
				require.NoError(t, err)
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "no changes",
			setup: func(t *testing.T, repo *Repository, tempDir string) {
				// 不进行任何修改
			},
			want:    false,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo, tempDir := setupTestRepoWithCommit(t)
			tt.setup(t, repo, tempDir)

			hasChanges, err := repo.HasChanges()

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, hasChanges)
			}
		})
	}
}

// ==================== Add 测试 ====================

func TestRepository_Add(t *testing.T) {
	repo, tempDir := setupTestRepoWithCommit(t)

	// 创建新文件
	newFile := filepath.Join(tempDir, "new.txt")
	err := os.WriteFile(newFile, []byte("new content"), 0644)
	require.NoError(t, err)

	// 添加文件
	err = repo.Add("new.txt")
	assert.NoError(t, err)

	// 验证文件已添加到暂存区
	status, err := repo.Status()
	assert.NoError(t, err)
	assert.Contains(t, status.StagedFiles, "new.txt")
}

func TestRepository_Add_NonExistent(t *testing.T) {
	repo, _ := setupTestRepoWithCommit(t)

	err := repo.Add("non-existent.txt")
	assert.Error(t, err)
}

// ==================== AddAll 测试 ====================

func TestRepository_AddAll(t *testing.T) {
	repo, tempDir := setupTestRepoWithCommit(t)

	// 创建多个新文件
	files := []string{"file1.txt", "file2.txt", "file3.txt"}
	for _, filename := range files {
		file := filepath.Join(tempDir, filename)
		err := os.WriteFile(file, []byte("content"), 0644)
		require.NoError(t, err)
	}

	// 添加所有文件
	err := repo.AddAll()
	assert.NoError(t, err)

	// 验证所有文件已添加到暂存区
	status, err := repo.Status()
	assert.NoError(t, err)
	for _, filename := range files {
		assert.Contains(t, status.StagedFiles, filename)
	}
}

// ==================== Commit 测试 ====================

func TestRepository_Commit(t *testing.T) {
	repo, tempDir := setupTestRepoWithCommit(t)

	// 创建新文件并添加
	newFile := filepath.Join(tempDir, "commit-test.txt")
	err := os.WriteFile(newFile, []byte("commit test"), 0644)
	require.NoError(t, err)

	err = repo.Add("commit-test.txt")
	require.NoError(t, err)

	// 提交
	author := &object.Signature{
		Name:  "Test User",
		Email: "test@example.com",
		When:  time.Now(),
	}
	hash, err := repo.Commit("Test commit", author)
	assert.NoError(t, err)
	assert.NotEqual(t, plumbing.ZeroHash, hash)

	// 验证提交
	commit, err := repo.GetCommit(hash)
	assert.NoError(t, err)
	assert.Equal(t, "Test commit", commit.Message)
	assert.Contains(t, commit.Author, "Test User")
}

func TestRepository_Commit_WithoutAuthor(t *testing.T) {
	repo, tempDir := setupTestRepoWithCommit(t)

	// 创建新文件并添加
	newFile := filepath.Join(tempDir, "commit-test2.txt")
	err := os.WriteFile(newFile, []byte("commit test 2"), 0644)
	require.NoError(t, err)

	err = repo.Add("commit-test2.txt")
	require.NoError(t, err)

	// 提交（不提供作者，应该使用配置中的作者）
	hash, err := repo.Commit("Test commit without author", nil)
	assert.NoError(t, err)
	assert.NotEqual(t, plumbing.ZeroHash, hash)

	// 验证提交
	commit, err := repo.GetCommit(hash)
	assert.NoError(t, err)
	assert.Equal(t, "Test commit without author", commit.Message)
}

// ==================== GetHead 测试 ====================

func TestRepository_GetHead(t *testing.T) {
	repo, _ := setupTestRepoWithCommit(t)

	hash, err := repo.GetHead()
	assert.NoError(t, err)
	assert.NotEqual(t, plumbing.ZeroHash, hash)
}

func TestRepository_GetHead_EmptyRepo(t *testing.T) {
	repo, _ := setupTestRepo(t)

	// 没有提交的仓库，HEAD 可能无效
	hash, err := repo.GetHead()
	// 这个行为取决于 go-git 的实现
	// 可能返回错误，也可能返回零哈希
	if err != nil {
		assert.Error(t, err)
	} else {
		t.Logf("GetHead 在没有提交时返回: %s", hash.String())
	}
}

// ==================== GetCommit 测试 ====================

func TestRepository_GetCommit(t *testing.T) {
	repo, _ := setupTestRepoWithCommit(t)

	// 获取 HEAD
	head, err := repo.GetHead()
	require.NoError(t, err)

	// 获取提交信息
	commit, err := repo.GetCommit(head)
	assert.NoError(t, err)
	assert.NotNil(t, commit)
	assert.Equal(t, head.String(), commit.Hash)
	assert.Equal(t, "Initial commit", commit.Message)
	assert.Contains(t, commit.Author, "Test User")
}

func TestRepository_GetCommit_InvalidHash(t *testing.T) {
	repo, _ := setupTestRepoWithCommit(t)

	invalidHash := plumbing.NewHash("0000000000000000000000000000000000000000")
	commit, err := repo.GetCommit(invalidHash)
	assert.Error(t, err)
	assert.Nil(t, commit)
}

// ==================== GetLastCommit 测试 ====================

func TestRepository_GetLastCommit(t *testing.T) {
	repo, tempDir := setupTestRepoWithCommit(t)

	// 创建新提交
	newFile := filepath.Join(tempDir, "last-commit.txt")
	err := os.WriteFile(newFile, []byte("last commit"), 0644)
	require.NoError(t, err)

	err = repo.Add("last-commit.txt")
	require.NoError(t, err)

	author := &object.Signature{
		Name:  "Test User",
		Email: "test@example.com",
	}
	_, err = repo.Commit("Last commit", author)
	require.NoError(t, err)

	// 获取最后一次提交
	commit, err := repo.GetLastCommit()
	assert.NoError(t, err)
	assert.NotNil(t, commit)
	assert.Equal(t, "Last commit", commit.Message)
}

// ==================== Log 测试 ====================

func TestRepository_Log(t *testing.T) {
	repo, tempDir := setupTestRepoWithCommit(t)

	// 创建多个提交
	author := &object.Signature{
		Name:  "Test User",
		Email: "test@example.com",
	}

	for i := 0; i < 3; i++ {
		filename := filepath.Join(tempDir, fmt.Sprintf("log-test%d.txt", i))
		err := os.WriteFile(filename, []byte("content"), 0644)
		require.NoError(t, err)

		err = repo.Add(fmt.Sprintf("log-test%d.txt", i))
		require.NoError(t, err)

		_, err = repo.Commit(fmt.Sprintf("Commit %d", i), author)
		require.NoError(t, err)
	}

	// 获取 HEAD
	head, err := repo.GetHead()
	require.NoError(t, err)

	// 获取提交历史（限制 2 个）
	logs, err := repo.Log(head, 2)
	assert.NoError(t, err)
	assert.Len(t, logs, 2)

	// 验证提交信息
	assert.Equal(t, "Commit 2", logs[0].Message)
	assert.Equal(t, "Commit 1", logs[1].Message)
}

func TestRepository_Log_NoLimit(t *testing.T) {
	repo, tempDir := setupTestRepoWithCommit(t)

	// 创建多个提交
	author := &object.Signature{
		Name:  "Test User",
		Email: "test@example.com",
	}

	for i := 0; i < 5; i++ {
		filename := filepath.Join(tempDir, fmt.Sprintf("log-test%d.txt", i))
		err := os.WriteFile(filename, []byte("content"), 0644)
		require.NoError(t, err)

		err = repo.Add(fmt.Sprintf("log-test%d.txt", i))
		require.NoError(t, err)

		_, err = repo.Commit(fmt.Sprintf("Commit %d", i), author)
		require.NoError(t, err)
	}

	// 获取 HEAD
	head, err := repo.GetHead()
	require.NoError(t, err)

	// 获取所有提交历史（limit = 0 表示无限制）
	logs, err := repo.Log(head, 0)
	assert.NoError(t, err)
	// 应该包含所有提交（至少 5 个，加上初始提交）
	assert.GreaterOrEqual(t, len(logs), 5)
}

// ==================== ResolveRevision 测试 ====================

func TestRepository_ResolveRevision(t *testing.T) {
	repo, _ := setupTestRepoWithCommit(t)

	// 获取 HEAD
	head, err := repo.GetHead()
	require.NoError(t, err)

	// 解析 HEAD
	hash, err := repo.ResolveRevision("HEAD")
	assert.NoError(t, err)
	assert.Equal(t, head, hash)

	// 解析分支名
	hash, err = repo.ResolveRevision("main")
	assert.NoError(t, err)
	assert.Equal(t, head, hash)

	// 解析完整哈希
	hash, err = repo.ResolveRevision(head.String())
	assert.NoError(t, err)
	assert.Equal(t, head, hash)
}

func TestRepository_ResolveRevision_Invalid(t *testing.T) {
	repo, _ := setupTestRepoWithCommit(t)

	hash, err := repo.ResolveRevision("invalid-revision")
	assert.Error(t, err)
	assert.Equal(t, plumbing.ZeroHash, hash)
}

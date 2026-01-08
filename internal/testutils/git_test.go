//go:build test

package testutils

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestGitTestRepo_Basic 测试基本的 Git 仓库创建
func TestGitTestRepo_Basic(t *testing.T) {
	repo := NewGitTestRepo().
		WithDefaultBranch("main").
		Build(t)

	assert.NotNil(t, repo)
	assert.NotEmpty(t, repo.Path())
	assert.False(t, repo.IsBare())
	assert.NotNil(t, repo.Repository())
}

// TestGitTestRepo_WithFile 测试带文件的 Git 仓库
func TestGitTestRepo_WithFile(t *testing.T) {
	repo := NewGitTestRepo().
		WithFileString("test.txt", "test content").
		Build(t)

	assert.NotNil(t, repo)

	// 验证文件存在
	repoPath := repo.Repository()
	assert.NotNil(t, repoPath)
}

// TestGitTestRepo_WithCommit 测试带提交的 Git 仓库
func TestGitTestRepo_WithCommit(t *testing.T) {
	repo := NewGitTestRepo().
		WithFileString("test.txt", "test content").
		WithCommit("Initial commit").
		Build(t)

	assert.NotNil(t, repo)

	// 验证有提交
	head, err := repo.Repository().GetHead()
	require.NoError(t, err)
	assert.NotEmpty(t, head.String())
}

// TestGitTestRepo_WithBranch 测试带分支的 Git 仓库
func TestGitTestRepo_WithBranch(t *testing.T) {
	repo := NewGitTestRepo().
		WithFileString("test.txt", "test content").
		WithCommit("Initial commit").
		WithBranch("feature/test").
		Build(t)

	assert.NotNil(t, repo)

	// 验证分支存在
	exists, err := repo.Repository().BranchExists("feature/test")
	require.NoError(t, err)
	assert.True(t, exists)
}

// TestGitTestRepo_WithBranchAndCheckout 测试创建并切换到分支
func TestGitTestRepo_WithBranchAndCheckout(t *testing.T) {
	repo := NewGitTestRepo().
		WithFileString("test.txt", "test content").
		WithCommit("Initial commit").
		WithBranchAndCheckout("feature/test").
		Build(t)

	assert.NotNil(t, repo)

	// 验证当前分支
	currentBranch, err := repo.Repository().CurrentBranch()
	require.NoError(t, err)
	assert.Equal(t, "feature/test", currentBranch)
}

// TestGitTestRepo_WithRemote 测试带远程仓库的 Git 仓库
func TestGitTestRepo_WithRemote(t *testing.T) {
	// 先创建一个 bare 仓库作为远程
	remoteRepo := NewBareGitTestRepo().
		Build(t)

	// 创建本地仓库并添加远程
	repo := NewGitTestRepo().
		WithFileString("test.txt", "test content").
		WithCommit("Initial commit").
		WithRemote("origin", remoteRepo.Path()).
		Build(t)

	assert.NotNil(t, repo)

	// 验证远程存在
	url, err := repo.Repository().GetRemoteURL("origin")
	require.NoError(t, err)
	assert.Equal(t, remoteRepo.Path(), url)
}

// TestGitTestRepo_WithTag 测试带 Tag 的 Git 仓库
func TestGitTestRepo_WithTag(t *testing.T) {
	repo := NewGitTestRepo().
		WithFileString("test.txt", "test content").
		WithCommit("Initial commit").
		WithTag("v1.0.0").
		Build(t)

	assert.NotNil(t, repo)

	// 验证 tag 存在
	exists, err := repo.Repository().TagExists("v1.0.0")
	require.NoError(t, err)
	assert.True(t, exists)
}

// TestGitTestRepo_WithTagMessage 测试带消息的 Tag
func TestGitTestRepo_WithTagMessage(t *testing.T) {
	repo := NewGitTestRepo().
		WithFileString("test.txt", "test content").
		WithCommit("Initial commit").
		WithTagMessage("v1.0.0", "Release version 1.0.0").
		Build(t)

	assert.NotNil(t, repo)

	// 验证 tag 存在
	exists, err := repo.Repository().TagExists("v1.0.0")
	require.NoError(t, err)
	assert.True(t, exists)
}

// TestGitTestRepo_Bare 测试 bare 仓库
func TestGitTestRepo_Bare(t *testing.T) {
	repo := NewBareGitTestRepo().
		WithDefaultBranch("main").
		Build(t)

	assert.NotNil(t, repo)
	assert.True(t, repo.IsBare())
	assert.NotEmpty(t, repo.Path())
	// bare 仓库没有 Repository 包装
	assert.Nil(t, repo.Repository())
}

// TestGitTestRepo_MultipleCommits 测试多个提交
func TestGitTestRepo_MultipleCommits(t *testing.T) {
	repo := NewGitTestRepo().
		WithFileString("file1.txt", "content 1").
		WithCommit("First commit").
		WithFileString("file2.txt", "content 2").
		WithCommit("Second commit").
		WithFileString("file3.txt", "content 3").
		WithCommit("Third commit").
		Build(t)

	assert.NotNil(t, repo)

	// 验证最后一次提交
	commit, err := repo.Repository().GetLastCommit()
	require.NoError(t, err)
	assert.Equal(t, "Third commit", commit.Message)
}

// TestGitTestRepo_WithUser 测试自定义用户信息
func TestGitTestRepo_WithUser(t *testing.T) {
	repo := NewGitTestRepo().
		WithUser("Custom User", "custom@example.com").
		WithFileString("test.txt", "test content").
		WithCommit("Initial commit").
		Build(t)

	assert.NotNil(t, repo)

	// 验证提交使用了自定义用户信息
	commit, err := repo.Repository().GetLastCommit()
	require.NoError(t, err)
	assert.Contains(t, commit.Author, "Custom User")
}

// TestGitTestRepo_WithCommitFiles 测试提交指定文件
func TestGitTestRepo_WithCommitFiles(t *testing.T) {
	repo := NewGitTestRepo().
		WithFileString("file1.txt", "content 1").
		WithFileString("file2.txt", "content 2").
		WithCommitFiles("Partial commit", "file1.txt").
		Build(t)

	assert.NotNil(t, repo)

	// 验证只有 file1.txt 被提交
	status, err := repo.Repository().Status()
	require.NoError(t, err)
	// file2.txt 应该在未跟踪列表中
	assert.Contains(t, status.UntrackedFiles, "file2.txt")
}


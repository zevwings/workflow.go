package git

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// setupTestRepo 在临时目录中创建并初始化一个 Git 仓库
//
// 返回:
//   - repo: 初始化的 Repository 实例
//   - tempDir: 临时目录路径
func setupTestRepo(t *testing.T) (*Repository, string) {
	t.Helper()

	// 创建临时目录（自动清理）
	tempDir := t.TempDir()

	// 初始化 Git 仓库
	repo, err := Init(tempDir, "main")
	require.NoError(t, err, "初始化 Git 仓库失败")

	// 配置 Git 用户（用于提交）
	config, err := repo.repo.Config()
	require.NoError(t, err)
	config.User.Name = "Test User"
	config.User.Email = "test@example.com"
	err = repo.repo.SetConfig(config)
	require.NoError(t, err)

	return repo, tempDir
}

// setupTestRepoWithCommit 创建带有一个初始提交的测试仓库
func setupTestRepoWithCommit(t *testing.T) (*Repository, string) {
	t.Helper()

	repo, tempDir := setupTestRepo(t)

	// 创建测试文件
	testFile := filepath.Join(tempDir, "test.txt")
	err := os.WriteFile(testFile, []byte("test content"), 0644)
	require.NoError(t, err)

	// 添加文件到暂存区
	err = repo.Add("test.txt")
	require.NoError(t, err)

	// 提交
	author := &object.Signature{
		Name:  "Test User",
		Email: "test@example.com",
	}
	_, err = repo.Commit("Initial commit", author)
	require.NoError(t, err)

	return repo, tempDir
}

// ==================== Open 测试 ====================

func TestOpen(t *testing.T) {
	tests := []struct {
		name    string
		setup   func(t *testing.T) string
		wantErr bool
		errMsg  string
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
			name: "non-existent path",
			setup: func(t *testing.T) string {
				return "/non/existent/path"
			},
			wantErr: true,
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
			path := tt.setup(t)

			repo, err := Open(path)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, repo)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, repo)
				assert.Equal(t, path, repo.Path())
			}
		})
	}
}

func TestOpenCurrent(t *testing.T) {
	// 创建临时目录并切换到该目录
	tempDir := t.TempDir()
	_, err := Init(tempDir, "main")
	require.NoError(t, err)

	// 保存当前工作目录
	oldWd, err := os.Getwd()
	require.NoError(t, err)
	defer os.Chdir(oldWd)

	// 切换到临时目录
	err = os.Chdir(tempDir)
	require.NoError(t, err)

	// 测试 OpenCurrent
	repo, err := OpenCurrent()
	require.NoError(t, err)
	assert.NotNil(t, repo)

	// 使用规范化的路径比较（处理 macOS 的 /var -> /private/var 符号链接）
	absTempDir, err := filepath.Abs(tempDir)
	require.NoError(t, err)
	canonicalTempDir, err := filepath.EvalSymlinks(absTempDir)
	require.NoError(t, err)

	canonicalRepoPath, err := filepath.EvalSymlinks(repo.Path())
	require.NoError(t, err)
	assert.Equal(t, canonicalTempDir, canonicalRepoPath)
}

// ==================== Init 测试 ====================

func TestInit(t *testing.T) {
	tests := []struct {
		name          string
		initialBranch string
		wantErr       bool
	}{
		{
			name:          "init with main branch",
			initialBranch: "main",
			wantErr:       false,
		},
		{
			name:          "init with master branch",
			initialBranch: "master",
			wantErr:       false,
		},
		{
			name:          "init with custom branch",
			initialBranch: "develop",
			wantErr:       false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tempDir := t.TempDir()

			repo, err := Init(tempDir, tt.initialBranch)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, repo)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, repo)
				assert.Equal(t, tempDir, repo.Path())

				// 注意：在没有提交的情况下，可能无法获取当前分支
				// 这里我们只验证仓库已成功初始化
				// 分支验证需要先创建提交
			}
		})
	}
}

// ==================== IsGitRepo 测试 ====================

func TestIsGitRepo(t *testing.T) {
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
		{
			name: "non-existent path",
			setup: func(t *testing.T) string {
				return "/non/existent/path"
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			path := tt.setup(t)

			result := IsGitRepo(path)

			assert.Equal(t, tt.want, result)
		})
	}
}

// ==================== Repository 方法测试 ====================

func TestRepository_Path(t *testing.T) {
	repo, tempDir := setupTestRepo(t)

	path := repo.Path()
	assert.Equal(t, tempDir, path)
}

func TestRepository_Repo(t *testing.T) {
	repo, _ := setupTestRepo(t)

	gitRepo := repo.Repo()
	assert.NotNil(t, gitRepo)
	assert.IsType(t, &git.Repository{}, gitRepo)
}

func TestRepository_Worktree(t *testing.T) {
	repo, _ := setupTestRepo(t)

	worktree := repo.Worktree()
	assert.NotNil(t, worktree)
	assert.IsType(t, &git.Worktree{}, worktree)
}

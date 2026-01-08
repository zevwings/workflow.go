package git

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
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

// setupMockRemoteRepo 创建一个 Mock 远程仓库（本地 bare 仓库）
//
// 返回:
//   - remoteRepo: 远程仓库的 Repository 实例
//   - remotePath: 远程仓库路径（可用作 remote URL）
func setupMockRemoteRepo(t *testing.T, defaultBranch string) (*Repository, string) {
	t.Helper()

	// 创建临时目录作为远程仓库（bare 仓库）
	remoteDir := t.TempDir()

	// 使用 go-git 创建 bare 仓库
	opts := &git.PlainInitOptions{
		InitOptions: git.InitOptions{
			DefaultBranch: plumbing.ReferenceName(fmt.Sprintf("refs/heads/%s", defaultBranch)),
		},
		Bare: true, // 创建 bare 仓库
	}

	remoteGitRepo, err := git.PlainInitWithOptions(remoteDir, opts)
	require.NoError(t, err, "创建 Mock 远程仓库失败")

	// 配置 Git 用户
	config, err := remoteGitRepo.Config()
	require.NoError(t, err)
	config.User.Name = "Test User"
	config.User.Email = "test@example.com"
	err = remoteGitRepo.SetConfig(config)
	require.NoError(t, err)

	// 创建远程仓库的 Repository 包装
	remoteRepo := &Repository{
		repo: remoteGitRepo,
		path: remoteDir,
	}

	return remoteRepo, remoteDir
}

// setupMockRemoteRepoWithCommit 创建带提交的 Mock 远程仓库
// 使用更简单的方法：先创建普通仓库，提交后转换为 bare 仓库
func setupMockRemoteRepoWithCommit(t *testing.T, defaultBranch string) (*Repository, string, plumbing.Hash) {
	t.Helper()

	// 先创建普通仓库并提交
	tempWorkDir := t.TempDir()
	workRepo, err := Init(tempWorkDir, defaultBranch)
	require.NoError(t, err)

	// 创建文件并提交
	testFile := filepath.Join(tempWorkDir, "test.txt")
	err = os.WriteFile(testFile, []byte("test content"), 0644)
	require.NoError(t, err)

	err = workRepo.Add("test.txt")
	require.NoError(t, err)

	author := &object.Signature{
		Name:  "Test User",
		Email: "test@example.com",
	}
	commitHash, err := workRepo.Commit("Initial commit", author)
	require.NoError(t, err)

	// 创建 bare 仓库目录
	remoteDir := t.TempDir()

	// 使用 go-git 创建 bare 仓库
	opts := &git.PlainInitOptions{
		InitOptions: git.InitOptions{
			DefaultBranch: plumbing.ReferenceName(fmt.Sprintf("refs/heads/%s", defaultBranch)),
		},
		Bare: true,
	}

	remoteGitRepo, err := git.PlainInitWithOptions(remoteDir, opts)
	require.NoError(t, err)

	// 将工作仓库的对象复制到 bare 仓库
	// 获取工作仓库的所有对象
	workIter, err := workRepo.repo.Objects()
	require.NoError(t, err)

	err = workIter.ForEach(func(obj object.Object) error {
		// 将 object.Object 转换为 EncodedObject
		encodedObj := &plumbing.MemoryObject{}
		err := obj.Encode(encodedObj)
		if err != nil {
			return err
		}
		// 复制对象到 bare 仓库
		_, err = remoteGitRepo.Storer.SetEncodedObject(encodedObj)
		return err
	})
	require.NoError(t, err)

	// 复制引用
	refs, err := workRepo.repo.References()
	require.NoError(t, err)

	err = refs.ForEach(func(ref *plumbing.Reference) error {
		// 只复制分支引用
		if ref.Name().IsBranch() {
			return remoteGitRepo.Storer.SetReference(ref)
		}
		return nil
	})
	require.NoError(t, err)

	// 创建 HEAD 引用指向默认分支
	headRef := plumbing.NewSymbolicReference(
		plumbing.HEAD,
		plumbing.ReferenceName(fmt.Sprintf("refs/heads/%s", defaultBranch)),
	)
	err = remoteGitRepo.Storer.SetReference(headRef)
	require.NoError(t, err)

	// 创建远程仓库的 Repository 包装
	remoteRepo := &Repository{
		repo: remoteGitRepo,
		path: remoteDir,
	}

	return remoteRepo, remoteDir, commitHash
}

// getRealTestRemoteURL 获取真实测试远程仓库 URL（如果配置了环境变量）
//
// 如果设置了 GIT_TEST_REMOTE_URL 环境变量，返回该 URL
// 否则返回空字符串，表示使用 Mock 远程仓库
func getRealTestRemoteURL(t *testing.T) string {
	t.Helper()
	return os.Getenv("GIT_TEST_REMOTE_URL")
}

// setupTestRemoteRepo 根据配置创建测试远程仓库
//
// 如果设置了 GIT_TEST_REMOTE_URL 环境变量，使用真实的远程仓库
// 否则使用 Mock 远程仓库（本地 bare 仓库）
//
// 返回:
//   - remotePath: 远程仓库路径或 URL
//   - isReal: 是否为真实远程仓库
func setupTestRemoteRepo(t *testing.T, defaultBranch string) (string, bool, plumbing.Hash) {
	t.Helper()

	// 检查是否配置了真实测试远程仓库
	realRemoteURL := getRealTestRemoteURL(t)
	if realRemoteURL != "" {
		// 使用真实远程仓库
		// 注意：真实远程仓库需要已经存在并包含提交
		// 这里我们返回 URL，实际的提交哈希需要从远程获取
		return realRemoteURL, true, plumbing.ZeroHash
	}

	// 使用 Mock 远程仓库
	_, remoteDir, commitHash := setupMockRemoteRepoWithCommit(t, defaultBranch)
	return remoteDir, false, commitHash
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

				// 验证仓库路径正确
				absPath, err := filepath.Abs(tempDir)
				require.NoError(t, err)
				canonicalPath, err := filepath.EvalSymlinks(absPath)
				require.NoError(t, err)

				canonicalRepoPath, err := filepath.EvalSymlinks(repo.Path())
				require.NoError(t, err)
				assert.Equal(t, canonicalPath, canonicalRepoPath)

				// 验证仓库可以打开
				openedRepo, err := Open(tempDir)
				assert.NoError(t, err)
				assert.NotNil(t, openedRepo)
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

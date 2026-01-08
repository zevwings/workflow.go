package git

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/stretchr/testify/suite"
)

// RepositoryTestSuite 是 Git Repository 测试套件
// 所有测试共享一个带初始提交的仓库设置
//
// 使用 Suite 的优势：
// 1. 减少重复的 setup 代码
// 2. 所有测试共享同一个仓库实例
// 3. 更好的测试组织和可读性
// 4. 支持 SetupSuite/TearDownSuite 用于一次性设置
type RepositoryTestSuite struct {
	suite.Suite
	repo    *Repository
	tempDir string
}

// SetupTest 在每个测试运行前执行
// 为每个测试创建新的仓库实例，确保测试之间的隔离
func (s *RepositoryTestSuite) SetupTest() {
	// 使用现有的 setup 函数创建带初始提交的测试仓库
	s.repo, s.tempDir = setupTestRepoWithCommit(s.T())
}

// TearDownTest 在每个测试运行后执行（通常不需要，因为使用了 t.TempDir()）
func (s *RepositoryTestSuite) TearDownTest() {
	// t.TempDir() 会自动清理，这里可以添加额外的清理逻辑
}

// TestRepositoryTestSuite 运行整个测试套件
func TestRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(RepositoryTestSuite))
}

// ==================== Repository 基础测试 ====================

func (s *RepositoryTestSuite) TestPath() {
	path := s.repo.Path()
	s.NotEmpty(path)
	s.Equal(s.tempDir, path)
}

func (s *RepositoryTestSuite) TestRepo() {
	gitRepo := s.repo.Repo()
	s.NotNil(gitRepo)
}

func (s *RepositoryTestSuite) TestWorktree() {
	worktree := s.repo.Worktree()
	s.NotNil(worktree)
}

// ==================== CurrentBranch 测试 ====================

func (s *RepositoryTestSuite) TestCurrentBranch() {
	branch, err := s.repo.CurrentBranch()
	s.NoError(err)
	s.NotEmpty(branch)
	s.Equal("main", branch)
}

// ==================== CreateBranch 测试 ====================

func (s *RepositoryTestSuite) TestCreateBranch() {
	err := s.repo.CreateBranch("feature/test")
	s.NoError(err)

	// 验证分支已创建
	exists, err := s.repo.BranchExists("feature/test")
	s.NoError(err)
	s.True(exists)
}

// ==================== CheckoutBranch 测试 ====================

func (s *RepositoryTestSuite) TestCheckoutBranch() {
	// 创建新分支
	branchName := "feature/test"
	err := s.repo.CreateBranch(branchName)
	s.NoError(err)

	// 切换到新分支
	err = s.repo.CheckoutBranch(branchName)
	s.NoError(err)

	// 验证当前分支
	currentBranch, err := s.repo.CurrentBranch()
	s.NoError(err)
	s.Equal(branchName, currentBranch)
}

func (s *RepositoryTestSuite) TestCheckoutBranch_NonExistent() {
	err := s.repo.CheckoutBranch("non-existent-branch")
	s.Error(err)
}

// ==================== CreateAndCheckoutBranch 测试 ====================

func (s *RepositoryTestSuite) TestCreateAndCheckoutBranch() {
	branchName := "feature/new"
	err := s.repo.CreateAndCheckoutBranch(branchName)
	s.NoError(err)

	// 验证分支已创建
	exists, err := s.repo.BranchExists(branchName)
	s.NoError(err)
	s.True(exists)

	// 验证已切换到新分支
	currentBranch, err := s.repo.CurrentBranch()
	s.NoError(err)
	s.Equal(branchName, currentBranch)
}

// ==================== DeleteBranch 测试 ====================

func (s *RepositoryTestSuite) TestDeleteBranch() {
	// 创建并切换到新分支
	branchName := "feature/to-delete"
	err := s.repo.CreateAndCheckoutBranch(branchName)
	s.NoError(err)

	// 切换回主分支
	err = s.repo.CheckoutBranch("main")
	s.NoError(err)

	// 删除分支
	err = s.repo.DeleteBranch(branchName)
	s.NoError(err)

	// 验证分支已删除
	exists, err := s.repo.BranchExists(branchName)
	s.NoError(err)
	s.False(exists)
}

func (s *RepositoryTestSuite) TestDeleteBranch_CurrentBranch() {
	currentBranch, err := s.repo.CurrentBranch()
	s.NoError(err)

	// 尝试删除当前分支
	err = s.repo.DeleteBranch(currentBranch)
	s.Error(err)
	s.Contains(err.Error(), "cannot delete current branch")
}

// ==================== BranchExists 测试 ====================

func (s *RepositoryTestSuite) TestBranchExists() {
	// 验证主分支存在
	exists, err := s.repo.BranchExists("main")
	s.NoError(err)
	s.True(exists)

	// 验证不存在的分支
	exists, err = s.repo.BranchExists("non-existent")
	s.NoError(err)
	s.False(exists)
}

// ==================== Status 测试 ====================

func (s *RepositoryTestSuite) TestStatus() {
	// 创建新文件（未跟踪）
	newFile := filepath.Join(s.tempDir, "new.txt")
	err := os.WriteFile(newFile, []byte("new content"), 0644)
	s.NoError(err)

	// 修改已存在的文件
	existingFile := filepath.Join(s.tempDir, "test.txt")
	err = os.WriteFile(existingFile, []byte("modified content"), 0644)
	s.NoError(err)

	// 获取状态
	status, err := s.repo.Status()
	s.NoError(err)
	s.NotNil(status)

	// 验证新文件在未跟踪列表中
	s.Contains(status.UntrackedFiles, "new.txt")
	// 验证修改的文件在修改列表中
	s.Contains(status.ModifiedFiles, "test.txt")
}

func (s *RepositoryTestSuite) TestStatus_Clean() {
	status, err := s.repo.Status()
	s.NoError(err)
	s.NotNil(status)

	// 干净的工作区应该没有更改
	s.Empty(status.ModifiedFiles)
	s.Empty(status.StagedFiles)
	s.Empty(status.UntrackedFiles)
}

// ==================== Add 测试 ====================

func (s *RepositoryTestSuite) TestAdd() {
	// 创建新文件
	newFile := filepath.Join(s.tempDir, "new.txt")
	err := os.WriteFile(newFile, []byte("new content"), 0644)
	s.NoError(err)

	// 添加文件
	err = s.repo.Add("new.txt")
	s.NoError(err)

	// 验证文件已添加到暂存区
	status, err := s.repo.Status()
	s.NoError(err)
	s.Contains(status.StagedFiles, "new.txt")
}

// ==================== Commit 测试 ====================

func (s *RepositoryTestSuite) TestCommit() {
	// 创建新文件并添加
	newFile := filepath.Join(s.tempDir, "commit-test.txt")
	err := os.WriteFile(newFile, []byte("commit test"), 0644)
	s.NoError(err)

	err = s.repo.Add("commit-test.txt")
	s.NoError(err)

	// 提交
	author := &object.Signature{
		Name:  "Test User",
		Email: "test@example.com",
	}
	hash, err := s.repo.Commit("Test commit", author)
	s.NoError(err)
	s.NotEqual(plumbing.ZeroHash, hash)

	// 验证提交
	commit, err := s.repo.GetCommit(hash)
	s.NoError(err)
	s.Equal("Test commit", commit.Message)
	s.Contains(commit.Author, "Test User")
}

// ==================== GetHead 测试 ====================

func (s *RepositoryTestSuite) TestGetHead() {
	hash, err := s.repo.GetHead()
	s.NoError(err)
	s.NotEqual(plumbing.ZeroHash, hash)
}

// ==================== GetLastCommit 测试 ====================

func (s *RepositoryTestSuite) TestGetLastCommit() {
	commit, err := s.repo.GetLastCommit()
	s.NoError(err)
	s.NotNil(commit)
	s.Equal("Initial commit", commit.Message)
}

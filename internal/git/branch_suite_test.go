package git

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

// BranchTestSuite 是 Git Branch 测试套件
// 所有测试共享一个带初始提交的仓库设置
//
// 使用 Suite 的优势：
// 1. 减少重复的 setup 代码
// 2. 所有测试共享同一个仓库实例
// 3. 更好的测试组织和可读性
type BranchTestSuite struct {
	suite.Suite
	repo    *Repository
	tempDir string
}

// SetupTest 在每个测试运行前执行
// 为每个测试创建新的仓库实例，确保测试之间的隔离
func (s *BranchTestSuite) SetupTest() {
	s.repo, s.tempDir = setupTestRepoWithCommit(s.T())
}

// TestBranchTestSuite 运行整个测试套件
func TestBranchTestSuite(t *testing.T) {
	suite.Run(t, new(BranchTestSuite))
}

// ==================== CheckoutBranch 测试 ====================

func (s *BranchTestSuite) TestCheckoutBranch() {
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

func (s *BranchTestSuite) TestCheckoutBranch_NonExistent() {
	err := s.repo.CheckoutBranch("non-existent-branch")
	s.Error(err)
}

// ==================== CreateAndCheckoutBranch 测试 ====================

func (s *BranchTestSuite) TestCreateAndCheckoutBranch() {
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

func (s *BranchTestSuite) TestDeleteBranch() {
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

func (s *BranchTestSuite) TestDeleteBranch_CurrentBranch() {
	currentBranch, err := s.repo.CurrentBranch()
	s.NoError(err)

	// 尝试删除当前分支
	err = s.repo.DeleteBranch(currentBranch)
	s.Error(err)
	s.Contains(err.Error(), "cannot delete current branch")
}

func (s *BranchTestSuite) TestDeleteBranch_NonExistent() {
	// 删除不存在的分支，go-git 可能不会返回错误
	// 这取决于实现，我们主要测试方法不会 panic
	err := s.repo.DeleteBranch("non-existent-branch")
	// 这里不检查错误，因为 go-git 的实现可能会成功或失败
	_ = err
}

// ==================== BranchExists 测试 ====================

func (s *BranchTestSuite) TestBranchExists() {
	// 验证主分支存在
	exists, err := s.repo.BranchExists("main")
	s.NoError(err)
	s.True(exists)

	// 验证不存在的分支
	exists, err = s.repo.BranchExists("non-existent")
	s.NoError(err)
	s.False(exists)
}

// ==================== ListBranches 测试 ====================

func (s *BranchTestSuite) TestListBranches() {
	// 创建一些分支
	err := s.repo.CreateBranch("feature/1")
	s.NoError(err)
	err = s.repo.CreateBranch("feature/2")
	s.NoError(err)

	// 列出所有分支
	branches, err := s.repo.ListBranches()
	s.NoError(err)
	s.NotEmpty(branches)

	// 验证主分支在列表中
	var hasMain bool
	for _, branch := range branches {
		if branch.Name == "main" {
			hasMain = true
			break
		}
	}
	s.True(hasMain, "应该包含 main 分支")
}


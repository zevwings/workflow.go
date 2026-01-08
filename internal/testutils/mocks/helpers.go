//go:build test

package mocks

import (
	"testing"

	"github.com/zevwings/workflow/internal/config"
)

// NewMockGitRepositoryWithDefaults 创建配置了默认期望的 MockGitRepository
//
// 这是一个辅助函数，简化常见的 mock 设置场景。
//
// 参数:
//   - t: 测试对象（用于注册清理）
//   - repoPath: 仓库路径
//   - remoteURL: 远程仓库 URL（如果为空，则不设置 remote）
//
// 返回:
//   - 配置好的 MockGitRepository 实例
//
// 示例:
//
//	mockRepo := mocks.NewMockGitRepositoryWithDefaults(t, "/path/to/repo", "https://github.com/owner/repo.git")
func NewMockGitRepositoryWithDefaults(t *testing.T, repoPath, remoteURL string) *MockGitRepository {
	mock := new(MockGitRepository)
	mock.On("GetRepoPath").Return(repoPath)
	mock.On("IsGitRepo", repoPath).Return(true)

	var mockGitRepo config.GitRepo
	if remoteURL != "" {
		mockGitRepoInstance := new(MockGitRepo)
		mockGitRepoInstance.On("GetRemoteURL", "origin").Return(remoteURL, nil)
		mockGitRepo = mockGitRepoInstance
	}

	mock.On("Open", repoPath).Return(mockGitRepo, nil)

	// 注册清理函数，确保所有期望都被验证
	t.Cleanup(func() {
		mock.AssertExpectations(t)
	})

	return mock
}

// NewMockGitRepository 创建基本的 MockGitRepository（不设置默认期望）
//
// 参数:
//   - t: 测试对象（用于注册清理）
//
// 返回:
//   - MockGitRepository 实例
//
// 示例:
//
//	mockRepo := mocks.NewMockGitRepository(t)
//	mockRepo.On("GetRepoPath").Return("/path/to/repo")
func NewMockGitRepository(t *testing.T) *MockGitRepository {
	mock := new(MockGitRepository)
	t.Cleanup(func() {
		mock.AssertExpectations(t)
	})
	return mock
}

// NewMockGitRepo 创建基本的 MockGitRepo（不设置默认期望）
//
// 参数:
//   - t: 测试对象（用于注册清理）
//
// 返回:
//   - MockGitRepo 实例
//
// 示例:
//
//	mockRepo := mocks.NewMockGitRepo(t)
//	mockRepo.On("GetRemoteURL", "origin").Return("https://github.com/owner/repo.git", nil)
func NewMockGitRepo(t *testing.T) *MockGitRepo {
	mock := new(MockGitRepo)
	t.Cleanup(func() {
		mock.AssertExpectations(t)
	})
	return mock
}


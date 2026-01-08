//go:build test

package mocks

import (
	"github.com/stretchr/testify/mock"
	"github.com/zevwings/workflow/internal/config"
)

// MockGitRepository 是 GitRepository 接口的 mock 实现
//
// 使用 testify/mock 框架创建，支持设置期望值和返回值。
//
// 使用示例:
//
//	mockRepo := new(mocks.MockGitRepository)
//	mockRepo.On("GetRepoPath").Return("/path/to/repo")
//	mockRepo.On("IsGitRepo", "/path/to/repo").Return(true)
type MockGitRepository struct {
	mock.Mock
}

// GetRepoPath 实现 GitRepository 接口
func (m *MockGitRepository) GetRepoPath() string {
	args := m.Called()
	return args.String(0)
}

// IsGitRepo 实现 GitRepository 接口
func (m *MockGitRepository) IsGitRepo(path string) bool {
	args := m.Called(path)
	return args.Bool(0)
}

// Open 实现 GitRepository 接口
func (m *MockGitRepository) Open(path string) (config.GitRepo, error) {
	args := m.Called(path)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(config.GitRepo), args.Error(1)
}

// MockGitRepo 是 GitRepo 接口的 mock 实现
//
// 使用示例:
//
//	mockRepo := new(mocks.MockGitRepo)
//	mockRepo.On("GetRemoteURL", "origin").Return("https://github.com/owner/repo.git", nil)
type MockGitRepo struct {
	mock.Mock
}

// GetRemoteURL 实现 GitRepo 接口
func (m *MockGitRepo) GetRemoteURL(name string) (string, error) {
	args := m.Called(name)
	return args.String(0), args.Error(1)
}


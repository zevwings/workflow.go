# Mock 生成工具使用指南

本文档介绍如何在项目中使用 mock 生成工具来自动化创建 mock 对象。

## 工具选择

我们推荐使用 `github.com/stretchr/testify/mock` 配合 `go generate` 和接口定义来创建 mock。

### 为什么选择 testify/mock？

1. **无需额外依赖**：项目已经使用了 `stretchr/testify`
2. **简单易用**：直接使用 Go 代码，无需外部工具
3. **灵活强大**：支持设置期望、返回值、调用次数等

## 使用方法

### 方法 1: 使用 testify/mock（推荐）

#### 步骤 1: 创建 mock 结构体

对于需要 mock 的接口，创建一个 mock 结构体：

```go
// internal/testutils/mocks/git_repository_mock.go
package mocks

import (
	"github.com/stretchr/testify/mock"
	"github.com/zevwings/workflow/internal/config"
	"github.com/zevwings/workflow/internal/git"
)

// MockGitRepository 是 GitRepository 接口的 mock 实现
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
type MockGitRepo struct {
	mock.Mock
}

// GetRemoteURL 实现 GitRepo 接口
func (m *MockGitRepo) GetRemoteURL(name string) (string, error) {
	args := m.Called(name)
	return args.String(0), args.Error(1)
}
```

#### 步骤 2: 在测试中使用

```go
func TestNewRepoManager_WithGitRepo(t *testing.T) {
	// 创建 mock
	mockGitRepo := new(mocks.MockGitRepository)
	mockGitRepoInstance := new(mocks.MockGitRepo)

	// 设置期望
	mockGitRepo.On("GetRepoPath").Return("/path/to/repo")
	mockGitRepo.On("IsGitRepo", "/path/to/repo").Return(true)
	mockGitRepo.On("Open", "/path/to/repo").Return(mockGitRepoInstance, nil)
	mockGitRepoInstance.On("GetRemoteURL", "origin").Return("https://github.com/owner/repo.git", nil)

	// 执行测试
	manager, err := config.NewRepoManager(mockGitRepo)

	// 验证
	require.NoError(t, err)
	assert.NotNil(t, manager)
	mockGitRepo.AssertExpectations(t)
}
```

### 方法 2: 使用 mockery（可选）

如果项目规模较大，可以考虑使用 [mockery](https://github.com/vektra/mockery) 自动生成 mock。

#### 安装 mockery

```bash
go install github.com/vektra/mockery/v2@latest
```

#### 配置 mockery

创建 `.mockery.yaml` 配置文件：

```yaml
with-expecter: true
packages:
  github.com/zevwings/workflow/internal/config:
    interfaces:
      GitRepository:
        config:
          dir: "internal/testutils/mocks"
      GitRepo:
        config:
          dir: "internal/testutils/mocks"
```

#### 生成 mock

```bash
mockery
```

#### 在代码中使用

在接口定义处添加 `//go:generate` 注释：

```go
//go:generate mockery --name=GitRepository --output=../../testutils/mocks

type GitRepository interface {
	GetRepoPath() string
	IsGitRepo(path string) bool
	Open(path string) (GitRepo, error)
}
```

## 最佳实践

### 1. 集中管理 Mock

将所有 mock 对象放在 `internal/testutils/mocks` 目录下，便于管理和查找。

### 2. 使用 Mock 辅助函数

为常用的 mock 设置模式创建辅助函数：

```go
// internal/testutils/mocks/helpers.go
func NewMockGitRepositoryWithDefaults(t *testing.T, repoPath, remoteURL string) *MockGitRepository {
	mock := new(MockGitRepository)
	mock.On("GetRepoPath").Return(repoPath)
	mock.On("IsGitRepo", repoPath).Return(true)

	mockRepo := new(MockGitRepo)
	mockRepo.On("GetRemoteURL", "origin").Return(remoteURL, nil)
	mock.On("Open", repoPath).Return(mockRepo, nil)

	return mock
}
```

### 3. 清理 Mock

在测试完成后调用 `AssertExpectations` 确保所有期望都被满足：

```go
defer mock.AssertExpectations(t)
```

## 迁移现有代码

### 从手动 mock 迁移到 testify/mock

**之前的代码**：
```go
type mockGitRepository struct {
	repoPath string
	isGitRepo bool
	remoteURL string
}

func (m *mockGitRepository) GetRepoPath() string {
	return m.repoPath
}
```

**迁移后的代码**：
```go
mockGitRepo := new(mocks.MockGitRepository)
mockGitRepo.On("GetRepoPath").Return("/path/to/repo")
```

### 迁移步骤

1. 在 `internal/testutils/mocks` 创建 mock 实现
2. 更新测试文件，使用新的 mock
3. 删除旧的 mock 代码
4. 运行测试确保一切正常

## 参考资料

- [testify/mock 文档](https://github.com/stretchr/testify#mock-package)
- [mockery 文档](https://github.com/vektra/mockery)


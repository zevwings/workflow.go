# testify/suite 使用指南

> 使用 `testify/suite` 简化测试代码，减少重复的 setup/teardown 逻辑。

---

## 📋 目录

- [概述](#概述)
- [为什么使用 Suite](#为什么使用-suite)
- [基本用法](#基本用法)
- [Suite 生命周期](#suite-生命周期)
- [使用示例](#使用示例)
- [最佳实践](#最佳实践)
- [何时使用 Suite](#何时使用-suite)

---

## 概述

`testify/suite` 是 testify 提供的测试套件框架，允许你：
- 共享 setup/teardown 逻辑
- 在测试之间共享状态
- 组织相关测试
- 减少重复代码

---

## 为什么使用 Suite

### Before（不使用 Suite）

```go
func TestCreateBranch(t *testing.T) {
    repo, _ := setupTestRepoWithCommit(t)
    // 测试代码...
}

func TestDeleteBranch(t *testing.T) {
    repo, _ := setupTestRepoWithCommit(t)  // 重复的 setup
    // 测试代码...
}

func TestCheckoutBranch(t *testing.T) {
    repo, _ := setupTestRepoWithCommit(t)  // 重复的 setup
    // 测试代码...
}
```

### After（使用 Suite）

```go
type RepositoryTestSuite struct {
    suite.Suite
    repo    *Repository
    tempDir string
}

func (s *RepositoryTestSuite) SetupTest() {
    s.repo, s.tempDir = setupTestRepoWithCommit(s.T())
}

func (s *RepositoryTestSuite) TestCreateBranch() {
    // 直接使用 s.repo，无需重复 setup
    err := s.repo.CreateBranch("feature/test")
    s.NoError(err)
}

func (s *RepositoryTestSuite) TestDeleteBranch() {
    // 直接使用 s.repo
    err := s.repo.DeleteBranch("feature/test")
    s.NoError(err)
}
```

**优势**：
- ✅ 减少重复的 setup 代码
- ✅ 所有测试共享同一个仓库实例
- ✅ 更好的测试组织和可读性
- ✅ 易于维护和扩展

---

## 基本用法

### 1. 定义 Suite 结构

```go
import (
    "testing"
    "github.com/stretchr/testify/suite"
)

type MyTestSuite struct {
    suite.Suite
    // 测试共享的字段
    repo *Repository
    tempDir string
}
```

**重要**：Suite 结构必须嵌入 `suite.Suite`。

### 2. 实现 Setup 和 TearDown 方法

```go
// SetupTest 在每个测试运行前执行
func (s *MyTestSuite) SetupTest() {
    s.repo, s.tempDir = setupTestRepo(s.T())
}

// TearDownTest 在每个测试运行后执行
func (s *MyTestSuite) TearDownTest() {
    // 清理资源（通常不需要，因为使用了 t.TempDir()）
}

// SetupSuite 在套件开始前执行一次（可选）
func (s *MyTestSuite) SetupSuite() {
    // 一次性设置（例如：创建共享资源）
}

// TearDownSuite 在套件结束后执行一次（可选）
func (s *MyTestSuite) TearDownSuite() {
    // 一次性清理
}
```

### 3. 编写测试方法

测试方法必须以 `Test` 开头：

```go
func (s *MyTestSuite) TestCreateBranch() {
    err := s.repo.CreateBranch("feature/test")
    s.NoError(err)  // 使用 s.NoError 而不是 assert.NoError
}
```

### 4. 运行 Suite

```go
func TestMyTestSuite(t *testing.T) {
    suite.Run(t, new(MyTestSuite))
}
```

---

## Suite 生命周期

```
TestMyTestSuite
  ├─ SetupSuite()           (执行一次)
  │
  ├─ SetupTest()            (每个测试前)
  ├─ TestMethod1()
  ├─ TearDownTest()         (每个测试后)
  │
  ├─ SetupTest()            (每个测试前)
  ├─ TestMethod2()
  ├─ TearDownTest()         (每个测试后)
  │
  └─ TearDownSuite()        (执行一次)
```

### 方法说明

| 方法 | 执行时机 | 用途 |
|------|---------|------|
| `SetupSuite()` | 套件开始前执行一次 | 创建共享资源、初始化全局状态 |
| `SetupTest()` | 每个测试运行前 | 为每个测试准备独立的环境 |
| `TearDownTest()` | 每个测试运行后 | 清理测试环境 |
| `TearDownSuite()` | 套件结束后执行一次 | 清理共享资源 |

---

## 使用示例

### 示例 1: Git Repository 测试套件

```go
package git

import (
    "testing"
    "github.com/stretchr/testify/suite"
)

type RepositoryTestSuite struct {
    suite.Suite
    repo    *Repository
    tempDir string
}

func (s *RepositoryTestSuite) SetupTest() {
    s.repo, s.tempDir = setupTestRepoWithCommit(s.T())
}

func (s *RepositoryTestSuite) TestCreateBranch() {
    err := s.repo.CreateBranch("feature/test")
    s.NoError(err)

    exists, err := s.repo.BranchExists("feature/test")
    s.NoError(err)
    s.True(exists)
}

func (s *RepositoryTestSuite) TestCheckoutBranch() {
    err := s.repo.CreateBranch("feature/test")
    s.NoError(err)

    err = s.repo.CheckoutBranch("feature/test")
    s.NoError(err)

    currentBranch, err := s.repo.CurrentBranch()
    s.NoError(err)
    s.Equal("feature/test", currentBranch)
}

func TestRepositoryTestSuite(t *testing.T) {
    suite.Run(t, new(RepositoryTestSuite))
}
```

### 示例 2: HTTP 客户端测试套件

```go
package http

import (
    "net/http"
    "net/http/httptest"
    "testing"
    "github.com/stretchr/testify/suite"
)

type HTTPClientTestSuite struct {
    suite.Suite
    client *Client
    server *httptest.Server
}

func (s *HTTPClientTestSuite) SetupTest() {
    s.client = NewClient()
    s.server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.WriteHeader(http.StatusOK)
        w.Write([]byte(`{"status": "ok"}`))
    }))
}

func (s *HTTPClientTestSuite) TearDownTest() {
    if s.server != nil {
        s.server.Close()
    }
}

func (s *HTTPClientTestSuite) TestGet() {
    resp, err := s.client.Get(s.server.URL)
    s.NoError(err)
    s.Equal(http.StatusOK, resp.StatusCode())
}

func (s *HTTPClientTestSuite) TestPost() {
    resp, err := s.client.Post(s.server.URL, nil)
    s.NoError(err)
    s.Equal(http.StatusOK, resp.StatusCode())
}

func TestHTTPClientTestSuite(t *testing.T) {
    suite.Run(t, new(HTTPClientTestSuite))
}
```

### 示例 3: 使用 SetupSuite 和 TearDownSuite

```go
type DatabaseTestSuite struct {
    suite.Suite
    db     *sql.DB
    connStr string
}

func (s *DatabaseTestSuite) SetupSuite() {
    // 只执行一次：创建测试数据库
    s.connStr = "postgres://test:test@localhost/testdb"
    var err error
    s.db, err = sql.Open("postgres", s.connStr)
    s.Require().NoError(err)
}

func (s *DatabaseTestSuite) TearDownSuite() {
    // 只执行一次：清理测试数据库
    if s.db != nil {
        s.db.Close()
    }
}

func (s *DatabaseTestSuite) SetupTest() {
    // 每个测试前：清理表
    _, err := s.db.Exec("TRUNCATE TABLE users")
    s.NoError(err)
}

func (s *DatabaseTestSuite) TestInsertUser() {
    // 测试代码...
}

func TestDatabaseTestSuite(t *testing.T) {
    suite.Run(t, new(DatabaseTestSuite))
}
```

---

## 最佳实践

### 1. 使用 `s.T()` 获取 `*testing.T`

在 setup/teardown 方法中，使用 `s.T()` 获取 `*testing.T`：

```go
func (s *MyTestSuite) SetupTest() {
    s.repo, s.tempDir = setupTestRepo(s.T())  // ✅ 正确
    // 不要使用 t *testing.T 参数
}
```

### 2. 使用 Suite 的断言方法

使用 `s.Assert()` 或 `s.Require()` 的快捷方法：

```go
func (s *MyTestSuite) TestSomething() {
    s.NoError(err)        // 而不是 assert.NoError(s.T(), err)
    s.Equal(expected, actual)
    s.True(condition)
    s.Contains(slice, item)

    // Require 在失败时会停止测试
    s.Require().NotNil(obj)  // 如果失败，后续代码不会执行
}
```

### 3. SetupTest vs SetupSuite

- **SetupTest**：为每个测试创建独立的环境（推荐用于大多数场景）
- **SetupSuite**：只在需要共享昂贵资源时使用（例如：数据库连接）

### 4. 测试隔离

每个测试都应该相互独立，即使它们共享同一个 Suite：

```go
func (s *MyTestSuite) SetupTest() {
    // ✅ 好：为每个测试创建新的仓库
    s.repo, s.tempDir = setupTestRepoWithCommit(s.T())
}

func (s *MyTestSuite) SetupSuite() {
    // ⚠️ 谨慎：所有测试共享同一个仓库
    // 只在确保测试不会相互影响时使用
}
```

### 5. 清理资源

大多数情况下，使用 `t.TempDir()` 自动清理：

```go
func (s *MyTestSuite) SetupTest() {
    tempDir := s.T().TempDir()  // 自动清理
    s.repo, _ = setupTestRepo(s.T(), tempDir)
}

func (s *MyTestSuite) TearDownTest() {
    // 通常不需要手动清理
}
```

---

## 何时使用 Suite

### ✅ 适合使用 Suite 的场景

1. **多个测试共享相同的 setup 逻辑**
   ```go
   // 10+ 个测试都需要 setupTestRepoWithCommit(t)
   ```

2. **需要在测试之间共享状态**
   ```go
   // 所有测试都需要同一个配置的仓库
   ```

3. **有复杂的测试前置条件**
   ```go
   // 需要多个步骤才能设置好测试环境
   ```

4. **需要组织相关测试**
   ```go
   // 同一功能的多个测试方法
   ```

### ❌ 不适合使用 Suite 的场景

1. **测试之间没有共享逻辑**
   ```go
   // 每个测试的 setup 都不同
   ```

2. **简单的单测试用例**
   ```go
   // 只有一个测试函数
   ```

3. **表驱动测试**
   ```go
   // 使用表驱动测试更合适
   ```

---

## 对比：Suite vs 普通测试 vs 表驱动测试

### Suite（适合共享 setup）

```go
type RepositoryTestSuite struct {
    suite.Suite
    repo *Repository
}

func (s *RepositoryTestSuite) SetupTest() {
    s.repo = setupTestRepo(s.T())
}

func (s *RepositoryTestSuite) TestCreateBranch() { ... }
func (s *RepositoryTestSuite) TestDeleteBranch() { ... }
```

### 表驱动测试（适合参数化测试）

```go
func TestCreateBranch(t *testing.T) {
    tests := []struct {
        name    string
        branch  string
        wantErr bool
    }{
        {"valid branch", "feature/test", false},
        {"invalid branch", "", true},
    }
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            repo := setupTestRepo(t)
            err := repo.CreateBranch(tt.branch)
            // ...
        })
    }
}
```

### 混合使用（推荐）

```go
// Suite 用于共享 setup
type RepositoryTestSuite struct {
    suite.Suite
    repo *Repository
}

// 在 Suite 中使用表驱动测试
func (s *RepositoryTestSuite) TestCreateBranch_TableDriven() {
    tests := []struct {
        name    string
        branch  string
        wantErr bool
    }{
        {"valid", "feature/test", false},
        {"invalid", "", true},
    }
    for _, tt := range tests {
        s.Run(tt.name, func() {
            err := s.repo.CreateBranch(tt.branch)
            if tt.wantErr {
                s.Error(err)
            } else {
                s.NoError(err)
            }
        })
    }
}
```

---

## 相关资源

- [testify/suite 官方文档](https://pkg.go.dev/github.com/stretchr/testify/suite)
- [testify 断言文档](../../../internal/testutils/README.md)

---

**最后更新**: 2025-01-28


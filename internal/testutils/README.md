# testutils 测试辅助工具包

> 提供测试辅助功能，简化测试代码编写，支持测试环境隔离和跨平台兼容性。

**⚠️ 重要提示**：此包使用构建标签 `test`，只在测试时可用，不会被打包到 release 中。

---

## 📋 目录

- [概述](#概述)
- [构建标签](#构建标签)
- [路径获取函数](#路径获取函数)
- [测试数据加载](#测试数据加载)
- [CLI 命令测试](#cli-命令测试)
- [HTTP 测试辅助](#http-测试辅助)
- [Git 测试辅助](#git-测试辅助)
- [使用示例](#使用示例)

---

## 概述

`testutils` 包提供以下功能：

- **路径获取函数**：统一的路径获取函数，支持测试环境隔离
- **测试数据加载**：从 `testdata/` 目录加载测试数据
- **CLI 命令测试**：简化 CLI 命令的执行和断言
- **HTTP 测试辅助**：简化 HTTP 测试服务器的创建和管理
- **Git 测试辅助**：简化 Git 仓库的创建和管理，减少测试设置代码

---

## 构建标签

`testutils` 包使用构建标签 `//go:build test`，确保只在测试时编译，不会被打包到 release 中。

### 工作原理

所有 `testutils` 包的文件都包含以下构建标签：

```go
//go:build test

package testutils
```

这意味着：
- ✅ **使用 `-tags=test` 时**：testutils 包会被编译，可以在测试中使用
- ❌ **不使用标签时**：testutils 包不会被编译，如果生产代码导入会失败

### 验证

```bash
# ✅ 正常构建（不包含 testutils）
go build ./cmd/workflow

# ✅ 使用标签构建 testutils
go build -tags=test ./internal/testutils

# ❌ 不使用标签构建 testutils（会失败）
go build ./internal/testutils
# 输出：package github.com/zevwings/workflow/internal/testutils: build constraints exclude all Go files
```

### 运行测试

使用 `-tags=test` 标签运行测试：

```bash
# 使用 Makefile（已自动包含 -tags=test）
make test
make test-coverage

# 手动运行测试
go test -tags=test ./...

# 运行特定包的测试
go test -tags=test ./internal/config
```

### 在测试代码中使用

```go
package config_test

import (
    "testing"
    "github.com/zevwings/workflow/internal/testutils"
)

func TestExample(t *testing.T) {
    // 使用 testutils
    homeDir := testutils.TestHomeDir(t)
    // ...
}
```

### 注意事项

1. **测试时必须使用 `-tags=test`**：否则 testutils 包不会被编译
2. **生产代码不能导入 testutils**：如果导入，编译会失败（这是好的，提醒开发者）
3. **Makefile 已更新**：`make test` 和 `make test-coverage` 已自动包含 `-tags=test`

### 为什么使用构建标签？

- ✅ **避免打包到 release**：testutils 是测试专用工具，不应该出现在生产代码中
- ✅ **编译时检查**：如果生产代码错误导入了 testutils，编译会失败，及时发现问题
- ✅ **清晰的职责分离**：明确区分测试代码和生产代码
- ✅ **减少二进制大小**：testutils 不会被包含在最终二进制文件中

---

## 路径获取函数

提供统一的路径获取函数，支持测试环境隔离和跨平台兼容性。

### 可用函数

| 函数 | 说明 | 环境变量优先级 |
|------|------|----------------|
| `TestHomeDir(t)` | 获取主目录 | `HOME` > `USERPROFILE` (Windows) |
| `TestConfigDir(t)` | 获取配置目录 | `XDG_CONFIG_HOME` > `HOME/.config` > `APPDATA` (Windows) |
| `TestDataDir(t)` | 获取数据目录 | `XDG_DATA_HOME` > `HOME/.local/share` > `APPDATA` (Windows) |
| `TestCacheDir(t)` | 获取缓存目录 | `XDG_CACHE_HOME` > `HOME/.cache` > `LOCALAPPDATA` (Windows) |
| `TestWorkflowConfigDir(t)` | 获取 Workflow 配置目录 | `HOME/.workflow` |

### 使用示例

```go
import (
    "testing"
    "path/filepath"
    "github.com/zevwings/workflow/internal/testutils"
)

func TestWithPaths(t *testing.T) {
    // 获取测试主目录（支持环境变量隔离）
    homeDir := testutils.TestHomeDir(t)
    configDir := testutils.TestConfigDir(t)
    dataDir := testutils.TestDataDir(t)
    cacheDir := testutils.TestCacheDir(t)
    workflowConfigDir := testutils.TestWorkflowConfigDir(t)

    // 使用测试目录
    configPath := filepath.Join(configDir, "config.toml")
    // ...
}
```

### 环境变量隔离

路径获取函数优先使用环境变量，支持测试环境隔离：

```go
func TestWithEnvIsolation(t *testing.T) {
    // 设置环境变量
    tempDir := t.TempDir()
    t.Setenv("HOME", tempDir)
    t.Setenv("XDG_CONFIG_HOME", filepath.Join(tempDir, ".config"))

    // 路径获取函数会使用环境变量
    homeDir := testutils.TestHomeDir(t)
    assert.Equal(t, tempDir, homeDir)

    configDir := testutils.TestConfigDir(t)
    assert.Contains(t, configDir, tempDir)
}
```

---

## 测试数据加载

从 `testdata/fixtures/` 目录加载测试数据。

### 可用函数

| 函数 | 说明 |
|------|------|
| `LoadFixture(t, filename)` | 加载测试数据文件（二进制） |
| `LoadTextFixture(t, filename)` | 加载测试文本文件 |
| `LoadBinaryFixture(t, filename)` | 加载测试二进制文件（别名） |

### 使用示例

```go
import (
    "encoding/json"
    "testing"
    "github.com/zevwings/workflow/internal/testutils"
    "github.com/stretchr/testify/assert"
)

func TestLoadFixture(t *testing.T) {
    // 加载测试数据
    data := testutils.LoadFixture(t, "sample_github_pr.json")

    // 使用测试数据
    var pr GitHubPR
    err := json.Unmarshal(data, &pr)
    assert.NoError(t, err)
    assert.NotNil(t, pr)
}

func TestLoadTextFixture(t *testing.T) {
    // 加载文本文件
    content := testutils.LoadTextFixture(t, "sample_pr_body.md")

    // 使用文本内容
    assert.Contains(t, content, "PR Title")
}
```

### 测试数据目录结构

```
testdata/
├── fixtures/
│   ├── sample_github_pr.json
│   ├── sample_jira_response.json
│   └── sample_pr_body.md
└── integration/
    └── workflow_scenarios.json
```

---

## CLI 命令测试

简化 CLI 命令的执行和断言。

### 可用函数

| 函数 | 说明 |
|------|------|
| `ExecuteCommand(t, command, args...)` | 执行 CLI 命令 |
| `ExecuteCommandWithEnv(t, env, command, args...)` | 执行 CLI 命令（带环境变量） |
| `ExecuteCommandWithDir(t, dir, command, args...)` | 执行 CLI 命令（带工作目录） |
| `ExecuteCommandCapture(t, command, args...)` | 执行 CLI 命令并捕获输出 |
| `ExecuteCommandCaptureWithEnv(t, env, command, args...)` | 执行 CLI 命令并捕获输出（带环境变量） |
| `ExecuteCommandCaptureWithDir(t, dir, command, args...)` | 执行 CLI 命令并捕获输出（带工作目录） |

### 使用示例

```go
import (
    "testing"
    "github.com/zevwings/workflow/internal/testutils"
    "github.com/stretchr/testify/assert"
)

func TestCLICommand(t *testing.T) {
    // 执行 CLI 命令
    output, err := testutils.ExecuteCommand(t, "workflow", "version")
    assert.NoError(t, err)
    assert.Contains(t, output, "version")
}

func TestCLICommandWithEnv(t *testing.T) {
    // 设置环境变量
    env := map[string]string{
        "HOME": t.TempDir(),
    }

    // 执行命令
    output, err := testutils.ExecuteCommandWithEnv(t, env, "workflow", "config", "show")
    assert.NoError(t, err)
    assert.Contains(t, output, "config")
}

func TestCLICommandCapture(t *testing.T) {
    // 执行命令并捕获输出
    result := testutils.ExecuteCommandCapture(t, "workflow", "version")

    assert.NoError(t, result.Err)
    assert.Contains(t, result.Stdout, "version")
    assert.Empty(t, result.Stderr)
}
```

---

## HTTP 测试辅助

简化 HTTP 测试服务器的创建和管理，减少样板代码。

### 可用函数和类型

| 类型/函数 | 说明 |
|----------|------|
| `HTTPTestServer` | HTTP 测试服务器封装 |
| `HTTPTestServerBuilder` | 构建器模式，用于配置测试服务器 |
| `NewHTTPTestServer()` | 创建新的构建器 |
| `ReadRequestBody(t, r, v)` | 读取并解析请求体为 JSON |
| `AssertRequestMethod(t, r, method)` | 断言请求方法 |
| `AssertRequestPath(t, r, path)` | 断言请求路径 |
| `AssertRequestHeader(t, r, key, value)` | 断言请求头 |
| `AssertRequestQuery(t, r, key, value)` | 断言查询参数 |

### 基本使用

```go
import (
    "net/http"
    "testing"
    "github.com/zevwings/workflow/internal/testutils"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
)

func TestSimpleHTTP(t *testing.T) {
    // 创建简单的测试服务器
    server := testutils.NewHTTPTestServer().
        WithStatus(http.StatusOK).
        WithStringBody("success").
        Build(t)

    // 使用服务器 URL 进行测试
    resp, err := http.Get(server.URL())
    require.NoError(t, err)
    defer resp.Body.Close()

    assert.Equal(t, http.StatusOK, resp.StatusCode)
}
```

### 设置响应

```go
func TestWithJSONResponse(t *testing.T) {
    // 设置 JSON 响应（自动设置 Content-Type）
    server := testutils.NewHTTPTestServer().
        WithStatus(http.StatusOK).
        WithJSONBody(map[string]string{
            "message": "success",
            "id": "123",
        }).
        Build(t)

    resp, err := http.Get(server.URL())
    require.NoError(t, err)
    defer resp.Body.Close()

    assert.Equal(t, "application/json", resp.Header.Get("Content-Type"))
}

func TestWithHeaders(t *testing.T) {
    // 设置响应头
    server := testutils.NewHTTPTestServer().
        WithStatus(http.StatusCreated).
        WithHeader("X-Custom-Header", "custom-value").
        WithHeader("Authorization", "Bearer token123").
        WithStringBody("created").
        Build(t)

    resp, err := http.Get(server.URL())
    require.NoError(t, err)
    defer resp.Body.Close()

    assert.Equal(t, "custom-value", resp.Header.Get("X-Custom-Header"))
    assert.Equal(t, "Bearer token123", resp.Header.Get("Authorization"))
}
```

### 验证请求

```go
func TestVerifyRequest(t *testing.T) {
    // 验证请求方法和路径
    server := testutils.NewHTTPTestServer().
        WithMethodCheck(http.MethodPost).
        WithPathCheck("/api/users").
        WithStatus(http.StatusOK).
        Build(t)

    req, err := http.NewRequest(http.MethodPost, server.URL()+"/api/users", nil)
    require.NoError(t, err)

    resp, err := http.DefaultClient.Do(req)
    require.NoError(t, err)
    defer resp.Body.Close()

    assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestCustomRequestCheck(t *testing.T) {
    // 自定义请求验证
    server := testutils.NewHTTPTestServer().
        WithRequestCheck(func(t *testing.T, r *http.Request) {
            assert.Equal(t, http.MethodPut, r.Method)
            assert.Equal(t, "/api/update", r.URL.Path)
            assert.Equal(t, "Bearer token", r.Header.Get("Authorization"))
        }).
        WithStatus(http.StatusOK).
        Build(t)

    req, err := http.NewRequest(http.MethodPut, server.URL()+"/api/update", nil)
    require.NoError(t, err)
    req.Header.Set("Authorization", "Bearer token")

    resp, err := http.DefaultClient.Do(req)
    require.NoError(t, err)
    defer resp.Body.Close()
}
```

### 自定义 Handler

```go
func TestCustomHandler(t *testing.T) {
    // 使用自定义 handler
    server := testutils.NewHTTPTestServer().
        WithHandler(func(w http.ResponseWriter, r *http.Request) {
            // 解析请求体
            var data map[string]interface{}
            testutils.ReadRequestBody(t, r, &data)

            // 验证请求
            testutils.AssertRequestMethod(t, r, http.MethodPost)
            testutils.AssertRequestPath(t, r, "/api/data")

            // 返回响应
            w.Header().Set("X-Custom", "value")
            w.WriteHeader(http.StatusCreated)
            w.Write([]byte(`{"id": 123}`))
        }).
        Build(t)

    jsonData := `{"name": "test", "value": 42}`
    req, err := http.NewRequest(http.MethodPost, server.URL()+"/api/data",
        strings.NewReader(jsonData))
    require.NoError(t, err)
    req.Header.Set("Content-Type", "application/json")

    resp, err := http.DefaultClient.Do(req)
    require.NoError(t, err)
    defer resp.Body.Close()

    assert.Equal(t, http.StatusCreated, resp.StatusCode)
    assert.Equal(t, "value", resp.Header.Get("X-Custom"))
}
```

### 便捷断言函数

```go
func TestUsingHelperFunctions(t *testing.T) {
    server := testutils.NewHTTPTestServer().
        WithHandler(func(w http.ResponseWriter, r *http.Request) {
            // 使用便捷断言函数
            testutils.AssertRequestMethod(t, r, http.MethodGet)
            testutils.AssertRequestPath(t, r, "/api/test")
            testutils.AssertRequestHeader(t, r, "X-API-Key", "secret")
            testutils.AssertRequestQuery(t, r, "page", "1")

            w.WriteHeader(http.StatusOK)
        }).
        Build(t)

    req, err := http.NewRequest(http.MethodGet, server.URL()+"/api/test?page=1", nil)
    require.NoError(t, err)
    req.Header.Set("X-API-Key", "secret")

    resp, err := http.DefaultClient.Do(req)
    require.NoError(t, err)
    defer resp.Body.Close()
}
```

### 与现有 HTTP 客户端集成

```go
import (
    "github.com/zevwings/workflow/internal/http"
)

func TestWithInternalHTTPClient(t *testing.T) {
    // 创建测试服务器
    server := testutils.NewHTTPTestServer().
        WithStatus(http.StatusOK).
        WithJSONBody(map[string]string{"message": "success"}).
        Build(t)

    // 使用内部 HTTP 客户端
    client := http.NewClient()
    resp, err := client.Get(server.URL())

    require.NoError(t, err)
    assert.Equal(t, http.StatusOK, resp.StatusCode())
}
```

### 自动清理

`Build()` 方法会自动使用 `t.Cleanup()` 注册清理函数，测试结束时服务器会自动关闭，无需手动调用 `defer server.Close()`。

---

## Git 测试辅助

简化 Git 仓库的创建和管理，减少测试设置代码。使用构建器模式，支持链式调用。

### 可用函数和类型

| 类型/函数 | 说明 |
|----------|------|
| `GitTestRepo` | Git 测试仓库封装 |
| `GitTestRepoBuilder` | 构建器模式，用于配置测试仓库 |
| `NewGitTestRepo()` | 创建新的构建器（普通仓库） |
| `NewBareGitTestRepo()` | 创建新的构建器（bare 仓库，用于模拟远程仓库） |

### 基本使用

```go
import (
    "testing"
    "github.com/zevwings/workflow/internal/testutils"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
)

func TestBasicRepo(t *testing.T) {
    // 创建简单的测试仓库
    repo := testutils.NewGitTestRepo().
        WithDefaultBranch("main").
        Build(t)

    assert.NotNil(t, repo)
    assert.NotEmpty(t, repo.Path())
}
```

### 添加文件和提交

```go
func TestRepoWithCommit(t *testing.T) {
    // 创建带文件、提交的测试仓库
    repo := testutils.NewGitTestRepo().
        WithFileString("test.txt", "test content").
        WithCommit("Initial commit").
        Build(t)

    // 验证提交存在
    head, err := repo.Repository().GetHead()
    require.NoError(t, err)
    assert.NotEmpty(t, head.String())
}

func TestRepoWithMultipleCommits(t *testing.T) {
    // 创建多个提交
    repo := testutils.NewGitTestRepo().
        WithFileString("file1.txt", "content 1").
        WithCommit("First commit").
        WithFileString("file2.txt", "content 2").
        WithCommit("Second commit").
        Build(t)

    // 验证最后一次提交
    commit, err := repo.Repository().GetLastCommit()
    require.NoError(t, err)
    assert.Equal(t, "Second commit", commit.Message)
}
```

### 创建分支

```go
func TestRepoWithBranch(t *testing.T) {
    // 创建分支
    repo := testutils.NewGitTestRepo().
        WithFileString("test.txt", "test content").
        WithCommit("Initial commit").
        WithBranch("feature/test").
        Build(t)

    // 验证分支存在
    exists, err := repo.Repository().BranchExists("feature/test")
    require.NoError(t, err)
    assert.True(t, exists)
}

func TestRepoWithBranchAndCheckout(t *testing.T) {
    // 创建并切换到分支
    repo := testutils.NewGitTestRepo().
        WithFileString("test.txt", "test content").
        WithCommit("Initial commit").
        WithBranchAndCheckout("feature/test").
        Build(t)

    // 验证当前分支
    currentBranch, err := repo.Repository().CurrentBranch()
    require.NoError(t, err)
    assert.Equal(t, "feature/test", currentBranch)
}
```

### 添加远程仓库

```go
func TestRepoWithRemote(t *testing.T) {
    // 先创建一个 bare 仓库作为远程
    remoteRepo := testutils.NewBareGitTestRepo().
        Build(t)

    // 创建本地仓库并添加远程
    repo := testutils.NewGitTestRepo().
        WithFileString("test.txt", "test content").
        WithCommit("Initial commit").
        WithRemote("origin", remoteRepo.Path()).
        Build(t)

    // 验证远程存在
    url, err := repo.Repository().GetRemoteURL("origin")
    require.NoError(t, err)
    assert.Equal(t, remoteRepo.Path(), url)
}
```

### 创建 Tag

```go
func TestRepoWithTag(t *testing.T) {
    // 创建 lightweight tag
    repo := testutils.NewGitTestRepo().
        WithFileString("test.txt", "test content").
        WithCommit("Initial commit").
        WithTag("v1.0.0").
        Build(t)

    // 验证 tag 存在
    exists, err := repo.Repository().TagExists("v1.0.0")
    require.NoError(t, err)
    assert.True(t, exists)
}

func TestRepoWithAnnotatedTag(t *testing.T) {
    // 创建 annotated tag（带消息）
    repo := testutils.NewGitTestRepo().
        WithFileString("test.txt", "test content").
        WithCommit("Initial commit").
        WithTagMessage("v1.0.0", "Release version 1.0.0").
        Build(t)

    // 验证 tag 存在
    exists, err := repo.Repository().TagExists("v1.0.0")
    require.NoError(t, err)
    assert.True(t, exists)
}
```

### 自定义用户信息

```go
func TestRepoWithCustomUser(t *testing.T) {
    // 设置自定义 Git 用户信息
    repo := testutils.NewGitTestRepo().
        WithUser("Custom User", "custom@example.com").
        WithFileString("test.txt", "test content").
        WithCommit("Initial commit").
        Build(t)

    // 验证提交使用了自定义用户信息
    commit, err := repo.Repository().GetLastCommit()
    require.NoError(t, err)
    assert.Contains(t, commit.Author, "Custom User")
}
```

### 提交指定文件

```go
func TestRepoWithCommitFiles(t *testing.T) {
    // 只提交指定文件
    repo := testutils.NewGitTestRepo().
        WithFileString("file1.txt", "content 1").
        WithFileString("file2.txt", "content 2").
        WithCommitFiles("Partial commit", "file1.txt").
        Build(t)

    // 验证只有 file1.txt 被提交
    status, err := repo.Repository().Status()
    require.NoError(t, err)
    // file2.txt 应该在未跟踪列表中
    assert.Contains(t, status.UntrackedFiles, "file2.txt")
}
```

### Bare 仓库（用于模拟远程仓库）

```go
func TestBareRepo(t *testing.T) {
    // 创建 bare 仓库（用于模拟远程仓库）
    remoteRepo := testutils.NewBareGitTestRepo().
        WithDefaultBranch("main").
        Build(t)

    assert.NotNil(t, remoteRepo)
    assert.True(t, remoteRepo.IsBare())
    assert.NotEmpty(t, remoteRepo.Path())
    // bare 仓库没有 Repository 包装
    assert.Nil(t, remoteRepo.Repository())
}
```

### 复杂示例：完整的测试场景

```go
func TestComplexRepo(t *testing.T) {
    // 创建远程仓库
    remoteRepo := testutils.NewBareGitTestRepo().
        WithDefaultBranch("main").
        Build(t)

    // 创建本地仓库，包含文件、提交、分支、远程和 tag
    repo := testutils.NewGitTestRepo().
        WithUser("Test User", "test@example.com").
        WithDefaultBranch("main").
        WithFileString("README.md", "# Project").
        WithFileString("main.go", "package main").
        WithCommit("Initial commit").
        WithTag("v1.0.0").
        WithBranchAndCheckout("feature/new").
        WithFileString("feature.go", "package main").
        WithCommit("Add feature").
        WithBranch("develop").
        WithRemote("origin", remoteRepo.Path()).
        Build(t)

    // 验证仓库状态
    currentBranch, err := repo.Repository().CurrentBranch()
    require.NoError(t, err)
    assert.Equal(t, "feature/new", currentBranch)

    // 验证 tag 存在
    exists, err := repo.Repository().TagExists("v1.0.0")
    require.NoError(t, err)
    assert.True(t, exists)

    // 验证远程存在
    url, err := repo.Repository().GetRemoteURL("origin")
    require.NoError(t, err)
    assert.Equal(t, remoteRepo.Path(), url)
}
```

### 自动清理

`Build()` 方法会自动使用 `t.Cleanup()` 注册清理函数，测试结束时仓库会自动清理，无需手动调用 `defer repo.Close()`。

### 与现有测试代码对比

**Before（使用 setupTestRepoWithCommit）：**
```go
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
}
```

**After（使用 GitTestRepoBuilder）：**
```go
func TestRepository_Commit(t *testing.T) {
    repo := testutils.NewGitTestRepo().
        WithFileString("commit-test.txt", "commit test").
        WithCommit("Initial commit").
        WithFileString("commit-test2.txt", "commit test 2").
        WithCommit("Test commit").
        Build(t)

    // 验证提交
    commit, err := repo.Repository().GetLastCommit()
    assert.NoError(t, err)
    assert.Equal(t, "Test commit", commit.Message)
}
```

---

## 使用示例

### 完整示例：测试配置管理

```go
package config_test

import (
    "path/filepath"
    "testing"
    "github.com/zevwings/workflow/internal/config"
    "github.com/zevwings/workflow/internal/testutils"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
)

func TestGlobalManager(t *testing.T) {
    // 设置测试环境
    tempDir := t.TempDir()
    t.Setenv("HOME", tempDir)

    // 使用 testutils 获取路径
    configDir := testutils.TestWorkflowConfigDir(t)
    configPath := filepath.Join(configDir, "config.toml")

    // 创建配置管理器
    manager, err := config.NewGlobalManager()
    require.NoError(t, err)

    // 测试配置加载
    err = manager.Load()
    assert.NoError(t, err)

    // 验证配置文件路径
    assert.Equal(t, configPath, manager.GetConfigPath())
}
```

### 完整示例：测试 HTTP 客户端（使用 testutils）

```go
package http_test

import (
    "net/http"
    "testing"
    "github.com/zevwings/workflow/internal/http"
    "github.com/zevwings/workflow/internal/testutils"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
)

func TestHTTPClient(t *testing.T) {
    // 使用 testutils 创建 Mock 服务器（自动清理）
    server := testutils.NewHTTPTestServer().
        WithStatus(http.StatusOK).
        WithJSONBody(map[string]int{"id": 123}).
        Build(t)

    // 创建 HTTP 客户端
    client := http.NewClient()

    // 执行请求
    resp, err := client.Get(server.URL())
    require.NoError(t, err)
    assert.Equal(t, http.StatusOK, resp.StatusCode())
}

func TestHTTPClientWithRequestCheck(t *testing.T) {
    // 验证请求方法和路径
    server := testutils.NewHTTPTestServer().
        WithMethodCheck(http.MethodPost).
        WithPathCheck("/api/users").
        WithStatus(http.StatusCreated).
        WithJSONBody(map[string]string{"message": "created"}).
        Build(t)

    client := http.NewClient()
    resp, err := client.Post(server.URL()+"/api/users", nil)

    require.NoError(t, err)
    assert.Equal(t, http.StatusCreated, resp.StatusCode())
}
```

---

## 最佳实践

### 1. 使用路径获取函数

```go
// ✅ 推荐：使用路径获取函数
func TestExample(t *testing.T) {
    homeDir := testutils.TestHomeDir(t)
    configDir := testutils.TestConfigDir(t)
    // 使用路径
}

// ❌ 不推荐：直接使用系统路径
func TestExample(t *testing.T) {
    homeDir, err := os.UserHomeDir()
    if err != nil {
        t.Fatal(err)
    }
    // 不支持测试隔离
}
```

### 2. 使用测试数据加载

```go
// ✅ 推荐：使用 LoadFixture
func TestExample(t *testing.T) {
    data := testutils.LoadFixture(t, "sample_github_pr.json")
    // 使用数据
}

// ❌ 不推荐：硬编码路径
func TestExample(t *testing.T) {
    data, err := os.ReadFile("testdata/fixtures/sample_github_pr.json")
    if err != nil {
        t.Fatal(err)
    }
    // 路径可能不存在
}
```

### 3. 测试环境隔离

```go
// ✅ 推荐：使用环境变量隔离
func TestExample(t *testing.T) {
    tempDir := t.TempDir()
    t.Setenv("HOME", tempDir)

    homeDir := testutils.TestHomeDir(t)
    assert.Equal(t, tempDir, homeDir)
}

// ❌ 不推荐：使用真实系统路径
func TestExample(t *testing.T) {
    homeDir := testutils.TestHomeDir(t)
    // 可能使用真实系统路径，污染系统
}
```

---

## 相关文档

- [测试环境工具指南](../../docs/testing/references/environments.md) - 测试环境工具详细使用方法
- [测试辅助工具指南](../../docs/testing/references/helpers.md) - 测试辅助工具详细使用方法
- [测试编写规范](../../docs/testing/writing.md) - 测试编写规范

---

**最后更新**: 2025-01-28


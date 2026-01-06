# 测试工具指南

> 本文档介绍常用测试工具的使用方法。

---

## 📋 目录

- [testify](#1-testify)
- [go-cmp](#2-go-cmp)
- [httptest](#3-httptest)
- [Mock对象使用规范](#4-mock对象使用规范)
- [测试环境工具](#5-测试环境工具)
- [测试辅助工具](#6-测试辅助工具)

---

## 1. testify

`testify` 是 Go 最流行的测试框架，提供断言、Mock 和测试套件功能。

### 安装

```bash
go get github.com/stretchr/testify
```

### 使用方式

#### 断言（assert）

```go
import (
    "testing"
    "github.com/stretchr/testify/assert"
)

func TestExample(t *testing.T) {
    result := ProcessData("input")

    // 基本断言
    assert.Equal(t, "expected", result)
    assert.NotNil(t, result)
    assert.NoError(t, err)

    // 带消息的断言
    assert.Equal(t, "expected", result, "ProcessData should return expected value")
}
```

#### 必须断言（require）

```go
import (
    "testing"
    "github.com/stretchr/testify/require"
)

func TestExample(t *testing.T) {
    config, err := LoadConfig("config.toml")
    require.NoError(t, err) // 失败时立即停止测试
    require.NotNil(t, config)

    // 继续测试
}
```

#### Mock 对象

```go
import (
    "testing"
    "github.com/stretchr/testify/mock"
)

// 定义接口
type HTTPClient interface {
    Get(url string) ([]byte, error)
}

// 创建 Mock 对象
type MockHTTPClient struct {
    mock.Mock
}

func (m *MockHTTPClient) Get(url string) ([]byte, error) {
    args := m.Called(url)
    return args.Get(0).([]byte), args.Error(1)
}

// 使用 Mock
func TestAPICall(t *testing.T) {
    mockClient := new(MockHTTPClient)
    mockClient.On("Get", "https://api.example.com").Return([]byte("response"), nil)

    // 使用 mockClient 进行测试
    result, err := mockClient.Get("https://api.example.com")
    assert.NoError(t, err)
    assert.Equal(t, []byte("response"), result)

    // 验证 Mock 被调用
    mockClient.AssertExpectations(t)
}
```

### 优势

- 清晰的断言输出
- 丰富的断言函数
- Mock 对象支持
- 测试套件支持

---

## 2. go-cmp

`go-cmp` 提供深度比较功能，特别适合比较复杂的数据结构。

### 安装

```bash
go get github.com/google/go-cmp/cmp
```

### 使用方式

```go
import (
    "testing"
    "github.com/google/go-cmp/cmp"
)

type Config struct {
    Host string
    Port int
}

func TestConfigEqual(t *testing.T) {
    want := &Config{Host: "localhost", Port: 8080}
    got := &Config{Host: "localhost", Port: 8080}

    if diff := cmp.Diff(want, got); diff != "" {
        t.Errorf("Config mismatch (-want +got):\n%s", diff)
    }
}
```

### 自定义比较选项

```go
import (
    "testing"
    "github.com/google/go-cmp/cmp"
    "github.com/google/go-cmp/cmp/cmpopts"
)

func TestConfigEqual_IgnoreFields(t *testing.T) {
    want := &Config{Host: "localhost", Port: 8080, CreatedAt: time.Now()}
    got := &Config{Host: "localhost", Port: 8080, CreatedAt: time.Now().Add(time.Hour)}

    // 忽略 CreatedAt 字段
    opts := cmpopts.IgnoreFields(Config{}, "CreatedAt")
    if diff := cmp.Diff(want, got, opts); diff != "" {
        t.Errorf("Config mismatch (-want +got):\n%s", diff)
    }
}
```

### 优势

- 深度比较复杂数据结构
- 清晰的差异输出
- 灵活的比较选项
- 适合比较结构体、切片、映射等

---

## 3. httptest

`httptest` 是 Go 标准库提供的 HTTP 测试工具，用于创建 Mock HTTP 服务器。

### 使用方式

```go
import (
    "net/http"
    "net/http/httptest"
    "testing"
    "github.com/stretchr/testify/assert"
)

func TestGitHubAPI(t *testing.T) {
    // 创建 Mock 服务器
    server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.WriteHeader(http.StatusOK)
        w.Write([]byte(`{"id": 123, "title": "Test PR"}`))
    }))
    defer server.Close()

    // 使用 Mock 服务器进行测试
    client := NewHTTPClient(server.URL)
    result, err := client.GetPR(123)

    assert.NoError(t, err)
    assert.Equal(t, 123, result.ID)
    assert.Equal(t, "Test PR", result.Title)
}
```

### 测试客户端请求

```go
func TestHTTPHandler(t *testing.T) {
    req := httptest.NewRequest("GET", "/api/pr/123", nil)
    w := httptest.NewRecorder()

    handler := http.HandlerFunc(PRHandler)
    handler.ServeHTTP(w, req)

    assert.Equal(t, http.StatusOK, w.Code)
    assert.Contains(t, w.Body.String(), "PR 123")
}
```

### 优势

- Go 标准库，无需额外依赖
- 简单易用
- 适合 HTTP API 测试
- 支持测试服务器和客户端

---

## 4. Mock对象使用规范

### 何时使用 Mock

- 测试需要调用外部 API（GitHub、Jira 等）
- 测试需要模拟网络请求和响应
- 测试需要避免依赖外部服务
- 测试需要模拟错误情况（网络超时、服务器错误等）

### Mock对象组织规范

```go
// ✅ 推荐：使用 testify/mock
import (
    "testing"
    "github.com/stretchr/testify/mock"
)

type MockHTTPClient struct {
    mock.Mock
}

func (m *MockHTTPClient) Get(url string) ([]byte, error) {
    args := m.Called(url)
    return args.Get(0).([]byte), args.Error(1)
}

func TestAPICall(t *testing.T) {
    mockClient := new(MockHTTPClient)
    mockClient.On("Get", "https://api.example.com").Return([]byte("response"), nil)

    // 执行测试
    result, err := mockClient.Get("https://api.example.com")
    assert.NoError(t, err)
    assert.Equal(t, []byte("response"), result)

    // 验证 Mock 被调用
    mockClient.AssertExpectations(t)
}
```

### Mock使用规则

- **每个测试独立 Mock**：每个测试应创建自己的 Mock 对象实例
- **明确 Mock 范围**：每个 Mock 应明确指定调用的参数和返回值
- **验证 Mock 调用**：重要测试应验证 Mock 是否被正确调用（使用 `AssertExpectations()`）

### 不推荐的用法

```go
// ❌ 不推荐：在测试之间共享 Mock 对象
var mockClient *MockHTTPClient

func Test1(t *testing.T) {
    mockClient = new(MockHTTPClient)
    // ...
}

func Test2(t *testing.T) {
    // 依赖 Test1 的 mockClient
    mockClient.On("Get", "url").Return([]byte("response"), nil)
}
```

---

## 5. 测试环境工具

项目提供了统一的测试环境工具，基于 Go 标准库构建，提供完全隔离的测试环境。

### 包含工具

- **临时目录管理**：使用 `t.TempDir()` 创建临时目录
- **环境变量隔离**：使用 `t.Setenv()` 设置环境变量（自动恢复）
- **测试清理**：使用 `t.Cleanup()` 注册清理函数

### 快速使用

```go
import (
    "testing"
    "github.com/your-org/workflow/testutils"
)

func TestCLICommand(t *testing.T) {
    // 使用 t.TempDir() 创建临时目录
    tempDir := t.TempDir()

    // 使用 t.Setenv() 设置环境变量（自动恢复）
    t.Setenv("HOME", tempDir)

    // 测试代码在完全隔离的环境中运行
}
```

### 详细文档

更多详细信息和使用示例，请参考：
- [测试环境工具指南](./environments.md) - 完整的使用指南和API参考

---

## 6. 测试辅助工具

项目提供了测试辅助工具，简化测试代码编写。

### 包含工具

- **路径获取函数**：`TestHomeDir()`, `TestConfigDir()`, `TestDataDir()`, `TestCacheDir()`
- **测试数据生成**：`LoadFixture()`, `GenerateTestData()`
- **CLI 命令测试**：`ExecuteCommand()`, `CaptureOutput()`

### 快速使用

```go
import (
    "testing"
    "github.com/your-org/workflow/testutils"
)

func TestCLICommand(t *testing.T) {
    // 使用路径获取函数
    homeDir := testutils.TestHomeDir(t)
    configDir := testutils.TestConfigDir(t)

    // 使用测试数据生成
    data := testutils.LoadFixture(t, "sample_github_pr.json")

    // 使用 CLI 命令测试
    output, err := testutils.ExecuteCommand(t, "workflow", "version")
    assert.NoError(t, err)
    assert.Contains(t, output, "version")
}
```

### 详细文档

更多详细信息和使用示例，请参考：
- [测试辅助工具指南](./helpers.md) - 完整的使用指南和API参考

---

## 相关文档

- [测试环境工具指南](./environments.md) - 测试环境工具详细使用方法
- [测试辅助工具指南](./helpers.md) - 测试辅助工具详细使用方法
- [Mock测试指南](./mock-server.md) - Mock测试详细使用方法
- [测试编写规范](../writing.md) - 测试编写规范

---

**最后更新**: 2025-01-28

---

## 📝 变更历史

### 2025-01-28
- **重写文档**：从 Rust 工具（pretty_assertions, rstest, mockito）完全重写为 Go 工具（testify, go-cmp, httptest）
- **更新测试工具**：更新所有代码示例为 Go 风格
- **新增测试环境工具文档**：添加测试环境工具章节
- **新增测试辅助工具文档**：添加测试辅助工具章节

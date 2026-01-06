# Mock 测试指南

> 本文档详细说明 Mock 测试的使用方法，包括 HTTP Mock、接口 Mock 和测试替身。

---

## 📋 目录

- [概述](#-概述)
- [HTTP Mock（httptest）](#1-http-mockhttptest)
- [接口 Mock（testify/mock）](#2-接口-mocktestifymock)
- [Mock 最佳实践](#3-mock-最佳实践)

---

## 📋 概述

Mock 测试允许我们在不依赖外部服务的情况下测试代码，提高测试的稳定性和速度。

### Mock 类型

- **HTTP Mock**：使用 `net/http/httptest` 创建 Mock HTTP 服务器
- **接口 Mock**：使用 `testify/mock` 创建接口 Mock 对象
- **测试替身**：使用简单的结构体实现接口

---

## 1. HTTP Mock（httptest）

`httptest` 是 Go 标准库提供的 HTTP 测试工具，用于创建 Mock HTTP 服务器。

### 1.1 基本使用

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

### 1.2 测试不同的 HTTP 方法

```go
func TestHTTPMethods(t *testing.T) {
    server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        switch r.Method {
        case "GET":
            w.WriteHeader(http.StatusOK)
            w.Write([]byte(`{"method": "GET"}`))
        case "POST":
            w.WriteHeader(http.StatusCreated)
            w.Write([]byte(`{"method": "POST"}`))
        case "PUT":
            w.WriteHeader(http.StatusOK)
            w.Write([]byte(`{"method": "PUT"}`))
        case "DELETE":
            w.WriteHeader(http.StatusNoContent)
        default:
            w.WriteHeader(http.StatusMethodNotAllowed)
        }
    }))
    defer server.Close()

    // 测试不同的 HTTP 方法
    // ...
}
```

### 1.3 测试请求头和参数

```go
func TestRequestHeaders(t *testing.T) {
    server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // 验证请求头
        auth := r.Header.Get("Authorization")
        if auth != "Bearer token123" {
            w.WriteHeader(http.StatusUnauthorized)
            return
        }

        // 验证查询参数
        id := r.URL.Query().Get("id")
        if id != "123" {
            w.WriteHeader(http.StatusBadRequest)
            return
        }

        w.WriteHeader(http.StatusOK)
        w.Write([]byte(`{"id": 123}`))
    }))
    defer server.Close()

    // 测试请求
    // ...
}
```

### 1.4 测试错误情况

```go
func TestHTTPErrors(t *testing.T) {
    tests := []struct {
        name       string
        statusCode int
        wantErr    bool
    }{
        {"success", http.StatusOK, false},
        {"not found", http.StatusNotFound, true},
        {"server error", http.StatusInternalServerError, true},
        {"unauthorized", http.StatusUnauthorized, true},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
                w.WriteHeader(tt.statusCode)
            }))
            defer server.Close()

            client := NewHTTPClient(server.URL)
            _, err := client.GetPR(123)

            if tt.wantErr {
                assert.Error(t, err)
            } else {
                assert.NoError(t, err)
            }
        })
    }
}
```

### 1.5 测试客户端请求

```go
func TestHTTPHandler(t *testing.T) {
    // 创建测试请求
    req := httptest.NewRequest("GET", "/api/pr/123", nil)
    req.Header.Set("Authorization", "Bearer token123")

    // 创建响应记录器
    w := httptest.NewRecorder()

    // 调用处理器
    handler := http.HandlerFunc(PRHandler)
    handler.ServeHTTP(w, req)

    // 验证响应
    assert.Equal(t, http.StatusOK, w.Code)
    assert.Contains(t, w.Body.String(), "PR 123")
}
```

---

## 2. 接口 Mock（testify/mock）

`testify/mock` 提供接口 Mock 功能，用于模拟依赖接口。

### 2.1 定义 Mock 对象

```go
import (
    "testing"
    "github.com/stretchr/testify/mock"
)

// 定义接口
type HTTPClient interface {
    Get(url string) ([]byte, error)
    Post(url string, data []byte) ([]byte, error)
}

// 创建 Mock 对象
type MockHTTPClient struct {
    mock.Mock
}

func (m *MockHTTPClient) Get(url string) ([]byte, error) {
    args := m.Called(url)
    return args.Get(0).([]byte), args.Error(1)
}

func (m *MockHTTPClient) Post(url string, data []byte) ([]byte, error) {
    args := m.Called(url, data)
    return args.Get(0).([]byte), args.Error(1)
}
```

### 2.2 使用 Mock 对象

```go
func TestWithMock(t *testing.T) {
    // 创建 Mock 对象
    mockClient := new(MockHTTPClient)

    // 设置期望
    mockClient.On("Get", "https://api.example.com/pr/123").
        Return([]byte(`{"id": 123}`), nil)

    // 使用 Mock 对象
    client := NewService(mockClient)
    result, err := client.GetPR(123)

    // 验证结果
    assert.NoError(t, err)
    assert.Equal(t, 123, result.ID)

    // 验证 Mock 被调用
    mockClient.AssertExpectations(t)
}
```

### 2.3 Mock 多个调用

```go
func TestMultipleCalls(t *testing.T) {
    mockClient := new(MockHTTPClient)

    // 设置多个期望
    mockClient.On("Get", "https://api.example.com/pr/123").
        Return([]byte(`{"id": 123}`), nil).Once()

    mockClient.On("Get", "https://api.example.com/pr/456").
        Return([]byte(`{"id": 456}`), nil).Once()

    // 使用 Mock 对象
    // ...

    // 验证所有期望都被满足
    mockClient.AssertExpectations(t)
}
```

### 2.4 Mock 错误情况

```go
func TestMockError(t *testing.T) {
    mockClient := new(MockHTTPClient)

    // 设置错误期望
    mockClient.On("Get", "https://api.example.com/pr/123").
        Return(nil, errors.New("network error"))

    // 使用 Mock 对象
    client := NewService(mockClient)
    _, err := client.GetPR(123)

    // 验证错误
    assert.Error(t, err)
    assert.Contains(t, err.Error(), "network error")

    mockClient.AssertExpectations(t)
}
```

### 2.5 Mock 参数匹配

```go
func TestMockMatchers(t *testing.T) {
    mockClient := new(MockHTTPClient)

    // 使用参数匹配器
    mockClient.On("Get", mock.Anything).
        Return([]byte(`{"id": 123}`), nil)

    mockClient.On("Post", mock.MatchedBy(func(url string) bool {
        return strings.Contains(url, "api.example.com")
    }), mock.Anything).
        Return([]byte(`{"success": true}`), nil)

    // 使用 Mock 对象
    // ...

    mockClient.AssertExpectations(t)
}
```

---

## 3. Mock 最佳实践

### 3.1 何时使用 Mock

✅ **适合使用 Mock 的场景**：
- 测试需要调用外部 API（GitHub、Jira 等）
- 测试需要模拟网络请求和响应
- 测试需要避免依赖外部服务
- 测试需要模拟错误情况

❌ **不适合使用 Mock 的场景**：
- 测试简单函数（不需要 Mock）
- 测试内部逻辑（可以使用真实对象）
- 集成测试（应该使用真实服务）

### 3.2 Mock 组织规范

```go
// ✅ 推荐：每个测试独立创建 Mock
func TestExample(t *testing.T) {
    mockClient := new(MockHTTPClient)
    mockClient.On("Get", "url").Return([]byte("response"), nil)

    // 使用 Mock
    // ...

    mockClient.AssertExpectations(t)
}

// ❌ 不推荐：在测试之间共享 Mock
var mockClient *MockHTTPClient

func Test1(t *testing.T) {
    mockClient = new(MockHTTPClient)
    // ...
}

func Test2(t *testing.T) {
    // 依赖 Test1 的 mockClient（不推荐）
}
```

### 3.3 验证 Mock 调用

```go
// ✅ 推荐：验证 Mock 调用
func TestExample(t *testing.T) {
    mockClient := new(MockHTTPClient)
    mockClient.On("Get", "url").Return([]byte("response"), nil)

    // 使用 Mock
    // ...

    // 验证 Mock 被调用
    mockClient.AssertExpectations(t)
}

// ❌ 不推荐：不验证 Mock 调用
func TestExample(t *testing.T) {
    mockClient := new(MockHTTPClient)
    mockClient.On("Get", "url").Return([]byte("response"), nil)

    // 使用 Mock，但不验证
    // ...
}
```

### 3.4 使用测试辅助函数

```go
// ✅ 推荐：使用测试辅助函数创建 Mock
func setupMockClient(t *testing.T) *MockHTTPClient {
    t.Helper()
    mockClient := new(MockHTTPClient)
    return mockClient
}

func TestExample(t *testing.T) {
    mockClient := setupMockClient(t)
    mockClient.On("Get", "url").Return([]byte("response"), nil)
    // ...
}
```

---

## 相关文档

- [测试工具指南](./tools.md) - 测试工具详细使用方法
- [测试编写规范](../writing.md) - 测试编写规范
- [测试组织规范](../organization.md) - 测试组织结构

---

**最后更新**: 2025-01-28

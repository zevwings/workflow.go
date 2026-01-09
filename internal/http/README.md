# HTTP 客户端模块

本模块提供了统一的 HTTP 请求接口，封装了底层 `go-resty/resty` 客户端，简化 HTTP 请求的使用。

## 文件说明

```
internal/http/
├── client.go              # HTTP 客户端接口和实现（449行）
├── request.go             # 请求配置（RequestConfig）（208行）
├── response.go           # 响应封装（HttpResponse）（246行）
├── method.go             # HTTP 方法枚举（40行）
├── auth.go               # 认证相关（Authorization）（29行）
├── retry.go              # 重试配置（RetryConfig）（165行）
├── multipart.go          # Multipart 请求配置（200行）
├── parser.go             # 响应解析器（JSON、Text）（124行）
└── *_test.go             # 测试文件
```

### 核心文件

- **`client.go`**：HTTP 客户端接口定义和实现，提供全局单例和实例创建方法
- **`request.go`**：请求配置结构体，支持链式调用配置请求参数
- **`response.go`**：响应封装结构体，提供延迟解析和多种解析方法
- **`method.go`**：HTTP 方法枚举（GET、POST、PUT、DELETE、PATCH）
- **`auth.go`**：Basic Authentication 认证信息结构体
- **`retry.go`**：重试配置结构体，支持自定义重试策略
- **`multipart.go`**：Multipart 请求配置，用于文件上传
- **`parser.go`**：响应解析器接口和实现（JSON、Text）

## 快速开始

### 使用全局客户端（推荐）

```go
import "github.com/zevwings/workflow/internal/http"

// 获取全局客户端
client := http.Global()

// 发送 GET 请求（新版 API）
resp, err := client.GetWithConfig("https://api.example.com/data", nil)
if err != nil {
    return err
}

// 解析 JSON 响应
var data map[string]interface{}
result, err := http.AsJSON[map[string]interface{}](resp)
if err != nil {
    return err
}
```

### 使用请求配置

```go
import "github.com/zevwings/workflow/internal/http"

client := http.Global()

// 创建请求配置
config := http.NewRequestConfig().
    WithQuery(map[string]string{"page": "1"}).
    WithHeader("X-API-Key", "your-key").
    WithAuth(http.NewAuthorization("user", "pass"))

// 发送 POST 请求
resp, err := client.PostWithConfig("https://api.example.com/data", config)
if err != nil {
    return err
}
```

### 流式请求

```go
import "github.com/zevwings/workflow/internal/http"

client := http.Global()

// 发送流式请求
stream, err := client.Stream(http.MethodGet, "https://api.example.com/stream", nil)
if err != nil {
    return err
}
defer stream.Close()

// 读取流数据
// ...
```

## 主要接口

### Client 接口

- `Get(url)` / `GetWithConfig(url, config)` - GET 请求
- `Post(url, body)` / `PostWithConfig(url, config)` - POST 请求
- `Put(url, body)` / `PutWithConfig(url, config)` - PUT 请求
- `Delete(url)` / `DeleteWithConfig(url, config)` - DELETE 请求
- `Patch(url, body)` / `PatchWithConfig(url, config)` - PATCH 请求
- `Stream(method, url, config)` - 流式请求
- `PostMultipart(url, config)` - Multipart 请求
- `SetAuth(token)` - 设置认证 Token
- `SetBasicAuth(username, password)` - 设置 Basic Auth
- `SetProxy(proxyURL)` - 设置代理

### RequestConfig

- `WithBody(body)` - 设置请求体
- `WithQuery(query)` - 设置查询参数
- `WithHeader(key, value)` - 设置单个 Header
- `WithHeaders(headers)` - 设置多个 Headers
- `WithAuth(auth)` - 设置认证信息
- `WithTimeout(timeout)` - 设置超时时间
- `WithRetry(retry)` - 设置重试配置

### HttpResponse

- `IsSuccess()` - 检查是否为成功响应
- `IsError()` - 检查是否为错误响应
- `AsJSON[T]()` - 解析为 JSON
- `AsText()` - 解析为文本
- `AsBytes()` - 获取原始字节
- `EnsureSuccess()` - 确保响应成功
- `ExtractErrorMessage()` - 提取错误消息

## 注意事项

1. **新旧 API 兼容**：模块同时提供旧版 API（返回 `*resty.Response`）和新版 API（返回 `*HttpResponse`），建议使用新版 API
2. **全局客户端**：`Global()` 返回单例客户端，适合大多数场景
3. **响应解析**：响应体采用延迟解析，可以多次解析为不同格式
4. **默认配置**：客户端默认超时 30 秒，重试 3 次，初始等待 1 秒

## 依赖

- `github.com/go-resty/resty/v2` - 底层 HTTP 客户端库

## 相关文档

- [详细架构文档](../../docs/architecture/http.md) - 模块设计思路和实现细节

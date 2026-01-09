# HTTP 模块架构文档

## 📋 概述

HTTP 模块是 Workflow CLI 的基础工具模块，提供统一的 HTTP 请求接口，封装了底层 `go-resty/resty` 客户端。该模块专注于 HTTP 客户端功能，不涉及命令层的业务逻辑。

HTTP 模块提供完整的 HTTP 请求功能，包括 GET、POST、PUT、DELETE、PATCH 等 HTTP 方法，支持流式请求、multipart 文件上传、认证、代理、重试等功能，总代码行数约 1454+ 行。

**模块统计：**
- 代码行数：约 1454+ 行（不含测试文件）
- 主要文件：8 个核心文件
- 主要结构体：`Client`、`RequestConfig`、`HttpResponse`、`RetryConfig`、`MultipartRequestConfig`
- 支持功能：HTTP 方法、流式请求、文件上传、认证、代理、重试、响应解析

**注意**：本模块是基础库模块，不包含 Commands 层，其他模块通过导入使用。

---

## 📁 Lib 层架构（核心业务逻辑）

HTTP 模块（`internal/http/`）是 Workflow CLI 的基础工具库模块，提供统一的 HTTP 请求接口。该模块专注于封装底层 HTTP 客户端，提供简洁易用的 API，不涉及命令层的业务逻辑。

### 模块结构

```
internal/http/
├── client.go              # HTTP 客户端接口和实现（449行）
├── request.go             # 请求配置（RequestConfig）（208行）
├── response.go            # 响应封装（HttpResponse）（246行）
├── method.go              # HTTP 方法枚举（40行）
├── auth.go                # 认证相关（Authorization）（29行）
├── retry.go               # 重试配置（RetryConfig）（165行）
├── multipart.go           # Multipart 请求配置（200行）
└── parser.go              # 响应解析器（JSON、Text）（124行）
```

**总计：约 1454+ 行代码**

### 依赖模块

- **`github.com/go-resty/resty/v2`**：底层 HTTP 客户端库
  - 提供 HTTP 请求发送、重试机制、连接池等功能

### 模块集成

- **`internal/llm/client/`**：LLM 客户端模块使用 HTTP 模块发送 API 请求
  - `http.Global()` - 获取全局客户端
  - `http.PostWithConfig()` - 发送 POST 请求
  - `http.AsJSON[T]()` - 解析 JSON 响应
- **`internal/testutils/`**：测试工具模块使用 HTTP 模块进行测试
  - `http.NewClient()` - 创建测试客户端
  - `http.RequestConfig` - 配置测试请求

---

## 🏗️ Lib 层架构设计

### 设计原则

1. **接口抽象**：通过 `Client` 接口隐藏实现细节，提供统一的 HTTP 请求接口
2. **链式配置**：`RequestConfig`、`RetryConfig` 等配置结构体支持链式调用，提高代码可读性
3. **延迟解析**：响应体采用延迟解析策略，可以多次解析为不同格式（JSON、Text、Bytes）
4. **向后兼容**：保留旧版 API（返回 `*resty.Response`），同时提供新版 API（返回 `*HttpResponse`）
5. **单例模式**：提供全局单例客户端 `Global()`，减少资源消耗，提高性能

### 核心组件

#### 1. Client 接口和实现 (`client.go`)

**职责**：提供统一的 HTTP 请求接口，封装底层 resty 客户端

**主要方法**：
- `Global()` - 获取全局 HTTP 客户端单例
- `GetWithConfig()`, `PostWithConfig()`, `PutWithConfig()`, `DeleteWithConfig()`, `PatchWithConfig()` - 新版 API，支持 RequestConfig
- `Get()`, `Post()`, `Put()`, `Delete()`, `Patch()` - 旧版 API，保持向后兼容
- `Stream()` - 流式请求
- `PostMultipart()` - Multipart 文件上传
- `SetAuth()`, `SetBasicAuth()`, `SetProxy()` - 客户端级别配置

**关键特性**：
- 全局单例：使用 `sync.Once` 确保线程安全的单例初始化
- 默认配置：超时 30 秒，重试 3 次，初始等待 1 秒
- 智能重试：自动处理 5xx 错误、429 错误和网络错误
- 接口封装：返回接口类型，隐藏实现细节

**使用场景**：
- 发送 HTTP 请求到外部 API
- 流式下载大文件
- 上传文件（multipart）

#### 2. RequestConfig (`request.go`)

**职责**：配置 HTTP 请求参数，支持链式调用

**主要方法**：
- `WithBody()` - 设置请求体
- `WithQuery()` - 设置查询参数
- `WithHeader()`, `WithHeaders()` - 设置 HTTP Headers
- `WithAuth()` - 设置 Basic Authentication
- `WithTimeout()` - 设置超时时间
- `WithRetry()` - 设置重试配置

**关键特性**：
- 链式调用：所有配置方法返回自身，支持链式调用
- 灵活查询参数：支持 `map[string]string`、`map[string]interface{}`、`[]string` 等多种格式
- 默认值处理：未设置的字段使用默认值

**使用场景**：
- 配置请求参数、Headers、认证信息
- 设置请求超时和重试策略

#### 3. HttpResponse (`response.go`)

**职责**：封装 HTTP 响应，提供延迟解析和多种解析方法

**主要方法**：
- `IsSuccess()`, `IsError()` - 检查响应状态
- `AsJSON[T]()` - 解析为 JSON（泛型方法）
- `AsText()` - 解析为文本
- `AsBytes()` - 获取原始字节
- `EnsureSuccess()`, `EnsureSuccessWith()` - 确保响应成功
- `ExtractErrorMessage()` - 提取错误消息
- `ParseWith()` - 使用自定义解析器解析

**关键特性**：
- 延迟解析：响应体字节缓存，可以多次解析为不同格式
- 泛型支持：`AsJSON[T]()` 使用 Go 泛型，类型安全
- 错误提取：自动从 JSON 响应中提取错误消息
- 状态检查：提供便捷的状态码检查方法

**使用场景**：
- 解析 JSON 响应
- 读取文本响应
- 处理错误响应

#### 4. RetryConfig (`retry.go`)

**职责**：配置 HTTP 请求的重试策略

**主要方法**：
- `WithRetryCount()` - 设置重试次数
- `WithRetryWaitTime()` - 设置初始等待时间
- `WithRetryMaxWaitTime()` - 设置最大等待时间
- `WithRetryCondition()` - 设置自定义重试条件
- `WithRetryAfter()` - 设置自定义重试延迟函数
- `DisableRetry()` - 禁用重试

**关键特性**：
- 默认重试条件：自动重试 5xx 错误、429 错误和网络错误
- 指数退避：支持指数退避策略
- Retry-After 支持：自动解析 `Retry-After` header
- 自定义策略：支持自定义重试条件和延迟函数

**使用场景**：
- 配置请求重试策略
- 处理不稳定的网络连接
- 应对服务器临时错误

#### 5. MultipartRequestConfig (`multipart.go`)

**职责**：配置 Multipart 文件上传请求

**主要方法**：
- `WithMultipartField()` - 添加 multipart 字段
- `WithMultipartFields()` - 设置多个 multipart 字段
- `WithQuery()`, `WithAuth()`, `WithHeader()`, `WithHeaders()` - 继承基础配置方法

**关键特性**：
- 文件上传：支持文件路径和流式上传
- 多字段支持：支持多个 multipart 字段
- 继承配置：继承 `baseRequestConfig` 的所有配置方法

**使用场景**：
- 上传文件到服务器
- 发送 multipart/form-data 请求

#### 6. 响应解析器 (`parser.go`)

**职责**：定义响应解析器接口和实现（JSON、Text）

**主要组件**：
- `ResponseParser` 接口 - 定义解析器接口
- `JsonParser` - JSON 解析器实现
- `TextParser` - 文本解析器实现
- `ParseJSON()`, `ParseText()` - 便捷解析函数

**关键特性**：
- 接口设计：通过接口支持扩展自定义解析器
- 错误处理：提供详细的解析错误信息
- 空响应处理：正确处理空响应和空白响应

**使用场景**：
- 解析 JSON 响应
- 解析文本响应
- 扩展自定义解析器（XML、YAML 等）

#### 7. HTTP 方法枚举 (`method.go`)

**职责**：定义 HTTP 方法枚举和解析函数

**主要组件**：
- `HttpMethod` 类型 - HTTP 方法枚举
- `MethodGet`, `MethodPost`, `MethodPut`, `MethodDelete`, `MethodPatch` - 方法常量
- `ParseHttpMethod()` - 从字符串解析 HTTP 方法

**关键特性**：
- 类型安全：使用枚举类型而非字符串，提高类型安全
- 验证：解析函数会验证方法是否有效

**使用场景**：
- 指定 HTTP 请求方法
- 流式请求方法选择

#### 8. 认证 (`auth.go`)

**职责**：定义 Basic Authentication 认证信息结构体

**主要组件**：
- `Authorization` 结构体 - Basic Auth 认证信息
- `NewAuthorization()` - 创建认证信息

**关键特性**：
- 简单封装：封装用户名和密码
- 类型安全：使用结构体而非参数列表

**使用场景**：
- 配置 Basic Authentication
- 设置 API 认证信息

### 设计模式

#### 1. 单例模式

**实现**：`Global()` 函数使用 `sync.Once` 确保全局客户端只初始化一次

**优势**：
- 资源复用：所有请求共享同一个连接池
- 性能优化：避免重复创建客户端实例
- 线程安全：可以在多线程环境中安全使用

#### 2. 建造者模式

**实现**：`RequestConfig`、`RetryConfig`、`MultipartRequestConfig` 支持链式调用

**优势**：
- 代码可读性：链式调用使代码更清晰
- 灵活性：可以按需配置，未设置的字段使用默认值
- 易用性：减少函数参数数量

#### 3. 策略模式

**实现**：`ResponseParser` 接口支持不同的解析策略（JSON、Text）

**优势**：
- 可扩展性：可以轻松添加新的解析器（XML、YAML 等）
- 解耦：解析逻辑与响应处理解耦
- 灵活性：可以为不同请求选择不同的解析策略

#### 4. 适配器模式

**实现**：`HttpResponse` 封装 `resty.Response`，提供统一的响应接口

**优势**：
- 接口统一：隐藏底层实现细节
- 延迟解析：响应体可以多次解析
- 向后兼容：保留对 `resty.Response` 的访问（通过 `GetRestyClient()`）

### 错误处理

#### 分层错误处理

1. **网络层错误**：连接超时、网络不可达等，自动重试
2. **HTTP 层错误**：4xx 客户端错误、5xx 服务器错误
   - 5xx 错误：自动重试
   - 4xx 错误：不重试，返回错误
   - 429 错误：自动重试，支持 `Retry-After` header
3. **解析层错误**：JSON 解析失败、文本编码错误等

#### 容错机制

- **自动重试**：网络错误和 5xx 错误自动重试，最多 3 次
- **指数退避**：重试延迟采用指数退避策略，避免服务器压力过大
- **Retry-After 支持**：自动解析并遵守服务器的 `Retry-After` header
- **错误提取**：从 JSON 响应中自动提取错误消息，提供更友好的错误信息

---

## 🔄 集成关系

### 模块使用关系

HTTP 模块被以下模块使用：

1. **`internal/llm/client/`**：LLM 客户端模块
   - 使用 `http.Global()` 获取全局客户端
   - 使用 `http.PostWithConfig()` 发送 API 请求
   - 使用 `http.AsJSON[T]()` 解析 JSON 响应
   - 使用 `http.RequestConfig` 配置请求参数

2. **`internal/testutils/`**：测试工具模块
   - 使用 `http.NewClient()` 创建测试客户端
   - 使用 `http.RequestConfig` 配置测试请求

### 调用流程

#### HTTP 请求流程

```
调用者代码
  ↓
http.Global() / http.NewClient() (获取客户端)
  ↓
client.GetWithConfig() / client.PostWithConfig() (发送请求)
  ↓
RequestConfig.applyToRequest() (应用配置)
  ↓
resty.Client (底层 HTTP 客户端)
  ↓
HTTP 服务器
  ↓
resty.Response (原始响应)
  ↓
FromRestyResponse() (封装响应)
  ↓
HttpResponse (封装后的响应)
  ↓
AsJSON[T]() / AsText() (解析响应)
  ↓
返回结果
```

#### 重试流程

```
发送请求
  ↓
检查响应/错误
  ↓
DefaultRetryCondition() (判断是否重试)
  ↓
是 → DefaultRetryAfter() (计算延迟时间)
  ↓
等待延迟时间
  ↓
重试请求
  ↓
否 → 返回错误
```

---

## 🎯 核心功能

### 1. HTTP 请求发送

**功能说明**：支持 GET、POST、PUT、DELETE、PATCH 等 HTTP 方法，提供新旧两套 API

**流程**：
1. 获取客户端（`Global()` 或 `NewClient()`）
2. 创建请求配置（`NewRequestConfig()`）
3. 配置请求参数（链式调用）
4. 发送请求（`GetWithConfig()` 等）
5. 处理响应（`AsJSON[T]()` 等）

**示例**：
```go
import "github.com/zevwings/workflow/internal/http"

client := http.Global()

config := http.NewRequestConfig().
    WithQuery(map[string]string{"page": "1"}).
    WithHeader("X-API-Key", "your-key")

resp, err := client.GetWithConfig("https://api.example.com/data", config)
if err != nil {
    return err
}

var data map[string]interface{}
result, err := http.AsJSON[map[string]interface{}](resp)
if err != nil {
    return err
}
```

### 2. 流式请求

**功能说明**：支持流式请求，用于处理大文件或流式数据

**流程**：
1. 获取客户端
2. 创建请求配置
3. 发送流式请求（`Stream()`）
4. 读取流数据
5. 关闭流

**示例**：
```go
import "github.com/zevwings/workflow/internal/http"

client := http.Global()

stream, err := client.Stream(http.MethodGet, "https://api.example.com/stream", nil)
if err != nil {
    return err
}
defer stream.Close()

// 读取流数据
buf := make([]byte, 4096)
for {
    n, err := stream.Read(buf)
    if err == io.EOF {
        break
    }
    if err != nil {
        return err
    }
    // 处理数据
    process(buf[:n])
}
```

### 3. 文件上传

**功能说明**：支持 Multipart 文件上传

**流程**：
1. 获取客户端
2. 创建 Multipart 请求配置（`NewMultipartRequestConfig()`）
3. 添加文件字段（`WithMultipartField()`）
4. 发送请求（`PostMultipart()`）
5. 处理响应

**示例**：
```go
import "github.com/zevwings/workflow/internal/http"

client := http.Global()

config := http.NewMultipartRequestConfig().
    WithMultipartField(http.MultipartField{
        ParamName: "file",
        FilePath:   "/path/to/file.txt",
        FileName:   "file.txt",
    })

resp, err := client.PostMultipart("https://api.example.com/upload", config)
if err != nil {
    return err
}
```

### 4. 重试机制

**功能说明**：支持自动重试和自定义重试策略

**流程**：
1. 创建重试配置（`NewRetryConfig()`）
2. 配置重试参数（链式调用）
3. 应用到请求配置（`WithRetry()`）
4. 发送请求，自动重试

**示例**：
```go
import "github.com/zevwings/workflow/internal/http"

client := http.Global()

retryConfig := http.NewRetryConfig().
    WithRetryCount(5).
    WithRetryWaitTime(2 * time.Second).
    WithRetryMaxWaitTime(60 * time.Second)

config := http.NewRequestConfig().
    WithRetry(retryConfig)

resp, err := client.GetWithConfig("https://api.example.com/data", config)
```

---

## 📋 使用示例

### 基础 GET 请求

```go
import "github.com/zevwings/workflow/internal/http"

client := http.Global()
resp, err := client.GetWithConfig("https://api.example.com/data", nil)
if err != nil {
    return err
}

var data map[string]interface{}
result, err := http.AsJSON[map[string]interface{}](resp)
```

### 带配置的 POST 请求

```go
import "github.com/zevwings/workflow/internal/http"

client := http.Global()

config := http.NewRequestConfig().
    WithBody(map[string]string{"key": "value"}).
    WithQuery(map[string]string{"page": "1"}).
    WithHeader("Content-Type", "application/json").
    WithAuth(http.NewAuthorization("user", "pass"))

resp, err := client.PostWithConfig("https://api.example.com/data", config)
if err != nil {
    return err
}

resp, err = resp.EnsureSuccess()
if err != nil {
    return err
}
```

### 错误处理

```go
import "github.com/zevwings/workflow/internal/http"

client := http.Global()
resp, err := client.GetWithConfig("https://api.example.com/data", nil)
if err != nil {
    return err
}

if resp.IsError() {
    errorMsg := resp.ExtractErrorMessage()
    return fmt.Errorf("request failed: %s", errorMsg)
}

// 或者使用 EnsureSuccess
resp, err = resp.EnsureSuccess()
if err != nil {
    return err
}
```

---

## 📝 扩展性

### 添加新的响应解析器

1. 实现 `ResponseParser` 接口
2. 在 `HttpResponse.ParseWith()` 中使用

**示例**：
```go
type XMLParser struct{}

func (p *XMLParser) Parse(bytes []byte, status int) (interface{}, error) {
    // 实现 XML 解析逻辑
    var result interface{}
    err := xml.Unmarshal(bytes, &result)
    return result, err
}

// 使用
resp, err := client.GetWithConfig(url, nil)
if err != nil {
    return err
}

result, err := resp.ParseWith(&XMLParser{})
```

### 添加新的 HTTP 方法

1. 在 `method.go` 中添加新的方法常量
2. 在 `client.go` 的 `doRequest()` 和 `Stream()` 中添加对应的 case
3. 在 `Client` 接口中添加对应的方法（如需要）

---

## 📚 相关文档

- [模块 README](../../internal/http/README.md) - 基础使用说明
- [LLM HTTP 适配器分析](./llm-http-adapter-analysis.md) - LLM 模块使用 HTTP 模块的分析

---

## ✅ 总结

HTTP 模块采用清晰的接口抽象和建造者模式设计：

1. **接口抽象**：通过 `Client` 接口隐藏实现细节，提供统一的 HTTP 请求接口
2. **链式配置**：`RequestConfig`、`RetryConfig` 等支持链式调用，提高代码可读性
3. **延迟解析**：响应体采用延迟解析策略，可以多次解析为不同格式
4. **智能重试**：自动处理网络错误和服务器错误，支持自定义重试策略
5. **向后兼容**：保留旧版 API，同时提供新版 API

**设计优势**：
- ✅ 接口封装：隐藏底层实现，易于测试和替换
- ✅ 类型安全：使用泛型和枚举类型，提高类型安全
- ✅ 易用性：链式调用和默认配置，降低使用门槛
- ✅ 可扩展性：支持自定义解析器和重试策略
- ✅ 性能优化：单例模式和连接池复用，提高性能

**当前实现状态**：
- ✅ 核心功能已实现：HTTP 方法、流式请求、文件上传、认证、代理、重试
- ✅ 响应解析已实现：JSON、Text 解析，支持自定义解析器
- ✅ 测试覆盖：包含完整的单元测试

---

**最后更新**: 2025-01-27

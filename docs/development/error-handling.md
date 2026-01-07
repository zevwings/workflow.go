# 错误处理规范

> 本文档定义了 Workflow CLI 项目的错误处理规范和最佳实践，所有贡献者都应遵循这些规范。

---

## 📋 目录

- [概述](#-概述)
- [错误类型](#-错误类型)
- [错误信息](#-错误信息)
- [错误消息格式规范](#-错误消息格式规范)
- [错误消息内容要求](#-错误消息内容要求)
- [错误消息管理](#-错误消息管理)
- [错误处理模式](#-错误处理模式)
- [分层错误处理](#-分层错误处理)
- [错误包装和上下文](#-错误包装和上下文)
- [相关文档](#-相关文档)

---

## 📋 概述

本文档定义了错误处理规范，包括错误类型、错误信息格式、错误处理模式和分层错误处理。

### 核心原则

- **统一性**：统一使用 Go 标准 `error` 接口作为函数返回类型
- **上下文**：为错误消息添加上下文信息
- **用户友好**：错误消息应清晰、可操作
- **错误包装**：使用 `fmt.Errorf` 和 `%w` 动词进行错误包装，保持错误链

### 使用场景

- 编写新代码时参考
- 错误处理代码审查时检查
- 调试和错误排查时使用

### 快速参考

| 操作 | 方法 | 说明 |
|------|------|------|
| **创建错误** | `fmt.Errorf()` | 创建格式化错误 |
| **包装错误** | `fmt.Errorf("...: %w", err)` | 包装错误并添加上下文 |
| **检查错误** | `errors.Is()` | 检查错误类型 |
| **展开错误** | `errors.Unwrap()` | 展开错误链 |
| **断言错误** | `errors.As()` | 错误类型断言 |

---

## 错误类型

统一使用 Go 标准 `error` 接口作为函数返回类型：

```go
// ✅ 好的做法
func DownloadLogs(ticketID string) ([]byte, error) {
    // 实现
}

// ✅ 返回多个值
func GetUser(id string) (*User, error) {
    // 实现
}
```

**规则**：
- 函数返回错误时，错误应该是最后一个返回值
- 如果函数可能失败，必须返回错误
- 使用 `error` 接口，而不是具体错误类型（除非需要特殊处理）

---

## 错误信息

提供清晰、有上下文的错误信息：

```go
// ✅ 好的做法
import (
    "fmt"
    "os"
)

func ParseConfig(path string) (*Config, error) {
    content, err := os.ReadFile(path)
    if err != nil {
        return nil, fmt.Errorf("读取配置文件失败 %s: %w", path, err)
    }

    var config Config
    if err := toml.Unmarshal(content, &config); err != nil {
        return nil, fmt.Errorf("解析 TOML 配置失败: %w", err)
    }

    return &config, nil
}

// ❌ 不好的做法
func ParseConfig(path string) (*Config, error) {
    content, err := os.ReadFile(path)
    if err != nil {
        return nil, err  // 错误信息不清晰
    }

    var config Config
    if err := toml.Unmarshal(content, &config); err != nil {
        return nil, err  // 缺少上下文
    }

    return &config, nil
}
```

---

## 错误消息格式规范

### 用户友好的错误消息格式

错误消息应遵循以下格式：

1. **包含操作上下文**：说明在做什么操作时出错
2. **包含目标信息**：文件路径、URL、ID 等
3. **包含可操作的指导**：告诉用户如何解决问题

```go
// ✅ 好的错误消息格式
return fmt.Errorf(
    "读取配置文件失败 %s，请检查文件权限或运行 'workflow setup' 创建配置文件",
    path,
)

// ❌ 不好的错误消息格式
return fmt.Errorf("错误: 失败")
```

### 使用统一的错误消息格式

对于常见错误，可以使用统一的格式化函数：

```go
// 定义错误格式化函数
func formatError(operation, target, reason string) string {
    return fmt.Sprintf("执行 %s 操作失败，目标: %s，原因: %s", operation, target, reason)
}

// 使用
return fmt.Errorf(formatError("读取", "config.toml", "权限被拒绝"))
// 输出: "执行 读取 操作失败，目标: config.toml，原因: 权限被拒绝"
```

---

## 错误消息内容要求

### 避免技术术语

错误消息应使用用户可理解的语言：

```go
// ✅ 好的做法：用户友好的语言
return fmt.Errorf("配置文件未找到，请运行 'workflow setup' 创建配置文件")

// ❌ 不好的做法：技术术语
return fmt.Errorf("FileNotFoundError: Config file missing")
```

### 提供解决方案

错误消息应包含解决方案或下一步操作建议：

```go
// ✅ 好的做法：提供解决方案
return fmt.Errorf("JIRA ID 格式无效: %s，期望格式: PROJ-123", input)

// ❌ 不好的做法：只说明问题
return fmt.Errorf("JIRA ID 格式无效")
```

### 区分用户错误和系统错误

- **用户错误**：输入验证失败、配置错误等，应提供清晰的指导
- **系统错误**：网络错误、文件系统错误等，应提供详细的错误信息

```go
// 用户错误：提供格式说明
if !isValidJiraID(input) {
    return fmt.Errorf(
        "JIRA ID 格式无效: %s\n\n期望格式:\n  • Ticket ID: PROJ-123\n  • 项目名称: PROJ",
        input,
    )
}

// 系统错误：提供详细错误信息
resp, err := client.Get(url)
if err != nil {
    return fmt.Errorf("从 %s 获取数据失败: %w", url, err)
}
```

---

## 错误消息管理

### 使用错误消息常量

使用错误消息常量统一管理，避免硬编码：

```go
// 定义错误消息常量
const (
    ErrReadConfigFailed = "读取配置文件失败"
    ErrParseConfigFailed = "解析配置文件失败"
    ErrJiraIDInvalid = "JIRA ID 格式无效"
)

// ✅ 好的做法：使用常量
return fmt.Errorf("%s: %s", ErrReadConfigFailed, path)

// ❌ 不好的做法：硬编码字符串
return fmt.Errorf("读取配置文件失败: %s", path)
```

### 错误消息模板

错误消息模板应包含格式说明：

```go
// 定义错误消息模板
const (
    JiraIDFormatHelp = `期望格式:
  • Ticket ID: PROJ-123
  • 项目名称: PROJ`
)

// 使用
return fmt.Errorf(
    "JIRA ID 格式无效\n%s\n\n错误详情: %s",
    JiraIDFormatHelp,
    input,
)
```

---

## 错误处理模式

### 1. 使用 `fmt.Errorf` 包装错误

```go
import "fmt"

result, err := operation()
if err != nil {
    return nil, fmt.Errorf("执行操作失败，ID: %s: %w", id, err)
}
```

### 2. 使用 `errors.New` 创建简单错误

```go
import "errors"

if value < 0 {
    return errors.New("值必须非负")
}
```

### 3. 使用 `fmt.Errorf` 创建格式化错误

```go
import "fmt"

if condition {
    return fmt.Errorf("错误消息，上下文: %s", value)
}
```

### 4. 使用 `errors.Is` 检查错误

```go
import (
    "errors"
    "os"
)

if errors.Is(err, os.ErrNotExist) {
    // 处理文件不存在错误
}
```

### 5. 使用 `errors.As` 进行错误类型断言

```go
import (
    "errors"
    "net/url"
)

var urlErr *url.Error
if errors.As(err, &urlErr) {
    // 处理 URL 错误
}
```

### 6. 使用 `errors.Unwrap` 展开错误链

```go
import "errors"

// 获取底层错误
underlyingErr := errors.Unwrap(err)
```

---

## 分层错误处理

不同层级使用不同的错误处理策略：

1. **CLI 层**：参数验证错误，使用 `cobra` 自动处理
2. **命令层**：用户交互错误、业务逻辑错误，提供友好的错误提示，可使用日志输出
3. **库层**：底层操作错误（文件、网络、API），提供详细的错误信息，使用错误包装添加上下文

```go
import (
    "fmt"
    "github.com/zevwings/workflow/internal/logging"
)

// 命令层：提供友好的错误提示
func DownloadCommand(ticketID string) error {
    if ticketID == "" {
        return fmt.Errorf("JIRA ticket ID 是必需的")
    }

    // 调用库层，传递详细错误
    logs, err := jira.DownloadFromJira(ticketID)
    if err != nil {
        logging.Errorf("下载日志失败: %v", err)
        return fmt.Errorf("下载日志失败: %w", err)
    }

    // 处理成功情况
    return nil
}

// 库层：提供详细的错误信息
func DownloadFromJira(ticketID string) ([]byte, error) {
    url := fmt.Sprintf("%s/api/ticket/%s", baseURL, ticketID)
    resp, err := client.Get(url)
    if err != nil {
        return nil, fmt.Errorf("从 JIRA 获取 ticket %s 失败: %w", ticketID, err)
    }

    body, err := resp.BodyBytes()
    if err != nil {
        return nil, fmt.Errorf("读取响应体失败: %w", err)
    }

    return body, nil
}
```

---

## 错误包装和上下文

### 使用 `%w` 动词包装错误

Go 1.13+ 支持使用 `%w` 动词包装错误，保持错误链：

```go
import "fmt"

// ✅ 好的做法：使用 %w 包装错误
func ReadConfig(path string) (*Config, error) {
    data, err := os.ReadFile(path)
    if err != nil {
        return nil, fmt.Errorf("读取配置文件 %s 失败: %w", path, err)
    }
    // ...
}

// 调用者可以使用 errors.Is 和 errors.Unwrap
if errors.Is(err, os.ErrNotExist) {
    // 处理文件不存在
}
```

### 添加上下文信息

在错误链的每一层添加上下文信息：

```go
// 底层：原始错误
func readFile(path string) ([]byte, error) {
    data, err := os.ReadFile(path)
    if err != nil {
        return nil, fmt.Errorf("读取文件失败: %w", err)
    }
    return data, nil
}

// 中层：添加业务上下文
func loadConfig(path string) (*Config, error) {
    data, err := readFile(path)
    if err != nil {
        return nil, fmt.Errorf("加载配置失败: %w", err)
    }
    // ...
}

// 上层：添加用户友好的上下文
func initConfig() error {
    config, err := loadConfig("config.toml")
    if err != nil {
        return fmt.Errorf("初始化配置失败，请运行 'workflow setup': %w", err)
    }
    // ...
}
```

### 错误链示例

```go
// 错误链：初始化配置失败，请运行 'workflow setup': 加载配置失败: 读取文件失败: open config.toml: no such file or directory
```

---

## 🔍 故障排除

### 问题 1：错误消息不清晰

**症状**：错误消息缺少上下文信息

**解决方案**：

1. 使用 `fmt.Errorf` 添加上下文
2. 使用统一的错误格式化函数
3. 确保错误消息包含操作上下文和目标信息

### 问题 2：错误链丢失

**症状**：无法追踪错误的原始原因

**解决方案**：

1. 使用 `%w` 动词包装错误，而不是 `%v` 或 `%s`
2. 在错误链的每一层添加上下文
3. 使用 `errors.Is` 和 `errors.Unwrap` 检查和处理错误

### 问题 3：错误类型检查失败

**症状**：使用 `==` 比较错误失败

**解决方案**：

1. 使用 `errors.Is` 检查错误值
2. 使用 `errors.As` 进行错误类型断言
3. 定义自定义错误类型时实现 `Is()` 和 `Unwrap()` 方法

---

## 📚 相关文档

### 开发规范

- [代码风格规范](./code-style.md) - 代码风格规范
- [日志和调试规范](./references/logging.md) - 日志和调试规范

### 检查工作流

- [提交前检查](./workflows/pre-commit.md) - 代码质量检查流程

### Go 官方文档

- [Error Handling](https://go.dev/blog/error-handling-and-go) - Go 错误处理最佳实践
- [Working with Errors](https://go.dev/doc/tutorial/handle-errors) - Go 错误处理教程
- [errors package](https://pkg.go.dev/errors) - errors 包文档

---

## ✅ 检查清单

使用本规范时，请确保：

- [ ] 统一使用 `error` 接口作为函数返回类型
- [ ] 为错误消息添加上下文信息
- [ ] 错误消息使用用户友好的语言
- [ ] 区分用户错误和系统错误
- [ ] 使用错误消息常量统一管理
- [ ] 使用 `%w` 动词包装错误，保持错误链
- [ ] 使用 `errors.Is` 和 `errors.As` 检查错误

---

**最后更新**: 2025-01-27

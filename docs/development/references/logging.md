# 日志和调试规范

> 本文档定义了 Workflow CLI 项目的日志和调试规范和最佳实践，所有贡献者都应遵循这些规范。

---

## 📋 目录

- [概述](#-概述)
- [日志系统架构](#-日志系统架构)
- [日志级别使用规则](#-日志级别使用规则)
- [日志配置](#-日志配置)
- [敏感信息过滤规则](#-敏感信息过滤规则)
- [日志输出规则](#-日志输出规则)
- [最佳实践](#-最佳实践)
- [相关文档](#-相关文档)

---

## 📋 概述

本文档定义了日志和调试规范，包括日志系统架构、日志级别使用规则、敏感信息过滤规则和最佳实践。

### 核心原则

- **职责分离**：Commands 层和 Lib 层使用不同的日志系统
- **敏感信息过滤**：所有日志输出前必须过滤敏感信息
- **日志级别**：根据重要性选择合适的日志级别

### 使用场景

- 编写日志代码时参考
- 调试代码时使用
- 代码审查时检查

---

## 日志系统架构

项目采用**分离的日志系统设计**，实现了职责分离：

- **Commands 层**：使用 `output` 包进行用户友好的控制台输出（带颜色、Emoji）
- **Lib 层**：使用 `logging` 包进行结构化日志记录（默认不输出到控制台，可配置启用）

### Commands 层日志

Commands 层使用 `output` 包，直接输出到控制台，用户可见：

```go
import "github.com/zevwings/workflow/internal/output"

out := output.NewOutput(false)

// 成功消息（总是输出）
out.Success("Operation completed")

// 错误消息
out.Error("Operation failed: %s", errorMsg)

// 警告消息
out.Warning("Retrying operation")

// 信息消息
out.Info("Processing data")

// 调试消息（仅在 verbose 模式下输出）
out.Debug("Debug information: %s", data)
```

### Lib 层日志

Lib 层使用 `logging` 包，默认输出到日志文件，不输出到控制台：

```go
import "github.com/zevwings/workflow/internal/logging"

// 获取带模块名的 logger（自动识别模块名）
logger := logging.GetLogger()

// 调试信息（输出到日志文件）
logger.Debugf("Processing data: %s", data)

// 信息日志（输出到日志文件）
logger.Infof("Operation completed")

// 警告日志（输出到日志文件）
logger.Warnf("Retrying operation")

// 错误日志（输出到日志文件）
logger.Errorf("Operation failed: %v", err)
```

**重要规则**：
- ❌ **禁止**在 Lib 层使用 `output` 包（会直接输出到控制台，影响用户体验）
- ✅ **必须**在 Lib 层使用 `logging` 包（输出到日志文件）
- ✅ **必须**使用 `GetLogger()` 获取 logger，自动包含模块信息

---

## 日志级别使用规则

### Commands 层日志级别

- **`Success()`**：成功消息，总是输出，不受日志级别限制
- **`Error()`**：系统错误，需要立即关注
- **`Warning()`**：警告信息，可能的问题
- **`Info()`**：重要操作信息（如命令执行、配置加载）
- **`Debug()`**：调试信息，仅在 verbose 模式下输出
- **`Print()` / `Println()`**：说明信息，总是输出，不受日志级别限制（用于 setup/check 等命令）

### Lib 层日志级别

- **`Error()` / `Errorf()`**：系统错误，需要立即关注
- **`Warn()` / `Warnf()`**：警告信息，可能的问题
- **`Info()` / `Infof()`**：重要操作信息（如 API 调用、文件操作）
- **`Debug()` / `Debugf()`**：调试信息，仅在开发时使用

---

## 日志配置

### 日志级别配置

日志级别从配置文件 `~/.workflow/config/workflow.toml` 中的 `log.level` 字段读取：

```toml
[log]
level = "info"  # 可选值：off, error, warn, info, debug, trace
enable_trace_console = false  # 是否同时输出 trace_*! 日志到控制台
```

### 日志文件位置

Lib 层的 `trace_*!` 日志默认输出到日志文件：
- 路径：`~/.workflow/logs/tracing/workflow-YYYY-MM-DD.log`
- 格式：按日期分割，每天一个文件
- 存储：强制本地存储（不使用 iCloud 同步）

### 初始化

在 `main()` 函数中初始化日志系统：

```go
import (
    "github.com/zevwings/workflow/internal/logging"
    "github.com/zevwings/workflow/internal/lib/config"
)

func main() error {
    // 初始化配置
    cfg, err := config.Load()
    if err != nil {
        return err
    }

    // 初始化日志系统（从配置文件读取日志级别）
    err = logging.Init(cfg.Log.Level)
    if err != nil {
        return err
    }

    // ... 其他初始化代码
    return nil
}
```

---

## 敏感信息过滤规则

**重要**：所有日志输出前必须过滤敏感信息。

### 使用 `MaskSensitiveValue` 函数

```go
import "github.com/zevwings/workflow/internal/lib/util"

// ❌ 不安全
logging.Infof("API token: %s", token)

// ✅ 安全
logging.Infof("API token: %s", util.MaskSensitiveValue(token))
```

### 使用敏感值包装器

```go
import "github.com/zevwings/workflow/internal/lib/util"

// ❌ 不安全
out.Info("API token: %s", token)

// ✅ 安全
maskedToken := util.MaskSensitiveValue(token)
out.Info("API token: %s", maskedToken)
```

### 敏感信息类型

以下信息必须过滤：
- API Token（GitHub、Jira、LLM 等）
- 密码和密钥
- 用户输入中的敏感信息
- 配置文件路径中的敏感信息（如包含用户名的路径）

---

## 日志输出规则

### Commands 层规则

- 使用 `output` 包，直接输出到控制台，用户可见
- 成功消息和说明信息总是输出，不受日志级别限制
- 调试消息仅在 verbose 模式下输出
- 所有敏感信息必须过滤

### Lib 层规则

- 使用 `logging` 包，默认输出到日志文件
- 不输出到控制台（除非配置启用控制台输出）
- 禁止使用 `output` 包（会直接输出到控制台，影响用户体验）
- 所有敏感信息必须过滤

---

## 最佳实践

### 1. 选择合适的日志级别

```go
logger := logging.GetLogger()

// ✅ 好的做法：使用合适的日志级别
logger.Infof("API request sent")  // 重要操作
logger.Debugf("Request payload: %s", payload)  // 调试信息

// ❌ 不好的做法：过度使用 debug 级别
logger.Debugf("API request sent")  // 应该是 info 级别
```

### 2. 提供有意义的日志消息

```go
logger := logging.GetLogger()

// ✅ 好的做法：提供上下文信息
logger.Infof("Downloading file from %s to %s", url, path)

// ❌ 不好的做法：日志消息不清晰
logger.Infof("Downloading")
```

### 3. 过滤敏感信息

```go
import (
    "github.com/zevwings/workflow/internal/logging"
    "github.com/zevwings/workflow/internal/lib/util"
)

logger := logging.GetLogger()

// ✅ 好的做法：过滤敏感信息
logger.Infof("API token: %s", util.MaskSensitiveValue(token))
logger.Infof("User: %s", util.MaskSensitiveValue(user))

// ❌ 不好的做法：直接输出敏感信息
logger.Infof("API token: %s", token)
```

### 4. 避免过度日志记录

```go
logger := logging.GetLogger()

// ✅ 好的做法：只在关键点记录日志
logger.Infof("Starting operation")
result, err := performOperation()
if err != nil {
    return err
}
logger.Infof("Operation completed")

// ❌ 不好的做法：记录过多细节
logger.Debugf("Step 1")
logger.Debugf("Step 2")
logger.Debugf("Step 3")
// ... 太多日志
```

---

## 🔍 故障排除

### 问题 1：日志未输出

**症状**：日志消息未显示

**解决方案**：

1. 检查日志级别配置
2. 确认日志级别是否允许该消息输出
3. 检查日志文件位置和权限

### 问题 2：敏感信息泄露

**症状**：日志中包含敏感信息

**解决方案**：

1. 使用 `MaskSensitiveValue()` 函数过滤敏感信息
2. 检查所有日志输出点
3. 定期审查日志文件

---

## 📚 相关文档

### 开发规范

- [错误处理规范](../error-handling.md) - 错误处理规范

### 架构文档

- [Logger 模块架构文档](../../../architecture/logger.md) - 详细的日志系统架构说明

### 代码实现

- `internal/logging/logger.go` - 日志系统实现
- `internal/output/output.go` - 控制台输出格式化

---

## ✅ 检查清单

使用本规范时，请确保：

- [ ] Commands 层使用 `output` 包
- [ ] Lib 层使用 `logging` 包
- [ ] 选择合适的日志级别
- [ ] 所有敏感信息已过滤
- [ ] 日志消息清晰有意义

---

**最后更新**: 2025-01-27


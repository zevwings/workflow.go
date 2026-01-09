# Logging 模块

本模块提供了统一的日志记录功能，支持模块级别的日志管理和自动模块识别。

## 文件说明

```
internal/logging/
├── logger.go              # 日志系统核心实现（437行）
│   ├── Init/InitWithFiles # 初始化日志系统
│   ├── GetLogger          # 获取带模块名的 logger
│   ├── getCallerModule    # 自动识别调用者模块名
│   ├── getModuleLogger    # 获取模块级别的 logger
│   ├── errorLogHook       # 错误日志 Hook
│   └── LoggerEntry        # 日志 Entry 封装类型
│
└── logger_test.go         # 测试文件（367行）
```

### 核心文件

- **`logger.go`**：日志系统核心实现
  - 提供日志初始化、模块识别、文件输出等功能
  - 支持按模块分别输出到不同文件
  - 支持日志轮转和统一错误日志收集

## 快速开始

### 基本使用

```go
import "github.com/zevwings/workflow/internal/logging"

// 初始化日志系统（支持文件输出）
// 使用 XDG_STATE_HOME（遵循 XDG Base Directory Specification）
// 默认位置：~/.local/state/workflow/logs
logDir := "$XDG_STATE_HOME/workflow/logs" // 或使用 config.StateDir()
logging.InitWithFiles("info", "text", nil, logDir, false)

// 获取 logger（自动识别模块名）
logger := logging.GetLogger()

// 记录日志
logger.Infof("Operation completed")
logger.Errorf("Operation failed: %v", err)
```

### 在 HTTP 模块中使用

```go
package http

import "github.com/zevwings/workflow/internal/logging"

func doRequest() {
    // 自动识别模块名为 "http"
    logger := logging.GetLogger()
    logger.Infof("HTTP request: %s", url)
    // 日志输出到: $XDG_STATE_HOME/workflow/logs/http.log（默认：~/.local/state/workflow/logs/http.log）
}
```

### 在 HTTP 适配器中使用（全局函数）

```go
package adapter

import "github.com/zevwings/workflow/internal/logging"

// HTTP 适配器场景，需要手动指定模块名
func (l *LogrusLogger) Errorf(format string, v ...interface{}) {
    // 使用全局函数，手动指定模块名
    logging.WithField("module", "http").Errorf(format, v...)
}
```

### 在 LLM 模块中使用

```go
package llm

import "github.com/zevwings/workflow/internal/logging"

func (c *llmClient) Call() {
    // 自动识别模块名为 "llm"
    logger := logging.GetLogger()
    logger.Infof("Calling LLM API: %s", url)
    // 日志输出到: $XDG_STATE_HOME/workflow/logs/llm.log（默认：~/.local/state/workflow/logs/llm.log）
}
```

### 添加额外字段

```go
logger := logging.GetLogger()

// 添加单个字段
logger.WithField("user_id", userID).Info("User logged in")

// 添加多个字段（使用 logging.Fields）
logger.WithFields(logging.Fields{
    "request_id": requestID,
    "duration":   duration,
}).Info("Request completed")

// 添加错误信息
logger.WithError(err).Errorf("Operation failed")
```

## 主要接口

### 初始化接口

- `InitWithFiles(level, format, output, logDir, consoleOut)` - 初始化日志系统（支持文件输出）
  - `level`: 日志级别 (debug, info, warn, error)
  - `format`: 日志格式 (text, json)
  - `output`: 输出目标（nil 则使用默认）
  - `logDir`: 日志目录（空则不创建文件日志）
  - `consoleOut`: 是否同时输出到控制台

- `Init(level, format, output)` - 初始化日志系统（向后兼容，不创建文件日志）

- `SetLevel(level)` - 设置日志级别

### 日志记录接口

- `GetLogger()` - 获取带模块名的 logger Entry（自动识别模块名）
  - 返回 `*LoggerEntry`，包含 `module` 字段
  - 如果启用了文件日志，每个模块输出到独立的文件
  - `LoggerEntry` 封装了 `logrus.Entry`，提供更好的封装性

- `WithField(key, value)` - 添加单个字段（全局函数，返回 `*logrus.Entry`）
- `WithFields(fields)` - 添加多个字段（全局函数，接受 `Fields` 类型）
- `WithError(err)` - 添加错误字段（全局函数）

**注意**：
- `GetLogger()` 返回的 `*LoggerEntry` 类型也提供了 `WithField`、`WithFields`、`WithError` 方法，推荐使用 `GetLogger()` 获取 logger 后使用其方法
- 全局函数 `WithField`、`WithFields`、`WithError` 主要用于适配器场景（如 HTTP 适配器），在这些场景中需要手动指定模块名

## 注意事项

1. **必须使用 GetLogger()**：所有模块必须使用 `GetLogger()` 获取 logger，自动包含模块信息

2. **模块名自动识别**：系统通过 `runtime.Caller()` 自动识别调用者模块名，无需手动指定

3. **日志文件结构**：
   - 每个模块输出到 `{logDir}/{module}.log`
   - 所有错误日志统一输出到 `{logDir}/error.log`

4. **日志轮转**：
   - 单个文件最大 10MB
   - 保留最近 5 个备份
   - 保留 30 天
   - 自动压缩旧日志

5. **线程安全**：模块 logger 的创建是线程安全的，支持并发场景

## 类型说明

### LoggerEntry

`LoggerEntry` 是日志 Entry 的封装类型，提供了以下方法：

- `WithField(key, value)` - 添加单个字段
- `WithFields(fields Fields)` - 添加多个字段（使用 `logging.Fields` 类型）
- `WithError(err)` - 添加错误字段
- `Debug/Debugf` - 记录 Debug 级别日志
- `Info/Infof` - 记录 Info 级别日志
- `Warn/Warnf` - 记录 Warn 级别日志
- `Error/Errorf` - 记录 Error 级别日志

### Fields

`Fields` 是 `map[string]interface{}` 的类型别名，用于结构化日志记录。使用 `logging.Fields` 而不是 `logrus.Fields` 可以避免直接依赖 logrus。

## 依赖

- `github.com/sirupsen/logrus` - 日志库
- `gopkg.in/natefinch/lumberjack.v2` - 日志轮转库

## 相关文档

- [详细架构文档](../../docs/architecture/logging.md) - 模块设计思路和实现细节
- [日志和调试规范](../../docs/development/references/logging.md) - 日志使用规范

# Logging 模块架构文档

## 📋 概述

Logging 模块是 Workflow CLI 的基础模块，提供统一的日志记录功能。该模块专注于日志系统的封装和管理，不涉及命令层的业务逻辑。

Logging 模块提供模块级别的日志管理、自动模块识别、文件日志输出、日志轮转等功能，总代码行数约 804 行。

**模块统计：**
- 代码行数：约 804 行（含测试文件：437 + 367）
- 主要文件：2 个核心文件
- 主要结构体：`errorLogHook`、`logConfig`、`LoggerEntry`
- 支持功能：模块级别日志、自动模块识别、文件日志输出、日志轮转、统一错误日志

**注意**：本模块是基础库模块，其他模块通过导入使用。所有 Lib 层模块必须使用此模块进行日志记录。

---

## 📁 模块架构（核心业务逻辑）

Logging 模块（`internal/logging/`）是 Workflow CLI 的基础库模块，提供统一的日志记录功能。该模块专注于日志系统的封装和管理，不涉及命令层的业务逻辑。

### 模块结构

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

**总计：约 804 行代码**

### 依赖模块

- **`github.com/sirupsen/logrus`**：日志库
  - 提供结构化日志记录功能
- **`gopkg.in/natefinch/lumberjack.v2`**：日志轮转库
  - 提供日志文件自动轮转功能（按大小、时间、备份数量）

### 模块集成

- **`internal/http/`**：HTTP 客户端模块
  - 使用 `GetLogger()` - 记录 HTTP 请求日志
- **`internal/infrastructure/http/`**：HTTP 基础设施模块
  - 使用 `WithField("module", "http")` - 标识 HTTP 模块日志
- **`internal/llm/`**：LLM 客户端模块
  - 使用 `GetLogger()` - 记录 LLM API 调用日志
- **`internal/config/`**：配置管理模块
  - 使用 `GetLogger()` - 记录配置操作日志

---

## 🏗️ 架构设计

### 设计原则

1. **模块隔离**：每个模块的日志输出到独立的文件，便于日志管理和排查
2. **自动识别**：通过 `runtime.Caller()` 自动识别调用者模块名，无需手动指定
3. **统一管理**：所有错误日志统一输出到 `error.log`，便于集中查看
4. **灵活配置**：支持多种日志格式（text/json）和输出目标（控制台/文件）
5. **线程安全**：使用 `sync.RWMutex` 保证模块 logger 创建的线程安全

### 核心组件

#### 1. 日志初始化 (logger.go)

**职责**：初始化日志系统，配置日志级别、格式和输出目标

**主要方法**：
- `InitWithFiles()` - 初始化日志系统（支持文件输出）
- `Init()` - 初始化日志系统（向后兼容，不创建文件日志）
- `SetLevel()` - 设置日志级别

**关键特性**：
- 支持文本和 JSON 两种格式
- 支持控制台和文件双重输出
- 自动创建日志目录
- 容错处理（创建目录失败不中断程序）

**使用场景**：
- 应用启动时初始化日志系统
- 动态调整日志级别

#### 2. 模块 Logger 管理 (logger.go)

**职责**：为每个模块创建独立的 logger 实例，管理模块级别的日志输出

**主要方法**：
- `GetLogger()` - 获取带模块名的 logger Entry（返回 `*LoggerEntry`）
- `getModuleLogger()` - 获取模块级别的 logger 实例
- `getCallerModule()` - 自动识别调用者模块名

**关键特性**：
- 自动识别模块名（从调用栈提取）
- 每个模块独立的 logger 实例
- 模块 logger 缓存（避免重复创建）
- 线程安全的双重检查锁定模式
- 返回 `*LoggerEntry` 封装类型，提供更好的封装性

**使用场景**：
- 各模块记录日志时自动获取模块名
- 按模块分别输出到不同文件

#### 3. 错误日志 Hook (logger.go)

**职责**：统一收集所有模块的错误日志到 `error.log`

**主要方法**：
- `Levels()` - 定义 Hook 处理的日志级别（ERROR/FATAL/PANIC）
- `Fire()` - 执行日志写入操作

**关键特性**：
- 只记录 ERROR 级别以上的日志
- 所有模块的错误统一输出到 `error.log`
- 使用独立的文件 writer，不影响模块日志文件

**使用场景**：
- 集中查看所有模块的错误日志
- 错误日志分析和监控

#### 4. 日志文件管理 (logger.go)

**职责**：管理日志文件的创建、轮转和清理

**关键特性**：
- 使用 `lumberjack` 实现日志轮转
- 支持按大小轮转（默认 10MB）
- 支持备份数量限制（默认 5 个）
- 支持保留时间限制（默认 30 天）
- 支持压缩旧日志文件

#### 5. LoggerEntry 封装类型 (logger.go)

**职责**：封装 `logrus.Entry`，提供更好的封装性和类型安全

**主要方法**：
- `WithField(key, value)` - 添加单个字段
- `WithFields(fields Fields)` - 添加多个字段（使用 `logging.Fields` 类型）
- `WithError(err)` - 添加错误字段
- `Debug/Debugf` - 记录 Debug 级别日志
- `Info/Infof` - 记录 Info 级别日志
- `Warn/Warnf` - 记录 Warn 级别日志
- `Error/Errorf` - 记录 Error 级别日志

**关键特性**：
- 封装 `logrus.Entry`，避免直接依赖 logrus
- 提供 `Fields` 类型别名，使用 `logging.Fields` 而不是 `logrus.Fields`
- 保持与 logrus API 的兼容性

**日志文件结构**：
```
~/.workflow/logs/
├── http.log          # HTTP 模块日志
├── llm.log           # LLM 模块日志
├── config.log        # Config 模块日志
├── error.log         # 统一错误日志
└── {module}.log      # 其他模块日志
```

### 设计模式

#### 1. 单例模式

**实现**：使用全局变量和互斥锁管理模块 logger 缓存

**优势**：
- 每个模块只有一个 logger 实例
- 避免重复创建，提高性能
- 线程安全

#### 2. 工厂模式

**实现**：`GetLogger()` 作为工厂方法，根据调用者自动创建对应的 logger

**优势**：
- 隐藏 logger 创建细节
- 自动识别模块名
- 统一的创建接口

#### 3. Hook 模式

**实现**：使用 logrus 的 Hook 机制实现错误日志统一收集

**优势**：
- 解耦错误日志收集逻辑
- 不影响模块日志的正常输出
- 易于扩展（可添加更多 Hook）

### 错误处理

#### 分层错误处理

1. **初始化层**：创建日志目录失败时，记录错误但不中断程序
2. **模块识别层**：无法识别模块名时，返回 "unknown"
3. **文件写入层**：文件写入失败时，由 logrus 和 lumberjack 处理

#### 容错机制

- **目录创建失败**：记录错误但不中断程序，继续使用控制台输出
- **模块识别失败**：返回 "unknown" 模块名，日志仍可正常输出
- **文件写入失败**：由底层库处理，不影响程序运行

---

## 🔄 集成关系

### 模块使用关系

Logging 模块被以下模块使用：

1. **`internal/http/`**：HTTP 客户端模块
   - 使用 `GetLogger()` - 记录 HTTP 请求和响应日志
   - 使用场景：请求发送前、响应后、错误处理

2. **`internal/infrastructure/http/`**：HTTP 适配器模块
   - 使用 `logging.WithField("module", "http")` - 标识 HTTP 模块日志（全局函数）
   - 使用场景：Resty 日志适配，将 Resty 的日志转发到 logging 包
   - 说明：适配器场景中无法使用 `GetLogger()` 自动识别模块名，因此使用全局函数手动指定

3. **`internal/llm/`**：LLM 客户端模块
   - 使用 `GetLogger()` - 记录 LLM API 调用日志
   - 使用场景：API 调用开始、成功、失败

4. **`internal/config/`**：配置管理模块
   - 使用 `GetLogger()` - 记录配置操作日志
   - 使用场景：配置加载、保存、错误处理

### 调用流程

#### 日志记录流程

```
调用模块代码
  ↓
logging.GetLogger()
  ↓
getCallerModule() - 自动识别模块名
  ↓
getModuleLogger(module) - 获取模块 logger
  ↓
检查缓存 → 存在则返回
  ↓
不存在则创建新 logger
  ↓
配置输出（控制台/文件）
  ↓
添加错误日志 Hook
  ↓
缓存 logger 并返回
  ↓
返回带 module 字段的 LoggerEntry
  ↓
调用者记录日志
```

#### 错误日志收集流程

```
模块记录错误日志
  ↓
模块 logger 处理
  ↓
写入模块日志文件
  ↓
触发 errorLogHook
  ↓
检查日志级别（ERROR/FATAL/PANIC）
  ↓
写入 error.log
```

---

## 🎯 核心功能

### 1. 自动模块识别

**功能说明**：通过 `runtime.Caller()` 自动识别调用者的模块名，无需手动指定

**流程**：
1. 获取调用栈信息
2. 跳过 logging 包本身的调用
3. 从函数名或文件路径提取模块名
4. 返回模块名（如 "http", "llm", "config"）

**示例**：
```go
import "github.com/zevwings/workflow/internal/logging"

// 在 http 模块中调用
logger := logging.GetLogger()  // 返回 *LoggerEntry
logger.Infof("HTTP request: %s", url)
// 自动识别模块名为 "http"
// 日志包含 module=http 字段
```

### 2. 模块级别日志文件

**功能说明**：每个模块的日志输出到独立的文件，便于日志管理和排查

**流程**：
1. 识别调用者模块名
2. 获取或创建模块 logger
3. 配置模块日志文件路径：`{logDir}/{module}.log`
4. 写入日志到模块文件

**示例**：
```go
// 初始化时指定日志目录
logging.InitWithFiles("info", "text", nil, "~/.workflow/logs", false)

// http 模块日志 → ~/.workflow/logs/http.log
// llm 模块日志 → ~/.workflow/logs/llm.log
// config 模块日志 → ~/.workflow/logs/config.log
```

### 3. 统一错误日志收集

**功能说明**：所有模块的错误日志统一输出到 `error.log`，便于集中查看

**流程**：
1. 模块 logger 记录日志
2. 错误日志 Hook 检查日志级别
3. 如果是 ERROR/FATAL/PANIC，写入 `error.log`
4. 同时写入模块日志文件

**示例**：
```go
logger := logging.GetLogger()
logger.Errorf("Operation failed: %v", err)
// 同时写入：
// - ~/.workflow/logs/{module}.log
// - ~/.workflow/logs/error.log
```

### 4. 日志轮转

**功能说明**：自动管理日志文件大小和保留时间，避免日志文件过大

**流程**：
1. 日志文件达到 10MB 时自动轮转
2. 保留最近 5 个备份文件
3. 删除超过 30 天的旧日志
4. 压缩旧日志文件

**配置**：
- MaxSize: 10MB
- MaxBackups: 5
- MaxAge: 30 天
- Compress: true

---

## 📋 使用示例

### 基本使用

```go
import "github.com/zevwings/workflow/internal/logging"

// 初始化日志系统
logging.InitWithFiles("info", "text", nil, "~/.workflow/logs", false)

// 获取 logger（自动识别模块名）
logger := logging.GetLogger()

// 记录日志
logger.Infof("Operation completed")
logger.Errorf("Operation failed: %v", err)
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
```

### 错误日志记录

```go
logger := logging.GetLogger()

err := doSomething()
if err != nil {
    // 使用 WithError 添加错误信息
    logger.WithError(err).Errorf("Operation failed")
}
```

---

## 📝 扩展性

### 添加新的日志格式

1. 在 `InitWithFiles()` 中添加新的格式判断
2. 创建对应的 Formatter
3. 设置到 logger

**示例**：
```go
case "custom":
    formatter = &CustomFormatter{
        TimestampFormat: "2006-01-02 15:04:05",
    }
```

### 添加新的日志 Hook

1. 实现 `logrus.Hook` 接口
2. 在 `getModuleLogger()` 中添加 Hook
3. 配置 Hook 的日志级别和处理逻辑

**示例**：
```go
type customHook struct {
    // ...
}

func (h *customHook) Levels() []logrus.Level {
    return []logrus.Level{logrus.InfoLevel}
}

func (h *customHook) Fire(entry *logrus.Entry) error {
    // 处理逻辑
    return nil
}
```

---

## 📚 相关文档

- [模块 README](../../internal/logging/README.md) - 基础使用说明
- [日志和调试规范](../development/references/logging.md) - 日志使用规范

---

## ✅ 总结

Logging 模块采用清晰的模块化设计：

1. **模块隔离**：每个模块独立的日志文件
2. **自动识别**：无需手动指定模块名
3. **统一管理**：错误日志集中收集
4. **灵活配置**：支持多种格式和输出方式
5. **线程安全**：保证并发安全

**设计优势**：
- ✅ 模块日志隔离，便于排查问题
- ✅ 自动识别模块名，使用简单
- ✅ 统一错误日志，便于监控
- ✅ 日志轮转，避免文件过大
- ✅ 线程安全，支持并发场景

**当前实现状态**：
- ✅ 模块级别日志文件
- ✅ 自动模块识别
- ✅ 统一错误日志收集
- ✅ 日志轮转功能
- ✅ 多种日志格式支持

---

**最后更新**: 2026-01-09

**更新说明**：
- 更新代码行数统计（logger.go: 437行，logger_test.go: 367行，总计 804行）
- 补充 `LoggerEntry` 封装类型的说明
- 更新 `GetLogger()` 返回类型说明（`*LoggerEntry` 而非 `*logrus.Entry`）
- 更新 `Fields` 类型使用说明（使用 `logging.Fields` 而非 `logrus.Fields`）
- 更新 HTTP 适配器使用方式说明

# 模块组织规范

> 本文档定义了 Workflow CLI 项目的模块组织规范和最佳实践，所有贡献者都应遵循这些规范。

---

## 📋 目录

- [概述](#-概述)
- [目录结构](#-目录结构)
- [模块职责](#-模块职责)
- [模块依赖规则](#-模块依赖规则)
- [平台特定代码组织](#-平台特定代码组织)
- [相关文档](#-相关文档)

---

## 📋 概述

本文档定义了模块组织规范，包括目录结构、模块职责、模块依赖规则和平台特定代码组织。

### 核心原则

- **标准布局**：遵循 Go 标准项目布局（[Standard Go Project Layout](https://github.com/golang-standards/project-layout)）
- **职责分离**：各层职责清晰，避免混淆
- **依赖管理**：避免循环依赖，保持依赖方向清晰

### 使用场景

- 创建新模块时参考
- 重构模块时使用
- 代码审查时检查

---

## 目录结构

遵循 Go 标准项目布局：

```
workflow/
├── cmd/
│   └── workflow/          # 主入口
│       └── main.go
├── internal/              # 内部包（不对外暴露）
│   ├── cli/               # CLI 根命令
│   │   └── root.go
│   ├── commands/          # 命令实现
│   │   ├── setup.go
│   │   ├── version.go
│   │   └── ...
│   ├── lib/               # 核心业务逻辑
│   │   ├── config/        # 配置管理
│   │   ├── http/          # HTTP 客户端
│   │   ├── git/           # Git 操作
│   │   ├── github/        # GitHub API
│   │   ├── jira/          # Jira API
│   │   └── ...
│   ├── logging/           # 日志系统
│   └── output/            # 输出格式化
├── pkg/                   # 公共包（可选，对外暴露）
├── go.mod
├── go.sum
└── Makefile
```

---

## 模块职责

### `cmd/` 目录

- **用途**：可执行文件的入口点
- **规则**：
  - 每个子目录代表一个可执行文件
  - 每个子目录应该只包含一个 `main.go` 文件
  - 不应该包含业务逻辑，只负责初始化

```go
// cmd/workflow/main.go
package main

import (
    "os"
    "github.com/zevwings/workflow/internal/cli"
)

func main() {
    if err := cli.Execute(); err != nil {
        os.Exit(1)
    }
}
```

### `internal/` 目录

- **用途**：内部包，不允许外部项目导入
- **规则**：
  - 所有业务逻辑应该放在 `internal/` 目录下
  - `internal/` 下的包只能被同一项目内的其他包导入
  - 外部项目无法导入 `internal/` 下的包

#### `internal/cli/` - CLI 根命令

- **职责**：CLI 根命令定义和初始化
- **包含**：根命令、命令注册、版本信息等

#### `internal/commands/` - 命令实现

- **职责**：CLI 命令的具体实现
- **包含**：各个子命令的实现（如 `setup`、`version`、`config` 等）

#### `internal/lib/` - 核心业务逻辑

- **职责**：可复用的业务逻辑模块
- **包含**：配置管理、HTTP 客户端、Git 操作、API 集成等

#### `internal/logging/` - 日志系统

- **职责**：日志系统的封装
- **包含**：日志初始化、日志级别管理、日志格式化等

#### `internal/output/` - 输出格式化

- **职责**：命令行输出的格式化
- **包含**：表格输出、颜色输出、格式化工具等

### `pkg/` 目录（可选）

- **用途**：公共包，允许外部项目导入
- **规则**：
  - 只有在需要被外部项目导入时才使用 `pkg/` 目录
  - 如果项目只是 CLI 工具，通常不需要 `pkg/` 目录

---

## 模块依赖规则

### 依赖方向

```
cmd/ → internal/cli → internal/commands → internal/lib
```

**规则**：
- **命令层** → **库层**：命令层可以依赖库层，但不能反向依赖
- **库层内部**：可以相互依赖，但避免循环依赖
- **基础模块**：`internal/lib/base/`（如果存在）不依赖其他业务模块

### 禁止的依赖

- ❌ `internal/lib/` 不能依赖 `internal/commands/`
- ❌ `internal/lib/` 不能依赖 `internal/cli/`
- ❌ `internal/lib/` 不能依赖 `cmd/`
- ❌ 避免循环依赖

### 示例

```go
// ✅ 好的依赖方向
// internal/commands/setup.go
package commands

import (
    "github.com/zevwings/workflow/internal/lib/config"  // ✅ 命令层依赖库层
)

// ❌ 不好的依赖方向
// internal/lib/config/manager.go
package config

import (
    "github.com/zevwings/workflow/internal/commands"  // ❌ 库层不能依赖命令层
)
```

---

## 平台特定代码组织

项目支持跨平台开发（macOS、Linux、Windows），需要正确处理平台特定代码。

### 使用构建标签组织平台特定代码

使用 Go 的构建标签（build tags）来组织平台特定代码：

```go
// +build darwin

package util

// macOS 特定的函数实现
func GetSystemPath() string {
    return "/usr/local/bin"
}
```

```go
// +build linux

package util

// Linux 特定的函数实现
func GetSystemPath() string {
    return "/usr/bin"
}
```

```go
// +build windows

package util

// Windows 特定的函数实现
func GetSystemPath() string {
    return "C:\\Program Files\\Workflow"
}
```

### 使用文件后缀组织平台特定代码

Go 支持使用文件后缀来组织平台特定代码：

```
internal/lib/util/
├── clipboard.go          # 跨平台接口
├── clipboard_darwin.go   # macOS 实现
├── clipboard_linux.go    # Linux 实现
└── clipboard_windows.go  # Windows 实现
```

**规则**：
- 使用 `_GOOS` 后缀（如 `_darwin.go`、`_linux.go`、`_windows.go`）
- 使用 `_GOARCH` 后缀（如 `_amd64.go`、`_arm64.go`）
- 使用 `_GOOS_GOARCH` 后缀（如 `_linux_arm64.go`）

### 平台特定导入规则

**一般规则**：所有导入语句都应该在文件顶部。

**例外情况**：如果导入只在特定平台使用，可以使用构建标签限制导入：

```go
// +build !windows

package util

import (
    "syscall"
)

// Windows 不支持 syscall，使用构建标签排除
```

### 平台特定代码的模块化组织

对于复杂的平台特定功能，建议使用独立的文件：

```
internal/lib/util/
├── clipboard.go          # 跨平台接口定义
├── clipboard_darwin.go  # macOS 实现
├── clipboard_linux.go   # Linux 实现
└── clipboard_windows.go # Windows 实现
```

在 `clipboard.go` 中定义接口：

```go
package util

// CopyToClipboard 复制文本到剪贴板
func CopyToClipboard(text string) error {
    return copyToClipboard(text)
}
```

在各个平台文件中实现：

```go
// clipboard_darwin.go
// +build darwin

package util

import (
    "github.com/atotto/clipboard"
)

func copyToClipboard(text string) error {
    return clipboard.WriteAll(text)
}
```

### 平台特定依赖管理

在 `go.mod` 中，Go 会自动处理平台特定依赖。某些依赖可能只在特定平台可用：

```go
// 使用构建标签限制导入
// +build !musl

package util

import (
    "github.com/atotto/clipboard"
)
```

**规则**：
- 使用构建标签明确说明平台限制
- 在文档中说明平台限制
- 确保所有平台的代码都能正常编译

### 平台检测工具

使用 `runtime` 包进行运行时平台检测：

```go
import (
    "runtime"
)

func DetectPlatform() string {
    return runtime.GOOS
}

func IsMacOS() bool {
    return runtime.GOOS == "darwin"
}

func IsLinux() bool {
    return runtime.GOOS == "linux"
}

func IsWindows() bool {
    return runtime.GOOS == "windows"
}
```

**使用场景**：
- **编译时检测**：使用构建标签，适用于编译期已知的平台差异
- **运行时检测**：使用 `runtime` 包，适用于需要在运行时判断平台的场景

### 平台特定功能的测试要求

**测试覆盖要求**：
- 平台特定功能必须添加单元测试
- 使用构建标签确保测试只在对应平台运行：

```go
// +build darwin

package util_test

import (
    "testing"
    "github.com/zevwings/workflow/internal/lib/util"
)

func TestMacOSSpecificFeature(t *testing.T) {
    // macOS 特定测试
}
```

**跨平台测试策略**：
- 在 CI/CD 中为所有支持的平台运行测试
- 使用 GitHub Actions 的矩阵构建策略测试多个平台
- 确保所有平台的测试都能通过

**测试注意事项**：
- 某些功能可能在某些平台上不可用
- 在测试中正确处理平台限制，避免在不支持的平台上运行相关测试

### 平台特定代码的最佳实践

1. **优先使用构建标签**：对于编译期已知的平台差异，使用构建标签
2. **运行时检测作为补充**：对于需要在运行时判断的场景，使用 `runtime` 包
3. **模块化组织**：将平台特定代码组织到独立的文件中，保持代码清晰
4. **文档说明**：在代码和文档中明确说明平台限制和平台特定行为
5. **测试覆盖**：确保所有平台的代码都有相应的测试覆盖

---

## 🔍 故障排除

### 问题 1：模块依赖混乱

**症状**：模块依赖方向不清晰，存在循环依赖

**解决方案**：

1. 检查模块依赖是否符合标准布局规则
2. 避免命令层和库层之间的反向依赖
3. 使用 `go mod graph` 分析依赖关系

### 问题 2：平台特定代码组织混乱

**症状**：平台特定代码分散在多个文件中，难以维护

**解决方案**：

1. 使用独立的文件组织平台特定代码
2. 使用构建标签明确平台限制
3. 使用文件后缀（`_GOOS.go`）组织平台特定代码

### 问题 3：循环依赖

**症状**：编译时出现循环依赖错误

**解决方案**：

1. 检查导入路径，确保依赖方向正确
2. 提取公共接口到独立包
3. 使用依赖注入避免直接依赖

---

## 📚 相关文档

### 开发规范

- [代码风格规范](./code-style.md) - 代码风格规范
- [命名规范](./naming.md) - 命名规范

### 架构文档

- [架构文档](../../architecture/README.md) - 项目架构总览

### Go 官方文档

- [Standard Go Project Layout](https://github.com/golang-standards/project-layout) - Go 标准项目布局
- [Build Constraints](https://pkg.go.dev/cmd/go#hdr-Build_constraints) - Go 构建标签文档

---

## ✅ 检查清单

使用本规范时，请确保：

- [ ] 遵循 Go 标准项目布局
- [ ] 模块职责清晰
- [ ] 模块依赖方向正确
- [ ] 平台特定代码组织清晰
- [ ] 平台特定导入使用构建标签标记
- [ ] 平台特定功能有测试覆盖

---

**最后更新**: 2025-01-27

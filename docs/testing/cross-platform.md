# 跨平台测试方案

> 本文档定义跨平台测试的策略和方法，确保代码在所有支持的平台上都能正确运行。

---

## 📋 目录

- [概述](#-概述)
- [支持的平台](#-支持的平台)
- [本地测试](#-本地测试)
- [CI/CD 跨平台测试](#-cicd-跨平台测试)

---

## 🎯 概述

Workflow CLI 项目支持多个平台和架构，需要确保代码在所有目标平台上都能正确运行。

### 测试目标

- ✅ **功能一致性**：确保所有平台上的功能行为一致
- ✅ **构建验证**：验证代码可以在所有目标平台上正确构建
- ✅ **平台兼容性**：发现和修复平台特定的问题

---

## 🖥️ 支持的平台

| 平台 | GOOS | GOARCH | 说明 |
|------|------|--------|------|
| macOS Intel | `darwin` | `amd64` | 原生支持 |
| macOS Apple Silicon | `darwin` | `arm64` | 原生支持 |
| Linux x86_64 | `linux` | `amd64` | glibc 动态链接 |
| Linux ARM64 | `linux` | `arm64` | ARM 架构 |
| Windows x86_64 | `windows` | `amd64` | MSVC 工具链 |
| Windows ARM64 | `windows` | `arm64` | ARM 架构 |

**注意**：Go 支持交叉编译，可以在一个平台上构建所有目标平台的二进制文件。

---

## 🔧 本地测试

### 前置要求

#### 安装 Go toolchain

确保已安装 Go 1.21+：

```bash
# 检查 Go 版本
go version

# 应该显示：go version go1.21.x ...
```

### 构建和测试命令

#### 构建可执行文件

```bash
# 在当前平台构建
go build -o bin/workflow ./cmd/workflow

# 指定目标平台（交叉编译）
GOOS=linux GOARCH=amd64 go build -o bin/workflow-linux-amd64 ./cmd/workflow
GOOS=windows GOARCH=amd64 go build -o bin/workflow-windows-amd64.exe ./cmd/workflow
GOOS=darwin GOARCH=arm64 go build -o bin/workflow-darwin-arm64 ./cmd/workflow
```

#### 运行测试

```bash
# 运行所有测试
go test ./...

# 运行单元测试
go test ./internal/...

# 运行集成测试
go test -tags=integration ./test/integration

# 运行特定包的测试
go test ./internal/lib/config
```

#### 跨平台构建测试

```bash
# 测试所有目标平台的构建
for os in darwin linux windows; do
  for arch in amd64 arm64; do
    echo "Building for $os/$arch..."
    GOOS=$os GOARCH=$arch go build -o bin/workflow-$os-$arch ./cmd/workflow
  done
done
```

---

## 🚀 CI/CD 跨平台测试

项目在 CI/CD 中使用矩阵策略并行运行所有原生平台的测试：

- **单元测试**：在所有平台上运行
- **集成测试**：在所有平台上运行（如适用）
- **构建验证**：验证所有目标平台可以成功构建

详细配置请参考 [`.github/workflows/ci.yml`](../../.github/workflows/ci.yml)（如果存在）。

### 支持的 CI 平台

- ✅ `ubuntu-latest` - Linux x86_64
- ✅ `macos-latest` - macOS (Intel 和 Apple Silicon)
- ✅ `windows-latest` - Windows x86_64

### CI/CD 测试矩阵示例

```yaml
# .github/workflows/ci.yml
name: CI

on: [push, pull_request]

jobs:
  test:
    strategy:
      matrix:
        os: [ubuntu-latest, macos-latest, windows-latest]
        go-version: ['1.21', '1.22']
    runs-on: ${{ matrix.os }}
    steps:
      - uses: actions/checkout@v3
      - uses: actions/setup-go@v4
        with:
          go-version: ${{ matrix.go-version }}
      - run: go test ./...
      - run: go build ./cmd/workflow
```

---

## 🔍 平台特定测试

### 使用构建标签

使用构建标签标记平台特定的测试：

```go
//go:build windows

package config

import "testing"

func TestWindowsSpecific(t *testing.T) {
    // Windows 特定测试
}
```

```go
//go:build !windows

package config

import "testing"

func TestUnixSpecific(t *testing.T) {
    // Unix 特定测试（Linux 和 macOS）
}
```

### 运行平台特定测试

```bash
# 运行所有测试（包括平台特定测试）
go test ./...

# 只运行 Windows 特定测试
go test -tags=windows ./...

# 只运行 Unix 特定测试
go test -tags=!windows ./...
```

---

## 🛠️ 跨平台测试工具

### 使用 Docker 测试 Linux

```bash
# 在 Docker 容器中运行 Linux 测试
docker run --rm -v $(pwd):/work -w /work golang:1.21 go test ./...
```

### 使用 GitHub Actions 测试多平台

GitHub Actions 提供免费的 macOS、Linux 和 Windows 运行器，可以轻松测试多平台。

---

## 📋 跨平台测试检查清单

### 功能测试

- [ ] 所有核心功能在所有平台上都能正常工作
- [ ] 文件路径处理在所有平台上都正确
- [ ] 环境变量处理在所有平台上都一致
- [ ] 命令行参数解析在所有平台上都正确

### 构建测试

- [ ] 所有目标平台都能成功构建
- [ ] 构建的二进制文件能在目标平台上运行
- [ ] 构建的二进制文件大小合理

### 集成测试

- [ ] 集成测试在所有平台上都能通过
- [ ] 外部依赖（如 Git）在所有平台上都能正常工作

---

## 🐛 常见问题和解决方案

### 问题 1：路径分隔符不一致

**问题**：Windows 使用 `\`，Unix 使用 `/`

**解决方案**：使用 `filepath.Join()` 和 `filepath.Separator`

```go
import "path/filepath"

// ✅ 正确：使用 filepath.Join
path := filepath.Join("dir", "file.txt")

// ❌ 错误：硬编码路径分隔符
path := "dir/file.txt" // Windows 上会失败
```

### 问题 2：环境变量差异

**问题**：不同平台的环境变量名称可能不同

**解决方案**：使用 `os.Getenv()` 并提供默认值

```go
import "os"

homeDir := os.Getenv("HOME") // Unix
if homeDir == "" {
    homeDir = os.Getenv("USERPROFILE") // Windows
}
```

### 问题 3：文件权限差异

**问题**：Windows 和 Unix 的文件权限模型不同

**解决方案**：使用 `os.Chmod()` 并提供平台特定的处理

```go
import (
    "os"
    "runtime"
)

if runtime.GOOS != "windows" {
    os.Chmod(file, 0755)
}
```

---

## 📚 相关文档

- [测试组织规范](./organization.md) - 测试组织结构和命名约定
- [测试编写规范](./writing.md) - 测试编写的具体规范
- [CI/CD 工作流](../guidelines/ci-workflow.md) - CI/CD 工作流说明

---

**最后更新**: 2025-01-28

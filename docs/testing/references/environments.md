# 测试环境工具指南

> 本文档介绍 Go 标准库提供的测试环境工具使用方法，包括临时目录管理、环境变量隔离和测试清理。

---

## 📋 目录

- [概述](#-概述)
- [Go 标准库测试环境](#-go-标准库测试环境)
- [最佳实践](#-最佳实践)

---

## 📋 概述

Go 标准库提供了强大的测试环境工具，可以创建完全隔离的测试环境：

- **临时目录管理**：使用 `t.TempDir()` 创建临时目录（自动清理）
- **环境变量隔离**：使用 `t.Setenv()` 设置环境变量（自动恢复）
- **测试清理**：使用 `t.Cleanup()` 注册清理函数（自动执行）

### 核心特性

- ✅ **完全隔离**：每个测试运行在独立的临时目录中，不会影响实际系统
- ✅ **自动清理**：测试结束后自动清理临时文件和恢复环境变量
- ✅ **线程安全**：支持并行测试执行
- ✅ **简单易用**：使用 Go 标准库，无需额外依赖

---

## Go 标准库测试环境

### 1. 临时目录管理

使用 `t.TempDir()` 创建临时目录，测试结束后自动清理：

```go
import (
    "os"
    "path/filepath"
    "testing"
)

func TestWithTempDir(t *testing.T) {
    // 创建临时目录（自动清理）
    tempDir := t.TempDir()

    // 在临时目录中创建文件
    filePath := filepath.Join(tempDir, "test.txt")
    err := os.WriteFile(filePath, []byte("test"), 0644)
    if err != nil {
        t.Fatalf("Failed to write file: %v", err)
    }

    // 测试代码
    // ...

    // 测试结束后，tempDir 会自动清理
}
```

### 2. 环境变量隔离

使用 `t.Setenv()` 设置环境变量，测试结束后自动恢复：

```go
import (
    "os"
    "testing"
)

func TestWithEnvVar(t *testing.T) {
    // 设置环境变量（自动恢复）
    t.Setenv("HOME", "/test/home")
    t.Setenv("CONFIG_DIR", "/test/config")

    // 测试代码可以使用环境变量
    homeDir := os.Getenv("HOME")
    if homeDir != "/test/home" {
        t.Errorf("Expected HOME=/test/home, got %s", homeDir)
    }

    // 测试结束后，环境变量会自动恢复
}
```

### 3. 测试清理

使用 `t.Cleanup()` 注册清理函数，测试结束后自动执行：

```go
import (
    "os"
    "testing"
)

func TestWithCleanup(t *testing.T) {
    // 创建临时文件
    tmpFile, err := os.CreateTemp("", "test-*.txt")
    if err != nil {
        t.Fatalf("Failed to create temp file: %v", err)
    }

    // 注册清理函数（自动执行）
    t.Cleanup(func() {
        os.Remove(tmpFile.Name())
    })

    // 测试代码
    // ...

    // 测试结束后，清理函数会自动执行
}
```

### 4. 组合使用

```go
import (
    "os"
    "path/filepath"
    "testing"
)

func TestCompleteIsolation(t *testing.T) {
    // 创建临时目录
    tempDir := t.TempDir()

    // 设置环境变量
    t.Setenv("HOME", tempDir)
    t.Setenv("CONFIG_DIR", filepath.Join(tempDir, ".config"))

    // 创建配置文件
    configPath := filepath.Join(tempDir, ".config", "config.toml")
    os.MkdirAll(filepath.Dir(configPath), 0755)
    os.WriteFile(configPath, []byte("[test]\nkey = \"value\""), 0644)

    // 测试代码
    // ...

    // 所有资源都会自动清理
}
```

### 5. Git 测试环境示例

使用标准库创建 Git 测试环境：

```go
import (
    "os/exec"
    "path/filepath"
    "testing"
)

func TestGitCommand(t *testing.T) {
    // 创建临时目录
    tempDir := t.TempDir()

    // 设置环境变量
    t.Setenv("HOME", tempDir)

    // 初始化 Git 仓库
    err := exec.Command("git", "init", tempDir).Run()
    if err != nil {
        t.Fatalf("Failed to init git repo: %v", err)
    }

    // 配置 Git 用户
    exec.Command("git", "config", "user.name", "Test User").Run()
    exec.Command("git", "config", "user.email", "test@example.com").Run()

    // 创建文件并提交
    filePath := filepath.Join(tempDir, "test.txt")
    os.WriteFile(filePath, []byte("content"), 0644)
    exec.Command("git", "add", "test.txt").Run()
    exec.Command("git", "commit", "-m", "Initial commit").Run()

    // 执行 Git 命令测试
    // ...
}
```

---

## 最佳实践

### 1. 使用 t.TempDir() 而不是 os.MkdirTemp()

```go
// ✅ 推荐：使用 t.TempDir()（自动清理）
func TestExample(t *testing.T) {
    tempDir := t.TempDir()
    // 使用 tempDir
}

// ❌ 不推荐：使用 os.MkdirTemp()（需要手动清理）
func TestExample(t *testing.T) {
    tempDir, err := os.MkdirTemp("", "test-*")
    if err != nil {
        t.Fatal(err)
    }
    defer os.RemoveAll(tempDir) // 容易忘记
    // 使用 tempDir
}
```

### 2. 使用 t.Setenv() 而不是 os.Setenv()

```go
// ✅ 推荐：使用 t.Setenv()（自动恢复）
func TestExample(t *testing.T) {
    t.Setenv("HOME", "/test/home")
    // 使用环境变量
}

// ❌ 不推荐：使用 os.Setenv()（需要手动恢复）
func TestExample(t *testing.T) {
    oldHome := os.Getenv("HOME")
    os.Setenv("HOME", "/test/home")
    defer os.Setenv("HOME", oldHome) // 容易忘记
    // 使用环境变量
}
```

### 3. 使用 t.Cleanup() 注册清理函数

```go
// ✅ 推荐：使用 t.Cleanup()（自动执行）
func TestExample(t *testing.T) {
    resource := setupResource(t)
    t.Cleanup(func() {
        resource.Cleanup()
    })
    // 使用 resource
}

// ❌ 不推荐：使用 defer（可能在某些情况下不执行）
func TestExample(t *testing.T) {
    resource := setupResource(t)
    defer resource.Cleanup() // 在某些情况下可能不执行
    // 使用 resource
}
```

### 4. 测试之间相互独立

```go
// ✅ 推荐：每个测试独立设置环境
func Test1(t *testing.T) {
    tempDir := t.TempDir()
    t.Setenv("HOME", tempDir)
    // 测试代码
}

func Test2(t *testing.T) {
    tempDir := t.TempDir()
    t.Setenv("HOME", tempDir)
    // 测试代码（不依赖 Test1）
}

// ❌ 不推荐：测试之间共享状态
var sharedDir string

func Test1(t *testing.T) {
    sharedDir = t.TempDir()
    // 测试代码
}

func Test2(t *testing.T) {
    // 依赖 Test1 的 sharedDir（不推荐）
}
```

### 5. 使用环境变量获取路径

```go
// ✅ 推荐：使用环境变量获取路径（支持测试隔离）
func TestExample(t *testing.T) {
    tempDir := t.TempDir()
    t.Setenv("HOME", tempDir)

    homeDir := os.Getenv("HOME")
    configDir := filepath.Join(homeDir, ".config")
    // 使用路径
}

// ❌ 不推荐：直接使用系统路径
func TestExample(t *testing.T) {
    homeDir, err := os.UserHomeDir()
    if err != nil {
        t.Fatal(err)
    }
    // 不支持测试隔离
}
```

---

## 相关文档

- [测试辅助工具指南](./helpers.md) - 测试辅助工具详细使用方法
- [测试编写规范](../writing.md) - 测试编写规范
- [测试组织规范](../organization.md) - 测试组织结构

---

**最后更新**: 2025-01-28

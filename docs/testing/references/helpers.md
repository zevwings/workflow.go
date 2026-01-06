# 测试辅助工具指南

> 本文档介绍测试辅助工具的使用方法，包括测试辅助函数、测试数据生成和路径获取函数。

---

## 📋 目录

- [概述](#-概述)
- [路径获取函数](#1-路径获取函数)
- [测试数据加载](#2-测试数据加载)
- [CLI 命令测试](#3-cli-命令测试)
- [最佳实践](#4-最佳实践)

---

## 📋 概述

测试辅助工具提供便捷的测试辅助功能，简化测试代码编写：

- **路径获取函数**：统一的路径获取函数，支持测试环境隔离
- **测试数据加载**：从 `testdata/` 目录加载测试数据
- **CLI 命令测试**：简化 CLI 命令的执行和断言

---

## 1. 路径获取函数

提供统一的路径获取函数，支持测试环境隔离和跨平台兼容性。

### 1.1 基本使用

```go
import (
    "testing"
    "github.com/your-org/workflow/testutils"
)

func TestWithPaths(t *testing.T) {
    // 获取测试主目录（支持环境变量隔离）
    homeDir := testutils.TestHomeDir(t)
    configDir := testutils.TestConfigDir(t)
    dataDir := testutils.TestDataDir(t)
    cacheDir := testutils.TestCacheDir(t)

    // 使用测试目录
    configPath := filepath.Join(configDir, "config.toml")
    // ...
}
```

### 1.2 可用函数

| 函数 | 说明 | 环境变量优先级 |
|------|------|----------------|
| `TestHomeDir(t)` | 获取主目录 | `HOME` > `USERPROFILE` (Windows) |
| `TestConfigDir(t)` | 获取配置目录 | `XDG_CONFIG_HOME` > `HOME/.config` > `APPDATA` (Windows) |
| `TestDataDir(t)` | 获取数据目录 | `XDG_DATA_HOME` > `HOME/.local/share` > `APPDATA` (Windows) |
| `TestCacheDir(t)` | 获取缓存目录 | `XDG_CACHE_HOME` > `HOME/.cache` > `LOCALAPPDATA` (Windows) |

### 1.3 环境变量隔离

路径获取函数优先使用环境变量，支持测试环境隔离：

```go
func TestWithEnvIsolation(t *testing.T) {
    // 设置环境变量
    tempDir := t.TempDir()
    t.Setenv("HOME", tempDir)
    t.Setenv("XDG_CONFIG_HOME", filepath.Join(tempDir, ".config"))

    // 路径获取函数会使用环境变量
    homeDir := testutils.TestHomeDir(t)
    assert.Equal(t, tempDir, homeDir)

    configDir := testutils.TestConfigDir(t)
    assert.Contains(t, configDir, tempDir)
}
```

### 1.4 跨平台兼容性

路径获取函数自动处理平台差异：

```go
func TestCrossPlatform(t *testing.T) {
    // 在不同平台上都能正常工作
    homeDir := testutils.TestHomeDir(t)
    configDir := testutils.TestConfigDir(t)

    // Windows: C:\Users\username\.config
    // Unix: /home/username/.config
    // 函数自动处理平台差异
}
```

---

## 2. 测试数据加载

### 2.1 加载 Fixtures

从 `testdata/fixtures/` 目录加载测试数据：

```go
import (
    "testing"
    "github.com/your-org/workflow/testutils"
)

func TestLoadFixture(t *testing.T) {
    // 加载测试数据
    data := testutils.LoadFixture(t, "sample_github_pr.json")

    // 使用测试数据
    var pr GitHubPR
    err := json.Unmarshal(data, &pr)
    assert.NoError(t, err)
    assert.NotNil(t, pr)
}
```

### 2.2 加载文本文件

```go
func TestLoadTextFixture(t *testing.T) {
    // 加载文本文件
    content := testutils.LoadTextFixture(t, "sample_pr_body.md")

    // 使用文本内容
    assert.Contains(t, content, "PR Title")
}
```

### 2.3 加载二进制文件

```go
func TestLoadBinaryFixture(t *testing.T) {
    // 加载二进制文件
    data := testutils.LoadBinaryFixture(t, "sample_image.png")

    // 使用二进制数据
    assert.NotEmpty(t, data)
}
```

### 2.4 测试数据目录结构

```
testdata/
├── fixtures/
│   ├── sample_github_pr.json
│   ├── sample_jira_response.json
│   └── sample_pr_body.md
└── integration/
    └── workflow_scenarios.json
```

---

## 3. CLI 命令测试

### 3.1 执行 CLI 命令

```go
import (
    "testing"
    "github.com/your-org/workflow/testutils"
)

func TestCLICommand(t *testing.T) {
    // 执行 CLI 命令
    output, err := testutils.ExecuteCommand(t, "workflow", "version")
    assert.NoError(t, err)
    assert.Contains(t, output, "version")
}
```

### 3.2 带参数的命令

```go
func TestCLICommandWithArgs(t *testing.T) {
    // 执行带参数的命令
    output, err := testutils.ExecuteCommand(t, "workflow", "config", "show")
    assert.NoError(t, err)
    assert.Contains(t, output, "config")
}
```

### 3.3 设置环境变量

```go
func TestCLICommandWithEnv(t *testing.T) {
    // 设置环境变量
    env := map[string]string{
        "HOME": t.TempDir(),
    }

    // 执行命令
    output, err := testutils.ExecuteCommandWithEnv(t, env, "workflow", "config", "show")
    assert.NoError(t, err)
}
```

### 3.4 设置工作目录

```go
func TestCLICommandWithDir(t *testing.T) {
    // 设置工作目录
    workDir := t.TempDir()

    // 执行命令
    output, err := testutils.ExecuteCommandWithDir(t, workDir, "workflow", "version")
    assert.NoError(t, err)
}
```

### 3.5 捕获标准输出和错误

```go
func TestCLICommandCaptureOutput(t *testing.T) {
    // 执行命令并捕获输出
    result := testutils.ExecuteCommandCapture(t, "workflow", "version")

    assert.NoError(t, result.Err)
    assert.Contains(t, result.Stdout, "version")
    assert.Empty(t, result.Stderr)
}
```

---

## 4. 最佳实践

### 4.1 使用路径获取函数

```go
// ✅ 推荐：使用路径获取函数
func TestExample(t *testing.T) {
    homeDir := testutils.TestHomeDir(t)
    configDir := testutils.TestConfigDir(t)
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

### 4.2 使用测试数据加载

```go
// ✅ 推荐：使用 LoadFixture
func TestExample(t *testing.T) {
    data := testutils.LoadFixture(t, "sample_github_pr.json")
    // 使用数据
}

// ❌ 不推荐：硬编码路径
func TestExample(t *testing.T) {
    data, err := os.ReadFile("testdata/fixtures/sample_github_pr.json")
    if err != nil {
        t.Fatal(err)
    }
    // 路径可能不存在
}
```

### 4.3 使用 CLI 命令测试

```go
// ✅ 推荐：使用 ExecuteCommand
func TestExample(t *testing.T) {
    output, err := testutils.ExecuteCommand(t, "workflow", "version")
    assert.NoError(t, err)
}

// ❌ 不推荐：手动执行命令
func TestExample(t *testing.T) {
    cmd := exec.Command("workflow", "version")
    output, err := cmd.Output()
    if err != nil {
        t.Fatal(err)
    }
    // 需要手动处理错误和输出
}
```

### 4.4 测试环境隔离

```go
// ✅ 推荐：使用环境变量隔离
func TestExample(t *testing.T) {
    tempDir := t.TempDir()
    t.Setenv("HOME", tempDir)

    homeDir := testutils.TestHomeDir(t)
    assert.Equal(t, tempDir, homeDir)
}

// ❌ 不推荐：使用真实系统路径
func TestExample(t *testing.T) {
    homeDir := testutils.TestHomeDir(t)
    // 可能使用真实系统路径，污染系统
}
```

---

## 相关文档

- [测试环境工具指南](./environments.md) - 测试环境工具详细使用方法
- [测试编写规范](../writing.md) - 测试编写规范
- [测试组织规范](../organization.md) - 测试组织结构

---

**最后更新**: 2025-01-28

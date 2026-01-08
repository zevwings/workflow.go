# testutils 测试辅助工具包

> 提供测试辅助功能，简化测试代码编写，支持测试环境隔离和跨平台兼容性。

**⚠️ 重要提示**：此包使用构建标签 `test`，只在测试时可用，不会被打包到 release 中。

---

## 📋 目录

- [概述](#概述)
- [构建标签](#构建标签)
- [路径获取函数](#路径获取函数)
- [测试数据加载](#测试数据加载)
- [CLI 命令测试](#cli-命令测试)
- [使用示例](#使用示例)

---

## 概述

`testutils` 包提供以下功能：

- **路径获取函数**：统一的路径获取函数，支持测试环境隔离
- **测试数据加载**：从 `testdata/` 目录加载测试数据
- **CLI 命令测试**：简化 CLI 命令的执行和断言

---

## 构建标签

`testutils` 包使用构建标签 `//go:build test`，确保只在测试时编译，不会被打包到 release 中。

### 运行测试

使用 `-tags=test` 标签运行测试：

```bash
# 运行所有测试
go test -tags=test ./...

# 运行特定包的测试
go test -tags=test ./internal/config

# 使用 Makefile（已自动包含 -tags=test）
make test
make test-coverage
```

### 验证不会被打包

正常构建时（不带 `-tags=test`），testutils 不会被包含：

```bash
# 正常构建（不包含 testutils）
go build ./cmd/workflow  # ✅ 成功

# 尝试编译 testutils（不带标签）
go build ./internal/testutils  # ❌ 失败：build constraints exclude all Go files
```

### 为什么使用构建标签？

- ✅ **避免打包到 release**：testutils 是测试专用工具，不应该出现在生产代码中
- ✅ **编译时检查**：如果生产代码错误导入了 testutils，编译会失败，及时发现问题
- ✅ **清晰的职责分离**：明确区分测试代码和生产代码

---

## 路径获取函数

提供统一的路径获取函数，支持测试环境隔离和跨平台兼容性。

### 可用函数

| 函数 | 说明 | 环境变量优先级 |
|------|------|----------------|
| `TestHomeDir(t)` | 获取主目录 | `HOME` > `USERPROFILE` (Windows) |
| `TestConfigDir(t)` | 获取配置目录 | `XDG_CONFIG_HOME` > `HOME/.config` > `APPDATA` (Windows) |
| `TestDataDir(t)` | 获取数据目录 | `XDG_DATA_HOME` > `HOME/.local/share` > `APPDATA` (Windows) |
| `TestCacheDir(t)` | 获取缓存目录 | `XDG_CACHE_HOME` > `HOME/.cache` > `LOCALAPPDATA` (Windows) |
| `TestWorkflowConfigDir(t)` | 获取 Workflow 配置目录 | `HOME/.workflow` |

### 使用示例

```go
import (
    "testing"
    "path/filepath"
    "github.com/zevwings/workflow/internal/testutils"
)

func TestWithPaths(t *testing.T) {
    // 获取测试主目录（支持环境变量隔离）
    homeDir := testutils.TestHomeDir(t)
    configDir := testutils.TestConfigDir(t)
    dataDir := testutils.TestDataDir(t)
    cacheDir := testutils.TestCacheDir(t)
    workflowConfigDir := testutils.TestWorkflowConfigDir(t)

    // 使用测试目录
    configPath := filepath.Join(configDir, "config.toml")
    // ...
}
```

### 环境变量隔离

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

---

## 测试数据加载

从 `testdata/fixtures/` 目录加载测试数据。

### 可用函数

| 函数 | 说明 |
|------|------|
| `LoadFixture(t, filename)` | 加载测试数据文件（二进制） |
| `LoadTextFixture(t, filename)` | 加载测试文本文件 |
| `LoadBinaryFixture(t, filename)` | 加载测试二进制文件（别名） |

### 使用示例

```go
import (
    "encoding/json"
    "testing"
    "github.com/zevwings/workflow/internal/testutils"
    "github.com/stretchr/testify/assert"
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

func TestLoadTextFixture(t *testing.T) {
    // 加载文本文件
    content := testutils.LoadTextFixture(t, "sample_pr_body.md")

    // 使用文本内容
    assert.Contains(t, content, "PR Title")
}
```

### 测试数据目录结构

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

## CLI 命令测试

简化 CLI 命令的执行和断言。

### 可用函数

| 函数 | 说明 |
|------|------|
| `ExecuteCommand(t, command, args...)` | 执行 CLI 命令 |
| `ExecuteCommandWithEnv(t, env, command, args...)` | 执行 CLI 命令（带环境变量） |
| `ExecuteCommandWithDir(t, dir, command, args...)` | 执行 CLI 命令（带工作目录） |
| `ExecuteCommandCapture(t, command, args...)` | 执行 CLI 命令并捕获输出 |
| `ExecuteCommandCaptureWithEnv(t, env, command, args...)` | 执行 CLI 命令并捕获输出（带环境变量） |
| `ExecuteCommandCaptureWithDir(t, dir, command, args...)` | 执行 CLI 命令并捕获输出（带工作目录） |

### 使用示例

```go
import (
    "testing"
    "github.com/zevwings/workflow/internal/testutils"
    "github.com/stretchr/testify/assert"
)

func TestCLICommand(t *testing.T) {
    // 执行 CLI 命令
    output, err := testutils.ExecuteCommand(t, "workflow", "version")
    assert.NoError(t, err)
    assert.Contains(t, output, "version")
}

func TestCLICommandWithEnv(t *testing.T) {
    // 设置环境变量
    env := map[string]string{
        "HOME": t.TempDir(),
    }

    // 执行命令
    output, err := testutils.ExecuteCommandWithEnv(t, env, "workflow", "config", "show")
    assert.NoError(t, err)
    assert.Contains(t, output, "config")
}

func TestCLICommandCapture(t *testing.T) {
    // 执行命令并捕获输出
    result := testutils.ExecuteCommandCapture(t, "workflow", "version")

    assert.NoError(t, result.Err)
    assert.Contains(t, result.Stdout, "version")
    assert.Empty(t, result.Stderr)
}
```

---

## 使用示例

### 完整示例：测试配置管理

```go
package config_test

import (
    "path/filepath"
    "testing"
    "github.com/zevwings/workflow/internal/config"
    "github.com/zevwings/workflow/internal/testutils"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
)

func TestGlobalManager(t *testing.T) {
    // 设置测试环境
    tempDir := t.TempDir()
    t.Setenv("HOME", tempDir)

    // 使用 testutils 获取路径
    configDir := testutils.TestWorkflowConfigDir(t)
    configPath := filepath.Join(configDir, "config.toml")

    // 创建配置管理器
    manager, err := config.NewGlobalManager()
    require.NoError(t, err)

    // 测试配置加载
    err = manager.Load()
    assert.NoError(t, err)

    // 验证配置文件路径
    assert.Equal(t, configPath, manager.GetConfigPath())
}
```

### 完整示例：测试 HTTP 客户端

```go
package http_test

import (
    "net/http"
    "net/http/httptest"
    "testing"
    "github.com/zevwings/workflow/internal/http"
    "github.com/stretchr/testify/assert"
)

func TestHTTPClient(t *testing.T) {
    // 创建 Mock 服务器
    server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.WriteHeader(http.StatusOK)
        w.Write([]byte(`{"id": 123}`))
    }))
    defer server.Close()

    // 创建 HTTP 客户端
    client := http.NewClient()

    // 执行请求
    resp, err := client.Get(server.URL)
    assert.NoError(t, err)
    assert.Equal(t, http.StatusOK, resp.StatusCode())
}
```

---

## 最佳实践

### 1. 使用路径获取函数

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

### 2. 使用测试数据加载

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

### 3. 测试环境隔离

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

- [测试环境工具指南](../../docs/testing/references/environments.md) - 测试环境工具详细使用方法
- [测试辅助工具指南](../../docs/testing/references/helpers.md) - 测试辅助工具详细使用方法
- [测试编写规范](../../docs/testing/writing.md) - 测试编写规范

---

**最后更新**: 2025-01-28


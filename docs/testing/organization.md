# 测试组织规范

> 本文档定义测试组织结构、命名约定和共享工具使用规范。

---

## 📋 目录

- [测试类型](#-测试类型)
- [测试组织结构](#-测试组织结构)
- [测试文件命名约定](#-测试文件命名约定)
- [共享测试工具](#-共享测试工具)
- [测试数据管理](#-测试数据管理)
- [测试组织最佳实践](#-测试组织最佳实践)

---

## 🎯 测试类型

### 1. 单元测试 (Unit Tests)

- **位置**：与源代码在同一包中，使用 `*_test.go` 文件
- **测试对象**：**测试包内的所有函数（包括私有和公开）**
- **特点**：快速执行，最小依赖
- **组织方式**：使用 `*_test.go` 文件，与源码文件同目录

**重要规则**：
- ✅ **可以测试包内的所有函数**（包括私有函数）
- ✅ **主要测试公共 API 和关键私有函数**
- ✅ 快速执行，不依赖外部环境
- ✅ 使用表驱动测试提高覆盖率

> 📖 **详细指南**：请参考 [单元测试指南](./references/unit-tests.md) 了解单元测试的详细编写规范、组织方式和最佳实践。

### 2. 集成测试 (Integration Tests)

- **位置**：单独的测试文件或 `test/integration/` 目录
- **测试对象**：**跨包的集成场景和端到端流程**
- **特点**：测试多个模块的交互，可能需要外部依赖
- **组织方式**：使用构建标签 `//go:build integration` 标记

```go
// test/integration/workflow_test.go
//go:build integration

package integration

import (
    "testing"
    "github.com/stretchr/testify/assert"
)

func TestWorkflowIntegration(t *testing.T) {
    // 测试完整的工作流程
    // 可能涉及多个包和外部服务
}
```

**重要规则**：
- ✅ **测试跨包的集成场景**
- ✅ **测试端到端流程**
- ✅ 可以使用外部依赖（数据库、API等）
- ✅ 使用构建标签区分单元测试和集成测试

### 3. 表驱动测试 (Table-Driven Tests)

Go 推荐使用表驱动测试，提高测试覆盖率和可维护性。

```go
func TestParseTicketID(t *testing.T) {
    tests := []struct {
        name     string
        input    string
        expected string
        wantErr  bool
    }{
        {
            name:     "valid ticket ID",
            input:    "PROJ-123",
            expected: "PROJ-123",
            wantErr:  false,
        },
        {
            name:     "invalid format",
            input:    "invalid",
            expected: "",
            wantErr:  true,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            result, err := ParseTicketID(tt.input)
            if tt.wantErr {
                assert.Error(t, err)
            } else {
                assert.NoError(t, err)
                assert.Equal(t, tt.expected, result)
            }
        })
    }
}
```

---

## 📁 测试组织结构

### 当前测试结构

本项目采用 **Go 标准测试结构**，测试文件与源码文件同目录：

```
internal/
├── lib/
│   ├── config/
│   │   ├── manager.go           # 源代码
│   │   └── manager_test.go      # 单元测试
│   ├── http/
│   │   ├── client.go
│   │   └── client_test.go
│   └── prompt/
│       ├── input.go
│       └── input_test.go
├── commands/
│   ├── check.go
│   └── check_test.go
├── cli/
│   ├── root.go
│   └── root_test.go
testdata/                          # 测试数据目录
├── fixtures/
│   ├── sample_github_pr.json
│   └── sample_jira_response.json
└── integration/                   # 集成测试数据
testutils/                         # 共享测试工具（可选）
├── helpers.go                     # 测试辅助函数
├── mock.go                        # Mock 工具
└── testdata.go                    # 测试数据生成
test/                              # 集成测试（可选）
└── integration/
    └── workflow_test.go
```

### 结构说明

- **测试文件位置**：测试文件与源码文件同目录，使用 `*_test.go` 后缀
- **测试包名**：测试文件使用与源码相同的包名（可以访问私有函数）
- **测试数据**：`testdata/` 目录存放测试用的示例数据（Go 会自动忽略此目录）
- **共享工具**：`testutils/` 或 `internal/testutils` 目录存放共享的测试辅助函数
- **集成测试**：`test/` 目录或使用构建标签标记的集成测试文件

---

## 📝 测试文件命名约定

### 命名规则

1. **与源码文件对应**：测试文件名 = 源码文件名 + `_test.go`
2. **使用下划线分隔**：使用下划线（`_`）分隔单词
3. **保持简洁**：避免不必要的后缀

### 命名示例

```go
// 源代码文件 → 测试文件
internal/lib/config/manager.go          → internal/lib/config/manager_test.go
internal/http/client.go                 → internal/http/client_test.go
internal/commands/check.go              → internal/commands/check_test.go
internal/cli/root.go                    → internal/cli/root_test.go
```

### 不推荐的命名

- ❌ `manager_testing.go` - 包含不必要的前缀
- ❌ `test_manager.go` - 不符合 Go 命名规范
- ❌ `manager_test_suite.go` - 过于复杂

---

## 🛠️ 共享测试工具

### testutils 目录结构

共享的测试工具应放在 `testutils/` 或 `internal/testutils/` 目录。该目录采用模块化组织，按功能分类：

```
testutils/
├── helpers.go              # 通用辅助函数
│   ├── TestHomeDir()       # 测试主目录
│   ├── TestConfigDir()     # 测试配置目录
│   └── TestDataDir()        # 测试数据目录
├── mock.go                  # Mock 工具
│   ├── MockHTTPServer()    # HTTP Mock 服务器
│   └── MockGitHubAPI()     # GitHub API Mock
├── testdata.go             # 测试数据生成
│   ├── LoadFixture()       # 加载测试数据
│   └── GenerateTestData()  # 生成测试数据
├── environment.go          # 测试环境管理
│   ├── SetupTestEnv()      # 设置测试环境
│   └── CleanupTestEnv()    # 清理测试环境
└── cli.go                  # CLI 测试辅助
    ├── ExecuteCommand()    # 执行 CLI 命令
    └── CaptureOutput()     # 捕获输出
```

### 核心模块说明

#### 1. 测试辅助函数 (`helpers.go`)

提供通用的测试辅助函数：

- **路径获取函数**：`TestHomeDir()`, `TestConfigDir()`, `TestDataDir()`, `TestCacheDir()`
- **文件操作函数**：`CreateTestFile()`, `ReadTestFile()`, `RemoveTestFile()`
- **环境变量函数**：`SetTestEnv()`, `UnsetTestEnv()`

#### 2. Mock 工具 (`mock.go`)

提供 HTTP Mock 和接口 Mock 功能：

- **HTTP Mock**：使用 `net/http/httptest` 创建 Mock HTTP 服务器
- **接口 Mock**：使用 `testify/mock` 创建接口 Mock
- **预设场景**：提供常见的 Mock 场景（GitHub API、Jira API 等）

#### 3. 测试数据生成 (`testdata.go`)

提供测试数据生成和管理：

- **加载 Fixtures**：从 `testdata/fixtures/` 加载测试数据
- **生成测试数据**：使用 Builder 模式生成测试数据
- **数据工厂**：提供常见数据类型的工厂函数

#### 4. 测试环境管理 (`environment.go`)

提供测试环境隔离和管理：

- **环境设置**：设置测试环境变量和临时目录
- **环境清理**：自动清理测试环境
- **环境隔离**：确保测试之间相互独立

### 使用示例

#### 使用测试辅助函数

```go
package config

import (
    "testing"
    "github.com/zevwings/workflow/testutils"
    "github.com/stretchr/testify/assert"
)

func TestLoadConfig(t *testing.T) {
    // 使用测试辅助函数获取测试目录
    testDir := testutils.TestHomeDir(t)
    configPath := filepath.Join(testDir, ".workflow", "config.toml")

    // 测试代码
    config, err := LoadConfig(configPath)
    assert.NoError(t, err)
    assert.NotNil(t, config)
}
```

#### 使用 Mock 工具

```go
package http

import (
    "testing"
    "github.com/zevwings/workflow/testutils"
    "github.com/stretchr/testify/assert"
)

func TestGitHubAPI(t *testing.T) {
    // 创建 Mock HTTP 服务器
    server := testutils.MockGitHubAPI(t)
    defer server.Close()

    // 使用 Mock 服务器进行测试
    client := NewClient(server.URL)
    result, err := client.GetPR("owner", "repo", 123)
    assert.NoError(t, err)
    assert.NotNil(t, result)
}
```

#### 使用测试数据生成

```go
package pr

import (
    "testing"
    "github.com/zevwings/workflow/testutils"
    "github.com/stretchr/testify/assert"
)

func TestParsePR(t *testing.T) {
    // 加载测试数据
    data := testutils.LoadFixture(t, "sample_github_pr.json")

    // 或使用数据工厂生成测试数据
    pr := testutils.NewGitHubPR().
        WithNumber(123).
        WithTitle("Test PR").
        Build()

    // 测试代码
    result, err := ParsePR(data)
    assert.NoError(t, err)
    assert.Equal(t, pr.Number, result.Number)
}
```

### 模块导入路径

```go
// 测试辅助函数
import "github.com/zevwings/workflow/testutils"

// 或使用内部包
import "github.com/zevwings/workflow/internal/testutils"
```

---

## 📦 测试数据管理

### testdata 目录

测试数据应放在 `testdata/` 目录。Go 会自动忽略此目录，不会将其编译到二进制文件中：

```
testdata/
├── fixtures/                    # 测试 Fixtures
│   ├── sample_github_pr.json
│   ├── sample_jira_response.json
│   └── sample_pr_body.md
└── integration/                # 集成测试数据
    └── workflow_scenarios.json
```

### 使用 testdata

```go
package http

import (
    "os"
    "path/filepath"
    "testing"
)

func TestParsePRResponse(t *testing.T) {
    // 读取 testdata 中的文件
    dataPath := filepath.Join("testdata", "fixtures", "sample_github_pr.json")
    data, err := os.ReadFile(dataPath)
    if err != nil {
        t.Fatalf("Failed to read fixture: %v", err)
    }

    // 使用测试数据
    // ...
}
```

### 使用 testutils 加载 Fixtures

```go
package http

import (
    "testing"
    "github.com/zevwings/workflow/testutils"
)

func TestParsePRResponse(t *testing.T) {
    // 使用 testutils 加载 Fixtures
    data := testutils.LoadFixture(t, "sample_github_pr.json")

    // 使用测试数据
    // ...
}
```

---

## 📋 测试组织最佳实践

### 1. 单元测试 vs 集成测试

**测试组织规则**：

- **单元测试（`*_test.go` 文件）**：
  - ✅ **测试包内的所有函数**（包括私有函数）
  - ✅ **主要测试公共 API 和关键私有函数**
  - ✅ 快速执行，最小依赖
  - ✅ 使用表驱动测试提高覆盖率

- **集成测试（`test/` 目录或构建标签）**：
  - ✅ **测试跨包的集成场景**
  - ✅ **测试端到端流程**
  - ✅ 可以使用外部依赖
  - ✅ 使用构建标签 `//go:build integration` 标记

**为什么要区分？**

1. **清晰的测试边界**：单元测试关注单个包，集成测试关注跨包交互
2. **更好的性能**：单元测试快速执行，集成测试可以单独运行
3. **独立编译**：集成测试使用构建标签，可以单独编译和运行
4. **测试覆盖率**：单元测试和集成测试各有侧重，共同提高覆盖率
5. **重构友好**：重构内部实现时，只需更新单元测试

> 📖 **详细指南**：
> - 单元测试：请参考 [单元测试指南](./references/unit-tests.md)
> - 集成测试：请参考 [集成测试指南](./references/integration-tests.md)

### 2. 表驱动测试

使用表驱动测试提高测试覆盖率和可维护性：

```go
func TestParseTicketID(t *testing.T) {
    tests := []struct {
        name     string
        input    string
        expected string
        wantErr  bool
    }{
        {
            name:     "valid ticket ID",
            input:    "PROJ-123",
            expected: "PROJ-123",
            wantErr:  false,
        },
        {
            name:     "invalid format",
            input:    "invalid",
            expected: "",
            wantErr:  true,
        },
        {
            name:     "empty string",
            input:    "",
            expected: "",
            wantErr:  true,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            result, err := ParseTicketID(tt.input)
            if tt.wantErr {
                assert.Error(t, err)
            } else {
                assert.NoError(t, err)
                assert.Equal(t, tt.expected, result)
            }
        })
    }
}
```

### 3. 测试函数命名

- 使用描述性的测试名称
- 使用 `Test` 前缀
- 测试名称应说明测试的内容和预期结果

```go
// ✅ 好的命名
func TestParseTicketID_ValidInput(t *testing.T) {}
func TestParseTicketID_InvalidInput_ReturnsError(t *testing.T) {}

// ❌ 不好的命名
func Test1(t *testing.T) {}
func TestParse(t *testing.T) {}
```

### 4. 测试分组

使用子测试（`t.Run()`）组织相关测试：

```go
func TestHTTPClient(t *testing.T) {
    t.Run("GET request", func(t *testing.T) {
        // 测试 GET 请求
    })

    t.Run("POST request", func(t *testing.T) {
        // 测试 POST 请求
    })

    t.Run("error handling", func(t *testing.T) {
        // 测试错误处理
    })
}
```

---

## 🎯 测试覆盖率

### 覆盖率目标

- **总体覆盖率**：> 80%
- **关键业务逻辑**：> 90%
- **工具函数**：> 70%
- **CLI 命令层**：> 75%

### 覆盖率检查

使用 `go test -cover` 检查覆盖率：

```bash
# 检查覆盖率
go test -cover ./...

# 生成覆盖率报告
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out -o coverage.html

# 使用 Makefile
make test-coverage
```

---

## 📦 测试迁移指南

如果你在其他位置发现了测试代码，请按以下步骤迁移：

### 迁移步骤

1. **识别测试类型**：
   - 单元测试应放在与源码同目录的 `*_test.go` 文件中
   - 集成测试应放在 `test/` 目录或使用构建标签

2. **创建测试文件**：
   - 在源码目录中创建对应的 `*_test.go` 文件
   - 例如：`internal/lib/config/manager.go` → `internal/lib/config/manager_test.go`

3. **迁移测试代码**：
   - 将测试代码复制到新的测试文件
   - 更新 import 语句
   - 添加适当的测试文档注释

4. **验证测试**：
   - 运行 `go test ./...` 确保所有测试通过
   - 检查测试覆盖率没有下降

### 迁移示例

```go
// ✅ 正确：在源码同目录创建测试文件
// internal/lib/config/manager_test.go
package config

import (
    "testing"
    "github.com/stretchr/testify/assert"
)

func TestLoadConfig(t *testing.T) {
    config, err := LoadConfig("testdata/config.toml")
    assert.NoError(t, err)
    assert.NotNil(t, config)
}

func TestParseConfig(t *testing.T) {
    // 可以测试私有函数（如果在同一包内）
    data := []byte(`{"key": "value"}`)
    result, err := parseConfig(data)
    assert.NoError(t, err)
    assert.NotNil(t, result)
}
```

---

## 相关文档

- [测试编写规范](./writing.md) - 测试编写的具体规范
- [测试命令参考](./commands.md) - 常用测试命令
- [测试工具指南](./references/tools.md) - 测试工具使用指南

---

**最后更新**: 2025-01-28

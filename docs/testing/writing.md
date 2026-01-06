# 测试编写规范

> 本文档定义测试编写的具体规范和最佳实践。

---

## 📋 目录

- [测试编写规范](#-测试编写规范)
- [编写测试最佳实践](#-编写测试最佳实践)
  - [1. 测试命名规范](#1-测试命名规范)
  - [2. 测试结构（AAA 模式）](#2-测试结构aaa-模式)
  - [3. 测试独立性](#3-测试独立性)
  - [4. 测试覆盖原则](#4-测试覆盖原则)
  - [5. 测试数据管理](#5-测试数据管理)
  - [6. Mock 使用原则](#6-mock-使用原则)
  - [7. 断言最佳实践](#7-断言最佳实践)
  - [8. 表驱动测试](#8-表驱动测试)
  - [9. 测试基础设施最佳实践](#9-测试基础设施最佳实践)
  - [10. 测试文档](#10-测试文档)
- [被忽略测试文档规范](#-被忽略测试文档规范)

---

## ✅ 测试编写规范

> 📖 **注意**：本文档定义**通用的测试编写规范**，适用于所有测试类型（单元测试、集成测试等）。如需了解单元测试的详细指南，请参考 [单元测试指南](./references/unit-tests.md)。

### 1. 测试结构

每个测试应包含：
- **Arrange**：准备测试数据和环境
- **Act**：执行被测试的功能
- **Assert**：验证结果

```go
func TestParseTicketID(t *testing.T) {
    // Arrange: 准备测试数据
    input := "PROJ-123"
    expected := "PROJ-123"

    // Act: 执行被测试的功能
    result := ParseTicketID(input)

    // Assert: 验证结果
    assert.Equal(t, expected, result)
}
```

### 2. 错误处理测试

为错误情况编写测试：

```go
func TestParseTicketID_InvalidInput(t *testing.T) {
    tests := []struct {
        name  string
        input string
    }{
        {"empty string", ""},
        {"invalid format", "invalid"},
        {"missing project", "-123"},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            result, err := ParseTicketID(tt.input)
            assert.Error(t, err)
            assert.Empty(t, result)
        })
    }
}
```

### 错误处理最佳实践

#### 使用 `t.Fatal` 和 `t.Fatalf` 处理致命错误

```go
// ✅ 推荐：使用 t.Fatal 处理致命错误
func TestExample(t *testing.T) {
    config, err := LoadConfig("testdata/config.toml")
    if err != nil {
        t.Fatalf("Failed to load config: %v", err)
    }
    // 继续测试
}

// ❌ 不推荐：使用 assert 处理致命错误
func TestExample(t *testing.T) {
    config, err := LoadConfig("testdata/config.toml")
    assert.NoError(t, err) // 如果失败，测试会继续执行
    // 可能导致 nil pointer panic
}
```

#### 使用 `require` 包处理必须成功的操作

```go
import (
    "testing"
    "github.com/stretchr/testify/require"
)

// ✅ 推荐：使用 require 处理必须成功的操作
func TestExample(t *testing.T) {
    config, err := LoadConfig("testdata/config.toml")
    require.NoError(t, err) // 失败时立即停止测试
    require.NotNil(t, config)

    // 继续测试
}
```

#### 使用 `assert` 包处理可继续的断言

```go
import (
    "testing"
    "github.com/stretchr/testify/assert"
)

// ✅ 推荐：使用 assert 处理可继续的断言
func TestExample(t *testing.T) {
    result, err := ProcessData("input")
    assert.NoError(t, err) // 失败时记录错误，但继续执行
    assert.Equal(t, "expected", result)
}
```

**选择建议**：
- **`require`**：用于必须成功的操作（配置加载、初始化等）
- **`assert`**：用于可继续的断言（结果验证等）

#### 测试辅助函数中的错误处理

```go
// ✅ 推荐：返回错误，让调用者处理
func LoadFixture(t *testing.T, name string) []byte {
    t.Helper()
    path := filepath.Join("testdata", "fixtures", name)
    data, err := os.ReadFile(path)
    if err != nil {
        t.Fatalf("Failed to load fixture %s: %v", name, err)
    }
    return data
}

// ❌ 不推荐：使用 panic
func LoadFixture(name string) []byte {
    data, err := os.ReadFile(path)
    if err != nil {
        panic(err) // 不应该在测试中使用 panic
    }
    return data
}
```

### 3. 边界条件测试

测试边界条件和极端情况：

```go
func TestParseTicketID_Boundary(t *testing.T) {
    tests := []struct {
        name     string
        input    string
        expected string
        wantErr  bool
    }{
        {
            name:     "minimum length",
            input:    "A-1",
            expected: "A-1",
            wantErr:  false,
        },
        {
            name:     "maximum length",
            input:    "VERY-LONG-PROJECT-NAME-123",
            expected: "VERY-LONG-PROJECT-NAME-123",
            wantErr:  false,
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

---

## ✍️ 编写测试最佳实践

### 1. 测试命名规范

**描述性命名**：
- ✅ 使用描述性的测试名称，说明测试的内容和预期结果
- ✅ 使用 `Test` 前缀
- ✅ 测试名称应包含：被测试的功能、输入条件、预期结果

```go
// ✅ 好的命名
func TestParseTicketID_ValidInput(t *testing.T) {}
func TestParseTicketID_InvalidInput_ReturnsError(t *testing.T) {}
func TestParseTicketID_EmptyString_ReturnsError(t *testing.T) {}

// ❌ 不好的命名
func Test1(t *testing.T) {}
func TestParse(t *testing.T) {}
func TestParseTicketID(t *testing.T) {} // 不够具体
```

**命名模式**：
- `TestFunctionName_Scenario_ExpectedResult`
- `TestFunctionName_InputCondition_Behavior`

### 2. 测试结构（AAA 模式）

**Arrange-Act-Assert 模式**：
```go
func TestExample(t *testing.T) {
    // Arrange: 准备测试数据和环境
    input := "PROJ-123"
    expected := "PROJ-123"

    // Act: 执行被测试的功能
    result := ParseTicketID(input)

    // Assert: 验证结果
    assert.Equal(t, expected, result)
}
```

### 3. 测试独立性

**每个测试应独立**：
- ✅ 每个测试应独立运行，不依赖其他测试
- ✅ 每个测试应使用独立的数据和环境
- ✅ 测试之间不应共享状态

```go
// ✅ 好的做法：每个测试独立
func TestParseTicketID_1(t *testing.T) {
    result := ParseTicketID("PROJ-123")
    assert.Equal(t, "PROJ-123", result)
}

func TestParseTicketID_2(t *testing.T) {
    result := ParseTicketID("PROJ-456")
    assert.Equal(t, "PROJ-456", result)
}

// ❌ 不好的做法：测试之间共享状态
var counter int

func Test1(t *testing.T) {
    counter++
    assert.Equal(t, 1, counter)
}

func Test2(t *testing.T) {
    counter++
    assert.Equal(t, 2, counter) // 依赖 Test1
}
```

### 4. 测试覆盖原则

**测试覆盖重点**：
- ✅ **成功路径**：测试正常流程
- ✅ **错误路径**：测试错误处理和边界条件
- ✅ **边界条件**：测试边界值和极端情况
- ✅ **集成场景**：测试模块间交互

### 5. 测试数据管理

**使用 testdata 目录**：
```go
// ✅ 使用 testdata 目录中的测试数据
func TestParsePRResponse(t *testing.T) {
    dataPath := filepath.Join("testdata", "fixtures", "sample_github_pr.json")
    data, err := os.ReadFile(dataPath)
    require.NoError(t, err)
    // 使用测试数据
}
```

**使用测试数据工厂**：
```go
// ✅ 使用测试数据工厂生成测试数据
import "github.com/your-org/workflow/testutils"

func TestWithFactory(t *testing.T) {
    pr := testutils.NewGitHubPR().
        WithNumber(123).
        WithTitle("Test PR").
        Build()
    // 使用生成的测试数据
}
```

### 6. Mock 使用原则

**何时使用 Mock**：
- ✅ 测试需要调用外部 API（GitHub、Jira 等）
- ✅ 测试需要模拟网络请求和响应
- ✅ 测试需要避免依赖外部服务
- ✅ 测试需要模拟错误情况

**Mock 使用规范**：
```go
// ✅ 使用 httptest 创建 Mock HTTP 服务器
import (
    "net/http/httptest"
    "testing"
)

func TestAPICall(t *testing.T) {
    // 创建 Mock 服务器
    server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.WriteHeader(http.StatusOK)
        w.Write([]byte(`{"result": "success"}`))
    }))
    defer server.Close()

    // 使用 Mock 服务器进行测试
    client := NewClient(server.URL)
    result, err := client.CallAPI()
    assert.NoError(t, err)
    assert.Equal(t, "success", result)
}
```

### 7. 断言最佳实践

**使用清晰的断言**：
```go
// ✅ 使用描述性的断言消息
assert.Equal(t, expected, result, "Failed to parse ticket ID: %s", input)

// ✅ 使用专门的断言工具
import "github.com/stretchr/testify/assert"

assert.Equal(t, expected, result)
assert.NoError(t, err)
assert.NotNil(t, obj)

// ❌ 避免模糊的断言
if result == nil {
    t.Fatal("result is nil") // 不够清晰
}
```

**使用 `t.Helper()` 标记辅助函数**：
```go
func loadFixture(t *testing.T, name string) []byte {
    t.Helper() // 标记为辅助函数，错误信息会指向调用者
    path := filepath.Join("testdata", "fixtures", name)
    data, err := os.ReadFile(path)
    if err != nil {
        t.Fatalf("Failed to load fixture %s: %v", name, err)
    }
    return data
}
```

### 8. 表驱动测试

表驱动测试允许你使用不同的输入值运行同一个测试函数，从而减少重复代码并提高测试覆盖率。

#### 何时使用表驱动测试

✅ **适合使用表驱动测试的场景**：
- 多个相似测试函数（测试相同的功能，只是输入不同）
- 表格驱动测试（需要测试多种输入组合）
- 边界值测试（测试多个边界值和正常值）
- 枚举值测试（测试枚举的所有变体）

❌ **不适合使用表驱动测试的场景**：
- 测试不同的错误场景（不同的错误需要不同的断言和验证逻辑）
- 需要不同设置的测试（每个测试需要不同的环境设置或fixture配置）
- 测试执行顺序重要（测试之间有依赖关系）

#### 基本用法

```go
func TestParseTicketID(t *testing.T) {
    tests := []struct {
        name     string
        input    string
        expected string
        wantErr  bool
    }{
        {
            name:     "valid input",
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

#### 表驱动测试最佳实践

**1. 测试函数命名**：
```go
// ✅ 好的命名
func TestParseTicketID_TableDriven(t *testing.T) {}

// ❌ 不好的命名
func TestParse(t *testing.T) {}
```

**2. 文档注释**：
```go
// TestParseTicketID tests the ParseTicketID function with various inputs.
//
// Test cases:
//   - Valid ticket IDs (PROJ-123)
//   - Invalid formats (invalid, empty)
//   - Boundary conditions (minimum/maximum length)
func TestParseTicketID(t *testing.T) {
    // ...
}
```

**3. Case 注释**：
```go
tests := []struct {
    name     string
    input    string
    expected string
    wantErr  bool
}{
    {
        name:     "valid input",
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
```

**4. 保持测试独立**：
```go
// ✅ 好的做法：每个 case 独立
for _, tt := range tests {
    t.Run(tt.name, func(t *testing.T) {
        // 每个 case 独立执行
    })
}

// ❌ 不好的做法：case 之间有依赖
var sharedState int
for _, tt := range tests {
    t.Run(tt.name, func(t *testing.T) {
        sharedState++ // 依赖其他 case
    })
}
```

#### 常见模式

**验证器测试**：
```go
func TestValidator(t *testing.T) {
    tests := []struct {
        name      string
        input     string
        wantValid bool
    }{
        {"valid", "valid", true},
        {"invalid", "invalid", false},
        {"empty", "", false},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            validator := NewValidator()
            result := validator.Validate(tt.input)
            assert.Equal(t, tt.wantValid, result.IsValid())
        })
    }
}
```

### 9. 测试基础设施最佳实践

测试基础设施提供了统一的测试环境隔离和路径获取机制，确保测试的可靠性和跨平台兼容性。

#### 9.1 使用测试隔离

**✅ 好的做法**：
```go
import (
    "testing"
    "github.com/your-org/workflow/testutils"
)

func TestExample(t *testing.T) {
    // 使用 t.TempDir() 创建临时目录
    tempDir := t.TempDir()

    // 使用 t.Setenv() 设置环境变量（自动恢复）
    t.Setenv("HOME", tempDir)

    // 测试代码在完全隔离的环境中运行
}
```

**❌ 不好的做法**：
```go
func TestExample(t *testing.T) {
    // 直接使用系统路径，可能污染系统
    configPath := filepath.Join(os.Getenv("HOME"), ".workflow", "config.toml")
    // ...
}
```

**优势**：
- ✅ 测试之间不会相互影响
- ✅ 测试不会污染真实系统路径
- ✅ 测试可以在并行环境中安全运行
- ✅ 测试结果可重复

#### 9.2 使用统一的路径获取函数

**✅ 好的做法**：
```go
import (
    "testing"
    "github.com/your-org/workflow/testutils"
)

func TestExample(t *testing.T) {
    // 使用测试辅助函数获取测试目录
    homeDir := testutils.TestHomeDir(t)
    configDir := testutils.TestConfigDir(t)

    // 使用测试目录进行测试
}
```

**❌ 不好的做法**：
```go
func TestExample(t *testing.T) {
    // 直接使用 os.UserHomeDir()，不支持测试隔离
    homeDir, err := os.UserHomeDir()
    if err != nil {
        t.Fatal(err)
    }
    // ...
}
```

**可用的路径获取函数**：
- `TestHomeDir(t)` - 获取主目录（测试环境感知）
- `TestConfigDir(t)` - 获取配置目录（测试环境感知）
- `TestDataDir(t)` - 获取数据目录（测试环境感知）
- `TestCacheDir(t)` - 获取缓存目录（测试环境感知）

**注意事项**：
- 这些函数优先使用环境变量（支持测试隔离），然后回退到系统路径
- 与源代码中的路径获取行为一致
- 临时目录应使用 `t.TempDir()`
- 当前目录应使用 `os.Getwd()`

详细说明请参考 [测试辅助工具指南 - 路径获取函数](./references/helpers.md#3-路径获取函数)。

#### 9.3 清理测试数据

测试环境使用 `t.Cleanup()` 自动清理，无需手动清理：

```go
func TestExample(t *testing.T) {
    tempDir := t.TempDir() // 自动清理

    // 测试代码

    // 不需要手动清理，测试结束后自动清理
}
```

**自动清理机制**：
- `t.TempDir()` 创建的临时目录在测试结束后自动清理
- `t.Setenv()` 设置的环境变量在测试结束后自动恢复
- `t.Cleanup()` 注册的清理函数在测试结束后自动执行

#### 9.4 平台特定测试

使用构建标签标记平台特定测试：

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
    // Unix 特定测试
}
```

**平台差异处理**：
- 使用统一的路径获取函数（`TestHomeDir()` 等）自动处理平台差异
- 使用构建标签标记平台特定代码
- 在 CI/CD 中运行跨平台测试
- 参考平台差异分析文档了解平台差异

#### 9.5 测试环境选择

根据测试需求选择合适的测试环境：

**基础隔离（t.TempDir）**：
```go
func TestBasicIsolation(t *testing.T) {
    tempDir := t.TempDir()
    t.Setenv("HOME", tempDir)

    // 使用 tempDir 进行测试
}
```

**CLI 测试环境（testutils）**：
```go
import "github.com/your-org/workflow/testutils"

func TestCLICommand(t *testing.T) {
    env := testutils.SetupTestEnv(t)
    defer env.Cleanup()

    // 使用 env 进行测试
}
```

**选择建议**：
- **需要 Git 仓库操作** → 使用 `testutils.GitTestEnv`
- **只需要基础隔离** → 使用 `t.TempDir()` 和 `t.Setenv()`
- **需要 Mock 服务器** → 使用 `httptest.NewServer()`

详细说明请参考 [测试环境工具指南](./references/environments.md)。

### 10. 测试文档

**为复杂测试添加注释**：
```go
// TestComplexScenario tests the complex scenario where a user inputs an invalid ticket ID.
// The system should return an error and log the error message.
func TestComplexScenario(t *testing.T) {
    input := "INVALID"
    result, err := ParseTicketID(input)

    assert.Error(t, err)
    assert.Empty(t, result)
    // 验证错误日志已记录
}
```

---

## 🚫 被忽略测试文档规范

对于使用 `t.Skip()` 跳过的测试，必须添加完整的文档注释。

### 统一文档格式

所有被跳过的测试都应该包含以下5个部分的文档注释：

```go
// TestFunctionName tests the function with a specific scenario.
//
// ## 测试目的
// 验证/测试...（说明测试验证什么功能）
//
// ## 为什么被跳过
// - **主要原因**: ...
// - **次要原因**: ...
// - **使用场景**: ...
//
// ## 如何手动运行
// ```bash
// go test -run TestFunctionName ./...
// ```
// （如适用）额外的运行说明或交互步骤
//
// ## 测试场景
// 1. ...
// 2. ...
// 3. ...
//
// ## 预期行为
// - ...
// - ...
func TestFunctionName(t *testing.T) {
    t.Skip("简短原因")
    // 测试代码
}
```

### 常见跳过原因

**1. 用户交互测试**：
- **需要用户交互**: 测试需要用户在终端中进行交互操作
- **CI环境不支持**: 自动化CI环境无法提供交互式输入

**2. 网络请求测试**：
- **需要网络连接**: 测试需要实际的网络连接到外部API
- **需要API密钥**: 需要有效的API密钥或认证凭据
- **CI成本考虑**: 避免在CI中产生API调用费用

**3. 时间相关测试**：
- **涉及真实时间延迟**: 测试需要等待实际的时间流逝
- **测试运行时间长**: 完整测试需要较长时间
- **CI时间限制**: 避免在CI中占用过多时间

**4. 修改系统配置的测试**：
- **修改系统文件**: 测试会修改用户的配置文件
- **安全风险**: 避免在CI或开发环境中意外修改配置

详细的被忽略测试规范请参考 [被忽略测试规范](./references/ignored-tests.md)。

---

## 相关文档

- [测试组织规范](./organization.md) - 测试组织结构和命名约定
- [单元测试指南](./references/unit-tests.md) - 单元测试的详细编写规范和组织方式
- [集成测试指南](./references/integration-tests.md) - 集成测试的环境配置和最佳实践
- [测试命令参考](./commands.md) - 常用测试命令
- [测试工具指南](./references/tools.md) - 测试工具使用
- [被忽略测试规范](./references/ignored-tests.md) - 被忽略测试的完整规范

---

**最后更新**: 2025-01-28

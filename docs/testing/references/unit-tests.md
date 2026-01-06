# 单元测试指南

> 本文档介绍单元测试的编写规范、组织方式和最佳实践。

---

## 📋 目录

- [概述](#-概述)
- [单元测试组织](#-单元测试组织)
- [编写规范](#-编写规范)
- [测试模式](#-测试模式)
- [最佳实践](#-最佳实践)
- [常见场景](#-常见场景)

---

## 📋 概述

单元测试验证单个函数或模块的行为，确保代码在隔离环境中按预期工作。

### 单元测试特点

- **单包测试**：测试包内的所有函数（包括私有和公开）
- **快速执行**：不依赖外部环境，执行速度快（< 100ms）
- **最小依赖**：使用 Mock 或测试替身隔离外部依赖
- **高覆盖率**：通过表驱动测试提高覆盖率
- **独立运行**：每个测试相互独立，不共享状态

### 单元测试 vs 集成测试

| 特性 | 单元测试 | 集成测试 |
|------|---------|---------|
| **测试范围** | 单个包内的函数 | 跨包的交互和端到端流程 |
| **执行速度** | 快速（< 100ms） | 较慢（< 1s） |
| **依赖** | 最小依赖，使用 Mock | 可以使用真实依赖 |
| **文件位置** | 与源码同目录的 `*_test.go` | `test/` 目录或构建标签 |
| **构建标签** | 不需要 | 使用 `//go:build integration` |

---

## 单元测试组织

### 测试文件位置

单元测试文件应与源码文件在同一目录，使用 `*_test.go` 后缀：

```
internal/
├── lib/
│   ├── config/
│   │   ├── manager.go           # 源代码
│   │   └── manager_test.go      # 单元测试
│   └── http/
│       ├── client.go
│       └── client_test.go
```

### 测试包名

单元测试使用与源码相同的包名，可以访问包内的私有函数：

```go
// internal/lib/config/manager.go
package config

func parseConfig(data []byte) (*Config, error) {
    // 私有函数实现
}

// internal/lib/config/manager_test.go
package config  // ✅ 相同包名，可以测试私有函数

import (
    "testing"
    "github.com/stretchr/testify/assert"
)

func TestParseConfig(t *testing.T) {
    // ✅ 可以测试私有函数 parseConfig
    data := []byte(`{"key": "value"}`)
    result, err := parseConfig(data)
    assert.NoError(t, err)
    assert.NotNil(t, result)
}
```

### 运行单元测试

```bash
# 运行所有单元测试
go test ./...

# 运行特定包的单元测试
go test ./internal/lib/config

# 运行特定测试函数
go test -run TestParseConfig ./internal/lib/config

# 运行测试并显示覆盖率
go test -cover ./internal/lib/config
```

---

## 编写规范

> 📖 **注意**：本节介绍单元测试特有的编写规范。通用的测试编写规范（AAA 模式、命名规范、表驱动测试等）请参考 [测试编写规范](../writing.md)。

### 1. 测试结构（AAA 模式）

每个测试应遵循 **Arrange-Act-Assert** 模式：

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

> 📖 更多关于 AAA 模式的说明，请参考 [测试编写规范 - 测试结构](../writing.md#1-测试结构)。

### 2. 测试命名规范

**命名模式**：`TestFunctionName_Scenario_ExpectedResult`

```go
// ✅ 好的命名
func TestParseTicketID_ValidInput_ReturnsTicketID(t *testing.T) {}
func TestParseTicketID_InvalidInput_ReturnsError(t *testing.T) {}
func TestParseTicketID_EmptyString_ReturnsError(t *testing.T) {}

// ❌ 不好的命名
func Test1(t *testing.T) {}
func TestParse(t *testing.T) {}
func TestParseTicketID(t *testing.T) {} // 不够具体
```

> 📖 更多关于测试命名规范的说明，请参考 [测试编写规范 - 测试命名规范](../writing.md#1-测试命名规范)。

### 3. 测试私有函数

单元测试可以测试包内的私有函数（小写开头的函数）：

```go
// internal/lib/config/manager.go
package config

func parseConfig(data []byte) (*Config, error) {
    // 私有函数实现
}

// internal/lib/config/manager_test.go
package config

func TestParseConfig(t *testing.T) {
    // ✅ 可以测试私有函数 parseConfig
    data := []byte(`{"key": "value"}`)
    result, err := parseConfig(data)
    assert.NoError(t, err)
    assert.NotNil(t, result)
}
```

**何时测试私有函数**：
- ✅ 复杂的私有函数逻辑
- ✅ 关键的业务逻辑（即使私有）
- ✅ 需要高覆盖率的函数
- ❌ 简单的 getter/setter 函数
- ❌ 通过公共 API 可以充分测试的函数

### 4. 测试独立性

每个测试应独立运行，不依赖其他测试：

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

---

## 测试模式

### 1. 表驱动测试

表驱动测试是 Go 推荐的测试模式，提高测试覆盖率和可维护性：

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
        {
            name:     "boundary: minimum length",
            input:    "A-1",
            expected: "A-1",
            wantErr:  false,
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

> 📖 更多关于表驱动测试的详细说明和最佳实践，请参考 [测试编写规范 - 表驱动测试](../writing.md#8-表驱动测试)。

**表驱动测试优势**：
- ✅ 减少重复代码
- ✅ 易于添加新的测试用例
- ✅ 清晰的测试场景展示
- ✅ 统一的测试结构

### 2. 错误处理测试

为错误情况编写测试：

```go
func TestParseTicketID_ErrorCases(t *testing.T) {
    tests := []struct {
        name  string
        input string
    }{
        {"empty string", ""},
        {"invalid format", "invalid"},
        {"missing project", "-123"},
        {"missing number", "PROJ-"},
        {"special characters", "PROJ-123@#$"},
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

### 3. 边界条件测试

测试边界值和极端情况：

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
            name:     "exceeds maximum length",
            input:    strings.Repeat("A", 100) + "-123",
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

### 4. Mock 测试

使用 Mock 隔离外部依赖：

```go
func TestHTTPClient_Get(t *testing.T) {
    // 创建 Mock HTTP 服务器
    server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        assert.Equal(t, "/api/v1/data", r.URL.Path)
        w.WriteHeader(http.StatusOK)
        w.Write([]byte(`{"result": "success"}`))
    }))
    defer server.Close()

    // 使用 Mock 服务器进行测试
    client := NewClient(server.URL)
    result, err := client.Get("/api/v1/data")

    assert.NoError(t, err)
    assert.Equal(t, "success", result.Result)
}
```

---

## 最佳实践

### 1. 使用测试辅助函数

使用 `t.Helper()` 标记辅助函数：

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

func TestParsePR(t *testing.T) {
    data := loadFixture(t, "sample_github_pr.json")
    // 使用测试数据
}
```

### 2. 使用 require 和 assert

**`require`**：用于必须成功的操作（失败时立即停止测试）
**`assert`**：用于可继续的断言（失败时记录错误，但继续执行）

```go
import (
    "testing"
    "github.com/stretchr/testify/require"
    "github.com/stretchr/testify/assert"
)

func TestExample(t *testing.T) {
    // ✅ 使用 require 处理必须成功的操作
    config, err := LoadConfig("testdata/config.toml")
    require.NoError(t, err) // 失败时立即停止
    require.NotNil(t, config)

    // ✅ 使用 assert 处理可继续的断言
    result, err := ProcessData("input")
    assert.NoError(t, err) // 失败时记录错误，但继续执行
    assert.Equal(t, "expected", result)
}
```

### 3. 测试文档注释

为复杂测试添加文档注释：

```go
// TestComplexScenario tests the complex scenario where a user inputs an invalid ticket ID.
// The system should return an error and log the error message.
//
// Test cases:
//   - Invalid ticket ID format
//   - Empty string input
//   - Special characters in input
func TestComplexScenario(t *testing.T) {
    // 测试代码
}
```

### 4. 使用测试数据

使用 `testdata/` 目录存放测试数据：

```go
func TestParsePRResponse(t *testing.T) {
    // 使用 testdata 目录中的测试数据
    dataPath := filepath.Join("testdata", "fixtures", "sample_github_pr.json")
    data, err := os.ReadFile(dataPath)
    require.NoError(t, err)

    // 使用测试数据
    pr, err := ParsePRResponse(data)
    assert.NoError(t, err)
    assert.NotNil(t, pr)
}
```

### 5. 并行测试

对于独立的测试，可以使用 `t.Parallel()` 并行执行：

```go
func TestParseTicketID_Parallel(t *testing.T) {
    t.Parallel()

    tests := []struct {
        name     string
        input    string
        expected string
    }{
        {"case1", "PROJ-123", "PROJ-123"},
        {"case2", "PROJ-456", "PROJ-456"},
    }

    for _, tt := range tests {
        tt := tt // 捕获循环变量
        t.Run(tt.name, func(t *testing.T) {
            t.Parallel()
            result := ParseTicketID(tt.input)
            assert.Equal(t, tt.expected, result)
        })
    }
}
```

---

## 常见场景

### 1. 测试配置加载

```go
func TestLoadConfig(t *testing.T) {
    // 创建临时配置文件
    tempDir := t.TempDir()
    configPath := filepath.Join(tempDir, "config.toml")
    configContent := `[github]
token = "test-token"
`
    err := os.WriteFile(configPath, []byte(configContent), 0644)
    require.NoError(t, err)

    // 测试配置加载
    config, err := LoadConfig(configPath)
    require.NoError(t, err)
    assert.Equal(t, "test-token", config.GitHub.Token)
}
```

### 2. 测试字符串处理

```go
func TestParseTicketID(t *testing.T) {
    tests := []struct {
        name     string
        input    string
        expected string
        wantErr  bool
    }{
        {"valid", "PROJ-123", "PROJ-123", false},
        {"invalid", "invalid", "", true},
        {"empty", "", "", true},
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

### 3. 测试 HTTP 客户端

```go
func TestHTTPClient_Get(t *testing.T) {
    // 创建 Mock 服务器
    server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.WriteHeader(http.StatusOK)
        w.Write([]byte(`{"id": 123}`))
    }))
    defer server.Close()

    // 测试客户端
    client := NewClient(server.URL)
    result, err := client.Get("/api/data")

    require.NoError(t, err)
    assert.Equal(t, 123, result.ID)
}
```

### 4. 测试错误处理

```go
func TestProcessData_ErrorHandling(t *testing.T) {
    tests := []struct {
        name        string
        input       string
        expectedErr string
    }{
        {
            name:        "invalid input",
            input:       "invalid",
            expectedErr: "invalid input format",
        },
        {
            name:        "empty input",
            input:       "",
            expectedErr: "input cannot be empty",
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            _, err := ProcessData(tt.input)
            assert.Error(t, err)
            assert.Contains(t, err.Error(), tt.expectedErr)
        })
    }
}
```

### 5. 测试数据结构转换

```go
func TestConvertPR(t *testing.T) {
    // 使用测试数据
    data := loadFixture(t, "sample_github_pr.json")

    var githubPR GitHubPR
    err := json.Unmarshal(data, &githubPR)
    require.NoError(t, err)

    // 测试转换
    pr := ConvertPR(githubPR)
    assert.Equal(t, githubPR.Number, pr.Number)
    assert.Equal(t, githubPR.Title, pr.Title)
}
```

---

## 单元测试检查清单

### 开发时

- [ ] 为新功能添加单元测试
- [ ] 测试文件放在与源码同目录
- [ ] 使用表驱动测试提高覆盖率
- [ ] 测试成功路径和错误路径
- [ ] 测试边界条件

### 代码审查时

- [ ] 检查单元测试覆盖主要功能
- [ ] 确保测试使用 AAA 模式
- [ ] 检查测试命名是否清晰
- [ ] 确保测试之间相互独立
- [ ] 检查是否使用了适当的 Mock

### 发布前

- [ ] 运行完整的单元测试套件
- [ ] 确保所有单元测试通过
- [ ] 检查测试覆盖率（目标 > 80%）
- [ ] 检查测试执行时间（单元测试 < 100ms）

---

## 相关文档

- [测试组织规范](../organization.md) - 测试组织结构和命名约定
- [测试编写规范](../writing.md) - 通用的测试编写规范（AAA 模式、命名规范、表驱动测试等）
- [集成测试指南](./integration-tests.md) - 集成测试指南
- [Mock测试指南](./mock-server.md) - Mock 测试详细说明
- [测试工具指南](./tools.md) - 测试工具使用指南

---

**最后更新**: 2025-01-28


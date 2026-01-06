# 集成测试指南

> 本文档介绍集成测试的环境配置和最佳实践。

---

## 📋 目录

- [概述](#-概述)
- [集成测试组织](#-集成测试组织)
- [环境配置](#-环境配置)
- [数据隔离](#-数据隔离)
- [清理机制](#-清理机制)
- [最佳实践](#-最佳实践)

---

## 📋 概述

集成测试验证多个模块之间的交互和端到端流程，确保系统作为一个整体正常工作。

### 集成测试特点

- **跨包测试**：测试多个包的交互
- **真实环境**：使用真实的依赖（如 Git、文件系统）
- **端到端流程**：测试完整的工作流程
- **独立运行**：使用构建标签 `//go:build integration` 标记

---

## 集成测试组织

### 测试文件位置

集成测试可以放在以下位置：

1. **单独的测试文件**：使用构建标签标记
2. **`test/integration/` 目录**：专门的集成测试目录

### 使用构建标签

```go
//go:build integration

package config

import (
    "testing"
)

func TestIntegration(t *testing.T) {
    // 集成测试代码
}
```

### 运行集成测试

```bash
# 运行集成测试
go test -tags=integration ./...

# 运行特定包的集成测试
go test -tags=integration ./test/integration
```

---

## 环境配置

### 基本环境设置

```go
//go:build integration

package integration

import (
    "os"
    "path/filepath"
    "testing"
    "github.com/stretchr/testify/require"
)

func TestWorkflowIntegration(t *testing.T) {
    // 创建临时目录
    tempDir := t.TempDir()

    // 设置环境变量
    t.Setenv("HOME", tempDir)
    t.Setenv("CONFIG_DIR", filepath.Join(tempDir, ".config"))

    // 初始化测试环境
    setupTestEnvironment(t, tempDir)

    // 执行集成测试
    // ...
}
```

### Git 环境设置

```go
func TestGitIntegration(t *testing.T) {
    // 创建临时目录
    tempDir := t.TempDir()

    // 初始化 Git 仓库
    err := exec.Command("git", "init", tempDir).Run()
    require.NoError(t, err)

    // 配置 Git 用户
    exec.Command("git", "config", "user.name", "Test User").Run()
    exec.Command("git", "config", "user.email", "test@example.com").Run()

    // 执行 Git 集成测试
    // ...
}
```

### 外部服务 Mock

```go
func TestAPIIntegration(t *testing.T) {
    // 创建 Mock HTTP 服务器
    server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // Mock 响应
        w.WriteHeader(http.StatusOK)
        w.Write([]byte(`{"id": 123}`))
    }))
    defer server.Close()

    // 设置 API URL
    t.Setenv("GITHUB_API_URL", server.URL)

    // 执行 API 集成测试
    // ...
}
```

---

## 数据隔离

### 每个测试独立数据

```go
func TestIntegration1(t *testing.T) {
    // 每个测试使用独立的临时目录
    tempDir := t.TempDir()
    t.Setenv("HOME", tempDir)

    // 测试代码
}

func TestIntegration2(t *testing.T) {
    // 独立的测试环境
    tempDir := t.TempDir()
    t.Setenv("HOME", tempDir)

    // 测试代码（不依赖 TestIntegration1）
}
```

### 测试数据目录

```go
func TestWithTestData(t *testing.T) {
    // 使用 testdata 目录
    dataPath := filepath.Join("testdata", "integration", "scenario1.json")
    data, err := os.ReadFile(dataPath)
    require.NoError(t, err)

    // 使用测试数据
    // ...
}
```

---

## 清理机制

### 自动清理

Go 测试框架提供自动清理机制：

```go
func TestIntegration(t *testing.T) {
    // 使用 t.TempDir() 创建的目录会自动清理
    tempDir := t.TempDir()

    // 使用 t.Setenv() 设置的环境变量会自动恢复
    t.Setenv("HOME", tempDir)

    // 使用 t.Cleanup() 注册的清理函数会自动执行
    t.Cleanup(func() {
        // 额外的清理工作
    })

    // 测试代码
}
```

### 手动清理

```go
func TestIntegration(t *testing.T) {
    // 创建需要手动清理的资源
    resource := setupResource(t)
    defer resource.Cleanup()

    // 测试代码
}
```

---

## 最佳实践

### 1. 使用构建标签

```go
//go:build integration

package integration

// 集成测试代码
```

### 2. 测试完整流程

```go
func TestWorkflowIntegration(t *testing.T) {
    // 1. 设置环境
    env := setupTestEnvironment(t)

    // 2. 执行操作
    result, err := ExecuteWorkflow(env)
    require.NoError(t, err)

    // 3. 验证结果
    assert.NotNil(t, result)
    assert.Equal(t, expected, result)
}
```

### 3. 测试错误处理

```go
func TestWorkflowIntegration_ErrorHandling(t *testing.T) {
    // 设置错误场景
    env := setupErrorScenario(t)

    // 执行操作
    _, err := ExecuteWorkflow(env)

    // 验证错误处理
    assert.Error(t, err)
    assert.Contains(t, err.Error(), "expected error")
}
```

### 4. 使用测试辅助函数

```go
func setupTestEnvironment(t *testing.T) *TestEnv {
    t.Helper()

    tempDir := t.TempDir()
    t.Setenv("HOME", tempDir)

    return &TestEnv{
        HomeDir: tempDir,
    }
}
```

### 5. 并行测试

```go
func TestIntegration_Parallel(t *testing.T) {
    t.Parallel()

    // 测试代码（确保相互独立）
}
```

---

## 集成测试检查清单

### 开发时

- [ ] 为新功能添加集成测试
- [ ] 确保集成测试使用构建标签
- [ ] 测试完整的工作流程

### 代码审查时

- [ ] 检查集成测试覆盖主要流程
- [ ] 确保测试环境隔离
- [ ] 检查清理机制

### 发布前

- [ ] 运行完整的集成测试套件
- [ ] 确保所有集成测试通过
- [ ] 检查测试执行时间

---

## 相关文档

- [测试组织规范](../organization.md) - 测试组织结构
- [测试编写规范](../writing.md) - 测试编写规范
- [测试环境工具指南](./environments.md) - 测试环境工具详细使用方法

---

**最后更新**: 2025-01-28


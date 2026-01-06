# 覆盖率测试指南

> 本文档介绍测试覆盖率的检查和提升方法。

---

## 📋 目录

- [覆盖率工具](#-覆盖率工具)
- [生成覆盖率报告](#-生成覆盖率报告)
- [覆盖率目标](#-覆盖率目标)
- [覆盖率提升技巧](#-覆盖率提升技巧)

---

## 覆盖率工具

Go 标准库提供了内置的覆盖率工具，无需额外安装。

### 基本使用

```bash
# 显示覆盖率
go test -cover ./...

# 显示每个包的覆盖率
go test -cover ./internal/lib/config

# 显示每个函数的覆盖率
go test -coverprofile=coverage.out ./...
go tool cover -func=coverage.out
```

---

## 生成覆盖率报告

### HTML 格式报告

```bash
# 生成覆盖率文件
go test -coverprofile=coverage.out ./...

# 生成 HTML 报告
go tool cover -html=coverage.out -o coverage.html

# 打开报告（macOS）
open coverage.html

# 打开报告（Linux）
xdg-open coverage.html

# 打开报告（Windows）
start coverage.html
```

### 使用 Makefile

```bash
# 生成覆盖率报告
make test-coverage

# 查看覆盖率报告
open coverage.html
```

### CI 环境覆盖率

```bash
# 生成覆盖率文件（CI 环境）
go test -coverprofile=coverage.out -covermode=atomic ./...

# 上传到覆盖率服务（如 Codecov）
# codecov -f coverage.out
```

---

## 覆盖率目标

- **总体覆盖率**：> 80%
- **关键业务逻辑**：> 90%
- **工具函数**：> 70%
- **CLI 命令层**：> 75%

### 检查覆盖率

```bash
# 检查总体覆盖率
go test -cover ./... | grep coverage

# 检查特定包的覆盖率
go test -cover ./internal/lib/config

# 生成详细覆盖率报告
go test -coverprofile=coverage.out ./...
go tool cover -func=coverage.out | grep total
```

---

## 覆盖率提升技巧

### 1. 识别低覆盖率区域

```bash
# 生成覆盖率报告
go test -coverprofile=coverage.out ./...

# 查看未覆盖的代码
go tool cover -html=coverage.out

# 在浏览器中查看，红色表示未覆盖的代码
```

### 2. 补充边界测试

为边界条件添加测试：

```go
func TestParseTicketID_Boundary(t *testing.T) {
    tests := []struct {
        name     string
        input    string
        expected string
        wantErr  bool
    }{
        {"minimum length", "A-1", "A-1", false},
        {"maximum length", "VERY-LONG-PROJECT-NAME-123", "VERY-LONG-PROJECT-NAME-123", false},
        {"empty string", "", "", true},
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

### 3. 添加错误处理测试

为错误情况添加测试：

```go
func TestLoadConfig_ErrorCases(t *testing.T) {
    tests := []struct {
        name    string
        path    string
        wantErr bool
    }{
        {"file not found", "/nonexistent/config.toml", true},
        {"invalid format", "testdata/invalid.toml", true},
        {"permission denied", "/root/config.toml", true},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            _, err := LoadConfig(tt.path)
            if tt.wantErr {
                assert.Error(t, err)
            } else {
                assert.NoError(t, err)
            }
        })
    }
}
```

### 4. 使用表驱动测试

使用表驱动测试提高覆盖率：

```go
func TestParseTicketID_TableDriven(t *testing.T) {
    tests := []struct {
        name     string
        input    string
        expected string
        wantErr  bool
    }{
        {"valid", "PROJ-123", "PROJ-123", false},
        {"invalid", "invalid", "", true},
        {"empty", "", "", true},
        // 添加更多测试用例
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

### 5. 测试所有分支

确保测试覆盖所有代码分支：

```go
func TestProcessData_AllBranches(t *testing.T) {
    // 测试成功路径
    result, err := ProcessData("valid")
    assert.NoError(t, err)
    assert.NotNil(t, result)

    // 测试错误路径
    _, err = ProcessData("invalid")
    assert.Error(t, err)

    // 测试边界条件
    result, err = ProcessData("")
    assert.Error(t, err)
}
```

### 6. 使用覆盖率工具分析

```bash
# 生成覆盖率报告
go test -coverprofile=coverage.out ./...

# 查看未覆盖的代码行
go tool cover -html=coverage.out

# 查看覆盖率统计
go tool cover -func=coverage.out | grep -v "100.0%"
```

---

## 覆盖率检查清单

### 开发时

- [ ] 运行 `go test -cover ./...` 检查覆盖率
- [ ] 查看覆盖率报告，识别低覆盖率区域
- [ ] 为新功能添加测试，确保覆盖率不下降

### 代码审查时

- [ ] 检查新代码的测试覆盖率
- [ ] 确保关键业务逻辑有充分的测试
- [ ] 确保错误处理路径有测试覆盖

### 发布前

- [ ] 运行完整的覆盖率检查
- [ ] 确保总体覆盖率 > 80%
- [ ] 确保关键业务逻辑覆盖率 > 90%

---

## 相关文档

- [测试组织规范](../organization.md) - 测试组织结构
- [测试编写规范](../writing.md) - 测试编写规范
- [测试命令参考](../commands.md) - 常用测试命令

---

**最后更新**: 2025-01-28

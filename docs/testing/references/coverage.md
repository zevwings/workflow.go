# 覆盖率测试指南

> 本文档介绍测试覆盖率的检查和提升方法。

---

## 📋 目录

- [覆盖率工具](#-覆盖率工具)
- [生成覆盖率报告](#-生成覆盖率报告)
- [美观的覆盖率 UI 工具](#-美观的覆盖率-ui-工具)
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

## 美观的覆盖率 UI 工具

Go 标准库的 `go tool cover` 生成的 HTML 报告功能完整，但界面较为简单。如果你需要类似 Rust `cargo-tarpaulin` 那样更美观、交互性更强的 UI，可以使用以下工具：

### 1. gocovsh（终端交互式 UI，推荐，类似 cargo-tarpaulin）

`gocovsh` 提供终端交互式 UI，类似 `cargo-tarpaulin` 的终端体验，无需浏览器即可查看覆盖率。

#### 安装

```bash
go install github.com/orlangure/gocovsh@latest
```

#### 使用

```bash
# 生成覆盖率文件
go test -tags=test -coverprofile=coverage/coverage.out ./...

# 启动交互式终端 UI
gocovsh coverage/coverage.out
```

#### 使用 Makefile

```bash
# 使用 gocovsh 交互式查看覆盖率
make test-coverage-interactive
```

**特点**：
- 🖥️ 终端交互式界面（最接近 cargo-tarpaulin 的体验）
- ⌨️ 键盘导航（方向键、搜索等）
- 📁 按包浏览覆盖率
- 🔎 实时搜索功能
- 🎨 彩色高亮显示
- 🎯 支持多种主题（mocha, latte, frappe, macchiato）

### 2. gocov + gocov-html（HTML 报告）

`gocov` 和 `gocov-html` 组合使用，提供另一种 HTML 报告格式。

#### 安装

```bash
go install github.com/axw/gocov/gocov@latest
go install github.com/matm/gocov-html@latest
```

#### 使用

```bash
# 生成覆盖率 JSON 报告
gocov test -tags=test ./... > coverage/coverage.json

# 转换为 HTML
gocov-html coverage/coverage.json > coverage/coverage-ui.html

# 打开报告
open coverage/coverage-ui.html  # macOS
```

#### 使用 Makefile

```bash
# 生成美观的覆盖率报告
make test-coverage-ui

# 打开报告
make open-coverage-ui
```

**特点**：
- 🎨 比标准 HTML 更美观的界面
- 📊 详细的覆盖率统计
- 🔍 按包、文件查看覆盖率
- 📄 JSON 格式便于集成其他工具

### 3. go-cover-treemap（树状图可视化）

`go-cover-treemap` 生成 SVG 树状图，直观展示各包的覆盖率情况。

#### 安装

```bash
go install github.com/nikolaydubina/go-cover-treemap@latest
```

#### 使用

```bash
# 生成覆盖率文件
go test -tags=test -coverprofile=coverage/coverage.out ./...

# 生成树状图
go-cover-treemap -coverprofile=coverage/coverage.out > coverage/coverage-treemap.svg

# 打开报告
open coverage/coverage-treemap.svg
```

#### 使用 Makefile

```bash
# 生成覆盖率树状图
make test-coverage-treemap
```

**特点**：
- 📊 树状图可视化，直观展示覆盖率分布
- 🎨 SVG 格式，可缩放
- 🔍 快速识别低覆盖率区域

### 工具对比

| 工具 | 类型 | 界面 | 交互性 | 推荐场景 |
|------|------|------|--------|----------|
| `go tool cover` | HTML | 简单 | 低 | 快速查看，CI/CD |
| `gocovsh` | 终端 | 美观 | 高 | **日常开发，最接近 cargo-tarpaulin** |
| `gocov + gocov-html` | HTML | 中等 | 中 | 代码审查，HTML 报告 |
| `go-cover-treemap` | SVG | 美观 | 低 | 覆盖率概览，可视化展示 |

### 推荐使用方式

1. **日常开发**（最推荐）：使用 `gocovsh` 快速查看覆盖率，体验最接近 `cargo-tarpaulin`
   ```bash
   make test-coverage-interactive
   ```
   使用方向键导航，Enter 选择，Esc 退出

2. **代码审查**：使用 `gocov + gocov-html` 生成 HTML 报告
   ```bash
   make test-coverage-ui
   make open-coverage-ui
   ```

3. **可视化概览**：使用 `go-cover-treemap` 生成树状图
   ```bash
   make test-coverage-treemap
   ```

4. **CI/CD**：使用标准 `go tool cover` 生成报告
   ```bash
   make test-coverage
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

---

## 快速参考

### 标准覆盖率报告

```bash
make test-coverage        # 生成标准 HTML 报告
make open-coverage       # 打开标准报告
```

### 美观的覆盖率报告

```bash
make test-coverage-interactive  # 使用 gocovsh 终端 UI（最推荐，类似 cargo-tarpaulin）
make test-coverage-ui           # 生成美观的 HTML 报告（gocov-html）
make open-coverage-ui           # 打开美观的报告
make test-coverage-treemap      # 生成覆盖率树状图可视化
```

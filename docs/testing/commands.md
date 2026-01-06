# 测试命令参考

> 本文档提供常用测试命令的快速参考。

---

## 📋 目录

- [基本测试命令](#-基本测试命令)
- [测试类型命令](#-测试类型命令)
- [Makefile测试命令](#makefile测试命令)
- [测试调试](#-测试调试)

---

## 🚀 基本测试命令

### 运行测试

**运行所有测试**：
```bash
# 使用 Go
go test ./...

# 使用 Makefile
make test
```

**运行特定包的测试**：
```bash
# 运行特定包的测试
go test ./internal/lib/config

# 运行特定包的测试（显示详细输出）
go test -v ./internal/lib/config

# 运行匹配模式的测试
go test -run TestParseTicketID ./internal/lib/config
```

**测试输出选项**：
```bash
# 显示详细输出
go test -v ./...

# 显示测试执行时间
go test -v -timeout 30s ./...

# 只运行失败的测试（需要先运行一次）
go test -run TestFailed ./...
```

---

## 🎯 测试类型命令

### 单元测试

```bash
# 运行所有单元测试
go test ./...

# 运行特定包的单元测试
go test ./internal/lib/config

# 运行特定测试函数
go test -run TestParseTicketID ./internal/lib/config
```

### 集成测试

```bash
# 运行集成测试（使用构建标签）
go test -tags=integration ./test/integration

# 运行所有测试（包括集成测试）
go test -tags=integration ./...
```

### 基准测试

```bash
# 运行所有基准测试
go test -bench=. ./...

# 运行特定包的基准测试
go test -bench=. ./internal/lib/config

# 运行基准测试并显示内存分配
go test -bench=. -benchmem ./...

# 运行基准测试并生成CPU profile
go test -bench=. -cpuprofile=cpu.prof ./...

# 运行基准测试并生成内存profile
go test -bench=. -memprofile=mem.prof ./...
```

### 示例测试

```bash
# 运行示例测试（Example functions）
go test -run Example ./...

# 运行特定包的示例测试
go test -run Example ./internal/lib/config
```

---

## Makefile测试命令

项目提供了便捷的 Makefile 命令：

```bash
# 运行所有测试
make test

# 运行测试并生成覆盖率报告
make test-coverage

# 查看覆盖率报告
open coverage.html  # macOS
xdg-open coverage.html  # Linux
```

### Makefile 测试命令详情

```makefile
# 运行测试
test:
	@echo "运行测试..."
	@go test -v ./...

# 运行测试并生成覆盖率报告
test-coverage:
	@echo "运行测试并生成覆盖率报告..."
	@go test -v -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html
	@echo "覆盖率报告已生成: coverage.html"
```

---

## 🐛 测试调试

### 运行单个测试

```bash
# 运行单个测试函数
go test -run TestParseTicketID ./internal/lib/config

# 运行单个测试并显示详细输出
go test -v -run TestParseTicketID ./internal/lib/config

# 运行单个测试并显示覆盖率
go test -cover -run TestParseTicketID ./internal/lib/config
```

### 测试失败时调试

```bash
# 显示失败的测试输出
go test -v ./...

# 只运行失败的测试（需要先运行一次，保存失败信息）
go test -run TestFailed ./...

# 显示测试执行时间（找出慢测试）
go test -v -timeout 30s ./...
```

### 测试覆盖率调试

```bash
# 显示覆盖率
go test -cover ./...

# 显示每个函数的覆盖率
go test -coverprofile=coverage.out ./...
go tool cover -func=coverage.out

# 生成HTML覆盖率报告
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out -o coverage.html
```

### 测试性能分析

```bash
# 生成CPU profile
go test -cpuprofile=cpu.prof -bench=. ./internal/lib/config
go tool pprof cpu.prof

# 生成内存profile
go test -memprofile=mem.prof -bench=. ./internal/lib/config
go tool pprof mem.prof

# 生成trace文件
go test -trace=trace.out ./internal/lib/config
go tool trace trace.out
```

---

## 📊 常用命令组合

### 开发时常用

```bash
# 快速测试（只运行当前包的测试）
go test ./internal/lib/config

# 详细测试输出
go test -v ./internal/lib/config

# 测试并显示覆盖率
go test -cover ./internal/lib/config

# 测试并生成覆盖率报告
make test-coverage && open coverage.html
```

### CI 环境常用

```bash
# 运行所有测试
go test ./...

# 运行所有测试（包括集成测试）
go test -tags=integration ./...

# 生成覆盖率报告
make test-coverage

# 运行基准测试
go test -bench=. ./...
```

### 调试时常用

```bash
# 运行单个测试
go test -v -run TestParseTicketID ./internal/lib/config

# 运行测试并显示覆盖率
go test -cover -run TestParseTicketID ./internal/lib/config

# 运行测试并生成profile
go test -cpuprofile=cpu.prof -run TestParseTicketID ./internal/lib/config
```

---

## 🔍 测试命令选项详解

### 基本选项

| 选项 | 说明 | 示例 |
|------|------|------|
| `-v` | 显示详细输出 | `go test -v ./...` |
| `-run` | 运行匹配模式的测试 | `go test -run TestParse ./...` |
| `-cover` | 显示覆盖率 | `go test -cover ./...` |
| `-coverprofile` | 生成覆盖率文件 | `go test -coverprofile=coverage.out ./...` |
| `-timeout` | 设置超时时间 | `go test -timeout 30s ./...` |
| `-count` | 运行测试的次数 | `go test -count=3 ./...` |
| `-parallel` | 并行运行测试的数量 | `go test -parallel=4 ./...` |

### 基准测试选项

| 选项 | 说明 | 示例 |
|------|------|------|
| `-bench` | 运行基准测试 | `go test -bench=. ./...` |
| `-benchmem` | 显示内存分配 | `go test -bench=. -benchmem ./...` |
| `-benchtime` | 设置基准测试时间 | `go test -bench=. -benchtime=5s ./...` |

### 性能分析选项

| 选项 | 说明 | 示例 |
|------|------|------|
| `-cpuprofile` | 生成CPU profile | `go test -cpuprofile=cpu.prof ./...` |
| `-memprofile` | 生成内存profile | `go test -memprofile=mem.prof ./...` |
| `-trace` | 生成trace文件 | `go test -trace=trace.out ./...` |

### 构建标签选项

| 选项 | 说明 | 示例 |
|------|------|------|
| `-tags` | 使用构建标签 | `go test -tags=integration ./...` |

---

## 📚 相关文档

- [测试组织规范](./organization.md) - 测试组织结构
- [测试编写规范](./writing.md) - 测试编写规范
- [覆盖率测试指南](./references/coverage.md) - 覆盖率工具详细使用

---

**最后更新**: 2025-01-28

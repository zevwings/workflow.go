# 开发工具规范

> 本文档定义了 Workflow CLI 项目开发工具的使用规范和最佳实践，所有贡献者都应遵循这些规范。

---

## 📋 目录

- [概述](#-概述)
- [必需工具](#-必需工具)
- [常用命令](#-常用命令)
- [工具配置](#-工具配置)
- [相关文档](#-相关文档)

---

## 📋 概述

本文档定义了开发工具的使用规范，包括必需工具、常用命令和工具配置。

### 核心原则

- **工具统一**：使用统一的开发工具和命令
- **自动化**：使用工具自动化检查和格式化
- **持续集成**：工具检查应集成到 CI/CD 流程

### 使用场景

- 开发环境设置时参考
- 日常开发时使用
- CI/CD 配置时参考

---

## 必需工具

### Go 工具链

- **go**：Go 编译器和工具链（Go 1.21 或更高版本）
- **gofmt**：代码格式化工具（Go 标准工具）
- **goimports**：导入语句管理工具（推荐）
- **golangci-lint**：代码检查工具（推荐）

### 安装方法

```bash
# 安装 Go（如果未安装）
# macOS
brew install go

# Linux
sudo apt-get install golang-go  # Ubuntu/Debian
sudo dnf install golang         # Fedora/RHEL

# 或从官网下载：https://go.dev/dl/

# 安装 goimports（推荐）
go install golang.org/x/tools/cmd/goimports@latest

# 安装 golangci-lint（推荐）
# macOS
brew install golangci-lint

# Linux
# Ubuntu/Debian
sudo apt-get install golangci-lint

# Fedora/RHEL
sudo dnf install golangci-lint

# 或使用 go install
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
```

### 开发工具（可选）

- **gofumpt**：更严格的格式化工具（可选）
- **govulncheck**：Go 漏洞检查工具（推荐）
- **gocov**：测试覆盖率工具（可选，Go 内置 `go test -cover` 已足够）
- **go tool pprof**：性能分析工具（Go 标准工具）

### 安装方法

```bash
# 安装 gofumpt（可选，更严格的格式化）
go install mvdan.cc/gofumpt@latest

# 安装 govulncheck（推荐，安全漏洞检查）
go install golang.org/x/vuln/cmd/govulncheck@latest

# 安装 gocov（可选，如果需要 HTML 报告）
go install github.com/axw/gocov/gocov@latest
go install github.com/AlekSi/gocov-xml@latest
```

---

## 常用命令

### 代码格式化

```bash
# 格式化代码
go fmt ./...

# 或使用 goimports（推荐，自动管理导入）
goimports -w .

# 检查代码格式（CI/CD 中使用）
gofmt -l .

# 或使用 gofumpt（可选，更严格的格式化）
gofumpt -w .
gofumpt -l .  # 检查格式
```

### 代码检查

```bash
# 运行 golangci-lint 检查
golangci-lint run

# 运行特定检查
golangci-lint run --enable-all --disable-all -E errcheck -E gosec

# 自动修复可修复的问题
golangci-lint run --fix

# 或使用 Makefile
make lint
```

### 测试

```bash
# 运行所有测试
go test ./...

# 运行特定包的测试
go test ./internal/lib/config

# 运行特定测试函数
go test -run TestParseTicketID ./internal/lib/config

# 运行测试并显示详细输出
go test -v ./...

# 检查测试覆盖率
go test -cover ./...

# 生成覆盖率报告
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out -o coverage.html

# 或使用 Makefile
make test
make test-coverage
```

### 构建

```bash
# 构建二进制文件
go build -o bin/workflow ./cmd/workflow

# 构建所有包（检查编译，不生成二进制）
go build ./...

# 检查编译（不生成二进制文件）
go build ./...

# 安装到 $GOPATH/bin
go install ./cmd/workflow

# 或使用 Makefile
make build
```

### 依赖管理

```bash
# 添加依赖
go get package-name

# 添加特定版本的依赖
go get package-name@v1.2.3

# 更新依赖到最新版本
go get -u package-name

# 更新所有依赖
go get -u ./...

# 整理依赖（移除未使用的依赖）
go mod tidy

# 下载依赖
go mod download

# 检查依赖安全漏洞（使用 govulncheck）
govulncheck ./...

# 查看依赖关系
go list -m all

# 查看特定包的依赖
go list -m -json package-name
```

### 性能分析

```bash
# 运行基准测试
go test -bench=. ./...

# 运行基准测试并显示内存分配
go test -bench=. -benchmem ./...

# CPU 性能分析
go test -bench=. -cpuprofile=cpu.prof ./internal/lib/module
go tool pprof cpu.prof

# 内存性能分析
go test -bench=. -memprofile=mem.prof ./internal/lib/module
go tool pprof mem.prof

# 分析二进制大小
go build -ldflags="-s -w" -o bin/workflow ./cmd/workflow
ls -lh bin/workflow
```

---

## 工具配置

### gofmt / goimports 配置

Go 的格式化工具使用 Go 官方代码风格，无需额外配置。如果需要更严格的格式化，可以使用 `gofumpt`。

### golangci-lint 配置

项目根目录的 `.golangci.yml` 文件配置 golangci-lint 检查规则：

```yaml
linters-settings:
  errcheck:
    check-type-assertions: true
    check-blank: true
  gosec:
    severity: medium
  govet:
    check-shadowing: true
  unused:
    check-exported: false

linters:
  enable:
    - errcheck
    - gosec
    - govet
    - staticcheck
    - unused
    - gofmt
    - goimports
    - gocritic
    - revive

issues:
  exclude-use-default: false
  max-issues-per-linter: 0
  max-same-issues: 0
```

### go.mod 版本配置

项目的 Go 版本要求在 `go.mod` 文件中指定：

```go
module github.com/zevwings/workflow

go 1.21
```

### Makefile 命令

项目提供了 Makefile 命令简化常用操作：

```bash
# 格式化代码
make fmt

# 运行 Lint 检查
make lint

# 运行测试
make test

# 运行所有检查
make check
```

---

## 🔍 故障排除

### 问题 1：工具未安装

**症状**：运行命令时提示工具未找到

**解决方案**：

1. 检查工具是否已安装
2. 检查 `PATH` 环境变量
3. 重新安装工具

### 问题 2：工具版本不兼容

**症状**：工具版本不兼容导致错误

**解决方案**：

1. 更新工具到最新版本
2. 检查项目要求的 Go 版本（`go.mod` 中的 `go` 指令）
3. 确保 Go 版本满足要求（Go 1.21 或更高版本）
4. 更新 golangci-lint：`go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest`

---

## 📚 相关文档

### 开发规范

- [代码风格规范](../code-style.md) - 代码风格规范（包含 gofmt 和 golangci-lint 使用）
- [依赖管理规范](./dependency-management.md) - 依赖管理规范（包含 govulncheck 使用）
- [性能优化规范](./performance.md) - 性能优化规范（包含基准测试工具）

### 工具文档

- [Go 官方文档](https://go.dev/doc/) - Go 语言和工具链文档
- [golangci-lint 文档](https://golangci-lint.run/) - golangci-lint 官方文档
- [go.mod 文档](https://go.dev/doc/modules/gomod-ref) - Go 模块管理文档

---

## ✅ 检查清单

使用本规范时，请确保：

- [ ] 必需工具已安装
- [ ] 工具配置已设置
- [ ] 常用命令已熟悉
- [ ] CI/CD 已集成工具检查

---

**最后更新**: 2025-01-27


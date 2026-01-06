# 代码风格规范

> 本文档定义了 Workflow CLI 项目的代码风格规范和最佳实践，所有贡献者都应遵循这些规范。

---

## 📋 目录

- [概述](#-概述)
- [代码格式化](#-代码格式化)
- [Lint 检查](#-lint-检查)
- [Go 命名约定](#-go-命名约定)
- [代码组织](#-代码组织)
- [相关文档](#-相关文档)

---

## 📋 概述

本文档定义了代码风格规范，包括代码格式化、Lint 检查、Go 命名约定和代码组织规范。

### 核心原则

- **一致性**：所有代码必须遵循统一的风格规范
- **自动化**：使用工具自动检查和格式化代码
- **可读性**：代码风格应提高代码可读性
- **遵循标准**：遵循 Go 官方代码规范和社区最佳实践

### 使用场景

- 编写新代码时参考
- 代码审查时检查
- 代码格式化时使用

### 快速参考

| 操作 | 命令 | 说明 |
|------|------|------|
| **格式化代码** | `go fmt ./...` | 自动格式化代码 |
| **格式化并整理导入** | `goimports -w .` | 格式化并自动管理导入 |
| **检查格式** | `gofmt -l .` | 检查代码格式（CI/CD） |
| **Lint 检查** | `golangci-lint run` | 运行 golangci-lint 检查 |
| **Lint 检查（Makefile）** | `make lint` | 使用 Makefile 运行 Lint |
| **格式化（Makefile）** | `make fmt` | 使用 Makefile 格式化代码 |

---

## 代码格式化

所有代码必须使用 Go 官方工具进行格式化：

### 使用 go fmt

```bash
# 自动格式化所有代码
go fmt ./...

# 或使用 Makefile
make fmt
```

**规则**：
- 提交前必须运行 `go fmt ./...`
- CI/CD 会检查代码格式，格式不正确会导致构建失败
- `go fmt` 会自动应用 Go 官方代码风格

### 使用 goimports

`goimports` 是 `gofmt` 的增强版本，会自动管理导入语句：

```bash
# 安装 goimports（如果未安装）
go install golang.org/x/tools/cmd/goimports@latest

# 格式化并整理导入
goimports -w .

# 检查导入（不修改文件）
goimports -l .
```

**规则**：
- 推荐使用 `goimports` 替代 `gofmt`，因为它会自动添加缺失的导入并移除未使用的导入
- 提交前确保所有导入都已正确整理

### 使用 gofumpt（可选，推荐）

`gofumpt` 是 `gofmt` 的严格版本，提供更严格的格式化规则：

```bash
# 安装 gofumpt（如果未安装）
go install mvdan.cc/gofumpt@latest

# 格式化代码
gofumpt -w .

# 检查格式
gofumpt -l .
```

**规则**：
- 如果项目使用 `gofumpt`，所有代码必须通过 `gofumpt` 检查
- `gofumpt` 与 `gofmt` 兼容，但规则更严格

---

## Lint 检查

使用 `golangci-lint` 进行代码质量检查：

### 安装 golangci-lint

```bash
# 安装最新版本
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# 或使用 Homebrew (macOS)
brew install golangci-lint

# 或使用包管理器 (Linux)
# Ubuntu/Debian
sudo apt-get install golangci-lint

# Fedora/RHEL
sudo dnf install golangci-lint
```

### 运行 Lint 检查

```bash
# 运行所有检查
golangci-lint run

# 运行特定检查
golangci-lint run --enable-all --disable-all -E errcheck -E gosec

# 或使用 Makefile
make lint
```

**规则**：
- 所有警告必须修复（除非有充分理由并添加注释说明）
- 禁止使用 `//nolint` 注释除非有充分理由，并添加注释说明原因
- 定期运行 `golangci-lint run` 检查代码质量
- CI/CD 会运行 Lint 检查，未通过的检查会导致构建失败

### 常用 Lint 规则

项目推荐启用以下 Lint 规则：

- **errcheck**：检查错误处理
- **gosec**：安全检查
- **govet**：Go 官方 vet 工具
- **staticcheck**：静态分析
- **unused**：未使用的代码检查
- **gofmt**：代码格式检查
- **goimports**：导入检查

### 配置 golangci-lint

可以在项目根目录创建 `.golangci.yml` 配置文件：

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

---

## Go 命名约定

遵循 Go 官方命名约定（[Effective Go](https://go.dev/doc/effective_go#names)）：

### 包名

- **包名**：小写单词，不使用下划线或混合大小写（如 `http`、`json`、`config`）
- **简短**：包名应该简短、清晰
- **避免冲突**：避免使用标准库已有的包名

```go
// ✅ 好的包名
package http
package config
package logging

// ❌ 不好的包名
package HTTPClient  // 应该使用小写
package my_package  // 不应该使用下划线
package MyPackage   // 不应该使用大写
```

### 导出标识符

- **导出标识符**：首字母大写（如 `Client`、`NewClient`、`GetUser`）
- **未导出标识符**：首字母小写（如 `client`、`newClient`、`getUser`）

```go
// ✅ 导出函数（公共 API）
func NewClient() *Client {
    return &Client{}
}

// ✅ 未导出函数（内部使用）
func newClient() *client {
    return &client{}
}
```

### 函数名

- **函数名**：使用驼峰命名（如 `GetUser`、`CreateTicket`、`DownloadLogs`）
- **Getter**：不需要 `Get` 前缀（如 `User()` 而不是 `GetUser()`）
- **Setter**：使用 `Set` 前缀（如 `SetTimeout`）

```go
// ✅ 好的函数名
func User() *User { }
func SetTimeout(d time.Duration) { }
func CreateTicket(id string) error { }

// ❌ 不好的函数名
func GetUser() *User { }  // Getter 不需要 Get 前缀
func create_ticket() { }  // 不应该使用下划线
```

### 变量名

- **变量名**：使用驼峰命名，首字母小写（如 `userID`、`apiToken`、`responseData`）
- **简短**：局部变量应该简短（如 `i`、`j`、`err`）
- **描述性**：包级别变量应该描述性（如 `defaultTimeout`、`maxRetries`）

```go
// ✅ 好的变量名
var defaultTimeout = 30 * time.Second
var maxRetries = 3

func processUser(userID string) {
    id := userID  // 局部变量可以简短
    err := doSomething()
}

// ❌ 不好的变量名
var DefaultTimeout = 30 * time.Second  // 包级别变量不应该导出（除非需要）
var max_retries = 3  // 不应该使用下划线
```

### 常量名

- **常量名**：使用驼峰命名，首字母大写（导出）或小写（未导出）
- **特殊常量**：可以使用 `SCREAMING_SNAKE_CASE`（如 `MAX_RETRIES`、`DEFAULT_TIMEOUT`），但不推荐

```go
// ✅ 好的常量名
const DefaultTimeout = 30 * time.Second
const maxRetries = 3

// 也可以使用 SCREAMING_SNAKE_CASE（但不推荐）
const MAX_RETRIES = 3
const DEFAULT_TIMEOUT = 30 * time.Second
```

### 类型名

- **类型名**：使用驼峰命名，首字母大写（如 `Client`、`HTTPClient`、`UserInfo`）
- **接口名**：通常以 `-er` 结尾（如 `Reader`、`Writer`、`Closer`），或使用描述性名称（如 `Client`、`Config`）

```go
// ✅ 好的类型名
type Client struct { }
type HTTPClient struct { }
type UserInfo struct { }

// ✅ 好的接口名
type Reader interface {
    Read([]byte) (int, error)
}

type Client interface {
    Get(url string) (*Response, error)
}
```

### 方法接收者名

- **接收者名**：应该简短，通常是类型名的首字母小写（如 `c *Client`、`h *HTTPClient`）
- **一致性**：同一类型的所有方法应该使用相同的接收者名

```go
// ✅ 好的接收者名
func (c *Client) Get(url string) (*Response, error) { }
func (c *Client) Post(url string, body interface{}) (*Response, error) { }

// ❌ 不好的接收者名
func (client *Client) Get(url string) (*Response, error) { }  // 应该使用简短名称
func (cl *Client) Get(url string) (*Response, error) { }     // 不一致
```

---

## 代码组织

### 导入顺序

Go 的导入语句应该按以下顺序组织：

1. 标准库导入
2. 第三方库导入
3. 项目内部导入

每组导入之间用空行分隔：

```go
package http

import (
    // 标准库
    "fmt"
    "net/http"
    "time"

    // 第三方库
    "github.com/go-resty/resty/v2"

    // 项目内部
    "github.com/your-org/workflow/internal/lib/config"
    "github.com/your-org/workflow/internal/logging"
)
```

**规则**：
- 使用 `goimports` 自动管理导入顺序
- 每组导入之间用空行分隔
- 导入路径按字母顺序排序

### 包声明

- **包名**：应该与目录名一致（如 `package http` 在 `internal/lib/http/` 目录中）
- **main 包**：只有 `main.go` 文件使用 `package main`

```go
// internal/lib/http/client.go
package http  // ✅ 包名与目录名一致

// cmd/workflow/main.go
package main  // ✅ main 包用于可执行文件
```

### 文件组织

Go 文件应该按以下顺序组织：

1. 包声明
2. 导入语句
3. 常量声明
4. 变量声明
5. 类型声明
6. 函数/方法实现

```go
package http

import (
    "time"
    "github.com/go-resty/resty/v2"
)

// 常量
const (
    DefaultTimeout = 30 * time.Second
    MaxRetries     = 3
)

// 变量
var (
    defaultClient *Client
)

// 类型
type Client struct {
    client *resty.Client
}

// 函数
func NewClient() *Client {
    return &Client{
        client: resty.New(),
    }
}

// 方法
func (c *Client) Get(url string) (*resty.Response, error) {
    return c.client.R().Get(url)
}
```

### 目录结构

遵循 Go 标准项目布局：

```
workflow/
├── cmd/
│   └── workflow/          # 主入口
│       └── main.go
├── internal/             # 内部包（不对外暴露）
│   ├── cli/              # CLI 根命令
│   ├── commands/         # 命令实现
│   ├── lib/              # 核心业务逻辑
│   │   ├── config/       # 配置管理
│   │   ├── http/         # HTTP 客户端
│   │   └── ...
│   └── logging/          # 日志系统
├── pkg/                  # 公共包（可选，对外暴露）
├── go.mod
├── go.sum
└── Makefile
```

**规则**：
- `cmd/`：可执行文件入口
- `internal/`：内部包，不允许外部导入
- `pkg/`：公共包，允许外部导入（如果项目需要）
- 每个目录一个包，包名与目录名一致

---

## 🔍 故障排除

### 问题 1：代码格式检查失败

**症状**：运行 `gofmt -l .` 时提示格式不正确

**解决方案**：

1. 运行 `go fmt ./...` 自动格式化代码
2. 运行 `goimports -w .` 整理导入语句
3. 确保使用最新版本的 Go 工具链
4. 检查是否有自定义的格式化配置

### 问题 2：golangci-lint 警告过多

**症状**：运行 `golangci-lint run` 时出现大量警告

**解决方案**：

1. 逐个修复警告（优先修复高优先级警告）
2. 对于确实需要忽略的警告，使用 `//nolint:linter-name` 并添加注释说明原因
3. 定期运行 `golangci-lint run` 保持代码质量
4. 检查 `.golangci.yml` 配置文件，调整规则

### 问题 3：导入顺序不正确

**症状**：导入语句顺序不符合规范

**解决方案**：

1. 使用 `goimports -w .` 自动整理导入
2. 确保标准库、第三方库、项目内部导入分组正确
3. 每组导入之间用空行分隔

### 问题 4：包名与目录名不一致

**症状**：包声明与目录名不一致

**解决方案**：

1. 确保包名与目录名一致
2. 检查是否有拼写错误
3. 使用 `goimports` 自动修复

---

## 📚 相关文档

### 开发规范

- [错误处理规范](./error-handling.md) - 错误处理规范
- [命名规范](./naming.md) - 命名规范
- [模块组织规范](./module-organization.md) - 模块组织规范

### 检查工作流

- [提交前检查](./workflows/pre-commit.md) - 代码质量检查流程

### Go 官方文档

- [Effective Go](https://go.dev/doc/effective_go) - Go 官方最佳实践
- [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments) - Go 代码审查注释
- [golangci-lint 文档](https://golangci-lint.run/) - golangci-lint 官方文档

---

## ✅ 检查清单

使用本规范时，请确保：

- [ ] 代码已格式化（`go fmt ./...` 或 `make fmt`）
- [ ] 导入已整理（`goimports -w .`）
- [ ] 通过 golangci-lint 检查（`golangci-lint run` 或 `make lint`）
- [ ] 遵循 Go 命名约定
- [ ] 导入顺序正确（标准库 → 第三方库 → 项目内部）
- [ ] 包名与目录名一致
- [ ] 文件组织符合规范

---

**最后更新**: 2025-01-27

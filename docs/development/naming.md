# 命名规范

> 本文档定义了 Workflow CLI 项目的命名规范和最佳实践，所有贡献者都应遵循这些规范。

---

## 📋 目录

- [概述](#-概述)
- [文件命名](#-文件命名)
- [函数命名](#-函数命名)
- [结构体命名](#-结构体命名)
- [常量命名](#-常量命名)
- [CLI 参数命名规范](#-cli-参数命名规范)
- [相关文档](#-相关文档)

---

## 📋 概述

本文档定义了命名规范，包括文件、函数、结构体、常量和 CLI 参数的命名规范。

### 核心原则

- **一致性**：相同语义的命名必须保持一致
- **可读性**：命名应清晰表达意图
- **遵循约定**：遵循 Go 官方命名约定（[Effective Go](https://go.dev/doc/effective_go#names)）

### 使用场景

- 编写新代码时参考
- 代码审查时检查
- 重构代码时使用

---

## 文件命名

- **Go 源文件**：`snake_case.go`（如 `jira_client.go`、`pr_helpers.go`）
- **测试文件**：`*_test.go`（如 `jira_client_test.go`），与源文件在同一目录
- **文档文件**：`kebab-case.md`（如 `testing/README.md`、`pr.md`）
  - **架构文档**：`{module}.md`（如 `pr.md`、`git.md`）
  - **指南文档**：`{topic}.md`（如 `testing/README.md`、`document.md`）
  - **需求文档**：`{topic}.md`（如 `jira.md`、`integration.md`，存放到 `docs/requirements/`）
  - **迁移文档**：`{version}-to-{version}.md`（如 `1.5.6-to-1.5.7.md`）

**规则**：
- Go 文件名使用 `snake_case`，与包名保持一致
- 测试文件必须以 `_test.go` 结尾
- 文档文件使用 `kebab-case`

---

## 函数命名

遵循 Go 官方命名约定：

- **导出函数**：首字母大写，使用驼峰命名（如 `DownloadLogs`、`CreateTicket`、`GetUser`）
- **未导出函数**：首字母小写，使用驼峰命名（如 `downloadLogs`、`createTicket`、`getUser`）
- **动作函数**：使用动词（如 `Download`、`Create`、`Merge`）
- **Getter 函数**：不需要 `Get` 前缀（如 `User()` 而不是 `GetUser()`）
- **Setter 函数**：使用 `Set` 前缀（如 `SetTimeout`）
- **检查函数**：使用 `Is` 或 `Has` 前缀（如 `IsValid`、`HasPermission`）
- **转换函数**：使用 `To` 前缀（如 `ToString`、`ToJSON`）

```go
// ✅ 好的函数名
func DownloadLogs(ticketID string) ([]byte, error) { }
func CreateTicket(title string) (*Ticket, error) { }
func User() *User { }  // Getter 不需要 Get 前缀
func SetTimeout(d time.Duration) { }  // Setter 使用 Set 前缀
func IsValid(id string) bool { }  // 检查函数使用 Is 前缀

// ❌ 不好的函数名
func get_user() { }  // 不应该使用下划线
func GetUser() { }  // Getter 不需要 Get 前缀（除非有特殊原因）
func create_ticket() { }  // 不应该使用下划线
```

---

## 结构体命名

- **导出结构体**：首字母大写，使用驼峰命名（如 `HTTPClient`、`JiraTicket`、`UserInfo`）
- **未导出结构体**：首字母小写，使用驼峰命名（如 `httpClient`、`jiraTicket`、`userInfo`）
- 使用名词或名词短语
- 避免使用 `Data`、`Info`、`Manager` 等泛化名称，使用具体名称

```go
// ✅ 好的结构体名
type HTTPClient struct { }
type JiraTicket struct { }
type UserInfo struct { }

// ❌ 不好的结构体名
type Data struct { }  // 太泛化
type Info struct { }  // 太泛化
type Manager struct { }  // 太泛化
```

---

## 常量命名

- **导出常量**：首字母大写，使用驼峰命名（如 `DefaultTimeout`、`MaxRetries`）
- **未导出常量**：首字母小写，使用驼峰命名（如 `defaultTimeout`、`maxRetries`）
- 也可以使用 `SCREAMING_SNAKE_CASE`（如 `MAX_RETRIES`、`DEFAULT_TIMEOUT`），但不推荐

```go
// ✅ 好的常量名（推荐）
const (
    DefaultTimeout = 30 * time.Second
    MaxRetries    = 3
)

// ✅ 也可以使用 SCREAMING_SNAKE_CASE（但不推荐）
const (
    MAX_RETRIES     = 3
    DEFAULT_TIMEOUT = 30 * time.Second
)

// ❌ 不好的常量名
const max_retries = 3  // 不应该使用下划线
```

---

## CLI 参数命名规范

CLI 参数命名需要遵循以下规范，确保一致性和可维护性。

### 结构体字段名

- 使用 `camelCase`（如 `jiraID`、`dryRun`、`outputFormat`）

```go
// ✅ 好的做法
type JiraIDArg struct {
    JiraID string `json:"jira_id"`  // camelCase
}

// ❌ 不好的做法
type JiraIDArg struct {
    Jira_id string  // 不应该使用下划线
    jiraID  string  // 字段名应该导出（如果需要）
}
```

### 参数长名规范

- 使用 `kebab-case`（如 `--jira-id`、`--dry-run`）
- Cobra 会自动从字段名转换（`JiraID` → `--jira-id`）

```go
// ✅ 好的做法
type Args struct {
    JiraID    string `mapstructure:"jira_id"`  // 自动生成 --jira-id
    DryRun    bool   `mapstructure:"dry_run"`    // 自动生成 --dry-run
    OutputFormat string `mapstructure:"output_format"`  // 自动生成 --output-format
}
```

### 参数短名规范

- 使用单个字符（如 `-n`、`-f`、`-v`）
- 优先使用常见的短名（如 `-n` 用于 dry-run，`-f` 用于 force）

```go
// ✅ 好的做法
cmd.Flags().BoolVarP(&dryRun, "dry-run", "n", false, "Dry run mode")
cmd.Flags().BoolVarP(&force, "force", "f", false, "Force operation")
```

### 参数类型规范

- **可选参数**：使用指针类型（如 `*string`、`*int`）或提供默认值
- **必需参数**：直接使用类型（如 `string`、`int`）
- **布尔标志**：使用 `bool` 类型

```go
// ✅ 可选参数
type Args struct {
    JiraID *string  // 可选
}

// ✅ 必需参数
type Args struct {
    BranchName string  // 必需
}

// ✅ 布尔标志
type Args struct {
    Force bool  // 布尔标志
}
```

### 文档注释规范

所有参数必须有文档注释，说明参数的用途、格式和默认行为：

```go
// JiraID 是 JIRA ticket ID（可选，如果未提供将交互式提示）
//
// 示例:
//   workflow jira info PROJ-123
//   workflow jira info  # 将提示输入 JIRA ID
type JiraIDArg struct {
    JiraID string
}
```

### 命名一致性规范

- **相同语义的参数必须使用相同的命名**：
  - ✅ 统一使用 `jiraID`（而不是 `jiraTicket`、`jira-id` 等）
  - ✅ 统一使用 `dryRun`（而不是 `dry-run`、`dryrun` 等）
  - ✅ 统一使用 `outputFormat`（而不是 `format`、`output` 等）

### 共用参数规范

对于在多个命令中重复使用的参数，应该提取为共用参数结构（见 [CLI 检查指南](./references/review-cli.md)）：

```go
// internal/lib/cli/args.go
// JiraIDArg 可选 JIRA ID 参数
type JiraIDArg struct {
    // JiraID 是 JIRA ticket ID（可选，如果未提供将交互式提示）
    JiraID string
}

// 在命令中使用
type MyCommand struct {
    JiraIDArg  // ✅ 使用共用参数
}
```

### 示例对比

```go
// ❌ 不好的做法
type CreateArgs struct {
    JiraTicket string  // 命名不一致（应该用 JiraID）
}

// ✅ 好的做法
type CreateArgs struct {
    JiraIDArg  // 使用共用参数，命名一致
}
```

**参考**：
- [CLI 检查指南](./references/review-cli.md) - 参数复用检查和参数提取指南
- [Cobra 文档](https://github.com/spf13/cobra) - Cobra 参数定义规范

---

## 🔍 故障排除

### 问题 1：命名不一致

**症状**：相同语义的命名在不同地方使用了不同的形式

**解决方案**：

1. 统一使用项目命名规范
2. 使用共用参数结构避免重复定义
3. 在代码审查时检查命名一致性

### 问题 2：CLI 参数命名混乱

**症状**：CLI 参数命名不规范，导致帮助信息不清晰

**解决方案**：

1. 遵循 CLI 参数命名规范
2. 使用清晰的字段名和文档注释
3. 添加文档注释说明参数用途

---

## 📚 相关文档

### 开发规范

- [代码风格规范](./code-style.md) - 代码风格规范
- [模块组织规范](./module-organization.md) - 模块组织规范

### 检查工作流

- [CLI 检查指南](./references/review-cli.md) - CLI 参数检查流程

### Go 官方文档

- [Effective Go - Names](https://go.dev/doc/effective_go#names) - Go 命名约定
- [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments) - Go 代码审查注释

---

## ✅ 检查清单

使用本规范时，请确保：

- [ ] 文件命名遵循规范
- [ ] 函数命名清晰表达意图
- [ ] 结构体命名使用具体名称
- [ ] 常量命名使用驼峰命名（或 `SCREAMING_SNAKE_CASE`）
- [ ] CLI 参数命名一致
- [ ] 使用共用参数结构避免重复

---

**最后更新**: 2025-01-27

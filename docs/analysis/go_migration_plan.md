# Workflow CLI Go 迁移方案

## 📋 文档概述

本文档详细规划了 Workflow CLI 从 Rust 迁移到 Go 的具体实施方案，包括命令结构设计、实现优先级、技术实现细节和迁移步骤。

**迁移策略：**
- **CLI 框架**：Cobra + Viper（推荐方案）
- **实现方式**：分阶段迁移，优先核心功能
- **兼容性**：保持与现有 Rust 版本的 API 兼容性

---

## 🎯 迁移目标

### 功能目标
1. **完整迁移**：所有命令功能对等迁移
2. **API 兼容**：保持命令行接口一致性
3. **配置兼容**：支持现有 TOML 配置文件格式

### 性能目标
- 启动速度：< 200ms
- 二进制体积：< 30MB
- 内存占用：< 50MB（运行时）

---

## 📐 命令结构设计

### 根命令结构

基于 Cobra 框架，命令结构如下：

```
workflow
├── setup              # 生命周期管理
├── update             # 生命周期管理
├── uninstall          # 生命周期管理
├── version            # 生命周期管理
├── config             # 配置管理
├── check              # 环境检查
├── github             # GitHub 账号管理
├── completion         # Shell Completion 管理
├── stash              # Stash 管理
├── repo               # 仓库管理
├── pr                 # PR 操作
└── jira               # Jira 操作
```

### 命令分组

#### 1. 生命周期管理组

**命令列表：**
- `workflow setup` - 初始化或更新配置（交互式）
- `workflow update [--version VERSION]` - 更新 Workflow CLI
- `workflow uninstall` - 卸载 Workflow CLI
- `workflow version` - 显示版本信息

**实现文件：**
```
internal/commands/
├── setup.go
├── update.go
├── uninstall.go
└── version.go
```

**技术要点：**
- `setup`：使用 `survey` 进行交互式配置
- `update`：实现版本检查和下载更新逻辑
- `uninstall`：清理配置文件和二进制文件
- `version`：显示版本、构建信息

---

#### 2. 配置管理组

**命令列表：**
- `workflow config` / `workflow config show` - 查看当前配置
- `workflow config validate [CONFIG_PATH] [--fix] [--strict]` - 验证配置文件
- `workflow config export <OUTPUT> [--section SECTION] [--no-secrets] [--toml|--json|--yaml]` - 导出配置
- `workflow config import <INPUT> [--overwrite] [--section SECTION] [--dry-run]` - 导入配置

**实现文件：**
```
internal/commands/
└── config.go          # 主命令
    ├── config_show.go
    ├── config_validate.go
    ├── config_export.go
    └── config_import.go
```

**技术要点：**
- 使用 Viper 管理配置
- 支持 TOML、JSON、YAML 格式
- 实现配置验证和修复逻辑
- 敏感信息过滤（`--no-secrets`）

---

#### 3. 环境检查组

**命令列表：**
- `workflow check` - 运行环境检查（Git 状态和网络连接）

**实现文件：**
```
internal/commands/
└── check.go
```

**技术要点：**
- 检查 Git 仓库状态
- 检查网络连接（GitHub、Jira）
- 检查配置文件完整性
- 使用表格显示检查结果

---

#### 4. GitHub 账号管理组

**命令列表：**
- `workflow github list` - 列出所有 GitHub 账号
- `workflow github current` - 显示当前激活的账号
- `workflow github add` - 添加新的 GitHub 账号
- `workflow github remove` - 删除 GitHub 账号
- `workflow github switch` - 切换当前 GitHub 账号
- `workflow github update` - 更新 GitHub 账号信息

**实现文件：**
```
internal/commands/
└── github.go          # 主命令
    ├── github_list.go
    ├── github_current.go
    ├── github_add.go
    ├── github_remove.go
    ├── github_switch.go
    └── github_update.go
```

**技术要点：**
- 使用 `google/go-github` 或自定义 HTTP 客户端
- 实现多账号管理逻辑
- 账号信息加密存储
- 交互式添加账号（使用 `survey`）

---

#### 5. Shell Completion 管理组

**命令列表：**
- `workflow completion generate` - 生成 completion 脚本
- `workflow completion check` - 检查 completion 状态
- `workflow completion remove` - 移除 completion 配置

**实现文件：**
```
internal/commands/
└── completion.go      # 主命令
    ├── completion_generate.go
    ├── completion_check.go
    └── completion_remove.go
```

**技术要点：**
- 使用 Cobra 内置的 completion 生成功能
- 支持 bash、zsh、fish、powershell
- 自动检测 Shell 类型
- 提供安装指导

---

#### 6. Stash 管理组

**命令列表：**
- `workflow stash list [--stat]` - 列出所有 stash
- `workflow stash apply` - 应用 stash（保留条目）
- `workflow stash drop` - 删除 stash
- `workflow stash pop` - 应用并删除 stash
- `workflow stash push` - 保存当前更改到 stash

**实现文件：**
```
internal/commands/
└── stash.go           # 主命令
    ├── stash_list.go
    ├── stash_apply.go
    ├── stash_drop.go
    ├── stash_pop.go
    └── stash_push.go
```

**技术要点：**
- 使用 `go-git` 或 `os/exec` 执行 git stash 命令
- 实现 stash 列表解析和显示
- 支持 `--stat` 选项显示统计信息

---

#### 7. 仓库管理组

**命令列表：**
- `workflow repo setup` - 配置项目级设置
- `workflow repo show` - 显示项目级配置
- `workflow repo clean [--dry-run]` - 清理本地分支和 tag

**实现文件：**
```
internal/commands/
└── repo.go            # 主命令
    ├── repo_setup.go
    ├── repo_show.go
    └── repo_clean.go
```

**技术要点：**
- 使用 `go-git` 进行仓库操作
- 实现分支和 tag 清理逻辑
- 支持 `--dry-run` 预览模式

---

#### 8. PR 操作组

**命令列表：**
- `workflow pr create [JIRA_TICKET] [--title TITLE] [--description DESC] [--dry-run]` - 创建 PR
- `workflow pr merge [PR_ID] [--force]` - 合并 PR
- `workflow pr close [PR_ID]` - 关闭 PR
- `workflow pr status [PR_ID_OR_BRANCH]` - 查看 PR 状态
- `workflow pr list [--state STATE] [--limit LIMIT]` - 列出 PR
- `workflow pr update` - 更新代码
- `workflow pr summarize [PR_ID] [--language LANG]` - 总结 PR
- `workflow pr approve [PR_ID]` - 批准 PR
- `workflow pr comment [PR_ID] <MESSAGE>` - 添加评论
- `workflow pr reword [PR_ID] [--title] [--description] [--dry-run]` - Reword PR 标题和描述

**实现文件：**
```
internal/commands/
└── pr.go              # 主命令
    ├── pr_create.go
    ├── pr_merge.go
    ├── pr_close.go
    ├── pr_status.go
    ├── pr_list.go
    ├── pr_update.go
    ├── pr_summarize.go
    ├── pr_approve.go
    ├── pr_comment.go
    └── pr_reword.go
```

**技术要点：**
- 使用 `google/go-github` 进行 GitHub API 调用
- 实现 Jira ticket 关联逻辑
- `pr summarize` 需要集成 LLM（OpenAI/DeepSeek）
- `pr reword` 需要集成 LLM 进行文本重写

---

#### 9. Jira 操作组

**命令列表：**
- `workflow jira info [PROJ-123] [--json|--markdown]` - 显示 ticket 信息
- `workflow jira related [PROJ-123] [--json|--markdown]` - 显示关联信息
- `workflow jira changelog [PROJ-123] [--json|--markdown]` - 显示变更历史
- `workflow jira comment [PROJ-123]` - 添加评论
- `workflow jira comments [PROJ-123] [--json|--markdown] [--limit LIMIT] [--offset OFFSET] [--author AUTHOR] [--since DATE]` - 显示评论
- `workflow jira attachments [PROJ-123]` - 下载所有附件
- `workflow jira clean [PROJ-123] [--all] [--dry-run] [--list]` - 清理日志目录

**实现文件：**
```
internal/commands/
└── jira.go            # 主命令
    ├── jira_info.go
    ├── jira_related.go
    ├── jira_changelog.go
    ├── jira_comment.go
    ├── jira_comments.go
    ├── jira_attachments.go
    └── jira_clean.go
```

**技术要点：**
- 使用 `andygrunwald/go-jira` 或自定义 HTTP 客户端
- 实现 Jira API 认证（API Token + Basic Auth）
- 支持 JSON 和 Markdown 输出格式
- 实现附件下载和日志清理逻辑

---

## 🏗️ 项目结构

### 目录结构

```
workflow.go/
├── cmd/
│   └── workflow/
│       └── main.go              # 主入口
├── internal/
│   ├── cli/
│   │   └── root.go              # Cobra 根命令定义
│   ├── commands/                # 命令实现
│   │   ├── setup.go
│   │   ├── update.go
│   │   ├── uninstall.go
│   │   ├── version.go
│   │   ├── config.go
│   │   ├── check.go
│   │   ├── github.go
│   │   ├── completion.go
│   │   ├── stash.go
│   │   ├── repo.go
│   │   ├── pr.go
│   │   └── jira.go
│   ├── lib/                     # 核心业务逻辑
│   │   ├── git/
│   │   │   ├── client.go        # Git 客户端封装
│   │   │   ├── branch.go        # 分支操作
│   │   │   ├── commit.go        # Commit 操作
│   │   │   └── stash.go         # Stash 操作
│   │   ├── github/
│   │   │   ├── client.go        # GitHub API 客户端
│   │   │   ├── pr.go            # PR 操作
│   │   │   └── account.go       # 账号管理
│   │   ├── jira/
│   │   │   ├── client.go        # Jira API 客户端
│   │   │   ├── ticket.go        # Ticket 操作
│   │   │   └── attachment.go    # 附件操作
│   │   ├── llm/
│   │   │   ├── interface.go     # LLM 接口定义
│   │   │   ├── openai.go        # OpenAI 实现
│   │   │   ├── deepseek.go      # DeepSeek 实现
│   │   │   └── proxy.go         # Proxy 实现
│   │   ├── http/
│   │   │   ├── client.go        # HTTP 客户端封装
│   │   │   └── retry.go         # 重试机制
│   │   └── config/
│   │       ├── manager.go       # 配置管理器
│   │       ├── validator.go     # 配置验证
│   │       └── migrator.go      # 配置迁移
│   └── utils/
│       ├── output.go            # 输出格式化
│       ├── table.go             # 表格显示
│       ├── spinner.go           # 进度指示器
│       └── file.go              # 文件操作
├── pkg/                         # 公共包（可选）
├── scripts/                     # 安装脚本
├── docs/                        # 文档
├── go.mod
├── go.sum
├── Makefile
└── README.md
```

---

## 📦 依赖清单

### 核心依赖

```go
// go.mod
module github.com/your-org/workflow

go 1.21

require (
    // CLI 框架
    github.com/spf13/cobra v1.8.0
    github.com/spf13/viper v1.18.0

    // Git 操作
    github.com/go-git/go-git/v5 v5.11.0

    // HTTP 客户端
    github.com/go-resty/resty/v2 v2.11.0

    // GitHub SDK
    github.com/google/go-github/v57 v57.0.0

    // Jira SDK
    github.com/andygrunwald/go-jira v1.16.0

    // LLM SDK
    github.com/sashabaranov/go-openai v1.20.0

    // 交互式输入
    github.com/AlecAivazis/survey/v2 v2.3.7

    // 表格显示
    github.com/olekukonko/tablewriter v0.0.5

    // 进度条
    github.com/cheggaaa/pb/v3 v3.1.4

    // 颜色输出
    github.com/fatih/color v1.16.0

    // 日志
    github.com/sirupsen/logrus v1.9.3

    // TOML 解析
    github.com/pelletier/go-toml/v2 v2.1.1
)
```

### 开发依赖

```go
require (
    // 测试框架
    github.com/stretchr/testify v1.8.4
)
```

---

## 🚀 迁移实施计划

### 阶段一：基础设施搭建（第 1-2 周）

**目标**：搭建项目基础架构和核心模块

#### 任务清单

1. **项目初始化**
   - [ ] 创建 Go 项目结构
   - [ ] 配置 `go.mod` 和依赖管理
   - [ ] 设置构建脚本（Makefile）
   - [ ] 配置 CI/CD（GitHub Actions）

2. **CLI 框架集成**
   - [ ] 集成 Cobra 框架
   - [ ] 实现根命令结构
   - [ ] 实现命令分组和帮助系统
   - [ ] 实现 Shell 补全生成（completion 命令）

3. **配置管理**
   - [ ] 集成 Viper 配置管理
   - [ ] 实现 TOML 配置文件读写
   - [ ] 实现配置验证逻辑
   - [ ] 实现配置导入/导出功能

4. **HTTP 客户端**
   - [ ] 实现统一的 HTTP 客户端封装
   - [ ] 实现重试机制（指数退避）
   - [ ] 实现认证支持（Bearer Token、Basic Auth）
   - [ ] 实现代理支持

5. **工具库**
   - [ ] 实现输出格式化工具
   - [ ] 实现表格显示工具
   - [ ] 实现进度条和 Spinner
   - [ ] 实现颜色输出

**交付物：**
- 可运行的 CLI 框架
- 配置管理模块
- HTTP 客户端封装
- 基础工具库

---

### 阶段二：核心命令实现（第 3-5 周）

**目标**：实现生命周期管理和基础命令

#### 任务清单

1. **生命周期管理命令**
   - [ ] `workflow setup` - 交互式配置初始化
   - [ ] `workflow update` - 版本检查和更新
   - [ ] `workflow uninstall` - 卸载逻辑
   - [ ] `workflow version` - 版本信息显示

2. **配置管理命令**
   - [ ] `workflow config show` - 显示配置
   - [ ] `workflow config validate` - 配置验证
   - [ ] `workflow config export` - 配置导出
   - [ ] `workflow config import` - 配置导入

3. **环境检查命令**
   - [ ] `workflow check` - 环境检查逻辑
   - [ ] Git 状态检查
   - [ ] 网络连接检查
   - [ ] 配置文件检查

4. **Shell Completion 命令**
   - [ ] `workflow completion generate` - 生成补全脚本
   - [ ] `workflow completion check` - 检查补全状态
   - [ ] `workflow completion remove` - 移除补全配置

**交付物：**
- 生命周期管理命令完整实现
- 配置管理命令完整实现
- 环境检查命令完整实现
- Shell 补全功能完整实现

---

### 阶段三：Git 操作实现（第 6-8 周）

**目标**：实现 Git 相关操作

#### 任务清单

1. **Git 客户端封装**
   - [ ] 集成 `go-git` 库
   - [ ] 实现 Git 操作统一接口
   - [ ] 实现错误处理和重试逻辑

2. **Stash 管理命令**
   - [ ] `workflow stash list` - 列出 stash
   - [ ] `workflow stash apply` - 应用 stash
   - [ ] `workflow stash drop` - 删除 stash
   - [ ] `workflow stash pop` - 应用并删除
   - [ ] `workflow stash push` - 保存到 stash

3. **仓库管理命令**
   - [ ] `workflow repo setup` - 项目级配置
   - [ ] `workflow repo show` - 显示项目配置
   - [ ] `workflow repo clean` - 清理分支和 tag

**交付物：**
- Git 操作模块完整实现
- Stash 管理命令完整实现
- 仓库管理命令完整实现

---

### 阶段四：GitHub 集成（第 9-11 周）

**目标**：实现 GitHub API 集成和 PR 操作

#### 任务清单

1. **GitHub API 客户端**
   - [ ] 集成 `google/go-github` 或自定义实现
   - [ ] 实现认证逻辑（Personal Access Token）
   - [ ] 实现错误处理和重试

2. **GitHub 账号管理命令**
   - [ ] `workflow github list` - 列出账号
   - [ ] `workflow github current` - 显示当前账号
   - [ ] `workflow github add` - 添加账号
   - [ ] `workflow github remove` - 删除账号
   - [ ] `workflow github switch` - 切换账号
   - [ ] `workflow github update` - 更新账号信息

3. **PR 操作命令**
   - [ ] `workflow pr create` - 创建 PR
   - [ ] `workflow pr merge` - 合并 PR
   - [ ] `workflow pr close` - 关闭 PR
   - [ ] `workflow pr status` - 查看 PR 状态
   - [ ] `workflow pr list` - 列出 PR
   - [ ] `workflow pr update` - 更新代码
   - [ ] `workflow pr approve` - 批准 PR
   - [ ] `workflow pr comment` - 添加评论

**交付物：**
- GitHub API 客户端完整实现
- GitHub 账号管理命令完整实现
- PR 操作命令完整实现（除 summarize 和 reword）

---

### 阶段五：LLM 集成和高级 PR 功能（第 12-13 周）

**目标**：实现 LLM 集成和高级 PR 功能

#### 任务清单

1. **LLM 集成模块**
   - [ ] 定义统一的 LLM 接口
   - [ ] 实现 OpenAI 提供者
   - [ ] 实现 DeepSeek 提供者
   - [ ] 实现 Proxy 提供者（自定义代理）
   - [ ] 实现多语言提示词生成

2. **高级 PR 功能**
   - [ ] `workflow pr summarize` - PR 总结（集成 LLM）
   - [ ] `workflow pr reword` - PR 标题和描述重写（集成 LLM）

**交付物：**
- LLM 集成模块完整实现
- 高级 PR 功能完整实现

---

### 阶段六：Jira 集成（第 14-16 周）

**目标**：实现 Jira API 集成

#### 任务清单

1. **Jira API 客户端**
   - [ ] 集成 `andygrunwald/go-jira` 或自定义实现
   - [ ] 实现认证逻辑（API Token + Basic Auth）
   - [ ] 实现错误处理和重试

2. **Jira 操作命令**
   - [ ] `workflow jira info` - 显示 ticket 信息
   - [ ] `workflow jira related` - 显示关联信息
   - [ ] `workflow jira changelog` - 显示变更历史
   - [ ] `workflow jira comment` - 添加评论
   - [ ] `workflow jira comments` - 显示评论列表
   - [ ] `workflow jira attachments` - 下载附件
   - [ ] `workflow jira clean` - 清理日志目录

**交付物：**
- Jira API 客户端完整实现
- Jira 操作命令完整实现

---

### 阶段七：测试和优化（第 17-19 周）

**目标**：完善测试、优化性能、完善文档

#### 任务清单

1. **单元测试**
   - [ ] 核心模块单元测试（覆盖率 > 80%）
   - [ ] 命令模块单元测试
   - [ ] API 客户端 Mock 测试

2. **集成测试**
   - [ ] 端到端测试场景
   - [ ] 配置文件迁移测试
   - [ ] 跨平台测试

3. **性能优化**
   - [ ] 启动速度优化（目标 < 200ms）
   - [ ] 二进制体积优化（目标 < 30MB）
   - [ ] 内存占用优化（目标 < 50MB）

4. **文档完善**
   - [ ] API 文档
   - [ ] 用户手册
   - [ ] 迁移指南
   - [ ] 开发文档

5. **发布准备**
   - [ ] 构建脚本优化
   - [ ] 安装脚本（macOS/Linux/Windows）
   - [ ] 发布流程自动化
   - [ ] 版本管理

**交付物：**
- 完整的测试套件
- 性能优化报告
- 完整的文档
- 可发布的版本

---

## 🔧 技术实现细节

### 1. Cobra 命令定义示例

```go
// internal/cli/root.go
package cli

import (
    "github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
    Use:   "workflow",
    Short: "Workflow CLI - Git 工作流自动化工具",
    Long:  `Workflow CLI 是一个功能强大的 Git 工作流自动化工具，支持 PR 管理、Jira 集成、LLM 集成等功能。`,
    Version: "1.0.0",
}

func Execute() error {
    return rootCmd.Execute()
}
```

```go
// internal/commands/setup.go
package commands

import (
    "github.com/spf13/cobra"
    "github.com/AlecAivazis/survey/v2"
)

var setupCmd = &cobra.Command{
    Use:   "setup",
    Short: "初始化或更新配置（交互式）",
    Long:  `交互式初始化或更新 Workflow CLI 配置。`,
    RunE: func(cmd *cobra.Command, args []string) error {
        return runSetup()
    },
}

func runSetup() error {
    // 使用 survey 进行交互式配置
    var config Config
    err := survey.Ask([]*survey.Question{
        {
            Name: "githubToken",
            Prompt: &survey.Input{
                Message: "请输入 GitHub Personal Access Token:",
            },
        },
        // ... 更多问题
    }, &config)

    if err != nil {
        return err
    }

    // 保存配置
    return saveConfig(config)
}
```

### 2. 配置管理实现

```go
// internal/lib/config/manager.go
package config

import (
    "github.com/spf13/viper"
    "github.com/pelletier/go-toml/v2"
)

type Manager struct {
    viper *viper.Viper
}

func NewManager() *Manager {
    v := viper.New()
    v.SetConfigName("config")
    v.SetConfigType("toml")
    v.AddConfigPath("$HOME/.workflow")
    v.AddConfigPath(".")

    return &Manager{viper: v}
}

func (m *Manager) Load() error {
    return m.viper.ReadInConfig()
}

func (m *Manager) Save(config interface{}) error {
    // 序列化为 TOML
    data, err := toml.Marshal(config)
    if err != nil {
        return err
    }

    // 写入文件
    configPath := m.viper.ConfigFileUsed()
    return os.WriteFile(configPath, data, 0644)
}
```

### 3. HTTP 客户端封装

```go
// internal/lib/http/client.go
package http

import (
    "net/http"
    "time"
    "github.com/go-resty/resty/v2"
)

type Client struct {
    client *resty.Client
}

func NewClient() *Client {
    client := resty.New()
    client.SetTimeout(30 * time.Second)
    client.SetRetryCount(3)
    client.SetRetryWaitTime(1 * time.Second)
    client.SetRetryMaxWaitTime(10 * time.Second)

    return &Client{client: client}
}

func (c *Client) SetAuth(token string) {
    c.client.SetAuthToken(token)
}

func (c *Client) Get(url string) (*resty.Response, error) {
    return c.client.R().Get(url)
}

func (c *Client) Post(url string, body interface{}) (*resty.Response, error) {
    return c.client.R().SetBody(body).Post(url)
}
```

### 4. GitHub API 客户端

```go
// internal/lib/github/client.go
package github

import (
    "context"
    "github.com/google/go-github/v57/github"
    "golang.org/x/oauth2"
)

type Client struct {
    client *github.Client
    ctx    context.Context
}

func NewClient(token string) *Client {
    ctx := context.Background()
    ts := oauth2.StaticTokenSource(
        &oauth2.Token{AccessToken: token},
    )
    tc := oauth2.NewClient(ctx, ts)

    return &Client{
        client: github.NewClient(tc),
        ctx:    ctx,
    }
}

func (c *Client) CreatePR(owner, repo string, pr *github.NewPullRequest) (*github.PullRequest, error) {
    return c.client.PullRequests.Create(c.ctx, owner, repo, pr)
}
```

### 5. LLM 集成实现

```go
// internal/lib/llm/interface.go
package llm

type Provider interface {
    GenerateText(prompt string) (string, error)
    SummarizePR(prContent string, language string) (string, error)
    RewordText(text string, instruction string) (string, error)
}

// internal/lib/llm/openai.go
package llm

import (
    "github.com/sashabaranov/go-openai"
)

type OpenAIProvider struct {
    client *openai.Client
}

func NewOpenAIProvider(apiKey string) *OpenAIProvider {
    return &OpenAIProvider{
        client: openai.NewClient(apiKey),
    }
}

func (p *OpenAIProvider) GenerateText(prompt string) (string, error) {
    resp, err := p.client.CreateChatCompletion(
        context.Background(),
        openai.ChatCompletionRequest{
            Model: openai.GPT4,
            Messages: []openai.ChatCompletionMessage{
                {Role: openai.ChatMessageRoleUser, Content: prompt},
            },
        },
    )

    if err != nil {
        return "", err
    }

    return resp.Choices[0].Message.Content, nil
}
```

---

## ✅ 测试策略

### 单元测试

```go
// internal/commands/setup_test.go
package commands

import (
    "testing"
    "github.com/stretchr/testify/assert"
)

func TestSetupCommand(t *testing.T) {
    cmd := setupCmd
    cmd.SetArgs([]string{})

    err := cmd.Execute()
    assert.NoError(t, err)
}
```

### 集成测试

```go
// internal/lib/github/client_test.go
package github

import (
    "testing"
    "github.com/stretchr/testify/assert"
)

func TestCreatePR(t *testing.T) {
    client := NewClient("test-token")
    // Mock HTTP 客户端
    // 测试 PR 创建逻辑
}
```

---

## 📊 优先级矩阵

### 高优先级（P0）- 必须实现

1. **生命周期管理**
   - `workflow setup` - 核心功能
   - `workflow version` - 基础功能
   - `workflow update` - 重要功能

2. **配置管理**
   - `workflow config show` - 基础功能
   - `workflow config validate` - 重要功能

3. **环境检查**
   - `workflow check` - 基础功能

4. **Git 操作**
   - `workflow stash list/apply/drop/pop/push` - 核心功能
   - `workflow repo clean` - 重要功能

5. **PR 操作**
   - `workflow pr create` - 核心功能
   - `workflow pr list` - 核心功能
   - `workflow pr status` - 核心功能
   - `workflow pr merge` - 核心功能
   - `workflow pr close` - 核心功能

### 中优先级（P1）- 重要功能

1. **配置管理**
   - `workflow config export/import` - 重要功能

2. **GitHub 账号管理**
   - `workflow github list/current/add/switch` - 重要功能

3. **PR 操作**
   - `workflow pr update` - 重要功能
   - `workflow pr approve` - 重要功能
   - `workflow pr comment` - 重要功能

4. **Jira 操作**
   - `workflow jira info` - 重要功能
   - `workflow jira comments` - 重要功能

### 低优先级（P2）- 增强功能

1. **生命周期管理**
   - `workflow uninstall` - 增强功能

2. **GitHub 账号管理**
   - `workflow github remove/update` - 增强功能

3. **PR 操作**
   - `workflow pr summarize` - 增强功能（需要 LLM）
   - `workflow pr reword` - 增强功能（需要 LLM）

4. **Jira 操作**
   - `workflow jira related/changelog/comment/attachments/clean` - 增强功能

5. **Shell Completion**
   - `workflow completion generate/check/remove` - 增强功能

---

## 🎯 成功标准

### 功能完整性
- [ ] 所有 P0 优先级命令实现完成
- [ ] 所有 P1 优先级命令实现完成
- [ ] 所有 P2 优先级命令实现完成（可选）

### 性能指标
- [ ] 启动速度 < 200ms
- [ ] 二进制体积 < 30MB
- [ ] 内存占用 < 50MB（运行时）

### 质量指标
- [ ] 单元测试覆盖率 > 80%
- [ ] 集成测试通过率 100%
- [ ] 代码审查通过

### 文档完整性
- [ ] API 文档完整
- [ ] 用户手册完整
- [ ] 迁移指南完整

---

## 📝 注意事项

### 1. 配置兼容性
- 确保 TOML 配置文件格式与 Rust 版本兼容
- 实现配置迁移工具，支持版本升级

### 2. API 兼容性
- 保持命令行接口与 Rust 版本一致
- 确保参数和选项名称一致

### 3. 错误处理
- 实现统一的错误处理机制
- 提供友好的错误消息

### 4. 跨平台支持
- 确保所有功能在 macOS、Linux、Windows 上正常工作
- 处理平台特定的路径和配置

### 5. 安全性
- 敏感信息（Token、密码）加密存储
- 实现安全的配置导入/导出（`--no-secrets`）

---

## 📚 参考资料

- [Cobra 官方文档](https://github.com/spf13/cobra)
- [Viper 官方文档](https://github.com/spf13/viper)
- [go-git 官方文档](https://github.com/go-git/go-git)
- [go-github 官方文档](https://github.com/google/go-github)
- [go-jira 官方文档](https://github.com/andygrunwald/go-jira)
- [go-openai 官方文档](https://github.com/sashabaranov/go-openai)

---

**最后更新**: 2025-12-28


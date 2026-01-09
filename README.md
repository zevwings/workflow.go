# Workflow CLI

Workflow CLI 是一个功能强大的 Git 工作流自动化工具，支持 PR 管理、Jira 集成、LLM 集成等功能。

## 功能特性

- 🔧 **生命周期管理**：setup、update、uninstall、version
- ⚙️ **配置管理**：配置查看、验证、导入、导出
- 🔍 **环境检查**：Git 状态和网络连接检查
- 🔐 **GitHub 账号管理**：多账号管理、切换
- 💾 **Stash 管理**：Git stash 操作
- 📦 **仓库管理**：项目级配置和清理
- 🔄 **PR 操作**：创建、合并、关闭、查询、总结
- 🎫 **Jira 集成**：Ticket 查询、评论、附件下载

## 快速开始

### 安装

```bash
# 克隆仓库
git clone https://github.com/zevwings/workflow.git
cd workflow

# 构建
make build

# 安装
make install
```

### 初始化配置

```bash
workflow setup
```

### 检查环境

```bash
workflow check
```

## 命令列表

### 生命周期管理

- `workflow setup` - 初始化或更新配置（交互式）
- `workflow update [--version VERSION]` - 更新 Workflow CLI
- `workflow uninstall` - 卸载 Workflow CLI
- `workflow version` - 显示版本信息

### 配置管理

- `workflow config show` - 查看当前配置并验证配置有效性
- `workflow config export <OUTPUT> [--section SECTION] [--no-secrets] [--toml|--json|--yaml]` - 导出配置
- `workflow config import <INPUT> [--overwrite] [--section SECTION] [--dry-run]` - 导入配置

### 环境检查

- `workflow check` - 运行环境检查（Git 状态和网络连接）

### GitHub 账号管理

- `workflow github list` - 列出所有 GitHub 账号
- `workflow github current` - 显示当前激活的账号
- `workflow github add` - 添加新的 GitHub 账号
- `workflow github remove` - 删除 GitHub 账号
- `workflow github switch` - 切换当前 GitHub 账号
- `workflow github update` - 更新 GitHub 账号信息

### Shell Completion 管理

- `workflow completion generate` - 生成 completion 脚本
- `workflow completion check` - 检查 completion 状态
- `workflow completion remove` - 移除 completion 配置

### Stash 管理

- `workflow stash list [--stat]` - 列出所有 stash
- `workflow stash apply` - 应用 stash（保留条目）
- `workflow stash drop` - 删除 stash
- `workflow stash pop` - 应用并删除 stash
- `workflow stash push` - 保存当前更改到 stash

### 仓库管理

- `workflow repo setup` - 配置项目级设置
- `workflow repo show` - 显示项目级配置
- `workflow repo clean [--dry-run]` - 清理本地分支和 tag

### PR 操作

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

### Jira 操作

- `workflow jira info [PROJ-123] [--json|--markdown]` - 显示 ticket 信息
- `workflow jira related [PROJ-123] [--json|--markdown]` - 显示关联信息
- `workflow jira changelog [PROJ-123] [--json|--markdown]` - 显示变更历史
- `workflow jira comment [PROJ-123]` - 添加评论
- `workflow jira comments [PROJ-123] [--json|--markdown] [--limit LIMIT] [--offset OFFSET] [--author AUTHOR] [--since DATE]` - 显示评论
- `workflow jira attachments [PROJ-123]` - 下载所有附件
- `workflow jira clean [PROJ-123] [--all] [--dry-run] [--list]` - 清理日志目录

## 开发

### 项目结构

```
workflow.go/
├── cmd/workflow/          # 主入口
├── internal/
│   ├── cli/               # CLI 根命令
│   ├── commands/          # 命令实现
│   ├── lib/               # 核心业务逻辑
│   │   ├── git/           # Git 操作
│   │   ├── github/        # GitHub API
│   │   ├── jira/          # Jira API
│   │   ├── llm/           # LLM 集成
│   │   ├── http/          # HTTP 客户端
│   │   └── config/        # 配置管理
│   └── utils/             # 工具函数
├── go.mod
├── go.sum
└── Makefile
```

### 构建

```bash
# 构建
make build

# 运行
make run

# 测试
make test

# 格式化代码
make fmt

# 代码检查
make lint
```

## 许可证

MIT License

## 贡献

欢迎提交 Issue 和 Pull Request！


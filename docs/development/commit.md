# 提交规范

> 本文档定义了 Workflow CLI 项目的提交规范和最佳实践，所有贡献者都应遵循这些规范。

---

## 📋 目录

- [概述](#-概述)
- [Conventional Commits](#-conventional-commits)
- [提交类型](#-提交类型)
- [提交示例](#-提交示例)
- [提交信息要求](#-提交信息要求)
- [相关文档](#-相关文档)

---

## 📋 概述

本文档定义了提交规范，使用 Conventional Commits 格式，确保提交信息清晰、一致。

### 核心原则

- **格式统一**：使用 Conventional Commits 格式
- **信息清晰**：提交信息清晰表达变更内容
- **类型明确**：提交类型明确表达变更性质

### 使用场景

- 提交代码时参考
- 代码审查时检查
- 生成 CHANGELOG 时使用

---

## Conventional Commits

使用 [Conventional Commits](https://www.conventionalcommits.org/) 格式：

```
<type>(<scope>): <subject>

<body>

<footer>
```

---

## 提交类型

- **`feat`**：新功能
- **`fix`**：修复 bug
- **`docs`**：文档更新
- **`style`**：代码格式调整（不影响功能）
- **`refactor`**：代码重构
- **`test`**：测试相关
- **`chore`**：构建过程或辅助工具的变动
- **`perf`**：性能优化
- **`ci`**：CI/CD 配置变更

---

## 提交示例

```bash
# 功能提交
feat(jira): add attachments download command

Add new command to download all attachments from a JIRA ticket.
The command supports filtering by file type and size.

Closes #123

# 修复提交
fix(pr): handle merge conflict error

Fix the issue where PR merge fails silently when there's a merge conflict.
Now the command will display a clear error message.

Fixes #456

# 文档提交
docs: update development guidelines

Add error handling best practices section.

# 重构提交
refactor(http): simplify retry logic

Extract retry logic into a separate module for better maintainability.
```

---

## 提交信息要求

- **主题行**：不超过 50 个字符，使用祈使语气
- **正文**：详细说明变更原因和方式，每行不超过 72 个字符
- **页脚**：引用相关 issue（如 `Closes #123`）

**注意**：Workflow CLI 支持通过模板系统自定义提交消息格式，包括是否使用 Conventional Commits 格式，可在 `.workflow/config.toml` 或全局配置文件中配置 `[template.commit]` 部分。

---

## 🔍 故障排除

### 问题 1：提交信息格式不正确

**症状**：提交信息不符合 Conventional Commits 格式

**解决方案**：

1. 使用正确的提交类型和格式
2. 确保主题行不超过 50 个字符
3. 使用祈使语气

### 问题 2：提交信息不够详细

**症状**：提交信息缺少必要的上下文

**解决方案**：

1. 在正文中详细说明变更原因和方式
2. 引用相关的 issue
3. 说明变更的影响范围

---

## 📚 相关文档

### 开发规范

- [Git 工作流规范](./git-workflow.md) - Git 工作流规范
- [代码审查规范](./code-review.md) - 代码审查规范

### 模板配置

提交消息模板可在 `.workflow/config.toml` 或全局配置文件中通过 `[template.commit]` 部分进行配置。

---

## ✅ 检查清单

使用本规范时，请确保：

- [ ] 提交信息符合 Conventional Commits 格式
- [ ] 提交类型正确
- [ ] 主题行不超过 50 个字符
- [ ] 正文详细说明变更内容
- [ ] 引用相关 issue（如适用）

---

**最后更新**: 2025-12-23


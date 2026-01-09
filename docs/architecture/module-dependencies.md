# 模块调用关系分析

> 本文档详细分析了 `internal/config`、`internal/git`、`internal/http`、`internal/jira`、`internal/llm`、`internal/logging`、`internal/pr`、`internal/prompt` 这几个核心模块之间的调用关系。

---

## 📋 目录

- [概述](#-概述)
- [模块依赖图](#-模块依赖图)
- [详细依赖关系](#-详细依赖关系)
- [依赖层次分析](#-依赖层次分析)
- [循环依赖检查](#-循环依赖检查)
- [设计模式应用](#-设计模式应用)

---

## 📋 概述

### 分析范围

本次分析涵盖以下 8 个核心模块：

1. **config** - 配置管理模块
2. **git** - Git 操作模块
3. **http** - HTTP 客户端模块
4. **jira** - Jira API 模块
5. **llm** - LLM 功能模块
6. **logging** - 日志系统模块
7. **pr** - PR 平台提供者模块
8. **prompt** - 用户交互模块

### 分析原则

- **依赖方向**：命令层 → 库层，库层内部可相互依赖，但避免循环依赖
- **基础模块**：不依赖其他业务模块
- **适配器模式**：通过适配器解耦模块间的直接依赖

---

## 模块依赖图

### 整体依赖关系

```
┌─────────────────────────────────────────────────────────────┐
│                      Commands 层                            │
│  (commands/, cmd/workflow/)                                 │
└────────────────────┬────────────────────────────────────────┘
                     │
                     ▼
┌─────────────────────────────────────────────────────────────┐
│                      库层 (internal/)                       │
│                                                             │
│  ┌──────────┐  ┌──────────┐  ┌──────────┐  ┌──────────┐  │
│  │  config  │  │   git    │  │   http   │  │  logging │  │
│  │ (基础)   │  │ (基础)   │  │ (基础)   │  │ (基础)   │  │
│  └──────────┘  └──────────┘  └──────────┘  └──────────┘  │
│                                                             │
│  ┌──────────┐  ┌──────────┐  ┌──────────┐  ┌──────────┐  │
│  │   jira   │  │   pr     │  │   llm    │  │  prompt  │  │
│  │ (基础)   │  │ (基础)   │  │ (依赖)   │  │ (独立)   │  │
│  └──────────┘  └──────────┘  └────┬─────┘  └──────────┘  │
│                                    │                       │
│                                    ▼                       │
│                              ┌──────────┐                 │
│                              │   http   │                 │
│                              └──────────┘                 │
└─────────────────────────────────────────────────────────────┘
```

### 详细依赖关系图

```
                    ┌─────────┐
                    │ config  │ (基础模块，无依赖)
                    └─────────┘
                         ▲
                         │ (通过 adapter)
                         │
                    ┌─────────┐
                    │   git   │ (基础模块，无依赖)
                    └─────────┘

                    ┌─────────┐
                    │  http   │ (基础模块，无依赖)
                    └─────────┘
                         ▲
                         │
                    ┌─────────┐
                    │   llm   │
                    └─────────┘

                    ┌─────────┐
                    │  jira   │ (基础模块，无依赖，使用外部 SDK)
                    └─────────┘

                    ┌─────────┐
                    │   pr    │ (基础模块，定义接口)
                    └─────────┘
                         ▲
                         │
                    ┌─────────┐
                    │ github  │ (pr/github 实现)
                    └─────────┘

                    ┌─────────┐
                    │ logging │ (基础模块，无依赖)
                    └─────────┘

                    ┌─────────┐
                    │ prompt  │ (独立模块，内部子模块相互依赖)
                    └─────────┘
```

---

## 详细依赖关系

### 1. config 模块

**职责**：配置管理，包括全局配置和仓库配置

**依赖关系**：
- **被依赖**：
  - `commands/` - 命令层使用配置
  - `adapter/config/git` - Git 适配器包装
  - `testutils/` - 测试工具
- **依赖**：
  - 无（基础模块）
  - 通过适配器模式间接使用 `git` 模块（避免直接依赖）

**关键文件**：
- `internal/config/global_manager.go` - 全局配置管理
- `internal/config/repo_manager.go` - 仓库配置管理
- `internal/config/types.go` - 配置类型定义

**设计模式**：
- 适配器模式：通过 `adapter/config/git` 解耦与 `git` 模块的直接依赖

---

### 2. git 模块

**职责**：Git 仓库操作封装

**依赖关系**：
- **被依赖**：
  - `adapter/config/git` - 配置模块的适配器
  - `testutils/` - 测试工具
  - `pr/provider/` - PR 提供者可能需要 Git 信息（未来）
- **依赖**：
  - 无（基础模块，只依赖外部 `go-git` 库）

**关键文件**：
- `internal/git/repository.go` - 仓库操作
- `internal/git/branch.go` - 分支操作
- `internal/git/adapter.go` - 适配器实现

**设计模式**：
- 适配器模式：提供 `GitAdapter` 供其他模块使用

---

### 3. http 模块

**职责**：HTTP 客户端封装，提供统一的 HTTP 请求接口

**依赖关系**：
- **被依赖**：
  - `llm/client/` - LLM 客户端使用 HTTP 客户端
  - `commands/` - 命令层直接使用（如网络检查）
  - `jira/` - 可能使用（但当前使用外部 SDK）
- **依赖**：
  - 无（基础模块，只依赖外部 `resty` 库）

**关键文件**：
- `internal/http/client.go` - HTTP 客户端实现
- `internal/http/response.go` - 响应封装
- `internal/http/config.go` - 请求配置

**设计模式**：
- 单例模式：提供 `Global()` 全局客户端
- 建造者模式：`RequestConfig` 使用链式调用

---

### 4. jira 模块

**职责**：Jira API 客户端封装

**依赖关系**：
- **被依赖**：
  - `commands/` - 命令层使用 Jira 功能
  - `testutils/` - 测试工具
- **依赖**：
  - 无（基础模块，只依赖外部 `go-jira` SDK）

**关键文件**：
- `internal/jira/client.go` - 底层客户端
- `internal/jira/jira_client.go` - 高级封装客户端
- `internal/jira/api/` - API 模块（issue, project, user）

**设计模式**：
- 分层设计：底层 `Client` + 高级 `JiraClient`
- API 模块化：按功能划分 API 子模块

---

### 5. llm 模块

**职责**：LLM 功能统一接口，包括 PR 生成、翻译等

**依赖关系**：
- **被依赖**：
  - `commands/` - 命令层使用 LLM 功能
- **依赖**：
  - `http` - LLM 客户端需要 HTTP 客户端发送请求
  - `llm/client/` - 核心客户端实现
  - `llm/pr/` - PR 相关功能
  - `llm/branch/` - 分支相关功能
  - `llm/prompt/` - Prompt 模板管理

**关键文件**：
- `internal/llm/llm.go` - 统一接口导出
- `internal/llm/client/client.go` - LLM 客户端实现
- `internal/llm/pr/client.go` - PR LLM 客户端
- `internal/llm/branch/client.go` - 分支 LLM 客户端

**依赖链**：
```
llm → http
llm/pr → llm/client → http
llm/branch → llm/client → http
llm/prompt → (独立，只管理模板)
```

**设计模式**：
- 依赖注入：LLM 客户端接收 HTTP 客户端和配置
- 单例模式：`Global()` 全局客户端
- 策略模式：通过配置区分不同 LLM 提供商

---

### 6. logging 模块

**职责**：日志系统封装

**依赖关系**：
- **被依赖**：
  - `commands/` - 命令层可能使用（但当前主要使用 prompt 输出）
  - 其他模块可能使用（但当前使用较少）
- **依赖**：
  - 无（基础模块，只依赖外部 `logrus` 库）

**关键文件**：
- `internal/logging/logger.go` - 日志实现

**设计模式**：
- 单例模式：全局 `Logger` 实例

---

### 7. pr 模块

**职责**：PR 平台提供者统一接口

**依赖关系**：
- **被依赖**：
  - `commands/` - 命令层使用 PR 功能
  - `pr/github/` - GitHub 平台实现
  - `pr/provider/` - 工厂函数
- **依赖**：
  - 无（基础模块，只定义接口和类型）

**关键文件**：
- `internal/pr/platform.go` - 平台提供者接口定义
- `internal/pr/types.go` - 类型定义
- `internal/pr/provider/factory.go` - 工厂函数
- `internal/pr/github/platform.go` - GitHub 实现

**依赖链**：
```
pr/github → pr (接口)
pr/provider → pr (接口) + pr/github (实现)
```

**设计模式**：
- 接口模式：定义 `PlatformProvider` 接口
- 工厂模式：`provider.NewPlatformProvider()` 创建实例
- 策略模式：不同平台实现统一接口

---

### 8. prompt 模块

**职责**：用户交互模块，提供各种交互式提示

**依赖关系**：
- **被依赖**：
  - `commands/` - 命令层使用用户交互
  - `llm/` - 可能使用（但当前未发现直接依赖）
- **依赖**：
  - 无（独立模块，内部子模块相互依赖）

**关键文件**：
- `internal/prompt/input.go` - 输入提示
- `internal/prompt/select.go` - 选择提示
- `internal/prompt/form.go` - 表单提示
- `internal/prompt/confirm.go` - 确认提示
- `internal/prompt/form/` - 表单子模块

**内部依赖**：
```
prompt → prompt/form
prompt/form → prompt/input
prompt/form → prompt/confirm
prompt/select → prompt/common
prompt/multiselect → prompt/common
prompt/confirm → prompt/common
prompt/common → prompt/io
```

**设计模式**：
- 建造者模式：`FormBuilder` 链式调用
- 配置注入：通过 `FormConfig` 注入格式化函数（避免循环依赖）

---

## 依赖层次分析

### 第一层：基础模块（无依赖）

这些模块不依赖其他业务模块，只依赖外部库：

1. **config** - 配置管理
2. **git** - Git 操作
3. **http** - HTTP 客户端
4. **jira** - Jira API
5. **logging** - 日志系统
6. **pr** - PR 接口定义
7. **prompt** - 用户交互

### 第二层：依赖基础模块

这些模块依赖第一层的基础模块：

1. **llm** - 依赖 `http`

### 第三层：命令层

命令层依赖所有库层模块：

- `commands/` - 依赖 `config`、`http`、`prompt`、`jira`、`llm`、`pr` 等

---

## 循环依赖检查

### ✅ 无循环依赖

经过分析，所有模块之间**不存在循环依赖**：

1. **config ↔ git**：通过适配器模式解耦，`adapter/config/git` 作为中间层
2. **prompt ↔ form**：通过配置注入避免循环依赖（`FormConfig` 注入函数）
3. **llm → http**：单向依赖，无循环

### 解耦机制

1. **适配器模式**：
   - `adapter/config/git` 解耦 `config` 和 `git` 的直接依赖
   - `git` 提供 `GitAdapter`，`config` 通过适配器使用

2. **配置注入**：
   - `prompt/form` 通过 `FormConfig` 注入函数，避免直接依赖 `prompt`

3. **接口隔离**：
   - `pr` 模块只定义接口，实现放在子模块（如 `pr/github`）

---

## 设计模式应用

### 1. 适配器模式

**应用场景**：解耦 `config` 和 `git` 模块

```go
// internal/adapter/config/git.go
// 包装 git 适配器，实现 config.GitRepository 接口
type gitAdapterWrapper struct {
    impl *git.GitAdapter
}
```

**优势**：
- 避免 `config` 直接依赖 `git`
- 保持模块独立性
- 便于测试和替换实现

### 2. 工厂模式

**应用场景**：创建 PR 平台提供者

```go
// internal/pr/provider/factory.go
func NewPlatformProvider(platform, token, owner, repo string) (pr.PlatformProvider, error)
```

**优势**：
- 统一创建接口
- 支持自动检测平台
- 便于扩展新平台

### 3. 依赖注入

**应用场景**：LLM 客户端接收 HTTP 客户端和配置

```go
// internal/llm/client/client.go
func NewClient(httpClient *http.Client, config *ProviderConfig) *LLMClient
```

**优势**：
- 解耦依赖关系
- 便于测试（可注入 Mock）
- 灵活配置

### 4. 单例模式

**应用场景**：全局客户端

```go
// internal/http/client.go
func Global() *Client

// internal/llm/client/client.go
func Global(httpClient *http.Client, config *ProviderConfig) *LLMClient
```

**优势**：
- 减少资源消耗
- 线程安全
- 统一管理

### 5. 建造者模式

**应用场景**：请求配置和表单构建

```go
// internal/http/config.go
reqConfig := http.NewRequestConfig().
    WithHeaders(headers).
    WithBody(payload).
    WithTimeout(60 * time.Second)

// internal/prompt/form/builder.go
prompt.Form().
    AddInput("name", "请输入姓名").
    AddSelect("type", "请选择类型", options).
    Run()
```

**优势**：
- 链式调用，代码简洁
- 灵活配置
- 易于扩展

---

## 总结

### 依赖特点

1. **清晰的层次结构**：基础模块 → 依赖模块 → 命令层
2. **无循环依赖**：通过适配器、配置注入等模式解耦
3. **模块独立性**：各模块职责清晰，依赖关系简单
4. **易于扩展**：接口定义清晰，便于添加新功能

### 设计优势

1. **解耦机制完善**：适配器模式、配置注入等
2. **设计模式应用合理**：工厂、单例、建造者等
3. **依赖方向清晰**：命令层 → 库层，库层内部单向依赖

### 改进建议

1. **日志使用**：当前 `logging` 模块使用较少，建议在库层统一使用
2. **错误处理**：各模块错误处理方式可以进一步统一
3. **测试覆盖**：增加模块间集成测试

---

## 相关文档

- [模块组织规范](../development/module-organization.md)
- [代码风格指南](../development/code-style.md)
- [架构一致性审查](../development/references/review-architecture-consistency.md)

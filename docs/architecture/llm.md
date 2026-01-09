# LLM 模块架构文档

## 📋 概述

LLM 模块是 Workflow CLI 的核心模块，提供统一的 LLM 功能接口，包括 PR 内容生成、PR 总结、翻译等功能。该模块专注于 LLM 客户端封装和业务逻辑处理，不涉及命令层的业务逻辑。

LLM 模块提供完整的 LLM 功能，包括 LLM 客户端管理、PR 相关功能（生成、总结、重写）、翻译功能、多语言 prompt 增强等，总代码行数约 1839+ 行。

**模块统计：**
- 代码行数：约 1839+ 行（不含测试文件）
- 主要文件：15+ 个核心文件
- 主要结构体：`LLMClient`、`PullRequestLLMClient`、`BranchLLMClient`、`ProviderConfig`、`SupportedLanguage`
- 支持功能：LLM API 调用、PR 内容生成、PR 总结、PR 重写、文件变更总结、文本翻译、多语言支持

**注意**：本模块是核心库模块，其他模块通过导入使用。通过 `LLMConfigProvider` 接口实现配置的解耦。

---

## 📁 模块架构（核心业务逻辑）

LLM 模块（`internal/llm/`）是 Workflow CLI 的核心库模块，提供统一的 LLM 功能接口。该模块专注于封装 LLM API 调用，提供简洁易用的业务接口，不涉及命令层的业务逻辑。

### 模块结构

```
internal/llm/
├── llm.go                    # 统一接口导出和构造函数（281行）
│
├── client/                    # LLM 客户端核心实现
│   ├── client.go              # LLM 客户端接口和实现（337行）
│   ├── types.go               # 类型定义（LLMRequestParams、ChatCompletionResponse等）（68行）
│   ├── provider.go            # 提供商配置（ProviderConfig）（13行）
│   └── language.go            # 语言支持（SupportedLanguage、GetLanguageRequirement）（65行）
│
├── pr/                        # PR 相关功能
│   ├── client.go              # PR LLM 客户端（569行）
│   └── types.go               # PR 相关类型定义（PullRequestContent、PullRequestReword、PullRequestSummary）（39行）
│
├── branch/                    # 分支相关功能
│   └── client.go              # 分支 LLM 客户端（翻译功能）（127行）
│
├── prompt/                    # Prompt 模板管理
│   ├── loader.go              # 模板加载器（从嵌入文件系统加载）（78行）
│   ├── branch.go              # 分支生成 prompt（8行）
│   ├── pr.go                  # PR 总结和重写 prompt（40行）
│   ├── file.go                # 文件变更总结 prompt（25行）
│   ├── translate.go           # 翻译 prompt（7行）
│   └── templates/             # Prompt 模板文件（嵌入文件系统）
│       ├── branch.md          # 分支生成模板
│       ├── pr-summary.md      # PR 总结模板
│       ├── pr-reword.md       # PR 重写模板
│       ├── file-summary.md    # 文件变更总结模板
│       └── translate.md       # 翻译模板
│
└── utils/                     # 工具函数
    ├── json.go                # JSON 处理工具（提取、修复转义问题）（137行）
    └── string.go              # 字符串处理工具（分支名清理、文件名清理）（58行）
```

**总计：约 1839+ 行代码**

### 依赖模块

- **`internal/http`**：HTTP 客户端模块
  - LLM 客户端使用 HTTP 模块发送 API 请求
  - 自动使用 `http.Global()` 获取全局 HTTP 客户端

### 模块集成

- **`internal/adapter/llm/`**：适配器层，提供便捷的构造函数
  - `NewLLMConfigProvider()` - 创建配置提供者
  - `NewPullRequestLLMClient()` - 创建 PR LLM 客户端
  - `NewBranchLLMClient()` - 创建分支 LLM 客户端
- **`cmd/`**：命令层使用 LLM 功能
  - 通过适配器层创建 LLM 客户端
  - 调用 PR 生成、总结、翻译等功能

---

## 🏗️ 架构设计

### 设计原则

1. **接口抽象**：通过 `LLMClient` 接口隐藏实现细节，提供统一的 LLM API 调用接口
2. **配置驱动**：通过 `ProviderConfig` 结构体配置不同的 LLM 提供商（OpenAI、DeepSeek、代理 API）
3. **单例模式**：使用 `Global()` 函数提供全局单例客户端，减少资源消耗，提高性能
4. **依赖注入**：通过 `LLMConfigProvider` 接口实现配置的解耦，支持从不同配置源获取配置
5. **模板管理**：使用嵌入文件系统管理 prompt 模板，支持编译时验证和运行时加载

### 核心组件

#### 1. LLMClient 接口和实现 (`client/client.go`)

**职责**：提供统一的 LLM API 调用接口，封装不同提供商的 API 差异

**主要方法**：
- `Call(params *LLMRequestParams) (string, error)` - 调用 LLM API，返回生成的文本内容

**关键特性**：
- 统一格式：所有提供商使用相同的请求格式（OpenAI Chat Completions API 标准）
- 自动重试：最多重试 3 次，处理网络错误和 5xx 错误
- 超时控制：默认 60 秒超时，适合 LLM API 的响应时间
- 错误处理：统一的错误处理和错误信息提取
- 响应解析：自动解析 OpenAI 标准格式的响应，提取消息内容

**使用场景**：
- 发送 LLM API 请求
- 获取 LLM 生成的文本内容

#### 2. PullRequestLLMClient (`pr/client.go`)

**职责**：封装所有 PR 相关的 LLM 操作，提供统一的业务接口

**主要方法**：
- `GenerateContent(commitTitle, existsBranches, gitDiff) (*PullRequestContent, error)` - 生成 PR 内容（分支名、标题、描述、scope）
- `Summarize(prTitle, prDiff) (*PullRequestSummary, error)` - 生成 PR 总结文档和文件名
- `Reword(prDiff, currentTitle) (*PullRequestReword, error)` - 重写 PR 标题和描述
- `SummarizeFileChange(filePath, fileDiff) (string, error)` - 总结单个文件变更

**关键特性**：
- 多语言支持：根据语言配置生成不同语言的 prompt 和输出
- JSON 解析：自动解析 LLM 返回的 JSON 响应，提取结构化数据
- 数据清理：自动清理分支名、文件名，确保符合规范
- 错误处理：详细的错误信息，包含上下文信息

**使用场景**：
- 根据 commit 标题生成 PR 内容
- 总结 PR 的变更内容
- 更新现有 PR 的标题和描述
- 总结单个文件的修改

#### 3. BranchLLMClient (`branch/client.go`)

**职责**：封装分支相关的 LLM 操作，主要是翻译功能

**主要方法**：
- `TranslateToEnglish(text) (string, error)` - 将文本翻译为英文

**关键特性**：
- 简单接口：专注于翻译功能，接口简洁
- 自动清理：自动清理响应中的引号和多余空白
- 错误处理：处理空响应和翻译失败的情况

**使用场景**：
- 将非英文文本（中文、俄文等）翻译为英文
- 清理和规范化分支名

#### 4. Prompt 模板管理 (`prompt/loader.go`, `prompt/*.go`)

**职责**：管理 LLM prompt 模板，支持从嵌入文件系统加载

**主要方法**：
- `LoadTemplate(name) (string, error)` - 加载模板文件
- `MustLoadTemplate(name) string` - 加载模板文件（失败时 panic）
- `ListTemplates() ([]string, error)` - 列出所有可用模板

**关键特性**：
- 嵌入文件系统：使用 `embed.FS` 将模板文件嵌入到二进制文件中
- 编译时验证：模板文件在编译时验证，确保存在
- 动态生成：支持根据语言配置动态生成 prompt（如 `GenerateSummarizePRSystemPrompt`）
- 语言增强：通过 `GetLanguageRequirement` 增强 prompt 中的语言要求

**使用场景**：
- 加载各种 LLM prompt 模板
- 根据语言配置生成定制的 prompt

#### 5. 工具函数 (`utils/json.go`, `utils/string.go`)

**职责**：提供 JSON 和字符串处理工具函数

**主要方法**：
- `ExtractAndFixJSON(response) string` - 从 markdown 代码块中提取并修复 JSON
- `SanitizeBranchName(name) string` - 清理分支名，确保只保留 ASCII 字符
- `CleanFilename(filename) string` - 清理文件名，确保只包含有效的文件名字符

**关键特性**：
- JSON 修复：自动修复 LLM 生成的 JSON 中的转义问题（如 Windows 路径中的反斜杠）
- Markdown 提取：从 markdown 代码块中提取 JSON 内容
- 字符串清理：清理分支名和文件名，确保符合规范

**使用场景**：
- 解析 LLM 返回的 JSON 响应
- 清理和规范化分支名和文件名

### 设计模式

#### 1. 单例模式

**实现**：使用 `sync.Once` 确保全局客户端单例的线程安全初始化

**优势**：
- 减少资源消耗：避免重复创建客户端实例
- 线程安全：可以在多线程环境中安全使用
- 统一管理：所有 LLM 调用使用同一个客户端实例

#### 2. 依赖注入

**实现**：通过 `LLMConfigProvider` 接口实现配置的解耦，客户端接收配置提供者而不是直接依赖配置

**优势**：
- 解耦配置：客户端不直接依赖配置模块
- 灵活扩展：支持从不同配置源获取配置
- 易于测试：可以轻松创建测试用的配置提供者

#### 3. 策略模式

**实现**：通过 `ProviderConfig` 结构体配置不同的 LLM 提供商，客户端根据配置自动适配

**优势**：
- 统一接口：所有提供商使用相同的客户端接口
- 易于扩展：添加新提供商只需配置，无需修改代码
- 配置驱动：通过配置区分不同的提供商

#### 4. 模板方法模式

**实现**：Prompt 模板定义了 LLM 调用的结构，具体的 prompt 内容通过模板文件定义

**优势**：
- 易于维护：Prompt 内容独立于代码，便于修改和优化
- 编译时验证：模板文件在编译时验证，确保存在
- 动态生成：支持根据语言配置动态生成 prompt

### 错误处理

#### 分层错误处理

1. **LLM API 调用层**：处理网络错误、HTTP 错误、超时等
   - 自动重试：最多重试 3 次
   - 错误信息提取：从 HTTP 响应中提取错误信息
   - 超时控制：60 秒超时，适合 LLM API 的响应时间

2. **响应解析层**：处理 JSON 解析错误、格式错误等
   - JSON 修复：自动修复 JSON 中的转义问题
   - Markdown 提取：从 markdown 代码块中提取 JSON
   - 字段验证：验证必需字段是否存在

3. **业务逻辑层**：处理业务相关的错误
   - 数据清理：清理和规范化数据
   - 上下文信息：在错误信息中包含上下文（如 commit title、PR title）

#### 容错机制

- **网络错误**：自动重试，最多 3 次
- **JSON 解析错误**：自动修复转义问题，从 markdown 代码块中提取
- **空响应**：检查并返回明确的错误信息
- **配置错误**：在初始化时检查配置，无效配置会导致 panic

---

## 🔄 集成关系

### 模块使用关系

LLM 模块被以下模块使用：

1. **`internal/adapter/llm/`**：适配器层，提供便捷的构造函数
   - 使用 `llm.NewPullRequestLLMClient()` - 创建 PR LLM 客户端
   - 使用 `llm.NewBranchLLMClient()` - 创建分支 LLM 客户端
   - 实现 `llm.LLMConfigProvider` 接口 - 从配置模块获取配置

2. **`cmd/`**：命令层使用 LLM 功能
   - 通过适配器层创建 LLM 客户端
   - 调用 `GenerateContent()` - 生成 PR 内容
   - 调用 `Summarize()` - 总结 PR
   - 调用 `Reword()` - 重写 PR
   - 调用 `TranslateToEnglish()` - 翻译文本

### 调用流程

#### PR 内容生成流程

```
命令层 (cmd/)
  ↓
适配器层 (adapter/llm/)
  ↓ NewPullRequestLLMClient(provider)
LLM 模块 (llm/)
  ↓ global(provider) → client.Global(providerConfig)
LLM 客户端 (llm/client/)
  ↓ Call(params)
HTTP 模块 (http/)
  ↓ PostWithConfig(url, config)
LLM API
  ↓ 返回 JSON 响应
LLM 客户端 (llm/client/)
  ↓ extractContent(response)
PR 客户端 (llm/pr/)
  ↓ parseCreateResponse(response)
  ↓ 返回 PullRequestContent
命令层 (cmd/)
```

#### PR 总结流程

```
命令层 (cmd/)
  ↓
适配器层 (adapter/llm/)
  ↓ NewPullRequestLLMClient(provider)
LLM 模块 (llm/)
  ↓ global(provider) → client.Global(providerConfig)
LLM 客户端 (llm/client/)
  ↓ Call(params)
HTTP 模块 (http/)
  ↓ PostWithConfig(url, config)
LLM API
  ↓ 返回 JSON 响应
LLM 客户端 (llm/client/)
  ↓ extractContent(response)
PR 客户端 (llm/pr/)
  ↓ parseSummaryResponse(response)
  ↓ 返回 PullRequestSummary
命令层 (cmd/)
```

#### 翻译流程

```
命令层 (cmd/)
  ↓
适配器层 (adapter/llm/)
  ↓ NewBranchLLMClient(provider)
LLM 模块 (llm/)
  ↓ global(provider) → client.Global(providerConfig)
LLM 客户端 (llm/client/)
  ↓ Call(params)
HTTP 模块 (http/)
  ↓ PostWithConfig(url, config)
LLM API
  ↓ 返回文本响应
LLM 客户端 (llm/client/)
  ↓ extractContent(response)
分支客户端 (llm/branch/)
  ↓ 清理响应（移除引号、多余空白）
  ↓ 返回翻译后的文本
命令层 (cmd/)
```

---

## 🎯 核心功能

### 1. PR 内容生成

**功能说明**：根据 commit 标题和 git diff 生成符合规范的分支名、PR 标题、描述和 scope

**流程**：
1. 构建 user prompt，包含 commit 标题、已存在分支列表和 git diff
2. 加载分支生成 system prompt 模板
3. 调用 LLM API，请求生成 JSON 格式的响应
4. 解析 JSON 响应，提取 `branch_name`、`pr_title`、`description` 和 `scope` 字段
5. 清理分支名，确保只保留 ASCII 字符
6. 返回 `PullRequestContent` 结构体

**示例**：
```go
import adapterllm "github.com/zevwings/workflow/internal/adapter/llm"

// 创建 PR LLM 客户端
prClient := adapterllm.NewPullRequestLLMClient()

// 生成 PR 内容
content, err := prClient.GenerateContent("fix: bug in authentication", nil, gitDiff)
if err != nil {
    return err
}

// 使用生成的内容
fmt.Println("Branch:", content.BranchName)
fmt.Println("Title:", content.PRTitle)
fmt.Println("Description:", *content.Description)
```

### 2. PR 总结

**功能说明**：根据 PR 的 diff 内容生成总结文档和合适的文件名

**流程**：
1. 构建 user prompt，包含 PR 标题和 PR diff
2. 根据语言配置生成 system prompt（支持多语言）
3. 调用 LLM API，请求生成 JSON 格式的响应
4. 解析 JSON 响应，提取 `summary` 和 `filename` 字段
5. 清理文件名，确保只包含有效的文件名字符
6. 返回 `PullRequestSummary` 结构体

**示例**：
```go
import adapterllm "github.com/zevwings/workflow/internal/adapter/llm"

// 创建 PR LLM 客户端
prClient := adapterllm.NewPullRequestLLMClient()

// 生成 PR 总结
summary, err := prClient.Summarize("Add user authentication", prDiff)
if err != nil {
    return err
}

// 使用生成的总结
fmt.Println("Filename:", summary.Filename)
fmt.Println("Summary:", summary.Summary)
```

### 3. PR 重写

**功能说明**：根据当前 PR 标题和 PR diff 生成更新的标题和完整的描述，用于更新现有 PR

**流程**：
1. 构建 user prompt，包含当前 PR 标题（主要输入）和 PR diff（用于验证和细化）
2. 加载 PR 重写 system prompt 模板
3. 调用 LLM API，请求生成 JSON 格式的响应
4. 解析 JSON 响应，提取 `pr_title` 和 `description` 字段
5. 返回 `PullRequestReword` 结构体

**示例**：
```go
import adapterllm "github.com/zevwings/workflow/internal/adapter/llm"

// 创建 PR LLM 客户端
prClient := adapterllm.NewPullRequestLLMClient()

// 重写 PR
currentTitle := "fix bug"
reword, err := prClient.Reword(prDiff, &currentTitle)
if err != nil {
    return err
}

// 使用重写的内容
fmt.Println("New Title:", reword.PRTitle)
fmt.Println("Description:", *reword.Description)
```

### 4. 文件变更总结

**功能说明**：根据文件的 diff 内容生成该文件的修改总结

**流程**：
1. 构建 user prompt，包含文件路径和文件 diff
2. 根据语言配置生成 system prompt（支持多语言）
3. 调用 LLM API，请求生成文本响应
4. 清理响应，移除可能的 markdown 代码块包装
5. 返回文件的修改总结（纯文本）

**示例**：
```go
import adapterllm "github.com/zevwings/workflow/internal/adapter/llm"

// 创建 PR LLM 客户端
prClient := adapterllm.NewPullRequestLLMClient()

// 总结文件变更
summary, err := prClient.SummarizeFileChange("src/auth/login.ts", fileDiff)
if err != nil {
    return err
}

// 使用生成的总结
fmt.Println("File Summary:", summary)
```

### 5. 文本翻译

**功能说明**：使用 LLM 将非英文文本（中文、俄文等）翻译为英文

**流程**：
1. 构建 user prompt，包含需要翻译的文本
2. 加载翻译 system prompt 模板
3. 调用 LLM API，请求生成文本响应
4. 清理响应，移除引号和多余空白
5. 返回翻译后的英文文本

**示例**：
```go
import adapterllm "github.com/zevwings/workflow/internal/adapter/llm"

// 创建分支 LLM 客户端
branchClient := adapterllm.NewBranchLLMClient()

// 翻译文本
translated, err := branchClient.TranslateToEnglish("你好")
if err != nil {
    return err
}

// 使用翻译后的文本
fmt.Println("Translated:", translated) // "Hello"
```

---

## 📋 使用示例

### 创建 PR LLM 客户端

```go
import (
    adapterllm "github.com/zevwings/workflow/internal/adapter/llm"
)

// 方式 1: 使用适配器层的便捷函数（推荐，最简单）
prClient := adapterllm.NewPullRequestLLMClient()

// 方式 2: 使用适配器创建 provider 并传入
provider := adapterllm.NewLLMConfigProvider()
prClient := llm.NewPullRequestLLMClient(provider)

// 方式 3: 手动实现 LLMConfigProvider 接口并传入
provider := yourCustomProvider // 实现 LLMConfigProvider 接口
prClient := llm.NewPullRequestLLMClient(provider)
```

### 生成 PR 内容

```go
import adapterllm "github.com/zevwings/workflow/internal/adapter/llm"

// 创建 PR LLM 客户端
prClient := adapterllm.NewPullRequestLLMClient()

// 生成 PR 内容
content, err := prClient.GenerateContent(
    "fix: authentication bug",
    []string{"main", "develop"}, // 已存在的分支列表
    gitDiff,                      // Git diff 内容
)
if err != nil {
    return err
}

// 使用生成的内容
fmt.Println("Branch:", content.BranchName)
fmt.Println("Title:", content.PRTitle)
if content.Description != nil {
    fmt.Println("Description:", *content.Description)
}
if content.Scope != nil {
    fmt.Println("Scope:", *content.Scope)
}
```

### 总结 PR

```go
import adapterllm "github.com/zevwings/workflow/internal/adapter/llm"

// 创建 PR LLM 客户端
prClient := adapterllm.NewPullRequestLLMClient()

// 生成 PR 总结
summary, err := prClient.Summarize("Add user authentication", prDiff)
if err != nil {
    return err
}

// 使用生成的总结
fmt.Println("Filename:", summary.Filename)
fmt.Println("Summary:", summary.Summary)
```

### 翻译文本

```go
import adapterllm "github.com/zevwings/workflow/internal/adapter/llm"

// 创建分支 LLM 客户端
branchClient := adapterllm.NewBranchLLMClient()

// 翻译文本
translated, err := branchClient.TranslateToEnglish("修复认证错误")
if err != nil {
    return err
}

// 使用翻译后的文本
fmt.Println("Translated:", translated) // "Fix authentication error"
```

---

## 📝 扩展性

### 添加新功能

1. 在相应的客户端（`pr/client.go` 或 `branch/client.go`）中添加新方法
2. 在 `prompt/` 目录中添加新的 prompt 模板文件（如果需要）
3. 在 `prompt/*.go` 中添加加载模板的代码
4. 实现业务逻辑，调用 `llmClient.Call()` 发送请求
5. 解析响应，返回结构化的数据

**示例**：
```go
// 在 pr/client.go 中添加新方法
func (c *PullRequestLLMClient) NewFeature(input string) (*NewFeatureResult, error) {
    // 构建 prompt
    userPrompt := buildNewFeaturePrompt(input)
    systemPrompt := prompt.MustLoadTemplate("new-feature.md")

    // 调用 LLM API
    params := &client.LLMRequestParams{
        SystemPrompt: systemPrompt,
        UserPrompt:   userPrompt,
        Temperature:  0.5,
    }

    response, err := c.llmClient.Call(params)
    if err != nil {
        return nil, err
    }

    // 解析响应
    result, err := parseNewFeatureResponse(response)
    if err != nil {
        return nil, err
    }

    return result, nil
}
```

### 添加新组件

1. 在 `internal/llm/` 目录下创建新的子包（如 `internal/llm/newcomponent/`）
2. 实现客户端结构体和方法
3. 在 `llm.go` 中导出类型和构造函数
4. 在适配器层添加便捷的构造函数（如果需要）

**示例**：
```go
// 在 internal/llm/newcomponent/client.go 中实现
package newcomponent

type NewComponentLLMClient struct {
    llmClient client.LLMClient
}

func Global(llmClient client.LLMClient) *NewComponentLLMClient {
    // 实现单例逻辑
}

// 在 llm.go 中导出
type NewComponentLLMClient = newcomponent.NewComponentLLMClient

func NewNewComponentLLMClient(provider LLMConfigProvider) *NewComponentLLMClient {
    llmClient, err := global(provider)
    if err != nil {
        panic(err)
    }
    return newcomponent.Global(llmClient)
}
```

### 添加新的 LLM 提供商

1. 在配置模块中添加新提供商的配置结构
2. 在适配器层实现配置提供者，支持新提供商
3. LLM 客户端会自动适配，因为所有提供商使用相同的 API 格式（OpenAI Chat Completions API 标准）

**注意**：LLM 客户端已经支持通过配置区分不同的提供商，无需修改客户端代码。

---

## 📚 相关文档

- [模块 README](../../internal/llm/README.md) - 基础使用说明
- [HTTP 模块架构文档](./http.md) - HTTP 模块设计思路和实现细节
- [模块依赖关系文档](./module-dependencies.md) - 模块之间的依赖关系

---

## ✅ 总结

LLM 模块采用清晰的依赖注入和单例模式设计：

1. **统一接口**：通过 `LLMClient` 接口提供统一的 LLM API 调用接口
2. **配置驱动**：通过 `ProviderConfig` 结构体配置不同的 LLM 提供商
3. **单例模式**：使用 `Global()` 函数提供全局单例客户端，减少资源消耗
4. **依赖注入**：通过 `LLMConfigProvider` 接口实现配置的解耦
5. **模板管理**：使用嵌入文件系统管理 prompt 模板，支持编译时验证

**设计优势**：
- ✅ 统一接口：所有提供商使用相同的客户端接口
- ✅ 易于扩展：添加新提供商只需配置，无需修改代码
- ✅ 配置解耦：通过接口实现配置的解耦，易于测试
- ✅ 资源优化：使用单例模式减少资源消耗
- ✅ 模板管理：Prompt 模板独立于代码，便于维护和优化

**当前实现状态**：
- ✅ LLM 客户端：支持 OpenAI、DeepSeek 和代理 API
- ✅ PR 内容生成：根据 commit 标题和 git diff 生成 PR 内容
- ✅ PR 总结：生成 PR 总结文档和文件名
- ✅ PR 重写：重写 PR 标题和描述
- ✅ 文件变更总结：总结单个文件的修改
- ✅ 文本翻译：将非英文文本翻译为英文
- ✅ 多语言支持：支持多语言的 prompt 增强

---

**最后更新**: 2024-12-19

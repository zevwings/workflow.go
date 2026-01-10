# LLM 模块

本模块提供了统一的 LLM 功能接口，包括 PR 内容生成、PR 总结、翻译等功能。

## 文件说明

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
│       ├── pr-summary.md       # PR 总结模板
│       ├── pr-reword.md        # PR 重写模板
│       ├── file-summary.md     # 文件变更总结模板
│       └── translate.md       # 翻译模板
│
└── utils/                     # 工具函数
    ├── json.go                # JSON 处理工具（提取、修复转义问题）（137行）
    └── string.go              # 字符串处理工具（分支名清理、文件名清理）（58行）
```

### 核心文件

- **`llm.go`**：统一接口导出和构造函数，提供 `NewPullRequestLLMClient()` 和 `NewBranchLLMClient()` 等构造函数
- **`client/client.go`**：LLM 客户端接口定义和实现，提供统一的 LLM API 调用接口
- **`client/types.go`**：类型定义，包括 `LLMRequestParams`、`ChatCompletionResponse` 等
- **`client/provider.go`**：提供商配置结构体，用于配置不同的 LLM 提供商
- **`client/language.go`**：语言支持，包括 `SupportedLanguage` 和 `GetLanguageRequirement()` 函数
- **`pr/client.go`**：PR LLM 客户端实现，提供 PR 内容生成、总结、重写等功能
- **`pr/types.go`**：PR 相关类型定义，包括 `PullRequestContent`、`PullRequestReword`、`PullRequestSummary`
- **`branch/client.go`**：分支 LLM 客户端实现，提供翻译功能
- **`prompt/loader.go`**：模板加载器，从嵌入文件系统加载 prompt 模板
- **`prompt/*.go`**：各种 prompt 模板的加载和生成函数
- **`utils/json.go`**：JSON 处理工具，包括从 markdown 代码块中提取 JSON、修复转义问题等
- **`utils/string.go`**：字符串处理工具，包括分支名清理、文件名清理等

## 快速开始

### 创建 PR LLM 客户端

```go
import infrastructurellm "github.com/zevwings/workflow/internal/infrastructure/llm"

// 使用基础设施层的便捷函数（推荐，最简单）
prClient := infrastructurellm.NewPullRequestLLMClient()

// 或者使用基础设施层创建 provider 并传入
provider := infrastructurellm.NewLLMConfigProvider()
prClient := llm.NewPullRequestLLMClient(provider)
```

### 生成 PR 内容

```go
import infrastructurellm "github.com/zevwings/workflow/internal/infrastructure/llm"

// 创建 PR LLM 客户端
prClient := infrastructurellm.NewPullRequestLLMClient()

// 生成 PR 内容
content, err := prClient.GenerateContent(
    "fix: authentication bug",
    []string{"main", "develop"}, // 已存在的分支列表（可选）
    gitDiff,                      // Git diff 内容（可选）
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
import infrastructurellm "github.com/zevwings/workflow/internal/infrastructure/llm"

// 创建 PR LLM 客户端
prClient := infrastructurellm.NewPullRequestLLMClient()

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
import infrastructurellm "github.com/zevwings/workflow/internal/infrastructure/llm"

// 创建分支 LLM 客户端
branchClient := infrastructurellm.NewBranchLLMClient()

// 翻译文本
translated, err := branchClient.TranslateToEnglish("你好")
if err != nil {
    return err
}

// 使用翻译后的文本
fmt.Println("Translated:", translated) // "Hello"
```

## 主要接口

### LLMConfigProvider

LLM 配置提供者接口，用于从外部配置源获取 LLM 提供商配置和语言配置。

- `GetProviderConfig() (*ProviderConfig, error)` - 获取提供商配置
- `GetLanguage() (*SupportedLanguage, error)` - 获取语言配置

### LLMClient

LLM 客户端接口，提供统一的 LLM API 调用接口。

- `Call(params *LLMRequestParams) (string, error)` - 调用 LLM API，返回生成的文本内容

### PullRequestLLMClient

PR LLM 客户端，封装所有 PR 相关的 LLM 操作。

- `GenerateContent(commitTitle, existsBranches, gitDiff) (*PullRequestContent, error)` - 生成 PR 内容（分支名、标题、描述、scope）
- `Summarize(prTitle, prDiff) (*PullRequestSummary, error)` - 生成 PR 总结文档和文件名
- `Reword(prDiff, currentTitle) (*PullRequestReword, error)` - 重写 PR 标题和描述
- `SummarizeFileChange(filePath, fileDiff) (string, error)` - 总结单个文件变更

### BranchLLMClient

分支 LLM 客户端，封装分支相关的 LLM 操作。

- `TranslateToEnglish(text) (string, error)` - 将文本翻译为英文

## 注意事项

1. **单例模式**：`NewPullRequestLLMClient()` 和 `NewBranchLLMClient()` 返回的是全局单例，首次调用时初始化，后续调用会复用同一个实例。首次调用时传入的参数会被保存，后续调用会忽略参数。

2. **配置验证**：如果配置无效，构造函数会 panic。建议在应用启动时进行配置验证。

3. **错误处理**：LLM API 调用可能因为网络错误、超时、API 错误等原因失败。建议在调用时进行适当的错误处理。

4. **超时控制**：LLM API 调用默认超时时间为 60 秒，适合 LLM API 的响应时间。如果需要在更短的时间内完成，可以考虑使用更短的超时时间（需要修改客户端实现）。

5. **重试机制**：LLM 客户端会自动重试失败的请求，最多重试 3 次。这有助于处理临时的网络错误。

6. **JSON 解析**：LLM 返回的 JSON 响应会自动修复转义问题，并从 markdown 代码块中提取。如果 JSON 格式不正确，会返回详细的错误信息。

7. **多语言支持**：PR 总结和文件变更总结支持多语言。语言配置通过 `LLMConfigProvider` 接口获取，如果为 nil 则使用默认英文配置。

## 依赖

- `internal/http` - HTTP 客户端模块，用于发送 LLM API 请求
- `internal/config` - 配置模块（通过适配器层间接使用）

## 相关文档

- [详细架构文档](../../docs/architecture/llm.md) - 模块设计思路和实现细节
- [HTTP 模块架构文档](../../docs/architecture/http.md) - HTTP 模块设计思路和实现细节

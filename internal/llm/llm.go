// Package llm 提供了统一的 LLM 功能接口
//
// 这个包提供了所有 LLM 相关的功能，包括：
//   - LLM 客户端：创建和管理 LLM 客户端实例
//   - PR 相关功能：生成 PR 内容、总结 PR、重写 PR 等
//   - 翻译功能：将文本翻译为英文
//   - 语言支持：多语言 prompt 增强
//
// 使用示例：
//
//	import (
//		"github.com/zevwings/workflow/internal/http"
//		"github.com/zevwings/workflow/internal/llm"
//	)
//
//	// 创建 HTTP 客户端
//	httpClient := http.Global()
//
//	// 从外部获取配置（例如从 config 包）
//	providerConfig := &llm.ProviderConfig{
//		APIKey: "your-api-key",
//		Model:  "gpt-3.5-turbo",
//		URL:    "https://api.openai.com/v1/chat/completions",
//	}
//
//	// 创建 LLM 客户端
//	llmClient := llm.NewClient(httpClient, providerConfig)
//
//	// 创建 PR LLM 客户端
//	prClient := llm.NewPullRequestLLMClient(llmClient, nil)
//	content, err := prClient.GenerateContent("fix: bug", nil, "")
//	if err != nil {
//		// 处理错误
//	}
package llm

import (
	"github.com/zevwings/workflow/internal/http"
	"github.com/zevwings/workflow/internal/llm/branch"
	"github.com/zevwings/workflow/internal/llm/client"
	"github.com/zevwings/workflow/internal/llm/pr"
	"github.com/zevwings/workflow/internal/llm/utils"
)

// ============================================================================
// 类型重新导出
// ============================================================================

// LLMClient LLM 客户端
//
// 所有 LLM 提供商使用同一个客户端实现，通过配置结构体区分不同的提供商。
type LLMClient = client.LLMClient

// LLMRequestParams LLM 请求参数
//
// 包含调用 LLM API 所需的所有参数。
type LLMRequestParams = client.LLMRequestParams

// PullRequestContent PR 内容，包含分支名、PR 标题、描述和 scope
//
// 由 LLM 生成的分支名、PR 标题、描述和 scope，用于创建 Pull Request。
type PullRequestContent = pr.PullRequestContent

// PullRequestReword PR Reword 结果，包含标题和描述
//
// 由 LLM 基于当前 PR 标题和 PR diff 生成的标题和完整描述，用于更新现有 PR。
type PullRequestReword = pr.PullRequestReword

// PullRequestSummary PR 总结结果，包含总结文档和文件名
//
// 由 LLM 生成的 PR 总结文档和对应的文件名。
type PullRequestSummary = pr.PullRequestSummary

// ProviderConfig 提供商配置
//
// 用于 LLM 客户端的基础配置，包含 API 密钥、模型名称和 URL。
// 配置应该从外部传入，例如从 config 包获取。
type ProviderConfig = client.ProviderConfig

// ============================================================================
// 客户端相关函数重新导出
// ============================================================================

// NewClient 创建新的 LLM 客户端
//
// 参数:
//   - httpClient: HTTP 客户端（不能为 nil）
//   - config: LLM 提供商配置（不能为 nil）
//
// 返回:
//   - *LLMClient: LLM 客户端实例
//
// 注意:
//
//	ProviderConfig 应该从外部传入，例如从 config 包获取。
func NewClient(httpClient *http.Client, config *ProviderConfig) *LLMClient {
	return client.NewClient(httpClient, config)
}

// Global 获取全局 LLMClient 单例
//
// 返回进程级别的 LLMClient 单例。
// 单例会在首次调用时初始化，后续调用会复用同一个实例。
//
// 参数:
//   - httpClient: HTTP 客户端（必须，不能为 nil）
//   - config: LLM 提供商配置（必须，不能为 nil）
//
// 返回:
//   - *LLMClient: LLM 客户端实例
//
// 注意:
//
//	ProviderConfig 应该从外部传入，例如从 config 包获取。
func Global(httpClient *http.Client, config *ProviderConfig) *LLMClient {
	return client.Global(httpClient, config)
}

// DefaultLLMRequestParams 返回默认的 LLM 请求参数
//
// 返回:
//   - *LLMRequestParams: 使用默认值的请求参数
func DefaultLLMRequestParams() *LLMRequestParams {
	return client.DefaultLLMRequestParams()
}

// ============================================================================
// 辅助函数重新导出
// ============================================================================

// ExtractJSONFromMarkdown 从 markdown 代码块中提取 JSON 字符串
//
// 支持以下格式：
// - ```json\n{...}\n```
// - ```\n{...}\n```
// - 纯 JSON 字符串
//
// 参数:
//   - response: 可能包含 markdown 代码块的响应字符串
//
// 返回:
//   - string: 提取的 JSON 字符串（已去除 markdown 代码块包装）
func ExtractJSONFromMarkdown(response string) string {
	return utils.ExtractJSONFromMarkdown(response)
}

// FixJSONEscapes 修复 JSON 字符串中的转义问题
//
// LLM 生成的 JSON 可能包含未转义的反斜杠（特别是在 Windows 路径中），
// 这会导致 JSON 解析失败。此函数尝试修复这些转义问题。
//
// 参数:
//   - jsonStr: 需要修复的 JSON 字符串
//
// 返回:
//   - string: 修复后的 JSON 字符串
func FixJSONEscapes(jsonStr string) string {
	return utils.FixJSONEscapes(jsonStr)
}

// ExtractAndFixJSON 从 markdown 代码块中提取并修复 JSON 字符串
//
// 这是 ExtractJSONFromMarkdown 的增强版本，会自动修复 JSON 字符串中的转义问题。
//
// 参数:
//   - response: 可能包含 markdown 代码块的响应字符串
//
// 返回:
//   - string: 提取并修复后的 JSON 字符串
func ExtractAndFixJSON(response string) string {
	return utils.ExtractAndFixJSON(response)
}

// TruncateDiff 安全截断 diff 内容
//
// 使用字符边界安全截取，避免在 UTF-8 字符中间截断。
// 会在最后一个换行符处截断，以保持 diff 格式的完整性。
//
// 参数:
//   - diff: 需要截断的 diff 内容
//   - maxLength: 最大字符数（UTF-8 rune 计数）
//
// 返回:
//   - string: 截断后的 diff 内容，如果超过长度会添加截断提示
func TruncateDiff(diff string, maxLength int) string {
	return utils.TruncateDiff(diff, maxLength)
}

// ============================================================================
// Client 构造函数
// ============================================================================

// PullRequestLLMClient PR LLM 客户端
//
// 封装所有 PR 相关的 LLM 操作，包括生成 PR 内容、总结 PR、重写 PR 和总结文件变更。
// 提供统一的接口和配置管理。
type PullRequestLLMClient = pr.PullRequestLLMClient

// NewPullRequestLLMClient 创建新的 PR LLM 客户端
//
// 参数:
//   - llmClient: LLM 客户端实例（不能为 nil）
//   - lang: 语言配置（如果为 nil，使用默认英文配置）
//
// 返回:
//   - *PullRequestLLMClient: PR LLM 客户端实例
//
// 使用示例:
//
//	prClient := llm.NewPullRequestLLMClient(llmClient, lang)
//	content, err := prClient.GenerateContent(commitTitle, branches, diff)
//	summary, err := prClient.Summarize(prTitle, prDiff)
func NewPullRequestLLMClient(llmClient *LLMClient, lang *client.SupportedLanguage) *PullRequestLLMClient {
	return pr.NewPullRequestLLMClient(llmClient, lang)
}

// BranchLLMClient 分支 LLM 客户端
//
// 封装所有分支相关的 LLM 操作，包括翻译功能。
// 提供统一的接口和配置管理。
type BranchLLMClient = branch.BranchLLMClient

// NewBranchLLMClient 创建新的分支 LLM 客户端
//
// 参数:
//   - llmClient: LLM 客户端实例（不能为 nil）
//
// 返回:
//   - *BranchLLMClient: 分支 LLM 客户端实例
//
// 使用示例:
//
//	branchClient := llm.NewBranchLLMClient(llmClient)
//	translated, err := branchClient.TranslateToEnglish("你好")
func NewBranchLLMClient(llmClient *LLMClient) *BranchLLMClient {
	return branch.NewBranchLLMClient(llmClient)
}

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
//		adapterllm "github.com/zevwings/workflow/internal/adapter/llm"
//	)
//
//	// 使用适配器层的便捷函数创建 PR LLM 客户端（内部自动创建配置提供者和 LLM 客户端）
//	prClient := adapterllm.NewPullRequestLLMClient()
//	content, err := prClient.GenerateContent("fix: bug", nil, "")
//	if err != nil {
//		// 处理错误
//	}
package llm

import (
	"fmt"

	"github.com/zevwings/workflow/internal/llm/branch"
	"github.com/zevwings/workflow/internal/llm/client"
	"github.com/zevwings/workflow/internal/llm/pr"
)

// ============================================================================
// 接口定义
// ============================================================================

// LLMConfigProvider LLM 配置提供者接口
//
// 用于从外部配置源获取 LLM 提供商配置和语言配置。
// 实现此接口的类型可以从 config.LLMConfig 或其他配置源获取配置。
//
// 接口方法返回的类型：
//   - GetProviderConfig() 返回 *ProviderConfig（即 client.ProviderConfig）
//   - GetLanguage() 返回 *SupportedLanguage（即 client.SupportedLanguage）
//
// 使用示例：
//
//	import adapterllm "github.com/zevwings/workflow/internal/adapter/llm"
//
//	// 方式 1: 使用适配器层的便捷函数（推荐，最简单）
//	prClient := adapterllm.NewPullRequestLLMClient()
//
//	// 方式 2: 使用适配器创建 provider 并传入
//	provider := adapterllm.NewLLMConfigProvider()
//	prClient := llm.NewPullRequestLLMClient(provider)
//
//	// 方式 3: 手动实现 LLMConfigProvider 接口并传入
//	provider := yourCustomProvider // 实现 LLMConfigProvider 接口
//	prClient := llm.NewPullRequestLLMClient(provider)
type LLMConfigProvider interface {
	// GetProviderConfig 获取提供商配置
	//
	// 返回:
	//   - *ProviderConfig: LLM 提供商配置（APIKey、Model、URL）
	//   - error: 如果配置无效，返回错误
	GetProviderConfig() (*ProviderConfig, error)

	// GetLanguage 获取语言配置
	//
	// 返回:
	//   - *SupportedLanguage: 语言配置，如果为 nil 表示使用默认英文配置
	//   - error: 如果配置无效，返回错误
	GetLanguage() (*SupportedLanguage, error)
}

// ============================================================================
// 类型重新导出
// ============================================================================

// ProviderConfig 提供商配置
//
// 用于 LLM 客户端的基础配置，包含 API 密钥、模型名称和 URL。
// 此类型是 client.ProviderConfig 的类型别名。
type ProviderConfig = client.ProviderConfig

// SupportedLanguage 支持的语言信息
//
// 用于 LLM prompt 生成时的语言配置。
// 此类型是 client.SupportedLanguage 的类型别名。
type SupportedLanguage = client.SupportedLanguage

// LLMClient LLM 客户端
//
// 所有 LLM 提供商使用同一个客户端实现，通过配置结构体区分不同的提供商。
// 此类型是 client.LLMClient 的类型别名。
type LLMClient = client.LLMClient

// LLMRequestParams LLM 请求参数
//
// 包含调用 LLM API 所需的所有参数。
// 此类型是 client.LLMRequestParams 的类型别名。
type LLMRequestParams = client.LLMRequestParams

// PullRequestContent PR 内容，包含分支名、PR 标题、描述和 scope
//
// 由 LLM 生成的分支名、PR 标题、描述和 scope，用于创建 Pull Request。
// 此类型是 pr.PullRequestContent 的类型别名。
type PullRequestContent = pr.PullRequestContent

// PullRequestReword PR Reword 结果，包含标题和描述
//
// 由 LLM 基于当前 PR 标题和 PR diff 生成的标题和完整描述，用于更新现有 PR。
// 此类型是 pr.PullRequestReword 的类型别名。
type PullRequestReword = pr.PullRequestReword

// PullRequestSummary PR 总结结果，包含总结文档和文件名
//
// 由 LLM 生成的 PR 总结文档和对应的文件名。
// 此类型是 pr.PullRequestSummary 的类型别名。
type PullRequestSummary = pr.PullRequestSummary

// PullRequestLLMClient PR LLM 客户端
//
// 封装所有 PR 相关的 LLM 操作，包括生成 PR 内容、总结 PR、重写 PR 和总结文件变更。
// 提供统一的接口和配置管理。
// 此类型是 pr.PullRequestLLMClient 的类型别名。
type PullRequestLLMClient = pr.PullRequestLLMClient

// BranchLLMClient 分支 LLM 客户端
//
// 封装所有分支相关的 LLM 操作，包括翻译功能。
// 提供统一的接口和配置管理。
// 此类型是 branch.BranchLLMClient 的类型别名。
type BranchLLMClient = branch.BranchLLMClient

// ============================================================================
// 内部函数
// ============================================================================

// global 获取全局 LLMClient 单例（内部函数，不导出）
//
// 从 LLMConfigProvider 接口获取配置，内部创建 HTTP 客户端和 LLM 客户端。
// 返回进程级别的 LLMClient 单例，首次调用时初始化，后续调用复用同一个实例。
//
// 参数:
//   - provider: LLM 配置提供者（不能为 nil）
//
// 返回:
//   - LLMClient: LLM 客户端实例
//   - error: 如果配置无效或创建失败，返回错误
func global(provider LLMConfigProvider) (LLMClient, error) {
	if provider == nil {
		return nil, fmt.Errorf("llm.global: LLMConfigProvider cannot be nil")
	}

	// 从接口获取提供商配置
	providerConfig, err := provider.GetProviderConfig()
	if err != nil {
		return nil, fmt.Errorf("llm.global: failed to get LLM provider configuration: %w", err)
	}

	// 创建全局 LLM 客户端单例（自动使用 http.Global()）
	llmClient := client.Global(providerConfig)

	return llmClient, nil
}

// ============================================================================
// 公开构造函数
// ============================================================================

// NewPullRequestLLMClient 创建新的 PR LLM 客户端
//
// 从 LLMConfigProvider 接口获取配置，内部自动创建 LLM 客户端和 HTTP 客户端。
// 语言配置从 provider 中获取，如果为 nil 则使用默认英文配置。
// 返回进程级别的 PullRequestLLMClient 单例，首次调用时初始化，后续调用复用同一个实例。
//
// 参数:
//   - provider: LLM 配置提供者（不能为 nil）
//
// 返回:
//   - *PullRequestLLMClient: PR LLM 客户端实例
//
// 注意:
//   - 如果配置无效，函数会 panic
//   - LLM 客户端和 HTTP 客户端在内部自动创建和管理
//   - 首次调用时传入的参数会被保存，后续调用会忽略参数
//   - 如果传入 nil，会在首次调用时 panic
//
// 优势:
//   - 减少资源消耗：避免重复创建客户端实例
//   - 线程安全：可以在多线程环境中安全使用
//   - 统一管理：所有 PR LLM 调用使用同一个客户端实例
//
// 使用示例:
//
//	import adapterllm "github.com/zevwings/workflow/internal/adapter/llm"
//
//	// 方式 1: 使用适配器层的便捷函数（推荐，最简单）
//	prClient := adapterllm.NewPullRequestLLMClient()
//	content, err := prClient.GenerateContent("fix: bug", nil, "")
//	summary, err := prClient.Summarize("PR Title", "PR Diff")
//
//	// 方式 2: 使用适配器创建 provider 并传入
//	provider := adapterllm.NewLLMConfigProvider()
//	prClient := llm.NewPullRequestLLMClient(provider)
//
//	// 方式 3: 手动实现 LLMConfigProvider 接口并传入
//	provider := yourCustomProvider // 实现 LLMConfigProvider 接口
//	prClient := llm.NewPullRequestLLMClient(provider)
func NewPullRequestLLMClient(provider LLMConfigProvider) *PullRequestLLMClient {
	if provider == nil {
		panic(fmt.Errorf("llm.NewPullRequestLLMClient: LLMConfigProvider cannot be nil"))
	}

	// 创建 LLM 客户端
	llmClient, err := global(provider)
	if err != nil {
		panic(fmt.Errorf("llm.NewPullRequestLLMClient: failed to create LLM client: %w", err))
	}

	// 获取语言配置
	lang, err := provider.GetLanguage()
	if err != nil {
		panic(fmt.Errorf("llm.NewPullRequestLLMClient: failed to get language configuration: %w", err))
	}

	// 使用 PR 包中的单例函数
	return pr.Global(llmClient, lang)
}

// NewBranchLLMClient 创建新的分支 LLM 客户端
//
// 从 LLMConfigProvider 接口获取配置，内部自动创建 LLM 客户端和 HTTP 客户端。
// 返回进程级别的 BranchLLMClient 单例，首次调用时初始化，后续调用复用同一个实例。
//
// 参数:
//   - provider: LLM 配置提供者（不能为 nil）
//
// 返回:
//   - *BranchLLMClient: 分支 LLM 客户端实例
//
// 注意:
//   - 如果配置无效，函数会 panic
//   - LLM 客户端和 HTTP 客户端在内部自动创建和管理
//   - 首次调用时传入的参数会被保存，后续调用会忽略参数
//   - 如果传入 nil，会在首次调用时 panic
//
// 优势:
//   - 减少资源消耗：避免重复创建客户端实例
//   - 线程安全：可以在多线程环境中安全使用
//   - 统一管理：所有分支 LLM 调用使用同一个客户端实例
//
// 使用示例:
//
//	import adapterllm "github.com/zevwings/workflow/internal/adapter/llm"
//
//	// 方式 1: 使用适配器层的便捷函数（推荐，最简单）
//	branchClient := adapterllm.NewBranchLLMClient()
//	translated, err := branchClient.TranslateToEnglish("你好")
//
//	// 方式 2: 使用适配器创建 provider 并传入
//	provider := adapterllm.NewLLMConfigProvider()
//	branchClient := llm.NewBranchLLMClient(provider)
//
//	// 方式 3: 手动实现 LLMConfigProvider 接口并传入
//	provider := yourCustomProvider // 实现 LLMConfigProvider 接口
//	branchClient := llm.NewBranchLLMClient(provider)
func NewBranchLLMClient(provider LLMConfigProvider) *BranchLLMClient {
	if provider == nil {
		panic(fmt.Errorf("llm.NewBranchLLMClient: LLMConfigProvider cannot be nil"))
	}

	// 创建 LLM 客户端
	llmClient, err := global(provider)
	if err != nil {
		panic(fmt.Errorf("llm.NewBranchLLMClient: failed to create LLM client: %w", err))
	}

	// 使用分支包中的单例函数
	return branch.Global(llmClient)
}

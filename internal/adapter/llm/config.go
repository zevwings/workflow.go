package llm

import (
	"fmt"

	"github.com/zevwings/workflow/internal/config"
	"github.com/zevwings/workflow/internal/llm"
)

// ============================================================================
// 类型定义
// ============================================================================

// llmConfigProvider 实现 llm.LLMConfigProvider 接口
//
// 包装 config.LLMConfig，提供接口实现。
// 将 LLM 配置的转换逻辑封装到 adapter 层，避免其他包直接依赖 config 包的具体类型。
type llmConfigProvider struct {
	llmConfig *config.LLMConfig
}

// ============================================================================
// 接口实现
// ============================================================================

// GetProviderConfig 获取提供商配置
//
// 从 config.LLMConfig 中获取当前 provider 的配置信息（APIKey、Model、URL），
// 并转换为 llm.ProviderConfig 格式返回。
//
// 返回:
//   - *llm.ProviderConfig: LLM 提供商配置（APIKey、Model、URL）
//   - error: 如果配置无效或获取失败，返回错误
func (p *llmConfigProvider) GetProviderConfig() (*llm.ProviderConfig, error) {
	apiKey, model, url, err := p.llmConfig.CurrentProvider()
	if err != nil {
		return nil, fmt.Errorf("获取 LLM provider 配置失败: %w", err)
	}

	return &llm.ProviderConfig{
		APIKey: apiKey,
		Model:  model,
		URL:    url,
	}, nil
}

// GetLanguage 获取语言配置
//
// 从 config.LLMConfig 中获取当前语言配置信息，
// 并转换为 llm.SupportedLanguage 格式返回。
// 如果语言未配置，返回默认英文配置。
//
// 返回:
//   - *llm.SupportedLanguage: 语言配置信息（Code、Name、NativeName、InstructionTemplate）
//   - error: 如果语言配置无效或获取失败，返回错误
func (p *llmConfigProvider) GetLanguage() (*llm.SupportedLanguage, error) {
	lang, err := p.llmConfig.CurrentLanguage()
	if err != nil {
		return nil, fmt.Errorf("获取语言配置失败: %w", err)
	}

	return &llm.SupportedLanguage{
		Code:                lang.Code,
		Name:                lang.Name,
		NativeName:          lang.NativeName,
		InstructionTemplate: lang.InstructionTemplate,
	}, nil
}

// ============================================================================
// 私有辅助方法
// ============================================================================

// getLLMConfig 从全局配置管理器获取 LLM 配置
//
// 内部辅助方法，用于获取全局配置管理器并返回 LLM 配置。
// 如果配置管理器初始化失败，会 panic。
//
// 返回:
//   - *config.LLMConfig: LLM 配置实例
func getLLMConfig() *config.LLMConfig {
	manager, err := config.NewGlobalManager()
	if err != nil {
		panic(fmt.Errorf("adapter/llm.getLLMConfig: failed to initialize global config manager: %w", err))
	}
	return manager.GetLLMConfig()
}

// ============================================================================
// 公开构造函数
// ============================================================================

// NewLLMConfigProvider 创建 LLM 配置提供者
//
// 从全局配置管理器获取 LLM 配置，并创建 llm.LLMConfigProvider 接口实现。
// 此函数将 config.LLMConfig 转换为 llm.LLMConfigProvider 接口实现，
// 将配置适配器的创建移到 adapter 层，避免其他包直接依赖 config 包的具体类型。
//
// 返回:
//   - llm.LLMConfigProvider: LLM 配置提供者接口实现
//
// 注意:
//   - 如果全局配置管理器初始化失败，函数会 panic
//
// 使用示例:
//
//	provider := adapterllm.NewLLMConfigProvider()
//	prClient := llm.NewPullRequestLLMClient(provider)
func NewLLMConfigProvider() llm.LLMConfigProvider {
	llmConfig := getLLMConfig()
	return &llmConfigProvider{llmConfig: llmConfig}
}

// NewBranchLLMClient 创建分支 LLM 客户端
//
// 从全局配置创建并返回分支 LLM 客户端实例。
// 内部自动创建配置提供者和 LLM 客户端，简化客户端创建流程。
//
// 返回:
//   - *llm.BranchLLMClient: 分支 LLM 客户端实例
//
// 注意:
//   - 如果配置无效或创建失败，函数会 panic
//
// 使用示例:
//
//	branchClient := adapterllm.NewBranchLLMClient()
//	translated, err := branchClient.TranslateToEnglish("你好")
func NewBranchLLMClient() *llm.BranchLLMClient {
	provider := NewLLMConfigProvider()
	return llm.NewBranchLLMClient(provider)
}

// NewPullRequestLLMClient 创建 PR LLM 客户端
//
// 从全局配置创建并返回 PR LLM 客户端实例。
// 内部自动创建配置提供者和 LLM 客户端，简化客户端创建流程。
//
// 返回:
//   - *llm.PullRequestLLMClient: PR LLM 客户端实例
//
// 注意:
//   - 如果配置无效或创建失败，函数会 panic
//
// 使用示例:
//
//	prClient := adapterllm.NewPullRequestLLMClient()
//	content, err := prClient.GenerateContent("fix: bug", nil, "")
//	summary, err := prClient.Summarize("PR Title", "PR Diff")
func NewPullRequestLLMClient() *llm.PullRequestLLMClient {
	provider := NewLLMConfigProvider()
	return llm.NewPullRequestLLMClient(provider)
}

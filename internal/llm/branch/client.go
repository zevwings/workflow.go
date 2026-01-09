package branch

import (
	"fmt"
	"strings"
	"sync"

	"github.com/zevwings/workflow/internal/llm/client"
	"github.com/zevwings/workflow/internal/llm/prompt"
)

var (
	// globalBranchClient 全局分支 LLM 客户端单例
	globalBranchClient *BranchLLMClient
	branchOnce         sync.Once
)

// BranchLLMClient 分支 LLM 客户端
//
// 封装所有分支相关的 LLM 操作，包括翻译功能。
// 提供统一的接口和配置管理。
type BranchLLMClient struct {
	llmClient client.LLMClient
}

// newBranchLLMClient 创建新的分支 LLM 客户端（内部函数，不导出）
//
// 参数:
//   - llmClient: LLM 客户端实例（不能为 nil）
//
// 返回:
//   - *BranchLLMClient: 分支 LLM 客户端实例
func newBranchLLMClient(llmClient client.LLMClient) *BranchLLMClient {
	if llmClient == nil {
		panic(fmt.Errorf("branch.newBranchLLMClient: llmClient cannot be nil"))
	}
	return &BranchLLMClient{
		llmClient: llmClient,
	}
}

// Global 获取全局 BranchLLMClient 单例
//
// 返回进程级别的 BranchLLMClient 单例。
// 单例会在首次调用时初始化，后续调用会复用同一个实例。
//
// 参数:
//   - llmClient: LLM 客户端实例（必须，不能为 nil）
//
// 返回:
//   - *BranchLLMClient: 分支 LLM 客户端实例
//
// 注意:
//   - LLM 客户端必须由调用者提供，分支模块不负责它的创建和生命周期
//   - 首次调用时传入的参数会被保存，后续调用会忽略参数
//   - 如果传入 nil，会在首次调用时 panic
//
// 优势:
//   - 减少资源消耗：避免重复创建客户端实例
//   - 线程安全：可以在多线程环境中安全使用
//   - 统一管理：所有分支 LLM 调用使用同一个客户端实例
func Global(llmClient client.LLMClient) *BranchLLMClient {
	if llmClient == nil {
		panic(fmt.Errorf("branch.Global: llmClient cannot be nil"))
	}
	branchOnce.Do(func() {
		globalBranchClient = newBranchLLMClient(llmClient)
	})
	return globalBranchClient
}

// TranslateToEnglish 翻译为英文
//
// 使用 LLM 将非英文文本（中文、俄文等）翻译为英文。
//
// 参数:
//   - text: 需要翻译的文本
//
// 返回:
//   - string: 翻译后的英文文本
//   - error: 如果 LLM API 调用失败或返回空结果，返回相应的错误信息
func (c *BranchLLMClient) TranslateToEnglish(text string) (string, error) {
	return TranslateToEnglish(text, c.llmClient)
}

// ============================================================================
// TranslateToEnglish 相关函数
// ============================================================================

// TranslateToEnglish 使用 LLM 将文本翻译为英文
//
// 使用 LLM 将非英文文本（中文、俄文等）翻译为英文。
//
// 参数:
//   - text: 需要翻译的文本
//   - llmClient: LLM 客户端实例
//
// 返回:
//   - string: 翻译后的英文文本
//   - error: 如果 LLM API 调用失败或返回空结果，返回相应的错误信息
func TranslateToEnglish(text string, llmClient client.LLMClient) (string, error) {
	userPrompt := fmt.Sprintf("Translate this text to English: %s", text)

	maxTokens := 100
	params := &client.LLMRequestParams{
		SystemPrompt: prompt.TranslateSystemPrompt,
		UserPrompt:   userPrompt,
		MaxTokens:    &maxTokens,
		Temperature:  0.3,
	}

	translated, err := llmClient.Call(params)
	if err != nil {
		return "", fmt.Errorf("调用 LLM API 翻译文本失败: %w", err)
	}

	// 清理响应（移除引号、多余空白等）
	cleaned := strings.TrimSpace(translated)
	cleaned = strings.Trim(cleaned, `"'`)

	if cleaned == "" {
		return "", fmt.Errorf("LLM 返回的翻译结果为空")
	}

	return cleaned, nil
}

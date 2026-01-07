package branch

import (
	"fmt"
	"strings"

	"github.com/zevwings/workflow/internal/llm/client"
	"github.com/zevwings/workflow/internal/llm/prompt"
)

// BranchLLMClient 分支 LLM 客户端
//
// 封装所有分支相关的 LLM 操作，包括翻译功能。
// 提供统一的接口和配置管理。
type BranchLLMClient struct {
	llmClient *client.LLMClient
}

// NewBranchLLMClient 创建新的分支 LLM 客户端
//
// 参数:
//   - llmClient: LLM 客户端实例（不能为 nil）
//
// 返回:
//   - *BranchLLMClient: 分支 LLM 客户端实例
func NewBranchLLMClient(llmClient *client.LLMClient) *BranchLLMClient {
	if llmClient == nil {
		panic("branch.NewBranchLLMClient: llmClient 不能为 nil")
	}
	return &BranchLLMClient{
		llmClient: llmClient,
	}
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
func TranslateToEnglish(text string, llmClient *client.LLMClient) (string, error) {
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

package pr

import (
	"fmt"

	"github.com/zevwings/workflow/internal/config"
	"github.com/zevwings/workflow/internal/llm"
	"github.com/zevwings/workflow/internal/llm/prompt"
)

// SummarizeFileChange 生成单个文件的修改总结
//
// 根据文件的 diff 内容生成该文件的修改总结。
//
// 参数:
//   - filePath: 文件路径
//   - fileDiff: 文件的 diff 内容
//   - cfgMgr: 配置管理器（用于语言设置）
//   - client: LLM 客户端实例
//
// 返回:
//   - string: 文件的修改总结（纯文本）
//   - error: 如果 LLM API 调用失败，返回相应的错误信息
func SummarizeFileChange(filePath, fileDiff string, cfgMgr *config.GlobalManager, client *llm.LLMClient) (string, error) {
	// 构建请求参数
	userPrompt := buildFileSummaryUserPrompt(filePath, fileDiff)
	// 根据语言生成 system prompt（语言选择逻辑在 prompt 生成函数内部处理）
	systemPrompt := prompt.GenerateSummarizeFileChangeSystemPrompt(cfgMgr)

	params := &llm.LLMRequestParams{
		SystemPrompt: systemPrompt,
		UserPrompt:   userPrompt,
		MaxTokens:    nil, // 单个文件的总结应该比较简短
		Temperature:  0.3,
	}

	// 调用 LLM API
	response, err := client.Call(params)
	if err != nil {
		return "", fmt.Errorf("调用 LLM API 总结文件修改失败 (file path: '%s'): %w", filePath, err)
	}

	// 清理响应（移除可能的 markdown 代码块包装）
	summary := cleanFileChangeSummaryResponse(response)

	return summary, nil
}

// buildFileSummaryUserPrompt 生成单个文件修改总结的 user prompt
func buildFileSummaryUserPrompt(filePath, fileDiff string) string {
	// 限制单个文件的 diff 长度，避免超过 LLM token 限制
	const maxFileDiffLength = 8000 // 单个文件的总结不需要太多上下文
	diffTrimmed := TruncateDiff(fileDiff, maxFileDiffLength)
	return fmt.Sprintf("File path: %s\n\nFile diff:\n%s", filePath, diffTrimmed)
}

// cleanFileChangeSummaryResponse 清理文件修改总结响应
//
// 移除可能的 markdown 代码块包装，返回纯文本。
func cleanFileChangeSummaryResponse(response string) string {
	return ExtractJSONFromMarkdown(response)
}

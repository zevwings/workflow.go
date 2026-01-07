package prompt

import (
	"github.com/zevwings/workflow/internal/llm/client"
)

// GenerateSummarizeFileChangeSystemPrompt 根据语言生成单个文件修改总结的 system prompt
//
// 参数:
//   - lang: 语言配置（如果为 nil，使用默认英文配置）
//
// 返回:
//   - string: 根据语言定制的 system prompt
//
// 说明:
//
//	如果 lang 为 nil，将使用默认的英文配置。
func GenerateSummarizeFileChangeSystemPrompt(lang *client.SupportedLanguage) string {
	// 从嵌入的模板文件中加载基础 prompt
	basePrompt := MustLoadTemplate("file-summary.md")

	// 使用语言增强功能
	return client.GetLanguageRequirement(basePrompt, lang)
}

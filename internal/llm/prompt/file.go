package prompt

import (
	"github.com/zevwings/workflow/internal/config"
	"github.com/zevwings/workflow/internal/llm"
)

// GenerateSummarizeFileChangeSystemPrompt 根据语言生成单个文件修改总结的 system prompt
//
// 参数:
//   - cfg: 配置管理器（用于获取语言设置）
//
// 返回:
//   - string: 根据语言定制的 system prompt
//
// 说明:
//
//	语言选择优先级：配置文件 > 默认值（"en"）
//	如果配置文件中的语言代码不在支持列表中，将使用英文作为默认语言。
func GenerateSummarizeFileChangeSystemPrompt(cfg *config.GlobalManager) string {
	// 从嵌入的模板文件中加载基础 prompt
	basePrompt := MustLoadTemplate("file-summary.md")

	// 使用 LLM 模块的语言增强功能
	return llm.GetLanguageRequirement(basePrompt, cfg)
}

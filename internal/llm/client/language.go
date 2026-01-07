package client

import (
	"fmt"
)

// SupportedLanguage 支持的语言信息
//
// 用于 LLM prompt 生成时的语言配置。
// 与 config.SupportedLanguage 结构相同，但类型不同，避免包之间的依赖。
type SupportedLanguage struct {
	// Code 语言代码（ISO 639-1 或 ISO 639-1 + ISO 3166-1，如 "en", "zh-CN"）
	Code string
	// Name 语言名称（英文）
	Name string
	// NativeName 语言名称（本地化，用于显示）
	NativeName string
	// InstructionTemplate 语言 instruction 模板
	InstructionTemplate string
}

// GetLanguageRequirement 增强 system prompt 中的语言要求
//
// 在给定的 system prompt 开头添加强化的语言要求，确保 LLM 严格按照指定语言生成内容。
//
// 参数:
//   - systemPrompt: 原始 system prompt
//   - lang: 语言配置（如果为 nil，使用默认英文配置）
//
// 返回:
//   - string: 增强后的 system prompt，包含强化的语言要求
//
// 说明:
//
//	如果 lang 为 nil，将使用默认的英文配置。
//	语言要求会被添加到 system prompt 的开头和结尾，确保 LLM 严格按照指定语言生成内容。
func GetLanguageRequirement(systemPrompt string, lang *SupportedLanguage) string {
	// 如果 lang 为 nil，使用默认英文配置
	if lang == nil {
		lang = &SupportedLanguage{
			Code:                "en",
			Name:                "English",
			NativeName:          "English",
			InstructionTemplate: "**All outputs MUST be in English only.** If the PR title or content contains non-English text (like Chinese), translate it to English in the summary.",
		}
	}

	return fmt.Sprintf(`## CRITICAL LANGUAGE REQUIREMENT

%s

**IMPORTANT REMINDER**: The entire output, including all sections, headings, content, and text MUST be written in %s only. This is a strict requirement. Do NOT use English or any other language. Every single word in the output must be in %s.

---

%s

---

## REMINDER: Language Requirement

Remember: ALL output must be in %s only. No exceptions.`,
		lang.InstructionTemplate, lang.NativeName, lang.NativeName, systemPrompt, lang.NativeName)
}

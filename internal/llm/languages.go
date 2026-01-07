package llm

import (
	"fmt"
	"strings"

	"github.com/zevwings/workflow/internal/config"
)

// SupportedLanguage 支持的语言信息
type SupportedLanguage struct {
	// Code 语言代码（ISO 639-1 或 ISO 639-1 + ISO 3166-1，如 "en", "zh-CN"）
	Code string
	// Name 语言名称（英文）
	Name string
	// NativeName 语言名称（本地化，用于显示）
	NativeName string
	// InstructionTemplate 语言 instruction 模板
	// 使用 {language_name} 作为占位符
	InstructionTemplate string
}

// SupportedLanguages 支持的语言列表
//
// 包含主流语言：英语、中文（简体/繁体）、日语、韩语、德语、法语、西班牙语等
var SupportedLanguages = []SupportedLanguage{
	{
		Code:                "en",
		Name:                "English",
		NativeName:          "English",
		InstructionTemplate: "**All outputs MUST be in English only.** If the PR title or content contains non-English text (like Chinese), translate it to English in the summary.",
	},
	{
		Code:                "zh-CN",
		Name:                "Simplified Chinese",
		NativeName:          "简体中文",
		InstructionTemplate: "**所有输出必须使用简体中文。** 如果 PR 标题或内容包含非中文文本（如英文），请在总结中翻译为中文。",
	},
	{
		Code:                "zh-TW",
		Name:                "Traditional Chinese",
		NativeName:          "繁體中文",
		InstructionTemplate: "**所有輸出必須使用繁體中文。** 如果 PR 標題或內容包含非中文文本（如英文），請在總結中翻譯為繁體中文。",
	},
	{
		Code:                "ja",
		Name:                "Japanese",
		NativeName:          "日本語",
		InstructionTemplate: "**すべての出力は日本語のみで行う必要があります。** PR タイトルまたはコンテンツに非日本語テキスト（英語など）が含まれている場合は、要約で日本語に翻訳してください。",
	},
	{
		Code:                "ko",
		Name:                "Korean",
		NativeName:          "한국어",
		InstructionTemplate: "**모든 출력은 한국어로만 작성해야 합니다.** PR 제목이나 내용에 비한국어 텍스트(예: 영어)가 포함된 경우 요약에서 한국어로 번역하세요.",
	},
	{
		Code:                "de",
		Name:                "German",
		NativeName:          "Deutsch",
		InstructionTemplate: "**Alle Ausgaben MÜSSEN ausschließlich auf Deutsch sein.** Wenn der PR-Titel oder Inhalt nicht-deutschen Text (z.B. Englisch) enthält, übersetzen Sie ihn in der Zusammenfassung ins Deutsche.",
	},
	{
		Code:                "fr",
		Name:                "French",
		NativeName:          "Français",
		InstructionTemplate: "**Toutes les sorties DOIVENT être uniquement en français.** Si le titre ou le contenu de la PR contient du texte non français (comme l'anglais), traduisez-le en français dans le résumé.",
	},
	{
		Code:                "es",
		Name:                "Spanish",
		NativeName:          "Español",
		InstructionTemplate: "**Todas las salidas DEBEN estar únicamente en español.** Si el título o el contenido de la PR contiene texto no español (como inglés), tradúzcalo al español en el resumen.",
	},
	{
		Code:                "pt",
		Name:                "Portuguese",
		NativeName:          "Português",
		InstructionTemplate: "**Todas as saídas DEVEM estar exclusivamente em português.** Se o título ou o conteúdo da PR contiver texto não português (como inglês), traduza-o para português no resumo.",
	},
	{
		Code:                "ru",
		Name:                "Russian",
		NativeName:          "Русский",
		InstructionTemplate: "**Все выходные данные ДОЛЖНЫ быть только на русском языке.** Если заголовок или содержимое PR содержит текст не на русском языке (например, английский), переведите его на русский в резюме.",
	},
}

// FindLanguage 根据语言代码查找支持的语言
//
// 参数:
//   - code: 语言代码（如 "en", "zh-CN", "zh" 等）
//
// 返回:
//   - *SupportedLanguage: 如果找到匹配的语言，返回语言信息，否则返回 nil
//
// 说明:
//
//	支持的语言代码变体：
//	- "zh" 和 "zh-CN" 都匹配简体中文
//	- "zh-TW" 匹配繁体中文
//	- 其他语言代码精确匹配
func FindLanguage(code string) *SupportedLanguage {
	codeLower := strings.ToLower(code)

	// 特殊处理：zh 和 zh-cn 都匹配简体中文
	if codeLower == "zh" || codeLower == "zh-cn" {
		for i := range SupportedLanguages {
			if SupportedLanguages[i].Code == "zh-CN" {
				return &SupportedLanguages[i]
			}
		}
	}

	// 精确匹配
	for i := range SupportedLanguages {
		if strings.ToLower(SupportedLanguages[i].Code) == codeLower {
			return &SupportedLanguages[i]
		}
	}

	return nil
}

// GetLanguageInstruction 获取语言的 instruction
//
// 参数:
//   - code: 语言代码
//
// 返回:
//   - string: 如果找到匹配的语言，返回对应的 instruction，否则返回英文的默认 instruction
func GetLanguageInstruction(code string) string {
	lang := FindLanguage(code)
	if lang != nil {
		return lang.InstructionTemplate
	}

	// 如果找不到匹配的语言，使用英文的默认 instruction
	if len(SupportedLanguages) > 0 {
		return SupportedLanguages[0].InstructionTemplate
	}

	return ""
}

// GetLanguageRequirement 增强 system prompt 中的语言要求
//
// 在给定的 system prompt 开头添加强化的语言要求，确保 LLM 严格按照指定语言生成内容。
//
// 参数:
//   - systemPrompt: 原始 system prompt
//   - cfg: 配置管理器（用于获取语言设置）
//
// 返回:
//   - string: 增强后的 system prompt，包含强化的语言要求
//
// 说明:
//
//	语言选择优先级：配置文件 > 默认值（"en"）
//	如果配置文件中的语言代码不在支持列表中，将使用英文作为默认语言。
func GetLanguageRequirement(systemPrompt string, cfg *config.GlobalManager) string {
	// 从配置文件读取语言设置
	languageCode := cfg.GetString("llm.language")
	if languageCode == "" {
		languageCode = "en"
	}

	languageInstruction := GetLanguageInstruction(languageCode)
	lang := FindLanguage(languageCode)
	languageInfo := "English"
	if lang != nil {
		languageInfo = lang.NativeName
	}

	return fmt.Sprintf(`## CRITICAL LANGUAGE REQUIREMENT

%s

**IMPORTANT REMINDER**: The entire output, including all sections, headings, content, and text MUST be written in %s only. This is a strict requirement. Do NOT use English or any other language. Every single word in the output must be in %s.

---

%s

---

## REMINDER: Language Requirement

Remember: ALL output must be in %s only. No exceptions.`,
		languageInstruction, languageInfo, languageInfo, systemPrompt, languageInfo)
}

// GetSupportedLanguageCodes 获取所有支持的语言代码列表
//
// 返回:
//   - []string: 所有支持的语言代码的切片
func GetSupportedLanguageCodes() []string {
	codes := make([]string, len(SupportedLanguages))
	for i, lang := range SupportedLanguages {
		codes[i] = lang.Code
	}
	return codes
}

// GetSupportedLanguageDisplayNames 获取所有支持的语言显示名称列表
//
// 格式："{native_name} ({name}) - {code}"
//
// 返回:
//   - []string: 格式化的语言名称列表
func GetSupportedLanguageDisplayNames() []string {
	names := make([]string, len(SupportedLanguages))
	for i, lang := range SupportedLanguages {
		names[i] = fmt.Sprintf("%s (%s) - %s", lang.NativeName, lang.Name, lang.Code)
	}
	return names
}

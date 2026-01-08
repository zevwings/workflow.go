package client

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// ==================== GetLanguageRequirement 测试 ====================

func TestGetLanguageRequirement(t *testing.T) {
	systemPrompt := "You are a helpful assistant."
	lang := &SupportedLanguage{
		Code:                "zh-CN",
		Name:                "Chinese",
		NativeName:          "中文",
		InstructionTemplate: "**所有输出必须仅使用中文。**",
	}

	result := GetLanguageRequirement(systemPrompt, lang)

	// 验证包含原始 system prompt
	assert.Contains(t, result, systemPrompt)

	// 验证包含语言要求
	assert.Contains(t, result, "CRITICAL LANGUAGE REQUIREMENT")
	assert.Contains(t, result, "REMINDER: Language Requirement")
	assert.Contains(t, result, lang.InstructionTemplate)
	assert.Contains(t, result, lang.NativeName)
}

func TestGetLanguageRequirement_NilLang(t *testing.T) {
	systemPrompt := "You are a helpful assistant."

	result := GetLanguageRequirement(systemPrompt, nil)

	// 应该使用默认英文配置
	assert.Contains(t, result, systemPrompt)
	assert.Contains(t, result, "CRITICAL LANGUAGE REQUIREMENT")
	assert.Contains(t, result, "REMINDER: Language Requirement")
	assert.Contains(t, result, "English")
	assert.Contains(t, result, "All outputs MUST be in English only")
}

func TestGetLanguageRequirement_EmptySystemPrompt(t *testing.T) {
	lang := &SupportedLanguage{
		Code:                "zh-CN",
		Name:                "Chinese",
		NativeName:          "中文",
		InstructionTemplate: "**所有输出必须仅使用中文。**",
	}

	result := GetLanguageRequirement("", lang)

	// 应该仍然包含语言要求
	assert.Contains(t, result, "CRITICAL LANGUAGE REQUIREMENT")
	assert.Contains(t, result, "REMINDER: Language Requirement")
	assert.Contains(t, result, lang.InstructionTemplate)
	assert.Contains(t, result, lang.NativeName)
}

func TestGetLanguageRequirement_MultipleLanguages(t *testing.T) {
	tests := []struct {
		name string
		lang *SupportedLanguage
	}{
		{
			name: "中文",
			lang: &SupportedLanguage{
				Code:                "zh-CN",
				Name:                "Chinese",
				NativeName:          "中文",
				InstructionTemplate: "**所有输出必须仅使用中文。**",
			},
		},
		{
			name: "日文",
			lang: &SupportedLanguage{
				Code:                "ja",
				Name:                "Japanese",
				NativeName:          "日本語",
				InstructionTemplate: "**すべての出力は日本語のみで行う必要があります。**",
			},
		},
		{
			name: "法文",
			lang: &SupportedLanguage{
				Code:                "fr",
				Name:                "French",
				NativeName:          "Français",
				InstructionTemplate: "**Toutes les sorties doivent être en français uniquement.**",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			systemPrompt := "You are a helpful assistant."
			result := GetLanguageRequirement(systemPrompt, tt.lang)

			assert.Contains(t, result, systemPrompt)
			assert.Contains(t, result, tt.lang.NativeName)
			assert.Contains(t, result, tt.lang.InstructionTemplate)
		})
	}
}


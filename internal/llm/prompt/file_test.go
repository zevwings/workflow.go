package prompt

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zevwings/workflow/internal/llm/client"
)

// ==================== GenerateSummarizeFileChangeSystemPrompt 测试 ====================

func TestGenerateSummarizeFileChangeSystemPrompt(t *testing.T) {
	tests := []struct {
		name string
		lang *client.SupportedLanguage
	}{
		{
			name: "nil 语言配置（使用默认英文）",
			lang: nil,
		},
		{
			name: "中文配置",
			lang: &client.SupportedLanguage{
				Code:                "zh-CN",
				Name:                "Chinese",
				NativeName:          "中文",
				InstructionTemplate: "**所有输出必须仅使用中文。**",
			},
		},
		{
			name: "英文配置",
			lang: &client.SupportedLanguage{
				Code:                "en",
				Name:                "English",
				NativeName:          "English",
				InstructionTemplate: "**All outputs MUST be in English only.**",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Act: 生成文件总结 prompt
			prompt := GenerateSummarizeFileChangeSystemPrompt(tt.lang)

			// Assert: 验证 prompt 不为空
			assert.NotEmpty(t, prompt, "prompt 不应为空")
			assert.Greater(t, len(prompt), 0, "prompt 长度应大于 0")
		})
	}
}

func TestGenerateSummarizeFileChangeSystemPrompt_ContainsBaseTemplate(t *testing.T) {
	// Act: 生成文件总结 prompt（使用默认语言）
	prompt := GenerateSummarizeFileChangeSystemPrompt(nil)

	// Assert: 验证包含基础模板内容
	// 由于基础模板是从 file-summary.md 加载的，我们应该验证 prompt 包含一些预期的内容
	assert.NotEmpty(t, prompt, "prompt 不应为空")
	assert.Greater(t, len(prompt), 50, "prompt 应该有足够的长度")
}

func TestGenerateSummarizeFileChangeSystemPrompt_LanguageRequirement(t *testing.T) {
	tests := []struct {
		name     string
		lang     *client.SupportedLanguage
		expected string
	}{
		{
			name: "中文配置应该包含中文要求",
			lang: &client.SupportedLanguage{
				Code:                "zh-CN",
				Name:                "Chinese",
				NativeName:          "中文",
				InstructionTemplate: "**所有输出必须仅使用中文。**",
			},
			expected: "中文",
		},
		{
			name: "英文配置应该包含英文要求",
			lang: &client.SupportedLanguage{
				Code:                "en",
				Name:                "English",
				NativeName:          "English",
				InstructionTemplate: "**All outputs MUST be in English only.**",
			},
			expected: "English",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Act: 生成文件总结 prompt
			prompt := GenerateSummarizeFileChangeSystemPrompt(tt.lang)

			// Assert: 验证包含语言要求
			assert.Contains(t, prompt, tt.expected, "应该包含语言要求: %s", tt.expected)
			assert.Contains(t, prompt, "CRITICAL LANGUAGE REQUIREMENT", "应该包含关键语言要求")
			assert.Contains(t, prompt, "REMINDER: Language Requirement", "应该包含语言要求提醒")
		})
	}
}

func TestGenerateSummarizeFileChangeSystemPrompt_NilLanguageUsesDefault(t *testing.T) {
	// Act: 生成文件总结 prompt（nil 语言）
	prompt := GenerateSummarizeFileChangeSystemPrompt(nil)

	// Assert: 验证使用默认英文配置
	assert.Contains(t, prompt, "English", "nil 语言应该使用默认英文配置")
	assert.Contains(t, prompt, "All outputs MUST be in English only", "应该包含默认英文指令")
}

func TestGenerateSummarizeFileChangeSystemPrompt_Formatting(t *testing.T) {
	// Act: 生成文件总结 prompt
	prompt := GenerateSummarizeFileChangeSystemPrompt(nil)

	// Assert: 验证格式正确
	// 验证包含必要的分隔符和结构
	assert.Contains(t, prompt, "---", "应该包含分隔符")

	// 验证 prompt 结构完整（包含开头和结尾）
	lines := strings.Split(prompt, "\n")
	assert.Greater(t, len(lines), 10, "prompt 应该有足够的行数")
}

func TestGenerateSummarizeFileChangeSystemPrompt_ConsistentOutput(t *testing.T) {
	// Act: 多次生成 prompt（相同配置）
	prompt1 := GenerateSummarizeFileChangeSystemPrompt(nil)
	prompt2 := GenerateSummarizeFileChangeSystemPrompt(nil)

	// Assert: 验证输出一致
	assert.Equal(t, prompt1, prompt2, "相同配置应该生成相同的 prompt")
}

func TestGenerateSummarizeFileChangeSystemPrompt_DifferentLanguages(t *testing.T) {
	// Arrange: 不同语言配置
	lang1 := &client.SupportedLanguage{
		Code:                "zh-CN",
		Name:                "Chinese",
		NativeName:          "中文",
		InstructionTemplate: "**所有输出必须仅使用中文。**",
	}
	lang2 := &client.SupportedLanguage{
		Code:                "en",
		Name:                "English",
		NativeName:          "English",
		InstructionTemplate: "**All outputs MUST be in English only.**",
	}

	// Act: 生成不同语言的 prompt
	prompt1 := GenerateSummarizeFileChangeSystemPrompt(lang1)
	prompt2 := GenerateSummarizeFileChangeSystemPrompt(lang2)

	// Assert: 验证输出不同
	assert.NotEqual(t, prompt1, prompt2, "不同语言配置应该生成不同的 prompt")

	// 验证各自包含正确的语言要求
	assert.Contains(t, prompt1, "中文", "中文配置应该包含中文要求")
	assert.Contains(t, prompt2, "English", "英文配置应该包含英文要求")
}


package prompt

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zevwings/workflow/internal/llm/client"
)

// ==================== GenerateSummarizePRSystemPrompt 测试 ====================

func TestGenerateSummarizePRSystemPrompt(t *testing.T) {
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
		{
			name: "日文配置",
			lang: &client.SupportedLanguage{
				Code:                "ja",
				Name:                "Japanese",
				NativeName:          "日本語",
				InstructionTemplate: "**すべての出力は日本語のみで行う必要があります。**",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Act: 生成 PR 总结 prompt
			prompt := GenerateSummarizePRSystemPrompt(tt.lang)

			// Assert: 验证 prompt 不为空
			assert.NotEmpty(t, prompt, "prompt 不应为空")

			// 验证包含预期的内容
			assert.Contains(t, prompt, "summary", "应该包含 summary 相关内容")
			assert.Contains(t, prompt, "filename", "应该包含 filename 相关内容")

			// 验证包含 JSON 响应示例
			assert.Contains(t, prompt, "summary", "应该包含 JSON 响应示例")
			assert.Contains(t, prompt, "filename", "应该包含 filename 字段")
		})
	}
}

func TestGenerateSummarizePRSystemPrompt_ContainsBaseTemplate(t *testing.T) {
	// Act: 生成 PR 总结 prompt（使用默认语言）
	prompt := GenerateSummarizePRSystemPrompt(nil)

	// Assert: 验证包含基础模板内容
	// 由于基础模板是从 pr-summary.md 加载的，我们应该验证 prompt 包含一些预期的内容
	assert.NotEmpty(t, prompt, "prompt 不应为空")
	assert.Greater(t, len(prompt), 100, "prompt 应该有足够的长度")
}

func TestGenerateSummarizePRSystemPrompt_ContainsJSONExample(t *testing.T) {
	// Act: 生成 PR 总结 prompt
	prompt := GenerateSummarizePRSystemPrompt(nil)

	// Assert: 验证包含 JSON 响应示例
	// 注意：prompt 中包含的是转义后的 JSON 示例，所以需要检查转义版本
	assert.Contains(t, prompt, `summary`, "应该包含 summary JSON 字段")
	assert.Contains(t, prompt, `filename`, "应该包含 filename JSON 字段")
	assert.Contains(t, prompt, "{", "应该包含 JSON 结构")
}

func TestGenerateSummarizePRSystemPrompt_LanguageRequirement(t *testing.T) {
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
			// Act: 生成 PR 总结 prompt
			prompt := GenerateSummarizePRSystemPrompt(tt.lang)

			// Assert: 验证包含语言要求
			assert.Contains(t, prompt, tt.expected, "应该包含语言要求: %s", tt.expected)
			assert.Contains(t, prompt, "CRITICAL LANGUAGE REQUIREMENT", "应该包含关键语言要求")
			assert.Contains(t, prompt, "REMINDER: Language Requirement", "应该包含语言要求提醒")
		})
	}
}

func TestGenerateSummarizePRSystemPrompt_NilLanguageUsesDefault(t *testing.T) {
	// Act: 生成 PR 总结 prompt（nil 语言）
	prompt := GenerateSummarizePRSystemPrompt(nil)

	// Assert: 验证使用默认英文配置
	assert.Contains(t, prompt, "English", "nil 语言应该使用默认英文配置")
	assert.Contains(t, prompt, "All outputs MUST be in English only", "应该包含默认英文指令")
}

func TestGenerateSummarizePRSystemPrompt_Formatting(t *testing.T) {
	// Act: 生成 PR 总结 prompt
	prompt := GenerateSummarizePRSystemPrompt(nil)

	// Assert: 验证格式正确
	// 验证包含必要的分隔符和结构
	assert.Contains(t, prompt, "---", "应该包含分隔符")

	// 验证 prompt 结构完整（包含开头和结尾）
	lines := strings.Split(prompt, "\n")
	assert.Greater(t, len(lines), 10, "prompt 应该有足够的行数")
}

func TestGenerateSummarizePRSystemPrompt_ConsistentOutput(t *testing.T) {
	// Act: 多次生成 prompt（相同配置）
	prompt1 := GenerateSummarizePRSystemPrompt(nil)
	prompt2 := GenerateSummarizePRSystemPrompt(nil)

	// Assert: 验证输出一致
	assert.Equal(t, prompt1, prompt2, "相同配置应该生成相同的 prompt")
}

func TestGenerateSummarizePRSystemPrompt_DifferentLanguages(t *testing.T) {
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
	prompt1 := GenerateSummarizePRSystemPrompt(lang1)
	prompt2 := GenerateSummarizePRSystemPrompt(lang2)

	// Assert: 验证输出不同
	assert.NotEqual(t, prompt1, prompt2, "不同语言配置应该生成不同的 prompt")

	// 验证各自包含正确的语言要求
	assert.Contains(t, prompt1, "中文", "中文配置应该包含中文要求")
	assert.Contains(t, prompt2, "English", "英文配置应该包含英文要求")
}

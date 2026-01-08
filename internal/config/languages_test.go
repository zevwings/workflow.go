package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// ==================== FindLanguage 测试 ====================

func TestFindLanguage(t *testing.T) {
	tests := []struct {
		name     string
		code     string
		expected *SupportedLanguage
	}{
		{
			name:     "精确匹配 - 英语",
			code:     "en",
			expected: &SupportedLanguages[0],
		},
		{
			name:     "精确匹配 - 简体中文",
			code:     "zh-CN",
			expected: &SupportedLanguages[1],
		},
		{
			name:     "精确匹配 - 繁体中文",
			code:     "zh-TW",
			expected: &SupportedLanguages[2],
		},
		{
			name:     "精确匹配 - 日语",
			code:     "ja",
			expected: &SupportedLanguages[3],
		},
		{
			name:     "精确匹配 - 韩语",
			code:     "ko",
			expected: &SupportedLanguages[4],
		},
		{
			name:     "精确匹配 - 德语",
			code:     "de",
			expected: &SupportedLanguages[5],
		},
		{
			name:     "精确匹配 - 法语",
			code:     "fr",
			expected: &SupportedLanguages[6],
		},
		{
			name:     "精确匹配 - 西班牙语",
			code:     "es",
			expected: &SupportedLanguages[7],
		},
		{
			name:     "精确匹配 - 葡萄牙语",
			code:     "pt",
			expected: &SupportedLanguages[8],
		},
		{
			name:     "精确匹配 - 俄语",
			code:     "ru",
			expected: &SupportedLanguages[9],
		},
		{
			name:     "zh 匹配 zh-CN",
			code:     "zh",
			expected: &SupportedLanguages[1], // zh-CN
		},
		{
			name:     "zh-cn 匹配 zh-CN（大小写不敏感）",
			code:     "zh-cn",
			expected: &SupportedLanguages[1], // zh-CN
		},
		{
			name:     "ZH-CN 匹配 zh-CN（大小写不敏感）",
			code:     "ZH-CN",
			expected: &SupportedLanguages[1], // zh-CN
		},
		{
			name:     "EN 匹配 en（大小写不敏感）",
			code:     "EN",
			expected: &SupportedLanguages[0],
		},
		{
			name:     "未找到语言",
			code:     "xx",
			expected: nil,
		},
		{
			name:     "空字符串",
			code:     "",
			expected: nil,
		},
		{
			name:     "无效代码",
			code:     "invalid-code",
			expected: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Act: 查找语言
			result := FindLanguage(tt.code)

			// Assert: 验证结果
			if tt.expected == nil {
				assert.Nil(t, result)
			} else {
				requireNotNil(t, result)
				assert.Equal(t, tt.expected.Code, result.Code)
				assert.Equal(t, tt.expected.Name, result.Name)
				assert.Equal(t, tt.expected.NativeName, result.NativeName)
			}
		})
	}
}

// ==================== GetLanguageInstruction 测试 ====================

func TestGetLanguageInstruction(t *testing.T) {
	tests := []struct {
		name     string
		code     string
		expected string
	}{
		{
			name:     "英语",
			code:     "en",
			expected: SupportedLanguages[0].InstructionTemplate,
		},
		{
			name:     "简体中文",
			code:     "zh-CN",
			expected: SupportedLanguages[1].InstructionTemplate,
		},
		{
			name:     "zh 匹配 zh-CN",
			code:     "zh",
			expected: SupportedLanguages[1].InstructionTemplate,
		},
		{
			name:     "未找到语言时返回英文默认值",
			code:     "xx",
			expected: SupportedLanguages[0].InstructionTemplate,
		},
		{
			name:     "空字符串返回英文默认值",
			code:     "",
			expected: SupportedLanguages[0].InstructionTemplate,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Act: 获取语言指令
			result := GetLanguageInstruction(tt.code)

			// Assert: 验证结果
			assert.Equal(t, tt.expected, result)
			assert.NotEmpty(t, result)
		})
	}
}

// ==================== GetLanguageRequirement 测试 ====================

func TestGetLanguageRequirement(t *testing.T) {
	tests := []struct {
		name         string
		systemPrompt string
		languageCode string
		shouldContain []string
	}{
		{
			name:         "英语",
			systemPrompt: "Original prompt",
			languageCode: "en",
			shouldContain: []string{
				"CRITICAL LANGUAGE REQUIREMENT",
				"English",
				"Original prompt",
				"REMINDER: Language Requirement",
			},
		},
		{
			name:         "简体中文",
			systemPrompt: "原始提示",
			languageCode: "zh-CN",
			shouldContain: []string{
				"CRITICAL LANGUAGE REQUIREMENT",
				"简体中文",
				"原始提示",
				"REMINDER: Language Requirement",
			},
		},
		{
			name:         "zh 匹配 zh-CN",
			systemPrompt: "Test prompt",
			languageCode: "zh",
			shouldContain: []string{
				"简体中文",
				"Test prompt",
			},
		},
		{
			name:         "空语言代码使用默认值 en",
			systemPrompt: "Test prompt",
			languageCode: "",
			shouldContain: []string{
				"English",
				"Test prompt",
			},
		},
		{
			name:         "未找到语言使用默认值 en",
			systemPrompt: "Test prompt",
			languageCode: "xx",
			shouldContain: []string{
				"English",
				"Test prompt",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Act: 获取语言要求
			result := GetLanguageRequirement(tt.systemPrompt, tt.languageCode)

			// Assert: 验证结果
			assert.NotEmpty(t, result)
			assert.Contains(t, result, tt.systemPrompt, "应该包含原始 system prompt")

			// 验证包含所有必需的字符串
			for _, str := range tt.shouldContain {
				assert.Contains(t, result, str, "应该包含: %s", str)
			}

			// 验证格式：应该包含三个分隔符（---）
			separatorCount := countSubstring(result, "---")
			assert.GreaterOrEqual(t, separatorCount, 2, "应该包含至少两个分隔符")
		})
	}
}

// ==================== GetSupportedLanguageCodes 测试 ====================

func TestGetSupportedLanguageCodes(t *testing.T) {
	// Act: 获取支持的语言代码列表
	codes := GetSupportedLanguageCodes()

	// Assert: 验证结果
	assert.NotEmpty(t, codes)
	assert.Equal(t, len(SupportedLanguages), len(codes), "代码数量应该等于支持的语言数量")

	// 验证包含所有支持的语言代码
	expectedCodes := make(map[string]bool)
	for _, lang := range SupportedLanguages {
		expectedCodes[lang.Code] = true
	}

	for _, code := range codes {
		assert.True(t, expectedCodes[code], "代码 %s 应该在支持列表中", code)
	}

	// 验证顺序一致
	for i, code := range codes {
		assert.Equal(t, SupportedLanguages[i].Code, code, "代码顺序应该与 SupportedLanguages 一致")
	}
}

// ==================== GetSupportedLanguageDisplayNames 测试 ====================

func TestGetSupportedLanguageDisplayNames(t *testing.T) {
	// Act: 获取支持的语言显示名称列表
	names := GetSupportedLanguageDisplayNames()

	// Assert: 验证结果
	assert.NotEmpty(t, names)
	assert.Equal(t, len(SupportedLanguages), len(names), "名称数量应该等于支持的语言数量")

	// 验证格式：每个名称应该包含 native name, name 和 code
	for i, name := range names {
		lang := SupportedLanguages[i]
		assert.Contains(t, name, lang.NativeName, "应该包含本地名称")
		assert.Contains(t, name, lang.Name, "应该包含英文名称")
		assert.Contains(t, name, lang.Code, "应该包含语言代码")

		// 验证格式："{native_name} ({name}) - {code}"
		expectedFormat := lang.NativeName + " (" + lang.Name + ") - " + lang.Code
		assert.Equal(t, expectedFormat, name, "格式应该匹配")
	}
}

// ==================== 辅助函数 ====================

// requireNotNil 是一个辅助函数，用于验证结果不为 nil
func requireNotNil(t *testing.T, result *SupportedLanguage) {
	if result == nil {
		t.Fatal("Expected result to be not nil")
	}
}

// countSubstring 计算子字符串在字符串中出现的次数
func countSubstring(s, substr string) int {
	if len(substr) == 0 {
		return 0
	}
	count := 0
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			count++
		}
	}
	return count
}


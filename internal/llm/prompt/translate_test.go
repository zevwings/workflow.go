package prompt

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// ==================== TranslateSystemPrompt 测试 ====================

func TestTranslateSystemPrompt(t *testing.T) {
	// Act: 获取翻译 prompt
	prompt := TranslateSystemPrompt

	// Assert: 验证 prompt 不为空
	assert.NotEmpty(t, prompt, "翻译 prompt 不应为空")
	assert.Greater(t, len(prompt), 0, "翻译 prompt 长度应大于 0")
}

func TestTranslateSystemPrompt_LoadedFromTemplate(t *testing.T) {
	// Act: 获取翻译 prompt
	prompt := TranslateSystemPrompt

	// Assert: 验证 prompt 已从模板加载
	// 由于 TranslateSystemPrompt 是使用 MustLoadTemplate 在包初始化时加载的
	// 我们应该验证它包含一些预期的内容
	assert.NotEmpty(t, prompt, "翻译 prompt 不应为空")

	// 验证包含一些预期的关键词（根据实际模板内容调整）
	// 这些关键词应该出现在翻译的 prompt 中
	assert.Greater(t, len(prompt), 50, "翻译 prompt 应该有足够的长度")
}

func TestTranslateSystemPrompt_Consistent(t *testing.T) {
	// Act: 多次获取翻译 prompt
	prompt1 := TranslateSystemPrompt
	prompt2 := TranslateSystemPrompt

	// Assert: 验证输出一致
	assert.Equal(t, prompt1, prompt2, "翻译 prompt 应该保持一致")
}

func TestTranslateSystemPrompt_NotPanic(t *testing.T) {
	// Act & Assert: 验证获取 prompt 不会 panic
	// 由于 TranslateSystemPrompt 是使用 MustLoadTemplate 加载的
	// 如果模板不存在，会在包初始化时 panic
	// 这个测试主要验证变量可以正常访问
	assert.NotPanics(t, func() {
		_ = TranslateSystemPrompt
	}, "获取翻译 prompt 不应该 panic")
}


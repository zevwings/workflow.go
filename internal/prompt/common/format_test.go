//go:build test

package common

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zevwings/workflow/internal/prompt/io"
)

// TestFormatResult_WithoutFormatAnswer 测试不使用格式化函数的结果显示
func TestFormatResult_WithoutFormatAnswer(t *testing.T) {
	mockTerminal := io.NewMockTerminal([]byte{})

	err := FormatResult(mockTerminal, "提示消息", "结果文本", nil)

	assert.NoError(t, err)
	assert.True(t, mockTerminal.RestoreCursorCalled, "应该调用了恢复光标")
	assert.True(t, mockTerminal.ClearToEndCalled, "应该调用了清除到末尾")
	assert.True(t, mockTerminal.ResetFormatCalled, "应该调用了重置格式")
	assert.Contains(t, mockTerminal.GetOutput(), "提示消息", "输出应该包含提示消息")
	assert.Contains(t, mockTerminal.GetOutput(), "结果文本", "输出应该包含结果文本")
}

// TestFormatResult_WithFormatAnswer 测试使用格式化函数的结果显示
func TestFormatResult_WithFormatAnswer(t *testing.T) {
	mockTerminal := io.NewMockTerminal([]byte{})

	formatAnswer := func(text string) string {
		return "[格式化]" + text
	}

	err := FormatResult(mockTerminal, "提示消息", "结果文本", formatAnswer)

	assert.NoError(t, err)
	assert.Contains(t, mockTerminal.GetOutput(), "提示消息", "输出应该包含提示消息")
	assert.Contains(t, mockTerminal.GetOutput(), "[格式化]结果文本", "输出应该包含格式化后的结果文本")
}

// TestFormatResult_FormatAnswerNil 测试 formatAnswer 为 nil 的情况
func TestFormatResult_FormatAnswerNil(t *testing.T) {
	mockTerminal := io.NewMockTerminal([]byte{})

	err := FormatResult(mockTerminal, "提示", "结果", nil)

	assert.NoError(t, err)
	output := mockTerminal.GetOutput()
	assert.Contains(t, output, "提示", "输出应该包含提示")
	assert.Contains(t, output, "结果", "输出应该包含结果")
	// 确保没有格式化标记
	assert.NotContains(t, output, "[格式化]", "不应该包含格式化标记")
}

// TestFormatPromptWithPrefix_WithCustomPrefix 测试使用自定义前缀
func TestFormatPromptWithPrefix_WithCustomPrefix(t *testing.T) {
	cfg := PromptConfig{
		FormatQuestionPrefix: func() string {
			return ">>> "
		},
	}

	result := FormatPromptWithPrefix("请选择", cfg)
	assert.Equal(t, ">>> 请选择", result)
}

// TestFormatPromptWithPrefix_WithoutCustomPrefix 测试使用默认前缀
func TestFormatPromptWithPrefix_WithoutCustomPrefix(t *testing.T) {
	cfg := PromptConfig{
		FormatQuestionPrefix: nil,
	}

	result := FormatPromptWithPrefix("请选择", cfg)
	assert.Equal(t, "? 请选择", result)
}

// TestFormatResultWithTitle_IncludePrompt 测试包含提示消息的情况
func TestFormatResultWithTitle_IncludePrompt(t *testing.T) {
	mockTerminal := io.NewMockTerminal([]byte{})

	err := FormatResultWithTitle(
		mockTerminal,
		"提示消息",
		"结果文本",
		nil,
		true,
		"",
		nil,
	)

	assert.NoError(t, err)
	assert.True(t, mockTerminal.RestoreCursorCalled)
	assert.True(t, mockTerminal.ClearToEndCalled)
	assert.True(t, mockTerminal.ResetFormatCalled)
	assert.Contains(t, mockTerminal.GetOutput(), "提示消息")
	assert.Contains(t, mockTerminal.GetOutput(), "结果文本")
}

// TestFormatResultWithTitle_NotIncludePrompt 测试不包含提示消息的情况
func TestFormatResultWithTitle_NotIncludePrompt(t *testing.T) {
	mockTerminal := io.NewMockTerminal([]byte{})

	err := FormatResultWithTitle(
		mockTerminal,
		"提示消息",
		"结果文本",
		nil,
		false,
		"",
		nil,
	)

	assert.NoError(t, err)
	assert.True(t, mockTerminal.RestoreCursorCalled)
	assert.Greater(t, mockTerminal.MoveToStartCalled, 0)
	assert.Greater(t, mockTerminal.ClearLineCalled, 0)
	assert.True(t, mockTerminal.ClearToEndCalled)
	assert.True(t, mockTerminal.ResetFormatCalled)
	assert.Contains(t, mockTerminal.GetOutput(), "\033[2A", "应该包含向上移动两行的控制码")
}

// TestFormatResultWithTitle_WithFormatResultTitle 测试使用 FormatResultTitle
func TestFormatResultWithTitle_WithFormatResultTitle(t *testing.T) {
	mockTerminal := io.NewMockTerminal([]byte{})

	formatResultTitle := func(originalMessage, resultValue string) string {
		return originalMessage + " -> " + resultValue
	}

	err := FormatResultWithTitle(
		mockTerminal,
		"提示消息",
		"结果文本",
		nil,
		true,
		"原始消息",
		formatResultTitle,
	)

	assert.NoError(t, err)
	output := mockTerminal.GetOutput()
	assert.Contains(t, output, "原始消息 -> 结果文本")
}

// TestFormatResultWithTitle_WithFormatAnswerPrefix 测试使用 FormatAnswerPrefix
func TestFormatResultWithTitle_WithFormatAnswerPrefix(t *testing.T) {
	mockTerminal := io.NewMockTerminal([]byte{})

	formatAnswerPrefix := func() string {
		return ">>> "
	}

	err := FormatResultWithTitle(
		mockTerminal,
		"提示消息",
		"结果文本",
		nil,
		false,
		"",
		nil,
		formatAnswerPrefix,
	)

	assert.NoError(t, err)
	output := mockTerminal.GetOutput()
	assert.Contains(t, output, ">>> ")
}

// TestFormatResultInline_Basic 测试基本的行内格式化
func TestFormatResultInline_Basic(t *testing.T) {
	mockTerminal := io.NewMockTerminal([]byte{})

	err := FormatResultInline(
		mockTerminal,
		"提示消息",
		"结果文本",
		nil,
		"",
		nil,
		nil,
	)

	assert.NoError(t, err)
	assert.True(t, mockTerminal.RestoreCursorCalled)
	assert.Greater(t, mockTerminal.MoveToStartCalled, 0)
	assert.Greater(t, mockTerminal.ClearLineCalled, 0)
	assert.True(t, mockTerminal.ResetFormatCalled)
	assert.Contains(t, mockTerminal.GetOutput(), "提示消息")
	assert.Contains(t, mockTerminal.GetOutput(), "结果文本")
}

// TestFormatResultInline_WithFormatAnswer 测试使用格式化答案函数
func TestFormatResultInline_WithFormatAnswer(t *testing.T) {
	mockTerminal := io.NewMockTerminal([]byte{})

	formatAnswer := func(text string) string {
		return "[答案]" + text
	}

	err := FormatResultInline(
		mockTerminal,
		"提示消息",
		"结果文本",
		formatAnswer,
		"",
		nil,
		nil,
	)

	assert.NoError(t, err)
	output := mockTerminal.GetOutput()
	assert.Contains(t, output, "[答案]结果文本")
}

// TestFormatResultInline_WithFormatResultTitle 测试使用 FormatResultTitle
func TestFormatResultInline_WithFormatResultTitle(t *testing.T) {
	mockTerminal := io.NewMockTerminal([]byte{})

	formatResultTitle := func(originalMessage, resultValue string) string {
		return "标题: " + originalMessage
	}

	err := FormatResultInline(
		mockTerminal,
		"提示消息",
		"结果文本",
		nil,
		"原始消息",
		formatResultTitle,
		nil,
	)

	assert.NoError(t, err)
	output := mockTerminal.GetOutput()
	assert.Contains(t, output, "标题: 原始消息")
}

// TestFormatResultInline_WithFormatAnswerPrefix 测试使用 FormatAnswerPrefix
func TestFormatResultInline_WithFormatAnswerPrefix(t *testing.T) {
	mockTerminal := io.NewMockTerminal([]byte{})

	formatAnswerPrefix := func() string {
		return ">>> "
	}

	err := FormatResultInline(
		mockTerminal,
		"提示消息",
		"结果文本",
		nil,
		"",
		nil,
		formatAnswerPrefix,
	)

	assert.NoError(t, err)
	output := mockTerminal.GetOutput()
	assert.Contains(t, output, ">>> ")
}

// TestFormatResultWithOptions_Basic 测试基本选项
func TestFormatResultWithOptions_Basic(t *testing.T) {
	mockTerminal := io.NewMockTerminal([]byte{})

	err := FormatResultWithOptions(
		mockTerminal,
		"提示消息",
		"结果文本",
		nil,
		true,
	)

	assert.NoError(t, err)
	assert.True(t, mockTerminal.RestoreCursorCalled)
	assert.Contains(t, mockTerminal.GetOutput(), "提示消息")
	assert.Contains(t, mockTerminal.GetOutput(), "结果文本")
}

// TestFormatResultWithOptions_NotIncludePrompt 测试不包含提示消息
func TestFormatResultWithOptions_NotIncludePrompt(t *testing.T) {
	mockTerminal := io.NewMockTerminal([]byte{})

	err := FormatResultWithOptions(
		mockTerminal,
		"提示消息",
		"结果文本",
		nil,
		false,
	)

	assert.NoError(t, err)
	assert.True(t, mockTerminal.RestoreCursorCalled)
	assert.Greater(t, mockTerminal.MoveToStartCalled, 0)
}

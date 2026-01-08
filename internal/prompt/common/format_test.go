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

//go:build test

package selectpkg

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zevwings/workflow/internal/prompt/common"
	"github.com/zevwings/workflow/internal/testutils"
)

func TestSelectHandler_ValidateAndAdjustDefaultIndex(t *testing.T) {
	cfg := testutils.NewPromptConfig()

	options := []string{"A", "B", "C"}

	// 有效索引
	handler := NewSelectHandler(options, 1, cfg)
	assert.Equal(t, 1, handler.ValidateAndAdjustDefaultIndex())

	// 无效索引（负数）
	handler = NewSelectHandler(options, -1, cfg)
	assert.Equal(t, 0, handler.ValidateAndAdjustDefaultIndex())

	// 无效索引（超出范围）
	handler = NewSelectHandler(options, 10, cfg)
	assert.Equal(t, 0, handler.ValidateAndAdjustDefaultIndex())
}

func TestSelectHandler_ProcessArrowKey(t *testing.T) {
	cfg := testutils.NewPromptConfig()

	options := []string{"A", "B", "C"}
	handler := NewSelectHandler(options, 0, cfg)

	// 上箭头
	newIndex, shouldRender := handler.ProcessArrowKey(1, "up")
	assert.True(t, shouldRender)
	assert.Equal(t, 0, newIndex)

	// 上箭头在顶部
	newIndex, shouldRender = handler.ProcessArrowKey(0, "up")
	assert.False(t, shouldRender)
	assert.Equal(t, 0, newIndex)

	// 下箭头
	newIndex, shouldRender = handler.ProcessArrowKey(1, "down")
	assert.True(t, shouldRender)
	assert.Equal(t, 2, newIndex)

	// 下箭头在底部
	newIndex, shouldRender = handler.ProcessArrowKey(2, "down")
	assert.False(t, shouldRender)
	assert.Equal(t, 2, newIndex)
}

func TestSelectHandler_FormatOptionLine(t *testing.T) {
	cfg := testutils.NewPromptConfig()

	options := []string{"A", "B", "C"}
	handler := NewSelectHandler(options, 0, cfg)

	// 当前行
	line, isHighlighted := handler.FormatOptionLine(0, 0)
	assert.True(t, isHighlighted)
	assert.Contains(t, line, ">")
	assert.Contains(t, line, "A")

	// 非当前行
	line, isHighlighted = handler.FormatOptionLine(1, 0)
	assert.False(t, isHighlighted)
	assert.Contains(t, line, "  ")
	assert.Contains(t, line, "B")
}

func TestSelectHandler_FormatSelectedOption(t *testing.T) {
	cfg := common.PromptConfig{
		FormatPrompt:        func(msg string) string { return msg },
		FormatAnswer:        func(v string) string { return "[" + v + "]" },
		FormatError:         nil,
		FormatHint:          func(msg string) string { return msg },
		FormatQuestionPrefix: nil,
		FormatAnswerPrefix:   nil,
		FormatResultTitle:   nil,
	}

	options := []string{"A", "B", "C"}
	handler := NewSelectHandler(options, 0, cfg)

	result := handler.FormatSelectedOption(1)
	assert.Equal(t, "[B]", result)
}

func TestSelectHandler_ParseNumericInput(t *testing.T) {
	cfg := testutils.NewPromptConfig()

	options := []string{"A", "B", "C"}
	handler := NewSelectHandler(options, 0, cfg)

	// 有效输入
	result := handler.ParseNumericInput(2)
	assert.Equal(t, 1, result)

	// 无效输入（小于1）
	result = handler.ParseNumericInput(0)
	assert.Equal(t, 0, result) // 返回默认值

	// 无效输入（超出范围）
	result = handler.ParseNumericInput(10)
	assert.Equal(t, 0, result) // 返回默认值
}

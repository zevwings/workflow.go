//go:build test

package multiselect

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zevwings/workflow/internal/prompt/common"
	"github.com/zevwings/workflow/internal/testutils"
)

func TestMultiSelectHandler_ValidateAndCleanDefaults(t *testing.T) {
	cfg := testutils.NewPromptConfig()

	options := []string{"A", "B", "C", "D"}
	// 包含无效索引：-1 (无效), 10 (超出范围), 1 (有效)
	defaultSelected := []int{-1, 10, 1}

	handler := NewMultiSelectHandler(options, defaultSelected, cfg)
	selected := handler.ValidateAndCleanDefaults()

	// 只有有效的索引 1 应该被保留
	assert.True(t, selected[1])
	assert.False(t, selected[-1])
	assert.False(t, selected[10])
}

func TestMultiSelectHandler_GetInitialCurrentIndex(t *testing.T) {
	cfg := testutils.NewPromptConfig()

	options := []string{"A", "B", "C"}

	// 有有效默认选中项
	handler := NewMultiSelectHandler(options, []int{2}, cfg)
	assert.Equal(t, 2, handler.GetInitialCurrentIndex())

	// 无默认选中项
	handler = NewMultiSelectHandler(options, []int{}, cfg)
	assert.Equal(t, 0, handler.GetInitialCurrentIndex())

	// 无效默认选中项
	handler = NewMultiSelectHandler(options, []int{10}, cfg)
	assert.Equal(t, 0, handler.GetInitialCurrentIndex())
}

func TestMultiSelectHandler_ProcessArrowKey(t *testing.T) {
	cfg := testutils.NewPromptConfig()

	options := []string{"A", "B", "C"}
	handler := NewMultiSelectHandler(options, []int{}, cfg)

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

func TestMultiSelectHandler_ToggleSelection(t *testing.T) {
	cfg := testutils.NewPromptConfig()

	options := []string{"A", "B", "C"}
	handler := NewMultiSelectHandler(options, []int{}, cfg)

	selected := make(map[int]bool)

	// 切换选择
	handler.ToggleSelection(selected, 1)
	assert.True(t, selected[1])

	// 再次切换（取消选择）
	handler.ToggleSelection(selected, 1)
	assert.False(t, selected[1])
}

func TestMultiSelectHandler_FormatOptionLine(t *testing.T) {
	cfg := testutils.NewPromptConfig()

	options := []string{"A", "B", "C"}
	handler := NewMultiSelectHandler(options, []int{}, cfg)

	selected := map[int]bool{1: true}

	// 当前行，未选中
	line, isHighlighted := handler.FormatOptionLine(0, 0, selected)
	assert.True(t, isHighlighted)
	assert.Contains(t, line, ">")
	assert.Contains(t, line, "[ ]")
	assert.Contains(t, line, "A")

	// 当前行，已选中
	line, isHighlighted = handler.FormatOptionLine(1, 1, selected)
	assert.True(t, isHighlighted)
	assert.Contains(t, line, ">")
	assert.Contains(t, line, "[x]")
	assert.Contains(t, line, "B")

	// 非当前行
	line, isHighlighted = handler.FormatOptionLine(2, 1, selected)
	assert.False(t, isHighlighted)
	assert.Contains(t, line, "  ")
	assert.Contains(t, line, "[ ]")
	assert.Contains(t, line, "C")
}

func TestMultiSelectHandler_FormatSelectedOptions(t *testing.T) {
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
	handler := NewMultiSelectHandler(options, []int{}, cfg)

	// 未选择
	result := handler.FormatSelectedOptions([]int{})
	assert.Equal(t, "(未选择)", result)

	// 选择多个
	result = handler.FormatSelectedOptions([]int{0, 2})
	assert.Contains(t, result, "A")
	assert.Contains(t, result, "C")
	assert.Contains(t, result, "[")
}

func TestMultiSelectHandler_ParseCommaSeparatedInput(t *testing.T) {
	cfg := testutils.NewPromptConfig()

	options := []string{"A", "B", "C", "D"}
	handler := NewMultiSelectHandler(options, []int{}, cfg)

	// 有效输入
	result := handler.ParseCommaSeparatedInput("1,3")
	assert.ElementsMatch(t, []int{0, 2}, result)

	// 包含无效数字
	result = handler.ParseCommaSeparatedInput("abc,5,2")
	assert.ElementsMatch(t, []int{1}, result)

	// 超出范围
	result = handler.ParseCommaSeparatedInput("0,10,1")
	assert.ElementsMatch(t, []int{0}, result)

	// 空输入
	result = handler.ParseCommaSeparatedInput("")
	assert.Equal(t, []int{}, result)
}

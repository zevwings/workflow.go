//go:build test

package multiselect

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zevwings/workflow/internal/prompt/io"
)

// ==================== MultiSelect 主函数测试 ====================

// TestMultiSelect_EmptyOptions 验证空选项时直接返回错误
func TestMultiSelect_EmptyOptions(t *testing.T) {
	cfg := Config{
		FormatPrompt: func(msg string) string { return msg },
		FormatAnswer: func(v string) string { return v },
		FormatHint:   func(msg string) string { return msg },
	}

	mockTerminal := io.NewMockTerminal([]byte{})
	indices, err := MultiSelect("请选择", []string{}, nil, cfg, mockTerminal)
	assert.Error(t, err)
	assert.Nil(t, indices)
	assert.Contains(t, err.Error(), "选项列表不能为空")
}

// TestMultiSelect_WithMockTerminal_EnterKey 验证直接回车时返回空选择
func TestMultiSelect_WithMockTerminal_EnterKey(t *testing.T) {
	cfg := Config{
		FormatPrompt: func(msg string) string { return msg },
		FormatAnswer: func(v string) string { return v },
		FormatHint:   func(msg string) string { return msg },
	}

	options := []string{"A", "B", "C"}
	mockTerminal := io.NewMockTerminal([]byte{'\r'}) // 直接回车

	indices, err := MultiSelect("请选择", options, []int{}, cfg, mockTerminal)
	assert.NoError(t, err)
	assert.Empty(t, indices) // 未选择任何选项
}

// TestMultiSelect_WithMockTerminal_SpaceThenEnter 验证空格选中后回车
func TestMultiSelect_WithMockTerminal_SpaceThenEnter(t *testing.T) {
	cfg := Config{
		FormatPrompt: func(msg string) string { return msg },
		FormatAnswer: func(v string) string { return v },
		FormatHint:   func(msg string) string { return msg },
	}

	options := []string{"A", "B", "C"}
	// 空格选中第一个，然后回车
	mockTerminal := io.NewMockTerminal([]byte{' ', '\r'})

	indices, err := MultiSelect("请选择", options, []int{}, cfg, mockTerminal)
	assert.NoError(t, err)
	assert.ElementsMatch(t, []int{0}, indices)
}

// TestMultiSelect_WithMockTerminal_ArrowKeys 验证箭头键导航
func TestMultiSelect_WithMockTerminal_ArrowKeys(t *testing.T) {
	cfg := Config{
		FormatPrompt: func(msg string) string { return msg },
		FormatAnswer: func(v string) string { return v },
		FormatHint:   func(msg string) string { return msg },
	}

	options := []string{"A", "B", "C"}
	// 下箭头（转义序列：ESC [ B），然后空格选中，然后回车
	mockTerminal := io.NewMockTerminal([]byte{0x1b, '[', 'B', ' ', '\r'})

	indices, err := MultiSelect("请选择", options, []int{}, cfg, mockTerminal)
	assert.NoError(t, err)
	assert.ElementsMatch(t, []int{1}, indices) // 选中第二个选项
}

// TestMultiSelect_WithMockTerminal_CtrlC 验证 Ctrl+C 取消
func TestMultiSelect_WithMockTerminal_CtrlC(t *testing.T) {
	cfg := Config{
		FormatPrompt: func(msg string) string { return msg },
		FormatAnswer: func(v string) string { return v },
		FormatHint:   func(msg string) string { return msg },
	}

	options := []string{"A", "B", "C"}
	mockTerminal := io.NewMockTerminal([]byte{3}) // Ctrl+C

	indices, err := MultiSelect("请选择", options, []int{}, cfg, mockTerminal)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "用户取消输入")
	assert.Nil(t, indices)
}

// ==================== multiselectFallback 测试 ====================

// Test_multiselectFallback_ParseInput 验证回退模式下解析逗号分隔的编号
func Test_multiselectFallback_ParseInput(t *testing.T) {
	options := []string{"A", "B", "C", "D"}
	cfg := Config{
		FormatPrompt: func(msg string) string { return msg },
		FormatAnswer: func(v string) string { return v },
		FormatHint:   func(msg string) string { return msg },
	}

	mockTerminal := io.NewMockTerminalWithLines([]string{"1,3"})
	indices, err := multiselectFallback("请选择", options, []int{}, cfg, mockTerminal)
	assert.NoError(t, err)
	assert.ElementsMatch(t, []int{0, 2}, indices)
}

// Test_multiselectFallback_EmptyInput 验证空输入时返回默认选中项
func Test_multiselectFallback_EmptyInput(t *testing.T) {
	options := []string{"A", "B", "C"}
	defaultSelected := []int{0, 2}
	cfg := Config{
		FormatPrompt: func(msg string) string { return msg },
		FormatAnswer: func(v string) string { return v },
		FormatHint:   func(msg string) string { return msg },
	}

	mockTerminal := io.NewMockTerminalWithLines([]string{""})
	indices, err := multiselectFallback("请选择", options, defaultSelected, cfg, mockTerminal)
	assert.NoError(t, err)
	assert.ElementsMatch(t, defaultSelected, indices)
}

// Test_multiselectFallback_InvalidNumbers 验证无效数字输入会被忽略
func Test_multiselectFallback_InvalidNumbers(t *testing.T) {
	options := []string{"A", "B", "C"}
	cfg := Config{
		FormatPrompt: func(msg string) string { return msg },
		FormatAnswer: func(v string) string { return v },
		FormatHint:   func(msg string) string { return msg },
	}

	// 输入包含无效数字：abc, 5 (超出范围), 2 (有效)
	mockTerminal := io.NewMockTerminalWithLines([]string{"abc,5,2"})
	indices, err := multiselectFallback("请选择", options, []int{}, cfg, mockTerminal)
	assert.NoError(t, err)
	// 只有有效的 2 应该被选中（索引 1）
	assert.ElementsMatch(t, []int{1}, indices)
}

// Test_multiselectFallback_OutOfRangeNumbers 验证超出范围的数字会被忽略
func Test_multiselectFallback_OutOfRangeNumbers(t *testing.T) {
	options := []string{"A", "B", "C"}
	cfg := Config{
		FormatPrompt: func(msg string) string { return msg },
		FormatAnswer: func(v string) string { return v },
		FormatHint:   func(msg string) string { return msg },
	}

	// 输入超出范围的数字：0 (无效，从1开始), 10 (超出范围), 1 (有效)
	mockTerminal := io.NewMockTerminalWithLines([]string{"0,10,1"})
	indices, err := multiselectFallback("请选择", options, []int{}, cfg, mockTerminal)
	assert.NoError(t, err)
	// 只有有效的 1 应该被选中（索引 0）
	assert.ElementsMatch(t, []int{0}, indices)
}

// Test_multiselectFallback_DefaultSelected 验证默认选中项会被正确显示
func Test_multiselectFallback_DefaultSelected(t *testing.T) {
	options := []string{"A", "B", "C", "D"}
	defaultSelected := []int{1, 3}
	cfg := Config{
		FormatPrompt: func(msg string) string { return msg },
		FormatAnswer: func(v string) string { return v },
		FormatHint:   func(msg string) string { return msg },
	}

	// 直接回车，使用默认值
	mockTerminal := io.NewMockTerminalWithLines([]string{""})
	indices, err := multiselectFallback("请选择", options, defaultSelected, cfg, mockTerminal)
	assert.NoError(t, err)
	assert.ElementsMatch(t, defaultSelected, indices)
}

// Test_multiselectFallback_InvalidDefaultIndices 验证无效的默认索引会被过滤
func Test_multiselectFallback_InvalidDefaultIndices(t *testing.T) {
	options := []string{"A", "B", "C"}
	// 包含无效索引：-1 (无效), 10 (超出范围), 1 (有效)
	defaultSelected := []int{-1, 10, 1}
	cfg := Config{
		FormatPrompt: func(msg string) string { return msg },
		FormatAnswer: func(v string) string { return v },
		FormatHint:   func(msg string) string { return msg },
	}

	// 空输入，会返回清理后的 defaultSelected
	mockTerminal := io.NewMockTerminalWithLines([]string{""})
	indices, err := multiselectFallback("请选择", options, defaultSelected, cfg, mockTerminal)
	assert.NoError(t, err)
	// 空输入时，multiselectFallback 返回清理后的 defaultSelected（只包含有效索引）
	assert.ElementsMatch(t, []int{1}, indices)
}

// Test_multiselectFallback_NoSelection 验证未选择任何选项时返回空切片
func Test_multiselectFallback_NoSelection(t *testing.T) {
	options := []string{"A", "B", "C"}
	cfg := Config{
		FormatPrompt: func(msg string) string { return msg },
		FormatAnswer: func(v string) string { return v },
		FormatHint:   func(msg string) string { return msg },
	}

	// 输入空字符串
	mockTerminal := io.NewMockTerminalWithLines([]string{""})
	indices, err := multiselectFallback("请选择", options, []int{}, cfg, mockTerminal)
	assert.NoError(t, err)
	assert.Empty(t, indices)
}

// ==================== mapToSlice 辅助函数测试 ====================

// Test_mapToSlice_Sorted 验证 mapToSlice 会返回排序后的结果
func Test_mapToSlice_Sorted(t *testing.T) {
	m := map[int]bool{
		3: true,
		1: true,
		2: true,
	}

	result := mapToSlice(m)
	assert.Equal(t, []int{1, 2, 3}, result)
}

// Test_mapToSlice_EmptyMap 验证空 map 返回空切片
func Test_mapToSlice_EmptyMap(t *testing.T) {
	m := map[int]bool{}
	result := mapToSlice(m)
	assert.Empty(t, result)
}

// Test_mapToSlice_SingleElement 验证单个元素的情况
func Test_mapToSlice_SingleElement(t *testing.T) {
	m := map[int]bool{5: true}
	result := mapToSlice(m)
	assert.Equal(t, []int{5}, result)
}

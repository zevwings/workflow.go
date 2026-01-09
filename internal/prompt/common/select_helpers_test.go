//go:build test

package common

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zevwings/workflow/internal/prompt/io"
)

// newTestPromptConfig 创建测试用的 PromptConfig（避免导入循环）
func newTestPromptConfig() PromptConfig {
	return PromptConfig{
		FormatPrompt:         func(msg string) string { return msg },
		FormatAnswer:         func(v string) string { return v },
		FormatError:          nil,
		FormatHint:           func(msg string) string { return msg },
		FormatQuestionPrefix: func() string { return "? " },
		FormatAnswerPrefix:   func() string { return "> " },
		FormatResultTitle:    func(originalMessage string, resultValue string) string { return originalMessage },
	}
}

// TestExecuteSelectFallback_ValidInput 测试有效输入
func TestExecuteSelectFallback_ValidInput(t *testing.T) {
	options := []string{"选项A", "选项B", "选项C"}
	cfg := newTestPromptConfig()
	mockTerminal := io.NewMockTerminalWithLines([]string{"2"})

	fallbackOptions := SelectFallbackOptions{
		FormatOptionLine: func(index int, option string, isDefault bool) string {
			marker := " "
			if isDefault {
				marker = "*"
			}
			return fmt.Sprintf("  %s %d. %s\n", marker, index+1, option)
		},
		GetDefaultIndex: func() int {
			return 0
		},
		ParseInput: func(input string) (int, bool) {
			var num int
			_, err := fmt.Sscanf(input, "%d", &num)
			if err != nil {
				return 0, false
			}
			if num < 1 || num > len(options) {
				return 0, false
			}
			return num - 1, true
		},
		FormatSelectedOption: func(index int) string {
			return options[index]
		},
		InputPrompt:   fmt.Sprintf("请选择 (1-%d): ", len(options)),
		ResultPrefix:  "已选择: ",
	}

	selectedIndex, err := ExecuteSelectFallback(
		mockTerminal,
		"请选择一个选项",
		cfg,
		options,
		fallbackOptions,
	)

	assert.NoError(t, err)
	assert.Equal(t, 1, selectedIndex) // 用户输入 2，对应索引 1
}

// TestExecuteSelectFallback_InvalidInput 测试无效输入（返回默认值）
func TestExecuteSelectFallback_InvalidInput(t *testing.T) {
	options := []string{"选项A", "选项B", "选项C"}
	cfg := newTestPromptConfig()
	mockTerminal := io.NewMockTerminalWithLines([]string{"invalid"})
	defaultIndex := 1

	fallbackOptions := SelectFallbackOptions{
		FormatOptionLine: func(index int, option string, isDefault bool) string {
			marker := " "
			if isDefault {
				marker = "*"
			}
			return fmt.Sprintf("  %s %d. %s\n", marker, index+1, option)
		},
		GetDefaultIndex: func() int {
			return defaultIndex
		},
		ParseInput: func(input string) (int, bool) {
			var num int
			_, err := fmt.Sscanf(input, "%d", &num)
			if err != nil {
				return 0, false
			}
			if num < 1 || num > len(options) {
				return 0, false
			}
			return num - 1, true
		},
		FormatSelectedOption: func(index int) string {
			return options[index]
		},
		InputPrompt:   fmt.Sprintf("请选择 (1-%d): ", len(options)),
		ResultPrefix:  "已选择: ",
	}

	selectedIndex, err := ExecuteSelectFallback(
		mockTerminal,
		"请选择一个选项",
		cfg,
		options,
		fallbackOptions,
	)

	assert.NoError(t, err)
	assert.Equal(t, defaultIndex, selectedIndex) // 无效输入应返回默认值
}

// TestExecuteSelectFallback_OutOfRangeInput 测试超出范围的输入（返回默认值）
func TestExecuteSelectFallback_OutOfRangeInput(t *testing.T) {
	options := []string{"选项A", "选项B", "选项C"}
	cfg := newTestPromptConfig()
	mockTerminal := io.NewMockTerminalWithLines([]string{"10"}) // 超出范围
	defaultIndex := 2

	fallbackOptions := SelectFallbackOptions{
		FormatOptionLine: func(index int, option string, isDefault bool) string {
			marker := " "
			if isDefault {
				marker = "*"
			}
			return fmt.Sprintf("  %s %d. %s\n", marker, index+1, option)
		},
		GetDefaultIndex: func() int {
			return defaultIndex
		},
		ParseInput: func(input string) (int, bool) {
			var num int
			_, err := fmt.Sscanf(input, "%d", &num)
			if err != nil {
				return 0, false
			}
			if num < 1 || num > len(options) {
				return 0, false
			}
			return num - 1, true
		},
		FormatSelectedOption: func(index int) string {
			return options[index]
		},
		InputPrompt:   fmt.Sprintf("请选择 (1-%d): ", len(options)),
		ResultPrefix:  "已选择: ",
	}

	selectedIndex, err := ExecuteSelectFallback(
		mockTerminal,
		"请选择一个选项",
		cfg,
		options,
		fallbackOptions,
	)

	assert.NoError(t, err)
	assert.Equal(t, defaultIndex, selectedIndex) // 超出范围应返回默认值
}

// TestExecuteSelectFallback_EmptyInput 测试空输入（返回默认值）
func TestExecuteSelectFallback_EmptyInput(t *testing.T) {
	options := []string{"选项A", "选项B", "选项C"}
	cfg := newTestPromptConfig()
	mockTerminal := io.NewMockTerminalWithLines([]string{""}) // 空输入
	defaultIndex := 0

	fallbackOptions := SelectFallbackOptions{
		FormatOptionLine: func(index int, option string, isDefault bool) string {
			marker := " "
			if isDefault {
				marker = "*"
			}
			return fmt.Sprintf("  %s %d. %s\n", marker, index+1, option)
		},
		GetDefaultIndex: func() int {
			return defaultIndex
		},
		ParseInput: func(input string) (int, bool) {
			if input == "" {
				return 0, false
			}
			var num int
			_, err := fmt.Sscanf(input, "%d", &num)
			if err != nil {
				return 0, false
			}
			if num < 1 || num > len(options) {
				return 0, false
			}
			return num - 1, true
		},
		FormatSelectedOption: func(index int) string {
			return options[index]
		},
		InputPrompt:   fmt.Sprintf("请选择 (1-%d): ", len(options)),
		ResultPrefix:  "已选择: ",
	}

	selectedIndex, err := ExecuteSelectFallback(
		mockTerminal,
		"请选择一个选项",
		cfg,
		options,
		fallbackOptions,
	)

	assert.NoError(t, err)
	assert.Equal(t, defaultIndex, selectedIndex) // 空输入应返回默认值
}

// TestExecuteSelectFallback_ReadError 测试读取失败（返回默认值，不返回错误）
func TestExecuteSelectFallback_ReadError(t *testing.T) {
	options := []string{"选项A", "选项B", "选项C"}
	cfg := newTestPromptConfig()
	// 创建一个会在 ReadLine 时返回错误的 terminal
	mockTerminal := io.NewMockTerminal([]byte{}) // 空输入，ReadLine 会返回错误
	defaultIndex := 1

	fallbackOptions := SelectFallbackOptions{
		FormatOptionLine: func(index int, option string, isDefault bool) string {
			marker := " "
			if isDefault {
				marker = "*"
			}
			return fmt.Sprintf("  %s %d. %s\n", marker, index+1, option)
		},
		GetDefaultIndex: func() int {
			return defaultIndex
		},
		ParseInput: func(input string) (int, bool) {
			var num int
			_, err := fmt.Sscanf(input, "%d", &num)
			if err != nil {
				return 0, false
			}
			if num < 1 || num > len(options) {
				return 0, false
			}
			return num - 1, true
		},
		FormatSelectedOption: func(index int) string {
			return options[index]
		},
		InputPrompt:   fmt.Sprintf("请选择 (1-%d): ", len(options)),
		ResultPrefix:  "已选择: ",
	}

	selectedIndex, err := ExecuteSelectFallback(
		mockTerminal,
		"请选择一个选项",
		cfg,
		options,
		fallbackOptions,
	)

	// 根据错误处理策略，读取失败应返回默认值，不返回错误
	assert.NoError(t, err)
	assert.Equal(t, defaultIndex, selectedIndex)
}

// TestExecuteSelectFallback_NoInputPrompt 测试没有输入提示的情况
func TestExecuteSelectFallback_NoInputPrompt(t *testing.T) {
	options := []string{"选项A", "选项B"}
	cfg := newTestPromptConfig()
	mockTerminal := io.NewMockTerminalWithLines([]string{"1"})

	fallbackOptions := SelectFallbackOptions{
		FormatOptionLine: func(index int, option string, isDefault bool) string {
			return fmt.Sprintf("  %d. %s\n", index+1, option)
		},
		GetDefaultIndex: func() int {
			return 0
		},
		ParseInput: func(input string) (int, bool) {
			var num int
			_, err := fmt.Sscanf(input, "%d", &num)
			if err != nil {
				return 0, false
			}
			if num < 1 || num > len(options) {
				return 0, false
			}
			return num - 1, true
		},
		FormatSelectedOption: func(index int) string {
			return options[index]
		},
		InputPrompt:   "", // 没有输入提示
		ResultPrefix:  "已选择: ",
	}

	selectedIndex, err := ExecuteSelectFallback(
		mockTerminal,
		"请选择",
		cfg,
		options,
		fallbackOptions,
	)

	assert.NoError(t, err)
	assert.Equal(t, 0, selectedIndex)
}

// ==================== ExecuteMultiSelectFallback 测试 ====================

// TestExecuteMultiSelectFallback_ValidInput 测试有效输入
func TestExecuteMultiSelectFallback_ValidInput(t *testing.T) {
	options := []string{"选项A", "选项B", "选项C", "选项D"}
	cfg := newTestPromptConfig()
	mockTerminal := io.NewMockTerminalWithLines([]string{"1,3"})

	fallbackOptions := MultiSelectFallbackOptions{
		FormatOptionLine: func(index int, option string, isSelected bool) string {
			marker := " "
			if isSelected {
				marker = "*"
			}
			return fmt.Sprintf("  %s %d. %s\n", marker, index+1, option)
		},
		GetDefaultSelected: func() map[int]bool {
			return map[int]bool{}
		},
		ParseInput: func(input string) []int {
			// 简单的逗号分隔解析
			var indices []int
			var num int
			for _, part := range []string{"1", "3"} { // 模拟解析 "1,3"
				if _, err := fmt.Sscanf(part, "%d", &num); err == nil {
					if num >= 1 && num <= len(options) {
						indices = append(indices, num-1)
					}
				}
			}
			return indices
		},
		FormatSelectedOptions: func(selectedIndices []int) string {
			var texts []string
			for _, idx := range selectedIndices {
				texts = append(texts, options[idx])
			}
			return fmt.Sprintf("%v", texts)
		},
		InputPrompt:   "请选择 (例如: 1,3,5): ",
		ResultPrefix:  "已选择: ",
		EmptyResultText: "未选择任何选项",
	}

	selectedIndices, err := ExecuteMultiSelectFallback(
		mockTerminal,
		"请选择多个选项",
		cfg,
		options,
		fallbackOptions,
	)

	assert.NoError(t, err)
	assert.ElementsMatch(t, []int{0, 2}, selectedIndices) // 用户输入 1,3，对应索引 0,2
}

// TestExecuteMultiSelectFallback_EmptyInput 测试空输入（ParseInput 返回空列表，不返回默认值）
func TestExecuteMultiSelectFallback_EmptyInput(t *testing.T) {
	options := []string{"选项A", "选项B", "选项C"}
	cfg := newTestPromptConfig()
	mockTerminal := io.NewMockTerminalWithLines([]string{""})
	defaultSelected := map[int]bool{1: true, 2: true}

	fallbackOptions := MultiSelectFallbackOptions{
		FormatOptionLine: func(index int, option string, isSelected bool) string {
			marker := " "
			if isSelected {
				marker = "*"
			}
			return fmt.Sprintf("  %s %d. %s\n", marker, index+1, option)
		},
		GetDefaultSelected: func() map[int]bool {
			return defaultSelected
		},
		ParseInput: func(input string) []int {
			// 空输入时，ParseInput 返回空列表（这是正常行为，不是使用默认值）
			// 只有在 ReadLine 失败时才会返回默认值
			return []int{}
		},
		FormatSelectedOptions: func(selectedIndices []int) string {
			var texts []string
			for _, idx := range selectedIndices {
				texts = append(texts, options[idx])
			}
			return fmt.Sprintf("%v", texts)
		},
		InputPrompt:   "请选择 (例如: 1,3): ",
		ResultPrefix:  "已选择: ",
		EmptyResultText: "未选择任何选项",
	}

	selectedIndices, err := ExecuteMultiSelectFallback(
		mockTerminal,
		"请选择多个选项",
		cfg,
		options,
		fallbackOptions,
	)

	assert.NoError(t, err)
	// 空输入时，ParseInput 返回空列表，所以结果应该是空列表
	// 只有在 ReadLine 失败时才会返回默认值
	assert.Empty(t, selectedIndices)
}

// TestExecuteMultiSelectFallback_NoSelection 测试未选择任何选项
func TestExecuteMultiSelectFallback_NoSelection(t *testing.T) {
	options := []string{"选项A", "选项B", "选项C"}
	cfg := newTestPromptConfig()
	mockTerminal := io.NewMockTerminalWithLines([]string{""}) // 空输入，且没有默认选中

	fallbackOptions := MultiSelectFallbackOptions{
		FormatOptionLine: func(index int, option string, isSelected bool) string {
			return fmt.Sprintf("  %d. %s\n", index+1, option)
		},
		GetDefaultSelected: func() map[int]bool {
			return map[int]bool{} // 没有默认选中
		},
		ParseInput: func(input string) []int {
			return []int{} // 空输入返回空列表
		},
		FormatSelectedOptions: func(selectedIndices []int) string {
			var texts []string
			for _, idx := range selectedIndices {
				texts = append(texts, options[idx])
			}
			return fmt.Sprintf("%v", texts)
		},
		InputPrompt:   "请选择 (例如: 1,3): ",
		ResultPrefix:  "已选择: ",
		EmptyResultText: "未选择任何选项",
	}

	selectedIndices, err := ExecuteMultiSelectFallback(
		mockTerminal,
		"请选择多个选项",
		cfg,
		options,
		fallbackOptions,
	)

	assert.NoError(t, err)
	assert.Empty(t, selectedIndices) // 应返回空列表
}

// TestExecuteMultiSelectFallback_ReadError 测试读取失败（返回默认值）
func TestExecuteMultiSelectFallback_ReadError(t *testing.T) {
	options := []string{"选项A", "选项B", "选项C"}
	cfg := newTestPromptConfig()
	mockTerminal := io.NewMockTerminal([]byte{}) // 空输入，ReadLine 会返回错误
	defaultSelected := map[int]bool{0: true}

	fallbackOptions := MultiSelectFallbackOptions{
		FormatOptionLine: func(index int, option string, isSelected bool) string {
			marker := " "
			if isSelected {
				marker = "*"
			}
			return fmt.Sprintf("  %s %d. %s\n", marker, index+1, option)
		},
		GetDefaultSelected: func() map[int]bool {
			return defaultSelected
		},
		ParseInput: func(input string) []int {
			return []int{}
		},
		FormatSelectedOptions: func(selectedIndices []int) string {
			return ""
		},
		InputPrompt:   "请选择: ",
		ResultPrefix:  "已选择: ",
		EmptyResultText: "未选择",
	}

	selectedIndices, err := ExecuteMultiSelectFallback(
		mockTerminal,
		"请选择多个选项",
		cfg,
		options,
		fallbackOptions,
	)

	// 根据错误处理策略，读取失败应返回默认值，不返回错误
	assert.NoError(t, err)
	assert.ElementsMatch(t, []int{0}, selectedIndices)
}

// TestExecuteMultiSelectFallback_WithInstructions 测试带说明文本的情况
func TestExecuteMultiSelectFallback_WithInstructions(t *testing.T) {
	options := []string{"选项A", "选项B"}
	cfg := newTestPromptConfig()
	mockTerminal := io.NewMockTerminalWithLines([]string{"1"})

	fallbackOptions := MultiSelectFallbackOptions{
		FormatOptionLine: func(index int, option string, isSelected bool) string {
			return fmt.Sprintf("  %d. %s\n", index+1, option)
		},
		GetDefaultSelected: func() map[int]bool {
			return map[int]bool{}
		},
		ParseInput: func(input string) []int {
			var indices []int
			var num int
			if _, err := fmt.Sscanf(input, "%d", &num); err == nil {
				if num >= 1 && num <= len(options) {
					indices = append(indices, num-1)
				}
			}
			return indices
		},
		FormatSelectedOptions: func(selectedIndices []int) string {
			return "测试"
		},
		Instructions:  "这是说明文本",
		InputPrompt:   "请选择: ",
		ResultPrefix:  "已选择: ",
		EmptyResultText: "未选择",
	}

	selectedIndices, err := ExecuteMultiSelectFallback(
		mockTerminal,
		"请选择",
		cfg,
		options,
		fallbackOptions,
	)

	assert.NoError(t, err)
	assert.Equal(t, []int{0}, selectedIndices)
}

// ==================== mapToSlice 辅助函数测试 ====================

// TestMapToSlice_Basic 测试基本的 map 转 slice
func TestMapToSlice_Basic(t *testing.T) {
	// mapToSlice 是私有函数，我们通过 ExecuteMultiSelectFallback 间接测试
	// 但我们可以创建一个简单的测试来验证其行为
	selected := map[int]bool{
		2: true,
		0: true,
		1: true,
	}
	
	// 通过 ExecuteMultiSelectFallback 的 ReadError 路径来测试 mapToSlice
	options := []string{"A", "B", "C"}
	cfg := newTestPromptConfig()
	mockTerminal := io.NewMockTerminal([]byte{}) // 空输入，ReadLine 会返回错误

	fallbackOptions := MultiSelectFallbackOptions{
		FormatOptionLine: func(index int, option string, isSelected bool) string {
			return fmt.Sprintf("  %d. %s\n", index+1, option)
		},
		GetDefaultSelected: func() map[int]bool {
			return selected
		},
		ParseInput: func(input string) []int {
			return []int{}
		},
		FormatSelectedOptions: func(selectedIndices []int) string {
			return ""
		},
		InputPrompt:   "请选择: ",
		ResultPrefix:  "已选择: ",
		EmptyResultText: "未选择",
	}

	selectedIndices, err := ExecuteMultiSelectFallback(
		mockTerminal,
		"请选择",
		cfg,
		options,
		fallbackOptions,
	)

	assert.NoError(t, err)
	// mapToSlice 应该返回排序后的索引
	assert.Equal(t, []int{0, 1, 2}, selectedIndices)
}

// TestMapToSlice_EmptyMap 测试空 map
func TestMapToSlice_EmptyMap(t *testing.T) {
	options := []string{"A", "B"}
	cfg := newTestPromptConfig()
	mockTerminal := io.NewMockTerminal([]byte{})

	fallbackOptions := MultiSelectFallbackOptions{
		FormatOptionLine: func(index int, option string, isSelected bool) string {
			return fmt.Sprintf("  %d. %s\n", index+1, option)
		},
		GetDefaultSelected: func() map[int]bool {
			return map[int]bool{} // 空 map
		},
		ParseInput: func(input string) []int {
			return []int{}
		},
		FormatSelectedOptions: func(selectedIndices []int) string {
			return ""
		},
		InputPrompt:   "请选择: ",
		ResultPrefix:  "已选择: ",
		EmptyResultText: "未选择",
	}

	selectedIndices, err := ExecuteMultiSelectFallback(
		mockTerminal,
		"请选择",
		cfg,
		options,
		fallbackOptions,
	)

	assert.NoError(t, err)
	// 空 map 应该返回空切片
	assert.Empty(t, selectedIndices)
}

// TestMapToSlice_SingleElement 测试单个元素
func TestMapToSlice_SingleElement(t *testing.T) {
	options := []string{"A", "B", "C"}
	cfg := newTestPromptConfig()
	mockTerminal := io.NewMockTerminal([]byte{})

	fallbackOptions := MultiSelectFallbackOptions{
		FormatOptionLine: func(index int, option string, isSelected bool) string {
			return fmt.Sprintf("  %d. %s\n", index+1, option)
		},
		GetDefaultSelected: func() map[int]bool {
			return map[int]bool{1: true} // 单个元素
		},
		ParseInput: func(input string) []int {
			return []int{}
		},
		FormatSelectedOptions: func(selectedIndices []int) string {
			return ""
		},
		InputPrompt:   "请选择: ",
		ResultPrefix:  "已选择: ",
		EmptyResultText: "未选择",
	}

	selectedIndices, err := ExecuteMultiSelectFallback(
		mockTerminal,
		"请选择",
		cfg,
		options,
		fallbackOptions,
	)

	assert.NoError(t, err)
	assert.Equal(t, []int{1}, selectedIndices)
}

// ==================== SetupInteractiveSelect 测试 ====================

// TestSetupInteractiveSelect_Basic 测试基本的设置
func TestSetupInteractiveSelect_Basic(t *testing.T) {
	cfg := newTestPromptConfig()
	mockTerminal := io.NewMockTerminal([]byte{})

	baseCfg := BasePromptConfig{
		Message:  "请选择",
		Config:   cfg,
		Terminal: mockTerminal,
	}

	setup := SetupInteractiveSelect(baseCfg)

	assert.NotNil(t, setup.RawModeMgr)
	assert.NotNil(t, setup.Parser)
	assert.NotNil(t, setup.Renderer)
	assert.Equal(t, "请选择", setup.PromptMsg)
	assert.Equal(t, "? 请选择", setup.PromptMsgWithPrefix)
}

// TestSetupInteractiveSelect_WithCustomPrefix 测试使用自定义前缀
func TestSetupInteractiveSelect_WithCustomPrefix(t *testing.T) {
	cfg := PromptConfig{
		FormatPrompt:         func(msg string) string { return msg },
		FormatAnswer:         func(v string) string { return v },
		FormatError:          nil,
		FormatHint:           func(msg string) string { return msg },
		FormatQuestionPrefix: func() string { return ">>> " },
		FormatAnswerPrefix:   func() string { return "> " },
		FormatResultTitle:    func(originalMessage string, resultValue string) string { return originalMessage },
	}
	mockTerminal := io.NewMockTerminal([]byte{})

	baseCfg := BasePromptConfig{
		Message:  "请选择",
		Config:   cfg,
		Terminal: mockTerminal,
	}

	setup := SetupInteractiveSelect(baseCfg)

	assert.Equal(t, ">>> 请选择", setup.PromptMsgWithPrefix)
}
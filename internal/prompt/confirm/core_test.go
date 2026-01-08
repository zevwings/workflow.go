//go:build test

package confirm

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zevwings/workflow/internal/prompt/io"
)

// ==================== Confirm 主函数测试（使用 MockTerminal） ====================

func TestConfirm_WithMockTerminal_YesInput(t *testing.T) {
	cfg := Config{
		FormatPrompt: func(msg string) string { return msg },
		FormatAnswer: func(v string) string { return v },
	}

	mockTerminal := io.NewMockTerminal([]byte{'y'})
	result, err := Confirm("是否继续？", false, cfg, mockTerminal)

	assert.NoError(t, err)
	assert.True(t, result)
	assert.Contains(t, mockTerminal.GetOutput(), "是否继续？")
	assert.Contains(t, mockTerminal.GetOutput(), "yes")
}

func TestConfirm_WithMockTerminal_NoInput(t *testing.T) {
	cfg := Config{
		FormatPrompt: func(msg string) string { return msg },
		FormatAnswer: func(v string) string { return v },
	}

	mockTerminal := io.NewMockTerminal([]byte{'n'})
	result, err := Confirm("是否继续？", true, cfg, mockTerminal)

	assert.NoError(t, err)
	assert.False(t, result)
	assert.Contains(t, mockTerminal.GetOutput(), "是否继续？")
	assert.Contains(t, mockTerminal.GetOutput(), "no")
}

func TestConfirm_WithMockTerminal_EnterKey(t *testing.T) {
	cfg := Config{
		FormatPrompt: func(msg string) string { return msg },
		FormatAnswer: func(v string) string { return v },
	}

	// 默认 yes，回车应返回 true
	mockTerminal := io.NewMockTerminal([]byte{'\r'})
	result, err := Confirm("是否继续？", true, cfg, mockTerminal)

	assert.NoError(t, err)
	assert.True(t, result)
	assert.Contains(t, mockTerminal.GetOutput(), "yes")

	// 默认 no，回车应返回 false
	mockTerminal = io.NewMockTerminal([]byte{'\n'})
	result, err = Confirm("是否继续？", false, cfg, mockTerminal)

	assert.NoError(t, err)
	assert.False(t, result)
	assert.Contains(t, mockTerminal.GetOutput(), "no")
}

func TestConfirm_WithMockTerminal_CtrlC(t *testing.T) {
	cfg := Config{
		FormatPrompt: func(msg string) string { return msg },
		FormatAnswer: func(v string) string { return v },
	}

	mockTerminal := io.NewMockTerminal([]byte{3}) // Ctrl+C
	result, err := Confirm("是否继续？", true, cfg, mockTerminal)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "用户取消输入")
	// 当用户取消时，result 的值取决于错误处理逻辑，但应该不是默认值
	// 由于 HandleCancel 返回错误，result 会保持为 defaultYes 的值（true）
	// 但这不是重点，重点是错误消息正确
	_ = result // 忽略 result 值，因为取消时重点是错误
}

func TestConfirm_WithMockTerminal_InvalidCharThenYes(t *testing.T) {
	cfg := Config{
		FormatPrompt: func(msg string) string { return msg },
		FormatAnswer: func(v string) string { return v },
	}

	// 先输入无效字符 'x'，然后输入 'y'
	mockTerminal := io.NewMockTerminal([]byte{'x', 'y'})
	result, err := Confirm("是否继续？", false, cfg, mockTerminal)

	assert.NoError(t, err)
	assert.True(t, result)
}

func TestConfirm_WithMockTerminal_TerminalControl(t *testing.T) {
	cfg := Config{
		FormatPrompt: func(msg string) string { return msg },
		FormatAnswer: func(v string) string { return v },
	}

	mockTerminal := io.NewMockTerminal([]byte{'y'})
	_, err := Confirm("test", false, cfg, mockTerminal)

	assert.NoError(t, err)
	// 验证终端控制被调用
	assert.True(t, mockTerminal.HideCursorCalled)
	assert.True(t, mockTerminal.ShowCursorCalled)
	// confirm 使用 SaveCursor、RestoreCursor、MoveToStart 和 ClearLine 来清除并显示结果
	assert.True(t, mockTerminal.SaveCursorCalled)
	assert.True(t, mockTerminal.RestoreCursorCalled)
	assert.True(t, mockTerminal.MoveToStartCalled > 0)
	assert.True(t, mockTerminal.ClearLineCalled > 0)
	assert.True(t, mockTerminal.ResetFormatCalled)
}

// ==================== confirmFallback 测试 ====================

// Test_confirmFallback_DefaultYes_NoInput 验证在默认值为 true 且无有效输入时返回默认值
func Test_confirmFallback_DefaultYes_NoInput(t *testing.T) {
	cfg := Config{
		FormatPrompt: func(msg string) string { return msg },
		FormatAnswer: func(v string) string { return v },
	}

	mockTerminal := io.NewMockTerminalWithLines([]string{""})
	result, err := confirmFallback("是否继续？", true, cfg, mockTerminal)
	assert.NoError(t, err)
	assert.True(t, result)
}

// Test_confirmFallback_DefaultNo_InvalidInput 验证非法输入时也会回落到默认值
func Test_confirmFallback_DefaultNo_InvalidInput(t *testing.T) {
	mockTerminal := io.NewMockTerminalWithLines([]string{"invalid"})

	cfg := Config{
		FormatPrompt: func(msg string) string { return msg },
		FormatAnswer: func(v string) string { return v },
	}

	result, err := confirmFallback("是否继续？", false, cfg, mockTerminal)
	assert.NoError(t, err)
	// 非法输入应回落到默认值 false
	assert.False(t, result)
}

// Test_confirmFallback_YesInput 验证输入 "y" 或 "yes" 时返回 true
func Test_confirmFallback_YesInput(t *testing.T) {
	testCases := []struct {
		name  string
		input string
	}{
		{"输入 y", "y"},
		{"输入 yes", "yes"},
		{"输入 Y", "Y"},
		{"输入 YES", "YES"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockTerminal := io.NewMockTerminalWithLines([]string{tc.input})

			cfg := Config{
				FormatPrompt: func(msg string) string { return msg },
				FormatAnswer: func(v string) string { return v },
			}

			result, err := confirmFallback("是否继续？", false, cfg, mockTerminal)
			assert.NoError(t, err)
			assert.True(t, result)
		})
	}
}

// Test_confirmFallback_NoInput 验证输入 "n" 或 "no" 时返回 false
func Test_confirmFallback_NoInput(t *testing.T) {
	testCases := []struct {
		name  string
		input string
	}{
		{"输入 n", "n"},
		{"输入 no", "no"},
		{"输入 N", "N"},
		{"输入 NO", "NO"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockTerminal := io.NewMockTerminalWithLines([]string{tc.input})

			cfg := Config{
				FormatPrompt: func(msg string) string { return msg },
				FormatAnswer: func(v string) string { return v },
			}

			result, err := confirmFallback("是否继续？", true, cfg, mockTerminal)
			assert.NoError(t, err)
			assert.False(t, result)
		})
	}
}

// Test_confirmFallback_EmptyInput 验证空输入时使用默认值
func Test_confirmFallback_EmptyInput(t *testing.T) {
	testCases := []struct {
		name       string
		defaultYes bool
		expected   bool
	}{
		{"默认 yes，空输入", true, true},
		{"默认 no，空输入", false, false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockTerminal := io.NewMockTerminalWithLines([]string{""})

			cfg := Config{
				FormatPrompt: func(msg string) string { return msg },
				FormatAnswer: func(v string) string { return v },
			}

			result, err := confirmFallback("是否继续？", tc.defaultYes, cfg, mockTerminal)
			assert.NoError(t, err)
			assert.Equal(t, tc.expected, result)
		})
	}
}

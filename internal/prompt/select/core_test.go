//go:build test

package selectpkg

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zevwings/workflow/internal/prompt/io"
)

// TestSelect_EmptyOptions 验证空选项时直接返回错误
func TestSelect_EmptyOptions(t *testing.T) {
	cfg := Config{
		FormatPrompt: func(msg string) string { return msg },
		FormatAnswer: func(v string) string { return v },
		FormatHint:   func(msg string) string { return msg },
	}

	mockTerminal := io.NewMockTerminal([]byte{})
	index, err := Select("请选择", []string{}, 0, cfg, mockTerminal)
	assert.Error(t, err)
	assert.Equal(t, -1, index)
}

// Test_selectFallback_ValidAndInvalidInput 验证回退模式下合法与非法输入行为
func Test_selectFallback_ValidAndInvalidInput(t *testing.T) {
	options := []string{"A", "B", "C"}
	cfg := Config{
		FormatPrompt: func(msg string) string { return msg },
		FormatAnswer: func(v string) string { return v },
		FormatHint:   func(msg string) string { return msg },
	}

	// 1) 非法输入时应返回默认索引
	mockTerminal := io.NewMockTerminalWithLines([]string{"invalid"})
	index, err := selectFallback("请选择", options, 1, cfg, mockTerminal)
	assert.NoError(t, err)
	assert.Equal(t, 1, index)

	// 2) 合法数字输入时应返回对应索引
	mockTerminal = io.NewMockTerminalWithLines([]string{"2"})
	index, err = selectFallback("请选择", options, 0, cfg, mockTerminal)
	assert.NoError(t, err)
	assert.Equal(t, 1, index)
}

// Test_selectFallback_OutOfRangeInput 验证超出范围的输入会返回默认索引
func Test_selectFallback_OutOfRangeInput(t *testing.T) {
	options := []string{"A", "B", "C"}
	cfg := Config{
		FormatPrompt: func(msg string) string { return msg },
		FormatAnswer: func(v string) string { return v },
		FormatHint:   func(msg string) string { return msg },
	}

	// 测试超出上限
	mockTerminal := io.NewMockTerminalWithLines([]string{"10"})
	index, err := selectFallback("请选择", options, 1, cfg, mockTerminal)
	assert.NoError(t, err)
	assert.Equal(t, 1, index) // 返回默认索引

	// 测试小于下限
	mockTerminal = io.NewMockTerminalWithLines([]string{"0"})
	index, err = selectFallback("请选择", options, 2, cfg, mockTerminal)
	assert.NoError(t, err)
	assert.Equal(t, 2, index) // 返回默认索引
}

// Test_selectFallback_EmptyInput 验证空输入时返回默认索引
func Test_selectFallback_EmptyInput(t *testing.T) {
	options := []string{"A", "B", "C"}
	cfg := Config{
		FormatPrompt: func(msg string) string { return msg },
		FormatAnswer: func(v string) string { return v },
		FormatHint:   func(msg string) string { return msg },
	}

	// 空输入会触发 ReadLine 错误，返回默认索引
	mockTerminal := io.NewMockTerminalWithLines([]string{})
	index, err := selectFallback("请选择", options, 1, cfg, mockTerminal)
	assert.NoError(t, err)
	assert.Equal(t, 1, index) // 返回默认索引
}

// TestSelect_InvalidDefaultIndex 验证无效的默认索引会被调整为 0
func TestSelect_InvalidDefaultIndex(t *testing.T) {
	cfg := Config{
		FormatPrompt: func(msg string) string { return msg },
		FormatAnswer: func(v string) string { return v },
		FormatHint:   func(msg string) string { return msg },
	}

	options := []string{"A", "B", "C"}

	// 测试负索引（使用 MockTerminal，直接回车）
	mockTerminal := io.NewMockTerminal([]byte{'\r'})
	index, err := Select("请选择", options, -1, cfg, mockTerminal)
	// 无效索引会被调整为 0，回车选择索引 0
	assert.NoError(t, err)
	assert.Equal(t, 0, index)

	// 测试超出范围的索引
	mockTerminal = io.NewMockTerminal([]byte{'\r'})
	index, err = Select("请选择", options, 10, cfg, mockTerminal)
	// 无效索引会被调整为 0，回车选择索引 0
	assert.NoError(t, err)
	assert.Equal(t, 0, index)
}



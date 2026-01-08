//go:build test

package input

import (
	"testing"

	"github.com/charmbracelet/lipgloss"
	"github.com/stretchr/testify/assert"
)

// ==================== FormatPlaceholder 测试 ====================

func TestFormatPlaceholder_DisableColor(t *testing.T) {
	// Arrange: 创建样式和禁用颜色
	hintStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	text := "example@domain.com"

	// Act: 格式化占位符（禁用颜色）
	result := FormatPlaceholder(text, hintStyle, false)

	// Assert: 应该返回原始文本
	assert.Equal(t, text, result, "禁用颜色时应返回原始文本")
}

func TestFormatPlaceholder_EnableColor(t *testing.T) {
	// Arrange: 创建样式和启用颜色
	hintStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	text := "example@domain.com"

	// Act: 格式化占位符（启用颜色）
	result := FormatPlaceholder(text, hintStyle, true)

	// Assert: 去除 ANSI 码后应还原为原始文本
	assert.Equal(t, text, StripAnsiCodes(result), "启用颜色时，去除 ANSI 码后应得到原始文本")
}

// ==================== StripAnsiCodes 测试 ====================

func TestStripAnsiCodes_NoAnsiCodes(t *testing.T) {
	// Arrange: 纯文本，无 ANSI 代码
	text := "plain text without codes"

	// Act: 去除 ANSI 代码
	result := StripAnsiCodes(text)

	// Assert: 应该返回原始文本
	assert.Equal(t, text, result, "无 ANSI 代码时应返回原始文本")
}

func TestStripAnsiCodes_WithAnsiCodes(t *testing.T) {
	// Arrange: 包含 ANSI 代码的文本
	text := "\x1b[31mred text\x1b[0m"

	// Act: 去除 ANSI 代码
	result := StripAnsiCodes(text)

	// Assert: 应该只返回纯文本
	assert.Equal(t, "red text", result, "应去除所有 ANSI 代码")
	assert.NotContains(t, result, "\x1b[", "结果不应包含 ANSI 转义码")
}

func TestStripAnsiCodes_ComplexAnsiCodes(t *testing.T) {
	// Arrange: 包含复杂 ANSI 代码的文本
	text := "\x1b[1;33;44mbold yellow on blue\x1b[0m normal text"

	// Act: 去除 ANSI 代码
	result := StripAnsiCodes(text)

	// Assert: 应该只返回纯文本
	assert.Equal(t, "bold yellow on blue normal text", result, "应去除所有 ANSI 代码")
}

func TestStripAnsiCodes_MultipleAnsiCodes(t *testing.T) {
	// Arrange: 包含多个 ANSI 代码的文本
	text := "\x1b[31mred\x1b[0m \x1b[32mgreen\x1b[0m \x1b[34mblue\x1b[0m"

	// Act: 去除 ANSI 代码
	result := StripAnsiCodes(text)

	// Assert: 应该只返回纯文本
	assert.Equal(t, "red green blue", result, "应去除所有 ANSI 代码")
}

func TestStripAnsiCodes_EmptyString(t *testing.T) {
	// Arrange: 空字符串
	text := ""

	// Act: 去除 ANSI 代码
	result := StripAnsiCodes(text)

	// Assert: 应该返回空字符串
	assert.Equal(t, "", result, "空字符串应返回空字符串")
}

func TestStripAnsiCodes_OnlyAnsiCodes(t *testing.T) {
	// Arrange: 只有 ANSI 代码，无文本
	text := "\x1b[31m\x1b[0m"

	// Act: 去除 ANSI 代码
	result := StripAnsiCodes(text)

	// Assert: 应该返回空字符串
	assert.Equal(t, "", result, "只有 ANSI 代码时应返回空字符串")
}

func TestStripAnsiCodes_WithPlaceholder(t *testing.T) {
	// Arrange: 使用 FormatPlaceholder 生成的文本（包含 ANSI 代码）
	hintStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	placeholder := "example@domain.com"
	formatted := FormatPlaceholder(placeholder, hintStyle, true)

	// Act: 去除 ANSI 代码
	result := StripAnsiCodes(formatted)

	// Assert: 应该返回原始占位符文本
	assert.Equal(t, placeholder, result, "应去除格式化后的 ANSI 代码，返回原始文本")
}



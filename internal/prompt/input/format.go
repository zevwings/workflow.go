package input

import (
	"regexp"

	"github.com/charmbracelet/lipgloss"
)

// FormatPlaceholder 格式化占位符文本（使用主题的 HintStyle，并添加斜体）
func FormatPlaceholder(text string, hintStyle lipgloss.Style, enableColor bool) string {
	if !enableColor {
		return text
	}
	// 使用主题的 HintStyle，并添加斜体效果
	placeholderStyle := hintStyle.Copy().Italic(true)
	return placeholderStyle.Render(text)
}

// StripAnsiCodes 去除 ANSI 转义码，返回纯文本（用于计算显示宽度）
func StripAnsiCodes(s string) string {
	// ANSI 转义码的正则表达式：\033[或\x1b[开头，后跟参数和命令字符
	ansiRegex := regexp.MustCompile(`\x1b\[[0-9;]*[a-zA-Z]`)
	return ansiRegex.ReplaceAllString(s, "")
}

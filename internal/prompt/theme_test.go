//go:build test

package prompt

import (
	"sync"
	"testing"

	"github.com/charmbracelet/lipgloss"
	"github.com/stretchr/testify/require"
)

func TestGetTheme_Default(t *testing.T) {
	// 获取默认主题
	theme := GetTheme()

	// 验证默认主题的基本属性
	require.True(t, theme.EnableColor)
	require.Equal(t, "!", theme.PrefixWarn)
	require.Equal(t, "x", theme.PrefixError)
	require.Equal(t, "[", theme.InputBracketLeft)
	require.Equal(t, "]", theme.InputBracketRight)
}

func TestSetTheme_GetTheme(t *testing.T) {
	// 保存原始主题
	originalTheme := GetTheme()

	// 恢复原始主题（测试结束后）
	defer SetTheme(originalTheme)

	// 创建自定义主题
	customTheme := Theme{
		EnableColor: false,
		PrefixWarn:  "WARN",
		PrefixError:  "ERROR",
		InfoStyle:    lipgloss.NewStyle().Foreground(lipgloss.Color("blue")),
	}

	// 设置新主题
	SetTheme(customTheme)

	// 验证主题已更新
	retrievedTheme := GetTheme()
	require.Equal(t, false, retrievedTheme.EnableColor)
	require.Equal(t, "WARN", retrievedTheme.PrefixWarn)
	require.Equal(t, "ERROR", retrievedTheme.PrefixError)
}

func TestSetTheme_GetTheme_ThreadSafe(t *testing.T) {
	// 保存原始主题
	originalTheme := GetTheme()
	defer SetTheme(originalTheme)

	// 并发测试线程安全性
	const numGoroutines = 10
	const numIterations = 100

	var wg sync.WaitGroup
	wg.Add(numGoroutines)

	// 启动多个 goroutine 同时读写主题
	for i := 0; i < numGoroutines; i++ {
		go func(id int) {
			defer wg.Done()
			for j := 0; j < numIterations; j++ {
				// 读取主题
				theme := GetTheme()
				require.NotNil(t, theme)

				// 设置新主题
				newTheme := Theme{
					EnableColor: id%2 == 0,
					PrefixWarn:  "WARN",
				}
				SetTheme(newTheme)

				// 再次读取验证
				retrieved := GetTheme()
				require.NotNil(t, retrieved)
			}
		}(i)
	}

	wg.Wait()
	// 测试完成后，主题应该处于某个有效状态
	finalTheme := GetTheme()
	require.NotNil(t, finalTheme)
}

func TestFormatTitle_EnableColor(t *testing.T) {
	// 保存原始主题
	originalTheme := GetTheme()
	defer SetTheme(originalTheme)

	// 测试启用颜色
	theme := GetTheme()
	theme.EnableColor = true
	SetTheme(theme)

	result := formatTitle("test message")
	require.Contains(t, result, "test message")
	// 注意：在非TTY环境下，lipgloss可能不会添加ANSI代码
	// 所以只检查内容存在即可

	// 测试禁用颜色
	theme.EnableColor = false
	SetTheme(theme)

	result = formatTitle("test message")
	require.Equal(t, "test message", result)
}

func TestFormatAnswer_EnableColor(t *testing.T) {
	// 保存原始主题
	originalTheme := GetTheme()
	defer SetTheme(originalTheme)

	// 测试启用颜色
	theme := GetTheme()
	theme.EnableColor = true
	SetTheme(theme)

	result := formatAnswer("answer value")
	require.Contains(t, result, "answer value")
	// 注意：在非TTY环境下，lipgloss可能不会添加ANSI代码
	// 所以只检查内容存在即可

	// 测试禁用颜色
	theme.EnableColor = false
	SetTheme(theme)

	result = formatAnswer("answer value")
	require.Equal(t, "answer value", result)
}

func TestFormatError_EnableColor(t *testing.T) {
	// 保存原始主题
	originalTheme := GetTheme()
	defer SetTheme(originalTheme)

	// 测试启用颜色
	theme := GetTheme()
	theme.EnableColor = true
	SetTheme(theme)

	result := formatError("error message")
	require.Contains(t, result, "* error message")
	// 注意：在非TTY环境下，lipgloss可能不会添加ANSI代码
	// 所以只检查内容存在即可

	// 测试禁用颜色
	theme.EnableColor = false
	SetTheme(theme)

	result = formatError("error message")
	require.Equal(t, "* error message", result)
}

func TestFormatHint_EnableColor(t *testing.T) {
	// 保存原始主题
	originalTheme := GetTheme()
	defer SetTheme(originalTheme)

	// 测试启用颜色
	theme := GetTheme()
	theme.EnableColor = true
	SetTheme(theme)

	result := formatHint("hint message")
	require.Contains(t, result, "hint message")
	// 注意：在非TTY环境下，lipgloss可能不会添加ANSI代码
	// 所以只检查内容存在即可

	// 测试禁用颜色
	theme.EnableColor = false
	SetTheme(theme)

	result = formatHint("hint message")
	require.Equal(t, "hint message", result)
}

func TestFormatError_AlwaysIncludesPrefix(t *testing.T) {
	// 保存原始主题
	originalTheme := GetTheme()
	defer SetTheme(originalTheme)

	theme := GetTheme()
	theme.EnableColor = false
	SetTheme(theme)

	// formatError 应该总是包含 "* " 前缀
	result := formatError("test error")
	require.Equal(t, "* test error", result)

	result = formatError("another error")
	require.Equal(t, "* another error", result)
}

func TestFormatFunctions_WithCustomTheme(t *testing.T) {
	// 保存原始主题
	originalTheme := GetTheme()
	defer SetTheme(originalTheme)

	// 创建自定义主题
	customTheme := Theme{
		EnableColor: true,
		TitleStyle: lipgloss.NewStyle().Foreground(lipgloss.Color("red")),
		AnswerStyle: lipgloss.NewStyle().Foreground(lipgloss.Color("green")),
		ErrorStyle:  lipgloss.NewStyle().Foreground(lipgloss.Color("yellow")),
		HintStyle:   lipgloss.NewStyle().Foreground(lipgloss.Color("blue")),
	}
	SetTheme(customTheme)

	// 测试所有格式化函数都使用自定义主题
	titleResult := formatTitle("title")
	require.Contains(t, titleResult, "title")

	answerResult := formatAnswer("answer")
	require.Contains(t, answerResult, "answer")

	errorResult := formatError("error")
	require.Contains(t, errorResult, "* error")

	hintResult := formatHint("hint")
	require.Contains(t, hintResult, "hint")
}


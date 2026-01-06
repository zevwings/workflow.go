package prompt

import (
	"fmt"
	"sync"

	"github.com/charmbracelet/lipgloss"
)

// Theme 表示交互式 UI 的主题配置（颜色、前缀、样式等）
type Theme struct {
	// 基础样式（使用 lipgloss.Style）
	InfoStyle   lipgloss.Style
	WarnStyle   lipgloss.Style
	ErrorStyle  lipgloss.Style
	PromptStyle lipgloss.Style
	AnswerStyle lipgloss.Style
	HintStyle   lipgloss.Style // 提示信息样式（如操作说明）

	// 前缀符号
	PrefixInfo  string // e.g. "INFO"
	PrefixWarn  string // e.g. "WARN"
	PrefixError string // e.g. "ERROR"

	// 输入提示样式
	InputBracketLeft  string // e.g. "["
	InputBracketRight string // e.g. "]"

	// 是否启用颜色（方便 CI / 非 TTY 环境关闭）
	EnableColor bool
}

var (
	// defaultTheme 默认主题（使用 lipgloss）
	defaultTheme = Theme{
		InfoStyle: lipgloss.NewStyle().
			Foreground(lipgloss.Color("51")). // HiCyan
			Bold(false),
		WarnStyle: lipgloss.NewStyle().
			Foreground(lipgloss.Color("226")). // HiYellow
			Bold(false),
		ErrorStyle: lipgloss.NewStyle().
			Foreground(lipgloss.Color("196")). // HiRed
			Bold(true),
		PromptStyle: lipgloss.NewStyle().
			Foreground(lipgloss.Color("51")). // HiCyan
			Bold(false),
		AnswerStyle: lipgloss.NewStyle().
			Foreground(lipgloss.Color("46")). // HiGreen
			Bold(false),
		HintStyle: lipgloss.NewStyle().
			Foreground(lipgloss.Color("240")). // HiBlack (gray)
			Bold(false),

		PrefixInfo:  "",
		PrefixWarn:  "!",
		PrefixError: "x",

		InputBracketLeft:  "[",
		InputBracketRight: "]",

		EnableColor: true,
	}

	// currentTheme 当前主题（使用互斥锁保护，虽然文档说线程不安全版本足够，但为了安全还是加锁）
	currentTheme Theme
	themeMutex   sync.RWMutex
)

func init() {
	// 初始化时使用默认主题
	currentTheme = defaultTheme
}

// SetTheme 设置全局主题（线程安全版本）
func SetTheme(t Theme) {
	themeMutex.Lock()
	defer themeMutex.Unlock()
	currentTheme = t
}

// GetTheme 获取当前主题（只读，线程安全）
func GetTheme() Theme {
	themeMutex.RLock()
	defer themeMutex.RUnlock()
	return currentTheme
}

// formatPrompt 格式化提示消息
func formatPrompt(message string) string {
	t := GetTheme()
	if !t.EnableColor {
		return message
	}
	return t.PromptStyle.Render(message)
}

// formatAnswer 格式化答案显示
func formatAnswer(value string) string {
	t := GetTheme()
	if !t.EnableColor {
		return value
	}
	return t.AnswerStyle.Render(value)
}

// formatError 格式化错误消息
func formatError(message string) string {
	t := GetTheme()
	// 统一错误提示文案格式为: "* 错误提示"
	formatted := fmt.Sprintf("* %s", message)
	if !t.EnableColor {
		return formatted
	}
	return t.ErrorStyle.Render(formatted)
}

// formatHint 格式化提示信息（如操作说明）
func formatHint(message string) string {
	t := GetTheme()
	if !t.EnableColor {
		return message
	}
	return t.HintStyle.Render(message)
}

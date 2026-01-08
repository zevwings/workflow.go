package prompt

import (
	"fmt"
	"os"

	"github.com/charmbracelet/lipgloss"
)

// Message 消息输出工具
type Message struct {
	verbose bool
}

// NewMessage 创建新的消息输出工具
func NewMessage(verbose bool) *Message {
	return &Message{verbose: verbose}
}

// Info 输出信息
func (m *Message) Info(format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	formatted := formatMessage("ℹ", msg, GetTheme().InfoStyle)
	fmt.Println(formatted)
}

// Success 输出成功信息
func (m *Message) Success(format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	formatted := formatMessage("✓", msg, GetTheme().SuccessStyle)
	fmt.Println(formatted)
}

// Warning 输出警告信息
func (m *Message) Warning(format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	t := GetTheme()
	prefix := "⚠"
	if t.PrefixWarn != "" {
		prefix = t.PrefixWarn
	}
	formatted := formatMessage(prefix, msg, t.WarnStyle)
	fmt.Println(formatted)
}

// Error 输出错误信息
func (m *Message) Error(format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	t := GetTheme()
	prefix := "✗"
	if t.PrefixError != "" {
		prefix = t.PrefixError
	}
	formatted := formatMessage(prefix, msg, t.ErrorStyle)
	fmt.Println(formatted)
}

// Fatal 输出致命错误并退出
func (m *Message) Fatal(format string, args ...interface{}) {
	m.Error(format, args...)
	os.Exit(1)
}

// Debug 输出调试信息
func (m *Message) Debug(format string, args ...interface{}) {
	if m.verbose {
		msg := fmt.Sprintf(format, args...)
		formatted := formatMessage("DEBUG:", msg, GetTheme().DebugStyle)
		fmt.Println(formatted)
	}
}

// Print 普通输出
func (m *Message) Print(format string, args ...interface{}) {
	fmt.Printf(format, args...)
}

// Println 普通输出并换行
func (m *Message) Println(format string, args ...interface{}) {
	fmt.Println(fmt.Sprintf(format, args...))
}

// formatMessage 格式化消息（使用主题样式）
func formatMessage(prefix, message string, style lipgloss.Style) string {
	t := GetTheme()
	// 拼接前缀和消息，确保前缀和消息之间有一个空格
	var text string
	if prefix != "" && message != "" {
		text = prefix + " " + message
	} else {
		text = prefix + message
	}
	if !t.EnableColor {
		return text
	}
	return style.Render(text)
}

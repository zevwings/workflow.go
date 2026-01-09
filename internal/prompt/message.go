package prompt

import (
	"fmt"
	"os"
	"strings"

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

// Print 格式化输出并换行
func (m *Message) Print(format string, args ...interface{}) {
	fmt.Println(fmt.Sprintf(format, args...))
}

// Break 打印分隔线或换行
//
// 根据参数的不同组合，提供多种使用方式：
//   - Break() - 输出换行符
//   - Break('-') - 使用默认分隔符（80个 '-'）
//   - Break('=', 100) - 指定分隔符字符和长度
//   - Break('=', 20, "flutter-api.log") - 在分隔线中间插入文本
//
// 参数:
//   - args: 可变参数，支持以下组合：
//   - 无参数: 输出换行符
//   - 1个参数 (rune): 使用该字符作为分隔符，默认长度80
//   - 2个参数 (rune, int): 使用指定字符和长度
//   - 3个参数 (rune, int, string): 在分隔线中间插入文本
//
// 使用示例:
//
//	msg := prompt.NewMessage(false)
//	// 输出换行符
//	msg.Break()
//
//	// 使用默认分隔符（80个 '-'）
//	msg.Break('-')
//
//	// 指定分隔符字符和长度
//	msg.Break('=', 100)
//
//	// 在分隔线中间插入文本
//	msg.Break('=', 20, "flutter-api.log")
//	// 输出: ===========  flutter-api.log ===========
func (m *Message) Break(args ...interface{}) {
	switch len(args) {
	case 0:
		// 无参数：输出换行符
		fmt.Println()
	case 1:
		// 1个参数：使用该字符作为分隔符，默认长度80
		char, ok := args[0].(rune)
		if !ok {
			// 尝试 string 类型
			if str, ok := args[0].(string); ok && len(str) > 0 {
				char = rune(str[0])
			} else {
				return
			}
		}
		separator := strings.Repeat(string(char), 80)
		fmt.Println(separator)
	case 2:
		// 2个参数：字符和长度
		var char rune
		var length int
		var ok bool

		// 处理第一个参数（字符）
		if c, ok := args[0].(rune); ok {
			char = c
		} else if str, ok := args[0].(string); ok && len(str) > 0 {
			char = rune(str[0])
		} else {
			return
		}

		// 处理第二个参数（长度）
		if length, ok = args[1].(int); !ok {
			return
		}

		if length <= 0 {
			length = 80
		}
		separator := strings.Repeat(string(char), length)
		fmt.Println(separator)
	case 3:
		// 3个参数：字符、长度和文本
		var char rune
		var length int
		var text string
		var ok bool

		// 处理第一个参数（字符）
		if c, ok := args[0].(rune); ok {
			char = c
		} else if str, ok := args[0].(string); ok && len(str) > 0 {
			char = rune(str[0])
		} else {
			return
		}

		// 处理第二个参数（长度）
		if length, ok = args[1].(int); !ok {
			return
		}

		// 处理第三个参数（文本）
		if text, ok = args[2].(string); !ok {
			return
		}

		if length <= 0 {
			length = 80
		}

		// 计算分隔符的长度
		// 格式: {separator}  {text}  {separator}
		// 文本前后各有两个空格，所以需要: textLen + 4
		textLen := len(text)
		separatorLen := (length - textLen - 4) / 2

		// 如果文本过长，至少每边保留 2 个字符
		if separatorLen < 2 {
			separatorLen = 2
		}

		separator := strings.Repeat(string(char), separatorLen)
		line := separator + "  " + text + "  " + separator

		// 如果总长度不够，在末尾补充字符
		if len(line) < length {
			line += strings.Repeat(string(char), length-len(line))
		}

		fmt.Println(line)
	default:
		// 参数过多，忽略
		return
	}
}

// formatMessage 格式化消息（使用主题样式）（私有函数）
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

package input

import (
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/mattn/go-runewidth"
	"golang.org/x/term"
)

// Config 编辑器配置
type Config struct {
	// 格式化函数
	FormatPlaceholder func(text string) string
	FormatError       func(message string) string
	// 主题样式（用于格式化）
	HintStyle   lipgloss.Style
	ErrorStyle  lipgloss.Style
	EnableColor bool
}

// ReadWithPlaceholder 读取输入（支持 placeholder 和实时错误提示）
func ReadWithPlaceholder(promptText string, placeholder string, validator Validator, config Config) (string, error) {
	fd := int(os.Stdin.Fd())
	oldState, err := term.MakeRaw(fd)
	if err != nil {
		return readInputFallback(placeholder)
	}
	defer term.Restore(fd, oldState)

	// 确保光标可见
	fmt.Print("\033[?25h")

	var value []byte
	var buf [1]byte
	hasPlaceholder := placeholder != ""
	placeholderDisplayed := false
	errorLineExists := false
	cursorPos := 0

	// 初始显示：promptText + placeholder
	fmt.Print(promptText)
	if hasPlaceholder {
		formattedPlaceholder := config.FormatPlaceholder(placeholder)
		fmt.Print(formattedPlaceholder)
		placeholderDisplayed = true
		placeholderText := StripAnsiCodes(formattedPlaceholder)
		placeholderWidth := runewidth.StringWidth(placeholderText)
		for i := 0; i < placeholderWidth; i++ {
			fmt.Print("\b")
		}
	}

	for {
		n, err := os.Stdin.Read(buf[:])
		if err != nil || n == 0 {
			term.Restore(fd, oldState)
			if len(value) > 0 {
				fmt.Println()
				return strings.TrimSpace(string(value)), nil
			}
			return "", fmt.Errorf("读取输入失败: %w", err)
		}

		char := buf[0]

		// 处理转义序列（方向键等）
		if char == 0x1b { // ESC
			n2, err2 := os.Stdin.Read(buf[:])
			if err2 != nil || n2 == 0 {
				continue
			}
			if buf[0] == '[' {
				n3, err3 := os.Stdin.Read(buf[:])
				if err3 != nil || n3 == 0 {
					continue
				}
				switch buf[0] {
				case 'C': // 右箭头
					if placeholderDisplayed {
						fmt.Print("\r")
						fmt.Print("\033[K")
						fmt.Print(promptText)
						placeholderDisplayed = false
						cursorPos = 0
					} else if cursorPos < len(value) {
						cursorPos++
						fmt.Print("\r")
						fmt.Print("\033[K")
						fmt.Print(promptText)
						fmt.Print(string(value))
						moveCursorToPosition(promptText, value, cursorPos)
					}
					continue
				case 'D': // 左箭头
					if placeholderDisplayed {
						continue
					}
					if cursorPos > 0 {
						cursorPos--
						fmt.Print("\r")
						fmt.Print("\033[K")
						fmt.Print(promptText)
						fmt.Print(string(value))
						moveCursorToPosition(promptText, value, cursorPos)
					}
					continue
				case 'A', 'B': // 上下箭头 - 忽略
					continue
				default:
					continue
				}
			} else {
				continue
			}
		}

		// 处理回车键（结束输入）
		if char == '\r' || char == '\n' {
			if validator != nil {
				if err := validator(string(value)); err != nil {
					if errorLineExists {
						clearErrorLine()
						fmt.Print("\r")
						fmt.Print("\033[K")
						fmt.Print(promptText)
						if len(value) > 0 {
							fmt.Print(string(value))
							moveCursorToPosition(promptText, value, cursorPos)
						} else if placeholderDisplayed {
							formattedPlaceholder := config.FormatPlaceholder(placeholder)
							fmt.Print(formattedPlaceholder)
							placeholderText := StripAnsiCodes(formattedPlaceholder)
							placeholderWidth := runewidth.StringWidth(placeholderText)
							for i := 0; i < placeholderWidth; i++ {
								fmt.Print("\b")
							}
						}
					}
					showError(err.Error(), promptText, value, cursorPos, config.FormatError)
					errorLineExists = true
					continue
				}
			}
			if errorLineExists {
				clearErrorLine()
				errorLineExists = false
			}
			fmt.Print("\r")
			fmt.Print("\033[K")
			fmt.Print(promptText)
			if len(value) > 0 {
				fmt.Print(string(value))
			} else if placeholderDisplayed {
				fmt.Print(" ")
			}
			fmt.Println()
			break
		}

		// 处理退格键或删除键
		if char == 127 || char == 8 {
			if placeholderDisplayed {
				continue
			}
			if cursorPos > 0 {
				value = append(value[:cursorPos-1], value[cursorPos:]...)
				cursorPos--
			} else if len(value) > 0 {
				value = value[1:]
				cursorPos = 0
			}

			fmt.Print("\r")
			fmt.Print("\033[K")
			fmt.Print(promptText)
			if len(value) > 0 {
				fmt.Print(string(value))
				moveCursorToPosition(promptText, value, cursorPos)
			} else if hasPlaceholder {
				formattedPlaceholder := config.FormatPlaceholder(placeholder)
				fmt.Print(formattedPlaceholder)
				placeholderDisplayed = true
				placeholderText := StripAnsiCodes(formattedPlaceholder)
				placeholderWidth := runewidth.StringWidth(placeholderText)
				for i := 0; i < placeholderWidth; i++ {
					fmt.Print("\b")
				}
			}

			if errorLineExists {
				clearErrorLine()
				fmt.Print("\r")
				fmt.Print("\033[K")
				fmt.Print(promptText)
				if len(value) == 0 && hasPlaceholder {
					formattedPlaceholder := config.FormatPlaceholder(placeholder)
					fmt.Print(formattedPlaceholder)
					placeholderText := StripAnsiCodes(formattedPlaceholder)
					placeholderWidth := runewidth.StringWidth(placeholderText)
					for i := 0; i < placeholderWidth; i++ {
						fmt.Print("\b")
					}
				} else if len(value) > 0 {
					fmt.Print(string(value))
					moveCursorToPosition(promptText, value, cursorPos)
				}
				errorLineExists = false
			}
			if validator != nil {
				if err := validator(string(value)); err != nil {
					showError(err.Error(), promptText, value, cursorPos, config.FormatError)
					errorLineExists = true
				}
			}
			continue
		}

		// 处理 Ctrl+C
		if char == 3 {
			term.Restore(fd, oldState)
			if errorLineExists {
				clearErrorLine()
			}
			fmt.Print("\r")
			fmt.Print("\033[K")
			fmt.Print(promptText)
			fmt.Println()
			return "", fmt.Errorf("用户取消输入")
		}

		// 跳过其他控制字符
		if char < 32 {
			continue
		}

		// 处理普通字符 - 在光标位置插入
		if placeholderDisplayed {
			fmt.Print("\r")
			fmt.Print("\033[K")
			fmt.Print(promptText)
			placeholderDisplayed = false
			cursorPos = 0
		}

		if cursorPos < len(value) {
			value = append(value[:cursorPos], append([]byte{char}, value[cursorPos:]...)...)
		} else {
			value = append(value, char)
		}
		cursorPos++

		fmt.Print("\r")
		fmt.Print("\033[K")
		fmt.Print(promptText)
		fmt.Print(string(value))
		moveCursorToPosition(promptText, value, cursorPos)

		if errorLineExists {
			clearErrorLine()
			fmt.Print("\r")
			fmt.Print("\033[K")
			fmt.Print(promptText)
			if len(value) > 0 {
				fmt.Print(string(value))
			} else if placeholderDisplayed {
				formattedPlaceholder := config.FormatPlaceholder(placeholder)
				fmt.Print(formattedPlaceholder)
				placeholderText := StripAnsiCodes(formattedPlaceholder)
				placeholderWidth := runewidth.StringWidth(placeholderText)
				for i := 0; i < placeholderWidth; i++ {
					fmt.Print("\b")
				}
			}
			if !placeholderDisplayed {
				moveCursorToPosition(promptText, value, cursorPos)
			}
			errorLineExists = false
		}

		if validator != nil {
			if err := validator(string(value)); err != nil {
				showError(err.Error(), promptText, value, cursorPos, config.FormatError)
				errorLineExists = true
			}
		}
	}

	term.Restore(fd, oldState)
	return strings.TrimSpace(string(value)), nil
}

// ReadLineCore 通用单行编辑内核（不支持 placeholder）
// 通过 echo 函数控制显示内容（明文 / 密文）
func ReadLineCore(promptText string, validator Validator, echo func([]byte) string, formatError func(string) string) (string, error) {
	fd := int(os.Stdin.Fd())
	oldState, err := term.MakeRaw(fd)
	if err != nil {
		return readInputFallback("")
	}
	defer term.Restore(fd, oldState)

	var value []byte
	var buf [1]byte
	cursorPos := 0
	errorLineExists := false

	fmt.Print(promptText)
	fmt.Print("\033[?25h")

	for {
		n, err := os.Stdin.Read(buf[:])
		if err != nil || n == 0 {
			term.Restore(fd, oldState)
			if len(value) > 0 {
				fmt.Println()
				return strings.TrimSpace(string(value)), nil
			}
			return "", fmt.Errorf("读取输入失败: %w", err)
		}

		char := buf[0]

		// 处理转义序列（方向键等）
		if char == 0x1b {
			n2, err2 := os.Stdin.Read(buf[:])
			if err2 != nil || n2 == 0 {
				continue
			}
			if buf[0] == '[' {
				n3, err3 := os.Stdin.Read(buf[:])
				if err3 != nil || n3 == 0 {
					continue
				}
				switch buf[0] {
				case 'C': // 右箭头
					if cursorPos < len(value) {
						cursorPos++
						fmt.Print("\r")
						fmt.Print("\033[K")
						fmt.Print(promptText)
						display := echo(value)
						fmt.Print(display)
						moveCursorToPosition(promptText, []byte(display), cursorPos)
					}
					continue
				case 'D': // 左箭头
					if cursorPos > 0 {
						cursorPos--
						fmt.Print("\r")
						fmt.Print("\033[K")
						fmt.Print(promptText)
						display := echo(value)
						fmt.Print(display)
						moveCursorToPosition(promptText, []byte(display), cursorPos)
					}
					continue
				case 'A', 'B': // 上下箭头 - 忽略
					continue
				default:
					continue
				}
			} else {
				continue
			}
		}

		// 处理回车键（结束输入）
		if char == '\r' || char == '\n' {
			if validator != nil {
				if err := validator(string(value)); err != nil {
					if errorLineExists {
						clearErrorLine()
						fmt.Print("\r")
						fmt.Print("\033[K")
						fmt.Print(promptText)
						display := echo(value)
						fmt.Print(display)
						moveCursorToPosition(promptText, []byte(display), cursorPos)
					}
					display := echo(value)
					showError(err.Error(), promptText, []byte(display), cursorPos, formatError)
					errorLineExists = true
					continue
				}
			}
			if errorLineExists {
				clearErrorLine()
				errorLineExists = false
			}
			fmt.Print("\r")
			fmt.Print("\033[K")
			fmt.Print(promptText)
			if len(value) > 0 {
				display := echo(value)
				fmt.Print(display)
			}
			fmt.Println()
			break
		}

		// 处理退格键或删除键
		if char == 127 || char == 8 {
			if cursorPos > 0 {
				value = append(value[:cursorPos-1], value[cursorPos:]...)
				cursorPos--
				fmt.Print("\r")
				fmt.Print("\033[K")
				fmt.Print(promptText)
				display := echo(value)
				fmt.Print(display)
				moveCursorToPosition(promptText, []byte(display), cursorPos)
			}
			if errorLineExists {
				clearErrorLine()
				fmt.Print("\r")
				fmt.Print("\033[K")
				fmt.Print(promptText)
				display := echo(value)
				fmt.Print(display)
				moveCursorToPosition(promptText, []byte(display), cursorPos)
				errorLineExists = false
			}
			if validator != nil {
				if err := validator(string(value)); err != nil {
					display := echo(value)
					showError(err.Error(), promptText, []byte(display), cursorPos, formatError)
					errorLineExists = true
				}
			}
			continue
		}

		// 处理 Ctrl+C
		if char == 3 {
			term.Restore(fd, oldState)
			if errorLineExists {
				clearErrorLine()
			}
			fmt.Print("\r")
			fmt.Print("\033[K")
			fmt.Print(promptText)
			fmt.Println()
			return "", fmt.Errorf("用户取消输入")
		}

		// 跳过其他控制字符
		if char < 32 {
			continue
		}

		// 处理普通字符 - 在光标位置插入
		if cursorPos < len(value) {
			value = append(value[:cursorPos], append([]byte{char}, value[cursorPos:]...)...)
		} else {
			value = append(value, char)
		}
		cursorPos++

		fmt.Print("\r")
		fmt.Print("\033[K")
		fmt.Print(promptText)
		display := echo(value)
		fmt.Print(display)
		moveCursorToPosition(promptText, []byte(display), cursorPos)

		if errorLineExists {
			clearErrorLine()
			fmt.Print("\r")
			fmt.Print("\033[K")
			fmt.Print(promptText)
			display := echo(value)
			fmt.Print(display)
			moveCursorToPosition(promptText, []byte(display), cursorPos)
			errorLineExists = false
		}

		if validator != nil {
			if err := validator(string(value)); err != nil {
				display := echo(value)
				showError(err.Error(), promptText, []byte(display), cursorPos, formatError)
				errorLineExists = true
			}
		}
	}

	term.Restore(fd, oldState)
	return strings.TrimSpace(string(value)), nil
}

// moveCursorToPosition 移动光标到指定位置
func moveCursorToPosition(promptText string, value []byte, pos int) {
	if pos < 0 {
		pos = 0
	}
	if pos > len(value) {
		pos = len(value)
	}

	if pos < len(value) {
		tail := value[pos:]
		tailWidth := runewidth.StringWidth(string(tail))
		for i := 0; i < tailWidth; i++ {
			fmt.Print("\b")
		}
	}
}

// clearErrorLine 清除错误提示行
func clearErrorLine() {
	fmt.Print("\n")
	fmt.Print("\r")
	fmt.Print("\033[K")
	fmt.Print("\033[A")
}

// showError 显示错误提示
func showError(message string, promptText string, value []byte, cursorPos int, formatError func(string) string) {
	fmt.Print("\n")
	fmt.Print("\r")
	fmt.Print(formatError(message))
	fmt.Print("\033[A")
	fmt.Print("\r")
	fmt.Print("\033[K")
	fmt.Print(promptText)
	if len(value) > 0 {
		fmt.Print(string(value))
	}
	moveCursorToPosition(promptText, value, cursorPos)
}

// readInputFallback 回退方案：如果无法设置原始模式，使用普通输入
func readInputFallback(placeholder string) (string, error) {
	var input string
	_, err := fmt.Scanln(&input)
	if err != nil {
		return "", fmt.Errorf("读取输入失败: %w", err)
	}
	return strings.TrimSpace(input), nil
}

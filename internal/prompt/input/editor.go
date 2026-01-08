package input

import (
	"fmt"
	"strings"

	"github.com/zevwings/workflow/internal/prompt/io"
)

// Config 编辑器配置
type Config struct {
	// 格式化函数
	FormatPlaceholder func(text string) string
	FormatError       func(message string) string
	// 主题样式（用于格式化）
	HintStyle   interface{} // lipgloss.Style，使用 interface{} 避免循环依赖
	ErrorStyle  interface{} // lipgloss.Style
	EnableColor bool
}

// ReadWithPlaceholder 读取输入（支持 placeholder 和实时错误提示）
func ReadWithPlaceholder(promptText string, placeholder string, validator Validator, config Config, terminal io.TerminalIO) (string, error) {
	handler := NewInputHandler(validator)
	placeholderHandler := NewPlaceholderHandler(placeholder, config)

	// 设置终端原始模式
	oldState, err := terminal.MakeRaw()
	if err != nil {
		return readInputFallback(placeholder, terminal)
	}
	defer terminal.Restore(oldState)

	// 确保光标可见
	terminal.ShowCursor()

	var value []byte
	hasPlaceholder := placeholderHandler.HasPlaceholder()
	placeholderDisplayed := false
	errorLineExists := false
	cursorPos := 0

	// 初始显示：promptText + placeholder
	terminal.Print(promptText)
	if hasPlaceholder {
		formattedPlaceholder := placeholderHandler.FormatPlaceholder()
		terminal.Print(formattedPlaceholder)
		placeholderDisplayed = true
		placeholderWidth := placeholderHandler.GetPlaceholderWidth()
		for i := 0; i < placeholderWidth; i++ {
			terminal.Print("\b")
		}
	}

	for {
		char, err := terminal.ReadByte()
		if err != nil {
			terminal.Restore(oldState)
			if len(value) > 0 {
				terminal.Println("")
				return handler.TrimInput(string(value)), nil
			}
			return "", fmt.Errorf("读取输入失败: %w", err)
		}

		// 处理转义序列（方向键等）
		if char == 0x1b { // ESC
			char2, err := terminal.ReadByte()
			if err != nil {
				continue
			}
			if char2 == '[' {
				char3, err := terminal.ReadByte()
				if err != nil {
					continue
				}
				direction, shouldUpdate := handler.ProcessEscapeSequence(char2, char3)
				if shouldUpdate {
					newPos, shouldClear := handler.ProcessArrowKey(cursorPos, len(value), direction, placeholderDisplayed)
					if shouldClear {
						terminal.MoveToStart()
						terminal.ClearLine()
						terminal.Print(promptText)
						placeholderDisplayed = false
						cursorPos = 0
					} else if newPos != cursorPos {
						cursorPos = newPos
						terminal.MoveToStart()
						terminal.ClearLine()
						terminal.Print(promptText)
						terminal.Print(string(value))
						backspaces := handler.CalculateCursorBackspaces(value, cursorPos)
						for i := 0; i < backspaces; i++ {
							terminal.Print("\b")
						}
					}
				}
				continue
			}
			continue
		}

		// 处理回车键（结束输入）
		if char == '\r' || char == '\n' {
			if err := handler.ValidateInput(string(value)); err != nil {
				if errorLineExists {
					clearErrorLine(terminal)
					terminal.MoveToStart()
					terminal.ClearLine()
					terminal.Print(promptText)
					if len(value) > 0 {
						terminal.Print(string(value))
						backspaces := handler.CalculateCursorBackspaces(value, cursorPos)
						for i := 0; i < backspaces; i++ {
							terminal.Print("\b")
						}
					} else if placeholderDisplayed {
						formattedPlaceholder := placeholderHandler.FormatPlaceholder()
						terminal.Print(formattedPlaceholder)
						placeholderWidth := placeholderHandler.GetPlaceholderWidth()
						for i := 0; i < placeholderWidth; i++ {
							terminal.Print("\b")
						}
					}
				}
				showError(err.Error(), promptText, value, cursorPos, config.FormatError, terminal)
				errorLineExists = true
				continue
			}
			if errorLineExists {
				clearErrorLine(terminal)
				errorLineExists = false
			}
			terminal.MoveToStart()
			terminal.ClearLine()
			terminal.Print(promptText)
			if len(value) > 0 {
				terminal.Print(string(value))
			} else if placeholderDisplayed {
				formattedPlaceholder := placeholderHandler.FormatPlaceholder()
				terminal.Print(formattedPlaceholder)
			}
			terminal.Println("")
			terminal.Restore(oldState)
			return handler.TrimInput(string(value)), nil
		}

		// 处理退格键或删除键
		if char == 127 || char == 8 {
			if cursorPos > 0 {
				value, cursorPos = handler.ProcessBackspace(value, cursorPos)
				terminal.MoveToStart()
				terminal.ClearLine()
				terminal.Print(promptText)
				terminal.Print(string(value))
				backspaces := handler.CalculateCursorBackspaces(value, cursorPos)
				for i := 0; i < backspaces; i++ {
					terminal.Print("\b")
				}
			}
			if errorLineExists {
				clearErrorLine(terminal)
				terminal.MoveToStart()
				terminal.ClearLine()
				terminal.Print(promptText)
				if len(value) > 0 {
					terminal.Print(string(value))
					backspaces := handler.CalculateCursorBackspaces(value, cursorPos)
					for i := 0; i < backspaces; i++ {
						terminal.Print("\b")
					}
				} else if placeholderDisplayed {
					formattedPlaceholder := placeholderHandler.FormatPlaceholder()
					terminal.Print(formattedPlaceholder)
					placeholderWidth := placeholderHandler.GetPlaceholderWidth()
					for i := 0; i < placeholderWidth; i++ {
						terminal.Print("\b")
					}
				}
				errorLineExists = false
			}
			if err := handler.ValidateInput(string(value)); err != nil {
				showError(err.Error(), promptText, value, cursorPos, config.FormatError, terminal)
				errorLineExists = true
			}
			continue
		}

		// 处理 Ctrl+C
		if char == 3 {
			terminal.Restore(oldState)
			if errorLineExists {
				clearErrorLine(terminal)
			}
			terminal.MoveToStart()
			terminal.ClearLine()
			terminal.Print(promptText)
			terminal.Println("")
			return "", fmt.Errorf("用户取消输入")
		}

		// 跳过其他控制字符
		if char < 32 {
			continue
		}

		// 处理普通字符 - 在光标位置插入
		if placeholderDisplayed {
			terminal.MoveToStart()
			terminal.ClearLine()
			terminal.Print(promptText)
			placeholderDisplayed = false
			cursorPos = 0
		}

		value, cursorPos = handler.ProcessChar(value, cursorPos, char)

		terminal.MoveToStart()
		terminal.ClearLine()
		terminal.Print(promptText)
		terminal.Print(string(value))
		backspaces := handler.CalculateCursorBackspaces(value, cursorPos)
		for i := 0; i < backspaces; i++ {
			terminal.Print("\b")
		}

		if errorLineExists {
			clearErrorLine(terminal)
			terminal.MoveToStart()
			terminal.ClearLine()
			terminal.Print(promptText)
			if len(value) > 0 {
				terminal.Print(string(value))
				backspaces := handler.CalculateCursorBackspaces(value, cursorPos)
				for i := 0; i < backspaces; i++ {
					terminal.Print("\b")
				}
			} else if placeholderDisplayed {
				formattedPlaceholder := placeholderHandler.FormatPlaceholder()
				terminal.Print(formattedPlaceholder)
				placeholderWidth := placeholderHandler.GetPlaceholderWidth()
				for i := 0; i < placeholderWidth; i++ {
					terminal.Print("\b")
				}
			}
			errorLineExists = false
		}

		if err := handler.ValidateInput(string(value)); err != nil {
			showError(err.Error(), promptText, value, cursorPos, config.FormatError, terminal)
			errorLineExists = true
		}
	}
}

// ReadWithPlaceholderDefault 向后兼容的 ReadWithPlaceholder 函数
func ReadWithPlaceholderDefault(promptText string, placeholder string, validator Validator, config Config) (string, error) {
	return ReadWithPlaceholder(promptText, placeholder, validator, config, io.NewStdTerminal())
}

// ReadLineCore 通用单行编辑内核（不支持 placeholder）
// 通过 echo 函数控制显示内容（明文 / 密文）
func ReadLineCore(promptText string, validator Validator, echo func([]byte) string, formatError func(string) string, terminal io.TerminalIO) (string, error) {
	handler := NewInputHandler(validator)

	// 设置终端原始模式
	oldState, err := terminal.MakeRaw()
	if err != nil {
		return readInputFallback("", terminal)
	}
	defer terminal.Restore(oldState)

	var value []byte
	cursorPos := 0
	errorLineExists := false

	terminal.Print(promptText)
	terminal.ShowCursor()

	for {
		char, err := terminal.ReadByte()
		if err != nil {
			terminal.Restore(oldState)
			if len(value) > 0 {
				terminal.Println("")
				return handler.TrimInput(string(value)), nil
			}
			return "", fmt.Errorf("读取输入失败: %w", err)
		}

		// 处理转义序列（方向键等）
		if char == 0x1b {
			char2, err := terminal.ReadByte()
			if err != nil {
				continue
			}
			if char2 == '[' {
				char3, err := terminal.ReadByte()
				if err != nil {
					continue
				}
				direction, shouldUpdate := handler.ProcessEscapeSequence(char2, char3)
				if shouldUpdate {
					newPos, _ := handler.ProcessArrowKey(cursorPos, len(value), direction, false)
					if newPos != cursorPos {
						cursorPos = newPos
						terminal.MoveToStart()
						terminal.ClearLine()
						terminal.Print(promptText)
						display := echo(value)
						terminal.Print(display)
						backspaces := handler.CalculateCursorBackspaces([]byte(display), cursorPos)
						for i := 0; i < backspaces; i++ {
							terminal.Print("\b")
						}
					}
				}
				continue
			}
			continue
		}

		// 处理回车键（结束输入）
		if char == '\r' || char == '\n' {
			if err := handler.ValidateInput(string(value)); err != nil {
				if errorLineExists {
					clearErrorLine(terminal)
					terminal.MoveToStart()
					terminal.ClearLine()
					terminal.Print(promptText)
					display := echo(value)
					terminal.Print(display)
					backspaces := handler.CalculateCursorBackspaces([]byte(display), cursorPos)
					for i := 0; i < backspaces; i++ {
						terminal.Print("\b")
					}
				}
				display := echo(value)
				showError(err.Error(), promptText, []byte(display), cursorPos, formatError, terminal)
				errorLineExists = true
				continue
			}
			if errorLineExists {
				clearErrorLine(terminal)
				errorLineExists = false
			}
			terminal.MoveToStart()
			terminal.ClearLine()
			terminal.Print(promptText)
			if len(value) > 0 {
				display := echo(value)
				terminal.Print(display)
			}
			terminal.Println("")
			terminal.Restore(oldState)
			return handler.TrimInput(string(value)), nil
		}

		// 处理退格键或删除键
		if char == 127 || char == 8 {
			if cursorPos > 0 {
				value, cursorPos = handler.ProcessBackspace(value, cursorPos)
				terminal.MoveToStart()
				terminal.ClearLine()
				terminal.Print(promptText)
				display := echo(value)
				terminal.Print(display)
				backspaces := handler.CalculateCursorBackspaces([]byte(display), cursorPos)
				for i := 0; i < backspaces; i++ {
					terminal.Print("\b")
				}
			}
			if errorLineExists {
				clearErrorLine(terminal)
				terminal.MoveToStart()
				terminal.ClearLine()
				terminal.Print(promptText)
				display := echo(value)
				terminal.Print(display)
				backspaces := handler.CalculateCursorBackspaces([]byte(display), cursorPos)
				for i := 0; i < backspaces; i++ {
					terminal.Print("\b")
				}
				errorLineExists = false
			}
			if err := handler.ValidateInput(string(value)); err != nil {
				display := echo(value)
				showError(err.Error(), promptText, []byte(display), cursorPos, formatError, terminal)
				errorLineExists = true
			}
			continue
		}

		// 处理 Ctrl+C
		if char == 3 {
			terminal.Restore(oldState)
			if errorLineExists {
				clearErrorLine(terminal)
			}
			terminal.MoveToStart()
			terminal.ClearLine()
			terminal.Print(promptText)
			terminal.Println("")
			return "", fmt.Errorf("用户取消输入")
		}

		// 跳过其他控制字符
		if char < 32 {
			continue
		}

		// 处理普通字符 - 在光标位置插入
		value, cursorPos = handler.ProcessChar(value, cursorPos, char)

		terminal.MoveToStart()
		terminal.ClearLine()
		terminal.Print(promptText)
		display := echo(value)
		terminal.Print(display)
		backspaces := handler.CalculateCursorBackspaces([]byte(display), cursorPos)
		for i := 0; i < backspaces; i++ {
			terminal.Print("\b")
		}

		if errorLineExists {
			clearErrorLine(terminal)
			terminal.MoveToStart()
			terminal.ClearLine()
			terminal.Print(promptText)
			display := echo(value)
			terminal.Print(display)
			backspaces := handler.CalculateCursorBackspaces([]byte(display), cursorPos)
			for i := 0; i < backspaces; i++ {
				terminal.Print("\b")
			}
			errorLineExists = false
		}

		if err := handler.ValidateInput(string(value)); err != nil {
			display := echo(value)
			showError(err.Error(), promptText, []byte(display), cursorPos, formatError, terminal)
			errorLineExists = true
		}
	}
}

// ReadLineCoreDefault 向后兼容的 ReadLineCore 函数
func ReadLineCoreDefault(promptText string, validator Validator, echo func([]byte) string, formatError func(string) string) (string, error) {
	return ReadLineCore(promptText, validator, echo, formatError, io.NewStdTerminal())
}

// moveCursorToPosition 移动光标到指定位置（保留用于向后兼容，但推荐使用 handler 的方法）
func moveCursorToPosition(promptText string, value []byte, pos int) {
	// 这个函数现在只是调用 handler 的逻辑，但为了向后兼容保留
	handler := NewInputHandler(nil)
	backspaces := handler.CalculateCursorBackspaces(value, pos)
	for i := 0; i < backspaces; i++ {
		fmt.Print("\b")
	}
}

// clearErrorLine 清除错误提示行
func clearErrorLine(terminal io.TerminalIO) {
	terminal.Print("\n")
	terminal.MoveToStart()
	terminal.ClearLine()
	terminal.Print("\033[A")
}

// showError 显示错误提示
func showError(message string, promptText string, value []byte, cursorPos int, formatError func(string) string, terminal io.TerminalIO) {
	terminal.Print("\n")
	terminal.MoveToStart()
	terminal.Print(formatError(message))
	terminal.Print("\033[A")
	terminal.MoveToStart()
	terminal.ClearLine()
	terminal.Print(promptText)
	if len(value) > 0 {
		terminal.Print(string(value))
	}
	handler := NewInputHandler(nil)
	backspaces := handler.CalculateCursorBackspaces(value, cursorPos)
	for i := 0; i < backspaces; i++ {
		terminal.Print("\b")
	}
}

// readInputFallback 回退方案：如果无法设置原始模式，使用普通输入
func readInputFallback(placeholder string, terminal io.TerminalIO) (string, error) {
	input, err := terminal.ReadLine()
	if err != nil {
		return "", fmt.Errorf("读取输入失败: %w", err)
	}
	return strings.TrimSpace(input), nil
}

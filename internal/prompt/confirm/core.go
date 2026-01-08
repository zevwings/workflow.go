package confirm

import (
	"fmt"

	"github.com/zevwings/workflow/internal/prompt/common"
	"github.com/zevwings/workflow/internal/prompt/io"
)

// Config 确认功能配置
type Config struct {
	// 格式化函数
	FormatPrompt func(message string) string
	FormatAnswer func(value string) string
	FormatHint   func(message string) string // 格式化提示信息（如 【y/N】）
}

// Confirm 内部实现函数（使用 TerminalIO 接口）
// 这是新的实现，接受 TerminalIO 参数以便测试
func Confirm(message string, defaultYes bool, config Config, terminal io.TerminalIO) (bool, error) {
	handler := NewConfirmHandler(defaultYes, config)

	// 格式化提示文本
	promptText := handler.FormatPromptText(message)

	// 保存光标位置（在提示行的开始）
	terminal.SaveCursor()
	terminal.Print(promptText)

	// 创建原始模式管理器和转义序列解析器
	rawModeMgr := io.NewRawModeManager(terminal)
	parser := io.NewEscapeSequenceParser(terminal)

	// 使用原始模式管理器执行交互逻辑
	var result bool
	var resultErr error

	err := rawModeMgr.WithRawModeAndFallback(
		func() error {
			// 读取单个字符，检测到 y/n 立即返回
			for {
				keyType, char, err := parser.ReadKey()
				if err != nil {
					return fmt.Errorf("读取输入失败: %w", err)
				}

				// 处理 Ctrl+C
				if keyType == io.KeyCtrlC {
					resultErr = common.HandleCancel(terminal)
					return resultErr
				}

				// 处理回车键（KeyEnter 时 char 为 0，需要传递实际字符）
				if keyType == io.KeyEnter {
					char = '\r' // ProcessInput 需要实际的回车字符
				}

				// 处理输入
				processResult, shouldContinue, err := handler.ProcessInput(char)
				if err != nil {
					resultErr = err
					terminal.Println("")
					return err
				}

				if !shouldContinue && processResult != nil {
					// 显示结果
					promptMsg := config.FormatPrompt(message)
					answerText := handler.FormatAnswer(*processResult)
					displayConfirmResult(terminal, promptMsg, answerText)

					result = *processResult
					return nil
				}

				// 继续等待有效输入
			}
		},
		func() error {
			// Fallback: 如果无法设置原始模式，使用普通输入
			fallbackResult, err := confirmFallback(message, defaultYes, config, terminal)
			result = fallbackResult
			return err
		},
	)

	if err != nil {
		// 如果 err 是取消错误，返回 false
		if resultErr != nil {
			return false, resultErr
		}
		return defaultYes, err
	}

	if resultErr != nil {
		return false, resultErr
	}

	return result, nil
}

// ConfirmDefault 向后兼容的 Confirm 函数
// 使用标准终端实现，保持现有 API 不变
func ConfirmDefault(message string, defaultYes bool, config Config) (bool, error) {
	return Confirm(message, defaultYes, config, io.NewStdTerminal())
}

// displayConfirmResult 显示确认结果（统一的结果显示逻辑）
func displayConfirmResult(terminal io.TerminalIO, promptMsg string, answerText string) {
	// 恢复到提示行的开始位置
	terminal.RestoreCursor()
	// 移动到行首（确保光标在行首）
	terminal.MoveToStart()
	// 清除从光标位置到行尾的内容
	terminal.ClearLine()

	// 显示提示消息和结果（在同一行）
	terminal.Print(promptMsg)
	terminal.Print(" ")
	terminal.Print(answerText)
	// 换行
	terminal.Print("\r\n")

	// 重置格式
	terminal.ResetFormat()
}

// confirmFallback 回退方案：如果无法设置原始模式，使用普通输入
func confirmFallback(message string, defaultYes bool, config Config, terminal io.TerminalIO) (bool, error) {
	handler := NewConfirmHandler(defaultYes, config)

	// 格式化提示文本
	promptText := handler.FormatPromptText(message)

	// 保存光标位置（在提示行的开始）
	terminal.SaveCursor()
	terminal.Print(promptText)

	// 读取一行输入
	input, err := terminal.ReadLine()
	if err != nil {
		// 如果读取失败（可能是空输入），返回默认值并显示格式化结果
		promptMsg := config.FormatPrompt(message)
		answerText := handler.FormatAnswer(defaultYes)
		displayConfirmResult(terminal, promptMsg, answerText)

		return defaultYes, nil
	}

	// 处理输入
	processResult, err := handler.ProcessLineInput(input)
	if err != nil {
		return defaultYes, err
	}

	// 显示格式化结果
	promptMsg := config.FormatPrompt(message)
	answerText := handler.FormatAnswer(*processResult)
	displayConfirmResult(terminal, promptMsg, answerText)

	return *processResult, nil
}

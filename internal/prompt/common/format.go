package common

import (
	"github.com/zevwings/workflow/internal/prompt/io"
)

// FormatResult 格式化并显示结果
// 统一的结果格式化显示逻辑，处理光标恢复、清除、格式化输出
//
// 参数:
//   - terminal: 终端接口
//   - promptMsg: 提示消息（已格式化）
//   - resultText: 结果文本
//   - formatAnswer: 格式化答案的函数（可选，如果为 nil 则直接显示 resultText）
//
// 返回:
//   - error: 格式化过程中的错误
func FormatResult(terminal io.TerminalIO, promptMsg string, resultText string, formatAnswer func(string) string) error {
	return FormatResultWithOptions(terminal, promptMsg, resultText, formatAnswer, true)
}

// FormatResultWithOptions 格式化并显示结果（带选项）
// 统一的结果格式化显示逻辑，处理光标恢复、清除、格式化输出
//
// 参数:
//   - terminal: 终端接口
//   - promptMsg: 提示消息（已格式化）
//   - resultText: 结果文本
//   - formatAnswer: 格式化答案的函数（可选，如果为 nil 则直接显示 resultText）
//   - includePrompt: 是否包含提示消息（如果为 false，则只显示结果，不显示提示消息）
//
// 返回:
//   - error: 格式化过程中的错误
func FormatResultWithOptions(terminal io.TerminalIO, promptMsg string, resultText string, formatAnswer func(string) string, includePrompt bool) error {
	if includePrompt {
		// 恢复光标位置并清除内容
		terminal.RestoreCursor()
		terminal.ClearToEnd()

		// 显示提示消息和结果
		terminal.Print(promptMsg)
		terminal.Print(" ")

		// 格式化结果文本
		if formatAnswer != nil {
			resultText = formatAnswer(resultText)
		}
		terminal.Println(resultText)

		// 重置格式
		terminal.ResetFormat()
	} else {
		// 当 includePrompt 为 false 时，需要回到提示消息那一行
		// RenderWithPrompt 输出了提示消息后换行了两次，所以需要向上移动两行
		// 恢复光标位置（在选项列表之前）
		terminal.RestoreCursor()
		// 向上移动两行（回到提示消息那一行）
		// ANSI 转义序列: \033[2A 向上移动两行
		terminal.Print("\033[2A")
		// 移动到行首并清除当前行
		terminal.MoveToStart()
		terminal.ClearLine()

		// 输出"提示消息 + 结果"
		terminal.Print(promptMsg)
		terminal.Print(" ")

		// 格式化结果文本
		if formatAnswer != nil {
			resultText = formatAnswer(resultText)
		}
		terminal.Print(resultText)
		// 清除从光标到屏幕底部的内容（清除选项列表和提示信息）
		terminal.ClearToEnd()
		// 换行
		terminal.Print("\r\n")

		// 重置格式
		terminal.ResetFormat()
	}

	return nil
}


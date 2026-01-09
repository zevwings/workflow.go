package common

import (
	"github.com/zevwings/workflow/internal/prompt/io"
)

// formatResultTitleHelper 格式化结果标题的辅助函数
// 如果 formatFn 为 nil，返回原始的 promptMsg
// 否则使用 formatFn 函数格式化
func formatResultTitleHelper(promptMsg, originalMessage, resultText string, formatFn func(string, string) string) string {
	if formatFn == nil {
		return promptMsg
	}
	originalMsg := originalMessage
	if originalMsg == "" {
		originalMsg = promptMsg
	}
	return formatFn(originalMsg, resultText)
}

// formatAnswerPrefixHelper 格式化答案前缀的辅助函数
// 如果 formatFn 为 nil，返回默认的 "> "
// 否则使用 formatFn 函数格式化
func formatAnswerPrefixHelper(formatFn func() string) string {
	if formatFn == nil {
		return "> "
	}
	return formatFn()
}

// formatAnswerText 格式化答案文本的辅助函数
// 如果 formatAnswer 为 nil，返回原始的 resultText
// 否则使用 formatAnswer 函数格式化
func formatAnswerText(resultText string, formatAnswer func(string) string) string {
	if formatAnswer != nil {
		return formatAnswer(resultText)
	}
	return resultText
}

// FormatPromptWithPrefix 格式化提示消息并添加前缀
// 统一的前缀处理逻辑，用于 select、multiselect 等交互式提示
//
// 参数:
//   - promptMsg: 已格式化的提示消息
//   - config: 提示配置（用于获取 FormatQuestionPrefix）
//
// 返回:
//   - promptMsgWithPrefix: 带前缀的提示消息
func FormatPromptWithPrefix(promptMsg string, config PromptConfig) string {
	if config.FormatQuestionPrefix != nil {
		return config.FormatQuestionPrefix() + promptMsg
	}
	return "? " + promptMsg
}

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
//   - formatAnswerPrefix: 格式化答案前缀 "> " 的函数（可选）
//   - formatResultTitle: 格式化完成后显示的 title 的函数（可选，参数为 originalMessage 和 resultValue）
//   - originalMessage: 原始的提示消息（用于 FormatResultTitle，如果为 "" 则使用 promptMsg）
//
// 返回:
//   - error: 格式化过程中的错误
func FormatResultWithOptions(terminal io.TerminalIO, promptMsg string, resultText string, formatAnswer func(string) string, includePrompt bool, formatAnswerPrefix ...func() string) error {
	return FormatResultWithTitle(terminal, promptMsg, resultText, formatAnswer, includePrompt, "", nil, formatAnswerPrefix...)
}

// FormatResultWithTitle 格式化并显示结果（支持动态 title）
// 统一的结果格式化显示逻辑，处理光标恢复、清除、格式化输出
//
// 参数:
//   - terminal: 终端接口
//   - promptMsg: 提示消息（已格式化）
//   - resultText: 结果文本
//   - formatAnswer: 格式化答案的函数（可选，如果为 nil 则直接显示 resultText）
//   - includePrompt: 是否包含提示消息（如果为 false，则只显示结果，不显示提示消息）
//   - originalMessage: 原始的提示消息（用于 FormatResultTitle，如果为 "" 则使用 promptMsg）
//   - formatResultTitle: 格式化完成后显示的 title 的函数（可选，参数为 originalMessage 和 resultValue）
//   - formatAnswerPrefix: 格式化答案前缀 "> " 的函数（可选）
//
// 返回:
//   - error: 格式化过程中的错误
func FormatResultWithTitle(terminal io.TerminalIO, promptMsg string, resultText string, formatAnswer func(string) string, includePrompt bool, originalMessage string, formatResultTitle func(string, string) string, formatAnswerPrefix ...func() string) error {
	// 提取 formatAnswerPrefix 函数（处理可变参数）
	var answerPrefixFn func() string
	if len(formatAnswerPrefix) > 0 {
		answerPrefixFn = formatAnswerPrefix[0]
	}

	if includePrompt {
		// 恢复光标位置并清除内容
		terminal.RestoreCursor()
		terminal.ClearToEnd()

		// 格式化并显示标题和结果
		displayTitle := formatResultTitleHelper(promptMsg, originalMessage, resultText, formatResultTitle)
		terminal.Print(displayTitle)
		terminal.Print(" ")

		formattedResult := formatAnswerText(resultText, formatAnswer)
		terminal.Println(formattedResult)

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

		// 输出"> 提示消息 + 结果"（已输入时使用 "> " 前缀）
		answerPrefix := formatAnswerPrefixHelper(answerPrefixFn)
		terminal.Print(answerPrefix)

		displayTitle := formatResultTitleHelper(promptMsg, originalMessage, resultText, formatResultTitle)
		terminal.Print(displayTitle)
		terminal.Print(" ")

		formattedResult := formatAnswerText(resultText, formatAnswer)
		terminal.Print(formattedResult)
		// 清除从光标到屏幕底部的内容（清除选项列表和提示信息）
		terminal.ClearToEnd()
		// 换行
		terminal.Print("\r\n")

		// 重置格式
		terminal.ResetFormat()
	}

	return nil
}

// FormatResultInline 格式化并显示结果（在同一行，用于 confirm 等场景）
// 统一的结果格式化显示逻辑，处理光标恢复、清除、格式化输出
// 与 FormatResultWithTitle 的区别是：不需要向上移动两行，因为 confirm 不使用 RenderWithPrompt
//
// 参数:
//   - terminal: 终端接口
//   - promptMsg: 提示消息（已格式化）
//   - resultText: 结果文本
//   - formatAnswer: 格式化答案的函数（可选，如果为 nil 则直接显示 resultText）
//   - originalMessage: 原始的提示消息（用于 FormatResultTitle，如果为 "" 则使用 promptMsg）
//   - formatResultTitle: 格式化完成后显示的 title 的函数（可选，参数为 originalMessage 和 resultValue）
//   - formatAnswerPrefix: 格式化答案前缀 "> " 的函数（可选）
//
// 返回:
//   - error: 格式化过程中的错误
func FormatResultInline(terminal io.TerminalIO, promptMsg string, resultText string, formatAnswer func(string) string, originalMessage string, formatResultTitle func(string, string) string, formatAnswerPrefix func() string) error {
	// 恢复到提示行的开始位置
	terminal.RestoreCursor()
	// 移动到行首（确保光标在行首）
	terminal.MoveToStart()
	// 清除从光标位置到行尾的内容
	terminal.ClearLine()

	// 显示提示消息和结果（在同一行），已输入时使用 "> " 前缀
	answerPrefix := formatAnswerPrefixHelper(formatAnswerPrefix)
	terminal.Print(answerPrefix)

	displayTitle := formatResultTitleHelper(promptMsg, originalMessage, resultText, formatResultTitle)
	terminal.Print(displayTitle)
	terminal.Print(" ")

	formattedResult := formatAnswerText(resultText, formatAnswer)
	terminal.Print(formattedResult)
	// 换行
	terminal.Print("\r\n")

	// 重置格式
	terminal.ResetFormat()

	return nil
}

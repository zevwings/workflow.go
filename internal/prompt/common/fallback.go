package common

import (
	"github.com/zevwings/workflow/internal/prompt/io"
)

// TypedFallbackHandler 类型安全的 fallback 处理器接口（泛型版本）
// 用于提供类型安全的 fallback 处理，避免类型断言
type TypedFallbackHandler[T any] interface {
	// FormatPromptText 格式化提示文本（用于显示）
	FormatPromptText(message string) string

	// FormatAnswer 格式化答案文本（用于显示结果）
	FormatAnswer(result T) string

	// ProcessLineInput 处理一行输入（用于 fallback 模式）
	// 返回处理结果和错误
	ProcessLineInput(input string) (T, error)

	// GetDefaultResult 获取默认结果（当输入为空或无效时使用）
	GetDefaultResult() T
}

// ExecuteFallbackTyped 执行 fallback 模式的通用框架（类型安全版本）
// 使用泛型提供类型安全，避免类型断言
//
// 参数:
//   - terminal: 终端接口
//   - message: 原始提示消息
//   - config: 提示配置
//   - handler: 类型安全的 fallback 处理器
//   - options: fallback 选项
//
// 返回:
//   - result: 处理结果（类型 T）
//   - error: 错误
func ExecuteFallbackTyped[T any](
	terminal io.TerminalIO,
	message string,
	config PromptConfig,
	handler TypedFallbackHandler[T],
	options FallbackOptionsTyped[T],
) (T, error) {
	var zero T

	// 格式化提示文本
	promptText := handler.FormatPromptText(message)

	// 保存光标位置（在提示行的开始）
	terminal.SaveCursor()
	terminal.Print(promptText)

	// 如果设置了显示选项，显示选项列表
	if options.ShowOptions && options.FormatOptions != nil {
		if err := options.FormatOptions(terminal); err != nil {
			return handler.GetDefaultResult(), err
		}
	}

	// 显示输入提示（如果有）
	if options.InputPrompt != "" {
		terminal.Print(options.InputPrompt)
	}

	// 读取一行输入
	input, err := terminal.ReadLine()
	if err != nil {
		// 如果读取失败（可能是空输入），返回默认值并显示格式化结果
		defaultResult := handler.GetDefaultResult()
		if options.ResultDisplay != nil {
			promptMsg := config.FormatPrompt(message)
			if err := options.ResultDisplay(terminal, promptMsg, defaultResult, handler, message, config); err != nil {
				return zero, err
			}
		}
		return defaultResult, nil
	}

	// 处理输入
	result, err := handler.ProcessLineInput(input)
	if err != nil {
		// 处理失败，返回默认值
		defaultResult := handler.GetDefaultResult()
		if options.ResultDisplay != nil {
			promptMsg := config.FormatPrompt(message)
			if err := options.ResultDisplay(terminal, promptMsg, defaultResult, handler, message, config); err != nil {
				return zero, err
			}
		}
		return zero, err
	}

	// 显示格式化结果
	if options.ResultDisplay != nil {
		promptMsg := config.FormatPrompt(message)
		if err := options.ResultDisplay(terminal, promptMsg, result, handler, message, config); err != nil {
			return zero, err
		}
	}

	return result, nil
}

// FallbackOptionsTyped 类型安全的 fallback 执行选项（泛型版本）
type FallbackOptionsTyped[T any] struct {
	// ShowOptions 是否显示选项列表（用于 select/multiselect）
	ShowOptions bool
	// FormatOptions 格式化选项列表的函数（如果 ShowOptions 为 true，必须提供）
	FormatOptions func(terminal io.TerminalIO) error
	// InputPrompt 输入提示文本（如 "请选择 (1-3): "）
	InputPrompt string
	// ResultDisplay 结果显示函数（用于显示最终结果）
	// 参数: terminal, promptMsg, result, handler, originalMessage, config
	ResultDisplay func(terminal io.TerminalIO, promptMsg string, result T, handler TypedFallbackHandler[T], originalMessage string, config PromptConfig) error
}

package confirm

import (
	"fmt"

	"github.com/zevwings/workflow/internal/prompt/common"
	"github.com/zevwings/workflow/internal/prompt/io"
)

// ConfirmConfig 确认功能配置
type ConfirmConfig struct {
	common.BasePromptConfig
	// DefaultYes 默认值（true 表示默认 Yes）
	DefaultYes bool
}

// Confirm 内部实现函数（使用配置结构体）
func Confirm(cfg ConfirmConfig) (bool, error) {
	handler := NewConfirmHandler(cfg.DefaultYes, cfg.Config)

	// 格式化提示文本，添加 "? " 前缀
	promptText := handler.FormatPromptText(cfg.Message)
	promptText = common.FormatPromptWithPrefix(promptText, cfg.Config)

	// 保存光标位置（在提示行的开始）
	cfg.Terminal.SaveCursor()
	cfg.Terminal.Print(promptText)

	// 创建原始模式管理器和转义序列解析器
	rawModeMgr := io.NewRawModeManager(cfg.Terminal)
	parser := io.NewEscapeSequenceParser(cfg.Terminal)

	// 使用原始模式管理器执行交互逻辑
	var result bool

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
					cfg.Terminal.Println("")
					return common.HandleCancel(cfg.Terminal)
				}

				// 处理回车键（KeyEnter 时 char 为 0，需要传递实际字符）
				if keyType == io.KeyEnter {
					char = '\r' // ProcessInput 需要实际的回车字符
				}

				// 处理输入
				processResult, shouldContinue, err := handler.ProcessInput(char)
				if err != nil {
					cfg.Terminal.Println("")
					return err
				}

				if !shouldContinue && processResult != nil {
					// 显示结果
					promptMsg := cfg.Config.FormatPrompt(cfg.Message)
					answerText := handler.FormatAnswer(*processResult)
					if err := common.FormatResultInline(cfg.Terminal, promptMsg, answerText, nil, cfg.Message, cfg.Config.FormatResultTitle, cfg.Config.FormatAnswerPrefix); err != nil {
						return err
					}

					result = *processResult
					return nil
				}

				// 继续等待有效输入
			}
		},
		func() error {
			// Fallback: 如果无法设置原始模式，使用普通输入
			fallbackResult, err := confirmFallback(cfg)
			result = fallbackResult
			return err
		},
	)

	if err != nil {
		// 如果 err 是取消错误，返回 false 和错误
		// 否则返回默认值
		return cfg.DefaultYes, err
	}

	return result, nil
}

// confirmFallback 回退方案：如果无法设置原始模式，使用普通输入
func confirmFallback(cfg ConfirmConfig) (bool, error) {
	handler := NewConfirmHandler(cfg.DefaultYes, cfg.Config)
	adapter := newConfirmFallbackAdapter(handler)

	// 使用类型安全的版本
	result, err := common.ExecuteFallbackTyped(
		cfg.Terminal,
		cfg.Message,
		cfg.Config,
		adapter,
		common.FallbackOptionsTyped[bool]{
			ShowOptions:   false,
			FormatOptions: nil,
			InputPrompt:   "",
			ResultDisplay: func(terminal io.TerminalIO, promptMsg string, result bool, handler common.TypedFallbackHandler[bool], originalMessage string, config common.PromptConfig) error {
				answerText := handler.FormatAnswer(result)
				return common.FormatResultInline(terminal, promptMsg, answerText, nil, originalMessage, config.FormatResultTitle, config.FormatAnswerPrefix)
			},
		},
	)
	if err != nil {
		return cfg.DefaultYes, err
	}

	return result, nil
}

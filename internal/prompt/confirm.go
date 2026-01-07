package prompt

import (
	"github.com/zevwings/workflow/internal/prompt/confirm"
)

// AskConfirm 函数式调用 Confirm（保持向后兼容）
// 使用方式：prompt.AskConfirm("消息", defaultYes)
func AskConfirm(message string, defaultYes bool) (bool, error) {
	return confirmFunc(message, defaultYes)
}

// confirmFunc 统一的确认函数
func confirmFunc(message string, defaultYes bool) (bool, error) {
	// 构建配置
	config := confirm.Config{
		FormatPrompt: formatPrompt,
		FormatAnswer: formatAnswer,
	}

	return confirm.Confirm(message, defaultYes, config)
}

// ConfirmBuilder Confirm 的链式构建器
type ConfirmBuilder struct {
	message    string
	defaultYes bool
}

// Confirm 创建一个 Confirm 构建器（链式调用）
// 使用方式：prompt.Confirm().Prompt("消息").Default(true).Run()
func Confirm() *ConfirmBuilder {
	return &ConfirmBuilder{
		defaultYes: false,
	}
}

// Prompt 设置提示消息
func (b *ConfirmBuilder) Prompt(message string) *ConfirmBuilder {
	b.message = message
	return b
}

// Default 设置默认值
func (b *ConfirmBuilder) Default(defaultYes bool) *ConfirmBuilder {
	b.defaultYes = defaultYes
	return b
}

// Run 执行确认并返回结果
func (b *ConfirmBuilder) Run() (bool, error) {
	return confirmFunc(b.message, b.defaultYes)
}

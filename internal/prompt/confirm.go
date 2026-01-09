package prompt

import (
	"github.com/zevwings/workflow/internal/prompt/common"
	"github.com/zevwings/workflow/internal/prompt/confirm"
	"github.com/zevwings/workflow/internal/prompt/io"
)

// ConfirmField 确认字段配置
type ConfirmField struct {
	// Message 提示消息
	Message string
	// DefaultYes 默认值（true 表示默认 Yes）
	DefaultYes bool
	// ResultTitle 确认完成后显示的 title（可选）
	// 如果设置，将优先于全局的 FormatResultTitle 使用
	ResultTitle string
}

// AskConfirm 使用配置结构体的确认函数
func AskConfirm(field ConfirmField) (bool, error) {
	config := newDefaultConfig()
	// 如果提供了 ResultTitle，设置 FormatResultTitle
	if field.ResultTitle != "" {
		titleStr := field.ResultTitle
		config.FormatResultTitle = func(originalMessage string, resultValue string) string {
			return titleStr
		}
	}
	return confirm.Confirm(confirm.ConfirmConfig{
		BasePromptConfig: common.BasePromptConfig{
			Message:  field.Message,
			Config:   config,
			Terminal: io.NewStdTerminal(),
		},
		DefaultYes: field.DefaultYes,
	})
}

// ConfirmBuilder Confirm 的链式构建器
type ConfirmBuilder struct {
	baseBuilder
	defaultYes bool
}

// Confirm 创建一个 Confirm 构建器（链式调用）
// 使用方式：prompt.Confirm().Prompt("消息").Default(true).Run()
func Confirm() *ConfirmBuilder {
	return &ConfirmBuilder{
		defaultYes: false,
	}
}

// Prompt 设置提示消息（覆盖基类方法以返回正确的类型）
func (b *ConfirmBuilder) Prompt(message string) *ConfirmBuilder {
	b.baseBuilder.Prompt(message)
	return b
}

// Default 设置默认值
func (b *ConfirmBuilder) Default(defaultYes bool) *ConfirmBuilder {
	b.defaultYes = defaultYes
	return b
}

// Run 执行确认并返回结果
func (b *ConfirmBuilder) Run() (bool, error) {
	return AskConfirm(ConfirmField{
		Message:    b.GetMessage(),
		DefaultYes: b.defaultYes,
	})
}

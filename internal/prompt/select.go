package prompt

import (
	selectpkg "github.com/zevwings/workflow/internal/prompt/select"
)

// AskSelect 函数式调用 Select（保持向后兼容）
// 使用方式：prompt.AskSelect("消息", options, defaultIndex)
func AskSelect(message string, options []string, defaultIndex int) (int, error) {
	return selectFunc(message, options, defaultIndex)
}

// selectFunc 统一的选择函数
func selectFunc(message string, options []string, defaultIndex int) (int, error) {
	// 构建配置
	config := selectpkg.Config{
		FormatPrompt: formatPrompt,
		FormatAnswer: formatAnswer,
		FormatHint:   formatHint,
	}

	return selectpkg.Select(message, options, defaultIndex, config)
}

// SelectBuilder Select 的链式构建器
type SelectBuilder struct {
	message      string
	options      []string
	defaultIndex int
}

// Select 创建一个 Select 构建器（链式调用）
// 使用方式：prompt.Select().Prompt("消息").Options(options).Default(0).Run()
func Select() *SelectBuilder {
	return &SelectBuilder{
		defaultIndex: 0,
	}
}

// Prompt 设置提示消息
func (b *SelectBuilder) Prompt(message string) *SelectBuilder {
	b.message = message
	return b
}

// Options 设置选项列表
func (b *SelectBuilder) Options(options []string) *SelectBuilder {
	b.options = options
	return b
}

// Default 设置默认选中的索引
func (b *SelectBuilder) Default(index int) *SelectBuilder {
	b.defaultIndex = index
	return b
}

// Run 执行选择并返回结果（返回选中的索引）
func (b *SelectBuilder) Run() (int, error) {
	return selectFunc(b.message, b.options, b.defaultIndex)
}

package prompt

import (
	"github.com/zevwings/workflow/internal/prompt/common"
	"github.com/zevwings/workflow/internal/prompt/io"
	selectpkg "github.com/zevwings/workflow/internal/prompt/select"
)

// SelectField 单选字段配置
type SelectField struct {
	// Message 提示消息
	Message string
	// Options 选项列表
	Options []string
	// DefaultIndex 默认选中的索引
	DefaultIndex int
	// ResultTitle 选择完成后显示的 title（可选）
	// 如果设置，将优先于全局的 FormatResultTitle 使用
	ResultTitle string
}

// AskSelect 使用配置结构体的选择函数
func AskSelect(field SelectField) (int, error) {
	config := newDefaultConfig()
	// 如果提供了 ResultTitle，设置 FormatResultTitle
	if field.ResultTitle != "" {
		titleStr := field.ResultTitle
		config.FormatResultTitle = func(originalMessage string, resultValue string) string {
			return titleStr
		}
	}
	return selectpkg.Select(selectpkg.SelectConfig{
		BasePromptConfig: common.BasePromptConfig{
			Message:  field.Message,
			Config:   config,
			Terminal: io.NewStdTerminal(),
		},
		Options:      field.Options,
		DefaultIndex: field.DefaultIndex,
	})
}

// SelectBuilder Select 的链式构建器
type SelectBuilder struct {
	baseBuilder
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

// Prompt 设置提示消息（覆盖基类方法以返回正确的类型）
func (b *SelectBuilder) Prompt(message string) *SelectBuilder {
	b.baseBuilder.Prompt(message)
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
	return AskSelect(SelectField{
		Message:      b.GetMessage(),
		Options:      b.options,
		DefaultIndex: b.defaultIndex,
	})
}

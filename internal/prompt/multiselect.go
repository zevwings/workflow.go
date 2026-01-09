package prompt

import (
	"github.com/zevwings/workflow/internal/prompt/common"
	"github.com/zevwings/workflow/internal/prompt/io"
	"github.com/zevwings/workflow/internal/prompt/multiselect"
)

// MultiSelectField 多选字段配置
type MultiSelectField struct {
	// Message 提示消息
	Message string
	// Options 选项列表
	Options []string
	// DefaultSelected 默认选中的索引列表
	DefaultSelected []int
	// ResultTitle 选择完成后显示的 title（可选）
	// 如果设置，将优先于全局的 FormatResultTitle 使用
	ResultTitle string
}

// AskMultiSelect 使用配置结构体的多选函数
func AskMultiSelect(field MultiSelectField) ([]int, error) {
	config := newDefaultConfig()
	// 如果提供了 ResultTitle，设置 FormatResultTitle
	if field.ResultTitle != "" {
		titleStr := field.ResultTitle
		config.FormatResultTitle = func(originalMessage string, resultValue string) string {
			return titleStr
		}
	}
	return multiselect.MultiSelect(multiselect.MultiSelectConfig{
		BasePromptConfig: common.BasePromptConfig{
			Message:  field.Message,
			Config:   config,
			Terminal: io.NewStdTerminal(),
		},
		Options:         field.Options,
		DefaultSelected: field.DefaultSelected,
	})
}

// MultiSelectBuilder MultiSelect 的链式构建器
type MultiSelectBuilder struct {
	baseBuilder
	options         []string
	defaultSelected []int
}

// MultiSelect 创建一个 MultiSelect 构建器（链式调用）
// 使用方式：prompt.MultiSelect().Prompt("消息").Options(options).Default([]int{0, 2}).Run()
func MultiSelect() *MultiSelectBuilder {
	return &MultiSelectBuilder{
		defaultSelected: []int{},
	}
}

// Prompt 设置提示消息（覆盖基类方法以返回正确的类型）
func (b *MultiSelectBuilder) Prompt(message string) *MultiSelectBuilder {
	b.baseBuilder.Prompt(message)
	return b
}

// Options 设置选项列表
func (b *MultiSelectBuilder) Options(options []string) *MultiSelectBuilder {
	b.options = options
	return b
}

// Default 设置默认选中的索引列表
func (b *MultiSelectBuilder) Default(indices []int) *MultiSelectBuilder {
	b.defaultSelected = indices
	return b
}

// Run 执行多选并返回结果（返回选中的索引列表）
func (b *MultiSelectBuilder) Run() ([]int, error) {
	return AskMultiSelect(MultiSelectField{
		Message:         b.GetMessage(),
		Options:         b.options,
		DefaultSelected: b.defaultSelected,
	})
}

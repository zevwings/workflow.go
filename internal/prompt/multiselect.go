package prompt

import (
	"github.com/zevwings/workflow/internal/prompt/multiselect"
)

// AskMultiSelect 函数式调用 MultiSelect（保持向后兼容）
// 使用方式：prompt.AskMultiSelect("消息", options, []int{0, 2})
func AskMultiSelect(message string, options []string, defaultSelected []int) ([]int, error) {
	return multiselectFunc(message, options, defaultSelected)
}

// multiselectFunc 统一的多选函数
func multiselectFunc(message string, options []string, defaultSelected []int) ([]int, error) {
	// 构建配置
	config := multiselect.Config{
		FormatPrompt: formatPrompt,
		FormatAnswer: formatAnswer,
		FormatHint:   formatHint,
	}

	return multiselect.MultiSelect(message, options, defaultSelected, config)
}

// MultiSelectBuilder MultiSelect 的链式构建器
type MultiSelectBuilder struct {
	BaseBuilder
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
	b.BaseBuilder.Prompt(message)
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
	return multiselectFunc(b.GetMessage(), b.options, b.defaultSelected)
}

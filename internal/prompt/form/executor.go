package form

import (
	"fmt"

	"github.com/zevwings/workflow/internal/prompt/common"
	"github.com/zevwings/workflow/internal/prompt/confirm"
	multiselectpkg "github.com/zevwings/workflow/internal/prompt/multiselect"
	selectpkg "github.com/zevwings/workflow/internal/prompt/select"
)

// FormExecutor 表单执行器
type FormExecutor struct {
	config     common.PromptConfig
	formConfig FormConfig
}

// NewFormExecutor 创建新的表单执行器
// formConfig 通过 SetFormConfig 设置，由 prompt 包在初始化时调用
var globalFormConfig FormConfig

// SetFormConfig 设置全局 Form 配置（由 prompt 包调用）
func SetFormConfig(config FormConfig) {
	globalFormConfig = config
}

// NewFormExecutor 创建新的表单执行器
func NewFormExecutor() *FormExecutor {
	return &FormExecutor{
		config:     newDefaultConfig(globalFormConfig),
		formConfig: globalFormConfig,
	}
}

// Execute 执行表单字段序列
func (e *FormExecutor) Execute(builder *FormBuilder) (*FormResult, error) {
	result := NewFormResult()
	fields := builder.GetFields()

	for _, field := range fields {
		// 评估条件（如果有）
		if field.Condition != nil {
			if !field.Condition(result) {
				// 条件不满足，跳过该字段
				continue
			}
		}

		// 执行字段
		value, err := e.executeField(field, result)
		if err != nil {
			return nil, fmt.Errorf("执行字段 %s 失败: %w", field.Key, err)
		}

		// 收集结果
		result.Set(field.Key, value)
	}

	return result, nil
}

// executeField 执行单个字段
func (e *FormExecutor) executeField(field FormField, currentResult *FormResult) (interface{}, error) {
	switch field.Type {
	case FieldTypeConfirm:
		return e.executeConfirm(field)
	case FieldTypeInput:
		return e.executeInput(field)
	case FieldTypePassword:
		return e.executePassword(field)
	case FieldTypeSelect:
		return e.executeSelect(field)
	case FieldTypeMultiSelect:
		return e.executeMultiSelect(field)
	case FieldTypeForm:
		return e.executeForm(field)
	default:
		return nil, fmt.Errorf("未知的字段类型: %s", field.Type)
	}
}

// executeConfirm 执行确认字段
func (e *FormExecutor) executeConfirm(field FormField) (bool, error) {
	defaultValue, ok := field.DefaultValue.(bool)
	if !ok {
		defaultValue = false
	}
	return confirm.ConfirmDefault(field.Prompt, defaultValue, e.config)
}

// executeInput 执行输入字段
// 使用 formConfig 中的 AskInputFunc 来保持格式一致
func (e *FormExecutor) executeInput(field FormField) (string, error) {
	defaultValue := ""
	if field.DefaultValue != nil {
		if str, ok := field.DefaultValue.(string); ok {
			defaultValue = str
		}
	}
	// 使用 formConfig 中的函数，保持格式一致
	if e.formConfig.AskInputFunc != nil {
		return e.formConfig.AskInputFunc(field.Prompt, defaultValue, field.Validator)
	}
	// 如果没有设置，返回错误
	return "", fmt.Errorf("AskInputFunc 未设置，请确保 prompt 包已正确初始化")
}

// executePassword 执行密码字段
// 使用 formConfig 中的 AskPasswordFunc 来保持格式一致
func (e *FormExecutor) executePassword(field FormField) (string, error) {
	// 使用 formConfig 中的函数，保持格式一致
	if e.formConfig.AskPasswordFunc != nil {
		return e.formConfig.AskPasswordFunc(field.Prompt, field.Validator)
	}
	// 如果没有设置，返回错误
	return "", fmt.Errorf("AskPasswordFunc 未设置，请确保 prompt 包已正确初始化")
}

// executeSelect 执行单选字段
func (e *FormExecutor) executeSelect(field FormField) (int, error) {
	return selectpkg.SelectDefault(field.Prompt, field.Options, field.DefaultIndex, e.config)
}

// executeMultiSelect 执行多选字段
func (e *FormExecutor) executeMultiSelect(field FormField) ([]int, error) {
	return multiselectpkg.MultiSelectDefault(field.Prompt, field.Options, field.DefaultSelected, e.config)
}

// executeForm 执行嵌套表单字段
func (e *FormExecutor) executeForm(field FormField) (*FormResult, error) {
	if field.NestedForm == nil {
		return nil, fmt.Errorf("嵌套表单不能为空")
	}
	return e.Execute(field.NestedForm)
}

// newDefaultConfig 创建默认配置（使用 formConfig 中的格式化函数）
func newDefaultConfig(formConfig FormConfig) common.PromptConfig {
	return common.PromptConfig{
		FormatPrompt: formConfig.FormatPrompt,
		FormatAnswer: formConfig.FormatAnswer,
		FormatHint:   formConfig.FormatHint,
	}
}

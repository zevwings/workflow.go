package form

import (
	"fmt"
)

// FormBuilder 表单构建器（链式 API）
type FormBuilder struct {
	fields    []FormField
	validator FormValidator
	title     string // 表单标题（用于显示在分割线中）
}

// NewFormBuilder 创建新的表单构建器
func NewFormBuilder() *FormBuilder {
	return &FormBuilder{
		fields: make([]FormField, 0),
	}
}

// AddConfirm 添加确认字段
func (b *FormBuilder) AddConfirm(config ConfirmFormField) *FormBuilder {
	b.fields = append(b.fields, FormField{
		Key:          config.Key,
		Type:         FieldTypeConfirm,
		Prompt:       config.Prompt,
		DefaultValue: config.DefaultValue,
		ResultTitle:  config.ResultTitle,
		Condition:    config.Condition,
	})
	return b
}

// AddInput 添加输入字段
func (b *FormBuilder) AddInput(config InputFormField) *FormBuilder {
	b.fields = append(b.fields, FormField{
		Key:          config.Key,
		Type:         FieldTypeInput,
		Prompt:       config.Prompt,
		DefaultValue: config.DefaultValue,
		Validator:    config.Validator,
		ResultTitle:  config.ResultTitle,
		Condition:    config.Condition,
	})
	return b
}

// AddPassword 添加密码字段
// DefaultValue 为空字符串表示无默认值
func (b *FormBuilder) AddPassword(config PasswordFormField) *FormBuilder {
	b.fields = append(b.fields, FormField{
		Key:          config.Key,
		Type:         FieldTypePassword,
		Prompt:       config.Prompt,
		DefaultValue: config.DefaultValue,
		Validator:    config.Validator,
		ResultTitle:  config.ResultTitle,
		Condition:    config.Condition,
	})
	return b
}

// AddSelect 添加单选字段
func (b *FormBuilder) AddSelect(config SelectFormField) *FormBuilder {
	b.fields = append(b.fields, FormField{
		Key:          config.Key,
		Type:         FieldTypeSelect,
		Prompt:       config.Prompt,
		Options:      config.Options,
		DefaultIndex: config.DefaultIndex,
		ResultTitle:  config.ResultTitle,
		Condition:    config.Condition,
	})
	return b
}

// AddMultiSelect 添加多选字段
func (b *FormBuilder) AddMultiSelect(config MultiSelectFormField) *FormBuilder {
	b.fields = append(b.fields, FormField{
		Key:             config.Key,
		Type:            FieldTypeMultiSelect,
		Prompt:          config.Prompt,
		Options:         config.Options,
		DefaultSelected: config.DefaultSelected,
		ResultTitle:     config.ResultTitle,
		Condition:       config.Condition,
	})
	return b
}

// AddForm 添加嵌套表单字段
func (b *FormBuilder) AddForm(config NestedFormField) *FormBuilder {
	b.fields = append(b.fields, FormField{
		Key:         config.Key,
		Type:        FieldTypeForm,
		Prompt:      config.Prompt,
		NestedForm:  config.NestedForm,
		ResultTitle: config.ResultTitle,
		Condition:   config.Condition,
	})
	return b
}

// Condition 为最后一个字段设置条件函数
// 注意：这个方法需要在添加字段后立即调用
func (b *FormBuilder) Condition(condition Condition) *FormBuilder {
	if len(b.fields) == 0 {
		return b
	}
	// 为最后一个字段设置条件
	lastField := &b.fields[len(b.fields)-1]
	lastField.Condition = condition
	return b
}

// Validate 设置表单级验证器
func (b *FormBuilder) Validate(validator FormValidator) *FormBuilder {
	b.validator = validator
	return b
}

// SetTitle 设置表单标题（用于显示在分割线中）
func (b *FormBuilder) SetTitle(title string) *FormBuilder {
	b.title = title
	return b
}

// GetTitle 获取表单标题（内部使用）
func (b *FormBuilder) GetTitle() string {
	return b.title
}

// Run 执行表单并返回结果
func (b *FormBuilder) Run() (*FormResult, error) {
	executor := NewFormExecutor()
	result, err := executor.Execute(b)
	if err != nil {
		return nil, err
	}
	// 执行表单级验证（如果有）
	if b.validator != nil {
		if err := b.validator(result); err != nil {
			return nil, fmt.Errorf("表单验证失败: %w", err)
		}
	}
	return result, nil
}

// GetFields 获取字段列表（内部使用）
func (b *FormBuilder) GetFields() []FormField {
	return b.fields
}

package form

import (
	"fmt"

	"github.com/zevwings/workflow/internal/prompt/input"
)

// FormBuilder 表单构建器（链式 API）
type FormBuilder struct {
	fields    []FormField
	validator FormValidator
}

// NewFormBuilder 创建新的表单构建器
func NewFormBuilder() *FormBuilder {
	return &FormBuilder{
		fields: make([]FormField, 0),
	}
}

// AddConfirm 添加确认字段
func (b *FormBuilder) AddConfirm(key, prompt string, defaultValue bool) *FormBuilder {
	b.fields = append(b.fields, FormField{
		Key:          key,
		Type:         FieldTypeConfirm,
		Prompt:       prompt,
		DefaultValue: defaultValue,
	})
	return b
}

// AddInput 添加输入字段
func (b *FormBuilder) AddInput(key, prompt, defaultValue string, validator input.Validator) *FormBuilder {
	b.fields = append(b.fields, FormField{
		Key:          key,
		Type:         FieldTypeInput,
		Prompt:       prompt,
		DefaultValue: defaultValue,
		Validator:    validator,
	})
	return b
}

// AddPassword 添加密码字段
func (b *FormBuilder) AddPassword(key, prompt string, validator input.Validator) *FormBuilder {
	b.fields = append(b.fields, FormField{
		Key:       key,
		Type:      FieldTypePassword,
		Prompt:    prompt,
		Validator: validator,
	})
	return b
}

// AddSelect 添加单选字段
func (b *FormBuilder) AddSelect(key, prompt string, options []string, defaultIndex int) *FormBuilder {
	b.fields = append(b.fields, FormField{
		Key:          key,
		Type:         FieldTypeSelect,
		Prompt:       prompt,
		Options:      options,
		DefaultIndex: defaultIndex,
	})
	return b
}

// AddMultiSelect 添加多选字段
func (b *FormBuilder) AddMultiSelect(key, prompt string, options []string, defaultSelected []int) *FormBuilder {
	b.fields = append(b.fields, FormField{
		Key:             key,
		Type:            FieldTypeMultiSelect,
		Prompt:          prompt,
		Options:         options,
		DefaultSelected: defaultSelected,
	})
	return b
}

// AddForm 添加嵌套表单字段
func (b *FormBuilder) AddForm(key, prompt string, nestedForm *FormBuilder) *FormBuilder {
	b.fields = append(b.fields, FormField{
		Key:        key,
		Type:       FieldTypeForm,
		Prompt:     prompt,
		NestedForm: nestedForm,
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

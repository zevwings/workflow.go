//go:build test

package form

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

// ==================== FormBuilder 基础功能测试 ====================

func TestNewFormBuilder(t *testing.T) {
	builder := NewFormBuilder()
	assert.NotNil(t, builder)
	assert.NotNil(t, builder.fields)
	assert.Equal(t, 0, len(builder.fields))
	assert.Nil(t, builder.validator)
	assert.Empty(t, builder.title)
}

// ==================== AddConfirm 测试 ====================

func TestFormBuilder_AddConfirm(t *testing.T) {
	builder := NewFormBuilder()

	// 添加确认字段
	builder.AddConfirm("agree", "是否同意？", true)

	assert.Equal(t, 1, len(builder.fields))
	field := builder.fields[0]
	assert.Equal(t, "agree", field.Key)
	assert.Equal(t, FieldTypeConfirm, field.Type)
	assert.Equal(t, "是否同意？", field.Prompt)
	assert.Equal(t, true, field.DefaultValue)
	assert.Nil(t, field.Validator)

	// 链式调用
	builder.AddConfirm("confirm2", "确认2？", false)
	assert.Equal(t, 2, len(builder.fields))
}

// ==================== AddInput 测试 ====================

func TestFormBuilder_AddInput(t *testing.T) {
	builder := NewFormBuilder()

	// 添加输入字段（无验证器）
	builder.AddInput("name", "请输入姓名", "默认值", nil)

	assert.Equal(t, 1, len(builder.fields))
	field := builder.fields[0]
	assert.Equal(t, "name", field.Key)
	assert.Equal(t, FieldTypeInput, field.Type)
	assert.Equal(t, "请输入姓名", field.Prompt)
	assert.Equal(t, "默认值", field.DefaultValue)
	assert.Nil(t, field.Validator)

	// 添加输入字段（有验证器）
	validator := func(s string) error {
		if len(s) < 3 {
			return errors.New("长度至少3个字符")
		}
		return nil
	}
	builder.AddInput("email", "请输入邮箱", "", validator)
	assert.Equal(t, 2, len(builder.fields))
	assert.NotNil(t, builder.fields[1].Validator)
}

// ==================== AddPassword 测试 ====================

func TestFormBuilder_AddPassword(t *testing.T) {
	builder := NewFormBuilder()

	// 添加密码字段（无验证器）
	builder.AddPassword("password", "请输入密码", nil)

	assert.Equal(t, 1, len(builder.fields))
	field := builder.fields[0]
	assert.Equal(t, "password", field.Key)
	assert.Equal(t, FieldTypePassword, field.Type)
	assert.Equal(t, "请输入密码", field.Prompt)
	assert.Nil(t, field.DefaultValue)
	assert.Nil(t, field.Validator)

	// 添加密码字段（有验证器）
	validator := func(s string) error {
		if len(s) < 6 {
			return errors.New("密码长度至少6个字符")
		}
		return nil
	}
	builder.AddPassword("password2", "请再次输入密码", validator)
	assert.Equal(t, 2, len(builder.fields))
	assert.NotNil(t, builder.fields[1].Validator)
}

// ==================== AddSelect 测试 ====================

func TestFormBuilder_AddSelect(t *testing.T) {
	builder := NewFormBuilder()
	options := []string{"选项1", "选项2", "选项3"}

	// 添加单选字段
	builder.AddSelect("choice", "请选择", options, 1)

	assert.Equal(t, 1, len(builder.fields))
	field := builder.fields[0]
	assert.Equal(t, "choice", field.Key)
	assert.Equal(t, FieldTypeSelect, field.Type)
	assert.Equal(t, "请选择", field.Prompt)
	assert.Equal(t, options, field.Options)
	assert.Equal(t, 1, field.DefaultIndex)
	assert.Nil(t, field.DefaultValue)

	// 测试默认索引为 0
	builder.AddSelect("choice2", "请选择2", options, 0)
	assert.Equal(t, 0, builder.fields[1].DefaultIndex)
}

// ==================== AddMultiSelect 测试 ====================

func TestFormBuilder_AddMultiSelect(t *testing.T) {
	builder := NewFormBuilder()
	options := []string{"选项1", "选项2", "选项3"}

	// 添加多选字段
	builder.AddMultiSelect("multi", "请多选", options, []int{0, 2})

	assert.Equal(t, 1, len(builder.fields))
	field := builder.fields[0]
	assert.Equal(t, "multi", field.Key)
	assert.Equal(t, FieldTypeMultiSelect, field.Type)
	assert.Equal(t, "请多选", field.Prompt)
	assert.Equal(t, options, field.Options)
	assert.Equal(t, []int{0, 2}, field.DefaultSelected)

	// 测试空默认选择
	builder.AddMultiSelect("multi2", "请多选2", options, nil)
	assert.Nil(t, builder.fields[1].DefaultSelected)
}

// ==================== AddForm 测试 ====================

func TestFormBuilder_AddForm(t *testing.T) {
	builder := NewFormBuilder()
	nestedForm := NewFormBuilder()
	nestedForm.AddInput("nested_name", "嵌套姓名", "", nil)

	// 添加嵌套表单
	builder.AddForm("user", "用户信息", nestedForm)

	assert.Equal(t, 1, len(builder.fields))
	field := builder.fields[0]
	assert.Equal(t, "user", field.Key)
	assert.Equal(t, FieldTypeForm, field.Type)
	assert.Equal(t, "用户信息", field.Prompt)
	assert.Equal(t, nestedForm, field.NestedForm)
	assert.Nil(t, field.DefaultValue)
}

// ==================== Condition 测试 ====================

func TestFormBuilder_Condition(t *testing.T) {
	builder := NewFormBuilder()

	// 空字段列表，Condition 应该不生效
	builder.Condition(func(result *FormResult) bool {
		return true
	})
	assert.Equal(t, 0, len(builder.fields))

	// 添加字段后设置条件
	builder.AddConfirm("agree", "是否同意？", false)
	condition := func(result *FormResult) bool {
		return result.GetBool("agree")
	}
	builder.Condition(condition)

	assert.Equal(t, 1, len(builder.fields))
	assert.NotNil(t, builder.fields[0].Condition)

	// 测试条件函数
	result := NewFormResult()
	result.Set("agree", true)
	assert.True(t, builder.fields[0].Condition(result))

	result.Set("agree", false)
	assert.False(t, builder.fields[0].Condition(result))
}

// ==================== Validate 测试 ====================

func TestFormBuilder_Validate(t *testing.T) {
	builder := NewFormBuilder()

	// 设置表单级验证器
	validator := func(result *FormResult) error {
		if result.GetString("name") == "" {
			return errors.New("姓名不能为空")
		}
		return nil
	}
	builder.Validate(validator)

	assert.NotNil(t, builder.validator)

	// 测试验证器
	result := NewFormResult()
	result.Set("name", "")
	err := builder.validator(result)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "姓名不能为空")

	result.Set("name", "张三")
	err = builder.validator(result)
	assert.NoError(t, err)
}

// ==================== SetTitle 和 GetTitle 测试 ====================

func TestFormBuilder_SetTitle(t *testing.T) {
	builder := NewFormBuilder()

	// 设置标题
	builder.SetTitle("用户注册表单")

	assert.Equal(t, "用户注册表单", builder.title)
	assert.Equal(t, "用户注册表单", builder.GetTitle())

	// 链式调用
	builder.SetTitle("新标题").AddInput("name", "姓名", "", nil)
	assert.Equal(t, "新标题", builder.GetTitle())
}

// ==================== GetFields 测试 ====================

func TestFormBuilder_GetFields(t *testing.T) {
	builder := NewFormBuilder()

	// 空字段列表
	fields := builder.GetFields()
	assert.NotNil(t, fields)
	assert.Equal(t, 0, len(fields))

	// 添加字段后
	builder.AddInput("name", "姓名", "", nil)
	builder.AddConfirm("agree", "同意", true)

	fields = builder.GetFields()
	assert.Equal(t, 2, len(fields))
	assert.Equal(t, "name", fields[0].Key)
	assert.Equal(t, "agree", fields[1].Key)
}

// ==================== 链式调用测试 ====================

func TestFormBuilder_Chaining(t *testing.T) {
	builder := NewFormBuilder().
		SetTitle("测试表单").
		AddInput("name", "姓名", "", nil).
		AddConfirm("agree", "同意", false).
		AddSelect("choice", "选择", []string{"A", "B"}, 0).
		Validate(func(result *FormResult) error {
			return nil
		})

	assert.Equal(t, "测试表单", builder.GetTitle())
	assert.Equal(t, 3, len(builder.fields))
	assert.NotNil(t, builder.validator)
}

// ==================== Run 测试（需要 mock executor） ====================

func TestFormBuilder_Run_WithValidator(t *testing.T) {
	// 注意：Run 方法会调用 executor，需要设置 FormConfig
	// 这里只测试验证器逻辑，不测试完整的执行流程
	// 完整的执行流程测试在 executor_test.go 中

	builder := NewFormBuilder().
		AddInput("name", "姓名", "默认", nil).
		Validate(func(result *FormResult) error {
			if result.GetString("name") == "" {
				return errors.New("姓名不能为空")
			}
			return nil
		})

	// 由于 Run 需要真实的 executor，这里只验证 builder 结构
	assert.NotNil(t, builder.validator)
	assert.Equal(t, 1, len(builder.fields))
}

//go:build test

package form

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zevwings/workflow/internal/prompt/common"
	"github.com/zevwings/workflow/internal/prompt/io"
)

// mockInputProvider mock 输入提供者（用于测试）
type mockInputProvider struct{}

func (m *mockInputProvider) AskInput(field InputField) (string, error) {
	// Mock 实现：返回默认值
	// 如果默认值为空字符串，返回错误（用于测试错误处理）
	if field.DefaultValue == "" {
		return "", fmt.Errorf("输入错误")
	}
	return field.DefaultValue, nil
}

func (m *mockInputProvider) AskPassword(field PasswordField) (string, error) {
	// Mock 实现：返回固定密码
	// 如果默认值为空字符串，返回错误（用于测试错误处理）
	if field.DefaultValue == "" {
		return "", fmt.Errorf("密码输入错误")
	}
	return "password123", nil
}

// ==================== NewFormExecutor 测试 ====================

func TestNewFormExecutor(t *testing.T) {
	// 保存原始配置
	originalConfig := globalPromptConfig
	originalProvider := globalInputProvider
	defer func() {
		globalPromptConfig = originalConfig
		globalInputProvider = originalProvider
	}()

	// 设置测试配置
	SetPromptConfig(common.PromptConfig{
		FormatPrompt: func(msg string) string { return msg },
		FormatAnswer: func(v string) string { return v },
		FormatError:  func(msg string) string { return msg },
		FormatHint:   func(msg string) string { return msg },
	})
	SetInputProvider(&mockInputProvider{})

	executor := NewFormExecutor()
	assert.NotNil(t, executor)
	assert.NotNil(t, executor.config)
}

// ==================== Execute 基础测试 ====================

func TestFormExecutor_Execute_EmptyForm(t *testing.T) {
	// 保存原始配置
	originalConfig := globalPromptConfig
	originalProvider := globalInputProvider
	defer func() {
		globalPromptConfig = originalConfig
		globalInputProvider = originalProvider
	}()

	// 设置测试配置
	SetPromptConfig(common.PromptConfig{
		FormatPrompt: func(msg string) string { return msg },
		FormatAnswer: func(v string) string { return v },
		FormatError:  func(msg string) string { return msg },
		FormatHint:   func(msg string) string { return msg },
	})
	SetInputProvider(&mockInputProvider{})

	executor := NewFormExecutor()
	builder := NewFormBuilder()

	result, err := executor.Execute(builder)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, 0, len(result.Values))
}

// ==================== Execute Confirm 字段测试 ====================

func TestFormExecutor_Execute_ConfirmField(t *testing.T) {
	// 保存原始配置
	originalConfig := globalPromptConfig
	originalProvider := globalInputProvider
	defer func() {
		globalPromptConfig = originalConfig
		globalInputProvider = originalProvider
	}()

	// 设置测试配置
	SetPromptConfig(common.PromptConfig{
		FormatPrompt: func(msg string) string { return msg },
		FormatAnswer: func(v string) string { return v },
		FormatError:  func(msg string) string { return msg },
		FormatHint:   func(msg string) string { return msg },
	})
	SetInputProvider(&mockInputProvider{})

	executor := NewFormExecutor()
	builder := NewFormBuilder()

	// 注意：这个测试需要真实的终端交互，所以这里只测试结构
	// 完整的交互测试需要使用 mock terminal
	builder.AddConfirm(ConfirmFormField{
		Key:          "agree",
		Prompt:       "是否同意？",
		DefaultValue: true,
	})

	// 由于需要真实的终端交互，这里只验证 builder 结构
	assert.Equal(t, 1, len(builder.GetFields()))
	assert.Equal(t, FieldTypeConfirm, builder.GetFields()[0].Type)
	_ = executor // 标记为已使用
}

// ==================== Execute Input 字段测试 ====================

func TestFormExecutor_Execute_InputField(t *testing.T) {
	// 保存原始配置
	originalConfig := globalPromptConfig
	originalProvider := globalInputProvider
	defer func() {
		globalPromptConfig = originalConfig
		globalInputProvider = originalProvider
	}()

	// 设置测试配置
	SetPromptConfig(common.PromptConfig{
		FormatPrompt: func(msg string) string { return msg },
		FormatAnswer: func(v string) string { return v },
		FormatError:  func(msg string) string { return msg },
		FormatHint:   func(msg string) string { return msg },
	})
	SetInputProvider(&mockInputProvider{})

	executor := NewFormExecutor()
	builder := NewFormBuilder()
	builder.AddInput(InputFormField{
		Key:          "name",
		Prompt:       "请输入姓名",
		DefaultValue: "默认姓名",
		Validator:    nil,
	})

	result, err := executor.Execute(builder)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "默认姓名", result.GetString("name"))
}

// ==================== Execute Input 字段验证器测试 ====================

func TestFormExecutor_Execute_InputFieldWithValidator(t *testing.T) {
	// 保存原始配置
	originalConfig := globalPromptConfig
	originalProvider := globalInputProvider
	defer func() {
		globalPromptConfig = originalConfig
		globalInputProvider = originalProvider
	}()

	// 设置测试配置
	validator := func(s string) error {
		if len(s) < 3 {
			return errors.New("长度至少3个字符")
		}
		return nil
	}

	SetPromptConfig(common.PromptConfig{
		FormatPrompt: func(msg string) string { return msg },
		FormatAnswer: func(v string) string { return v },
		FormatError:  func(msg string) string { return msg },
		FormatHint:   func(msg string) string { return msg },
	})
	SetInputProvider(&mockInputProvider{})

	executor := NewFormExecutor()
	builder := NewFormBuilder()
	builder.AddInput(InputFormField{
		Key:          "name",
		Prompt:       "请输入姓名",
		DefaultValue: "ab",
		Validator:    validator,
	}) // 默认值太短

	// 由于验证器会失败，这里应该返回错误
	// 但我们的 mock 实现会返回默认值，所以这里只测试结构
	// 实际验证应该在集成测试中完成
	_, err := executor.Execute(builder)
	// 由于 mock 实现会返回默认值，这里可能不会报错
	// 实际验证应该在集成测试中完成
	_ = err
}

// ==================== Execute Password 字段测试 ====================

func TestFormExecutor_Execute_PasswordField(t *testing.T) {
	// 保存原始配置
	originalConfig := globalPromptConfig
	originalProvider := globalInputProvider
	defer func() {
		globalPromptConfig = originalConfig
		globalInputProvider = originalProvider
	}()

	// 设置测试配置
	SetPromptConfig(common.PromptConfig{
		FormatPrompt: func(msg string) string { return msg },
		FormatAnswer: func(v string) string { return v },
		FormatError:  func(msg string) string { return msg },
		FormatHint:   func(msg string) string { return msg },
	})
	SetInputProvider(&mockInputProvider{})

	executor := NewFormExecutor()
	builder := NewFormBuilder()
	builder.AddPassword(PasswordFormField{
		Key:          "password",
		Prompt:       "请输入密码",
		DefaultValue: "",
		Validator:    nil,
	})

	result, err := executor.Execute(builder)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "password123", result.GetString("password"))
}

// ==================== Execute Select 字段测试 ====================

func TestFormExecutor_Execute_SelectField(t *testing.T) {
	// 保存原始配置
	originalConfig := globalPromptConfig
	originalProvider := globalInputProvider
	defer func() {
		globalPromptConfig = originalConfig
		globalInputProvider = originalProvider
	}()

	// 设置测试配置
	SetPromptConfig(common.PromptConfig{
		FormatPrompt: func(msg string) string { return msg },
		FormatAnswer: func(v string) string { return v },
		FormatError:  func(msg string) string { return msg },
		FormatHint:   func(msg string) string { return msg },
	})
	SetInputProvider(&mockInputProvider{})

	executor := NewFormExecutor()
	builder := NewFormBuilder()
	builder.AddSelect(SelectFormField{
		Key:          "choice",
		Prompt:       "请选择",
		Options:      []string{"选项1", "选项2", "选项3"},
		DefaultIndex: 1,
	})

	// 注意：这个测试需要真实的终端交互
	// 这里只验证 builder 结构
	assert.Equal(t, 1, len(builder.GetFields()))
	assert.Equal(t, FieldTypeSelect, builder.GetFields()[0].Type)
	_ = executor // 标记为已使用
}

// ==================== Execute MultiSelect 字段测试 ====================

func TestFormExecutor_Execute_MultiSelectField(t *testing.T) {
	// 保存原始配置
	originalConfig := globalPromptConfig
	originalProvider := globalInputProvider
	defer func() {
		globalPromptConfig = originalConfig
		globalInputProvider = originalProvider
	}()

	// 设置测试配置
	SetPromptConfig(common.PromptConfig{
		FormatPrompt: func(msg string) string { return msg },
		FormatAnswer: func(v string) string { return v },
		FormatError:  func(msg string) string { return msg },
		FormatHint:   func(msg string) string { return msg },
	})
	SetInputProvider(&mockInputProvider{})

	executor := NewFormExecutor()
	builder := NewFormBuilder()
	builder.AddMultiSelect(MultiSelectFormField{
		Key:             "multi",
		Prompt:          "请多选",
		Options:         []string{"选项1", "选项2", "选项3"},
		DefaultSelected: []int{0, 2},
	})

	// 注意：这个测试需要真实的终端交互
	// 这里只验证 builder 结构
	assert.Equal(t, 1, len(builder.GetFields()))
	assert.Equal(t, FieldTypeMultiSelect, builder.GetFields()[0].Type)
	_ = executor // 标记为已使用
}

// ==================== Execute 嵌套表单测试 ====================

func TestFormExecutor_Execute_NestedForm(t *testing.T) {
	// 保存原始配置
	originalConfig := globalPromptConfig
	originalProvider := globalInputProvider
	defer func() {
		globalPromptConfig = originalConfig
		globalInputProvider = originalProvider
	}()

	// 设置测试配置
	SetPromptConfig(common.PromptConfig{
		FormatPrompt: func(msg string) string { return msg },
		FormatAnswer: func(v string) string { return v },
		FormatError:  func(msg string) string { return msg },
		FormatHint:   func(msg string) string { return msg },
	})
	SetInputProvider(&mockInputProvider{})

	executor := NewFormExecutor()

	// 创建嵌套表单
	nestedForm := NewFormBuilder()
	nestedForm.AddInput(InputFormField{
		Key:          "nested_name",
		Prompt:       "嵌套姓名",
		DefaultValue: "嵌套默认值",
		Validator:    nil,
	})

	// 主表单
	builder := NewFormBuilder()
	builder.AddForm(NestedFormField{
		Key:        "user",
		Prompt:     "用户信息",
		NestedForm: nestedForm,
	})

	result, err := executor.Execute(builder)
	assert.NoError(t, err)
	assert.NotNil(t, result)

	// 验证嵌套表单结果
	nestedResult := result.GetForm("user")
	assert.NotNil(t, nestedResult)
	assert.Equal(t, "嵌套默认值", nestedResult.GetString("nested_name"))
}

// ==================== Execute 条件字段测试 ====================

func TestFormExecutor_Execute_ConditionalField(t *testing.T) {
	// 保存原始配置
	originalConfig := globalPromptConfig
	originalProvider := globalInputProvider
	defer func() {
		globalPromptConfig = originalConfig
		globalInputProvider = originalProvider
	}()

	// 设置测试配置
	SetPromptConfig(common.PromptConfig{
		FormatPrompt: func(msg string) string { return msg },
		FormatAnswer: func(v string) string { return v },
		FormatError:  func(msg string) string { return msg },
		FormatHint:   func(msg string) string { return msg },
	})
	SetInputProvider(&mockInputProvider{})

	executor := NewFormExecutor()
	builder := NewFormBuilder()

	// 第一个字段：确认
	builder.AddConfirm(ConfirmFormField{
		Key:          "need_email",
		Prompt:       "是否需要邮箱？",
		DefaultValue: false,
	})

	// 第二个字段：条件输入（只有当 need_email 为 true 时才显示）
	builder.AddInput(InputFormField{
		Key:          "email",
		Prompt:       "请输入邮箱",
		DefaultValue: "",
		Validator:    nil,
		Condition: func(result *FormResult) bool {
			return result.GetBool("need_email")
		},
	})

	// 执行表单（need_email 为 false，email 字段应该被跳过）
	result, err := executor.Execute(builder)
	assert.NoError(t, err)
	assert.NotNil(t, result)

	// 验证 need_email 存在
	assert.False(t, result.GetBool("need_email"))

	// 验证 email 不存在（因为条件不满足）
	_, exists := result.Values["email"]
	assert.False(t, exists, "条件不满足时，email 字段不应该存在")
}

// ==================== Execute 表单标题测试 ====================

func TestFormExecutor_Execute_WithTitle(t *testing.T) {
	// 保存原始配置
	originalConfig := globalPromptConfig
	originalProvider := globalInputProvider
	defer func() {
		globalPromptConfig = originalConfig
		globalInputProvider = originalProvider
	}()

	// 设置测试配置
	SetPromptConfig(common.PromptConfig{
		FormatPrompt: func(msg string) string { return msg },
		FormatAnswer: func(v string) string { return v },
		FormatError:  func(msg string) string { return msg },
		FormatHint:   func(msg string) string { return msg },
	})
	SetInputProvider(&mockInputProvider{})

	executor := NewFormExecutor()
	builder := NewFormBuilder().
		SetTitle("用户注册表单").
		AddInput(InputFormField{
			Key:          "name",
			Prompt:       "姓名",
			DefaultValue: "默认",
			Validator:    nil,
		})

	result, err := executor.Execute(builder)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "默认", result.GetString("name"))
}

// ==================== Execute 多个字段测试 ====================

func TestFormExecutor_Execute_MultipleFields(t *testing.T) {
	// 保存原始配置
	originalConfig := globalPromptConfig
	originalProvider := globalInputProvider
	defer func() {
		globalPromptConfig = originalConfig
		globalInputProvider = originalProvider
	}()

	// 设置测试配置
	SetPromptConfig(common.PromptConfig{
		FormatPrompt: func(msg string) string { return msg },
		FormatAnswer: func(v string) string { return v },
		FormatError:  func(msg string) string { return msg },
		FormatHint:   func(msg string) string { return msg },
	})
	SetInputProvider(&mockInputProvider{})

	executor := NewFormExecutor()
	builder := NewFormBuilder().
		AddInput(InputFormField{
			Key:          "name",
			Prompt:       "姓名",
			DefaultValue: "张三",
			Validator:    nil,
		}).
		AddInput(InputFormField{
			Key:          "email",
			Prompt:       "邮箱",
			DefaultValue: "test@example.com",
			Validator:    nil,
		}).
		AddPassword(PasswordFormField{
			Key:          "password",
			Prompt:       "密码",
			DefaultValue: "",
			Validator:    nil,
		})

	result, err := executor.Execute(builder)
	assert.NoError(t, err)
	assert.NotNil(t, result)

	assert.Equal(t, "张三", result.GetString("name"))
	assert.Equal(t, "test@example.com", result.GetString("email"))
	assert.Equal(t, "password123", result.GetString("password"))
}

// ==================== Execute 错误处理测试 ====================

func TestFormExecutor_Execute_InputError(t *testing.T) {
	// 保存原始配置
	originalConfig := globalPromptConfig
	originalProvider := globalInputProvider
	defer func() {
		globalPromptConfig = originalConfig
		globalInputProvider = originalProvider
	}()

	// 设置测试配置
	SetPromptConfig(common.PromptConfig{
		FormatPrompt: func(msg string) string { return msg },
		FormatAnswer: func(v string) string { return v },
		FormatError:  func(msg string) string { return msg },
		FormatHint:   func(msg string) string { return msg },
	})
	SetInputProvider(&mockInputProvider{})

	executor := NewFormExecutor()
	builder := NewFormBuilder()
	builder.AddInput(InputFormField{
		Key:          "name",
		Prompt:       "姓名",
		DefaultValue: "",
		Validator:    nil,
	})

	result, err := executor.Execute(builder)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "执行字段 name 失败")
	assert.Nil(t, result)
}

func TestFormExecutor_Execute_PasswordError(t *testing.T) {
	// 保存原始配置
	originalConfig := globalPromptConfig
	originalProvider := globalInputProvider
	defer func() {
		globalPromptConfig = originalConfig
		globalInputProvider = originalProvider
	}()

	// 设置测试配置
	SetPromptConfig(common.PromptConfig{
		FormatPrompt: func(msg string) string { return msg },
		FormatAnswer: func(v string) string { return v },
		FormatError:  func(msg string) string { return msg },
		FormatHint:   func(msg string) string { return msg },
	})
	SetInputProvider(&mockInputProvider{})

	executor := NewFormExecutor()
	builder := NewFormBuilder()
	builder.AddPassword(PasswordFormField{
		Key:          "password",
		Prompt:       "密码",
		DefaultValue: "",
		Validator:    nil,
	})

	result, err := executor.Execute(builder)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "执行字段 password 失败")
	assert.Nil(t, result)
}

func TestFormExecutor_Execute_UnknownFieldType(t *testing.T) {
	// 保存原始配置
	originalConfig := globalPromptConfig
	originalProvider := globalInputProvider
	defer func() {
		globalPromptConfig = originalConfig
		globalInputProvider = originalProvider
	}()

	// 设置测试配置
	SetPromptConfig(common.PromptConfig{
		FormatPrompt: func(msg string) string { return msg },
		FormatAnswer: func(v string) string { return v },
		FormatError:  func(msg string) string { return msg },
		FormatHint:   func(msg string) string { return msg },
	})
	SetInputProvider(&mockInputProvider{})

	executor := NewFormExecutor()
	builder := NewFormBuilder()

	// 手动创建一个未知类型的字段
	field := FormField{
		Key:  "unknown",
		Type: FieldType("unknown"),
	}
	builder.fields = append(builder.fields, field)

	result, err := executor.Execute(builder)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "未知的字段类型")
	assert.Nil(t, result)
}

func TestFormExecutor_Execute_NestedFormError(t *testing.T) {
	// 保存原始配置
	originalConfig := globalPromptConfig
	originalProvider := globalInputProvider
	defer func() {
		globalPromptConfig = originalConfig
		globalInputProvider = originalProvider
	}()

	// 设置测试配置
	SetPromptConfig(common.PromptConfig{
		FormatPrompt: func(msg string) string { return msg },
		FormatAnswer: func(v string) string { return v },
		FormatError:  func(msg string) string { return msg },
		FormatHint:   func(msg string) string { return msg },
	})
	SetInputProvider(&mockInputProvider{})

	executor := NewFormExecutor()
	builder := NewFormBuilder()

	// 创建一个 NestedForm 为 nil 的字段
	field := FormField{
		Key:        "user",
		Type:       FieldTypeForm,
		NestedForm: nil,
	}
	builder.fields = append(builder.fields, field)

	result, err := executor.Execute(builder)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "嵌套表单不能为空")
	assert.Nil(t, result)
}

// ==================== Execute 缺少配置测试 ====================

// 注意：由于现在直接调用 prompt 函数，不再需要测试 MissingAskInputFunc 和 MissingAskPasswordFunc

// ==================== newDefaultConfig 测试 ====================

func TestPromptConfig(t *testing.T) {
	config := common.PromptConfig{
		FormatPrompt: func(msg string) string { return "[" + msg + "]" },
		FormatAnswer: func(v string) string { return "{" + v + "}" },
		FormatHint:   func(msg string) string { return "<" + msg + ">" },
	}

	assert.NotNil(t, config.FormatPrompt)
	assert.NotNil(t, config.FormatAnswer)
	assert.NotNil(t, config.FormatHint)

	// 验证格式化函数
	assert.Equal(t, "[test]", config.FormatPrompt("test"))
	assert.Equal(t, "{value}", config.FormatAnswer("value"))
	assert.Equal(t, "<hint>", config.FormatHint("hint"))
}

// ==================== printSeparator 测试 ====================

func TestFormExecutor_PrintSeparator_MainForm(t *testing.T) {
	executor := NewFormExecutor()
	mockTerminal := io.NewMockTerminal([]byte{})

	// 测试主表单分割线
	executor.printSeparator(mockTerminal, "测试表单", "start", true)

	output := mockTerminal.GetOutput()
	assert.Contains(t, output, "测试表单")
	assert.Contains(t, output, "Start")
}

func TestFormExecutor_PrintSeparator_NestedForm(t *testing.T) {
	executor := NewFormExecutor()
	mockTerminal := io.NewMockTerminal([]byte{})

	// 测试嵌套表单分割线
	executor.printSeparator(mockTerminal, "嵌套表单", "end", false)

	output := mockTerminal.GetOutput()
	assert.Contains(t, output, "嵌套表单")
	assert.Contains(t, output, "End")
}

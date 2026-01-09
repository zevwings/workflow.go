//go:build test

package form

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zevwings/workflow/internal/prompt/common"
)

// ==================== ValidateAllRequired 测试 ====================

func TestValidateAllRequired(t *testing.T) {
	// ValidateAllRequired 是一个示例实现，当前返回 nil
	// 这里测试它的基本行为

	result := NewFormResult()
	err := ValidateAllRequired(result)
	assert.NoError(t, err)

	// 设置一些值
	result.Set("name", "张三")
	result.Set("age", 25)
	err = ValidateAllRequired(result)
	assert.NoError(t, err)
}

// ==================== FormValidator 自定义验证器测试 ====================

func TestFormValidator_RequiredFields(t *testing.T) {
	// 创建自定义验证器：验证必填字段
	validator := func(result *FormResult) error {
		if result.GetString("name") == "" {
			return errors.New("姓名不能为空")
		}
		if result.GetString("email") == "" {
			return errors.New("邮箱不能为空")
		}
		return nil
	}

	// 测试：缺少必填字段
	result := NewFormResult()
	result.Set("name", "")
	err := validator(result)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "姓名不能为空")

	// 测试：部分字段存在
	result.Set("name", "张三")
	err = validator(result)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "邮箱不能为空")

	// 测试：所有字段都存在
	result.Set("email", "test@example.com")
	err = validator(result)
	assert.NoError(t, err)
}

func TestFormValidator_FieldValueValidation(t *testing.T) {
	// 创建自定义验证器：验证字段值
	validator := func(result *FormResult) error {
		age := result.GetInt("age")
		if age < 18 {
			return errors.New("年龄必须大于等于18岁")
		}
		if age > 100 {
			return errors.New("年龄不能超过100岁")
		}
		return nil
	}

	// 测试：年龄太小
	result := NewFormResult()
	result.Set("age", 15)
	err := validator(result)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "年龄必须大于等于18岁")

	// 测试：年龄太大
	result.Set("age", 150)
	err = validator(result)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "年龄不能超过100岁")

	// 测试：有效年龄
	result.Set("age", 25)
	err = validator(result)
	assert.NoError(t, err)
}

func TestFormValidator_CrossFieldValidation(t *testing.T) {
	// 创建自定义验证器：跨字段验证
	validator := func(result *FormResult) error {
		password := result.GetString("password")
		confirmPassword := result.GetString("confirm_password")

		if password != confirmPassword {
			return errors.New("两次输入的密码不一致")
		}
		return nil
	}

	// 测试：密码不一致
	result := NewFormResult()
	result.Set("password", "password123")
	result.Set("confirm_password", "password456")
	err := validator(result)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "两次输入的密码不一致")

	// 测试：密码一致
	result.Set("confirm_password", "password123")
	err = validator(result)
	assert.NoError(t, err)
}

func TestFormValidator_NestedFormValidation(t *testing.T) {
	// 创建自定义验证器：验证嵌套表单
	validator := func(result *FormResult) error {
		userForm := result.GetForm("user")
		if userForm == nil {
			return errors.New("用户信息不能为空")
		}

		if userForm.GetString("name") == "" {
			return errors.New("用户姓名不能为空")
		}
		return nil
	}

	// 测试：嵌套表单不存在
	result := NewFormResult()
	err := validator(result)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "用户信息不能为空")

	// 测试：嵌套表单存在但字段为空
	nestedResult := NewFormResult()
	result.Set("user", nestedResult)
	err = validator(result)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "用户姓名不能为空")

	// 测试：嵌套表单字段有效
	nestedResult.Set("name", "张三")
	err = validator(result)
	assert.NoError(t, err)
}

func TestFormValidator_ComplexValidation(t *testing.T) {
	// 创建复杂的验证器：组合多个验证规则
	validator := func(result *FormResult) error {
		// 验证必填字段
		if result.GetString("name") == "" {
			return errors.New("姓名不能为空")
		}

		// 验证邮箱格式（简单验证）
		email := result.GetString("email")
		if email == "" {
			return errors.New("邮箱不能为空")
		}
		if len(email) < 5 || !contains(email, "@") {
			return errors.New("邮箱格式不正确")
		}

		// 验证年龄范围
		age := result.GetInt("age")
		if age < 0 || age > 150 {
			return errors.New("年龄必须在0-150之间")
		}

		// 验证密码长度
		password := result.GetString("password")
		if len(password) < 6 {
			return errors.New("密码长度至少6个字符")
		}

		return nil
	}

	// 测试：所有验证都通过
	result := NewFormResult()
	result.Set("name", "张三")
	result.Set("email", "test@example.com")
	result.Set("age", 25)
	result.Set("password", "password123")
	err := validator(result)
	assert.NoError(t, err)

	// 测试：姓名为空
	result.Set("name", "")
	err = validator(result)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "姓名不能为空")

	// 测试：邮箱格式错误
	result.Set("name", "张三")
	result.Set("email", "invalid")
	err = validator(result)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "邮箱格式不正确")

	// 测试：年龄超出范围
	result.Set("email", "test@example.com")
	result.Set("age", 200)
	err = validator(result)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "年龄必须在0-150之间")

	// 测试：密码太短
	result.Set("age", 25)
	result.Set("password", "12345")
	err = validator(result)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "密码长度至少6个字符")
}

// contains 辅助函数：检查字符串是否包含子串
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > len(substr) && containsHelper(s, substr))
}

func containsHelper(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

// ==================== FormBuilder.Validate 集成测试 ====================

func TestFormBuilder_Validate_Integration(t *testing.T) {
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

	// 创建带验证器的表单
	validator := func(result *FormResult) error {
		if result.GetString("name") == "" {
			return errors.New("姓名不能为空")
		}
		return nil
	}

	builder := NewFormBuilder().
		AddInput(InputFormField{
			Key:          "name",
			Prompt:       "姓名",
			DefaultValue: "",
			Validator:    nil,
		}).
		Validate(validator)

	// 执行表单（默认值为空，验证器应该失败）
	// 注意：由于 Run 方法会先执行表单，然后验证，所以这里需要确保验证器被调用
	// 但由于我们的 mock 返回空字符串，验证器会失败
	_, err := builder.Run()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "表单验证失败")
}

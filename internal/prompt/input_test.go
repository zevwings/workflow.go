//go:build test

package prompt_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
	"github.com/zevwings/workflow/internal/prompt"
	"github.com/zevwings/workflow/internal/testutils"
)

// InputTestSuite 是 Input 提示测试套件
// 所有测试共享相同的测试设置和辅助方法
type InputTestSuite struct {
	testutils.BasePromptTestSuite
}

// TestInputTestSuite 运行整个测试套件
func TestInputTestSuite(t *testing.T) {
	suite.Run(t, new(InputTestSuite))
}

// ==================== AskInput 测试 ====================

func (s *InputTestSuite) TestAskInput_Basic() {
	result, err := s.RunPromptTest(
		func(pt *testutils.PromptTestBuilder) *testutils.PromptTestBuilder {
			return pt.WithInput("test@example.com")
		},
		func() (interface{}, error) {
			return prompt.AskInput(prompt.InputField{
				Message:      "请输入邮箱",
				DefaultValue: "",
			})
		},
	)
	s.NoError(err)
	s.Equal("test@example.com", result)
}

func (s *InputTestSuite) TestAskInput_WithDefaultValue() {
	result, err := s.RunPromptTest(
		func(pt *testutils.PromptTestBuilder) *testutils.PromptTestBuilder {
			return pt.WithInput("") // 空输入表示回车使用默认值
		},
		func() (interface{}, error) {
			return prompt.AskInput(prompt.InputField{
				Message:      "请输入邮箱",
				DefaultValue: "default@example.com",
			})
		},
	)
	s.NoError(err)
	s.Equal("default@example.com", result)
}

func (s *InputTestSuite) TestAskInput_WithValidator() {
	result, err := s.RunPromptTest(
		func(pt *testutils.PromptTestBuilder) *testutils.PromptTestBuilder {
			return pt.WithInput("test@example.com")
		},
		func() (interface{}, error) {
			validator := prompt.ValidateEmail()
			return prompt.AskInput(prompt.InputField{
				Message:      "请输入邮箱",
				DefaultValue: "",
				Validator:    validator,
			})
		},
	)
	s.NoError(err)
	s.Equal("test@example.com", result)
}

func (s *InputTestSuite) TestAskInput_WithPlaceholder() {
	result, err := s.RunPromptTest(
		func(pt *testutils.PromptTestBuilder) *testutils.PromptTestBuilder {
			return pt.WithInput("user@example.com")
		},
		func() (interface{}, error) {
			return prompt.Input().
				Prompt("请输入邮箱").
				Placeholder("example@domain.com").
				Run()
		},
	)
	s.NoError(err)
	s.Equal("user@example.com", result)
}

// ==================== AskPassword 测试 ====================

func (s *InputTestSuite) TestAskPassword_Input() {
	result, err := s.RunPromptTest(
		func(pt *testutils.PromptTestBuilder) *testutils.PromptTestBuilder {
			return pt.WithInput("secret123")
		},
		func() (interface{}, error) {
			return prompt.AskPassword(prompt.PasswordField{
				Message:      "请输入密码",
				DefaultValue: "",
			})
		},
	)
	s.NoError(err)
	s.Equal("secret123", result)
	// 注意：密码输入时输出应该是隐藏的，但这里我们主要验证功能是否正常
}

// ==================== Input Validator 测试 ====================

func (s *InputTestSuite) TestInput_Validator_Email() {
	testCases := []struct {
		name         string
		inputs       []string // 多个输入，第一个可能无效
		expectResult string
	}{
		{"有效邮箱", []string{"test@example.com"}, "test@example.com"},
		{"有效邮箱2", []string{"user@domain.co.uk"}, "user@domain.co.uk"},
		{"无效邮箱1", []string{"invalid", "test@example.com"}, "test@example.com"},
		{"无效邮箱2", []string{"test@", "test@example.com"}, "test@example.com"},
		{"无效邮箱3", []string{"@example.com", "test@example.com"}, "test@example.com"},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			result, err := s.RunPromptTest(
				func(pt *testutils.PromptTestBuilder) *testutils.PromptTestBuilder {
					// 对于多个输入的情况，在输入之间添加延迟，确保验证失败后重新提示时输入被正确处理
					builder := pt.WithStartDelay(100 * time.Millisecond)
					for i, input := range tc.inputs {
						builder = builder.WithInput(input)
						// 如果不是最后一个输入，添加延迟以确保验证失败后的重新提示被正确处理
						if i < len(tc.inputs)-1 {
							builder = builder.WithDelay(200 * time.Millisecond)
						}
					}
					return builder
				},
				func() (interface{}, error) {
					validator := prompt.ValidateEmail()
					return prompt.AskInput(prompt.InputField{
						Message:      "请输入邮箱",
						DefaultValue: "",
						Validator:    validator,
					})
				},
			)
			s.NoError(err)
			s.Equal(tc.expectResult, result)
		})
	}
}

func (s *InputTestSuite) TestInput_Validator_Required() {
	result, err := s.RunPromptTest(
		func(pt *testutils.PromptTestBuilder) *testutils.PromptTestBuilder {
			return pt.WithInputs("", "non-empty") // 先输入空值（会被验证器拒绝），然后输入有效值
		},
		func() (interface{}, error) {
			validator := prompt.ValidateRequired()
			return prompt.AskInput(prompt.InputField{
				Message:      "请输入内容",
				DefaultValue: "",
				Validator:    validator,
			})
		},
	)
	s.NoError(err)
	s.Equal("non-empty", result)
}

func (s *InputTestSuite) TestInput_Validator_MinLength() {
	result, err := s.RunPromptTest(
		func(pt *testutils.PromptTestBuilder) *testutils.PromptTestBuilder {
			return pt.
				WithStartDelay(100 * time.Millisecond).
				WithInput("ab").
				WithDelay(200 * time.Millisecond).
				WithInput("abcde") // 先输入过短的值（会被验证器拒绝），然后输入有效值
		},
		func() (interface{}, error) {
			validator := prompt.ValidateMinLength(5)
			return prompt.AskInput(prompt.InputField{
				Message:      "请输入至少5个字符",
				DefaultValue: "",
				Validator:    validator,
			})
		},
	)
	s.NoError(err)
	s.Equal("abcde", result)
}

func (s *InputTestSuite) TestInput_Validator_MaxLength() {
	result, err := s.RunPromptTest(
		func(pt *testutils.PromptTestBuilder) *testutils.PromptTestBuilder {
			return pt.
				WithStartDelay(100 * time.Millisecond).
				WithInput("toolongvalue").
				WithDelay(200 * time.Millisecond).
				WithInput("valid") // 先输入过长的值（会被验证器拒绝），然后输入有效值
		},
		func() (interface{}, error) {
			validator := prompt.ValidateMaxLength(10)
			return prompt.AskInput(prompt.InputField{
				Message:      "请输入最多10个字符",
				DefaultValue: "",
				Validator:    validator,
			})
		},
	)
	s.NoError(err)
	s.Equal("valid", result)
}

func (s *InputTestSuite) TestInput_Validator_Length() {
	result, err := s.RunPromptTest(
		func(pt *testutils.PromptTestBuilder) *testutils.PromptTestBuilder {
			return pt.
				WithStartDelay(100 * time.Millisecond).
				WithInput("ab").
				WithDelay(200 * time.Millisecond).
				WithInput("abcde") // 先输入无效值（会被验证器拒绝），然后输入有效值
		},
		func() (interface{}, error) {
			validator := prompt.ValidateLength(3, 10)
			return prompt.AskInput(prompt.InputField{
				Message:      "请输入3-10个字符",
				DefaultValue: "",
				Validator:    validator,
			})
		},
	)
	s.NoError(err)
	s.Equal("abcde", result)
}

func (s *InputTestSuite) TestInput_Validator_Regex() {
	result, err := s.RunPromptTest(
		func(pt *testutils.PromptTestBuilder) *testutils.PromptTestBuilder {
			return pt.
				WithStartDelay(100 * time.Millisecond).
				WithInput("invalid").
				WithDelay(200 * time.Millisecond).
				WithInput("ABC123") // 先输入无效值（会被验证器拒绝），然后输入有效值
		},
		func() (interface{}, error) {
			validator := prompt.ValidateRegex(`^[A-Z0-9]+$`, "只能包含大写字母和数字")
			return prompt.AskInput(prompt.InputField{
				Message:      "请输入大写字母和数字",
				DefaultValue: "",
				Validator:    validator,
			})
		},
	)
	s.NoError(err)
	s.Equal("ABC123", result)
}

func (s *InputTestSuite) TestInput_Validator_URL() {
	result, err := s.RunPromptTest(
		func(pt *testutils.PromptTestBuilder) *testutils.PromptTestBuilder {
			return pt.
				WithStartDelay(100 * time.Millisecond).
				WithInput("not-a-url").
				WithDelay(200 * time.Millisecond).
				WithInput("https://example.com") // 先输入无效值（会被验证器拒绝），然后输入有效值
		},
		func() (interface{}, error) {
			validator := prompt.ValidateURL()
			return prompt.AskInput(prompt.InputField{
				Message:      "请输入URL",
				DefaultValue: "",
				Validator:    validator,
			})
		},
	)
	s.NoError(err)
	s.Equal("https://example.com", result)
}

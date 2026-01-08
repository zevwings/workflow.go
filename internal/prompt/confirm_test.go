//go:build test

package prompt_test

import (
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/zevwings/workflow/internal/prompt"
	"github.com/zevwings/workflow/internal/testutils"
)

// ConfirmTestSuite 是 Confirm 提示测试套件
// 所有测试共享相同的测试设置和辅助方法
type ConfirmTestSuite struct {
	testutils.BasePromptTestSuite
}

// TestConfirmTestSuite 运行整个测试套件
func TestConfirmTestSuite(t *testing.T) {
	suite.Run(t, new(ConfirmTestSuite))
}

// ==================== AskConfirm 测试 ====================

func (s *ConfirmTestSuite) TestAskConfirm_Yes() {
	result, err := s.RunPromptTest(
		func(pt *testutils.PromptTestBuilder) *testutils.PromptTestBuilder {
			return pt.WithInput("y")
		},
		func() (interface{}, error) {
			return prompt.AskConfirm("是否继续？", false)
		},
	)
	s.NoError(err)
	s.True(result.(bool))
}

func (s *ConfirmTestSuite) TestAskConfirm_No() {
	result, err := s.RunPromptTest(
		func(pt *testutils.PromptTestBuilder) *testutils.PromptTestBuilder {
			return pt.WithInput("n")
		},
		func() (interface{}, error) {
			return prompt.AskConfirm("是否继续？", true)
		},
	)
	s.NoError(err)
	s.False(result.(bool))
}

func (s *ConfirmTestSuite) TestAskConfirm_DefaultYes() {
	result, err := s.RunPromptTest(
		func(pt *testutils.PromptTestBuilder) *testutils.PromptTestBuilder {
			return pt.WithInput("") // 空输入表示回车使用默认值
		},
		func() (interface{}, error) {
			return prompt.AskConfirm("是否继续？", true)
		},
	)
	s.NoError(err)
	s.True(result.(bool)) // 使用默认值 true
}

func (s *ConfirmTestSuite) TestAskConfirm_DefaultNo() {
	result, err := s.RunPromptTest(
		func(pt *testutils.PromptTestBuilder) *testutils.PromptTestBuilder {
			return pt.WithInput("") // 空输入表示回车使用默认值
		},
		func() (interface{}, error) {
			return prompt.AskConfirm("是否继续？", false)
		},
	)
	s.NoError(err)
	s.False(result.(bool)) // 使用默认值 false
}

func (s *ConfirmTestSuite) TestAskConfirm_UpperCase() {
	result, err := s.RunPromptTest(
		func(pt *testutils.PromptTestBuilder) *testutils.PromptTestBuilder {
			return pt.WithInput("Y")
		},
		func() (interface{}, error) {
			return prompt.AskConfirm("是否继续？", false)
		},
	)
	s.NoError(err)
	s.True(result.(bool))
}

// ==================== Confirm Builder 测试 ====================

func (s *ConfirmTestSuite) TestConfirm_Builder() {
	result, err := s.RunPromptTest(
		func(pt *testutils.PromptTestBuilder) *testutils.PromptTestBuilder {
			return pt.WithInput("y")
		},
		func() (interface{}, error) {
			return prompt.Confirm().
				Prompt("是否继续？").
				Default(false).
				Run()
		},
	)
	s.NoError(err)
	s.True(result.(bool))
}

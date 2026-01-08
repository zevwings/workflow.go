//go:build test

package prompt_test

import (
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/zevwings/workflow/internal/prompt"
	"github.com/zevwings/workflow/internal/testutils"
)

// SelectTestSuite 是 Select 提示测试套件
// 所有测试共享相同的测试设置和辅助方法
type SelectTestSuite struct {
	testutils.BasePromptTestSuite
	options []string // 共享的测试选项
}

// SetupTest 在每个测试运行前执行
func (s *SelectTestSuite) SetupTest() {
	s.options = []string{"选项1", "选项2", "选项3"}
}

// TestSelectTestSuite 运行整个测试套件
func TestSelectTestSuite(t *testing.T) {
	suite.Run(t, new(SelectTestSuite))
}

// ==================== AskSelect 测试 ====================

func (s *SelectTestSuite) TestAskSelect_Basic() {
	result, err := s.RunPromptTest(
		func(pt *testutils.PromptTestBuilder) *testutils.PromptTestBuilder {
			return pt.
				WithRawInput("\x1b[B"). // 下箭头转义序列: ESC [ B
				WithInput("")
		},
		func() (interface{}, error) {
			return prompt.AskSelect("请选择", s.options, 0)
		},
	)
	s.NoError(err)
	s.Equal(1, result.(int)) // 选中第二个选项（索引1）
}

func (s *SelectTestSuite) TestAskSelect_ArrowKeys() {
	options := []string{"选项1", "选项2", "选项3", "选项4"}
	result, err := s.RunPromptTest(
		func(pt *testutils.PromptTestBuilder) *testutils.PromptTestBuilder {
			return pt.
				WithRawInput("\x1b[B"). // 下移
				WithRawInput("\x1b[B"). // 下移
				WithRawInput("\x1b[A"). // 上移
				WithInput("")           // 确认
		},
		func() (interface{}, error) {
			return prompt.AskSelect("请选择", options, 0)
		},
	)
	s.NoError(err)
	s.Equal(1, result.(int)) // 最终选中第二个选项
}

func (s *SelectTestSuite) TestAskSelect_DefaultIndex() {
	result, err := s.RunPromptTest(
		func(pt *testutils.PromptTestBuilder) *testutils.PromptTestBuilder {
			return pt.WithInput("") // 空输入表示回车使用默认值
		},
		func() (interface{}, error) {
			return prompt.AskSelect("请选择", s.options, 2)
		},
	)
	s.NoError(err)
	s.Equal(2, result.(int)) // 使用默认索引
}

func (s *SelectTestSuite) TestAskSelect_EmptyOptions() {
	options := []string{}
	_, err := prompt.AskSelect("请选择", options, 0)
	s.Error(err)
	s.Contains(err.Error(), "选项列表不能为空")
}

func (s *SelectTestSuite) TestAskSelect_InvalidDefaultIndex() {
	result, err := s.RunPromptTest(
		func(pt *testutils.PromptTestBuilder) *testutils.PromptTestBuilder {
			return pt.WithInput("") // 空输入表示回车使用默认值
		},
		func() (interface{}, error) {
			return prompt.AskSelect("请选择", s.options, 10)
		},
	)
	s.NoError(err)
	s.Equal(0, result.(int)) // 自动调整为第一个选项
}

// ==================== Select Builder 测试 ====================

func (s *SelectTestSuite) TestSelect_Builder() {
	result, err := s.RunPromptTest(
		func(pt *testutils.PromptTestBuilder) *testutils.PromptTestBuilder {
			return pt.
				WithRawInput("\x1b[B").
				WithInput("")
		},
		func() (interface{}, error) {
			return prompt.Select().
				Prompt("请选择").
				Options(s.options).
				Default(0).
				Run()
		},
	)
	s.NoError(err)
	s.Equal(1, result.(int))
}

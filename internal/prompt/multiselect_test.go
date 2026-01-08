//go:build test

package prompt_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
	"github.com/zevwings/workflow/internal/prompt"
	"github.com/zevwings/workflow/internal/testutils"
)

// MultiSelectTestSuite 是 MultiSelect 提示测试套件
// 所有测试共享相同的测试设置和辅助方法
type MultiSelectTestSuite struct {
	testutils.BasePromptTestSuite
	options []string // 共享的测试选项
}

// SetupTest 在每个测试运行前执行
func (s *MultiSelectTestSuite) SetupTest() {
	s.options = []string{"选项1", "选项2", "选项3"}
}

// TestMultiSelectTestSuite 运行整个测试套件
func TestMultiSelectTestSuite(t *testing.T) {
	suite.Run(t, new(MultiSelectTestSuite))
}

// ==================== AskMultiSelect 测试 ====================

func (s *MultiSelectTestSuite) TestAskMultiSelect_Basic() {
	result, err := s.RunPromptTest(
		func(pt *testutils.PromptTestBuilder) *testutils.PromptTestBuilder {
			return pt.
				WithStartDelay(100 * time.Millisecond).
				WithRawInput("\x1b[B").           // 下移到选项2
				WithDelay(50 * time.Millisecond). // 延迟确保移动完成
				WithInput(" ").                   // 选中选项2
				WithDelay(50 * time.Millisecond). // 延迟确保选中完成
				WithRawInput("\x1b[B").           // 下移到选项3
				WithDelay(50 * time.Millisecond). // 延迟确保移动完成
				WithInput(" ").                   // 选中选项3
				WithDelay(50 * time.Millisecond). // 延迟确保选中完成
				WithInput("")                     // 确认
		},
		func() (interface{}, error) {
			return prompt.AskMultiSelect("请选择", s.options, []int{})
		},
	)
	s.NoError(err)
	s.Equal([]int{1, 2}, result.([]int)) // 选中选项2和选项3
}

func (s *MultiSelectTestSuite) TestAskMultiSelect_DefaultSelected() {
	result, err := s.RunPromptTest(
		func(pt *testutils.PromptTestBuilder) *testutils.PromptTestBuilder {
			return pt.WithInput("") // 空输入表示回车使用默认值
		},
		func() (interface{}, error) {
			return prompt.AskMultiSelect("请选择", s.options, []int{0, 2})
		},
	)
	s.NoError(err)
	s.Equal([]int{0, 2}, result.([]int))
}

func (s *MultiSelectTestSuite) TestAskMultiSelect_ToggleSelection() {
	result, err := s.RunPromptTest(
		func(pt *testutils.PromptTestBuilder) *testutils.PromptTestBuilder {
			return pt.
				WithStartDelay(100 * time.Millisecond).
				WithInput(" ").                   // 选中选项1
				WithDelay(50 * time.Millisecond). // 延迟确保选中完成
				WithInput(" ").                   // 取消选项1
				WithDelay(50 * time.Millisecond). // 延迟确保取消完成
				WithRawInput("\x1b[B").           // 下移到选项2
				WithDelay(50 * time.Millisecond). // 延迟确保移动完成
				WithInput(" ").                   // 选中选项2
				WithDelay(50 * time.Millisecond). // 延迟确保选中完成
				WithInput("")                     // 确认
		},
		func() (interface{}, error) {
			return prompt.AskMultiSelect("请选择", s.options, []int{})
		},
	)
	s.NoError(err)
	s.Equal([]int{1}, result.([]int)) // 只选中选项2
}

func (s *MultiSelectTestSuite) TestAskMultiSelect_EmptySelection() {
	result, err := s.RunPromptTest(
		func(pt *testutils.PromptTestBuilder) *testutils.PromptTestBuilder {
			return pt.WithInput("") // 空输入表示不选择任何选项，直接回车
		},
		func() (interface{}, error) {
			return prompt.AskMultiSelect("请选择", s.options, []int{})
		},
	)
	s.NoError(err)
	// 注意：可能返回 nil 或空 slice，都表示未选择
	if result == nil {
		s.Nil(result) // 接受 nil
	} else {
		actualResult, ok := result.([]int)
		s.True(ok, "结果应该是 []int 类型")
		s.Empty(actualResult) // 或者空 slice
	}
}

func (s *MultiSelectTestSuite) TestAskMultiSelect_EmptyOptions() {
	options := []string{}
	_, err := prompt.AskMultiSelect("请选择", options, []int{})
	s.Error(err)
	s.Contains(err.Error(), "选项列表不能为空")
}

func (s *MultiSelectTestSuite) TestAskMultiSelect_InvalidDefaultIndices() {
	result, err := s.RunPromptTest(
		func(pt *testutils.PromptTestBuilder) *testutils.PromptTestBuilder {
			return pt.WithInput("") // 空输入表示回车使用默认值
		},
		func() (interface{}, error) {
			return prompt.AskMultiSelect("请选择", s.options, []int{0, 10, 2})
		},
	)
	s.NoError(err)
	s.Equal([]int{0, 2}, result.([]int)) // 无效索引 10 应该被过滤
}

// ==================== MultiSelect Builder 测试 ====================

func (s *MultiSelectTestSuite) TestMultiSelect_Builder() {
	result, err := s.RunPromptTest(
		func(pt *testutils.PromptTestBuilder) *testutils.PromptTestBuilder {
			return pt.
				WithRawInput("\x1b[B").
				WithInput(" ").
				WithInput("")
		},
		func() (interface{}, error) {
			return prompt.MultiSelect().
				Prompt("请选择").
				Options(s.options).
				Default([]int{}).
				Run()
		},
	)
	s.NoError(err)
	s.Equal([]int{1}, result.([]int))
}

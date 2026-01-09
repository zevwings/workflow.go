//go:build test

package selectpkg

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zevwings/workflow/internal/prompt/common"
	"github.com/zevwings/workflow/internal/prompt/io"
	"github.com/zevwings/workflow/internal/testutils"
)

// newSelectConfig 创建 SelectConfig（测试辅助函数）
// 用于简化测试代码中的配置创建
func newSelectConfig(message string, options []string, defaultIndex int, terminal io.TerminalIO) SelectConfig {
	return SelectConfig{
		BasePromptConfig: common.BasePromptConfig{
			Message:  message,
			Config:   testutils.NewPromptConfig(),
			Terminal: terminal,
		},
		Options:      options,
		DefaultIndex: defaultIndex,
	}
}

// TestSelect_EmptyOptions 验证空选项时直接返回错误
func TestSelect_EmptyOptions(t *testing.T) {
	cfg := testutils.NewPromptConfig()

	mockTerminal := io.NewMockTerminal([]byte{})
	index, err := Select(SelectConfig{
		BasePromptConfig: common.BasePromptConfig{
			Message:  "请选择",
			Config:   cfg,
			Terminal: mockTerminal,
		},
		Options:      []string{},
		DefaultIndex: 0,
	})
	assert.Error(t, err)
	assert.Equal(t, -1, index)
}

// Test_selectFallback_ValidAndInvalidInput 验证回退模式下合法与非法输入行为
func Test_selectFallback_ValidAndInvalidInput(t *testing.T) {
	options := []string{"A", "B", "C"}
	cfg := testutils.NewPromptConfig()

	// 1) 非法输入时应返回默认索引
	mockTerminal := io.NewMockTerminalWithLines([]string{"invalid"})
	selectCfg := SelectConfig{
		BasePromptConfig: common.BasePromptConfig{
			Message:  "请选择",
			Config:   cfg,
			Terminal: mockTerminal,
		},
		Options:      options,
		DefaultIndex: 1,
	}
	index, err := selectFallback(selectCfg)
	assert.NoError(t, err)
	assert.Equal(t, 1, index)

	// 2) 合法数字输入时应返回对应索引
	mockTerminal = io.NewMockTerminalWithLines([]string{"2"})
	selectCfg = SelectConfig{
		BasePromptConfig: common.BasePromptConfig{
			Message:  "请选择",
			Config:   cfg,
			Terminal: mockTerminal,
		},
		Options:      options,
		DefaultIndex: 0,
	}
	index, err = selectFallback(selectCfg)
	assert.NoError(t, err)
	assert.Equal(t, 1, index)
}

// Test_selectFallback_OutOfRangeInput 验证超出范围的输入会返回默认索引
func Test_selectFallback_OutOfRangeInput(t *testing.T) {
	options := []string{"A", "B", "C"}
	cfg := testutils.NewPromptConfig()

	// 测试超出上限
	mockTerminal := io.NewMockTerminalWithLines([]string{"10"})
	selectCfg := SelectConfig{
		BasePromptConfig: common.BasePromptConfig{
			Message:  "请选择",
			Config:   cfg,
			Terminal: mockTerminal,
		},
		Options:      options,
		DefaultIndex: 1,
	}
	index, err := selectFallback(selectCfg)
	assert.NoError(t, err)
	assert.Equal(t, 1, index) // 返回默认索引

	// 测试小于下限
	mockTerminal = io.NewMockTerminalWithLines([]string{"0"})
	selectCfg = SelectConfig{
		BasePromptConfig: common.BasePromptConfig{
			Message:  "请选择",
			Config:   cfg,
			Terminal: mockTerminal,
		},
		Options:      options,
		DefaultIndex: 2,
	}
	index, err = selectFallback(selectCfg)
	assert.NoError(t, err)
	assert.Equal(t, 2, index) // 返回默认索引
}

// Test_selectFallback_EmptyInput 验证空输入时返回默认索引
func Test_selectFallback_EmptyInput(t *testing.T) {
	options := []string{"A", "B", "C"}
	cfg := testutils.NewPromptConfig()

	// 空输入会触发 ReadLine 错误，返回默认索引
	mockTerminal := io.NewMockTerminalWithLines([]string{})
	selectCfg := SelectConfig{
		BasePromptConfig: common.BasePromptConfig{
			Message:  "请选择",
			Config:   cfg,
			Terminal: mockTerminal,
		},
		Options:      options,
		DefaultIndex: 1,
	}
	index, err := selectFallback(selectCfg)
	assert.NoError(t, err)
	assert.Equal(t, 1, index) // 返回默认索引
}

// TestSelect_InvalidDefaultIndex 验证无效的默认索引会被调整为 0
func TestSelect_InvalidDefaultIndex(t *testing.T) {
	cfg := testutils.NewPromptConfig()

	options := []string{"A", "B", "C"}

	// 测试负索引（使用 MockTerminal，直接回车）
	mockTerminal := io.NewMockTerminal([]byte{'\r'})
	index, err := Select(SelectConfig{
		BasePromptConfig: common.BasePromptConfig{
			Message:  "请选择",
			Config:   cfg,
			Terminal: mockTerminal,
		},
		Options:      options,
		DefaultIndex: -1,
	})
	// 无效索引会被调整为 0，回车选择索引 0
	assert.NoError(t, err)
	assert.Equal(t, 0, index)

	// 测试超出范围的索引
	mockTerminal = io.NewMockTerminal([]byte{'\r'})
	index, err = Select(SelectConfig{
		BasePromptConfig: common.BasePromptConfig{
			Message:  "请选择",
			Config:   cfg,
			Terminal: mockTerminal,
		},
		Options:      options,
		DefaultIndex: 10,
	})
	// 无效索引会被调整为 0，回车选择索引 0
	assert.NoError(t, err)
	assert.Equal(t, 0, index)
}

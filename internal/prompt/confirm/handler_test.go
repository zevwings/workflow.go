//go:build test

package confirm

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// ==================== ConfirmHandler 测试 ====================

func TestConfirmHandler_ProcessInput_Enter(t *testing.T) {
	cfg := Config{
		FormatPrompt: func(msg string) string { return msg },
		FormatAnswer: func(v string) string { return v },
	}

	// 默认 yes，回车应返回 true
	handler := NewConfirmHandler(true, cfg)
	result, shouldContinue, err := handler.ProcessInput('\r')
	assert.NoError(t, err)
	assert.False(t, shouldContinue)
	assert.NotNil(t, result)
	assert.True(t, *result)

	// 默认 no，回车应返回 false
	handler = NewConfirmHandler(false, cfg)
	result, shouldContinue, err = handler.ProcessInput('\n')
	assert.NoError(t, err)
	assert.False(t, shouldContinue)
	assert.NotNil(t, result)
	assert.False(t, *result)
}

func TestConfirmHandler_ProcessInput_Yes(t *testing.T) {
	cfg := Config{
		FormatPrompt: func(msg string) string { return msg },
		FormatAnswer: func(v string) string { return v },
	}

	handler := NewConfirmHandler(false, cfg)

	// 测试 'y'
	result, shouldContinue, err := handler.ProcessInput('y')
	assert.NoError(t, err)
	assert.False(t, shouldContinue)
	assert.NotNil(t, result)
	assert.True(t, *result)

	// 测试 'Y'
	result, shouldContinue, err = handler.ProcessInput('Y')
	assert.NoError(t, err)
	assert.False(t, shouldContinue)
	assert.NotNil(t, result)
	assert.True(t, *result)
}

func TestConfirmHandler_ProcessInput_No(t *testing.T) {
	cfg := Config{
		FormatPrompt: func(msg string) string { return msg },
		FormatAnswer: func(v string) string { return v },
	}

	handler := NewConfirmHandler(true, cfg)

	// 测试 'n'
	result, shouldContinue, err := handler.ProcessInput('n')
	assert.NoError(t, err)
	assert.False(t, shouldContinue)
	assert.NotNil(t, result)
	assert.False(t, *result)

	// 测试 'N'
	result, shouldContinue, err = handler.ProcessInput('N')
	assert.NoError(t, err)
	assert.False(t, shouldContinue)
	assert.NotNil(t, result)
	assert.False(t, *result)
}

func TestConfirmHandler_ProcessInput_CtrlC(t *testing.T) {
	cfg := Config{
		FormatPrompt: func(msg string) string { return msg },
		FormatAnswer: func(v string) string { return v },
	}

	handler := NewConfirmHandler(true, cfg)
	result, shouldContinue, err := handler.ProcessInput(3) // Ctrl+C
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "用户取消输入")
	assert.False(t, shouldContinue)
	assert.Nil(t, result)
}

func TestConfirmHandler_ProcessInput_InvalidChar(t *testing.T) {
	cfg := Config{
		FormatPrompt: func(msg string) string { return msg },
		FormatAnswer: func(v string) string { return v },
	}

	handler := NewConfirmHandler(true, cfg)

	// 测试无效字符（应继续等待）
	result, shouldContinue, err := handler.ProcessInput('x')
	assert.NoError(t, err)
	assert.True(t, shouldContinue)
	assert.Nil(t, result)

	result, shouldContinue, err = handler.ProcessInput('1')
	assert.NoError(t, err)
	assert.True(t, shouldContinue)
	assert.Nil(t, result)
}

func TestConfirmHandler_FormatPromptText(t *testing.T) {
	cfg := Config{
		FormatPrompt: func(msg string) string { return "[" + msg + "]" },
		FormatAnswer: func(v string) string { return v },
	}

	// 默认 yes
	handler := NewConfirmHandler(true, cfg)
	text := handler.FormatPromptText("test")
	assert.Contains(t, text, "[test]")
	assert.Contains(t, text, "[Y/n]")

	// 默认 no
	handler = NewConfirmHandler(false, cfg)
	text = handler.FormatPromptText("test")
	assert.Contains(t, text, "[test]")
	assert.Contains(t, text, "[y/N]")
}

func TestConfirmHandler_FormatAnswer(t *testing.T) {
	cfg := Config{
		FormatPrompt: func(msg string) string { return msg },
		FormatAnswer: func(v string) string { return "[" + v + "]" },
	}

	handler := NewConfirmHandler(true, cfg)

	answer := handler.FormatAnswer(true)
	assert.Equal(t, "[yes]", answer)

	answer = handler.FormatAnswer(false)
	assert.Equal(t, "[no]", answer)
}

func TestConfirmHandler_ProcessLineInput(t *testing.T) {
	cfg := Config{
		FormatPrompt: func(msg string) string { return msg },
		FormatAnswer: func(v string) string { return v },
	}

	handler := NewConfirmHandler(true, cfg)

	// 空输入
	result, err := handler.ProcessLineInput("")
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.True(t, *result) // 使用默认值

	// "y"
	result, err = handler.ProcessLineInput("y")
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.True(t, *result)

	// "yes"
	result, err = handler.ProcessLineInput("yes")
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.True(t, *result)

	// "Y"
	result, err = handler.ProcessLineInput("Y")
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.True(t, *result)

	// "n"
	handler = NewConfirmHandler(true, cfg)
	result, err = handler.ProcessLineInput("n")
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.False(t, *result)

	// "no"
	result, err = handler.ProcessLineInput("no")
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.False(t, *result)

	// 非法输入
	handler = NewConfirmHandler(false, cfg)
	result, err = handler.ProcessLineInput("invalid")
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.False(t, *result) // 使用默认值

	// 带空格的输入
	result, err = handler.ProcessLineInput("  y  ")
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.True(t, *result)
}

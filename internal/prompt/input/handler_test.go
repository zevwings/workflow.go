//go:build test

package input

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

// ==================== InputHandler 测试 ====================

func TestInputHandler_ProcessEscapeSequence(t *testing.T) {
	handler := NewInputHandler(nil)

	// 右箭头
	direction, shouldUpdate := handler.ProcessEscapeSequence('[', 'C')
	assert.Equal(t, "right", direction)
	assert.True(t, shouldUpdate)

	// 左箭头
	direction, shouldUpdate = handler.ProcessEscapeSequence('[', 'D')
	assert.Equal(t, "left", direction)
	assert.True(t, shouldUpdate)

	// 上箭头（忽略）
	direction, shouldUpdate = handler.ProcessEscapeSequence('[', 'A')
	assert.Equal(t, "none", direction)
	assert.False(t, shouldUpdate)

	// 下箭头（忽略）
	direction, shouldUpdate = handler.ProcessEscapeSequence('[', 'B')
	assert.Equal(t, "none", direction)
	assert.False(t, shouldUpdate)

	// 无效序列
	direction, shouldUpdate = handler.ProcessEscapeSequence('X', 'Y')
	assert.Equal(t, "none", direction)
	assert.False(t, shouldUpdate)
}

func TestInputHandler_ProcessArrowKey(t *testing.T) {
	handler := NewInputHandler(nil)

	// 右箭头，有 placeholder
	newPos, shouldClear := handler.ProcessArrowKey(0, 5, "right", true)
	assert.Equal(t, 0, newPos)
	assert.True(t, shouldClear)

	// 右箭头，无 placeholder，可以移动
	newPos, shouldClear = handler.ProcessArrowKey(2, 5, "right", false)
	assert.Equal(t, 3, newPos)
	assert.False(t, shouldClear)

	// 右箭头，无 placeholder，已在末尾
	newPos, shouldClear = handler.ProcessArrowKey(5, 5, "right", false)
	assert.Equal(t, 5, newPos)
	assert.False(t, shouldClear)

	// 左箭头，有 placeholder（无效）
	newPos, shouldClear = handler.ProcessArrowKey(0, 5, "left", true)
	assert.Equal(t, 0, newPos)
	assert.False(t, shouldClear)

	// 左箭头，无 placeholder，可以移动
	newPos, shouldClear = handler.ProcessArrowKey(3, 5, "left", false)
	assert.Equal(t, 2, newPos)
	assert.False(t, shouldClear)

	// 左箭头，无 placeholder，已在开头
	newPos, shouldClear = handler.ProcessArrowKey(0, 5, "left", false)
	assert.Equal(t, 0, newPos)
	assert.False(t, shouldClear)
}

func TestInputHandler_ProcessBackspace(t *testing.T) {
	handler := NewInputHandler(nil)

	// 正常退格
	value := []byte("hello")
	newValue, newPos := handler.ProcessBackspace(value, 3)
	assert.Equal(t, []byte("helo"), newValue)
	assert.Equal(t, 2, newPos)

	// 光标在开头，无法退格
	value = []byte("hello")
	newValue, newPos = handler.ProcessBackspace(value, 0)
	assert.Equal(t, []byte("hello"), newValue)
	assert.Equal(t, 0, newPos)
}

func TestInputHandler_ProcessChar(t *testing.T) {
	handler := NewInputHandler(nil)

	// 在末尾插入
	value := []byte("hel")
	newValue, newPos := handler.ProcessChar(value, 3, 'l')
	assert.Equal(t, []byte("hell"), newValue)
	assert.Equal(t, 4, newPos)

	// 在中间插入
	value = []byte("helo")
	newValue, newPos = handler.ProcessChar(value, 3, 'l')
	assert.Equal(t, []byte("hello"), newValue)
	assert.Equal(t, 4, newPos)
}

func TestInputHandler_ValidateInput(t *testing.T) {
	// 无验证器
	handler := NewInputHandler(nil)
	err := handler.ValidateInput("test")
	assert.NoError(t, err)

	// 有验证器，通过
	validator := func(s string) error {
		if len(s) < 5 {
			return errors.New("太短")
		}
		return nil
	}
	handler = NewInputHandler(validator)
	err = handler.ValidateInput("hello")
	assert.NoError(t, err)

	// 有验证器，失败
	err = handler.ValidateInput("hi")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "太短")
}

func TestInputHandler_CalculateCursorBackspaces(t *testing.T) {
	handler := NewInputHandler(nil)

	// 正常情况
	value := []byte("hello")
	backspaces := handler.CalculateCursorBackspaces(value, 2)
	assert.Equal(t, 3, backspaces) // "llo" 的宽度

	// 光标在末尾
	backspaces = handler.CalculateCursorBackspaces(value, 5)
	assert.Equal(t, 0, backspaces)

	// 光标超出范围
	backspaces = handler.CalculateCursorBackspaces(value, 10)
	assert.Equal(t, 0, backspaces)

	// 光标在开头
	backspaces = handler.CalculateCursorBackspaces(value, 0)
	assert.Equal(t, 5, backspaces) // "hello" 的宽度
}

func TestInputHandler_TrimInput(t *testing.T) {
	handler := NewInputHandler(nil)

	// 正常情况
	result := handler.TrimInput("  hello  ")
	assert.Equal(t, "hello", result)

	// 无空白
	result = handler.TrimInput("hello")
	assert.Equal(t, "hello", result)

	// 只有空白
	result = handler.TrimInput("   ")
	assert.Equal(t, "", result)
}

// ==================== PlaceholderHandler 测试 ====================

func TestPlaceholderHandler_HasPlaceholder(t *testing.T) {
	cfg := Config{
		FormatPlaceholder: func(text string) string { return text },
		FormatError:       func(msg string) string { return msg },
	}

	// 有 placeholder
	handler := NewPlaceholderHandler("example@domain.com", cfg)
	assert.True(t, handler.HasPlaceholder())

	// 无 placeholder
	handler = NewPlaceholderHandler("", cfg)
	assert.False(t, handler.HasPlaceholder())
}

func TestPlaceholderHandler_FormatPlaceholder(t *testing.T) {
	cfg := Config{
		FormatPlaceholder: func(text string) string { return "[" + text + "]" },
		FormatError:       func(msg string) string { return msg },
	}

	handler := NewPlaceholderHandler("test", cfg)
	formatted := handler.FormatPlaceholder()
	assert.Equal(t, "[test]", formatted)
}

func TestPlaceholderHandler_GetPlaceholderText(t *testing.T) {
	cfg := Config{
		FormatPlaceholder: func(text string) string { return "\x1b[31m" + text + "\x1b[0m" },
		FormatError:       func(msg string) string { return msg },
	}

	handler := NewPlaceholderHandler("test", cfg)
	text := handler.GetPlaceholderText()
	assert.Equal(t, "test", text) // 应该去除 ANSI 代码
}

func TestPlaceholderHandler_GetPlaceholderWidth(t *testing.T) {
	cfg := Config{
		FormatPlaceholder: func(text string) string { return text },
		FormatError:       func(msg string) string { return msg },
	}

	handler := NewPlaceholderHandler("hello", cfg)
	width := handler.GetPlaceholderWidth()
	assert.Equal(t, 5, width)
}

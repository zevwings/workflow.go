//go:build test

package input

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zevwings/workflow/internal/prompt/io"
)

// ==================== ReadWithPlaceholder 完整流程测试 ====================

func TestReadWithPlaceholder_WithMockTerminal_BasicInput(t *testing.T) {
	cfg := Config{
		FormatPlaceholder: func(text string) string { return text },
		FormatError:       func(msg string) string { return msg },
	}

	// 输入 "hello" 然后回车
	mockTerminal := io.NewMockTerminal([]byte{'h', 'e', 'l', 'l', 'o', '\r'})
	result, err := ReadWithPlaceholder("> ", "", nil, cfg, mockTerminal)

	assert.NoError(t, err)
	assert.Equal(t, "hello", result)
}

func TestReadWithPlaceholder_WithMockTerminal_Placeholder(t *testing.T) {
	cfg := Config{
		FormatPlaceholder: func(text string) string { return text },
		FormatError:       func(msg string) string { return msg },
	}

	// 有 placeholder，输入右箭头清除 placeholder，然后输入 "test" 回车
	mockTerminal := io.NewMockTerminal([]byte{0x1b, '[', 'C', 't', 'e', 's', 't', '\r'})
	result, err := ReadWithPlaceholder("> ", "example@domain.com", nil, cfg, mockTerminal)

	assert.NoError(t, err)
	assert.Equal(t, "test", result)
}

func TestReadWithPlaceholder_WithMockTerminal_Backspace(t *testing.T) {
	cfg := Config{
		FormatPlaceholder: func(text string) string { return text },
		FormatError:       func(msg string) string { return msg },
	}

	// 输入 "helo"，然后左移，插入 'l'，回车
	mockTerminal := io.NewMockTerminal([]byte{'h', 'e', 'l', 'o', 0x1b, '[', 'D', 'l', '\r'})
	result, err := ReadWithPlaceholder("> ", "", nil, cfg, mockTerminal)

	assert.NoError(t, err)
	assert.Equal(t, "hello", result)
}

func TestReadWithPlaceholder_WithMockTerminal_ValidatorError(t *testing.T) {
	cfg := Config{
		FormatPlaceholder: func(text string) string { return text },
		FormatError:       func(msg string) string { return "ERR: " + msg },
	}

	validator := func(s string) error {
		if len(s) < 5 {
			return errors.New("太短")
		}
		return nil
	}

	// 输入 "hi" 回车（验证失败，继续等待），然后退格两次清除，输入 "hello" 回车（验证通过）
	mockTerminal := io.NewMockTerminal([]byte{'h', 'i', '\r', 127, 127, 'h', 'e', 'l', 'l', 'o', '\r'})
	result, err := ReadWithPlaceholder("> ", "", validator, cfg, mockTerminal)

	assert.NoError(t, err)
	assert.Equal(t, "hello", result)
}

func TestReadWithPlaceholder_WithMockTerminal_CtrlC(t *testing.T) {
	cfg := Config{
		FormatPlaceholder: func(text string) string { return text },
		FormatError:       func(msg string) string { return msg },
	}

	mockTerminal := io.NewMockTerminal([]byte{3}) // Ctrl+C
	result, err := ReadWithPlaceholder("> ", "", nil, cfg, mockTerminal)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "用户取消输入")
	assert.Equal(t, "", result)
}

// ==================== ReadLineCore 完整流程测试 ====================

func TestReadLineCore_WithMockTerminal_BasicInput(t *testing.T) {
	echo := func(b []byte) string { return string(b) }
	formatError := func(msg string) string { return msg }

	// 输入 "hello" 然后回车
	mockTerminal := io.NewMockTerminal([]byte{'h', 'e', 'l', 'l', 'o', '\r'})
	result, err := ReadLineCore("> ", nil, echo, formatError, mockTerminal)

	assert.NoError(t, err)
	assert.Equal(t, "hello", result)
}

func TestReadLineCore_WithMockTerminal_PasswordMode(t *testing.T) {
	echo := func(b []byte) string {
		return "****" // 固定长度的掩码
	}
	formatError := func(msg string) string { return msg }

	// 输入 "secret" 然后回车
	mockTerminal := io.NewMockTerminal([]byte{'s', 'e', 'c', 'r', 'e', 't', '\r'})
	result, err := ReadLineCore("> ", nil, echo, formatError, mockTerminal)

	assert.NoError(t, err)
	assert.Equal(t, "secret", result)
	// 验证输出包含掩码
	assert.Contains(t, mockTerminal.GetOutput(), "****")
}

func TestReadLineCore_WithMockTerminal_ArrowKeys(t *testing.T) {
	echo := func(b []byte) string { return string(b) }
	formatError := func(msg string) string { return msg }

	// 输入 "helo"，左移，插入 'l'，回车
	mockTerminal := io.NewMockTerminal([]byte{'h', 'e', 'l', 'o', 0x1b, '[', 'D', 'l', '\r'})
	result, err := ReadLineCore("> ", nil, echo, formatError, mockTerminal)

	assert.NoError(t, err)
	assert.Equal(t, "hello", result)
}

func TestReadLineCore_WithMockTerminal_Backspace(t *testing.T) {
	echo := func(b []byte) string { return string(b) }
	formatError := func(msg string) string { return msg }

	// 输入 "hello"，退格两次，回车
	mockTerminal := io.NewMockTerminal([]byte{'h', 'e', 'l', 'l', 'o', 127, 127, '\r'})
	result, err := ReadLineCore("> ", nil, echo, formatError, mockTerminal)

	assert.NoError(t, err)
	assert.Equal(t, "hel", result)
}

func TestReadLineCore_WithMockTerminal_ValidatorError(t *testing.T) {
	echo := func(b []byte) string { return string(b) }
	formatError := func(msg string) string { return "ERR: " + msg }

	validator := func(s string) error {
		if len(s) < 5 {
			return errors.New("太短")
		}
		return nil
	}

	// 输入 "hi" 回车（验证失败，继续等待），然后退格两次清除，输入 "hello" 回车（验证通过）
	mockTerminal := io.NewMockTerminal([]byte{'h', 'i', '\r', 127, 127, 'h', 'e', 'l', 'l', 'o', '\r'})
	result, err := ReadLineCore("> ", validator, echo, formatError, mockTerminal)

	assert.NoError(t, err)
	assert.Equal(t, "hello", result)
}

func TestReadLineCore_WithMockTerminal_CtrlC(t *testing.T) {
	echo := func(b []byte) string { return string(b) }
	formatError := func(msg string) string { return msg }

	mockTerminal := io.NewMockTerminal([]byte{3}) // Ctrl+C
	result, err := ReadLineCore("> ", nil, echo, formatError, mockTerminal)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "用户取消输入")
	assert.Equal(t, "", result)
}

// ==================== readInputFallback 测试 ====================

// Test_readInputFallback_Basic 验证回退输入逻辑会读取并去除首尾空白
func Test_readInputFallback_Basic(t *testing.T) {
	mockTerminal := io.NewMockTerminalWithLines([]string{"value"})
	result, err := readInputFallback("", mockTerminal)
	assert.NoError(t, err)
	assert.Equal(t, "value", result)
}

// Test_readInputFallback_WithWhitespace 验证会去除首尾空白
func Test_readInputFallback_WithWhitespace(t *testing.T) {
	mockTerminal := io.NewMockTerminalWithLines([]string{"  value  "})
	result, err := readInputFallback("", mockTerminal)
	assert.NoError(t, err)
	assert.Equal(t, "value", result) // 应该去除首尾空白
}

// Test_readInputFallback_EmptyInput 验证空输入时返回错误
func Test_readInputFallback_EmptyInput(t *testing.T) {
	mockTerminal := io.NewMockTerminalWithLines([]string{})
	result, err := readInputFallback("", mockTerminal)
	// ReadLine 在空输入时会返回 EOF 错误
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "读取输入失败")
	assert.Equal(t, "", result)
}

// Test_readInputFallback_OnlyWhitespace 验证只有空白时也会去除空白
func Test_readInputFallback_OnlyWhitespace(t *testing.T) {
	mockTerminal := io.NewMockTerminalWithLines([]string{"   "})
	result, err := readInputFallback("", mockTerminal)
	assert.NoError(t, err)
	assert.Equal(t, "", result) // 去除空白后为空
}

// ==================== 辅助函数测试 ====================

// Test_moveCursorToPosition_NoPanic 仅验证不会发生 panic（依赖终端输出，不做内容断言）
func Test_moveCursorToPosition_NoPanic(t *testing.T) {
	// 不同位置调用一次，确保不会 panic
	moveCursorToPosition("prompt: ", []byte("hello"), -1)
	moveCursorToPosition("prompt: ", []byte("hello"), 0)
	moveCursorToPosition("prompt: ", []byte("hello"), 10)
}

// Test_moveCursorToPosition_EdgeCases 验证边界情况
func Test_moveCursorToPosition_EdgeCases(t *testing.T) {
	// 测试空值
	moveCursorToPosition("prompt: ", []byte{}, 0)

	// 测试位置等于长度
	moveCursorToPosition("prompt: ", []byte("hello"), 5)

	// 测试位置大于长度
	moveCursorToPosition("prompt: ", []byte("hello"), 10)
}

// Test_clearErrorLine_MultipleCalls 验证多次调用不会出错
func Test_clearErrorLine_MultipleCalls(t *testing.T) {
	mockTerminal := io.NewMockTerminal([]byte{})

	clearErrorLine(mockTerminal)
	clearErrorLine(mockTerminal)
	clearErrorLine(mockTerminal)
}

// Test_showError_EdgeCases 验证错误显示的边界情况
func Test_showError_EdgeCases(t *testing.T) {
	mockTerminal := io.NewMockTerminal([]byte{})

	formatError := func(msg string) string { return "ERR: " + msg }

	// 空值
	showError("error", "prompt: ", []byte{}, 0, formatError, mockTerminal)

	// 光标位置在末尾
	showError("error", "prompt: ", []byte("hello"), 5, formatError, mockTerminal)

	// 光标位置超出
	showError("error", "prompt: ", []byte("hello"), 10, formatError, mockTerminal)
}

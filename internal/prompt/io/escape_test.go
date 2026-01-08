//go:build test

package io

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/term"
)

// TestEscapeSequenceParser_ReadKey_ArrowKeys 测试箭头键识别
func TestEscapeSequenceParser_ReadKey_ArrowKeys(t *testing.T) {
	testCases := []struct {
		name     string
		input    []byte
		expected KeyType
	}{
		{"上箭头 [A", []byte{'\x1b', '[', 'A'}, KeyUp},
		{"上箭头 O A", []byte{'\x1b', 'O', 'A'}, KeyUp},
		{"下箭头 [B", []byte{'\x1b', '[', 'B'}, KeyDown},
		{"下箭头 O B", []byte{'\x1b', 'O', 'B'}, KeyDown},
		{"右箭头 [C", []byte{'\x1b', '[', 'C'}, KeyRight},
		{"右箭头 O C", []byte{'\x1b', 'O', 'C'}, KeyRight},
		{"左箭头 [D", []byte{'\x1b', '[', 'D'}, KeyLeft},
		{"左箭头 O D", []byte{'\x1b', 'O', 'D'}, KeyLeft},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockTerminal := NewMockTerminal(tc.input)
			parser := NewEscapeSequenceParser(mockTerminal)

			keyType, char, err := parser.ReadKey()

			assert.NoError(t, err)
			assert.Equal(t, tc.expected, keyType)
			assert.Equal(t, byte(0), char)
		})
	}
}

// TestEscapeSequenceParser_ReadKey_Enter 测试回车键识别
func TestEscapeSequenceParser_ReadKey_Enter(t *testing.T) {
	testCases := []struct {
		name     string
		input    []byte
		expected KeyType
	}{
		{"回车 \\r", []byte{'\r'}, KeyEnter},
		{"回车 \\n", []byte{'\n'}, KeyEnter},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockTerminal := NewMockTerminal(tc.input)
			parser := NewEscapeSequenceParser(mockTerminal)

			keyType, char, err := parser.ReadKey()

			assert.NoError(t, err)
			assert.Equal(t, tc.expected, keyType)
			assert.Equal(t, byte(0), char)
		})
	}
}

// TestEscapeSequenceParser_ReadKey_Space 测试空格键识别
func TestEscapeSequenceParser_ReadKey_Space(t *testing.T) {
	mockTerminal := NewMockTerminal([]byte{' '})
	parser := NewEscapeSequenceParser(mockTerminal)

	keyType, char, err := parser.ReadKey()

	assert.NoError(t, err)
	assert.Equal(t, KeySpace, keyType)
	assert.Equal(t, byte(0), char)
}

// TestEscapeSequenceParser_ReadKey_CtrlC 测试 Ctrl+C 识别
func TestEscapeSequenceParser_ReadKey_CtrlC(t *testing.T) {
	mockTerminal := NewMockTerminal([]byte{3})
	parser := NewEscapeSequenceParser(mockTerminal)

	keyType, char, err := parser.ReadKey()

	assert.NoError(t, err)
	assert.Equal(t, KeyCtrlC, keyType)
	assert.Equal(t, byte(0), char)
}

// TestEscapeSequenceParser_ReadKey_Char 测试普通字符识别
func TestEscapeSequenceParser_ReadKey_Char(t *testing.T) {
	testCases := []struct {
		name     string
		input    []byte
		expected byte
	}{
		{"字母 a", []byte{'a'}, 'a'},
		{"字母 A", []byte{'A'}, 'A'},
		{"数字 1", []byte{'1'}, '1'},
		{"符号 @", []byte{'@'}, '@'},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockTerminal := NewMockTerminal(tc.input)
			parser := NewEscapeSequenceParser(mockTerminal)

			keyType, char, err := parser.ReadKey()

			assert.NoError(t, err)
			assert.Equal(t, KeyChar, keyType)
			assert.Equal(t, tc.expected, char)
		})
	}
}

// TestEscapeSequenceParser_ReadKey_InvalidEscape 测试无效转义序列
func TestEscapeSequenceParser_ReadKey_InvalidEscape(t *testing.T) {
	// 无效的转义序列应该返回 KeyUnknown
	mockTerminal := NewMockTerminal([]byte{'\x1b', 'X', 'Y'})
	parser := NewEscapeSequenceParser(mockTerminal)

	keyType, char, err := parser.ReadKey()

	assert.NoError(t, err)
	assert.Equal(t, KeyUnknown, keyType)
	assert.Equal(t, byte(0), char)
}

// TestEscapeSequenceParser_ReadKey_IncompleteEscape 测试不完整的转义序列
func TestEscapeSequenceParser_ReadKey_IncompleteEscape(t *testing.T) {
	// 不完整的转义序列（只有 ESC，没有后续字符）
	mockTerminal := NewMockTerminal([]byte{'\x1b'})
	parser := NewEscapeSequenceParser(mockTerminal)

	// 模拟 ReadByte 返回 EOF
	mockTerminal.SetInputBytes([]byte{'\x1b'})
	// 第二次 ReadByte 会返回 EOF
	keyType, char, err := parser.ReadKey()

	// 当读取失败时，应该返回 KeyUnknown
	assert.NoError(t, err) // parseEscapeSequence 内部处理了错误，返回 KeyUnknown
	assert.Equal(t, KeyUnknown, keyType)
	assert.Equal(t, byte(0), char)
}

// TestEscapeSequenceParser_ReadKey_ReadError 测试读取错误
func TestEscapeSequenceParser_ReadKey_ReadError(t *testing.T) {
	// 创建一个会在 ReadByte 时返回错误的终端
	errorTerminal := &errorTerminal{MockTerminal: NewMockTerminal([]byte{})}
	parser := NewEscapeSequenceParser(errorTerminal)

	keyType, char, err := parser.ReadKey()

	assert.Error(t, err)
	assert.Equal(t, KeyUnknown, keyType)
	assert.Equal(t, byte(0), char)
	assert.Contains(t, err.Error(), "读取输入失败")
}

// errorTerminal 用于测试读取错误的终端
type errorTerminal struct {
	MockTerminal *MockTerminal
}

func (e *errorTerminal) ReadByte() (byte, error) {
	return 0, errors.New("读取错误")
}

func (e *errorTerminal) ReadLine() (string, error) {
	return "", errors.New("读取错误")
}

func (e *errorTerminal) Print(s string) {
	if e.MockTerminal != nil {
		e.MockTerminal.Print(s)
	}
}

func (e *errorTerminal) Println(s string) {
	if e.MockTerminal != nil {
		e.MockTerminal.Println(s)
	}
}

func (e *errorTerminal) Write(data []byte) (int, error) {
	if e.MockTerminal != nil {
		return e.MockTerminal.Write(data)
	}
	return 0, nil
}

func (e *errorTerminal) WriteString(s string) (int, error) {
	if e.MockTerminal != nil {
		return e.MockTerminal.WriteString(s)
	}
	return 0, nil
}

func (e *errorTerminal) MakeRaw() (*term.State, error) {
	if e.MockTerminal != nil {
		return e.MockTerminal.MakeRaw()
	}
	return nil, nil
}

func (e *errorTerminal) Restore(state *term.State) error {
	if e.MockTerminal != nil {
		return e.MockTerminal.Restore(state)
	}
	return nil
}

func (e *errorTerminal) GetFd() int {
	if e.MockTerminal != nil {
		return e.MockTerminal.GetFd()
	}
	return 0
}

func (e *errorTerminal) HideCursor() {
	if e.MockTerminal != nil {
		e.MockTerminal.HideCursor()
	}
}

func (e *errorTerminal) ShowCursor() {
	if e.MockTerminal != nil {
		e.MockTerminal.ShowCursor()
	}
}

func (e *errorTerminal) ClearLine() {
	if e.MockTerminal != nil {
		e.MockTerminal.ClearLine()
	}
}

func (e *errorTerminal) MoveToStart() {
	if e.MockTerminal != nil {
		e.MockTerminal.MoveToStart()
	}
}

func (e *errorTerminal) SaveCursor() {
	if e.MockTerminal != nil {
		e.MockTerminal.SaveCursor()
	}
}

func (e *errorTerminal) RestoreCursor() {
	if e.MockTerminal != nil {
		e.MockTerminal.RestoreCursor()
	}
}

func (e *errorTerminal) ClearToEnd() {
	if e.MockTerminal != nil {
		e.MockTerminal.ClearToEnd()
	}
}

func (e *errorTerminal) ResetFormat() {
	if e.MockTerminal != nil {
		e.MockTerminal.ResetFormat()
	}
}

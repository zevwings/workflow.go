//go:build test

package io

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// ==================== MockTerminal 测试 ====================

func TestMockTerminal_ReadByte(t *testing.T) {
	mock := NewMockTerminal([]byte{'a', 'b', 'c'})

	b1, err := mock.ReadByte()
	assert.NoError(t, err)
	assert.Equal(t, byte('a'), b1)

	b2, err := mock.ReadByte()
	assert.NoError(t, err)
	assert.Equal(t, byte('b'), b2)

	b3, err := mock.ReadByte()
	assert.NoError(t, err)
	assert.Equal(t, byte('c'), b3)

	// 读取完毕应返回 EOF
	_, err = mock.ReadByte()
	assert.Error(t, err)
}

func TestMockTerminal_ReadLine(t *testing.T) {
	mock := NewMockTerminalWithLines([]string{"line1", "line2"})

	line1, err := mock.ReadLine()
	assert.NoError(t, err)
	assert.Equal(t, "line1", line1)

	line2, err := mock.ReadLine()
	assert.NoError(t, err)
	assert.Equal(t, "line2", line2)

	// 读取完毕应返回 EOF
	_, err = mock.ReadLine()
	assert.Error(t, err)
}

func TestMockTerminal_Print(t *testing.T) {
	mock := NewMockTerminal([]byte{})

	mock.Print("hello")
	mock.Print(" world")

	assert.Equal(t, "hello world", mock.GetOutput())
}

func TestMockTerminal_Println(t *testing.T) {
	mock := NewMockTerminal([]byte{})

	mock.Println("hello")
	mock.Println("world")

	assert.Equal(t, "hello\nworld\n", mock.GetOutput())
}

func TestMockTerminal_ANSICommands(t *testing.T) {
	mock := NewMockTerminal([]byte{})

	mock.HideCursor()
	assert.True(t, mock.HideCursorCalled)
	assert.Contains(t, mock.GetOutput(), "\033[?25l")

	mock.ShowCursor()
	assert.True(t, mock.ShowCursorCalled)
	assert.Contains(t, mock.GetOutput(), "\033[?25h")

	mock.ClearLine()
	assert.Equal(t, 1, mock.ClearLineCalled)
	assert.Contains(t, mock.GetOutput(), "\033[K")

	mock.MoveToStart()
	assert.Equal(t, 1, mock.MoveToStartCalled)
	assert.Contains(t, mock.GetOutput(), "\r")

	mock.SaveCursor()
	assert.True(t, mock.SaveCursorCalled)
	assert.Contains(t, mock.GetOutput(), "\033[s")

	mock.RestoreCursor()
	assert.True(t, mock.RestoreCursorCalled)
	assert.Contains(t, mock.GetOutput(), "\033[u")

	mock.ClearToEnd()
	assert.True(t, mock.ClearToEndCalled)
	assert.Contains(t, mock.GetOutput(), "\033[J")

	mock.ResetFormat()
	assert.True(t, mock.ResetFormatCalled)
	assert.Contains(t, mock.GetOutput(), "\033[0m")
}

func TestMockTerminal_MakeRawAndRestore(t *testing.T) {
	mock := NewMockTerminal([]byte{})

	state, err := mock.MakeRaw()
	assert.NoError(t, err)
	assert.True(t, mock.RawModeEnabled)

	err = mock.Restore(state)
	assert.NoError(t, err)
	assert.False(t, mock.RawModeEnabled)
	assert.Equal(t, state, mock.SavedState)
}

func TestMockTerminal_Reset(t *testing.T) {
	mock := NewMockTerminal([]byte{'a'})
	mock.Print("test")
	mock.HideCursor()

	mock.Reset()

	assert.Equal(t, 0, mock.ReadIndex)
	assert.Equal(t, "", mock.GetOutput())
	assert.False(t, mock.HideCursorCalled)
}

func TestMockTerminal_SetInputBytes(t *testing.T) {
	mock := NewMockTerminal([]byte{'a'})
	mock.ReadByte() // 消耗第一个字节

	mock.SetInputBytes([]byte{'b', 'c'})
	b, err := mock.ReadByte()
	assert.NoError(t, err)
	assert.Equal(t, byte('b'), b)
}

func TestMockTerminal_SetInputLines(t *testing.T) {
	mock := NewMockTerminalWithLines([]string{"line1"})
	mock.ReadLine() // 消耗第一行

	mock.SetInputLines([]string{"line2", "line3"})
	line, err := mock.ReadLine()
	assert.NoError(t, err)
	assert.Equal(t, "line2", line)
}

func TestMockTerminal_VerifyOutputContains(t *testing.T) {
	mock := NewMockTerminal([]byte{})
	mock.Print("hello world")

	assert.True(t, mock.VerifyOutputContains("hello"))
	assert.True(t, mock.VerifyOutputContains("world"))
	assert.False(t, mock.VerifyOutputContains("foo"))
}


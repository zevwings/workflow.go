//go:build test

package io

import (
	"bytes"
	"io"
	"strings"

	"golang.org/x/term"
)

// MockTerminal Mock 终端实现，用于测试
// 可以预设输入序列，捕获所有输出，验证调用
type MockTerminal struct {
	// 输入相关
	InputBytes   []byte
	ReadIndex    int
	InputLines   []string
	ReadLineIndex int

	// 输出相关
	OutputBuffer *bytes.Buffer

	// 终端状态
	RawModeEnabled bool
	SavedState     *term.State

	// 调用记录（用于验证）
	HideCursorCalled   bool
	ShowCursorCalled   bool
	ClearLineCalled    int
	MoveToStartCalled  int
	SaveCursorCalled   bool
	RestoreCursorCalled bool
	ClearToEndCalled   bool
	ResetFormatCalled  bool
}

// NewMockTerminal 创建 Mock 终端实例
// inputBytes: 预设的输入字节序列（用于 ReadByte）
func NewMockTerminal(inputBytes []byte) *MockTerminal {
	return &MockTerminal{
		InputBytes:    inputBytes,
		ReadIndex:     0,
		OutputBuffer:  &bytes.Buffer{},
		InputLines:     []string{},
		ReadLineIndex: 0,
	}
}

// NewMockTerminalWithLines 创建 Mock 终端实例（支持 ReadLine）
// inputLines: 预设的输入行序列（用于 ReadLine）
func NewMockTerminalWithLines(inputLines []string) *MockTerminal {
	return &MockTerminal{
		InputBytes:    []byte{},
		ReadIndex:     0,
		OutputBuffer:  &bytes.Buffer{},
		InputLines:    inputLines,
		ReadLineIndex: 0,
	}
}

// ==================== 输入操作 ====================

// ReadByte 读取单个字节
func (m *MockTerminal) ReadByte() (byte, error) {
	if m.ReadIndex >= len(m.InputBytes) {
		return 0, io.EOF
	}
	b := m.InputBytes[m.ReadIndex]
	m.ReadIndex++
	return b, nil
}

// ReadLine 读取一行输入
func (m *MockTerminal) ReadLine() (string, error) {
	if m.ReadLineIndex >= len(m.InputLines) {
		return "", io.EOF
	}
	line := m.InputLines[m.ReadLineIndex]
	m.ReadLineIndex++
	return line, nil
}

// ==================== 输出操作 ====================

// Print 输出字符串（不换行）
func (m *MockTerminal) Print(s string) {
	m.OutputBuffer.WriteString(s)
}

// Println 输出字符串并换行
func (m *MockTerminal) Println(s string) {
	m.OutputBuffer.WriteString(s)
	m.OutputBuffer.WriteString("\n")
}

// Write 写入字节数据
func (m *MockTerminal) Write(data []byte) (int, error) {
	return m.OutputBuffer.Write(data)
}

// WriteString 写入字符串
func (m *MockTerminal) WriteString(s string) (int, error) {
	return m.OutputBuffer.WriteString(s)
}

// ==================== 终端控制 ====================

// MakeRaw 设置终端为原始模式
func (m *MockTerminal) MakeRaw() (*term.State, error) {
	m.RawModeEnabled = true
	// 返回 nil state，模拟成功
	return nil, nil
}

// Restore 恢复终端状态
func (m *MockTerminal) Restore(state *term.State) error {
	m.RawModeEnabled = false
	m.SavedState = state
	return nil
}

// GetFd 获取文件描述符（Mock 返回固定值）
func (m *MockTerminal) GetFd() int {
	return 0
}

// ==================== ANSI 控制码 ====================

// HideCursor 隐藏光标
func (m *MockTerminal) HideCursor() {
	m.HideCursorCalled = true
	m.OutputBuffer.WriteString("\033[?25l")
}

// ShowCursor 显示光标
func (m *MockTerminal) ShowCursor() {
	m.ShowCursorCalled = true
	m.OutputBuffer.WriteString("\033[?25h")
}

// ClearLine 清除当前行
func (m *MockTerminal) ClearLine() {
	m.ClearLineCalled++
	m.OutputBuffer.WriteString("\033[K")
}

// MoveToStart 移动到行首
func (m *MockTerminal) MoveToStart() {
	m.MoveToStartCalled++
	m.OutputBuffer.WriteString("\r")
}

// SaveCursor 保存光标位置
func (m *MockTerminal) SaveCursor() {
	m.SaveCursorCalled = true
	m.OutputBuffer.WriteString("\033[s")
}

// RestoreCursor 恢复光标位置
func (m *MockTerminal) RestoreCursor() {
	m.RestoreCursorCalled = true
	m.OutputBuffer.WriteString("\033[u")
}

// ClearToEnd 清除从光标到屏幕底部的内容
func (m *MockTerminal) ClearToEnd() {
	m.ClearToEndCalled = true
	m.OutputBuffer.WriteString("\033[J")
}

// ResetFormat 重置所有 ANSI 格式
func (m *MockTerminal) ResetFormat() {
	m.ResetFormatCalled = true
	m.OutputBuffer.WriteString("\033[0m")
}

// ==================== 测试辅助方法 ====================

// GetOutput 获取所有输出内容
func (m *MockTerminal) GetOutput() string {
	return m.OutputBuffer.String()
}

// Reset 重置 Mock 状态（用于测试多个场景）
func (m *MockTerminal) Reset() {
	m.ReadIndex = 0
	m.ReadLineIndex = 0
	m.OutputBuffer.Reset()
	m.RawModeEnabled = false
	m.SavedState = nil
	m.HideCursorCalled = false
	m.ShowCursorCalled = false
	m.ClearLineCalled = 0
	m.MoveToStartCalled = 0
	m.SaveCursorCalled = false
	m.RestoreCursorCalled = false
	m.ClearToEndCalled = false
	m.ResetFormatCalled = false
}

// SetInputBytes 设置输入字节序列
func (m *MockTerminal) SetInputBytes(bytes []byte) {
	m.InputBytes = bytes
	m.ReadIndex = 0
}

// SetInputLines 设置输入行序列
func (m *MockTerminal) SetInputLines(lines []string) {
	m.InputLines = lines
	m.ReadLineIndex = 0
}

// VerifyOutputContains 验证输出包含指定字符串
func (m *MockTerminal) VerifyOutputContains(s string) bool {
	return strings.Contains(m.GetOutput(), s)
}


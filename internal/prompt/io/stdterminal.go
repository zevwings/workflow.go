package io

import (
	"fmt"
	"io"
	"os"

	"golang.org/x/term"
)

// StdTerminal 标准终端实现
// 使用系统的标准输入输出和终端控制
type StdTerminal struct {
	stdin  *os.File
	stdout *os.File
}

// NewStdTerminal 创建标准终端实例
func NewStdTerminal() *StdTerminal {
	return &StdTerminal{
		stdin:  os.Stdin,
		stdout: os.Stdout,
	}
}

// ==================== 输入操作 ====================

// ReadByte 读取单个字节
func (t *StdTerminal) ReadByte() (byte, error) {
	buf := make([]byte, 1)
	n, err := t.stdin.Read(buf)
	if err != nil {
		return 0, err
	}
	if n == 0 {
		return 0, io.EOF
	}
	return buf[0], nil
}

// ReadLine 读取一行输入
func (t *StdTerminal) ReadLine() (string, error) {
	var input string
	_, err := fmt.Fscanln(t.stdin, &input)
	if err != nil {
		return "", err
	}
	return input, nil
}

// ==================== 输出操作 ====================

// Print 输出字符串（不换行）
func (t *StdTerminal) Print(s string) {
	fmt.Fprint(t.stdout, s)
}

// Println 输出字符串并换行
func (t *StdTerminal) Println(s string) {
	fmt.Fprintln(t.stdout, s)
}

// Write 写入字节数据
func (t *StdTerminal) Write(data []byte) (int, error) {
	return t.stdout.Write(data)
}

// WriteString 写入字符串
func (t *StdTerminal) WriteString(s string) (int, error) {
	return t.stdout.WriteString(s)
}

// ==================== 终端控制 ====================

// MakeRaw 设置终端为原始模式
func (t *StdTerminal) MakeRaw() (*term.State, error) {
	fd := int(t.stdin.Fd())
	return term.MakeRaw(fd)
}

// Restore 恢复终端状态
func (t *StdTerminal) Restore(state *term.State) error {
	fd := int(t.stdin.Fd())
	return term.Restore(fd, state)
}

// GetFd 获取文件描述符
func (t *StdTerminal) GetFd() int {
	return int(t.stdin.Fd())
}

// ==================== ANSI 控制码 ====================

// HideCursor 隐藏光标
func (t *StdTerminal) HideCursor() {
	fmt.Fprint(t.stdout, "\033[?25l")
}

// ShowCursor 显示光标
func (t *StdTerminal) ShowCursor() {
	fmt.Fprint(t.stdout, "\033[?25h")
}

// ClearLine 清除当前行
func (t *StdTerminal) ClearLine() {
	fmt.Fprint(t.stdout, "\033[K")
}

// MoveToStart 移动到行首
func (t *StdTerminal) MoveToStart() {
	fmt.Fprint(t.stdout, "\r")
}

// SaveCursor 保存光标位置
func (t *StdTerminal) SaveCursor() {
	fmt.Fprint(t.stdout, "\033[s")
}

// RestoreCursor 恢复光标位置
func (t *StdTerminal) RestoreCursor() {
	fmt.Fprint(t.stdout, "\033[u")
}

// ClearToEnd 清除从光标到屏幕底部的内容
func (t *StdTerminal) ClearToEnd() {
	fmt.Fprint(t.stdout, "\033[J")
}

// ResetFormat 重置所有 ANSI 格式
func (t *StdTerminal) ResetFormat() {
	fmt.Fprint(t.stdout, "\033[0m")
}


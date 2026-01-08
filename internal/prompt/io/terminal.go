package io

import "golang.org/x/term"

// TerminalIO 终端 I/O 接口
// 抽象了所有终端输入输出和终端控制操作，便于测试和扩展
type TerminalIO interface {
	// ==================== 输入操作 ====================

	// ReadByte 读取单个字节（用于交互式输入）
	ReadByte() (byte, error)

	// ReadLine 读取一行输入（用于 fallback 模式）
	ReadLine() (string, error)

	// ==================== 输出操作 ====================

	// Print 输出字符串（不换行）
	Print(s string)

	// Println 输出字符串并换行
	Println(s string)

	// Write 写入字节数据
	Write(data []byte) (int, error)

	// WriteString 写入字符串
	WriteString(s string) (int, error)

	// ==================== 终端控制 ====================

	// MakeRaw 设置终端为原始模式
	// 返回终端状态，用于后续恢复
	MakeRaw() (*term.State, error)

	// Restore 恢复终端状态
	Restore(state *term.State) error

	// GetFd 获取文件描述符
	GetFd() int

	// ==================== ANSI 控制码 ====================

	// HideCursor 隐藏光标
	HideCursor()

	// ShowCursor 显示光标
	ShowCursor()

	// ClearLine 清除当前行
	ClearLine()

	// MoveToStart 移动到行首
	MoveToStart()

	// SaveCursor 保存光标位置
	SaveCursor()

	// RestoreCursor 恢复光标位置
	RestoreCursor()

	// ClearToEnd 清除从光标到屏幕底部的内容
	ClearToEnd()

	// ResetFormat 重置所有 ANSI 格式
	ResetFormat()
}


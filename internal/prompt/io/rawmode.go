package io

import (
	"fmt"
)

// RawModeManager 管理终端原始模式
// 封装 MakeRaw() 和 Restore() 的调用，统一处理错误和 Fallback
type RawModeManager struct {
	terminal TerminalIO
}

// NewRawModeManager 创建原始模式管理器
func NewRawModeManager(terminal TerminalIO) *RawModeManager {
	return &RawModeManager{
		terminal: terminal,
	}
}

// WithRawMode 在原始模式下执行函数
// 自动管理终端状态的设置和恢复，包括光标的显示/隐藏
//
// 参数:
//   - fn: 在原始模式下执行的函数
//
// 返回:
//   - error: 执行过程中的错误
func (m *RawModeManager) WithRawMode(fn func() error) error {
	// 设置终端原始模式
	oldState, err := m.terminal.MakeRaw()
	if err != nil {
		return fmt.Errorf("设置终端原始模式失败: %w", err)
	}

	// 确保恢复终端状态
	defer func() {
		m.terminal.ShowCursor()
		m.terminal.Restore(oldState)
	}()

	// 隐藏光标
	m.terminal.HideCursor()

	// 执行函数
	return fn()
}

// WithRawModeAndFallback 在原始模式下执行函数，失败时执行 fallback
// 如果无法设置原始模式，会执行 fallback 函数作为回退方案
//
// 参数:
//   - fn: 在原始模式下执行的函数
//   - fallback: 当无法设置原始模式时执行的回退函数
//
// 返回:
//   - error: 执行过程中的错误
func (m *RawModeManager) WithRawModeAndFallback(fn func() error, fallback func() error) error {
	// 尝试设置终端原始模式
	oldState, err := m.terminal.MakeRaw()
	if err != nil {
		// 如果无法设置原始模式，执行 fallback
		return fallback()
	}

	// 确保恢复终端状态
	defer func() {
		m.terminal.ShowCursor()
		m.terminal.Restore(oldState)
	}()

	// 隐藏光标
	m.terminal.HideCursor()

	// 执行函数
	return fn()
}

// GetTerminal 获取终端接口（用于需要直接访问终端的场景）
func (m *RawModeManager) GetTerminal() TerminalIO {
	return m.terminal
}

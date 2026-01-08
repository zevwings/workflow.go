package io

import (
	"fmt"
)

// KeyType 按键类型
type KeyType string

const (
	KeyUp      KeyType = "up"      // 上箭头
	KeyDown    KeyType = "down"    // 下箭头
	KeyLeft    KeyType = "left"    // 左箭头
	KeyRight   KeyType = "right"   // 右箭头
	KeyEnter   KeyType = "enter"   // 回车键
	KeySpace   KeyType = "space"   // 空格键
	KeyCtrlC   KeyType = "ctrl+c"  // Ctrl+C
	KeyChar    KeyType = "char"    // 普通字符
	KeyUnknown KeyType = "unknown" // 未知按键
)

// EscapeSequenceParser 解析转义序列
// 用于识别 ANSI 转义序列，特别是箭头键和其他特殊键
type EscapeSequenceParser struct {
	terminal TerminalIO
}

// NewEscapeSequenceParser 创建转义序列解析器
func NewEscapeSequenceParser(terminal TerminalIO) *EscapeSequenceParser {
	return &EscapeSequenceParser{
		terminal: terminal,
	}
}

// ReadKey 读取并解析按键
// 返回按键类型和字符（如果是普通字符）
//
// 返回:
//   - keyType: 按键类型（"up", "down", "enter", "ctrl+c", "char" 等）
//   - char: 如果是普通字符，返回该字符；否则为 0
//   - err: 读取错误
func (p *EscapeSequenceParser) ReadKey() (keyType KeyType, char byte, err error) {
	char, err = p.terminal.ReadByte()
	if err != nil {
		return KeyUnknown, 0, fmt.Errorf("读取输入失败: %w", err)
	}

	// 处理转义序列（箭头键）
	if char == '\x1b' {
		return p.parseEscapeSequence()
	}

	// 处理回车键
	if char == '\r' || char == '\n' {
		return KeyEnter, 0, nil
	}

	// 处理空格键
	if char == ' ' {
		return KeySpace, 0, nil
	}

	// 处理 Ctrl+C
	if char == 3 {
		return KeyCtrlC, 0, nil
	}

	// 普通字符
	return KeyChar, char, nil
}

// parseEscapeSequence 解析转义序列
// 处理 ANSI 转义序列，识别箭头键
func (p *EscapeSequenceParser) parseEscapeSequence() (KeyType, byte, error) {
	// 读取转义序列的后续字符
	char2, err := p.terminal.ReadByte()
	if err != nil {
		// 如果读取失败，忽略这个转义序列
		return KeyUnknown, 0, nil
	}

	// 检查第二个字符
	if char2 == '[' || char2 == 'O' {
		char3, err := p.terminal.ReadByte()
		if err != nil {
			// 如果读取失败，忽略这个转义序列
			return KeyUnknown, 0, nil
		}

		// 判断箭头方向
		if char3 == 'A' || (char2 == 'O' && char3 == 'A') {
			// 上箭头
			return KeyUp, 0, nil
		}
		if char3 == 'B' || (char2 == 'O' && char3 == 'B') {
			// 下箭头
			return KeyDown, 0, nil
		}
		if char3 == 'C' || (char2 == 'O' && char3 == 'C') {
			// 右箭头
			return KeyRight, 0, nil
		}
		if char3 == 'D' || (char2 == 'O' && char3 == 'D') {
			// 左箭头
			return KeyLeft, 0, nil
		}
	}

	// 不是有效的箭头键序列，忽略
	return KeyUnknown, 0, nil
}


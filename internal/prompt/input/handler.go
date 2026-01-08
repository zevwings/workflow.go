package input

import (
	"strings"

	"github.com/mattn/go-runewidth"
)

// InputHandler 处理输入逻辑（纯业务逻辑，无 I/O 操作）
type InputHandler struct {
	validator Validator
}

// NewInputHandler 创建输入处理器
func NewInputHandler(validator Validator) *InputHandler {
	return &InputHandler{
		validator: validator,
	}
}

// ProcessEscapeSequence 处理转义序列（方向键等）
// 返回：方向（"left"/"right"/"none"），是否需要更新光标位置
func (h *InputHandler) ProcessEscapeSequence(char2, char3 byte) (direction string, shouldUpdate bool) {
	if char2 == '[' {
		switch char3 {
		case 'C': // 右箭头
			return "right", true
		case 'D': // 左箭头
			return "left", true
		case 'A', 'B': // 上下箭头 - 忽略
			return "none", false
		default:
			return "none", false
		}
	}
	return "none", false
}

// ProcessArrowKey 处理箭头键移动光标
// 返回：新的光标位置
func (h *InputHandler) ProcessArrowKey(currentPos int, valueLen int, direction string, placeholderDisplayed bool) (newPos int, shouldClearPlaceholder bool) {
	if placeholderDisplayed {
		if direction == "right" {
			return 0, true // 清除 placeholder，光标位置设为 0
		}
		return 0, false // 左箭头在 placeholder 时无效
	}

	if direction == "right" {
		if currentPos < valueLen {
			return currentPos + 1, false
		}
		return currentPos, false
	}

	if direction == "left" {
		if currentPos > 0 {
			return currentPos - 1, false
		}
		return currentPos, false
	}

	return currentPos, false
}

// ProcessBackspace 处理退格键
// 返回：新的值和光标位置
func (h *InputHandler) ProcessBackspace(value []byte, cursorPos int) (newValue []byte, newPos int) {
	if cursorPos > 0 {
		newValue = append(value[:cursorPos-1], value[cursorPos:]...)
		return newValue, cursorPos - 1
	}
	return value, cursorPos
}

// ProcessChar 处理普通字符输入
// 返回：新的值和光标位置
func (h *InputHandler) ProcessChar(value []byte, cursorPos int, char byte) (newValue []byte, newPos int) {
	if cursorPos < len(value) {
		newValue = append(value[:cursorPos], append([]byte{char}, value[cursorPos:]...)...)
	} else {
		newValue = append(value, char)
	}
	return newValue, cursorPos + 1
}

// ValidateInput 验证输入
func (h *InputHandler) ValidateInput(value string) error {
	if h.validator == nil {
		return nil
	}
	return h.validator(value)
}

// CalculateCursorBackspaces 计算需要回退的光标位置（用于 moveCursorToPosition）
func (h *InputHandler) CalculateCursorBackspaces(value []byte, pos int) int {
	if pos < 0 {
		pos = 0
	}
	if pos > len(value) {
		pos = len(value)
	}

	if pos < len(value) {
		tail := value[pos:]
		return runewidth.StringWidth(string(tail))
	}
	return 0
}

// CalculatePlaceholderBackspaces 计算 placeholder 需要回退的字符数
func (h *InputHandler) CalculatePlaceholderBackspaces(placeholderText string) int {
	return runewidth.StringWidth(placeholderText)
}

// TrimInput 清理输入（去除首尾空白）
func (h *InputHandler) TrimInput(value string) string {
	return strings.TrimSpace(value)
}

// PlaceholderHandler 处理 placeholder 相关逻辑
type PlaceholderHandler struct {
	placeholder string
	config      Config
}

// NewPlaceholderHandler 创建 placeholder 处理器
func NewPlaceholderHandler(placeholder string, config Config) *PlaceholderHandler {
	return &PlaceholderHandler{
		placeholder: placeholder,
		config:      config,
	}
}

// HasPlaceholder 检查是否有 placeholder
func (p *PlaceholderHandler) HasPlaceholder() bool {
	return p.placeholder != ""
}

// FormatPlaceholder 格式化 placeholder
func (p *PlaceholderHandler) FormatPlaceholder() string {
	return p.config.FormatPlaceholder(p.placeholder)
}

// GetPlaceholderText 获取去除 ANSI 代码的 placeholder 文本
func (p *PlaceholderHandler) GetPlaceholderText() string {
	formatted := p.FormatPlaceholder()
	return StripAnsiCodes(formatted)
}

// GetPlaceholderWidth 获取 placeholder 的显示宽度
func (p *PlaceholderHandler) GetPlaceholderWidth() int {
	text := p.GetPlaceholderText()
	return runewidth.StringWidth(text)
}

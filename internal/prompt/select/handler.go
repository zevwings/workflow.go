package selectpkg

import (
	"github.com/zevwings/workflow/internal/prompt/common"
)

// SelectHandler 处理选择逻辑（纯业务逻辑，无 I/O 操作）
type SelectHandler struct {
	options      []string
	defaultIndex int
	config       Config
	navigator    *common.NavigationHandler
}

// NewSelectHandler 创建选择处理器
func NewSelectHandler(options []string, defaultIndex int, config Config) *SelectHandler {
	return &SelectHandler{
		options:      options,
		defaultIndex: defaultIndex,
		config:       config,
		navigator:    common.NewNavigationHandler(len(options), false),
	}
}

// ValidateAndAdjustDefaultIndex 验证并调整默认索引
func (h *SelectHandler) ValidateAndAdjustDefaultIndex() int {
	return h.navigator.ValidateIndex(h.defaultIndex)
}

// ProcessArrowKey 处理箭头键输入
// 返回新的 currentIndex 和是否需要重新渲染
func (h *SelectHandler) ProcessArrowKey(currentIndex int, direction string) (newIndex int, shouldRender bool) {
	return h.navigator.ProcessArrowKey(currentIndex, direction)
}

// FormatOptionLine 格式化选项行
// 返回格式化后的行文本和是否需要高亮
func (h *SelectHandler) FormatOptionLine(index int, currentIndex int) (line string, isHighlighted bool) {
	if index == currentIndex {
		return "> " + h.options[index], true
	}
	return "  " + h.options[index], false
}

// FormatSelectedOption 格式化选中的选项（用于显示结果）
func (h *SelectHandler) FormatSelectedOption(index int) string {
	return h.config.FormatAnswer(h.options[index])
}

// ParseNumericInput 解析数字输入（用于 fallback 模式）
// 返回选中的索引，如果无效则返回默认索引
func (h *SelectHandler) ParseNumericInput(input int) int {
	if input < 1 || input > len(h.options) {
		return h.defaultIndex
	}
	return input - 1
}

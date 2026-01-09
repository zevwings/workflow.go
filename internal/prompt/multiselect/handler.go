package multiselect

import (
	"fmt"
	"strings"

	"github.com/zevwings/workflow/internal/prompt/common"
)

// MultiSelectHandler 处理多选逻辑（纯业务逻辑，无 I/O 操作）
type MultiSelectHandler struct {
	options         []string
	defaultSelected []int
	config          common.PromptConfig
	navigator       *common.NavigationHandler
}

// NewMultiSelectHandler 创建多选处理器
func NewMultiSelectHandler(options []string, defaultSelected []int, config common.PromptConfig) *MultiSelectHandler {
	return &MultiSelectHandler{
		options:         options,
		defaultSelected: defaultSelected,
		config:          config,
		navigator:       common.NewNavigationHandler(len(options), false),
	}
}

// ValidateAndCleanDefaults 验证并清理默认选中项
// 返回清理后的选中项 map
func (h *MultiSelectHandler) ValidateAndCleanDefaults() map[int]bool {
	selected := make(map[int]bool)
	for _, idx := range h.defaultSelected {
		if idx >= 0 && idx < len(h.options) {
			selected[idx] = true
		}
	}
	return selected
}

// GetInitialCurrentIndex 获取初始光标位置
func (h *MultiSelectHandler) GetInitialCurrentIndex() int {
	if len(h.defaultSelected) > 0 && h.defaultSelected[0] >= 0 && h.defaultSelected[0] < len(h.options) {
		return h.defaultSelected[0]
	}
	return 0
}

// ProcessArrowKey 处理箭头键输入
// 返回新的 currentIndex 和是否需要重新渲染
func (h *MultiSelectHandler) ProcessArrowKey(currentIndex int, direction string) (newIndex int, shouldRender bool) {
	return h.navigator.ProcessArrowKey(currentIndex, direction)
}

// ToggleSelection 切换指定索引的选择状态
func (h *MultiSelectHandler) ToggleSelection(selected map[int]bool, index int) {
	if selected[index] {
		delete(selected, index)
	} else {
		selected[index] = true
	}
}

// FormatOptionLine 格式化选项行
// 返回格式化后的行文本和是否需要高亮
func (h *MultiSelectHandler) FormatOptionLine(index int, currentIndex int, selected map[int]bool) (line string, isHighlighted bool) {
	prefix := "  "
	if index == currentIndex {
		prefix = "> "
	}

	marker := "[ ]"
	if selected[index] {
		marker = "[x]"
	}

	line = fmt.Sprintf("%s%s %s", prefix, marker, h.options[index])
	isHighlighted = (index == currentIndex)
	return
}

// FormatSelectedOptions 格式化选中的选项（用于显示结果）
func (h *MultiSelectHandler) FormatSelectedOptions(selectedIndices []int) string {
	if len(selectedIndices) == 0 {
		return "(未选择)"
	}

	var selectedOptions []string
	for _, idx := range selectedIndices {
		// 验证索引有效性
		if idx >= 0 && idx < len(h.options) {
			selectedOptions = append(selectedOptions, h.options[idx])
		}
	}

	if len(selectedOptions) == 0 {
		return "(未选择)"
	}

	return h.config.FormatAnswer(strings.Join(selectedOptions, ", "))
}

// ParseCommaSeparatedInput 解析逗号分隔的输入（用于 fallback 模式）
func (h *MultiSelectHandler) ParseCommaSeparatedInput(input string) []int {
	input = strings.TrimSpace(input)
	if input == "" {
		// 返回清理后的默认值
		return h.GetDefaultSelectedForFallback()
	}

	parts := strings.Split(input, ",")
	selectedIndices := make(map[int]bool)

	for _, part := range parts {
		part = strings.TrimSpace(part)
		var num int
		_, err := fmt.Sscanf(part, "%d", &num)
		if err != nil {
			continue
		}
		if num >= 1 && num <= len(h.options) {
			selectedIndices[num-1] = true
		}
	}

	return mapToSlice(selectedIndices)
}

// GetDefaultSelectedForFallback 获取 fallback 模式的默认选中项（已清理）
func (h *MultiSelectHandler) GetDefaultSelectedForFallback() []int {
	selected := h.ValidateAndCleanDefaults()
	return mapToSlice(selected)
}

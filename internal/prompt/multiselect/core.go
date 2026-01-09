package multiselect

import (
	"fmt"
	"sort"

	"github.com/zevwings/workflow/internal/prompt/common"
)

// MultiSelectConfig 多选功能配置
type MultiSelectConfig struct {
	common.BasePromptConfig
	// Options 选项列表
	Options []string
	// DefaultSelected 默认选中的索引列表
	DefaultSelected []int
}

// MultiSelect 多选功能（使用配置结构体）
func MultiSelect(cfg MultiSelectConfig) ([]int, error) {
	if len(cfg.Options) == 0 {
		return nil, fmt.Errorf("选项列表不能为空")
	}

	handler := NewMultiSelectHandler(cfg.Options, cfg.DefaultSelected, cfg.Config)

	// 验证并清理默认选中项
	selected := handler.ValidateAndCleanDefaults()

	// 设置交互式选择的通用组件
	setup := common.SetupInteractiveSelect(cfg.BasePromptConfig)
	promptMsg := setup.PromptMsg
	promptMsgWithPrefix := setup.PromptMsgWithPrefix
	rawModeMgr := setup.RawModeMgr
	parser := setup.Parser
	renderer := setup.Renderer

	currentIndex := handler.GetInitialCurrentIndex()

	// 使用原始模式管理器执行交互逻辑
	var result []int

	err := rawModeMgr.WithRawModeAndFallback(
		func() error {
			// 使用通用渲染函数渲染选项列表
			// 通过闭包捕获 selected 和 currentIndex，创建符合 FormatOptionLineFunc 签名的函数
			formatLine := func(index int, currentIdx int) (string, bool) {
				return handler.FormatOptionLine(index, currentIdx, selected)
			}
			// 使用闭包捕获 currentIndex 的引用，支持动态更新
			getCurrentIndex := func() int {
				return currentIndex
			}
			renderMultiSelect := common.RenderOptions(
				cfg.Terminal,
				renderer,
				len(cfg.Options),
				getCurrentIndex,
				formatLine,
				"使用 ↑/↓ 导航，空格键切换选择，回车确认",
				cfg.Config,
			)

			// 使用渲染器渲染提示和初始界面（未输入时使用 "? " 前缀）
			if err := renderer.RenderWithPrompt(promptMsgWithPrefix, renderMultiSelect); err != nil {
				return err
			}

			// 使用通用输入处理函数
			return common.HandleInteractiveInput(
				parser,
				cfg.Terminal,
				&currentIndex,
				func(idx int, dir string) (int, bool) {
					return handler.ProcessArrowKey(idx, dir)
				},
				func() (bool, error) {
					selectedIndices := mapToSlice(selected)
					selectedText := handler.FormatSelectedOptions(selectedIndices)
					if err := common.FormatResultWithTitle(cfg.Terminal, promptMsg, selectedText, nil, false, cfg.Message, cfg.Config.FormatResultTitle, cfg.Config.FormatAnswerPrefix); err != nil {
						return false, err
					}
					result = selectedIndices
					return true, nil
				},
				func() bool {
					handler.ToggleSelection(selected, currentIndex)
					return true
				},
				func() {
					// 使用渲染器的 ReRender 方法确保正确清除之前的内容
					renderer.ReRender(renderMultiSelect)
				},
			)
		},
		func() error {
			// Fallback: 如果无法设置原始模式，使用简单多选
			selectedSlice, err := multiselectFallback(cfg)
			result = selectedSlice
			return err
		},
	)

	if err != nil {
		// 如果 err 是取消错误，返回 nil 和错误
		return nil, err
	}

	return result, nil
}

// mapToSlice 将 map[int]bool 转换为排序后的 []int
func mapToSlice(selected map[int]bool) []int {
	if len(selected) == 0 {
		return []int{} // 返回空切片而不是 nil
	}
	indices := make([]int, 0, len(selected))
	for idx := range selected {
		indices = append(indices, idx)
	}
	// 使用标准库的快速排序算法（O(n log n)）
	sort.Ints(indices)
	return indices
}

// multiselectFallback 回退方案：如果无法设置原始模式，使用简单多选
func multiselectFallback(cfg MultiSelectConfig) ([]int, error) {
	handler := NewMultiSelectHandler(cfg.Options, cfg.DefaultSelected, cfg.Config)
	defaultSelected := handler.ValidateAndCleanDefaults()

	// 使用通用的 fallback 框架
	return common.ExecuteMultiSelectFallback(
		cfg.Terminal,
		cfg.Message,
		cfg.Config,
		cfg.Options,
		common.MultiSelectFallbackOptions{
			FormatOptionLine: func(index int, option string, isSelected bool) string {
				marker := "[ ]"
				if isSelected {
					marker = "[x]"
				}
				return fmt.Sprintf("  %s %d. %s\n", marker, index+1, option)
			},
			GetDefaultSelected: func() map[int]bool {
				return defaultSelected
			},
			ParseInput: func(input string) []int {
				return handler.ParseCommaSeparatedInput(input)
			},
			FormatSelectedOptions: func(selectedIndices []int) string {
				return handler.FormatSelectedOptions(selectedIndices)
			},
			Instructions:    "请输入选项编号（多个选项用逗号分隔，如：1,3,5）",
			InputPrompt:     "请选择 (例如: 1,3,5): ",
			ResultPrefix:    "已选择: ",
			EmptyResultText: "未选择任何选项",
		},
	)
}

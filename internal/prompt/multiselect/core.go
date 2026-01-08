package multiselect

import (
	"fmt"

	"github.com/zevwings/workflow/internal/prompt/common"
	"github.com/zevwings/workflow/internal/prompt/io"
)

// Config 多选功能配置（使用 common.PromptConfig 的别名，保持向后兼容）
type Config = common.PromptConfig

// MultiSelect 多选功能（使用 TerminalIO 接口）
func MultiSelect(message string, options []string, defaultSelected []int, config Config, terminal io.TerminalIO) ([]int, error) {
	if len(options) == 0 {
		return nil, fmt.Errorf("选项列表不能为空")
	}

	handler := NewMultiSelectHandler(options, defaultSelected, config)

	// 验证并清理默认选中项
	selected := handler.ValidateAndCleanDefaults()

	// 格式化提示消息
	promptMsg := config.FormatPrompt(message)

	// 创建原始模式管理器、转义序列解析器和渲染器
	rawModeMgr := io.NewRawModeManager(terminal)
	parser := io.NewEscapeSequenceParser(terminal)
	renderer := io.NewInteractiveRenderer(terminal)

	currentIndex := handler.GetInitialCurrentIndex()

	// 使用原始模式管理器执行交互逻辑
	var result []int
	var resultErr error

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
				terminal,
				renderer,
				len(options),
				getCurrentIndex,
				formatLine,
				"使用 ↑/↓ 导航，空格键切换选择，回车确认",
				config,
			)

			// 使用渲染器渲染提示和初始界面
			if err := renderer.RenderWithPrompt(promptMsg, renderMultiSelect); err != nil {
				return err
			}

			// 使用通用输入处理函数
			err := common.HandleInteractiveInput(
				parser,
				terminal,
				&currentIndex,
				func(idx int, dir string) (int, bool) {
					return handler.ProcessArrowKey(idx, dir)
				},
				func() (bool, error) {
					selectedIndices := mapToSlice(selected)
					selectedText := handler.FormatSelectedOptions(selectedIndices)
					if err := common.FormatResultWithOptions(terminal, promptMsg, selectedText, nil, false); err != nil {
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
					renderMultiSelect(false)
				},
			)
			if err != nil {
				resultErr = err
				return err
			}
			return nil
		},
		func() error {
			// Fallback: 如果无法设置原始模式，使用简单多选
			selectedSlice, err := multiselectFallback(message, options, defaultSelected, config, terminal)
			result = selectedSlice
			return err
		},
	)

	if err != nil {
		return result, err
	}

	if resultErr != nil {
		return nil, resultErr
	}

	return result, nil
}

// MultiSelectDefault 向后兼容的 MultiSelect 函数
func MultiSelectDefault(message string, options []string, defaultSelected []int, config Config) ([]int, error) {
	return MultiSelect(message, options, defaultSelected, config, io.NewStdTerminal())
}

// mapToSlice 将 map[int]bool 转换为排序后的 []int
func mapToSlice(selected map[int]bool) []int {
	if len(selected) == 0 {
		return []int{} // 返回空切片而不是 nil
	}
	var indices []int
	for idx := range selected {
		indices = append(indices, idx)
	}
	// 简单排序（冒泡排序）
	for i := 0; i < len(indices); i++ {
		for j := i + 1; j < len(indices); j++ {
			if indices[i] > indices[j] {
				indices[i], indices[j] = indices[j], indices[i]
			}
		}
	}
	return indices
}

// multiselectFallback 回退方案：如果无法设置原始模式，使用简单多选
func multiselectFallback(message string, options []string, defaultSelected []int, config Config, terminal io.TerminalIO) ([]int, error) {
	handler := NewMultiSelectHandler(options, defaultSelected, config)

	// 格式化提示消息
	promptMsg := config.FormatPrompt(message)

	// 显示选项列表
	terminal.Println(promptMsg)
	terminal.Println("请输入选项编号（多个选项用逗号分隔，如：1,3,5）")
	terminal.Println("")

	selectedMap := handler.ValidateAndCleanDefaults()

	for i, option := range options {
		marker := "[ ]"
		if selectedMap[i] {
			marker = "[x]"
		}
		terminal.Print(fmt.Sprintf("  %s %d. %s\n", marker, i+1, option))
	}

	// 提示输入
	terminal.Print("请选择 (例如: 1,3,5): ")

	// 读取输入
	input, err := terminal.ReadLine()
	if err != nil {
		// 如果读取失败，返回默认值
		return handler.GetDefaultSelectedForFallback(), nil
	}

	// 解析输入
	selectedSlice := handler.ParseCommaSeparatedInput(input)

	// 显示选择结果
	if len(selectedSlice) > 0 {
		selectedText := handler.FormatSelectedOptions(selectedSlice)
		terminal.Println("已选择: " + selectedText)
	} else {
		terminal.Println("未选择任何选项")
	}

	return selectedSlice, nil
}

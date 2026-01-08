package multiselect

import (
	"fmt"

	"github.com/zevwings/workflow/internal/prompt/common"
	"github.com/zevwings/workflow/internal/prompt/io"
)

// Config 多选功能配置
type Config struct {
	// 格式化函数
	FormatPrompt func(message string) string
	FormatAnswer func(value string) string
	FormatHint   func(message string) string
}

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
			// 渲染选项列表
			renderMultiSelect := func(isFirst bool) error {
				if !isFirst {
					renderer.ReRender(func(bool) error {
						// 渲染选项
						for i := range options {
							terminal.MoveToStart()
							line, isHighlighted := handler.FormatOptionLine(i, currentIndex, selected)
							if isHighlighted {
								highlightedLine := config.FormatAnswer(line)
								terminal.Print(highlightedLine)
							} else {
								terminal.Print(line)
							}
							terminal.ClearLine()
							terminal.Println("")
						}

						// 空行
						terminal.MoveToStart()
						terminal.Println("")

						// 显示提示信息
						hintMsg := config.FormatHint("使用 ↑/↓ 导航，空格键切换选择，回车确认")
						terminal.MoveToStart()
						terminal.Print(hintMsg)
						terminal.Print("\r\n")

						terminal.HideCursor()
						return nil
					})
				} else {
					// 首次渲染：渲染选项
					for i := range options {
						terminal.MoveToStart()
						line, isHighlighted := handler.FormatOptionLine(i, currentIndex, selected)
						if isHighlighted {
							highlightedLine := config.FormatAnswer(line)
							terminal.Print(highlightedLine)
						} else {
							terminal.Print(line)
						}
						terminal.ClearLine()
						terminal.Println("")
					}

					// 空行
					terminal.MoveToStart()
					terminal.Println("")

					// 显示提示信息
					hintMsg := config.FormatHint("使用 ↑/↓ 导航，空格键切换选择，回车确认")
					terminal.MoveToStart()
					terminal.Print(hintMsg)
					terminal.Print("\r\n")

					terminal.HideCursor()
				}
				return nil
			}

			// 使用渲染器渲染提示和初始界面
			if err := renderer.RenderWithPrompt(promptMsg, renderMultiSelect); err != nil {
				return err
			}

			// 读取输入
			for {
				keyType, _, err := parser.ReadKey()
				if err != nil {
					return fmt.Errorf("读取输入失败: %w", err)
				}

				// 处理箭头键
				if keyType == io.KeyUp {
					newIndex, shouldRender := handler.ProcessArrowKey(currentIndex, "up")
					if shouldRender {
						currentIndex = newIndex
						renderMultiSelect(false)
					}
					continue
				}
				if keyType == io.KeyDown {
					newIndex, shouldRender := handler.ProcessArrowKey(currentIndex, "down")
					if shouldRender {
						currentIndex = newIndex
						renderMultiSelect(false)
					}
					continue
				}

				// 处理空格键（切换选择状态）
				if keyType == io.KeySpace {
					handler.ToggleSelection(selected, currentIndex)
					renderMultiSelect(false)
					continue
				}

				// 处理回车键（确认选择）
				if keyType == io.KeyEnter {
					selectedIndices := mapToSlice(selected)
					selectedText := handler.FormatSelectedOptions(selectedIndices)
					// handler.FormatSelectedOptions 已经返回格式化后的文本，所以不需要再次格式化
					// 由于 RenderWithPrompt 已经输出了提示消息，所以这里不需要再次输出
					if err := common.FormatResultWithOptions(terminal, promptMsg, selectedText, nil, false); err != nil {
						return err
					}
					result = selectedIndices
					return nil
				}

				// 处理 Ctrl+C
				if keyType == io.KeyCtrlC {
					resultErr = common.HandleCancel(terminal)
					return resultErr
				}

				// 其他字符：静默忽略
			}
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

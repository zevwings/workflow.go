package selectpkg

import (
	"fmt"

	"github.com/zevwings/workflow/internal/prompt/common"
	"github.com/zevwings/workflow/internal/prompt/io"
)

// Config 选择功能配置
type Config struct {
	// 格式化函数
	FormatPrompt func(message string) string
	FormatAnswer func(value string) string
	FormatHint   func(message string) string
}

// Select 选择选项（使用 TerminalIO 接口）
func Select(message string, options []string, defaultIndex int, config Config, terminal io.TerminalIO) (int, error) {
	if len(options) == 0 {
		return -1, fmt.Errorf("选项列表不能为空")
	}

	handler := NewSelectHandler(options, defaultIndex, config)

	// 确保默认索引有效
	currentIndex := handler.ValidateAndAdjustDefaultIndex()

	// 格式化提示消息
	promptMsg := config.FormatPrompt(message)

	// 创建原始模式管理器、转义序列解析器和渲染器
	rawModeMgr := io.NewRawModeManager(terminal)
	parser := io.NewEscapeSequenceParser(terminal)
	renderer := io.NewInteractiveRenderer(terminal)

	// 使用原始模式管理器执行交互逻辑
	var result int
	var resultErr error

	err := rawModeMgr.WithRawModeAndFallback(
		func() error {
			// 渲染选项列表
			renderSelect := func(isFirst bool) error {
				if !isFirst {
					renderer.ReRender(func(bool) error {
						// 渲染选项
						for i := range options {
							terminal.MoveToStart()
							line, isHighlighted := handler.FormatOptionLine(i, currentIndex)
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
						hintMsg := config.FormatHint("使用 ↑/↓ 导航，回车确认")
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
						line, isHighlighted := handler.FormatOptionLine(i, currentIndex)
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
					hintMsg := config.FormatHint("使用 ↑/↓ 导航，回车确认")
					terminal.MoveToStart()
					terminal.Print(hintMsg)
					terminal.Print("\r\n")

					terminal.HideCursor()
				}
				return nil
			}

			// 使用渲染器渲染提示和初始界面
			if err := renderer.RenderWithPrompt(promptMsg, renderSelect); err != nil {
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
						renderSelect(false)
					}
					continue
				}
				if keyType == io.KeyDown {
					newIndex, shouldRender := handler.ProcessArrowKey(currentIndex, "down")
					if shouldRender {
						currentIndex = newIndex
						renderSelect(false)
					}
					continue
				}

				// 处理回车键（确认选择）
				if keyType == io.KeyEnter {
					selectedText := handler.FormatSelectedOption(currentIndex)
					// handler.FormatSelectedOption 已经返回格式化后的文本，所以不需要再次格式化
					// 由于 RenderWithPrompt 已经输出了提示消息，所以这里不需要再次输出
					if err := common.FormatResultWithOptions(terminal, promptMsg, selectedText, nil, false); err != nil {
						return err
					}
					result = currentIndex
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
			// Fallback: 如果无法设置原始模式，使用简单编号选择
			selectedIndex, err := selectFallback(message, options, defaultIndex, config, terminal)
			result = selectedIndex
			return err
		},
	)

	if err != nil {
		return result, err
	}

	if resultErr != nil {
		return -1, resultErr
	}

	return result, nil
}

// SelectDefault 向后兼容的 Select 函数
func SelectDefault(message string, options []string, defaultIndex int, config Config) (int, error) {
	return Select(message, options, defaultIndex, config, io.NewStdTerminal())
}

// selectFallback 回退方案：如果无法设置原始模式，使用简单编号选择
func selectFallback(message string, options []string, defaultIndex int, config Config, terminal io.TerminalIO) (int, error) {
	handler := NewSelectHandler(options, defaultIndex, config)

	// 格式化提示消息
	promptMsg := config.FormatPrompt(message)

	// 显示选项列表
	terminal.Println(promptMsg)
	for i, option := range options {
		marker := " "
		if i == handler.ValidateAndAdjustDefaultIndex() {
			marker = "*"
		}
		terminal.Print(fmt.Sprintf("  %s %d. %s\n", marker, i+1, option))
	}

	// 提示输入
	terminal.Print(fmt.Sprintf("请选择 (1-%d): ", len(options)))

	// 读取输入
	inputLine, err := terminal.ReadLine()
	if err != nil {
		// 如果读取失败，返回默认值
		return handler.ValidateAndAdjustDefaultIndex(), nil
	}

	// 解析输入
	var input int
	_, err = fmt.Sscanf(inputLine, "%d", &input)
	if err != nil {
		// 如果解析失败，返回默认值
		return handler.ValidateAndAdjustDefaultIndex(), nil
	}

	// 验证输入范围
	selectedIndex := handler.ParseNumericInput(input)

	// 显示选择结果
	selectedText := handler.FormatSelectedOption(selectedIndex)
	terminal.Println("已选择: " + selectedText)

	return selectedIndex, nil
}

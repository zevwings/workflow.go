package selectpkg

import (
	"fmt"

	"github.com/zevwings/workflow/internal/prompt/common"
	"github.com/zevwings/workflow/internal/prompt/io"
)

// Config 选择功能配置（使用 common.PromptConfig 的别名，保持向后兼容）
type Config = common.PromptConfig

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
			// 使用通用渲染函数渲染选项列表
			formatLine := func(index int, currentIdx int) (string, bool) {
				return handler.FormatOptionLine(index, currentIdx)
			}
			// 使用闭包捕获 currentIndex 的引用，支持动态更新
			getCurrentIndex := func() int {
				return currentIndex
			}
			renderSelect := common.RenderOptions(
				terminal,
				renderer,
				len(options),
				getCurrentIndex,
				formatLine,
				"使用 ↑/↓ 导航，回车确认",
				config,
			)

			// 使用渲染器渲染提示和初始界面
			if err := renderer.RenderWithPrompt(promptMsg, renderSelect); err != nil {
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
					selectedText := handler.FormatSelectedOption(currentIndex)
					if err := common.FormatResultWithOptions(terminal, promptMsg, selectedText, nil, false); err != nil {
						return false, err
					}
					result = currentIndex
					return true, nil
				},
				nil, // select 不需要处理空格键
				func() {
					renderSelect(false)
				},
			)
			if err != nil {
				resultErr = err
				return err
			}
			return nil
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

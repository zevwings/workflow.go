package common

import (
	"github.com/zevwings/workflow/internal/prompt/io"
)

// FormatOptionLineFunc 格式化选项行的函数类型
// 参数: index - 选项索引, currentIndex - 当前光标位置
// 返回: line - 格式化后的行文本, isHighlighted - 是否需要高亮
type FormatOptionLineFunc func(index int, currentIndex int) (line string, isHighlighted bool)

// RenderOptions 渲染选项列表的通用函数
// 用于 select 和 multiselect 等需要渲染选项列表的场景
//
// 参数:
//   - terminal: 终端接口
//   - renderer: 交互式渲染器
//   - optionsCount: 选项数量
//   - getCurrentIndex: 获取当前光标位置的函数（用于支持动态更新）
//   - formatLine: 格式化选项行的函数
//   - hintText: 提示信息文本
//   - config: 提示配置（用于格式化）
//
// 返回:
//   - renderFn: 渲染函数，可用于 RenderWithPrompt
func RenderOptions(
	terminal io.TerminalIO,
	renderer *io.InteractiveRenderer,
	optionsCount int,
	getCurrentIndex func() int,
	formatLine FormatOptionLineFunc,
	hintText string,
	config PromptConfig,
) io.RenderCallback {
	return func(isFirst bool) error {
		// 渲染选项的通用逻辑
		renderOptions := func() error {
			currentIndex := getCurrentIndex()

			// 渲染所有选项行
			for i := 0; i < optionsCount; i++ {
				line, isHighlighted := formatLine(i, currentIndex)
				if isHighlighted {
					highlightedLine := config.FormatAnswer(line)
					terminal.Print(highlightedLine)
				} else {
					terminal.Print(line)
				}
				// 重置格式，确保每行独立，避免格式影响后续行
				terminal.ResetFormat()
				terminal.Print("\r\n")
			}

			// 空行
			terminal.Println("")

			// 显示提示信息
			hintMsg := config.FormatHint(hintText)
			terminal.Print(hintMsg)
			terminal.Print("\r\n")

			terminal.HideCursor()
			return nil
		}

		if !isFirst {
			// 重新渲染：先恢复光标位置并清除内容
			terminal.RestoreCursor()
			// 重置所有 ANSI 格式，避免之前的格式影响新内容
			terminal.ResetFormat()
			// 移动到行首，确保清除操作从行首开始
			terminal.MoveToStart()
			// 清除从光标到屏幕底部的所有内容
			terminal.ClearToEnd()
			// 然后渲染新内容
			return renderOptions()
		}

		// 首次渲染：直接渲染
		return renderOptions()
	}
}

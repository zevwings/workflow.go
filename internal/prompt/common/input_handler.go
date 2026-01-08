package common

import (
	"fmt"

	"github.com/zevwings/workflow/internal/prompt/io"
)

// ArrowKeyHandler 处理箭头键的回调函数
// 参数: currentIndex - 当前索引, direction - 方向 ("up" 或 "down")
// 返回: newIndex - 新索引, shouldRender - 是否需要重新渲染
type ArrowKeyHandler func(currentIndex int, direction string) (newIndex int, shouldRender bool)

// EnterKeyHandler 处理回车键的回调函数
// 返回: shouldExit - 是否应该退出循环, err - 错误
type EnterKeyHandler func() (shouldExit bool, err error)

// SpaceKeyHandler 处理空格键的回调函数（可选）
// 返回: shouldRender - 是否需要重新渲染
type SpaceKeyHandler func() (shouldRender bool)

// RenderCallback 重新渲染的回调函数
type RenderCallback func()

// HandleInteractiveInput 处理交互式输入的通用循环
// 用于 select、multiselect 等需要处理键盘输入的交互式提示
//
// 参数:
//   - parser: 转义序列解析器
//   - terminal: 终端接口
//   - currentIndex: 当前索引的指针（用于动态更新）
//   - onArrowKey: 处理箭头键的回调函数
//   - onEnter: 处理回车键的回调函数
//   - onSpace: 处理空格键的回调函数（可选，如果为 nil 则忽略空格键）
//   - onRender: 重新渲染的回调函数
//
// 返回:
//   - error: 处理过程中的错误（如用户取消）
func HandleInteractiveInput(
	parser *io.EscapeSequenceParser,
	terminal io.TerminalIO,
	currentIndex *int,
	onArrowKey ArrowKeyHandler,
	onEnter EnterKeyHandler,
	onSpace SpaceKeyHandler,
	onRender RenderCallback,
) error {
	for {
		keyType, _, err := parser.ReadKey()
		if err != nil {
			return fmt.Errorf("读取输入失败: %w", err)
		}

		// 处理箭头键
		if keyType == io.KeyUp {
			newIndex, shouldRender := onArrowKey(*currentIndex, "up")
			if shouldRender {
				*currentIndex = newIndex
				onRender()
			}
			continue
		}
		if keyType == io.KeyDown {
			newIndex, shouldRender := onArrowKey(*currentIndex, "down")
			if shouldRender {
				*currentIndex = newIndex
				onRender()
			}
			continue
		}

		// 处理空格键（可选）
		if keyType == io.KeySpace && onSpace != nil {
			shouldRender := onSpace()
			if shouldRender {
				onRender()
			}
			continue
		}

		// 处理回车键（确认选择）
		if keyType == io.KeyEnter {
			shouldExit, err := onEnter()
			if err != nil {
				return err
			}
			if shouldExit {
				return nil
			}
			continue
		}

		// 处理 Ctrl+C
		if keyType == io.KeyCtrlC {
			return HandleCancel(terminal)
		}

		// 其他字符：静默忽略
	}
}

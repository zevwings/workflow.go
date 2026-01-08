package common

import (
	"fmt"

	"github.com/zevwings/workflow/internal/prompt/io"
)

// HandleCancel 处理用户取消（Ctrl+C）
// 统一的取消处理逻辑，清理终端状态，返回标准错误
//
// 参数:
//   - terminal: 终端接口
//
// 返回:
//   - error: 用户取消的错误
func HandleCancel(terminal io.TerminalIO) error {
	// 恢复光标位置并清除内容
	terminal.RestoreCursor()
	terminal.ClearToEnd()

	// 输出换行（保持终端整洁）
	terminal.Println("")

	// 返回标准错误
	return fmt.Errorf("用户取消输入")
}


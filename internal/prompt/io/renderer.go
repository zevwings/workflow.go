package io

// RenderCallback 渲染回调函数
// 用于在渲染循环中自定义渲染逻辑
type RenderCallback func(isFirst bool) error

// InteractiveRenderer 交互式渲染框架
// 封装渲染循环的通用逻辑，管理光标位置保存/恢复
type InteractiveRenderer struct {
	terminal TerminalIO
}

// NewInteractiveRenderer 创建交互式渲染器
func NewInteractiveRenderer(terminal TerminalIO) *InteractiveRenderer {
	return &InteractiveRenderer{
		terminal: terminal,
	}
}

// RenderWithPrompt 渲染带提示消息的交互界面
// 先输出提示消息，然后保存光标位置，执行渲染回调
//
// 参数:
//   - promptMsg: 提示消息（已格式化）
//   - renderFn: 渲染函数，isFirst 表示是否为首次渲染
//
// 返回:
//   - error: 渲染过程中的错误
func (r *InteractiveRenderer) RenderWithPrompt(promptMsg string, renderFn RenderCallback) error {
	// 先输出提示消息，确保换行并重置颜色
	r.terminal.Print(promptMsg)
	r.terminal.ResetFormat()
	r.terminal.Print("\r\n")
	r.terminal.Print("\r\n")

	// 保存光标位置
	r.terminal.SaveCursor()

	// 执行渲染（首次渲染）
	return renderFn(true)
}

// ReRender 重新渲染
// 恢复光标位置，清除内容，然后执行渲染回调
//
// 参数:
//   - renderFn: 渲染函数，isFirst 始终为 false
//
// 返回:
//   - error: 渲染过程中的错误
func (r *InteractiveRenderer) ReRender(renderFn RenderCallback) error {
	// 恢复光标位置并清除内容
	r.terminal.RestoreCursor()
	// 重置所有 ANSI 格式，避免之前的格式影响新内容
	r.terminal.ResetFormat()
	// 清除从光标到屏幕底部的内容（不需要 MoveToStart，因为 RestoreCursor 已经恢复到正确位置）
	r.terminal.ClearToEnd()

	// 执行渲染（非首次渲染）
	return renderFn(false)
}

// GetTerminal 获取终端接口（用于需要直接访问终端的场景）
func (r *InteractiveRenderer) GetTerminal() TerminalIO {
	return r.terminal
}

package prompt

// baseBuilder Builder 模式的基础结构
//
// 提供通用的 Builder 功能，减少各具体 Builder 的代码重复。
// 各具体 Builder 可以嵌入此结构以复用 Prompt() 方法。
type baseBuilder struct {
	message string
}

// Prompt 设置提示消息
//
// 这是所有 Builder 共有的方法，用于设置提示消息。
//
// 参数:
//   - message: 提示消息文本
//
// 返回:
//   - *baseBuilder: 返回自身，支持链式调用
func (b *baseBuilder) Prompt(message string) *baseBuilder {
	b.message = message
	return b
}

// GetMessage 获取提示消息（供子类使用）
func (b *baseBuilder) GetMessage() string {
	return b.message
}

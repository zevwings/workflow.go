package common

// PromptConfig 提示功能的通用配置
// 用于 select、multiselect、confirm 等交互式提示功能
type PromptConfig struct {
	// FormatPrompt 格式化提示消息的函数
	FormatPrompt func(message string) string

	// FormatAnswer 格式化答案的函数
	FormatAnswer func(value string) string

	// FormatHint 格式化提示信息（如操作说明）的函数
	FormatHint func(message string) string
}

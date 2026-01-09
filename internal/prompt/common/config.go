package common

import (
	"github.com/zevwings/workflow/internal/prompt/io"
)

// PromptConfig 提示功能的通用配置
// 用于 select、multiselect、confirm、form 等交互式提示功能
type PromptConfig struct {
	// FormatPrompt 格式化提示消息的函数
	FormatPrompt func(message string) string

	// FormatAnswer 格式化答案的函数
	FormatAnswer func(value string) string

	// FormatError 格式化错误消息的函数（可选）
	// 用于输入验证错误显示
	FormatError func(message string) string

	// FormatHint 格式化提示信息（如操作说明）的函数
	FormatHint func(message string) string

	// FormatQuestionPrefix 格式化问题前缀 "? " 的函数（可选）
	FormatQuestionPrefix func() string

	// FormatAnswerPrefix 格式化答案前缀 "> " 的函数（可选）
	FormatAnswerPrefix func() string

	// FormatResultTitle 格式化完成后显示的 title 的函数（可选）
	// 参数: originalMessage - 原始的提示消息, resultValue - 用户输入/选择的值
	// 返回: 格式化后的 title 文本
	// 如果为 nil，则使用原始的 message
	FormatResultTitle func(originalMessage string, resultValue string) string
}

// MergeConfig 合并配置，将 base 和 override 合并，override 中的非 nil 字段会覆盖 base
// 如果 base 为 nil，则返回 override 的副本
// 如果 override 为 nil，则返回 base 的副本
func MergeConfig(base, override *PromptConfig) PromptConfig {
	if base == nil && override == nil {
		return PromptConfig{}
	}
	if base == nil {
		return *override
	}
	if override == nil {
		return *base
	}

	merged := *base
	if override.FormatPrompt != nil {
		merged.FormatPrompt = override.FormatPrompt
	}
	if override.FormatAnswer != nil {
		merged.FormatAnswer = override.FormatAnswer
	}
	if override.FormatError != nil {
		merged.FormatError = override.FormatError
	}
	if override.FormatHint != nil {
		merged.FormatHint = override.FormatHint
	}
	if override.FormatQuestionPrefix != nil {
		merged.FormatQuestionPrefix = override.FormatQuestionPrefix
	}
	if override.FormatAnswerPrefix != nil {
		merged.FormatAnswerPrefix = override.FormatAnswerPrefix
	}
	if override.FormatResultTitle != nil {
		merged.FormatResultTitle = override.FormatResultTitle
	}
	return merged
}

// WithResultTitle 为配置添加或覆盖 FormatResultTitle
// 如果 resultTitle 为空字符串，返回原始配置
// 否则返回新配置，其中 FormatResultTitle 返回固定的 resultTitle 字符串
func WithResultTitle(config PromptConfig, resultTitle string) PromptConfig {
	if resultTitle == "" {
		return config
	}

	titleStr := resultTitle
	config.FormatResultTitle = func(originalMessage string, resultValue string) string {
		return titleStr
	}
	return config
}

// FillDefaults 填充配置的默认值
// 如果 config 中的某个字段为 nil，则使用 defaultConfig 中对应的字段
func FillDefaults(config PromptConfig, defaultConfig PromptConfig) PromptConfig {
	if config.FormatPrompt == nil {
		config.FormatPrompt = defaultConfig.FormatPrompt
	}
	if config.FormatAnswer == nil {
		config.FormatAnswer = defaultConfig.FormatAnswer
	}
	if config.FormatError == nil {
		config.FormatError = defaultConfig.FormatError
	}
	if config.FormatHint == nil {
		config.FormatHint = defaultConfig.FormatHint
	}
	if config.FormatQuestionPrefix == nil {
		config.FormatQuestionPrefix = defaultConfig.FormatQuestionPrefix
	}
	if config.FormatAnswerPrefix == nil {
		config.FormatAnswerPrefix = defaultConfig.FormatAnswerPrefix
	}
	// FormatResultTitle 不填充默认值，保持为 nil 表示使用原始 message
	return config
}

// BasePromptConfig 基础提示配置（通用参数）
// 包含所有 prompt 函数都需要的通用参数
type BasePromptConfig struct {
	// Message 提示消息
	Message string
	// Config 提示功能配置
	Config PromptConfig
	// Terminal 终端接口
	Terminal io.TerminalIO
}

// BuildConfigWithDefaults 构建配置（带默认值填充）
// 这是一个统一的配置构建函数，用于统一配置构建逻辑
// 按照优先级：localConfig > defaultConfig
//
// 参数:
//   - localConfig: 局部配置（可选，如果为 nil 则只使用 defaultConfig）
//   - defaultConfig: 默认配置（必须提供）
//
// 返回:
//   - 合并后的配置
func BuildConfigWithDefaults(localConfig *PromptConfig, defaultConfig PromptConfig) PromptConfig {
	if localConfig == nil {
		return defaultConfig
	}
	return FillDefaults(*localConfig, defaultConfig)
}

// BuildConfigWithResultTitle 构建配置并设置 ResultTitle
// 这是一个便捷函数，用于构建配置并同时设置 ResultTitle
//
// 参数:
//   - localConfig: 局部配置（可选）
//   - defaultConfig: 默认配置
//   - resultTitle: 结果标题（可选，如果为空字符串则忽略）
//
// 返回:
//   - 合并后的配置
func BuildConfigWithResultTitle(localConfig *PromptConfig, defaultConfig PromptConfig, resultTitle string) PromptConfig {
	config := BuildConfigWithDefaults(localConfig, defaultConfig)
	return WithResultTitle(config, resultTitle)
}

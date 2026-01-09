package prompt

import (
	"strings"

	"github.com/zevwings/workflow/internal/prompt/common"
	"github.com/zevwings/workflow/internal/prompt/form"
)

// FormBuilder Form 构建器类型（重新导出）
type FormBuilder = form.FormBuilder

// FormResult Form 结果类型（重新导出）
type FormResult = form.FormResult

// Condition Form 条件函数类型（重新导出）
type Condition = form.Condition

// promptInputProvider prompt 包的输入提供者实现
type promptInputProvider struct{}

// convertValidator 转换 Validator 类型的辅助函数
// form 包中的 Validator 已经是 input.Validator 类型，可以直接使用
func convertValidator(validator interface{}) Validator {
	if validator == nil {
		return nil
	}
	if v, ok := validator.(Validator); ok {
		return v
	}
	return nil
}

func (p *promptInputProvider) AskInput(field form.InputField) (string, error) {
	return AskInput(InputField{
		Message:      field.Message,
		DefaultValue: field.DefaultValue,
		Validator:    convertValidator(field.Validator),
		ResultTitle:  field.ResultTitle,
		Config:       field.Config,
	})
}

func (p *promptInputProvider) AskPassword(field form.PasswordField) (string, error) {
	return AskPassword(PasswordField{
		Message:      field.Message,
		DefaultValue: field.DefaultValue,
		Validator:    convertValidator(field.Validator),
		ResultTitle:  field.ResultTitle,
		Config:       field.Config,
	})
}

// init 初始化 Form 模块配置
func init() {
	// 设置 Form 模块的配置，使其能够使用 prompt 包的格式化函数
	form.SetPromptConfig(common.PromptConfig{
		FormatPrompt:         formatTitle,
		FormatAnswer:         formatAnswer,
		FormatError:          formatError,
		FormatHint:           formatHint,
		FormatQuestionPrefix: formatQuestionPrefix,
		FormatAnswerPrefix:   formatAnswerPrefix,
		// FormatResultTitle 不设置默认值，需要手动设置
	})
	form.SetInputProvider(&promptInputProvider{})
}

// Form 创建一个 Form 构建器（链式调用）
// 使用方式：prompt.Form().AddXxx().Run()
func Form() *FormBuilder {
	return form.NewFormBuilder()
}

// AskForm 函数式调用 Form（保持向后兼容）
// 使用方式：prompt.AskForm(builder)
func AskForm(builder *FormBuilder) (*FormResult, error) {
	executor := form.NewFormExecutor()
	return executor.Execute(builder)
}

// SetFormFormatResultTitle 设置 Form 的 FormatResultTitle 函数
// 用于自定义输入完成后显示的 title 格式
// 如果设置为 nil，则不更新 title（使用原始 message）
func SetFormFormatResultTitle(formatFunc func(originalMessage string, resultValue string) string) {
	currentConfig := form.GetPromptConfig()
	// 使用 MergeConfig 合并配置，只更新 FormatResultTitle
	updatedConfig := common.MergeConfig(&currentConfig, &common.PromptConfig{
		FormatResultTitle: formatFunc,
	})
	form.SetPromptConfig(updatedConfig)
}

// FormatResultTitleForForm 为 Form 格式化完成后显示的 title 的辅助函数
// 将 "Please enter your X" 转换为 "Your X"
// 可以在 setup.go 中使用：prompt.SetFormFormatResultTitle(prompt.FormatResultTitleForForm)
func FormatResultTitleForForm(originalMessage string, resultValue string) string {
	// 移除常见的提示前缀
	title := originalMessage
	prefixes := []string{
		"Please enter your ",
		"Please enter ",
		"Enter your ",
		"Enter ",
	}
	for _, prefix := range prefixes {
		if strings.HasPrefix(title, prefix) {
			title = strings.TrimPrefix(title, prefix)
			// 首字母大写
			if len(title) > 0 {
				title = strings.ToUpper(title[:1]) + title[1:]
			}
			break
		}
	}
	// 移除后缀中的提示信息
	suffixes := []string{
		" (required)",
		" (press Enter to keep)",
		" [required]",
		" [press Enter to keep]",
	}
	for _, suffix := range suffixes {
		title = strings.TrimSuffix(title, suffix)
	}
	return title
}

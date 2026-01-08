package prompt

import (
	"github.com/zevwings/workflow/internal/prompt/form"
)

// FormBuilder Form 构建器类型（重新导出）
type FormBuilder = form.FormBuilder

// FormResult Form 结果类型（重新导出）
type FormResult = form.FormResult

// Condition Form 条件函数类型（重新导出）
type Condition = form.Condition

// init 初始化 Form 模块配置
func init() {
	// 设置 Form 模块的配置，使其能够使用 prompt 包的格式化函数
	form.SetFormConfig(form.FormConfig{
		FormatPrompt: formatTitle,
		FormatAnswer: formatAnswer,
		FormatError:  formatError,
		FormatHint:   formatHint,
		AskInputFunc: func(message string, defaultValue string, validator interface{}) (string, error) {
			var v Validator
			if validator != nil {
				if val, ok := validator.(Validator); ok {
					v = val
				}
			}
			return AskInput(message, defaultValue, v)
		},
		AskPasswordFunc: func(message string, validator interface{}) (string, error) {
			var v Validator
			if validator != nil {
				if val, ok := validator.(Validator); ok {
					v = val
				}
			}
			return AskPassword(message, v)
		},
	})
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

package confirm

import (
	"github.com/zevwings/workflow/internal/prompt/common"
)

// confirmFallbackAdapter ConfirmHandler 的 TypedFallbackHandler 适配器（类型安全版本）
// 用于将 ConfirmHandler 适配为 common.TypedFallbackHandler[bool] 接口
type confirmFallbackAdapter struct {
	handler *ConfirmHandler
}

// FormatPromptText 格式化提示文本
func (a *confirmFallbackAdapter) FormatPromptText(message string) string {
	return a.handler.FormatPromptText(message)
}

// FormatAnswer 格式化答案文本
func (a *confirmFallbackAdapter) FormatAnswer(result bool) string {
	return a.handler.FormatAnswer(result)
}

// ProcessLineInput 处理一行输入（用于 fallback 模式）
func (a *confirmFallbackAdapter) ProcessLineInput(input string) (bool, error) {
	result, err := a.handler.ProcessLineInput(input)
	if err != nil {
		return a.handler.defaultYes, err
	}
	if result == nil {
		return a.handler.defaultYes, nil
	}
	return *result, nil
}

// GetDefaultResult 获取默认结果
func (a *confirmFallbackAdapter) GetDefaultResult() bool {
	return a.handler.defaultYes
}

// newConfirmFallbackAdapter 创建 confirm fallback 适配器
func newConfirmFallbackAdapter(handler *ConfirmHandler) common.TypedFallbackHandler[bool] {
	return &confirmFallbackAdapter{handler: handler}
}

package confirm

import (
	"fmt"
	"strings"
)

// ConfirmHandler 处理确认逻辑（纯业务逻辑，无 I/O 操作）
// 负责处理输入字符、格式化提示和答案
type ConfirmHandler struct {
	defaultYes bool
	config     Config
}

// NewConfirmHandler 创建确认处理器
func NewConfirmHandler(defaultYes bool, config Config) *ConfirmHandler {
	return &ConfirmHandler{
		defaultYes: defaultYes,
		config:     config,
	}
}

// ProcessInput 处理单个字符输入
// 返回：
//   - result: 处理结果（如果确定则返回 bool 指针，否则为 nil）
//   - shouldContinue: 是否需要继续等待输入
//   - err: 错误（如用户取消）
func (h *ConfirmHandler) ProcessInput(char byte) (result *bool, shouldContinue bool, err error) {
	// 处理回车键（使用默认值）
	if char == '\r' || char == '\n' {
		return &h.defaultYes, false, nil
	}

	// 处理 Ctrl+C
	if char == 3 { // Ctrl+C
		return nil, false, fmt.Errorf("用户取消输入")
	}

	// 转换为小写进行比较
	charLower := strings.ToLower(string(char))[0]

	// 处理 yes 输入：y 或 Y
	if charLower == 'y' {
		yes := true
		return &yes, false, nil
	}

	// 处理 no 输入：n 或 N
	if charLower == 'n' {
		no := false
		return &no, false, nil
	}

	// 其他字符：静默忽略，继续等待有效输入
	return nil, true, nil
}

// FormatPromptText 格式化提示文本
// 根据 defaultYes 构建不同的提示（显示 【Y/n】 或 【y/N】）
func (h *ConfirmHandler) FormatPromptText(message string) string {
	promptMsg := h.config.FormatPrompt(message)
	var hintText string
	if h.defaultYes {
		hintText = "[Y/n]"
	} else {
		hintText = "[y/N]"
	}

	// 如果配置了 FormatHint，使用它来格式化提示文本
	if h.config.FormatHint != nil {
		hintText = h.config.FormatHint(hintText)
	}

	return fmt.Sprintf("%s %s ", promptMsg, hintText)
}

// FormatAnswer 格式化答案
// 根据布尔值返回格式化的 "yes" 或 "no"
func (h *ConfirmHandler) FormatAnswer(value bool) string {
	if value {
		return h.config.FormatAnswer("yes")
	}
	return h.config.FormatAnswer("no")
}

// ProcessLineInput 处理一行输入（用于 fallback 模式）
// 返回：
//   - result: 处理结果（如果确定则返回 bool 指针，否则为 nil）
//   - err: 错误
func (h *ConfirmHandler) ProcessLineInput(input string) (result *bool, err error) {
	// 清理输入
	input = strings.TrimSpace(strings.ToLower(input))

	// 处理空输入（使用默认值）
	if input == "" {
		return &h.defaultYes, nil
	}

	// 处理 yes 输入
	if input == "y" || input == "yes" {
		yes := true
		return &yes, nil
	}

	// 处理 no 输入
	if input == "n" || input == "no" {
		no := false
		return &no, nil
	}

	// 非法输入，返回默认值
	return &h.defaultYes, nil
}

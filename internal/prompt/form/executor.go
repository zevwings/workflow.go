package form

import (
	"fmt"
	"strings"

	"github.com/mattn/go-runewidth"
	"github.com/zevwings/workflow/internal/prompt/common"
	"github.com/zevwings/workflow/internal/prompt/confirm"
	"github.com/zevwings/workflow/internal/prompt/io"
	multiselectpkg "github.com/zevwings/workflow/internal/prompt/multiselect"
	selectpkg "github.com/zevwings/workflow/internal/prompt/select"
)

// FormExecutor 表单执行器
type FormExecutor struct {
	config common.PromptConfig
}

// 全局配置（由 prompt 包在初始化时设置）
var globalPromptConfig common.PromptConfig
var globalInputProvider InputProvider

// SetPromptConfig 设置全局 Prompt 配置（由 prompt 包调用）
func SetPromptConfig(config common.PromptConfig) {
	globalPromptConfig = config
}

// SetInputProvider 设置全局 InputProvider（由 prompt 包调用）
func SetInputProvider(provider InputProvider) {
	globalInputProvider = provider
}

// GetPromptConfig 获取全局 Prompt 配置（用于访问 FormatResultTitle 等）
func GetPromptConfig() common.PromptConfig {
	return globalPromptConfig
}

// NewFormExecutor 创建新的表单执行器
func NewFormExecutor() *FormExecutor {
	return &FormExecutor{
		config: globalPromptConfig,
	}
}

// Execute 执行表单字段序列
func (e *FormExecutor) Execute(builder *FormBuilder) (*FormResult, error) {
	return e.executeWithLevel(builder, 0)
}

// executeWithLevel 执行表单字段序列（带层级信息）
func (e *FormExecutor) executeWithLevel(builder *FormBuilder, level int) (*FormResult, error) {
	// 获取终端接口用于输出分割线
	terminal := io.NewStdTerminal()
	title := builder.GetTitle()

	// 判断是主表单（level == 0）还是嵌套表单（level > 0）
	isMainForm := level == 0

	// 输出开始分割线
	if title != "" {
		e.printSeparator(terminal, title, "start", isMainForm)
	}

	result := NewFormResult()
	fields := builder.GetFields()

	for _, field := range fields {
		// 评估条件（如果有）
		if field.Condition != nil {
			if !field.Condition(result) {
				// 条件不满足，跳过该字段
				continue
			}
		}

		// 执行字段
		value, err := e.executeField(field, result, level)
		if err != nil {
			return nil, fmt.Errorf("执行字段 %s 失败: %w", field.Key, err)
		}

		// 收集结果
		result.Set(field.Key, value)
	}

	// 输出结束分割线
	if title != "" {
		e.printSeparator(terminal, title, "end", isMainForm)
	}

	return result, nil
}

// executeField 执行单个字段
func (e *FormExecutor) executeField(field FormField, currentResult *FormResult, level int) (interface{}, error) {
	switch field.Type {
	case FieldTypeConfirm:
		return e.executeConfirm(field)
	case FieldTypeInput:
		return e.executeInput(field)
	case FieldTypePassword:
		return e.executePassword(field)
	case FieldTypeSelect:
		return e.executeSelect(field)
	case FieldTypeMultiSelect:
		return e.executeMultiSelect(field)
	case FieldTypeForm:
		return e.executeForm(field, level)
	default:
		return nil, fmt.Errorf("未知的字段类型: %s", field.Type)
	}
}

// executeConfirm 执行确认字段
func (e *FormExecutor) executeConfirm(field FormField) (bool, error) {
	defaultValue, ok := field.DefaultValue.(bool)
	if !ok {
		defaultValue = false
	}
	config := e.buildConfigWithResultTitle(field.ResultTitle)
	return confirm.Confirm(confirm.ConfirmConfig{
		BasePromptConfig: common.BasePromptConfig{
			Message:  field.Prompt,
			Config:   config,
			Terminal: io.NewStdTerminal(),
		},
		DefaultYes: defaultValue,
	})
}

// executeInput 执行输入字段
func (e *FormExecutor) executeInput(field FormField) (string, error) {
	defaultValue := ""
	if field.DefaultValue != nil {
		if str, ok := field.DefaultValue.(string); ok {
			defaultValue = str
		}
	}
	config := e.buildConfigWithResultTitle(field.ResultTitle)
	if globalInputProvider == nil {
		return "", fmt.Errorf("InputProvider 未设置，请确保 prompt 包已正确初始化")
	}
	return globalInputProvider.AskInput(InputField{
		Message:      field.Prompt,
		DefaultValue: defaultValue,
		Validator:    field.Validator,
		ResultTitle:  "", // 使用 Config 中的 FormatResultTitle
		Config:       &config,
	})
}

// executePassword 执行密码字段
func (e *FormExecutor) executePassword(field FormField) (string, error) {
	defaultValue := ""
	if field.DefaultValue != nil {
		if str, ok := field.DefaultValue.(string); ok {
			defaultValue = str
		}
	}
	config := e.buildConfigWithResultTitle(field.ResultTitle)
	if globalInputProvider == nil {
		return "", fmt.Errorf("InputProvider 未设置，请确保 prompt 包已正确初始化")
	}
	return globalInputProvider.AskPassword(PasswordField{
		Message:      field.Prompt,
		DefaultValue: defaultValue,
		Validator:    field.Validator,
		ResultTitle:  "", // 使用 Config 中的 FormatResultTitle
		Config:       &config,
	})
}

// executeSelect 执行单选字段
func (e *FormExecutor) executeSelect(field FormField) (int, error) {
	config := e.buildConfigWithResultTitle(field.ResultTitle)
	return selectpkg.Select(selectpkg.SelectConfig{
		BasePromptConfig: common.BasePromptConfig{
			Message:  field.Prompt,
			Config:   config,
			Terminal: io.NewStdTerminal(),
		},
		Options:      field.Options,
		DefaultIndex: field.DefaultIndex,
	})
}

// executeMultiSelect 执行多选字段
func (e *FormExecutor) executeMultiSelect(field FormField) ([]int, error) {
	config := e.buildConfigWithResultTitle(field.ResultTitle)
	return multiselectpkg.MultiSelect(multiselectpkg.MultiSelectConfig{
		BasePromptConfig: common.BasePromptConfig{
			Message:  field.Prompt,
			Config:   config,
			Terminal: io.NewStdTerminal(),
		},
		Options:         field.Options,
		DefaultSelected: field.DefaultSelected,
	})
}

// executeForm 执行嵌套表单字段
func (e *FormExecutor) executeForm(field FormField, level int) (*FormResult, error) {
	if field.NestedForm == nil {
		return nil, fmt.Errorf("嵌套表单不能为空")
	}
	return e.executeWithLevel(field.NestedForm, level+1)
}

// 注意：newDefaultConfig 已移除，直接使用 globalPromptConfig

// buildConfigWithResultTitle 构建带 ResultTitle 的配置
// 如果 resultTitle 为空，返回原始配置；否则创建新配置并设置 FormatResultTitle
func (e *FormExecutor) buildConfigWithResultTitle(resultTitle string) common.PromptConfig {
	return common.WithResultTitle(e.config, resultTitle)
}

// printSeparator 打印分割线
// isMainForm: true 表示主表单（完整格式），false 表示子表单（单行格式）
func (e *FormExecutor) printSeparator(terminal io.TerminalIO, title string, suffix string, isMainForm bool) {
	const separatorChar = "─"
	const separatorLength = 72 // 分割线长度（显示宽度）

	// 构建文本：title + " " + suffix（首字母大写）
	var suffixCapitalized string
	if len(suffix) > 0 {
		suffixCapitalized = strings.ToUpper(suffix[:1]) + suffix[1:]
	} else {
		suffixCapitalized = suffix
	}
	text := fmt.Sprintf("%s %s", title, suffixCapitalized)

	if isMainForm {
		// 主表单：完整格式（3行：分割线+标题+分割线）
		e.printMainFormSeparator(terminal, text, separatorChar, separatorLength)
	} else {
		// 子表单：单行格式（─...─ 标题 ─...─）
		e.printNestedFormSeparator(terminal, text, separatorChar, separatorLength)
	}
}

// printMainFormSeparator 打印主表单分割线（完整格式）
// 格式：
// ────────────────────────────────────────────────────────────────
//
//	Form Title Start
//
// ────────────────────────────────────────────────────────────────
func (e *FormExecutor) printMainFormSeparator(terminal io.TerminalIO, text string, separatorChar string, totalWidth int) {
	// 计算文本的显示宽度
	textDisplayWidth := runewidth.StringWidth(text)

	// 计算居中位置（基于显示宽度）
	remainingWidth := totalWidth - textDisplayWidth
	leftPadding := remainingWidth / 2
	rightPadding := remainingWidth - leftPadding

	// 构建分割线
	separatorLine := strings.Repeat(separatorChar, totalWidth)

	// 构建文本行（居中，使用空格填充）
	textLine := strings.Repeat(" ", leftPadding) + text + strings.Repeat(" ", rightPadding)

	// 输出分割线（前后各加一个空行以提高可读性）
	terminal.ResetFormat()
	terminal.Println("")
	terminal.Println(separatorLine)
	terminal.Println(textLine)
	terminal.Println(separatorLine)
	terminal.Println("")
	terminal.ResetFormat()
}

// printNestedFormSeparator 打印嵌套表单分割线（单行格式）
// 格式：───────────────────────── 用户信息 Start ────────────────────────────
func (e *FormExecutor) printNestedFormSeparator(terminal io.TerminalIO, text string, separatorChar string, totalWidth int) {
	// 计算文本的显示宽度
	textDisplayWidth := runewidth.StringWidth(text)

	// 计算左右两侧的分割线长度（基于显示宽度）
	// 文本前后各有一个空格，所以需要减去 2
	remainingWidth := totalWidth - textDisplayWidth - 2
	leftDashes := remainingWidth / 2
	rightDashes := remainingWidth - leftDashes

	// 构建单行分割线：─...─ 标题 ─...─
	separatorLine := strings.Repeat(separatorChar, leftDashes) + " " + text + " " + strings.Repeat(separatorChar, rightDashes)

	// 输出分割线（前后各加一个空行以提高可读性）
	terminal.ResetFormat()
	terminal.Println("")
	terminal.Println(separatorLine)
	terminal.Println("")
	terminal.ResetFormat()
}

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
	config     common.PromptConfig
	formConfig FormConfig
}

// NewFormExecutor 创建新的表单执行器
// formConfig 通过 SetFormConfig 设置，由 prompt 包在初始化时调用
var globalFormConfig FormConfig

// SetFormConfig 设置全局 Form 配置（由 prompt 包调用）
func SetFormConfig(config FormConfig) {
	globalFormConfig = config
}

// NewFormExecutor 创建新的表单执行器
func NewFormExecutor() *FormExecutor {
	return &FormExecutor{
		config:     newDefaultConfig(globalFormConfig),
		formConfig: globalFormConfig,
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
	return confirm.ConfirmDefault(field.Prompt, defaultValue, e.config)
}

// executeInput 执行输入字段
// 使用 formConfig 中的 AskInputFunc 来保持格式一致
func (e *FormExecutor) executeInput(field FormField) (string, error) {
	defaultValue := ""
	if field.DefaultValue != nil {
		if str, ok := field.DefaultValue.(string); ok {
			defaultValue = str
		}
	}
	// 使用 formConfig 中的函数，保持格式一致
	if e.formConfig.AskInputFunc != nil {
		return e.formConfig.AskInputFunc(field.Prompt, defaultValue, field.Validator)
	}
	// 如果没有设置，返回错误
	return "", fmt.Errorf("AskInputFunc 未设置，请确保 prompt 包已正确初始化")
}

// executePassword 执行密码字段
// 使用 formConfig 中的 AskPasswordFunc 来保持格式一致
func (e *FormExecutor) executePassword(field FormField) (string, error) {
	// 使用 formConfig 中的函数，保持格式一致
	if e.formConfig.AskPasswordFunc != nil {
		return e.formConfig.AskPasswordFunc(field.Prompt, field.Validator)
	}
	// 如果没有设置，返回错误
	return "", fmt.Errorf("AskPasswordFunc 未设置，请确保 prompt 包已正确初始化")
}

// executeSelect 执行单选字段
func (e *FormExecutor) executeSelect(field FormField) (int, error) {
	return selectpkg.SelectDefault(field.Prompt, field.Options, field.DefaultIndex, e.config)
}

// executeMultiSelect 执行多选字段
func (e *FormExecutor) executeMultiSelect(field FormField) ([]int, error) {
	return multiselectpkg.MultiSelectDefault(field.Prompt, field.Options, field.DefaultSelected, e.config)
}

// executeForm 执行嵌套表单字段
func (e *FormExecutor) executeForm(field FormField, level int) (*FormResult, error) {
	if field.NestedForm == nil {
		return nil, fmt.Errorf("嵌套表单不能为空")
	}
	return e.executeWithLevel(field.NestedForm, level+1)
}

// newDefaultConfig 创建默认配置（使用 formConfig 中的格式化函数）
func newDefaultConfig(formConfig FormConfig) common.PromptConfig {
	return common.PromptConfig{
		FormatPrompt: formConfig.FormatPrompt,
		FormatAnswer: formConfig.FormatAnswer,
		FormatHint:   formConfig.FormatHint,
	}
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

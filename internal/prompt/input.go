package prompt

import (
	"fmt"
	"strings"

	"github.com/zevwings/workflow/internal/prompt/common"
	"github.com/zevwings/workflow/internal/prompt/input"
)

// Validator 验证函数类型（重新导出，保持向后兼容）
type Validator = input.Validator

// 常量定义
const (
	// MaxPasswordDisplayLength 密码显示的最大长度（用于掩码显示）
	MaxPasswordDisplayLength = 20
	// PasswordMask 密码掩码显示文本
	PasswordMask = "****"
	// DefaultQuestionPrefix 默认问题前缀
	DefaultQuestionPrefix = "? "
	// DefaultInputPrefix 默认输入前缀
	DefaultInputPrefix = "> "
)

// ANSI 转义序列常量
const (
	// ANSI 光标控制
	ansiCarriageReturn = "\r"      // 回到行首
	ansiClearToEnd     = "\033[K"  // 清除到行尾
	ansiMoveUp         = "\033[A"  // 上移一行
	ansiResetFormat    = "\033[0m" // 重置所有格式
)

// InputField 输入字段配置
type InputField struct {
	// Message 提示消息
	Message string
	// DefaultValue 默认值（可选）
	DefaultValue string
	// Validator 验证器（可选）
	Validator Validator
	// ResultTitle 输入完成后显示的 title（可选）
	// 如果设置，将优先于全局的 FormatResultTitle 使用
	ResultTitle string
	// Config 自定义配置（可选）
	// 如果设置，将使用此配置而不是默认配置
	Config *common.PromptConfig
}

// PasswordField 密码字段配置
type PasswordField struct {
	// Message 提示消息
	Message string
	// DefaultValue 默认值（可选，空字符串表示无默认值）
	DefaultValue string
	// Validator 验证器（可选）
	Validator Validator
	// ResultTitle 输入完成后显示的 title（可选）
	// 如果设置，将优先于全局的 FormatResultTitle 使用
	ResultTitle string
	// Config 自定义配置（可选）
	// 如果设置，将使用此配置而不是默认配置
	Config *common.PromptConfig
}

// inputFunc 统一的输入函数，通过 isPassword 参数控制是否使用密文模式（私有函数）
// placeholder 参数用于显示占位符文本（仅在非密码模式下有效）
func inputFunc(message string, defaultValue string, placeholder string, isPassword bool, validator Validator) (string, error) {
	return inputFuncWithConfig(message, defaultValue, placeholder, isPassword, validator, nil)
}

// inputFuncWithConfig 统一的输入函数（支持自定义配置）
// placeholder 参数用于显示占位符文本（仅在非密码模式下有效）
// config 参数为可选的 PromptConfig，如果为 nil 则使用默认配置
func inputFuncWithConfig(message string, defaultValue string, placeholder string, isPassword bool, validator Validator, config *common.PromptConfig) (string, error) {
	// 构建配置
	promptConfig := buildPromptConfig(config)
	promptMsg := promptConfig.FormatPrompt(message)

	// 显示提示和默认值
	displayPromptWithDefault(promptMsg, defaultValue, isPassword, promptConfig)

	// 构建编辑器配置
	editorConfig := buildEditorConfig()
	promptText := DefaultInputPrefix

	hasError := false
	for {
		// 如果有错误，清除上一轮的错误提示和输入行
		if hasError {
			clearErrorAndInputLines()
		}

		// 创建有效验证器（处理默认值场景）
		effectiveValidator := createEffectiveValidator(validator, defaultValue)

		// 读取输入值
		value, err := readInputValue(promptText, placeholder, isPassword, effectiveValidator, editorConfig)
		if err != nil {
			return "", err
		}

		// 应用默认值（如果用户未输入）
		if defaultValue != "" && strings.TrimSpace(value) == "" {
			value = defaultValue
		}

		// 验证输入（密码模式已在 ReadLineCore 中验证，这里只验证普通输入）
		if validator != nil && !isPassword {
			if err := validator(value); err != nil {
				handleInputError(err)
				hasError = true
				continue
			}
		}

		// 验证通过，显示格式化的结果
		hadError := hasError
		hasError = false

		clearAndDisplayResult(hadError, message, value, isPassword, promptConfig)
		return value, nil
	}
}

// buildPromptConfig 构建提示配置
func buildPromptConfig(config *common.PromptConfig) common.PromptConfig {
	if config != nil {
		return common.FillDefaults(*config, newDefaultConfig())
	}
	return newDefaultConfig()
}

// displayPromptWithDefault 显示提示信息和默认值
// 参考 huh + bubbletea 的样式：
//  1. Title（提示信息）单独显示一行，并在其后显示默认值
//     示例：请输入您的邮箱[user@example.com]
//  2. 输入框使用 "> " 前缀（类似 huh 的样式）
//  3. Placeholder 显示在输入框内（灰色斜体）
func displayPromptWithDefault(promptMsg, defaultValue string, isPassword bool, config common.PromptConfig) {
	questionPrefix := DefaultQuestionPrefix
	if config.FormatQuestionPrefix != nil {
		questionPrefix = config.FormatQuestionPrefix()
	}

	if defaultValue != "" {
		if !isPassword {
			// 普通输入：直接显示真实默认值
			titleText := fmt.Sprintf("%s%s[%s]", questionPrefix, promptMsg, defaultValue)
			fmt.Println(titleText)
		} else {
			// 密码输入：永远只显示固定长度的掩码，避免泄露真实长度
			titleText := fmt.Sprintf("%s%s[%s]", questionPrefix, promptMsg, PasswordMask)
			fmt.Println(titleText)
		}
	} else {
		titleText := fmt.Sprintf("%s%s", questionPrefix, promptMsg)
		fmt.Println(titleText)
	}
}

// buildEditorConfig 构建编辑器配置
func buildEditorConfig() input.Config {
	t := GetTheme()
	return input.Config{
		FormatPlaceholder: func(text string) string {
			return input.FormatPlaceholder(text, t.HintStyle, t.EnableColor)
		},
		FormatError: formatError,
		HintStyle:   t.HintStyle,
		ErrorStyle:  t.ErrorStyle,
		EnableColor: t.EnableColor,
	}
}

// createEffectiveValidator 创建有效验证器
// 针对有默认值的场景，为底层输入函数构造一个"宽松版"验证器：
// - 实时验证 / 回车验证时：空字符串视为合法（不报错），让用户可以直接回车退出输入循环
// - 真实的默认值验证放在后面统一处理
func createEffectiveValidator(validator Validator, defaultValue string) Validator {
	if defaultValue == "" || validator == nil {
		return validator
	}
	return func(v string) error {
		if strings.TrimSpace(v) == "" {
			return nil
		}
		return validator(v)
	}
}

// readInputValue 读取输入值
func readInputValue(promptText, placeholder string, isPassword bool, validator Validator, editorConfig input.Config) (string, error) {
	if isPassword {
		// 密码模式：使用通用编辑内核，但通过 echo 函数以 * 方式显示输入内容
		// 限制显示的星号数量，避免长 token 导致终端显示混乱
		return input.ReadLineCoreDefault(promptText, validator, func(b []byte) string {
			if len(b) <= MaxPasswordDisplayLength {
				return strings.Repeat("*", len(b))
			}
			// 如果输入很长，显示固定数量的星号，避免终端混乱
			return strings.Repeat("*", MaxPasswordDisplayLength)
		}, formatError)
	}

	// 普通输入模式：如果有 placeholder，使用字符级输入；否则也使用字符级输入（支持光标移动）
	if placeholder != "" {
		return input.ReadWithPlaceholderDefault(promptText, placeholder, validator, editorConfig)
	}

	// 字符级输入（支持光标移动，类似 ReadWithPlaceholder 但没有 placeholder）
	return input.ReadLineCoreDefault(promptText, validator, func(b []byte) string {
		return string(b)
	}, formatError)
}

// handleInputError 处理输入错误
// 清除用户输入行，在下一行显示红色错误提示
func handleInputError(err error) {
	// 回到用户输入行，清除输入
	fmt.Print(ansiMoveUp)         // 上移一行（回到用户输入行）
	fmt.Print(ansiCarriageReturn) // 回到行首
	fmt.Print(ansiClearToEnd)     // 清除到行尾（清除用户输入）

	// 在下一行显示红色错误提示
	errorMsg := formatError(err.Error())
	fmt.Print("\n") // 换行到新行
	fmt.Print(errorMsg)
	// 注意：不在这里打印额外的换行符，保持光标在错误行的末尾
	// 重置所有 ANSI 格式
	fmt.Print(ansiResetFormat)
}

// clearErrorAndInputLines 清除错误提示和输入行
// 当前光标在上一轮错误提示行的末尾
func clearErrorAndInputLines() {
	// 1. 清除错误提示行
	fmt.Print(ansiCarriageReturn) // 回到错误提示行行首
	fmt.Print(ansiClearToEnd)     // 清除到行尾（清除错误提示）
	// 2. 上移一行到上一轮输入行并清除
	fmt.Print(ansiMoveUp)         // 上移一行（到上一轮输入行）
	fmt.Print(ansiCarriageReturn) // 回到行首
	fmt.Print(ansiClearToEnd)     // 清除到行尾（清除上一轮输入）
}

// clearAndDisplayResult 清除输入区域并显示格式化结果
// 处理光标位置和清除逻辑：
// - 输入完成后，光标在输入行的下一行（ReadLineCore/ReadWithPlaceholder 都会打印换行）
// - 如果有错误提示，错误提示在输入行的下一行，当前光标在错误提示的下一行
// - 如果没有错误提示，当前光标在输入行的下一行
func clearAndDisplayResult(hadError bool, message, value string, isPassword bool, config common.PromptConfig) {
	if hadError {
		// 有错误提示的情况：
		// 当前光标在错误提示的下一行
		// 1. 清除当前行（空行）
		fmt.Print(ansiCarriageReturn) // 回到行首
		fmt.Print(ansiClearToEnd)     // 清除到行尾
		// 2. 上移一行到错误提示行
		fmt.Print(ansiMoveUp)
		fmt.Print(ansiCarriageReturn) // 回到行首
		fmt.Print(ansiClearToEnd)     // 清除错误提示行
		// 3. 上移一行到输入行
		fmt.Print(ansiMoveUp)
		fmt.Print(ansiCarriageReturn) // 回到行首
		fmt.Print(ansiClearToEnd)     // 清除输入行
		// 4. 上移一行到提示行
		fmt.Print(ansiMoveUp)
	} else {
		// 没有错误提示的情况：
		// 当前光标在输入行的下一行
		// 1. 清除当前行（空行）
		fmt.Print(ansiCarriageReturn) // 回到行首
		fmt.Print(ansiClearToEnd)     // 清除到行尾
		// 2. 上移一行到输入行
		fmt.Print(ansiMoveUp)
		fmt.Print(ansiCarriageReturn) // 回到行首
		fmt.Print(ansiClearToEnd)     // 清除输入行
		// 3. 上移一行到提示行
		fmt.Print(ansiMoveUp)
	}

	// 清除提示行（当前光标在提示行）
	fmt.Print(ansiCarriageReturn) // 回到行首
	fmt.Print(ansiClearToEnd)     // 清除到行尾

	// 格式化并显示结果
	displayFormattedResult(message, value, isPassword, config)
}

// displayFormattedResult 显示格式化的结果
// 已输入时显示："> [title] [value]"
func displayFormattedResult(message, value string, isPassword bool, config common.PromptConfig) {
	// 使用原始的 message，不包含默认值
	formattedPrompt := config.FormatPrompt(message)

	// 对于密码模式，显示固定长度的掩码（****）而不是实际值
	displayValue := value
	if isPassword {
		displayValue = PasswordMask
	}

	// 使用统一的配置来格式化答案
	formattedAnswer := displayValue
	if config.FormatAnswer != nil {
		formattedAnswer = config.FormatAnswer(displayValue)
	}

	// 如果提供了 FormatResultTitle，使用它来格式化 title
	displayTitle := formattedPrompt
	if config.FormatResultTitle != nil {
		displayTitle = config.FormatResultTitle(message, displayValue)
	}

	// 显示："> " + 提示文本 + " " + 输入值（已输入时使用 "> " 前缀）
	answerPrefix := DefaultInputPrefix
	if config.FormatAnswerPrefix != nil {
		answerPrefix = config.FormatAnswerPrefix()
	}

	fmt.Print(answerPrefix)
	fmt.Print(displayTitle)
	fmt.Print(" ")
	fmt.Print(formattedAnswer)
	fmt.Println() // 换行

	// 重置所有 ANSI 格式，确保后续输出格式正确
	fmt.Print(ansiResetFormat)
}

// buildInputConfig 构建输入配置的通用函数
// 处理配置合并和 ResultTitle 设置
func buildInputConfig(customConfig *common.PromptConfig, resultTitle string) common.PromptConfig {
	var baseConfig common.PromptConfig
	if customConfig != nil {
		baseConfig = *customConfig
	} else {
		baseConfig = newDefaultConfig()
	}
	return common.WithResultTitle(baseConfig, resultTitle)
}

// AskInput 使用配置结构体的输入函数
func AskInput(field InputField) (string, error) {
	config := buildInputConfig(field.Config, field.ResultTitle)
	return inputFuncWithConfig(field.Message, field.DefaultValue, "", false, field.Validator, &config)
}

// AskPassword 使用配置结构体的密码函数
func AskPassword(field PasswordField) (string, error) {
	config := buildInputConfig(field.Config, field.ResultTitle)
	return inputFuncWithConfig(field.Message, field.DefaultValue, "", true, field.Validator, &config)
}

// ValidateRegex 创建一个基于正则表达式的验证器（重新导出）
func ValidateRegex(pattern string, errorMsg string) Validator {
	return input.ValidateRegex(pattern, errorMsg)
}

// ValidateEmail 验证邮箱格式（重新导出）
func ValidateEmail() Validator {
	return input.ValidateEmail()
}

// ValidateURL 验证 URL 格式（重新导出）
func ValidateURL() Validator {
	return input.ValidateURL()
}

// ValidateRequired 验证输入不能为空（重新导出）
func ValidateRequired() Validator {
	return input.ValidateRequired()
}

// ValidateMinLength 验证最小长度（重新导出）
func ValidateMinLength(minLen int) Validator {
	return input.ValidateMinLength(minLen)
}

// ValidateMaxLength 验证最大长度（重新导出）
func ValidateMaxLength(maxLen int) Validator {
	return input.ValidateMaxLength(maxLen)
}

// ValidateLength 验证长度范围（重新导出）
func ValidateLength(minLen, maxLen int) Validator {
	return input.ValidateLength(minLen, maxLen)
}

// InputBuilder Input 的链式构建器
type InputBuilder struct {
	baseBuilder
	defaultValue string
	placeholder  string
	isPassword   bool
	validator    Validator
}

// Input 创建一个 Input 构建器（链式调用）
// 使用方式：prompt.Input().Prompt("消息").DefaultValue("默认值").Validate(validator).Run()
func Input() *InputBuilder {
	return &InputBuilder{
		isPassword: false,
	}
}

// Password 创建一个 Password 构建器（链式调用）
// 使用方式：prompt.Password().Prompt("消息").Validate(validator).Run()
func Password() *InputBuilder {
	return &InputBuilder{
		isPassword: true,
	}
}

// Prompt 设置提示消息（覆盖基类方法以返回正确的类型）
func (b *InputBuilder) Prompt(message string) *InputBuilder {
	b.baseBuilder.Prompt(message)
	return b
}

// DefaultValue 设置默认值
func (b *InputBuilder) DefaultValue(value string) *InputBuilder {
	b.defaultValue = value
	return b
}

// Placeholder 设置占位符文本（仅在非密码模式下有效）
// 占位符会在输入框为空时显示，用户开始输入时自动消失
func (b *InputBuilder) Placeholder(text string) *InputBuilder {
	b.placeholder = text
	return b
}

// Validate 设置验证器
func (b *InputBuilder) Validate(validator Validator) *InputBuilder {
	b.validator = validator
	return b
}

// Run 执行输入并返回结果
func (b *InputBuilder) Run() (string, error) {
	if b.isPassword {
		return inputFunc(b.GetMessage(), b.defaultValue, "", true, b.validator)
	}
	return inputFunc(b.GetMessage(), b.defaultValue, b.placeholder, false, b.validator)
}

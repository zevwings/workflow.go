package prompt

import (
	"fmt"
	"strings"

	"github.com/zevwings/workflow/internal/prompt/input"
)

// Validator 验证函数类型（重新导出，保持向后兼容）
type Validator = input.Validator

// input 统一的输入函数，通过 isPassword 参数控制是否使用密文模式
// placeholder 参数用于显示占位符文本（仅在非密码模式下有效）
func inputFunc(message string, defaultValue string, placeholder string, isPassword bool, validator Validator) (string, error) {
	// 格式化提示消息（类似 huh 的 Title，单独显示）
	promptMsg := formatPrompt(message)

	// 参考 huh + bubbletea 的样式：
	// 1. Title（提示信息）单独显示一行，并在其后显示默认值
	//    示例：请输入您的邮箱[user@example.com]
	// 2. 输入框使用 "> " 前缀（类似 huh 的样式）
	// 3. Placeholder 显示在输入框内（灰色斜体）

	// 显示 Title（提示信息），如果有默认值，则在 Title 后面显示
	if defaultValue != "" {
		// 普通输入：直接显示真实默认值
		if !isPassword {
			titleText := fmt.Sprintf("%s[%s]", promptMsg, defaultValue)
			fmt.Println(titleText)
		} else {
			// 密码输入：永远只显示固定长度的掩码，避免泄露真实长度
			masked := "****"
			titleText := fmt.Sprintf("%s[%s]", promptMsg, masked)
			fmt.Println(titleText)
		}
	} else {
		fmt.Println(promptMsg)
	}

	// 构建输入框前缀（类似 huh 的 "> " 样式）
	inputPrefix := "> "

	// 构建完整的提示文本（用于输入时显示）
	// 这里不再显示默认值，默认值仅显示在 Title 行
	promptText := inputPrefix

	// 获取主题配置
	t := GetTheme()

	// 构建编辑器配置
	config := input.Config{
		FormatPlaceholder: func(text string) string {
			return input.FormatPlaceholder(text, t.HintStyle, t.EnableColor)
		},
		FormatError: formatError,
		HintStyle:   t.HintStyle,
		ErrorStyle:  t.ErrorStyle,
		EnableColor: t.EnableColor,
	}

	hasError := false
	for {
		// 注意：对于有 placeholder 的情况，ReadWithPlaceholder 会自己打印 promptText
		// 对于没有 placeholder 的情况，我们需要在这里打印
		if hasError {
			// 重新开始新一轮输入前，清除上一轮的错误提示和输入行
			// 当前光标在上一轮错误提示行的末尾
			// 1. 清除错误提示行
			fmt.Print("\r")     // 回到错误提示行行首
			fmt.Print("\033[K") // 清除到行尾（清除错误提示）
			// 2. 上移一行到上一轮输入行并清除
			fmt.Print("\033[A") // 上移一行（到上一轮输入行）
			fmt.Print("\r")     // 回到行首
			fmt.Print("\033[K") // 清除到行尾（清除上一轮输入）
		}

		var value string
		var err error

		// 针对有默认值的场景，为底层输入函数构造一个"宽松版"验证器（仅非密码）：
		// - 实时验证 / 回车验证时：空字符串视为合法（不报错），让用户可以直接回车退出输入循环
		// - 真实的默认值验证放在后面统一处理
		effectiveValidator := validator
		if !isPassword && defaultValue != "" && validator != nil {
			effectiveValidator = func(v string) error {
				if strings.TrimSpace(v) == "" {
					return nil
				}
				return validator(v)
			}
		}

		if isPassword {
			// 密码模式：使用通用编辑内核，但通过 echo 函数以 * 方式显示输入内容
			value, err = input.ReadLineCore(promptText, validator, func(b []byte) string {
				return strings.Repeat("*", len(b))
			}, formatError)
		} else {
			// 普通输入模式：如果有 placeholder，使用字符级输入；否则也使用字符级输入（支持光标移动）
			if placeholder != "" {
				value, err = input.ReadWithPlaceholder(promptText, placeholder, effectiveValidator, config)
			} else {
				// 字符级输入（支持光标移动，类似 ReadWithPlaceholder 但没有 placeholder）
				value, err = input.ReadLineCore(promptText, effectiveValidator, func(b []byte) string {
					return string(b)
				}, formatError)
			}
		}

		if err != nil {
			return "", err
		}

		// 如果有默认值且用户未输入内容，直接使用默认值（仅非密码模式）
		if !isPassword && defaultValue != "" && strings.TrimSpace(value) == "" {
			value = defaultValue
		}

		// 如果提供了验证器，进行验证（密码模式已在 ReadLineCore 中进行实时验证和回车验证，这里无需重复）
		if validator != nil && !isPassword {
			if err := validator(value); err != nil {
				// 验证失败，清除用户输入行，在下一行显示红色错误提示
				if !isPassword {
					// 普通输入模式：回到用户输入行，清除输入
					fmt.Print("\033[A") // 上移一行（回到用户输入行）
					fmt.Print("\r")     // 回到行首
					fmt.Print("\033[K") // 清除到行尾（清除用户输入）
				}
				// 在下一行显示红色错误提示
				errorMsg := formatError(err.Error())
				fmt.Print("\n") // 换行到新行
				fmt.Print(errorMsg)
				// 注意：不在这里打印额外的换行符，保持光标在错误行的末尾
				// 重置所有 ANSI 格式
				fmt.Print("\033[0m")
				// 标记已显示错误，继续循环
				hasError = true
				continue
			}
		}

		// 能执行到这里说明本轮输入已通过验证，重置错误状态
		// 注意：我们需要在重置 hasError 之前保存是否有错误提示的状态
		hadError := hasError
		hasError = false

		// 验证通过，显示格式化的结果
		// 新布局：提示文本 + 输入值（同一行，无 > 前缀）

		// 处理光标位置和清除逻辑：
		// 输入完成后，光标在输入行的下一行（ReadLineCore/ReadWithPlaceholder 都会打印换行）
		// 如果有错误提示，错误提示在输入行的下一行，当前光标在错误提示的下一行
		// 如果没有错误提示，当前光标在输入行的下一行

		if hadError {
			// 有错误提示的情况：
			// 当前光标在错误提示的下一行
			// 1. 清除当前行（空行）
			fmt.Print("\r")     // 回到行首
			fmt.Print("\033[K") // 清除到行尾
			// 2. 上移一行到错误提示行
			fmt.Print("\033[A")
			fmt.Print("\r")     // 回到行首
			fmt.Print("\033[K") // 清除错误提示行
			// 3. 上移一行到输入行
			fmt.Print("\033[A")
			fmt.Print("\r")     // 回到行首
			fmt.Print("\033[K") // 清除输入行
			// 4. 上移一行到提示行
			fmt.Print("\033[A")
		} else {
			// 没有错误提示的情况：
			// 当前光标在输入行的下一行
			// 1. 清除当前行（空行）
			fmt.Print("\r")     // 回到行首
			fmt.Print("\033[K") // 清除到行尾
			// 2. 上移一行到输入行
			fmt.Print("\033[A")
			fmt.Print("\r")     // 回到行首
			fmt.Print("\033[K") // 清除输入行
			// 3. 上移一行到提示行
			fmt.Print("\033[A")
		}

		// 清除提示行（当前光标在提示行）
		fmt.Print("\r")     // 回到行首
		fmt.Print("\033[K") // 清除到行尾

		// 在同一行显示：formatPrompt(message) + " " + formatAnswer(displayValue)
		// 注意：使用原始的 message，不包含默认值
		formattedPrompt := formatPrompt(message)

		// 对于密码模式，显示星号而不是实际值
		var displayValue string
		if isPassword {
			displayValue = strings.Repeat("*", len(value))
		} else {
			displayValue = value
		}

		formattedAnswer := formatAnswer(displayValue)

		// 显示：提示文本 + 空格 + 输入值
		fmt.Print(formattedPrompt)
		fmt.Print(" ")
		fmt.Print(formattedAnswer)
		fmt.Println() // 换行

		// 重置所有 ANSI 格式，确保后续输出格式正确
		fmt.Print("\033[0m")

		return value, nil
	}
}

// AskInput 函数式调用 Input（保持向后兼容）
// 使用方式：prompt.AskInput("消息", "默认值", validator...)
func AskInput(message string, defaultValue string, validator ...Validator) (string, error) {
	var v Validator
	if len(validator) > 0 {
		v = validator[0]
	}
	return inputFunc(message, defaultValue, "", false, v)
}

// AskPassword 函数式调用 Password（保持向后兼容）
// 使用方式：prompt.AskPassword("消息", validator...)
func AskPassword(message string, validator ...Validator) (string, error) {
	var v Validator
	if len(validator) > 0 {
		v = validator[0]
	}
	return inputFunc(message, "", "", true, v)
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
	BaseBuilder
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
	b.BaseBuilder.Prompt(message)
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
		return inputFunc(b.GetMessage(), "", "", true, b.validator)
	}
	return inputFunc(b.GetMessage(), b.defaultValue, b.placeholder, false, b.validator)
}

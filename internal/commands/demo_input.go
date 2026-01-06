package commands

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/your-org/workflow/internal/output"
	"github.com/your-org/workflow/internal/prompt"
)

// NewDemoInputCmd 创建一个演示 Input 和 Password 功能的命令
func NewDemoInputCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "demo-input",
		Short: "演示 Input 和 Password 的交互功能",
		Long: `演示 Input 和 Password 的各种交互功能：
- Input: 文本输入（支持默认值、占位符、验证）
- Password: 密码输入（隐藏输入，支持验证）

这个 demo 会依次展示所有功能，帮助您了解 Input 和 Password 的用法。`,
		RunE: runDemoInput,
	}

	return cmd
}

func runDemoInput(cmd *cobra.Command, args []string) error {
	out := output.NewOutput(false)

	out.Info("欢迎使用 Input 和 Password 功能演示")
	out.Println("")
	out.Info("本演示将展示以下功能：")
	out.Println("  1. Input - 文本输入（无默认值，带验证）")
	out.Println("  2. Input - 文本输入（带占位符）")
	out.Println("  3. Input - 文本输入（有默认值，带验证）")
	out.Println("  4. Password - 密码输入（带验证）")
	out.Println("  5. Input 验证 - 邮箱格式")
	out.Println("  6. Input 验证 - URL 格式")
	out.Println("  7. Input 验证 - 正则表达式")
	out.Println("  8. Password 验证 - 长度验证")
	out.Println("  9. Input 验证 - 必填项")
	out.Println("")

	// 1. 演示 Input（无默认值，带验证）
	out.Info("=== 演示 1: Input（无默认值，带验证）===")
	name, err := prompt.Input().
		Prompt("请输入您的姓名").
		Validate(prompt.ValidateRequired()).
		Run()
	if err != nil {
		return fmt.Errorf("输入失败: %w", err)
	}
	out.Success("您输入的姓名是: %s", name)
	out.Println("")

	// 2. 演示 Input（带 Placeholder）
	out.Info("=== 演示 2: Input（带 Placeholder）===")
	out.Println("提示：注意输入框中的灰色占位符文本")
	username, err := prompt.Input().
		Prompt("请输入用户名").
		Placeholder("例如：admin").
		Validate(prompt.ValidateRequired()).
		Run()
	if err != nil {
		return fmt.Errorf("输入失败: %w", err)
	}
	out.Success("您输入的用户名是: %s", username)
	out.Println("")

	// 3. 演示 Input（有默认值，带验证）
	out.Info("=== 演示 3: Input（有默认值，带验证）===")
	email, err := prompt.Input().
		Prompt("请输入您的邮箱").
		DefaultValue("user@example.com").
		Validate(prompt.ValidateEmail()).
		Run()
	if err != nil {
		return fmt.Errorf("输入失败: %w", err)
	}
	out.Success("您输入的邮箱是: %s", email)
	out.Println("")

	// 4. 演示 Password（带验证）
	out.Info("=== 演示 4: Password（密码输入，带验证）===")
	out.Println("注意：密码输入时不会显示明文")
	token, err := prompt.Password().
		Prompt("请输入一个测试 Token（不会显示）").
		Validate(prompt.ValidateMinLength(8)).
		Run()
	if err != nil {
		return fmt.Errorf("输入失败: %w", err)
	}
	// 只显示前 4 个字符，其余用 * 替代
	maskedToken := maskToken(token)
	out.Success("Token 已输入（长度: %d 字符，显示: %s）", len(token), maskedToken)
	out.Println("")

	// 5. 演示验证功能（邮箱验证）
	out.Info("=== 演示 5: Input 验证（邮箱格式）===")
	out.Println("提示：请输入一个有效的邮箱地址，如果格式不正确会提示重试")
	emailValidated, err := prompt.Input().
		Prompt("请输入邮箱地址").
		Validate(prompt.ValidateEmail()).
		Run()
	if err != nil {
		return fmt.Errorf("输入失败: %w", err)
	}
	out.Success("邮箱验证通过: %s", emailValidated)
	out.Println("")

	// 6. 演示验证功能（URL 验证）
	out.Info("=== 演示 6: Input 验证（URL 格式）===")
	out.Println("提示：请输入一个有效的 URL 地址")
	urlValidated, err := prompt.Input().
		Prompt("请输入 URL 地址").
		DefaultValue("https://example.com").
		Validate(prompt.ValidateURL()).
		Run()
	if err != nil {
		return fmt.Errorf("输入失败: %w", err)
	}
	out.Success("URL 验证通过: %s", urlValidated)
	out.Println("")

	// 7. 演示验证功能（正则验证）
	out.Info("=== 演示 7: Input 验证（正则表达式）===")
	out.Println("提示：用户名只能包含字母、数字和下划线，长度 3-20 个字符")
	usernameValidated, err := prompt.Input().
		Prompt("请输入用户名").
		Validate(prompt.ValidateRegex(`^[a-zA-Z0-9_]{3,20}$`, "用户名只能包含字母、数字和下划线，长度 3-20 个字符")).
		Run()
	if err != nil {
		return fmt.Errorf("输入失败: %w", err)
	}
	out.Success("用户名验证通过: %s", usernameValidated)
	out.Println("")

	// 8. 演示验证功能（长度验证）
	out.Info("=== 演示 8: Password 验证（长度验证）===")
	out.Println("提示：密码长度必须在 8 到 32 个字符之间")
	passwordValidated, err := prompt.Password().
		Prompt("请输入密码（8-32 个字符）").
		Validate(prompt.ValidateLength(8, 32)).
		Run()
	if err != nil {
		return fmt.Errorf("输入失败: %w", err)
	}
	maskedPassword := maskToken(passwordValidated)
	out.Success("密码验证通过（长度: %d 字符，显示: %s）", len(passwordValidated), maskedPassword)
	out.Println("")

	// 9. 演示验证功能（必填验证）
	out.Info("=== 演示 9: Input 验证（必填项）===")
	out.Println("提示：此项为必填项，不能为空")
	requiredInput, err := prompt.Input().
		Prompt("请输入必填项（不能为空）").
		Validate(prompt.ValidateRequired()).
		Run()
	if err != nil {
		return fmt.Errorf("输入失败: %w", err)
	}
	out.Success("必填项验证通过: %s", requiredInput)
	out.Println("")

	out.Success("演示完成！感谢使用 Input 和 Password 功能。")

	return nil
}

// maskToken 掩码显示 Token（只显示前 4 个字符）
func maskToken(token string) string {
	if len(token) <= 4 {
		return strings.Repeat("*", len(token))
	}
	return token[:4] + strings.Repeat("*", len(token)-4)
}

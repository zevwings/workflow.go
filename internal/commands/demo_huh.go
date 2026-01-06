package commands

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/huh"
	"github.com/spf13/cobra"
	"github.com/your-org/workflow/internal/output"
)

// NewDemoHuhCmd 创建一个演示 huh 库功能的命令
func NewDemoHuhCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "demo-huh",
		Short: "演示 huh 库的交互功能",
		Long: `演示 charmbracelet/huh 库的各种交互功能：
- Input: 文本输入（支持 Placeholder 和验证）
- Text: 多行文本输入
- Select: 单选（键盘导航）
- MultiSelect: 多选（空格切换）
- Confirm: 是/否确认
- Form: 表单组合

huh 是基于 bubbletea 的高级表单库，提供流畅的链式 API 和自动验证功能。`,
		RunE: runDemoHuh,
	}

	return cmd
}

func runDemoHuh(cmd *cobra.Command, args []string) error {
	out := output.NewOutput(false)

	out.Info("欢迎使用 huh 库功能演示")
	out.Println("")
	out.Info("本演示将展示以下功能：")
	out.Println("  1. Input - 文本输入（支持 Placeholder 和验证）")
	out.Println("  2. Text - 多行文本输入")
	out.Println("  3. Select - 单选（键盘导航）")
	out.Println("  4. MultiSelect - 多选（空格切换）")
	out.Println("  5. Confirm - 是/否确认")
	out.Println("  6. Form - 表单组合（多个字段一起填写）")
	out.Println("")

	// 1. 演示 Input（基本用法）
	out.Info("=== 演示 1: Input（基本用法 + Placeholder）===")
	var name string
	input1 := huh.NewInput().
		Title("请输入您的姓名").
		Placeholder("例如：张三").
		Value(&name).
		Validate(huh.ValidateNotEmpty())

	if err := input1.Run(); err != nil {
		return fmt.Errorf("输入失败: %w", err)
	}
	out.Success("您输入的姓名是: %s", name)
	out.Println("")

	// 2. 演示 Input（带邮箱验证）
	out.Info("=== 演示 2: Input（邮箱验证）===")
	var email string
	input2 := huh.NewInput().
		Title("请输入您的邮箱").
		Placeholder("user@example.com").
		Value(&email).
		Validate(huh.ValidateNotEmpty()).
		Validate(func(s string) error {
			if !strings.Contains(s, "@") || !strings.Contains(s, ".") {
				return fmt.Errorf("请输入有效的邮箱地址")
			}
			return nil
		})

	if err := input2.Run(); err != nil {
		return fmt.Errorf("输入失败: %w", err)
	}
	out.Success("您输入的邮箱是: %s", email)
	out.Println("")

	// 3. 演示 Input（带正则验证）
	out.Info("=== 演示 3: Input（正则表达式验证）===")
	var username string
	input3 := huh.NewInput().
		Title("请输入用户名").
		Placeholder("只能包含字母、数字和下划线，长度 3-20 个字符").
		Value(&username).
		Validate(huh.ValidateNotEmpty()).
		Validate(huh.ValidateMinLength(3)).
		Validate(huh.ValidateMaxLength(20)).
		Validate(func(s string) error {
			for _, r := range s {
				if !((r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') || r == '_') {
					return fmt.Errorf("用户名只能包含字母、数字和下划线")
				}
			}
			return nil
		})

	if err := input3.Run(); err != nil {
		return fmt.Errorf("输入失败: %w", err)
	}
	out.Success("用户名验证通过: %s", username)
	out.Println("")

	// 4. 演示 Text（多行文本输入）
	out.Info("=== 演示 4: Text（多行文本输入）===")
	var description string
	text1 := huh.NewText().
		Title("请输入项目描述").
		Placeholder("请详细描述您的项目...").
		CharLimit(200).
		Value(&description).
		Validate(huh.ValidateMinLength(10))

	if err := text1.Run(); err != nil {
		return fmt.Errorf("输入失败: %w", err)
	}
	out.Success("项目描述: %s", description)
	out.Println("")

	// 5. 演示 Select（单选）
	out.Info("=== 演示 5: Select（单选 - 键盘导航）===")
	var environment string
	select1 := huh.NewSelect[string]().
		Title("请选择部署环境").
		Options(
			huh.NewOption("开发环境 (dev)", "dev"),
			huh.NewOption("测试环境 (test)", "test"),
			huh.NewOption("预发布环境 (staging)", "staging"),
			huh.NewOption("生产环境 (production)", "production"),
		).
		Value(&environment)

	if err := select1.Run(); err != nil {
		return fmt.Errorf("选择失败: %w", err)
	}
	out.Success("您选择的环境是: %s", environment)
	out.Println("")

	// 6. 演示 Select（带默认值）
	out.Info("=== 演示 6: Select（带默认值）===")
	var language string
	select2 := huh.NewSelect[string]().
		Title("请选择您最喜欢的编程语言").
		Options(
			huh.NewOption("Go", "go"),
			huh.NewOption("Python", "python"),
			huh.NewOption("JavaScript", "javascript"),
			huh.NewOption("Rust", "rust"),
			huh.NewOption("Java", "java"),
			huh.NewOption("C++", "cpp"),
		).
		Value(&language).
		Description("使用上下箭头键导航，回车确认")

	if err := select2.Run(); err != nil {
		return fmt.Errorf("选择失败: %w", err)
	}
	out.Success("您选择的语言是: %s", language)
	out.Println("")

	// 7. 演示 MultiSelect（多选）
	out.Info("=== 演示 7: MultiSelect（多选 - 空格切换）===")
	var tools []string
	multiselect1 := huh.NewMultiSelect[string]().
		Title("请选择您使用的 DevOps 工具（可多选）").
		Options(
			huh.NewOption("Git", "git"),
			huh.NewOption("Docker", "docker"),
			huh.NewOption("Kubernetes", "kubernetes"),
			huh.NewOption("Terraform", "terraform"),
			huh.NewOption("Ansible", "ansible"),
			huh.NewOption("Jenkins", "jenkins"),
		).
		Value(&tools).
		Description("使用空格键切换选择，回车确认")

	if err := multiselect1.Run(); err != nil {
		return fmt.Errorf("选择失败: %w", err)
	}
	if len(tools) == 0 {
		out.Warning("您没有选择任何工具")
	} else {
		out.Success("您选择的工具: %s", strings.Join(tools, ", "))
	}
	out.Println("")

	// 8. 演示 MultiSelect（带默认值）
	out.Info("=== 演示 8: MultiSelect（带默认值）===")
	var features []string
	multiselect2 := huh.NewMultiSelect[string]().
		Title("请选择要启用的功能（可多选）").
		Options(
			huh.NewOption("用户认证", "auth"),
			huh.NewOption("数据加密", "encryption"),
			huh.NewOption("日志记录", "logging"),
			huh.NewOption("监控告警", "monitoring"),
			huh.NewOption("备份恢复", "backup"),
			huh.NewOption("负载均衡", "loadbalancer"),
		).
		Value(&features).
		Description("使用空格键切换选择，回车确认")

	if err := multiselect2.Run(); err != nil {
		return fmt.Errorf("选择失败: %w", err)
	}
	if len(features) == 0 {
		out.Warning("您没有选择任何功能")
	} else {
		out.Success("启用的功能: %s", strings.Join(features, ", "))
	}
	out.Println("")

	// 9. 演示 Confirm（确认）
	out.Info("=== 演示 9: Confirm（是/否确认）===")
	var confirm1 bool
	confirm1Field := huh.NewConfirm().
		Title("是否继续演示？").
		Affirmative("是").
		Negative("否").
		Value(&confirm1)

	if err := confirm1Field.Run(); err != nil {
		return fmt.Errorf("确认失败: %w", err)
	}
	if confirm1 {
		out.Success("您选择了: 是")
	} else {
		out.Warning("您选择了: 否")
	}
	out.Println("")

	// 10. 演示 Confirm（带默认值）
	out.Info("=== 演示 10: Confirm（带默认值）===")
	var confirm2 bool
	confirm2Field := huh.NewConfirm().
		Title("是否退出演示？（默认: 否）").
		Affirmative("退出").
		Negative("继续").
		Value(&confirm2)

	if err := confirm2Field.Run(); err != nil {
		return fmt.Errorf("确认失败: %w", err)
	}
	if confirm2 {
		out.Warning("您选择了: 退出（演示结束）")
		return nil
	} else {
		out.Success("您选择了: 继续演示")
	}
	out.Println("")

	// 11. 演示 Form（表单组合 - 综合场景）
	out.Info("=== 演示 11: Form（表单组合 - 综合场景）===")
	out.Println("提示：这是一个完整的表单，包含多个字段，可以一次性填写完成")
	out.Println("")

	// 表单数据
	var (
		formName      string
		formEmail     string
		formRole      string
		formSkills    []string
		formBio       string
		formSubscribe bool
		formAgree     bool
	)

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("姓名").
				Placeholder("请输入您的姓名").
				Value(&formName).
				Validate(huh.ValidateNotEmpty()),
			huh.NewInput().
				Title("邮箱").
				Placeholder("user@example.com").
				Value(&formEmail).
				Validate(huh.ValidateNotEmpty()).
				Validate(func(s string) error {
					if !strings.Contains(s, "@") || !strings.Contains(s, ".") {
						return fmt.Errorf("请输入有效的邮箱地址")
					}
					return nil
				}),
			huh.NewSelect[string]().
				Title("角色").
				Options(
					huh.NewOption("开发工程师", "developer"),
					huh.NewOption("运维工程师", "ops"),
					huh.NewOption("产品经理", "pm"),
					huh.NewOption("设计师", "designer"),
				).
				Value(&formRole),
		),
		huh.NewGroup(
			huh.NewMultiSelect[string]().
				Title("技能（可多选）").
				Options(
					huh.NewOption("Go", "go"),
					huh.NewOption("Python", "python"),
					huh.NewOption("JavaScript", "javascript"),
					huh.NewOption("Docker", "docker"),
					huh.NewOption("Kubernetes", "k8s"),
				).
				Value(&formSkills).
				Description("使用空格键选择"),
			huh.NewText().
				Title("个人简介").
				Placeholder("请简单介绍一下自己...").
				CharLimit(150).
				Value(&formBio).
				Validate(huh.ValidateMinLength(10)),
		),
		huh.NewGroup(
			huh.NewConfirm().
				Title("是否订阅邮件通知？").
				Affirmative("订阅").
				Negative("不订阅").
				Value(&formSubscribe),
			huh.NewConfirm().
				Title("我已阅读并同意服务条款").
				Affirmative("同意").
				Negative("不同意").
				Value(&formAgree).
				Validate(func(b bool) error {
					if !b {
						return fmt.Errorf("必须同意服务条款才能继续")
					}
					return nil
				}),
		),
	).
		WithTheme(huh.ThemeBase16()).
		WithWidth(80)

	if err := form.Run(); err != nil {
		return fmt.Errorf("表单填写失败: %w", err)
	}

	// 显示表单摘要
	out.Println("")
	out.Info("=== 表单提交摘要 ===")
	out.Success("姓名: %s", formName)
	out.Success("邮箱: %s", formEmail)
	out.Success("角色: %s", formRole)
	if len(formSkills) > 0 {
		out.Success("技能: %s", strings.Join(formSkills, ", "))
	} else {
		out.Warning("技能: 未选择")
	}
	if formBio != "" {
		out.Success("个人简介: %s", formBio)
	} else {
		out.Warning("个人简介: 未填写")
	}
	if formSubscribe {
		out.Success("邮件通知: 已订阅")
	} else {
		out.Info("邮件通知: 未订阅")
	}
	if formAgree {
		out.Success("服务条款: 已同意")
	} else {
		out.Warning("服务条款: 未同意")
	}
	out.Println("")

	// 最终确认
	var finalConfirm bool
	finalConfirmField := huh.NewConfirm().
		Title("表单信息确认无误？（仅演示，不会真正保存）").
		Affirmative("确认").
		Negative("取消").
		Value(&finalConfirm)

	if err := finalConfirmField.Run(); err != nil {
		return fmt.Errorf("确认失败: %w", err)
	}
	if finalConfirm {
		out.Success("✓ 表单已提交（模拟）")
	} else {
		out.Info("表单提交已取消")
	}

	out.Println("")
	out.Success("演示完成！感谢使用 huh 库。")
	out.Println("")
	out.Info("💡 提示：")
	out.Println("  - huh 提供了流畅的链式 API 和自动验证")
	out.Println("  - 支持多种输入类型：Input、Text、Select、MultiSelect、Confirm")
	out.Println("  - 可以组合多个字段形成 Form，一次性填写")
	out.Println("  - 自动处理键盘导航和界面更新")

	return nil
}

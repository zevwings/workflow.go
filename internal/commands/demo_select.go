package commands

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/zevwings/workflow/internal/prompt"
)

// NewDemoSelectCmd 创建一个演示 Select 功能的命令
func NewDemoSelectCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "demo-select",
		Short: "演示 Select 的交互功能",
		Long: `演示 Select 的各种交互功能：
- Select: 单选（支持上下箭头导航，回车确认）
- 函数式 API: AskSelect()
- 链式 Builder API: Select()

这个 demo 会依次展示所有功能，帮助您了解 Select 的用法。`,
		RunE: runDemoSelect,
	}

	return cmd
}

func runDemoSelect(cmd *cobra.Command, args []string) error {
	out := prompt.NewMessage(false)

	out.Info("欢迎使用 Select 功能演示")
	out.Println("")
	out.Info("本演示将展示以下功能：")
	out.Println("  1. Select（函数式 API）")
	out.Println("  2. Select（链式 Builder API）")
	out.Println("  3. Select（带默认选项）")
	out.Println("  4. Select（不同场景示例）")
	out.Println("")
	out.Println("提示：使用 ↑/↓ 箭头键导航，回车键确认选择")
	out.Println("")

	// 1. 演示 Select（函数式 API）
	out.Info("=== 演示 1: Select（函数式 API）===")
	options := []string{"选项 A - 开发环境", "选项 B - 测试环境", "选项 C - 生产环境", "选项 D - 预发布环境"}
	selectedIndex, err := prompt.AskSelect("请选择部署环境", options, 0)
	if err != nil {
		return fmt.Errorf("选择失败: %w", err)
	}
	out.Success("您选择了: %s (索引: %d)", options[selectedIndex], selectedIndex)
	out.Println("")

	// 2. 演示 Select（链式 Builder API）
	out.Info("=== 演示 2: Select（链式 Builder API）===")
	programmingLanguages := []string{"Go", "Python", "JavaScript", "Rust", "Java", "C++"}
	langIndex, err := prompt.Select().
		Prompt("请选择您最喜欢的编程语言").
		Options(programmingLanguages).
		Default(0).
		Run()
	if err != nil {
		return fmt.Errorf("选择失败: %w", err)
	}
	out.Success("您选择了: %s (索引: %d)", programmingLanguages[langIndex], langIndex)
	out.Println("")

	// 3. 演示 Select（带默认选项）
	out.Info("=== 演示 3: Select（带默认选项）===")
	envOptions := []string{"开发环境", "测试环境", "预发布环境", "生产环境"}
	envIndex, err := prompt.Select().
		Prompt("请选择部署环境（默认: 开发环境）").
		Options(envOptions).
		Default(0).
		Run()
	if err != nil {
		return fmt.Errorf("选择失败: %w", err)
	}
	out.Success("您选择了: %s (索引: %d)", envOptions[envIndex], envIndex)
	out.Println("")

	// 4. 演示 Select（不同场景示例）
	out.Info("=== 演示 4: Select（不同场景示例）===")
	frameworkOptions := []string{
		"React - 用于构建用户界面的 JavaScript 库",
		"Vue - 渐进式 JavaScript 框架",
		"Angular - 用于构建移动和桌面 Web 应用的平台",
		"Svelte - 编译时优化的 UI 框架",
		"Next.js - React 生产框架",
	}
	frameworkIndex, err := prompt.Select().
		Prompt("请选择前端框架").
		Options(frameworkOptions).
		Default(0).
		Run()
	if err != nil {
		return fmt.Errorf("选择失败: %w", err)
	}
	out.Success("您选择了: %s (索引: %d)", frameworkOptions[frameworkIndex], frameworkIndex)
	out.Println("")

	out.Success("演示完成！感谢使用 Select 功能。")

	return nil
}

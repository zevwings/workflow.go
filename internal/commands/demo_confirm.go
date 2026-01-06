package commands

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/your-org/workflow/internal/prompt"
	"github.com/your-org/workflow/internal/output"
)

// NewDemoConfirmCmd 创建一个演示 Confirm 功能的命令
func NewDemoConfirmCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "demo-confirm",
		Short: "演示 Confirm 的交互功能",
		Long: `演示 Confirm 的各种交互功能：
- Confirm: 是/否确认（支持默认值）
- 函数式 API: AskConfirm()
- 链式 Builder API: Confirm()

这个 demo 会依次展示所有功能，帮助您了解 Confirm 的用法。`,
		RunE: runDemoConfirm,
	}

	return cmd
}

func runDemoConfirm(cmd *cobra.Command, args []string) error {
	out := output.NewOutput(false)

	out.Info("欢迎使用 Confirm 功能演示")
	out.Println("")
	out.Info("本演示将展示以下功能：")
	out.Println("  1. Confirm（默认 Yes）- 函数式 API")
	out.Println("  2. Confirm（默认 No）- 函数式 API")
	out.Println("  3. Confirm（测试非法输入）")
	out.Println("  4. Confirm（链式 Builder API）")
	out.Println("")

	// 1. 演示 Confirm（默认 Yes）
	out.Info("=== 演示 1: Confirm（默认 Yes）- 函数式 API ===")
	confirm1, err := prompt.AskConfirm("是否继续演示？（默认: Yes）", true)
	if err != nil {
		return fmt.Errorf("确认失败: %w", err)
	}
	if confirm1 {
		out.Success("您选择了: 是")
	} else {
		out.Warning("您选择了: 否")
	}
	out.Println("")

	// 2. 演示 Confirm（默认 No）
	out.Info("=== 演示 2: Confirm（默认 No）- 函数式 API ===")
	confirm2, err := prompt.AskConfirm("是否退出演示？（默认: No）", false)
	if err != nil {
		return fmt.Errorf("确认失败: %w", err)
	}
	if confirm2 {
		out.Warning("您选择了: 是（将退出）")
		return nil
	} else {
		out.Success("您选择了: 否（继续演示）")
	}
	out.Println("")

	// 3. 演示非法输入处理（Confirm）
	out.Info("=== 演示 3: Confirm（测试非法输入）===")
	out.Println("提示：您可以输入 'xyz' 等非法值，系统会提示重试")
	confirm3, err := prompt.AskConfirm("请输入 y 或 n（或直接回车使用默认值）", true)
	if err != nil {
		return fmt.Errorf("确认失败: %w", err)
	}
	if confirm3 {
		out.Success("最终选择: 是")
	} else {
		out.Success("最终选择: 否")
	}
	out.Println("")

	// 4. 演示 Confirm（链式 Builder API）
	out.Info("=== 演示 4: Confirm（链式 Builder API）===")
	confirm4, err := prompt.Confirm().
		Prompt("是否保存配置？").
		Default(true).
		Run()
	if err != nil {
		return fmt.Errorf("确认失败: %w", err)
	}
	if confirm4 {
		out.Success("您选择了: 是")
	} else {
		out.Warning("您选择了: 否")
	}
	out.Println("")

	// 5. 演示 Confirm（链式 Builder API，默认 No）
	out.Info("=== 演示 5: Confirm（链式 Builder API，默认 No）===")
	confirm5, err := prompt.Confirm().
		Prompt("是否删除数据？（此操作不可恢复）").
		Default(false).
		Run()
	if err != nil {
		return fmt.Errorf("确认失败: %w", err)
	}
	if confirm5 {
		out.Warning("您选择了: 是（将删除数据）")
	} else {
		out.Success("您选择了: 否（已取消）")
	}
	out.Println("")

	out.Success("演示完成！感谢使用 Confirm 功能。")

	return nil
}

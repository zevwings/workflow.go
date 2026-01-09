//go:build example

package example

import (
	"fmt"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"github.com/zevwings/workflow/internal/prompt"
)

// NewDemoCmd 创建一个演示所有 Prompt 组件功能的命令
func NewDemoCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "demo-prompt",
		Short: "演示所有 Prompt 组件的交互功能",
		Long: `演示所有 Prompt 组件的功能：
- Message: 消息输出（Info、Success、Warning、Error）
- Confirm: 是/否确认
- Input: 文本输入
- Password: 密码输入
- Select: 单选
- MultiSelect: 多选
- Spinner: 加载指示器
- Table: 表格显示

这个 demo 会依次展示所有组件的基本用法。`,
		RunE: runDemo,
	}

	return cmd
}

func runDemo(cmd *cobra.Command, args []string) error {
	msg := prompt.NewMessage(false)

	msg.Info("欢迎使用 Prompt 组件演示")
	msg.Println("")
	msg.Info("本演示将依次展示以下组件：")
	msg.Println("  1. Message - 消息输出")
	msg.Println("  2. Confirm - 确认对话框")
	msg.Println("  3. Input - 文本输入")
	msg.Println("  4. Password - 密码输入")
	msg.Println("  5. Select - 单选")
	msg.Println("  6. MultiSelect - 多选")
	msg.Println("  7. Spinner - 加载指示器")
	msg.Println("  8. Table - 表格显示")
	msg.Println("")

	// 1. Message 演示
	msg.Info("=== 1. Message 消息输出 ===")
	msg.Success("成功消息示例")
	msg.Warning("警告消息示例")
	msg.Error("错误消息示例")
	msg.Println("")

	// 2. Confirm 演示
	msg.Info("=== 2. Confirm 确认对话框 ===")
	confirm, err := prompt.AskConfirm("是否继续演示？", true)
	if err != nil {
		return fmt.Errorf("确认失败: %w", err)
	}
	if confirm {
		msg.Success("您选择了: 是")
	} else {
		msg.Warning("您选择了: 否")
	}
	msg.Println("")

	// 3. Input 演示
	msg.Info("=== 3. Input 文本输入 ===")
	name, err := prompt.Input().
		Prompt("请输入您的姓名").
		DefaultValue("张三").
		Validate(prompt.ValidateRequired()).
		Run()
	if err != nil {
		return fmt.Errorf("输入失败: %w", err)
	}
	msg.Success("您输入的姓名是: %s", name)
	msg.Println("")

	// 4. Password 演示
	msg.Info("=== 4. Password 密码输入 ===")
	password, err := prompt.Password().
		Prompt("请输入密码（至少 6 个字符）").
		Validate(func(s string) error {
			if len(s) < 6 {
				return fmt.Errorf("长度至少需要 6 个字符")
			}
			return nil
		}).
		Run()
	if err != nil {
		return fmt.Errorf("输入失败: %w", err)
	}
	maskedPassword := maskPassword(password)
	msg.Success("密码已输入（长度: %d 字符，显示: %s）", len(password), maskedPassword)
	msg.Println("")

	// 5. Select 演示
	msg.Info("=== 5. Select 单选 ===")
	msg.Println("提示：使用 ↑/↓ 箭头键导航，回车键确认")
	options := []string{"选项 A", "选项 B", "选项 C", "选项 D"}
	selectedIndex, err := prompt.Select().
		Prompt("请选择一个选项").
		Options(options).
		Default(0).
		Run()
	if err != nil {
		return fmt.Errorf("选择失败: %w", err)
	}
	msg.Success("您选择了: %s (索引: %d)", options[selectedIndex], selectedIndex)
	msg.Println("")

	// 6. MultiSelect 演示
	msg.Info("=== 6. MultiSelect 多选 ===")
	msg.Println("提示：使用 ↑/↓ 箭头键导航，空格键切换选择，回车键确认")
	features := []string{"功能 A", "功能 B", "功能 C", "功能 D"}
	selectedIndices, err := prompt.MultiSelect().
		Prompt("请选择要启用的功能（可多选）").
		Options(features).
		Default([]int{0}).
		Run()
	if err != nil {
		return fmt.Errorf("选择失败: %w", err)
	}
	if len(selectedIndices) == 0 {
		msg.Warning("您没有选择任何功能")
	} else {
		var selectedNames []string
		for _, idx := range selectedIndices {
			selectedNames = append(selectedNames, features[idx])
		}
		msg.Success("您选择了: %s (索引: %v)", strings.Join(selectedNames, ", "), selectedIndices)
	}
	msg.Println("")

	// 7. Spinner 演示
	msg.Info("=== 7. Spinner 加载指示器 ===")
	spinner := prompt.NewSpinner("正在处理中...")
	spinner.Start()
	time.Sleep(2 * time.Second)
	spinner.WithSuccess("处理完成")
	msg.Println("")

	// 8. Table 演示
	msg.Info("=== 8. Table 表格显示 ===")
	table := prompt.NewTable([]string{"组件", "状态", "说明"})
	table.AddRow([]string{"Message", "✓", "消息输出功能正常"})
	table.AddRow([]string{"Confirm", "✓", "确认对话框功能正常"})
	table.AddRow([]string{"Input", "✓", "文本输入功能正常"})
	table.AddRow([]string{"Password", "✓", "密码输入功能正常"})
	table.AddRow([]string{"Select", "✓", "单选功能正常"})
	table.AddRow([]string{"MultiSelect", "✓", "多选功能正常"})
	table.AddRow([]string{"Spinner", "✓", "加载指示器功能正常"})
	table.AddRow([]string{"Table", "✓", "表格显示功能正常"})
	table.Render()
	msg.Println("")

	msg.Success("演示完成！所有组件功能正常。")

	return nil
}

// maskPassword 掩码显示密码（只显示前 2 个字符）
func maskPassword(password string) string {
	if len(password) <= 2 {
		return strings.Repeat("*", len(password))
	}
	return password[:2] + strings.Repeat("*", len(password)-2)
}

package commands

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/your-org/workflow/internal/prompt"
	"github.com/your-org/workflow/internal/output"
)

// NewDemoMultiSelectCmd 创建一个演示 MultiSelect 功能的命令
func NewDemoMultiSelectCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "demo-multiselect",
		Short: "演示 MultiSelect 的交互功能",
		Long: `演示 MultiSelect 的各种交互功能：
- MultiSelect: 多选（支持空格切换选择，回车确认）
- 函数式 API: AskMultiSelect()
- 链式 Builder API: MultiSelect()

这个 demo 会依次展示所有功能，帮助您了解 MultiSelect 的用法。`,
		RunE: runDemoMultiSelect,
	}

	return cmd
}

func runDemoMultiSelect(cmd *cobra.Command, args []string) error {
	out := output.NewOutput(false)

	out.Info("欢迎使用 MultiSelect 功能演示")
	out.Println("")
	out.Info("本演示将展示以下功能：")
	out.Println("  1. MultiSelect（函数式 API）")
	out.Println("  2. MultiSelect（链式 Builder API）")
	out.Println("  3. MultiSelect（带默认选择）")
	out.Println("  4. MultiSelect（不同场景示例）")
	out.Println("")
	out.Println("提示：使用 ↑/↓ 箭头键导航，空格键切换选择，回车键确认")
	out.Println("")

	// 1. 演示 MultiSelect（函数式 API）
	out.Info("=== 演示 1: MultiSelect（函数式 API）===")
	tools := []string{"Git", "Docker", "Kubernetes", "Terraform", "Ansible", "Jenkins"}
	selectedTools, err := prompt.AskMultiSelect("请选择您使用的 DevOps 工具（可多选）", tools, []int{0, 1})
	if err != nil {
		return fmt.Errorf("选择失败: %w", err)
	}
	if len(selectedTools) == 0 {
		out.Warning("您没有选择任何工具")
	} else {
		var selectedNames []string
		for _, idx := range selectedTools {
			selectedNames = append(selectedNames, tools[idx])
		}
		out.Success("您选择了: %s (索引: %v)", strings.Join(selectedNames, ", "), selectedTools)
	}
	out.Println("")

	// 2. 演示 MultiSelect（链式 Builder API）
	out.Info("=== 演示 2: MultiSelect（链式 Builder API）===")
	features := []string{"用户认证", "数据加密", "日志记录", "监控告警", "备份恢复", "负载均衡"}
	selectedFeatures, err := prompt.MultiSelect().
		Prompt("请选择要启用的功能（可多选）").
		Options(features).
		Default([]int{0, 2}).
		Run()
	if err != nil {
		return fmt.Errorf("选择失败: %w", err)
	}
	if len(selectedFeatures) == 0 {
		out.Warning("您没有选择任何功能")
	} else {
		var selectedNames []string
		for _, idx := range selectedFeatures {
			selectedNames = append(selectedNames, features[idx])
		}
		out.Success("您选择了: %s (索引: %v)", strings.Join(selectedNames, ", "), selectedFeatures)
	}
	out.Println("")

	// 3. 演示 MultiSelect（带默认选择）
	out.Info("=== 演示 3: MultiSelect（带默认选择）===")
	plugins := []string{"插件 A", "插件 B", "插件 C", "插件 D", "插件 E"}
	selectedPlugins, err := prompt.MultiSelect().
		Prompt("请选择要安装的插件（可多选，默认已选择前两个）").
		Options(plugins).
		Default([]int{0, 1}).
		Run()
	if err != nil {
		return fmt.Errorf("选择失败: %w", err)
	}
	if len(selectedPlugins) == 0 {
		out.Warning("您没有选择任何插件")
	} else {
		var selectedNames []string
		for _, idx := range selectedPlugins {
			selectedNames = append(selectedNames, plugins[idx])
		}
		out.Success("您选择了: %s (索引: %v)", strings.Join(selectedNames, ", "), selectedPlugins)
	}
	out.Println("")

	// 4. 演示 MultiSelect（不同场景示例）
	out.Info("=== 演示 4: MultiSelect（不同场景示例）===")
	featureOptions := []string{"自动备份", "监控告警", "日志聚合", "性能分析", "安全扫描", "负载均衡", "CDN 加速"}
	featureIndices, err := prompt.MultiSelect().
		Prompt("请选择要启用的系统功能（可多选）").
		Options(featureOptions).
		Default([]int{0, 1}).
		Run()
	if err != nil {
		return fmt.Errorf("选择失败: %w", err)
	}
	if len(featureIndices) == 0 {
		out.Warning("您没有选择任何功能")
	} else {
		var selectedNames []string
		for _, idx := range featureIndices {
			selectedNames = append(selectedNames, featureOptions[idx])
		}
		out.Success("您选择了: %s (索引: %v)", strings.Join(selectedNames, ", "), featureIndices)
	}
	out.Println("")

	out.Success("演示完成！感谢使用 MultiSelect 功能。")

	return nil
}

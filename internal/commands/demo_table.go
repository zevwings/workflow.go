package commands

import (
	"github.com/spf13/cobra"
	"github.com/zevwings/workflow/internal/prompt"
)

// NewDemoTableCmd 创建一个演示 Table 功能的命令
func NewDemoTableCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "demo-table",
		Short: "演示 Table 的表格显示功能",
		Long: `演示 Table 的各种表格显示功能：
- Table: 表格显示（支持边框、行分隔线、对齐方式）
- 基本用法: NewTable() + AddRow() + Render()
- 配置选项: SetBorder(), SetRowLine(), SetAlignment()

这个 demo 会依次展示所有功能，帮助您了解 Table 的用法。`,
		RunE: runDemoTable,
	}

	return cmd
}

func runDemoTable(cmd *cobra.Command, args []string) error {
	out := prompt.NewMessage(false)

	out.Info("欢迎使用 Table 功能演示")
	out.Println("")
	out.Info("本演示将展示以下功能：")
	out.Println("  1. Table - 基本表格（带边框和行分隔线）")
	out.Println("  2. Table - 无边框表格")
	out.Println("  3. Table - 无行分隔线表格")
	out.Println("  4. Table - 不同对齐方式（左对齐、居中、右对齐）")
	out.Println("  5. Table - 实际应用场景示例")
	out.Println("")

	// 1. 演示基本表格（带边框和行分隔线）
	out.Info("=== 演示 1: Table（基本表格 - 带边框和行分隔线）===")
	table1 := prompt.NewTable([]string{"项目", "状态", "进度", "说明"})
	table1.AddRow([]string{"功能 A", "✓ 完成", "100%", "已上线"})
	table1.AddRow([]string{"功能 B", "🔄 进行中", "75%", "预计下周完成"})
	table1.AddRow([]string{"功能 C", "⏳ 待开始", "0%", "等待资源"})
	table1.AddRow([]string{"功能 D", "✓ 完成", "100%", "已上线"})
	table1.Render()
	out.Println("")

	// 2. 演示无边框表格
	out.Info("=== 演示 2: Table（无边框表格）===")
	table2 := prompt.NewTable([]string{"姓名", "年龄", "职位", "部门"})
	table2.SetBorder(false)
	table2.SetRowLine(false) // 无边框模式下，也不显示行分隔线
	table2.AddRow([]string{"张三", "28", "高级工程师", "研发部"})
	table2.AddRow([]string{"李四", "32", "架构师", "技术部"})
	table2.AddRow([]string{"王五", "25", "初级工程师", "研发部"})
	table2.AddRow([]string{"赵六", "35", "技术总监", "技术部"})
	table2.Render()
	out.Println("")

	// 3. 演示无行分隔线表格
	out.Info("=== 演示 3: Table（无行分隔线表格）===")
	table3 := prompt.NewTable([]string{"命令", "描述", "示例"})
	table3.SetRowLine(false)
	table3.AddRow([]string{"git clone", "克隆仓库", "git clone https://github.com/user/repo.git"})
	table3.AddRow([]string{"git status", "查看状态", "git status"})
	table3.AddRow([]string{"git commit", "提交更改", "git commit -m \"message\""})
	table3.AddRow([]string{"git push", "推送更改", "git push origin main"})
	table3.Render()
	out.Println("")

	// 4. 演示不同对齐方式
	out.Info("=== 演示 4: Table（不同对齐方式）===")

	// 左对齐（默认）
	out.Println("左对齐（默认）：")
	table4a := prompt.NewTable([]string{"项目", "数量", "金额"})
	table4a.SetAlignment(prompt.ALIGN_LEFT)
	table4a.AddRow([]string{"商品 A", "10", "$100.00"})
	table4a.AddRow([]string{"商品 B", "5", "$50.00"})
	table4a.AddRow([]string{"商品 C", "20", "$200.00"})
	table4a.Render()
	out.Println("")

	// 居中
	out.Println("居中对齐：")
	table4b := prompt.NewTable([]string{"项目", "数量", "金额"})
	table4b.SetAlignment(prompt.ALIGN_CENTER)
	table4b.AddRow([]string{"商品 A", "10", "$100.00"})
	table4b.AddRow([]string{"商品 B", "5", "$50.00"})
	table4b.AddRow([]string{"商品 C", "20", "$200.00"})
	table4b.Render()
	out.Println("")

	// 右对齐
	out.Println("右对齐：")
	table4c := prompt.NewTable([]string{"项目", "数量", "金额"})
	table4c.SetAlignment(prompt.ALIGN_RIGHT)
	table4c.AddRow([]string{"商品 A", "10", "$100.00"})
	table4c.AddRow([]string{"商品 B", "5", "$50.00"})
	table4c.AddRow([]string{"商品 C", "20", "$200.00"})
	table4c.Render()
	out.Println("")

	// 5. 演示实际应用场景
	out.Info("=== 演示 5: Table（实际应用场景示例）===")

	// 场景 1: 系统检查结果
	out.Println("场景 1: 系统检查结果")
	table5a := prompt.NewTable([]string{"检查项", "状态", "说明"})
	table5a.AddRow([]string{"Git", "✓", "Git 已安装"})
	table5a.AddRow([]string{"Docker", "✓", "Docker 已安装"})
	table5a.AddRow([]string{"Kubernetes", "✗", "Kubernetes 未安装"})
	table5a.AddRow([]string{"网络连接", "✓", "网络连接正常"})
	table5a.Render()
	out.Println("")

	// 场景 2: 依赖包列表
	out.Println("场景 2: 依赖包列表")
	table5b := prompt.NewTable([]string{"包名", "版本", "状态", "说明"})
	table5b.AddRow([]string{"github.com/spf13/cobra", "v1.8.0", "✓", "已安装"})
	table5b.AddRow([]string{"github.com/charmbracelet/lipgloss", "v1.1.0", "✓", "已安装"})
	table5b.AddRow([]string{"github.com/charmbracelet/bubbletea", "v1.3.6", "✓", "已安装"})
	table5b.AddRow([]string{"github.com/olekukonko/tablewriter", "v0.0.5", "✗", "已移除（使用自定义实现）"})
	table5b.Render()
	out.Println("")

	// 场景 3: 性能指标
	out.Println("场景 3: 性能指标")
	table5c := prompt.NewTable([]string{"指标", "当前值", "目标值", "状态"})
	table5c.AddRow([]string{"响应时间", "120ms", "< 200ms", "✓ 正常"})
	table5c.AddRow([]string{"吞吐量", "1000 req/s", "> 800 req/s", "✓ 正常"})
	table5c.AddRow([]string{"错误率", "0.1%", "< 1%", "✓ 正常"})
	table5c.AddRow([]string{"CPU 使用率", "65%", "< 80%", "✓ 正常"})
	table5c.AddRow([]string{"内存使用率", "85%", "< 90%", "⚠ 警告"})
	table5c.Render()
	out.Println("")

	// 场景 4: 任务列表
	out.Println("场景 4: 任务列表")
	table5d := prompt.NewTable([]string{"任务 ID", "任务名称", "优先级", "状态", "负责人"})
	table5d.AddRow([]string{"TASK-001", "实现用户认证", "高", "进行中", "张三"})
	table5d.AddRow([]string{"TASK-002", "优化数据库查询", "中", "待开始", "李四"})
	table5d.AddRow([]string{"TASK-003", "修复登录 Bug", "高", "已完成", "王五"})
	table5d.AddRow([]string{"TASK-004", "编写单元测试", "中", "进行中", "赵六"})
	table5d.Render()
	out.Println("")

	out.Success("演示完成！感谢使用 Table 功能。")
	out.Println("")
	out.Info("💡 提示：")
	out.Println("  - Table 支持边框、行分隔线、对齐方式等配置")
	out.Println("  - 表格会自动计算列宽，确保内容对齐")
	out.Println("  - 表头会自动应用主题样式（PromptStyle）")
	out.Println("  - 边框颜色会跟随主题配置")

	return nil
}

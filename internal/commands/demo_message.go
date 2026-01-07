package commands

import (
	"github.com/spf13/cobra"
	"github.com/zevwings/workflow/internal/prompt"
)

// NewDemoMessageCmd 创建一个演示 Message 功能的命令
func NewDemoMessageCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "demo-message",
		Short: "演示 Message 的消息输出功能",
		Long: `演示 Message 的各种消息输出功能：
- Info: 信息消息
- Success: 成功消息
- Warning: 警告消息
- Error: 错误消息
- Debug: 调试消息（需要 verbose 模式）
- Print/Println: 普通输出

这个 demo 会依次展示所有功能，帮助您了解 Message 的用法。`,
		RunE: runDemoMessage,
	}

	// 添加 verbose 标志
	cmd.Flags().BoolP("verbose", "v", false, "启用 verbose 模式以显示 Debug 消息")

	return cmd
}

func runDemoMessage(cmd *cobra.Command, args []string) error {
	// 获取 verbose 标志
	verbose, _ := cmd.Flags().GetBool("verbose")
	out := prompt.NewMessage(verbose)

	out.Info("欢迎使用 Message 功能演示")
	out.Println("")
	out.Info("本演示将展示以下功能：")
	out.Println("  1. Info - 信息消息")
	out.Println("  2. Success - 成功消息")
	out.Println("  3. Warning - 警告消息")
	out.Println("  4. Error - 错误消息")
	out.Println("  5. Debug - 调试消息（需要 --verbose 标志）")
	out.Println("  6. Print - 普通输出（不换行）")
	out.Println("  7. Println - 普通输出（换行）")
	out.Println("")

	// 1. 演示 Info
	out.Info("=== 演示 1: Info（信息消息）===")
	out.Info("这是一条信息消息")
	out.Info("支持格式化输出: %s, %d", "示例", 123)
	out.Println("")

	// 2. 演示 Success
	out.Info("=== 演示 2: Success（成功消息）===")
	out.Success("操作成功完成")
	out.Success("成功创建了 %d 个文件", 5)
	out.Println("")

	// 3. 演示 Warning
	out.Info("=== 演示 3: Warning（警告消息）===")
	out.Warning("这是一个警告消息")
	out.Warning("磁盘空间不足: 剩余 %d%%", 15)
	out.Println("")

	// 4. 演示 Error
	out.Info("=== 演示 4: Error（错误消息）===")
	out.Error("这是一个错误消息")
	out.Error("文件读取失败: %s", "permission denied")
	out.Println("")

	// 5. 演示 Debug（需要 verbose 模式）
	out.Info("=== 演示 5: Debug（调试消息）===")
	if verbose {
		out.Debug("这是调试消息（verbose 模式已启用）")
		out.Debug("调试信息: 变量值 = %v", map[string]interface{}{
			"key1": "value1",
			"key2": 42,
		})
	} else {
		out.Println("提示: Debug 消息需要启用 verbose 模式")
		out.Println("     使用 --verbose 或 -v 标志来启用")
		out.Println("     例如: workflow demo-message --verbose")
	}
	out.Println("")

	// 6. 演示 Print
	out.Info("=== 演示 6: Print（普通输出，不换行）===")
	out.Print("这是第一行")
	out.Print("，这是第二行")
	out.Print("，这是第三行")
	out.Println("") // 手动换行
	out.Println("")

	// 7. 演示 Println
	out.Info("=== 演示 7: Println（普通输出，换行）===")
	out.Println("这是第一行")
	out.Println("这是第二行")
	out.Println("这是第三行")
	out.Println("")

	// 8. 演示组合使用
	out.Info("=== 演示 8: 组合使用示例 ===")
	out.Println("模拟一个工作流程：")
	out.Info("步骤 1: 检查配置文件...")
	out.Success("配置文件检查通过")
	out.Info("步骤 2: 验证依赖项...")
	out.Warning("发现 2 个过时的依赖项")
	out.Info("步骤 3: 执行构建...")
	out.Success("构建成功完成")
	out.Println("")

	// 9. 演示格式化输出
	out.Info("=== 演示 9: 格式化输出示例 ===")
	out.Success("用户 %s 已成功登录", "alice")
	out.Info("当前时间: %s", "2024-01-01 12:00:00")
	out.Warning("API 响应时间: %dms（超过阈值）", 1500)
	out.Error("连接失败: %s (错误代码: %d)", "timeout", 408)
	if verbose {
		out.Debug("请求详情: URL=%s, Method=%s", "https://api.example.com", "GET")
	}
	out.Println("")

	out.Success("演示完成！感谢使用 Message 功能。")
	out.Println("")
	out.Println("提示: 使用 --verbose 或 -v 标志可以查看 Debug 消息")

	return nil
}

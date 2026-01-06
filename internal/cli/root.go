package cli

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/your-org/workflow/internal/commands"
	"github.com/your-org/workflow/internal/lib/config"
	"github.com/your-org/workflow/internal/logging"
)

var (
	version   = "0.1.0"
	buildDate = "unknown"
	gitCommit = "unknown"
)

var rootCmd = &cobra.Command{
	Use:   "workflow",
	Short: "Workflow CLI - Git 工作流自动化工具",
	Long: `Workflow CLI 是一个功能强大的 Git 工作流自动化工具，
支持 PR 管理、Jira 集成、LLM 集成等功能。`,
	Version: version,
}

// Execute 执行根命令
func Execute() error {
	// 初始化日志系统
	initLogging()

	return rootCmd.Execute()
}

// initLogging 初始化日志系统
func initLogging() {
	// 尝试加载配置以获取日志级别
	manager, err := config.NewManager()
	if err != nil {
		// 如果配置管理器创建失败，使用默认设置
		logging.Init("info", "text", nil)
		return
	}

	// 尝试加载配置
	if err := manager.Load(); err != nil {
		// 配置文件不存在时使用默认设置
		logging.Init("info", "text", nil)
		return
	}

	// 从配置中获取日志级别和格式
	logLevel := manager.GetString("log.level")
	if logLevel == "" {
		logLevel = "info"
	}

	logFormat := "text" // 默认文本格式
	logging.Init(logLevel, logFormat, nil)
}

func init() {
	// 添加子命令
	rootCmd.AddCommand(commands.NewSetupCmd())
	rootCmd.AddCommand(commands.NewVersionCmd(version, buildDate, gitCommit))
	rootCmd.AddCommand(commands.NewConfigCmd())
	rootCmd.AddCommand(commands.NewCheckCmd())
	rootCmd.AddCommand(commands.NewDemoHuhCmd())
	// 添加分离的演示命令
	rootCmd.AddCommand(commands.NewDemoInputCmd())
	rootCmd.AddCommand(commands.NewDemoConfirmCmd())
	rootCmd.AddCommand(commands.NewDemoSelectCmd())
	rootCmd.AddCommand(commands.NewDemoMultiSelectCmd())

	// 设置版本模板
	rootCmd.SetVersionTemplate(fmt.Sprintf("workflow version %s\nBuild Date: %s\nGit Commit: %s\n", version, buildDate, gitCommit))
}

// GetRootCmd 返回根命令（用于测试）
func GetRootCmd() *cobra.Command {
	return rootCmd
}

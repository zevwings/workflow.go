package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/zevwings/workflow/internal/commands"
	"github.com/zevwings/workflow/internal/config"
	"github.com/zevwings/workflow/internal/logging"
)

var (
	version   string = "0.1.0"
	buildDate string = "unknown"
	gitCommit string = "unknown"
)

func main() {
	// 初始化日志系统
	initLogging()

	// 创建根命令
	rootCmd := &cobra.Command{
		Use:   "workflow",
		Short: "Workflow CLI - Git 工作流自动化工具",
		Long: `Workflow CLI 是一个功能强大的 Git 工作流自动化工具，
支持 PR 管理、Jira 集成、LLM 集成等功能。`,
		Version: version,
	}

	// 注册子命令
	rootCmd.AddCommand(commands.NewSetupCmd())
	rootCmd.AddCommand(commands.NewVersionCmd(version, buildDate, gitCommit))
	rootCmd.AddCommand(commands.NewConfigCmd())
	rootCmd.AddCommand(commands.NewCheckCmd())

	// 注册示例命令（仅在构建时使用 -tags=example 才会包含）
	registerExampleCommands(rootCmd)

	// 设置版本模板
	rootCmd.SetVersionTemplate(fmt.Sprintf("workflow version %s\nBuild Date: %s\nGit Commit: %s\n", version, buildDate, gitCommit))

	// 执行命令
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

// initLogging 初始化日志系统
func initLogging() {
	// 尝试加载配置以获取日志级别
	manager, err := config.NewGlobalManager()
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

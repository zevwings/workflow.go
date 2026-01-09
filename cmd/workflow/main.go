package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/sirupsen/logrus"
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

	// 设置版本模板
	rootCmd.SetVersionTemplate(fmt.Sprintf("workflow version %s\nBuild Date: %s\nGit Commit: %s\n", version, buildDate, gitCommit))

	// 执行命令
	if err := rootCmd.Execute(); err != nil {
		logger := logging.GetLogger()
		logger.WithError(err).Error("Command execution failed")
		os.Exit(1)
	}
}

// initLogging 初始化日志系统
func initLogging() {
	// 使用全局 logger（此时模块 logger 还未初始化）
	// 如果 Logger 还未初始化，创建一个临时 logger
	var logger *logrus.Logger
	if logging.Logger != nil {
		logger = logging.Logger
	} else {
		logger = logrus.New()
		logger.SetOutput(os.Stdout)
		logger.SetLevel(logrus.InfoLevel)
	}
	logger.Info("Initializing logging system")

	// 获取日志目录（使用 XDG_STATE_HOME）
	var logDir string
	workflowStateDir, err := config.StateDir()
	if err == nil {
		logDir = filepath.Join(workflowStateDir, "logs")
	} else {
		logger.WithError(err).Warn("Failed to get XDG state directory, using empty log directory")
		logDir = ""
	}

	// 尝试加载配置以获取日志级别
	manager, err := config.Global()
	if err != nil {
		// 如果配置管理器创建失败，使用默认设置
		logger.WithError(err).Warn("Failed to create config manager, using default logging settings")
		logging.InitWithFiles("info", "text", nil, logDir, false)
		return
	}

	// 尝试加载配置
	if err := manager.Load(); err != nil {
		// 配置文件不存在时使用默认设置
		logger.Debug("Config file not found, using default logging settings")
		logging.InitWithFiles("info", "text", nil, logDir, false)
		return
	}

	// 从配置中获取日志级别和格式
	// 可以直接使用 manager.LogConfig.Level 或 manager.Config.Log.Level
	logLevel := manager.LogConfig.Level
	if logLevel == "" {
		logLevel = "info"
	}

	logFormat := "text" // 默认文本格式
	logging.InitWithFiles(logLevel, logFormat, nil, logDir, false)

	// 记录初始化完成
	logger.WithFields(logrus.Fields{
		"level":   logLevel,
		"format":  logFormat,
		"log_dir": logDir,
	}).Info("Logging system initialized")
}

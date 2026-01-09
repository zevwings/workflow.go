//go:build example

package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/zevwings/workflow/internal/commands/example"
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
		Use:   "example",
		Short: "Workflow CLI Example - 演示和测试工具",
		Long: `Workflow CLI Example 是一个演示和测试工具，
用于展示 Workflow CLI 的各种功能和组件。

此工具包含以下演示命令：
- demo-prompt: 演示所有 Prompt 组件的交互功能
- demo-form: 演示 Form 模块的交互式表单功能`,
		Version: version,
	}

	// 注册示例命令
	rootCmd.AddCommand(example.NewDemoCmd())
	rootCmd.AddCommand(example.NewDemoFormCmd())

	// 设置版本模板
	rootCmd.SetVersionTemplate(fmt.Sprintf("example version %s\nBuild Date: %s\nGit Commit: %s\n", version, buildDate, gitCommit))

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
	logger := logging.Logger
	logger.Info("Initializing logging system")

	// 获取日志目录
	homeDir, err := os.UserHomeDir()
	var logDir string
	if err == nil {
		logDir = filepath.Join(homeDir, ".workflow", "logs")
	} else {
		logger.Warn("Failed to get user home directory, using default log directory")
	}

	// 尝试加载配置以获取日志级别
	manager, err := config.NewGlobalManager()
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
	logLevel := manager.GetString("log.level")
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

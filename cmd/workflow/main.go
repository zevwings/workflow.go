package main

import (
	"fmt"
	"io"
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
	// 先获取日志目录和配置，再决定临时 logger 的输出目标
	var logDir string
	workflowStateDir, err := config.StateDir()
	if err == nil {
		logDir = filepath.Join(workflowStateDir, "logs")
	} else {
		logDir = ""
	}

	// 尝试加载配置以获取日志级别（静默加载，不输出日志）
	var logLevel string
	manager, err := config.Global()
	if err == nil {
		// 尝试加载配置，忽略错误（配置文件可能不存在）
		_ = manager.Load()
		logLevel = manager.LogConfig.Level
	}

	// 在初始化日志系统前，先清理旧的日志文件
	if logDir != "" {
		_ = clearLogFiles(logDir)
	}

	// 根据配置决定临时 logger 的输出目标
	// 如果 level 为空，临时 logger 应该输出到文件（如果 logDir 不为空）或不输出
	var tempLoggerOutput io.Writer
	if logLevel == "" {
		// level 为空，只写入文件，不输出到控制台
		if logDir != "" {
			// 确保日志目录存在
			_ = os.MkdirAll(logDir, 0755)
			// 创建临时日志文件用于初始化期间的日志
			tempLogFile := filepath.Join(logDir, "init.log")
			file, err := os.OpenFile(tempLogFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
			if err == nil {
				tempLoggerOutput = file
			} else {
				// 如果无法创建文件，使用 io.Discard 丢弃输出
				tempLoggerOutput = io.Discard
			}
		} else {
			tempLoggerOutput = io.Discard
		}
	} else {
		// level 不为空，输出到控制台
		tempLoggerOutput = os.Stdout
	}

	// 创建临时 logger
	var logger *logrus.Logger
	if logging.Logger != nil {
		logger = logging.Logger
	} else {
		logger = logrus.New()
		logger.SetOutput(tempLoggerOutput)
		logger.SetLevel(logrus.InfoLevel)
	}

	logger.Info("Initializing logging system")

	// 如果配置管理器创建失败，使用默认设置（level 为空，只写入文件）
	if manager == nil {
		logger.WithError(err).Warn("Failed to create config manager, using default logging settings")
		logging.InitWithFiles("", "text", nil, logDir, false)
		return
	}

	// 从配置中获取日志级别和格式
	// 如果 level 为空，保持为空字符串，这样会只写入文件且只记录 info/warn/error 级别
	// 如果 level 不为空，根据 level 控制输出到控制台的等级
	logFormat := "text" // 默认文本格式

	// 如果 level 不为空，输出到控制台；如果为空，不输出到控制台（InitWithFiles 会自动处理）
	consoleOut := logLevel != ""
	logging.InitWithFiles(logLevel, logFormat, nil, logDir, consoleOut)

	// 记录初始化完成（使用初始化后的 logger）
	initLogger := logging.GetLogger()
	initLogger.WithFields(logging.Fields{
		"level":   logLevel,
		"format":  logFormat,
		"log_dir": logDir,
	}).Info("Logging system initialized")
}

// clearLogFiles 删除日志目录下的所有日志文件
func clearLogFiles(logDir string) error {
	// 检查目录是否存在
	if _, err := os.Stat(logDir); os.IsNotExist(err) {
		// 目录不存在，无需清理
		return nil
	}

	// 读取目录内容
	entries, err := os.ReadDir(logDir)
	if err != nil {
		return err
	}

	// 删除所有 .log 文件
	for _, entry := range entries {
		if !entry.IsDir() && filepath.Ext(entry.Name()) == ".log" {
			logFile := filepath.Join(logDir, entry.Name())
			if err := os.Remove(logFile); err != nil {
				// 记录错误但继续删除其他文件
				_ = err
			}
		}
	}

	return nil
}

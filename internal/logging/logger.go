package logging

import (
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"

	"github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
)

var (
	// Logger 全局日志实例（向后兼容）
	Logger *logrus.Logger

	// moduleLoggers 模块级别的 logger 缓存
	moduleLoggers = make(map[string]*logrus.Logger)
	moduleMutex   sync.RWMutex

	// logConfig 日志配置
	logConfig struct {
		level       string
		format      string
		logDir      string
		consoleOut  bool
		initialized bool
	}
)

// InitWithFiles 初始化日志系统（支持按模块分别输出到文件）
//
// 参数:
//   - level: 日志级别 (debug, info, warn, error)
//   - format: 日志格式 (text, json)
//   - output: 输出目标，如果为 nil 且 logDir 为空则输出到控制台
//   - logDir: 日志目录，如果为空则不创建文件日志。如果指定，每个模块会输出到 {logDir}/{module}.log
//   - consoleOut: 是否同时输出到控制台（仅在 logDir 不为空时有效）
func InitWithFiles(level string, format string, output io.Writer, logDir string, consoleOut bool) {
	// 保存配置
	logConfig.level = level
	logConfig.format = format
	logConfig.logDir = logDir
	logConfig.consoleOut = consoleOut
	logConfig.initialized = true

	// 创建日志目录
	if logDir != "" {
		if err := os.MkdirAll(logDir, 0755); err != nil {
			// 如果创建目录失败，记录错误但不中断程序
			_ = err
		}
	}

	// 创建全局 Logger（向后兼容）
	Logger = logrus.New()

	// 设置日志级别
	logLevel, err := logrus.ParseLevel(strings.ToLower(level))
	if err != nil {
		logLevel = logrus.InfoLevel
	}
	Logger.SetLevel(logLevel)

	// 设置输出格式
	switch strings.ToLower(format) {
	case "json":
		Logger.SetFormatter(&logrus.JSONFormatter{
			TimestampFormat: "2006-01-02 15:04:05",
		})
	default:
		Logger.SetFormatter(&logrus.TextFormatter{
			FullTimestamp:   true,
			TimestampFormat: "2006-01-02 15:04:05",
			ForceColors:     true,
		})
	}

	// 设置输出目标
	if output == nil {
		output = os.Stdout
	}
	Logger.SetOutput(output)
}

// Init 初始化日志系统（向后兼容，不创建文件日志）
//
// 参数:
//   - level: 日志级别 (debug, info, warn, error)
//   - format: 日志格式 (text, json)
//   - output: 输出目标，如果为 nil 则输出到控制台
func Init(level string, format string, output io.Writer) {
	InitWithFiles(level, format, output, "", true)
}

// SetLevel 设置日志级别
func SetLevel(level string) {
	logLevel, err := logrus.ParseLevel(strings.ToLower(level))
	if err != nil {
		logLevel = logrus.InfoLevel
	}
	if Logger != nil {
		Logger.SetLevel(logLevel)
	}
}

// Fields 字段映射类型，用于结构化日志
// 这是 logrus.Fields 的封装，避免外部代码直接依赖 logrus
type Fields map[string]interface{}

// WithField 添加字段
func WithField(key string, value interface{}) *logrus.Entry {
	if Logger != nil {
		return Logger.WithField(key, value)
	}
	return logrus.NewEntry(logrus.New())
}

// WithFields 添加多个字段
// 接受 Fields 类型（map[string]interface{}），内部转换为 logrus.Fields
func WithFields(fields Fields) *logrus.Entry {
	if Logger != nil {
		return Logger.WithFields(logrus.Fields(fields))
	}
	return logrus.NewEntry(logrus.New())
}

// WithError 添加错误字段
func WithError(err error) *logrus.Entry {
	if Logger != nil {
		return Logger.WithError(err)
	}
	return logrus.NewEntry(logrus.New())
}

// getCallerModule 获取调用者的模块名
//
// 通过 runtime.Caller 获取调用栈信息，提取模块名。
// 跳过 logging 包本身的调用，返回第一个外部调用者的模块名。
//
// 返回:
//   - string: 模块名，如 "http", "llm", "config" 等
func getCallerModule() string {
	// 跳过当前函数和 GetLogger 的调用栈
	skip := 2
	for i := skip; i < 10; i++ {
		pc, file, _, ok := runtime.Caller(i)
		if !ok {
			break
		}

		// 获取包路径
		fn := runtime.FuncForPC(pc)
		if fn == nil {
			continue
		}

		// 从函数名中提取包路径
		// 格式: github.com/zevwings/workflow/internal/http.(*httpClient).doRequest
		fullName := fn.Name()

		// 跳过 logging 包本身的调用
		if strings.Contains(fullName, "/internal/logging") {
			continue
		}

		// 提取包路径
		// 查找 "internal/" 或项目根路径
		parts := strings.Split(fullName, "/")
		for j, part := range parts {
			if part == "internal" && j+1 < len(parts) {
				// 提取 internal 后面的模块名
				modulePath := parts[j+1]
				// 如果还有子模块，取第一个子模块
				if dotIdx := strings.Index(modulePath, "."); dotIdx > 0 {
					modulePath = modulePath[:dotIdx]
				}
				return modulePath
			}
		}

		// 如果没找到 internal，尝试从文件路径提取
		if strings.Contains(file, "/internal/") {
			relPath := strings.Split(file, "/internal/")
			if len(relPath) > 1 {
				modulePath := strings.Split(relPath[1], "/")[0]
				return modulePath
			}
		}

		// 从文件名提取（作为后备方案）
		if file != "" {
			base := filepath.Base(file)
			if base != "logger.go" && base != "logger_test.go" {
				moduleName := strings.TrimSuffix(base, ".go")
				moduleName = strings.TrimSuffix(moduleName, "_test")
				return moduleName
			}
		}
	}

	return "unknown"
}

// getModuleLogger 获取模块级别的 logger
//
// 为每个模块创建独立的 logger 实例，每个模块的日志输出到独立的文件。
//
// 参数:
//   - module: 模块名，如 "http", "llm", "config" 等
//
// 返回:
//   - *logrus.Logger: 模块级别的 logger 实例
func getModuleLogger(module string) *logrus.Logger {
	moduleMutex.RLock()
	if logger, exists := moduleLoggers[module]; exists {
		moduleMutex.RUnlock()
		return logger
	}
	moduleMutex.RUnlock()

	// 创建新的 logger
	moduleMutex.Lock()
	defer moduleMutex.Unlock()

	// 双重检查
	if logger, exists := moduleLoggers[module]; exists {
		return logger
	}

	logger := logrus.New()

	// 设置日志级别
	logLevel, err := logrus.ParseLevel(strings.ToLower(logConfig.level))
	if err != nil {
		logLevel = logrus.InfoLevel
	}
	logger.SetLevel(logLevel)

	// 设置输出格式
	var formatter logrus.Formatter
	switch strings.ToLower(logConfig.format) {
	case "json":
		formatter = &logrus.JSONFormatter{
			TimestampFormat: "2006-01-02 15:04:05",
		}
	default:
		formatter = &logrus.TextFormatter{
			FullTimestamp:   true,
			TimestampFormat: "2006-01-02 15:04:05",
			ForceColors:     false, // 文件输出不使用颜色
		}
	}

	// 创建 handlers
	var handlers []io.Writer

	// 1. 控制台输出（如果启用）
	if logConfig.consoleOut {
		handlers = append(handlers, os.Stdout)
	}

	// 2. 模块文件输出
	if logConfig.logDir != "" {
		// 模块日志文件：{module}.log
		moduleFile := filepath.Join(logConfig.logDir, module+".log")
		fileWriter := &lumberjack.Logger{
			Filename:   moduleFile,
			MaxSize:    10, // 10MB
			MaxBackups: 5,  // 保留 5 个备份
			MaxAge:     30, // 保留 30 天
			Compress:   true,
		}
		handlers = append(handlers, fileWriter)
	}

	// 3. 统一错误日志文件（只记录 ERROR 级别以上）
	if logConfig.logDir != "" {
		errorFile := filepath.Join(logConfig.logDir, "error.log")
		errorWriter := &lumberjack.Logger{
			Filename:   errorFile,
			MaxSize:    10,
			MaxBackups: 5,
			MaxAge:     30,
			Compress:   true,
		}
		// 创建错误日志的 hook
		errorHook := &errorLogHook{
			writer:    errorWriter,
			formatter: formatter,
		}
		logger.AddHook(errorHook)
	}

	// 设置输出（多个 writer）
	if len(handlers) > 0 {
		if len(handlers) == 1 {
			logger.SetOutput(handlers[0])
		} else {
			// 多个 writer 使用 MultiWriter
			logger.SetOutput(io.MultiWriter(handlers...))
		}
	} else {
		// 如果没有 handlers，输出到 stdout
		logger.SetOutput(os.Stdout)
	}

	// 设置格式化器
	logger.SetFormatter(formatter)

	// 缓存 logger
	moduleLoggers[module] = logger

	return logger
}

// errorLogHook 错误日志 Hook，只记录 ERROR 级别以上的日志到 error.log
type errorLogHook struct {
	writer    io.Writer
	formatter logrus.Formatter
}

func (h *errorLogHook) Levels() []logrus.Level {
	return []logrus.Level{logrus.ErrorLevel, logrus.FatalLevel, logrus.PanicLevel}
}

func (h *errorLogHook) Fire(entry *logrus.Entry) error {
	formatted, err := h.formatter.Format(entry)
	if err != nil {
		return err
	}
	_, err = h.writer.Write(formatted)
	return err
}

// LoggerEntry 日志 Entry 的包装类型，提供封装的方法避免直接依赖 logrus
type LoggerEntry struct {
	entry *logrus.Entry
}

// WithFields 添加多个字段（使用 logging.Fields 类型）
func (e *LoggerEntry) WithFields(fields Fields) *LoggerEntry {
	return &LoggerEntry{entry: e.entry.WithFields(logrus.Fields(fields))}
}

// WithField 添加单个字段
func (e *LoggerEntry) WithField(key string, value interface{}) *LoggerEntry {
	return &LoggerEntry{entry: e.entry.WithField(key, value)}
}

// WithError 添加错误字段
func (e *LoggerEntry) WithError(err error) *LoggerEntry {
	return &LoggerEntry{entry: e.entry.WithError(err)}
}

// Debug 记录 Debug 级别日志
func (e *LoggerEntry) Debug(args ...interface{}) {
	e.entry.Debug(args...)
}

// Debugf 记录 Debug 级别日志（格式化）
func (e *LoggerEntry) Debugf(format string, args ...interface{}) {
	e.entry.Debugf(format, args...)
}

// Info 记录 Info 级别日志
func (e *LoggerEntry) Info(args ...interface{}) {
	e.entry.Info(args...)
}

// Infof 记录 Info 级别日志（格式化）
func (e *LoggerEntry) Infof(format string, args ...interface{}) {
	e.entry.Infof(format, args...)
}

// Warn 记录 Warn 级别日志
func (e *LoggerEntry) Warn(args ...interface{}) {
	e.entry.Warn(args...)
}

// Warnf 记录 Warn 级别日志（格式化）
func (e *LoggerEntry) Warnf(format string, args ...interface{}) {
	e.entry.Warnf(format, args...)
}

// Error 记录 Error 级别日志
func (e *LoggerEntry) Error(args ...interface{}) {
	e.entry.Error(args...)
}

// Errorf 记录 Error 级别日志（格式化）
func (e *LoggerEntry) Errorf(format string, args ...interface{}) {
	e.entry.Errorf(format, args...)
}

// GetLogger 获取带模块名的 logger Entry
//
// 自动获取调用者的模块名，并返回一个带有 "module" 字段的 LoggerEntry。
// 如果启用了文件日志，每个模块的日志会输出到独立的文件：{module}.log
//
// 使用方式:
//
//	logger := logging.GetLogger()
//	logger.Info("This is a log message")
//	// 输出到: {logDir}/http.log (如果启用了文件日志)
//	// 输出: time="..." level=info msg="This is a log message" module=http
//
//	logger := logging.GetLogger()
//	logger.WithField("key", "value").Debug("Debug message")
//	// 输出到: {logDir}/llm.log (如果启用了文件日志)
//	// 输出: time="..." level=debug msg="Debug message" module=llm key=value
//
//	logger := logging.GetLogger()
//	logger.WithFields(logging.Fields{"key": "value"}).Info("Structured log")
//	// 使用 logging.Fields 而不是 logrus.Fields，避免直接依赖 logrus
//
// 返回:
//   - *LoggerEntry: 带有 "module" 字段的日志 Entry（封装了 logrus.Entry）
func GetLogger() *LoggerEntry {
	var entry *logrus.Entry
	if !logConfig.initialized {
		// 如果未初始化，返回全局 Logger
		if Logger == nil {
			entry = logrus.NewEntry(logrus.New())
		} else {
			module := getCallerModule()
			entry = Logger.WithField("module", module)
		}
	} else {
		module := getCallerModule()
		moduleLogger := getModuleLogger(module)
		entry = moduleLogger.WithField("module", module)
	}
	return &LoggerEntry{entry: entry}
}

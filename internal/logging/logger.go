package logging

import (
	"io"
	"os"
	"strings"

	"github.com/sirupsen/logrus"
)

var (
	// Logger 全局日志实例
	Logger *logrus.Logger
)

// Init 初始化日志系统
func Init(level string, format string, output io.Writer) {
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

// Debug 记录 Debug 级别日志
func Debug(args ...interface{}) {
	if Logger != nil {
		Logger.Debug(args...)
	}
}

// Debugf 记录 Debug 级别日志（格式化）
func Debugf(format string, args ...interface{}) {
	if Logger != nil {
		Logger.Debugf(format, args...)
	}
}

// Info 记录 Info 级别日志
func Info(args ...interface{}) {
	if Logger != nil {
		Logger.Info(args...)
	}
}

// Infof 记录 Info 级别日志（格式化）
func Infof(format string, args ...interface{}) {
	if Logger != nil {
		Logger.Infof(format, args...)
	}
}

// Warn 记录 Warn 级别日志
func Warn(args ...interface{}) {
	if Logger != nil {
		Logger.Warn(args...)
	}
}

// Warnf 记录 Warn 级别日志（格式化）
func Warnf(format string, args ...interface{}) {
	if Logger != nil {
		Logger.Warnf(format, args...)
	}
}

// Error 记录 Error 级别日志
func Error(args ...interface{}) {
	if Logger != nil {
		Logger.Error(args...)
	}
}

// Errorf 记录 Error 级别日志（格式化）
func Errorf(format string, args ...interface{}) {
	if Logger != nil {
		Logger.Errorf(format, args...)
	}
}

// Fatal 记录 Fatal 级别日志并退出
func Fatal(args ...interface{}) {
	if Logger != nil {
		Logger.Fatal(args...)
	}
}

// Fatalf 记录 Fatal 级别日志并退出（格式化）
func Fatalf(format string, args ...interface{}) {
	if Logger != nil {
		Logger.Fatalf(format, args...)
	}
}

// WithField 添加字段
func WithField(key string, value interface{}) *logrus.Entry {
	if Logger != nil {
		return Logger.WithField(key, value)
	}
	return logrus.NewEntry(logrus.New())
}

// WithFields 添加多个字段
func WithFields(fields logrus.Fields) *logrus.Entry {
	if Logger != nil {
		return Logger.WithFields(fields)
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

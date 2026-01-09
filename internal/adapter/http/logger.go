package http

import (
	"github.com/go-resty/resty/v2"
	"github.com/zevwings/workflow/internal/logging"
)

// LogrusLogger 实现 resty.Logger 接口，将 Resty 的日志转发到 Logrus
//
// 此适配器将 go-resty 的日志接口转换为项目内部的 logging 包。
// 通过此适配器，Resty 的所有日志输出都会通过 logging 包进行记录。
//
// 使用方式:
//
//	import adapterhttp "github.com/zevwings/workflow/internal/adapter/http"
//
//	client := resty.New()
//	client.SetLogger(adapterhttp.NewLogrusLogger())
type LogrusLogger struct{}

// NewLogrusLogger 创建新的 LogrusLogger 实例
//
// 返回:
//   - *LogrusLogger: LogrusLogger 实例，实现了 resty.Logger 接口
func NewLogrusLogger() *LogrusLogger {
	return &LogrusLogger{}
}

// Errorf 记录错误日志
//
// 实现 resty.Logger 接口的 Errorf 方法。
// 将 Resty 的错误日志转发到 logging，并标识模块为 "http"。
//
// 参数:
//   - format: 日志格式字符串
//   - v: 格式化参数
func (l *LogrusLogger) Errorf(format string, v ...interface{}) {
	logging.WithField("module", "http").Errorf(format, v...)
}

// Warnf 记录警告日志
//
// 实现 resty.Logger 接口的 Warnf 方法。
// 将 Resty 的警告日志转发到 logging，并标识模块为 "http"。
//
// 参数:
//   - format: 日志格式字符串
//   - v: 格式化参数
func (l *LogrusLogger) Warnf(format string, v ...interface{}) {
	logging.WithField("module", "http").Warnf(format, v...)
}

// Debugf 记录调试日志
//
// 实现 resty.Logger 接口的 Debugf 方法。
// 将 Resty 的调试日志转发到 logging，并标识模块为 "http"。
//
// 参数:
//   - format: 日志格式字符串
//   - v: 格式化参数
func (l *LogrusLogger) Debugf(format string, v ...interface{}) {
	logging.WithField("module", "http").Debugf(format, v...)
}

// Verify that LogrusLogger implements resty.Logger interface
var _ resty.Logger = (*LogrusLogger)(nil)

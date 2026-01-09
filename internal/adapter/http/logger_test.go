package http

import (
	"testing"

	"github.com/go-resty/resty/v2"
	"github.com/zevwings/workflow/internal/logging"
)

// TestLogrusLogger_ImplementsRestyLogger 测试 LogrusLogger 是否实现了 resty.Logger 接口
func TestLogrusLogger_ImplementsRestyLogger(t *testing.T) {
	var _ resty.Logger = (*LogrusLogger)(nil)
}

// TestNewLogrusLogger 测试 NewLogrusLogger 函数
func TestNewLogrusLogger(t *testing.T) {
	logger := NewLogrusLogger()
	if logger == nil {
		t.Fatal("NewLogrusLogger() returned nil")
	}
}

// TestLogrusLogger_Methods 测试 LogrusLogger 的所有方法
func TestLogrusLogger_Methods(t *testing.T) {
	// 初始化日志系统（测试需要）
	logging.Init("debug", "text", nil)

	logger := NewLogrusLogger()

	// 测试 Errorf
	logger.Errorf("Test error message: %s", "error")

	// 测试 Warnf
	logger.Warnf("Test warning message: %s", "warning")

	// 测试 Debugf
	logger.Debugf("Test debug message: %s", "debug")
}

// TestLogrusLogger_WithRestyClient 测试 LogrusLogger 与 Resty 客户端的集成
func TestLogrusLogger_WithRestyClient(t *testing.T) {
	// 初始化日志系统（测试需要）
	logging.Init("debug", "text", nil)

	// 创建 Resty 客户端并设置 Logger
	client := resty.New()
	logger := NewLogrusLogger()
	client.SetLogger(logger)

	// 验证 Logger 已设置（通过调用 SetLogger 不会返回错误即可）
	// Resty 内部会使用我们设置的 Logger
	if logger == nil {
		t.Error("Logger is nil")
	}
}

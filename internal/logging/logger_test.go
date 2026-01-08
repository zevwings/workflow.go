package logging

import (
	"bytes"
	"encoding/json"
	"io"
	"strings"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ==================== Init 测试 ====================

func TestInit(t *testing.T) {
	tests := []struct {
		name   string
		level  string
		format string
		output io.Writer
	}{
		{
			name:   "默认配置",
			level:  "info",
			format: "text",
			output: nil,
		},
		{
			name:   "JSON 格式",
			level:  "debug",
			format: "json",
			output: nil,
		},
		{
			name:   "自定义输出",
			level:  "warn",
			format: "text",
			output: &bytes.Buffer{},
		},
		{
			name:   "无效日志级别",
			level:  "invalid",
			format: "text",
			output: nil,
		},
		{
			name:   "大写日志级别",
			level:  "DEBUG",
			format: "text",
			output: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange: 重置 Logger
			Logger = nil

			// Act: 初始化日志系统
			Init(tt.level, tt.format, tt.output)

			// Assert: 验证 Logger 已创建
			require.NotNil(t, Logger)
			assert.NotNil(t, Logger.Out)

			// 验证日志级别
			expectedLevel, err := parseLevel(tt.level)
			if err != nil {
				expectedLevel = "info"
			}
			// logrus 返回 "warning" 而不是 "warn"
			actualLevel := Logger.GetLevel().String()
			if expectedLevel == "warn" {
				assert.Equal(t, "warning", actualLevel)
			} else {
				assert.Equal(t, expectedLevel, actualLevel)
			}
		})
	}
}

func parseLevel(level string) (string, error) {
	level = strings.ToLower(level)
	validLevels := map[string]string{
		"debug": "debug",
		"info":  "info",
		"warn":  "warn",
		"error": "error",
		"fatal": "fatal",
		"panic": "panic",
	}
	if val, ok := validLevels[level]; ok {
		return val, nil
	}
	return "info", nil
}

func TestInit_JSONFormat(t *testing.T) {
	// Arrange: 重置 Logger
	Logger = nil
	buf := &bytes.Buffer{}

	// Act: 初始化 JSON 格式日志
	Init("info", "json", buf)

	// Assert: 验证格式为 JSON
	require.NotNil(t, Logger)
	_, ok := Logger.Formatter.(*logrus.JSONFormatter)
	assert.True(t, ok, "应该使用 JSON 格式化器")

	// 验证输出
	Logger.Info("test message")
	output := buf.String()
	assert.Contains(t, output, "test message")
	assert.Contains(t, output, "\"level\"")
}

func TestInit_TextFormat(t *testing.T) {
	// Arrange: 重置 Logger
	Logger = nil
	buf := &bytes.Buffer{}

	// Act: 初始化文本格式日志
	Init("info", "text", buf)

	// Assert: 验证格式为文本
	require.NotNil(t, Logger)
	_, ok := Logger.Formatter.(*logrus.TextFormatter)
	assert.True(t, ok, "应该使用文本格式化器")

	// 验证输出
	Logger.Info("test message")
	output := buf.String()
	assert.Contains(t, output, "test message")
}

// ==================== SetLevel 测试 ====================

func TestSetLevel(t *testing.T) {
	tests := []struct {
		name     string
		level    string
		expected string
	}{
		{
			name:     "设置 debug 级别",
			level:    "debug",
			expected: "debug",
		},
		{
			name:     "设置 info 级别",
			level:    "info",
			expected: "info",
		},
		{
			name:     "设置 warn 级别",
			level:    "warn",
			expected: "warning", // logrus 返回 "warning" 而不是 "warn"
		},
		{
			name:     "设置 error 级别",
			level:    "error",
			expected: "error",
		},
		{
			name:     "无效级别使用默认",
			level:    "invalid",
			expected: "info",
		},
		{
			name:     "大写级别",
			level:    "DEBUG",
			expected: "debug",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange: 初始化 Logger
			Init("info", "text", nil)

			// Act: 设置日志级别
			SetLevel(tt.level)

			// Assert: 验证级别已设置
			assert.Equal(t, tt.expected, Logger.GetLevel().String())
		})
	}
}

func TestSetLevel_NilLogger(t *testing.T) {
	// Arrange: 确保 Logger 为 nil
	Logger = nil

	// Act: 设置日志级别（不应该 panic）
	SetLevel("debug")

	// Assert: Logger 仍为 nil（函数应该安全处理 nil）
	assert.Nil(t, Logger)
}

// ==================== Debug/Debugf 测试 ====================

func TestDebug(t *testing.T) {
	// Arrange: 初始化 Logger 并设置 debug 级别
	buf := &bytes.Buffer{}
	Init("debug", "text", buf)

	// Act: 记录 debug 日志
	Debug("debug message")

	// Assert: 验证日志已输出
	output := buf.String()
	assert.Contains(t, output, "debug message")
}

func TestDebugf(t *testing.T) {
	// Arrange: 初始化 Logger 并设置 debug 级别
	buf := &bytes.Buffer{}
	Init("debug", "text", buf)

	// Act: 记录格式化 debug 日志
	Debugf("debug message: %s", "value")

	// Assert: 验证日志已输出
	output := buf.String()
	assert.Contains(t, output, "debug message: value")
}

func TestDebug_LevelFilter(t *testing.T) {
	// Arrange: 初始化 Logger 并设置 info 级别（高于 debug）
	buf := &bytes.Buffer{}
	Init("info", "text", buf)

	// Act: 记录 debug 日志
	Debug("debug message")

	// Assert: 验证日志未输出（被过滤）
	output := buf.String()
	assert.NotContains(t, output, "debug message")
}

func TestDebug_NilLogger(t *testing.T) {
	// Arrange: 确保 Logger 为 nil
	Logger = nil

	// Act: 记录 debug 日志（不应该 panic）
	Debug("debug message")
	Debugf("debug message: %s", "value")

	// Assert: 不应该 panic
	assert.Nil(t, Logger)
}

// ==================== Info/Infof 测试 ====================

func TestInfo(t *testing.T) {
	// Arrange: 初始化 Logger
	buf := &bytes.Buffer{}
	Init("info", "text", buf)

	// Act: 记录 info 日志
	Info("info message")

	// Assert: 验证日志已输出
	output := buf.String()
	assert.Contains(t, output, "info message")
}

func TestInfof(t *testing.T) {
	// Arrange: 初始化 Logger
	buf := &bytes.Buffer{}
	Init("info", "text", buf)

	// Act: 记录格式化 info 日志
	Infof("info message: %s", "value")

	// Assert: 验证日志已输出
	output := buf.String()
	assert.Contains(t, output, "info message: value")
}

func TestInfo_NilLogger(t *testing.T) {
	// Arrange: 确保 Logger 为 nil
	Logger = nil

	// Act: 记录 info 日志（不应该 panic）
	Info("info message")
	Infof("info message: %s", "value")

	// Assert: 不应该 panic
	assert.Nil(t, Logger)
}

// ==================== Warn/Warnf 测试 ====================

func TestWarn(t *testing.T) {
	// Arrange: 初始化 Logger
	buf := &bytes.Buffer{}
	Init("warn", "text", buf)

	// Act: 记录 warn 日志
	Warn("warn message")

	// Assert: 验证日志已输出
	output := buf.String()
	assert.Contains(t, output, "warn message")
}

func TestWarnf(t *testing.T) {
	// Arrange: 初始化 Logger
	buf := &bytes.Buffer{}
	Init("warn", "text", buf)

	// Act: 记录格式化 warn 日志
	Warnf("warn message: %s", "value")

	// Assert: 验证日志已输出
	output := buf.String()
	assert.Contains(t, output, "warn message: value")
}

func TestWarn_NilLogger(t *testing.T) {
	// Arrange: 确保 Logger 为 nil
	Logger = nil

	// Act: 记录 warn 日志（不应该 panic）
	Warn("warn message")
	Warnf("warn message: %s", "value")

	// Assert: 不应该 panic
	assert.Nil(t, Logger)
}

// ==================== Error/Errorf 测试 ====================

func TestError(t *testing.T) {
	// Arrange: 初始化 Logger
	buf := &bytes.Buffer{}
	Init("error", "text", buf)

	// Act: 记录 error 日志
	Error("error message")

	// Assert: 验证日志已输出
	output := buf.String()
	assert.Contains(t, output, "error message")
}

func TestErrorf(t *testing.T) {
	// Arrange: 初始化 Logger
	buf := &bytes.Buffer{}
	Init("error", "text", buf)

	// Act: 记录格式化 error 日志
	Errorf("error message: %s", "value")

	// Assert: 验证日志已输出
	output := buf.String()
	assert.Contains(t, output, "error message: value")
}

func TestError_NilLogger(t *testing.T) {
	// Arrange: 确保 Logger 为 nil
	Logger = nil

	// Act: 记录 error 日志（不应该 panic）
	Error("error message")
	Errorf("error message: %s", "value")

	// Assert: 不应该 panic
	assert.Nil(t, Logger)
}

// ==================== Fatal/Fatalf 测试 ====================

// 注意：Fatal 函数会调用 os.Exit，无法直接测试
// 这里只测试函数不会 panic，以及当 Logger 为 nil 时的行为
func TestFatal_NilLogger(t *testing.T) {
	// Arrange: 确保 Logger 为 nil
	Logger = nil

	// Act: 记录 fatal 日志（不应该 panic，但也不会退出，因为 Logger 为 nil）
	// 注意：由于 Fatal 会调用 os.Exit，我们无法完全测试它
	// 这里只验证函数不会 panic
	Fatal("fatal message")
	Fatalf("fatal message: %s", "value")

	// Assert: 不应该 panic
	assert.Nil(t, Logger)
}

// ==================== WithField/WithFields/WithError 测试 ====================

func TestWithField(t *testing.T) {
	// Arrange: 初始化 Logger
	buf := &bytes.Buffer{}
	Init("info", "json", buf)

	// Act: 添加字段并记录日志
	entry := WithField("key", "value")
	entry.Info("message with field")

	// Assert: 验证字段已添加
	output := buf.String()
	var logData map[string]interface{}
	err := json.Unmarshal([]byte(output), &logData)
	require.NoError(t, err)
	assert.Equal(t, "value", logData["key"])
	assert.Equal(t, "message with field", logData["msg"])
}

func TestWithFields(t *testing.T) {
	// Arrange: 初始化 Logger
	buf := &bytes.Buffer{}
	Init("info", "json", buf)

	// Act: 添加多个字段并记录日志
	entry := WithFields(logrus.Fields{
		"key1": "value1",
		"key2": "value2",
	})
	entry.Info("message with fields")

	// Assert: 验证字段已添加
	output := buf.String()
	var logData map[string]interface{}
	err := json.Unmarshal([]byte(output), &logData)
	require.NoError(t, err)
	assert.Equal(t, "value1", logData["key1"])
	assert.Equal(t, "value2", logData["key2"])
}

func TestWithError(t *testing.T) {
	// Arrange: 初始化 Logger
	buf := &bytes.Buffer{}
	Init("error", "json", buf)

	// Act: 添加错误并记录日志
	err := assert.AnError
	entry := WithError(err)
	entry.Error("message with error")

	// Assert: 验证错误已添加
	output := buf.String()
	var logData map[string]interface{}
	err2 := json.Unmarshal([]byte(output), &logData)
	require.NoError(t, err2)
	assert.NotNil(t, logData["error"])
}

func TestWithField_NilLogger(t *testing.T) {
	// Arrange: 确保 Logger 为 nil
	Logger = nil

	// Act: 添加字段（应该返回一个空的 Entry）
	entry := WithField("key", "value")

	// Assert: 验证返回了 Entry（即使 Logger 为 nil）
	assert.NotNil(t, entry)
}

func TestWithFields_NilLogger(t *testing.T) {
	// Arrange: 确保 Logger 为 nil
	Logger = nil

	// Act: 添加字段（应该返回一个空的 Entry）
	entry := WithFields(logrus.Fields{
		"key1": "value1",
	})

	// Assert: 验证返回了 Entry（即使 Logger 为 nil）
	assert.NotNil(t, entry)
}

func TestWithError_NilLogger(t *testing.T) {
	// Arrange: 确保 Logger 为 nil
	Logger = nil

	// Act: 添加错误（应该返回一个空的 Entry）
	entry := WithError(assert.AnError)

	// Assert: 验证返回了 Entry（即使 Logger 为 nil）
	assert.NotNil(t, entry)
}

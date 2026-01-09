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
	entry := WithFields(Fields{
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
	entry := WithFields(Fields{
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

// ==================== GetLogger 测试 ====================

func TestGetLogger(t *testing.T) {
	// Arrange: 初始化 Logger
	buf := &bytes.Buffer{}
	Init("info", "json", buf)

	// Act: 获取带模块名的 logger
	logger := GetLogger()
	logger.Info("test message")

	// Assert: 验证日志包含模块字段
	output := buf.String()
	var logData map[string]interface{}
	err := json.Unmarshal([]byte(output), &logData)
	require.NoError(t, err)
	assert.Contains(t, logData, "module")
	assert.Equal(t, "test message", logData["msg"])
	
	// 验证模块名不是 "unknown"（说明成功获取了模块名）
	moduleName, ok := logData["module"].(string)
	assert.True(t, ok)
	assert.NotEqual(t, "unknown", moduleName)
}

func TestGetLogger_TextFormat(t *testing.T) {
	// Arrange: 初始化 Logger（文本格式）
	buf := &bytes.Buffer{}
	Init("info", "text", buf)

	// Act: 获取带模块名的 logger
	logger := GetLogger()
	logger.Info("test message")

	// Assert: 验证日志包含模块信息
	output := buf.String()
	assert.Contains(t, output, "test message")
	// 文本格式中，模块信息可能包含颜色代码，所以检查 "module" 关键字
	assert.Contains(t, strings.ToLower(output), "module")
}

func TestGetLogger_NilLogger(t *testing.T) {
	// Arrange: 确保 Logger 为 nil
	Logger = nil

	// Act: 获取 logger（应该返回一个空的 Entry）
	logger := GetLogger()

	// Assert: 验证返回了 Entry（即使 Logger 为 nil）
	assert.NotNil(t, logger)
}

func TestGetLogger_WithFields(t *testing.T) {
	// Arrange: 初始化 Logger
	buf := &bytes.Buffer{}
	Init("info", "json", buf)

	// Act: 获取 logger 并添加额外字段
	logger := GetLogger()
	logger.WithFields(Fields{"key": "value"}).Info("test message")

	// Assert: 验证日志包含模块字段和额外字段
	output := buf.String()
	var logData map[string]interface{}
	err := json.Unmarshal([]byte(output), &logData)
	require.NoError(t, err)
	assert.Contains(t, logData, "module")
	assert.Equal(t, "value", logData["key"])
	assert.Equal(t, "test message", logData["msg"])
}

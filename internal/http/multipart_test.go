//go:build test

package http

import (
	"strings"
	"testing"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/stretchr/testify/assert"
)

// ==================== NewMultipartRequestConfig 测试 ====================

func TestNewMultipartRequestConfig(t *testing.T) {
	// Act: 创建新的 MultipartRequestConfig
	config := NewMultipartRequestConfig()

	// Assert: 验证默认值
	assert.NotNil(t, config)
	assert.NotNil(t, config.Headers)
	assert.Empty(t, config.MultipartFields)
	assert.Nil(t, config.Query)
	assert.Nil(t, config.Auth)
	assert.Zero(t, config.Timeout)
	assert.Nil(t, config.Retry)
}

// ==================== WithMultipartField 测试 ====================

func TestMultipartRequestConfig_WithMultipartField(t *testing.T) {
	// Arrange: 创建配置
	config := NewMultipartRequestConfig()

	// Act: 添加单个字段
	field := MultipartField{
		ParamName: "test-field",
		FileName:  "test.txt",
		Reader:    strings.NewReader("test content"),
	}
	result := config.WithMultipartField(field)

	// Assert: 验证链式调用和字段添加
	assert.Equal(t, config, result)
	assert.Len(t, config.MultipartFields, 1)
	assert.Equal(t, "test-field", config.MultipartFields[0].ParamName)
	assert.Equal(t, "test.txt", config.MultipartFields[0].FileName)
}

func TestMultipartRequestConfig_WithMultipartField_MultipleFields(t *testing.T) {
	// Arrange: 创建配置
	config := NewMultipartRequestConfig()

	// Act: 添加多个字段
	config.WithMultipartField(MultipartField{
		ParamName: "field1",
		FileName:  "file1.txt",
	}).
		WithMultipartField(MultipartField{
			ParamName: "field2",
			FileName:  "file2.txt",
		})

	// Assert: 验证所有字段都已添加
	assert.Len(t, config.MultipartFields, 2)
	assert.Equal(t, "field1", config.MultipartFields[0].ParamName)
	assert.Equal(t, "field2", config.MultipartFields[1].ParamName)
}

// ==================== WithMultipartFields 测试 ====================

func TestMultipartRequestConfig_WithMultipartFields(t *testing.T) {
	// Arrange: 创建配置和字段列表
	config := NewMultipartRequestConfig()
	fields := []MultipartField{
		{
			ParamName: "field1",
			FileName:  "file1.txt",
		},
		{
			ParamName: "field2",
			FileName:  "file2.txt",
		},
		{
			ParamName: "field3",
			FileName:  "file3.txt",
		},
	}

	// Act: 设置多个字段
	result := config.WithMultipartFields(fields)

	// Assert: 验证链式调用和字段设置
	assert.Equal(t, config, result)
	assert.Len(t, config.MultipartFields, 3)
	assert.Equal(t, fields, config.MultipartFields)
}

func TestMultipartRequestConfig_WithMultipartFields_Empty(t *testing.T) {
	// Arrange: 创建配置并添加一些字段
	config := NewMultipartRequestConfig()
	config.WithMultipartField(MultipartField{
		ParamName: "field1",
		FileName:  "file1.txt",
	})

	// Act: 设置为空列表
	config.WithMultipartFields([]MultipartField{})

	// Assert: 验证字段列表已清空
	assert.Len(t, config.MultipartFields, 0)
}

// ==================== WithQuery 测试 ====================

func TestMultipartRequestConfig_WithQuery(t *testing.T) {
	// Arrange: 创建配置
	config := NewMultipartRequestConfig()

	// Act: 设置查询参数（map）
	query := map[string]string{
		"key1": "value1",
		"key2": "value2",
	}
	result := config.WithQuery(query)

	// Assert: 验证链式调用和查询参数设置
	assert.Equal(t, config, result)
	assert.Equal(t, query, config.Query)
}

func TestMultipartRequestConfig_WithQuery_Struct(t *testing.T) {
	// Arrange: 创建配置
	config := NewMultipartRequestConfig()

	// Act: 设置查询参数（结构体）
	type QueryParams struct {
		Key1 string
		Key2 string
	}
	query := QueryParams{
		Key1: "value1",
		Key2: "value2",
	}
	result := config.WithQuery(query)

	// Assert: 验证链式调用和查询参数设置
	assert.Equal(t, config, result)
	assert.Equal(t, query, config.Query)
}

// ==================== WithAuth 测试 ====================

func TestMultipartRequestConfig_WithAuth(t *testing.T) {
	// Arrange: 创建配置
	config := NewMultipartRequestConfig()
	auth := &Authorization{
		Username: "user",
		Password: "pass",
	}

	// Act: 设置认证信息
	result := config.WithAuth(auth)

	// Assert: 验证链式调用和认证信息设置
	assert.Equal(t, config, result)
	assert.Equal(t, auth, config.Auth)
}

func TestMultipartRequestConfig_WithAuth_Nil(t *testing.T) {
	// Arrange: 创建配置并设置认证信息
	config := NewMultipartRequestConfig()
	config.WithAuth(&Authorization{
		Username: "user",
		Password: "pass",
	})

	// Act: 设置为 nil
	config.WithAuth(nil)

	// Assert: 验证认证信息已清除
	assert.Nil(t, config.Auth)
}

// ==================== WithHeader 测试 ====================

func TestMultipartRequestConfig_WithHeader(t *testing.T) {
	// Arrange: 创建配置
	config := NewMultipartRequestConfig()

	// Act: 设置单个 Header
	result := config.WithHeader("X-Custom-Header", "custom-value")

	// Assert: 验证链式调用和 Header 设置
	assert.Equal(t, config, result)
	assert.Equal(t, "custom-value", config.Headers["X-Custom-Header"])
}

func TestMultipartRequestConfig_WithHeader_Multiple(t *testing.T) {
	// Arrange: 创建配置
	config := NewMultipartRequestConfig()

	// Act: 设置多个 Headers
	config.WithHeader("Header1", "value1").
		WithHeader("Header2", "value2").
		WithHeader("Header3", "value3")

	// Assert: 验证所有 Headers 都已设置
	assert.Equal(t, "value1", config.Headers["Header1"])
	assert.Equal(t, "value2", config.Headers["Header2"])
	assert.Equal(t, "value3", config.Headers["Header3"])
}

// ==================== WithHeaders 测试 ====================

func TestMultipartRequestConfig_WithHeaders(t *testing.T) {
	// Arrange: 创建配置
	config := NewMultipartRequestConfig()
	headers := map[string]string{
		"Header1": "value1",
		"Header2": "value2",
		"Header3": "value3",
	}

	// Act: 设置多个 Headers
	result := config.WithHeaders(headers)

	// Assert: 验证链式调用和 Headers 设置
	assert.Equal(t, config, result)
	assert.Equal(t, headers["Header1"], config.Headers["Header1"])
	assert.Equal(t, headers["Header2"], config.Headers["Header2"])
	assert.Equal(t, headers["Header3"], config.Headers["Header3"])
}

func TestMultipartRequestConfig_WithHeaders_Merge(t *testing.T) {
	// Arrange: 创建配置并设置一个 Header
	config := NewMultipartRequestConfig()
	config.WithHeader("Existing", "existing-value")

	// Act: 设置多个 Headers（应该合并）
	headers := map[string]string{
		"New1": "new-value1",
		"New2": "new-value2",
	}
	config.WithHeaders(headers)

	// Assert: 验证 Headers 已合并
	assert.Equal(t, "existing-value", config.Headers["Existing"])
	assert.Equal(t, "new-value1", config.Headers["New1"])
	assert.Equal(t, "new-value2", config.Headers["New2"])
}

// ==================== WithTimeout 测试 ====================

func TestMultipartRequestConfig_WithTimeout(t *testing.T) {
	// Arrange: 创建配置
	config := NewMultipartRequestConfig()
	timeout := 10 * time.Second

	// Act: 设置超时时间
	result := config.WithTimeout(timeout)

	// Assert: 验证链式调用和超时时间设置
	assert.Equal(t, config, result)
	assert.Equal(t, timeout, config.Timeout)
}

func TestMultipartRequestConfig_WithTimeout_Zero(t *testing.T) {
	// Arrange: 创建配置并设置超时时间
	config := NewMultipartRequestConfig()
	config.WithTimeout(10 * time.Second)

	// Act: 设置为零值
	config.WithTimeout(0)

	// Assert: 验证超时时间已清除
	assert.Zero(t, config.Timeout)
}

// ==================== WithRetry 测试 ====================

func TestMultipartRequestConfig_WithRetry(t *testing.T) {
	// Arrange: 创建配置
	config := NewMultipartRequestConfig()
	retry := NewRetryConfig().
		WithRetryCount(3).
		WithRetryWaitTime(1 * time.Second)

	// Act: 设置重试配置
	result := config.WithRetry(retry)

	// Assert: 验证链式调用和重试配置设置
	assert.Equal(t, config, result)
	assert.Equal(t, retry, config.Retry)
}

func TestMultipartRequestConfig_WithRetry_Nil(t *testing.T) {
	// Arrange: 创建配置并设置重试配置
	config := NewMultipartRequestConfig()
	config.WithRetry(NewRetryConfig())

	// Act: 设置为 nil
	config.WithRetry(nil)

	// Assert: 验证重试配置已清除
	assert.Nil(t, config.Retry)
}

// ==================== ensureHeaders 测试 ====================

func TestMultipartRequestConfig_ensureHeaders(t *testing.T) {
	// Arrange: 创建配置，Headers 为 nil
	config := &MultipartRequestConfig{
		baseRequestConfig: baseRequestConfig{
			Headers: nil,
		},
	}

	// Act: 调用 ensureHeaders（通过 WithHeader 间接调用）
	config.WithHeader("Test-Header", "test-value")

	// Assert: 验证 Headers 已初始化
	assert.NotNil(t, config.Headers)
	assert.Equal(t, "test-value", config.Headers["Test-Header"])
}

// ==================== applyToRequest 测试 ====================

func TestMultipartRequestConfig_applyToRequest(t *testing.T) {
	// Arrange: 创建配置和 resty 请求
	config := NewMultipartRequestConfig().
		WithMultipartField(MultipartField{
			ParamName: "field1",
			FileName:  "file1.txt",
			Reader:    strings.NewReader("content1"),
		}).
		WithMultipartField(MultipartField{
			ParamName: "field2",
			FileName:  "file2.txt",
			Reader:    strings.NewReader("content2"),
		}).
		WithHeader("X-Custom", "custom-value").
		WithQuery(map[string]string{"key": "value"}).
		WithAuth(&Authorization{
			Username: "user",
			Password: "pass",
		})

	client := resty.New()
	req := client.R()

	// Act: 应用配置到请求
	result := config.applyToRequest(req)

	// Assert: 验证返回的是同一个请求对象（链式调用）
	assert.Equal(t, req, result)
	// 注意：resty 的内部状态无法直接验证，但可以确认方法执行成功
}

func TestMultipartRequestConfig_applyToRequest_FilePath(t *testing.T) {
	// Arrange: 创建配置，使用 FilePath
	config := NewMultipartRequestConfig().
		WithMultipartField(MultipartField{
			ParamName: "file",
			FileName:  "test.txt",
			FilePath:  "/path/to/file",
			Reader:     strings.NewReader("file content"),
		})

	client := resty.New()
	req := client.R()

	// Act: 应用配置到请求
	result := config.applyToRequest(req)

	// Assert: 验证方法执行成功
	assert.Equal(t, req, result)
}

func TestMultipartRequestConfig_applyToRequest_FormData(t *testing.T) {
	// Arrange: 创建配置，使用普通表单数据（无 Reader）
	config := NewMultipartRequestConfig().
		WithMultipartField(MultipartField{
			ParamName: "field",
			FileName:  "value", // 作为表单值
		})

	client := resty.New()
	req := client.R()

	// Act: 应用配置到请求
	result := config.applyToRequest(req)

	// Assert: 验证方法执行成功
	assert.Equal(t, req, result)
}

func TestMultipartRequestConfig_applyToRequest_EmptyFields(t *testing.T) {
	// Arrange: 创建配置，无 multipart 字段
	config := NewMultipartRequestConfig().
		WithHeader("X-Test", "test-value")

	client := resty.New()
	req := client.R()

	// Act: 应用配置到请求
	result := config.applyToRequest(req)

	// Assert: 验证方法执行成功
	assert.Equal(t, req, result)
}

func TestMultipartRequestConfig_applyToRequest_NoConfig(t *testing.T) {
	// Arrange: 创建空配置
	config := NewMultipartRequestConfig()

	client := resty.New()
	req := client.R()

	// Act: 应用配置到请求
	result := config.applyToRequest(req)

	// Assert: 验证方法执行成功（即使配置为空）
	assert.Equal(t, req, result)
}

// ==================== 链式调用测试 ====================
// 注意：TestMultipartRequestConfig_Chain 已在 config_test.go 中定义

// ==================== MultipartField 结构体测试 ====================

func TestMultipartField_AllFields(t *testing.T) {
	// Arrange & Act: 创建包含所有字段的 MultipartField
	reader := strings.NewReader("test content")
	field := MultipartField{
		ParamName:  "test-param",
		FileName:   "test.txt",
		FilePath:   "/path/to/file",
		ContentType: "text/plain",
		Reader:     reader,
	}

	// Assert: 验证所有字段都已设置
	assert.Equal(t, "test-param", field.ParamName)
	assert.Equal(t, "test.txt", field.FileName)
	assert.Equal(t, "/path/to/file", field.FilePath)
	assert.Equal(t, "text/plain", field.ContentType)
	assert.Equal(t, reader, field.Reader)
}

func TestMultipartField_Minimal(t *testing.T) {
	// Arrange & Act: 创建最小字段的 MultipartField
	field := MultipartField{
		ParamName: "test-param",
		FileName:  "test.txt",
	}

	// Assert: 验证基本字段已设置
	assert.Equal(t, "test-param", field.ParamName)
	assert.Equal(t, "test.txt", field.FileName)
	assert.Empty(t, field.FilePath)
	assert.Empty(t, field.ContentType)
	assert.Nil(t, field.Reader)
}

// ==================== 边界情况测试 ====================

func TestMultipartRequestConfig_WithMultipartField_NilReader(t *testing.T) {
	// Arrange: 创建配置
	config := NewMultipartRequestConfig()

	// Act: 添加无 Reader 的字段
	field := MultipartField{
		ParamName: "field",
		FileName:  "value",
		Reader:    nil,
	}
	config.WithMultipartField(field)

	// Assert: 验证字段已添加（Reader 为 nil 是允许的）
	assert.Len(t, config.MultipartFields, 1)
	assert.Nil(t, config.MultipartFields[0].Reader)
}

func TestMultipartRequestConfig_WithMultipartField_EmptyParamName(t *testing.T) {
	// Arrange: 创建配置
	config := NewMultipartRequestConfig()

	// Act: 添加空 ParamName 的字段
	field := MultipartField{
		ParamName: "",
		FileName:  "value",
	}
	config.WithMultipartField(field)

	// Assert: 验证字段已添加（空 ParamName 由调用方验证）
	assert.Len(t, config.MultipartFields, 1)
	assert.Empty(t, config.MultipartFields[0].ParamName)
}

func TestMultipartRequestConfig_applyToRequest_ReaderOnly(t *testing.T) {
	// Arrange: 创建配置，只有 Reader，无 FilePath
	config := NewMultipartRequestConfig().
		WithMultipartField(MultipartField{
			ParamName: "file",
			FileName:  "test.txt",
			Reader:    strings.NewReader("content"),
		})

	client := resty.New()
	req := client.R()

	// Act: 应用配置到请求
	result := config.applyToRequest(req)

	// Assert: 验证方法执行成功
	assert.Equal(t, req, result)
}

func TestMultipartRequestConfig_applyToRequest_FilePathAndReader(t *testing.T) {
	// Arrange: 创建配置，同时有 FilePath 和 Reader（FilePath 优先）
	config := NewMultipartRequestConfig().
		WithMultipartField(MultipartField{
			ParamName: "file",
			FileName:  "test.txt",
			FilePath:  "/path/to/file",
			Reader:    strings.NewReader("content"),
		})

	client := resty.New()
	req := client.R()

	// Act: 应用配置到请求
	result := config.applyToRequest(req)

	// Assert: 验证方法执行成功
	assert.Equal(t, req, result)
}

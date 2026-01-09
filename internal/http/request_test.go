//go:build test

package http

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// ==================== NewRequestConfig 测试 ====================

func TestNewRequestConfig(t *testing.T) {
	// Act: 创建新的 RequestConfig
	config := NewRequestConfig()

	// Assert: 验证默认值
	assert.NotNil(t, config)
	assert.NotNil(t, config.Headers)
	assert.Nil(t, config.Body)
	assert.Nil(t, config.Query)
	assert.Nil(t, config.Auth)
	assert.Zero(t, config.Timeout)
	assert.Nil(t, config.Retry)
}

// ==================== WithBody 测试 ====================

func TestRequestConfig_WithBody(t *testing.T) {
	// Arrange: 创建配置
	config := NewRequestConfig()
	body := map[string]string{
		"key": "value",
	}

	// Act: 设置请求体
	result := config.WithBody(body)

	// Assert: 验证链式调用和请求体设置
	assert.Equal(t, config, result)
	assert.Equal(t, body, config.Body)
}

func TestRequestConfig_WithBody_Struct(t *testing.T) {
	// Arrange: 创建配置
	config := NewRequestConfig()
	type RequestBody struct {
		Name  string
		Value int
	}
	body := RequestBody{
		Name:  "test",
		Value: 123,
	}

	// Act: 设置请求体（结构体）
	result := config.WithBody(body)

	// Assert: 验证链式调用和请求体设置
	assert.Equal(t, config, result)
	assert.Equal(t, body, config.Body)
}

func TestRequestConfig_WithBody_Nil(t *testing.T) {
	// Arrange: 创建配置并设置请求体
	config := NewRequestConfig()
	config.WithBody(map[string]string{"key": "value"})

	// Act: 设置为 nil
	config.WithBody(nil)

	// Assert: 验证请求体已清除
	assert.Nil(t, config.Body)
}

// ==================== WithQuery 测试 ====================
// 注意：TestRequestConfig_WithQuery 已在 config_test.go 中定义

// ==================== WithAuth 测试 ====================
// 注意：TestRequestConfig_WithAuth 已在 config_test.go 中定义

// ==================== WithHeader 测试 ====================

func TestRequestConfig_WithHeader(t *testing.T) {
	// Arrange: 创建配置
	config := NewRequestConfig()

	// Act: 设置单个 Header
	result := config.WithHeader("X-Custom-Header", "custom-value")

	// Assert: 验证链式调用和 Header 设置
	assert.Equal(t, config, result)
	assert.Equal(t, "custom-value", config.Headers["X-Custom-Header"])
}

// ==================== WithHeaders 测试 ====================
// 注意：TestRequestConfig_WithHeaders 已在 config_test.go 中定义

// ==================== WithTimeout 测试 ====================

func TestRequestConfig_WithTimeout(t *testing.T) {
	// Arrange: 创建配置
	config := NewRequestConfig()
	timeout := 10 * time.Second

	// Act: 设置超时时间
	result := config.WithTimeout(timeout)

	// Assert: 验证链式调用和超时时间设置
	assert.Equal(t, config, result)
	assert.Equal(t, timeout, config.Timeout)
}

// ==================== WithRetry 测试 ====================

func TestRequestConfig_WithRetry(t *testing.T) {
	// Arrange: 创建配置
	config := NewRequestConfig()
	retry := NewRetryConfig().WithRetryCount(3)

	// Act: 设置重试配置
	result := config.WithRetry(retry)

	// Assert: 验证链式调用和重试配置设置
	assert.Equal(t, config, result)
	assert.Equal(t, retry, config.Retry)
}

// ==================== ensureHeaders 测试 ====================

func TestRequestConfig_ensureHeaders(t *testing.T) {
	// Arrange: 创建配置，Headers 为 nil
	config := &RequestConfig{
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

// ==================== convertToQueryParams 测试 ====================

func TestConvertToQueryParams_MapStringString(t *testing.T) {
	// Arrange: map[string]string 类型
	query := map[string]string{
		"key1": "value1",
		"key2": "value2",
		"key3": "value3",
	}

	// Act: 转换查询参数
	result := convertToQueryParams(query)

	// Assert: 验证返回的是同一个 map
	assert.Equal(t, query, result)
	assert.Equal(t, "value1", result["key1"])
	assert.Equal(t, "value2", result["key2"])
	assert.Equal(t, "value3", result["key3"])
}

func TestConvertToQueryParams_MapStringInterface(t *testing.T) {
	// Arrange: map[string]interface{} 类型
	query := map[string]interface{}{
		"key1": "value1",
		"key2": 123,
		"key3": true,
		"key4": 45.67,
	}

	// Act: 转换查询参数
	result := convertToQueryParams(query)

	// Assert: 验证所有值都转换为字符串
	assert.Equal(t, "value1", result["key1"])
	assert.Equal(t, "123", result["key2"])
	assert.Equal(t, "true", result["key3"])
	assert.Equal(t, "45.67", result["key4"])
}

func TestConvertToQueryParams_SliceString(t *testing.T) {
	// Arrange: []string 类型，格式为 "key=value"
	query := []string{
		"key1=value1",
		"key2=value2",
		"key3=value3",
	}

	// Act: 转换查询参数
	result := convertToQueryParams(query)

	// Assert: 验证正确解析
	assert.Equal(t, "value1", result["key1"])
	assert.Equal(t, "value2", result["key2"])
	assert.Equal(t, "value3", result["key3"])
}

func TestConvertToQueryParams_SliceString_InvalidFormat(t *testing.T) {
	// Arrange: []string 类型，包含无效格式
	query := []string{
		"key1=value1",
		"invalid-format",        // 没有等号
		"key2=value2",
		"=no-key",              // 没有键
		"no-value=",            // 没有值
		"key3=value3",
	}

	// Act: 转换查询参数
	result := convertToQueryParams(query)

	// Assert: 验证只解析有效格式
	assert.Equal(t, "value1", result["key1"])
	assert.Equal(t, "value2", result["key2"])
	assert.Equal(t, "value3", result["key3"])
	// 无效格式应该被忽略
	assert.NotContains(t, result, "invalid-format")
}

func TestConvertToQueryParams_SliceString_WithEqualsInValue(t *testing.T) {
	// Arrange: []string 类型，值中包含等号
	query := []string{
		"key1=value=with=equals",
		"key2=normal-value",
	}

	// Act: 转换查询参数（使用 SplitN，只分割第一个等号）
	result := convertToQueryParams(query)

	// Assert: 验证正确解析（只分割第一个等号）
	assert.Equal(t, "value=with=equals", result["key1"])
	assert.Equal(t, "normal-value", result["key2"])
}

func TestConvertToQueryParams_EmptyMap(t *testing.T) {
	// Arrange: 空 map
	query := map[string]string{}

	// Act: 转换查询参数
	result := convertToQueryParams(query)

	// Assert: 验证返回空 map
	assert.Empty(t, result)
}

func TestConvertToQueryParams_EmptySlice(t *testing.T) {
	// Arrange: 空 slice
	query := []string{}

	// Act: 转换查询参数
	result := convertToQueryParams(query)

	// Assert: 验证返回空 map
	assert.Empty(t, result)
}

func TestConvertToQueryParams_UnsupportedType(t *testing.T) {
	// Arrange: 不支持的类型
	query := 123

	// Act: 转换查询参数
	result := convertToQueryParams(query)

	// Assert: 验证返回空 map（不支持的类型）
	assert.Empty(t, result)
}

func TestConvertToQueryParams_StringType(t *testing.T) {
	// Arrange: string 类型（不支持）
	query := "key=value"

	// Act: 转换查询参数
	result := convertToQueryParams(query)

	// Assert: 验证返回空 map（string 类型不支持）
	assert.Empty(t, result)
}

func TestConvertToQueryParams_StructType(t *testing.T) {
	// Arrange: 结构体类型（不支持）
	type QueryStruct struct {
		Key1 string
		Key2 string
	}
	query := QueryStruct{
		Key1: "value1",
		Key2: "value2",
	}

	// Act: 转换查询参数
	result := convertToQueryParams(query)

	// Assert: 验证返回空 map（结构体类型不支持）
	assert.Empty(t, result)
}

func TestConvertToQueryParams_Nil(t *testing.T) {
	// Arrange: nil
	var query interface{} = nil

	// Act: 转换查询参数
	result := convertToQueryParams(query)

	// Assert: 验证返回空 map
	assert.Empty(t, result)
}

func TestConvertToQueryParams_MapStringInterface_Empty(t *testing.T) {
	// Arrange: 空 map[string]interface{}
	query := map[string]interface{}{}

	// Act: 转换查询参数
	result := convertToQueryParams(query)

	// Assert: 验证返回空 map
	assert.Empty(t, result)
}

func TestConvertToQueryParams_MapStringInterface_NilValues(t *testing.T) {
	// Arrange: map[string]interface{} 包含 nil 值
	query := map[string]interface{}{
		"key1": "value1",
		"key2": nil,
		"key3": "value3",
	}

	// Act: 转换查询参数
	result := convertToQueryParams(query)

	// Assert: 验证 nil 值转换为 "nil" 字符串
	assert.Equal(t, "value1", result["key1"])
	assert.Equal(t, "<nil>", result["key2"]) // fmt.Sprintf("%v", nil) 返回 "<nil>"
	assert.Equal(t, "value3", result["key3"])
}

// ==================== 链式调用测试 ====================
// 注意：TestRequestConfig_Chain 已在 config_test.go 中定义

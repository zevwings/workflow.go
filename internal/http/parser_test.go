package http

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ==================== JsonParser 测试 ====================

// TestJsonParser_Parse 测试 JSON 解析器
func TestJsonParser_Parse(t *testing.T) {
	parser := &JsonParser{}

	// 测试正常 JSON
	result, err := parser.Parse([]byte(`{"key": "value"}`), http.StatusOK)
	require.NoError(t, err)
	assert.NotNil(t, result)

	// 测试空响应
	result, err = parser.Parse([]byte(""), http.StatusOK)
	require.NoError(t, err)

	// 测试空白字符
	result, err = parser.Parse([]byte("   "), http.StatusOK)
	require.NoError(t, err)

	// 测试无效 JSON
	_, err = parser.Parse([]byte("invalid json"), http.StatusOK)
	assert.Error(t, err)
}

// TestParseJSON 测试 ParseJSON 函数
func TestParseJSON(t *testing.T) {
	type ResponseData struct {
		Message string `json:"message"`
		ID      int    `json:"id"`
	}

	// 测试正常 JSON
	var data ResponseData
	err := ParseJSON([]byte(`{"message": "success", "id": 123}`), http.StatusOK, &data)
	require.NoError(t, err)
	assert.Equal(t, "success", data.Message)
	assert.Equal(t, 123, data.ID)

	// 测试空响应
	err = ParseJSON([]byte(""), http.StatusOK, &data)
	require.NoError(t, err)

	// 测试无效 JSON
	err = ParseJSON([]byte("invalid json"), http.StatusOK, &data)
	assert.Error(t, err)
}

// ==================== TextParser 测试 ====================

// TestTextParser_Parse 测试文本解析器
func TestTextParser_Parse(t *testing.T) {
	parser := &TextParser{}

	// 测试正常文本
	result, err := parser.Parse([]byte("Hello, World!"), http.StatusOK)
	require.NoError(t, err)
	assert.Equal(t, "Hello, World!", result)

	// 测试错误状态码
	_, err = parser.Parse([]byte("Error"), http.StatusBadRequest)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "HTTP request failed")
}

// TestParseText 测试 ParseText 函数
func TestParseText(t *testing.T) {
	// 测试正常文本
	text, err := ParseText([]byte("Hello, World!"), http.StatusOK)
	require.NoError(t, err)
	assert.Equal(t, "Hello, World!", text)

	// 测试错误状态码
	_, err = ParseText([]byte("Error"), http.StatusBadRequest)
	assert.Error(t, err)
}


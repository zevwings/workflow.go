//go:build test

package http

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zevwings/workflow/internal/testutils"
)

// ==================== HttpResponse 基础测试 ====================

// TestHttpResponse_AsJSON 测试 JSON 响应解析
func TestHttpResponse_AsJSON(t *testing.T) {
	type ResponseData struct {
		Message string `json:"message"`
		ID      int    `json:"id"`
	}

	server := testutils.NewHTTPTestServer().
		WithContentType("application/json").
		WithStatus(http.StatusOK).
		WithJSONBody(map[string]interface{}{
			"message": "success",
			"id":      123,
		}).
		Build(t)

	client := NewClient()
	resp, err := client.GetWithConfig(server.URL(), nil)
	require.NoError(t, err)

	var data ResponseData
	data, err = AsJSON[ResponseData](resp)
	require.NoError(t, err)
	assert.Equal(t, "success", data.Message)
	assert.Equal(t, 123, data.ID)
}

// TestHttpResponse_AsText 测试文本响应解析
func TestHttpResponse_AsText(t *testing.T) {
	server := testutils.NewHTTPTestServer().
		WithStatus(http.StatusOK).
		WithStringBody("Hello, World!").
		Build(t)

	client := NewClient()
	resp, err := client.GetWithConfig(server.URL(), nil)
	require.NoError(t, err)

	text, err := resp.AsText()
	require.NoError(t, err)
	assert.Equal(t, "Hello, World!", text)
}

// TestHttpResponse_AsBytes 测试字节响应解析
func TestHttpResponse_AsBytes(t *testing.T) {
	expectedBytes := []byte("test data")
	server := testutils.NewHTTPTestServer().
		WithStatus(http.StatusOK).
		WithBody(expectedBytes).
		Build(t)

	client := NewClient()
	resp, err := client.GetWithConfig(server.URL(), nil)
	require.NoError(t, err)

	bytes := resp.AsBytes()
	assert.Equal(t, expectedBytes, bytes)
}

// TestHttpResponse_IsSuccess 测试成功响应判断
func TestHttpResponse_IsSuccess(t *testing.T) {
	testCases := []struct {
		name       string
		statusCode int
		expected   bool
	}{
		{"200 OK", http.StatusOK, true},
		{"201 Created", http.StatusCreated, true},
		{"204 No Content", http.StatusNoContent, true},
		{"299 Success", 299, true},
		{"400 Bad Request", http.StatusBadRequest, false},
		{"500 Internal Server Error", http.StatusInternalServerError, false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			server := testutils.NewHTTPTestServer().
				WithStatus(tc.statusCode).
				Build(t)

			client := NewClient()
			resp, err := client.GetWithConfig(server.URL(), nil)
			require.NoError(t, err)
			assert.Equal(t, tc.expected, resp.IsSuccess())
			assert.Equal(t, !tc.expected, resp.IsError())
		})
	}
}

// TestHttpResponse_EnsureSuccess 测试成功响应检查
func TestHttpResponse_EnsureSuccess(t *testing.T) {
	server := testutils.NewHTTPTestServer().
		WithStatus(http.StatusOK).
		WithJSONBody(map[string]string{"message": "success"}).
		Build(t)

	client := NewClient()
	resp, err := client.GetWithConfig(server.URL(), nil)
	require.NoError(t, err)

	successResp, err := resp.EnsureSuccess()
	require.NoError(t, err)
	assert.Equal(t, resp, successResp)
}

// TestHttpResponse_EnsureSuccess_Error 测试错误响应检查
func TestHttpResponse_EnsureSuccess_Error(t *testing.T) {
	server := testutils.NewHTTPTestServer().
		WithStatus(http.StatusBadRequest).
		WithJSONBody(map[string]string{"error": "Bad Request"}).
		Build(t)

	client := NewClient()
	resp, err := client.GetWithConfig(server.URL(), nil)
	require.NoError(t, err)

	_, err = resp.EnsureSuccess()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "HTTP request failed")
}

// TestHttpResponse_EnsureSuccessWith 测试自定义错误处理器
func TestHttpResponse_EnsureSuccessWith(t *testing.T) {
	// 测试成功响应
	server := testutils.NewHTTPTestServer().
		WithStatus(http.StatusOK).
		WithJSONBody(map[string]string{"message": "success"}).
		Build(t)

	client := NewClient()
	resp, err := client.GetWithConfig(server.URL(), nil)
	require.NoError(t, err)

	customError := func(r *HttpResponse) error {
		return &ConfigError{Message: "custom error"}
	}

	successResp, err := resp.EnsureSuccessWith(customError)
	require.NoError(t, err)
	assert.Equal(t, resp, successResp)

	// 测试错误响应
	server2 := testutils.NewHTTPTestServer().
		WithStatus(http.StatusBadRequest).
		WithJSONBody(map[string]string{"error": "Bad Request"}).
		Build(t)

	resp2, err := client.GetWithConfig(server2.URL(), nil)
	require.NoError(t, err)

	_, err = resp2.EnsureSuccessWith(customError)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "custom error")
}

// TestHttpResponse_ExtractErrorMessage 测试错误消息提取
func TestHttpResponse_ExtractErrorMessage(t *testing.T) {
	testCases := []struct {
		name           string
		statusCode     int
		body           string
		expectedSubstr string
	}{
		{
			name:           "JSON error with error.message",
			statusCode:     http.StatusBadRequest,
			body:           `{"error": {"message": "Invalid input"}}`,
			expectedSubstr: "Invalid input",
		},
		{
			name:           "JSON error with error string",
			statusCode:     http.StatusBadRequest,
			body:           `{"error": "Invalid input"}`,
			expectedSubstr: "Invalid input",
		},
		{
			name:           "JSON error with message",
			statusCode:     http.StatusBadRequest,
			body:           `{"message": "Invalid input"}`,
			expectedSubstr: "Invalid input",
		},
		{
			name:           "Plain text error",
			statusCode:     http.StatusBadRequest,
			body:           "Invalid input",
			expectedSubstr: "Invalid input",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			server := testutils.NewHTTPTestServer().
				WithStatus(tc.statusCode).
				WithStringBody(tc.body).
				Build(t)

			client := NewClient()
			resp, err := client.GetWithConfig(server.URL(), nil)
			require.NoError(t, err)

			errorMsg := resp.ExtractErrorMessage()
			assert.Contains(t, errorMsg, tc.expectedSubstr)
		})
	}
}

// TestHttpResponse_GetHeader 测试获取 Header
func TestHttpResponse_GetHeader(t *testing.T) {
	server := testutils.NewHTTPTestServer().
		WithHeader("X-Custom-Header", "custom-value").
		WithContentType("application/json").
		WithStatus(http.StatusOK).
		WithJSONBody(map[string]string{"message": "success"}).
		Build(t)

	client := NewClient()
	resp, err := client.GetWithConfig(server.URL(), nil)
	require.NoError(t, err)

	value, ok := resp.GetHeader("X-Custom-Header")
	assert.True(t, ok)
	assert.Equal(t, "custom-value", value)

	value, ok = resp.GetHeader("x-custom-header") // 不区分大小写
	assert.True(t, ok)
	assert.Equal(t, "custom-value", value)

	_, ok = resp.GetHeader("Non-Existent-Header")
	assert.False(t, ok)
}

// TestHttpResponse_ParseWith 测试使用自定义解析器
func TestHttpResponse_ParseWith(t *testing.T) {
	server := testutils.NewHTTPTestServer().
		WithStatus(http.StatusOK).
		WithJSONBody(map[string]string{"message": "success"}).
		Build(t)

	client := NewClient()
	resp, err := client.GetWithConfig(server.URL(), nil)
	require.NoError(t, err)

	// 使用 JSON 解析器
	parser := &JsonParser{}
	result, err := resp.ParseWith(parser)
	require.NoError(t, err)
	assert.NotNil(t, result)

	// 使用文本解析器
	textParser := &TextParser{}
	result, err = resp.ParseWith(textParser)
	require.NoError(t, err)
	// JSON 序列化可能压缩空格，所以只检查包含关键内容
	assert.Contains(t, result, "message")
	assert.Contains(t, result, "success")
}

// ==================== 边界情况测试 ====================

// TestHttpResponse_EmptyBody 测试空响应体
func TestHttpResponse_EmptyBody(t *testing.T) {
	server := testutils.NewHTTPTestServer().
		WithStatus(http.StatusNoContent).
		Build(t)

	client := NewClient()
	resp, err := client.GetWithConfig(server.URL(), nil)
	require.NoError(t, err)

	text, err := resp.AsText()
	require.NoError(t, err)
	assert.Empty(t, text)

	bytes := resp.AsBytes()
	assert.Empty(t, bytes)
}

// TestHttpResponse_JSONEmptyBody 测试空 JSON 响应体
func TestHttpResponse_JSONEmptyBody(t *testing.T) {
	type ResponseData struct {
		Message string `json:"message"`
	}

	server := testutils.NewHTTPTestServer().
		WithStatus(http.StatusOK).
		WithStringBody("").
		Build(t)

	client := NewClient()
	resp, err := client.GetWithConfig(server.URL(), nil)
	require.NoError(t, err)

	_, err = AsJSON[ResponseData](resp)
	require.NoError(t, err)
}

// TestHttpResponse_InvalidJSON 测试无效 JSON
func TestHttpResponse_InvalidJSON(t *testing.T) {
	type ResponseData struct {
		Message string `json:"message"`
	}

	server := testutils.NewHTTPTestServer().
		WithStatus(http.StatusOK).
		WithStringBody("invalid json").
		Build(t)

	client := NewClient()
	resp, err := client.GetWithConfig(server.URL(), nil)
	require.NoError(t, err)

	_, err = AsJSON[ResponseData](resp)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to parse JSON")
}

//go:build test

package testutils

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestNewHTTPTestServer 测试创建 HTTP 测试服务器构建器
func TestNewHTTPTestServer(t *testing.T) {
	builder := NewHTTPTestServer()
	assert.NotNil(t, builder)
}

// TestHTTPTestServerBuilder_Build 测试构建服务器
func TestHTTPTestServerBuilder_Build(t *testing.T) {
	server := NewHTTPTestServer().
		WithStatus(http.StatusOK).
		WithStringBody("test response").
		Build(t)

	assert.NotNil(t, server)
	assert.NotEmpty(t, server.URL())
	assert.NotNil(t, server.Server())
}

// TestHTTPTestServerBuilder_WithStatus 测试设置状态码
func TestHTTPTestServerBuilder_WithStatus(t *testing.T) {
	server := NewHTTPTestServer().
		WithStatus(http.StatusNotFound).
		Build(t)

	resp, err := http.Get(server.URL())
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusNotFound, resp.StatusCode)
}

// TestHTTPTestServerBuilder_WithStringBody 测试设置字符串响应体
func TestHTTPTestServerBuilder_WithStringBody(t *testing.T) {
	server := NewHTTPTestServer().
		WithStringBody("hello world").
		Build(t)

	resp, err := http.Get(server.URL())
	require.NoError(t, err)
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	require.NoError(t, err)
	assert.Equal(t, "hello world", string(body))
}

// TestHTTPTestServerBuilder_WithBody 测试设置字节响应体
func TestHTTPTestServerBuilder_WithBody(t *testing.T) {
	bodyData := []byte("test bytes")
	server := NewHTTPTestServer().
		WithBody(bodyData).
		Build(t)

	resp, err := http.Get(server.URL())
	require.NoError(t, err)
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	require.NoError(t, err)
	assert.Equal(t, bodyData, body)
}

// TestHTTPTestServerBuilder_WithJSONBody 测试设置 JSON 响应体
func TestHTTPTestServerBuilder_WithJSONBody(t *testing.T) {
	data := map[string]interface{}{
		"message": "success",
		"id":      123,
	}

	server := NewHTTPTestServer().
		WithJSONBody(data).
		Build(t)

	resp, err := http.Get(server.URL())
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, "application/json", resp.Header.Get("Content-Type"))

	var result map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&result)
	require.NoError(t, err)
	assert.Equal(t, "success", result["message"])
	assert.Equal(t, float64(123), result["id"]) // JSON 数字解析为 float64
}

// TestHTTPTestServerBuilder_WithHeader 测试设置响应头
func TestHTTPTestServerBuilder_WithHeader(t *testing.T) {
	server := NewHTTPTestServer().
		WithHeader("X-Custom-Header", "custom-value").
		WithHeader("Authorization", "Bearer token123").
		Build(t)

	resp, err := http.Get(server.URL())
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, "custom-value", resp.Header.Get("X-Custom-Header"))
	assert.Equal(t, "Bearer token123", resp.Header.Get("Authorization"))
}

// TestHTTPTestServerBuilder_WithContentType 测试设置 Content-Type
func TestHTTPTestServerBuilder_WithContentType(t *testing.T) {
	server := NewHTTPTestServer().
		WithContentType("application/xml").
		Build(t)

	resp, err := http.Get(server.URL())
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, "application/xml", resp.Header.Get("Content-Type"))
}

// TestHTTPTestServerBuilder_WithHandler 测试自定义处理函数
func TestHTTPTestServerBuilder_WithHandler(t *testing.T) {
	customHandler := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Custom", "handler")
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte("custom handler response"))
	}

	server := NewHTTPTestServer().
		WithHandler(customHandler).
		Build(t)

	resp, err := http.Get(server.URL())
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusCreated, resp.StatusCode)
	assert.Equal(t, "handler", resp.Header.Get("X-Custom"))

	body, err := io.ReadAll(resp.Body)
	require.NoError(t, err)
	assert.Equal(t, "custom handler response", string(body))
}

// TestHTTPTestServerBuilder_WithMethodCheck 测试请求方法验证
func TestHTTPTestServerBuilder_WithMethodCheck(t *testing.T) {
	server := NewHTTPTestServer().
		WithMethodCheck(http.MethodPost).
		WithStatus(http.StatusOK).
		Build(t)

	// 使用 POST 方法（应该成功）
	req, err := http.NewRequest(http.MethodPost, server.URL(), nil)
	require.NoError(t, err)

	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

// TestHTTPTestServerBuilder_WithPathCheck 测试请求路径验证
func TestHTTPTestServerBuilder_WithPathCheck(t *testing.T) {
	server := NewHTTPTestServer().
		WithPathCheck("/api/test").
		WithStatus(http.StatusOK).
		Build(t)

	req, err := http.NewRequest(http.MethodGet, server.URL()+"/api/test", nil)
	require.NoError(t, err)

	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

// TestHTTPTestServerBuilder_WithRequestCheck 测试自定义请求验证
func TestHTTPTestServerBuilder_WithRequestCheck(t *testing.T) {
	server := NewHTTPTestServer().
		WithRequestCheck(func(t *testing.T, r *http.Request) {
			assert.Equal(t, http.MethodPut, r.Method)
			assert.Equal(t, "/api/update", r.URL.Path)
			assert.Equal(t, "Bearer token", r.Header.Get("Authorization"))
		}).
		WithStatus(http.StatusOK).
		Build(t)

	req, err := http.NewRequest(http.MethodPut, server.URL()+"/api/update", nil)
	require.NoError(t, err)
	req.Header.Set("Authorization", "Bearer token")

	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

// TestReadRequestBody 测试读取请求体
func TestReadRequestBody(t *testing.T) {
	requestData := map[string]interface{}{
		"name":  "test",
		"value": 42,
	}

	server := NewHTTPTestServer().
		WithHandler(func(w http.ResponseWriter, r *http.Request) {
			var data map[string]interface{}
			ReadRequestBody(t, r, &data)
			assert.Equal(t, "test", data["name"])
			assert.Equal(t, float64(42), data["value"])
			w.WriteHeader(http.StatusOK)
		}).
		Build(t)

	jsonData, err := json.Marshal(requestData)
	require.NoError(t, err)

	req, err := http.NewRequest(http.MethodPost, server.URL(), strings.NewReader(string(jsonData)))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

// TestAssertRequestMethod 测试断言请求方法
func TestAssertRequestMethod(t *testing.T) {
	server := NewHTTPTestServer().
		WithHandler(func(w http.ResponseWriter, r *http.Request) {
			AssertRequestMethod(t, r, http.MethodDelete)
			w.WriteHeader(http.StatusOK)
		}).
		Build(t)

	req, err := http.NewRequest(http.MethodDelete, server.URL(), nil)
	require.NoError(t, err)

	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

// TestAssertRequestPath 测试断言请求路径
func TestAssertRequestPath(t *testing.T) {
	server := NewHTTPTestServer().
		WithHandler(func(w http.ResponseWriter, r *http.Request) {
			AssertRequestPath(t, r, "/api/users")
			w.WriteHeader(http.StatusOK)
		}).
		Build(t)

	req, err := http.NewRequest(http.MethodGet, server.URL()+"/api/users", nil)
	require.NoError(t, err)

	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

// TestAssertRequestHeader 测试断言请求头
func TestAssertRequestHeader(t *testing.T) {
	server := NewHTTPTestServer().
		WithHandler(func(w http.ResponseWriter, r *http.Request) {
			AssertRequestHeader(t, r, "X-API-Key", "secret-key")
			w.WriteHeader(http.StatusOK)
		}).
		Build(t)

	req, err := http.NewRequest(http.MethodGet, server.URL(), nil)
	require.NoError(t, err)
	req.Header.Set("X-API-Key", "secret-key")

	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

// TestAssertRequestQuery 测试断言查询参数
func TestAssertRequestQuery(t *testing.T) {
	server := NewHTTPTestServer().
		WithHandler(func(w http.ResponseWriter, r *http.Request) {
			AssertRequestQuery(t, r, "page", "1")
			AssertRequestQuery(t, r, "limit", "10")
			w.WriteHeader(http.StatusOK)
		}).
		Build(t)

	req, err := http.NewRequest(http.MethodGet, server.URL()+"?page=1&limit=10", nil)
	require.NoError(t, err)

	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}


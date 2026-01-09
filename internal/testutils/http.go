//go:build test

package testutils

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// HTTPTestServer 封装了 httptest.Server，提供便捷的访问方法
type HTTPTestServer struct {
	server *httptest.Server
	url    string
}

// URL 返回服务器的 URL
func (s *HTTPTestServer) URL() string {
	return s.url
}

// Server 返回底层的 httptest.Server（用于高级用法）
func (s *HTTPTestServer) Server() *httptest.Server {
	return s.server
}

// Close 关闭服务器（通常不需要手动调用，Build 时会自动注册清理）
func (s *HTTPTestServer) Close() {
	s.server.Close()
}

// HTTPTestServerBuilder 用于构建 HTTP 测试服务器的构建器
type HTTPTestServerBuilder struct {
	handler      http.HandlerFunc
	status       int
	body         []byte
	headers      map[string]string
	contentType  string
	requestCheck func(*testing.T, *http.Request) // 用于验证请求
}

// NewHTTPTestServer 创建新的 HTTP 测试服务器构建器
func NewHTTPTestServer() *HTTPTestServerBuilder {
	return &HTTPTestServerBuilder{
		status:  http.StatusOK,
		headers: make(map[string]string),
	}
}

// WithHandler 设置自定义请求处理函数
// 如果设置了 handler，其他响应配置（status, body, headers）将被忽略
func (b *HTTPTestServerBuilder) WithHandler(handler http.HandlerFunc) *HTTPTestServerBuilder {
	b.handler = handler
	return b
}

// WithStatus 设置响应状态码
func (b *HTTPTestServerBuilder) WithStatus(status int) *HTTPTestServerBuilder {
	b.status = status
	return b
}

// WithBody 设置响应体（原始字节）
func (b *HTTPTestServerBuilder) WithBody(body []byte) *HTTPTestServerBuilder {
	b.body = body
	return b
}

// WithStringBody 设置响应体（字符串）
func (b *HTTPTestServerBuilder) WithStringBody(body string) *HTTPTestServerBuilder {
	b.body = []byte(body)
	return b
}

// WithJSONBody 设置 JSON 响应体
// 会自动设置 Content-Type 为 application/json
func (b *HTTPTestServerBuilder) WithJSONBody(data interface{}) *HTTPTestServerBuilder {
	jsonData, err := json.Marshal(data)
	if err != nil {
		// 在构建时无法处理错误，延迟到 Build 时处理
		panic(fmt.Errorf("JSON serialization failed: %w", err))
	}
	b.body = jsonData
	b.contentType = "application/json"
	return b
}

// WithHeader 设置响应头
func (b *HTTPTestServerBuilder) WithHeader(key, value string) *HTTPTestServerBuilder {
	b.headers[key] = value
	return b
}

// WithContentType 设置 Content-Type 响应头
func (b *HTTPTestServerBuilder) WithContentType(contentType string) *HTTPTestServerBuilder {
	b.contentType = contentType
	return b
}

// WithRequestCheck 设置请求验证函数
// 用于验证请求的方法、路径、头部、体等
func (b *HTTPTestServerBuilder) WithRequestCheck(check func(*testing.T, *http.Request)) *HTTPTestServerBuilder {
	b.requestCheck = check
	return b
}

// WithMethodCheck 便捷方法：验证请求方法
func (b *HTTPTestServerBuilder) WithMethodCheck(expectedMethod string) *HTTPTestServerBuilder {
	b.requestCheck = func(t *testing.T, r *http.Request) {
		assert.Equal(t, expectedMethod, r.Method, "请求方法不匹配")
	}
	return b
}

// WithPathCheck 便捷方法：验证请求路径
func (b *HTTPTestServerBuilder) WithPathCheck(expectedPath string) *HTTPTestServerBuilder {
	originalCheck := b.requestCheck
	b.requestCheck = func(t *testing.T, r *http.Request) {
		if originalCheck != nil {
			originalCheck(t, r)
		}
		assert.Equal(t, expectedPath, r.URL.Path, "请求路径不匹配")
	}
	return b
}

// Build 构建并启动 HTTP 测试服务器
// 服务器会在测试结束时自动关闭（使用 t.Cleanup）
func (b *HTTPTestServerBuilder) Build(t *testing.T) *HTTPTestServer {
	t.Helper()

	handler := b.handler
	if handler == nil {
		// 使用默认 handler
		handler = func(w http.ResponseWriter, r *http.Request) {
			// 执行请求验证（如果有）
			if b.requestCheck != nil {
				b.requestCheck(t, r)
			}

			// 设置响应头
			if b.contentType != "" {
				w.Header().Set("Content-Type", b.contentType)
			}
			for k, v := range b.headers {
				w.Header().Set(k, v)
			}

			// 写入状态码和响应体
			w.WriteHeader(b.status)
			if len(b.body) > 0 {
				w.Write(b.body)
			}
		}
	} else if b.requestCheck != nil {
		// 如果同时设置了自定义 handler 和请求验证，包装 handler
		originalHandler := handler
		handler = func(w http.ResponseWriter, r *http.Request) {
			b.requestCheck(t, r)
			originalHandler(w, r)
		}
	}

	server := httptest.NewServer(handler)
	t.Cleanup(func() { server.Close() })

	return &HTTPTestServer{
		server: server,
		url:    server.URL,
	}
}

// ReadRequestBody 便捷函数：读取并解析请求体为 JSON
func ReadRequestBody(t *testing.T, r *http.Request, v interface{}) {
	t.Helper()
	body, err := io.ReadAll(r.Body)
	require.NoError(t, err, "读取请求体失败")
	err = json.Unmarshal(body, v)
	require.NoError(t, err, "解析请求体 JSON 失败")
}

// AssertRequestMethod 便捷函数：断言请求方法
func AssertRequestMethod(t *testing.T, r *http.Request, expectedMethod string) {
	t.Helper()
	assert.Equal(t, expectedMethod, r.Method, "请求方法不匹配")
}

// AssertRequestPath 便捷函数：断言请求路径
func AssertRequestPath(t *testing.T, r *http.Request, expectedPath string) {
	t.Helper()
	assert.Equal(t, expectedPath, r.URL.Path, "请求路径不匹配")
}

// AssertRequestHeader 便捷函数：断言请求头
func AssertRequestHeader(t *testing.T, r *http.Request, key, expectedValue string) {
	t.Helper()
	actualValue := r.Header.Get(key)
	assert.Equal(t, expectedValue, actualValue, "请求头 %s 不匹配", key)
}

// AssertRequestQuery 便捷函数：断言查询参数
func AssertRequestQuery(t *testing.T, r *http.Request, key, expectedValue string) {
	t.Helper()
	actualValue := r.URL.Query().Get(key)
	assert.Equal(t, expectedValue, actualValue, "查询参数 %s 不匹配", key)
}

// NewJSONServer 创建返回 JSON 响应的测试服务器
// 这是对 NewHTTPTestServer 的便捷包装，用于常见场景
//
// 参数:
//   - t: 测试对象
//   - method: HTTP 方法（用于验证请求方法）
//   - status: HTTP 状态码
//   - body: JSON 响应体（可以是任何可序列化的对象）
//
// 返回:
//   - 配置好的 HTTPTestServer
//
// 示例:
//
//	server := testutils.NewJSONServer(t, http.MethodGet, http.StatusOK, map[string]string{"message": "success"})
//	client := http.NewClient()
//	resp, err := client.Get(server.URL())
func NewJSONServer(t *testing.T, method string, status int, body interface{}) *HTTPTestServer {
	t.Helper()
	return NewHTTPTestServer().
		WithMethodCheck(method).
		WithStatus(status).
		WithJSONBody(body).
		Build(t)
}

// NewStringServer 创建返回字符串响应的测试服务器
//
// 参数:
//   - t: 测试对象
//   - method: HTTP 方法（用于验证请求方法）
//   - status: HTTP 状态码
//   - body: 字符串响应体
//
// 返回:
//   - 配置好的 HTTPTestServer
//
// 示例:
//
//	server := testutils.NewStringServer(t, http.MethodGet, http.StatusOK, "hello world")
func NewStringServer(t *testing.T, method string, status int, body string) *HTTPTestServer {
	t.Helper()
	return NewHTTPTestServer().
		WithMethodCheck(method).
		WithStatus(status).
		WithStringBody(body).
		Build(t)
}

// NewEmptyServer 创建空响应的测试服务器（仅状态码）
//
// 参数:
//   - t: 测试对象
//   - method: HTTP 方法（用于验证请求方法）
//   - status: HTTP 状态码
//
// 返回:
//   - 配置好的 HTTPTestServer
//
// 示例:
//
//	server := testutils.NewEmptyServer(t, http.MethodDelete, http.StatusNoContent)
func NewEmptyServer(t *testing.T, method string, status int) *HTTPTestServer {
	t.Helper()
	return NewHTTPTestServer().
		WithMethodCheck(method).
		WithStatus(status).
		Build(t)
}

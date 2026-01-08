package http

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ==================== 客户端基础功能测试 ====================

// TestNewClient 测试创建新的 HTTP 客户端
func TestNewClient(t *testing.T) {
	client := NewClient()
	assert.NotNil(t, client)
	assert.NotNil(t, client.client)
}

// TestGlobal 测试全局客户端单例
func TestGlobal(t *testing.T) {
	client1 := Global()
	client2 := Global()

	// 应该返回同一个实例
	assert.Equal(t, client1, client2)
}

// TestClient_GetRestyClient 测试获取底层 resty 客户端
func TestClient_GetRestyClient(t *testing.T) {
	client := NewClient()
	restyClient := client.GetRestyClient()
	assert.NotNil(t, restyClient)
}

// ==================== HTTP 方法测试（旧版 API）====================

// TestClient_Get 测试 GET 请求（旧版 API）
func TestClient_Get(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"message": "success"}`))
	}))
	defer server.Close()

	client := NewClient()
	resp, err := client.Get(server.URL)

	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode())
}

// TestClient_Post 测试 POST 请求（旧版 API）
func TestClient_Post(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)

		var body map[string]interface{}
		json.NewDecoder(r.Body).Decode(&body)
		assert.Equal(t, "test", body["key"])

		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(`{"id": 123}`))
	}))
	defer server.Close()

	client := NewClient()
	resp, err := client.Post(server.URL, map[string]string{"key": "test"})

	require.NoError(t, err)
	assert.Equal(t, http.StatusCreated, resp.StatusCode())
}

// TestClient_Put 测试 PUT 请求（旧版 API）
func TestClient_Put(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPut, r.Method)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"updated": true}`))
	}))
	defer server.Close()

	client := NewClient()
	resp, err := client.Put(server.URL, map[string]string{"key": "value"})

	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode())
}

// TestClient_Delete 测试 DELETE 请求（旧版 API）
func TestClient_Delete(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodDelete, r.Method)
		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	client := NewClient()
	resp, err := client.Delete(server.URL)

	require.NoError(t, err)
	assert.Equal(t, http.StatusNoContent, resp.StatusCode())
}

// TestClient_Patch 测试 PATCH 请求（旧版 API）
func TestClient_Patch(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPatch, r.Method)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"patched": true}`))
	}))
	defer server.Close()

	client := NewClient()
	resp, err := client.Patch(server.URL, map[string]string{"key": "value"})

	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode())
}

// ==================== HTTP 方法测试（新版 API）====================

// TestClient_GetWithConfig 测试 GET 请求（新版 API）
func TestClient_GetWithConfig(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Equal(t, "test-value", r.Header.Get("X-Test-Header"))
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"message": "success"}`))
	}))
	defer server.Close()

	client := NewClient()
	config := NewRequestConfig().
		WithHeader("X-Test-Header", "test-value")

	resp, err := client.GetWithConfig(server.URL, config)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.Status)
	assert.True(t, resp.IsSuccess())
}

// TestClient_PostWithConfig 测试 POST 请求（新版 API）
func TestClient_PostWithConfig(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)

		var body map[string]interface{}
		json.NewDecoder(r.Body).Decode(&body)
		assert.Equal(t, "test", body["key"])

		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(`{"id": 123}`))
	}))
	defer server.Close()

	client := NewClient()
	config := NewRequestConfig().
		WithBody(map[string]string{"key": "test"})

	resp, err := client.PostWithConfig(server.URL, config)
	require.NoError(t, err)
	assert.Equal(t, http.StatusCreated, resp.Status)
}

// TestClient_PutWithConfig 测试 PUT 请求（新版 API）
func TestClient_PutWithConfig(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPut, r.Method)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"updated": true}`))
	}))
	defer server.Close()

	client := NewClient()
	config := NewRequestConfig().
		WithBody(map[string]string{"key": "value"})

	resp, err := client.PutWithConfig(server.URL, config)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.Status)
}

// TestClient_DeleteWithConfig 测试 DELETE 请求（新版 API）
func TestClient_DeleteWithConfig(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodDelete, r.Method)
		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	client := NewClient()
	resp, err := client.DeleteWithConfig(server.URL, nil)
	require.NoError(t, err)
	assert.Equal(t, http.StatusNoContent, resp.Status)
}

// TestClient_PatchWithConfig 测试 PATCH 请求（新版 API）
func TestClient_PatchWithConfig(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPatch, r.Method)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"patched": true}`))
	}))
	defer server.Close()

	client := NewClient()
	config := NewRequestConfig().
		WithBody(map[string]string{"key": "value"})

	resp, err := client.PatchWithConfig(server.URL, config)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.Status)
}

// ==================== 认证测试 ====================

// TestClient_SetAuth 测试设置认证 Token
func TestClient_SetAuth(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		auth := r.Header.Get("Authorization")
		assert.Equal(t, "Bearer test-token", auth)
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client := NewClient()
	client.SetAuth("test-token")

	resp, err := client.Get(server.URL)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode())
}

// TestClient_SetBasicAuth 测试设置 Basic Auth
func TestClient_SetBasicAuth(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		username, password, ok := r.BasicAuth()
		assert.True(t, ok)
		assert.Equal(t, "user", username)
		assert.Equal(t, "pass", password)
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client := NewClient()
	client.SetBasicAuth("user", "pass")

	resp, err := client.Get(server.URL)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode())
}

// TestClient_SetProxy 测试设置代理
func TestClient_SetProxy(t *testing.T) {
	client := NewClient()
	// 设置代理（不会实际测试代理连接，只测试设置是否成功）
	client.SetProxy("http://proxy.example.com:8080")

	// 验证客户端已创建（代理设置是内部状态，无法直接验证）
	assert.NotNil(t, client.client)
}

// ==================== 流式请求测试 ====================

// TestClient_Stream 测试流式请求
func TestClient_Stream(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("streaming data"))
	}))
	defer server.Close()

	client := NewClient()
	stream, err := client.Stream(MethodGet, server.URL, nil)
	require.NoError(t, err)
	defer stream.Close()

	// 读取流数据
	buf := make([]byte, 1024)
	n, err := stream.Read(buf)
	// EOF 是正常的，表示读取完成
	if err != nil && err.Error() != "EOF" {
		require.NoError(t, err)
	}
	if n > 0 {
		assert.Contains(t, string(buf[:n]), "streaming data")
	}
}

// TestClient_Stream_AllMethods 测试所有 HTTP 方法的流式请求
func TestClient_Stream_AllMethods(t *testing.T) {
	methods := []HttpMethod{MethodGet, MethodPost, MethodPut, MethodDelete, MethodPatch}

	for _, method := range methods {
		t.Run(string(method), func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, string(method), r.Method)
				w.WriteHeader(http.StatusOK)
				w.Write([]byte("stream data"))
			}))
			defer server.Close()

			client := NewClient()
			stream, err := client.Stream(method, server.URL, nil)
			require.NoError(t, err)
			defer stream.Close()

			buf := make([]byte, 1024)
			n, err := stream.Read(buf)
			// EOF 是正常的，表示读取完成
			if err != nil && err.Error() != "EOF" {
				require.NoError(t, err)
			}
			if n > 0 {
				assert.Contains(t, string(buf[:n]), "stream data")
			}
		})
	}
}

// ==================== Multipart 请求测试 ====================

// TestClient_PostMultipart 测试 Multipart 请求
func TestClient_PostMultipart(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// resty 使用 SetFileReader 时会使用 multipart/form-data
		err := r.ParseMultipartForm(10 << 20) // 10MB
		require.NoError(t, err)

		// 检查是否有文件上传
		file, header, err := r.FormFile("test-field")
		if err == nil {
			defer file.Close()
			assert.Equal(t, "test.txt", header.Filename)
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"message": "uploaded"}`))
	}))
	defer server.Close()

	client := NewClient()
	// 对于普通字段，使用 SetFormData（resty 会自动处理）
	// 这里我们使用 Reader 来确保是 multipart 请求
	config := NewMultipartRequestConfig().
		WithMultipartField(MultipartField{
			ParamName: "test-field",
			FileName:  "test.txt",
			Reader:    strings.NewReader("test-value"),
		})

	resp, err := client.PostMultipart(server.URL, config)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.Status)
}

// TestClient_PostMultipart_File 测试 Multipart 文件上传
func TestClient_PostMultipart_File(t *testing.T) {
	fileContent := "test file content"
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Contains(t, r.Header.Get("Content-Type"), "multipart/form-data")

		err := r.ParseMultipartForm(10 << 20)
		require.NoError(t, err)

		file, header, err := r.FormFile("file")
		require.NoError(t, err)
		defer file.Close()

		assert.Equal(t, "test.txt", header.Filename)

		buf := make([]byte, len(fileContent))
		file.Read(buf)
		assert.Equal(t, fileContent, string(buf))

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"message": "file uploaded"}`))
	}))
	defer server.Close()

	client := NewClient()
	config := NewMultipartRequestConfig().
		WithMultipartField(MultipartField{
			ParamName: "file",
			FileName:  "test.txt",
			Reader:    strings.NewReader(fileContent),
		})

	resp, err := client.PostMultipart(server.URL, config)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.Status)
}

// TestClient_PostMultipart_Error 测试 Multipart 配置错误
func TestClient_PostMultipart_Error(t *testing.T) {
	client := NewClient()
	_, err := client.PostMultipart("http://example.com", nil)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "MultipartRequestConfig is required")
}

// ==================== HTTP 状态码测试 ====================

// TestClient_StatusCodes 测试各种 HTTP 状态码
func TestClient_StatusCodes(t *testing.T) {
	testCases := []struct {
		name       string
		statusCode int
		isSuccess  bool
	}{
		{"200 OK", http.StatusOK, true},
		{"201 Created", http.StatusCreated, true},
		{"204 No Content", http.StatusNoContent, true},
		{"400 Bad Request", http.StatusBadRequest, false},
		{"401 Unauthorized", http.StatusUnauthorized, false},
		{"404 Not Found", http.StatusNotFound, false},
		{"500 Internal Server Error", http.StatusInternalServerError, false},
		{"502 Bad Gateway", http.StatusBadGateway, false},
		{"503 Service Unavailable", http.StatusServiceUnavailable, false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tc.statusCode)
				w.Write([]byte(`{"message": "test"}`))
			}))
			defer server.Close()

			client := NewClient()
			resp, err := client.GetWithConfig(server.URL, nil)
			require.NoError(t, err)
			assert.Equal(t, tc.statusCode, resp.Status)
			assert.Equal(t, tc.isSuccess, resp.IsSuccess())
			assert.Equal(t, !tc.isSuccess, resp.IsError())
		})
	}
}

// ==================== 重试机制测试 ====================

// TestRetry_ServerError 测试服务器错误重试
func TestRetry_ServerError(t *testing.T) {
	retryCount := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		retryCount++
		if retryCount < 3 {
			w.WriteHeader(http.StatusInternalServerError)
		} else {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"message": "success"}`))
		}
	}))
	defer server.Close()

	client := NewClient()
	resp, err := client.GetWithConfig(server.URL, nil)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.Status)
	assert.GreaterOrEqual(t, retryCount, 3) // 应该至少重试了 3 次
}

// TestRetry_429TooManyRequests 测试 429 状态码重试
func TestRetry_429TooManyRequests(t *testing.T) {
	retryCount := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		retryCount++
		if retryCount < 2 {
			w.Header().Set("Retry-After", "1")
			w.WriteHeader(http.StatusTooManyRequests)
		} else {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"message": "success"}`))
		}
	}))
	defer server.Close()

	client := NewClient()
	resp, err := client.GetWithConfig(server.URL, nil)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.Status)
	assert.GreaterOrEqual(t, retryCount, 2)
}

// TestRetry_CustomConfig 测试自定义重试配置
func TestRetry_CustomConfig(t *testing.T) {
	retryCount := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		retryCount++
		if retryCount < 2 {
			w.WriteHeader(http.StatusInternalServerError)
		} else {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"message": "success"}`))
		}
	}))
	defer server.Close()

	client := NewClient()
	retryConfig := NewRetryConfig().
		WithRetryCount(2).
		WithRetryWaitTime(100 * time.Millisecond)
	config := NewRequestConfig().WithRetry(retryConfig)

	resp, err := client.GetWithConfig(server.URL, config)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.Status)
	assert.GreaterOrEqual(t, retryCount, 2)
}

// TestRetry_DisableRetry 测试禁用重试
func TestRetry_DisableRetry(t *testing.T) {
	retryCount := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		retryCount++
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	client := NewClient()
	retryConfig := NewRetryConfig().DisableRetry()
	config := NewRequestConfig().WithRetry(retryConfig)

	resp, err := client.GetWithConfig(server.URL, config)
	require.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, resp.Status)
	assert.Equal(t, 1, retryCount) // 不应该重试
}

// TestRetry_CustomCondition 测试自定义重试条件
func TestRetry_CustomCondition(t *testing.T) {
	retryCount := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		retryCount++
		if retryCount < 2 {
			w.WriteHeader(http.StatusBadRequest) // 通常 4xx 不重试
		} else {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"message": "success"}`))
		}
	}))
	defer server.Close()

	client := NewClient()
	// 自定义重试条件：即使是 400 也重试
	retryConfig := NewRetryConfig().
		WithRetryCount(2).
		WithRetryWaitTime(100 * time.Millisecond).
		WithRetryCondition(func(resp *resty.Response, err error) bool {
			return resp.StatusCode() == http.StatusBadRequest
		})
	config := NewRequestConfig().WithRetry(retryConfig)

	resp, err := client.GetWithConfig(server.URL, config)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.Status)
	assert.GreaterOrEqual(t, retryCount, 2)
}

// TestRetry_NoRetryOn4xx 测试 4xx 错误不重试
func TestRetry_NoRetryOn4xx(t *testing.T) {
	retryCount := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		retryCount++
		w.WriteHeader(http.StatusBadRequest)
	}))
	defer server.Close()

	client := NewClient()
	resp, err := client.GetWithConfig(server.URL, nil)
	require.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, resp.Status)
	assert.Equal(t, 1, retryCount) // 4xx 不应该重试
}

// TestDefaultRetryCondition 测试默认重试条件
func TestDefaultRetryCondition(t *testing.T) {
	testCases := []struct {
		name     string
		status   int
		err      error
		expected bool
	}{
		{"200 OK", http.StatusOK, nil, false},
		{"400 Bad Request", http.StatusBadRequest, nil, false},
		{"429 Too Many Requests", http.StatusTooManyRequests, nil, true},
		{"500 Internal Server Error", http.StatusInternalServerError, nil, true},
		{"502 Bad Gateway", http.StatusBadGateway, nil, true},
		{"503 Service Unavailable", http.StatusServiceUnavailable, nil, true},
		{"Network timeout error", 0, fmt.Errorf("timeout"), true},
		{"Connection refused", 0, fmt.Errorf("connection refused"), true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var resp *resty.Response
			if tc.status > 0 {
				server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(tc.status)
				}))
				defer server.Close()

				client := NewClient()
				resp, _ = client.Get(server.URL)
			}

			result := DefaultRetryCondition(resp, tc.err)
			assert.Equal(t, tc.expected, result)
		})
	}
}

// ==================== 错误处理测试 ====================

// TestClient_InvalidURL 测试无效 URL
func TestClient_InvalidURL(t *testing.T) {
	client := NewClient()
	_, err := client.GetWithConfig("invalid-url", nil)
	assert.Error(t, err)
}

// TestClient_InvalidMethod 测试无效的 HTTP 方法
func TestClient_InvalidMethod(t *testing.T) {
	client := NewClient()
	_, err := client.Stream(HttpMethod("INVALID"), "http://example.com", nil)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid HTTP method")
}

// TestClient_ConnectionRefused 测试连接被拒绝（需要模拟）
func TestClient_ConnectionRefused(t *testing.T) {
	// 使用一个不存在的端口
	client := NewClient()
	// 设置较短的超时时间以便快速失败
	client.client.SetTimeout(1 * time.Second)

	_, err := client.GetWithConfig("http://127.0.0.1:99999/invalid", nil)
	assert.Error(t, err)
	// 应该包含网络错误相关的关键词
	assert.True(t, isRetryableNetworkError(err))
}

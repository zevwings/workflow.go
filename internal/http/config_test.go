package http

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ==================== RequestConfig 测试 ====================

// TestRequestConfig_WithQuery 测试查询参数配置
func TestRequestConfig_WithQuery(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "value1", r.URL.Query().Get("key1"))
		assert.Equal(t, "value2", r.URL.Query().Get("key2"))
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client := NewClient()
	config := NewRequestConfig().
		WithQuery(map[string]string{
			"key1": "value1",
			"key2": "value2",
		})

	resp, err := client.GetWithConfig(server.URL, config)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.Status)
}

// TestRequestConfig_WithHeaders 测试多个 Headers 配置
func TestRequestConfig_WithHeaders(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "value1", r.Header.Get("X-Header-1"))
		assert.Equal(t, "value2", r.Header.Get("X-Header-2"))
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client := NewClient()
	config := NewRequestConfig().
		WithHeaders(map[string]string{
			"X-Header-1": "value1",
			"X-Header-2": "value2",
		})

	resp, err := client.GetWithConfig(server.URL, config)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.Status)
}

// TestRequestConfig_WithAuth 测试通过配置设置 Basic Auth
func TestRequestConfig_WithAuth(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		username, password, ok := r.BasicAuth()
		assert.True(t, ok)
		assert.Equal(t, "testuser", username)
		assert.Equal(t, "testpass", password)
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client := NewClient()
	config := NewRequestConfig().
		WithAuth(NewAuthorization("testuser", "testpass"))

	resp, err := client.GetWithConfig(server.URL, config)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.Status)
}

// TestRequestConfig_WithBody_JSON 测试 JSON Body 配置
func TestRequestConfig_WithBody_JSON(t *testing.T) {
	type RequestBody struct {
		Name  string `json:"name"`
		Value int    `json:"value"`
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var body RequestBody
		json.NewDecoder(r.Body).Decode(&body)
		assert.Equal(t, "test", body.Name)
		assert.Equal(t, 123, body.Value)
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client := NewClient()
	config := NewRequestConfig().
		WithBody(RequestBody{
			Name:  "test",
			Value: 123,
		})

	resp, err := client.PostWithConfig(server.URL, config)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.Status)
}

// TestRequestConfig_Chain 测试请求配置链式调用
func TestRequestConfig_Chain(t *testing.T) {
	config := NewRequestConfig().
		WithBody(map[string]string{"key": "value"}).
		WithQuery(map[string]string{"q": "test"}).
		WithHeader("X-Header", "value").
		WithAuth(NewAuthorization("user", "pass"))

	assert.NotNil(t, config.Body)
	assert.NotNil(t, config.Query)
	assert.Equal(t, "value", config.Headers["X-Header"])
	assert.NotNil(t, config.Auth)
}

// TestRequestConfig_DefaultValues 测试请求配置默认值
func TestRequestConfig_DefaultValues(t *testing.T) {
	config := NewRequestConfig()
	assert.NotNil(t, config.Headers)
	assert.Nil(t, config.Body)
	assert.Nil(t, config.Query)
	assert.Nil(t, config.Auth)
	assert.Nil(t, config.Retry)
}

// TestRequestConfig_NilConfig 测试 nil 配置
func TestRequestConfig_NilConfig(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"message": "success"}`))
	}))
	defer server.Close()

	client := NewClient()
	resp, err := client.GetWithConfig(server.URL, nil)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.Status)
}

// TestRequestConfig_QueryParams 测试查询参数转换
func TestRequestConfig_QueryParams(t *testing.T) {
	testCases := []struct {
		name     string
		query    interface{}
		expected map[string]string
	}{
		{
			name:     "map[string]string",
			query:    map[string]string{"key1": "value1", "key2": "value2"},
			expected: map[string]string{"key1": "value1", "key2": "value2"},
		},
		{
			name:     "map[string]interface{}",
			query:    map[string]interface{}{"key1": "value1", "key2": 123},
			expected: map[string]string{"key1": "value1", "key2": "123"},
		},
		{
			name:     "[]string with key=value format",
			query:    []string{"key1=value1", "key2=value2"},
			expected: map[string]string{"key1": "value1", "key2": "value2"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				for k, v := range tc.expected {
					assert.Equal(t, v, r.URL.Query().Get(k))
				}
				w.WriteHeader(http.StatusOK)
			}))
			defer server.Close()

			client := NewClient()
			config := NewRequestConfig().WithQuery(tc.query)
			resp, err := client.GetWithConfig(server.URL, config)
			require.NoError(t, err)
			assert.Equal(t, http.StatusOK, resp.Status)
		})
	}
}

// ==================== RetryConfig 测试 ====================

// TestRetryConfig_Chain 测试重试配置链式调用
func TestRetryConfig_Chain(t *testing.T) {
	config := NewRetryConfig().
		WithRetryCount(5).
		WithRetryWaitTime(2 * time.Second).
		WithRetryMaxWaitTime(60 * time.Second)

	assert.Equal(t, 5, config.Count)
	assert.Equal(t, 2*time.Second, config.WaitTime)
	assert.Equal(t, 60*time.Second, config.MaxWaitTime)
}

// TestRetryConfig_DefaultValues 测试重试配置默认值
func TestRetryConfig_DefaultValues(t *testing.T) {
	config := NewRetryConfig()
	assert.Equal(t, 3, config.Count)
	assert.Equal(t, 1*time.Second, config.WaitTime)
	assert.Equal(t, 30*time.Second, config.MaxWaitTime)
}

// ==================== MultipartRequestConfig 测试 ====================

// TestMultipartRequestConfig_Chain 测试 Multipart 配置链式调用
func TestMultipartRequestConfig_Chain(t *testing.T) {
	config := NewMultipartRequestConfig().
		WithMultipartField(MultipartField{
			ParamName: "field1",
			FileName:  "value1",
		}).
		WithMultipartField(MultipartField{
			ParamName: "field2",
			FileName:  "value2",
		}).
		WithHeader("X-Header", "value")

	assert.Len(t, config.MultipartFields, 2)
	assert.Equal(t, "value", config.Headers["X-Header"])
}


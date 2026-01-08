package client

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	httppkg "github.com/zevwings/workflow/internal/http"
)

// ==================== NewClient 测试 ====================

func TestNewClient(t *testing.T) {
	httpClient := httppkg.NewClient()
	config := &ProviderConfig{
		APIKey: "test-api-key",
		Model:  "gpt-3.5-turbo",
		URL:    "https://api.openai.com/v1/chat/completions",
	}

	client := NewClient(httpClient, config)
	assert.NotNil(t, client)
	assert.Equal(t, httpClient, client.httpClient)
	assert.Equal(t, config, client.config)
}

func TestNewClient_NilHTTPClient(t *testing.T) {
	config := &ProviderConfig{
		APIKey: "test-api-key",
		Model:  "gpt-3.5-turbo",
		URL:    "https://api.openai.com/v1/chat/completions",
	}

	assert.Panics(t, func() {
		NewClient(nil, config)
	}, "应该 panic 当 httpClient 为 nil")
}

func TestNewClient_NilConfig(t *testing.T) {
	httpClient := httppkg.NewClient()

	assert.Panics(t, func() {
		NewClient(httpClient, nil)
	}, "应该 panic 当 config 为 nil")
}

// ==================== Global 测试 ====================

func TestGlobal(t *testing.T) {
	httpClient1 := httppkg.NewClient()
	httpClient2 := httppkg.NewClient()
	config1 := &ProviderConfig{
		APIKey: "test-api-key-1",
		Model:  "gpt-3.5-turbo",
		URL:    "https://api.openai.com/v1/chat/completions",
	}
	config2 := &ProviderConfig{
		APIKey: "test-api-key-2",
		Model:  "gpt-4",
		URL:    "https://api.openai.com/v1/chat/completions",
	}

	// 首次调用
	client1 := Global(httpClient1, config1)
	assert.NotNil(t, client1)

	// 后续调用应该返回同一个实例
	client2 := Global(httpClient2, config2)
	assert.Equal(t, client1, client2)
	assert.Equal(t, config1, client2.config, "后续调用的参数应该被忽略")
}

func TestGlobal_NilHTTPClient(t *testing.T) {
	config := &ProviderConfig{
		APIKey: "test-api-key",
		Model:  "gpt-3.5-turbo",
		URL:    "https://api.openai.com/v1/chat/completions",
	}

	assert.Panics(t, func() {
		Global(nil, config)
	}, "应该 panic 当 httpClient 为 nil")
}

func TestGlobal_NilConfig(t *testing.T) {
	httpClient := httppkg.NewClient()

	assert.Panics(t, func() {
		Global(httpClient, nil)
	}, "应该 panic 当 config 为 nil")
}

// ==================== Call 测试 ====================

func TestLLMClient_Call(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 验证请求方法
		assert.Equal(t, http.MethodPost, r.Method)

		// 验证请求头
		assert.Equal(t, "Bearer test-api-key", r.Header.Get("Authorization"))
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))

		// 验证请求体
		var body map[string]interface{}
		err := json.NewDecoder(r.Body).Decode(&body)
		require.NoError(t, err)
		assert.Equal(t, "gpt-3.5-turbo", body["model"])
		assert.InDelta(t, 0.7, body["temperature"].(float64), 0.0001)

		// 返回响应
		response := map[string]interface{}{
			"id":      "chatcmpl-123",
			"object":  "chat.completion",
			"created": 1677652288,
			"model":   "gpt-3.5-turbo",
			"choices": []map[string]interface{}{
				{
					"index": 0,
					"message": map[string]interface{}{
						"role":    "assistant",
						"content": "Hello, how can I help you?",
					},
					"finish_reason": "stop",
				},
			},
			"usage": map[string]interface{}{
				"prompt_tokens":     10,
				"completion_tokens": 10,
				"total_tokens":      20,
			},
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	httpClient := httppkg.NewClient()
	config := &ProviderConfig{
		APIKey: "test-api-key",
		Model:  "gpt-3.5-turbo",
		URL:    server.URL,
	}

	client := NewClient(httpClient, config)
	params := &LLMRequestParams{
		SystemPrompt: "You are a helpful assistant.",
		UserPrompt:   "Hello",
		Temperature:  0.7,
	}

	response, err := client.Call(params)
	require.NoError(t, err)
	assert.Equal(t, "Hello, how can I help you?", response)
}

func TestLLMClient_Call_WithMaxTokens(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var body map[string]interface{}
		json.NewDecoder(r.Body).Decode(&body)

		// 验证 max_tokens 字段存在
		assert.Contains(t, body, "max_tokens")
		assert.Equal(t, float64(100), body["max_tokens"])

		response := map[string]interface{}{
			"id":      "chatcmpl-123",
			"object":  "chat.completion",
			"created": 1677652288,
			"model":   "gpt-3.5-turbo",
			"choices": []map[string]interface{}{
				{
					"index": 0,
					"message": map[string]interface{}{
						"role":    "assistant",
						"content": "Response",
					},
					"finish_reason": "stop",
				},
			},
			"usage": map[string]interface{}{
				"prompt_tokens":     10,
				"completion_tokens": 10,
				"total_tokens":      20,
			},
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	httpClient := httppkg.NewClient()
	config := &ProviderConfig{
		APIKey: "test-api-key",
		Model:  "gpt-3.5-turbo",
		URL:    server.URL,
	}

	client := NewClient(httpClient, config)
	maxTokens := 100
	params := &LLMRequestParams{
		SystemPrompt: "You are a helpful assistant.",
		UserPrompt:   "Hello",
		Temperature:  0.7,
		MaxTokens:    &maxTokens,
	}

	response, err := client.Call(params)
	require.NoError(t, err)
	assert.Equal(t, "Response", response)
}

func TestLLMClient_Call_WithCustomModel(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var body map[string]interface{}
		json.NewDecoder(r.Body).Decode(&body)

		// 验证使用 params 中的 model 而不是 config 中的
		assert.Equal(t, "gpt-4", body["model"])

		response := map[string]interface{}{
			"id":      "chatcmpl-123",
			"object":  "chat.completion",
			"created": 1677652288,
			"model":   "gpt-4",
			"choices": []map[string]interface{}{
				{
					"index": 0,
					"message": map[string]interface{}{
						"role":    "assistant",
						"content": "Response",
					},
					"finish_reason": "stop",
				},
			},
			"usage": map[string]interface{}{
				"prompt_tokens":     10,
				"completion_tokens": 10,
				"total_tokens":      20,
			},
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	httpClient := httppkg.NewClient()
	config := &ProviderConfig{
		APIKey: "test-api-key",
		Model:  "gpt-3.5-turbo",
		URL:    server.URL,
	}

	client := NewClient(httpClient, config)
	params := &LLMRequestParams{
		SystemPrompt: "You are a helpful assistant.",
		UserPrompt:   "Hello",
		Temperature:  0.7,
		Model:        "gpt-4",
	}

	response, err := client.Call(params)
	require.NoError(t, err)
	assert.Equal(t, "Response", response)
}

func TestLLMClient_Call_EmptyResponse(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response := map[string]interface{}{
			"id":      "chatcmpl-123",
			"object":  "chat.completion",
			"created": 1677652288,
			"model":   "gpt-3.5-turbo",
			"choices": []map[string]interface{}{
				{
					"index": 0,
					"message": map[string]interface{}{
						"role":    "assistant",
						"content": nil,
					},
					"finish_reason": "stop",
				},
			},
			"usage": map[string]interface{}{
				"prompt_tokens":     10,
				"completion_tokens": 0,
				"total_tokens":      10,
			},
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	httpClient := httppkg.NewClient()
	config := &ProviderConfig{
		APIKey: "test-api-key",
		Model:  "gpt-3.5-turbo",
		URL:    server.URL,
	}

	client := NewClient(httpClient, config)
	params := &LLMRequestParams{
		SystemPrompt: "You are a helpful assistant.",
		UserPrompt:   "Hello",
		Temperature:  0.7,
	}

	_, err := client.Call(params)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "content 为空")
}

func TestLLMClient_Call_NoChoices(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response := map[string]interface{}{
			"id":      "chatcmpl-123",
			"object":  "chat.completion",
			"created": 1677652288,
			"model":   "gpt-3.5-turbo",
			"choices": []map[string]interface{}{},
			"usage": map[string]interface{}{
				"prompt_tokens":     10,
				"completion_tokens": 0,
				"total_tokens":      10,
			},
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	httpClient := httppkg.NewClient()
	config := &ProviderConfig{
		APIKey: "test-api-key",
		Model:  "gpt-3.5-turbo",
		URL:    server.URL,
	}

	client := NewClient(httpClient, config)
	params := &LLMRequestParams{
		SystemPrompt: "You are a helpful assistant.",
		UserPrompt:   "Hello",
		Temperature:  0.7,
	}

	_, err := client.Call(params)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "choices 数组或数组为空")
}

func TestLLMClient_Call_HTTPError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error": {"message": "Internal server error"}}`))
	}))
	defer server.Close()

	httpClient := httppkg.NewClient()
	config := &ProviderConfig{
		APIKey: "test-api-key",
		Model:  "gpt-3.5-turbo",
		URL:    server.URL,
	}

	client := NewClient(httpClient, config)
	params := &LLMRequestParams{
		SystemPrompt: "You are a helpful assistant.",
		UserPrompt:   "Hello",
		Temperature:  0.7,
	}

	_, err := client.Call(params)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "LLM API 请求失败")
}

func TestLLMClient_Call_InvalidJSON(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("invalid json"))
	}))
	defer server.Close()

	httpClient := httppkg.NewClient()
	config := &ProviderConfig{
		APIKey: "test-api-key",
		Model:  "gpt-3.5-turbo",
		URL:    server.URL,
	}

	client := NewClient(httpClient, config)
	params := &LLMRequestParams{
		SystemPrompt: "You are a helpful assistant.",
		UserPrompt:   "Hello",
		Temperature:  0.7,
	}

	_, err := client.Call(params)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "解析响应 JSON 失败")
}

// ==================== buildURL 测试 ====================

func TestLLMClient_buildURL(t *testing.T) {
	httpClient := httppkg.NewClient()
	config := &ProviderConfig{
		APIKey: "test-api-key",
		Model:  "gpt-3.5-turbo",
		URL:    "https://api.openai.com/v1/chat/completions",
	}

	client := NewClient(httpClient, config)
	url, err := client.buildURL()
	require.NoError(t, err)
	assert.Equal(t, "https://api.openai.com/v1/chat/completions", url)
}

func TestLLMClient_buildURL_EmptyURL(t *testing.T) {
	httpClient := httppkg.NewClient()
	config := &ProviderConfig{
		APIKey: "test-api-key",
		Model:  "gpt-3.5-turbo",
		URL:    "",
	}

	client := NewClient(httpClient, config)
	_, err := client.buildURL()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "URL 未配置")
}

// ==================== buildHeaders 测试 ====================

func TestLLMClient_buildHeaders(t *testing.T) {
	httpClient := httppkg.NewClient()
	config := &ProviderConfig{
		APIKey: "test-api-key",
		Model:  "gpt-3.5-turbo",
		URL:    "https://api.openai.com/v1/chat/completions",
	}

	client := NewClient(httpClient, config)
	headers, err := client.buildHeaders()
	require.NoError(t, err)
	assert.Equal(t, "Bearer test-api-key", headers["Authorization"])
	assert.Equal(t, "application/json", headers["Content-Type"])
}

func TestLLMClient_buildHeaders_EmptyAPIKey(t *testing.T) {
	httpClient := httppkg.NewClient()
	config := &ProviderConfig{
		APIKey: "",
		Model:  "gpt-3.5-turbo",
		URL:    "https://api.openai.com/v1/chat/completions",
	}

	client := NewClient(httpClient, config)
	_, err := client.buildHeaders()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "LLM API key 未配置")
}

// ==================== buildModel 测试 ====================

func TestLLMClient_buildModel(t *testing.T) {
	httpClient := httppkg.NewClient()
	config := &ProviderConfig{
		APIKey: "test-api-key",
		Model:  "gpt-3.5-turbo",
		URL:    "https://api.openai.com/v1/chat/completions",
	}

	client := NewClient(httpClient, config)
	model, err := client.buildModel()
	require.NoError(t, err)
	assert.Equal(t, "gpt-3.5-turbo", model)
}

func TestLLMClient_buildModel_EmptyModel(t *testing.T) {
	httpClient := httppkg.NewClient()
	config := &ProviderConfig{
		APIKey: "test-api-key",
		Model:  "",
		URL:    "https://api.openai.com/v1/chat/completions",
	}

	client := NewClient(httpClient, config)
	_, err := client.buildModel()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "model 未配置")
}

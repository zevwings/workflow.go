package client

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ==================== newClient 测试 ====================

func TestNewClient(t *testing.T) {
	config := &ProviderConfig{
		APIKey: "test-api-key",
		Model:  "gpt-3.5-turbo",
		URL:    "https://api.openai.com/v1/chat/completions",
	}

	client := newClient(config)
	assert.NotNil(t, client)
	// 注意：client 现在是接口类型，无法直接访问字段
	// 可以通过调用方法来验证客户端是否正常工作
}

// TestNewClient_NilHTTPClient 已移除，因为 newClient 不再接受 httpClient 参数

func TestNewClient_NilConfig(t *testing.T) {
	assert.Panics(t, func() {
		newClient(nil)
	}, "应该 panic 当 config 为 nil")
}

// ==================== Global 测试 ====================

func TestGlobal(t *testing.T) {
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
	client1 := Global(config1)
	assert.NotNil(t, client1)

	// 后续调用应该返回同一个实例
	client2 := Global(config2)
	assert.Equal(t, client1, client2)
	// 注意：client 现在是接口类型，无法直接访问 config 字段
	// 单例行为已通过 client1 == client2 验证
}

// TestGlobal_NilHTTPClient 已移除，因为 Global 不再接受 httpClient 参数

func TestGlobal_NilConfig(t *testing.T) {
	assert.Panics(t, func() {
		Global(nil)
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

	config := &ProviderConfig{
		APIKey: "test-api-key",
		Model:  "gpt-3.5-turbo",
		URL:    server.URL,
	}

	client := newClient(config)
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

	config := &ProviderConfig{
		APIKey: "test-api-key",
		Model:  "gpt-3.5-turbo",
		URL:    server.URL,
	}

	client := newClient(config)
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

	config := &ProviderConfig{
		APIKey: "test-api-key",
		Model:  "gpt-3.5-turbo",
		URL:    server.URL,
	}

	client := newClient(config)
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

	config := &ProviderConfig{
		APIKey: "test-api-key",
		Model:  "gpt-3.5-turbo",
		URL:    server.URL,
	}

	client := newClient(config)
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

	config := &ProviderConfig{
		APIKey: "test-api-key",
		Model:  "gpt-3.5-turbo",
		URL:    server.URL,
	}

	client := newClient(config)
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

	config := &ProviderConfig{
		APIKey: "test-api-key",
		Model:  "gpt-3.5-turbo",
		URL:    server.URL,
	}

	client := newClient(config)
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

	config := &ProviderConfig{
		APIKey: "test-api-key",
		Model:  "gpt-3.5-turbo",
		URL:    server.URL,
	}

	client := newClient(config)
	params := &LLMRequestParams{
		SystemPrompt: "You are a helpful assistant.",
		UserPrompt:   "Hello",
		Temperature:  0.7,
	}

	_, err := client.Call(params)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "解析响应 JSON 失败")
}

// 注意：buildURL、buildHeaders、buildModel 等未导出方法的测试已移除
// 因为这些是内部实现细节，现在通过 LLMClient 接口隐藏。
// 这些功能已经通过 Call 方法的测试间接覆盖。

package branch

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zevwings/workflow/internal/llm/client"
)

// ==================== NewBranchLLMClient 测试 ====================

func TestNewBranchLLMClient(t *testing.T) {
	config := &client.ProviderConfig{
		APIKey: "test-api-key",
		Model:  "gpt-3.5-turbo",
		URL:    "https://api.openai.com/v1/chat/completions",
	}
	llmClient := client.Global(config)

	branchClient := newBranchLLMClient(llmClient)
	assert.NotNil(t, branchClient)
	assert.Equal(t, llmClient, branchClient.llmClient)
}

func TestNewBranchLLMClient_NilLLMClient(t *testing.T) {
	assert.Panics(t, func() {
		newBranchLLMClient(nil)
	}, "应该 panic 当 llmClient 为 nil")
}

// ==================== TranslateToEnglish 测试 ====================

func TestBranchLLMClient_TranslateToEnglish(t *testing.T) {
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
						"content": "Add user login",
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

	config := &client.ProviderConfig{
		APIKey: "test-api-key",
		Model:  "gpt-3.5-turbo",
		URL:    server.URL,
	}
	llmClient := client.Global(config)
	branchClient := newBranchLLMClient(llmClient)

	translated, err := branchClient.TranslateToEnglish("添加用户登录")
	require.NoError(t, err)
	assert.Equal(t, "Add user login", translated)
}

func TestBranchLLMClient_TranslateToEnglish_WithQuotes(t *testing.T) {
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
						"content": `"Add user login"`,
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

	config := &client.ProviderConfig{
		APIKey: "test-api-key",
		Model:  "gpt-3.5-turbo",
		URL:    server.URL,
	}
	llmClient := client.Global(config)
	branchClient := newBranchLLMClient(llmClient)

	translated, err := branchClient.TranslateToEnglish("添加用户登录")
	require.NoError(t, err)
	assert.Equal(t, "Add user login", translated, "应该移除引号")
}

func TestBranchLLMClient_TranslateToEnglish_EmptyResponse(t *testing.T) {
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
						"content": "",
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

	config := &client.ProviderConfig{
		APIKey: "test-api-key",
		Model:  "gpt-3.5-turbo",
		URL:    server.URL,
	}
	llmClient := client.Global(config)
	branchClient := newBranchLLMClient(llmClient)

	_, err := branchClient.TranslateToEnglish("添加用户登录")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "content 为空")
}

func TestBranchLLMClient_TranslateToEnglish_WhitespaceOnly(t *testing.T) {
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
						"content": "   ",
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

	config := &client.ProviderConfig{
		APIKey: "test-api-key",
		Model:  "gpt-3.5-turbo",
		URL:    server.URL,
	}
	llmClient := client.Global(config)
	branchClient := newBranchLLMClient(llmClient)

	_, err := branchClient.TranslateToEnglish("添加用户登录")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "content 为空")
}

func TestBranchLLMClient_TranslateToEnglish_HTTPError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error": {"message": "Internal server error"}}`))
	}))
	defer server.Close()

	config := &client.ProviderConfig{
		APIKey: "test-api-key",
		Model:  "gpt-3.5-turbo",
		URL:    server.URL,
	}
	llmClient := client.Global(config)
	branchClient := newBranchLLMClient(llmClient)

	_, err := branchClient.TranslateToEnglish("添加用户登录")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "调用 LLM API 翻译文本失败")
}

package pr

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	httppkg "github.com/zevwings/workflow/internal/http"
	"github.com/zevwings/workflow/internal/llm/client"
)

// ==================== NewPullRequestLLMClient 测试 ====================

func TestNewPullRequestLLMClient(t *testing.T) {
	httpClient := httppkg.NewClient()
	config := &client.ProviderConfig{
		APIKey: "test-api-key",
		Model:  "gpt-3.5-turbo",
		URL:    "https://api.openai.com/v1/chat/completions",
	}
	llmClient := client.NewClient(httpClient, config)

	prClient := NewPullRequestLLMClient(llmClient, nil)
	assert.NotNil(t, prClient)
	assert.Equal(t, llmClient, prClient.llmClient)
}

func TestNewPullRequestLLMClient_NilLLMClient(t *testing.T) {
	assert.Panics(t, func() {
		NewPullRequestLLMClient(nil, nil)
	}, "应该 panic 当 llmClient 为 nil")
}

func TestNewPullRequestLLMClient_WithLang(t *testing.T) {
	httpClient := httppkg.NewClient()
	config := &client.ProviderConfig{
		APIKey: "test-api-key",
		Model:  "gpt-3.5-turbo",
		URL:    "https://api.openai.com/v1/chat/completions",
	}
	llmClient := client.NewClient(httpClient, config)
	lang := &client.SupportedLanguage{
		Code:                "zh-CN",
		Name:                "Chinese",
		NativeName:          "中文",
		InstructionTemplate: "**所有输出必须仅使用中文。**",
	}

	prClient := NewPullRequestLLMClient(llmClient, lang)
	assert.NotNil(t, prClient)
	assert.Equal(t, lang, prClient.lang)
}

// ==================== GenerateContent 测试 ====================

func TestPullRequestLLMClient_GenerateContent(t *testing.T) {
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
						"content": `{"branch_name": "feature/add-user-login", "pr_title": "Add user login", "description": "This PR adds user login functionality", "scope": "auth"}`,
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
	config := &client.ProviderConfig{
		APIKey: "test-api-key",
		Model:  "gpt-3.5-turbo",
		URL:    server.URL,
	}
	llmClient := client.NewClient(httpClient, config)
	prClient := NewPullRequestLLMClient(llmClient, nil)

	content, err := prClient.GenerateContent("fix: add user login", []string{"main", "develop"}, "diff content")
	require.NoError(t, err)
	assert.Equal(t, "featureadd-user-login", content.BranchName)
	assert.Equal(t, "Add user login", content.PRTitle)
	require.NotNil(t, content.Description)
	assert.Equal(t, "This PR adds user login functionality", *content.Description)
	require.NotNil(t, content.Scope)
	assert.Equal(t, "auth", *content.Scope)
}

func TestPullRequestLLMClient_GenerateContent_WithoutOptionalFields(t *testing.T) {
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
						"content": `{"branch_name": "feature/add-user-login", "pr_title": "Add user login"}`,
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
	config := &client.ProviderConfig{
		APIKey: "test-api-key",
		Model:  "gpt-3.5-turbo",
		URL:    server.URL,
	}
	llmClient := client.NewClient(httpClient, config)
	prClient := NewPullRequestLLMClient(llmClient, nil)

	content, err := prClient.GenerateContent("fix: add user login", nil, "")
	require.NoError(t, err)
	assert.Equal(t, "featureadd-user-login", content.BranchName)
	assert.Equal(t, "Add user login", content.PRTitle)
	assert.Nil(t, content.Description)
	assert.Nil(t, content.Scope)
}

func TestPullRequestLLMClient_GenerateContent_InvalidJSON(t *testing.T) {
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
						"content": "invalid json",
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
	config := &client.ProviderConfig{
		APIKey: "test-api-key",
		Model:  "gpt-3.5-turbo",
		URL:    server.URL,
	}
	llmClient := client.NewClient(httpClient, config)
	prClient := NewPullRequestLLMClient(llmClient, nil)

	_, err := prClient.GenerateContent("fix: add user login", nil, "")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "解析 LLM 响应失败")
}

func TestPullRequestLLMClient_GenerateContent_MissingFields(t *testing.T) {
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
						"content": `{"pr_title": "Add user login"}`,
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
	config := &client.ProviderConfig{
		APIKey: "test-api-key",
		Model:  "gpt-3.5-turbo",
		URL:    server.URL,
	}
	llmClient := client.NewClient(httpClient, config)
	prClient := NewPullRequestLLMClient(llmClient, nil)

	_, err := prClient.GenerateContent("fix: add user login", nil, "")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "缺少 'branch_name' 字段")
}

// ==================== Summarize 测试 ====================

func TestPullRequestLLMClient_Summarize(t *testing.T) {
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
						"content": `{"summary": "# PR Summary\n\nThis PR adds user authentication.", "filename": "add-user-authentication"}`,
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
	config := &client.ProviderConfig{
		APIKey: "test-api-key",
		Model:  "gpt-3.5-turbo",
		URL:    server.URL,
	}
	llmClient := client.NewClient(httpClient, config)
	prClient := NewPullRequestLLMClient(llmClient, nil)

	summary, err := prClient.Summarize("Add user authentication", "diff content")
	require.NoError(t, err)
	assert.Contains(t, summary.Summary, "PR Summary")
	assert.Equal(t, "add-user-authentication", summary.Filename)
}

func TestPullRequestLLMClient_Summarize_MissingFields(t *testing.T) {
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
						"content": `{"summary": "# PR Summary"}`,
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
	config := &client.ProviderConfig{
		APIKey: "test-api-key",
		Model:  "gpt-3.5-turbo",
		URL:    server.URL,
	}
	llmClient := client.NewClient(httpClient, config)
	prClient := NewPullRequestLLMClient(llmClient, nil)

	_, err := prClient.Summarize("Add user authentication", "diff content")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "缺少 'filename' 字段")
}

// ==================== Reword 测试 ====================

func TestPullRequestLLMClient_Reword(t *testing.T) {
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
						"content": `{"pr_title": "Updated: Add user login", "description": "This PR adds user login functionality with JWT authentication"}`,
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
	config := &client.ProviderConfig{
		APIKey: "test-api-key",
		Model:  "gpt-3.5-turbo",
		URL:    server.URL,
	}
	llmClient := client.NewClient(httpClient, config)
	prClient := NewPullRequestLLMClient(llmClient, nil)

	currentTitle := "Add user login"
	reword, err := prClient.Reword("diff content", &currentTitle)
	require.NoError(t, err)
	assert.Equal(t, "Updated: Add user login", reword.PRTitle)
	require.NotNil(t, reword.Description)
	assert.Equal(t, "This PR adds user login functionality with JWT authentication", *reword.Description)
}

func TestPullRequestLLMClient_Reword_NilTitle(t *testing.T) {
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
						"content": `{"pr_title": "Add user login"}`,
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
	config := &client.ProviderConfig{
		APIKey: "test-api-key",
		Model:  "gpt-3.5-turbo",
		URL:    server.URL,
	}
	llmClient := client.NewClient(httpClient, config)
	prClient := NewPullRequestLLMClient(llmClient, nil)

	reword, err := prClient.Reword("diff content", nil)
	require.NoError(t, err)
	assert.Equal(t, "Add user login", reword.PRTitle)
}

// ==================== SummarizeFileChange 测试 ====================

func TestPullRequestLLMClient_SummarizeFileChange(t *testing.T) {
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
						"content": "Added user authentication logic with JWT token generation",
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
	config := &client.ProviderConfig{
		APIKey: "test-api-key",
		Model:  "gpt-3.5-turbo",
		URL:    server.URL,
	}
	llmClient := client.NewClient(httpClient, config)
	prClient := NewPullRequestLLMClient(llmClient, nil)

	summary, err := prClient.SummarizeFileChange("auth/login.go", "diff content")
	require.NoError(t, err)
	assert.Contains(t, summary, "authentication")
}


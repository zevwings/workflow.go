// Package client provides unified configuration-driven LLM client implementation, supporting OpenAI, DeepSeek, and proxy APIs.
//
// Main features:
//   - LLMClient: LLM client supporting multiple providers
//   - LLMRequestParams: LLM request parameters
//   - Language support: Multi-language prompt enhancement
//   - Configuration-driven: Pass configuration through configuration struct
//
// Usage example:
//
//	import (
//		"github.com/zevwings/workflow/internal/http"
//		"github.com/zevwings/workflow/internal/llm/client"
//	)
//
//	// Method 1: Use convenience function to create configuration (recommended)
//	cfg := client.NewOpenAIConfig("your-api-key", "gpt-3.5-turbo")
//	llmClient := client.Global(cfg.OpenAI) // Automatically uses http.Global()
//
//	// Method 2: Manually create configuration
//	cfg := &client.ProviderConfig{
//		APIKey: "your-api-key",
//		Model:  "gpt-3.5-turbo",
//		URL:    "https://api.openai.com/v1/chat/completions",
//	}
//	llmClient := client.Global(cfg) // Automatically uses http.Global()
//
//	// Prepare request parameters
//	params := &client.LLMRequestParams{
//		SystemPrompt: "You are a helpful assistant.",
//		UserPrompt:   "What is Go?",
//		Temperature:  0.7,
//	}
//
//	// Call LLM API
//	response, err := llmClient.Call(params)
//	if err != nil {
//		// Handle error
//	}
//
//	fmt.Println(response)
package client

import (
	"encoding/json"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/zevwings/workflow/internal/http"
	"github.com/zevwings/workflow/internal/logging"
)

var (
	// globalClient global LLM client singleton
	globalClient *llmClient
	globalOnce   sync.Once
)

// LLMClient LLM client interface
//
// All LLM providers use the same client implementation, distinguished by configuration struct.
// All configuration (URL, API key, model, response_format) is obtained from configuration struct.
type LLMClient interface {
	// Call calls LLM API
	//
	// Parameters:
	//   - params: LLM request parameters
	//
	// Returns:
	//   - string: LLM generated text content (trimmed of leading and trailing whitespace)
	//   - error: Returns corresponding error message if API call fails or response format is incorrect
	Call(params *LLMRequestParams) (string, error)
}

// llmClient LLM client implementation
//
// All LLM providers use the same client implementation, distinguished by configuration struct.
// All configuration (URL, API key, model, response_format) is obtained from configuration struct.
type llmClient struct {
	httpClient http.Client
	config     *ProviderConfig
}

// newClient creates a new LLM client (private function)
//
// Automatically uses http.Global() to get global HTTP client.
//
// Parameters:
//   - config: LLM configuration struct (cannot be nil)
//
// Returns:
//   - LLMClient: LLM client instance
func newClient(config *ProviderConfig) LLMClient {
	if config == nil {
		panic(fmt.Errorf("llm/client.newClient: config cannot be nil"))
	}
	return &llmClient{
		httpClient: http.Global(),
		config:     config,
	}
}

// Global gets global LLMClient singleton
//
// Returns process-level LLMClient singleton.
// Singleton is initialized on first call, subsequent calls reuse the same instance.
// Automatically uses http.Global() to get global HTTP client.
//
// Parameters:
//   - config: LLM configuration struct (required, cannot be nil)
//
// Returns:
//   - LLMClient: LLM client interface instance
//
// Note:
//   - Configuration must be provided by caller, LLM module is not responsible for configuration creation and lifecycle
//   - Parameters passed on first call will be saved, subsequent calls will ignore parameters
//   - If config is nil, will panic on first call
//   - Returns interface type, external code cannot directly access struct fields or unexported methods
//
// Advantages:
//   - Reduce resource consumption: Avoid duplicate client instance creation
//   - Thread-safe: Can be safely used in multi-threaded environments
//   - Unified management: All LLM calls use the same client instance
//   - Convenience: Automatically uses global HTTP client, no need to manually pass
//   - Encapsulation: Returns interface type, hides implementation details, prevents external direct access to internal structure
func Global(config *ProviderConfig) LLMClient {
	if config == nil {
		panic(fmt.Errorf("llm/client.Global: config cannot be nil"))
	}
	globalOnce.Do(func() {
		globalClient = newClient(config).(*llmClient)
	})
	return globalClient
}

// Call calls LLM API
//
// Parameters:
//   - params: LLM request parameters
//
// Returns:
//   - string: LLM generated text content (trimmed of leading and trailing whitespace)
//   - error: Returns corresponding error message if API call fails or response format is incorrect
func (c *llmClient) Call(params *LLMRequestParams) (string, error) {
	logger := logging.GetLogger()

	// Build URL (unified format)
	url, err := c.buildURL()
	if err != nil {
		logger.WithError(err).Error("Failed to build LLM API URL: URL not configured")
		return "", fmt.Errorf("failed to build URL: %w", err)
	}

	// Record LLM API call start
	logger.WithFields(logging.Fields{
		"model": c.config.Model,
		"url":   url,
	}).Info("Starting LLM API call")

	// Build request body (unified format)
	payload, err := c.buildPayload(params)
	if err != nil {
		logger.WithError(err).Error("Failed to build LLM request payload")
		return "", fmt.Errorf("failed to build request body: %w", err)
	}

	// Build request headers (unified format)
	headers, err := c.buildHeaders()
	if err != nil {
		logger.WithError(err).Error("Failed to build LLM request headers: API key not configured")
		return "", fmt.Errorf("failed to build request headers: %w", err)
	}

	// Record request parameter details (Debug level)
	model := c.config.Model
	if params.Model != "" {
		model = params.Model
	}
	maxTokens := "nil"
	if params.MaxTokens != nil {
		maxTokens = fmt.Sprintf("%d", *params.MaxTokens)
	}
	logger.WithFields(logging.Fields{
		"model":       model,
		"temperature": params.Temperature,
		"max_tokens":  maxTokens,
		"url":         url,
	}).Debug("LLM request parameters")

	// Build request configuration
	reqConfig := http.NewRequestConfig().
		WithHeaders(headers).
		WithBody(payload).
		WithTimeout(60 * time.Second).                     // LLM API usually requires longer timeout
		WithRetry(http.NewRetryConfig().WithRetryCount(3)) // Retry up to 3 times

	// Record before sending HTTP request
	logger.WithFields(logging.Fields{
		"url":     url,
		"payload": payload,
		"headers": headers,
		"timeout": 180 * time.Second,
		"retries": 3,
	}).Info("Sending LLM HTTP POST request (timeout: 60s, retries: 3)")

	// Send request
	resp, err := c.httpClient.PostWithConfig(url, reqConfig)
	if err != nil {
		logger.WithError(err).WithField("url", url).Error("LLM HTTP request failed")
		return "", fmt.Errorf("failed to send LLM request to %s: %w", url, err)
	}

	// Check error (use EnsureSuccessWith for unified handling)
	resp, err = resp.EnsureSuccessWith(func(r *http.HttpResponse) error {
		errorMessage := r.ExtractErrorMessage()
		logger.Warnf("LLM API returned error status: url=%s, status=%d, error=%s",
			url, r.Status, errorMessage)
		return fmt.Errorf("LLM API request failed (%s): %d - %s", url, r.Status, errorMessage)
	})
	if err != nil {
		return "", err
	}

	// Parse JSON response
	var data map[string]interface{}
	data, err = http.AsJSON[map[string]interface{}](resp)
	if err != nil {
		logger.WithError(err).Error("Failed to parse LLM response as JSON")
		return "", fmt.Errorf("failed to parse response JSON: %w", err)
	}

	// Extract content based on configured response format
	content, err := c.extractContent(data)
	if err != nil {
		logger.WithError(err).Error("Failed to extract content from LLM response")
		return "", fmt.Errorf("failed to extract response content: %w", err)
	}

	// Record response content summary (Debug level)
	contentPreview := content
	if len(contentPreview) > 100 {
		contentPreview = contentPreview[:100] + "..."
	}
	logger.Debugf("LLM response content preview: %s", contentPreview)

	// Record LLM API call success
	logger.WithFields(logging.Fields{
		"model":          c.config.Model,
		"url":            url,
		"content_length": len(content),
	}).Info("LLM API call succeeded")

	return content, nil
}

// buildURL builds API URL
//
// Gets URL directly from configuration struct.
//
// Returns:
//   - string: API URL
//   - error: Returns error if provider is not configured or invalid
func (c *llmClient) buildURL() (string, error) {
	url := c.config.URL
	if url == "" {
		return "", fmt.Errorf("URL not configured")
	}
	return fmt.Sprintf("%s/chat/completions", url), nil
}

// buildHeaders builds request headers
//
// Returns:
//   - map[string]string: Request header map
//   - error: Returns error if API key is not configured
func (c *llmClient) buildHeaders() (map[string]string, error) {
	headers := make(map[string]string)

	var apiKey string = c.config.APIKey

	if apiKey == "" {
		return nil, fmt.Errorf("LLM API key not configured")
	}

	headers["Authorization"] = fmt.Sprintf("Bearer %s", apiKey)
	headers["Content-Type"] = "application/json"

	return headers, nil
}

// buildModel builds model name
//
// Gets model name directly from configuration struct.
//
// Returns:
//   - string: Model name
//   - error: Returns error if model is not configured
func (c *llmClient) buildModel() (string, error) {
	var model = c.config.Model
	if model == "" {
		return "", fmt.Errorf("model not configured")
	}
	return model, nil
}

// buildPayload builds request body
//
// Parameters:
//   - params: LLM request parameters
//
// Returns:
//   - map[string]interface{}: Request body data
//   - error: Returns error if build fails
func (c *llmClient) buildPayload(params *LLMRequestParams) (map[string]interface{}, error) {
	model, err := c.buildModel()
	if err != nil {
		return nil, err
	}

	// If model is specified in params, prioritize using the one in params
	if params.Model != "" {
		model = params.Model
	}

	payload := map[string]interface{}{
		"model":  model,
		"stream": false, // Explicitly set to false to avoid proxy server using streaming response causing timeout
		"messages": []map[string]interface{}{
			{
				"role":    "system",
				"content": params.SystemPrompt,
			},
			{
				"role":    "user",
				"content": params.UserPrompt,
			},
		},
		"temperature": params.Temperature,
	}

	// Only add max_tokens to request body if it has a value
	// If nil, don't include this field, let API use model's default maximum value
	if params.MaxTokens != nil {
		payload["max_tokens"] = *params.MaxTokens
	}

	return payload, nil
}

// extractContent extracts content from response
//
// Uses OpenAI standard format to parse response and extract message content.
// Supports all response formats that follow OpenAI Chat Completions API standard.
//
// Parameters:
//   - response: JSON response data (map[string]interface{})
//
// Returns:
//   - string: Extracted content (trimmed of leading and trailing whitespace)
//   - error: Returns error if response format is incorrect or content is empty
func (c *llmClient) extractContent(response map[string]interface{}) (string, error) {
	logger := logging.GetLogger()

	// Parse as standard struct
	// First convert map to JSON string, then deserialize
	jsonBytes, err := json.Marshal(response)
	if err != nil {
		logger.WithError(err).Error("Failed to serialize LLM response to JSON string")
		return "", fmt.Errorf("failed to serialize response to JSON string: %w", err)
	}

	var completion ChatCompletionResponse
	if err := json.Unmarshal(jsonBytes, &completion); err != nil {
		logger.WithError(err).Error("Failed to parse LLM response as OpenAI ChatCompletion format")
		return "", fmt.Errorf("failed to parse response as OpenAI ChatCompletion format: %w", err)
	}

	// Extract content
	if len(completion.Choices) == 0 {
		logger.WithField("response", response).Error("LLM response has no choices")
		return "", fmt.Errorf("response has no choices array or array is empty")
	}

	choice := completion.Choices[0]
	if choice.Message.Content == nil {
		logger.WithField("choice", choice).Error("LLM response content is empty")
		return "", fmt.Errorf("response content is empty")
	}

	content := strings.TrimSpace(*choice.Message.Content)
	if content == "" {
		logger.WithField("choice", choice).Error("LLM response content is empty string")
		return "", fmt.Errorf("response content is empty string")
	}

	return content, nil
}

package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// ==================== CurrentProvider Tests ====================

func TestLLMConfig_CurrentProvider(t *testing.T) {
	tests := []struct {
		name       string
		config     LLMConfig
		wantAPIKey string
		wantModel  string
		wantURL    string
		wantErr    bool
		errMsg     string
	}{
		{
			name: "OpenAI provider - complete configuration",
			config: LLMConfig{
				Provider: "openai",
				OpenAI: struct {
					APIKey string `toml:"api_key,omitempty"`
					Model  string `toml:"model,omitempty"`
				}{
					APIKey: "sk-test-key",
					Model:  "gpt-4",
				},
			},
			wantAPIKey: "sk-test-key",
			wantModel:  "gpt-4",
			wantURL:    "https://api.openai.com/v1",
			wantErr:    false,
		},
		{
			name: "OpenAI provider - use default model",
			config: LLMConfig{
				Provider: "openai",
				OpenAI: struct {
					APIKey string `toml:"api_key,omitempty"`
					Model  string `toml:"model,omitempty"`
				}{
					APIKey: "sk-test-key",
					Model:  "", // Empty model, should use default value
				},
			},
			wantAPIKey: "sk-test-key",
			wantModel:  "gpt-3.5-turbo", // Default value
			wantURL:    "https://api.openai.com/v1",
			wantErr:    false,
		},
		{
			name: "DeepSeek provider - complete configuration",
			config: LLMConfig{
				Provider: "deepseek",
				DeepSeek: struct {
					APIKey string `toml:"api_key,omitempty"`
					Model  string `toml:"model,omitempty"`
				}{
					APIKey: "sk-deepseek-key",
					Model:  "deepseek-chat-v2",
				},
			},
			wantAPIKey: "sk-deepseek-key",
			wantModel:  "deepseek-chat-v2",
			wantURL:    "https://api.deepseek.com/v1",
			wantErr:    false,
		},
		{
			name: "DeepSeek provider - use default model",
			config: LLMConfig{
				Provider: "deepseek",
				DeepSeek: struct {
					APIKey string `toml:"api_key,omitempty"`
					Model  string `toml:"model,omitempty"`
				}{
					APIKey: "sk-deepseek-key",
					Model:  "", // Empty model, should use default value
				},
			},
			wantAPIKey: "sk-deepseek-key",
			wantModel:  "deepseek-chat", // Default value
			wantURL:    "https://api.deepseek.com/v1",
			wantErr:    false,
		},
		{
			name: "Proxy provider - complete configuration",
			config: LLMConfig{
				Provider: "proxy",
				Proxy: struct {
					URL    string `toml:"url,omitempty"`
					APIKey string `toml:"api_key,omitempty"`
					Model  string `toml:"model,omitempty"`
				}{
					URL:    "https://api.example.com/v1",
					APIKey: "proxy-key",
					Model:  "custom-model",
				},
			},
			wantAPIKey: "proxy-key",
			wantModel:  "custom-model",
			wantURL:    "https://api.example.com/v1",
			wantErr:    false,
		},
		{
			name: "Proxy provider - missing model",
			config: LLMConfig{
				Provider: "proxy",
				Proxy: struct {
					URL    string `toml:"url,omitempty"`
					APIKey string `toml:"api_key,omitempty"`
					Model  string `toml:"model,omitempty"`
				}{
					URL:    "https://api.example.com/v1",
					APIKey: "proxy-key",
					Model:  "", // Missing model
				},
			},
			wantAPIKey: "",
			wantModel:  "",
			wantURL:    "",
			wantErr:    true,
			errMsg:     "model is required for proxy provider",
		},
		{
			name: "Proxy provider - missing URL",
			config: LLMConfig{
				Provider: "proxy",
				Proxy: struct {
					URL    string `toml:"url,omitempty"`
					APIKey string `toml:"api_key,omitempty"`
					Model  string `toml:"model,omitempty"`
				}{
					URL:    "", // Missing URL
					APIKey: "proxy-key",
					Model:  "custom-model",
				},
			},
			wantAPIKey: "",
			wantModel:  "",
			wantURL:    "",
			wantErr:    true,
			errMsg:     "URL is required for proxy provider",
		},
		{
			name: "Provider not configured",
			config: LLMConfig{
				Provider: "", // Empty provider
			},
			wantAPIKey: "",
			wantModel:  "",
			wantURL:    "",
			wantErr:    true,
			errMsg:     "LLM provider is not configured",
		},
		{
			name: "Invalid provider",
			config: LLMConfig{
				Provider: "invalid-provider",
			},
			wantAPIKey: "",
			wantModel:  "",
			wantURL:    "",
			wantErr:    true,
			errMsg:     "unsupported LLM provider: invalid-provider",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Act: Get current provider configuration
			apiKey, model, url, err := tt.config.CurrentProvider()

			// Assert: Verify results
			if tt.wantErr {
				assert.Error(t, err)
				if tt.errMsg != "" {
					assert.Contains(t, err.Error(), tt.errMsg)
				}
				assert.Empty(t, apiKey)
				assert.Empty(t, model)
				assert.Empty(t, url)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.wantAPIKey, apiKey)
				assert.Equal(t, tt.wantModel, model)
				assert.Equal(t, tt.wantURL, url)
			}
		})
	}
}

// ==================== Edge Case Tests ====================

func TestLLMConfig_CurrentProvider_EdgeCases(t *testing.T) {
	t.Run("OpenAI provider - APIKey is empty", func(t *testing.T) {
		config := LLMConfig{
			Provider: "openai",
			OpenAI: struct {
				APIKey string `toml:"api_key,omitempty"`
				Model  string `toml:"model,omitempty"`
			}{
				APIKey: "", // Empty APIKey (allowed, validated by caller)
				Model:  "gpt-4",
			},
		}

		apiKey, model, url, err := config.CurrentProvider()
		assert.NoError(t, err)
		assert.Empty(t, apiKey)
		assert.Equal(t, "gpt-4", model)
		assert.Equal(t, "https://api.openai.com/v1", url)
	})

	t.Run("DeepSeek provider - APIKey is empty", func(t *testing.T) {
		config := LLMConfig{
			Provider: "deepseek",
			DeepSeek: struct {
				APIKey string `toml:"api_key,omitempty"`
				Model  string `toml:"model,omitempty"`
			}{
				APIKey: "", // Empty APIKey (allowed, validated by caller)
				Model:  "deepseek-chat-v2",
			},
		}

		apiKey, model, url, err := config.CurrentProvider()
		assert.NoError(t, err)
		assert.Empty(t, apiKey)
		assert.Equal(t, "deepseek-chat-v2", model)
		assert.Equal(t, "https://api.deepseek.com/v1", url)
	})

	t.Run("Proxy provider - APIKey is empty but model and URL exist", func(t *testing.T) {
		config := LLMConfig{
			Provider: "proxy",
			Proxy: struct {
				URL    string `toml:"url,omitempty"`
				APIKey string `toml:"api_key,omitempty"`
				Model  string `toml:"model,omitempty"`
			}{
				URL:    "https://api.example.com/v1",
				APIKey: "", // Empty APIKey (allowed, validated by caller)
				Model:  "custom-model",
			},
		}

		apiKey, model, url, err := config.CurrentProvider()
		assert.NoError(t, err)
		assert.Empty(t, apiKey)
		assert.Equal(t, "custom-model", model)
		assert.Equal(t, "https://api.example.com/v1", url)
	})
}

// ==================== CurrentLanguage Tests ====================

func TestLLMConfig_CurrentLanguage(t *testing.T) {
	tests := []struct {
		name     string
		config   LLMConfig
		wantCode string
		wantErr  bool
		errMsg   string
	}{
		{
			name: "Language not set - returns default English",
			config: LLMConfig{
				Language: "", // Language not set
			},
			wantCode: "en",
			wantErr:  false,
		},
		{
			name: "Set English language",
			config: LLMConfig{
				Language: "en",
			},
			wantCode: "en",
			wantErr:  false,
		},
		{
			name: "Set Simplified Chinese",
			config: LLMConfig{
				Language: "zh-CN",
			},
			wantCode: "zh-CN",
			wantErr:  false,
		},
		{
			name: "Set Traditional Chinese",
			config: LLMConfig{
				Language: "zh-TW",
			},
			wantCode: "zh-TW",
			wantErr:  false,
		},
		{
			name: "Set Japanese",
			config: LLMConfig{
				Language: "ja",
			},
			wantCode: "ja",
			wantErr:  false,
		},
		{
			name: "Set Korean",
			config: LLMConfig{
				Language: "ko",
			},
			wantCode: "ko",
			wantErr:  false,
		},
		{
			name: "Invalid language code",
			config: LLMConfig{
				Language: "invalid-lang",
			},
			wantCode: "",
			wantErr:  true,
			errMsg:   "unsupported language code",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Act: Get current language configuration
			lang, err := tt.config.CurrentLanguage()

			// Assert: Verify results
			if tt.wantErr {
				assert.Error(t, err)
				if tt.errMsg != "" {
					assert.Contains(t, err.Error(), tt.errMsg)
				}
				assert.Nil(t, lang)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, lang)
				assert.Equal(t, tt.wantCode, lang.Code)
			}
		})
	}
}

func TestLLMConfig_CurrentLanguage_DefaultEnglish(t *testing.T) {
	// Arrange: Configuration with language not set
	config := LLMConfig{
		Language: "",
	}

	// Act: Get current language configuration
	lang, err := config.CurrentLanguage()

	// Assert: Should return default English configuration
	assert.NoError(t, err)
	assert.NotNil(t, lang)
	assert.Equal(t, "en", lang.Code)
	assert.Equal(t, "English", lang.Name)
	assert.Equal(t, "English", lang.NativeName)
}

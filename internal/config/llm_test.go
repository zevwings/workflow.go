package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// ==================== CurrentProvider 测试 ====================

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
			name: "OpenAI provider - 完整配置",
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
			name: "OpenAI provider - 使用默认 model",
			config: LLMConfig{
				Provider: "openai",
				OpenAI: struct {
					APIKey string `toml:"api_key,omitempty"`
					Model  string `toml:"model,omitempty"`
				}{
					APIKey: "sk-test-key",
					Model:  "", // 空 model，应该使用默认值
				},
			},
			wantAPIKey: "sk-test-key",
			wantModel:  "gpt-3.5-turbo", // 默认值
			wantURL:    "https://api.openai.com/v1",
			wantErr:    false,
		},
		{
			name: "DeepSeek provider - 完整配置",
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
			name: "DeepSeek provider - 使用默认 model",
			config: LLMConfig{
				Provider: "deepseek",
				DeepSeek: struct {
					APIKey string `toml:"api_key,omitempty"`
					Model  string `toml:"model,omitempty"`
				}{
					APIKey: "sk-deepseek-key",
					Model:  "", // 空 model，应该使用默认值
				},
			},
			wantAPIKey: "sk-deepseek-key",
			wantModel:  "deepseek-chat", // 默认值
			wantURL:    "https://api.deepseek.com/v1",
			wantErr:    false,
		},
		{
			name: "Proxy provider - 完整配置",
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
			name: "Proxy provider - 缺少 model",
			config: LLMConfig{
				Provider: "proxy",
				Proxy: struct {
					URL    string `toml:"url,omitempty"`
					APIKey string `toml:"api_key,omitempty"`
					Model  string `toml:"model,omitempty"`
				}{
					URL:    "https://api.example.com/v1",
					APIKey: "proxy-key",
					Model:  "", // 缺少 model
				},
			},
			wantAPIKey: "",
			wantModel:  "",
			wantURL:    "",
			wantErr:    true,
			errMsg:     "model is required for proxy provider",
		},
		{
			name: "Proxy provider - 缺少 URL",
			config: LLMConfig{
				Provider: "proxy",
				Proxy: struct {
					URL    string `toml:"url,omitempty"`
					APIKey string `toml:"api_key,omitempty"`
					Model  string `toml:"model,omitempty"`
				}{
					URL:    "", // 缺少 URL
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
			name: "Provider 未配置",
			config: LLMConfig{
				Provider: "", // 空 provider
			},
			wantAPIKey: "",
			wantModel:  "",
			wantURL:    "",
			wantErr:    true,
			errMsg:     "LLM provider is not configured",
		},
		{
			name: "无效的 provider",
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
			// Act: 获取当前 provider 配置
			apiKey, model, url, err := tt.config.CurrentProvider()

			// Assert: 验证结果
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

// ==================== 边界情况测试 ====================

func TestLLMConfig_CurrentProvider_EdgeCases(t *testing.T) {
	t.Run("OpenAI provider - APIKey 为空", func(t *testing.T) {
		config := LLMConfig{
			Provider: "openai",
			OpenAI: struct {
				APIKey string `toml:"api_key,omitempty"`
				Model  string `toml:"model,omitempty"`
			}{
				APIKey: "", // 空 APIKey（允许，由调用方验证）
				Model:  "gpt-4",
			},
		}

		apiKey, model, url, err := config.CurrentProvider()
		assert.NoError(t, err)
		assert.Empty(t, apiKey)
		assert.Equal(t, "gpt-4", model)
		assert.Equal(t, "https://api.openai.com/v1", url)
	})

	t.Run("DeepSeek provider - APIKey 为空", func(t *testing.T) {
		config := LLMConfig{
			Provider: "deepseek",
			DeepSeek: struct {
				APIKey string `toml:"api_key,omitempty"`
				Model  string `toml:"model,omitempty"`
			}{
				APIKey: "", // 空 APIKey（允许，由调用方验证）
				Model:  "deepseek-chat-v2",
			},
		}

		apiKey, model, url, err := config.CurrentProvider()
		assert.NoError(t, err)
		assert.Empty(t, apiKey)
		assert.Equal(t, "deepseek-chat-v2", model)
		assert.Equal(t, "https://api.deepseek.com/v1", url)
	})

	t.Run("Proxy provider - APIKey 为空但 model 和 URL 存在", func(t *testing.T) {
		config := LLMConfig{
			Provider: "proxy",
			Proxy: struct {
				URL    string `toml:"url,omitempty"`
				APIKey string `toml:"api_key,omitempty"`
				Model  string `toml:"model,omitempty"`
			}{
				URL:    "https://api.example.com/v1",
				APIKey: "", // 空 APIKey（允许，由调用方验证）
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

// ==================== CurrentLanguage 测试 ====================

func TestLLMConfig_CurrentLanguage(t *testing.T) {
	tests := []struct {
		name     string
		config   LLMConfig
		wantCode string
		wantErr  bool
		errMsg   string
	}{
		{
			name: "语言未设置 - 返回默认英文",
			config: LLMConfig{
				Language: "", // 未设置语言
			},
			wantCode: "en",
			wantErr:  false,
		},
		{
			name: "设置英文语言",
			config: LLMConfig{
				Language: "en",
			},
			wantCode: "en",
			wantErr:  false,
		},
		{
			name: "设置简体中文",
			config: LLMConfig{
				Language: "zh-CN",
			},
			wantCode: "zh-CN",
			wantErr:  false,
		},
		{
			name: "设置繁体中文",
			config: LLMConfig{
				Language: "zh-TW",
			},
			wantCode: "zh-TW",
			wantErr:  false,
		},
		{
			name: "设置日语",
			config: LLMConfig{
				Language: "ja",
			},
			wantCode: "ja",
			wantErr:  false,
		},
		{
			name: "设置韩语",
			config: LLMConfig{
				Language: "ko",
			},
			wantCode: "ko",
			wantErr:  false,
		},
		{
			name: "无效的语言代码",
			config: LLMConfig{
				Language: "invalid-lang",
			},
			wantCode: "",
			wantErr:  true,
			errMsg:   "不支持的语言代码",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Act: 获取当前语言配置
			lang, err := tt.config.CurrentLanguage()

			// Assert: 验证结果
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
	// Arrange: 语言未设置的配置
	config := LLMConfig{
		Language: "",
	}

	// Act: 获取当前语言配置
	lang, err := config.CurrentLanguage()

	// Assert: 应该返回默认英文配置
	assert.NoError(t, err)
	assert.NotNil(t, lang)
	assert.Equal(t, "en", lang.Code)
	assert.Equal(t, "English", lang.Name)
	assert.Equal(t, "English", lang.NativeName)
}

package config

import "fmt"

// LLMConfig LLM 配置
type LLMConfig struct {
	Provider string `toml:"provider,omitempty"`
	Language string `toml:"language,omitempty"`
	OpenAI   struct {
		APIKey string `toml:"api_key,omitempty"`
		Model  string `toml:"model,omitempty"`
	} `toml:"openai,omitempty"`
	DeepSeek struct {
		APIKey string `toml:"api_key,omitempty"`
		Model  string `toml:"model,omitempty"`
	} `toml:"deepseek,omitempty"`
	Proxy struct {
		URL    string `toml:"url,omitempty"`
		APIKey string `toml:"api_key,omitempty"`
		Model  string `toml:"model,omitempty"`
	} `toml:"proxy,omitempty"`
}

// CurrentProvider 获取当前 provider 的配置
//
// 根据 LLMConfig.Provider 返回对应的配置信息（APIKey、Model、URL）。
// 如果 provider 未设置或不存在，返回错误。
//
// 返回:
//   - APIKey: API 密钥
//   - Model: 模型名称（如果未设置，返回默认值）
//   - URL: API URL（仅 proxy provider 需要）
//   - error: 如果 provider 未配置或无效，返回错误
func (c *LLMConfig) CurrentProvider() (apiKey, model, url string, err error) {
	switch c.Provider {
	case "openai":
		apiKey = c.OpenAI.APIKey
		model = c.OpenAI.Model
		if model == "" {
			model = "gpt-3.5-turbo"
		}
		return apiKey, model, "", nil
	case "deepseek":
		apiKey = c.DeepSeek.APIKey
		model = c.DeepSeek.Model
		if model == "" {
			model = "deepseek-chat"
		}
		return apiKey, model, "", nil
	case "proxy":
		apiKey = c.Proxy.APIKey
		model = c.Proxy.Model
		url = c.Proxy.URL
		if model == "" {
			return "", "", "", fmt.Errorf("model is required for proxy provider")
		}
		if url == "" {
			return "", "", "", fmt.Errorf("URL is required for proxy provider")
		}
		return apiKey, model, url, nil
	default:
		if c.Provider == "" {
			return "", "", "", fmt.Errorf("LLM provider is not configured")
		}
		return "", "", "", fmt.Errorf("unsupported LLM provider: %s", c.Provider)
	}
}


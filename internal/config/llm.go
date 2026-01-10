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
//   - URL: API URL（openai/deepseek 使用默认 URL，proxy 需要配置）
//   - error: 如果 provider 未配置或无效，返回错误
func (c *LLMConfig) CurrentProvider() (apiKey, model, url string, err error) {
	switch c.Provider {
	case "openai":
		apiKey = c.OpenAI.APIKey
		model = c.OpenAI.Model
		if model == "" {
			model = "gpt-3.5-turbo"
		}
		return apiKey, model, "https://api.openai.com/v1", nil
	case "deepseek":
		apiKey = c.DeepSeek.APIKey
		model = c.DeepSeek.Model
		if model == "" {
			model = "deepseek-chat"
		}
		return apiKey, model, "https://api.deepseek.com/v1", nil
	case "proxy":
		apiKey = c.Proxy.APIKey
		model = c.Proxy.Model
		url = c.Proxy.URL
		if model == "" || url == "" || apiKey == "" {
			return "", "", "", fmt.Errorf("model, URL and API key are required for proxy provider")
		}
		return apiKey, model, url, nil
	default:
		return "", "", "", fmt.Errorf("unsupported LLM provider: %s", c.Provider)
	}
}

// CurrentLanguage 获取当前语言的配置
//
// 根据 LLMConfig.Language 返回对应的语言配置信息。
// 如果 language 未设置，返回默认英文配置。
// 如果 language 设置了但找不到匹配的语言，返回错误。
//
// 返回:
//   - *SupportedLanguage: 语言配置信息（如果未设置，返回默认英文配置）
//   - error: 如果 language 配置了但无效，返回错误
func (c *LLMConfig) CurrentLanguage() (*SupportedLanguage, error) {
	// 如果语言未设置，返回默认英文配置
	if c.Language == "" {
		lang := FindLanguage("en")
		if lang == nil {
			return nil, fmt.Errorf("默认语言（英文）配置不可用")
		}
		return lang, nil
	}

	// 查找匹配的语言
	lang := FindLanguage(c.Language)
	if lang == nil {
		return nil, fmt.Errorf("不支持的语言代码: %s", c.Language)
	}

	return lang, nil
}

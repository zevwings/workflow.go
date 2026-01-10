package config

import "fmt"

// LLMConfig LLM configuration
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

// CurrentProvider gets current provider configuration
//
// Returns corresponding configuration information (APIKey, Model, URL) based on LLMConfig.Provider.
// If provider is not set or does not exist, returns error.
//
// Returns:
//   - APIKey: API key
//   - Model: Model name (if not set, returns default value)
//   - URL: API URL (openai/deepseek use default URL, proxy needs configuration)
//   - error: Returns error if provider is not configured or invalid
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

// CurrentLanguage gets current language configuration
//
// Returns corresponding language configuration information based on LLMConfig.Language.
// If language is not set, returns default English configuration.
// If language is set but no matching language is found, returns error.
//
// Returns:
//   - *SupportedLanguage: Language configuration information (if not set, returns default English configuration)
//   - error: Returns error if language is configured but invalid
func (c *LLMConfig) CurrentLanguage() (*SupportedLanguage, error) {
	// If language is not set, return default English configuration
	if c.Language == "" {
		lang := FindLanguage("en")
		if lang == nil {
			return nil, fmt.Errorf("default language (English) configuration is not available")
		}
		return lang, nil
	}

	// Find matching language
	lang := FindLanguage(c.Language)
	if lang == nil {
		return nil, fmt.Errorf("unsupported language code: %s", c.Language)
	}

	return lang, nil
}

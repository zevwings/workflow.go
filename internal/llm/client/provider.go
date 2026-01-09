package client

// ProviderConfig 提供商配置
//
// 用于 LLM 客户端的基础配置，包含 API 密钥、模型名称和 URL。
type ProviderConfig struct {
	// APIKey API 密钥
	APIKey string
	// Model 模型名称
	Model string
	// URL API URL（仅 proxy provider 需要，openai/deepseek 使用固定 URL）
	URL string
}

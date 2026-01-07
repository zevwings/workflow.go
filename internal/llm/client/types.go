package client

// LLMRequestParams LLM 请求参数
//
// 包含调用 LLM API 所需的所有参数。
type LLMRequestParams struct {
	// SystemPrompt 系统提示词
	SystemPrompt string `json:"system_prompt"`
	// UserPrompt 用户提示词
	UserPrompt string `json:"user_prompt"`
	// MaxTokens 最大 token 数（nil 表示不限制，使用模型默认最大值）
	MaxTokens *int `json:"max_tokens,omitempty"`
	// Temperature 温度参数（控制输出的随机性，范围 0.0-2.0）
	Temperature float32 `json:"temperature"`
	// Model 模型名称（可选，如果为空则从配置获取）
	Model string `json:"model,omitempty"`
}

// DefaultLLMRequestParams 返回默认的 LLM 请求参数
//
// 返回:
//   - *LLMRequestParams: 使用默认值的请求参数
func DefaultLLMRequestParams() *LLMRequestParams {
	temperature := float32(0.5)
	return &LLMRequestParams{
		SystemPrompt: "",
		UserPrompt:   "",
		MaxTokens:    nil,
		Temperature:  temperature,
		Model:        "",
	}
}

// ChatCompletionResponse OpenAI Chat Completions API 响应
//
// 完整的 OpenAI 标准响应格式，支持所有标准字段和扩展字段。
type ChatCompletionResponse struct {
	// ID 响应唯一标识符
	ID string `json:"id"`
	// Object 对象类型，固定为 "chat.completion"
	Object string `json:"object"`
	// Created 创建时间戳（Unix 时间戳）
	Created int64 `json:"created"`
	// Model 使用的模型名称
	Model string `json:"model"`
	// SystemFingerprint 系统指纹（可选）
	SystemFingerprint *string `json:"system_fingerprint,omitempty"`
	// Choices 选择列表
	Choices []ChatCompletionChoice `json:"choices"`
	// Usage Token 使用统计
	Usage Usage `json:"usage"`
}

// ChatCompletionChoice Chat Completion 选择项
type ChatCompletionChoice struct {
	// Index 选择索引
	Index int `json:"index"`
	// Message 消息对象
	Message ChatMessage `json:"message"`
	// Logprobs 对数概率（可选）
	Logprobs interface{} `json:"logprobs,omitempty"`
	// FinishReason 完成原因
	FinishReason string `json:"finish_reason"`
}

// ChatMessage Chat 消息
type ChatMessage struct {
	// Role 消息角色
	Role string `json:"role"`
	// Content 消息内容（可能为 null）
	Content *string `json:"content,omitempty"`
}

// Usage Token 使用统计
type Usage struct {
	// PromptTokens 提示词 token 数
	PromptTokens int `json:"prompt_tokens"`
	// CompletionTokens 完成 token 数
	CompletionTokens int `json:"completion_tokens"`
	// TotalTokens 总 token 数
	TotalTokens int `json:"total_tokens"`
}

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

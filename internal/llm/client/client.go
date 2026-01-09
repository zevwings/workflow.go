// Package client 提供了统一配置驱动的 LLM 客户端实现，支持 OpenAI、DeepSeek 和代理 API。
//
// 主要功能：
//   - LLMClient: LLM 客户端，支持多种提供商
//   - LLMRequestParams: LLM 请求参数
//   - 语言支持：支持多种语言的 prompt 增强
//   - 配置驱动：通过配置结构体传入配置
//
// 使用示例：
//
//	import (
//		"github.com/zevwings/workflow/internal/http"
//		"github.com/zevwings/workflow/internal/llm/client"
//	)
//
//	// 方式 1: 使用便捷函数创建配置（推荐）
//	cfg := client.NewOpenAIConfig("your-api-key", "gpt-3.5-turbo")
//	llmClient := client.Global(cfg.OpenAI) // 自动使用 http.Global()
//
//	// 方式 2: 手动创建配置
//	cfg := &client.ProviderConfig{
//		APIKey: "your-api-key",
//		Model:  "gpt-3.5-turbo",
//		URL:    "https://api.openai.com/v1/chat/completions",
//	}
//	llmClient := client.Global(cfg) // 自动使用 http.Global()
//
//	// 准备请求参数
//	params := &client.LLMRequestParams{
//		SystemPrompt: "You are a helpful assistant.",
//		UserPrompt:   "What is Go?",
//		Temperature:  0.7,
//	}
//
//	// 调用 LLM API
//	response, err := llmClient.Call(params)
//	if err != nil {
//		// 处理错误
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
)

var (
	// globalClient 全局 LLM 客户端单例
	globalClient *llmClient
	globalOnce   sync.Once
)

// LLMClient LLM 客户端接口
//
// 所有 LLM 提供商使用同一个客户端实现，通过配置结构体区分不同的提供商。
// 所有配置（URL、API key、model、response_format）都从配置结构体获取。
type LLMClient interface {
	// Call 调用 LLM API
	//
	// 参数:
	//   - params: LLM 请求参数
	//
	// 返回:
	//   - string: LLM 生成的文本内容（去除首尾空白）
	//   - error: 如果 API 调用失败或响应格式不正确，返回相应的错误信息
	Call(params *LLMRequestParams) (string, error)
}

// llmClient LLM 客户端实现
//
// 所有 LLM 提供商使用同一个客户端实现，通过配置结构体区分不同的提供商。
// 所有配置（URL、API key、model、response_format）都从配置结构体获取。
type llmClient struct {
	httpClient http.Client
	config     *ProviderConfig
}

// newClient 创建新的 LLM 客户端（私有函数）
//
// 自动使用 http.Global() 获取全局 HTTP 客户端。
//
// 参数:
//   - config: LLM 配置结构体（不能为 nil）
//
// 返回:
//   - LLMClient: LLM 客户端实例
func newClient(config *ProviderConfig) LLMClient {
	if config == nil {
		panic("llm/client.newClient: config 不能为 nil")
	}
	return &llmClient{
		httpClient: http.Global(),
		config:     config,
	}
}

// Global 获取全局 LLMClient 单例
//
// 返回进程级别的 LLMClient 单例。
// 单例会在首次调用时初始化，后续调用会复用同一个实例。
// 自动使用 http.Global() 获取全局 HTTP 客户端。
//
// 参数:
//   - config: LLM 配置结构体（必须，不能为 nil）
//
// 返回:
//   - LLMClient: LLM 客户端接口实例
//
// 注意:
//   - 配置必须由调用者提供，LLM 模块不负责配置的创建和生命周期
//   - 首次调用时传入的参数会被保存，后续调用会忽略参数
//   - 如果 config 为 nil，会在首次调用时 panic
//   - 返回的是接口类型，外部无法直接访问结构体的字段或未导出的方法
//
// 优势:
//   - 减少资源消耗：避免重复创建客户端实例
//   - 线程安全：可以在多线程环境中安全使用
//   - 统一管理：所有 LLM 调用使用同一个客户端实例
//   - 便捷性：自动使用全局 HTTP 客户端，无需手动传入
//   - 封装性：返回接口类型，隐藏实现细节，防止外部直接访问内部结构
func Global(config *ProviderConfig) LLMClient {
	if config == nil {
		panic("llm/client.Global: config 不能为 nil")
	}
	globalOnce.Do(func() {
		globalClient = newClient(config).(*llmClient)
	})
	return globalClient
}

// Call 调用 LLM API
//
// 参数:
//   - params: LLM 请求参数
//
// 返回:
//   - string: LLM 生成的文本内容（去除首尾空白）
//   - error: 如果 API 调用失败或响应格式不正确，返回相应的错误信息
func (c *llmClient) Call(params *LLMRequestParams) (string, error) {
	// 构建请求体（统一格式）
	payload, err := c.buildPayload(params)
	if err != nil {
		return "", fmt.Errorf("构建请求体失败: %w", err)
	}

	// 构建请求头（统一格式）
	headers, err := c.buildHeaders()
	if err != nil {
		return "", fmt.Errorf("构建请求头失败: %w", err)
	}

	// 构建 URL（统一格式）
	url, err := c.buildURL()
	if err != nil {
		return "", fmt.Errorf("构建 URL 失败: %w", err)
	}

	// 构建请求配置
	reqConfig := http.NewRequestConfig().
		WithHeaders(headers).
		WithBody(payload).
		WithTimeout(60 * time.Second).                     // LLM API 通常需要更长的超时时间
		WithRetry(http.NewRetryConfig().WithRetryCount(3)) // 最多重试 3 次

	// 发送请求
	resp, err := c.httpClient.PostWithConfig(url, reqConfig)
	if err != nil {
		return "", fmt.Errorf("发送 LLM 请求到 %s 失败: %w", url, err)
	}

	// 检查错误（使用 EnsureSuccessWith 统一处理）
	resp, err = resp.EnsureSuccessWith(func(r *http.HttpResponse) error {
		errorMessage := r.ExtractErrorMessage()
		return fmt.Errorf("LLM API 请求失败 (%s): %d - %s", url, r.Status, errorMessage)
	})
	if err != nil {
		return "", err
	}

	// 解析 JSON 响应
	var data map[string]interface{}
	data, err = http.AsJSON[map[string]interface{}](resp)
	if err != nil {
		return "", fmt.Errorf("解析响应 JSON 失败: %w", err)
	}

	// 根据配置的响应格式提取内容
	content, err := c.extractContent(data)
	if err != nil {
		return "", fmt.Errorf("提取响应内容失败: %w", err)
	}

	return content, nil
}

// buildURL 构建 API URL
//
// 直接从配置结构体获取 URL。
//
// 返回:
//   - string: API URL
//   - error: 如果 provider 未配置或无效，返回错误
func (c *llmClient) buildURL() (string, error) {
	url := c.config.URL
	if url == "" {
		return "", fmt.Errorf("URL 未配置")
	}
	return url, nil
}

// buildHeaders 构建请求头
//
// 返回:
//   - map[string]string: 请求头映射
//   - error: 如果 API key 未配置，返回错误
func (c *llmClient) buildHeaders() (map[string]string, error) {
	headers := make(map[string]string)

	var apiKey string = c.config.APIKey

	if apiKey == "" {
		return nil, fmt.Errorf("LLM API key 未配置")
	}

	headers["Authorization"] = fmt.Sprintf("Bearer %s", apiKey)
	headers["Content-Type"] = "application/json"

	return headers, nil
}

// buildModel 构建模型名称
//
// 直接从配置结构体获取模型名称。
//
// 返回:
//   - string: 模型名称
//   - error: 如果 model 未配置，返回错误
func (c *llmClient) buildModel() (string, error) {
	var model = c.config.Model
	if model == "" {
		return "", fmt.Errorf("model 未配置")
	}
	return model, nil
}

// buildPayload 构建请求体
//
// 参数:
//   - params: LLM 请求参数
//
// 返回:
//   - map[string]interface{}: 请求体数据
//   - error: 如果构建失败，返回错误
func (c *llmClient) buildPayload(params *LLMRequestParams) (map[string]interface{}, error) {
	model, err := c.buildModel()
	if err != nil {
		return nil, err
	}

	// 如果 params 中指定了 model，优先使用 params 中的
	if params.Model != "" {
		model = params.Model
	}

	payload := map[string]interface{}{
		"model": model,
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

	// 只有当 max_tokens 有值时才添加到请求体中
	// 如果为 nil，则不包含该字段，让 API 使用模型默认的最大值
	if params.MaxTokens != nil {
		payload["max_tokens"] = *params.MaxTokens
	}

	return payload, nil
}

// extractContent 从响应中提取内容
//
// 使用 OpenAI 标准格式解析响应，提取消息内容。
// 支持所有遵循 OpenAI Chat Completions API 标准的响应格式。
//
// 参数:
//   - response: JSON 响应数据（map[string]interface{}）
//
// 返回:
//   - string: 提取的内容（去除首尾空白）
//   - error: 如果响应格式不正确或内容为空，返回错误
func (c *llmClient) extractContent(response map[string]interface{}) (string, error) {
	// 解析为标准结构体
	// 先将 map 转换为 JSON 字符串，再反序列化
	jsonBytes, err := json.Marshal(response)
	if err != nil {
		return "", fmt.Errorf("序列化响应为 JSON 字符串失败: %w", err)
	}

	var completion ChatCompletionResponse
	if err := json.Unmarshal(jsonBytes, &completion); err != nil {
		return "", fmt.Errorf("解析响应为 OpenAI ChatCompletion 格式失败: %w", err)
	}

	// 提取内容
	if len(completion.Choices) == 0 {
		return "", fmt.Errorf("响应中没有 choices 数组或数组为空")
	}

	choice := completion.Choices[0]
	if choice.Message.Content == nil {
		return "", fmt.Errorf("响应中 content 为空")
	}

	content := strings.TrimSpace(*choice.Message.Content)
	if content == "" {
		return "", fmt.Errorf("响应中 content 为空字符串")
	}

	return content, nil
}

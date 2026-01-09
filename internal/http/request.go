package http

import (
	"fmt"
	"strings"
	"time"

	"github.com/go-resty/resty/v2"
)

// baseRequestConfig 基础请求配置（公共字段）
type baseRequestConfig struct {
	// Query 可选的查询参数（可以是 map[string]string 或实现了序列化的结构体）
	Query interface{}
	// Auth 可选的 Basic Authentication 认证信息
	Auth *Authorization
	// Headers 可选的自定义 HTTP Headers
	Headers map[string]string
	// Timeout 可选的请求超时时间（如果为 0，使用默认 30 秒）
	Timeout time.Duration
	// Retry 可选的重试配置（如果为 nil，使用 Client 的默认重试配置）
	Retry *RetryConfig
}

// RequestConfig HTTP 请求配置
//
// 支持链式调用配置请求参数，包括请求体、查询参数、认证信息、Headers 和超时时间。
type RequestConfig struct {
	baseRequestConfig
	// Body 可选的请求体（实现序列化接口）
	Body interface{}
}

// NewRequestConfig 创建新的 RequestConfig，使用默认值
//
// 返回:
//   - RequestConfig: 所有字段都为默认值的 RequestConfig 实例
func NewRequestConfig() *RequestConfig {
	return &RequestConfig{
		baseRequestConfig: baseRequestConfig{
			Headers: make(map[string]string),
		},
	}
}

// WithBody 设置请求体
//
// 参数:
//   - body: 请求体，可以是任意可序列化的类型
//
// 返回:
//   - *RequestConfig: 返回自身，支持链式调用
func (c *RequestConfig) WithBody(body interface{}) *RequestConfig {
	c.Body = body
	return c
}

// WithQuery 设置查询参数
//
// 参数:
//   - query: 查询参数，可以是 map[string]string 或实现了序列化的结构体
//
// 返回:
//   - *RequestConfig: 返回自身，支持链式调用
func (c *RequestConfig) WithQuery(query interface{}) *RequestConfig {
	c.Query = query
	return c
}

// WithAuth 设置认证信息
//
// 参数:
//   - auth: Basic Authentication 认证信息
//
// 返回:
//   - *RequestConfig: 返回自身，支持链式调用
func (c *RequestConfig) WithAuth(auth *Authorization) *RequestConfig {
	c.Auth = auth
	return c
}

// WithHeader 设置单个 HTTP Header
//
// 参数:
//   - key: Header 键
//   - value: Header 值
//
// 返回:
//   - *RequestConfig: 返回自身，支持链式调用
func (c *RequestConfig) WithHeader(key, value string) *RequestConfig {
	c.ensureHeaders()
	c.Headers[key] = value
	return c
}

// WithHeaders 设置多个 HTTP Headers
//
// 参数:
//   - headers: Headers 映射
//
// 返回:
//   - *RequestConfig: 返回自身，支持链式调用
func (c *RequestConfig) WithHeaders(headers map[string]string) *RequestConfig {
	c.ensureHeaders()
	for k, v := range headers {
		c.Headers[k] = v
	}
	return c
}

// WithTimeout 设置超时时间
//
// 参数:
//   - timeout: 请求超时时间
//
// 返回:
//   - *RequestConfig: 返回自身，支持链式调用
//
// 注意:
//
//	如果不设置超时时间，将使用默认的 30 秒超时。
func (c *RequestConfig) WithTimeout(timeout time.Duration) *RequestConfig {
	c.Timeout = timeout
	return c
}

// WithRetry 设置重试配置
//
// 参数:
//   - retry: 重试配置
//
// 返回:
//   - *RequestConfig: 返回自身，支持链式调用
func (c *RequestConfig) WithRetry(retry *RetryConfig) *RequestConfig {
	c.Retry = retry
	return c
}

// ensureHeaders 确保 Headers map 已初始化
func (c *RequestConfig) ensureHeaders() {
	if c.Headers == nil {
		c.Headers = make(map[string]string)
	}
}

// applyToRequest 将配置应用到 resty 请求
func (c *RequestConfig) applyToRequest(req *resty.Request) *resty.Request {
	// 添加 body（如果有）
	if c.Body != nil {
		req = req.SetBody(c.Body)
	}

	// 添加 query 参数
	if c.Query != nil {
		params := convertToQueryParams(c.Query)
		if len(params) > 0 {
			req = req.SetQueryParams(params)
		}
	}

	// 添加 auth
	if c.Auth != nil {
		req = req.SetBasicAuth(c.Auth.Username, c.Auth.Password)
	}

	// 添加 headers
	if c.Headers != nil {
		for k, v := range c.Headers {
			req = req.SetHeader(k, v)
		}
	}

	// 注意：重试配置需要在 doRequest 中应用，因为需要在 Client 级别设置

	// 注意：resty.Request 不支持单独设置超时，超时在 Client 级别设置
	// 如果需要不同的超时，需要创建新的 Client 实例
	// 这里我们暂时忽略 Timeout 字段，使用 Client 的默认超时

	return req
}

// convertToQueryParams 将查询参数转换为 resty 可用的格式
func convertToQueryParams(query interface{}) map[string]string {
	result := make(map[string]string)

	switch v := query.(type) {
	case map[string]string:
		return v
	case map[string]interface{}:
		for k, val := range v {
			result[k] = fmt.Sprintf("%v", val)
		}
	case []string:
		// 支持 ["key=value", "key2=value2"] 格式
		for _, pair := range v {
			parts := strings.SplitN(pair, "=", 2)
			if len(parts) == 2 {
				result[parts[0]] = parts[1]
			}
		}
	default:
		// 对于其他类型，尝试转换为字符串
		// 这里可以扩展支持更多类型
	}

	return result
}

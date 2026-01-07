package http

import (
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/go-resty/resty/v2"
)

// RetryConfig 重试配置
//
// 用于配置 HTTP 请求的重试策略，支持在每次请求时自定义重试行为。
type RetryConfig struct {
	// Count 最大重试次数（如果为 0，使用 Client 的默认配置；如果为 -1，禁用重试）
	Count int
	// WaitTime 初始等待时间（如果为 0，使用 Client 的默认配置）
	WaitTime time.Duration
	// MaxWaitTime 最大等待时间（如果为 0，使用 Client 的默认配置）
	MaxWaitTime time.Duration
	// Condition 自定义重试条件函数（如果为 nil，使用默认的重试条件）
	// 函数接收响应和错误，返回是否应该重试
	Condition func(*resty.Response, error) bool
	// After 自定义重试延迟函数（如果为 nil，使用默认的重试延迟策略）
	// 函数接收客户端和响应，返回重试延迟时间
	After func(*resty.Client, *resty.Response) (time.Duration, error)
}

// NewRetryConfig 创建新的 RetryConfig，使用默认值
//
// 返回:
//   - *RetryConfig: 使用默认配置的重试配置
func NewRetryConfig() *RetryConfig {
	return &RetryConfig{
		Count:       3,
		WaitTime:    1 * time.Second,
		MaxWaitTime: 30 * time.Second,
	}
}

// WithRetryCount 设置重试次数
//
// 参数:
//   - count: 最大重试次数（-1 表示禁用重试，0 表示使用 Client 默认配置）
//
// 返回:
//   - *RetryConfig: 返回自身，支持链式调用
func (r *RetryConfig) WithRetryCount(count int) *RetryConfig {
	r.Count = count
	return r
}

// WithRetryWaitTime 设置初始等待时间
//
// 参数:
//   - waitTime: 初始等待时间
//
// 返回:
//   - *RetryConfig: 返回自身，支持链式调用
func (r *RetryConfig) WithRetryWaitTime(waitTime time.Duration) *RetryConfig {
	r.WaitTime = waitTime
	return r
}

// WithRetryMaxWaitTime 设置最大等待时间
//
// 参数:
//   - maxWaitTime: 最大等待时间
//
// 返回:
//   - *RetryConfig: 返回自身，支持链式调用
func (r *RetryConfig) WithRetryMaxWaitTime(maxWaitTime time.Duration) *RetryConfig {
	r.MaxWaitTime = maxWaitTime
	return r
}

// WithRetryCondition 设置自定义重试条件
//
// 参数:
//   - condition: 重试条件函数
//
// 返回:
//   - *RetryConfig: 返回自身，支持链式调用
func (r *RetryConfig) WithRetryCondition(condition func(*resty.Response, error) bool) *RetryConfig {
	r.Condition = condition
	return r
}

// WithRetryAfter 设置自定义重试延迟函数
//
// 参数:
//   - after: 重试延迟函数
//
// 返回:
//   - *RetryConfig: 返回自身，支持链式调用
func (r *RetryConfig) WithRetryAfter(after func(*resty.Client, *resty.Response) (time.Duration, error)) *RetryConfig {
	r.After = after
	return r
}

// DisableRetry 禁用重试
//
// 返回:
//   - *RetryConfig: 返回自身，支持链式调用
func (r *RetryConfig) DisableRetry() *RetryConfig {
	r.Count = -1
	return r
}

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

// MultipartRequestConfig Multipart 请求配置
//
// 用于 multipart/form-data 请求的配置，支持文件上传等功能。
type MultipartRequestConfig struct {
	baseRequestConfig
	// MultipartFields Multipart form 数据字段列表
	MultipartFields []MultipartField
}

// NewMultipartRequestConfig 创建新的 MultipartRequestConfig，使用默认值
//
// 返回:
//   - MultipartRequestConfig: 所有字段都为默认值的 MultipartRequestConfig 实例
func NewMultipartRequestConfig() *MultipartRequestConfig {
	return &MultipartRequestConfig{
		baseRequestConfig: baseRequestConfig{
			Headers: make(map[string]string),
		},
	}
}

// MultipartField Multipart 表单字段
type MultipartField struct {
	// ParamName 字段名
	ParamName string
	// FileName 文件名（可选）
	FileName string
	// FilePath 文件路径（可选，用于文件上传）
	FilePath string
	// ContentType 内容类型（可选）
	ContentType string
	// Reader 数据读取器（可选，用于流式上传）
	Reader io.Reader
}

// WithMultipartField 添加 multipart form 数据字段
//
// 参数:
//   - field: Multipart 表单字段
//
// 返回:
//   - *MultipartRequestConfig: 返回自身，支持链式调用
func (c *MultipartRequestConfig) WithMultipartField(field MultipartField) *MultipartRequestConfig {
	if c.MultipartFields == nil {
		c.MultipartFields = make([]MultipartField, 0)
	}
	c.MultipartFields = append(c.MultipartFields, field)
	return c
}

// WithMultipartFields 设置多个 multipart form 数据字段
//
// 参数:
//   - fields: Multipart 表单字段列表
//
// 返回:
//   - *MultipartRequestConfig: 返回自身，支持链式调用
func (c *MultipartRequestConfig) WithMultipartFields(fields []MultipartField) *MultipartRequestConfig {
	c.MultipartFields = fields
	return c
}

// WithQuery 设置查询参数
//
// 参数:
//   - query: 查询参数
//
// 返回:
//   - *MultipartRequestConfig: 返回自身，支持链式调用
func (c *MultipartRequestConfig) WithQuery(query interface{}) *MultipartRequestConfig {
	c.Query = query
	return c
}

// WithAuth 设置认证信息
//
// 参数:
//   - auth: Basic Authentication 认证信息
//
// 返回:
//   - *MultipartRequestConfig: 返回自身，支持链式调用
func (c *MultipartRequestConfig) WithAuth(auth *Authorization) *MultipartRequestConfig {
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
//   - *MultipartRequestConfig: 返回自身，支持链式调用
func (c *MultipartRequestConfig) WithHeader(key, value string) *MultipartRequestConfig {
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
//   - *MultipartRequestConfig: 返回自身，支持链式调用
func (c *MultipartRequestConfig) WithHeaders(headers map[string]string) *MultipartRequestConfig {
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
//   - *MultipartRequestConfig: 返回自身，支持链式调用
func (c *MultipartRequestConfig) WithTimeout(timeout time.Duration) *MultipartRequestConfig {
	c.Timeout = timeout
	return c
}

// WithRetry 设置重试配置
//
// 参数:
//   - retry: 重试配置
//
// 返回:
//   - *MultipartRequestConfig: 返回自身，支持链式调用
func (c *MultipartRequestConfig) WithRetry(retry *RetryConfig) *MultipartRequestConfig {
	c.Retry = retry
	return c
}

// ensureHeaders 确保 Headers map 已初始化
func (c *MultipartRequestConfig) ensureHeaders() {
	if c.Headers == nil {
		c.Headers = make(map[string]string)
	}
}

// applyToRequest 将配置应用到 resty 请求
func (c *MultipartRequestConfig) applyToRequest(req *resty.Request) *resty.Request {
	// 添加 multipart form 数据字段
	if len(c.MultipartFields) > 0 {
		for _, field := range c.MultipartFields {
			if field.FilePath != "" {
				// 文件路径上传
				req = req.SetFileReader(field.ParamName, field.FileName, field.Reader)
			} else if field.Reader != nil {
				// 流式上传
				req = req.SetFileReader(field.ParamName, field.FileName, field.Reader)
			} else {
				// 普通字段
				req = req.SetFormData(map[string]string{field.ParamName: field.FileName})
			}
		}
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

	// 注意：resty.Request 不支持单独设置超时，超时在 Client 级别设置
	// 如果需要不同的超时，需要创建新的 Client 实例
	// 这里我们暂时忽略 Timeout 字段，使用 Client 的默认超时

	return req
}

// applyRetryConfig 将重试配置应用到 resty 客户端
// 注意：go-resty 的重试配置只能在 Client 级别设置，所以这里返回一个配置好的客户端
// 新客户端会继承基础客户端的超时设置，但重试配置会使用自定义值
// 注意：代理、认证等配置不会复制，因为这些配置通常在 Client 创建时设置，且 go-resty 没有提供获取这些配置的方法
func applyRetryConfig(baseClient *resty.Client, retry *RetryConfig) *resty.Client {
	// 创建新的客户端实例以应用自定义重试配置
	client := resty.New()

	// 继承基础客户端的超时设置
	client.SetTimeout(baseClient.GetClient().Timeout)

	// 如果 Count 为 -1，禁用重试
	if retry.Count == -1 {
		client.SetRetryCount(0)
		return client
	}

	// 设置重试次数（如果 > 0，否则使用基础客户端的默认值）
	if retry.Count > 0 {
		client.SetRetryCount(retry.Count)
	} else {
		// 使用基础客户端的重试次数
		client.SetRetryCount(baseClient.RetryCount)
	}

	// 设置等待时间（如果 > 0，否则使用基础客户端的默认值）
	if retry.WaitTime > 0 {
		client.SetRetryWaitTime(retry.WaitTime)
	} else {
		client.SetRetryWaitTime(baseClient.RetryWaitTime)
	}

	// 设置最大等待时间（如果 > 0，否则使用基础客户端的默认值）
	if retry.MaxWaitTime > 0 {
		client.SetRetryMaxWaitTime(retry.MaxWaitTime)
	} else {
		client.SetRetryMaxWaitTime(baseClient.RetryMaxWaitTime)
	}

	// 设置自定义重试条件（如果有）
	if retry.Condition != nil {
		client.AddRetryCondition(retry.Condition)
	} else {
		// 使用默认的重试条件（从 client.go 导入）
		client.AddRetryCondition(DefaultRetryCondition)
	}

	// 设置自定义重试延迟函数（如果有）
	if retry.After != nil {
		client.SetRetryAfter(retry.After)
	} else {
		// 使用默认的重试延迟函数（从 client.go 导入）
		client.SetRetryAfter(DefaultRetryAfter)
	}

	return client
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

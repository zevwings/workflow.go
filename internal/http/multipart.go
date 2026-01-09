package http

import (
	"io"
	"time"

	"github.com/go-resty/resty/v2"
)

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

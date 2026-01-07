package http

import (
	"fmt"
	"io"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/go-resty/resty/v2"
)

var (
	// globalClient 全局 HTTP 客户端单例
	globalClient *Client
	globalOnce   sync.Once
)

// Client HTTP 客户端封装
type Client struct {
	client *resty.Client
}

// DefaultRetryCondition 默认的重试条件函数（导出以供其他包使用）
func DefaultRetryCondition(r *resty.Response, err error) bool {
	// 网络错误或连接错误应该重试
	if err != nil {
		return isRetryableNetworkError(err)
	}

	// HTTP 状态码判断
	statusCode := r.StatusCode()
	// 5xx 服务器错误和 429 Too Many Requests 可重试
	if statusCode >= 500 && statusCode < 600 {
		return true
	}
	if statusCode == 429 {
		return true
	}

	// 4xx 客户端错误不可重试
	return false
}

// DefaultRetryAfter 默认的重试延迟函数（导出以供其他包使用）
func DefaultRetryAfter(client *resty.Client, resp *resty.Response) (time.Duration, error) {
	// 如果响应包含 Retry-After header，使用它
	if retryAfter := resp.Header().Get("Retry-After"); retryAfter != "" {
		if duration, err := parseRetryAfter(retryAfter); err == nil {
			return duration, nil
		}
	}

	// 否则使用指数退避：根据重试次数计算延迟
	// resty 会自动处理指数退避，这里返回 0 让 resty 使用默认的指数退避
	return 0, nil
}

// NewClient 创建新的 HTTP 客户端
func NewClient() *Client {
	client := resty.New()
	client.SetTimeout(30 * time.Second)

	// 配置默认重试策略
	client.SetRetryCount(3)
	client.SetRetryWaitTime(1 * time.Second)
	client.SetRetryMaxWaitTime(30 * time.Second)
	client.AddRetryCondition(DefaultRetryCondition)
	client.SetRetryAfter(DefaultRetryAfter)

	return &Client{client: client}
}

// isRetryableNetworkError 判断网络错误是否可重试
func isRetryableNetworkError(err error) bool {
	if err == nil {
		return false
	}

	errStr := err.Error()

	// 网络连接错误可重试
	retryableKeywords := []string{
		"timeout",
		"connection",
		"network",
		"dial",
		"connection refused",
		"connection reset",
		"no such host",
		"temporary failure",
	}

	for _, keyword := range retryableKeywords {
		if contains(errStr, keyword) {
			return true
		}
	}

	return false
}

// parseRetryAfter 解析 Retry-After header 值
func parseRetryAfter(value string) (time.Duration, error) {
	// Retry-After 可以是秒数（数字）或 HTTP 日期
	// 这里简化处理，只支持秒数
	if seconds, err := parseInt(value); err == nil {
		return time.Duration(seconds) * time.Second, nil
	}

	// 如果无法解析，返回错误
	return 0, fmt.Errorf("invalid Retry-After value: %s", value)
}

// contains 检查字符串是否包含子串（不区分大小写）
func contains(s, substr string) bool {
	return strings.Contains(strings.ToLower(s), strings.ToLower(substr))
}

// parseInt 尝试将字符串解析为整数
func parseInt(s string) (int, error) {
	return strconv.Atoi(strings.TrimSpace(s))
}

// Global 获取全局 HttpClient 单例
//
// 返回进程级别的 HttpClient 单例，使用默认配置。
// 单例会在首次调用时初始化，后续调用会复用同一个实例。
//
// 优势:
//   - 复用连接池：所有请求共享同一个连接池，提高性能
//   - 减少资源消耗：避免重复创建客户端实例
//   - 线程安全：可以在多线程环境中安全使用
//
// 返回:
//   - *Client: 全局 HTTP 客户端实例
func Global() *Client {
	globalOnce.Do(func() {
		globalClient = NewClient()
	})
	return globalClient
}

// SetAuth 设置认证 Token
func (c *Client) SetAuth(token string) {
	c.client.SetAuthToken(token)
}

// SetBasicAuth 设置 Basic Auth
func (c *Client) SetBasicAuth(username, password string) {
	c.client.SetBasicAuth(username, password)
}

// SetProxy 设置代理
func (c *Client) SetProxy(proxyURL string) {
	c.client.SetProxy(proxyURL)
}

// Get 发送 GET 请求（旧版 API，保持向后兼容）
func (c *Client) Get(url string) (*resty.Response, error) {
	return c.client.R().Get(url)
}

// Post 发送 POST 请求（旧版 API，保持向后兼容）
func (c *Client) Post(url string, body interface{}) (*resty.Response, error) {
	return c.client.R().SetBody(body).Post(url)
}

// Put 发送 PUT 请求（旧版 API，保持向后兼容）
func (c *Client) Put(url string, body interface{}) (*resty.Response, error) {
	return c.client.R().SetBody(body).Put(url)
}

// Delete 发送 DELETE 请求（旧版 API，保持向后兼容）
func (c *Client) Delete(url string) (*resty.Response, error) {
	return c.client.R().Delete(url)
}

// Patch 发送 PATCH 请求（旧版 API，保持向后兼容）
func (c *Client) Patch(url string, body interface{}) (*resty.Response, error) {
	return c.client.R().SetBody(body).Patch(url)
}

// doRequest 执行 HTTP 请求的通用方法
//
// 参数:
//   - method: HTTP 方法
//   - url: 请求 URL
//   - config: 请求配置（可选，如果为 nil 则使用默认配置）
//
// 返回:
//   - *HttpResponse: 封装后的 HTTP 响应
//   - error: 如果请求失败，返回错误
func (c *Client) doRequest(method HttpMethod, url string, config *RequestConfig) (*HttpResponse, error) {
	if config == nil {
		config = NewRequestConfig()
	}

	// 如果提供了自定义重试配置，创建临时客户端
	var client *resty.Client
	if config.Retry != nil {
		client = applyRetryConfig(c.client, config.Retry)
	} else {
		client = c.client
	}

	req := client.R()
	req = config.applyToRequest(req)

	var resp *resty.Response
	var err error

	switch method {
	case MethodGet:
		resp, err = req.Get(url)
	case MethodPost:
		resp, err = req.Post(url)
	case MethodPut:
		resp, err = req.Put(url)
	case MethodDelete:
		resp, err = req.Delete(url)
	case MethodPatch:
		resp, err = req.Patch(url)
	default:
		return nil, &InvalidMethodError{Method: string(method)}
	}

	if err != nil {
		return nil, err
	}

	return FromRestyResponse(resp)
}

// GetWithConfig 发送 GET 请求（新版 API，支持 RequestConfig）
//
// 参数:
//   - url: 请求 URL
//   - config: 请求配置（可选，如果为 nil 则使用默认配置）
//
// 返回:
//   - *HttpResponse: 封装后的 HTTP 响应
//   - error: 如果请求失败，返回错误
func (c *Client) GetWithConfig(url string, config *RequestConfig) (*HttpResponse, error) {
	return c.doRequest(MethodGet, url, config)
}

// PostWithConfig 发送 POST 请求（新版 API，支持 RequestConfig）
//
// 参数:
//   - url: 请求 URL
//   - config: 请求配置（可选，如果为 nil 则使用默认配置）
//
// 返回:
//   - *HttpResponse: 封装后的 HTTP 响应
//   - error: 如果请求失败，返回错误
func (c *Client) PostWithConfig(url string, config *RequestConfig) (*HttpResponse, error) {
	return c.doRequest(MethodPost, url, config)
}

// PutWithConfig 发送 PUT 请求（新版 API，支持 RequestConfig）
//
// 参数:
//   - url: 请求 URL
//   - config: 请求配置（可选，如果为 nil 则使用默认配置）
//
// 返回:
//   - *HttpResponse: 封装后的 HTTP 响应
//   - error: 如果请求失败，返回错误
func (c *Client) PutWithConfig(url string, config *RequestConfig) (*HttpResponse, error) {
	return c.doRequest(MethodPut, url, config)
}

// DeleteWithConfig 发送 DELETE 请求（新版 API，支持 RequestConfig）
//
// 参数:
//   - url: 请求 URL
//   - config: 请求配置（可选，如果为 nil 则使用默认配置）
//
// 返回:
//   - *HttpResponse: 封装后的 HTTP 响应
//   - error: 如果请求失败，返回错误
func (c *Client) DeleteWithConfig(url string, config *RequestConfig) (*HttpResponse, error) {
	return c.doRequest(MethodDelete, url, config)
}

// PatchWithConfig 发送 PATCH 请求（新版 API，支持 RequestConfig）
//
// 参数:
//   - url: 请求 URL
//   - config: 请求配置（可选，如果为 nil 则使用默认配置）
//
// 返回:
//   - *HttpResponse: 封装后的 HTTP 响应
//   - error: 如果请求失败，返回错误
func (c *Client) PatchWithConfig(url string, config *RequestConfig) (*HttpResponse, error) {
	return c.doRequest(MethodPatch, url, config)
}

// Stream 流式请求
//
// 发送请求并返回响应流，用于处理大文件或流式数据。
//
// 参数:
//   - method: HTTP 方法
//   - url: 请求 URL
//   - config: 请求配置（可选，如果为 nil 则使用默认配置）
//
// 返回:
//   - io.ReadCloser: 响应流
//   - error: 如果请求失败，返回错误
func (c *Client) Stream(method HttpMethod, url string, config *RequestConfig) (io.ReadCloser, error) {
	if config == nil {
		config = NewRequestConfig()
	}

	req := c.client.R()
	req = config.applyToRequest(req)

	var resp *resty.Response
	var err error

	switch method {
	case MethodGet:
		resp, err = req.Get(url)
	case MethodPost:
		resp, err = req.Post(url)
	case MethodPut:
		resp, err = req.Put(url)
	case MethodDelete:
		resp, err = req.Delete(url)
	case MethodPatch:
		resp, err = req.Patch(url)
	default:
		return nil, &InvalidMethodError{Method: string(method)}
	}

	if err != nil {
		return nil, err
	}

	return resp.RawBody(), nil
}

// PostMultipart POST Multipart 请求
//
// 发送 multipart/form-data 请求，通常用于文件上传。
//
// 参数:
//   - url: 请求 URL
//   - config: Multipart 请求配置
//
// 返回:
//   - *HttpResponse: 封装后的 HTTP 响应
//   - error: 如果请求失败，返回错误
func (c *Client) PostMultipart(url string, config *MultipartRequestConfig) (*HttpResponse, error) {
	if config == nil {
		return nil, &ConfigError{Message: "MultipartRequestConfig is required for multipart requests"}
	}

	// 如果提供了自定义重试配置，创建临时客户端
	var client *resty.Client
	if config.Retry != nil {
		client = applyRetryConfig(c.client, config.Retry)
	} else {
		client = c.client
	}

	req := client.R()
	req = config.applyToRequest(req)

	resp, err := req.Post(url)
	if err != nil {
		return nil, err
	}

	return FromRestyResponse(resp)
}

// GetRestyClient 获取底层 resty 客户端（用于高级用法）
func (c *Client) GetRestyClient() *resty.Client {
	return c.client
}

// InvalidMethodError 无效的 HTTP 方法错误
type InvalidMethodError struct {
	Method string
}

func (e *InvalidMethodError) Error() string {
	return "invalid HTTP method: " + e.Method
}

// ConfigError 配置错误
type ConfigError struct {
	Message string
}

func (e *ConfigError) Error() string {
	return "config error: " + e.Message
}

package http

import (
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

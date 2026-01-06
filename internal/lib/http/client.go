package http

import (
	"time"

	"github.com/go-resty/resty/v2"
)

// Client HTTP 客户端封装
type Client struct {
	client *resty.Client
}

// NewClient 创建新的 HTTP 客户端
func NewClient() *Client {
	client := resty.New()
	client.SetTimeout(30 * time.Second)
	client.SetRetryCount(3)
	client.SetRetryWaitTime(1 * time.Second)
	client.SetRetryMaxWaitTime(10 * time.Second)
	client.SetRetryAfter(func(client *resty.Client, resp *resty.Response) (time.Duration, error) {
		return time.Duration(resp.StatusCode()/100) * time.Second, nil
	})

	return &Client{client: client}
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

// Get 发送 GET 请求
func (c *Client) Get(url string) (*resty.Response, error) {
	return c.client.R().Get(url)
}

// Post 发送 POST 请求
func (c *Client) Post(url string, body interface{}) (*resty.Response, error) {
	return c.client.R().SetBody(body).Post(url)
}

// Put 发送 PUT 请求
func (c *Client) Put(url string, body interface{}) (*resty.Response, error) {
	return c.client.R().SetBody(body).Put(url)
}

// Delete 发送 DELETE 请求
func (c *Client) Delete(url string) (*resty.Response, error) {
	return c.client.R().Delete(url)
}

// Patch 发送 PATCH 请求
func (c *Client) Patch(url string, body interface{}) (*resty.Response, error) {
	return c.client.R().SetBody(body).Patch(url)
}

// GetRestyClient 获取底层 resty 客户端（用于高级用法）
func (c *Client) GetRestyClient() *resty.Client {
	return c.client
}


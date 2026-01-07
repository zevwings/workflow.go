package jira

import (
	"context"
	"fmt"

	"github.com/andygrunwald/go-jira/v2/cloud"
)

// Config Jira 配置
//
// 包含 Jira 客户端所需的所有配置信息。
type Config struct {
	// URL Jira 服务器 URL（如 "https://your-domain.atlassian.net"）
	URL string
	// Username 用户名（Email 地址）
	Username string
	// Token API Token
	Token string
}

// Client Jira 客户端封装
//
// 封装 go-jira SDK，提供统一的 Jira API 访问接口。
// 使用 Basic Auth（Email + API Token）进行认证。
type Client struct {
	jira *cloud.Client
	ctx  context.Context
}

// NewClient 创建新的 Jira 客户端
//
// 使用传入的配置创建并初始化 Jira 客户端。
//
// 参数:
//   - config: Jira 配置结构体（不能为 nil）
//
// 返回:
//   - *Client: Jira 客户端实例
//   - error: 如果配置不完整或客户端创建失败，返回错误
func NewClient(config *Config) (*Client, error) {
	if config == nil {
		return nil, fmt.Errorf("config 不能为 nil")
	}

	if config.URL == "" || config.Username == "" || config.Token == "" {
		return nil, fmt.Errorf("Jira 配置不完整，请设置 URL、Username 和 Token")
	}

	// 创建 Jira 客户端（使用 Basic Auth）
	tp := cloud.BasicAuthTransport{
		Username: config.Username,
		APIToken: config.Token,
	}

	jiraClient, err := cloud.NewClient(config.URL, tp.Client())
	if err != nil {
		return nil, fmt.Errorf("创建 Jira 客户端失败: %w", err)
	}

	return &Client{
		jira: jiraClient,
		ctx:  context.Background(),
	}, nil
}

// WithContext 使用指定的 context 创建客户端副本
//
// 参数:
//   - ctx: 上下文对象
//
// 返回:
//   - *Client: 新的客户端实例（使用指定的 context）
func (c *Client) WithContext(ctx context.Context) *Client {
	return &Client{
		jira: c.jira,
		ctx:  ctx,
	}
}

// GetJiraClient 获取底层 go-jira 客户端（用于高级用法）
//
// 返回:
//   - *cloud.Client: 底层 Jira 客户端
func (c *Client) GetJiraClient() *cloud.Client {
	return c.jira
}

// GetContext 获取当前使用的 context
//
// 返回:
//   - context.Context: 当前上下文
func (c *Client) GetContext() context.Context {
	return c.ctx
}

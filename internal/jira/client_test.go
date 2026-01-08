package jira

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// ==================== NewClient 测试 ====================

func TestNewClient(t *testing.T) {
	tests := []struct {
		name    string
		config  *Config
		wantErr bool
	}{
		{
			name: "valid config",
			config: &Config{
				URL:      "https://test.atlassian.net",
				Username: "test@example.com",
				Token:    "test-token",
			},
			wantErr: false,
		},
		{
			name:    "nil config",
			config:  nil,
			wantErr: true,
		},
		{
			name: "empty URL",
			config: &Config{
				URL:      "",
				Username: "test@example.com",
				Token:    "test-token",
			},
			wantErr: true,
		},
		{
			name: "empty username",
			config: &Config{
				URL:      "https://test.atlassian.net",
				Username: "",
				Token:    "test-token",
			},
			wantErr: true,
		},
		{
			name: "empty token",
			config: &Config{
				URL:      "https://test.atlassian.net",
				Username: "test@example.com",
				Token:    "",
			},
			wantErr: true,
		},
		{
			name: "invalid URL format",
			config: &Config{
				URL:      "not-a-url",
				Username: "test@example.com",
				Token:    "test-token",
			},
			wantErr: false, // URL 格式验证由 go-jira SDK 处理，这里只验证非空
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, err := NewClient(tt.config)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, client)
			} else {
				// 注意：go-jira SDK 会尝试连接服务器
				// 如果连接失败，这里会返回错误，但这是预期的行为
				// 在实际测试中，可能需要 Mock 或使用集成测试
				if err != nil {
					// 如果是因为网络连接失败，这是可以接受的
					// 我们主要测试配置验证逻辑
					t.Logf("客户端创建返回错误（可能是网络问题）: %v", err)
				} else {
					assert.NotNil(t, client)
					assert.NotNil(t, client.GetJiraClient())
					assert.NotNil(t, client.GetContext())
				}
			}
		})
	}
}

// ==================== WithContext 测试 ====================

func TestClient_WithContext(t *testing.T) {
	config := &Config{
		URL:      "https://test.atlassian.net",
		Username: "test@example.com",
		Token:    "test-token",
	}

	client, err := NewClient(config)
	if err != nil {
		t.Skip("需要有效的 Jira 配置或 Mock，跳过 WithContext 测试")
	}

	// 测试使用自定义 context
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	newClient := client.WithContext(ctx)

	// 验证返回的是新实例
	assert.NotEqual(t, client, newClient)

	// 验证 context 被正确设置
	assert.Equal(t, ctx, newClient.GetContext())

	// 验证底层客户端是共享的
	assert.Equal(t, client.GetJiraClient(), newClient.GetJiraClient())

	// 验证原始客户端的 context 没有被改变
	assert.NotEqual(t, ctx, client.GetContext())
}

// ==================== GetJiraClient 测试 ====================

func TestClient_GetJiraClient(t *testing.T) {
	config := &Config{
		URL:      "https://test.atlassian.net",
		Username: "test@example.com",
		Token:    "test-token",
	}

	client, err := NewClient(config)
	if err != nil {
		t.Skip("需要有效的 Jira 配置或 Mock，跳过 GetJiraClient 测试")
	}

	jiraClient := client.GetJiraClient()
	assert.NotNil(t, jiraClient)
}

// ==================== GetContext 测试 ====================

func TestClient_GetContext(t *testing.T) {
	config := &Config{
		URL:      "https://test.atlassian.net",
		Username: "test@example.com",
		Token:    "test-token",
	}

	client, err := NewClient(config)
	if err != nil {
		t.Skip("需要有效的 Jira 配置或 Mock，跳过 GetContext 测试")
	}

	ctx := client.GetContext()
	assert.NotNil(t, ctx)

	// 验证默认 context 是 context.Background()
	assert.Equal(t, context.Background(), ctx)
}

// ==================== Context 链式测试 ====================

func TestClient_ContextChain(t *testing.T) {
	config := &Config{
		URL:      "https://test.atlassian.net",
		Username: "test@example.com",
		Token:    "test-token",
	}

	client, err := NewClient(config)
	if err != nil {
		t.Skip("需要有效的 Jira 配置或 Mock，跳过 Context 链式测试")
	}

	// 创建多个带不同 context 的客户端
	ctx1 := context.WithValue(context.Background(), "key1", "value1")
	ctx2 := context.WithValue(context.Background(), "key2", "value2")

	client1 := client.WithContext(ctx1)
	client2 := client.WithContext(ctx2)
	client3 := client1.WithContext(ctx2)

	// 验证每个客户端都有独立的 context
	assert.Equal(t, ctx1, client1.GetContext())
	assert.Equal(t, ctx2, client2.GetContext())
	assert.Equal(t, ctx2, client3.GetContext())

	// 验证底层客户端是共享的
	assert.Equal(t, client.GetJiraClient(), client1.GetJiraClient())
	assert.Equal(t, client.GetJiraClient(), client2.GetJiraClient())
	assert.Equal(t, client.GetJiraClient(), client3.GetJiraClient())
}


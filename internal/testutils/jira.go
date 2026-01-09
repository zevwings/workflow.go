//go:build test

package testutils

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/zevwings/workflow/internal/jira"
	"github.com/zevwings/workflow/internal/jira/api"
)

// NewJiraTestClient 创建用于测试的 Jira 客户端
// 这是对常见测试模式的便捷包装，简化客户端创建流程
//
// 参数:
//   - t: 测试对象
//   - handler: 可选的 HTTP 请求处理函数，如果为 nil 则使用默认的 mock handler
//
// 返回:
//   - *jira.JiraClient: Jira 客户端实例
//   - *httptest.Server: Mock 服务器实例（用于访问 URL 或关闭）
//
// 注意: 服务器会在测试结束时自动关闭（使用 t.Cleanup）
//
// 示例:
//
//	client, server := testutils.NewJiraTestClient(t, nil)
//	// 使用 client 进行测试
//	// server.URL 可以用于需要访问服务器 URL 的场景
func NewJiraTestClient(t *testing.T, handler http.HandlerFunc) (*jira.JiraClient, *httptest.Server) {
	t.Helper()

	// 创建 Mock 服务器
	server := api.SetupMockJiraServer(t, handler)
	// 确保服务器在测试结束时关闭
	t.Cleanup(func() {
		server.Close()
	})

	// 创建配置
	config := &jira.Config{
		ServiceAddress: server.URL,
		Email:          "test@example.com",
		APIToken:       "test-token",
	}

	// 创建客户端
	client, err := jira.NewJiraClient(config)
	if err != nil {
		t.Fatalf("创建 Jira 测试客户端失败: %v", err)
	}

	return client, server
}

// NewJiraTestClientWithCreds 创建用于测试的 Jira 客户端（自定义认证信息）
//
// 参数:
//   - t: 测试对象
//   - username: 用户名（Email）
//   - token: API Token
//   - handler: 可选的 HTTP 请求处理函数，如果为 nil 则使用默认的 mock handler
//
// 返回:
//   - *jira.JiraClient: Jira 客户端实例
//   - *httptest.Server: Mock 服务器实例
//
// 示例:
//
//	client, server := testutils.NewJiraTestClientWithCreds(t, "custom@example.com", "custom-token", nil)
func NewJiraTestClientWithCreds(t *testing.T, username, token string, handler http.HandlerFunc) (*jira.JiraClient, *httptest.Server) {
	t.Helper()

	// 创建 Mock 服务器
	server := api.SetupMockJiraServer(t, handler)
	// 确保服务器在测试结束时关闭
	t.Cleanup(func() {
		server.Close()
	})

	// 创建配置
	config := &jira.Config{
		ServiceAddress: server.URL,
		Email:          username,
		APIToken:       token,
	}

	// 创建客户端
	client, err := jira.NewJiraClient(config)
	if err != nil {
		t.Fatalf("创建 Jira 测试客户端失败: %v", err)
	}

	return client, server
}


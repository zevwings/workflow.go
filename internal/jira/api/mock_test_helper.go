package api

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/andygrunwald/go-jira/v2/cloud"
)

// SetupMockJiraServer 创建 Mock Jira API 服务器（导出以供其他包使用）
//
// 参数:
//   - t: 测试对象
//   - handler: 自定义请求处理函数（可选，如果为 nil 则使用默认处理）
//
// 返回:
//   - *httptest.Server: Mock 服务器实例
func SetupMockJiraServer(t *testing.T, handler http.HandlerFunc) *httptest.Server {
	t.Helper()

	if handler == nil {
		handler = DefaultMockHandler
	}

	server := httptest.NewServer(handler)
	return server
}

// DefaultMockHandler 默认的 Mock 处理函数（导出以供其他包使用）
func DefaultMockHandler(w http.ResponseWriter, r *http.Request) {
	// 根据路径和方法返回不同的响应
	path := r.URL.Path

	// 使用更灵活的路径匹配，处理可能的路径变体
	normalizePath := func(p string) string {
		// 移除尾部斜杠
		if len(p) > 1 && p[len(p)-1] == '/' {
			return p[:len(p)-1]
		}
		return p
	}
	normalizedPath := normalizePath(path)

	// 支持 API v2 和 v3
	isIssuePath := normalizedPath == "/rest/api/2/issue/PROJ-123" || normalizedPath == "/rest/api/3/issue/PROJ-123"
	isTransitionsPath := normalizedPath == "/rest/api/2/issue/PROJ-123/transitions" || normalizedPath == "/rest/api/3/issue/PROJ-123/transitions"
	isAssigneePath := normalizedPath == "/rest/api/2/issue/PROJ-123/assignee" || normalizedPath == "/rest/api/3/issue/PROJ-123/assignee"
	isCommentPath := normalizedPath == "/rest/api/2/issue/PROJ-123/comment" || normalizedPath == "/rest/api/3/issue/PROJ-123/comment"
	isProjectPath := normalizedPath == "/rest/api/2/project/PROJ" || normalizedPath == "/rest/api/3/project/PROJ" ||
		normalizedPath == "/rest/api/2/project/proj" || normalizedPath == "/rest/api/3/project/proj"
	isMyselfPath := normalizedPath == "/rest/api/2/myself" || normalizedPath == "/rest/api/3/myself"
	isUserSearchPath := normalizedPath == "/rest/api/2/user/search" || normalizedPath == "/rest/api/3/user/search"

	switch {
	case isIssuePath && r.Method == http.MethodGet:
		// 获取 Issue
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"id":  "10000",
			"key": "PROJ-123",
			"fields": map[string]interface{}{
				"summary": "Test Issue",
				"status": map[string]interface{}{
					"id":   "1",
					"name": "To Do",
				},
			},
		})

	case isTransitionsPath && r.Method == http.MethodGet:
		// 获取 Transitions
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"transitions": []map[string]interface{}{
				{"id": "11", "name": "In Progress"},
				{"id": "21", "name": "Done"},
			},
		})

	case isTransitionsPath && r.Method == http.MethodPost:
		// 执行 Transition
		w.WriteHeader(http.StatusNoContent)

	case isAssigneePath && r.Method == http.MethodPut:
		// 分配 Issue
		w.WriteHeader(http.StatusNoContent)

	case isCommentPath && r.Method == http.MethodPost:
		// 添加评论
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"id":   "10000",
			"body": "Test comment",
		})

	case isProjectPath && r.Method == http.MethodGet:
		// 获取项目（支持大小写）
		projectKey := "PROJ"
		if normalizedPath == "/rest/api/2/project/proj" || normalizedPath == "/rest/api/3/project/proj" {
			projectKey = "proj"
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"id":   "10000",
			"key":  projectKey,
			"name": "Test Project",
		})

	case isMyselfPath && r.Method == http.MethodGet:
		// 获取当前用户
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"accountId":    "12345",
			"displayName":  "Test User",
			"emailAddress": "test@example.com",
		})

	case isUserSearchPath && r.Method == http.MethodGet:
		// 搜索用户
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode([]map[string]interface{}{
			{
				"accountId":    "123",
				"displayName":  "John Doe",
				"emailAddress": "john@example.com",
			},
		})

	// GetUser API: /rest/api/2/user?accountId=xxx 或 /rest/api/3/user?accountId=xxx
	case (normalizedPath == "/rest/api/2/user" || normalizedPath == "/rest/api/3/user") && r.Method == http.MethodGet:
		// 根据 Account ID 获取用户
		accountID := r.URL.Query().Get("accountId")
		if accountID == "" {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"errorMessages": []string{"Account ID is required"},
			})
			return
		}
		// 对 "invalid" account ID 返回错误
		if accountID == "invalid" {
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"errorMessages": []string{"User not found"},
			})
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"accountId":    accountID,
			"displayName":  "Test User",
			"emailAddress": "test@example.com",
		})

	default:
		// 调试：打印未匹配的请求
		// 在实际测试中，可以取消注释来调试
		// fmt.Printf("Mock Server: Unmatched request - %s %s\n", r.Method, path)
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`{"errorMessages": ["Not Found"]}`))
	}
}

// CreateTestClient 创建用于测试的 Jira 客户端（导出以供其他包使用）
//
// 参数:
//   - t: 测试对象
//   - serverURL: Mock 服务器 URL
//
// 返回:
//   - *cloud.Client: go-jira 客户端实例
//   - context.Context: 上下文
func CreateTestClient(t *testing.T, serverURL string) (*cloud.Client, context.Context) {
	t.Helper()

	// 解析 Mock 服务器 URL
	mockURL, err := url.Parse(serverURL)
	if err != nil {
		t.Fatalf("解析 Mock 服务器 URL 失败: %v", err)
	}

	// 创建 Basic Auth Transport（使用默认 Transport）
	tp := cloud.BasicAuthTransport{
		Username: "test@example.com",
		APIToken: "test-token",
	}

	// 创建自定义 Transport，包装 Basic Auth Transport
	// 这样可以在 Basic Auth 处理之后修改 URL 指向 Mock 服务器
	mockTransport := &mockTransport{
		mockURL:       mockURL,
		baseTransport: &tp,
	}

	// 创建 HTTP 客户端，使用自定义 Transport
	httpClient := &http.Client{
		Transport: mockTransport,
	}

	// 创建 go-jira 客户端，使用 Mock 服务器 URL 和自定义 HTTP 客户端
	client, err := cloud.NewClient(serverURL, httpClient)
	if err != nil {
		t.Fatalf("创建测试客户端失败: %v", err)
	}

	return client, context.Background()
}

// mockTransport 自定义 Transport，将请求转发到 Mock 服务器
// 它包装了 BasicAuthTransport，确保在 Basic Auth 处理之后修改 URL
type mockTransport struct {
	mockURL       *url.URL
	baseTransport http.RoundTripper
}

func (m *mockTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	// 创建新请求，避免修改原始请求
	newReq := req.Clone(req.Context())

	// 修改请求 URL，指向 Mock 服务器，但保留原始路径和查询参数
	newReq.URL.Scheme = m.mockURL.Scheme
	newReq.URL.Host = m.mockURL.Host
	// 路径和查询参数保持不变

	return m.baseTransport.RoundTrip(newReq)
}

// createTestIssueAPI 创建用于测试的 IssueAPI 实例
func createTestIssueAPI(t *testing.T, server *httptest.Server) *IssueAPI {
	t.Helper()

	client, ctx := CreateTestClient(t, server.URL)
	return NewIssueAPI(client, ctx)
}

// createTestProjectAPI 创建用于测试的 ProjectAPI 实例
func createTestProjectAPI(t *testing.T, server *httptest.Server) *ProjectAPI {
	t.Helper()

	client, ctx := CreateTestClient(t, server.URL)
	return NewProjectAPI(client, ctx)
}

// createTestUserAPI 创建用于测试的 UserAPI 实例
func createTestUserAPI(t *testing.T, server *httptest.Server) *UserAPI {
	t.Helper()

	client, ctx := CreateTestClient(t, server.URL)
	return NewUserAPI(client, ctx)
}

package jira

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zevwings/workflow/internal/jira/api"
)

// ==================== NewJiraClient 测试 ====================

func TestNewJiraClient(t *testing.T) {
	// 使用 Mock 服务器测试客户端创建
	server := api.SetupMockJiraServer(t, nil)
	defer server.Close()

	config := &Config{
		URL:      server.URL,
		Username: "test@example.com",
		Token:    "test-token",
	}

	client, err := NewJiraClient(config)
	require.NoError(t, err)
	assert.NotNil(t, client)
	assert.NotNil(t, client.GetClient())
	assert.NotNil(t, client.GetIssueAPI())
	assert.NotNil(t, client.GetProjectAPI())
	assert.NotNil(t, client.GetUserAPI())
}

func TestNewJiraClient_InvalidConfig(t *testing.T) {
	tests := []struct {
		name   string
		config *Config
	}{
		{
			name:   "nil config",
			config: nil,
		},
		{
			name: "empty URL",
			config: &Config{
				URL:      "",
				Username: "test@example.com",
				Token:    "test-token",
			},
		},
		{
			name: "empty username",
			config: &Config{
				URL:      "https://test.atlassian.net",
				Username: "",
				Token:    "test-token",
			},
		},
		{
			name: "empty token",
			config: &Config{
				URL:      "https://test.atlassian.net",
				Username: "test@example.com",
				Token:    "",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, err := NewJiraClient(tt.config)
			assert.Error(t, err)
			assert.Nil(t, client)
		})
	}
}

// ==================== GetUserInfo 测试 ====================

func TestJiraClient_GetUserInfo(t *testing.T) {
	server := api.SetupMockJiraServer(t, nil)
	defer server.Close()

	config := &Config{
		URL:      server.URL,
		Username: "test@example.com",
		Token:    "test-token",
	}

	client, err := NewJiraClient(config)
	require.NoError(t, err)

	user, err := client.GetUserInfo()
	require.NoError(t, err)
	assert.NotNil(t, user)
	assert.NotEmpty(t, user.AccountID)
	assert.NotEmpty(t, user.DisplayName)
}

// ==================== GetTicketInfo 测试 ====================

func TestJiraClient_GetTicketInfo(t *testing.T) {
	server := api.SetupMockJiraServer(t, nil)
	defer server.Close()

	config := &Config{
		URL:      server.URL,
		Username: "test@example.com",
		Token:    "test-token",
	}

	client, err := NewJiraClient(config)
	require.NoError(t, err)

	tests := []struct {
		name    string
		ticket  string
		wantErr bool
	}{
		{
			name:    "valid ticket",
			ticket:  "PROJ-123",
			wantErr: false,
		},
		{
			name:    "lowercase ticket",
			ticket:  "proj-123",
			wantErr: false, // 应该被规范化
		},
		{
			name:    "invalid format",
			ticket:  "invalid",
			wantErr: true,
		},
		{
			name:    "empty ticket",
			ticket:  "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			issue, err := client.GetTicketInfo(tt.ticket)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, issue)
			} else {
				require.NoError(t, err)
				assert.NotNil(t, issue)
				assert.Equal(t, "PROJ-123", issue.Key)
			}
		})
	}
}

// ==================== GetAttachments 测试 ====================

func TestJiraClient_GetAttachments(t *testing.T) {
	server := api.SetupMockJiraServer(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 返回带附件的 Issue
		if r.URL.Path == "/rest/api/2/issue/PROJ-123" && r.Method == http.MethodGet {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"id":  "10000",
				"key": "PROJ-123",
				"fields": map[string]interface{}{
					"summary": "Test Issue",
					"attachment": []map[string]interface{}{
						{"id": "10000", "filename": "test.txt"},
					},
				},
			})
			return
		}
		api.DefaultMockHandler(w, r)
	}))
	defer server.Close()

	config := &Config{
		URL:      server.URL,
		Username: "test@example.com",
		Token:    "test-token",
	}

	client, err := NewJiraClient(config)
	require.NoError(t, err)

	attachments, err := client.GetAttachments("PROJ-123")
	require.NoError(t, err)
	assert.NotEmpty(t, attachments)
	assert.Equal(t, "test.txt", attachments[0].Filename)
}

// ==================== MoveTicket 测试 ====================

func TestJiraClient_MoveTicket(t *testing.T) {
	server := api.SetupMockJiraServer(t, nil)
	defer server.Close()

	config := &Config{
		URL:      server.URL,
		Username: "test@example.com",
		Token:    "test-token",
	}

	client, err := NewJiraClient(config)
	require.NoError(t, err)

	tests := []struct {
		name    string
		ticket  string
		status  string
		wantErr bool
	}{
		{
			name:    "valid transition",
			ticket:  "PROJ-123",
			status:  "In Progress",
			wantErr: false,
		},
		{
			name:    "status not found",
			ticket:  "PROJ-123",
			status:  "Invalid Status",
			wantErr: true,
		},
		{
			name:    "invalid ticket format",
			ticket:  "invalid",
			status:  "In Progress",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := client.MoveTicket(tt.ticket, tt.status)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

// ==================== AssignTicket 测试 ====================

func TestJiraClient_AssignTicket(t *testing.T) {
	server := api.SetupMockJiraServer(t, nil)
	defer server.Close()

	config := &Config{
		URL:      server.URL,
		Username: "test@example.com",
		Token:    "test-token",
	}

	client, err := NewJiraClient(config)
	require.NoError(t, err)

	tests := []struct {
		name      string
		ticket    string
		accountID *string
		wantErr   bool
	}{
		{
			name:      "assign to user",
			ticket:    "PROJ-123",
			accountID: stringPtr("12345"),
			wantErr:   false,
		},
		{
			name:      "unassign (nil account ID)",
			ticket:    "PROJ-123",
			accountID: nil,
			wantErr:   false,
		},
		{
			name:      "unassign (empty account ID)",
			ticket:    "PROJ-123",
			accountID: stringPtr(""),
			wantErr:   false,
		},
		{
			name:    "invalid ticket format",
			ticket:  "invalid",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := client.AssignTicket(tt.ticket, tt.accountID)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

// ==================== AddComment 测试 ====================

func TestJiraClient_AddComment(t *testing.T) {
	server := api.SetupMockJiraServer(t, nil)
	defer server.Close()

	config := &Config{
		URL:      server.URL,
		Username: "test@example.com",
		Token:    "test-token",
	}

	client, err := NewJiraClient(config)
	require.NoError(t, err)

	tests := []struct {
		name    string
		ticket  string
		comment string
		wantErr bool
	}{
		{
			name:    "valid comment",
			ticket:  "PROJ-123",
			comment: "Test comment",
			wantErr: false,
		},
		{
			name:    "invalid ticket format",
			ticket:  "invalid",
			comment: "Test comment",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := client.AddComment(tt.ticket, tt.comment)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

// ==================== GetComments 测试 ====================

func TestJiraClient_GetComments(t *testing.T) {
	server := api.SetupMockJiraServer(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/rest/api/2/issue/PROJ-123" && r.Method == http.MethodGet {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"id":  "10000",
				"key": "PROJ-123",
				"fields": map[string]interface{}{
					"summary": "Test Issue",
					"comment": map[string]interface{}{
						"comments": []map[string]interface{}{
							{"id": "10000", "body": "First comment"},
							{"id": "10001", "body": "Second comment"},
						},
					},
				},
			})
			return
		}
		api.DefaultMockHandler(w, r)
	}))
	defer server.Close()

	config := &Config{
		URL:      server.URL,
		Username: "test@example.com",
		Token:    "test-token",
	}

	client, err := NewJiraClient(config)
	require.NoError(t, err)

	comments, err := client.GetComments("PROJ-123")
	require.NoError(t, err)
	assert.NotEmpty(t, comments)
	assert.Equal(t, 2, len(comments))
}

// ==================== GetTransitions 测试 ====================

func TestJiraClient_GetTransitions(t *testing.T) {
	server := api.SetupMockJiraServer(t, nil)
	defer server.Close()

	config := &Config{
		URL:      server.URL,
		Username: "test@example.com",
		Token:    "test-token",
	}

	client, err := NewJiraClient(config)
	require.NoError(t, err)

	transitions, err := client.GetTransitions("PROJ-123")
	require.NoError(t, err)
	assert.NotEmpty(t, transitions)
	assert.Equal(t, 2, len(transitions))
}

// ==================== GetChangelog 测试 ====================

func TestJiraClient_GetChangelog(t *testing.T) {
	server := api.SetupMockJiraServer(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 返回带 changelog 的 Issue
		if r.URL.Path == "/rest/api/2/issue/PROJ-123" && r.Method == http.MethodGet {
			// 检查是否有 expand=changelog 参数
			expand := r.URL.Query().Get("expand")
			if expand == "changelog" {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				json.NewEncoder(w).Encode(map[string]interface{}{
					"id":  "10000",
					"key": "PROJ-123",
					"changelog": map[string]interface{}{
						"histories": []map[string]interface{}{},
					},
				})
				return
			}
		}
		api.DefaultMockHandler(w, r)
	}))
	defer server.Close()

	config := &Config{
		URL:      server.URL,
		Username: "test@example.com",
		Token:    "test-token",
	}

	client, err := NewJiraClient(config)
	require.NoError(t, err)

	changelog, err := client.GetChangelog("PROJ-123")
	require.NoError(t, err)
	assert.NotNil(t, changelog)
}

// ==================== GetProject 测试 ====================

func TestJiraClient_GetProject(t *testing.T) {
	server := api.SetupMockJiraServer(t, nil)
	defer server.Close()

	config := &Config{
		URL:      server.URL,
		Username: "test@example.com",
		Token:    "test-token",
	}

	client, err := NewJiraClient(config)
	require.NoError(t, err)

	project, err := client.GetProject("PROJ")
	require.NoError(t, err)
	assert.NotNil(t, project)
	assert.Equal(t, "PROJ", project.Key)
}

// ==================== GetProjectStatuses 测试 ====================

func TestJiraClient_GetProjectStatuses(t *testing.T) {
	server := api.SetupMockJiraServer(t, nil)
	defer server.Close()

	config := &Config{
		URL:      server.URL,
		Username: "test@example.com",
		Token:    "test-token",
	}

	client, err := NewJiraClient(config)
	require.NoError(t, err)

	statuses, err := client.GetProjectStatuses("PROJ")
	require.NoError(t, err)
	// 当前实现返回空列表
	assert.Empty(t, statuses)
}

// ==================== FindUsers 测试 ====================

func TestJiraClient_FindUsers(t *testing.T) {
	server := api.SetupMockJiraServer(t, nil)
	defer server.Close()

	config := &Config{
		URL:      server.URL,
		Username: "test@example.com",
		Token:    "test-token",
	}

	client, err := NewJiraClient(config)
	require.NoError(t, err)

	users, err := client.FindUsers("john")
	require.NoError(t, err)
	assert.NotEmpty(t, users)
	assert.Equal(t, "John Doe", users[0].DisplayName)
}

// ==================== Getter 方法测试 ====================

func TestJiraClient_Getters(t *testing.T) {
	server := api.SetupMockJiraServer(t, nil)
	defer server.Close()

	config := &Config{
		URL:      server.URL,
		Username: "test@example.com",
		Token:    "test-token",
	}

	client, err := NewJiraClient(config)
	require.NoError(t, err)

	// 测试 Getter 方法
	assert.NotNil(t, client.GetClient())
	assert.NotNil(t, client.GetIssueAPI())
	assert.NotNil(t, client.GetProjectAPI())
	assert.NotNil(t, client.GetUserAPI())

	// 验证返回的是正确的类型（已经是具体类型，直接验证非 nil）
	assert.IsType(t, &api.IssueAPI{}, client.GetIssueAPI())
	assert.IsType(t, &api.ProjectAPI{}, client.GetProjectAPI())
	assert.IsType(t, &api.UserAPI{}, client.GetUserAPI())
}

// stringPtr 辅助函数，用于创建字符串指针
func stringPtr(s string) *string {
	return &s
}

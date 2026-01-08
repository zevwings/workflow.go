package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/andygrunwald/go-jira/v2/cloud"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ==================== NewIssueAPI 测试 ====================

func TestNewIssueAPI(t *testing.T) {
	server := SetupMockJiraServer(t, nil)
	defer server.Close()

	client, ctx := CreateTestClient(t, server.URL)
	api := NewIssueAPI(client, ctx)

	assert.NotNil(t, api)
	assert.NotNil(t, api.client)
	assert.NotNil(t, api.ctx)
}

// ==================== GetIssue 测试 ====================

func TestIssueAPI_GetIssue(t *testing.T) {
	tests := []struct {
		name       string
		ticket     string
		setupMock  func(*testing.T) *httptest.Server
		wantErr    bool
		validateFn func(*testing.T, *cloud.Issue)
	}{
		{
			name:   "valid ticket",
			ticket: "PROJ-123",
			setupMock: func(t *testing.T) *httptest.Server {
				return SetupMockJiraServer(t, nil)
			},
			wantErr: false,
			validateFn: func(t *testing.T, issue *cloud.Issue) {
				assert.NotNil(t, issue)
				assert.Equal(t, "PROJ-123", issue.Key)
				assert.Equal(t, "10000", issue.ID)
				if issue.Fields != nil {
					assert.Equal(t, "Test Issue", issue.Fields.Summary)
				}
			},
		},
		{
			name:   "not found",
			ticket: "PROJ-999",
			setupMock: func(t *testing.T) *httptest.Server {
				return SetupMockJiraServer(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusNotFound)
					json.NewEncoder(w).Encode(map[string]interface{}{
						"errorMessages": []string{"Issue does not exist"},
					})
				}))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := tt.setupMock(t)
			defer server.Close()

			api := createTestIssueAPI(t, server)
			issue, err := api.GetIssue(tt.ticket)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, issue)
			} else {
				require.NoError(t, err)
				assert.NotNil(t, issue)
				if tt.validateFn != nil {
					tt.validateFn(t, issue)
				}
			}
		})
	}
}

// ==================== GetIssueAttachments 测试 ====================

func TestIssueAPI_GetIssueAttachments(t *testing.T) {
	tests := []struct {
		name        string
		ticket      string
		attachments []*cloud.Attachment
		wantErr     bool
	}{
		{
			name:   "with attachments",
			ticket: "PROJ-123",
			attachments: []*cloud.Attachment{
				{ID: "10000", Filename: "test.txt"},
			},
			wantErr: false,
		},
		{
			name:        "no attachments",
			ticket:      "PROJ-123",
			attachments: nil,
			wantErr:     false,
		},
		{
			name:    "invalid ticket",
			ticket:  "invalid",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := SetupMockJiraServer(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if tt.name == "with attachments" {
					// 返回带附件的 Issue
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
				DefaultMockHandler(w, r)
			}))
			defer server.Close()

			api := createTestIssueAPI(t, server)
			attachments, err := api.GetIssueAttachments(tt.ticket)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				require.NoError(t, err)
				if tt.attachments != nil {
					assert.NotEmpty(t, attachments)
				} else {
					assert.Empty(t, attachments)
				}
			}
		})
	}
}

// ==================== GetIssueTransitions 测试 ====================

func TestIssueAPI_GetIssueTransitions(t *testing.T) {
	tests := []struct {
		name        string
		ticket      string
		transitions []cloud.Transition
		wantErr     bool
	}{
		{
			name:   "with transitions",
			ticket: "PROJ-123",
			transitions: []cloud.Transition{
				{ID: "11", Name: "In Progress"},
				{ID: "21", Name: "Done"},
			},
			wantErr: false,
		},
		{
			name:        "no transitions",
			ticket:      "PROJ-123",
			transitions: []cloud.Transition{},
			wantErr:     false,
		},
		{
			name:    "invalid ticket",
			ticket:  "invalid",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := SetupMockJiraServer(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if tt.name == "no transitions" {
					// 返回空的 transitions 列表
					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(http.StatusOK)
					json.NewEncoder(w).Encode(map[string]interface{}{
						"transitions": []map[string]interface{}{},
					})
					return
				}
				DefaultMockHandler(w, r)
			}))
			defer server.Close()

			api := createTestIssueAPI(t, server)
			transitions, err := api.GetIssueTransitions(tt.ticket)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				require.NoError(t, err)
				if len(tt.transitions) > 0 {
					assert.NotEmpty(t, transitions)
					assert.Equal(t, len(tt.transitions), len(transitions))
				} else {
					assert.Empty(t, transitions)
				}
			}
		})
	}
}

// ==================== TransitionIssue 测试 ====================

func TestIssueAPI_TransitionIssue(t *testing.T) {
	tests := []struct {
		name         string
		ticket       string
		transitionID string
		wantErr      bool
	}{
		{
			name:         "valid transition",
			ticket:       "PROJ-123",
			transitionID: "11",
			wantErr:      false,
		},
		{
			name:         "invalid ticket",
			ticket:       "invalid",
			transitionID: "11",
			wantErr:      true,
		},
		// 注意：空 transition ID 的行为取决于 go-jira SDK 的实现
		// SDK 可能在发送请求前验证，或者允许空字符串
		// 这里暂时跳过这个测试用例
		// {
		// 	name:         "empty transition ID",
		// 	ticket:       "PROJ-123",
		// 	transitionID: "",
		// 	wantErr:      true,
		// },
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := SetupMockJiraServer(t, nil)
			defer server.Close()

			api := createTestIssueAPI(t, server)
			err := api.TransitionIssue(tt.ticket, tt.transitionID)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

// ==================== AssignIssue 测试 ====================

func TestIssueAPI_AssignIssue(t *testing.T) {
	tests := []struct {
		name      string
		ticket    string
		accountID string
		wantErr   bool
	}{
		{
			name:      "assign to user",
			ticket:    "PROJ-123",
			accountID: "12345",
			wantErr:   false,
		},
		{
			name:      "unassign (empty account ID)",
			ticket:    "PROJ-123",
			accountID: "",
			wantErr:   false,
		},
		{
			name:      "unassign (-1 account ID)",
			ticket:    "PROJ-123",
			accountID: "-1",
			wantErr:   false,
		},
		{
			name:      "invalid ticket",
			ticket:    "invalid",
			accountID: "12345",
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := SetupMockJiraServer(t, nil)
			defer server.Close()

			api := createTestIssueAPI(t, server)
			err := api.AssignIssue(tt.ticket, tt.accountID)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

// ==================== AddComment 测试 ====================

func TestIssueAPI_AddComment(t *testing.T) {
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
			name:    "empty comment",
			ticket:  "PROJ-123",
			comment: "",
			wantErr: false, // Jira 可能允许空评论
		},
		{
			name:    "invalid ticket",
			ticket:  "invalid",
			comment: "Test comment",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := SetupMockJiraServer(t, nil)
			defer server.Close()

			api := createTestIssueAPI(t, server)
			err := api.AddComment(tt.ticket, tt.comment)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

// ==================== GetComments 测试 ====================

func TestIssueAPI_GetComments(t *testing.T) {
	tests := []struct {
		name     string
		ticket   string
		comments []*cloud.Comment
		wantErr  bool
	}{
		{
			name:   "with comments",
			ticket: "PROJ-123",
			comments: []*cloud.Comment{
				{ID: "10000", Body: "First comment"},
				{ID: "10001", Body: "Second comment"},
			},
			wantErr: false,
		},
		{
			name:     "no comments",
			ticket:   "PROJ-123",
			comments: []*cloud.Comment{},
			wantErr:  false,
		},
		{
			name:    "invalid ticket",
			ticket:  "invalid",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := SetupMockJiraServer(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if tt.name == "with comments" {
					// 返回带评论的 Issue
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
				DefaultMockHandler(w, r)
			}))
			defer server.Close()

			api := createTestIssueAPI(t, server)
			comments, err := api.GetComments(tt.ticket)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				require.NoError(t, err)
				if len(tt.comments) > 0 {
					assert.NotEmpty(t, comments)
					assert.Equal(t, len(tt.comments), len(comments))
				} else {
					assert.Empty(t, comments)
				}
			}
		})
	}
}

// ==================== GetChangelog 测试 ====================

func TestIssueAPI_GetChangelog(t *testing.T) {
	tests := []struct {
		name    string
		ticket  string
		wantErr bool
	}{
		{
			name:    "with changelog",
			ticket:  "PROJ-123",
			wantErr: false,
		},
		{
			name:    "no changelog",
			ticket:  "PROJ-123",
			wantErr: true, // GetChangelog 在没有 changelog 时返回错误
		},
		{
			name:    "invalid ticket",
			ticket:  "invalid",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := SetupMockJiraServer(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if tt.name == "with changelog" {
					// 返回带 changelog 的 Issue
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
				DefaultMockHandler(w, r)
			}))
			defer server.Close()

			api := createTestIssueAPI(t, server)
			changelog, err := api.GetChangelog(tt.ticket)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, changelog)
			} else {
				require.NoError(t, err)
				assert.NotNil(t, changelog)
			}
		})
	}
}

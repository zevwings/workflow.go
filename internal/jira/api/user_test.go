package api

import (
	"testing"

	"github.com/andygrunwald/go-jira/v2/cloud"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ==================== NewUserAPI 测试 ====================

func TestNewUserAPI(t *testing.T) {
	server := SetupMockJiraServer(t, nil)
	defer server.Close()

	client, ctx := CreateTestClient(t, server.URL)
	api := NewUserAPI(client, ctx)

	assert.NotNil(t, api)
	assert.NotNil(t, api.client)
	assert.NotNil(t, api.ctx)
}

// ==================== GetCurrentUser 测试 ====================

func TestUserAPI_GetCurrentUser(t *testing.T) {
	tests := []struct {
		name    string
		wantErr bool
	}{
		{
			name:    "get current user",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := SetupMockJiraServer(t, nil)
			defer server.Close()

			api := createTestUserAPI(t, server)
			user, err := api.GetCurrentUser()

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, user)
			} else {
				require.NoError(t, err)
				assert.NotNil(t, user)
				assert.NotEmpty(t, user.AccountID)
			}
		})
	}
}

// ==================== GetUser 测试 ====================

func TestUserAPI_GetUser(t *testing.T) {
	tests := []struct {
		name      string
		accountID string
		wantErr   bool
	}{
		{
			name:      "valid account ID",
			accountID: "12345",
			wantErr:   false,
		},
		{
			name:      "empty account ID",
			accountID: "",
			wantErr:   true,
		},
		{
			name:      "invalid account ID",
			accountID: "invalid",
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := SetupMockJiraServer(t, nil)
			defer server.Close()

			api := createTestUserAPI(t, server)
			user, err := api.GetUser(tt.accountID)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, user)
			} else {
				require.NoError(t, err)
				assert.NotNil(t, user)
				assert.Equal(t, tt.accountID, user.AccountID)
			}
		})
	}
}

// ==================== FindUsers 测试 ====================

func TestUserAPI_FindUsers(t *testing.T) {
	tests := []struct {
		name    string
		query   string
		users   []*cloud.User
		wantErr bool
	}{
		{
			name:  "find users by name",
			query: "john",
			users: []*cloud.User{
				{AccountID: "123", DisplayName: "John Doe"},
			},
			wantErr: false,
		},
		{
			name:  "find users by email",
			query: "john@example.com",
			users: []*cloud.User{
				{AccountID: "123", DisplayName: "John Doe", EmailAddress: "john@example.com"},
			},
			wantErr: false,
		},
		{
			name:    "empty query",
			query:   "",
			users:   nil,
			wantErr: false, // 可能返回空列表或错误，取决于 Jira API
		},
		{
			name:    "no results",
			query:   "nonexistent",
			users:   []*cloud.User{},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := SetupMockJiraServer(t, nil)
			defer server.Close()

			api := createTestUserAPI(t, server)
			users, err := api.FindUsers(tt.query)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				require.NoError(t, err)
				if len(tt.users) > 0 {
					assert.NotEmpty(t, users)
					assert.Equal(t, len(tt.users), len(users))
				} else {
					// 空查询可能返回空列表或错误，取决于 Jira API
					// 这里我们只验证不会 panic
					_ = users
				}
			}
		})
	}
}

package api

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ==================== NewProjectAPI 测试 ====================

func TestNewProjectAPI(t *testing.T) {
	server := SetupMockJiraServer(t, nil)
	defer server.Close()

	client, ctx := CreateTestClient(t, server.URL)
	projectAPI := NewProjectAPI(client, ctx)

	assert.NotNil(t, projectAPI)
	assert.NotNil(t, projectAPI.client)
	assert.NotNil(t, projectAPI.ctx)
}

// ==================== GetProject 测试 ====================

func TestProjectAPI_GetProject(t *testing.T) {
	tests := []struct {
		name       string
		projectKey string
		wantErr    bool
	}{
		{
			name:       "valid project key",
			projectKey: "PROJ",
			wantErr:    false,
		},
		{
			name:       "empty project key",
			projectKey: "",
			wantErr:    true,
		},
		{
			name:       "lowercase project key",
			projectKey: "proj",
			wantErr:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := SetupMockJiraServer(t, nil)
			defer server.Close()

			projectAPI := createTestProjectAPI(t, server)
			project, err := projectAPI.GetProject(tt.projectKey)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, project)
			} else {
				require.NoError(t, err)
				assert.NotNil(t, project)
				if tt.projectKey != "" {
					assert.Equal(t, tt.projectKey, project.Key)
				}
			}
		})
	}
}

// ==================== GetProjectStatuses 测试 ====================

func TestProjectAPI_GetProjectStatuses(t *testing.T) {
	tests := []struct {
		name       string
		projectKey string
		wantErr    bool
	}{
		{
			name:       "valid project key",
			projectKey: "PROJ",
			wantErr:    false,
		},
		{
			name:       "empty project key",
			projectKey: "",
			wantErr:    false, // 当前实现返回空列表，不返回错误
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := SetupMockJiraServer(t, nil)
			defer server.Close()

			projectAPI := createTestProjectAPI(t, server)
			statuses, err := projectAPI.GetProjectStatuses(tt.projectKey)

			// 当前实现总是返回空列表
			require.NoError(t, err)
			assert.Empty(t, statuses)
		})
	}
}

// ==================== ListProjects 测试 ====================

func TestProjectAPI_ListProjects(t *testing.T) {
	tests := []struct {
		name    string
		wantErr bool
	}{
		{
			name:    "list all projects",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := SetupMockJiraServer(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				// 返回项目列表
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				json.NewEncoder(w).Encode([]map[string]interface{}{
					{
						"id":   "10000",
						"key":  "PROJ",
						"name": "Test Project",
					},
				})
			}))
			defer server.Close()

			projectAPI := createTestProjectAPI(t, server)
			projects, err := projectAPI.ListProjects()

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.NotEmpty(t, projects)
			}
		})
	}
}

// ==================== GetProjectStatuses 返回值测试 ====================

func TestProjectAPI_GetProjectStatuses_ReturnsEmpty(t *testing.T) {
	// 测试当前实现返回空列表的行为
	// 根据代码注释，当前实现总是返回空列表
	// 这个测试验证行为是否符合预期

	t.Run("returns empty list", func(t *testing.T) {
		server := SetupMockJiraServer(t, nil)
		defer server.Close()

		projectAPI := createTestProjectAPI(t, server)
		statuses, err := projectAPI.GetProjectStatuses("PROJ")

		require.NoError(t, err)
		assert.Empty(t, statuses)
	})
}

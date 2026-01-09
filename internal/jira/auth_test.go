package jira

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zevwings/workflow/internal/jira/api"
)

// ==================== ValidateAuth 测试 ====================

func TestValidateAuth(t *testing.T) {
	tests := []struct {
		name       string
		config     *Config
		setupMock  func() *httptest.Server
		wantValid  bool
		wantErr    bool
		checkError func(t *testing.T, result *AuthResult)
	}{
		{
			name:   "nil config",
			config: nil,
			setupMock: func() *httptest.Server {
				return nil
			},
			wantValid: false,
			wantErr:   false,
			checkError: func(t *testing.T, result *AuthResult) {
				assert.False(t, result.Valid)
				assert.Equal(t, "Jira 配置为空", result.Message)
			},
		},
		{
			name: "empty service address",
			config: &Config{
				ServiceAddress: "",
				Email:          "test@example.com",
				APIToken:       "test-token",
			},
			setupMock: func() *httptest.Server {
				return nil
			},
			wantValid: false,
			wantErr:   false,
			checkError: func(t *testing.T, result *AuthResult) {
				assert.False(t, result.Valid)
				assert.Equal(t, "Jira Service Address 未配置", result.Message)
			},
		},
		{
			name: "empty email",
			config: &Config{
				ServiceAddress: "https://test.atlassian.net",
				Email:          "",
				APIToken:       "test-token",
			},
			setupMock: func() *httptest.Server {
				return nil
			},
			wantValid: false,
			wantErr:   false,
			checkError: func(t *testing.T, result *AuthResult) {
				assert.False(t, result.Valid)
				assert.Equal(t, "Jira Email 未配置", result.Message)
			},
		},
		{
			name: "empty api token",
			config: &Config{
				ServiceAddress: "https://test.atlassian.net",
				Email:          "test@example.com",
				APIToken:       "",
			},
			setupMock: func() *httptest.Server {
				return nil
			},
			wantValid: false,
			wantErr:   false,
			checkError: func(t *testing.T, result *AuthResult) {
				assert.False(t, result.Valid)
				assert.Equal(t, "Jira API Token 未配置", result.Message)
			},
		},
		{
			name: "valid config with successful auth",
			config: &Config{
				ServiceAddress: "https://test.atlassian.net",
				Email:          "test@example.com",
				APIToken:       "test-token",
			},
			setupMock: func() *httptest.Server {
				server := api.SetupMockJiraServer(t, nil)
				return server
			},
			wantValid: true,
			wantErr:   false,
			checkError: func(t *testing.T, result *AuthResult) {
				assert.True(t, result.Valid)
				assert.Equal(t, "Jira 认证成功", result.Message)
				assert.NotNil(t, result.Details)
				// 检查用户信息
				if displayName, ok := result.Details["display_name"].(string); ok {
					assert.NotEmpty(t, displayName)
				}
			},
		},
		{
			name: "valid config with auth failure",
			config: &Config{
				ServiceAddress: "https://test.atlassian.net",
				Email:          "test@example.com",
				APIToken:       "invalid-token",
			},
			setupMock: func() *httptest.Server {
				// 创建返回 401 错误的 mock 服务器
				server := api.SetupMockJiraServer(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					if r.URL.Path == "/rest/api/2/myself" || r.URL.Path == "/rest/api/3/myself" {
						w.WriteHeader(http.StatusUnauthorized)
						json.NewEncoder(w).Encode(map[string]interface{}{
							"errorMessages": []string{"Authentication failed"},
						})
						return
					}
					api.DefaultMockHandler(w, r)
				}))
				return server
			},
			wantValid: false,
			wantErr:   false,
			checkError: func(t *testing.T, result *AuthResult) {
				assert.False(t, result.Valid)
				assert.Equal(t, "Jira 认证失败", result.Message)
				assert.NotNil(t, result.Error)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var server *httptest.Server
			if tt.setupMock != nil {
				mockServer := tt.setupMock()
				if mockServer != nil {
					server = mockServer
					defer server.Close()
					// 更新 config 的 ServiceAddress 为 mock 服务器地址
					tt.config.ServiceAddress = server.URL
				}
			}

			result, err := ValidateAuth(tt.config)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				require.NoError(t, err)
				require.NotNil(t, result)
				assert.Equal(t, tt.wantValid, result.Valid)
				if tt.checkError != nil {
					tt.checkError(t, result)
				}
			}
		})
	}
}

func TestValidateAuth_Success(t *testing.T) {
	server := api.SetupMockJiraServer(t, nil)
	defer server.Close()

	config := &Config{
		ServiceAddress: server.URL,
		Email:          "test@example.com",
		APIToken:       "test-token",
	}

	result, err := ValidateAuth(config)
	require.NoError(t, err)
	require.NotNil(t, result)

	assert.True(t, result.Valid)
	assert.Equal(t, "Jira 认证成功", result.Message)
	assert.Nil(t, result.Error)
	assert.NotNil(t, result.Details)

	// 验证用户信息
	if displayName, ok := result.Details["display_name"].(string); ok {
		assert.NotEmpty(t, displayName)
	}
	if accountID, ok := result.Details["account_id"].(string); ok {
		assert.NotEmpty(t, accountID)
	}
}

func TestValidateAuth_InvalidToken(t *testing.T) {
	server := api.SetupMockJiraServer(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 返回 401 未授权错误
		if r.URL.Path == "/rest/api/2/myself" || r.URL.Path == "/rest/api/3/myself" {
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"errorMessages": []string{"Authentication failed"},
			})
			return
		}
		api.DefaultMockHandler(w, r)
	}))
	defer server.Close()

	config := &Config{
		ServiceAddress: server.URL,
		Email:          "test@example.com",
		APIToken:       "invalid-token",
	}

	result, err := ValidateAuth(config)
	require.NoError(t, err)
	require.NotNil(t, result)

	assert.False(t, result.Valid)
	assert.Equal(t, "Jira 认证失败", result.Message)
	assert.NotNil(t, result.Error)
}

func TestValidateAuth_ClientCreationFailure(t *testing.T) {
	// 使用无效的 URL 导致客户端创建失败
	config := &Config{
		ServiceAddress: "invalid-url",
		Email:          "test@example.com",
		APIToken:       "test-token",
	}

	result, err := ValidateAuth(config)
	require.NoError(t, err)
	require.NotNil(t, result)

	// 客户端创建可能失败，但验证方法应该返回结果而不是错误
	// 实际行为取决于 go-jira SDK 如何处理无效 URL
	// 这里主要验证方法不会 panic
	assert.NotNil(t, result)
}

//go:build integration

package jira

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ==================== 集成测试环境设置 ====================

// getIntegrationConfig 从环境变量获取 Jira 配置
func getIntegrationConfig(t *testing.T) *Config {
	t.Helper()

	url := os.Getenv("JIRA_URL")
	username := os.Getenv("JIRA_USERNAME")
	token := os.Getenv("JIRA_API_TOKEN")

	if url == "" || username == "" || token == "" {
		t.Skip("跳过集成测试：需要设置环境变量 JIRA_URL, JIRA_USERNAME, JIRA_API_TOKEN")
	}

	return &Config{
		URL:      url,
		Username: username,
		Token:    token,
	}
}

// ==================== NewJiraClient 集成测试 ====================

func TestNewJiraClient_Integration(t *testing.T) {
	config := getIntegrationConfig(t)

	client, err := NewJiraClient(config)
	require.NoError(t, err)
	assert.NotNil(t, client)
	assert.NotNil(t, client.GetClient())
	assert.NotNil(t, client.GetIssueAPI())
	assert.NotNil(t, client.GetProjectAPI())
	assert.NotNil(t, client.GetUserAPI())
}

// ==================== GetUserInfo 集成测试 ====================

func TestJiraClient_GetUserInfo_Integration(t *testing.T) {
	config := getIntegrationConfig(t)

	client, err := NewJiraClient(config)
	require.NoError(t, err)

	user, err := client.GetUserInfo()
	require.NoError(t, err)
	assert.NotNil(t, user)
	assert.NotEmpty(t, user.AccountID)
	assert.NotEmpty(t, user.DisplayName)
}

// ==================== GetTicketInfo 集成测试 ====================

func TestJiraClient_GetTicketInfo_Integration(t *testing.T) {
	config := getIntegrationConfig(t)

	client, err := NewJiraClient(config)
	require.NoError(t, err)

	// 注意：需要提供一个真实存在的 ticket key
	// 可以通过环境变量配置，或者使用测试项目中的固定 ticket
	testTicket := os.Getenv("JIRA_TEST_TICKET")
	if testTicket == "" {
		t.Skip("跳过测试：需要设置环境变量 JIRA_TEST_TICKET")
	}

	issue, err := client.GetTicketInfo(testTicket)
	require.NoError(t, err)
	assert.NotNil(t, issue)
	assert.Equal(t, testTicket, issue.Key)
	assert.NotEmpty(t, issue.Fields.Summary)
}

// ==================== GetAttachments 集成测试 ====================

func TestJiraClient_GetAttachments_Integration(t *testing.T) {
	config := getIntegrationConfig(t)

	client, err := NewJiraClient(config)
	require.NoError(t, err)

	testTicket := os.Getenv("JIRA_TEST_TICKET")
	if testTicket == "" {
		t.Skip("跳过测试：需要设置环境变量 JIRA_TEST_TICKET")
	}

	attachments, err := client.GetAttachments(testTicket)
	require.NoError(t, err)
	// 附件可能为空，这是正常的
	assert.NotNil(t, attachments)
}

// ==================== GetTransitions 集成测试 ====================

func TestJiraClient_GetTransitions_Integration(t *testing.T) {
	config := getIntegrationConfig(t)

	client, err := NewJiraClient(config)
	require.NoError(t, err)

	testTicket := os.Getenv("JIRA_TEST_TICKET")
	if testTicket == "" {
		t.Skip("跳过测试：需要设置环境变量 JIRA_TEST_TICKET")
	}

	transitions, err := client.GetTransitions(testTicket)
	require.NoError(t, err)
	// Transitions 可能为空，取决于 ticket 的当前状态
	assert.NotNil(t, transitions)
}

// ==================== GetComments 集成测试 ====================

func TestJiraClient_GetComments_Integration(t *testing.T) {
	config := getIntegrationConfig(t)

	client, err := NewJiraClient(config)
	require.NoError(t, err)

	testTicket := os.Getenv("JIRA_TEST_TICKET")
	if testTicket == "" {
		t.Skip("跳过测试：需要设置环境变量 JIRA_TEST_TICKET")
	}

	comments, err := client.GetComments(testTicket)
	require.NoError(t, err)
	// 评论可能为空，这是正常的
	assert.NotNil(t, comments)
}

// ==================== GetProject 集成测试 ====================

func TestJiraClient_GetProject_Integration(t *testing.T) {
	config := getIntegrationConfig(t)

	client, err := NewJiraClient(config)
	require.NoError(t, err)

	// 从 test ticket 提取项目 key，或使用环境变量
	testTicket := os.Getenv("JIRA_TEST_TICKET")
	if testTicket == "" {
		t.Skip("跳过测试：需要设置环境变量 JIRA_TEST_TICKET")
	}

	projectKey := ExtractProjectKey(testTicket)
	require.NotEmpty(t, projectKey)

	project, err := client.GetProject(projectKey)
	require.NoError(t, err)
	assert.NotNil(t, project)
	assert.Equal(t, projectKey, project.Key)
	assert.NotEmpty(t, project.Name)
}

// ==================== GetProjectStatuses 集成测试 ====================

func TestJiraClient_GetProjectStatuses_Integration(t *testing.T) {
	config := getIntegrationConfig(t)

	client, err := NewJiraClient(config)
	require.NoError(t, err)

	testTicket := os.Getenv("JIRA_TEST_TICKET")
	if testTicket == "" {
		t.Skip("跳过测试：需要设置环境变量 JIRA_TEST_TICKET")
	}

	projectKey := ExtractProjectKey(testTicket)
	require.NotEmpty(t, projectKey)

	statuses, err := client.GetProjectStatuses(projectKey)
	require.NoError(t, err)
	// 当前实现返回空列表
	assert.Empty(t, statuses)
}

// ==================== FindUsers 集成测试 ====================

func TestJiraClient_FindUsers_Integration(t *testing.T) {
	config := getIntegrationConfig(t)

	client, err := NewJiraClient(config)
	require.NoError(t, err)

	// 搜索当前用户（应该能找到）
	username := config.Username
	users, err := client.FindUsers(username)
	require.NoError(t, err)
	// 应该至少找到一个用户（当前用户）
	assert.NotEmpty(t, users)
}

// ==================== 端到端流程测试 ====================

func TestJiraClient_EndToEnd_Integration(t *testing.T) {
	config := getIntegrationConfig(t)

	client, err := NewJiraClient(config)
	require.NoError(t, err)

	// 1. 获取当前用户信息
	user, err := client.GetUserInfo()
	require.NoError(t, err)
	assert.NotNil(t, user)

	// 2. 获取测试 ticket 信息
	testTicket := os.Getenv("JIRA_TEST_TICKET")
	if testTicket == "" {
		t.Skip("跳过测试：需要设置环境变量 JIRA_TEST_TICKET")
	}

	issue, err := client.GetTicketInfo(testTicket)
	require.NoError(t, err)
	assert.NotNil(t, issue)

	// 3. 获取可用状态转换
	transitions, err := client.GetTransitions(testTicket)
	require.NoError(t, err)
	assert.NotNil(t, transitions)

	// 4. 获取项目信息
	projectKey := ExtractProjectKey(testTicket)
	project, err := client.GetProject(projectKey)
	require.NoError(t, err)
	assert.NotNil(t, project)

	// 5. 搜索用户
	users, err := client.FindUsers(user.DisplayName)
	require.NoError(t, err)
	assert.NotEmpty(t, users)
}

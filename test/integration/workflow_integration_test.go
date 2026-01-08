//go:build integration

package integration

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/zevwings/workflow/internal/jira"
)

// TestWorkflow_JiraClient_Integration
//
// 说明：
// - 这是一个示例性的跨包集成测试入口
// - 当前仅做占位和环境检查，避免影响日常开发
// - 后续可以在这里串联 Jira、GitHub、PR 等完整工作流
func TestWorkflow_JiraClient_Integration(t *testing.T) {
	// 环境变量依赖与 internal/jira/jira_client_integration_test.go 保持一致
	config := getIntegrationConfigForTest(t)

	// 验证可以创建 JiraClient（跨包集成入口）
	client, err := jira.NewJiraClient(config)
	require.NoError(t, err)
	require.NotNil(t, client)
}

// getIntegrationConfigForTest 为跨包集成测试提供最小配置
//
// 这里直接复用 jira.Config，保持与包内集成测试一致的配置结构。
func getIntegrationConfigForTest(t *testing.T) *jira.Config {
	t.Helper()

	// 这里不直接读取环境变量，避免与包内集成测试产生强耦合。
	// 先作为占位：如果环境未配置，直接跳过。
	t.Skip("跨包集成测试占位：需要根据实际场景补充配置和环境依赖")

	return &jira.Config{}
}



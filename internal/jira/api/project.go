package api

import (
	"context"
	"fmt"

	"github.com/andygrunwald/go-jira/v2/cloud"
	"github.com/zevwings/workflow/internal/logging"
)

// ProjectAPI 提供 Project 相关的 REST API 方法
type ProjectAPI struct {
	client *cloud.Client
	ctx    context.Context
}

// NewProjectAPI 创建新的 Project API 实例
func NewProjectAPI(client *cloud.Client, ctx context.Context) *ProjectAPI {
	return &ProjectAPI{
		client: client,
		ctx:    ctx,
	}
}

// GetProject 获取项目信息
//
// 参数:
//   - projectKey: 项目 Key（如 "PROJ"）
//
// 返回:
//   - *cloud.Project: 项目信息
//   - error: 如果获取失败，返回错误
func (api *ProjectAPI) GetProject(projectKey string) (*cloud.Project, error) {
	logger := logging.GetLogger()
	logger.Infof("Jira API call: GetProject(%s)", projectKey)

	project, _, err := api.client.Project.Get(api.ctx, projectKey)
	if err != nil {
		logger.WithError(err).Errorf("Jira API call failed: GetProject(%s)", projectKey)
		return nil, fmt.Errorf("获取项目 %s 失败: %w", projectKey, err)
	}

	return project, nil
}

// GetProjectStatuses 获取项目的状态列表
//
// 注意：此方法通过获取项目信息来获取状态列表。
// Jira Cloud API v2 不直接提供获取项目状态的接口，
// 需要通过获取项目的 Issue 来推断可用的状态。
//
// 参数:
//   - projectKey: 项目 Key（如 "PROJ"）
//
// 返回:
//   - []cloud.Status: 状态列表（空列表表示无法获取）
//   - error: 如果获取失败，返回错误
func (api *ProjectAPI) GetProjectStatuses(projectKey string) ([]cloud.Status, error) {
	// Jira Cloud API v2 不直接提供获取项目状态的接口
	// 这里返回空列表，实际使用时可以通过 Issue 的 transitions 来获取
	return []cloud.Status{}, nil
}

// ListProjects 列出所有项目
//
// 返回:
//   - []*cloud.Project: 项目列表
//   - error: 如果获取失败，返回错误
func (api *ProjectAPI) ListProjects() ([]*cloud.Project, error) {
	projectList, _, err := api.client.Project.GetAll(api.ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("获取项目列表失败: %w", err)
	}

	// ProjectList 是一个匿名结构体的 slice，需要手动转换为 Project
	result := make([]*cloud.Project, len(*projectList))
	for i, p := range *projectList {
		result[i] = &cloud.Project{
			Expand:          p.Expand,
			Self:            p.Self,
			ID:              p.ID,
			Key:             p.Key,
			Name:            p.Name,
			AvatarUrls:      p.AvatarUrls,
			ProjectCategory: p.ProjectCategory,
			IssueTypes:      p.IssueTypes,
		}
	}

	return result, nil
}


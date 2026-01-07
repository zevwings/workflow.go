package jira

import (
	"fmt"

	"github.com/andygrunwald/go-jira/v2/cloud"
	"github.com/zevwings/workflow/internal/jira/api"
)

// JiraClient Jira REST API 客户端（向后兼容包装器）
//
// 提供高级封装，简化常用操作。
// 内部使用 go-jira SDK 和 API 模块。
type JiraClient struct {
	client   *Client
	issueAPI *api.IssueAPI
	projectAPI *api.ProjectAPI
	userAPI *api.UserAPI
}

// NewJiraClient 创建新的 JiraClient 实例
//
// 参数:
//   - config: Jira 配置结构体（不能为 nil）
//
// 返回:
//   - *JiraClient: JiraClient 实例
//   - error: 如果创建失败，返回错误
func NewJiraClient(config *Config) (*JiraClient, error) {
	client, err := NewClient(config)
	if err != nil {
		return nil, err
	}

	jiraClient := client.GetJiraClient()
	ctx := client.GetContext()

	return &JiraClient{
		client:     client,
		issueAPI:   api.NewIssueAPI(jiraClient, ctx),
		projectAPI: api.NewProjectAPI(jiraClient, ctx),
		userAPI:    api.NewUserAPI(jiraClient, ctx),
	}, nil
}

// GetUserInfo 获取当前 Jira 用户信息
//
// 返回:
//   - *cloud.User: 当前用户信息
//   - error: 如果获取失败，返回错误
func (c *JiraClient) GetUserInfo() (*cloud.User, error) {
	return c.userAPI.GetCurrentUser()
}

// GetTicketInfo 获取 ticket 信息
//
// 参数:
//   - ticket: Ticket Key（如 "PROJ-123"）
//
// 返回:
//   - *cloud.Issue: Issue 信息
//   - error: 如果获取失败，返回错误
func (c *JiraClient) GetTicketInfo(ticket string) (*cloud.Issue, error) {
	if err := ValidateTicketKey(ticket); err != nil {
		return nil, err
	}

	ticket = NormalizeTicketKey(ticket)
	return c.issueAPI.GetIssue(ticket)
}

// GetAttachments 获取 ticket 的附件列表
//
// 参数:
//   - ticket: Ticket Key（如 "PROJ-123"）
//
// 返回:
//   - []*cloud.Attachment: 附件列表
//   - error: 如果获取失败，返回错误
func (c *JiraClient) GetAttachments(ticket string) ([]*cloud.Attachment, error) {
	if err := ValidateTicketKey(ticket); err != nil {
		return nil, err
	}

	ticket = NormalizeTicketKey(ticket)
	return c.issueAPI.GetIssueAttachments(ticket)
}

// MoveTicket 更新 ticket 状态
//
// 参数:
//   - ticket: Ticket Key（如 "PROJ-123"）
//   - status: 目标状态名称（如 "In Progress"）
//
// 返回:
//   - error: 如果更新失败，返回错误
func (c *JiraClient) MoveTicket(ticket, status string) error {
	if err := ValidateTicketKey(ticket); err != nil {
		return err
	}

	ticket = NormalizeTicketKey(ticket)

	// 获取可用的 transitions
	transitions, err := c.issueAPI.GetIssueTransitions(ticket)
	if err != nil {
		return err
	}

	// 查找匹配的状态转换
	var transitionID string
	for _, t := range transitions {
		if t.Name == status {
			transitionID = t.ID
			break
		}
	}

	if transitionID == "" {
		return fmt.Errorf("未找到状态转换: %s", status)
	}

	return c.issueAPI.TransitionIssue(ticket, transitionID)
}

// AssignTicket 分配 ticket 给用户
//
// 参数:
//   - ticket: Ticket Key（如 "PROJ-123"）
//   - accountID: 用户 Account ID（如果为 nil 或空字符串，则取消分配）
//
// 返回:
//   - error: 如果分配失败，返回错误
func (c *JiraClient) AssignTicket(ticket string, accountID *string) error {
	if err := ValidateTicketKey(ticket); err != nil {
		return err
	}

	ticket = NormalizeTicketKey(ticket)

	var id string
	if accountID != nil {
		id = *accountID
	}

	return c.issueAPI.AssignIssue(ticket, id)
}

// AddComment 添加评论到 ticket
//
// 参数:
//   - ticket: Ticket Key（如 "PROJ-123"）
//   - comment: 评论内容
//
// 返回:
//   - error: 如果添加失败，返回错误
func (c *JiraClient) AddComment(ticket, comment string) error {
	if err := ValidateTicketKey(ticket); err != nil {
		return err
	}

	ticket = NormalizeTicketKey(ticket)
	return c.issueAPI.AddComment(ticket, comment)
}

// GetComments 获取 ticket 的评论列表
//
// 参数:
//   - ticket: Ticket Key（如 "PROJ-123"）
//
// 返回:
//   - []*cloud.Comment: 评论列表
//   - error: 如果获取失败，返回错误
func (c *JiraClient) GetComments(ticket string) ([]*cloud.Comment, error) {
	if err := ValidateTicketKey(ticket); err != nil {
		return nil, err
	}

	ticket = NormalizeTicketKey(ticket)
	return c.issueAPI.GetComments(ticket)
}

// UploadAttachment 上传附件到 ticket
//
// 参数:
//   - ticket: Ticket Key（如 "PROJ-123"）
//   - filePath: 文件路径
//
// 返回:
//   - []*cloud.Attachment: 上传后的附件列表
//   - error: 如果上传失败，返回错误
func (c *JiraClient) UploadAttachment(ticket, filePath string) ([]*cloud.Attachment, error) {
	if err := ValidateTicketKey(ticket); err != nil {
		return nil, err
	}

	ticket = NormalizeTicketKey(ticket)
	return c.issueAPI.UploadAttachment(ticket, filePath)
}

// GetTransitions 获取 ticket 的可用状态转换
//
// 参数:
//   - ticket: Ticket Key（如 "PROJ-123"）
//
// 返回:
//   - []cloud.Transition: 可用的状态转换列表
//   - error: 如果获取失败，返回错误
func (c *JiraClient) GetTransitions(ticket string) ([]cloud.Transition, error) {
	if err := ValidateTicketKey(ticket); err != nil {
		return nil, err
	}

	ticket = NormalizeTicketKey(ticket)
	return c.issueAPI.GetIssueTransitions(ticket)
}

// GetChangelog 获取 ticket 的变更历史
//
// 参数:
//   - ticket: Ticket Key（如 "PROJ-123"）
//
// 返回:
//   - *cloud.Changelog: 变更历史
//   - error: 如果获取失败，返回错误
func (c *JiraClient) GetChangelog(ticket string) (*cloud.Changelog, error) {
	if err := ValidateTicketKey(ticket); err != nil {
		return nil, err
	}

	ticket = NormalizeTicketKey(ticket)
	return c.issueAPI.GetChangelog(ticket)
}

// GetProject 获取项目信息
//
// 参数:
//   - projectKey: 项目 Key（如 "PROJ"）
//
// 返回:
//   - *cloud.Project: 项目信息
//   - error: 如果获取失败，返回错误
func (c *JiraClient) GetProject(projectKey string) (*cloud.Project, error) {
	return c.projectAPI.GetProject(projectKey)
}

// GetProjectStatuses 获取项目的状态列表
//
// 参数:
//   - projectKey: 项目 Key（如 "PROJ"）
//
// 返回:
//   - []cloud.Status: 状态列表
//   - error: 如果获取失败，返回错误
func (c *JiraClient) GetProjectStatuses(projectKey string) ([]cloud.Status, error) {
	return c.projectAPI.GetProjectStatuses(projectKey)
}

// FindUsers 搜索用户
//
// 参数:
//   - query: 搜索关键词（用户名或邮箱）
//
// 返回:
//   - []*cloud.User: 用户列表
//   - error: 如果搜索失败，返回错误
func (c *JiraClient) FindUsers(query string) ([]*cloud.User, error) {
	return c.userAPI.FindUsers(query)
}

// GetClient 获取底层 Client（用于高级用法）
//
// 返回:
//   - *Client: 底层客户端
func (c *JiraClient) GetClient() *Client {
	return c.client
}

// GetIssueAPI 获取 Issue API（用于高级用法）
//
// 返回:
//   - *api.IssueAPI: Issue API 实例
func (c *JiraClient) GetIssueAPI() *api.IssueAPI {
	return c.issueAPI
}

// GetProjectAPI 获取 Project API（用于高级用法）
//
// 返回:
//   - *api.ProjectAPI: Project API 实例
func (c *JiraClient) GetProjectAPI() *api.ProjectAPI {
	return c.projectAPI
}

// GetUserAPI 获取 User API（用于高级用法）
//
// 返回:
//   - *api.UserAPI: User API 实例
func (c *JiraClient) GetUserAPI() *api.UserAPI {
	return c.userAPI
}


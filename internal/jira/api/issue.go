package api

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/andygrunwald/go-jira/v2/cloud"
)

// IssueAPI 提供 Issue/Ticket 相关的 REST API 方法
type IssueAPI struct {
	client *cloud.Client
	ctx    context.Context
}

// NewIssueAPI 创建新的 Issue API 实例
func NewIssueAPI(client *cloud.Client, ctx context.Context) *IssueAPI {
	return &IssueAPI{
		client: client,
		ctx:    ctx,
	}
}

// GetIssue 获取 issue 信息
//
// 参数:
//   - ticket: Issue Key（如 "PROJ-123"）
//
// 返回:
//   - *cloud.Issue: Issue 信息
//   - error: 如果获取失败，返回错误
func (api *IssueAPI) GetIssue(ticket string) (*cloud.Issue, error) {
	issue, _, err := api.client.Issue.Get(api.ctx, ticket, nil)
	if err != nil {
		return nil, fmt.Errorf("获取 issue %s 失败: %w", ticket, err)
	}

	return issue, nil
}

// GetIssueAttachments 获取 issue 的附件列表
//
// 参数:
//   - ticket: Issue Key（如 "PROJ-123"）
//
// 返回:
//   - []*cloud.Attachment: 附件列表
//   - error: 如果获取失败，返回错误
func (api *IssueAPI) GetIssueAttachments(ticket string) ([]*cloud.Attachment, error) {
	issue, err := api.GetIssue(ticket)
	if err != nil {
		return nil, err
	}

	if issue.Fields.Attachments == nil {
		return []*cloud.Attachment{}, nil
	}

	return issue.Fields.Attachments, nil
}

// GetIssueTransitions 获取 issue 的可用 transitions
//
// 参数:
//   - ticket: Issue Key（如 "PROJ-123"）
//
// 返回:
//   - []cloud.Transition: 可用的状态转换列表
//   - error: 如果获取失败，返回错误
func (api *IssueAPI) GetIssueTransitions(ticket string) ([]cloud.Transition, error) {
	transitions, _, err := api.client.Issue.GetTransitions(api.ctx, ticket)
	if err != nil {
		return nil, fmt.Errorf("获取 issue %s 的状态转换失败: %w", ticket, err)
	}

	return transitions, nil
}

// TransitionIssue 更新 issue 状态
//
// 参数:
//   - ticket: Issue Key（如 "PROJ-123"）
//   - transitionID: 状态转换 ID
//
// 返回:
//   - error: 如果更新失败，返回错误
func (api *IssueAPI) TransitionIssue(ticket, transitionID string) error {
	_, err := api.client.Issue.DoTransition(api.ctx, ticket, transitionID)
	if err != nil {
		return fmt.Errorf("更新 issue %s 状态失败: %w", ticket, err)
	}

	return nil
}

// AssignIssue 分配 issue 给用户
//
// 参数:
//   - ticket: Issue Key（如 "PROJ-123"）
//   - accountID: 用户 Account ID（如果为空字符串，则取消分配）
//
// 返回:
//   - error: 如果分配失败，返回错误
func (api *IssueAPI) AssignIssue(ticket, accountID string) error {
	var assignee *cloud.User
	if accountID != "" && accountID != "-1" {
		assignee = &cloud.User{
			AccountID: accountID,
		}
	}
	// accountID 为空或 "-1" 时，assignee 为 nil，表示取消分配

	_, err := api.client.Issue.UpdateAssignee(api.ctx, ticket, assignee)
	if err != nil {
		return fmt.Errorf("分配 issue %s 失败: %w", ticket, err)
	}

	return nil
}

// AddComment 添加评论到 issue
//
// 参数:
//   - ticket: Issue Key（如 "PROJ-123"）
//   - comment: 评论内容
//
// 返回:
//   - error: 如果添加失败，返回错误
func (api *IssueAPI) AddComment(ticket, comment string) error {
	newComment := &cloud.Comment{
		Body: comment,
	}

	_, _, err := api.client.Issue.AddComment(api.ctx, ticket, newComment)
	if err != nil {
		return fmt.Errorf("添加评论到 issue %s 失败: %w", ticket, err)
	}

	return nil
}

// GetComments 获取 issue 的评论列表
//
// 参数:
//   - ticket: Issue Key（如 "PROJ-123"）
//
// 返回:
//   - []*cloud.Comment: 评论列表
//   - error: 如果获取失败，返回错误
func (api *IssueAPI) GetComments(ticket string) ([]*cloud.Comment, error) {
	issue, err := api.GetIssue(ticket)
	if err != nil {
		return nil, err
	}

	if issue.Fields.Comments == nil || len(issue.Fields.Comments.Comments) == 0 {
		return []*cloud.Comment{}, nil
	}

	return issue.Fields.Comments.Comments, nil
}

// UploadAttachment 上传附件到 issue
//
// 参数:
//   - ticket: Issue Key（如 "PROJ-123"）
//   - filePath: 文件路径
//
// 返回:
//   - []*cloud.Attachment: 上传后的附件列表
//   - error: 如果上传失败，返回错误
func (api *IssueAPI) UploadAttachment(ticket, filePath string) ([]*cloud.Attachment, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("打开文件失败: %w", err)
	}
	defer file.Close()

	fileName := filepath.Base(filePath)
	attachments, _, err := api.client.Issue.PostAttachment(api.ctx, ticket, file, fileName)
	if err != nil {
		return nil, fmt.Errorf("上传附件到 issue %s 失败: %w", ticket, err)
	}

	// 转换 []Attachment 为 []*Attachment
	result := make([]*cloud.Attachment, len(*attachments))
	for i := range *attachments {
		result[i] = &(*attachments)[i]
	}

	return result, nil
}

// DownloadAttachment 下载附件
//
// 参数:
//   - attachment: 附件对象
//
// 返回:
//   - io.ReadCloser: 附件内容流
//   - error: 如果下载失败，返回错误
func (api *IssueAPI) DownloadAttachment(attachment *cloud.Attachment) (io.ReadCloser, error) {
	resp, err := api.client.Issue.DownloadAttachment(api.ctx, attachment.ID)
	if err != nil {
		return nil, fmt.Errorf("下载附件 %s 失败: %w", attachment.Filename, err)
	}

	return resp.Body, nil
}

// GetChangelog 获取 issue 的变更历史
//
// 注意：此方法通过获取完整的 Issue 信息来获取变更历史。
// 需要在 GetIssue 时使用 expand=changelog 选项。
//
// 参数:
//   - ticket: Issue Key（如 "PROJ-123"）
//
// 返回:
//   - *cloud.Changelog: 变更历史（如果存在）
//   - error: 如果获取失败，返回错误
func (api *IssueAPI) GetChangelog(ticket string) (*cloud.Changelog, error) {
	// 使用 expand 选项获取 changelog
	options := &cloud.GetQueryOptions{
		Expand: "changelog",
	}

	issue, _, err := api.client.Issue.Get(api.ctx, ticket, options)
	if err != nil {
		return nil, fmt.Errorf("获取 issue %s 的变更历史失败: %w", ticket, err)
	}

	if issue.Changelog == nil {
		return nil, fmt.Errorf("issue %s 没有变更历史", ticket)
	}

	return issue.Changelog, nil
}


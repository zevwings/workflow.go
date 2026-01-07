package pr

import (
	"context"
)

// PlatformProvider 平台提供者接口
//
// 定义统一的 PR 操作接口，支持多种平台（GitHub、GitLab、Bitbucket 等）
// 未来可以轻松扩展支持其他平台
type PlatformProvider interface {
	// CreatePullRequest 创建 Pull Request
	//
	// 参数:
	//   - title: PR 标题
	//   - body: PR 描述
	//   - sourceBranch: 源分支名
	//   - targetBranch: 目标分支名（可选，nil 时使用默认分支）
	//
	// 返回:
	//   - string: PR URL
	//   - error: 错误信息
	CreatePullRequest(ctx context.Context, title, body, sourceBranch string, targetBranch *string) (string, error)

	// MergePullRequest 合并 Pull Request
	//
	// 参数:
	//   - prID: PR ID（可以是数字、URL 等，由平台实现解析）
	//   - mergeMethod: 合并方法（"merge", "squash", "rebase"）
	//   - deleteBranch: 是否删除源分支
	//
	// 返回:
	//   - error: 错误信息
	MergePullRequest(ctx context.Context, prID string, mergeMethod string, deleteBranch bool) error

	// ClosePullRequest 关闭 Pull Request
	//
	// 参数:
	//   - prID: PR ID
	//
	// 返回:
	//   - error: 错误信息
	ClosePullRequest(ctx context.Context, prID string) error

	// GetPullRequestStatus 获取 PR 状态
	//
	// 参数:
	//   - prID: PR ID
	//
	// 返回:
	//   - *PullRequestStatus: PR 状态信息
	//   - error: 错误信息
	GetPullRequestStatus(ctx context.Context, prID string) (*PullRequestStatus, error)

	// ListPullRequests 列出 Pull Requests
	//
	// 参数:
	//   - state: 状态过滤（"open", "closed", "all"）
	//   - limit: 返回数量限制
	//
	// 返回:
	//   - []*PullRequestInfo: PR 列表
	//   - error: 错误信息
	ListPullRequests(ctx context.Context, state string, limit int) ([]*PullRequestInfo, error)

	// UpdatePullRequest 更新 Pull Request
	//
	// 参数:
	//   - prID: PR ID
	//   - title: 新标题（可选）
	//   - body: 新描述（可选）
	//   - state: 新状态（可选，"open" 或 "closed"）
	//
	// 返回:
	//   - error: 错误信息
	UpdatePullRequest(ctx context.Context, prID string, title, body *string, state *string) error

	// AddComment 添加评论
	//
	// 参数:
	//   - prID: PR ID
	//   - body: 评论内容
	//
	// 返回:
	//   - error: 错误信息
	AddComment(ctx context.Context, prID string, body string) error

	// ApprovePullRequest 批准 Pull Request
	//
	// 参数:
	//   - prID: PR ID
	//
	// 返回:
	//   - error: 错误信息
	ApprovePullRequest(ctx context.Context, prID string) error

	// GetPlatformName 获取平台名称
	//
	// 返回:
	//   - string: 平台名称（如 "github", "gitlab"）
	GetPlatformName() string
}

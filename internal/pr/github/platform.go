package github

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/go-github/v57/github"
	"github.com/zevwings/workflow/internal/pr"
	"golang.org/x/oauth2"
)

// GitHub 实现 PlatformProvider 接口
//
// 封装 google/go-github 库，提供统一的 PR 操作接口
type GitHub struct {
	client *github.Client
	ctx    context.Context
	owner  string
	repo   string
}

// NewGitHub 创建新的 GitHub 实例
//
// 使用传入的 token 和 owner/repo 参数创建 GitHub 客户端。
// 业务逻辑（如获取 owner/repo）应由调用方（commands 层）完成。
//
// 参数:
//   - token: GitHub Personal Access Token
//   - owner: 仓库所有者（如 "zevwings"）
//   - repo: 仓库名称（如 "workflow"）
//
// 返回:
//   - *GitHub: GitHub 实例
//   - error: 如果创建失败，返回错误
func NewGitHub(token, owner, repo string) (*GitHub, error) {
	if token == "" {
		return nil, fmt.Errorf("token is required")
	}
	if owner == "" || repo == "" {
		return nil, fmt.Errorf("owner and repo are required")
	}

	// 创建 OAuth 客户端
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(ctx, ts)

	// 创建 GitHub 客户端
	client := github.NewClient(tc)

	return &GitHub{
		client: client,
		ctx:    ctx,
		owner:  owner,
		repo:   repo,
	}, nil
}

// GetPlatformName 获取平台名称
func (g *GitHub) GetPlatformName() string {
	return "github"
}

// GetDefaultBranch 获取仓库的默认分支
//
// 通过 GitHub API 获取仓库的默认分支信息。
//
// 返回:
//   - string: 默认分支名（如 "main", "master"）
//   - error: 如果获取失败，返回错误
func (g *GitHub) GetDefaultBranch(ctx context.Context) (string, error) {
	repo, _, err := g.client.Repositories.Get(ctx, g.owner, g.repo)
	if err != nil {
		return "", fmt.Errorf("failed to get repository: %w", err)
	}

	if repo.DefaultBranch == nil || *repo.DefaultBranch == "" {
		return "", fmt.Errorf("repository has no default branch")
	}

	return *repo.DefaultBranch, nil
}

// CreatePullRequest 创建 Pull Request
func (g *GitHub) CreatePullRequest(ctx context.Context, title, body, sourceBranch string, targetBranch *string) (string, error) {
	// 1. 如果没有指定目标分支，获取默认分支
	baseBranch := "main"
	if targetBranch != nil && *targetBranch != "" {
		baseBranch = *targetBranch
	} else {
		// 优先通过 GitHub API 获取默认分支，失败则回退到本地 Git 仓库
		defaultBranch, err := g.GetDefaultBranch(ctx)
		if err == nil {
			baseBranch = defaultBranch
		}
		// 如果获取失败，使用默认值 "main"
	}

	// 2. 构建 PR 请求
	newPR := &github.NewPullRequest{
		Title: github.String(title),
		Body:  github.String(body),
		Head:  github.String(sourceBranch),
		Base:  github.String(baseBranch),
	}

	// 3. 创建 PR
	pr, _, err := g.client.PullRequests.Create(ctx, g.owner, g.repo, newPR)
	if err != nil {
		return "", fmt.Errorf("failed to create PR: %w", err)
	}

	return pr.GetHTMLURL(), nil
}

// MergePullRequest 合并 Pull Request
func (g *GitHub) MergePullRequest(ctx context.Context, prID string, mergeMethod string, deleteBranch bool) error {
	prNumber, err := parsePRNumber(prID)
	if err != nil {
		return err
	}

	// 验证合并方法
	validMethods := map[string]bool{
		"merge":  true,
		"squash": true,
		"rebase": true,
	}
	if !validMethods[mergeMethod] {
		mergeMethod = "squash" // 默认使用 squash
	}

	// 合并 PR
	_, _, err = g.client.PullRequests.Merge(ctx, g.owner, g.repo, prNumber, "", &github.PullRequestOptions{
		MergeMethod: mergeMethod,
	})
	if err != nil {
		return fmt.Errorf("failed to merge PR: %w", err)
	}

	// 如果需要，删除分支
	if deleteBranch {
		// 先获取 PR 信息以获取分支名
		ghPR, _, err := g.client.PullRequests.Get(ctx, g.owner, g.repo, prNumber)
		if err != nil {
			return fmt.Errorf("failed to get PR info: %w", err)
		}

		branchName := ghPR.GetHead().GetRef()
		ref := fmt.Sprintf("heads/%s", branchName)
		_, err = g.client.Git.DeleteRef(ctx, g.owner, g.repo, ref)
		if err != nil {
			return fmt.Errorf("failed to delete branch: %w", err)
		}
	}

	return nil
}

// ClosePullRequest 关闭 Pull Request
func (g *GitHub) ClosePullRequest(ctx context.Context, prID string) error {
	prNumber, err := parsePRNumber(prID)
	if err != nil {
		return err
	}

	state := "closed"
	_, _, err = g.client.PullRequests.Edit(ctx, g.owner, g.repo, prNumber, &github.PullRequest{
		State: &state,
	})
	if err != nil {
		return fmt.Errorf("failed to close PR: %w", err)
	}

	return nil
}

// GetPullRequestStatus 获取 PR 状态
func (g *GitHub) GetPullRequestStatus(ctx context.Context, prID string) (*pr.PullRequestStatus, error) {
	prNumber, err := parsePRNumber(prID)
	if err != nil {
		return nil, err
	}

	ghPR, _, err := g.client.PullRequests.Get(ctx, g.owner, g.repo, prNumber)
	if err != nil {
		return nil, fmt.Errorf("failed to get PR: %w", err)
	}

	status := &pr.PullRequestStatus{
		State:     ghPR.GetState(),
		Merged:    ghPR.GetMerged(),
		Mergeable: ghPR.Mergeable,
		UpdatedAt: ghPR.GetUpdatedAt().Time,
	}

	return status, nil
}

// ListPullRequests 列出 Pull Requests
func (g *GitHub) ListPullRequests(ctx context.Context, state string, limit int) ([]*pr.PullRequestInfo, error) {
	// 验证状态
	if state == "" {
		state = "open"
	}
	validStates := map[string]bool{
		"open":   true,
		"closed": true,
		"all":    true,
	}
	if !validStates[state] {
		state = "open"
	}

	// 设置分页
	opts := &github.PullRequestListOptions{
		State: state,
		ListOptions: github.ListOptions{
			PerPage: limit,
		},
	}

	ghPRs, _, err := g.client.PullRequests.List(ctx, g.owner, g.repo, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to list PRs: %w", err)
	}

	// 转换为统一格式
	result := make([]*pr.PullRequestInfo, 0, len(ghPRs))
	for _, ghPR := range ghPRs {
		info := &pr.PullRequestInfo{
			Number:    ghPR.GetNumber(),
			Title:     ghPR.GetTitle(),
			State:     ghPR.GetState(),
			HTMLURL:   ghPR.GetHTMLURL(),
			CreatedAt: ghPR.GetCreatedAt().Time,
			UpdatedAt: ghPR.GetUpdatedAt().Time,
			Author:    ghPR.GetUser().GetLogin(),
		}
		result = append(result, info)
	}

	return result, nil
}

// UpdatePullRequest 更新 Pull Request
func (g *GitHub) UpdatePullRequest(ctx context.Context, prID string, title, body *string, state *string) error {
	prNumber, err := parsePRNumber(prID)
	if err != nil {
		return err
	}

	update := &github.PullRequest{
		Title: title,
		Body:  body,
		State: state,
	}

	_, _, err = g.client.PullRequests.Edit(ctx, g.owner, g.repo, prNumber, update)
	if err != nil {
		return fmt.Errorf("failed to update PR: %w", err)
	}

	return nil
}

// AddComment 添加评论
func (g *GitHub) AddComment(ctx context.Context, prID string, body string) error {
	prNumber, err := parsePRNumber(prID)
	if err != nil {
		return err
	}

	comment := &github.IssueComment{
		Body: &body,
	}

	_, _, err = g.client.Issues.CreateComment(ctx, g.owner, g.repo, prNumber, comment)
	if err != nil {
		return fmt.Errorf("failed to add comment: %w", err)
	}

	return nil
}

// ApprovePullRequest 批准 Pull Request
func (g *GitHub) ApprovePullRequest(ctx context.Context, prID string) error {
	prNumber, err := parsePRNumber(prID)
	if err != nil {
		return err
	}

	review := &github.PullRequestReviewRequest{
		Event: github.String("APPROVE"),
	}

	_, _, err = g.client.PullRequests.CreateReview(ctx, g.owner, g.repo, prNumber, review)
	if err != nil {
		return fmt.Errorf("failed to approve PR: %w", err)
	}

	return nil
}

// parsePRNumber 解析 PR ID（支持数字、URL 等格式）
func parsePRNumber(prID string) (int, error) {
	// 使用 helpers 包解析
	prNumber, err := parsePRNumberInternal(prID)
	if err != nil {
		return 0, err
	}
	return prNumber, nil
}

// parsePRNumberInternal 内部解析函数（避免循环依赖）
func parsePRNumberInternal(prID string) (int, error) {
	// 如果是 URL，提取数字
	if strings.Contains(prID, "/pull/") {
		parts := strings.Split(prID, "/pull/")
		if len(parts) == 2 {
			prID = strings.TrimSuffix(parts[1], "/")
		}
	}

	// 尝试解析为数字
	var number int
	_, err := fmt.Sscanf(prID, "%d", &number)
	if err != nil {
		return 0, fmt.Errorf("invalid PR ID format: %s", prID)
	}

	return number, nil
}

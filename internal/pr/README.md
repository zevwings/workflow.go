# PR 平台提供者模块

本模块提供了统一的 Pull Request 操作接口，支持多种平台（GitHub、GitLab、Bitbucket 等）。

## 架构设计

### 平台提供者接口

`PlatformProvider` 接口定义了统一的 PR 操作接口，所有平台实现都需要实现此接口：

```go
type PlatformProvider interface {
    CreatePullRequest(ctx context.Context, title, body, sourceBranch string, targetBranch *string) (string, error)
    MergePullRequest(ctx context.Context, prID string, mergeMethod string, deleteBranch bool) error
    ClosePullRequest(ctx context.Context, prID string) error
    GetPullRequestStatus(ctx context.Context, prID string) (*PullRequestStatus, error)
    ListPullRequests(ctx context.Context, state string, limit int) ([]*PullRequestInfo, error)
    UpdatePullRequest(ctx context.Context, prID string, title, body *string, state *string) error
    AddComment(ctx context.Context, prID string, body string) error
    ApprovePullRequest(ctx context.Context, prID string) error
    GetPlatformName() string
}
```

### 目录结构

```
internal/pr/
├── platform.go          # 平台提供者接口定义
├── types.go             # 类型定义（PullRequestStatus, PullRequestInfo）
├── provider/
│   └── factory.go        # 工厂函数（创建平台提供者实例）
├── github/
│   ├── platform.go      # GitHub 平台实现
│   └── errors.go        # GitHub 错误处理
└── helpers/
    ├── resolution.go     # PR ID 解析辅助
    └── url.go            # URL 处理辅助
```

## 使用方法

### 1. 创建平台提供者

#### 自动检测平台

```go
import "github.com/zevwings/workflow/internal/pr/provider"

// 自动检测当前仓库的平台类型
platform, err := provider.NewPlatformProviderAuto()
if err != nil {
    log.Fatal(err)
}
```

#### 指定平台

```go
import "github.com/zevwings/workflow/internal/pr/provider"

// 创建 GitHub 平台提供者
platform, err := provider.NewPlatformProvider("github")
if err != nil {
    log.Fatal(err)
}
```

### 2. 创建 Pull Request

```go
ctx := context.Background()

// 创建 PR
prURL, err := platform.CreatePullRequest(
    ctx,
    "Fix bug in authentication",
    "This PR fixes a critical bug in the authentication flow.",
    "feature/fix-auth-bug",
    nil, // 使用默认分支
)
if err != nil {
    log.Fatal(err)
}

fmt.Printf("PR created: %s\n", prURL)
```

### 3. 合并 Pull Request

```go
ctx := context.Background()

// 合并 PR（使用 squash 方法，并删除源分支）
err := platform.MergePullRequest(
    ctx,
    "123",              // PR ID（支持数字、URL 等格式）
    "squash",           // 合并方法：merge, squash, rebase
    true,                // 删除源分支
)
if err != nil {
    log.Fatal(err)
}
```

### 4. 获取 PR 状态

```go
ctx := context.Background()

status, err := platform.GetPullRequestStatus(ctx, "123")
if err != nil {
    log.Fatal(err)
}

fmt.Printf("PR State: %s\n", status.State)
fmt.Printf("Merged: %v\n", status.Merged)
fmt.Printf("Mergeable: %v\n", status.Mergeable)
```

### 5. 列出 Pull Requests

```go
ctx := context.Background()

// 列出所有打开的 PR（最多 10 个）
prs, err := platform.ListPullRequests(ctx, "open", 10)
if err != nil {
    log.Fatal(err)
}

for _, pr := range prs {
    fmt.Printf("#%d: %s (%s)\n", pr.Number, pr.Title, pr.State)
}
```

### 6. 更新 Pull Request

```go
ctx := context.Background()

newTitle := "Updated: Fix bug in authentication"
err := platform.UpdatePullRequest(
    ctx,
    "123",
    &newTitle,  // 新标题
    nil,         // 不更新描述
    nil,         // 不更新状态
)
if err != nil {
    log.Fatal(err)
}
```

### 7. 添加评论

```go
ctx := context.Background()

err := platform.AddComment(
    ctx,
    "123",
    "This looks good! Please add tests.",
)
if err != nil {
    log.Fatal(err)
}
```

### 8. 批准 Pull Request

```go
ctx := context.Background()

err := platform.ApprovePullRequest(ctx, "123")
if err != nil {
    log.Fatal(err)
}
```

## 扩展支持新平台

要支持新平台（如 GitLab），需要：

1. 在 `internal/pr/` 下创建新目录（如 `gitlab/`）
2. 实现 `PlatformProvider` 接口
3. 在 `internal/pr/provider/factory.go` 中添加新平台的支持

示例：

```go
// internal/pr/gitlab/platform.go
package gitlab

import (
    "context"
    "github.com/zevwings/workflow/internal/pr"
)

type GitLab struct {
    // ... 实现细节
}

func (g *GitLab) CreatePullRequest(ctx context.Context, title, body, sourceBranch string, targetBranch *string) (string, error) {
    // 实现创建 PR 的逻辑
}

// ... 实现其他接口方法
```

然后在 `factory.go` 中添加：

```go
case "gitlab":
    return gitlab.NewGitLab()
```

## 配置要求

### GitHub

需要在 `$XDG_CONFIG_HOME/workflow/config.toml`（默认：`~/.config/workflow/config.toml`）中配置 GitHub token：

```toml
[github]
current = "default"

[[github.accounts]]
name = "default"
token = "ghp_xxxxxxxxxxxxxxxxxxxx"
```

## 错误处理

GitHub 平台提供了专门的错误处理函数：

```go
import "github.com/zevwings/workflow/internal/pr/github"

if github.IsNotFoundError(err) {
    // 处理 404 错误
}

if github.IsUnauthorizedError(err) {
    // 处理 401 错误
}

if github.IsForbiddenError(err) {
    // 处理 403 错误
}

// 格式化错误信息
formattedErr := github.FormatError(err)
```

## 辅助功能

### PR ID 解析

```go
import "github.com/zevwings/workflow/internal/pr/helpers"

// 支持多种格式：
// - "123"
// - "https://github.com/owner/repo/pull/123"
// - "owner/repo#123"
prNumber, err := helpers.ParsePRNumber("123")
if err != nil {
    log.Fatal(err)
}
```

### URL 处理

```go
import "github.com/zevwings/workflow/internal/pr/helpers"

// 从 URL 提取仓库信息
owner, repo, err := helpers.ExtractRepoFromURL("https://github.com/owner/repo.git")
if err != nil {
    log.Fatal(err)
}

// 构建 PR URL
prURL := helpers.BuildPRURL("https://github.com", owner, repo, 123)
```


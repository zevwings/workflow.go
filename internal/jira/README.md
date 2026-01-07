# Jira 客户端封装

本模块提供了对 `go-jira` SDK 的封装，简化 Jira API 的使用。

## 快速开始

### 1. 配置 Jira

在 `~/.workflow/config.toml` 中配置 Jira 信息：

```toml
[jira]
url = "https://your-domain.atlassian.net"
username = "your-email@example.com"
token = "your-api-token"
```

### 2. 使用 JiraClient（推荐）

`JiraClient` 提供了高级封装，简化常用操作：

```go
package main

import (
    "fmt"
    "github.com/zevwings/workflow/internal/jira"
)

func main() {
    // 创建客户端
    client, err := jira.NewJiraClient()
    if err != nil {
        panic(err)
    }

    // 获取 Ticket 信息
    issue, err := client.GetTicketInfo("PROJ-123")
    if err != nil {
        panic(err)
    }
    fmt.Printf("Issue: %s - %s\n", issue.Key, issue.Fields.Summary)

    // 获取附件列表
    attachments, err := client.GetAttachments("PROJ-123")
    if err != nil {
        panic(err)
    }
    fmt.Printf("Attachments: %d\n", len(attachments))

    // 添加评论
    err = client.AddComment("PROJ-123", "This is a comment")
    if err != nil {
        panic(err)
    }

    // 更新状态
    err = client.MoveTicket("PROJ-123", "In Progress")
    if err != nil {
        panic(err)
    }
}
```

### 3. 使用底层 API（高级用法）

如果需要更细粒度的控制，可以直接使用 API 模块：

```go
package main

import (
    "context"
    "github.com/zevwings/workflow/internal/jira"
    "github.com/zevwings/workflow/internal/jira/api"
)

func main() {
    // 创建底层客户端
    client, err := jira.NewClient()
    if err != nil {
        panic(err)
    }

    // 使用自定义 context
    ctx := context.WithTimeout(context.Background(), 10*time.Second)
    client = client.WithContext(ctx)

    // 创建 API 实例
    issueAPI := api.NewIssueAPI(client.GetJiraClient(), client.GetContext())

    // 使用 API
    issue, err := issueAPI.GetIssue("PROJ-123")
    if err != nil {
        panic(err)
    }
}
```

## API 参考

### JiraClient 方法

- `GetUserInfo()` - 获取当前用户信息
- `GetTicketInfo(ticket)` - 获取 Ticket 信息
- `GetAttachments(ticket)` - 获取附件列表
- `GetComments(ticket)` - 获取评论列表
- `AddComment(ticket, comment)` - 添加评论
- `MoveTicket(ticket, status)` - 更新状态（通过状态名称）
- `AssignTicket(ticket, accountID)` - 分配 Ticket
- `UploadAttachment(ticket, filePath)` - 上传附件
- `GetTransitions(ticket)` - 获取可用的状态转换
- `GetChangelog(ticket)` - 获取变更历史
- `GetProject(projectKey)` - 获取项目信息
- `GetProjectStatuses(projectKey)` - 获取项目状态列表
- `FindUsers(query)` - 搜索用户

### IssueAPI 方法

- `GetIssue(ticket)` - 获取 Issue 信息
- `GetIssueAttachments(ticket)` - 获取附件列表
- `GetIssueTransitions(ticket)` - 获取可用的状态转换
- `TransitionIssue(ticket, transitionID)` - 更新状态（通过转换 ID）
- `AssignIssue(ticket, accountID)` - 分配 Issue
- `AddComment(ticket, comment)` - 添加评论
- `GetComments(ticket)` - 获取评论列表
- `UploadAttachment(ticket, filePath)` - 上传附件
- `DownloadAttachment(attachment)` - 下载附件
- `GetChangelog(ticket)` - 获取变更历史

### ProjectAPI 方法

- `GetProject(projectKey)` - 获取项目信息
- `GetProjectStatuses(projectKey)` - 获取项目状态列表（注意：Jira Cloud API v2 不直接支持）
- `ListProjects()` - 列出所有项目

### UserAPI 方法

- `GetCurrentUser()` - 获取当前用户信息
- `GetUser(accountID)` - 根据 Account ID 获取用户信息
- `FindUsers(query)` - 搜索用户

## 辅助函数

- `ValidateTicketKey(ticket)` - 验证 Ticket Key 格式
- `NormalizeTicketKey(ticket)` - 规范化 Ticket Key（转大写）
- `ExtractProjectKey(ticket)` - 从 Ticket Key 中提取项目 Key
- `ExtractTicketNumber(ticket)` - 从 Ticket Key 中提取 Ticket 编号

## 注意事项

1. **认证方式**：使用 Basic Auth（Email + API Token）
2. **Ticket Key 格式**：必须是 `PROJECT-NUMBER` 格式（如 "PROJ-123"）
3. **错误处理**：所有方法都会返回详细的错误信息
4. **Context 支持**：底层客户端支持自定义 context，用于超时控制等

## 依赖

- `github.com/andygrunwald/go-jira/v2/cloud` - Jira SDK


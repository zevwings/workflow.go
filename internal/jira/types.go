package jira

// JiraIssue Issue 完整信息
type JiraIssue struct {
	Key     string          `json:"key"`
	ID      string          `json:"id"`
	SelfURL string          `json:"self"`
	Fields  JiraIssueFields `json:"fields"`
}

// JiraIssueFields Issue 字段
type JiraIssueFields struct {
	Summary      string            `json:"summary"`
	Description  *string           `json:"description,omitempty"`
	Status       JiraStatus        `json:"status"`
	Attachment   []*JiraAttachment `json:"attachment,omitempty"`
	Comment      *JiraComments     `json:"comment,omitempty"`
	Priority     *JiraPriority     `json:"priority,omitempty"`
	Created      *string           `json:"created,omitempty"`
	Updated      *string           `json:"updated,omitempty"`
	Reporter     *JiraUser         `json:"reporter,omitempty"`
	Assignee     *JiraUser         `json:"assignee,omitempty"`
	Labels       []string          `json:"labels,omitempty"`
	Components   []*JiraComponent  `json:"components,omitempty"`
	FixVersions  []*JiraVersion    `json:"fixVersions,omitempty"`
	IssueLinks   []*JiraIssueLink  `json:"issuelinks,omitempty"`
	Subtasks     []*JiraSubtask    `json:"subtasks,omitempty"`
	TimeTracking *JiraTimeTracking `json:"timeTracking,omitempty"`
}

// JiraAttachment 附件信息
type JiraAttachment struct {
	Filename   string  `json:"filename"`
	ContentURL string  `json:"content"`
	MimeType   *string `json:"mimeType,omitempty"`
	Size       *int64  `json:"size,omitempty"`
}

// JiraComments 评论容器
type JiraComments struct {
	Comments   []*JiraComment `json:"comments"`
	MaxResults *int64         `json:"maxResults,omitempty"`
	StartAt    *int64         `json:"startAt,omitempty"`
	Total      *int64         `json:"total,omitempty"`
}

// JiraComment 评论信息
type JiraComment struct {
	ID           string    `json:"id"`
	Body         string    `json:"body"`
	Created      string    `json:"created"`
	Updated      *string   `json:"updated,omitempty"`
	Author       *JiraUser `json:"author,omitempty"`
	UpdateAuthor *JiraUser `json:"updateAuthor,omitempty"`
}

// JiraStatus 状态信息
type JiraStatus struct {
	ID      string  `json:"id"`
	Name    string  `json:"name"`
	SelfURL *string `json:"self,omitempty"`
}

// JiraUser 用户信息
type JiraUser struct {
	AccountID    string  `json:"accountId"`
	DisplayName  string  `json:"displayName"`
	EmailAddress *string `json:"emailAddress,omitempty"`
}

// JiraTransition Transition 信息
type JiraTransition struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// JiraPriority 优先级信息
type JiraPriority struct {
	ID      string  `json:"id"`
	Name    string  `json:"name"`
	IconURL *string `json:"iconUrl,omitempty"`
}

// JiraComponent 组件信息
type JiraComponent struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Description *string `json:"description,omitempty"`
}

// JiraVersion 版本信息
type JiraVersion struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Released    bool    `json:"released"`
	ReleaseDate *string `json:"releaseDate,omitempty"`
}

// JiraIssueLink Issue 链接
type JiraIssueLink struct {
	ID           string            `json:"id"`
	Type         JiraIssueLinkType `json:"type"`
	InwardIssue  *JiraIssueRef     `json:"inwardIssue,omitempty"`
	OutwardIssue *JiraIssueRef     `json:"outwardIssue,omitempty"`
}

// JiraIssueLinkType Issue 链接类型
type JiraIssueLinkType struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Inward  string `json:"inward"`
	Outward string `json:"outward"`
}

// JiraIssueRef Issue 引用
type JiraIssueRef struct {
	Key    string `json:"key"`
	Fields struct {
		Summary string     `json:"summary"`
		Status  JiraStatus `json:"status"`
	} `json:"fields"`
}

// JiraSubtask 子任务
type JiraSubtask struct {
	ID     string            `json:"id"`
	Key    string            `json:"key"`
	Fields JiraSubtaskFields `json:"fields"`
}

// JiraSubtaskFields 子任务字段
type JiraSubtaskFields struct {
	Summary string     `json:"summary"`
	Status  JiraStatus `json:"status"`
}

// JiraTimeTracking 时间跟踪
type JiraTimeTracking struct {
	OriginalEstimate  *string `json:"originalEstimate,omitempty"`
	RemainingEstimate *string `json:"remainingEstimate,omitempty"`
	TimeSpent         *string `json:"timeSpent,omitempty"`
}

// JiraChangelog 变更日志
type JiraChangelog struct {
	ID      string               `json:"id"`
	Items   []*JiraChangelogItem `json:"items"`
	Author  *JiraUser            `json:"author,omitempty"`
	Created string               `json:"created"`
}

// JiraChangelogItem 变更日志项
type JiraChangelogItem struct {
	Field      string  `json:"field"`
	FieldType  string  `json:"fieldtype"`
	From       *string `json:"from,omitempty"`
	FromString *string `json:"fromString,omitempty"`
	To         *string `json:"to,omitempty"`
	ToString   *string `json:"toString,omitempty"`
}

// JiraChangelogHistory 变更日志历史
type JiraChangelogHistory struct {
	Histories  []*JiraChangelog `json:"histories"`
	MaxResults int              `json:"maxResults"`
	StartAt    int              `json:"startAt"`
	Total      int              `json:"total"`
}

// JiraProject 项目信息
type JiraProject struct {
	ID          string  `json:"id"`
	Key         string  `json:"key"`
	Name        string  `json:"name"`
	Description *string `json:"description,omitempty"`
	AvatarURLs  *struct {
		Small  *string `json:"48x48,omitempty"`
		Medium *string `json:"24x24,omitempty"`
		Large  *string `json:"16x16,omitempty"`
	} `json:"avatarUrls,omitempty"`
}

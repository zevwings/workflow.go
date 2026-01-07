package pr

import "time"

// PullRequestStatus PR 状态信息
type PullRequestStatus struct {
	State     string    // 状态（"open", "closed", "merged"）
	Merged    bool      // 是否已合并
	Mergeable *bool     // 是否可合并（nil 表示未知）
	UpdatedAt time.Time // 更新时间
}

// PullRequestInfo PR 信息
type PullRequestInfo struct {
	Number    int       // PR 编号
	Title     string    // 标题
	State     string    // 状态
	HTMLURL   string    // PR URL
	CreatedAt time.Time // 创建时间
	UpdatedAt time.Time // 更新时间
	Author    string    // 作者
}


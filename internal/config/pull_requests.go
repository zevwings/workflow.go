package config

// PullRequestsConfig PR 配置（个人偏好）
//
// 仓库级别的个人偏好配置，不提交到 Git。
type PullRequestsConfig struct {
	// AutoAcceptChangeType 自动接受变更类型选择（个人偏好）
	AutoAcceptChangeType *bool `toml:"auto_accept_change_type,omitempty"`
}

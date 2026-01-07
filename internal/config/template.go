package config

// TemplateConfig 模板配置
//
// 项目级别的模板配置，提交到 Git。
// 包含 commit、branch、pull_requests 的模板配置。
type TemplateConfig struct {
	// Commit 提交消息模板配置
	Commit map[string]interface{} `toml:"commit,omitempty"`
	// Branch 分支命名模板配置
	Branch map[string]interface{} `toml:"branch,omitempty"`
	// PullRequests PR 模板配置
	PullRequests map[string]interface{} `toml:"pull_requests,omitempty"`
}

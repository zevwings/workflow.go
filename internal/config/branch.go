package config

// BranchConfig 分支配置（个人偏好）
//
// 仓库级别的个人偏好配置，不提交到 Git。
type BranchConfig struct {
	// Prefix 分支前缀（个人偏好）
	Prefix *string `toml:"prefix,omitempty"`
	// Ignore 忽略的分支列表（个人偏好）
	Ignore []string `toml:"ignore,omitempty"`
}

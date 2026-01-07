package config

// GlobalConfig 全局配置结构
//
// 统一配置结构，包含所有子配置模块。
// 用于全局配置：~/.workflow/config.toml
type GlobalConfig struct {
	User   UserConfig   `toml:"user,omitempty"`
	Jira   JiraConfig   `toml:"jira,omitempty"`
	GitHub GitHubConfig `toml:"github,omitempty"`
	Log    LogConfig    `toml:"log,omitempty"`
	LLM    LLMConfig    `toml:"llm,omitempty"`
	Proxy  ProxyConfig  `toml:"proxy,omitempty"`
}

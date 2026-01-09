package config

// GlobalConfig 全局配置结构
//
// 统一配置结构，包含所有子配置模块。
// 用于全局配置：遵循 XDG Base Directory Specification
// 配置文件位置：$XDG_CONFIG_HOME/Workflow/config.toml（默认：~/.config/Workflow/config.toml）
type GlobalConfig struct {
	Jira   JiraConfig   `toml:"jira,omitempty"`
	GitHub GitHubConfig `toml:"github,omitempty"`
	Log    LogConfig    `toml:"log,omitempty"`
	LLM    LLMConfig    `toml:"llm,omitempty"`
	Proxy  ProxyConfig  `toml:"proxy,omitempty"`
}

// RepoConfig 仓库配置结构
//
// 统一配置结构，包含所有仓库级别的公共配置模块。
// 用于仓库公共配置：.workflow/config.toml（项目根目录，提交到 Git）
type RepoConfig struct {
	Template TemplateConfig `toml:"template,omitempty"`
}

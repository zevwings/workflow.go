package config

// JiraConfig Jira 配置
type JiraConfig struct {
	URL      string `toml:"url,omitempty"`
	Username string `toml:"username,omitempty"`
	Token    string `toml:"token,omitempty"`
}


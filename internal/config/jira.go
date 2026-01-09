package config

// JiraConfig Jira 配置
type JiraConfig struct {
	Email          string `toml:"email,omitempty"`
	APIToken       string `toml:"api_token,omitempty"`
	ServiceAddress string `toml:"service_address,omitempty"`
}


package config

// GitHubConfig GitHub 配置
type GitHubConfig struct {
	Accounts []GitHubAccount `toml:"accounts,omitempty"`
	Current  string          `toml:"current,omitempty"`
}

// GitHubAccount GitHub 账号
type GitHubAccount struct {
	Name  string `toml:"name"`
	Token string `toml:"token"`
}


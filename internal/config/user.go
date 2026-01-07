package config

// UserConfig 用户配置
type UserConfig struct {
	Name  string `toml:"name,omitempty"`
	Email string `toml:"email,omitempty"`
}


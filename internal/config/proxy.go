package config

// ProxyConfig 代理配置
type ProxyConfig struct {
	Enabled bool   `toml:"enabled,omitempty"`
	HTTP    string `toml:"http,omitempty"`
	HTTPS   string `toml:"https,omitempty"`
}

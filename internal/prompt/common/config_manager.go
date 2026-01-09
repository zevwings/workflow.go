package common

// ConfigManager 配置管理器
// 提供统一的配置管理接口，支持默认配置、全局配置和局部配置的层次结构
type ConfigManager struct {
	// defaultConfig 默认配置（系统默认值）
	defaultConfig PromptConfig
	// globalConfig 全局配置（用户设置的全局配置）
	globalConfig PromptConfig
}

// NewConfigManager 创建配置管理器
func NewConfigManager(defaultConfig PromptConfig) *ConfigManager {
	return &ConfigManager{
		defaultConfig: defaultConfig,
		globalConfig:  PromptConfig{},
	}
}

// SetGlobalConfig 设置全局配置
// 全局配置会与默认配置合并，非 nil 字段会覆盖默认配置
func (m *ConfigManager) SetGlobalConfig(config PromptConfig) {
	m.globalConfig = config
}

// GetGlobalConfig 获取全局配置
func (m *ConfigManager) GetGlobalConfig() PromptConfig {
	return m.globalConfig
}

// BuildConfig 构建最终配置
// 按照优先级合并：defaultConfig < globalConfig < localConfig
// 返回合并后的配置
func (m *ConfigManager) BuildConfig(localConfig *PromptConfig) PromptConfig {
	// 首先合并默认配置和全局配置
	merged := FillDefaults(m.globalConfig, m.defaultConfig)

	// 如果有局部配置，继续合并
	if localConfig != nil {
		merged = MergeConfig(&merged, localConfig)
	}

	return merged
}

// ResetGlobalConfig 重置全局配置为空配置
func (m *ConfigManager) ResetGlobalConfig() {
	m.globalConfig = PromptConfig{}
}

// GetDefaultConfig 获取默认配置（只读）
func (m *ConfigManager) GetDefaultConfig() PromptConfig {
	return m.defaultConfig
}

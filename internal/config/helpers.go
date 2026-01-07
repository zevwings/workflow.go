package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/pelletier/go-toml/v2"
)

// SaveConfigToFile 保存配置到文件
//
// 这是一个通用的配置保存辅助函数，用于消除配置管理器中的重复代码。
// 它会自动创建目录（如果不存在），序列化配置为 TOML 格式并写入文件。
//
// 参数:
//   - path: 配置文件路径
//   - config: 要保存的配置对象（可以是任意可序列化的类型）
//
// 返回:
//   - error: 如果保存失败，返回错误
func SaveConfigToFile(path string, config interface{}) error {
	// 序列化配置为 TOML
	data, err := toml.Marshal(config)
	if err != nil {
		return fmt.Errorf("序列化配置失败: %w", err)
	}

	// 确保目录存在
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return fmt.Errorf("创建配置目录失败: %w", err)
	}

	// 写入文件
	if err := os.WriteFile(path, data, 0644); err != nil {
		return fmt.Errorf("写入配置文件失败: %w", err)
	}

	return nil
}

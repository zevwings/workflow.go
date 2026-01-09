package config

import (
	"path/filepath"

	"github.com/adrg/xdg"
)

// ConfigDir 获取配置目录
//
// 返回 XDG 配置目录下的 workflow 子目录路径。
// 遵循 XDG Base Directory Specification：
//   - Unix: ~/.config/workflow
//   - Windows: %APPDATA%\workflow
//   - macOS: ~/Library/Application Support/workflow
//
// 返回:
//   - string: 配置目录路径
//   - error: 如果获取失败，返回错误
func ConfigDir() (string, error) {
	configDir := xdg.ConfigHome
	return filepath.Join(configDir, "workflow"), nil
}

// DataDir 获取数据目录
//
// 返回 XDG 数据目录下的 workflow 子目录路径。
// 遵循 XDG Base Directory Specification：
//   - Unix: ~/.local/share/workflow
//   - Windows: %APPDATA%\workflow
//   - macOS: ~/Library/Application Support/workflow
//
// 返回:
//   - string: 数据目录路径
//   - error: 如果获取失败，返回错误
func DataDir() (string, error) {
	dataDir := xdg.DataHome
	return filepath.Join(dataDir, "workflow"), nil
}

// StateDir 获取状态目录
//
// 返回 XDG 状态目录下的 workflow 子目录路径。
// 遵循 XDG Base Directory Specification：
//   - Unix: ~/.local/state/workflow
//   - Windows: %LOCALAPPDATA%\workflow
//   - macOS: ~/Library/Application Support/workflow
//
// 返回:
//   - string: 状态目录路径
//   - error: 如果获取失败，返回错误
func StateDir() (string, error) {
	stateDir := xdg.StateHome
	return filepath.Join(stateDir, "workflow"), nil
}

// CacheDir 获取缓存目录
//
// 返回 XDG 缓存目录下的 workflow 子目录路径。
// 遵循 XDG Base Directory Specification：
//   - Unix: ~/.cache/workflow
//   - Windows: %LOCALAPPDATA%\workflow
//   - macOS: ~/Library/Caches/workflow
//
// 返回:
//   - string: 缓存目录路径
//   - error: 如果获取失败，返回错误
func CacheDir() (string, error) {
	cacheDir := xdg.CacheHome
	return filepath.Join(cacheDir, "workflow"), nil
}

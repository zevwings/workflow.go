package config

import (
	"os"
	"path/filepath"
	"runtime"

	"github.com/adrg/xdg"
)

// ConfigDir 获取配置目录
//
// 返回 XDG 配置目录下的 Workflow 子目录路径。
// 遵循 XDG Base Directory Specification：
//   - Unix: ~/.config/Workflow
//   - Windows: %APPDATA%\Workflow
//   - macOS: ~/Library/Application Support/Workflow
//
// 在 macOS 上，默认使用 iCloud Drive 路径：
//   - macOS (iCloud): ~/Library/Mobile Documents/com~apple~CloudDocs/Workflow
//
// 可以通过环境变量 WORKFLOW_DISABLE_ICLOUD_CONFIG=1 禁用 iCloud 支持。
//
// 返回:
//   - string: 配置目录路径
//   - error: 如果获取失败，返回错误
func ConfigDir() (string, error) {
	// 检查是否禁用 iCloud（环境变量优先级最高）
	disableICloud := os.Getenv("WORKFLOW_DISABLE_ICLOUD_CONFIG") == "1"

	// macOS 上默认使用 iCloud Drive（除非被禁用）
	if runtime.GOOS == "darwin" && !disableICloud {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}
		iCloudPath := filepath.Join(homeDir, "Library", "Mobile Documents", "com~apple~CloudDocs", "Workflow")
		return iCloudPath, nil
	}

	// 默认使用 XDG 配置目录
	configDir := xdg.ConfigHome
	return filepath.Join(configDir, "Workflow"), nil
}

// DataDir 获取数据目录
//
// 返回 XDG 数据目录下的 Workflow 子目录路径。
// 遵循 XDG Base Directory Specification：
//   - Unix: ~/.local/share/Workflow
//   - Windows: %APPDATA%\Workflow
//   - macOS: ~/Library/Application Support/Workflow
//
// 返回:
//   - string: 数据目录路径
//   - error: 如果获取失败，返回错误
func DataDir() (string, error) {
	dataDir := xdg.DataHome
	return filepath.Join(dataDir, "Workflow"), nil
}

// StateDir 获取状态目录
//
// 返回 XDG 状态目录下的 Workflow 子目录路径。
// 遵循 XDG Base Directory Specification：
//   - Unix: ~/.local/state/Workflow
//   - Windows: %LOCALAPPDATA%\Workflow
//   - macOS: ~/Library/Application Support/Workflow
//
// 返回:
//   - string: 状态目录路径
//   - error: 如果获取失败，返回错误
func StateDir() (string, error) {
	stateDir := xdg.StateHome
	return filepath.Join(stateDir, "Workflow"), nil
}

// CacheDir 获取缓存目录
//
// 返回 XDG 缓存目录下的 Workflow 子目录路径。
// 遵循 XDG Base Directory Specification：
//   - Unix: ~/.cache/Workflow
//   - Windows: %LOCALAPPDATA%\Workflow
//   - macOS: ~/Library/Caches/Workflow
//
// 返回:
//   - string: 缓存目录路径
//   - error: 如果获取失败，返回错误
func CacheDir() (string, error) {
	cacheDir := xdg.CacheHome
	return filepath.Join(cacheDir, "Workflow"), nil
}

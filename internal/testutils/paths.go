//go:build test

package testutils

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/stretchr/testify/require"
)

// TestHomeDir 获取测试主目录
//
// 优先使用环境变量（HOME 或 USERPROFILE），支持测试环境隔离。
// 如果环境变量未设置，则使用系统默认主目录。
//
// 参数:
//   - t: 测试对象
//
// 返回:
//   - string: 主目录路径
func TestHomeDir(t *testing.T) string {
	// 优先使用环境变量（支持测试隔离）
	if home := os.Getenv("HOME"); home != "" {
		return home
	}
	if runtime.GOOS == "windows" {
		if home := os.Getenv("USERPROFILE"); home != "" {
			return home
		}
	}

	// 回退到系统默认主目录
	homeDir, err := os.UserHomeDir()
	require.NoError(t, err, "获取主目录失败")
	return homeDir
}

// TestConfigDir 获取测试配置目录
//
// 优先使用环境变量（XDG_CONFIG_HOME 或 APPDATA），支持测试环境隔离。
// 如果环境变量未设置，则使用默认配置目录：
//   - Unix: ~/.config
//   - Windows: %APPDATA%
//
// 参数:
//   - t: 测试对象
//
// 返回:
//   - string: 配置目录路径（已创建）
func TestConfigDir(t *testing.T) string {
	var configDir string

	// 优先使用环境变量（支持测试隔离）
	if xdgConfig := os.Getenv("XDG_CONFIG_HOME"); xdgConfig != "" {
		configDir = xdgConfig
	} else if runtime.GOOS == "windows" {
		if appData := os.Getenv("APPDATA"); appData != "" {
			configDir = appData
		} else {
			homeDir := TestHomeDir(t)
			configDir = filepath.Join(homeDir, "AppData", "Roaming")
		}
	} else {
		// Unix 默认：~/.config
		homeDir := TestHomeDir(t)
		configDir = filepath.Join(homeDir, ".config")
	}

	// 确保目录存在
	err := os.MkdirAll(configDir, 0755)
	require.NoError(t, err, "创建配置目录失败: %s", configDir)

	return configDir
}

// TestDataDir 获取测试数据目录
//
// 优先使用环境变量（XDG_DATA_HOME），支持测试环境隔离。
// 如果环境变量未设置，则使用默认数据目录：
//   - Unix: ~/.local/share
//   - Windows: %APPDATA%
//
// 参数:
//   - t: 测试对象
//
// 返回:
//   - string: 数据目录路径（已创建）
func TestDataDir(t *testing.T) string {
	var dataDir string

	// 优先使用环境变量（支持测试隔离）
	if xdgData := os.Getenv("XDG_DATA_HOME"); xdgData != "" {
		dataDir = xdgData
	} else if runtime.GOOS == "windows" {
		if appData := os.Getenv("APPDATA"); appData != "" {
			dataDir = appData
		} else {
			homeDir := TestHomeDir(t)
			dataDir = filepath.Join(homeDir, "AppData", "Roaming")
		}
	} else {
		// Unix 默认：~/.local/share
		homeDir := TestHomeDir(t)
		dataDir = filepath.Join(homeDir, ".local", "share")
	}

	// 确保目录存在
	err := os.MkdirAll(dataDir, 0755)
	require.NoError(t, err, "创建数据目录失败: %s", dataDir)

	return dataDir
}

// TestCacheDir 获取测试缓存目录
//
// 优先使用环境变量（XDG_CACHE_HOME），支持测试环境隔离。
// 如果环境变量未设置，则使用默认缓存目录：
//   - Unix: ~/.cache
//   - Windows: %LOCALAPPDATA%
//
// 参数:
//   - t: 测试对象
//
// 返回:
//   - string: 缓存目录路径（已创建）
func TestCacheDir(t *testing.T) string {
	var cacheDir string

	// 优先使用环境变量（支持测试隔离）
	if xdgCache := os.Getenv("XDG_CACHE_HOME"); xdgCache != "" {
		cacheDir = xdgCache
	} else if runtime.GOOS == "windows" {
		if localAppData := os.Getenv("LOCALAPPDATA"); localAppData != "" {
			cacheDir = localAppData
		} else {
			homeDir := TestHomeDir(t)
			cacheDir = filepath.Join(homeDir, "AppData", "Local")
		}
	} else {
		// Unix 默认：~/.cache
		homeDir := TestHomeDir(t)
		cacheDir = filepath.Join(homeDir, ".cache")
	}

	// 确保目录存在
	err := os.MkdirAll(cacheDir, 0755)
	require.NoError(t, err, "创建缓存目录失败: %s", cacheDir)

	return cacheDir
}

// TestWorkflowConfigDir 获取测试 Workflow 配置目录
//
// 返回 XDG 配置目录下的 workflow 子目录路径（已创建）。
// 遵循 XDG Base Directory Specification，使用 TestConfigDir 获取基础配置目录。
// 支持测试环境隔离（通过环境变量）。
//
// 参数:
//   - t: 测试对象
//
// 返回:
//   - string: Workflow 配置目录路径（已创建）
func TestWorkflowConfigDir(t *testing.T) string {
	configDir := TestConfigDir(t)
	workflowConfigDir := filepath.Join(configDir, "workflow")

	// 确保目录存在
	err := os.MkdirAll(workflowConfigDir, 0755)
	require.NoError(t, err, "创建 Workflow 配置目录失败: %s", workflowConfigDir)

	return workflowConfigDir
}

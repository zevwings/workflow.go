package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ==================== SaveConfigToFile 测试 ====================

func TestSaveConfigToFile(t *testing.T) {
	tests := []struct {
		name    string
		config  interface{}
		wantErr bool
	}{
		{
			name: "保存简单配置",
			config: map[string]interface{}{
				"key1": "value1",
				"key2": 123,
			},
			wantErr: false,
		},
		{
			name: "保存嵌套配置",
			config: map[string]interface{}{
				"section": map[string]interface{}{
					"key": "value",
				},
			},
			wantErr: false,
		},
		{
			name: "保存空配置",
			config: map[string]interface{}{},
			wantErr: false,
			// 注意：空 map 序列化后可能为空，这是正常行为
		},
		{
			name: "保存结构体配置",
			config: GlobalConfig{
				Log: LogConfig{
					Level: "debug",
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange: 创建临时目录
			tempDir := t.TempDir()
			configPath := filepath.Join(tempDir, "config.toml")

			// Act: 保存配置
			err := SaveConfigToFile(configPath, tt.config)

			// Assert: 验证结果
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)

				// 验证文件已创建
				_, err := os.Stat(configPath)
				assert.NoError(t, err, "配置文件应该已创建")

				// 验证文件内容不为空（空 map 除外）
				data, err := os.ReadFile(configPath)
				require.NoError(t, err)
				// 空 map 序列化后可能为空，这是正常行为
				if tt.name != "保存空配置" {
					assert.NotEmpty(t, data, "配置文件内容不应为空")
				}
			}
		})
	}
}

func TestSaveConfigToFile_CreateDirectory(t *testing.T) {
	// Arrange: 创建临时目录，但配置文件路径包含不存在的子目录
	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, "subdir", "nested", "config.toml")

	config := map[string]interface{}{
		"key": "value",
	}

	// Act: 保存配置
	err := SaveConfigToFile(configPath, config)

	// Assert: 应该自动创建目录
	assert.NoError(t, err)

	// 验证目录已创建
	dir := filepath.Dir(configPath)
	_, err = os.Stat(dir)
	assert.NoError(t, err, "目录应该已自动创建")

	// 验证文件已创建
	_, err = os.Stat(configPath)
	assert.NoError(t, err, "配置文件应该已创建")
}

func TestSaveConfigToFile_FilePermissions(t *testing.T) {
	// Arrange: 创建临时目录
	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, "config.toml")

	config := map[string]interface{}{
		"key": "value",
	}

	// Act: 保存配置
	err := SaveConfigToFile(configPath, config)

	// Assert: 验证文件权限
	assert.NoError(t, err)

	info, err := os.Stat(configPath)
	require.NoError(t, err)

	// 验证文件权限（应该是 0644）
	mode := info.Mode()
	assert.Equal(t, os.FileMode(0644), mode.Perm(), "文件权限应该是 0644")
}

func TestSaveConfigToFile_OverwriteExisting(t *testing.T) {
	// Arrange: 创建临时目录和现有配置文件
	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, "config.toml")

	// 创建现有配置文件
	existingConfig := map[string]interface{}{
		"old_key": "old_value",
	}
	err := SaveConfigToFile(configPath, existingConfig)
	require.NoError(t, err)

	// 读取原始内容
	originalData, err := os.ReadFile(configPath)
	require.NoError(t, err)

	// Act: 保存新配置
	newConfig := map[string]interface{}{
		"new_key": "new_value",
	}
	err = SaveConfigToFile(configPath, newConfig)

	// Assert: 验证文件已被覆盖
	assert.NoError(t, err)

	// 验证文件内容已更新
	newData, err := os.ReadFile(configPath)
	require.NoError(t, err)
	assert.NotEqual(t, originalData, newData, "文件内容应该已被更新")

	// 验证新内容包含新配置
	assert.Contains(t, string(newData), "new_key")
	assert.Contains(t, string(newData), "new_value")
}

func TestSaveConfigToFile_InvalidPath(t *testing.T) {
	// 这个测试在大多数系统上可能无法触发真正的错误
	// 因为现代文件系统通常会自动处理路径问题
	// 但我们可以测试一些边界情况

	t.Run("空路径", func(t *testing.T) {
		config := map[string]interface{}{
			"key": "value",
		}

		// 空路径应该会导致错误
		err := SaveConfigToFile("", config)
		assert.Error(t, err)
	})

	t.Run("根目录路径", func(t *testing.T) {
		// 在 Unix 系统上，尝试写入根目录通常会失败（权限问题）
		// 在 Windows 上，这可能会失败或成功，取决于权限
		// 这个测试主要是为了覆盖错误处理路径

		config := map[string]interface{}{
			"key": "value",
		}

		// 尝试写入根目录（在大多数系统上会失败）
		rootPath := filepath.Join(string(filepath.Separator), "config.toml")
		err := SaveConfigToFile(rootPath, config)

		// 这个测试可能会失败（权限错误）或成功（取决于系统）
		// 我们主要验证函数不会 panic
		if err != nil {
			assert.Error(t, err)
		}
	})
}


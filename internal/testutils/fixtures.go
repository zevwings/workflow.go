//go:build test

package testutils

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

// LoadFixture 加载测试数据文件（二进制）
//
// 从 testdata/fixtures/ 目录加载测试数据文件。
//
// 参数:
//   - t: 测试对象
//   - filename: 文件名（相对于 testdata/fixtures/）
//
// 返回:
//   - []byte: 文件内容
func LoadFixture(t *testing.T, filename string) []byte {
	path := getFixturePath(t, filename)
	data, err := os.ReadFile(path)
	require.NoError(t, err, "读取测试数据文件失败: %s", path)
	return data
}

// LoadTextFixture 加载测试文本文件
//
// 从 testdata/fixtures/ 目录加载文本文件。
//
// 参数:
//   - t: 测试对象
//   - filename: 文件名（相对于 testdata/fixtures/）
//
// 返回:
//   - string: 文件内容（文本）
func LoadTextFixture(t *testing.T, filename string) string {
	data := LoadFixture(t, filename)
	return string(data)
}

// LoadBinaryFixture 加载测试二进制文件
//
// 从 testdata/fixtures/ 目录加载二进制文件。
// 这是 LoadFixture 的别名，用于语义清晰。
//
// 参数:
//   - t: 测试对象
//   - filename: 文件名（相对于 testdata/fixtures/）
//
// 返回:
//   - []byte: 文件内容（二进制）
func LoadBinaryFixture(t *testing.T, filename string) []byte {
	return LoadFixture(t, filename)
}

// getFixturePath 获取测试数据文件路径
//
// 从项目根目录的 testdata/fixtures/ 目录查找文件。
// 如果当前目录不在项目根目录，会向上查找。
//
// 参数:
//   - t: 测试对象
//   - filename: 文件名
//
// 返回:
//   - string: 文件路径
func getFixturePath(t *testing.T, filename string) string {
	// 获取当前工作目录
	wd, err := os.Getwd()
	require.NoError(t, err, "获取当前工作目录失败")

	// 从当前目录向上查找 testdata 目录
	dir := wd
	for {
		testdataPath := filepath.Join(dir, "testdata", "fixtures", filename)
		if _, err := os.Stat(testdataPath); err == nil {
			return testdataPath
		}

		// 向上查找
		parent := filepath.Dir(dir)
		if parent == dir {
			// 已到达根目录，停止查找
			break
		}
		dir = parent
	}

	// 如果找不到，尝试直接使用相对路径
	testdataPath := filepath.Join("testdata", "fixtures", filename)
	require.FileExists(t, testdataPath, "测试数据文件不存在: %s", testdataPath)
	return testdataPath
}

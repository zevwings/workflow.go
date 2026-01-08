//go:build test

package testutils

import (
	"testing"
)

// SetupTestEnv 设置测试环境
// 创建临时目录并设置为 HOME 环境变量
//
// 返回:
//   - tempDir: 临时目录路径
//
// 示例:
//   tempDir := testutils.SetupTestEnv(t)
//   // 使用 tempDir 进行测试
func SetupTestEnv(t *testing.T) string {
	t.Helper()
	tempDir := t.TempDir()
	t.Setenv("HOME", tempDir)
	return tempDir
}

// WithTestEnv 在测试环境中执行函数
// 自动设置测试环境并在函数执行后恢复（如果需要）
//
// 参数:
//   - t: 测试对象
//   - fn: 在测试环境中执行的函数，接收临时目录路径作为参数
//
// 示例:
//   testutils.WithTestEnv(t, func(tempDir string) {
//       // 使用 tempDir 进行测试
//   })
func WithTestEnv(t *testing.T, fn func(tempDir string)) {
	t.Helper()
	tempDir := SetupTestEnv(t)
	fn(tempDir)
}


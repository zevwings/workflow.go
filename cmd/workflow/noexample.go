//go:build !example

package main

import "github.com/spf13/cobra"

// registerExampleCommands 空实现（正常构建时不包含示例命令）
// 当使用 -tags=example 构建时，此文件会被排除，使用 example.go 中的实现
func registerExampleCommands(rootCmd *cobra.Command) {
	// 正常构建时不注册示例命令
}

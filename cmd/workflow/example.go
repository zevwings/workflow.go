//go:build example

package main

import (
	"github.com/spf13/cobra"
	"github.com/zevwings/workflow/internal/commands/example"
)

// registerExampleCommands 注册示例命令
// 此函数在 main() 中调用，确保 rootCmd 已经创建
// 仅在构建时使用 -tags=example 才会包含此实现
func registerExampleCommands(rootCmd *cobra.Command) {
	rootCmd.AddCommand(example.NewDemoCmd())
	rootCmd.AddCommand(example.NewDemoFormCmd())
}

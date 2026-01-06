package output

import (
	"fmt"
	"os"

	"github.com/fatih/color"
)

// Output 输出工具
type Output struct {
	verbose bool
}

// NewOutput 创建新的输出工具
func NewOutput(verbose bool) *Output {
	return &Output{verbose: verbose}
}

// Info 输出信息
func (o *Output) Info(format string, args ...interface{}) {
	color.New(color.FgCyan).Printf("ℹ %s\n", fmt.Sprintf(format, args...))
}

// Success 输出成功信息
func (o *Output) Success(format string, args ...interface{}) {
	color.New(color.FgGreen).Printf("✓ %s\n", fmt.Sprintf(format, args...))
}

// Warning 输出警告信息
func (o *Output) Warning(format string, args ...interface{}) {
	color.New(color.FgYellow).Printf("⚠ %s\n", fmt.Sprintf(format, args...))
}

// Error 输出错误信息
func (o *Output) Error(format string, args ...interface{}) {
	color.New(color.FgRed).Printf("✗ %s\n", fmt.Sprintf(format, args...))
}

// Fatal 输出致命错误并退出
func (o *Output) Fatal(format string, args ...interface{}) {
	o.Error(format, args...)
	os.Exit(1)
}

// Debug 输出调试信息
func (o *Output) Debug(format string, args ...interface{}) {
	if o.verbose {
		color.New(color.FgMagenta).Printf("DEBUG: %s\n", fmt.Sprintf(format, args...))
	}
}

// Print 普通输出
func (o *Output) Print(format string, args ...interface{}) {
	fmt.Printf(format, args...)
}

// Println 普通输出并换行
func (o *Output) Println(format string, args ...interface{}) {
	fmt.Println(fmt.Sprintf(format, args...))
}

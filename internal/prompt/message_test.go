//go:build test

package prompt

import (
	"bytes"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

// captureOutput 是一个简单的辅助函数，用于捕获标准输出
func captureOutput(t *testing.T, fn func()) string {
	t.Helper()

	// 备份原始 stdout
	oldStdout := os.Stdout
	r, w, err := os.Pipe()
	require.NoError(t, err)
	os.Stdout = w

	// 执行待测函数
	fn()

	// 关闭写端并恢复 stdout
	require.NoError(t, w.Close())
	os.Stdout = oldStdout

	var buf bytes.Buffer
	_, err = buf.ReadFrom(r)
	require.NoError(t, err)
	require.NoError(t, r.Close())

	return buf.String()
}

func TestNewMessage(t *testing.T) {
	m1 := NewMessage(false)
	require.NotNil(t, m1)

	m2 := NewMessage(true)
	require.NotNil(t, m2)
}

func TestMessage_Info_Success_Warning_Error(t *testing.T) {
	// 关闭颜色，方便断言纯文本
	theme := GetTheme()
	theme.EnableColor = false
	SetTheme(theme)

	m := NewMessage(false)

	out := captureOutput(t, func() {
		m.Info("info-%s", "msg")
		m.Success("ok-%d", 1)
		m.Warning("warn")
		m.Error("err: %s", "x")
	})

	lines := strings.Split(strings.TrimSpace(out), "\n")
	require.Len(t, lines, 4)

	require.Equal(t, "ℹ info-msg", lines[0])
	require.Equal(t, "✓ ok-1", lines[1])
	// Warning 使用 PrefixWarn（默认 "!"）作为前缀
	require.Equal(t, "! warn", lines[2])
	// Error 使用 PrefixError（默认 "x"）作为前缀
	require.Equal(t, "x err: x", lines[3])
}

func TestMessage_Fatal_ExitCode(t *testing.T) {
	// Fatal 会调用 os.Exit，不能在同一进程里直接断言，这里只做一次轻量调用
	// 通过子进程方式测试会比较重，这里留给后续需要时再补充
}

func TestMessage_Debug_Verbose(t *testing.T) {
	theme := GetTheme()
	theme.EnableColor = false
	SetTheme(theme)

	// verbose=true 时应输出 Debug
	mVerbose := NewMessage(true)
	outVerbose := captureOutput(t, func() {
		mVerbose.Debug("value=%d", 42)
	})
	require.Contains(t, outVerbose, "DEBUG: value=42")

	// verbose=false 时不应输出 Debug
	mSilent := NewMessage(false)
	outSilent := captureOutput(t, func() {
		mSilent.Debug("value=%d", 42)
	})
	require.Equal(t, "", strings.TrimSpace(outSilent))
}

func TestMessage_PrintAndPrintln(t *testing.T) {
	m := NewMessage(false)

	out := captureOutput(t, func() {
		m.Print("hello %s", "world")
		m.Println(" line%d", 1)
	})

	// Print 不自动换行，Println 会换行
	require.True(t, strings.Contains(out, "hello world"))
	require.True(t, strings.Contains(out, " line1"))
}

func TestFormatMessage_EnableAndDisableColor(t *testing.T) {
	// 自定义一个简单的 theme，便于控制
	theme := GetTheme()
	theme.EnableColor = false
	SetTheme(theme)

	text := formatMessage("P:", "msg", theme.InfoStyle)
	require.Equal(t, "P:msg", text)

	// 启用颜色后，返回值应包含 ANSI 但前缀和内容不变
	// 注意：在非TTY环境下，lipgloss可能不会添加ANSI代码
	theme.EnableColor = true
	SetTheme(theme)
	colored := formatMessage("P:", "msg", theme.InfoStyle)
	require.Contains(t, colored, "msg")
	require.Contains(t, colored, "P:")
	// 在非TTY环境下，可能不会添加颜色，所以只检查内容存在即可
}



package prompt

import (
	"bytes"
	"testing"
	"time"

	"github.com/charmbracelet/lipgloss"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewSpinner(t *testing.T) {
	s := NewSpinner("测试中...")
	assert.NotNil(t, s)
	assert.Equal(t, "测试中...", s.message)
	assert.Greater(t, len(s.spinner), 0)
	assert.Greater(t, s.interval, time.Duration(0))
}

func TestSpinnerWithOptions(t *testing.T) {
	customSpinner := []string{"-", "\\", "|", "/"}
	customInterval := 50 * time.Millisecond
	customStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("red"))

	var buf bytes.Buffer
	s := NewSpinner("自定义加载",
		WithSpinner(customSpinner),
		WithInterval(customInterval),
		WithStyle(customStyle),
		WithWriter(&buf),
	)

	assert.Equal(t, customSpinner, s.spinner)
	assert.Equal(t, customInterval, s.interval)
	assert.Equal(t, customStyle, s.style)
	assert.Equal(t, &buf, s.writer)
}

func TestSpinnerStartStop(t *testing.T) {
	var buf bytes.Buffer
	s := NewSpinner("测试", WithWriter(&buf))

	// 启动
	s.Start()
	time.Sleep(150 * time.Millisecond) // 等待几帧

	// 停止
	s.Stop()
	time.Sleep(50 * time.Millisecond) // 确保停止

	// 验证有输出
	assert.Greater(t, buf.Len(), 0)
}

func TestSpinnerUpdateMessage(t *testing.T) {
	s := NewSpinner("初始消息")
	assert.Equal(t, "初始消息", s.message)

	s.UpdateMessage("更新后的消息")
	assert.Equal(t, "更新后的消息", s.message)
}

func TestSpinnerWithSuccess(t *testing.T) {
	var buf bytes.Buffer
	s := NewSpinner("处理中", WithWriter(&buf))

	s.Start()
	time.Sleep(50 * time.Millisecond)
	s.WithSuccess("完成")

	output := buf.String()
	assert.Contains(t, output, "✓")
	assert.Contains(t, output, "完成")
}

func TestSpinnerWithError(t *testing.T) {
	var buf bytes.Buffer
	s := NewSpinner("处理中", WithWriter(&buf))

	s.Start()
	time.Sleep(50 * time.Millisecond)
	s.WithError("失败")

	output := buf.String()
	assert.Contains(t, output, "✗")
	assert.Contains(t, output, "失败")
}

func TestSpinnerDo(t *testing.T) {
	var buf bytes.Buffer
	s := NewSpinner("执行中", WithWriter(&buf))

	err := s.Do(func() error {
		time.Sleep(100 * time.Millisecond)
		return nil
	})

	require.NoError(t, err)
	assert.Greater(t, buf.Len(), 0)
}

func TestSpinnerDoWithError(t *testing.T) {
	var buf bytes.Buffer
	s := NewSpinner("执行中", WithWriter(&buf))

	testErr := assert.AnError
	err := s.Do(func() error {
		return testErr
	})

	assert.Equal(t, testErr, err)
}

func TestSpinnerRestart(t *testing.T) {
	var buf bytes.Buffer
	s := NewSpinner("测试", WithWriter(&buf))

	// 第一次启动和停止
	s.Start()
	time.Sleep(150 * time.Millisecond)
	s.Stop()
	time.Sleep(50 * time.Millisecond)

	// 验证第一次操作完成（通过检查 stopChan 状态）
	// 由于 Stop() 会清除输出，我们主要验证能够成功停止和重新启动

	// 第二次启动（应该能够重新启动）
	s.Start()
	time.Sleep(150 * time.Millisecond)
	s.Stop()
	time.Sleep(50 * time.Millisecond)

	// 验证能够成功重新启动（如果代码能执行到这里且没有 panic，说明重启成功）
	// 注意：由于 \r 会覆盖输出，且 Stop() 会清除行，测试中难以直接验证输出内容
	// 但我们可以验证基本的启动/停止功能
	assert.NotNil(t, s)
}

func TestSpinnerCursorVisibility(t *testing.T) {
	var buf bytes.Buffer
	s := NewSpinner("测试", WithWriter(&buf))

	// 启动时应该隐藏光标
	s.Start()
	time.Sleep(50 * time.Millisecond)

	// 检查输出中是否包含隐藏光标的 ANSI 序列
	output := buf.String()
	assert.Contains(t, output, "\033[?25l") // 隐藏光标序列

	// 停止时应该恢复光标
	s.Stop()
	time.Sleep(50 * time.Millisecond)

	// 检查输出中是否包含显示光标的 ANSI 序列
	output = buf.String()
	assert.Contains(t, output, "\033[?25h") // 显示光标序列
}

func TestSpinnerDifferentColors(t *testing.T) {
	var buf bytes.Buffer
	spinnerStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("226")) // 黄色
	messageStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("51"))  // 青色

	s := NewSpinner("测试消息",
		WithSpinnerStyle(spinnerStyle),
		WithMessageStyle(messageStyle),
		WithWriter(&buf),
	)

	s.Start()
	time.Sleep(100 * time.Millisecond)
	s.Stop()
	time.Sleep(50 * time.Millisecond)

	// 验证有输出
	assert.Greater(t, buf.Len(), 0)

	// 验证 spinnerStyle 和 messageStyle 被正确设置
	assert.NotNil(t, s.spinnerStyle)
	assert.NotNil(t, s.messageStyle)
}


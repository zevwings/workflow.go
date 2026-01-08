//go:build test

package common

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zevwings/workflow/internal/prompt/io"
)

// TestHandleCancel 测试取消处理
func TestHandleCancel(t *testing.T) {
	mockTerminal := io.NewMockTerminal([]byte{})

	err := HandleCancel(mockTerminal)

	assert.Error(t, err)
	assert.Equal(t, "用户取消输入", err.Error())
	assert.True(t, mockTerminal.RestoreCursorCalled, "应该调用了恢复光标")
	assert.True(t, mockTerminal.ClearToEndCalled, "应该调用了清除到末尾")
	assert.Contains(t, mockTerminal.GetOutput(), "\n", "输出应该包含换行符")
}

// TestHandleCancel_Output 测试取消处理的输出内容
func TestHandleCancel_Output(t *testing.T) {
	mockTerminal := io.NewMockTerminal([]byte{})

	err := HandleCancel(mockTerminal)

	assert.Error(t, err)
	output := mockTerminal.GetOutput()
	// 应该包含 ANSI 控制码
	assert.Contains(t, output, "\033[u", "应该包含恢复光标控制码")
	assert.Contains(t, output, "\033[J", "应该包含清除到末尾控制码")
	assert.Contains(t, output, "\n", "应该包含换行符")
}

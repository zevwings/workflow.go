//go:build test

package io

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestInteractiveRenderer_RenderWithPrompt 测试渲染带提示消息的交互界面
func TestInteractiveRenderer_RenderWithPrompt(t *testing.T) {
	mockTerminal := NewMockTerminal([]byte{})
	renderer := NewInteractiveRenderer(mockTerminal)

	renderCalled := false
	renderIsFirst := false
	err := renderer.RenderWithPrompt("测试提示", func(isFirst bool) error {
		renderCalled = true
		renderIsFirst = isFirst
		return nil
	})

	assert.NoError(t, err)
	assert.True(t, renderCalled, "渲染函数应该被调用")
	assert.True(t, renderIsFirst, "首次渲染时 isFirst 应该为 true")
	assert.True(t, mockTerminal.SaveCursorCalled, "应该调用了保存光标")
	assert.True(t, mockTerminal.ResetFormatCalled, "应该调用了重置格式")
	assert.Contains(t, mockTerminal.GetOutput(), "测试提示", "输出应该包含提示消息")
}

// TestInteractiveRenderer_RenderWithPrompt_RenderError 测试渲染函数返回错误
func TestInteractiveRenderer_RenderWithPrompt_RenderError(t *testing.T) {
	mockTerminal := NewMockTerminal([]byte{})
	renderer := NewInteractiveRenderer(mockTerminal)

	testErr := errors.New("渲染错误")
	err := renderer.RenderWithPrompt("测试提示", func(isFirst bool) error {
		return testErr
	})

	assert.Error(t, err)
	assert.Equal(t, testErr, err)
}

// TestInteractiveRenderer_ReRender 测试重新渲染
func TestInteractiveRenderer_ReRender(t *testing.T) {
	mockTerminal := NewMockTerminal([]byte{})
	renderer := NewInteractiveRenderer(mockTerminal)

	renderCalled := false
	renderIsFirst := false
	err := renderer.ReRender(func(isFirst bool) error {
		renderCalled = true
		renderIsFirst = isFirst
		return nil
	})

	assert.NoError(t, err)
	assert.True(t, renderCalled, "渲染函数应该被调用")
	assert.False(t, renderIsFirst, "重新渲染时 isFirst 应该为 false")
	assert.True(t, mockTerminal.RestoreCursorCalled, "应该调用了恢复光标")
	assert.True(t, mockTerminal.ClearToEndCalled, "应该调用了清除到末尾")
	assert.True(t, mockTerminal.MoveToStartCalled > 0, "应该调用了移动到行首")
}

// TestInteractiveRenderer_ReRender_Error 测试重新渲染时返回错误
func TestInteractiveRenderer_ReRender_Error(t *testing.T) {
	mockTerminal := NewMockTerminal([]byte{})
	renderer := NewInteractiveRenderer(mockTerminal)

	testErr := errors.New("重新渲染错误")
	err := renderer.ReRender(func(isFirst bool) error {
		return testErr
	})

	assert.Error(t, err)
	assert.Equal(t, testErr, err)
}

// TestInteractiveRenderer_GetTerminal 测试获取终端接口
func TestInteractiveRenderer_GetTerminal(t *testing.T) {
	mockTerminal := NewMockTerminal([]byte{})
	renderer := NewInteractiveRenderer(mockTerminal)

	terminal := renderer.GetTerminal()
	assert.Equal(t, mockTerminal, terminal, "应该返回相同的终端实例")
}


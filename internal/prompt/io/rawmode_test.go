//go:build test

package io

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/term"
)

// TestRawModeManager_WithRawMode_Success 测试成功设置原始模式并执行函数
func TestRawModeManager_WithRawMode_Success(t *testing.T) {
	mockTerminal := NewMockTerminal([]byte{})
	mgr := NewRawModeManager(mockTerminal)

	executed := false
	err := mgr.WithRawMode(func() error {
		executed = true
		assert.True(t, mockTerminal.RawModeEnabled, "原始模式应该已启用")
		assert.True(t, mockTerminal.HideCursorCalled, "应该调用了隐藏光标")
		return nil
	})

	assert.NoError(t, err)
	assert.True(t, executed, "函数应该被执行")
	assert.False(t, mockTerminal.RawModeEnabled, "原始模式应该已恢复")
	assert.True(t, mockTerminal.ShowCursorCalled, "应该调用了显示光标")
	// MockTerminal 的 MakeRaw 返回 nil state，这是正常的（用于测试）
	// 实际使用中会返回真实的终端状态
}

// TestRawModeManager_WithRawMode_Error 测试函数执行错误时仍能正确恢复
func TestRawModeManager_WithRawMode_Error(t *testing.T) {
	mockTerminal := NewMockTerminal([]byte{})
	mgr := NewRawModeManager(mockTerminal)

	testErr := errors.New("测试错误")
	err := mgr.WithRawMode(func() error {
		return testErr
	})

	assert.Error(t, err)
	assert.Equal(t, testErr, err)
	assert.False(t, mockTerminal.RawModeEnabled, "原始模式应该已恢复")
	assert.True(t, mockTerminal.ShowCursorCalled, "应该调用了显示光标")
}

// failingTerminal 用于测试 MakeRaw 失败的终端
type failingTerminal struct {
	*MockTerminal
}

func (f *failingTerminal) MakeRaw() (*term.State, error) {
	return nil, errors.New("MakeRaw 失败")
}

// TestRawModeManager_WithRawMode_MakeRawError 测试 MakeRaw 失败的情况
func TestRawModeManager_WithRawMode_MakeRawError(t *testing.T) {
	mockTerminal := NewMockTerminal([]byte{})
	failingTerm := &failingTerminal{MockTerminal: mockTerminal}
	mgr := NewRawModeManager(failingTerm)

	executed := false
	err := mgr.WithRawMode(func() error {
		executed = true
		return nil
	})

	assert.Error(t, err)
	assert.False(t, executed, "函数不应该被执行")
	assert.Contains(t, err.Error(), "设置终端原始模式失败")
}

// TestRawModeManager_WithRawModeAndFallback_Success 测试成功设置原始模式
func TestRawModeManager_WithRawModeAndFallback_Success(t *testing.T) {
	mockTerminal := NewMockTerminal([]byte{})
	mgr := NewRawModeManager(mockTerminal)

	executed := false
	fallbackExecuted := false
	err := mgr.WithRawModeAndFallback(
		func() error {
			executed = true
			return nil
		},
		func() error {
			fallbackExecuted = true
			return nil
		},
	)

	assert.NoError(t, err)
	assert.True(t, executed, "主函数应该被执行")
	assert.False(t, fallbackExecuted, "Fallback 不应该被执行")
	assert.False(t, mockTerminal.RawModeEnabled, "原始模式应该已恢复")
}

// TestRawModeManager_WithRawModeAndFallback_Fallback 测试 MakeRaw 失败时执行 Fallback
func TestRawModeManager_WithRawModeAndFallback_Fallback(t *testing.T) {
	mockTerminal := NewMockTerminal([]byte{})
	failingTerm := &failingTerminal{MockTerminal: mockTerminal}
	mgr := NewRawModeManager(failingTerm)

	executed := false
	fallbackExecuted := false
	err := mgr.WithRawModeAndFallback(
		func() error {
			executed = true
			return nil
		},
		func() error {
			fallbackExecuted = true
			return nil
		},
	)

	assert.NoError(t, err)
	assert.False(t, executed, "主函数不应该被执行")
	assert.True(t, fallbackExecuted, "Fallback 应该被执行")
}

// TestRawModeManager_GetTerminal 测试获取终端接口
func TestRawModeManager_GetTerminal(t *testing.T) {
	mockTerminal := NewMockTerminal([]byte{})
	mgr := NewRawModeManager(mockTerminal)

	terminal := mgr.GetTerminal()
	assert.Equal(t, mockTerminal, terminal, "应该返回相同的终端实例")
}

//go:build test

package common

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestNavigationHandler_ProcessArrowKey_Up 测试向上导航
func TestNavigationHandler_ProcessArrowKey_Up(t *testing.T) {
	testCases := []struct {
		name          string
		itemCount     int
		currentIndex  int
		expectedIndex int
		shouldRender  bool
	}{
		{"从中间向上", 5, 3, 2, true},
		{"从顶部向上（非循环）", 5, 0, 0, false},
		{"从第二个向上", 5, 1, 0, true},
		{"单个选项", 1, 0, 0, false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			handler := NewNavigationHandler(tc.itemCount, false)
			newIndex, shouldRender := handler.ProcessArrowKey(tc.currentIndex, "up")

			assert.Equal(t, tc.expectedIndex, newIndex)
			assert.Equal(t, tc.shouldRender, shouldRender)
		})
	}
}

// TestNavigationHandler_ProcessArrowKey_Down 测试向下导航
func TestNavigationHandler_ProcessArrowKey_Down(t *testing.T) {
	testCases := []struct {
		name          string
		itemCount     int
		currentIndex  int
		expectedIndex int
		shouldRender  bool
	}{
		{"从中间向下", 5, 2, 3, true},
		{"从底部向下（非循环）", 5, 4, 4, false},
		{"从第一个向下", 5, 0, 1, true},
		{"单个选项", 1, 0, 0, false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			handler := NewNavigationHandler(tc.itemCount, false)
			newIndex, shouldRender := handler.ProcessArrowKey(tc.currentIndex, "down")

			assert.Equal(t, tc.expectedIndex, newIndex)
			assert.Equal(t, tc.shouldRender, shouldRender)
		})
	}
}

// TestNavigationHandler_ProcessArrowKey_Cyclic 测试循环导航
func TestNavigationHandler_ProcessArrowKey_Cyclic(t *testing.T) {
	handler := NewNavigationHandler(5, true)

	// 从顶部向上应该循环到底部
	newIndex, shouldRender := handler.ProcessArrowKey(0, "up")
	assert.Equal(t, 4, newIndex)
	assert.True(t, shouldRender)

	// 从底部向下应该循环到顶部
	newIndex, shouldRender = handler.ProcessArrowKey(4, "down")
	assert.Equal(t, 0, newIndex)
	assert.True(t, shouldRender)
}

// TestNavigationHandler_ProcessArrowKey_InvalidDirection 测试无效方向
func TestNavigationHandler_ProcessArrowKey_InvalidDirection(t *testing.T) {
	handler := NewNavigationHandler(5, false)

	newIndex, shouldRender := handler.ProcessArrowKey(2, "left")
	assert.Equal(t, 2, newIndex)
	assert.False(t, shouldRender)

	newIndex, shouldRender = handler.ProcessArrowKey(2, "right")
	assert.Equal(t, 2, newIndex)
	assert.False(t, shouldRender)

	newIndex, shouldRender = handler.ProcessArrowKey(2, "invalid")
	assert.Equal(t, 2, newIndex)
	assert.False(t, shouldRender)
}

// TestNavigationHandler_ProcessArrowKey_EmptyItems 测试空选项列表
func TestNavigationHandler_ProcessArrowKey_EmptyItems(t *testing.T) {
	handler := NewNavigationHandler(0, false)

	newIndex, shouldRender := handler.ProcessArrowKey(0, "up")
	assert.Equal(t, 0, newIndex)
	assert.False(t, shouldRender)

	newIndex, shouldRender = handler.ProcessArrowKey(0, "down")
	assert.Equal(t, 0, newIndex)
	assert.False(t, shouldRender)
}

// TestNavigationHandler_ValidateIndex 测试索引验证
func TestNavigationHandler_ValidateIndex(t *testing.T) {
	handler := NewNavigationHandler(5, false)

	testCases := []struct {
		name          string
		index         int
		expectedIndex int
	}{
		{"有效索引", 2, 2},
		{"有效索引（第一个）", 0, 0},
		{"有效索引（最后一个）", 4, 4},
		{"负索引", -1, 0},
		{"超出范围索引", 10, 0},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			validIndex := handler.ValidateIndex(tc.index)
			assert.Equal(t, tc.expectedIndex, validIndex)
		})
	}
}

// TestNavigationHandler_ValidateIndex_EmptyItems 测试空选项列表的索引验证
func TestNavigationHandler_ValidateIndex_EmptyItems(t *testing.T) {
	handler := NewNavigationHandler(0, false)

	validIndex := handler.ValidateIndex(0)
	assert.Equal(t, 0, validIndex)

	validIndex = handler.ValidateIndex(5)
	assert.Equal(t, 0, validIndex)
}

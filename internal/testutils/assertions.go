//go:build test

package testutils

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// AssertErrorContains 断言错误包含指定消息
//
// 参数:
//   - t: 测试对象
//   - err: 错误对象（可以为 nil）
//   - msg: 期望的错误消息片段
//
// 示例:
//
//	err := someFunction()
//	testutils.AssertErrorContains(t, err, "invalid")
func AssertErrorContains(t *testing.T, err error, msg string) {
	t.Helper()
	require.Error(t, err, "期望有错误，但没有错误")
	assert.Contains(t, err.Error(), msg, "错误消息应该包含: %s", msg)
}

// AssertNoError 断言没有错误（增强版的 require.NoError）
//
// 参数:
//   - t: 测试对象
//   - err: 错误对象
//
// 如果 err 不为 nil，会立即终止测试
//
// 示例:
//
//	result, err := someFunction()
//	testutils.AssertNoError(t, err)
func AssertNoError(t *testing.T, err error) {
	t.Helper()
	require.NoError(t, err)
}

// AssertNotEmpty 断言字符串不为空
//
// 参数:
//   - t: 测试对象
//   - value: 要检查的值
//   - msgAndArgs: 可选的错误消息和参数
//
// 示例:
//
//	id := getID()
//	testutils.AssertNotEmpty(t, id, "ID 不应该为空")
func AssertNotEmpty(t *testing.T, value string, msgAndArgs ...interface{}) {
	t.Helper()
	assert.NotEmpty(t, value, msgAndArgs...)
}

// AssertEmpty 断言字符串为空
//
// 参数:
//   - t: 测试对象
//   - value: 要检查的值
//   - msgAndArgs: 可选的错误消息和参数
//
// 示例:
//
//	result := someFunction()
//	testutils.AssertEmpty(t, result)
func AssertEmpty(t *testing.T, value string, msgAndArgs ...interface{}) {
	t.Helper()
	assert.Empty(t, value, msgAndArgs...)
}

// AssertEqualStrings 断言两个字符串相等（提供更清晰的错误消息）
//
// 参数:
//   - t: 测试对象
//   - expected: 期望值
//   - actual: 实际值
//   - msgAndArgs: 可选的错误消息和参数
//
// 示例:
//
//	result := someFunction()
//	testutils.AssertEqualStrings(t, "expected", result)
func AssertEqualStrings(t *testing.T, expected, actual string, msgAndArgs ...interface{}) {
	t.Helper()
	if len(msgAndArgs) == 0 {
		msgAndArgs = []interface{}{fmt.Sprintf("期望: %q, 实际: %q", expected, actual)}
	}
	assert.Equal(t, expected, actual, msgAndArgs...)
}

// AssertNotNilAndNotEmpty 断言值不为 nil 且字符串不为空
//
// 参数:
//   - t: 测试对象
//   - value: 要检查的值
//   - msgAndArgs: 可选的错误消息和参数
//
// 示例:
//
//	result := getResult()
//	testutils.AssertNotNilAndNotEmpty(t, result)
func AssertNotNilAndNotEmpty(t *testing.T, value *string, msgAndArgs ...interface{}) {
	t.Helper()
	require.NotNil(t, value, msgAndArgs...)
	assert.NotEmpty(t, *value, msgAndArgs...)
}

// AssertInSlice 断言值在切片中
//
// 参数:
//   - t: 测试对象
//   - slice: 切片
//   - value: 要查找的值
//   - msgAndArgs: 可选的错误消息和参数
//
// 示例:
//
//	options := []string{"a", "b", "c"}
//	testutils.AssertInSlice(t, options, "b")
func AssertInSlice(t *testing.T, slice []string, value string, msgAndArgs ...interface{}) {
	t.Helper()
	found := false
	for _, v := range slice {
		if v == value {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("期望 %q 在切片中，但未找到。切片: %v", value, slice)
	}
}

// AssertSliceContains 断言切片包含指定元素（支持任意类型）
//
// 参数:
//   - t: 测试对象
//   - slice: 切片
//   - element: 要查找的元素
//   - msgAndArgs: 可选的错误消息和参数
//
// 示例:
//
//	numbers := []int{1, 2, 3}
//	testutils.AssertSliceContains(t, numbers, 2)
func AssertSliceContains(t *testing.T, slice interface{}, element interface{}, msgAndArgs ...interface{}) {
	t.Helper()
	assert.Contains(t, slice, element, msgAndArgs...)
}

// AssertSliceNotEmpty 断言切片不为空
//
// 参数:
//   - t: 测试对象
//   - slice: 切片
//   - msgAndArgs: 可选的错误消息和参数
//
// 示例:
//
//	results := getResults()
//	testutils.AssertSliceNotEmpty(t, results)
func AssertSliceNotEmpty(t *testing.T, slice interface{}, msgAndArgs ...interface{}) {
	t.Helper()
	assert.NotEmpty(t, slice, msgAndArgs...)
}

// AssertMapContainsKey 断言 map 包含指定键
//
// 参数:
//   - t: 测试对象
//   - m: map 对象
//   - key: 要查找的键
//   - msgAndArgs: 可选的错误消息和参数
//
// 示例:
//
//	config := map[string]string{"key1": "value1"}
//	testutils.AssertMapContainsKey(t, config, "key1")
func AssertMapContainsKey(t *testing.T, m interface{}, key interface{}, msgAndArgs ...interface{}) {
	t.Helper()
	// 使用 Contains 可以检查 map 的键
	switch v := m.(type) {
	case map[string]interface{}:
		_, ok := v[key.(string)]
		if !ok {
			t.Errorf("期望 map 包含键 %v，但未找到。Map: %v", key, v)
		}
	case map[string]string:
		_, ok := v[key.(string)]
		if !ok {
			t.Errorf("期望 map 包含键 %v，但未找到。Map: %v", key, v)
		}
	default:
		// 对于其他类型，使用通用的 Contains
		assert.Contains(t, m, key, msgAndArgs...)
	}
}

// AssertBetween 断言值在指定范围内（包含边界）
//
// 参数:
//   - t: 测试对象
//   - value: 要检查的值
//   - min: 最小值（包含）
//   - max: 最大值（包含）
//   - msgAndArgs: 可选的错误消息和参数
//
// 示例:
//
//	length := len(someString)
//	testutils.AssertBetween(t, length, 1, 100)
func AssertBetween(t *testing.T, value, min, max int, msgAndArgs ...interface{}) {
	t.Helper()
	if value < min || value > max {
		if len(msgAndArgs) == 0 {
			t.Errorf("期望值在 [%d, %d] 范围内，实际: %d", min, max, value)
		} else {
			t.Errorf("%s", fmt.Sprint(msgAndArgs...))
		}
	}
}


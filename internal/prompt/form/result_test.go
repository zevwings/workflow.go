//go:build test

package form

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// ==================== NewFormResult 测试 ====================

func TestNewFormResult(t *testing.T) {
	result := NewFormResult()
	assert.NotNil(t, result)
	assert.NotNil(t, result.Values)
	assert.Equal(t, 0, len(result.Values))
}

// ==================== Set 和 GetString 测试 ====================

func TestFormResult_SetAndGetString(t *testing.T) {
	result := NewFormResult()

	// 设置字符串值
	result.Set("name", "张三")
	assert.Equal(t, "张三", result.GetString("name"))

	// 获取不存在的键
	assert.Empty(t, result.GetString("nonexistent"))

	// 设置不同类型的值，GetString 应该能转换
	result.Set("number", 123)
	assert.Equal(t, "123", result.GetString("number"))

	result.Set("bool", true)
	assert.Equal(t, "true", result.GetString("bool"))
}

// ==================== GetBool 测试 ====================

func TestFormResult_GetBool(t *testing.T) {
	result := NewFormResult()

	// 设置布尔值
	result.Set("agree", true)
	assert.True(t, result.GetBool("agree"))

	result.Set("agree", false)
	assert.False(t, result.GetBool("agree"))

	// 获取不存在的键
	assert.False(t, result.GetBool("nonexistent"))

	// 设置非布尔值，应该返回 false
	result.Set("name", "张三")
	assert.False(t, result.GetBool("name"))

	result.Set("number", 123)
	assert.False(t, result.GetBool("number"))
}

// ==================== GetInt 测试 ====================

func TestFormResult_GetInt(t *testing.T) {
	result := NewFormResult()

	// 设置整数值
	result.Set("age", 25)
	assert.Equal(t, 25, result.GetInt("age"))

	result.Set("count", 0)
	assert.Equal(t, 0, result.GetInt("count"))

	result.Set("negative", -10)
	assert.Equal(t, -10, result.GetInt("negative"))

	// 获取不存在的键
	assert.Equal(t, 0, result.GetInt("nonexistent"))

	// 设置非整数值，应该返回 0
	result.Set("name", "张三")
	assert.Equal(t, 0, result.GetInt("name"))

	result.Set("bool", true)
	assert.Equal(t, 0, result.GetInt("bool"))
}

// ==================== GetIntSlice 测试 ====================

func TestFormResult_GetIntSlice(t *testing.T) {
	result := NewFormResult()

	// 设置整数切片
	result.Set("selected", []int{0, 2, 4})
	slice := result.GetIntSlice("selected")
	assert.NotNil(t, slice)
	assert.Equal(t, []int{0, 2, 4}, slice)

	// 空切片
	result.Set("empty", []int{})
	slice = result.GetIntSlice("empty")
	assert.NotNil(t, slice)
	assert.Equal(t, 0, len(slice))

	// 获取不存在的键
	assert.Nil(t, result.GetIntSlice("nonexistent"))

	// 设置非切片值，应该返回 nil
	result.Set("name", "张三")
	assert.Nil(t, result.GetIntSlice("name"))

	result.Set("number", 123)
	assert.Nil(t, result.GetIntSlice("number"))
}

// ==================== GetForm 测试 ====================

func TestFormResult_GetForm(t *testing.T) {
	result := NewFormResult()

	// 设置嵌套表单结果
	nestedResult := NewFormResult()
	nestedResult.Set("nested_name", "嵌套值")
	result.Set("user", nestedResult)

	formResult := result.GetForm("user")
	assert.NotNil(t, formResult)
	assert.Equal(t, "嵌套值", formResult.GetString("nested_name"))

	// 获取不存在的键
	assert.Nil(t, result.GetForm("nonexistent"))

	// 设置非 FormResult 值，应该返回 nil
	result.Set("name", "张三")
	assert.Nil(t, result.GetForm("name"))

	result.Set("number", 123)
	assert.Nil(t, result.GetForm("number"))
}

// ==================== 综合测试 ====================

func TestFormResult_ComplexScenario(t *testing.T) {
	result := NewFormResult()

	// 设置多种类型的值
	result.Set("name", "张三")
	result.Set("age", 25)
	result.Set("agree", true)
	result.Set("selected", []int{1, 3})

	// 设置嵌套表单
	nestedResult := NewFormResult()
	nestedResult.Set("email", "test@example.com")
	result.Set("contact", nestedResult)

	// 验证所有值
	assert.Equal(t, "张三", result.GetString("name"))
	assert.Equal(t, 25, result.GetInt("age"))
	assert.True(t, result.GetBool("agree"))
	assert.Equal(t, []int{1, 3}, result.GetIntSlice("selected"))

	nested := result.GetForm("contact")
	assert.NotNil(t, nested)
	assert.Equal(t, "test@example.com", nested.GetString("email"))
}

// ==================== 覆盖值测试 ====================

func TestFormResult_OverwriteValue(t *testing.T) {
	result := NewFormResult()

	// 设置值
	result.Set("name", "张三")
	assert.Equal(t, "张三", result.GetString("name"))

	// 覆盖值
	result.Set("name", "李四")
	assert.Equal(t, "李四", result.GetString("name"))

	// 覆盖为不同类型（int）
	result.Set("name", 123)
	assert.Equal(t, "123", result.GetString("name"))
	assert.Equal(t, 123, result.GetInt("name")) // int 类型可以正确获取

	// 覆盖为字符串类型
	result.Set("name", "test")
	assert.Equal(t, "test", result.GetString("name"))
	assert.Equal(t, 0, result.GetInt("name")) // 字符串类型，GetInt 返回 0
}

// ==================== 空值测试 ====================

func TestFormResult_EmptyValues(t *testing.T) {
	result := NewFormResult()

	// 设置空字符串
	result.Set("empty", "")
	assert.Empty(t, result.GetString("empty"))

	// 设置 nil（通过 interface{}）
	result.Set("nil_value", nil)
	value, ok := result.Values["nil_value"]
	assert.True(t, ok)
	assert.Nil(t, value)

	// 设置空切片
	result.Set("empty_slice", []int{})
	slice := result.GetIntSlice("empty_slice")
	assert.NotNil(t, slice)
	assert.Equal(t, 0, len(slice))
}

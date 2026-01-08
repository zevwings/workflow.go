package form

import (
	"fmt"
)

// FormResult 表单执行结果
type FormResult struct {
	// Values 字段值映射
	Values map[string]interface{}
}

// NewFormResult 创建新的表单结果
func NewFormResult() *FormResult {
	return &FormResult{
		Values: make(map[string]interface{}),
	}
}

// GetString 获取字符串值
func (r *FormResult) GetString(key string) string {
	value, ok := r.Values[key]
	if !ok {
		return ""
	}
	if str, ok := value.(string); ok {
		return str
	}
	return fmt.Sprintf("%v", value)
}

// GetBool 获取布尔值
func (r *FormResult) GetBool(key string) bool {
	value, ok := r.Values[key]
	if !ok {
		return false
	}
	if b, ok := value.(bool); ok {
		return b
	}
	return false
}

// GetInt 获取整数值
func (r *FormResult) GetInt(key string) int {
	value, ok := r.Values[key]
	if !ok {
		return 0
	}
	if i, ok := value.(int); ok {
		return i
	}
	return 0
}

// GetIntSlice 获取整数切片
func (r *FormResult) GetIntSlice(key string) []int {
	value, ok := r.Values[key]
	if !ok {
		return nil
	}
	if slice, ok := value.([]int); ok {
		return slice
	}
	return nil
}

// GetForm 获取嵌套表单结果
func (r *FormResult) GetForm(key string) *FormResult {
	value, ok := r.Values[key]
	if !ok {
		return nil
	}
	if formResult, ok := value.(*FormResult); ok {
		return formResult
	}
	return nil
}

// Set 设置字段值（内部使用）
func (r *FormResult) Set(key string, value interface{}) {
	r.Values[key] = value
}

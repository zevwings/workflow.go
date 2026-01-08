package form

import (
	"github.com/zevwings/workflow/internal/prompt/input"
)

// FieldType 字段类型
type FieldType string

const (
	// FieldTypeConfirm 确认字段（bool）
	FieldTypeConfirm FieldType = "confirm"
	// FieldTypeInput 文本输入字段（string）
	FieldTypeInput FieldType = "input"
	// FieldTypePassword 密码输入字段（string）
	FieldTypePassword FieldType = "password"
	// FieldTypeSelect 单选字段（int）
	FieldTypeSelect FieldType = "select"
	// FieldTypeMultiSelect 多选字段（[]int）
	FieldTypeMultiSelect FieldType = "multiselect"
	// FieldTypeForm 嵌套表单字段（FormResult）
	FieldTypeForm FieldType = "form"
)

// Condition 条件函数类型
// 基于前面字段的值决定是否执行当前字段
type Condition func(result *FormResult) bool

// FormField 表单字段定义
type FormField struct {
	// Key 字段键名（用于结果映射）
	Key string
	// Type 字段类型
	Type FieldType
	// Prompt 提示消息
	Prompt string
	// DefaultValue 默认值（可选）
	DefaultValue interface{}
	// Validator 验证器（可选，仅用于 input/password 字段）
	Validator input.Validator
	// Condition 条件函数（可选，基于前面字段的值决定是否执行）
	Condition Condition
	// NestedForm 嵌套表单（仅用于 FieldTypeForm）
	NestedForm *FormBuilder
	// Options 选项列表（仅用于 FieldTypeSelect 和 FieldTypeMultiSelect）
	Options []string
	// DefaultIndex 默认选中的索引（仅用于 FieldTypeSelect）
	DefaultIndex int
	// DefaultSelected 默认选中的索引列表（仅用于 FieldTypeMultiSelect）
	DefaultSelected []int
}

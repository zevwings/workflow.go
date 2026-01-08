package form

// FormValidator 表单级验证器函数类型
// 基于整个表单的结果进行验证
type FormValidator func(result *FormResult) error

// ValidateAllRequired 验证所有必填字段（示例验证器）
// 注意：这是一个示例实现，实际使用时需要根据具体需求实现
func ValidateAllRequired(result *FormResult) error {
	// 这里可以根据具体需求实现验证逻辑
	// 例如：检查某些字段是否存在且不为空
	return nil
}

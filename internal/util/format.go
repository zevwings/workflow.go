package util

// MaskSensitiveValue 掩码显示敏感值
//
// 用于在日志或输出中隐藏敏感信息（如 API key、密码等）。
// - 短值（长度 ≤ 12）：完全隐藏，显示为 `***`
// - 长值（长度 > 12）：显示前 4 个字符和后 4 个字符，中间用 `***` 代替
//
// 参数:
//   - value: 需要掩码的值
//
// 返回:
//   - string: 掩码后的值
//
// 示例:
//
//	MaskSensitiveValue("short")           // "***"
//	MaskSensitiveValue("verylongapikey123456") // "very***3456"
func MaskSensitiveValue(value string) string {
	if value == "" {
		return ""
	}

	length := len(value)
	if length <= 12 {
		// 如果值较短，完全隐藏
		return "***"
	}

	// 显示前4个字符和后4个字符，中间用 *** 代替
	return value[:4] + "***" + value[length-4:]
}

// FormatBool 格式化布尔值为 Yes/No
//
// 参数:
//   - b: 布尔值
//
// 返回:
//   - string: "Yes" 或 "No"
func FormatBool(b bool) string {
	if b {
		return "Yes"
	}
	return "No"
}

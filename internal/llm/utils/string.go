package utils

import (
	"strings"
)

// SanitizeBranchName 清理分支名，确保只保留 ASCII 字符
//
// 只保留字母、数字、连字符和下划线，移除所有其他字符。
//
// 参数:
//   - name: 原始分支名
//
// 返回:
//   - string: 清理后的分支名
func SanitizeBranchName(name string) string {
	var result strings.Builder
	for _, r := range name {
		// 只保留字母、数字、连字符和下划线
		if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') || r == '-' || r == '_' {
			result.WriteRune(r)
		}
	}
	return result.String()
}

// CleanFilename 清理文件名，确保只包含有效的文件名字符
//
// 将文件名转换为小写，替换空格为连字符，只保留字母、数字、连字符和下划线，
// 并移除 .md 扩展名（如果存在）。
//
// 参数:
//   - filename: 原始文件名
//
// 返回:
//   - string: 清理后的文件名
func CleanFilename(filename string) string {
	// 转小写
	cleaned := strings.ToLower(strings.TrimSpace(filename))
	// 替换空格为连字符
	cleaned = strings.ReplaceAll(cleaned, " ", "-")

	// 只保留字母、数字、连字符和下划线
	var result strings.Builder
	for _, r := range cleaned {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') || r == '-' || r == '_' {
			result.WriteRune(r)
		}
	}

	cleaned = result.String()

	// 移除 .md 扩展名（如果存在），因为我们会自动添加
	cleaned = strings.TrimSuffix(cleaned, ".md")

	return cleaned
}

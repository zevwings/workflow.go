package utils

import (
	"fmt"
	"strings"
	"unicode/utf8"
)

// TruncateDiff 安全截断 diff 内容
//
// 使用字符边界安全截取，避免在 UTF-8 字符中间截断。
// 会在最后一个换行符处截断，以保持 diff 格式的完整性。
//
// 参数:
//   - diff: 需要截断的 diff 内容
//   - maxLength: 最大字符数（UTF-8 rune 计数）
//
// 返回:
//   - string: 截断后的 diff 内容，如果超过长度会添加截断提示
func TruncateDiff(diff string, maxLength int) string {
	charCount := utf8.RuneCountInString(diff)
	if charCount <= maxLength {
		return diff
	}

	// 使用字符边界安全截取
	var charBoundary int
	count := 0
	for i := range diff {
		if count >= maxLength {
			charBoundary = i
			break
		}
		count++
	}

	if charBoundary == 0 {
		charBoundary = len(diff)
	}

	truncated := diff[:charBoundary]
	// 尝试在最后一个换行符处截断
	lastNewline := strings.LastIndex(truncated, "\n")
	if lastNewline > 0 {
		truncated = diff[:lastNewline]
	}

	return fmt.Sprintf("%s\n... (diff truncated, %d characters total)", truncated, charCount)
}

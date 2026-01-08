package utils

import (
	"strings"
	"unicode/utf8"
)

// ExtractJSONFromMarkdown 从 markdown 代码块中提取 JSON 字符串
//
// 支持以下格式：
// - ```json\n{...}\n```
// - ```\n{...}\n```
// - 纯 JSON 字符串
//
// 参数:
//   - response: 可能包含 markdown 代码块的响应字符串
//
// 返回:
//   - string: 提取的 JSON 字符串（已去除 markdown 代码块包装）
func ExtractJSONFromMarkdown(response string) string {
	trimmed := strings.TrimSpace(response)

	// 尝试提取 JSON（可能包含 markdown 代码块）
	if strings.HasPrefix(trimmed, "```json") {
		// 移除 ```json 开头和 ``` 结尾
		start := strings.Index(trimmed, "\n")
		if start == -1 {
			// 没有换行，直接跳过 ```json 前缀
			start = len("```json")
		} else {
			start++ // 跳过换行符
		}
		end := strings.LastIndex(trimmed, "```")
		if end == -1 || end <= start {
			// 没有找到结束标记，或者结束标记在开始位置之前，返回剩余部分
			return strings.TrimSpace(trimmed[start:])
		}
		return strings.TrimSpace(trimmed[start:end])
	} else if strings.HasPrefix(trimmed, "```") {
		// 移除 ``` 开头和 ``` 结尾
		start := strings.Index(trimmed, "\n")
		if start == -1 {
			// 没有换行，直接跳过 ``` 前缀
			start = len("```")
		} else {
			start++ // 跳过换行符
		}
		end := strings.LastIndex(trimmed, "```")
		if end == -1 || end <= start {
			// 没有找到结束标记，或者结束标记在开始位置之前，返回剩余部分
			return strings.TrimSpace(trimmed[start:])
		}
		return strings.TrimSpace(trimmed[start:end])
	}

	return trimmed
}

// FixJSONEscapes 修复 JSON 字符串中的转义问题
//
// LLM 生成的 JSON 可能包含未转义的反斜杠（特别是在 Windows 路径中），
// 这会导致 JSON 解析失败。此函数尝试修复这些转义问题。
//
// 参数:
//   - jsonStr: 需要修复的 JSON 字符串
//
// 返回:
//   - string: 修复后的 JSON 字符串
func FixJSONEscapes(jsonStr string) string {
	var result strings.Builder
	result.Grow(len(jsonStr) * 2)

	inString := false
	escapeNext := false

	for i, ch := range jsonStr {
		if escapeNext {
			escapeNext = false
			result.WriteRune(ch)
			continue
		}

		switch ch {
		case '"':
			if !escapeNext {
				inString = !inString
			}
			result.WriteRune(ch)
		case '\\':
			if inString {
				// 检查下一个字符是否是有效的转义序列
				if i+1 < len(jsonStr) {
					nextRune, _ := utf8.DecodeRuneInString(jsonStr[i+1:])
					// 检查是否是有效的转义字符
					isValidEscape := nextRune == '"' || nextRune == '\\' || nextRune == '/' ||
						nextRune == 'b' || nextRune == 'f' || nextRune == 'n' ||
						nextRune == 'r' || nextRune == 't' || nextRune == 'u'
					if isValidEscape {
						// 有效的转义序列，保留原样
						result.WriteRune(ch)
						escapeNext = true
					} else if nextRune >= 0 && nextRune <= 127 {
						// 无效的转义序列（如 \s, \d），需要转义反斜杠
						result.WriteString("\\\\")
						// 下一个字符会正常处理（不设置 escapeNext）
					} else {
						result.WriteRune(ch)
					}
				} else {
					// 字符串末尾的反斜杠，需要转义
					result.WriteString("\\\\")
				}
			} else {
				result.WriteRune(ch)
			}
		default:
			result.WriteRune(ch)
		}
	}

	return result.String()
}

// ExtractAndFixJSON 从 markdown 代码块中提取并修复 JSON 字符串
//
// 这是 ExtractJSONFromMarkdown 的增强版本，会自动修复 JSON 字符串中的转义问题。
//
// 参数:
//   - response: 可能包含 markdown 代码块的响应字符串
//
// 返回:
//   - string: 提取并修复后的 JSON 字符串
func ExtractAndFixJSON(response string) string {
	extracted := ExtractJSONFromMarkdown(response)
	return FixJSONEscapes(extracted)
}

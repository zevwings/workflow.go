package helpers

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// ParsePRID 解析 PR ID
//
// 支持多种格式：
//   - 数字：123
//   - URL：https://github.com/owner/repo/pull/123
//   - 短 URL：owner/repo#123
//
// 返回:
//   - string: 解析后的 PR ID（数字字符串）
//   - error: 如果解析失败，返回错误
func ParsePRID(prID string) (string, error) {
	// 移除前后空格
	prID = strings.TrimSpace(prID)

	// 如果是纯数字，直接返回
	if matched, _ := regexp.MatchString(`^\d+$`, prID); matched {
		return prID, nil
	}

	// 尝试从 URL 中提取
	// https://github.com/owner/repo/pull/123
	// https://github.com/owner/repo/pull/123/
	if strings.Contains(prID, "/pull/") {
		parts := strings.Split(prID, "/pull/")
		if len(parts) == 2 {
			number := strings.TrimSuffix(strings.TrimPrefix(parts[1], "/"), "/")
			if matched, _ := regexp.MatchString(`^\d+$`, number); matched {
				return number, nil
			}
		}
	}

	// 尝试从短格式中提取
	// owner/repo#123
	if strings.Contains(prID, "#") {
		parts := strings.Split(prID, "#")
		if len(parts) == 2 {
			number := strings.TrimSpace(parts[1])
			if matched, _ := regexp.MatchString(`^\d+$`, number); matched {
				return number, nil
			}
		}
	}

	return "", fmt.Errorf("invalid PR ID format: %s", prID)
}

// ParsePRNumber 解析 PR 编号（返回整数）
//
// 参数:
//   - prID: PR ID（支持多种格式）
//
// 返回:
//   - int: PR 编号
//   - error: 如果解析失败，返回错误
func ParsePRNumber(prID string) (int, error) {
	id, err := ParsePRID(prID)
	if err != nil {
		return 0, err
	}

	number, err := strconv.Atoi(id)
	if err != nil {
		return 0, fmt.Errorf("failed to parse PR number: %w", err)
	}

	return number, nil
}


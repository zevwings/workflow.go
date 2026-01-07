package jira

import (
	"fmt"
	"strings"
)

// ValidateTicketKey 验证 Ticket Key 格式
//
// Jira Ticket Key 格式：PROJECT-NUMBER（如 "PROJ-123"）
//
// 参数:
//   - ticket: Ticket Key
//
// 返回:
//   - error: 如果格式无效，返回错误
func ValidateTicketKey(ticket string) error {
	if ticket == "" {
		return fmt.Errorf("ticket key 不能为空")
	}

	parts := strings.Split(ticket, "-")
	if len(parts) != 2 {
		return fmt.Errorf("ticket key 格式无效，应为 PROJECT-NUMBER（如 PROJ-123）")
	}

	if parts[0] == "" || parts[1] == "" {
		return fmt.Errorf("ticket key 格式无效，项目名和编号不能为空")
	}

	return nil
}

// NormalizeTicketKey 规范化 Ticket Key（转大写）
//
// 参数:
//   - ticket: Ticket Key
//
// 返回:
//   - string: 规范化后的 Ticket Key
func NormalizeTicketKey(ticket string) string {
	return strings.ToUpper(strings.TrimSpace(ticket))
}

// ExtractProjectKey 从 Ticket Key 中提取项目 Key
//
// 参数:
//   - ticket: Ticket Key（如 "PROJ-123"）
//
// 返回:
//   - string: 项目 Key（如 "PROJ"）
func ExtractProjectKey(ticket string) string {
	parts := strings.Split(ticket, "-")
	if len(parts) > 0 {
		return parts[0]
	}
	return ""
}

// ExtractTicketNumber 从 Ticket Key 中提取 Ticket 编号
//
// 参数:
//   - ticket: Ticket Key（如 "PROJ-123"）
//
// 返回:
//   - string: Ticket 编号（如 "123"）
func ExtractTicketNumber(ticket string) string {
	parts := strings.Split(ticket, "-")
	if len(parts) > 1 {
		return parts[1]
	}
	return ""
}


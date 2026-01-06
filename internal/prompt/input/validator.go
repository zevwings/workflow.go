package input

import (
	"fmt"
	"regexp"
	"strings"
)

// Validator 验证函数类型
type Validator func(value string) error

// ValidateRegex 创建一个基于正则表达式的验证器
func ValidateRegex(pattern string, errorMsg string) Validator {
	re := regexp.MustCompile(pattern)
	return func(value string) error {
		if !re.MatchString(value) {
			if errorMsg != "" {
				return fmt.Errorf(errorMsg)
			}
			return fmt.Errorf("输入格式不正确，必须匹配: %s", pattern)
		}
		return nil
	}
}

// ValidateEmail 验证邮箱格式
func ValidateEmail() Validator {
	return ValidateRegex(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`, "请输入有效的邮箱地址")
}

// ValidateURL 验证 URL 格式
func ValidateURL() Validator {
	return ValidateRegex(`^https?://[^\s/$.?#].[^\s]*$`, "请输入有效的 URL 地址")
}

// ValidateRequired 验证输入不能为空
func ValidateRequired() Validator {
	return func(value string) error {
		if strings.TrimSpace(value) == "" {
			return fmt.Errorf("此项为必填项，不能为空")
		}
		return nil
	}
}

// ValidateMinLength 验证最小长度
func ValidateMinLength(minLen int) Validator {
	return func(value string) error {
		if len(value) < minLen {
			return fmt.Errorf("输入长度至少需要 %d 个字符", minLen)
		}
		return nil
	}
}

// ValidateMaxLength 验证最大长度
func ValidateMaxLength(maxLen int) Validator {
	return func(value string) error {
		if len(value) > maxLen {
			return fmt.Errorf("输入长度不能超过 %d 个字符", maxLen)
		}
		return nil
	}
}

// ValidateLength 验证长度范围
func ValidateLength(minLen, maxLen int) Validator {
	return func(value string) error {
		length := len(value)
		if length < minLen || length > maxLen {
			return fmt.Errorf("输入长度必须在 %d 到 %d 个字符之间", minLen, maxLen)
		}
		return nil
	}
}

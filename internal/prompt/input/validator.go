package input

import (
	"fmt"
	"net/mail"
	"net/url"
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
				return fmt.Errorf("%s", errorMsg)
			}
			return fmt.Errorf("输入格式不正确，必须匹配: %s", pattern)
		}
		return nil
	}
}

// ValidateEmail 验证邮箱格式
func ValidateEmail() Validator {
	return func(value string) error {
		if value == "" {
			return fmt.Errorf("请输入有效的邮箱地址")
		}
		addr, err := mail.ParseAddress(value)
		if err != nil {
			return fmt.Errorf("请输入有效的邮箱地址")
		}
		// 确保邮箱地址包含 @ 符号
		email := addr.Address
		parts := strings.Split(email, "@")
		if len(parts) != 2 {
			return fmt.Errorf("请输入有效的邮箱地址")
		}
		// 确保域名部分包含至少一个点（顶级域名）
		domain := parts[1]
		if !strings.Contains(domain, ".") {
			return fmt.Errorf("请输入有效的邮箱地址")
		}
		return nil
	}
}

// ValidateURL 验证 URL 格式
func ValidateURL() Validator {
	return func(value string) error {
		if value == "" {
			return fmt.Errorf("请输入有效的 URL 地址")
		}
		// 检查是否包含空格（URL 不应该包含未编码的空格）
		if strings.Contains(value, " ") {
			return fmt.Errorf("请输入有效的 URL 地址")
		}
		parsedURL, err := url.Parse(value)
		if err != nil {
			return fmt.Errorf("请输入有效的 URL 地址")
		}
		// 确保 URL 有有效的 scheme（http 或 https）
		if parsedURL.Scheme != "http" && parsedURL.Scheme != "https" {
			return fmt.Errorf("请输入有效的 URL 地址")
		}
		// 确保 URL 有有效的 host
		if parsedURL.Host == "" {
			return fmt.Errorf("请输入有效的 URL 地址")
		}
		return nil
	}
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

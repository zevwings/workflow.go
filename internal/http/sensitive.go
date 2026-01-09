package http

import (
	"net/url"
	"regexp"
	"strings"
)

// sensitiveHeaderNames 需要过滤的敏感请求头名称（不区分大小写）
var sensitiveHeaderNames = map[string]bool{
	"authorization": true,
	"api-key":       true,
	"apikey":        true,
	"x-api-key":     true,
	"x-auth-token":  true,
	"x-token":       true,
	"cookie":        true,
	"set-cookie":    true,
}

// sensitiveQueryParams 需要过滤的敏感查询参数名称（不区分大小写）
var sensitiveQueryParams = map[string]bool{
	"api_key":    true,
	"apikey":     true,
	"api-key":    true,
	"token":      true,
	"access_key": true,
	"secret":     true,
	"password":   true,
	"auth":       true,
}

// apiKeyPattern 匹配常见的 API Key 格式
var apiKeyPattern = regexp.MustCompile(`(?i)(api[_-]?key|token|auth|secret|password|access[_-]?key)\s*[=:]\s*([a-zA-Z0-9_\-]{10,})`)

// maskValue 掩码显示敏感值
//
// 如果值长度 <= 4，全部掩码
// 如果值长度 > 4，只显示前 2 个字符和后 2 个字符，中间用 *** 替代
//
// 参数:
//   - value: 需要掩码的值
//
// 返回:
//   - string: 掩码后的值
func maskValue(value string) string {
	if value == "" {
		return ""
	}

	length := len(value)
	if length <= 4 {
		return "****"
	}

	// 显示前 2 个字符和后 2 个字符
	return value[:2] + "***" + value[length-2:]
}

// FilterSensitiveURL 过滤 URL 中的敏感信息
//
// 过滤 URL 查询参数中的敏感信息（如 API Key、Token 等）。
// 保留 URL 的基本结构，只掩码敏感参数的值。
//
// 参数:
//   - rawURL: 原始 URL 字符串
//
// 返回:
//   - string: 过滤后的 URL
func FilterSensitiveURL(rawURL string) string {
	if rawURL == "" {
		return ""
	}

	// 解析 URL
	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		// 如果解析失败，尝试使用正则表达式过滤
		return filterURLWithRegex(rawURL)
	}

	// 过滤查询参数
	if parsedURL.RawQuery != "" {
		query := parsedURL.Query()
		for key := range query {
			keyLower := strings.ToLower(key)
			if sensitiveQueryParams[keyLower] {
				// 掩码敏感参数的值
				query.Set(key, maskValue(query.Get(key)))
			}
		}
		parsedURL.RawQuery = query.Encode()
	}

	// 过滤路径中的敏感信息（如果路径中包含类似 API Key 的模式）
	parsedURL.Path = filterPath(parsedURL.Path)

	return parsedURL.String()
}

// filterURLWithRegex 使用正则表达式过滤 URL（当 URL 解析失败时使用）
func filterURLWithRegex(rawURL string) string {
	// 使用正则表达式匹配并替换敏感信息
	result := apiKeyPattern.ReplaceAllStringFunc(rawURL, func(match string) string {
		// 提取键和值
		parts := strings.SplitN(match, "=", 2)
		if len(parts) != 2 {
			parts = strings.SplitN(match, ":", 2)
		}
		if len(parts) == 2 {
			key := strings.TrimSpace(parts[0])
			value := strings.TrimSpace(parts[1])
			return key + "=" + maskValue(value)
		}
		return "***"
	})
	return result
}

// filterPath 过滤路径中的敏感信息
func filterPath(path string) string {
	// 检查路径中是否包含类似 API Key 的模式
	if apiKeyPattern.MatchString(path) {
		return apiKeyPattern.ReplaceAllStringFunc(path, func(match string) string {
			return "***"
		})
	}
	return path
}

// FilterSensitiveHeaders 过滤请求头中的敏感信息
//
// 过滤敏感请求头（如 Authorization、API-Key 等）的值。
// 返回一个新的请求头映射，敏感头部的值会被掩码。
//
// 参数:
//   - headers: 原始请求头映射
//
// 返回:
//   - map[string]string: 过滤后的请求头映射
func FilterSensitiveHeaders(headers map[string]string) map[string]string {
	if headers == nil {
		return nil
	}

	filtered := make(map[string]string, len(headers))
	for key, value := range headers {
		keyLower := strings.ToLower(key)
		if sensitiveHeaderNames[keyLower] {
			// 掩码敏感请求头的值
			filtered[key] = maskValue(value)
		} else {
			// 保留非敏感请求头
			filtered[key] = value
		}
	}

	return filtered
}

// FilterSensitiveHeaderValue 过滤单个请求头值
//
// 如果请求头名称是敏感的，则掩码其值；否则返回原值。
//
// 参数:
//   - headerName: 请求头名称
//   - headerValue: 请求头值
//
// 返回:
//   - string: 过滤后的请求头值
func FilterSensitiveHeaderValue(headerName, headerValue string) string {
	if headerName == "" || headerValue == "" {
		return headerValue
	}

	keyLower := strings.ToLower(headerName)
	if sensitiveHeaderNames[keyLower] {
		return maskValue(headerValue)
	}

	return headerValue
}

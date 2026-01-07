package http

import (
	"fmt"
	"strings"
)

// HttpMethod HTTP 方法枚举
type HttpMethod string

const (
	MethodGet    HttpMethod = "GET"
	MethodPost   HttpMethod = "POST"
	MethodPut    HttpMethod = "PUT"
	MethodDelete HttpMethod = "DELETE"
	MethodPatch  HttpMethod = "PATCH"
)

// String 返回 HTTP 方法的字符串表示
func (m HttpMethod) String() string {
	return string(m)
}

// ParseHttpMethod 从字符串解析 HTTP 方法
//
// 参数:
//   - s: HTTP 方法字符串（如 "GET", "POST"）
//
// 返回:
//   - HttpMethod: 解析后的 HTTP 方法
//   - error: 如果方法无效，返回错误
func ParseHttpMethod(s string) (HttpMethod, error) {
	method := HttpMethod(strings.ToUpper(s))
	switch method {
	case MethodGet, MethodPost, MethodPut, MethodDelete, MethodPatch:
		return method, nil
	default:
		return "", fmt.Errorf("invalid HTTP method: %s", s)
	}
}

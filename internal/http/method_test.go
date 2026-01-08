package http

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// ==================== HttpMethod 测试 ====================

// TestParseHttpMethod 测试 HTTP 方法解析
func TestParseHttpMethod(t *testing.T) {
	testCases := []struct {
		input    string
		expected HttpMethod
		hasError bool
	}{
		{"GET", MethodGet, false},
		{"POST", MethodPost, false},
		{"PUT", MethodPut, false},
		{"DELETE", MethodDelete, false},
		{"PATCH", MethodPatch, false},
		{"get", MethodGet, false}, // 不区分大小写
		{"post", MethodPost, false},
		{"INVALID", "", true},
		{"", "", true},
	}

	for _, tc := range testCases {
		t.Run(tc.input, func(t *testing.T) {
			method, err := ParseHttpMethod(tc.input)
			if tc.hasError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expected, method)
			}
		})
	}
}

// TestHttpMethod_String 测试 HTTP 方法字符串转换
func TestHttpMethod_String(t *testing.T) {
	assert.Equal(t, "GET", MethodGet.String())
	assert.Equal(t, "POST", MethodPost.String())
	assert.Equal(t, "PUT", MethodPut.String())
	assert.Equal(t, "DELETE", MethodDelete.String())
	assert.Equal(t, "PATCH", MethodPatch.String())
}


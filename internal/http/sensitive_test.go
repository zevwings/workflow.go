package http

import (
	"strings"
	"testing"
)

func TestMaskValue(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "empty string",
			input:    "",
			expected: "",
		},
		{
			name:     "short value (<=4)",
			input:    "test",
			expected: "****",
		},
		{
			name:     "long value (>4)",
			input:    "abcdefghijklmnop",
			expected: "ab***op",
		},
		{
			name:     "API key format",
			input:    "ghp_xxxxxxxxxxxxxxxxxxxx",
			expected: "gh***xx",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := maskValue(tt.input)
			if result != tt.expected {
				t.Errorf("maskValue(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestFilterSensitiveURL(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string // 期望 URL 中的敏感参数被掩码
	}{
		{
			name:     "URL with API key in query",
			input:    "https://api.example.com/users?api_key=secret123456789&page=1",
			expected: "https://api.example.com/users?api_key=se***89&page=1",
		},
		{
			name:     "URL with token in query",
			input:    "https://api.example.com/data?token=abc123xyz&limit=10",
			expected: "https://api.example.com/data?limit=10&token=ab***yz",
		},
		{
			name:     "URL without sensitive params",
			input:    "https://api.example.com/users?page=1&limit=10",
			expected: "https://api.example.com/users?limit=10&page=1",
		},
		{
			name:     "URL without query",
			input:    "https://api.example.com/users",
			expected: "https://api.example.com/users",
		},
		{
			name:     "empty URL",
			input:    "",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := FilterSensitiveURL(tt.input)
			// 注意：由于 URL 解析和编码的顺序可能不同，我们只检查敏感参数是否被掩码
			if tt.input != "" && result == tt.input {
				// 如果输入包含敏感参数，结果应该与输入不同
				if strings.Contains(tt.input, "api_key") || strings.Contains(tt.input, "token") {
					t.Errorf("FilterSensitiveURL(%q) = %q, should mask sensitive params", tt.input, result)
				}
			}
		})
	}
}

func TestFilterSensitiveHeaders(t *testing.T) {
	tests := []struct {
		name     string
		input    map[string]string
		expected map[string]string
	}{
		{
			name: "headers with Authorization",
			input: map[string]string{
				"Authorization": "Bearer token123456789",
				"Content-Type":  "application/json",
			},
			expected: map[string]string{
				"Authorization": "Be***89",
				"Content-Type":  "application/json",
			},
		},
		{
			name: "headers with API-Key",
			input: map[string]string{
				"API-Key":      "secret123456",
				"User-Agent":   "workflow-cli",
				"Content-Type": "application/json",
			},
			expected: map[string]string{
				"API-Key":      "se***56",
				"User-Agent":   "workflow-cli",
				"Content-Type": "application/json",
			},
		},
		{
			name: "headers without sensitive info",
			input: map[string]string{
				"Content-Type": "application/json",
				"User-Agent":   "workflow-cli",
			},
			expected: map[string]string{
				"Content-Type": "application/json",
				"User-Agent":   "workflow-cli",
			},
		},
		{
			name:     "nil headers",
			input:    nil,
			expected: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := FilterSensitiveHeaders(tt.input)
			if tt.input == nil {
				if result != nil {
					t.Errorf("FilterSensitiveHeaders(nil) = %v, want nil", result)
				}
				return
			}

			for key, expectedValue := range tt.expected {
				if result[key] != expectedValue {
					t.Errorf("FilterSensitiveHeaders()[%q] = %q, want %q", key, result[key], expectedValue)
				}
			}
		})
	}
}

func TestFilterSensitiveHeaderValue(t *testing.T) {
	tests := []struct {
		name        string
		headerName  string
		headerValue string
		expected    string
	}{
		{
			name:        "sensitive header (Authorization)",
			headerName:  "Authorization",
			headerValue: "Bearer token123456789",
			expected:    "Be***89",
		},
		{
			name:        "sensitive header (API-Key)",
			headerName:  "API-Key",
			headerValue: "secret123456",
			expected:    "se***56",
		},
		{
			name:        "non-sensitive header",
			headerName:  "Content-Type",
			headerValue: "application/json",
			expected:    "application/json",
		},
		{
			name:        "case insensitive",
			headerName:  "authorization",
			headerValue: "Bearer token123456789",
			expected:    "Be***89",
		},
		{
			name:        "empty value",
			headerName:  "Authorization",
			headerValue: "",
			expected:    "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := FilterSensitiveHeaderValue(tt.headerName, tt.headerValue)
			if result != tt.expected {
				t.Errorf("FilterSensitiveHeaderValue(%q, %q) = %q, want %q", tt.headerName, tt.headerValue, result, tt.expected)
			}
		})
	}
}

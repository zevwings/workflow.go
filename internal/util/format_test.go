package util

import (
	"testing"
)

func TestMaskSensitiveValue(t *testing.T) {
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
			name:     "short value (<=12)",
			input:    "short",
			expected: "***",
		},
		{
			name:     "medium value (<=12)",
			input:    "test12345678",
			expected: "***",
		},
		{
			name:     "exactly 12 characters",
			input:    "123456789012",
			expected: "***",
		},
		{
			name:     "long value (>12)",
			input:    "verylongapikey123456",
			expected: "very***3456",
		},
		{
			name:     "API key format",
			input:    "ghp_xxxxxxxxxxxxxxxxxxxx",
			expected: "ghp_***xxxx",
		},
		{
			name:     "very long value",
			input:    "abcdefghijklmnopqrstuvwxyz",
			expected: "abcd***wxyz",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := MaskSensitiveValue(tt.input)
			if result != tt.expected {
				t.Errorf("MaskSensitiveValue(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestFormatBool(t *testing.T) {
	tests := []struct {
		name     string
		input    bool
		expected string
	}{
		{
			name:     "true",
			input:    true,
			expected: "Yes",
		},
		{
			name:     "false",
			input:    false,
			expected: "No",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := FormatBool(tt.input)
			if result != tt.expected {
				t.Errorf("FormatBool(%v) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

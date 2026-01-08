package jira

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// ==================== ValidateTicketKey 测试 ====================

func TestValidateTicketKey(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{
			name:    "valid ticket key",
			input:   "PROJ-123",
			wantErr: false,
		},
		{
			name:    "valid ticket key lowercase",
			input:   "proj-123",
			wantErr: false,
		},
		{
			name:    "empty string",
			input:   "",
			wantErr: true,
		},
		{
			name:    "invalid format - no dash",
			input:   "PROJ123",
			wantErr: true,
		},
		{
			name:    "invalid format - multiple dashes",
			input:   "PROJ-123-456",
			wantErr: true,
		},
		{
			name:    "missing project key",
			input:   "-123",
			wantErr: true,
		},
		{
			name:    "missing ticket number",
			input:   "PROJ-",
			wantErr: true,
		},
		{
			name:    "only dash",
			input:   "-",
			wantErr: true,
		},
		{
			name:    "minimum valid length",
			input:   "A-1",
			wantErr: false,
		},
		{
			name:    "long project key",
			input:   "VERY-LONG-PROJECT-KEY-123",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateTicketKey(tt.input)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

// ==================== NormalizeTicketKey 测试 ====================

func TestNormalizeTicketKey(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "uppercase",
			input:    "PROJ-123",
			expected: "PROJ-123",
		},
		{
			name:     "lowercase",
			input:    "proj-123",
			expected: "PROJ-123",
		},
		{
			name:     "mixed case",
			input:    "Proj-123",
			expected: "PROJ-123",
		},
		{
			name:     "with spaces",
			input:    "  proj-123  ",
			expected: "PROJ-123",
		},
		{
			name:     "with leading spaces",
			input:    "  PROJ-123",
			expected: "PROJ-123",
		},
		{
			name:     "with trailing spaces",
			input:    "PROJ-123  ",
			expected: "PROJ-123",
		},
		{
			name:     "only spaces",
			input:    "   ",
			expected: "",
		},
		{
			name:     "empty string",
			input:    "",
			expected: "",
		},
		{
			name:     "single character project",
			input:    "a-1",
			expected: "A-1",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := NormalizeTicketKey(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// ==================== ExtractProjectKey 测试 ====================

func TestExtractProjectKey(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "valid ticket key",
			input:    "PROJ-123",
			expected: "PROJ",
		},
		{
			name:     "lowercase",
			input:    "proj-123",
			expected: "proj",
		},
		{
			name:     "long project key",
			input:    "VERY-LONG-PROJECT-123",
			expected: "VERY",
		},
		{
			name:     "single character",
			input:    "A-1",
			expected: "A",
		},
		{
			name:     "invalid format - no dash",
			input:    "invalid",
			expected: "invalid",
		},
		{
			name:     "empty string",
			input:    "",
			expected: "",
		},
		{
			name:     "only dash",
			input:    "-",
			expected: "",
		},
		{
			name:     "multiple dashes",
			input:    "PROJ-123-456",
			expected: "PROJ",
		},
		{
			name:     "with spaces",
			input:    " PROJ-123",
			expected: " PROJ",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ExtractProjectKey(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// ==================== ExtractTicketNumber 测试 ====================

func TestExtractTicketNumber(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "valid ticket key",
			input:    "PROJ-123",
			expected: "123",
		},
		{
			name:     "large number",
			input:    "PROJ-999999",
			expected: "999999",
		},
		{
			name:     "single digit",
			input:    "PROJ-1",
			expected: "1",
		},
		{
			name:     "invalid format - no dash",
			input:    "invalid",
			expected: "",
		},
		{
			name:     "empty string",
			input:    "",
			expected: "",
		},
		{
			name:     "only dash",
			input:    "-",
			expected: "",
		},
		{
			name:     "missing ticket number",
			input:    "PROJ-",
			expected: "",
		},
		{
			name:     "multiple dashes",
			input:    "PROJ-123-456",
			expected: "123",
		},
		{
			name:     "with spaces",
			input:    "PROJ- 123",
			expected: " 123",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ExtractTicketNumber(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}


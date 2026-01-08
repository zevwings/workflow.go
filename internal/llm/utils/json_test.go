package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// ==================== ExtractJSONFromMarkdown 测试 ====================

func TestExtractJSONFromMarkdown(t *testing.T) {
	tests := []struct {
		name     string
		response string
		want     string
	}{
		{
			name:     "纯 JSON 字符串",
			response: `{"key": "value"}`,
			want:     `{"key": "value"}`,
		},
		{
			name:     "带 json 代码块的 markdown",
			response: "```json\n{\"key\": \"value\"}\n```",
			want:     `{"key": "value"}`,
		},
		{
			name:     "带普通代码块的 markdown",
			response: "```\n{\"key\": \"value\"}\n```",
			want:     `{"key": "value"}`,
		},
		{
			name:     "带空格的纯 JSON",
			response: "  {\"key\": \"value\"}  ",
			want:     `{"key": "value"}`,
		},
		{
			name:     "带空格的 markdown 代码块",
			response: "  ```json\n{\"key\": \"value\"}\n```  ",
			want:     `{"key": "value"}`,
		},
		{
			name:     "代码块没有换行（边界情况）",
			response: "```json{\"key\": \"value\"}```",
			want:     `{"key": "value"}`,
		},
		{
			name:     "代码块只有开始标记（边界情况）",
			response: "```json\n{\"key\": \"value\"}",
			want:     `{"key": "value"}`,
		},
		{
			name:     "空字符串",
			response: "",
			want:     "",
		},
		{
			name:     "只有空白字符",
			response: "   \n\t  ",
			want:     "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ExtractJSONFromMarkdown(tt.response)
			assert.Equal(t, tt.want, got)
		})
	}
}

// ==================== FixJSONEscapes 测试 ====================

func TestFixJSONEscapes(t *testing.T) {
	tests := []struct {
		name    string
		jsonStr string
		want    string
	}{
		{
			name:    "正常 JSON，无需修复",
			jsonStr: `{"path": "C:\\Users\\test"}`,
			want:    `{"path": "C:\\Users\\test"}`,
		},
		{
			name:    "Windows 路径未转义",
			jsonStr: `{"path": "C:\Users\test"}`,
			want:    `{"path": "C:\\Users\test"}`,
		},
		{
			name:    "包含有效转义序列",
			jsonStr: `{"text": "Hello\nWorld"}`,
			want:    `{"text": "Hello\nWorld"}`,
		},
		{
			name:    "包含无效转义序列",
			jsonStr: `{"text": "Hello\sWorld"}`,
			want:    `{"text": "Hello\\sWorld"}`,
		},
		{
			name:    "字符串末尾的反斜杠（实际是转义引号）",
			jsonStr: `{"path": "C:\Users\"}`,
			want:    `{"path": "C:\\Users\"}`,
		},
		{
			name:    "多个无效转义（\t 是有效转义）",
			jsonStr: `{"path": "C:\Users\test\file.txt"}`,
			want:    `{"path": "C:\\Users\test\file.txt"}`,
		},
		{
			name:    "混合有效和无效转义",
			jsonStr: `{"text": "Line1\nLine2\sLine3"}`,
			want:    `{"text": "Line1\nLine2\\sLine3"}`,
		},
		{
			name:    "空字符串",
			jsonStr: "",
			want:    "",
		},
		{
			name:    "不包含反斜杠",
			jsonStr: `{"key": "value"}`,
			want:    `{"key": "value"}`,
		},
		{
			name:    "转义的引号",
			jsonStr: `{"text": "He said \"Hello\""}`,
			want:    `{"text": "He said \"Hello\""}`,
		},
		{
			name:    "转义的反斜杠",
			jsonStr: `{"text": "Path: \\\\server\\share"}`,
			want:    `{"text": "Path: \\\\server\\share"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := FixJSONEscapes(tt.jsonStr)
			assert.Equal(t, tt.want, got)
		})
	}
}

// ==================== ExtractAndFixJSON 测试 ====================

func TestExtractAndFixJSON(t *testing.T) {
	tests := []struct {
		name     string
		response string
		want     string
	}{
		{
			name:     "从 markdown 提取并修复",
			response: "```json\n{\"path\": \"C:\\Users\\test\"}\n```",
			want:     `{"path": "C:\\Users\test"}`,
		},
		{
			name:     "纯 JSON 需要修复",
			response: `{"path": "C:\Users\test"}`,
			want:     `{"path": "C:\\Users\test"}`,
		},
		{
			name:     "正常 JSON 无需修复",
			response: `{"key": "value"}`,
			want:     `{"key": "value"}`,
		},
		{
			name:     "带 markdown 和需要修复的 JSON",
			response: "```\n{\"path\": \"C:\\Users\\test\", \"text\": \"Hello\\sWorld\"}\n```",
			want:     `{"path": "C:\\Users\test", "text": "Hello\\sWorld"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ExtractAndFixJSON(tt.response)
			assert.Equal(t, tt.want, got)
		})
	}
}

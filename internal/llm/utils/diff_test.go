package utils

import (
	"strings"
	"testing"
	"unicode/utf8"

	"github.com/stretchr/testify/assert"
)

// ==================== TruncateDiff 测试 ====================

func TestTruncateDiff(t *testing.T) {
	tests := []struct {
		name          string
		diff          string
		maxLength     int
		want          string
		wantTruncated bool
	}{
		{
			name:          "不需要截断",
			diff:          "line1\nline2\nline3",
			maxLength:     100,
			want:          "line1\nline2\nline3",
			wantTruncated: false,
		},
		{
			name:          "需要截断",
			diff:          strings.Repeat("line\n", 100),
			maxLength:     20,
			want:          strings.Repeat("line\n", 4) + "\n... (diff truncated, 500 characters total)",
			wantTruncated: true,
		},
		{
			name:          "空字符串",
			diff:          "",
			maxLength:     10,
			want:          "",
			wantTruncated: false,
		},
		{
			name:          "单行超过长度",
			diff:          strings.Repeat("a", 100),
			maxLength:     10,
			want:          strings.Repeat("a", 10) + "\n... (diff truncated, 100 characters total)",
			wantTruncated: true,
		},
		{
			name:          "在换行符处截断",
			diff:          "line1\nline2\nline3\nline4\nline5",
			maxLength:     15,
			want:          "line1\nline2\n... (diff truncated, 35 characters total)",
			wantTruncated: true,
		},
		{
			name:          "没有换行符",
			diff:          strings.Repeat("a", 100),
			maxLength:     20,
			want:          strings.Repeat("a", 20) + "\n... (diff truncated, 100 characters total)",
			wantTruncated: true,
		},
		{
			name:          "UTF-8 字符",
			diff:          "中文测试\n" + strings.Repeat("a", 50),
			maxLength:     10,
			want:          "中文测试\n... (diff truncated, 58 characters total)",
			wantTruncated: true,
		},
		{
			name:          "精确长度",
			diff:          "line1\nline2",
			maxLength:     12,
			want:          "line1\nline2",
			wantTruncated: false,
		},
		{
			name:          "maxLength 为 0（边界情况）",
			diff:          "line1\nline2",
			maxLength:     0,
			want:          "\n... (diff truncated, 12 characters total)",
			wantTruncated: true,
		},
		{
			name:          "多行 diff 格式",
			diff:          "diff --git a/file.go b/file.go\n@@ -1,3 +1,4 @@\n line1\n+line2\n line3",
			maxLength:     30,
			want:          "diff --git a/file.go b/file.go\n@@ -1,3 +1,4 @@\n... (diff truncated, 58 characters total)",
			wantTruncated: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := TruncateDiff(tt.diff, tt.maxLength)

			// 检查是否包含截断提示
			if tt.wantTruncated {
				assert.Contains(t, got, "... (diff truncated")
			} else {
				assert.NotContains(t, got, "... (diff truncated")
			}

			// 检查长度（考虑截断提示）
			if tt.wantTruncated {
				// 截断后的内容应该小于等于 maxLength + 提示信息长度
				truncatedPart := strings.Split(got, "\n... (diff truncated")[0]
				truncatedLength := utf8.RuneCountInString(truncatedPart)
				if tt.maxLength > 0 {
					assert.LessOrEqual(t, truncatedLength, tt.maxLength)
				}
			} else {
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

// ==================== TruncateDiff UTF-8 边界测试 ====================

func TestTruncateDiff_UTF8Boundary(t *testing.T) {
	// 测试 UTF-8 字符边界处理
	chineseText := "中文测试"
	diff := chineseText + "\n" + strings.Repeat("a", 100)

	result := TruncateDiff(diff, 10)

	// 应该不会在 UTF-8 字符中间截断
	assert.NotContains(t, result, "\ufffd") // 不应该包含替换字符
}

// ==================== TruncateDiff 换行符处理测试 ====================

func TestTruncateDiff_NewlineHandling(t *testing.T) {
	// 测试在换行符处截断
	diff := "line1\nline2\nline3\nline4\nline5"
	result := TruncateDiff(diff, 15)

	// 应该在最后一个换行符处截断
	assert.Contains(t, result, "line2")
	assert.NotContains(t, result, "line3")
}

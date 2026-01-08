package helpers

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ==================== ParsePRID 测试 ====================

func TestParsePRID(t *testing.T) {
	tests := []struct {
		name    string
		prID    string
		want    string
		wantErr bool
	}{
		{
			name:    "纯数字",
			prID:    "123",
			want:    "123",
			wantErr: false,
		},
		{
			name:    "GitHub URL 格式",
			prID:    "https://github.com/owner/repo/pull/123",
			want:    "123",
			wantErr: false,
		},
		{
			name:    "GitHub URL 格式（带尾部斜杠）",
			prID:    "https://github.com/owner/repo/pull/123/",
			want:    "123",
			wantErr: false,
		},
		{
			name:    "短格式 owner/repo#123",
			prID:    "owner/repo#123",
			want:    "123",
			wantErr: false,
		},
		{
			name:    "短格式（带空格）",
			prID:    "owner/repo# 123 ",
			want:    "123",
			wantErr: false,
		},
		{
			name:    "无效格式 - 非数字",
			prID:    "abc",
			want:    "",
			wantErr: true,
		},
		{
			name:    "无效格式 - URL 但无数字",
			prID:    "https://github.com/owner/repo/pull/abc",
			want:    "",
			wantErr: true,
		},
		{
			name:    "无效格式 - 短格式但无数字",
			prID:    "owner/repo#abc",
			want:    "",
			wantErr: true,
		},
		{
			name:    "空字符串",
			prID:    "",
			want:    "",
			wantErr: true,
		},
		{
			name:    "只有空格",
			prID:    "   ",
			want:    "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParsePRID(tt.prID)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Empty(t, got)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

// ==================== ParsePRNumber 测试 ====================

func TestParsePRNumber(t *testing.T) {
	tests := []struct {
		name    string
		prID    string
		want    int
		wantErr bool
	}{
		{
			name:    "纯数字",
			prID:    "123",
			want:    123,
			wantErr: false,
		},
		{
			name:    "GitHub URL 格式",
			prID:    "https://github.com/owner/repo/pull/456",
			want:    456,
			wantErr: false,
		},
		{
			name:    "短格式",
			prID:    "owner/repo#789",
			want:    789,
			wantErr: false,
		},
		{
			name:    "无效格式",
			prID:    "abc",
			want:    0,
			wantErr: true,
		},
		{
			name:    "空字符串",
			prID:    "",
			want:    0,
			wantErr: true,
		},
		{
			name:    "ParsePRID 返回错误",
			prID:    "invalid-format",
			want:    0,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParsePRNumber(tt.prID)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Equal(t, 0, got)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

// ==================== 边界情况测试 ====================

func TestParsePRID_EdgeCases(t *testing.T) {
	tests := []struct {
		name    string
		prID    string
		want    string
		wantErr bool
	}{
		{
			name:    "大数字",
			prID:    "999999999",
			want:    "999999999",
			wantErr: false,
		},
		{
			name:    "单个数字",
			prID:    "1",
			want:    "1",
			wantErr: false,
		},
		{
			name:    "URL 包含多个 /pull/",
			prID:    "https://github.com/owner/repo/pull/123/pull/456",
			want:    "",
			wantErr: true, // 实际解析逻辑只处理第一个 /pull/，但后续部分不是纯数字
		},
		{
			name:    "短格式包含多个 #",
			prID:    "owner/repo#123#456",
			want:    "",
			wantErr: true, // 实际解析逻辑只处理第一个 #，但后续部分不是纯数字
		},
		{
			name:    "URL 格式但路径不完整",
			prID:    "https://github.com/owner/repo/pull/",
			want:    "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParsePRID(tt.prID)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

// ==================== 集成测试 ====================

func TestParsePRID_ThenParsePRNumber(t *testing.T) {
	// 测试 ParsePRID 和 ParsePRNumber 的集成使用
	testCases := []string{
		"123",
		"https://github.com/owner/repo/pull/456",
		"owner/repo#789",
	}

	for _, prID := range testCases {
		t.Run(prID, func(t *testing.T) {
			// 先解析为字符串
			idStr, err := ParsePRID(prID)
			require.NoError(t, err)

			// 再解析为数字
			idNum, err := ParsePRNumber(idStr)
			require.NoError(t, err)

			// 验证结果一致
			idNumFromOriginal, err := ParsePRNumber(prID)
			require.NoError(t, err)
			assert.Equal(t, idNumFromOriginal, idNum)
		})
	}
}

package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// ==================== SanitizeBranchName 测试 ====================

func TestSanitizeBranchName(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "正常分支名",
			input: "feature/add-user-login",
			want:  "featureadd-user-login",
		},
		{
			name:  "包含中文",
			input: "feature/添加用户登录",
			want:  "feature",
		},
		{
			name:  "包含特殊字符",
			input: "feature/add@user#login",
			want:  "featureadduserlogin",
		},
		{
			name:  "包含数字",
			input: "bugfix/issue-123",
			want:  "bugfixissue-123",
		},
		{
			name:  "包含下划线",
			input: "feature/add_user_login",
			want:  "featureadd_user_login",
		},
		{
			name:  "包含空格",
			input: "feature/add user login",
			want:  "featureadduserlogin",
		},
		{
			name:  "空字符串",
			input: "",
			want:  "",
		},
		{
			name:  "只有特殊字符",
			input: "@#$%^&*()",
			want:  "",
		},
		{
			name:  "混合大小写字母",
			input: "Feature/AddUserLogin",
			want:  "FeatureAddUserLogin",
		},
		{
			name:  "包含多个连字符",
			input: "feature---add---user---login",
			want:  "feature---add---user---login",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := SanitizeBranchName(tt.input)
			assert.Equal(t, tt.want, got)
		})
	}
}

// ==================== CleanFilename 测试 ====================

func TestCleanFilename(t *testing.T) {
	tests := []struct {
		name     string
		filename string
		want     string
	}{
		{
			name:     "正常文件名",
			filename: "pr-summary",
			want:     "pr-summary",
		},
		{
			name:     "包含空格",
			filename: "PR Summary",
			want:     "pr-summary",
		},
		{
			name:     "包含 .md 扩展名",
			filename: "pr-summary.md",
			want:     "pr-summarymd",
		},
		{
			name:     "包含特殊字符",
			filename: "PR@Summary#2024",
			want:     "prsummary2024",
		},
		{
			name:     "包含中文",
			filename: "PR总结",
			want:     "pr",
		},
		{
			name:     "包含下划线",
			filename: "pr_summary",
			want:     "pr_summary",
		},
		{
			name:     "包含多个空格",
			filename: "PR  Summary  2024",
			want:     "pr--summary--2024",
		},
		{
			name:     "全大写",
			filename: "PR-SUMMARY",
			want:     "pr-summary",
		},
		{
			name:     "前后有空格",
			filename: "  pr-summary  ",
			want:     "pr-summary",
		},
		{
			name:     "包含 .md 和空格",
			filename: "PR Summary.md",
			want:     "pr-summarymd",
		},
		{
			name:     "空字符串",
			filename: "",
			want:     "",
		},
		{
			name:     "只有特殊字符",
			filename: "@#$%^&*()",
			want:     "",
		},
		{
			name:     "包含数字",
			filename: "PR Summary 2024",
			want:     "pr-summary-2024",
		},
		{
			name:     "多个 .md",
			filename: "pr-summary.md.md",
			want:     "pr-summarymdmd",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := CleanFilename(tt.filename)
			assert.Equal(t, tt.want, got)
		})
	}
}

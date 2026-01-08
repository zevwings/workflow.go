package helpers

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// ==================== ExtractRepoFromURL 测试 ====================

func TestExtractRepoFromURL(t *testing.T) {
	tests := []struct {
		name      string
		url       string
		wantOwner string
		wantRepo  string
		wantErr   bool
	}{
		{
			name:      "HTTPS 格式",
			url:       "https://github.com/owner/repo",
			wantOwner: "owner",
			wantRepo:  "repo",
			wantErr:   false,
		},
		{
			name:      "HTTPS 格式（带 .git）",
			url:       "https://github.com/owner/repo.git",
			wantOwner: "owner",
			wantRepo:  "repo",
			wantErr:   false,
		},
		{
			name:      "HTTPS 格式（带 www）",
			url:       "https://www.github.com/owner/repo.git",
			wantOwner: "owner",
			wantRepo:  "repo",
			wantErr:   false,
		},
		{
			name:      "SSH 格式",
			url:       "git@github.com:owner/repo.git",
			wantOwner: "owner",
			wantRepo:  "repo",
			wantErr:   false,
		},
		{
			name:      "SSH 格式（不带 .git）",
			url:       "git@github.com:owner/repo",
			wantOwner: "owner",
			wantRepo:  "repo",
			wantErr:   false,
		},
		{
			name:      "SSH 格式（完整路径）",
			url:       "ssh://git@github.com/owner/repo.git",
			wantOwner: "owner",
			wantRepo:  "repo",
			wantErr:   false,
		},
		{
			name:      "GitLab HTTPS 格式",
			url:       "https://gitlab.com/owner/repo.git",
			wantOwner: "owner",
			wantRepo:  "repo",
			wantErr:   false,
		},
		{
			name:      "GitLab SSH 格式",
			url:       "git@gitlab.com:owner/repo.git",
			wantOwner: "owner",
			wantRepo:  "repo",
			wantErr:   false,
		},
		{
			name:      "无效格式 - 无协议",
			url:       "github.com/owner/repo",
			wantOwner: "",
			wantRepo:  "",
			wantErr:   true,
		},
		{
			name:      "无效格式 - 路径不完整",
			url:       "https://github.com/owner",
			wantOwner: "",
			wantRepo:  "",
			wantErr:   true,
		},
		{
			name:      "无效格式 - 空字符串",
			url:       "",
			wantOwner: "",
			wantRepo:  "",
			wantErr:   true,
		},
		{
			name:      "无效格式 - 只有域名",
			url:       "https://github.com",
			wantOwner: "",
			wantRepo:  "",
			wantErr:   true,
		},
		{
			name:      "多级路径（取前两级）",
			url:       "https://github.com/owner/repo/subdir",
			wantOwner: "owner",
			wantRepo:  "repo",
			wantErr:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			owner, repo, err := ExtractRepoFromURL(tt.url)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Empty(t, owner)
				assert.Empty(t, repo)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.wantOwner, owner)
				assert.Equal(t, tt.wantRepo, repo)
			}
		})
	}
}

// ==================== BuildPRURL 测试 ====================

func TestBuildPRURL(t *testing.T) {
	tests := []struct {
		name     string
		baseURL  string
		owner    string
		repo     string
		prNumber int
		want     string
	}{
		{
			name:     "GitHub 标准格式",
			baseURL:  "https://github.com",
			owner:    "owner",
			repo:     "repo",
			prNumber: 123,
			want:     "https://github.com/owner/repo/pull/123",
		},
		{
			name:     "GitHub 带尾部斜杠",
			baseURL:  "https://github.com/",
			owner:    "owner",
			repo:     "repo",
			prNumber: 456,
			want:     "https://github.com/owner/repo/pull/456",
		},
		{
			name:     "GitLab 格式",
			baseURL:  "https://gitlab.com",
			owner:    "owner",
			repo:     "repo",
			prNumber: 789,
			want:     "https://gitlab.com/owner/repo/pull/789",
		},
		{
			name:     "自定义域名",
			baseURL:  "https://git.example.com",
			owner:    "owner",
			repo:     "repo",
			prNumber: 999,
			want:     "https://git.example.com/owner/repo/pull/999",
		},
		{
			name:     "PR 编号为 1",
			baseURL:  "https://github.com",
			owner:    "owner",
			repo:     "repo",
			prNumber: 1,
			want:     "https://github.com/owner/repo/pull/1",
		},
		{
			name:     "大 PR 编号",
			baseURL:  "https://github.com",
			owner:    "owner",
			repo:     "repo",
			prNumber: 999999,
			want:     "https://github.com/owner/repo/pull/999999",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := BuildPRURL(tt.baseURL, tt.owner, tt.repo, tt.prNumber)
			assert.Equal(t, tt.want, got)
		})
	}
}

// ==================== 集成测试 ====================

func TestExtractRepoFromURL_ThenBuildPRURL(t *testing.T) {
	// 测试 ExtractRepoFromURL 和 BuildPRURL 的集成使用
	testCases := []struct {
		name     string
		url      string
		baseURL  string
		prNumber int
		wantURL  string
		wantErr  bool
	}{
		{
			name:     "GitHub HTTPS",
			url:      "https://github.com/owner/repo.git",
			baseURL:  "https://github.com",
			prNumber: 123,
			wantURL:  "https://github.com/owner/repo/pull/123",
			wantErr:  false,
		},
		{
			name:     "GitHub SSH",
			url:      "git@github.com:owner/repo.git",
			baseURL:  "https://github.com",
			prNumber: 456,
			wantURL:  "https://github.com/owner/repo/pull/456",
			wantErr:  false,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			owner, repo, err := ExtractRepoFromURL(tt.url)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			prURL := BuildPRURL(tt.baseURL, owner, repo, tt.prNumber)
			assert.Equal(t, tt.wantURL, prURL)
		})
	}
}


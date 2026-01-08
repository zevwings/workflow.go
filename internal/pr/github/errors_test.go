package github

import (
	"errors"
	"net/http"
	"testing"

	"github.com/google/go-github/v57/github"
	"github.com/stretchr/testify/assert"
)

// ==================== IsNotFoundError 测试 ====================

func TestIsNotFoundError(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want bool
	}{
		{
			name: "404 GitHub 错误",
			err: &github.ErrorResponse{
				Response: &http.Response{
					StatusCode: http.StatusNotFound,
				},
			},
			want: true,
		},
		{
			name: "非 404 GitHub 错误",
			err: &github.ErrorResponse{
				Response: &http.Response{
					StatusCode: http.StatusInternalServerError,
				},
			},
			want: false,
		},
		{
			name: "普通错误",
			err:  errors.New("some error"),
			want: false,
		},
		{
			name: "nil 错误",
			err:  nil,
			want: false,
		},
		{
			name: "GitHub 错误但 Response 为 nil",
			err: &github.ErrorResponse{
				Response: nil,
			},
			want: false,
		},
		{
			name: "401 错误（不是 404）",
			err: &github.ErrorResponse{
				Response: &http.Response{
					StatusCode: http.StatusUnauthorized,
				},
			},
			want: false,
		},
		{
			name: "403 错误（不是 404）",
			err: &github.ErrorResponse{
				Response: &http.Response{
					StatusCode: http.StatusForbidden,
				},
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := IsNotFoundError(tt.err)
			assert.Equal(t, tt.want, got)
		})
	}
}

// ==================== IsUnauthorizedError 测试 ====================

func TestIsUnauthorizedError(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want bool
	}{
		{
			name: "401 GitHub 错误",
			err: &github.ErrorResponse{
				Response: &http.Response{
					StatusCode: http.StatusUnauthorized,
				},
			},
			want: true,
		},
		{
			name: "非 401 GitHub 错误",
			err: &github.ErrorResponse{
				Response: &http.Response{
					StatusCode: http.StatusInternalServerError,
				},
			},
			want: false,
		},
		{
			name: "普通错误",
			err:  errors.New("some error"),
			want: false,
		},
		{
			name: "nil 错误",
			err:  nil,
			want: false,
		},
		{
			name: "GitHub 错误但 Response 为 nil",
			err: &github.ErrorResponse{
				Response: nil,
			},
			want: false,
		},
		{
			name: "404 错误（不是 401）",
			err: &github.ErrorResponse{
				Response: &http.Response{
					StatusCode: http.StatusNotFound,
				},
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := IsUnauthorizedError(tt.err)
			assert.Equal(t, tt.want, got)
		})
	}
}

// ==================== IsForbiddenError 测试 ====================

func TestIsForbiddenError(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want bool
	}{
		{
			name: "403 GitHub 错误",
			err: &github.ErrorResponse{
				Response: &http.Response{
					StatusCode: http.StatusForbidden,
				},
			},
			want: true,
		},
		{
			name: "非 403 GitHub 错误",
			err: &github.ErrorResponse{
				Response: &http.Response{
					StatusCode: http.StatusInternalServerError,
				},
			},
			want: false,
		},
		{
			name: "普通错误",
			err:  errors.New("some error"),
			want: false,
		},
		{
			name: "nil 错误",
			err:  nil,
			want: false,
		},
		{
			name: "GitHub 错误但 Response 为 nil",
			err: &github.ErrorResponse{
				Response: nil,
			},
			want: false,
		},
		{
			name: "404 错误（不是 403）",
			err: &github.ErrorResponse{
				Response: &http.Response{
					StatusCode: http.StatusNotFound,
				},
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := IsForbiddenError(tt.err)
			assert.Equal(t, tt.want, got)
		})
	}
}

// ==================== FormatError 测试 ====================

func TestFormatError(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want string
	}{
		{
			name: "GitHub 错误（有 Message）",
			err: &github.ErrorResponse{
				Message:          "Not Found",
				DocumentationURL: "https://docs.github.com/rest",
				Response: &http.Response{
					StatusCode: http.StatusNotFound,
				},
			},
			want: "Not Found: https://docs.github.com/rest",
		},
		{
			name: "GitHub 错误（无 Message，有 Response）",
			err: &github.ErrorResponse{
				Response: &http.Response{
					StatusCode: http.StatusInternalServerError,
					Status:     "500 Internal Server Error",
				},
			},
			want: "GitHub API error: 500 Internal Server Error",
		},
		{
			name: "普通错误",
			err:  errors.New("some error message"),
			want: "some error message",
		},
		{
			name: "nil 错误",
			err:  nil,
			want: "",
		},
		{
			name: "GitHub 错误（无 Message 和 Response）",
			err:  &github.ErrorResponse{},
			want: "GitHub API error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := FormatError(tt.err)
			assert.Equal(t, tt.want, got)
		})
	}
}

// ==================== 集成测试 ====================

func TestErrorHandling_Integration(t *testing.T) {
	// 测试错误处理函数的组合使用
	t.Run("404 错误处理流程", func(t *testing.T) {
		err := &github.ErrorResponse{
			Message:          "Not Found",
			DocumentationURL: "https://docs.github.com/rest",
			Response: &http.Response{
				StatusCode: http.StatusNotFound,
			},
		}

		// 验证错误类型检测
		assert.True(t, IsNotFoundError(err))
		assert.False(t, IsUnauthorizedError(err))
		assert.False(t, IsForbiddenError(err))

		// 验证错误格式化
		formatted := FormatError(err)
		assert.Contains(t, formatted, "Not Found")
		assert.Contains(t, formatted, "docs.github.com")
	})

	t.Run("401 错误处理流程", func(t *testing.T) {
		err := &github.ErrorResponse{
			Message:          "Bad credentials",
			DocumentationURL: "https://docs.github.com/rest",
			Response: &http.Response{
				StatusCode: http.StatusUnauthorized,
			},
		}

		// 验证错误类型检测
		assert.False(t, IsNotFoundError(err))
		assert.True(t, IsUnauthorizedError(err))
		assert.False(t, IsForbiddenError(err))

		// 验证错误格式化
		formatted := FormatError(err)
		assert.Contains(t, formatted, "Bad credentials")
	})

	t.Run("403 错误处理流程", func(t *testing.T) {
		err := &github.ErrorResponse{
			Message:          "API rate limit exceeded",
			DocumentationURL: "https://docs.github.com/rest",
			Response: &http.Response{
				StatusCode: http.StatusForbidden,
			},
		}

		// 验证错误类型检测
		assert.False(t, IsNotFoundError(err))
		assert.False(t, IsUnauthorizedError(err))
		assert.True(t, IsForbiddenError(err))

		// 验证错误格式化
		formatted := FormatError(err)
		assert.Contains(t, formatted, "API rate limit exceeded")
	})
}


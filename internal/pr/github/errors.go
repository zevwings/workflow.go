package github

import (
	"fmt"
	"net/http"

	"github.com/google/go-github/v57/github"
)

// IsNotFoundError 检查是否是 404 错误
func IsNotFoundError(err error) bool {
	if err == nil {
		return false
	}

	ghErr, ok := err.(*github.ErrorResponse)
	if ok && ghErr.Response != nil {
		return ghErr.Response.StatusCode == http.StatusNotFound
	}

	return false
}

// IsUnauthorizedError 检查是否是 401 错误
func IsUnauthorizedError(err error) bool {
	if err == nil {
		return false
	}

	ghErr, ok := err.(*github.ErrorResponse)
	if ok && ghErr.Response != nil {
		return ghErr.Response.StatusCode == http.StatusUnauthorized
	}

	return false
}

// IsForbiddenError 检查是否是 403 错误
func IsForbiddenError(err error) bool {
	if err == nil {
		return false
	}

	ghErr, ok := err.(*github.ErrorResponse)
	if ok && ghErr.Response != nil {
		return ghErr.Response.StatusCode == http.StatusForbidden
	}

	return false
}

// FormatError 格式化错误信息
func FormatError(err error) string {
	if err == nil {
		return ""
	}

	ghErr, ok := err.(*github.ErrorResponse)
	if ok {
		if ghErr.Message != "" {
			return fmt.Sprintf("%s: %s", ghErr.Message, ghErr.DocumentationURL)
		}
		return fmt.Sprintf("GitHub API error: %s", ghErr.Response.Status)
	}

	return err.Error()
}

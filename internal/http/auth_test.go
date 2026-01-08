package http

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// ==================== Authorization 测试 ====================

// TestNewAuthorization 测试创建认证信息
func TestNewAuthorization(t *testing.T) {
	auth := NewAuthorization("user", "pass")
	assert.NotNil(t, auth)
	assert.Equal(t, "user", auth.Username)
	assert.Equal(t, "pass", auth.Password)
}


package http

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zevwings/workflow/internal/testutils"
)

// ==================== HttpMethod 测试 ====================

// TestParseHttpMethod 测试 HTTP 方法解析
func TestParseHttpMethod(t *testing.T) {
	testCases := []testutils.TableTestCase{
		{
			Name:  "GET",
			Input: "GET",
			Want:  MethodGet,
		},
		{
			Name:  "POST",
			Input: "POST",
			Want:  MethodPost,
		},
		{
			Name:  "PUT",
			Input: "PUT",
			Want:  MethodPut,
		},
		{
			Name:  "DELETE",
			Input: "DELETE",
			Want:  MethodDelete,
		},
		{
			Name:  "PATCH",
			Input: "PATCH",
			Want:  MethodPatch,
		},
		{
			Name:  "get (小写)",
			Input: "get",
			Want:  MethodGet, // 不区分大小写
		},
		{
			Name:  "post (小写)",
			Input: "post",
			Want:  MethodPost,
		},
		{
			Name:    "INVALID",
			Input:   "INVALID",
			WantErr: true,
		},
		{
			Name:    "空字符串",
			Input:   "",
			WantErr: true,
		},
	}

	testutils.RunTableTest(t, testCases, func(input interface{}) (interface{}, error) {
		return ParseHttpMethod(input.(string))
	})
}

// TestHttpMethod_String 测试 HTTP 方法字符串转换
func TestHttpMethod_String(t *testing.T) {
	assert.Equal(t, "GET", MethodGet.String())
	assert.Equal(t, "POST", MethodPost.String())
	assert.Equal(t, "PUT", MethodPut.String())
	assert.Equal(t, "DELETE", MethodDelete.String())
	assert.Equal(t, "PATCH", MethodPatch.String())
}

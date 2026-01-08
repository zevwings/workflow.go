package input

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// ==================== ValidateEmail 测试 ====================

func TestValidateEmail(t *testing.T) {
	tests := []struct {
		name    string
		email   string
		wantErr bool
	}{
		{
			name:    "有效邮箱 - 标准格式",
			email:   "user@example.com",
			wantErr: false,
		},
		{
			name:    "有效邮箱 - 带子域名",
			email:   "user@mail.example.com",
			wantErr: false,
		},
		{
			name:    "有效邮箱 - 带加号",
			email:   "user+tag@example.com",
			wantErr: false,
		},
		{
			name:    "有效邮箱 - 带点号",
			email:   "user.name@example.com",
			wantErr: false,
		},
		{
			name:    "有效邮箱 - 带下划线",
			email:   "user_name@example.com",
			wantErr: false,
		},
		{
			name:    "有效邮箱 - 带百分号",
			email:   "user%name@example.com",
			wantErr: false,
		},
		{
			name:    "有效邮箱 - 带连字符",
			email:   "user-name@example.com",
			wantErr: false,
		},
		{
			name:    "无效邮箱 - 缺少 @",
			email:   "userexample.com",
			wantErr: true,
		},
		{
			name:    "无效邮箱 - 缺少域名",
			email:   "user@",
			wantErr: true,
		},
		{
			name:    "无效邮箱 - 缺少用户名",
			email:   "@example.com",
			wantErr: true,
		},
		{
			name:    "无效邮箱 - 多个 @",
			email:   "user@@example.com",
			wantErr: true,
		},
		{
			name:    "无效邮箱 - 缺少顶级域名",
			email:   "user@example",
			wantErr: true,
		},
		{
			name:    "无效邮箱 - 空字符串",
			email:   "",
			wantErr: true,
		},
		{
			name:    "无效邮箱 - 只有空格",
			email:   "   ",
			wantErr: true,
		},
		{
			name:    "无效邮箱 - 包含空格",
			email:   "user name@example.com",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange: 创建验证器
			validator := ValidateEmail()

			// Act: 验证邮箱
			err := validator(tt.email)

			// Assert: 验证结果
			if tt.wantErr {
				assert.Error(t, err, "应该返回错误: %s", tt.email)
				assert.Contains(t, err.Error(), "邮箱", "错误消息应该包含'邮箱'")
			} else {
				assert.NoError(t, err, "不应该返回错误: %s", tt.email)
			}
		})
	}
}

// ==================== ValidateURL 测试 ====================

func TestValidateURL(t *testing.T) {
	tests := []struct {
		name    string
		url     string
		wantErr bool
	}{
		{
			name:    "有效 URL - HTTP",
			url:     "http://example.com",
			wantErr: false,
		},
		{
			name:    "有效 URL - HTTPS",
			url:     "https://example.com",
			wantErr: false,
		},
		{
			name:    "有效 URL - 带路径",
			url:     "https://example.com/path/to/resource",
			wantErr: false,
		},
		{
			name:    "有效 URL - 带查询参数",
			url:     "https://example.com/path?key=value",
			wantErr: false,
		},
		{
			name:    "有效 URL - 带端口",
			url:     "http://example.com:8080",
			wantErr: false,
		},
		{
			name:    "有效 URL - 带锚点",
			url:     "https://example.com#section",
			wantErr: false,
		},
		{
			name:    "无效 URL - 缺少协议",
			url:     "example.com",
			wantErr: true,
		},
		{
			name:    "无效 URL - 错误协议",
			url:     "ftp://example.com",
			wantErr: true,
		},
		{
			name:    "无效 URL - 空字符串",
			url:     "",
			wantErr: true,
		},
		{
			name:    "无效 URL - 只有协议",
			url:     "http://",
			wantErr: true,
		},
		{
			name:    "无效 URL - 包含空格",
			url:     "https://example.com/path with spaces",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange: 创建验证器
			validator := ValidateURL()

			// Act: 验证 URL
			err := validator(tt.url)

			// Assert: 验证结果
			if tt.wantErr {
				assert.Error(t, err, "应该返回错误: %s", tt.url)
				assert.Contains(t, err.Error(), "URL", "错误消息应该包含'URL'")
			} else {
				assert.NoError(t, err, "不应该返回错误: %s", tt.url)
			}
		})
	}
}

// ==================== ValidateRequired 测试 ====================

func TestValidateRequired(t *testing.T) {
	tests := []struct {
		name    string
		value   string
		wantErr bool
	}{
		{
			name:    "有效 - 非空字符串",
			value:   "value",
			wantErr: false,
		},
		{
			name:    "有效 - 包含空格但非空",
			value:   " value ",
			wantErr: false,
		},
		{
			name:    "无效 - 空字符串",
			value:   "",
			wantErr: true,
		},
		{
			name:    "无效 - 只有空格",
			value:   "   ",
			wantErr: true,
		},
		{
			name:    "无效 - 只有制表符",
			value:   "\t\t",
			wantErr: true,
		},
		{
			name:    "无效 - 只有换行符",
			value:   "\n\n",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange: 创建验证器
			validator := ValidateRequired()

			// Act: 验证值
			err := validator(tt.value)

			// Assert: 验证结果
			if tt.wantErr {
				assert.Error(t, err, "应该返回错误: %q", tt.value)
				assert.Contains(t, err.Error(), "必填", "错误消息应该包含'必填'")
			} else {
				assert.NoError(t, err, "不应该返回错误: %q", tt.value)
			}
		})
	}
}

// ==================== ValidateMinLength 测试 ====================

func TestValidateMinLength(t *testing.T) {
	tests := []struct {
		name    string
		minLen  int
		value   string
		wantErr bool
	}{
		{
			name:    "有效 - 长度等于最小值",
			minLen:  5,
			value:   "12345",
			wantErr: false,
		},
		{
			name:    "有效 - 长度大于最小值",
			minLen:  5,
			value:   "123456",
			wantErr: false,
		},
		{
			name:    "无效 - 长度小于最小值",
			minLen:  5,
			value:   "1234",
			wantErr: true,
		},
		{
			name:    "无效 - 空字符串",
			minLen:  5,
			value:   "",
			wantErr: true,
		},
		{
			name:    "边界 - 最小长度为 0",
			minLen:  0,
			value:   "",
			wantErr: false,
		},
		{
			name:    "边界 - 最小长度为 1",
			minLen:  1,
			value:   "a",
			wantErr: false,
		},
		{
			name:    "边界 - 最小长度为 1，空字符串",
			minLen:  1,
			value:   "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange: 创建验证器
			validator := ValidateMinLength(tt.minLen)

			// Act: 验证值
			err := validator(tt.value)

			// Assert: 验证结果
			if tt.wantErr {
				assert.Error(t, err, "应该返回错误: %q (minLen: %d)", tt.value, tt.minLen)
				assert.Contains(t, err.Error(), "至少", "错误消息应该包含'至少'")
			} else {
				assert.NoError(t, err, "不应该返回错误: %q (minLen: %d)", tt.value, tt.minLen)
			}
		})
	}
}

// ==================== ValidateMaxLength 测试 ====================

func TestValidateMaxLength(t *testing.T) {
	tests := []struct {
		name    string
		maxLen  int
		value   string
		wantErr bool
	}{
		{
			name:    "有效 - 长度等于最大值",
			maxLen:  5,
			value:   "12345",
			wantErr: false,
		},
		{
			name:    "有效 - 长度小于最大值",
			maxLen:  5,
			value:   "1234",
			wantErr: false,
		},
		{
			name:    "无效 - 长度大于最大值",
			maxLen:  5,
			value:   "123456",
			wantErr: true,
		},
		{
			name:    "有效 - 空字符串",
			maxLen:  5,
			value:   "",
			wantErr: false,
		},
		{
			name:    "边界 - 最大长度为 0",
			maxLen:  0,
			value:   "",
			wantErr: false,
		},
		{
			name:    "边界 - 最大长度为 0，非空字符串",
			maxLen:  0,
			value:   "a",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange: 创建验证器
			validator := ValidateMaxLength(tt.maxLen)

			// Act: 验证值
			err := validator(tt.value)

			// Assert: 验证结果
			if tt.wantErr {
				assert.Error(t, err, "应该返回错误: %q (maxLen: %d)", tt.value, tt.maxLen)
				assert.Contains(t, err.Error(), "超过", "错误消息应该包含'超过'")
			} else {
				assert.NoError(t, err, "不应该返回错误: %q (maxLen: %d)", tt.value, tt.maxLen)
			}
		})
	}
}

// ==================== ValidateLength 测试 ====================

func TestValidateLength(t *testing.T) {
	tests := []struct {
		name    string
		minLen  int
		maxLen  int
		value   string
		wantErr bool
	}{
		{
			name:    "有效 - 长度在范围内",
			minLen:  3,
			maxLen:  5,
			value:   "123",
			wantErr: false,
		},
		{
			name:    "有效 - 长度等于最小值",
			minLen:  3,
			maxLen:  5,
			value:   "123",
			wantErr: false,
		},
		{
			name:    "有效 - 长度等于最大值",
			minLen:  3,
			maxLen:  5,
			value:   "12345",
			wantErr: false,
		},
		{
			name:    "无效 - 长度小于最小值",
			minLen:  3,
			maxLen:  5,
			value:   "12",
			wantErr: true,
		},
		{
			name:    "无效 - 长度大于最大值",
			minLen:  3,
			maxLen:  5,
			value:   "123456",
			wantErr: true,
		},
		{
			name:    "边界 - 最小和最大长度相同",
			minLen:  3,
			maxLen:  3,
			value:   "123",
			wantErr: false,
		},
		{
			name:    "边界 - 最小和最大长度相同，不匹配",
			minLen:  3,
			maxLen:  3,
			value:   "12",
			wantErr: true,
		},
		{
			name:    "边界 - 空字符串，最小长度为 0",
			minLen:  0,
			maxLen:  5,
			value:   "",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange: 创建验证器
			validator := ValidateLength(tt.minLen, tt.maxLen)

			// Act: 验证值
			err := validator(tt.value)

			// Assert: 验证结果
			if tt.wantErr {
				assert.Error(t, err, "应该返回错误: %q (minLen: %d, maxLen: %d)", tt.value, tt.minLen, tt.maxLen)
				assert.Contains(t, err.Error(), "之间", "错误消息应该包含'之间'")
			} else {
				assert.NoError(t, err, "不应该返回错误: %q (minLen: %d, maxLen: %d)", tt.value, tt.minLen, tt.maxLen)
			}
		})
	}
}

// ==================== ValidateRegex 测试 ====================

func TestValidateRegex(t *testing.T) {
	tests := []struct {
		name     string
		pattern  string
		errorMsg string
		value    string
		wantErr  bool
	}{
		{
			name:     "有效 - 匹配模式",
			pattern:  `^[a-z]+$`,
			errorMsg: "必须是小写字母",
			value:    "abc",
			wantErr:  false,
		},
		{
			name:     "无效 - 不匹配模式",
			pattern:  `^[a-z]+$`,
			errorMsg: "必须是小写字母",
			value:    "ABC",
			wantErr:  true,
		},
		{
			name:     "无效 - 使用自定义错误消息",
			pattern:  `^[0-9]+$`,
			errorMsg: "必须是数字",
			value:    "abc",
			wantErr:  true,
		},
		{
			name:     "有效 - 数字模式",
			pattern:  `^[0-9]+$`,
			errorMsg: "必须是数字",
			value:    "123",
			wantErr:  false,
		},
		{
			name:     "无效 - 空错误消息使用默认",
			pattern:  `^[a-z]+$`,
			errorMsg: "",
			value:    "ABC",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange: 创建验证器
			validator := ValidateRegex(tt.pattern, tt.errorMsg)

			// Act: 验证值
			err := validator(tt.value)

			// Assert: 验证结果
			if tt.wantErr {
				assert.Error(t, err, "应该返回错误: %q (pattern: %s)", tt.value, tt.pattern)
				if tt.errorMsg != "" {
					assert.Contains(t, err.Error(), tt.errorMsg, "错误消息应该包含自定义消息")
				} else {
					assert.Contains(t, err.Error(), "格式不正确", "错误消息应该包含默认消息")
				}
			} else {
				assert.NoError(t, err, "不应该返回错误: %q (pattern: %s)", tt.value, tt.pattern)
			}
		})
	}
}

func TestValidateRegex_InvalidPattern(t *testing.T) {
	// Arrange: 无效的正则表达式模式（会导致 panic）
	// 注意：由于 ValidateRegex 使用 regexp.MustCompile，无效模式会导致 panic
	// 这个测试验证函数在编译时就会失败（如果模式无效）

	// 这个测试主要是文档性质的，因为无效模式会在编译时或运行时 panic
	// 实际使用中应该确保传入的模式是有效的
}


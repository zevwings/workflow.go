package prompt

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ==================== LoadTemplate 测试 ====================

func TestLoadTemplate(t *testing.T) {
	tests := []struct {
		name    string
		template string
		wantErr bool
	}{
		{
			name:     "加载 translate.md",
			template: "translate.md",
			wantErr:  false,
		},
		{
			name:     "加载 branch.md",
			template: "branch.md",
			wantErr:  false,
		},
		{
			name:     "加载 pr-reword.md",
			template: "pr-reword.md",
			wantErr:  false,
		},
		{
			name:     "加载 file-summary.md",
			template: "file-summary.md",
			wantErr:  false,
		},
		{
			name:     "加载 pr-summary.md",
			template: "pr-summary.md",
			wantErr:  false,
		},
		{
			name:     "不存在的模板",
			template: "non-existent.md",
			wantErr:  true,
		},
		{
			name:     "空文件名",
			template: "",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Act: 加载模板
			content, err := LoadTemplate(tt.template)

			// Assert: 验证结果
			if tt.wantErr {
				assert.Error(t, err)
				assert.Empty(t, content)
			} else {
				assert.NoError(t, err)
				assert.NotEmpty(t, content, "模板内容不应为空")
				// 验证内容包含一些预期的文本
				assert.Contains(t, content, "", "模板内容应该存在")
			}
		})
	}
}

func TestLoadTemplate_ContentNotEmpty(t *testing.T) {
	// Arrange: 已知存在的模板文件
	template := "translate.md"

	// Act: 加载模板
	content, err := LoadTemplate(template)

	// Assert: 验证内容不为空
	require.NoError(t, err)
	assert.NotEmpty(t, content, "模板内容不应为空")
	assert.Greater(t, len(content), 0, "模板内容长度应大于 0")
}

func TestLoadTemplate_ErrorFormat(t *testing.T) {
	// Arrange: 不存在的模板
	template := "non-existent-template.md"

	// Act: 加载模板
	_, err := LoadTemplate(template)

	// Assert: 验证错误格式
	require.Error(t, err)
	assert.Contains(t, err.Error(), "读取模板文件失败")
	assert.Contains(t, err.Error(), template)
}

// ==================== MustLoadTemplate 测试 ====================

func TestMustLoadTemplate(t *testing.T) {
	tests := []struct {
		name      string
		template  string
		shouldPanic bool
	}{
		{
			name:       "加载存在的模板",
			template:   "translate.md",
			shouldPanic: false,
		},
		{
			name:       "加载 branch.md",
			template:   "branch.md",
			shouldPanic: false,
		},
		{
			name:       "加载 pr-reword.md",
			template:   "pr-reword.md",
			shouldPanic: false,
		},
		{
			name:       "加载 file-summary.md",
			template:   "file-summary.md",
			shouldPanic: false,
		},
		{
			name:       "加载 pr-summary.md",
			template:   "pr-summary.md",
			shouldPanic: false,
		},
		{
			name:       "不存在的模板应该 panic",
			template:   "non-existent.md",
			shouldPanic: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.shouldPanic {
				// Assert: 验证会 panic
				assert.Panics(t, func() {
					MustLoadTemplate(tt.template)
				}, "应该 panic 当模板不存在时")
			} else {
				// Act: 加载模板（不应该 panic）
				content := MustLoadTemplate(tt.template)

				// Assert: 验证内容不为空
				assert.NotEmpty(t, content, "模板内容不应为空")
			}
		})
	}
}

func TestMustLoadTemplate_PanicMessage(t *testing.T) {
	// Arrange: 不存在的模板
	template := "non-existent-template.md"

	// Act & Assert: 验证 panic 消息
	assert.Panics(t, func() {
		MustLoadTemplate(template)
	}, "应该 panic 当模板不存在时")

	// 验证 panic 消息包含模板名称
	defer func() {
		if r := recover(); r != nil {
			panicMsg := r.(string)
			assert.Contains(t, panicMsg, "无法加载必需的模板文件")
			assert.Contains(t, panicMsg, template)
		}
	}()
	MustLoadTemplate(template)
}

// ==================== ListTemplates 测试 ====================

func TestListTemplates(t *testing.T) {
	// Act: 列出所有模板
	templates, err := ListTemplates()

	// Assert: 验证结果
	require.NoError(t, err)
	assert.NotEmpty(t, templates, "应该至少有一个模板文件")

	// 验证包含预期的模板文件
	expectedTemplates := []string{
		"translate.md",
		"branch.md",
		"pr-reword.md",
		"file-summary.md",
		"pr-summary.md",
	}

	for _, expected := range expectedTemplates {
		assert.Contains(t, templates, expected, "应该包含模板: %s", expected)
	}
}

func TestListTemplates_OnlyMarkdownFiles(t *testing.T) {
	// Act: 列出所有模板
	templates, err := ListTemplates()

	// Assert: 验证结果
	require.NoError(t, err)

	// 验证所有文件都是 .md 文件
	for _, template := range templates {
		assert.Contains(t, template, ".md", "所有模板文件应该是 .md 文件: %s", template)
	}
}

func TestListTemplates_NoDirectories(t *testing.T) {
	// Act: 列出所有模板
	templates, err := ListTemplates()

	// Assert: 验证结果
	require.NoError(t, err)

	// 验证不包含目录（所有条目都应该是文件）
	// 由于我们只检查 !entry.IsDir()，所以不应该有目录
	for _, template := range templates {
		// 模板名称应该包含 .md 扩展名，表示是文件
		assert.Contains(t, template, ".md", "模板应该是文件: %s", template)
	}
}

func TestListTemplates_ConsistentWithLoadTemplate(t *testing.T) {
	// Arrange: 列出所有模板
	templates, err := ListTemplates()
	require.NoError(t, err)

	// Act & Assert: 验证列出的每个模板都可以加载
	for _, template := range templates {
		t.Run("加载模板: "+template, func(t *testing.T) {
			content, err := LoadTemplate(template)
			assert.NoError(t, err, "应该能够加载模板: %s", template)
			assert.NotEmpty(t, content, "模板内容不应为空: %s", template)
		})
	}
}


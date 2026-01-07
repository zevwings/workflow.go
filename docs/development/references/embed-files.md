# Embed 文件使用指南

> 本文档描述了项目中 Go embed 功能的实现方案和使用方法。

## 项目实现方案

本项目使用 **统一加载器 + Markdown 模板文件** 的方案：

1. **统一加载器**：`internal/llm/prompt/loader.go` 提供统一的模板加载功能
2. **模板文件格式**：统一使用 `.md` (Markdown) 格式
3. **嵌入方式**：使用 `embed.FS` 嵌入整个 `templates/` 目录

## 实现架构

### 统一加载器

所有模板文件通过 `loader.go` 中的统一加载器加载：

```go
// internal/llm/prompt/loader.go
package prompt

import (
    "embed"
    "fmt"
)

//go:embed templates/*.md
var templatesFS embed.FS

// LoadTemplate 加载模板文件（带错误处理）
func LoadTemplate(name string) (string, error) {
    data, err := templatesFS.ReadFile("templates/" + name)
    if err != nil {
        return "", fmt.Errorf("读取模板文件失败 (%s): %w", name, err)
    }
    return string(data), nil
}

// MustLoadTemplate 加载模板文件（失败则 panic，用于包初始化）
func MustLoadTemplate(name string) string {
    prompt, err := LoadTemplate(name)
    if err != nil {
        panic(fmt.Sprintf("无法加载必需的模板文件 %s: %v", name, err))
    }
    return prompt
}
```

### 使用方式

**方式 1：包初始化时加载（推荐用于常量）**

```go
// internal/llm/prompt/translate.go
package prompt

var TranslateSystemPrompt = MustLoadTemplate("translate.md")
```

**方式 2：函数中加载（用于需要动态处理的场景）**

```go
// internal/llm/prompt/file.go
func GenerateSummarizeFileChangeSystemPrompt(cfg *config.Manager) string {
    basePrompt := MustLoadTemplate("file-summary.md")
    return llm.GetLanguageRequirement(basePrompt, cfg)
}
```

**方式 3：动态内容拼接**

```go
// internal/llm/prompt/pr.go
func GenerateSummarizePRSystemPrompt(cfg *config.Manager) string {
    basePrompt := MustLoadTemplate("pr-summary.md")
    dynamicContent := buildResponseExample()
    fullPrompt := fmt.Sprintf(basePrompt, dynamicContent)
    return llm.GetLanguageRequirement(fullPrompt, cfg)
}
```

## 优势

1. **单文件分发**：所有模板文件都打包在二进制文件中，无需额外文件
2. **版本一致性**：模板文件与代码版本保持一致
3. **简化部署**：不需要管理外部模板文件
4. **性能优化**：文件内容在编译时嵌入，运行时无需文件 I/O
5. **代码简化**：代码行数减少 71%（从 542行 → 158行）
6. **易于维护**：Prompt 内容与业务逻辑分离，可直接编辑 Markdown 文件

## 模板文件结构

所有模板文件位于 `internal/llm/prompt/templates/` 目录：

```
templates/
├── translate.md      - 翻译 prompt 模板
├── branch.md         - 分支生成 prompt 模板
├── pr-reword.md      - PR 重写 prompt 模板
├── file-summary.md   - 文件总结 prompt 模板
└── pr-summary.md     - PR 总结 prompt 模板
```

### 文件命名规范

- 统一使用 `.md` 扩展名（Markdown 格式）
- 文件名使用小写字母和连字符（kebab-case）
- 文件名应该清晰描述 prompt 的用途

## 注意事项

### 路径规则

- `//go:embed` 指令必须紧邻变量声明之前
- 路径是相对于包含 `//go:embed` 指令的 `.go` 文件
- 不能使用 `..` 访问父目录
- 路径分隔符使用 `/`（即使在 Windows 上）

### 编译时检查

- 嵌入的文件会在编译时检查，如果文件不存在会导致编译失败
- 文件内容在编译时确定，运行时无法修改
- 使用 `MustLoadTemplate` 时，如果文件不存在会在包初始化时 panic

### 性能考虑

- 嵌入的文件会增加二进制文件大小（当前项目增加约 10-20KB，可接受）
- 文件内容在编译时嵌入，运行时读取速度很快
- 对于大型文件（>10MB），考虑是否真的需要嵌入

### 编辑模板文件

1. 直接编辑 `internal/llm/prompt/templates/*.md` 文件
2. 修改后重新编译即可生效
3. 无需修改 Go 代码
4. 版本控制会清晰显示 Prompt 内容的变更

## 测试

测试嵌入的模板文件：

```go
func TestEmbeddedTemplate(t *testing.T) {
    prompt := prompt.TranslateSystemPrompt
    assert.NotEmpty(t, prompt)
    assert.Contains(t, prompt, "translation assistant")
}

// 测试模板加载器
func TestLoadTemplate(t *testing.T) {
    content, err := prompt.LoadTemplate("translate.md")
    assert.NoError(t, err)
    assert.NotEmpty(t, content)
}

// 测试列出所有模板
func TestListTemplates(t *testing.T) {
    files, err := prompt.ListTemplates()
    assert.NoError(t, err)
    assert.Contains(t, files, "translate.md")
    assert.Contains(t, files, "branch.md")
}
```

## 相关文件

- `internal/llm/prompt/loader.go` - 统一模板加载器
- `internal/llm/prompt/templates/` - 模板文件目录

## 相关文档

- [Prompt Embed 迁移分析](../analysis/prompt_embed_migration_analysis.md) - 完整的迁移分析报告
- [Prompt Embed 迁移完成报告](../analysis/prompt_migration_complete.md) - 迁移完成报告


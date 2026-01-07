# Prompt 模板文件

这个目录包含所有 LLM prompt 模板文件，这些文件会在编译时嵌入到二进制文件中。

## 文件列表

- `translate.md` - 翻译 prompt 模板
- `branch.md` - 分支生成 prompt 模板
- `pr-reword.md` - PR 重写 prompt 模板
- `file-summary.md` - 文件总结 prompt 模板
- `pr-summary.md` - PR 总结 prompt 模板

## 使用方法

所有模板文件通过 `loader.go` 中的统一加载器加载：

```go
// 使用 MustLoadTemplate 在包初始化时加载（推荐）
var TranslateSystemPrompt = MustLoadTemplate("translate.md")

// 或使用 LoadTemplate 在函数中加载（带错误处理）
func GetPrompt() (string, error) {
    return LoadTemplate("translate.md")
}
```

## 文件命名规范

- 统一使用 `.md` 扩展名（Markdown 格式）
- 文件名使用小写字母和连字符（kebab-case）
- 文件名应该清晰描述 prompt 的用途

## 优势

1. **单文件分发**：所有模板文件都打包在二进制文件中，无需额外文件
2. **版本一致性**：模板文件与代码版本保持一致
3. **简化部署**：不需要管理外部模板文件
4. **性能优化**：文件内容在编译时嵌入，运行时无需文件 I/O

## 注意事项

- `//go:embed` 指令必须紧邻变量声明之前
- 路径是相对于包含 `//go:embed` 指令的 `.go` 文件
- 不能使用 `..` 访问父目录
- 嵌入的文件会在编译时检查，如果文件不存在会导致编译失败


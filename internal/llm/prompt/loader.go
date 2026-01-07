package prompt

import (
	"embed"
	"fmt"
)

//go:embed templates/*.md
var templatesFS embed.FS

// LoadTemplate 从嵌入的文件系统中加载模板文件
//
// 参数:
//   - name: 模板文件名（不包含路径，如 "translate.md"）
//
// 返回:
//   - string: 模板内容
//   - error: 如果读取失败，返回错误
//
// 示例:
//   prompt, err := LoadTemplate("translate.md")
//   if err != nil {
//       log.Fatal(err)
//   }
func LoadTemplate(name string) (string, error) {
	// 注意：路径是相对于 embed 指令中指定的路径
	// 如果使用 templates/*.md，则读取时需要包含 templates/ 前缀
	data, err := templatesFS.ReadFile("templates/" + name)
	if err != nil {
		return "", fmt.Errorf("读取模板文件失败 (%s): %w", name, err)
	}
	return string(data), nil
}

// MustLoadTemplate 从嵌入的文件系统中加载模板文件，如果失败则 panic
//
// 这个方法用于在编译时确保模板文件存在，如果文件不存在会导致编译失败。
// 适用于在包初始化时加载模板。
//
// 参数:
//   - name: 模板文件名（不包含路径，如 "translate.md"）
//
// 返回:
//   - string: 模板内容
//
// 示例:
//   var translatePrompt = MustLoadTemplate("translate.md")
func MustLoadTemplate(name string) string {
	prompt, err := LoadTemplate(name)
	if err != nil {
		panic(fmt.Sprintf("无法加载必需的模板文件 %s: %v", name, err))
	}
	return prompt
}

// ListTemplates 列出所有可用的模板文件
//
// 返回:
//   - []string: 模板文件名列表
//   - error: 如果列出失败，返回错误
func ListTemplates() ([]string, error) {
	entries, err := templatesFS.ReadDir("templates")
	if err != nil {
		return nil, fmt.Errorf("列出模板文件失败: %w", err)
	}

	var files []string
	for _, entry := range entries {
		if !entry.IsDir() {
			files = append(files, entry.Name())
		}
	}

	return files, nil
}


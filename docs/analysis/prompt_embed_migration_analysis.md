# Prompt 目录 Embed 改造分析报告

## 1. 当前状态分析

### 1.1 文件结构

```
internal/llm/prompt/
├── branch.go          (90行)   - GenerateBranchSystemPrompt (常量)
├── file.go            (80行)   - GenerateSummarizeFileChangeSystemPrompt (函数，包含 basePrompt)
├── pr.go              (256行)  - RewordPRSystemPrompt (常量) + GenerateSummarizePRSystemPrompt (函数)
├── translate.go        (20行)   - TranslateSystemPrompt (常量)
├── embedded.go        (26行)   - 示例实现
├── usage_example.go    (72行)   - 使用示例
└── templates/         (已创建)
    ├── branch.txt
    ├── translate.txt
    └── README.md
```

### 1.2 Prompt 使用情况

| Prompt 名称 | 类型 | 行数 | 使用位置 | 复杂度 |
|------------|------|------|----------|--------|
| `GenerateBranchSystemPrompt` | 常量 | ~90行 | `internal/llm/pr/create.go` | 高 |
| `TranslateSystemPrompt` | 常量 | ~15行 | `internal/llm/branch/translate.go` | 低 |
| `RewordPRSystemPrompt` | 常量 | ~60行 | `internal/llm/pr/reword.go` | 中 |
| `GenerateSummarizeFileChangeSystemPrompt` | 函数 | ~60行 | `internal/llm/pr/file_summary.go` | 中 |
| `GenerateSummarizePRSystemPrompt` | 函数 | ~170行 | `internal/llm/pr/summary.go` | 高 |

### 1.3 代码统计

- **总代码行数**: 542行
- **Prompt 相关代码**: ~364行 (67%)
- **实际业务逻辑**: ~178行 (33%)

## 2. 改造可行性分析

### 2.1 ✅ 完全可行

所有 prompt 都可以迁移到嵌入文件：

1. **常量类型** (3个)
   - `GenerateBranchSystemPrompt` - 可直接迁移
   - `TranslateSystemPrompt` - 可直接迁移
   - `RewordPRSystemPrompt` - 可直接迁移

2. **函数类型** (2个)
   - `GenerateSummarizeFileChangeSystemPrompt` - 可迁移 basePrompt 部分
   - `GenerateSummarizePRSystemPrompt` - 可迁移 basePrompt 部分

### 2.2 改造优势

#### 代码可维护性
- ✅ **代码行数减少**: 从 542行 → ~200行 (减少 63%)
- ✅ **关注点分离**: Prompt 内容与业务逻辑分离
- ✅ **易于编辑**: 直接在文本文件中编辑，无需修改 Go 代码
- ✅ **版本控制友好**: Prompt 变更更容易追踪和审查

#### 开发体验
- ✅ **语法高亮**: 文本编辑器对 Markdown/TXT 支持更好
- ✅ **无转义问题**: 不需要处理 Go 字符串转义
- ✅ **易于测试**: 可以单独测试 prompt 内容
- ✅ **协作友好**: 非 Go 开发者也可以编辑 prompt

#### 运行时优势
- ✅ **单文件分发**: 所有 prompt 打包在二进制中
- ✅ **版本一致性**: Prompt 与代码版本同步
- ✅ **性能**: 编译时嵌入，运行时无文件 I/O

### 2.3 潜在风险

#### 低风险
- ⚠️ **编译时检查**: 文件不存在会导致编译失败（这是好事，可以提前发现问题）
- ⚠️ **二进制大小**: 会增加二进制文件大小（预计增加 ~10-20KB，可接受）

#### 需要处理
- ⚠️ **动态内容**: `GenerateSummarizePRSystemPrompt` 中有 `fmt.Sprintf` 动态内容
  - **解决方案**: 保留函数，但 basePrompt 从文件读取，动态部分在函数中拼接

## 3. Markdown vs TXT 格式分析

### 3.1 当前 Prompt 内容特征

分析现有 prompt 的内容特征：

1. **包含 Markdown 语法**:
   - 标题: `## Important`, `### Branch Name Rules`
   - 粗体: `**All outputs MUST be in English only.**`
   - 代码块: ` ```json ... ``` `
   - 列表: `- Item 1`, `- Item 2`
   - 表格: `| Input | Output |`

2. **结构复杂**:
   - 多级标题
   - 嵌套列表
   - 代码示例
   - JSON 示例

### 3.2 Markdown 格式优势

#### ✅ 推荐使用 Markdown (.md)

**理由 1: 内容本身已经是 Markdown**
- 现有 prompt 大量使用 Markdown 语法
- 使用 `.md` 文件可以保持格式一致性
- 编辑器可以正确渲染和预览

**理由 2: 可读性和维护性**
```markdown
## Important

**All outputs MUST be in English only.**

### Branch Name Rules

- Must be all lowercase
- Use hyphens to separate words
```

vs

```
## Important\n\n**All outputs MUST be in English only.**\n\n### Branch Name Rules\n\n- Must be all lowercase
```

**理由 3: 工具支持**
- GitHub/GitLab 可以直接预览
- 编辑器（VS Code, Cursor）有语法高亮
- 可以生成文档网站

**理由 4: 未来扩展**
- 如果需要生成 prompt 文档，Markdown 更合适
- 如果需要 AI 辅助编辑 prompt，Markdown 格式更友好

### 3.3 TXT 格式优势

#### 适用场景
- ✅ 纯文本内容，无格式需求
- ✅ 简单 prompt，不需要复杂结构
- ✅ 最小化依赖

#### 当前项目不推荐
- ❌ 现有 prompt 已包含大量 Markdown 语法
- ❌ 使用 TXT 会丢失格式信息
- ❌ 可读性较差

### 3.4 混合方案

可以考虑混合使用：

```
templates/
├── translate.md          # 简单 prompt，但用 .md 保持一致性
├── branch.md             # 复杂 prompt，包含大量 Markdown
├── pr-reword.md          # 中等复杂度
├── file-summary.md       # 中等复杂度
└── pr-summary.md         # 最复杂，包含大量格式
```

**建议**: 统一使用 `.md` 格式，保持一致性。

## 4. 改造方案设计

### 4.1 文件结构

```
internal/llm/prompt/
├── branch.go                    # 简化为函数，从文件读取
├── file.go                      # 简化为函数，从文件读取
├── pr.go                        # 简化为函数，从文件读取
├── translate.go                 # 简化为函数，从文件读取
├── loader.go                    # 统一的文件加载器（新增）
└── templates/
    ├── branch.md
    ├── file-summary.md
    ├── pr-reword.md
    ├── pr-summary.md
    └── translate.md
```

### 4.2 实现方案

#### 方案 A: 统一加载器（推荐）

```go
// loader.go
package prompt

import (
    "embed"
    "fmt"
)

//go:embed templates/*.md
var templatesFS embed.FS

// LoadTemplate 加载模板文件
func LoadTemplate(name string) (string, error) {
    data, err := templatesFS.ReadFile("templates/" + name)
    if err != nil {
        return "", fmt.Errorf("读取模板失败 (%s): %w", name, err)
    }
    return string(data), nil
}

// branch.go
func GenerateBranchSystemPrompt() string {
    prompt, err := LoadTemplate("branch.md")
    if err != nil {
        // 降级到硬编码（可选）
        return defaultBranchPrompt
    }
    return prompt
}
```

#### 方案 B: 直接嵌入（简单场景）

```go
// translate.go
//go:embed templates/translate.md
var translatePrompt string

const TranslateSystemPrompt = translatePrompt
```

### 4.3 动态内容处理

对于需要动态拼接的 prompt（如 `GenerateSummarizePRSystemPrompt`）：

```go
func GenerateSummarizePRSystemPrompt(cfg *config.Manager) string {
    // 从文件加载基础 prompt
    basePrompt, err := LoadTemplate("pr-summary.md")
    if err != nil {
        // 降级处理
        return defaultPrompt
    }

    // 动态内容（JSON 示例等）
    summarizeResponseExample := buildResponseExample()

    // 使用 fmt.Sprintf 拼接
    fullPrompt := fmt.Sprintf(basePrompt, summarizeResponseExample)

    // 应用语言要求
    return llm.GetLanguageRequirement(fullPrompt, cfg)
}
```

## 5. 迁移计划

### 5.1 阶段 1: 准备（已完成 ✅）
- [x] 创建 templates 目录
- [x] 创建示例文件
- [x] 编写使用文档

### 5.2 阶段 2: 简单迁移
- [ ] 迁移 `TranslateSystemPrompt` (最简单)
- [ ] 迁移 `RewordPRSystemPrompt` (中等)
- [ ] 测试验证

### 5.3 阶段 3: 复杂迁移
- [ ] 迁移 `GenerateBranchSystemPrompt`
- [ ] 迁移 `GenerateSummarizeFileChangeSystemPrompt`
- [ ] 迁移 `GenerateSummarizePRSystemPrompt` (最复杂)
- [ ] 完整测试

### 5.4 阶段 4: 优化
- [ ] 统一加载器实现
- [ ] 错误处理优化
- [ ] 性能测试
- [ ] 文档更新

## 6. 推荐方案总结

### 6.1 格式选择: **Markdown (.md)** ✅

**理由**:
1. 现有 prompt 已大量使用 Markdown 语法
2. 更好的可读性和维护性
3. 工具支持更好
4. 未来扩展性更强

### 6.2 实现方案: **统一加载器** ✅

**理由**:
1. 代码复用，减少重复
2. 统一的错误处理
3. 易于扩展和维护
4. 支持动态内容拼接

### 6.3 迁移策略: **渐进式迁移** ✅

**理由**:
1. 降低风险
2. 可以逐步验证
3. 不影响现有功能
4. 可以回滚

## 7. 预期收益

### 7.1 代码质量
- **代码行数**: 减少 ~63% (542行 → ~200行)
- **可维护性**: 显著提升
- **可读性**: 显著提升

### 7.2 开发效率
- **编辑 prompt**: 无需修改 Go 代码
- **版本控制**: Prompt 变更更清晰
- **协作**: 非 Go 开发者也可以参与

### 7.3 运行时
- **二进制大小**: 增加 ~10-20KB (可接受)
- **性能**: 无影响（编译时嵌入）
- **部署**: 更简单（单文件）

## 8. 风险评估

### 8.1 低风险 ✅
- 编译时检查确保文件存在
- 可以保留降级方案
- 可以逐步迁移

### 8.2 需要关注
- 动态内容的处理
- 错误处理策略
- 向后兼容性

## 9. 结论

### ✅ 强烈推荐进行改造

**理由**:
1. ✅ 技术可行性: 100% 可行
2. ✅ 收益明显: 代码质量、可维护性、开发效率显著提升
3. ✅ 风险可控: 可以渐进式迁移，有降级方案
4. ✅ 格式选择: Markdown 更适合当前项目

### 📋 下一步行动

1. **立即开始**: 迁移最简单的 `TranslateSystemPrompt`
2. **验证方案**: 确保加载器方案可行
3. **逐步迁移**: 按照迁移计划逐步完成
4. **文档更新**: 更新相关文档和示例

---

**分析日期**: 2025-01-07
**分析人**: AI Assistant
**状态**: ✅ 建议执行


# Placeholder 和错误提示实现总结

## 概述

已成功实现与 `bubbletea/huh` 一致的 **Placeholder（占位符）** 和 **错误提示** 功能，使用 `lipgloss` 进行样式渲染，不引入 `bubbletea` 依赖。

---

## 1. 实现的功能

### 1.1 Placeholder（占位符）

**功能**：
- ✅ 输入框为空时显示灰色、斜体的占位符文本
- ✅ 用户输入第一个字符时自动清除 placeholder
- ✅ 用户删除所有字符时重新显示 placeholder
- ✅ 使用 `lipgloss` 进行样式渲染（灰色、斜体）

**视觉效果**：
```
请输入用户名: 例如：admin  ← 灰色、斜体的 placeholder
```

### 1.2 错误提示

**功能**：
- ✅ 验证失败时在输入框下方显示红色、加粗的错误提示
- ✅ 用户输入时自动清除错误提示
- ✅ 实时验证（每次输入都进行验证）
- ✅ 使用 `lipgloss` 进行样式渲染（红色、加粗）

**视觉效果**：
```
请输入邮箱: invalid@
请输入有效的邮箱地址  ← 红色、加粗的错误提示
```

### 1.3 状态交互

**与 `bubbletea/huh` 一致的状态交互**：
- ✅ 输入第一个字符 → 清除 placeholder
- ✅ 删除所有字符 → 重新显示 placeholder
- ✅ 输入时 → 清除错误提示
- ✅ 验证失败 → 显示错误提示
- ✅ 验证通过 → 清除错误提示

---

## 2. 引入的依赖

### 2.1 新增依赖

**直接依赖**：
- `github.com/charmbracelet/lipgloss v1.1.0` - 样式库

**自动引入的传递依赖**：
- `github.com/charmbracelet/x/ansi` - ANSI 支持
- `github.com/muesli/termenv` - 终端环境检测
- `github.com/mattn/go-runewidth` - 字符宽度计算（已升级）
- 其他 `lipgloss` 的依赖

### 2.2 依赖对比

**引入前**：
- 直接依赖：7 个
- 传递依赖：~30 个

**引入后**：
- 直接依赖：8 个（+1）
- 传递依赖：~35 个（+5）

**影响**：✅ 影响很小，只增加一个直接依赖

---

## 3. 实现细节

### 3.1 核心函数

**`readInputWithPlaceholder`**：
- 基于 `readPasswordInput` 的实现
- 使用 `golang.org/x/term` 进行字符级输入
- 支持 placeholder 显示和清除
- 支持实时验证和错误提示

**关键特性**：
- 字符级输入（实时响应）
- 状态管理（placeholder、错误状态）
- 终端控制（ANSI 转义码）
- 样式渲染（lipgloss）

### 3.2 样式实现

**Placeholder 样式**：
```go
placeholderStyle := lipgloss.NewStyle().
    Foreground(lipgloss.Color("240")).  // 灰色
    Italic(true)                         // 斜体
```

**错误提示样式**：
```go
errorStyle := lipgloss.NewStyle().
    Foreground(lipgloss.Color("196")).  // 红色
    Bold(true)                           // 加粗
```

---

## 4. 使用方法

### 4.1 链式 API（推荐）

```go
// 带 placeholder 的输入
username, err := prompt.Input().
    Prompt("请输入用户名").
    Placeholder("例如：admin").
    Validate(prompt.ValidateRequired()).
    Run()

// 带 placeholder 和验证的输入
email, err := prompt.Input().
    Prompt("请输入邮箱").
    Placeholder("user@example.com").
    Validate(prompt.ValidateEmail()).
    Run()
```

### 4.2 函数式 API（向后兼容）

```go
// 不带 placeholder（向后兼容）
name, err := prompt.AskInput("请输入您的姓名", "", prompt.ValidateRequired())
```

---

## 5. 代码变更

### 5.1 新增文件

无（基于现有代码扩展）

### 5.2 修改的文件

1. **`internal/lib/prompt/input.go`**：
   - 新增 `readInputWithPlaceholder` 函数
   - 新增 `formatPlaceholder` 函数（使用 lipgloss）
   - 新增 `formatErrorWithLipgloss` 函数
   - 新增 `clearLine`、`clearError`、`showError` 辅助函数
   - 更新 `input` 函数支持 placeholder 参数

2. **`internal/lib/prompt/input_builder.go`**：
   - 新增 `Placeholder` 方法
   - 更新 `InputBuilder` 结构体添加 `placeholder` 字段

3. **`go.mod`**：
   - 新增 `github.com/charmbracelet/lipgloss v1.1.0`

### 5.3 代码量

- 新增代码：约 200 行
- 修改代码：约 50 行
- 总计：约 250 行

---

## 6. 功能对比

### 6.1 与 Bubbletea/Huh 对比

| 特性 | `bubbletea/huh` | 我们的实现 |
|------|-----------------|-----------|
| **Placeholder** | ✅ | ✅ |
| **实时错误提示** | ✅ | ✅ |
| **状态交互** | ✅ | ✅（一致） |
| **样式渲染** | ✅（lipgloss） | ✅（lipgloss） |
| **依赖** | bubbletea + lipgloss | lipgloss（不依赖 bubbletea） |

### 6.2 与当前实现对比

| 特性 | 之前 | 现在 |
|------|------|------|
| **Placeholder** | ❌ | ✅ |
| **实时错误提示** | ✅（回车后） | ✅（实时） |
| **样式** | fatih/color | lipgloss |
| **状态交互** | 简单 | 与 bubbletea/huh 一致 |

---

## 7. 使用示例

### 7.1 基本使用

```go
// 带 placeholder 的输入
result, err := prompt.Input().
    Prompt("请输入您的姓名").
    Placeholder("例如：张三").
    Run()
```

### 7.2 带验证的使用

```go
// 带 placeholder 和验证的输入
email, err := prompt.Input().
    Prompt("请输入邮箱").
    Placeholder("user@example.com").
    Validate(prompt.ValidateEmail()).
    Run()
```

### 7.3 综合使用

```go
// 带 placeholder、默认值和验证的输入
username, err := prompt.Input().
    Prompt("请输入用户名").
    Placeholder("3-20 个字符").
    DefaultValue("admin").
    Validate(prompt.ValidateLength(3, 20)).
    Run()
```

---

## 8. 测试

### 8.1 测试场景

1. **Placeholder 显示**：
   - 输入框为空时显示 placeholder
   - 输入第一个字符时清除 placeholder
   - 删除所有字符时重新显示 placeholder

2. **错误提示**：
   - 验证失败时显示错误
   - 输入时清除错误
   - 验证通过时清除错误

3. **状态交互**：
   - 各种输入场景的状态切换
   - 退格键处理
   - Ctrl+C 处理

### 8.2 测试命令

```bash
# 运行 demo 测试
go run cmd/workflow/main.go demo-prompt
```

---

## 9. 总结

### 9.1 实现成果

1. ✅ **成功引入 `lipgloss`**：只增加一个直接依赖
2. ✅ **实现 Placeholder**：与 `bubbletea/huh` 一致的效果
3. ✅ **实现实时错误提示**：与 `bubbletea/huh` 一致的效果
4. ✅ **状态交互一致**：与 `bubbletea/huh` 完全一致
5. ✅ **不依赖 `bubbletea`**：保持项目轻量级

### 9.2 优势

1. ✅ **功能完整**：支持 placeholder 和实时错误提示
2. ✅ **样式美观**：使用 `lipgloss` 提供更好的视觉效果
3. ✅ **状态一致**：与 `bubbletea/huh` 的状态交互完全一致
4. ✅ **依赖少**：只引入 `lipgloss`，不引入 `bubbletea`
5. ✅ **向后兼容**：不影响现有代码

### 9.3 下一步

- ✅ 功能已实现
- ✅ 编译通过
- ✅ 可以开始测试和使用

---

## 参考资料

- [Lipgloss GitHub](https://github.com/charmbracelet/lipgloss)
- 当前实现：`internal/lib/prompt/input.go`
- Demo：`internal/commands/demo_prompt.go`


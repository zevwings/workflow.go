# Prompt 和 Placeholder 实现分析

## 概述

本文档分析 `charmbracelet/huh` 和 `charmbracelet/bubbletea` 两个库如何实现**提示（Prompt）**和**占位符（Placeholder）**功能，并与当前代码库的实现进行对比。

---

## 1. 概念区分

### 1.1 Prompt（提示）
- **定义**：在输入字段上方或旁边显示的**问题或说明文字**
- **作用**：告诉用户需要输入什么内容
- **示例**：`"请输入您的姓名:"`、`"What's for lunch?"`

### 1.2 Placeholder（占位符）
- **定义**：在输入框**内部**显示的**灰色提示文本**，当用户开始输入时自动消失
- **作用**：提供输入格式示例或提示
- **示例**：`"Enter your name here"`、`"例如：张三"`

### 1.3 Default Value（默认值）
- **定义**：如果用户不输入任何内容，直接回车时使用的**预设值**
- **作用**：提供快速输入选项
- **显示方式**：通常在提示文本旁边显示，如 `"请输入姓名 [Zev]:"`

---

## 2. `charmbracelet/huh` 的实现方式

### 2.1 API 设计

```go
huh.NewInput().
    Title("What's for lunch?").        // Prompt：提示信息
    Placeholder("Enter food name").    // Placeholder：占位符文本
    Value(&lunch)                      // 存储用户输入的值
```

### 2.2 关键特性

1. **分离的 API**：
   - `Title()` - 设置提示信息（Prompt）
   - `Placeholder()` - 设置占位符文本
   - 两者可以独立使用

2. **渲染方式**：
   - 基于 `bubbletea` 构建，使用 TUI（终端用户界面）模式
   - 在输入框**内部**显示 placeholder（灰色文本）
   - 当用户开始输入时，placeholder 自动消失
   - 如果用户清空输入，placeholder 重新显示

3. **视觉表现**：
   ```
   What's for lunch?
   ┌─────────────────────────┐
   │ Enter food name          │  ← Placeholder（灰色，在输入框内）
   └─────────────────────────┘
   ```

### 2.3 实现原理（推测）

基于 `bubbletea` 的 TUI 架构：
- 使用 `lipgloss` 进行样式渲染
- 在 `View()` 方法中判断：
  - 如果 `value == ""` → 显示 placeholder（灰色）
  - 如果 `value != ""` → 显示实际输入值
- 实时更新界面，提供流畅的交互体验

---

## 3. `charmbracelet/bubbletea` + `bubbles/textinput` 的实现方式

### 3.1 API 设计

```go
ti := textinput.New()
ti.Placeholder = "Enter your name..."  // Placeholder 字段
ti.Focus()                              // 获得焦点
ti.Width = 20                           // 设置宽度
```

### 3.2 关键特性

1. **简单的字段设置**：
   - `Placeholder` 是一个字符串字段
   - 直接赋值即可

2. **渲染逻辑**（在 `textinput` 组件的 `View()` 方法中）：
   ```go
   // 伪代码示例
   func (m Model) View() string {
       if m.Value() == "" {
           // 显示 placeholder（灰色样式）
           return lipgloss.NewStyle().
               Foreground(lipgloss.Color("240")).  // 灰色
               Render(m.Placeholder)
       }
       // 显示实际输入值
       return m.Value()
   }
   ```

3. **视觉表现**：
   ```
   Please enter your name:
   ┌──────────────────────┐
   │ Enter your name...    │  ← Placeholder（灰色，光标闪烁）
   └──────────────────────┘
   ```

### 3.3 实现细节

- **状态管理**：使用 `bubbletea` 的 Model 管理输入状态
- **样式渲染**：使用 `lipgloss` 进行颜色和样式控制
- **光标处理**：支持光标闪烁效果（`textinput.Blink`）
- **实时更新**：每次按键都会触发界面更新

---

## 4. 当前代码库的实现方式

### 4.1 当前实现（`internal/lib/prompt/input.go`）

```go
// 当前实现方式
promptText := fmt.Sprintf("%s %s%s%s: ",
    promptMsg,
    GetTheme().InputBracketLeft,  // "["
    defaultValue,                  // 默认值
    GetTheme().InputBracketRight) // "]"
```

### 4.2 关键特性

1. **Prompt 实现**：
   - ✅ 使用 `formatPrompt()` 格式化提示消息
   - ✅ 支持主题颜色（`PromptColor`）
   - ✅ 显示格式：`"提示消息: "`

2. **Default Value 实现**：
   - ✅ 在提示文本中显示：`"提示消息 [默认值]: "`
   - ✅ 用户直接回车时使用默认值
   - ✅ 逻辑处理正确

3. **Placeholder 实现**：
   - ❌ **未实现**：当前代码库**没有**真正的 placeholder 功能
   - 当前只有 default value，显示在提示文本旁边，而不是在输入框内部

### 4.3 当前实现的视觉表现

```
请输入您的姓名 [Zev]: _  ← 提示 + 默认值（在提示文本中）
```

**对比 huh 的实现**：
```
What's for lunch?
┌─────────────────────────┐
│ Enter food name          │  ← Placeholder（在输入框内，灰色）
└─────────────────────────┘
```

---

## 5. 实现方式对比

| 特性 | `huh` | `bubbletea/textinput` | 当前代码库 |
|------|-------|----------------------|-----------|
| **Prompt** | ✅ `Title()` | ✅ 在 `View()` 中自定义 | ✅ `formatPrompt()` |
| **Placeholder** | ✅ `Placeholder()` | ✅ `Placeholder` 字段 | ❌ 未实现 |
| **Default Value** | ✅ `Default()` | ❌ 不支持 | ✅ `defaultValue` 参数 |
| **渲染方式** | TUI（实时更新） | TUI（实时更新） | 行式输出（一次性） |
| **样式支持** | ✅ `lipgloss` | ✅ `lipgloss` | ✅ `fatih/color` |
| **交互体验** | 流畅（TUI） | 流畅（TUI） | 简单（命令行） |

---

## 6. 核心差异分析

### 6.1 架构差异

**huh / bubbletea**：
- 基于 **TUI（Terminal User Interface）** 架构
- 使用 `bubbletea` 的事件循环和状态管理
- 支持**实时界面更新**（每次按键都更新界面）
- 可以显示**光标闪烁**、**占位符高亮**等效果

**当前代码库**：
- 基于 **命令行交互** 架构
- 使用 `bufio.Reader` 读取标准输入
- **一次性输出**（打印提示，等待输入，显示结果）
- 不支持实时界面更新

### 6.2 Placeholder 实现的关键点

**huh / bubbletea 的实现**：
```go
// 伪代码
func renderInput(value, placeholder string) string {
    if value == "" {
        // 显示 placeholder（灰色）
        return grayStyle.Render(placeholder)
    }
    // 显示实际值
    return value
}
```

**关键点**：
1. **条件渲染**：根据输入值是否为空决定显示内容
2. **样式区分**：placeholder 使用灰色，实际输入使用正常颜色
3. **实时更新**：每次按键都重新渲染界面
4. **位置**：placeholder 显示在**输入框内部**，而不是提示文本旁边

### 6.3 为什么当前代码库难以实现 Placeholder？

1. **架构限制**：
   - 当前使用 `bufio.Reader.ReadString('\n')` 读取整行
   - 无法在用户输入过程中实时更新界面
   - 只有在用户按回车后才能获取输入

2. **终端控制**：
   - 要实现 placeholder，需要：
     - 实时读取单个字符（而不是整行）
     - 实时更新终端显示
     - 处理光标位置和退格键
   - 这需要类似 `golang.org/x/term` 的原始模式（类似密码输入的实现）

3. **复杂度**：
   - 当前密码输入已经实现了字符级读取（`readPasswordInput`）
   - 但普通输入使用的是行级读取，需要重构才能支持 placeholder

---

## 7. 实现 Placeholder 的可行方案

### 7.1 方案一：基于现有密码输入逻辑扩展

**思路**：复用 `readPasswordInput` 的字符级读取逻辑，但改为显示实际字符（而不是 `*`）

**优点**：
- 可以复用现有代码
- 不需要引入新的依赖

**缺点**：
- 需要处理更复杂的编辑逻辑（退格、删除、光标移动等）
- 代码复杂度增加

### 7.2 方案二：集成 `bubbletea` 或 `huh`

**思路**：将输入组件迁移到 `bubbletea` 或直接使用 `huh`

**优点**：
- 功能完整，支持 placeholder、光标闪烁等
- 交互体验更好

**缺点**：
- 引入新的依赖
- 需要重构现有代码
- 可能与当前简单的命令行风格不一致

### 7.3 方案三：混合方案（推荐）

**思路**：
- 保持当前简单输入方式（无 placeholder）
- 对于需要 placeholder 的场景，使用 `huh` 或 `bubbletea`
- 提供配置选项，让用户选择交互方式

**优点**：
- 保持向后兼容
- 在需要时可以使用更高级的功能
- 灵活性高

---

## 8. 总结

### 8.1 huh 和 bubbletea 的实现特点

1. **分离的概念**：
   - `Title()` / Prompt：提示信息（在输入框外部）
   - `Placeholder`：占位符文本（在输入框内部，灰色）

2. **TUI 架构**：
   - 基于事件循环和状态管理
   - 支持实时界面更新
   - 提供流畅的交互体验

3. **样式渲染**：
   - 使用 `lipgloss` 进行样式控制
   - 支持颜色、边框、光标等效果

### 8.2 当前代码库的特点

1. **简单实用**：
   - 基于标准库，零外部依赖（除了 `fatih/color`）
   - 适合简单的命令行交互场景

2. **功能完整**：
   - ✅ Prompt（提示信息）
   - ✅ Default Value（默认值）
   - ❌ Placeholder（占位符）

3. **架构限制**：
   - 行级输入，无法实时更新界面
   - 需要重构才能支持 placeholder

### 8.3 建议

如果需要在当前代码库中实现 placeholder：

1. **短期方案**：
   - 保持当前实现，不添加 placeholder
   - 通过改进 default value 的显示方式来提供类似体验

2. **长期方案**：
   - 考虑引入 `bubbletea` 或 `huh` 作为可选依赖
   - 提供两种模式：简单模式（当前实现）和 TUI 模式（huh/bubbletea）
   - 让用户根据需求选择

---

## 参考资料

- [huh GitHub 仓库](https://github.com/charmbracelet/huh)
- [bubbletea GitHub 仓库](https://github.com/charmbracelet/bubbletea)
- [bubbles/textinput 文档](https://github.com/charmbracelet/bubbles/tree/master/textinput)


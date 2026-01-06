# Bubbles/TextInput 独立使用分析

## 概述

分析是否可以单独获取 `bubbles/textinput` 组件并维护到当前项目中，而不引入整个 `bubbletea` 框架。

---

## 1. Bubbles/TextInput 的依赖关系

### 1.1 依赖检查

**`bubbles/textinput` 的依赖**：
- `github.com/charmbracelet/bubbletea` - **必需依赖**
- `github.com/charmbracelet/lipgloss` - 样式渲染
- `github.com/mattn/go-runewidth` - 字符宽度计算
- `github.com/charmbracelet/x/ansi` - ANSI 支持
- `github.com/charmbracelet/x/term` - 终端控制

**关键发现**：
- ❌ `bubbles/textinput` **依赖 `bubbletea`**
- ❌ 不能完全独立使用
- ✅ 但可以提取核心逻辑

---

## 2. 可行性分析

### 2.1 方案一：直接引入 bubbles/textinput（不推荐）

**问题**：
- 会引入 `bubbletea` 依赖
- 需要实现 `bubbletea.Model` 接口
- 需要 `bubbletea` 的事件循环

**结论**：❌ 不可行（会引入 `bubbletea`）

### 2.2 方案二：提取核心逻辑（推荐）

**思路**：
- 参考 `bubbles/textinput` 的实现
- 提取核心逻辑（字符读取、光标管理、placeholder 显示）
- 移除对 `bubbletea` 的依赖
- 适配到当前项目的架构

**优势**：
- ✅ 不引入 `bubbletea` 依赖
- ✅ 可以复用核心逻辑
- ✅ 适配当前项目架构

---

## 3. Bubbles/TextInput 核心功能分析

### 3.1 核心功能

1. **字符级输入**：
   - 使用 `golang.org/x/term` 进行原始模式输入
   - 实时读取单个字符
   - 处理退格、删除、光标移动等

2. **Placeholder 显示**：
   - 条件渲染：`if value == "" { 显示 placeholder }`
   - 样式控制：使用灰色显示

3. **光标管理**：
   - 光标位置跟踪
   - 光标闪烁效果

4. **状态管理**：
   - 输入值状态
   - 错误状态
   - 焦点状态

### 3.2 当前项目已有能力

**已有实现**：
- ✅ `readPasswordInput` - 字符级输入（已实现）
- ✅ `golang.org/x/term` - 终端控制（已引入）
- ✅ 主题系统 - 样式控制（已实现）

**缺失功能**：
- ❌ Placeholder 显示
- ❌ 实时错误提示
- ❌ 光标管理（可选）

---

## 4. 实现方案

### 4.1 基于现有代码扩展

**核心思路**：
- 基于 `readPasswordInput` 的实现
- 添加 placeholder 支持
- 添加实时验证和错误提示
- 保持与 `bubbletea/huh` 一致的状态交互

### 4.2 实现步骤

**第一步：创建 `readInputWithPlaceholder` 函数**

```go
// readInputWithPlaceholder 读取输入（支持 placeholder 和实时错误提示）
func readInputWithPlaceholder(
    promptText string,
    placeholder string,
    validator Validator,
) (string, error) {
    fd := int(os.Stdin.Fd())
    oldState, err := term.MakeRaw(fd)
    if err != nil {
        return readInputFallback(placeholder)
    }
    defer term.Restore(fd, oldState)

    var value []byte
    var buf [1]byte
    var currentError error
    hasPlaceholder := placeholder != ""

    // 初始显示 placeholder
    if hasPlaceholder {
        fmt.Print(formatPlaceholder(placeholder))
    }

    for {
        n, err := os.Stdin.Read(buf[:])
        if err != nil || n == 0 {
            break
        }

        char := buf[0]

        // 处理回车键
        if char == '\r' || char == '\n' {
            // 验证输入
            if validator != nil {
                if err := validator(string(value)); err != nil {
                    currentError = err
                    // 显示错误，继续输入
                    continue
                }
            }
            fmt.Println()
            break
        }

        // 处理退格键
        if char == 127 || char == 8 {
            if len(value) > 0 {
                value = value[:len(value)-1]

                if len(value) == 0 && hasPlaceholder {
                    // 删除到空，重新显示 placeholder
                    clearLine()
                    fmt.Print(formatPlaceholder(placeholder))
                    hasPlaceholder = true
                } else {
                    fmt.Print("\b \b")
                }
            } else if hasPlaceholder {
                // 如果正在显示 placeholder，退格应该清除 placeholder
                clearLine()
                hasPlaceholder = false
            }

            // 清除错误提示
            if currentError != nil {
                clearError()
                currentError = nil
            }
            continue
        }

        // 处理 Ctrl+C
        if char == 3 {
            term.Restore(fd, oldState)
            fmt.Println()
            return "", fmt.Errorf("用户取消输入")
        }

        // 跳过控制字符
        if char < 32 {
            continue
        }

        // 处理普通字符
        if hasPlaceholder {
            // 清除 placeholder
            clearLine()
            hasPlaceholder = false
        }

        // 清除错误提示
        if currentError != nil {
            clearError()
            currentError = nil
        }

        value = append(value, char)
        fmt.Print(string(char))

        // 实时验证（可选）
        if validator != nil {
            if err := validator(string(value)); err != nil {
                currentError = err
                showError(err.Error())
            }
        }
    }

    term.Restore(fd, oldState)
    return strings.TrimSpace(string(value)), nil
}
```

**第二步：实现辅助函数**

```go
// clearLine 清除当前行
func clearLine() {
    fmt.Print("\r")      // 回到行首
    fmt.Print("\033[K")  // 清除到行尾
}

// clearError 清除错误提示（在下一行）
func clearError() {
    fmt.Print("\033[A")  // 上移一行
    fmt.Print("\r")      // 回到行首
    fmt.Print("\033[K")  // 清除到行尾
    fmt.Print("\033[B")  // 下移一行（回到输入行）
}

// showError 显示错误提示（在下一行）
func showError(message string) {
    fmt.Print("\n")  // 换行
    fmt.Print(formatError(message))
    fmt.Print("\033[A")  // 上移一行（回到输入行）
}
```

---

## 5. 与 Bubbles/TextInput 的对比

### 5.1 功能对比

| 特性 | `bubbles/textinput` | 我们的实现 |
|------|---------------------|-----------|
| **字符级输入** | ✅ | ✅ |
| **Placeholder** | ✅ | ✅ |
| **实时错误提示** | ✅ | ✅ |
| **光标管理** | ✅ | ⚠️ 简化版 |
| **状态管理** | ✅ (bubbletea) | ✅ (简单状态) |
| **依赖** | bubbletea + lipgloss | golang.org/x/term |

### 5.2 架构对比

**bubbles/textinput**：
```
bubbletea (事件循环)
  └── textinput.Model (状态管理)
      └── View() (渲染)
```

**我们的实现**：
```
字符级输入循环
  └── 状态变量 (value, error, placeholder)
      └── 实时渲染
```

---

## 6. 推荐方案

### 6.1 方案：自研实现（推荐）

**理由**：
1. ✅ 不引入 `bubbletea` 依赖
2. ✅ 可以复用现有代码（`readPasswordInput`）
3. ✅ 适配当前项目架构
4. ✅ 功能完整（placeholder + 错误提示）

**实现复杂度**：
- 代码量：约 200-300 行
- 难度：中等（基于现有代码扩展）

### 6.2 不推荐：直接引入 bubbles/textinput

**理由**：
1. ❌ 会引入 `bubbletea` 依赖
2. ❌ 需要重构现有架构
3. ❌ 过度设计（不需要完整的 TUI 框架）

---

## 7. 实现建议

### 7.1 分阶段实现

**第一阶段：Placeholder 支持**
- 基于 `readPasswordInput` 实现
- 添加 placeholder 显示和清除逻辑

**第二阶段：实时错误提示**
- 添加实时验证
- 在输入框下方显示错误

**第三阶段：优化和测试**
- 处理边界情况
- 测试各种场景

### 7.2 关键点

1. **状态交互一致**：
   - 输入第一个字符 → 清除 placeholder
   - 删除所有字符 → 重新显示 placeholder
   - 输入时 → 清除错误提示
   - 验证失败 → 显示错误提示

2. **终端控制**：
   - 使用 ANSI 转义码控制光标
   - 使用 `golang.org/x/term` 进行原始模式输入

3. **样式统一**：
   - 使用现有的主题系统
   - placeholder 使用 `HintColor`（灰色）
   - 错误使用 `ErrorColor`（红色）

---

## 8. 总结

### 8.1 结论

**不能直接使用 `bubbles/textinput`**：
- ❌ 它依赖 `bubbletea`
- ❌ 需要完整的 TUI 框架

**可以提取核心逻辑**：
- ✅ 参考其实现思路
- ✅ 基于现有代码扩展
- ✅ 实现类似功能

### 8.2 推荐方案

**自研实现**：
- 基于 `readPasswordInput` 扩展
- 添加 placeholder 和错误提示支持
- 保持与 `bubbletea/huh` 一致的状态交互
- 不引入额外依赖

### 8.3 优势

1. ✅ **零额外依赖**：只使用已有的 `golang.org/x/term`
2. ✅ **功能完整**：支持 placeholder 和错误提示
3. ✅ **状态一致**：与 `bubbletea/huh` 交互一致
4. ✅ **易于维护**：基于现有代码，逻辑清晰

---

## 参考资料

- [Bubbles TextInput GitHub](https://github.com/charmbracelet/bubbles/tree/master/textinput)
- 当前项目：`internal/lib/prompt/input.go`


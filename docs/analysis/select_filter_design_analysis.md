# Select / MultiSelect Filter 功能设计分析

## 概述

分析如何在 `select` 和 `multiselect` 组件中添加 **Filter（过滤）** 功能，参考 `charmbracelet/bubbletea` 和 `charmbracelet/huh` 的设计思路，结合当前基于 ANSI 转义码的实现方式。

---

## 1. 需求分析

### 1.1 功能需求

**核心功能**：
- 用户可以通过输入文本过滤选项列表
- 实时过滤：输入时立即更新显示的选项
- 保持原有的导航和选择功能
- 支持中文和英文过滤

**交互需求**：
- 输入 `/` 或直接开始输入进入过滤模式
- 显示过滤输入框（在选项列表上方或下方）
- 显示过滤后的选项数量
- 支持清除过滤（Backspace 删除所有字符或 Ctrl+U）

### 1.2 使用场景

1. **选项较多时**：当选项列表超过 20 项时，快速定位目标选项
2. **模糊搜索**：用户只记得选项的部分内容
3. **多语言支持**：支持中文、英文混合搜索

---

## 2. Bubbletea / Huh 的设计参考

### 2.1 Bubbles List 组件的 Filter 实现

**核心组件**：`github.com/charmbracelet/bubbles/list`

**实现方式**：

```go
type Model struct {
    Items       []Item        // 所有选项
    FilterValue string        // 过滤文本
    // ... 其他字段
}

// 过滤逻辑
func (m Model) FilteredItems() []Item {
    if m.FilterValue == "" {
        return m.Items
    }

    var filtered []Item
    for _, item := range m.Items {
        if strings.Contains(
            strings.ToLower(item.FilterValue()),
            strings.ToLower(m.FilterValue),
        ) {
            filtered = append(filtered, item)
        }
    }
    return filtered
}

// 渲染逻辑
func (m Model) View() string {
    var s strings.Builder

    // 显示过滤输入框
    if m.FilteringEnabled {
        s.WriteString("Filter: ")
        s.WriteString(m.FilterValue)
        s.WriteString("\n")
    }

    // 显示过滤后的选项
    filtered := m.FilteredItems()
    for i, item := range filtered {
        if i == m.Cursor {
            s.WriteString("> " + item.Title())
        } else {
            s.WriteString("  " + item.Title())
        }
        s.WriteString("\n")
    }

    return s.String()
}
```

**关键特性**：
1. **状态管理**：使用 `FilterValue` 存储过滤文本
2. **实时过滤**：每次输入都重新计算过滤结果
3. **大小写不敏感**：统一转换为小写进行比较
4. **模糊匹配**：使用 `strings.Contains` 进行子串匹配

### 2.2 Huh Select 的 Filter 实现

**核心组件**：`github.com/charmbracelet/huh`

**API 设计**：

```go
select := huh.NewSelect[string]().
    Title("Choose a fruit").
    Options(
        huh.NewOption("Apple", "apple"),
        huh.NewOption("Banana", "banana"),
        // ...
    ).
    Filterable(true).  // 启用过滤
    Value(&fruit)
```

**内部实现**：

```go
type Select struct {
    options      []Option
    filterValue  string
    filtering    bool
    // ... 其他字段
}

func (s *Select) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.KeyMsg:
        // 进入过滤模式
        if msg.Type == tea.KeyRunes && s.filtering {
            s.filterValue += string(msg.Runes)
            return s, nil
        }

        // 开始过滤（输入 '/' 或直接输入）
        if msg.Type == tea.KeyRunes && !s.filtering {
            s.filtering = true
            s.filterValue = string(msg.Runes)
            return s, nil
        }

        // 退出过滤模式
        if msg.Type == tea.KeyEsc && s.filtering {
            s.filtering = false
            s.filterValue = ""
            return s, nil
        }
    }
    // ... 其他处理
}

func (s *Select) filteredOptions() []Option {
    if s.filterValue == "" {
        return s.options
    }

    var filtered []Option
    filterLower := strings.ToLower(s.filterValue)
    for _, opt := range s.options {
        if strings.Contains(
            strings.ToLower(opt.Title),
            filterLower,
        ) {
            filtered = append(filtered, opt)
        }
    }
    return filtered
}
```

**关键特性**：
1. **模式切换**：通过 `filtering` 标志控制是否处于过滤模式
2. **触发方式**：输入 `/` 或直接输入字符进入过滤模式
3. **退出方式**：按 `Esc` 退出过滤模式
4. **链式 API**：提供 `Filterable(true)` 方法启用过滤

---

## 3. 当前实现分析

### 3.1 Select 组件当前实现

**核心特点**：
- 使用 `golang.org/x/term` 进行原始输入处理
- 使用 ANSI 转义码进行界面渲染
- 支持上下箭头导航
- 支持回车确认

**代码结构**：

```go
func selectOption(message string, options []string, defaultIndex int) (int, error) {
    // 1. 设置原始模式
    oldState, err := term.MakeRaw(fd)

    // 2. 渲染函数
    renderSelect := func(isFirst bool) {
        // 清除屏幕
        // 渲染选项列表
        // 显示提示信息
    }

    // 3. 输入循环
    for {
        // 读取输入
        // 处理箭头键
        // 处理回车键
        // 处理 Ctrl+C
    }
}
```

### 3.2 MultiSelect 组件当前实现

**核心特点**：
- 与 Select 类似的结构
- 支持空格键切换选择状态
- 使用 `map[int]bool` 存储选中状态

---

## 4. Filter 功能设计方案

### 4.1 设计目标

1. **最小侵入**：尽量不改动现有代码结构
2. **向后兼容**：默认不启用过滤，需要时通过参数启用
3. **用户体验**：流畅的交互，清晰的视觉反馈
4. **性能考虑**：过滤算法高效，支持大量选项

### 4.2 状态管理设计

**新增状态字段**：

```go
type selectState struct {
    // 原有字段
    currentIndex int
    options      []string

    // 新增字段
    filterValue    string    // 过滤文本
    isFiltering    bool      // 是否处于过滤模式
    filteredOptions []string  // 过滤后的选项
    originalIndices []int     // 过滤后选项对应的原始索引
}
```

**状态转换**：

```
正常模式 → 输入 '/' 或字符 → 过滤模式
过滤模式 → 输入字符 → 更新过滤文本
过滤模式 → 按 Esc → 退出过滤模式（清除过滤）
过滤模式 → Backspace 删除所有 → 退出过滤模式
```

### 4.3 过滤算法设计

**基础过滤**（字符串包含匹配）：

```go
func filterOptions(options []string, filterText string) ([]string, []int) {
    if filterText == "" {
        // 返回所有选项和原始索引
        indices := make([]int, len(options))
        for i := range options {
            indices[i] = i
        }
        return options, indices
    }

    filterLower := strings.ToLower(filterText)
    var filtered []string
    var indices []int

    for i, opt := range options {
        if strings.Contains(strings.ToLower(opt), filterLower) {
            filtered = append(filtered, opt)
            indices = append(indices, i)
        }
    }

    return filtered, indices
}
```

**高级过滤**（可选功能）：
- 支持正则表达式
- 支持多关键词搜索（空格分隔）
- 支持高亮匹配部分

### 4.4 界面渲染设计

**布局结构**：

```
[提示消息]
[过滤输入框] (如果启用过滤)
  > 选项1
    选项2
    选项3
[提示信息：使用 ↑/↓ 导航，回车确认，/ 过滤]
```

**过滤输入框样式**：

```go
func renderFilterInput(filterValue string, isActive bool) string {
    prefix := "Filter: "
    if isActive {
        // 高亮显示
        return formatAnswer(prefix + filterValue + "_")  // 下划线表示光标
    }
    return formatHint(prefix + filterValue)
}
```

**选项渲染**：

```go
func renderOptions(options []string, currentIndex int, selected map[int]bool) {
    for i, option := range options {
        prefix := "  "
        if i == currentIndex {
            prefix = "> "
        }

        marker := ""
        if selected != nil && selected[i] {
            marker = "[x] "
        }

        line := prefix + marker + option
        if i == currentIndex {
            fmt.Print(formatAnswer(line))
        } else {
            fmt.Print(line)
        }
        fmt.Print("\033[K\n")
    }
}
```

### 4.5 输入处理设计

**输入处理流程**：

```go
for {
    char := readChar()

    // 1. 处理过滤模式
    if isFiltering {
        // 处理字符输入
        if isPrintable(char) {
            filterValue += string(char)
            updateFilteredOptions()
            render()
            continue
        }

        // 处理 Backspace
        if char == '\b' || char == 127 {
            if len(filterValue) > 0 {
                filterValue = filterValue[:len(filterValue)-1]
                updateFilteredOptions()
                render()
            } else {
                // 删除完所有字符，退出过滤模式
                isFiltering = false
                filterValue = ""
                updateFilteredOptions()
                render()
            }
            continue
        }

        // 处理 Esc（退出过滤模式）
        if char == '\x1b' {
            isFiltering = false
            filterValue = ""
            updateFilteredOptions()
            render()
            continue
        }

        // 处理 Ctrl+U（清除过滤）
        if char == 21 { // Ctrl+U
            filterValue = ""
            updateFilteredOptions()
            render()
            continue
        }
    }

    // 2. 处理正常模式
    // 进入过滤模式：输入 '/' 或直接输入字符
    if !isFiltering && (char == '/' || isPrintable(char)) {
        isFiltering = true
        if char != '/' {
            filterValue = string(char)
        } else {
            filterValue = ""
        }
        updateFilteredOptions()
        render()
        continue
    }

    // 3. 处理导航（在过滤后的选项列表中）
    if isArrowKey(char) {
        handleNavigation()
        render()
        continue
    }

    // 4. 处理确认和取消
    // ...
}
```

**关键点**：
1. **模式识别**：区分过滤模式和正常模式
2. **字符处理**：正确处理可打印字符和特殊键
3. **状态同步**：过滤后更新 `currentIndex` 和选项列表
4. **索引映射**：维护过滤后选项到原始选项的索引映射

---

## 5. 实现方案

### 5.1 Select 组件 Filter 实现

**函数签名**：

```go
func selectOptionWithFilter(
    message string,
    options []string,
    defaultIndex int,
    enableFilter bool,  // 是否启用过滤
) (int, error)
```

**核心实现**：

```go
func selectOptionWithFilter(
    message string,
    options []string,
    defaultIndex int,
    enableFilter bool,
) (int, error) {
    // 状态初始化
    state := &selectState{
        options:        options,
        currentIndex:    defaultIndex,
        filterValue:    "",
        isFiltering:    false,
        filteredOptions: options,
        originalIndices: make([]int, len(options)),
    }
    for i := range options {
        state.originalIndices[i] = i
    }

    // 过滤函数
    updateFilter := func() {
        if !enableFilter || state.filterValue == "" {
            state.filteredOptions = state.options
            state.originalIndices = make([]int, len(state.options))
            for i := range state.options {
                state.originalIndices[i] = i
            }
            return
        }

        filterLower := strings.ToLower(state.filterValue)
        state.filteredOptions = nil
        state.originalIndices = nil

        for i, opt := range state.options {
            if strings.Contains(strings.ToLower(opt), filterLower) {
                state.filteredOptions = append(state.filteredOptions, opt)
                state.originalIndices = append(state.originalIndices, i)
            }
        }

        // 确保 currentIndex 在有效范围内
        if state.currentIndex >= len(state.filteredOptions) {
            state.currentIndex = len(state.filteredOptions) - 1
        }
        if state.currentIndex < 0 {
            state.currentIndex = 0
        }
    }

    // 渲染函数
    renderSelect := func(isFirst bool) {
        if !isFirst {
            fmt.Print("\033[u")
            fmt.Print("\r")
            fmt.Print("\033[J")
        }

        // 显示过滤输入框
        if enableFilter {
            filterLine := renderFilterInput(state.filterValue, state.isFiltering)
            fmt.Print("\r")
            fmt.Print(filterLine)
            fmt.Print("\033[K\n")
        }

        // 显示选项列表
        for i, option := range state.filteredOptions {
            fmt.Print("\r")
            if i == state.currentIndex {
                selectedText := formatAnswer("> " + option)
                fmt.Print(selectedText)
            } else {
                fmt.Print("  " + option)
            }
            fmt.Print("\033[K\n")
        }

        // 显示提示信息
        hintMsg := formatHint("使用 ↑/↓ 导航，回车确认")
        if enableFilter {
            hintMsg += "，/ 或输入字符过滤"
        }
        fmt.Print("\r")
        fmt.Print(hintMsg)
        fmt.Print("\r\n")

        fmt.Print("\033[?25l")
    }

    // 输入处理循环
    // ... (详细实现见下文)
}
```

### 5.2 MultiSelect 组件 Filter 实现

**核心差异**：
- 需要维护选中状态（基于原始索引）
- 渲染时需要显示选中标记
- 确认时返回原始索引列表

**选中状态管理**：

```go
// 使用原始索引存储选中状态
selected := make(map[int]bool)  // key 是原始索引

// 切换选择时
originalIdx := state.originalIndices[state.currentIndex]
if selected[originalIdx] {
    delete(selected, originalIdx)
} else {
    selected[originalIdx] = true
}

// 渲染时检查选中状态
for i, option := range state.filteredOptions {
    originalIdx := state.originalIndices[i]
    marker := "[ ]"
    if selected[originalIdx] {
        marker = "[x]"
    }
    // ...
}
```

### 5.3 输入处理详细实现

**字符识别**：

```go
func isPrintable(char byte) bool {
    return char >= 32 && char < 127 && char != '\r' && char != '\n'
}

func isArrowKey(char byte) bool {
    return char == '\x1b'  // 需要进一步读取转义序列
}
```

**转义序列处理**：

```go
// 处理转义序列（箭头键）
if char == '\x1b' {
    // 如果在过滤模式，Esc 退出过滤
    if state.isFiltering {
        state.isFiltering = false
        state.filterValue = ""
        updateFilter()
        renderSelect(false)
        continue
    }

    // 否则处理箭头键
    seq := readEscapeSequence()
    if seq == "[A" || seq == "OA" {
        // 上箭头
        if state.currentIndex > 0 {
            state.currentIndex--
            renderSelect(false)
        }
        continue
    }
    if seq == "[B" || seq == "OB" {
        // 下箭头
        if state.currentIndex < len(state.filteredOptions)-1 {
            state.currentIndex++
            renderSelect(false)
        }
        continue
    }
}
```

**过滤模式处理**：

```go
// 进入过滤模式
if enableFilter && !state.isFiltering {
    if char == '/' {
        state.isFiltering = true
        state.filterValue = ""
        renderSelect(false)
        continue
    }
    if isPrintable(char) {
        state.isFiltering = true
        state.filterValue = string(char)
        updateFilter()
        renderSelect(false)
        continue
    }
}

// 过滤模式下的输入
if enableFilter && state.isFiltering {
    if isPrintable(char) {
        state.filterValue += string(char)
        updateFilter()
        renderSelect(false)
        continue
    }

    // Backspace
    if char == '\b' || char == 127 {
        if len(state.filterValue) > 0 {
            state.filterValue = state.filterValue[:len(state.filterValue)-1]
            updateFilter()
            renderSelect(false)
        } else {
            // 删除完所有字符，退出过滤模式
            state.isFiltering = false
            state.filterValue = ""
            updateFilter()
            renderSelect(false)
        }
        continue
    }

    // Ctrl+U 清除过滤
    if char == 21 {
        state.filterValue = ""
        updateFilter()
        renderSelect(false)
        continue
    }
}
```

---

## 6. API 设计

### 6.1 函数式 API（向后兼容）

```go
// 原有函数保持不变
func AskSelect(message string, options []string, defaultIndex int) (int, error)

// 新增带过滤的函数
func AskSelectWithFilter(
    message string,
    options []string,
    defaultIndex int,
    enableFilter bool,
) (int, error)
```

### 6.2 Builder API（推荐）

```go
// Select Builder
type SelectBuilder struct {
    message      string
    options      []string
    defaultIndex int
    enableFilter bool
}

func NewSelect() *SelectBuilder {
    return &SelectBuilder{
        enableFilter: false,
    }
}

func (b *SelectBuilder) Message(msg string) *SelectBuilder {
    b.message = msg
    return b
}

func (b *SelectBuilder) Options(opts []string) *SelectBuilder {
    b.options = opts
    return b
}

func (b *SelectBuilder) DefaultIndex(idx int) *SelectBuilder {
    b.defaultIndex = idx
    return b
}

func (b *SelectBuilder) Filterable(enable bool) *SelectBuilder {
    b.enableFilter = enable
    return b
}

func (b *SelectBuilder) Run() (int, error) {
    return selectOptionWithFilter(
        b.message,
        b.options,
        b.defaultIndex,
        b.enableFilter,
    )
}

// 使用示例
idx, err := NewSelect().
    Message("选择水果").
    Options([]string{"苹果", "香蕉", "橙子"}).
    DefaultIndex(0).
    Filterable(true).
    Run()
```

---

## 7. 实现细节

### 7.1 光标位置管理

**问题**：过滤输入框和选项列表需要正确管理光标位置

**解决方案**：

```go
// 保存光标位置（在过滤输入框之前）
fmt.Print("\033[s")

// 渲染时
// 1. 恢复光标位置
fmt.Print("\033[u")
fmt.Print("\r")
// 2. 清除到屏幕底部
fmt.Print("\033[J")
// 3. 渲染内容
// 4. 隐藏光标
fmt.Print("\033[?25l")
```

### 7.2 索引映射维护

**问题**：过滤后需要正确映射到原始选项

**解决方案**：

```go
// 维护两个数组
filteredOptions []string  // 过滤后的选项
originalIndices []int     // 对应的原始索引

// 返回结果时
originalIndex := originalIndices[currentIndex]
return originalIndex, nil
```

### 7.3 性能优化

**问题**：大量选项时过滤可能较慢

**优化方案**：

1. **延迟过滤**：只在用户停止输入一段时间后过滤（需要实现输入缓冲）
2. **索引缓存**：缓存过滤结果
3. **并行过滤**：使用 goroutine 进行过滤（对于超大量选项）

**简单优化**：

```go
// 使用 strings.ToLower 的缓存
var filterCache = make(map[string]string)

func toLowerCached(s string) string {
    if cached, ok := filterCache[s]; ok {
        return cached
    }
    lower := strings.ToLower(s)
    filterCache[s] = lower
    return lower
}
```

### 7.4 中文支持

**问题**：中文字符的输入和处理

**解决方案**：

```go
// 使用 UTF-8 编码处理
import "unicode/utf8"

func isPrintableUTF8(char byte) bool {
    // 检查是否是 UTF-8 字符的开始字节
    return char >= 32 && char < 127 ||
           (char & 0x80) != 0  // UTF-8 多字节字符
}

// 读取完整的 UTF-8 字符
func readUTF8Char(buf []byte, pos int) (rune, int) {
    return utf8.DecodeRune(buf[pos:])
}
```

---

## 8. 测试方案

### 8.1 单元测试

```go
func TestFilterOptions(t *testing.T) {
    options := []string{"Apple", "Banana", "Orange", "Grape"}

    tests := []struct {
        filter   string
        expected []string
    }{
        {"", []string{"Apple", "Banana", "Orange", "Grape"}},
        {"a", []string{"Apple", "Banana", "Orange", "Grape"}},
        {"ap", []string{"Apple", "Grape"}},
        {"z", []string{}},
    }

    for _, tt := range tests {
        filtered, _ := filterOptions(options, tt.filter)
        if !reflect.DeepEqual(filtered, tt.expected) {
            t.Errorf("filterOptions(%v, %q) = %v, want %v",
                options, tt.filter, filtered, tt.expected)
        }
    }
}
```

### 8.2 集成测试

```go
func TestSelectWithFilter(t *testing.T) {
    // 模拟输入
    // 测试过滤功能
    // 测试导航功能
    // 测试确认功能
}
```

---

## 9. 迁移计划

### 9.1 阶段一：基础实现

1. 实现 `selectOptionWithFilter` 函数
2. 实现 `multiselectOptionsWithFilter` 函数
3. 添加基础测试

### 9.2 阶段二：API 完善

1. 实现 Builder API
2. 更新文档
3. 添加示例代码

### 9.3 阶段三：优化和增强

1. 性能优化
2. 中文支持优化
3. 高级过滤功能（可选）

---

## 10. 参考资源

### 10.1 相关库

- **bubbles/list**: https://github.com/charmbracelet/bubbles/tree/master/list
- **huh**: https://github.com/charmbracelet/huh
- **bubbletea**: https://github.com/charmbracelet/bubbletea

### 10.2 实现参考

- Bubbles List Filter 实现
- Huh Select Filter 实现
- ANSI 转义码文档

---

## 11. 总结

### 11.1 核心设计要点

1. **状态管理**：清晰的状态转换和索引映射
2. **输入处理**：区分过滤模式和正常模式
3. **界面渲染**：清晰的视觉反馈和布局
4. **向后兼容**：保持原有 API 不变

### 11.2 实现难点

1. **输入模式切换**：正确处理过滤模式和正常模式的切换
2. **索引映射**：维护过滤后选项到原始选项的正确映射
3. **光标管理**：正确管理过滤输入框和选项列表的光标位置
4. **中文支持**：正确处理 UTF-8 编码的中文字符

### 11.3 后续优化方向

1. **高级过滤**：正则表达式、多关键词搜索
2. **性能优化**：大量选项时的过滤性能
3. **用户体验**：高亮匹配部分、过滤统计信息
4. **可配置性**：过滤触发键、过滤算法可配置

---

**文档创建日期**：2025-01-XX
**最后更新日期**：2025-01-XX


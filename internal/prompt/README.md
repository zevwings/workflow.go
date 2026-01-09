# prompt 模块

本模块提供了统一的用户交互功能，包括输入提示、选择提示、表单提示、消息输出、加载指示器、表格显示等。

## 文件说明

```
internal/prompt/
├── builder.go              # Builder 模式基础结构（29行）
├── input.go                 # 输入提示（Input/Password）（329行）
├── confirm.go               # 确认提示（49行）
├── select.go                # 单选提示（56行）
├── multiselect.go           # 多选提示（56行）
├── form.go                  # 表单提示（57行）
├── message.go               # 消息输出工具（98行）
├── spinner.go              # 加载指示器（292行）
├── table.go                 # 表格显示工具（307行）
├── theme.go                 # 主题配置（157行）
│
├── common/                   # 通用功能模块
│   ├── config.go            # 提示功能通用配置（PromptConfig、BasePromptConfig）
│   ├── config_manager.go    # 配置管理器（ConfigManager，支持默认/全局/局部配置）
│   ├── format.go            # 格式化函数（FormatResult、FormatResultWithTitle 等）
│   ├── render.go            # 渲染功能（RenderOptions 等）
│   ├── navigation.go        # 导航功能（NavigationHandler，键盘方向键处理）
│   ├── input_handler.go     # 输入处理（HandleInteractiveInput，键盘事件处理）
│   ├── fallback.go          # Fallback 机制（TypedFallbackHandler、ExecuteFallbackTyped）
│   ├── select_helpers.go    # 选择辅助函数（ExecuteSelectFallback、ExecuteMultiSelectFallback）
│   └── cancel.go            # 取消功能（Ctrl+C 处理）
│
├── input/                    # 输入子模块
│   ├── editor.go            # 输入编辑器（字符级输入、光标移动）
│   ├── handler.go           # 输入处理器（键盘事件处理）
│   ├── format.go            # 格式化函数（占位符、错误提示）
│   └── validator.go         # 验证器（邮箱、URL、长度等）
│
├── confirm/                  # 确认子模块
│   ├── core.go              # 确认核心逻辑
│   ├── handler.go          # 确认处理器（键盘事件处理）
│   └── adapter.go          # Fallback 适配器（confirmFallbackAdapter）
│
├── select/                   # 单选子模块
│   ├── core.go              # 选择核心逻辑
│   └── handler.go          # 选择处理器（键盘事件处理）
│
├── multiselect/              # 多选子模块
│   ├── core.go              # 多选核心逻辑
│   └── handler.go          # 多选处理器（键盘事件处理）
│
├── form/                     # 表单子模块
│   ├── builder.go           # 表单构建器（链式 API）
│   ├── executor.go          # 表单执行器（执行表单流程）
│   ├── field.go             # 表单字段定义
│   ├── result.go            # 表单结果定义
│   ├── validator.go         # 表单验证器
│   └── config.go           # 表单配置（格式化函数注入）
│
└── io/                       # I/O 抽象模块
    ├── terminal.go          # 终端 I/O 接口定义
    ├── stdterminal.go       # 标准终端实现
    ├── mockterminal.go      # Mock 终端实现（用于测试）
    ├── rawmode.go           # 原始模式控制
    ├── renderer.go          # 渲染器（ANSI 转义序列）
    └── escape.go            # ANSI 转义序列工具
```

### 核心文件

- **`input.go`**：输入提示功能，提供 `Input()` 和 `Password()` 构建器，支持默认值、占位符、验证等
- **`confirm.go`**：确认提示功能，提供 `Confirm()` 构建器，支持 Yes/No 选择
- **`select.go`**：单选功能，提供 `Select()` 构建器，支持从多个选项中选择一个
- **`multiselect.go`**：多选功能，提供 `MultiSelect()` 构建器，支持从多个选项中选择多个
- **`form.go`**：表单功能，提供 `Form()` 构建器，支持组合多个字段进行输入
- **`message.go`**：消息输出工具，提供 `Message` 结构体，支持不同级别的消息输出
- **`spinner.go`**：加载指示器，提供 `Spinner` 结构体，支持加载动画显示
- **`table.go`**：表格显示工具，提供 `Table` 结构体，支持表格渲染
- **`theme.go`**：主题配置，提供 `Theme` 结构体和全局主题管理

## 快速开始

### 输入提示

```go
import "github.com/zevwings/workflow/internal/prompt"

// 基础输入 - 使用配置结构体
value, err := prompt.AskInput(prompt.InputField{
    Message:      "请输入您的姓名",
    DefaultValue: "",
    Validator:    nil,
    ResultTitle:  "姓名",  // 可选
})

// 带默认值和验证的输入 - 使用配置结构体
email, err := prompt.AskInput(prompt.InputField{
    Message:      "请输入邮箱",
    DefaultValue: "user@example.com",
    Validator:    prompt.ValidateEmail(),
    ResultTitle:  "邮箱",  // 可选
})

// 密码输入 - 使用配置结构体
password, err := prompt.AskPassword(prompt.PasswordField{
    Message:      "请输入密码",
    DefaultValue: "",  // 空字符串表示无默认值
    Validator:    prompt.ValidateMinLength(8),
    ResultTitle:  "密码",  // 可选
})

// Builder 模式调用
value, err := prompt.Input().
    Prompt("请输入您的姓名").
    Run()

email, err := prompt.Input().
    Prompt("请输入邮箱").
    DefaultValue("user@example.com").
    Validate(prompt.ValidateEmail()).
    Run()

password, err := prompt.Password().
    Prompt("请输入密码").
    Validate(prompt.ValidateMinLength(8)).
    Run()
```

### 确认提示

```go
import "github.com/zevwings/workflow/internal/prompt"

// 使用配置结构体
confirmed, err := prompt.AskConfirm(prompt.ConfirmField{
    Message:     "是否继续？",
    DefaultYes:  true,
    ResultTitle: "继续",  // 可选
})

// Builder 模式调用
confirmed, err := prompt.Confirm().
    Prompt("是否继续？").
    Default(true).
    Run()
```

### 选择提示

```go
import "github.com/zevwings/workflow/internal/prompt"

options := []string{"选项1", "选项2", "选项3"}

// 单选 - 使用配置结构体
index, err := prompt.AskSelect(prompt.SelectField{
    Message:      "请选择一个选项",
    Options:      options,
    DefaultIndex: 0,
    ResultTitle:  "选择的选项",  // 可选
})

// 单选 - Builder 模式调用
index, err := prompt.Select().
    Prompt("请选择一个选项").
    Options(options).
    Default(0).
    Run()

// 多选 - 使用配置结构体
selected, err := prompt.AskMultiSelect(prompt.MultiSelectField{
    Message:         "请选择多个选项",
    Options:         options,
    DefaultSelected: []int{0, 2},
    ResultTitle:     "选择的选项",  // 可选
})

// 多选 - Builder 模式调用
selected, err := prompt.MultiSelect().
    Prompt("请选择多个选项").
    Options(options).
    Default([]int{0, 2}).
    Run()
```

### 表单提示

```go
import (
    "github.com/zevwings/workflow/internal/prompt"
    "github.com/zevwings/workflow/internal/prompt/form"
)

result, err := prompt.Form().
    AddInput(form.InputFormField{
        Key:          "name",
        Prompt:       "姓名",
        DefaultValue: "",
        Validator:    prompt.ValidateRequired(),
    }).
    AddInput(form.InputFormField{
        Key:          "email",
        Prompt:       "邮箱",
        DefaultValue: "",
        Validator:    prompt.ValidateEmail(),
    }).
    AddPassword(form.PasswordFormField{
        Key:          "password",
        Prompt:       "密码",
        DefaultValue: "",
        Validator:    prompt.ValidateMinLength(8),
    }).
    AddSelect(form.SelectFormField{
        Key:          "role",
        Prompt:       "角色",
        Options:      []string{"开发者", "测试"},
        DefaultIndex: 0,
    }).
    AddConfirm(form.ConfirmFormField{
        Key:          "agree",
        Prompt:       "同意协议",
        DefaultValue: false,
    }).
    Run()

if err != nil {
    return err
}

name := result.GetString("name")
email := result.GetString("email")
roleIndex := result.GetInt("role")
agree := result.GetBool("agree")
```

### 表单结果标题格式化

可以通过 `SetFormFormatResultTitle` 自定义表单字段完成后显示的 title 格式：

```go
// 使用内置的格式化函数
prompt.SetFormFormatResultTitle(prompt.FormatResultTitleForForm)

// 或自定义格式化函数
prompt.SetFormFormatResultTitle(func(originalMessage string, resultValue string) string {
    // 自定义逻辑
    return "自定义标题: " + resultValue
})
```

`FormatResultTitleForForm` 是一个辅助函数，可以将 "Please enter your X" 转换为 "Your X"。

### 消息输出

```go
import "github.com/zevwings/workflow/internal/prompt"

msg := prompt.NewMessage(true) // verbose 模式

msg.Info("这是一条信息")
msg.Success("操作成功")
msg.Warning("这是一条警告")
msg.Error("这是一条错误")
msg.Debug("这是调试信息") // 仅在 verbose 模式下显示
```

### 加载指示器

```go
import "github.com/zevwings/workflow/internal/prompt"
import "time"

// 基础使用
spinner := prompt.NewSpinner("正在处理...")
spinner.Start()
defer spinner.Stop()

// 执行操作
time.Sleep(2 * time.Second)

// 停止并显示成功消息
spinner.WithSuccess("处理完成")

// 使用 Do 方法
spinner := prompt.NewSpinner("正在处理...")
err := spinner.Do(func() error {
    // 执行操作
    return nil
})

// 使用选项自定义 Spinner
spinner := prompt.NewSpinner(
    "正在处理...",
    prompt.WithSpinner([]string{"⣾", "⣽", "⣻", "⢿", "⡿", "⣟", "⣯", "⣷"}),
    prompt.WithInterval(50 * time.Millisecond),
    prompt.WithSpinnerStyle(lipgloss.NewStyle().Foreground(lipgloss.Color("blue"))),
    prompt.WithMessageStyle(lipgloss.NewStyle().Foreground(lipgloss.Color("gray"))),
)
```

**Spinner 选项**：
- `WithSpinner(frames []string)` - 设置自定义 spinner 字符序列
- `WithInterval(d time.Duration)` - 设置更新间隔
- `WithStyle(style lipgloss.Style)` - 设置自定义样式（同时应用于 spinner 和文本）
- `WithSpinnerStyle(style lipgloss.Style)` - 设置 spinner 字符的样式
- `WithMessageStyle(style lipgloss.Style)` - 设置消息文本的样式
- `WithWriter(w io.Writer)` - 设置输出流

### 表格显示

```go
import "github.com/zevwings/workflow/internal/prompt"

table := prompt.NewTable([]string{"姓名", "年龄", "邮箱"})
table.AddRow([]string{"张三", "25", "zhangsan@example.com"})
table.AddRow([]string{"李四", "30", "lisi@example.com"})
table.SetBorder(true).
    SetRowLine(true).
    SetAlignment(prompt.ALIGN_LEFT).
    Render()
```

## 主要接口

### InputBuilder

- `Prompt(message string)` - 设置提示消息
- `DefaultValue(value string)` - 设置默认值
- `Placeholder(text string)` - 设置占位符文本
- `Validate(validator Validator)` - 设置验证器
- `Run()` - 执行输入并返回结果

### ConfirmBuilder

- `Prompt(message string)` - 设置提示消息
- `Default(defaultYes bool)` - 设置默认值
- `Run()` - 执行确认并返回结果

### SelectBuilder

- `Prompt(message string)` - 设置提示消息
- `Options(options []string)` - 设置选项列表
- `Default(index int)` - 设置默认选中的索引
- `Run()` - 执行选择并返回结果（返回选中索引）

### MultiSelectBuilder

- `Prompt(message string)` - 设置提示消息
- `Options(options []string)` - 设置选项列表
- `Default(indices []int)` - 设置默认选中的索引列表
- `Run()` - 执行多选并返回结果（返回选中索引列表）

### FormBuilder

- `AddInput(config InputFormField)` - 添加输入字段
- `AddPassword(config PasswordFormField)` - 添加密码字段
- `AddConfirm(config ConfirmFormField)` - 添加确认字段
- `AddSelect(config SelectFormField)` - 添加单选字段
- `AddMultiSelect(config MultiSelectFormField)` - 添加多选字段
- `AddForm(config NestedFormField)` - 添加嵌套表单字段
- `Condition(condition Condition)` - 为最后一个字段设置条件函数（已废弃，建议在配置中设置）
- `Validate(validator FormValidator)` - 设置表单级验证器

### 配置结构体

所有模块都支持使用配置结构体，提供更清晰的参数传递方式：

#### ConfirmField - 确认字段配置

```go
type ConfirmField struct {
    Message     string  // 提示消息
    DefaultYes  bool    // 默认值（true 表示默认 Yes）
    ResultTitle string  // 确认完成后显示的 title（可选）
}
```

#### InputField - 输入字段配置

```go
type InputField struct {
    Message      string                // 提示消息
    DefaultValue string                // 默认值（可选）
    Validator    Validator             // 验证器（可选）
    ResultTitle  string                // 输入完成后显示的 title（可选）
    Config       *common.PromptConfig  // 自定义配置（可选）
}
```

#### PasswordField - 密码字段配置

```go
type PasswordField struct {
    Message      string                // 提示消息
    DefaultValue string                // 默认值（可选，空字符串表示无默认值）
    Validator    Validator             // 验证器（可选）
    ResultTitle  string                // 输入完成后显示的 title（可选）
    Config       *common.PromptConfig  // 自定义配置（可选）
}
```

#### SelectField - 单选字段配置

```go
type SelectField struct {
    Message      string   // 提示消息
    Options      []string // 选项列表
    DefaultIndex int      // 默认选中的索引
    ResultTitle  string   // 选择完成后显示的 title（可选）
}
```

#### MultiSelectField - 多选字段配置

```go
type MultiSelectField struct {
    Message         string // 提示消息
    Options         []string // 选项列表
    DefaultSelected []int  // 默认选中的索引列表
    ResultTitle     string // 选择完成后显示的 title（可选）
}
```

### 验证器

prompt 模块提供了多种内置验证器：

- `ValidateRegex(pattern string, errorMsg string)` - 基于正则表达式的验证器
- `ValidateEmail()` - 验证邮箱格式
- `ValidateURL()` - 验证 URL 格式
- `ValidateRequired()` - 验证输入不能为空
- `ValidateMinLength(minLen int)` - 验证最小长度
- `ValidateMaxLength(maxLen int)` - 验证最大长度
- `ValidateLength(minLen, maxLen int)` - 验证长度范围

**使用示例**：
```go
// 邮箱验证
email, err := prompt.Input().
    Prompt("请输入邮箱").
    Validate(prompt.ValidateEmail()).
    Run()

// 最小长度验证
password, err := prompt.Password().
    Prompt("请输入密码").
    Validate(prompt.ValidateMinLength(8)).
    Run()

// 正则表达式验证
code, err := prompt.Input().
    Prompt("请输入验证码").
    Validate(prompt.ValidateRegex(`^\d{6}$`, "验证码必须是6位数字")).
    Run()
```

### 函数式 API

所有模块都提供了函数式调用和配置结构体调用两种方式：

#### 函数式 API

- `AskConfirm(field ConfirmField)` - 确认提示
- `AskInput(field InputField)` - 输入提示
- `AskPassword(field PasswordField)` - 密码提示
- `AskSelect(field SelectField)` - 单选提示
- `AskMultiSelect(field MultiSelectField)` - 多选提示

**优势**：
- ✅ 参数清晰，通过字段名访问
- ✅ 可选参数通过零值处理，无需传递
- ✅ 类型安全
- ✅ 易于扩展新字段

### 表单相关函数

- `SetFormFormatResultTitle(formatFunc)` - 设置 Form 的 FormatResultTitle 函数，用于自定义输入完成后显示的 title 格式
- `FormatResultTitleForForm(originalMessage, resultValue)` - 为 Form 格式化完成后显示的 title 的辅助函数，将 "Please enter your X" 转换为 "Your X"

所有字段类型都使用配置结构体：

- `InputFormField` - 输入字段配置
  - `Key string` - 字段键名
  - `Prompt string` - 提示消息
  - `DefaultValue string` - 默认值（可选）
  - `Validator Validator` - 验证器（可选）
  - `ResultTitle string` - 结果标题（可选）
  - `Condition Condition` - 条件函数（可选）

- `PasswordFormField` - 密码字段配置（字段同 InputFormField，DefaultValue 为空字符串表示无默认值）
- `ConfirmFormField` - 确认字段配置（DefaultValue 为 bool 类型）
- `SelectFormField` - 单选字段配置（包含 Options 和 DefaultIndex）
- `MultiSelectFormField` - 多选字段配置（包含 Options 和 DefaultSelected）
- `NestedFormField` - 嵌套表单配置（包含 NestedForm）
- `Run()` - 执行表单并返回结果

### Message

- `Info(format string, args ...interface{})` - 输出信息
- `Success(format string, args ...interface{})` - 输出成功信息
- `Warning(format string, args ...interface{})` - 输出警告信息
- `Error(format string, args ...interface{})` - 输出错误信息
- `Fatal(format string, args ...interface{})` - 输出致命错误并退出
- `Debug(format string, args ...interface{})` - 输出调试信息

### Spinner

- `Start()` - 启动加载动画
- `Stop()` - 停止加载动画
- `UpdateMessage(message string)` - 更新消息文本
- `WithSuccess(message string)` - 停止并显示成功消息
- `WithError(message string)` - 停止并显示错误消息
- `WithInfo(message string)` - 停止并显示信息消息
- `Do(fn func() error)` - 执行函数并显示加载状态

**Spinner 选项函数**：
- `WithSpinner(frames []string)` - 设置自定义 spinner 字符序列
- `WithInterval(d time.Duration)` - 设置更新间隔
- `WithStyle(style lipgloss.Style)` - 设置自定义样式（同时应用于 spinner 和文本）
- `WithSpinnerStyle(style lipgloss.Style)` - 设置 spinner 字符的样式
- `WithMessageStyle(style lipgloss.Style)` - 设置消息文本的样式
- `WithWriter(w io.Writer)` - 设置输出流

### Table

- `AddRow(row []string)` - 添加行
- `SetHeader(headers []string)` - 设置表头
- `SetBorder(border bool)` - 设置边框
- `SetRowLine(rowLine bool)` - 设置行线
- `SetAlignment(align Alignment)` - 设置对齐方式
- `Render()` - 渲染表格

## 配置管理

### PromptConfig

`PromptConfig` 是提示功能的通用配置结构，包含各种格式化函数：

```go
type PromptConfig struct {
    FormatPrompt         func(message string) string
    FormatAnswer         func(value string) string
    FormatError          func(message string) string
    FormatHint           func(message string) string
    FormatQuestionPrefix func() string
    FormatAnswerPrefix   func() string
    FormatResultTitle    func(originalMessage string, resultValue string) string
}
```

### ConfigManager

`ConfigManager` 提供统一的配置管理，支持三层配置优先级：

1. **默认配置**（defaultConfig）：系统默认值
2. **全局配置**（globalConfig）：用户设置的全局配置
3. **局部配置**（localConfig）：每次调用时的局部配置

优先级：`defaultConfig < globalConfig < localConfig`

```go
manager := common.NewConfigManager(defaultConfig)
manager.SetGlobalConfig(globalConfig)
finalConfig := manager.BuildConfig(localConfig)
```

### 配置合并

- `MergeConfig()`：合并两个配置，override 中的非 nil 字段会覆盖 base
- `FillDefaults()`：填充配置的默认值，如果 config 中的字段为 nil，则使用 defaultConfig 中对应的字段
- `BuildConfigWithDefaults()`：构建配置（带默认值填充）
- `BuildConfigWithResultTitle()`：构建配置并设置 ResultTitle

## Fallback 机制

### TypedFallbackHandler

类型安全的 fallback 处理器接口（泛型版本），用于提供类型安全的 fallback 处理：

```go
type TypedFallbackHandler[T any] interface {
    FormatPromptText(message string) string
    FormatAnswer(result T) string
    ProcessLineInput(input string) (T, error)
    GetDefaultResult() T
}
```

### ExecuteFallbackTyped

执行 fallback 模式的通用框架（类型安全版本），使用泛型提供类型安全，避免类型断言：

```go
result, err := common.ExecuteFallbackTyped(
    terminal,
    message,
    config,
    handler,
    options,
)
```

### 选择功能的 Fallback

- `ExecuteSelectFallback()`：执行选择 fallback 的通用框架
- `ExecuteMultiSelectFallback()`：执行多选 fallback 的通用框架

这些函数处理通用的 fallback 流程：格式化提示、显示选项列表、读取输入、解析输入、显示结果。

## 通用辅助函数

### 格式化函数（common/format.go）

- `FormatPromptWithPrefix()`：格式化提示消息并添加前缀
- `FormatResult()`：格式化并显示结果
- `FormatResultWithTitle()`：格式化并显示结果（支持动态 title）
- `FormatResultInline()`：格式化并显示结果（在同一行，用于 confirm 等场景）

### 渲染函数（common/render.go）

- `RenderOptions()`：渲染选项列表的通用函数，用于 select 和 multiselect

### 导航处理（common/navigation.go）

- `NavigationHandler`：导航处理器，提供通用的上下箭头键导航逻辑，支持边界处理

### 输入处理（common/input_handler.go）

- `HandleInteractiveInput()`：处理交互式输入的通用循环，用于 select、multiselect 等需要处理键盘输入的交互式提示

### 选择辅助函数（common/select_helpers.go）

- `SelectSetup` 和 `SetupInteractiveSelect()`：设置交互式选择的通用组件
- `ExecuteSelectFallback()`：执行选择 fallback 的通用框架
- `ExecuteMultiSelectFallback()`：执行多选 fallback 的通用框架

## 注意事项

1. **终端原始模式**：某些功能（如选择提示）需要终端原始模式，如果设置失败会自动 fallback 到普通模式
2. **颜色支持**：默认启用颜色，可通过 `SetTheme()` 关闭颜色（适用于 CI/非 TTY 环境）
3. **验证器**：输入验证支持实时验证和回车验证，验证失败会显示错误并允许重新输入
4. **取消操作**：所有交互式提示都支持 Ctrl+C 取消，会正确处理终端状态恢复
5. **非 TTY 环境**：在非 TTY 环境下，颜色会自动关闭，使用纯文本输出
6. **配置优先级**：配置合并遵循 `defaultConfig < globalConfig < localConfig` 的优先级
7. **Fallback 机制**：所有交互式提示都支持 fallback 到普通输入模式，确保在非交互式环境下的可用性

## 依赖

- `github.com/charmbracelet/lipgloss` - 终端样式库
- `github.com/mattn/go-runewidth` - 字符宽度计算
- `golang.org/x/term` - 终端控制

## 相关文档

- [详细架构文档](../../docs/architecture/prompt.md) - 模块设计思路和实现细节
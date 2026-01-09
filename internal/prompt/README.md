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
│   ├── config.go            # 提示功能通用配置（14行）
│   ├── format.go            # 格式化函数（私有）
│   ├── render.go            # 渲染功能（私有）
│   ├── navigation.go        # 导航功能（键盘方向键处理）
│   ├── input_handler.go     # 输入处理（键盘事件处理）
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
│   └── handler.go          # 确认处理器（键盘事件处理）
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

// 基础输入
value, err := prompt.Input().
    Prompt("请输入您的姓名").
    Run()

// 带默认值和验证的输入
email, err := prompt.Input().
    Prompt("请输入邮箱").
    DefaultValue("user@example.com").
    Validate(prompt.ValidateEmail()).
    Run()

// 密码输入
password, err := prompt.Password().
    Prompt("请输入密码").
    Validate(prompt.ValidateMinLength(8)).
    Run()
```

### 确认提示

```go
import "github.com/zevwings/workflow/internal/prompt"

// 函数式调用
confirmed, err := prompt.AskConfirm("是否继续？", true)

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

// 单选
index, err := prompt.Select().
    Prompt("请选择一个选项").
    Options(options).
    Default(0).
    Run()

// 多选
selected, err := prompt.MultiSelect().
    Prompt("请选择多个选项").
    Options(options).
    Default([]int{0, 2}).
    Run()
```

### 表单提示

```go
import "github.com/zevwings/workflow/internal/prompt"

result, err := prompt.Form().
    AddInput("name", "姓名", "", prompt.ValidateRequired()).
    AddInput("email", "邮箱", "", prompt.ValidateEmail()).
    AddPassword("password", "密码", prompt.ValidateMinLength(8)).
    AddSelect("role", "角色", []string{"开发者", "测试"}, 0).
    AddConfirm("agree", "同意协议", false).
    Run()

if err != nil {
    return err
}

name := result.GetString("name")
email := result.GetString("email")
roleIndex := result.GetInt("role")
agree := result.GetBool("agree")
```

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

spinner := prompt.NewSpinner("正在处理...")
spinner.Start()
defer spinner.Stop()

// 执行操作
time.Sleep(2 * time.Second)

// 停止并显示成功消息
spinner.WithSuccess("处理完成")

// 或使用 Do 方法
spinner := prompt.NewSpinner("正在处理...")
err := spinner.Do(func() error {
    // 执行操作
    return nil
})
```

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

- `AddInput(key, prompt, defaultValue string, validator Validator)` - 添加输入字段
- `AddPassword(key, prompt string, validator Validator)` - 添加密码字段
- `AddConfirm(key, prompt string, defaultValue bool)` - 添加确认字段
- `AddSelect(key, prompt string, options []string, defaultIndex int)` - 添加单选字段
- `AddMultiSelect(key, prompt string, options []string, defaultSelected []int)` - 添加多选字段
- `AddForm(key, prompt string, nestedForm *FormBuilder)` - 添加嵌套表单字段
- `Condition(condition Condition)` - 为最后一个字段设置条件函数
- `Validate(validator FormValidator)` - 设置表单级验证器
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

### Table

- `AddRow(row []string)` - 添加行
- `SetHeader(headers []string)` - 设置表头
- `SetBorder(border bool)` - 设置边框
- `SetRowLine(rowLine bool)` - 设置行线
- `SetAlignment(align Alignment)` - 设置对齐方式
- `Render()` - 渲染表格

## 注意事项

1. **终端原始模式**：某些功能（如选择提示）需要终端原始模式，如果设置失败会自动 fallback 到普通模式
2. **颜色支持**：默认启用颜色，可通过 `SetTheme()` 关闭颜色（适用于 CI/非 TTY 环境）
3. **验证器**：输入验证支持实时验证和回车验证，验证失败会显示错误并允许重新输入
4. **取消操作**：所有交互式提示都支持 Ctrl+C 取消，会正确处理终端状态恢复
5. **非 TTY 环境**：在非 TTY 环境下，颜色会自动关闭，使用纯文本输出

## 依赖

- `github.com/charmbracelet/lipgloss` - 终端样式库
- `github.com/mattn/go-runewidth` - 字符宽度计算
- `golang.org/x/term` - 终端控制

## 相关文档

- [详细架构文档](../../docs/architecture/prompt.md) - 模块设计思路和实现细节
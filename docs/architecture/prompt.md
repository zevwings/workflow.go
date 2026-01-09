# prompt 模块架构文档

## 📋 概述

prompt 模块是 Workflow CLI 的基础模块，提供统一的用户交互功能，包括输入提示、选择提示、表单提示、消息输出等。该模块专注于终端交互式 UI 的实现，不涉及命令层的业务逻辑。

prompt 模块提供完整的交互式提示功能，包括输入（Input/Password）、确认（Confirm）、选择（Select/MultiSelect）、表单（Form）、消息输出（Message）、加载指示器（Spinner）、表格显示（Table）等，总代码行数约 4624+ 行。

**模块统计：**
- 代码行数：约 4624+ 行（不含测试文件）
- 主要文件：20+ 个核心文件
- 主要结构体：`InputBuilder`、`ConfirmBuilder`、`SelectBuilder`、`MultiSelectBuilder`、`FormBuilder`、`Message`、`Spinner`、`Table`、`Theme`
- 支持功能：输入提示、密码输入、确认提示、单选、多选、表单、消息输出、加载指示器、表格显示、主题配置

**注意**：本模块是基础库模块，其他模块通过导入使用。模块内部采用分层设计，通过子模块实现功能解耦。

---

## 📁 模块架构（核心业务逻辑）

prompt 模块（`internal/prompt/`）是 Workflow CLI 的基础库模块，提供统一的用户交互功能。该模块专注于终端交互式 UI 的实现，提供简洁易用的 API，不涉及命令层的业务逻辑。

### 模块结构

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
│   ├── handler.go           # 确认处理器（键盘事件处理）
│   └── adapter.go          # Fallback 适配器（confirmFallbackAdapter）
│
├── select/                   # 单选子模块
│   ├── core.go              # 选择核心逻辑
│   └── handler.go           # 选择处理器（键盘事件处理）
│
├── multiselect/              # 多选子模块
│   ├── core.go              # 多选核心逻辑
│   └── handler.go           # 多选处理器（键盘事件处理）
│
├── form/                     # 表单子模块
│   ├── builder.go           # 表单构建器（链式 API）
│   ├── executor.go          # 表单执行器（执行表单流程）
│   ├── field.go             # 表单字段定义
│   ├── result.go            # 表单结果定义
│   ├── validator.go         # 表单验证器
│   └── config.go            # 表单配置（格式化函数注入）
│
└── io/                       # I/O 抽象模块
    ├── terminal.go          # 终端 I/O 接口定义
    ├── stdterminal.go       # 标准终端实现
    ├── mockterminal.go      # Mock 终端实现（用于测试）
    ├── rawmode.go           # 原始模式控制
    ├── renderer.go          # 渲染器（ANSI 转义序列）
    └── escape.go            # ANSI 转义序列工具
```

**总计：约 4624+ 行代码**

### 依赖模块

- **`github.com/charmbracelet/lipgloss`**：终端样式库
  - 用于颜色、样式渲染
- **`github.com/mattn/go-runewidth`**：字符宽度计算
  - 用于表格列宽计算（支持中文等多字节字符）
- **`golang.org/x/term`**：终端控制
  - 用于原始模式设置、终端状态管理

### 模块集成

- **`internal/commands/`**：命令层使用 prompt 功能
  - `prompt.Input()` - 输入提示
  - `prompt.Confirm()` - 确认提示
  - `prompt.Select()` - 选择提示
  - `prompt.Form()` - 表单提示
  - `prompt.Message` - 消息输出
  - `prompt.Spinner` - 加载指示器
  - `prompt.Table` - 表格显示

---

## 🏗️ 架构设计

### 设计原则

1. **Builder 模式**：所有提示功能都支持链式调用，提供流畅的 API
2. **函数式与 Builder 并存**：既支持函数式调用（向后兼容），也支持 Builder 模式（推荐）
3. **配置注入**：通过配置注入避免循环依赖，实现模块解耦
4. **I/O 抽象**：通过 `TerminalIO` 接口抽象终端操作，便于测试和扩展
5. **主题统一**：通过 `Theme` 统一管理所有 UI 样式，支持颜色开关

### 核心组件

#### 0. 配置管理组件 (`common/config.go`, `common/config_manager.go`)

**职责**：提供统一的配置管理和格式化功能

**主要结构**：
- `PromptConfig`：提示功能的通用配置结构
- `BasePromptConfig`：基础提示配置（通用参数）
- `ConfigManager`：配置管理器，支持三层配置优先级

**关键特性**：
- 支持默认配置、全局配置、局部配置的层次结构
- 配置合并和默认值填充
- 灵活的配置覆盖机制

**使用场景**：
- 统一管理所有提示功能的配置
- 支持全局样式配置
- 支持局部配置覆盖

#### 1. Input/Password 组件 (`input.go`)

**职责**：提供文本输入和密码输入功能

**主要方法**：
- `Input()` - 创建输入构建器
- `Password()` - 创建密码构建器
- `AskInput()` - 函数式输入调用
- `AskPassword()` - 函数式密码调用

**关键特性**：
- 支持默认值、占位符
- 支持实时验证和回车验证
- 支持光标移动（字符级输入）
- 密码模式使用星号掩码
- 错误提示自动清除和重试
- 支持自定义配置（`Config` 字段）
- 支持结果标题（`ResultTitle` 字段）

**使用场景**：
- 用户输入配置信息
- 密码输入
- 带验证的输入（邮箱、URL 等）

#### 2. Confirm 组件 (`confirm.go`)

**职责**：提供确认提示功能（Yes/No）

**主要方法**：
- `Confirm()` - 创建确认构建器
- `AskConfirm()` - 函数式确认调用

**关键特性**：
- 支持默认值（默认 Yes 或 No）
- 键盘导航（Y/N 键、方向键）
- 实时显示选择状态
- 支持 fallback 适配器（`confirmFallbackAdapter`）

**使用场景**：
- 操作确认
- 危险操作二次确认

#### 3. Select 组件 (`select.go`)

**职责**：提供单选功能

**主要方法**：
- `Select()` - 创建选择构建器
- `AskSelect()` - 函数式选择调用

**关键特性**：
- 支持选项列表
- 支持默认选中索引
- 键盘导航（方向键、Enter 确认）
- 实时高亮显示
- 支持 fallback 模式（`ExecuteSelectFallback`）
- 使用通用辅助函数（`SelectSetup`、`RenderOptions`）

**使用场景**：
- 从多个选项中选择一个
- 配置项选择

#### 4. MultiSelect 组件 (`multiselect.go`)

**职责**：提供多选功能

**主要方法**：
- `MultiSelect()` - 创建多选构建器
- `AskMultiSelect()` - 函数式多选调用

**关键特性**：
- 支持多选项选择
- 支持默认选中索引列表
- 键盘导航（方向键移动、空格切换、Enter 确认）
- 实时显示选中状态
- 支持 fallback 模式（`ExecuteMultiSelectFallback`）
- 使用通用辅助函数（`SelectSetup`、`RenderOptions`）

**使用场景**：
- 从多个选项中选择多个
- 批量操作选择

#### 5. Form 组件 (`form.go`)

**职责**：提供表单功能（组合多个字段）

**主要方法**：
- `Form()` - 创建表单构建器
- `AskForm()` - 函数式表单调用
- `SetFormFormatResultTitle()` - 设置 Form 的 FormatResultTitle 函数
- `FormatResultTitleForForm()` - 格式化结果标题的辅助函数

**关键特性**：
- 支持多种字段类型（Input、Password、Confirm、Select、MultiSelect、嵌套 Form）
- 支持字段条件显示（Condition）
- 支持表单级验证
- 支持嵌套表单
- 支持自定义结果标题格式化

**使用场景**：
- 复杂配置表单
- 多步骤输入流程
- 条件字段显示

#### 6. Message 组件 (`message.go`)

**职责**：提供消息输出功能

**主要方法**：
- `Info()` - 输出信息
- `Success()` - 输出成功信息
- `Warning()` - 输出警告信息
- `Error()` - 输出错误信息
- `Fatal()` - 输出致命错误并退出
- `Debug()` - 输出调试信息（需 verbose 模式）

**关键特性**：
- 支持不同级别的消息（Info、Success、Warning、Error、Debug）
- 使用主题样式渲染
- 支持格式化字符串

**使用场景**：
- 操作结果提示
- 错误信息显示
- 调试信息输出

#### 7. Spinner 组件 (`spinner.go`)

**职责**：提供加载指示器功能

**主要方法**：
- `NewSpinner()` - 创建加载指示器
- `Start()` - 启动动画
- `Stop()` - 停止动画
- `WithSuccess()` - 停止并显示成功消息
- `WithError()` - 停止并显示错误消息
- `Do()` - 执行函数并显示加载状态

**关键特性**：
- 支持自定义动画帧
- 支持自定义样式（spinner 和消息可分别设置）
- 自动隐藏/显示光标
- 支持后台运行（goroutine）

**使用场景**：
- 长时间操作提示
- 异步任务状态显示

#### 8. Table 组件 (`table.go`)

**职责**：提供表格显示功能

**主要方法**：
- `NewTable()` - 创建表格
- `AddRow()` - 添加行
- `Render()` - 渲染表格
- `SetBorder()` - 设置边框
- `SetAlignment()` - 设置对齐方式

**关键特性**：
- 支持边框显示/隐藏
- 支持行分隔线
- 支持对齐方式（左、中、右）
- 自动计算列宽（支持多字节字符）
- 支持 ANSI 代码（颜色）在单元格中

**使用场景**：
- 数据列表显示
- 配置信息展示
- 结果汇总

#### 9. Theme 组件 (`theme.go`)

**职责**：提供主题配置功能

**主要方法**：
- `SetTheme()` - 设置全局主题
- `GetTheme()` - 获取当前主题

**关键特性**：
- 统一管理所有 UI 样式
- 支持颜色开关（EnableColor）
- 线程安全（使用互斥锁）
- 支持多种样式（Info、Success、Warning、Error、Debug、Title、Answer、Hint、Border）

**使用场景**：
- 全局样式配置
- CI/非 TTY 环境关闭颜色

### 设计模式

#### 1. Builder 模式

**实现**：所有提示功能都提供 Builder 模式，支持链式调用

**优势**：
- 提供流畅的 API
- 支持可选参数
- 代码可读性强

**示例**：
```go
result, err := prompt.Input().
    Prompt("请输入邮箱").
    DefaultValue("user@example.com").
    Validate(prompt.ValidateEmail()).
    Run()
```

#### 2. 配置管理模式

**实现**：通过 `ConfigManager` 统一管理配置，支持三层配置优先级

**优势**：
- 支持默认配置、全局配置、局部配置的层次结构
- 配置合并和默认值填充
- 灵活的配置覆盖机制

**示例**：
```go
manager := common.NewConfigManager(defaultConfig)
manager.SetGlobalConfig(globalConfig)
finalConfig := manager.BuildConfig(localConfig)
```

**配置优先级**：`defaultConfig < globalConfig < localConfig`

#### 3. Fallback 模式（类型安全）

**实现**：通过 `TypedFallbackHandler` 接口和 `ExecuteFallbackTyped` 提供类型安全的 fallback 处理

**优势**：
- 类型安全，避免类型断言
- 统一的 fallback 处理框架
- 支持泛型，代码复用

**示例**：
```go
type TypedFallbackHandler[T any] interface {
    FormatPromptText(message string) string
    FormatAnswer(result T) string
    ProcessLineInput(input string) (T, error)
    GetDefaultResult() T
}

result, err := common.ExecuteFallbackTyped(
    terminal,
    message,
    config,
    handler,
    options,
)
```

#### 4. 适配器模式

**实现**：通过适配器将特定 Handler 适配为通用接口

**优势**：
- 解耦具体实现和通用框架
- 便于扩展新的提示类型
- 代码复用

**示例**：
```go
// confirmFallbackAdapter 将 ConfirmHandler 适配为 TypedFallbackHandler[bool]
adapter := newConfirmFallbackAdapter(handler)
result, err := common.ExecuteFallbackTyped(terminal, message, config, adapter, options)
```

#### 5. 配置注入模式

**实现**：通过 `PromptConfig` 注入格式化函数，避免循环依赖

**优势**：
- 解耦模块依赖
- 支持自定义格式化
- 便于测试

**示例**：
```go
form.SetPromptConfig(common.PromptConfig{
    FormatPrompt:         formatTitle,
    FormatAnswer:         formatAnswer,
    FormatError:          formatError,
    FormatHint:           formatHint,
    FormatQuestionPrefix: formatQuestionPrefix,
    FormatAnswerPrefix:   formatAnswerPrefix,
})
```

#### 6. I/O 抽象模式

**实现**：通过 `TerminalIO` 接口抽象终端操作

**优势**：
- 便于测试（Mock 实现）
- 支持不同终端实现
- 解耦终端操作

**示例**：
```go
type TerminalIO interface {
    ReadByte() (byte, error)
    Print(s string)
    MakeRaw() (*term.State, error)
    // ...
}
```

#### 7. 函数式与 Builder 并存

**实现**：既提供函数式调用（`AskXxx`），也提供 Builder 模式（`Xxx()`）

**优势**：
- 向后兼容
- 灵活使用
- 渐进式迁移

### 错误处理

#### 分层错误处理

1. **输入验证层**：处理输入验证错误
   - 实时验证：输入时即时反馈
   - 回车验证：按 Enter 时验证
   - 错误提示：显示红色错误信息
   - 自动重试：验证失败后继续输入

2. **终端操作层**：处理终端操作错误
   - 原始模式设置失败：fallback 到普通模式
   - 读取错误：返回错误信息
   - 恢复终端状态：确保终端状态正确恢复

3. **业务逻辑层**：处理业务相关错误
   - 表单验证失败：显示错误并允许重试
   - 取消操作：Ctrl+C 正确处理

#### 容错机制

- **终端原始模式失败**：fallback 到普通输入模式
- **输入验证失败**：清除错误提示，允许重新输入
- **取消操作**：正确处理 Ctrl+C，恢复终端状态
- **非 TTY 环境**：自动关闭颜色，使用纯文本输出

#### Fallback 机制

所有交互式提示都支持 fallback 机制，确保在非交互式环境下的可用性：

1. **类型安全的 Fallback**：
   - 使用 `TypedFallbackHandler[T]` 接口提供类型安全
   - 通过 `ExecuteFallbackTyped` 执行 fallback 流程
   - 支持泛型，避免类型断言

2. **选择功能的 Fallback**：
   - `ExecuteSelectFallback`：处理单选 fallback
   - `ExecuteMultiSelectFallback`：处理多选 fallback
   - 统一的 fallback 流程：格式化提示、显示选项、读取输入、解析输入、显示结果

3. **适配器模式**：
   - 通过适配器将特定 Handler 适配为通用接口
   - 例如：`confirmFallbackAdapter` 将 `ConfirmHandler` 适配为 `TypedFallbackHandler[bool]`

---

## 🔄 集成关系

### 模块使用关系

prompt 模块被以下模块使用：

1. **`internal/commands/`**：命令层使用 prompt 功能
   - 使用 `prompt.Input()` - 用户输入配置
   - 使用 `prompt.Confirm()` - 操作确认
   - 使用 `prompt.Select()` - 选项选择
   - 使用 `prompt.Form()` - 复杂表单
   - 使用 `prompt.Message` - 消息输出
   - 使用 `prompt.Spinner` - 加载提示
   - 使用 `prompt.Table` - 数据展示

### 调用流程

#### 输入提示流程

```
命令层 (commands/)
  ↓
prompt.Input().Prompt("消息").DefaultValue("默认值").Run()
  ↓
inputFunc() - 统一输入函数
  ↓
input.ReadLineCoreDefault() / input.ReadWithPlaceholderDefault()
  ↓
io.TerminalIO - 终端 I/O 抽象
  ↓
input.Editor - 字符级输入编辑器
  ↓
input.Handler - 键盘事件处理
  ↓
返回输入结果
```

#### 表单执行流程

```
命令层 (commands/)
  ↓
prompt.Form().AddInput().AddSelect().Run()
  ↓
form.NewFormExecutor().Execute(builder)
  ↓
form.Executor - 遍历字段执行
  ↓
根据字段类型调用对应提示：
  - Input → prompt.AskInput()
  - Password → prompt.AskPassword()
  - Confirm → prompt.AskConfirm()
  - Select → prompt.AskSelect()
  - MultiSelect → prompt.AskMultiSelect()
  - Form → 递归执行嵌套表单
  ↓
收集结果并验证
  ↓
返回 FormResult
```

#### 选择提示流程

```
命令层 (commands/)
  ↓
prompt.Select().Prompt("消息").Options(options).Run()
  ↓
select.SelectDefault() - 选择核心逻辑
  ↓
select.Handler - 键盘事件处理
  ↓
common.Navigation - 方向键导航
  ↓
common.Render - 渲染选项列表
  ↓
io.TerminalIO - 终端 I/O 操作
  ↓
返回选中索引
```

---

## 🎯 核心功能

### 1. 输入提示功能

**功能说明**：提供文本输入和密码输入功能，支持默认值、占位符、验证等

**流程**：
1. 显示提示消息和默认值（如果有）
2. 显示输入框（带 "> " 前缀）
3. 用户输入（支持光标移动）
4. 实时验证或回车验证
5. 验证失败显示错误并重试
6. 验证通过显示结果并返回

**示例**：
```go
import "github.com/zevwings/workflow/internal/prompt"

// 函数式调用
value, err := prompt.AskInput("请输入邮箱", "user@example.com", prompt.ValidateEmail())

// Builder 模式调用
value, err := prompt.Input().
    Prompt("请输入邮箱").
    DefaultValue("user@example.com").
    Placeholder("example@domain.com").
    Validate(prompt.ValidateEmail()).
    Run()
```

### 2. 确认提示功能

**功能说明**：提供 Yes/No 确认功能

**流程**：
1. 显示提示消息
2. 显示选项（Yes/No）和默认值
3. 用户使用键盘选择（Y/N 键或方向键）
4. 实时高亮显示选择
5. 按 Enter 确认并返回

**示例**：
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

### 3. 选择提示功能

**功能说明**：提供单选功能，从多个选项中选择一个

**流程**：
1. 显示提示消息
2. 显示选项列表（带默认选中）
3. 用户使用方向键导航
4. 实时高亮显示当前选项
5. 按 Enter 确认并返回选中索引

**示例**：
```go
import "github.com/zevwings/workflow/internal/prompt"

options := []string{"选项1", "选项2", "选项3"}

// 函数式调用
index, err := prompt.AskSelect("请选择", options, 0)

// Builder 模式调用
index, err := prompt.Select().
    Prompt("请选择").
    Options(options).
    Default(0).
    Run()
```

### 4. 表单功能

**功能说明**：提供表单功能，组合多个字段进行输入

**流程**：
1. 创建表单构建器
2. 添加字段（Input、Password、Confirm、Select、MultiSelect、嵌套 Form）
3. 设置字段条件（Condition）
4. 执行表单
5. 按顺序执行每个字段（根据条件决定是否显示）
6. 收集结果并验证
7. 返回表单结果

**示例**：
```go
import "github.com/zevwings/workflow/internal/prompt"

result, err := prompt.Form().
    AddInput("name", "姓名", "", prompt.ValidateRequired()).
    AddInput("email", "邮箱", "", prompt.ValidateEmail()).
    AddSelect("role", "角色", []string{"开发者", "测试", "产品"}, 0).
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

---

## 📋 使用示例

### 输入提示示例

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

### 选择提示示例

```go
import "github.com/zevwings/workflow/internal/prompt"

// 单选
options := []string{"选项1", "选项2", "选项3"}
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

### 表单示例

```go
import "github.com/zevwings/workflow/internal/prompt"

result, err := prompt.Form().
    AddInput("name", "姓名", "", prompt.ValidateRequired()).
    AddInput("email", "邮箱", "", prompt.ValidateEmail()).
    AddPassword("password", "密码", prompt.ValidateMinLength(8)).
    AddSelect("role", "角色", []string{"开发者", "测试"}, 0).
    AddConfirm("agree", "同意协议", false).
    Condition(func(r *prompt.FormResult) bool {
        // 只有同意协议才显示角色选择
        return r.GetBool("agree")
    }).
    Run()
```

### 消息输出示例

```go
import "github.com/zevwings/workflow/internal/prompt"

msg := prompt.NewMessage(true) // verbose 模式

msg.Info("这是一条信息")
msg.Success("操作成功")
msg.Warning("这是一条警告")
msg.Error("这是一条错误")
msg.Debug("这是调试信息") // 仅在 verbose 模式下显示
```

### 加载指示器示例

```go
import "github.com/zevwings/workflow/internal/prompt"

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

### 表格显示示例

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

---

## 📝 扩展性

### 添加新提示类型

1. 在 `prompt/` 目录下创建新的提示文件（如 `custom.go`）
2. 实现 Builder 结构体（嵌入 `baseBuilder`）
3. 实现核心逻辑（可创建子模块）
4. 提供函数式调用和 Builder 模式
5. 实现 fallback 处理（可选，推荐）

**示例**：
```go
// custom.go
type CustomBuilder struct {
    baseBuilder
    // 自定义字段
}

func Custom() *CustomBuilder {
    return &CustomBuilder{}
}

func (b *CustomBuilder) Run() (string, error) {
    // 实现逻辑
    return "", nil
}
```

### 实现 Fallback 支持

1. 创建 Handler 结构体，实现业务逻辑
2. 创建 Fallback 适配器，实现 `TypedFallbackHandler[T]` 接口
3. 在核心逻辑中调用 `ExecuteFallbackTyped`

**示例**：
```go
// handler.go
type CustomHandler struct {
    // 字段
}

func (h *CustomHandler) FormatPromptText(message string) string {
    // 实现
}

func (h *CustomHandler) FormatAnswer(result string) string {
    // 实现
}

func (h *CustomHandler) ProcessLineInput(input string) (string, error) {
    // 实现
}

func (h *CustomHandler) GetDefaultResult() string {
    // 实现
}

// adapter.go
type customFallbackAdapter struct {
    handler *CustomHandler
}

func (a *customFallbackAdapter) FormatPromptText(message string) string {
    return a.handler.FormatPromptText(message)
}

// ... 实现其他方法

// core.go
func Custom(cfg CustomConfig) (string, error) {
    // 尝试交互式模式
    if rawModeMgr.MakeRaw() == nil {
        // 交互式逻辑
    }

    // Fallback 模式
    adapter := &customFallbackAdapter{handler: handler}
    return common.ExecuteFallbackTyped(
        terminal,
        message,
        config,
        adapter,
        options,
    )
}
```

### 添加新验证器

1. 在 `input/validator.go` 中添加验证函数
2. 在 `input.go` 中重新导出（如果需要）

**示例**：
```go
// input/validator.go
func ValidateCustom(pattern string) Validator {
    return func(v string) error {
        // 验证逻辑
        return nil
    }
}

// input.go
func ValidateCustom(pattern string) Validator {
    return input.ValidateCustom(pattern)
}
```

### 内置验证器

prompt 模块提供了以下内置验证器：

- `ValidateRegex(pattern, errorMsg)` - 基于正则表达式的验证器
- `ValidateEmail()` - 验证邮箱格式
- `ValidateURL()` - 验证 URL 格式
- `ValidateRequired()` - 验证输入不能为空
- `ValidateMinLength(minLen)` - 验证最小长度
- `ValidateMaxLength(maxLen)` - 验证最大长度
- `ValidateLength(minLen, maxLen)` - 验证长度范围

### 自定义主题

```go
import "github.com/zevwings/workflow/internal/prompt"
import "github.com/charmbracelet/lipgloss"

customTheme := prompt.Theme{
    InfoStyle: lipgloss.NewStyle().Foreground(lipgloss.Color("blue")),
    SuccessStyle: lipgloss.NewStyle().Foreground(lipgloss.Color("green")),
    EnableColor: true,
}

prompt.SetTheme(customTheme)
```

---

## 📚 相关文档

- [模块 README](../../internal/prompt/README.md) - 基础使用说明
- [开发规范](../../development/references/embed-files.md) - 开发相关文档

---

## ✅ 总结

prompt 模块采用清晰的分层设计和多种设计模式：

1. **分层设计**：通过子模块实现功能解耦（common、input、confirm、select、form、io）
2. **Builder 模式**：所有提示功能都支持链式调用，提供流畅的 API
3. **配置管理**：通过 `ConfigManager` 统一管理配置，支持三层配置优先级
4. **Fallback 机制**：类型安全的 fallback 处理，确保非交互式环境下的可用性
5. **配置注入**：通过配置注入避免循环依赖，实现模块解耦
6. **I/O 抽象**：通过 `TerminalIO` 接口抽象终端操作，便于测试和扩展
7. **主题统一**：通过 `Theme` 统一管理所有 UI 样式，支持颜色开关
8. **通用辅助函数**：提供格式化、渲染、导航、输入处理等通用功能

**设计优势**：
- ✅ 提供流畅的链式 API
- ✅ 支持函数式和 Builder 两种调用方式
- ✅ 模块解耦，便于测试和扩展
- ✅ 统一的配置管理和主题配置
- ✅ 类型安全的 fallback 机制
- ✅ 完善的错误处理和容错机制
- ✅ 丰富的通用辅助函数，减少代码重复

**当前实现状态**：
- ✅ 输入提示（Input/Password）
- ✅ 确认提示（Confirm）
- ✅ 选择提示（Select/MultiSelect）
- ✅ 表单功能（Form）
- ✅ 消息输出（Message）
- ✅ 加载指示器（Spinner）
- ✅ 表格显示（Table）
- ✅ 主题配置（Theme）
- ✅ 配置管理（ConfigManager）
- ✅ Fallback 机制（类型安全）
- ✅ 通用辅助函数（格式化、渲染、导航、输入处理）

---

**最后更新**: 2024-12-19
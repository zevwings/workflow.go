## 自研交互式 CLI UI 设计文档（替代 `survey`）

### 1. 目标与范围

- **目标**：
  - 用一套我们自己维护的、基于标准库的轻量交互库，替代 `github.com/AlecAivazis/survey/v2` 在本项目中的使用。
  - 保持现有命令（尤其是 `setup`）的交互体验大体不变，迁移成本可控。
- **范围（v1）**：
  - 只覆盖当前项目实际用到的交互类型（Input / Password / Confirm / 简单 Select）。
  - 不追求 100% 覆盖 `survey` 的所有高级特性（多选、分页、自动补全、复杂验证等可以留到 v2+）。
  - 依赖尽量为 **零外部依赖**（仅使用标准库和已有基础库，如 `color` / `output`）。

---

### 2. `survey` 已用功能梳理

根据当前代码（特别是 `internal/commands/setup.go`）与规划文档，本项目目前/近期实际需要的 `survey` 功能主要包括：

#### 2.1 Input（文本输入）

- 用途：
  - 录入用户姓名、邮箱（`setup`）
  - 录入 Jira URL、用户名等
- 关键能力：
  - 显示一行提示 `Message: "请输入 xxx:"`
  - 支持 **默认值**（已有配置时）
  - 简单必填校验（目前代码中是“尽量填”，没有强校验）

在 `setup.go` 中的典型用法：

```go
prompt := &survey.Input{
    Message: "请输入您的姓名:",
    Default: cfg.User.Name,
}
```

#### 2.2 Password（隐藏输入）

- 用途：
  - 录入 GitHub Token
  - 录入 Jira API Token
- 关键能力：
  - 输入时不回显明文（或仅显示 `*`）
  - 不需要复杂校验（只要拿到字符串即可）

当前用法：

```go
prompt := &survey.Password{
    Message: "请输入 GitHub Personal Access Token:",
}
```

#### 2.3 Confirm（是/否确认）

- 用途：
  - 检测到已有配置时，询问“是否更新？”
  - 是否配置 GitHub / 是否配置 Jira
- 关键能力：
  - 显示一行提示，带默认值（是/否）
  - 支持直接回车使用默认值

当前用法：

```go
prompt := &survey.Confirm{
    Message: "检测到现有配置文件，是否要更新？",
    Default: false,
}
```

#### 2.4 Select / MultiSelect（暂不在 v1 强依赖）

- 目前代码还未使用 Select / MultiSelect。
- 未来可能的用途：
  - 从 GitHub 账号列表中选择当前账号
  - 选择 LLM Provider
  - 多选分支 / 多选 Jira 任务
- v1 可以不实现或仅预留 API，等有实际需求再补。

---

### 3. 计划支持的功能子集（v1）

#### 3.1 v1 必须支持

- **Input：文本输入（带默认值）**
  - 功能：
    - 显示一行 message
    - 显示默认值（如有）
    - 用户按回车时：
      - 如果输入非空：取用户输入
      - 如果输入为空且存在默认值：取默认值
  - 不做复杂校验逻辑，校验由调用方实现。

- **Password：密码输入**
  - 功能：
    - 显示一行 message
    - 输入过程不显示明文（可以简单用 `*` 替代，或者完全不回显）
    - 支持空输入

- **Confirm：确认（是/否）**
  - 功能：
    - 显示一行 `[Y/n]` 或 `[y/N]` 风格的提示，取决于 `Default`
    - 支持：
      - `y` / `yes` / `Y` / 空输入（配合默认值） → true/false
      - `n` / `no` / `N`
    - 对非法输入可以再次提示一次，超过次数直接使用默认值（防止死循环）

#### 3.2 v1 可选支持（若实现成本低）

- **简单 Select：单选列表**
  - 功能：
    - 打印一个编号列表（1. 选项A、2. 选项B...）
    - 用户输入编号（或直接输入值）选择
  - 初期只需要“普通列表 + 编号输入”，不做复杂的上下键选择渲染。

---

### 4. 自研 UI 包设计（API 草案）

包位置建议：`internal/lib/prompt`

#### 4.1 对外类型与函数

设计目标：**简单函数式 API，而不是复杂结构体**，降低接入成本。

```go
package prompt

// Input 显示一条输入提示，返回用户输入的字符串。
// defaultValue 为空字符串时表示无默认值。
func Input(message string, defaultValue string) (string, error)

// Password 显示一条密码输入提示，输入过程不回显明文。
func Password(message string) (string, error)

// Confirm 显示一个是/否确认，返回 true/false。
// defaultYes 表示回车时的默认选择。
func Confirm(message string, defaultYes bool) (bool, error)

// Select 显示一个简单的编号选择列表，返回选中的索引和值。
// v1 可以是可选实现。
func Select(message string, options []string, defaultIndex int) (index int, value string, err error)

// Theme 表示交互式 UI 的主题配置（颜色、前缀、样式等）
type Theme struct {
    // 基础颜色（可复用 fatih/color 提供的样式）
    InfoColor    func(format string, a ...any) string
    WarnColor    func(format string, a ...any) string
    ErrorColor   func(format string, a ...any) string
    PromptColor  func(format string, a ...any) string
    AnswerColor  func(format string, a ...any) string

    // 前缀符号
    PrefixInfo   string // e.g. "INFO"
    PrefixWarn   string // e.g. "WARN"
    PrefixError  string // e.g. "ERROR"

    // 输入提示样式
    InputBracketLeft   string // e.g. "["
    InputBracketRight  string // e.g. "]"

    // 是否启用颜色（方便 CI / 非 TTY 环境关闭）
    EnableColor bool
}

// SetTheme 设置全局主题（线程不安全版本，CLI 场景通常足够）
func SetTheme(t Theme)

// GetTheme 获取当前主题（只读）
func GetTheme() Theme
```

#### 4.2 行为细节约定

- 所有函数：
  - 读写均通过 `os.Stdin` / `os.Stdout`
  - 返回 `error`，以便上层命令处理 Ctrl+C / I/O 失败等情况

**Input：**
- 展示格式示例：
  - 无默认值：`请输入您的姓名: `
  - 有默认值：`请输入您的姓名 [Zev]: `
- 回车逻辑：
  - 用户有输入 → 返回输入
  - 用户无输入 + 有默认值 → 返回默认值
  - 用户无输入 + 无默认值 → 返回空字符串（调用方自己判断）

**Password：**
- 展示：`请输入 GitHub Personal Access Token: `
- 输入：不回显（或仅回显 `*`，实现时可权衡复杂度）
- 行为：
  - 用户输入后按回车 → 返回字符串
  - 不进行长度校验，调用方自行判断是否为空

**Confirm：**
- 根据 `defaultYes` 渲染不同提示：
  - `defaultYes = true` → `"是否配置 GitHub？ [Y/n] "`
  - `defaultYes = false` → `"是否配置 Jira？ [y/N] "`
- 接受的输入：
  - Yes: `""`（空，使用默认）、`y` / `Y` / `yes` / `YES`
  - No: `n` / `N` / `no` / `NO`
- 对非法输入的策略：
  - 提示一次 `"请输入 y 或 n（回车使用默认值）"`，允许重输
  - 连续 2 次非法输入后，直接返回默认值（并不报错）

**Select（若实现）：**
- 展示示例：

```text
请选择当前 GitHub 账号:
  1) default
  2) work
  3) personal
请输入编号 [1]:
```

- 行为：
  - 用户输入编号 → 返回对应项
  - 用户直接回车 → 使用 `defaultIndex`
  - 非法编号 → 提示一次重输，超过次数返回错误或默认值

---

### 5. 主题（Theme）设计

#### 5.1 设计目标

- 在不引入复杂 TUI 框架的前提下，提供一套 **统一、可配置** 的终端交互风格。
- 避免在各个命令中直接调用 `color.HiGreenString` 等散落的样式逻辑，把风格集中到 `prompt` 的 Theme 中管理。
- 支持后续根据用户偏好或环境（暗色/浅色终端、CI 模式）切换主题。

#### 5.2 主题作用范围

主题主要控制以下元素：

- **提示文案样式**：
  - Input/Password/Confirm/Select 的问题提示颜色与前缀。
- **答案显示样式**：
  - 回显用户选择/输入时的颜色（例如高亮显示选择结果）。
- **状态消息样式（可选）**：
  - 成功 / 警告 / 错误的颜色与前缀（可与 `internal/output` 协调）。
- **交互符号**：
  - 例如：`[Y/n]`、`[default]` 外框、列表前缀 `➜`、`*` 等的符号样式。

#### 5.3 与现有 `output` 的关系

- 现有的 `internal/output` 已经封装了 Info / Warning / Error / Success 等输出。
- 为避免重复设计，`prompt` 的 Theme 可以：
  - **复用 `color` 层级**（比如高亮颜色方案）
  - 但在逻辑上保持独立：`prompt` 专注“交互问题/回答”，`output` 专注“普通日志/提示”。
- 后续可以考虑在更高层统一注入 Theme：
  - 例如在 CLI 初始化时根据环境变量或配置选择“暗色主题 / 浅色主题 / 无色主题”。

#### 5.4 默认主题示意

默认主题（Pseudo-code）：

```go
var defaultTheme = Theme{
    InfoColor:   color.HiCyanString,
    WarnColor:   color.HiYellowString,
    ErrorColor:  color.HiRedString,
    PromptColor: color.HiCyanString,
    AnswerColor: color.HiGreenString,

    PrefixInfo:  "",
    PrefixWarn:  "!",
    PrefixError: "x",

    InputBracketLeft:  "[",
    InputBracketRight: "]",

    EnableColor: true,
}
```

在 `prompt` 内部使用方式示例：

```go
func formatPrompt(message string) string {
    t := GetTheme()
    if !t.EnableColor || t.PromptColor == nil {
        return message
    }
    return t.PromptColor("%s", message)
}
```

#### 5.5 主题切换策略

- v1：只提供 **全局唯一主题**：
  - CLI 启动时根据配置/环境变量设置一次 `prompt.SetTheme(...)`。
  - 后续所有交互都使用同一主题。
- v2（可选）：支持**按上下文覆盖**：
  - 提供 `WithTheme(ctx, theme)` 风格的 API 或在函数参数中可选传入自定义 Theme。
  - 适合在个别命令中临时改变风格（例如“高危操作”使用更醒目的颜色）。

---

### 6. 与现有代码的映射关系

以 `internal/commands/setup.go` 为例：

- **原来：**

```go
var update bool
prompt := &survey.Confirm{
    Message: "检测到现有配置文件，是否要更新？",
    Default: false,
}
if err := survey.AskOne(prompt, &update); err != nil {
    return err
}
```

- **迁移后（目标）：**

```go
update, err := prompt.Confirm("检测到现有配置文件，是否要更新？", false)
if err != nil {
    return err
}
```

类似地：

- `survey.Input` → `prompt.Input`
- `survey.Password` → `prompt.Password`
- 未来的 `survey.Select` → `prompt.Select`

通过这种方式，我们可以：

- **集中管理交互逻辑**（样式、颜色、国际化等都可以在 `prompt` 包里逐步演进）
- 给未来的替换留出空间（如果以后想换成 `promptui` 或其他库，只需在 `prompt` 包内部调整）

---

### 7. 实现细节草案（v1）

#### 7.1 Input / Password 基本实现思路

- 使用 `bufio.NewReader(os.Stdin)` 读取一行输入：

```go
reader := bufio.NewReader(os.Stdin)
line, err := reader.ReadString('\n')
line = strings.TrimRight(line, "\r\n")
```

- Password：
  - v1 可以先用简单方案：直接读输入，不回显（macOS/Linux 下可用 `terminal` 或 `golang.org/x/term` 暂时已在依赖中）
  - 为保持简单，可以优先实现“关闭回显”的版本，不做复杂编辑行为。

#### 7.2 Confirm 实现思路

- 循环最多 2 次：
  - 打印 `[Y/n]` 风格提示
  - 读取输入，转换为小写，匹配 `""/y/yes/n/no`
  - 返回 true/false
  - 否则打印错误提示，重试
- 超过次数：返回默认值，不返回错误。

---

### 8. 演进计划

#### v1（当前阶段）

- 实现 `prompt.Input` / `prompt.Password` / `prompt.Confirm`
- 将 `setup` 命令中的 `survey` 使用全部替换为 `prompt` 包
- 保持用户体验基本一致

#### v2（视需求而定）

- 增加 `prompt.Select`（简单编号选择）
- 考虑 MultiSelect（多选）和更丰富的校验机制
- 与 `color` / `output` 更紧密集成（统一风格、支持静默/非交互模式等）

#### v3（长期规划）

- 如果需要更复杂的交互（例如分页、多步表单），可以：
  - 在 `prompt` 下增加高级组件（但仍保持对调用方的简单 API）
  - 或在某些高交互场景下，单独启用 TUI（如 Bubble Tea / tview），但不影响日常命令。

---

### 9. 决策小结

- **不再依赖 `survey` 的维护状态**：通过自研轻量 `prompt` 包替代核心功能。
- **优先覆盖当前使用场景**：Input / Password / Confirm。
- **通过 `internal/lib/prompt` 进行封装**：
  - 调用方只依赖我们自己的 API
  - 未来可在内部按需调整实现（自研、标准库、第三方库等）



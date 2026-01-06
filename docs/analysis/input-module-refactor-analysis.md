## `input` / `input_builder` 模块结构与重构分析

### 1. 背景与目标

`internal/lib/prompt/input.go` 与 `input_builder.go` 目前承担了终端输入交互、占位符处理、密码输入、错误提示以及验证器等多种职责。随着功能不断增加，`input.go` 文件长度已经超过 1000 行，可读性与维护成本明显上升。

本分析文档的目标：

- **梳理现状职责与结构**；
- **识别主要复杂度与风险点**；
- **提出可以逐步落地的简化 / 重构方向**。

---

### 2. 现有结构与职责划分

#### 2.1 `input_builder.go`

- 提供链式构建接口：
  - `Input()` / `Password()` 返回 `*InputBuilder`。
  - `Prompt` / `DefaultValue` / `Placeholder` / `Validate` 设置参数。
  - `Run()` 内部调用统一的 `input(...)` 函数。
- 职责：
  - 面向调用方的易用 API，用来组装输入配置。
  - 不直接处理终端 IO 与光标控制。

整体来看，该文件职责单一、结构清晰，是当前设计中较为稳定的一层。

#### 2.2 `input.go`

`input.go` 集中了几乎所有输入相关的实现细节，可分为几大块：

- **对外接口与入口逻辑**
  - `input(message, defaultValue, placeholder, isPassword, validator)`：
    - 打印提示（Title），处理默认值展示（普通输入与密码输入不同逻辑）。
    - 控制重试逻辑与错误提示清理（通过 `hasError`）。
    - 根据模式选择不同的底层输入函数：
      - 普通输入：`readInputWithPlaceholder` / `readInputWithoutPlaceholder`。
      - 密码输入：`readLineCore`，通过 `echo` 函数实现密文显示。
    - 验证成功后，格式化并打印最终答案。
  - `AskInput` / `AskPassword`：包装函数式调用，保持向后兼容。

- **终端输入与编辑内核**
  - `readInputWithPlaceholder`：
    - 支持 placeholder 的输入框。
    - 实时错误提示：在输入行下方渲染/清除错误信息。
    - 字节级读取、解析 ESC 序列（方向键）、Backspace、Ctrl+C。
    - 管理 `cursorPos` / `placeholderDisplayed` / `errorLineExists` 等状态。
  - `readLineCore`：
    - 通用单行编辑内核，不带 placeholder。
    - 通过 `echo([]byte) string` 控制显示内容（明文或密文）。
    - 支持左右方向键、删除、错误提示行管理等。
  - `readInputWithoutPlaceholder`：
    - `readLineCore` 的明文封装。
  - `readPasswordInput`：
    - 使用 `term.MakeRaw` + `os.Stdin.Read` 的老式密码输入实现（显示 `*`）。
    - 当前主路径已改为基于 `readLineCore` 实现，存在遗留嫌疑。

- **光标与行控制辅助函数**
  - `moveCursorToPosition`：按显示宽度向左移动光标到指定位置。
  - `clearLine`：清除当前整行。
  - `clearErrorLine`：假定错误行在输入行下一行，负责清理。
  - `showError`：在输入行下方打印错误，并重绘输入行与光标。

- **格式化与主题相关**
  - `formatPlaceholder`：基于主题 `HintStyle` 渲染 placeholder，带斜体。
  - `formatErrorWithLipgloss`：目前简单转发到 `formatError`。
  - 同时依赖外部的 `formatPrompt` / `formatAnswer` 等函数。

- **校验器工具集合**
  - `Validator` 类型：`func(value string) error`。
  - 一组通用校验器：
    - `ValidateRegex` / `ValidateEmail` / `ValidateURL`；
    - `ValidateRequired`；
    - `ValidateMinLength` / `ValidateMaxLength` / `ValidateLength`。
  - 这些逻辑与终端 IO 无关，更偏向业务/工具层。

---

### 3. 主要复杂度与问题点

#### 3.1 终端控制逻辑分散且重复

- `readInputWithPlaceholder` 与 `readLineCore` 都在做：
  - 切换终端到 raw 模式（`term.MakeRaw` / `term.Restore`）。
  - 循环从 `os.Stdin` 读取字节。
  - 解析 ESC 序列处理方向键。
  - 处理回车、Backspace、Ctrl+C 等控制键。
  - 显示/清理错误提示行。
- 两者在整体流程上高度相似，但各自维护一套状态：
  - `cursorPos`；
  - `errorLineExists`；
  - `placeholderDisplayed`（仅 placeholder 版本）。
- 结果：
  - **功能调整需要多处修改**，非常容易出现一处更新、另一处漏改的问题。

#### 3.2 错误提示与重试状态跨层分散

- 外层 `input` 使用 `hasError` 控制“新一轮输入前先清理上一轮错误提示和输入行”。
- 内层编辑循环中又分别维护 `errorLineExists`（`readInputWithPlaceholder` / `readLineCore` 内部）。
- 错误的展示/清除路径涉及：
  - `input` 里的光标上移与行清理；
  - `showError` / `clearErrorLine` 中的换行、上移、重绘。
- 这些逻辑构成了一个隐式的“行布局状态机”，分散在多个函数里，理解和修改都比较困难。

#### 3.3 编辑内核抽象不彻底

- 目前已经有 `readLineCore(promptText, validator, echo)` 作为“通用单行编辑器”雏形：
  - 通过 `echo` 控制显示内容，可以支持明文/密文等模式。
  - 内部包含完整的按键处理和错误提示逻辑。
- 但 placeholder 模式并没有基于该内核实现，而是重新实现了一个大循环：
  - 再次处理 ESC、Backspace、Ctrl+C、错误行渲染等。
  - 导致行为不易统一，维护成本翻倍。

#### 3.4 校验器与 IO 混在同一大文件

- `Validator` 及 `ValidateXxx` 系列位于 `input.go` 尾部，与终端交互代码紧挨。
- 校验器本质上是通用业务逻辑，完全可以独立存在并被其他 prompt 类型复用。
- 混在一起使 `input.go` 文件行数暴涨，阅读时不易聚焦。

#### 3.5 可测试性不足

- 所有核心逻辑都直接依赖：
  - `os.Stdin` / `fmt.Print`；
  - `term.MakeRaw` / `term.Restore`。
- 没有抽象出 `Terminal` 或 IO 接口，导致：
  - 无法在测试中注入虚拟输入与捕获输出；
  - 只能依赖手工运行 demo 命令观察行为。

#### 3.6 可能存在遗留实现

- `readPasswordInput` 与 `readPasswordFallback`：
  - 当前密码逻辑主路径已经改为 `readLineCore` + `echo("*")`。
  - 需要确认是否仍有调用；若无，则属于遗留代码，应考虑移除或收敛。

---

### 4. 重构与简化建议

本节侧重结构和演进方向，尽量控制在“低风险、可分步执行”的前提下。

#### 4.1 按职责拆分文件

在不改变对外 API 的前提下，可以先做“物理拆分”：

- `input_public.go`
  - 对外接口与类型定义：
    - `Input` / `Password` / `AskInput` / `AskPassword`；
    - `Validator` 类型别名。
  - 保持调用方使用方式不变。

- `input_editor.go`
  - 核心编辑内核：
    - `readLineCore` 及其内部按键处理逻辑；
    - 与错误提示行渲染直接相关的函数（可视情况保留或内聚）。

- `input_placeholder.go`
  - 与 placeholder 相关的输入逻辑：
    - 当前的 `readInputWithPlaceholder`；
    - 或未来基于统一内核的 placeholder 封装。

- `input_validators.go`
  - `ValidateRegex` / `ValidateEmail` / `ValidateURL`；
  - `ValidateRequired` / `ValidateMinLength` / `ValidateMaxLength` / `ValidateLength`。

预计收益：

- 单文件行数明显下降，阅读体验更好；
- 修改终端行为时更容易只关注 `input_editor.go`；
- validators 可以被其他 prompt 类型自然复用。

#### 4.2 引入显式“编辑器状态机”

可以考虑将当前分散在多个函数中的状态集中到一个小 struct，例如：

```go
type lineEditor struct {
    value       []rune         // 当前输入内容
    cursorPos   int            // 光标位置（基于 rune 索引）
    hasError    bool           // 是否存在错误提示
    placeholder string         // 可选的 placeholder 文本
    echo        func([]rune) string // 显示转换函数（明文 / 密文）
    validator   Validator      // 可选验证器
}
```

配合方法：

- `HandleKey(b byte) (done bool, err error)`：
  - 处理单个按键（包括 ESC 序列、控制键）。
  - 更新内部状态，返回是否结束输入（回车或取消）。
- `View() (inputLine string, errorLine string)`：
  - 根据当前状态生成“输入行”和“错误行”的文本表示。
  - 区分占位符 / 正常输入 / 密文显示。

再由一个统一的终端渲染器将 `View()` 的结果转为具体的 ANSI 输出与光标位置控制。这样可以：

- 把“逻辑状态机”与“具体终端操作”解耦；
- 减少 `\033[A` / `\033[K` 等控制码在多处手写的重复与出错机会；
- 使后续行为修改集中在状态机与渲染两个清晰层次里。

#### 4.3 统一 placeholder 与普通输入的编辑内核

在有了统一 `lineEditor` 的前提下：

- placeholder 与普通输入的区别仅体现在配置上：
  - `placeholder` 字段是否为空；
  - 初始渲染时是否绘制 placeholder。
- 密码与普通输入的区别仅体现在 `echo` 函数：
  - 明文：`echo = func(rs []rune) string { return string(rs) }`；
  - 密文：`echo = func(rs []rune) string { return strings.Repeat("*", len(rs)) }`。

这样可以让：

- `readInputWithPlaceholder` 和 `readLineCore` 共用同一套按键处理与错误管理逻辑；
- 差异只体现在 editor 的初始化配置上，减少代码重复与行为不一致的风险。

#### 4.4 抽离与组合 Validators

在 `input_validators.go` 中：

- 保留现有的 `ValidateXxx` 系列，并确保命名与行为语义清晰；
- 额外提供一个组合器：

```go
func CombineValidators(validators ...Validator) Validator {
    return func(value string) error {
        for _, v := range validators {
            if v == nil {
                continue
            }
            if err := v(value); err != nil {
                return err
            }
        }
        return nil
    }
}
```

这样可以在调用侧更容易写出既简洁又易读的校验逻辑，例如：

```go
prompt.Input().
    Prompt("请输入用户名").
    Validate(CombineValidators(
        ValidateRequired(),
        ValidateMinLength(3),
        ValidateMaxLength(20),
    )).
    Run()
```

#### 4.5 Terminal 抽象与可测试性提升（中长期）

可以考虑引入一个简化后的终端接口，例如：

```go
type Terminal interface {
    ReadByte() (byte, error)
    Print(s string)
    Println(s string)
    SetRawMode() (restore func(), err error)
}
```

并提供：

- `defaultTerminal`：使用 `os.Stdin` / `fmt.Print` / `term.MakeRaw` 实现，供生产环境使用；
- `mockTerminal`：基于缓冲区的实现，方便单元测试：
  - 预置输入字节序列；
  - 捕获输出字符串；
  - 验证特定按键序列下的行为是否符合预期。

预计收益：

- 可以对编辑内核和状态机编写细粒度单元测试；
- 在 CI 或非 TTY 环境下更容易控制行为（例如自动回退到非交互实现）。

#### 4.6 清理遗留实现与降级路径

- 对 `readPasswordInput` / `readPasswordFallback` / `readInputFallback` 做一次调用树检查：
  - 若主路径已经完全统一到 `readLineCore`：
    - 可考虑移除旧实现，减少维护成本与行为分歧；
  - 若仍需保留降级逻辑：
    - 建议放入独立的 `fallback` 文件；
    - 显式标注“仅在无法进入 raw 模式时使用”的场景，并确保行为简单清晰。

---

### 5. 建议的分步落地方案

为了降低风险，可以采用渐进式重构策略：

1. **第一步：拆文件 + 抽出 validators**
   - 不改变任何对外 API 与行为；
   - 将 validators 与公共 API、编辑内核分离到不同文件；
   - 方便后续修改时集中关注。

2. **第二步：合并 placeholder 与普通输入的编辑逻辑**
   - 在保留现有对外行为的前提下，尽量复用 `readLineCore` 中的逻辑；
   - 优先去除显而易见的重复代码（如 ESC 处理、错误提示处理的重复片段）。

3. **第三步：引入 `lineEditor` 状态机并逐步迁移**
   - 先在内部实现状态机和渲染逻辑；
   - 通过小步重构，将原有逻辑迁移到基于状态机的实现上；
   - 每一步都保持对外行为不变，配合 demo 命令与手工测试验证。

4. **第四步（可选）：引入 `Terminal` 抽象并补上单元测试**
   - 把 `os.Stdin` / `fmt.Print` / `term.MakeRaw` 收敛到一个实现里；
   - 为典型场景（简单输入、带 placeholder、错误重试、密码输入）编写针对编辑内核的测试。

---

### 6. 小结

- 当前 `input` 模块已经具备较强的功能与良好的用户体验，但实现上存在：
  - 终端控制逻辑重复、状态分散；
  - 校验器与 IO 强耦合；
  - 可测试性较弱、存在遗留实现等问题。
- 通过**拆文件、抽象编辑状态机、统一 placeholder 与普通输入的内核、抽离 validators 与终端接口**等步骤，可以在不影响现有调用方式的前提下，逐步显著降低复杂度，提高维护性与扩展性。

---

**最后更新**: 2026-01-06



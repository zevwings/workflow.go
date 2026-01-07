# 代码审查指南

> 专为 AI 助手设计的代码检查指南，提供系统化的代码质量和复用性检查方法，识别重复代码和优化机会。

## 🎯 核心原则

**检查重点**：识别代码重复、工具复用和优化机会。

**关键目标**：
- ✅ 消除重复：识别并抽取重复代码为公共方法
- ✅ 工具复用：使用已封装的工具函数替换重复实现
- ✅ 依赖优化：识别可用第三方库替换的自实现代码
- ✅ 规范遵循：确保代码风格和导入顺序符合项目规范
- ✅ **避免过度设计**：保持代码简洁，避免不必要的抽象和复杂性

---

## 📋 目录

- [检查目标](#-检查目标)
- [检查流程](#-检查流程)
- [重复代码检查](#-重复代码检查)
- [已封装工具检查](#-已封装工具检查)
- [第三方工具检查](#-第三方工具检查)
- [检查清单](#-检查清单)
- [示例分析](#-示例分析)

---

## 🎯 检查目标

代码检查的主要目标：

1. **消除重复代码**：识别并抽取重复的代码模式为公共方法
2. **复用已有工具**：识别可以使用已封装工具函数替换的代码
3. **优化依赖使用**：识别可以使用第三方库替换的自实现代码
4. **规范代码风格**：确保代码遵循项目规范，包括导入顺序等
5. **避免过度设计**：保持代码简洁实用，避免不必要的抽象层和复杂性

---

## 🔄 检查流程

### 步骤 1：扫描代码库

1. **全量扫描**：使用 `grep` 或语义搜索工具扫描 `internal/` 目录
2. **分类统计**：按模块、功能分类统计代码模式
3. **模式识别**：识别常见的代码模式和重复片段

### 步骤 2：重复代码分析

1. **相似度检测**：查找结构相似、逻辑重复的代码块
2. **参数化分析**：分析重复代码的差异点，确定可参数化的部分
3. **抽取建议**：提出公共方法抽取方案

### 步骤 3：工具函数检查

1. **工具函数清单**：列出所有已封装的工具函数
2. **使用情况分析**：检查是否有代码可以直接使用这些工具
3. **替换建议**：提出使用工具函数替换的方案

### 步骤 4：第三方库检查

1. **功能识别**：识别自实现的功能
2. **库匹配**：查找是否有成熟的第三方库可以替换
3. **替换评估**：评估替换的可行性和收益

### 步骤 5：过度设计检查

1. **抽象层评估**：检查是否有不必要的抽象层或中间层
2. **复杂度评估**：评估代码复杂度是否与实际需求匹配
3. **实用性评估**：确保设计满足当前需求，避免为未来可能的需求过度设计

### 步骤 6：生成报告

1. **问题汇总**：汇总所有发现的问题
2. **优先级排序**：按影响范围和收益排序
3. **改进建议**：提供具体的改进方案

---

## 🔍 重复代码检查

### 检查方法

#### 1. 文件操作模式

**检查目标**：查找重复的文件读写操作

**搜索模式**：
```bash
# 查找文件读取操作
grep -r "os.ReadFile" internal/
grep -r "os.Open" internal/
grep -r "ioutil.ReadFile" internal/  # 已废弃，但可能仍在使用

# 查找文件写入操作
grep -r "os.WriteFile" internal/
grep -r "os.Mkdir" internal/
grep -r "os.MkdirAll" internal/
```

**常见重复模式**：
- 文件读取 + 错误处理
- 文件写入 + 目录创建
- 文件打开 + 读取

**已封装工具**：
- `internal/lib/util/file.go` - 文件操作工具函数（如适用）
- 建议封装：`ReadFileWithContext()`, `WriteFileWithContext()`

**检查清单**：
- [ ] 是否有多个地方使用相同的文件读取模式？
- [ ] 错误处理是否一致？
- [ ] 是否可以抽取为公共函数？

#### 2. Git 命令执行模式

**检查目标**：查找重复的 Git 命令执行代码

**搜索模式**：
```bash
# 查找 Git 命令执行
grep -r "exec.Command.*git" internal/
grep -r "os/exec.*git" internal/
```

**常见重复模式**：
- Git 命令执行 + 输出读取
- Git 命令执行 + 错误处理
- Git 命令执行 + 静默模式

**已封装工具**：
- `internal/lib/git/helpers.go`（如适用）：
  - `ExecuteGitCommand()` - 执行 Git 命令并读取输出
  - `RunGitCommand()` - 执行 Git 命令（不读取输出）
  - `CheckGitSuccess()` - 静默执行并检查成功
  - `CheckRefExists()` - 检查 Git 引用是否存在

**检查清单**：
- [ ] 是否直接使用 `exec.Command("git", ...)` 而不是封装的函数？
- [ ] 是否有重复的错误处理逻辑？
- [ ] 是否可以复用 Git 辅助函数？

#### 3. 错误处理模式

**检查目标**：查找重复的错误处理代码

**搜索模式**：
```bash
# 查找错误处理模式
grep -r "fmt.Errorf" internal/
grep -r "errors.Wrap" internal/
grep -r "errors.New" internal/
grep -r "errors.Is" internal/
grep -r "errors.As" internal/
```

**常见重复模式**：
- 相同的错误消息格式
- 相同的上下文添加方式
- 相同的错误转换逻辑

**已封装工具**：
- `fmt.Errorf` - 格式化错误并包装
- `errors.Wrap` - 包装错误（如使用 pkg/errors）
- `errors.Is` - 检查错误类型
- `errors.As` - 错误类型断言

**检查清单**：
- [ ] 是否有重复的错误消息？
- [ ] 错误上下文是否一致？
- [ ] 是否可以统一错误处理函数？

#### 4. 字符串处理模式

**检查目标**：查找重复的字符串处理代码

**搜索模式**：
```bash
# 查找字符串处理
grep -r "strings\.Trim" internal/
grep -r "strings\.ToLower" internal/
grep -r "fmt\.Sprintf\|fmt\.Sprintf" internal/
```

**常见重复模式**：
- 字符串清理和规范化
- 敏感值隐藏
- 字符串格式化

**已封装工具**：
- `internal/lib/util/string.go`（如适用）：
  - `MaskSensitiveValue()` - 隐藏敏感值

**检查清单**：
- [ ] 是否有重复的字符串处理逻辑？
- [ ] 是否可以使用 `MaskSensitiveValue()` 隐藏敏感信息？
- [ ] 是否可以抽取为公共函数？

#### 5. 路径处理模式

**检查目标**：查找重复的路径处理代码

**搜索模式**：
```bash
# 查找路径处理
grep -r "filepath\.Join" internal/
grep -r "filepath\.Dir" internal/
grep -r "os\.Stat\|os\.IsNotExist" internal/
```

**常见重复模式**：
- 路径拼接和规范化
- 父目录获取和创建
- 路径存在性检查

**检查清单**：
- [ ] 是否有重复的路径处理逻辑？
- [ ] 是否可以抽取为公共函数？

#### 6. HTTP 请求模式

**检查目标**：查找重复的 HTTP 请求代码

**搜索模式**：
```bash
# 查找 HTTP 请求
grep -r "http\.NewRequest\|http\.Get\|http\.Post" internal/
grep -r "http\.Client" internal/
```

**常见重复模式**：
- HTTP 客户端创建
- 请求构建和发送
- 响应处理和错误处理

**已封装工具**：
- `internal/http/client.go` - `HttpClient` - 统一的 HTTP 客户端（内置 go-resty 重试机制）

**检查清单**：
- [ ] 是否直接使用 `net/http` 而不是 `HttpClient`？
- [ ] 是否有重复的请求构建逻辑？
- [ ] HTTP 客户端已内置智能重试机制（自动处理 5xx 错误和网络错误）

#### 7. 配置读取模式

**检查目标**：查找重复的配置读取代码

**搜索模式**：
```bash
# 查找配置读取
grep -r "toml\.Unmarshal\|toml\.Marshal" internal/
grep -r "os\.ReadFile.*toml" internal/
```

**常见重复模式**：
- 配置文件读取和解析
- 配置验证
- 默认值处理

**已封装工具**：
- `internal/lib/config/manager.go` - `ConfigManager` - 统一配置管理

**检查清单**：
- [ ] 是否直接读取配置文件而不是使用 `ConfigManager`？
- [ ] 是否有重复的配置解析逻辑？
- [ ] 是否可以使用 `ConfigManager` API？

#### 8. 日志输出模式

**检查目标**：查找重复的日志输出代码

**搜索模式**：
```bash
# 查找日志输出
grep -r "fmt\.Print\|fmt\.Println\|fmt\.Printf" internal/
grep -r "log\.Print\|log\.Println\|log\.Printf" internal/
```

**常见重复模式**：
- 直接使用 `fmt.Print*` 而不是日志工具
- 重复的日志格式
- 不一致的日志级别

**已封装工具**：
- `internal/logging/logger.go`：
  - `Info()` - 信息日志
  - `Error()` - 错误日志
  - `Warn()` - 警告日志
  - `Debug()` - 调试日志

**检查清单**：
- [ ] 是否使用 `fmt.Print*` 而不是日志工具？
- [ ] 日志格式是否一致？
- [ ] 是否可以使用统一的日志函数？

#### 9. 用户交互模式

**检查目标**：查找重复的用户交互代码

**搜索模式**：
```bash
# 查找用户交互
grep -r "promptui\|survey" internal/
grep -r "Input\|Select\|Confirm" internal/
```

**常见重复模式**：
- 重复的对话框配置
- 重复的用户输入验证
- 重复的选项处理

**已封装工具**：
- `internal/lib/prompt/`：
  - `Input` - 输入对话框
  - `Select` - 选择对话框
  - `MultiSelect` - 多选对话框
  - `Confirm` - 确认对话框

**检查清单**：
- [ ] 是否直接使用 `promptui` 或 `survey` 而不是封装的 Prompt？
- [ ] 是否有重复的对话框配置？
- [ ] 是否可以使用封装的 Prompt 类型？

#### 10. 进度指示器模式

**检查目标**：查找重复的进度指示器代码

**搜索模式**：
```bash
# 查找进度指示器
grep -r "spinner\|progress" internal/
```

**常见重复模式**：
- 重复的 Spinner 配置
- 重复的 Progress 配置
- 重复的样式设置

**已封装工具**：
- `internal/lib/output/` 或 `internal/lib/progress/`（如适用）：
  - `Spinner` - 旋转指示器
  - `Progress` - 进度条

**检查清单**：
- [ ] 是否直接使用第三方库而不是封装的 Spinner/Progress？
- [ ] 是否有重复的样式配置？
- [ ] 是否可以使用封装的指示器？

#### 11. 导入顺序模式

**检查目标**：检查导入语句是否遵循从顶部导入的规范

**导入顺序规范**：
1. **标准库导入**：Go 标准库（如 `fmt`、`os`、`net/http`）
2. **第三方库导入**：外部包（如 `github.com/spf13/cobra`）
3. **项目内部导入**：项目内部包（如 `github.com/zevwings/workflow/internal/lib/config`）
4. **每个分组之间用空行分隔**

**平台特定导入例外**：
- **一般规则**：所有导入语句都应该在文件顶部
- **例外情况**：如果导入只在特定平台使用，可以使用构建标签 `//go:build` 标记平台特定的文件
- 平台特定的代码应该在单独的文件中，使用构建标签限制其生效的平台

**搜索模式**：
```bash
# 查找文件中的导入语句
grep -r "^import" internal/ -A 5

# 查找可能违反导入顺序的文件（导入语句不在文件顶部）
# 使用 goimports 自动检查和修复
goimports -l .
```

**常见问题模式**：
- 导入语句分散在代码中间
- 导入顺序不符合规范（标准库、第三方库、项目内部）
- 导入分组之间没有空行分隔

**检查方法**：
1. **手动检查**：查看每个文件的导入部分
2. **自动化检查**：使用 `goimports` 或 `golangci-lint` 检查导入顺序
3. **代码审查**：在代码审查时重点关注导入顺序

**正确的导入顺序示例**：
```go
package config

import (
    // 标准库导入
    "fmt"
    "os"
    "path/filepath"

    // 第三方库导入
    "github.com/spf13/cobra"
    "github.com/spf13/viper"

    // 项目内部导入
    "github.com/zevwings/workflow/internal/lib/util/file"
    "github.com/zevwings/workflow/internal/lib/config"
)
```

**错误的导入顺序示例**：
```go
// ❌ 不好的做法：导入顺序混乱
package config

import (
    "github.com/zevwings/workflow/internal/lib/util/file"  // 项目内部导入应该在最后
    "fmt"  // 标准库导入应该在前面
    "github.com/spf13/cobra"  // 第三方库导入
)
```

**平台特定导入的正确示例**：
```go
// ✅ 好的做法：平台特定的代码在单独文件中，使用构建标签
//go:build darwin

package config

import (
    // 标准库导入
    "os"
    "path/filepath"

    // 第三方库导入（跨平台）
    "github.com/spf13/cobra"

    // 项目内部导入
    "github.com/zevwings/workflow/internal/lib/util/file"
)
```

**检查清单**：
- [ ] 所有导入语句是否都在文件顶部？
- [ ] 导入顺序是否符合规范（标准库 → 第三方库 → 项目内部）？
- [ ] 导入分组之间是否有空行分隔？
- [ ] 是否使用了 `goimports` 自动管理导入？
- [ ] 平台特定的代码是否使用了构建标签（`//go:build`）？
- [ ] 是否使用了 `goimports` 格式化导入顺序？

#### 12. 过度设计模式

**检查目标**：识别不必要的抽象和过度设计

**检查原则**：
1. **YAGNI（You Aren't Gonna Need It）**：不要为未来可能的需求添加功能
2. **KISS（Keep It Simple, Stupid）**：保持简单，避免不必要的复杂性
3. **实用主义**：优先满足当前需求，避免过度抽象

**搜索模式**：
```bash
# 查找可能过度设计的模式
# - 过多的接口定义
grep -r "type.*interface" internal/ | wc -l

# - 过多的抽象层
grep -r "type.*struct" internal/ | wc -l

# - 不必要的泛型
grep -r "\[.*\]" internal/ | grep -v "test"

# - 过多的中间层或包装器
grep -r "Wrapper\|Adapter\|Facade" internal/
```

**常见过度设计模式**：

1. **不必要的抽象层**：
   - 为单一实现创建接口
   - 创建多层抽象但只有一层实现
   - 过度使用设计模式（如为简单场景使用策略模式）

2. **过早优化**：
   - 为未来可能的需求添加功能
   - 创建过于灵活的配置系统但实际只需要简单配置
   - 添加不必要的泛型参数

3. **过度封装**：
   - 将简单函数包装成复杂的结构体
   - 创建不必要的中间层
   - 过度使用 Builder 模式

4. **不必要的复杂性**：
   - 使用复杂的数据结构处理简单数据
   - 创建复杂的错误处理系统处理简单错误
   - 过度使用宏或代码生成

**检查清单**：
- [ ] 是否有为单一实现创建的接口？是否可以简化为直接实现？
- [ ] 是否有不必要的抽象层？是否可以减少中间层？
- [ ] 是否有为未来需求添加的功能？是否可以移除？
- [ ] 是否有过度封装的简单函数？是否可以简化为直接函数调用？
- [ ] 代码复杂度是否与实际需求匹配？
- [ ] 是否使用了不必要的设计模式？
- [ ] 是否有不必要的泛型参数？

**示例：过度设计 vs 简洁设计**

**过度设计示例**：
```go
// ❌ 过度设计：为单一实现创建接口
type FileReader interface {
    Read(path string) (string, error)
}

type SimpleFileReader struct{}

func (r *SimpleFileReader) Read(path string) (string, error) {
    data, err := os.ReadFile(path)
    if err != nil {
        return "", fmt.Errorf("failed to read file: %w", err)
    }
    return string(data), nil
}

// 使用
reader := &SimpleFileReader{}
content, err := reader.Read(path)
```

**简洁设计示例**：
```go
// ✅ 简洁设计：直接使用函数
func ReadFile(path string) (string, error) {
    data, err := os.ReadFile(path)
    if err != nil {
        return "", fmt.Errorf("failed to read file: %w", err)
    }
    return string(data), nil
}

// 使用
content, err := ReadFile(path)
```

**过度设计示例**：
```go
// ❌ 过度设计：为简单场景使用复杂的 Builder 模式
type ConfigBuilder struct {
    host *string
    port *int
}

func NewConfigBuilder() *ConfigBuilder {
    return &ConfigBuilder{}
}

func (b *ConfigBuilder) Host(host string) *ConfigBuilder {
    b.host = &host
    return b
}

func (b *ConfigBuilder) Port(port int) *ConfigBuilder {
    b.port = &port
    return b
}

func (b *ConfigBuilder) Build() (*Config, error) {
    host := "localhost"
    if b.host != nil {
        host = *b.host
    }
    port := 8080
    if b.port != nil {
        port = *b.port
    }
    return &Config{Host: host, Port: port}, nil
}
```

**简洁设计示例**：
```go
// ✅ 简洁设计：直接使用结构体初始化
type Config struct {
    Host string
    Port int
}

func NewConfig(host string, port int) *Config {
    return &Config{Host: host, Port: port}
}

func DefaultConfig() *Config {
    return &Config{
        Host: "localhost",
        Port: 8080,
    }
}
```

---

## 🛠️ 已封装工具检查

### 工具函数清单

#### 文件操作工具

**位置**：`internal/lib/util/file.go`（如适用）

- `ReadFile(path)` - 读取文件内容

**检查方法**：
```bash
# 查找直接使用 os.Open 或 os.ReadFile 的地方
grep -r "os\.Open\|os\.ReadFile" internal/ | grep -v "internal/lib/util/file.go"
```

#### 字符串处理工具

**位置**：`internal/lib/util/string.go`（如适用）

- `MaskSensitiveValue(value)` - 隐藏敏感值（如密码、token）

**检查方法**：
```bash
# 查找可能包含敏感信息的字符串处理
grep -r "password\|token\|secret\|key" internal/ -i
```

#### 格式化工具

**位置**：`internal/lib/util/format.go`（如适用）

- `FormatSize(bytes)` - 格式化文件大小（B, KB, MB, GB）

**检查方法**：
```bash
# 查找文件大小格式化
grep -r "1024\|bytes\|KB\|MB\|GB" internal/ -i
```

#### 平台检测工具

**位置**：`internal/lib/util/platform.go`（如适用）

- `DetectReleasePlatform()` - 检测发布平台（操作系统和架构）

**检查方法**：
```bash
# 查找平台检测代码
grep -r "runtime\.GOOS\|runtime\.GOARCH\|build\." internal/
```

#### 浏览器和剪贴板工具

**位置**：`internal/lib/util/browser.go`, `internal/lib/util/clipboard.go`（如适用）

- `OpenBrowser(url)` - 打开浏览器
- `CopyToClipboard(text)` - 复制到剪贴板

**检查方法**：
```bash
# 查找浏览器和剪贴板操作
grep -r "browser\|clipboard" internal/ -i
```

#### 解压和校验和工具

**位置**：`internal/lib/util/unzip.go`, `internal/lib/util/checksum.go`（如适用）

- `ExtractTarGz()` - 解压 tar.gz 文件
- `VerifyChecksum()` - 验证 SHA256 校验和

**检查方法**：
```bash
# 查找解压和校验和操作
grep -r "tar\.gz\|sha256\|crypto/sha256" internal/ -i
```

#### 日期格式化工具

**位置**：`internal/lib/util/date.go`（如适用）

- `FormatDocumentTimestamp()` - 格式化文档时间戳
- `FormatLastUpdated()` - 格式化最后更新时间
- `FormatLastUpdatedWithTime()` - 格式化最后更新时间（带时间）

**检查方法**：
```bash
# 查找日期格式化
grep -r "time\.Now\|time\.Format\|time\.Parse" internal/
```

#### 表格输出工具

**位置**：`internal/lib/util/table.go`

- `TableBuilder` - 表格构建器
- `TableStyle` - 表格样式

**检查方法**：
```bash
# 查找表格输出
grep -r "tabled\|Table" src/
```

#### Git 操作工具

**位置**：`internal/lib/git/helpers.go`（如适用）

- `ExecuteGitCommand(args)` - 执行 Git 命令并读取输出
- `RunGitCommand(args)` - 执行 Git 命令（不读取输出）
- `CheckGitSuccess(args)` - 静默执行并检查成功
- `CheckRefExists(refPath)` - 检查 Git 引用是否存在
- `RemoveBranchPrefix(branch)` - 移除分支名称前缀

**检查方法**：
```bash
# 查找直接使用 exec.Command("git", ...) 的地方
grep -r 'exec.Command.*"git"' internal/ | grep -v "internal/lib/git/helpers.go"
```

#### HTTP 客户端工具

**位置**：`internal/http/`

- `HttpClient` - 统一的 HTTP 客户端（内置 go-resty 重试机制）
- `HttpResponse`（如适用） - HTTP 响应处理

**检查方法**：
```bash
# 查找直接使用 net/http 的地方
grep -r "http\.Client\|http\.NewRequest" internal/ | grep -v "internal/http"
```

#### 日志工具

**位置**：`internal/logging/logger.go`

- `Info(...)` - 信息日志
- `Error(...)` - 错误日志
- `Warn(...)` - 警告日志
- `Debug(...)` - 调试日志

**检查方法**：
```bash
# 查找直接使用 fmt.Print* 的地方
grep -r "fmt\.Print\|fmt\.Println\|fmt\.Printf" internal/ | grep -v "internal/logging"
```

#### 对话框工具

**位置**：`internal/lib/prompt/`

- `Input` - 输入对话框
- `Select` - 选择对话框
- `MultiSelect` - 多选对话框
- `Confirm` - 确认对话框

**检查方法**：
```bash
# 查找直接使用 promptui 或 survey 的地方
grep -r "promptui\|survey" internal/ | grep -v "internal/lib/prompt"
```

#### 进度指示器工具

**位置**：`internal/lib/output/` 或 `internal/lib/progress/`（如适用）

- `Spinner` - 旋转指示器
- `Progress` - 进度条

**检查方法**：
```bash
# 查找直接使用第三方 spinner/progress 库的地方
grep -r "spinner\|progress" internal/ | grep -v "internal/lib/output\|internal/lib/progress"
```

#### 配置管理工具

**位置**：`internal/lib/config/`

- `ConfigManager` - 统一配置管理
- `Paths`（如适用） - 路径管理
- `LLMSettings`（如适用） - LLM 配置

**检查方法**：
```bash
# 查找直接读取配置文件的地方
grep -r "os\.ReadFile.*toml\|toml\.Unmarshal" internal/ | grep -v "internal/lib/config"
```

---

## 📦 第三方工具检查

### 检查方法

#### 1. 正则表达式处理

**检查目标**：查找自实现的正则表达式处理

**搜索模式**：
```bash
# 查找正则表达式使用
grep -r "regexp\.Compile\|regexp\.MustCompile" internal/
```

**第三方工具**：
- `regexp` - Go 标准正则表达式库（已使用）
- 检查是否有重复的正则表达式编译

**检查清单**：
- [ ] 是否重复编译相同的正则表达式？
- [ ] 是否可以缓存编译后的正则表达式（使用 `regexp.MustCompile` 或包级别变量）？

#### 2. JSON 处理

**检查目标**：查找 JSON 序列化/反序列化代码

**搜索模式**：
```bash
# 查找 JSON 处理
grep -r "json\.Marshal\|json\.Unmarshal" internal/
```

**第三方工具**：
- `encoding/json` - Go 标准 JSON 库（已使用）

**检查清单**：
- [ ] 是否使用统一的序列化方式？
- [ ] 是否有重复的 JSON 解析逻辑？

#### 3. TOML 处理

**检查目标**：查找 TOML 配置文件处理

**搜索模式**：
```bash
# 查找 TOML 处理
grep -r "toml\.Unmarshal\|toml\.Marshal" internal/
```

**第三方工具**：
- `github.com/pelletier/go-toml` 或类似库 - TOML 解析（已使用）

**检查清单**：
- [ ] 是否使用统一的 TOML 解析方式？
- [ ] 是否需要 TOML 编辑功能？

#### 4. 命令行参数解析

**检查目标**：查找命令行参数处理

**搜索模式**：
```bash
# 查找命令行参数
grep -r "cobra\|spf13/cobra" internal/
```

**第三方工具**：
- `github.com/spf13/cobra` - 命令行参数解析（已使用）

**检查清单**：
- [ ] 是否使用 `cobra` 进行参数解析？
- [ ] 是否有重复的参数定义？

#### 5. 并发处理

**检查目标**：查找并发代码

**搜索模式**：
```bash
# 查找并发代码
grep -r "go func\|goroutine\|channel\|sync\." internal/
```

**第三方工具**：
- Go 标准库 `goroutine` 和 `channel` - 并发处理（已使用）
- `sync` 包 - 同步原语（已使用）

**检查清单**：
- [ ] 是否需要并发处理？
- [ ] 是否可以使用 goroutine 提高性能？

#### 6. 并发处理

**检查目标**：查找并发处理代码

**搜索模式**：
```bash
# 查找并发处理
grep -r "go func\|goroutine\|sync\.WaitGroup\|sync\.Mutex" internal/
```

**已封装工具**：
- `internal/lib/concurrent/executor.go`（如适用） - `ConcurrentExecutor` - 并发执行器

**第三方工具**：
- Go 标准库 `sync` 包 - 同步原语（已使用）
- `golang.org/x/sync` - 高级并发工具（如果需要）

**检查清单**：
- [ ] 是否可以使用 `ConcurrentExecutor`？
- [ ] 是否需要更高级的并发工具？

#### 7. 时间处理

**检查目标**：查找时间处理代码

**搜索模式**：
```bash
# 查找时间处理
grep -r "time\." internal/ | grep -v "_test.go"
grep -r "time\.(Now\|Parse\|Format\|Unix\|Duration)" internal/
```

**第三方工具**：
- `time` - Go 标准库，满足大部分时间处理需求
  - 支持时间创建、格式化、解析、计算等
  - 时区支持：`time.LoadLocation()`
  - 示例：`time.Now()`, `time.Parse()`, `time.Format()`
- 第三方库（仅在标准库不足时考虑）：
  - `github.com/jinzhu/now` - 更友好的时间解析和操作
  - `github.com/araddon/dateparse` - 灵活的时间字符串解析

**检查清单**：
- [ ] 是否需要更复杂的时间操作？
- [ ] 是否可以复用已有的时间处理工具函数？
- [ ] 时间格式是否统一？是否使用了项目标准格式？
- [ ] 时区处理是否正确？
- [ ] 时间比较和计算是否有边界情况需要处理？
- [ ] 是否可以使用标准库 `time` 满足需求？

#### 8. 路径处理

**检查目标**：查找路径处理代码

**搜索模式**：
```bash
# 查找路径处理
grep -r "filepath\.Join\|filepath\.Dir\|filepath\.Walk" internal/
```

**第三方工具**：
- `path/filepath` - Go 标准库路径处理（已使用）
- `github.com/karrick/godirwalk` - 高性能目录遍历（如果需要）

**检查清单**：
- [ ] 是否需要路径差异计算？
- [ ] 是否需要目录遍历？

#### 9. 环境变量处理

**检查目标**：查找环境变量处理

**搜索模式**：
```bash
# 查找环境变量
grep -r "os\.Getenv\|os\.Setenv\|os\.LookupEnv" internal/
```

**第三方工具**：
- `os` 包 - Go 标准库环境变量支持（已使用）
- `github.com/joho/godotenv` - .env 文件支持（如果需要）
- `github.com/caarlos0/env` - 环境变量到结构体（如果需要）

**检查清单**：
- [ ] 是否需要 .env 文件支持？
- [ ] 是否需要环境变量到结构体的转换？

#### 10. 错误处理

**检查目标**：查找错误处理代码

**搜索模式**：
```bash
# 查找错误处理
grep -r "error\|fmt\.Errorf\|errors\." internal/
```

**第三方工具**：
- `errors` 包 - Go 标准库错误处理（已使用）
- `fmt.Errorf` - 错误格式化（已使用）
- `github.com/pkg/errors` - 增强错误处理（如果需要）

**检查清单**：
- [ ] 是否统一使用 `error` 接口？
- [ ] 是否使用 `fmt.Errorf` 和 `%w` 动词添加上下文？

---

## ✅ 检查清单

### 重复代码检查清单

- [ ] **文件操作**：检查是否有重复的文件读写模式
- [ ] **Git 命令**：检查是否直接使用 `cmd("git", ...)` 而不是 `cmd-_read()` 或 `cmd-_run()`
- [ ] **错误处理**：检查是否有重复的错误处理逻辑
- [ ] **字符串处理**：检查是否有重复的字符串处理逻辑
- [ ] **路径处理**：检查是否有重复的路径处理逻辑
- [ ] **HTTP 请求**：检查是否直接使用 `net/http` 而不是 `HttpClient`
- [ ] **配置读取**：检查是否直接读取配置文件而不是使用 `ConfigManager`
- [ ] **日志输出**：检查是否使用 `fmt.Print*` 而不是日志工具
- [ ] **用户交互**：检查是否直接使用 `promptui` 或 `survey` 而不是封装的 Prompt
- [ ] **进度指示器**：检查是否直接使用第三方库而不是封装的 Spinner/Progress
- [ ] **导入顺序**：检查导入语句是否遵循从顶部导入的规范（标准库 → 第三方库 → 项目内部）
- [ ] **过度设计**：检查是否有不必要的抽象层、过度封装或过早优化

### 已封装工具检查清单

- [ ] **文件操作**：是否可以使用封装的文件读取函数（如适用）？
- [ ] **字符串处理**：是否可以使用 `MaskSensitiveValue()`？
- [ ] **格式化**：是否可以使用 `FormatSize()`？
- [ ] **平台检测**：是否可以使用 `DetectReleasePlatform()`？
- [ ] **浏览器/剪贴板**：是否可以使用 `OpenBrowser()` 或 `CopyToClipboard()`（如适用）？
- [ ] **解压/校验和**：是否可以使用 `ExtractTarGz()` 或 `VerifyChecksum()`（如适用）？
- [ ] **日期格式化**：是否可以使用日期格式化函数？
- [ ] **表格输出**：是否可以使用 `TableBuilder`？
- [ ] **Git 操作**：是否可以使用 `git/helpers.go` 中的函数？
- [ ] **HTTP 客户端**：是否可以使用 `HttpClient`（已内置智能重试机制）？
- [ ] **日志**：是否可以使用日志函数？
- [ ] **对话框**：是否可以使用封装的 Prompt 类型？
- [ ] **进度指示器**：是否可以使用封装的 Spinner/Progress？
- [ ] **配置管理**：是否可以使用 `ConfigManager` API？

### 第三方工具检查清单

- [ ] **正则表达式**：是否重复编译正则表达式？是否可以缓存？
- [ ] **JSON 处理**：是否使用统一的序列化方式？
- [ ] **TOML 处理**：是否使用统一的 TOML 解析方式？
- [ ] **命令行参数**：是否使用 `cobra` 进行参数解析？
- [ ] **并发处理**：是否需要并发处理？是否可以使用 goroutine 提高性能？
- [ ] **并发处理**：是否可以使用 `ConcurrentExecutor`？
- [ ] **时间处理**：是否需要更复杂的时间操作？是否可以使用标准库 `time` 满足需求？
- [ ] **路径处理**：是否需要路径差异计算或目录遍历？
- [ ] **环境变量**：是否需要 .env 文件支持或环境变量到结构体的转换？
- [ ] **错误处理**：是否统一使用 `error` 接口？是否使用 `fmt.Errorf` 和 `%w` 动词添加上下文？

---

## 📝 示例分析

### 示例 1：文件读取重复代码

**问题代码**：
```go
// 文件 A
data, err := os.ReadFile(configPath)
if err != nil {
    return "", fmt.Errorf("failed to read config file: %s: %w", configPath, err)
}

// 文件 B
data, err := os.ReadFile(inputPath)
if err != nil {
    return "", fmt.Errorf("failed to read file: %s: %w", inputPath, err)
}
```

**改进方案**：
```go
// 抽取为公共函数
func ReadFileWithContext(path, context string) ([]byte, error) {
    data, err := os.ReadFile(path)
    if err != nil {
        return nil, fmt.Errorf("%s: %s: %w", context, path, err)
    }
    return data, nil
}

// 使用
data, err := ReadFileWithContext(configPath, "failed to read config file")
```

### 示例 2：Git 命令执行重复代码

**问题代码**：
```go
// 直接使用 exec.Command
cmd := exec.Command("git", "rev-parse", "--abbrev-ref", "HEAD")
output, err := cmd.Output()
if err != nil {
    return "", fmt.Errorf("failed to get current branch: %w", err)
}
```

**改进方案**：
```go
// 使用已封装的工具函数
import "github.com/zevwings/workflow/internal/lib/git"

output, err := git.ExecuteCommand("rev-parse", "--abbrev-ref", "HEAD")
if err != nil {
    return "", err
}
```

### 示例 3：日志输出重复代码

**问题代码**：
```go
// 直接使用 fmt.Printf
fmt.Printf("✓ Successfully created branch: %s\n", branchName)
fmt.Printf("✗ Failed to create branch: %v\n", err)
```

**改进方案**：
```go
// 使用封装的输出工具
import "github.com/zevwings/workflow/internal/output"

output := output.NewOutput(true)
output.Success("Successfully created branch: %s", branchName)
output.Error("Failed to create branch: %v", err)
```

### 示例 4：HTTP 请求重复代码

**问题代码**：
```go
// 直接使用 net/http
req, err := http.NewRequest("GET", url, nil)
if err != nil {
    return nil, err
}
req.Header.Set("Authorization", "Bearer "+token)
resp, err := http.DefaultClient.Do(req)
```

**改进方案**：
```go
// 使用封装的 HTTP 客户端
import "github.com/zevwings/workflow/internal/http"

client := http.NewClient()
resp, err := client.Get(url, http.WithAuthToken(token))
```

### 示例 5：用户交互重复代码

**问题代码**：
```go
// 直接使用 promptui
prompt := promptui.Prompt{
    Label: "Enter branch name",
}
input, err := prompt.Run()
```

**改进方案**：
```go
// 使用封装好的输入工具
import "github.com/zevwings/workflow/internal/lib/prompt"

input, err := prompt.Input("Enter branch name")
```

### 示例 6：导入顺序不规范

**问题代码**：
```go
// ❌ 不好的做法：导入顺序混乱
package config

import (
    "github.com/zevwings/workflow/internal/lib/util/file"  // 项目内部导入应该在最后
    "fmt"  // 标准库导入应该在前面
    "github.com/spf13/cobra"  // 第三方库导入
)
```

**改进方案**：
```go
// ✅ 好的做法：导入顺序规范，都在文件顶部
package config

import (
    // 标准库导入
    "fmt"
    "os"

    // 第三方库导入
    "github.com/spf13/cobra"

    // 项目内部导入
    "github.com/zevwings/workflow/internal/lib/util/file"
)
```

---

## 🔧 检查工具推荐

### 代码搜索工具

1. **ripgrep (rg)**：快速文本搜索
   ```bash
   rg "pattern" src/
   ```

2. **语义搜索**：使用 AI 工具进行语义搜索
   - 查找相似功能的代码
   - 识别重复模式

### 代码分析工具

1. **golangci-lint**：Go 代码检查
   ```bash
   golangci-lint run
   ```

2. **gofmt / goimports**：代码格式化检查
   ```bash
   gofmt -l .
   goimports -l .
   ```

### 重复代码检测

1. **人工审查**：逐文件检查
2. **模式匹配**：使用 grep 查找常见模式
3. **语义分析**：使用 AI 工具进行语义分析

---

## 📚 参考文档

### 项目文档

- [代码风格规范](../../development/code-style.md) - 代码风格规范
- [错误处理规范](../../development/error-handling.md) - 错误处理规范
- [命名规范](../../development/naming.md) - 命名规范
- [架构文档](../architecture/) - 各模块架构文档
- [工具函数模块架构文档](../../../architecture/) - 各工具函数模块架构（fs.md、system.md、zip.md、checksum.md、format.md 等）

### 工具函数位置

- `internal/lib/util/` - 通用工具函数
- `internal/lib/git/helpers.go` - Git 操作辅助函数
- `internal/http/` - HTTP 客户端工具
- `internal/lib/prompt/` - 用户交互工具
- `internal/lib/output/` - 输出格式化工具
- `internal/logging/` - 日志工具
- `internal/lib/config/` - 配置管理工具

---

## 📝 检查报告模板

检查完成后，应生成检查报告，包含：

1. **问题汇总**：列出所有发现的问题
2. **优先级排序**：按影响范围和收益排序
3. **改进建议**：提供具体的改进方案
4. **代码示例**：提供改进前后的代码对比

**报告格式**：
```markdown
# 代码检查报告

## 问题汇总

### 1. 重复代码问题
- [问题描述]
- [位置]
- [改进方案]

### 2. 工具函数替换
- [问题描述]
- [位置]
- [改进方案]

### 3. 第三方工具替换
- [问题描述]
- [位置]
- [改进方案]

## 优先级排序

1. [高优先级问题]
2. [中优先级问题]
3. [低优先级问题]

## 改进建议

[具体的改进建议和代码示例]
```

---

**最后更新**: 2025-01-27

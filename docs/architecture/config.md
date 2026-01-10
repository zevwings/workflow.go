# config 模块架构文档

## 📋 概述

config 模块是 Workflow CLI 的基础模块，提供配置管理功能。该模块专注于配置的读取、保存和管理，不涉及命令层的业务逻辑。

config 模块提供全局配置管理、仓库配置管理、多语言支持等功能，总代码行数约 1367+ 行。

**模块统计：**
- 代码行数：约 1367+ 行（不含测试文件）
- 主要文件：18 个核心文件
- 主要结构体：`GlobalManager`、`RepoManager`、`GlobalConfig`、`RepoConfig`、`LLMConfig`、`TemplateConfig` 等
- 支持功能：全局配置管理、仓库配置管理、多语言支持、LLM 配置管理

**注意**：本模块是基础库模块，其他模块通过导入使用。配置分为全局配置（用户级别）和仓库配置（项目级别）。

---

## 📁 模块架构（核心业务逻辑）

config 模块（`internal/config/`）是 Workflow CLI 的基础库模块，提供配置管理功能。该模块专注于配置的读取、保存和管理，不涉及命令层的业务逻辑。

### 模块结构

```
internal/config/
├── global_manager.go          # 全局配置管理器（363行）
├── repo_manager.go            # 仓库配置管理器（560行）
├── types.go                   # 全局配置结构定义（24行）
├── helpers.go                 # 配置辅助函数（40行）
├── paths.go                   # XDG 路径工具（72行）
│
├── 配置结构体文件
│   ├── user.go                # 用户配置结构（8行）
│   ├── github.go              # GitHub 配置结构（14行）
│   ├── jira.go                # Jira 配置结构（9行）
│   ├── log.go                 # 日志配置结构（7行）
│   ├── proxy.go               # 代理配置结构（8行）
│   ├── llm.go                 # LLM 配置结构和方法（95行）
│   ├── template.go            # 模板配置结构（14行）
│   ├── branch.go              # 分支配置结构（11行）
│   └── pull_requests.go       # PR 配置结构（9行）
│
└── languages.go               # 语言支持（215行）
```

**总计：约 1440+ 行代码**

### 依赖模块

- **`github.com/spf13/viper`**：配置文件读取和管理
  - 用于读取和解析 TOML 配置文件
- **`github.com/pelletier/go-toml/v2`**：TOML 格式解析和序列化
  - 用于配置文件的序列化和反序列化
- **`github.com/adrg/xdg`**：XDG Base Directory Specification 实现
  - 用于获取符合 XDG 规范的配置、数据、状态、缓存目录路径
  - 支持 Unix、Windows、macOS 等平台
- **`github.com/zevwings/workflow/internal/logging`**：日志记录
  - 用于记录配置操作的日志

### 模块集成

- **`internal/commands/`**：命令层使用配置管理器
  - `config.Global()` - 获取全局配置管理器
  - `repoManager.TemplateConfig` - 直接访问模板配置（推荐）
  - `repoManager.GetTemplateConfig()` - 获取模板配置（向后兼容）
- **`internal/infrastructure/llm/`**：LLM 基础设施层使用配置
  - `config.Global()` - 获取全局配置管理器
  - `manager.LLMConfig` - 访问 LLM 配置
- **`internal/infrastructure/config/`**：配置基础设施层
  - `config.GlobalRepoManager()` - 获取仓库配置管理器

---

## 🏗️ 架构设计

### 设计原则

1. **单例模式**：`GlobalManager` 和 `RepoManager` 采用单例模式，确保进程内配置一致性
2. **直接字段访问**：`GlobalManager` 和 `RepoManager` 都提供公开字段，支持直接访问配置，简化使用
3. **依赖注入**：`RepoManager` 通过接口实现依赖注入，解耦对 git 模块的依赖
4. **配置分离**：区分全局配置（用户级别）和仓库配置（项目级别）
5. **延迟加载**：私有配置采用延迟加载机制，提高性能
6. **统一配置结构**：`GlobalConfig` 和 `RepoConfig` 统一管理所有子配置模块

### 核心组件

#### 1. GlobalManager (global_manager.go)

**职责**：管理用户级别的全局配置（遵循 XDG 规范：`$XDG_CONFIG_HOME/workflow/config.toml`，默认 `~/.config/workflow/config.toml`）

**主要方法**：
- `Global()` - 获取全局配置管理器单例
- `Load()` - 从文件加载配置到内存
- `Save()` - 保存当前配置到文件
- `SaveDefault()` - 保存默认配置
- `GetLLMConfig()` - 获取 LLM 配置（向后兼容）
- `GetGitHubConfig()` - 获取 GitHub 配置（向后兼容）
- `GetUserConfig()` - 获取用户配置（向后兼容）
- `GetJiraConfig()` - 获取 Jira 配置（向后兼容）
- `GetLogConfig()` - 获取日志配置（向后兼容）
- `GetProxyConfig()` - 获取代理配置（向后兼容）

**关键特性**：
- 单例模式：使用 `sync.Once` 确保线程安全的单例初始化
- 直接字段访问：提供 `Config`、`LLMConfig`、`GitHubConfig` 等公开字段
- 便捷字段：提供指向 `Config` 子配置的便捷字段，简化访问
- 自动同步：`Save()` 后自动重新加载以同步 viper

**使用场景**：
- 读取和修改用户级别的全局配置
- 管理 LLM、GitHub、Jira 等服务的配置
- 管理日志和代理配置

#### 2. RepoManager (repo_manager.go)

**职责**：管理仓库级别的配置（项目公共配置和项目私有配置）

**主要方法**：
- `GlobalRepoManager(gitRepo GitRepository)` - 获取仓库配置管理器单例
- `Load()` - 加载仓库配置
- `Save()` - 保存当前配置到文件
- `GetTemplateConfig()` - 获取模板配置（向后兼容）
- `GetBranchPrefix()` - 获取分支前缀（个人偏好）
- `GetIgnoreBranches()` - 获取忽略的分支列表（个人偏好）
- `GetAutoAcceptChangeType()` - 获取自动接受变更类型设置（个人偏好）
- `SaveTemplateConfig(cfg *TemplateConfig)` - 保存模板配置（已废弃，请使用 `Save()`）
- `GetRepoID()` - 获取仓库 ID
- `GetPublicConfigPath()` - 获取公共配置文件路径
- `GetPrivateConfigPath()` - 获取私有配置文件路径

**关键特性**：
- 单例模式：使用 `sync.Once` 确保线程安全的单例初始化
- 直接字段访问：提供 `Config`、`TemplateConfig` 等公开字段
- 便捷字段：提供指向 `Config` 子配置的便捷字段，简化访问
- 依赖注入：通过 `GitRepository` 接口解耦对 git 模块的依赖
- 配置分离：区分项目公共配置（提交到 Git）和项目私有配置（不提交）
- 延迟加载：私有配置采用延迟加载机制，带缓存
- 自动同步：`Save()` 后自动重新加载以同步 publicViper
- 仓库 ID：基于 Git remote URL 生成唯一仓库标识符

**使用场景**：
- 读取和修改项目级别的模板配置
- 管理个人偏好配置（分支前缀、忽略分支等）
- 获取仓库相关信息（仓库 ID、配置路径等）

#### 3. LLMConfig (llm.go)

**职责**：管理 LLM 配置，提供 provider 和 language 的访问方法

**主要方法**：
- `CurrentProvider()` - 获取当前 provider 的配置（APIKey、Model、URL）
- `CurrentLanguage()` - 获取当前语言配置

**关键特性**：
- 多 provider 支持：支持 OpenAI、DeepSeek、Proxy 等多种 LLM 提供商
- 默认值处理：provider 未设置 model 时返回默认值
- 语言支持：与 `languages.go` 集成，提供多语言支持

**使用场景**：
- 获取当前配置的 LLM provider 信息
- 获取当前配置的语言信息

#### 4. 语言支持 (languages.go)

**职责**：提供多语言支持功能

**主要函数**：
- `FindLanguage(code string)` - 查找语言（支持大小写不敏感和部分匹配）
- `GetLanguageInstruction(code string)` - 获取语言指令模板
- `GetLanguageRequirement(systemPrompt, languageCode string)` - 获取语言要求
- `GetSupportedLanguageCodes()` - 获取支持的语言代码列表
- `GetSupportedLanguageDisplayNames()` - 获取支持的语言显示名称列表

**关键特性**：
- 多语言支持：支持英语、中文（简体/繁体）、日语、韩语、德语、法语等
- 智能匹配：支持大小写不敏感和部分匹配（如 "zh" 匹配 "zh-CN"）
- 指令模板：为每种语言提供 LLM 指令模板

**使用场景**：
- 根据语言代码查找语言信息
- 生成多语言的 LLM 指令
- 获取支持的语言列表

#### 5. 配置辅助函数 (helpers.go)

**职责**：提供通用的配置保存辅助函数

**主要函数**：
- `SaveConfigToFile(path string, config interface{})` - 保存配置到文件

**关键特性**：
- 自动创建目录：如果目录不存在，自动创建
- TOML 序列化：自动将配置序列化为 TOML 格式
- 错误处理：提供详细的错误信息

**使用场景**：
- 保存配置到文件（被 `GlobalManager` 和 `RepoManager` 使用）

#### 6. 路径工具 (paths.go)

**职责**：提供 XDG Base Directory Specification 路径工具函数

**主要函数**：
- `ConfigDir()` - 获取配置目录（`$XDG_CONFIG_HOME/workflow`）
- `DataDir()` - 获取数据目录（`$XDG_DATA_HOME/workflow`）
- `StateDir()` - 获取状态目录（`$XDG_STATE_HOME/workflow`）
- `CacheDir()` - 获取缓存目录（`$XDG_CACHE_HOME/workflow`）

**关键特性**：
- XDG 规范：遵循 XDG Base Directory Specification
- 跨平台支持：支持 Unix、Windows、macOS 等平台
- 第三方库：使用 `github.com/adrg/xdg` 实现，减少维护成本
- 统一接口：提供统一的路径获取接口，简化使用

**使用场景**：
- 获取配置目录（被 `GlobalManager` 和 `RepoManager` 使用）
- 获取状态目录（被日志系统使用）
- 获取数据目录和缓存目录（供其他模块使用）

#### 7. 配置结构体 (types.go)

**职责**：定义统一的配置结构体，包含所有子配置模块

**主要结构体**：

1. **`GlobalConfig`**：全局配置结构
   - 包含：`User`、`Jira`、`GitHub`、`Log`、`LLM`、`Proxy`
   - 用于：`~/.workflow/config.toml`（用户级别配置）

2. **`RepoConfig`**：仓库配置结构
   - 包含：`Template`（模板配置）
   - 用于：`.workflow/config.toml`（项目公共配置，提交到 Git）

**关键特性**：
- 统一管理：所有子配置模块统一在一个结构体中
- 类型安全：使用结构体而非 map，提供类型安全
- 易于扩展：添加新配置只需在结构体中添加字段

**使用场景**：
- `GlobalManager.Config` 字段使用 `GlobalConfig`
- `RepoManager.Config` 字段使用 `RepoConfig`
- 配置序列化和反序列化

### 设计模式

#### 1. 单例模式

**实现**：使用 `sync.Once` 确保线程安全的单例初始化

```go
var (
    globalManager *GlobalManager
    globalOnce    sync.Once
    globalErr     error
)

func Global() (*GlobalManager, error) {
    globalOnce.Do(func() {
        globalManager, globalErr = newGlobalManager()
    })
    return globalManager, globalErr
}
```

**优势**：
- 线程安全：可以在多线程环境中安全使用
- 资源优化：避免重复创建管理器实例
- 配置一致性：确保整个进程使用相同的配置状态

#### 2. 依赖注入

**实现**：通过接口实现依赖注入，解耦模块依赖

```go
type GitRepository interface {
    GetRepoPath() string
    IsGitRepo(path string) bool
    Open(path string) (GitRepo, error)
}

func GlobalRepoManager(gitRepo GitRepository) (*RepoManager, error) {
    // 使用接口而不是直接依赖 git 模块
}
```

**优势**：
- 解耦：config 模块不直接依赖 git 模块
- 可测试性：可以轻松创建 mock 实现进行测试
- 灵活性：可以替换不同的 Git 实现

#### 3. 直接字段访问

**实现**：提供公开字段和便捷字段，支持直接访问配置

```go
type GlobalManager struct {
    Config *GlobalConfig
    LLMConfig    *LLMConfig    // 指向 Config.LLM
    GitHubConfig *GitHubConfig // 指向 Config.GitHub
    // ...
}
```

**优势**：
- 简洁：`manager.LLMConfig.Provider` 比 `manager.GetLLMConfig().Provider` 更简洁
- 类型安全：直接访问字段，编译时检查
- 直观：代码更易读易写

### 错误处理

#### 分层错误处理

1. **文件操作层**：文件不存在、权限错误、IO 错误
   - 处理方式：返回详细的错误信息，支持自动创建默认配置
2. **配置解析层**：TOML 格式错误、类型转换错误
   - 处理方式：返回解析错误，记录日志
3. **业务逻辑层**：配置验证失败、必需字段缺失
   - 处理方式：返回业务错误，提供默认值

#### 容错机制

- **配置文件不存在**：自动创建默认配置文件
- **配置字段缺失**：使用默认值或返回空值
- **配置解析失败**：记录错误日志，返回错误信息
- **Git 仓库检测失败**：使用基于路径的简单 ID

---

## 🔄 集成关系

### 模块使用关系

config 模块被以下模块使用：

1. **`internal/commands/`**：命令层使用配置管理器
   - 使用 `config.Global()` - 获取全局配置管理器
   - 使用 `manager.Load()` - 加载配置
   - 使用 `manager.Config` 或便捷字段 - 访问配置
   - 使用 `manager.Save()` - 保存配置

2. **`internal/infrastructure/llm/`**：LLM 适配器使用配置
   - 使用 `config.Global()` - 获取全局配置管理器
   - 使用 `manager.LLMConfig` - 访问 LLM 配置
   - 使用 `llmConfig.CurrentProvider()` - 获取 provider 配置
   - 使用 `llmConfig.CurrentLanguage()` - 获取语言配置

3. **`internal/infrastructure/config/`**：配置适配器层
   - 使用 `config.GlobalRepoManager()` - 获取仓库配置管理器
   - 包装 git 模块，实现 `GitRepository` 接口

4. **`cmd/workflow/main.go`**：主程序初始化
   - 使用 `config.Global()` - 获取全局配置管理器
   - 使用 `manager.Load()` - 加载配置
   - 使用 `manager.LogConfig.Level` - 获取日志级别

### 调用流程

#### 全局配置加载流程

```
cmd/workflow/main.go
  ↓
config.Global()  // 获取单例
  ↓
manager.Load()   // 加载配置
  ↓
viper.ReadInConfig()  // 读取文件
  ↓
getGlobalConfig()  // 从 viper 解析配置
  ↓
更新 Config 字段和便捷字段
  ↓
返回配置管理器
```

#### 仓库配置加载流程

```
adapter/config.NewRepoManagerWithDefaultGit()
  ↓
config.GlobalRepoManager(gitRepo)  // 获取单例
  ↓
newRepoManager(gitRepo)  // 创建管理器
  ↓
generateRepoIDWithGit()  // 生成仓库 ID
  ↓
repoManager.Load()  // 加载配置
  ↓
加载公共配置和私有配置
  ↓
返回配置管理器
```

---

## 🎯 核心功能

### 1. 全局配置管理

**功能说明**：管理用户级别的全局配置，包括用户信息、认证配置、工具配置等。

**流程**：
1. 调用 `config.Global()` 获取全局配置管理器单例
2. 调用 `manager.Load()` 从文件加载配置
3. 直接访问配置字段（如 `manager.LLMConfig.Provider`）
4. 修改配置字段（如 `manager.LogConfig.Level = "debug"`）
5. 调用 `manager.Save()` 保存配置到文件

**示例**：
```go
import "github.com/zevwings/workflow/internal/config"

// 获取全局配置管理器
manager, err := config.Global()
if err != nil {
    return err
}

// 加载配置
if err := manager.Load(); err != nil {
    return err
}

// 访问配置
logLevel := manager.LogConfig.Level
llmProvider := manager.LLMConfig.Provider

// 修改配置
manager.LogConfig.Level = "debug"

// 保存配置
if err := manager.Save(); err != nil {
    return err
}
```

### 2. 仓库配置管理

**功能说明**：管理仓库级别的配置，包括项目公共配置和个人偏好配置。

**流程**：
1. 调用 `infrastructureconfig.NewRepoManagerWithDefaultGit()` 获取仓库配置管理器
2. 调用 `repoManager.Load()` 加载配置
3. 通过方法访问配置（如 `repoManager.GetBranchPrefix()`）
4. 调用 `repoManager.SaveTemplateConfig()` 保存模板配置

**示例**：
```go
import infrastructureconfig "github.com/zevwings/workflow/internal/infrastructure/config"

// 获取仓库配置管理器
repoManager, err := infrastructureconfig.NewRepoManagerWithDefaultGit("")
if err != nil {
    return err
}

// 加载配置
if err := repoManager.Load(); err != nil {
    return err
}

// 直接访问配置字段（推荐）
templateConfig := repoManager.TemplateConfig
commitFormat := repoManager.TemplateConfig.Commit["format"]
branchPrefix := repoManager.Config.Template.Branch["prefix"]

// 或者使用便捷方法（向后兼容）
branchPrefix := repoManager.GetBranchPrefix()
templateConfig := repoManager.GetTemplateConfig()

// 修改并保存配置
repoManager.TemplateConfig.Commit["type"] = "feat|fix|docs|style|refactor|test|chore"
if err := repoManager.Save(); err != nil {
    return err
}
```

### 3. LLM 配置管理

**功能说明**：管理 LLM 配置，支持多种 provider 和语言配置。

**流程**：
1. 从 `GlobalManager` 获取 `LLMConfig`
2. 调用 `CurrentProvider()` 获取当前 provider 的配置
3. 调用 `CurrentLanguage()` 获取当前语言配置

**示例**：
```go
manager, _ := config.Global()
manager.Load()

// 获取 provider 配置
apiKey, model, url, err := manager.LLMConfig.CurrentProvider()
if err != nil {
    return err
}

// 获取语言配置
lang, err := manager.LLMConfig.CurrentLanguage()
if err != nil {
    return err
}
```

### 4. 多语言支持

**功能说明**：提供多语言支持，包括语言查找、指令模板生成等功能。

**流程**：
1. 调用 `FindLanguage(code)` 查找语言
2. 调用 `GetLanguageInstruction(code)` 获取指令模板
3. 调用 `GetLanguageRequirement()` 获取语言要求

**示例**：
```go
import "github.com/zevwings/workflow/internal/config"

// 查找语言
lang := config.FindLanguage("zh-CN")
if lang != nil {
    fmt.Printf("Language: %s\n", lang.NativeName)
}

// 获取指令模板
instruction := config.GetLanguageInstruction("zh-CN")

// 获取支持的语言列表
codes := config.GetSupportedLanguageCodes()
```

---

## 📋 使用示例

### 示例 1: 读取全局配置

```go
import "github.com/zevwings/workflow/internal/config"

manager, err := config.Global()
if err != nil {
    panic(err)
}

if err := manager.Load(); err != nil {
    panic(err)
}

// 访问配置
fmt.Printf("Log Level: %s\n", manager.LogConfig.Level)
fmt.Printf("LLM Provider: %s\n", manager.LLMConfig.Provider)
fmt.Printf("GitHub Current: %s\n", manager.GitHubConfig.Current)
```

### 示例 2: 修改并保存配置

```go
import "github.com/zevwings/workflow/internal/config"

manager, _ := config.Global()
manager.Load()

// 修改配置
manager.LogConfig.Level = "debug"
manager.LLMConfig.Provider = "openai"
manager.LLMConfig.OpenAI.APIKey = "sk-xxx"

// 保存配置
if err := manager.Save(); err != nil {
    panic(err)
}
```

### 示例 3: 使用仓库配置

```go
import (
    "github.com/zevwings/workflow/internal/config"
    infrastructureconfig "github.com/zevwings/workflow/internal/infrastructure/config"
)

repoManager, err := infrastructureconfig.NewRepoManagerWithDefaultGit("")
if err != nil {
    panic(err)
}

repoManager.Load()

// 获取分支前缀
branchPrefix := repoManager.GetBranchPrefix()

// 获取模板配置
templateConfig := repoManager.GetTemplateConfig()
```

---

## 📝 扩展性

### 添加新配置字段

1. 在对应的配置结构体中添加字段（如 `LLMConfig`、`GitHubConfig` 等）
2. 在 `getGlobalConfig()` 或 `getRepoConfig()` 方法中添加字段的读取逻辑
3. 如果需要在 `GlobalManager` 中添加便捷字段，更新初始化代码

**示例**：
```go
// 1. 在 LLMConfig 中添加新字段
type LLMConfig struct {
    Provider string
    // 新增字段
    Timeout int `toml:"timeout,omitempty"`
}

// 2. 在 getGlobalConfig() 或 getRepoConfig() 中添加读取逻辑
cfg.LLM.Timeout = m.viper.GetInt("llm.timeout")
```

### 添加新配置类型

1. 创建新的配置结构体文件（如 `newconfig.go`）
2. 在 `GlobalConfig` 或 `RepoConfig` 中添加新字段
3. 在 `GlobalManager` 中添加便捷字段
4. 在 `getGlobalConfig()` 或 `getRepoConfig()` 中添加读取逻辑

**示例**：
```go
// 1. 创建 newconfig.go
type NewConfig struct {
    Field1 string `toml:"field1,omitempty"`
    Field2 int    `toml:"field2,omitempty"`
}

// 2. 在 GlobalConfig 或 RepoConfig 中添加
type GlobalConfig struct {
    // ...
    New NewConfig `toml:"new,omitempty"`
}

// 3. 在 GlobalManager 中添加便捷字段
type GlobalManager struct {
    // ...
    NewConfig *NewConfig
}

// 4. 在 getGlobalConfig() 或 getRepoConfig() 中添加读取逻辑
cfg.New.Field1 = m.viper.GetString("new.field1")
cfg.New.Field2 = m.viper.GetInt("new.field2")
```

---

## 📚 相关文档

- [模块 README](../../internal/config/README.md) - 基础使用说明

---

## ✅ 总结

config 模块采用清晰的单例模式和直接字段访问设计：

1. **单例模式**：`GlobalManager` 和 `RepoManager` 都是进程单例，确保配置一致性
2. **直接字段访问**：`GlobalManager` 提供公开字段，简化配置访问
3. **依赖注入**：`RepoManager` 通过接口实现依赖注入，解耦模块依赖
4. **配置分离**：区分全局配置和仓库配置，支持公共配置和私有配置
5. **多语言支持**：提供完整的多语言支持功能

**设计优势**：
- ✅ 线程安全：使用 `sync.Once` 确保线程安全的单例初始化
- ✅ 简洁易用：直接字段访问，代码更简洁直观
- ✅ 解耦设计：通过接口实现依赖注入，降低模块耦合
- ✅ 灵活扩展：易于添加新的配置类型和字段
- ✅ 性能优化：延迟加载和缓存机制提高性能

**当前实现状态**：
- ✅ 全局配置管理（GlobalManager）
- ✅ 仓库配置管理（RepoManager）
- ✅ LLM 配置管理
- ✅ 多语言支持
- ✅ 配置持久化
- ✅ 单例模式实现

---

**最后更新**: 2026-01-09

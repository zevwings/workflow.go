# config 模块

本模块提供了 Workflow CLI 的配置管理功能，包括全局配置和仓库配置的管理。

## 文件说明

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

### 核心文件

- **`global_manager.go`**：全局配置管理器，管理用户级别的全局配置（遵循 XDG 规范：`$XDG_CONFIG_HOME/workflow/config.toml`，默认 `~/.config/workflow/config.toml`），包含用户信息、认证配置（GitHub、Jira）、工具配置（LLM、Proxy、Log）。采用单例模式，提供直接字段访问和便捷方法。
- **`repo_manager.go`**：仓库配置管理器，管理仓库级别的配置，包括项目公共配置（`.workflow/config.toml`）和项目私有配置（遵循 XDG 规范：`$XDG_CONFIG_HOME/workflow/config/repository.toml`）。采用单例模式，通过依赖注入解耦 Git 模块依赖。
- **`types.go`**：定义 `GlobalConfig` 和 `RepoConfig` 结构体，统一所有子配置模块。
- **`helpers.go`**：提供通用的配置保存辅助函数 `SaveConfigToFile`。
- **`paths.go`**：提供 XDG Base Directory Specification 路径工具函数，包括 `ConfigDir()`、`DataDir()`、`StateDir()`、`CacheDir()` 等。
- **`languages.go`**：提供多语言支持，包括语言查找、指令模板生成等功能。
- **`llm.go`**：定义 LLM 配置结构体，提供 `CurrentProvider()` 和 `CurrentLanguage()` 方法。

## 快速开始

### 使用全局配置管理器

```go
import "github.com/zevwings/workflow/internal/config"

// 获取全局配置管理器单例
manager, err := config.Global()
if err != nil {
    // 处理错误
}

// 加载配置
err = manager.Load()
if err != nil {
    // 处理错误
}

// 直接访问配置字段（推荐）
logLevel := manager.LogConfig.Level
llmProvider := manager.LLMConfig.Provider
githubCurrent := manager.GitHubConfig.Current

// 修改配置
manager.LogConfig.Level = "debug"
manager.LLMConfig.Provider = "openai"

// 保存配置
err = manager.Save()
```

### 使用仓库配置管理器

```go
import (
    adapterconfig "github.com/zevwings/workflow/internal/adapter/config"
)

// 获取仓库配置管理器单例
repoManager, err := adapterconfig.NewRepoManagerWithDefaultGit("")
if err != nil {
    // 处理错误
}

// 加载配置
err = repoManager.Load()
if err != nil {
    // 处理错误
}

// 直接访问配置字段（推荐）
templateConfig := repoManager.TemplateConfig
commitFormat := repoManager.TemplateConfig.Commit["format"]
branchPrefix := repoManager.Config.Template.Branch["prefix"]

// 或者使用便捷方法（向后兼容）
branchPrefix := repoManager.GetBranchPrefix()
ignoreBranches := repoManager.GetIgnoreBranches()
templateConfig := repoManager.GetTemplateConfig()

// 修改配置
repoManager.TemplateConfig.Commit["format"] = "conventional"
repoManager.Config.Template.Branch["prefix"] = "feature/"

// 保存配置
err = repoManager.Save()
```

### 使用语言支持

```go
import "github.com/zevwings/workflow/internal/config"

// 查找语言
lang := config.FindLanguage("zh-CN")
if lang != nil {
    fmt.Printf("Language: %s\n", lang.NativeName)
}

// 获取语言指令模板
instruction := config.GetLanguageInstruction("zh-CN")

// 获取支持的语言代码列表
codes := config.GetSupportedLanguageCodes()
```

## 主要接口

### GlobalManager（全局配置管理器）

- `Global()` - 获取全局配置管理器单例
- `Load()` - 从文件加载配置到内存
- `Save()` - 保存当前配置到文件
- `SaveDefault()` - 保存默认配置
- `GetConfigPath()` - 获取配置文件路径
- `GetLLMConfig()` - 获取 LLM 配置（向后兼容）
- `GetGitHubConfig()` - 获取 GitHub 配置（向后兼容）
- `GetUserConfig()` - 获取用户配置（向后兼容）
- `GetJiraConfig()` - 获取 Jira 配置（向后兼容）
- `GetLogConfig()` - 获取日志配置（向后兼容）
- `GetProxyConfig()` - 获取代理配置（向后兼容）

**公开字段**（可直接访问）：
- `Config *GlobalConfig` - 完整配置
- `LLMConfig *LLMConfig` - LLM 配置
- `GitHubConfig *GitHubConfig` - GitHub 配置
- `UserConfig *UserConfig` - 用户配置
- `JiraConfig *JiraConfig` - Jira 配置
- `LogConfig *LogConfig` - 日志配置
- `ProxyConfig *ProxyConfig` - 代理配置

### RepoManager（仓库配置管理器）

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

**公开字段**（可直接访问）：
- `Config *RepoConfig` - 完整仓库公共配置
- `TemplateConfig *TemplateConfig` - 模板配置（指向 `Config.Template`）

### LLMConfig（LLM 配置）

- `CurrentProvider()` - 获取当前 provider 的配置（APIKey、Model、URL）
- `CurrentLanguage()` - 获取当前语言配置

### 语言支持函数

- `FindLanguage(code string)` - 查找语言（支持大小写不敏感和部分匹配）
- `GetLanguageInstruction(code string)` - 获取语言指令模板
- `GetLanguageRequirement(systemPrompt, languageCode string)` - 获取语言要求
- `GetSupportedLanguageCodes()` - 获取支持的语言代码列表
- `GetSupportedLanguageDisplayNames()` - 获取支持的语言显示名称列表

### 路径工具函数

- `ConfigDir()` - 获取配置目录（`$XDG_CONFIG_HOME/workflow`）
- `DataDir()` - 获取数据目录（`$XDG_DATA_HOME/workflow`）
- `StateDir()` - 获取状态目录（`$XDG_STATE_HOME/workflow`）
- `CacheDir()` - 获取缓存目录（`$XDG_CACHE_HOME/workflow`）

这些函数遵循 XDG Base Directory Specification，支持 Unix、Windows、macOS 等平台。

## 注意事项

1. **单例模式**：`GlobalManager` 和 `RepoManager` 都是进程单例，整个进程共享同一个实例。
2. **配置持久化**：直接修改字段只会修改内存中的值，必须显式调用 `Save()` 方法才会持久化到文件。
3. **配置加载**：访问配置前需要先调用 `Load()` 方法加载配置。
4. **字段访问**：`GlobalManager` 和 `RepoManager` 的配置字段都是公开的，可以直接访问和修改（推荐方式）。便捷方法（如 `GetTemplateConfig()`）保留以保持向后兼容。
5. **依赖注入**：`RepoManager` 通过 `GitRepository` 接口实现依赖注入，解耦对 git 模块的直接依赖。

## 依赖

- `github.com/spf13/viper` - 配置文件读取和管理
- `github.com/pelletier/go-toml/v2` - TOML 格式解析和序列化
- `github.com/adrg/xdg` - XDG Base Directory Specification 实现
- `github.com/zevwings/workflow/internal/logging` - 日志记录

## 相关文档

- [详细架构文档](../../docs/architecture/config.md) - 模块设计思路和实现细节

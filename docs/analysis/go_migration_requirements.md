# Workflow CLI Go 语言迁移需求文档

## 📋 文档概述

本文档基于现有的 Rust 实现（`workflow.rs`），详细梳理将项目迁移到 Go 语言所需的功能需求、技术选型和实现方案。

**项目信息：**
- **当前实现**：Rust (workflow.rs v1.6.9)
- **目标语言**：Go
- **项目类型**：跨平台 CLI 工具
- **主要功能**：Git 工作流自动化、PR 管理、Jira 集成、LLM 集成

---

## 🎯 迁移目标

### 核心目标
1. **功能对等**：完整迁移所有现有功能，保持 API 兼容性
2. **性能要求**：启动速度 < 200ms，二进制体积 < 30MB
3. **跨平台支持**：macOS、Linux、Windows（x86_64、ARM64）
4. **用户体验**：保持与 Rust 版本相同的用户体验

### 非功能性目标
- **开发效率**：利用 Go 的快速编译和开发体验
- **并发能力**：充分利用 Go 的 goroutine 进行并发处理
- **标准库优先**：尽可能使用 Go 标准库，减少外部依赖
- **代码可维护性**：清晰的代码结构和良好的错误处理

---

## 📦 功能模块清单

### 1. CLI 框架和命令系统

#### 1.1 核心命令结构
```
workflow
├── branch      # 分支管理
├── commit      # Commit 操作
├── pr          # Pull Request 操作
├── jira        # Jira 集成
├── stash       # Git stash 管理
├── tag         # Tag 管理
├── config      # 配置管理
├── github      # GitHub 账号管理
├── llm         # LLM 配置管理
├── log         # 日志级别管理
├── proxy       # 代理管理
├── repo        # 仓库管理
├── alias       # 别名管理
├── completion  # Shell 补全
├── check       # 环境检查
├── setup       # 初始化配置
├── migrate     # 配置迁移
├── update      # 更新工具
├── uninstall   # 卸载工具
└── version     # 版本信息
```

#### 1.2 技术选型
- **CLI 框架**：`cobra`（Go 生态最成熟的 CLI 框架）
- **参数解析**：`cobra` 内置支持
- **Shell 补全**：`cobra` 的 `cobra-cli` 工具生成
- **配置文件解析**：`viper`（支持 TOML、JSON、YAML）

**依赖映射：**
| Rust 依赖 | Go 替代方案 | 说明 |
|-----------|------------|------|
| `clap` | `cobra` | CLI 框架 |
| `clap_complete` | `cobra` 内置 | Shell 补全 |
| `toml` | `viper` + `pelletier/go-toml` | TOML 解析 |

---

### 2. Git 操作模块

#### 2.1 功能需求
- **分支管理**：创建、删除、切换、重命名、同步、忽略列表
- **Commit 操作**：amend、reword、squash
- **Stash 管理**：list、apply、drop、pop、push
- **Tag 管理**：删除（本地/远程）
- **仓库操作**：状态检查、配置管理、清理

#### 2.2 技术选型
- **Git 库**：`go-git`（纯 Go 实现，功能完善）
- **Git 命令包装**：`os/exec` 执行 git 命令（作为备选方案）

**依赖映射：**
| Rust 依赖 | Go 替代方案 | 说明 |
|-----------|------------|------|
| `git2` | `go-git` | Git 操作库 |
| `duct` | `os/exec` | 命令执行 |

**注意事项：**
- `go-git` 是纯 Go 实现，不依赖 libgit2，但某些底层操作可能不如 `git2` 灵活
- 对于复杂操作，可能需要结合 `os/exec` 执行 git 命令
- 需要实现 SSH 和 HTTPS 认证支持

---

### 3. HTTP 客户端模块

#### 3.1 功能需求
- **HTTP 请求**：GET、POST、PUT、DELETE、PATCH
- **认证支持**：Bearer Token、Basic Auth
- **重试机制**：可配置的重试策略（指数退避）
- **代理支持**：HTTP/HTTPS 代理
- **超时控制**：连接超时、读取超时
- **响应解析**：JSON、文本、二进制

#### 3.2 技术选型
- **HTTP 客户端**：`net/http`（Go 标准库）+ `resty`（可选，提供更高级功能）
- **重试机制**：自定义实现或使用 `go-retryablehttp`
- **JSON 解析**：`encoding/json`（标准库）

**依赖映射：**
| Rust 依赖 | Go 替代方案 | 说明 |
|-----------|------------|------|
| `reqwest` | `net/http` + `resty` | HTTP 客户端 |
| `serde_json` | `encoding/json` | JSON 序列化 |

**实现建议：**
- 优先使用 `net/http` 标准库
- 对于复杂需求（如自动重试、请求日志），考虑使用 `resty`
- 实现统一的 HTTP 客户端封装，支持配置化

---

### 4. 配置文件管理

#### 4.1 功能需求
- **配置文件格式**：TOML（主要）、JSON、YAML（导入/导出）
- **配置验证**：类型检查、必填项验证、格式验证
- **配置迁移**：版本迁移、格式转换
- **配置导入/导出**：支持部分导入、敏感信息过滤

#### 4.2 技术选型
- **TOML 解析**：`pelletier/go-toml` 或 `viper`（支持 TOML）
- **配置管理**：`viper`（统一配置管理框架）
- **验证**：自定义验证逻辑或使用 `go-playground/validator`

**依赖映射：**
| Rust 依赖 | Go 替代方案 | 说明 |
|-----------|------------|------|
| `toml` | `pelletier/go-toml` | TOML 解析 |
| `serde` | `viper` + 结构体标签 | 配置序列化 |

**配置文件结构：**
```go
type Settings struct {
    User    UserConfig    `toml:"user"`
    Jira    JiraConfig    `toml:"jira"`
    GitHub  GitHubConfig  `toml:"github"`
    Log     LogConfig     `toml:"log"`
    LLM     LLMConfig     `toml:"llm"`
    Codeup  CodeupConfig  `toml:"codeup"`
    Proxy   ProxyConfig   `toml:"proxy"`
}
```

---

### 5. API 集成模块

#### 5.1 GitHub API 集成
- **功能**：PR 创建、合并、关闭、查询、评论、批准
- **认证**：Personal Access Token
- **多账号管理**：账号列表、切换、更新

**技术选型：**
- **GitHub SDK**：`google/go-github`（官方推荐）或直接使用 REST API
- **API 客户端**：基于 `net/http` 封装

#### 5.2 Jira API 集成
- **功能**：Ticket 查询、评论、附件下载、变更历史
- **认证**：API Token + Basic Auth
- **日志操作**：下载、搜索、查找

**技术选型：**
- **Jira SDK**：`andygrunwald/go-jira` 或自定义实现
- **API 客户端**：基于 `net/http` 封装

#### 5.3 Codeup API 集成
- **功能**：PR 操作（部分支持）
- **认证**：CSRF Token + Cookie

**技术选型：**
- **自定义实现**：基于 `net/http` 封装（Codeup 没有官方 Go SDK）

---

### 6. LLM 集成模块

#### 6.1 功能需求
- **多提供者支持**：OpenAI、DeepSeek、Proxy（自定义代理）
- **统一接口**：配置驱动的提供者切换
- **功能**：PR 标题生成、PR 总结、分支名生成、Commit 消息重写
- **语言支持**：多语言提示词生成

#### 6.2 技术选型
- **OpenAI SDK**：`sashabaranov/go-openai`（非官方但成熟）
- **HTTP 客户端**：基于 `net/http` 封装统一接口

**依赖映射：**
| Rust 依赖 | Go 替代方案 | 说明 |
|-----------|------------|------|
| 自定义实现 | `sashabaranov/go-openai` + 自定义 | LLM 客户端 |

**实现建议：**
- 定义统一的 LLM 接口
- 实现多个提供者的适配器
- 支持配置驱动的提供者切换

---

### 7. 文件操作模块

#### 7.1 功能需求
- **ZIP 解压**：支持 ZIP 格式
- **TAR 解压**：支持 TAR、TAR.GZ 格式
- **文件校验**：SHA256 校验和验证
- **日志下载**：从 Jira 下载日志文件

#### 7.2 技术选型
- **ZIP 处理**：`archive/zip`（标准库）
- **TAR 处理**：`archive/tar` + `compress/gzip`（标准库）
- **校验和**：`crypto/sha256`（标准库）

**依赖映射：**
| Rust 依赖 | Go 替代方案 | 说明 |
|-----------|------------|------|
| `zip` | `archive/zip` | ZIP 解压 |
| `tar` + `flate2` | `archive/tar` + `compress/gzip` | TAR 解压 |
| `sha2` | `crypto/sha256` | SHA256 校验 |

---

### 8. 用户交互模块

#### 8.1 功能需求
- **交互式对话框**：确认、输入、选择、多选、表单
- **进度条**：文件下载、操作进度
- **Spinner**：长时间操作的加载指示器
- **表格显示**：格式化表格输出

#### 8.2 技术选型
- **交互式输入**：`survey`（功能完善）或 `promptui`（轻量）
- **进度条**：`cheggaaa/pb` 或 `schollz/progressbar`
- **Spinner**：`briandowns/spinner`
- **表格**：`olekukonko/tablewriter`

**依赖映射：**
| Rust 依赖 | Go 替代方案 | 说明 |
|-----------|------------|------|
| `inquire` / `dialoguer` | `survey` | 交互式对话框 |
| `indicatif` | `cheggaaa/pb` + `briandowns/spinner` | 进度条和 Spinner |
| `tabled` | `olekukonko/tablewriter` | 表格显示 |

---

### 9. 系统集成模块

#### 9.1 功能需求
- **剪贴板操作**：复制到剪贴板
- **浏览器打开**：打开 URL
- **Shell 检测**：检测当前 Shell 类型
- **Shell 重载**：重新加载 Shell 配置
- **代理管理**：系统代理开关

#### 9.2 技术选型
- **剪贴板**：`atotto/clipboard`（跨平台）
- **浏览器**：`pkg/browser`（标准库扩展）或自定义实现
- **Shell 检测**：自定义实现（检查环境变量）
- **代理管理**：平台特定实现（macOS/Linux/Windows）

**依赖映射：**
| Rust 依赖 | Go 替代方案 | 说明 |
|-----------|------------|------|
| `clipboard` | `atotto/clipboard` | 剪贴板操作 |
| `open` | `pkg/browser` 或自定义 | 浏览器打开 |
| 自定义实现 | 自定义实现 | Shell 检测和代理管理 |

---

### 10. 日志和追踪模块

#### 10.1 功能需求
- **日志级别**：none、error、warn、info、debug
- **日志输出**：控制台、文件
- **结构化日志**：支持 JSON 格式
- **追踪**：操作追踪和调试

#### 10.2 技术选型
- **日志库**：`log`（标准库）+ `logrus` 或 `zap`（结构化日志）
- **追踪**：自定义实现或使用 `opentelemetry-go`

**依赖映射：**
| Rust 依赖 | Go 替代方案 | 说明 |
|-----------|------------|------|
| `tracing` + `tracing-subscriber` | `log` + `logrus`/`zap` | 日志和追踪 |

---

### 11. 模板引擎

#### 11.1 功能需求
- **模板支持**：Handlebars 风格模板
- **变量替换**：分支名、Commit 消息、PR 标题等
- **条件渲染**：条件判断和循环

#### 11.2 技术选型
- **模板引擎**：`text/template`（标准库，Go 模板语法）或 `aymerick/raymond`（Handlebars 实现）

**依赖映射：**
| Rust 依赖 | Go 替代方案 | 说明 |
|-----------|------------|------|
| `handlebars` | `text/template` 或 `aymerick/raymond` | 模板引擎 |

**注意事项：**
- Go 标准库的 `text/template` 语法与 Handlebars 不同
- 如果需要保持模板兼容性，需要使用 `aymerick/raymond`
- 建议评估现有模板数量，决定是否需要 Handlebars 兼容

---

### 12. 错误处理

#### 12.1 功能需求
- **错误类型**：定义清晰的错误类型
- **错误链**：错误包装和上下文信息
- **用户友好**：友好的错误消息

#### 12.2 技术选型
- **错误处理**：Go 标准错误处理 + `pkg/errors`（错误包装）
- **错误类型**：自定义错误类型和错误码

**实现建议：**
- 定义统一的错误接口
- 使用 `fmt.Errorf` 和 `errors.Wrap` 进行错误包装
- 提供错误码和错误消息映射

---

### 13. 并发处理

#### 13.1 功能需求
- **并发下载**：多个文件并发下载
- **并发请求**：多个 API 请求并发执行
- **任务队列**：任务调度和执行

#### 13.2 技术选型
- **并发模型**：`goroutine` + `channel`（Go 原生支持）
- **任务池**：`sync.WaitGroup` 或自定义实现

**优势：**
- Go 的 goroutine 比 Rust 的 async/await 更简单
- Channel 提供了优雅的并发通信机制
- 无需复杂的异步运行时（如 tokio）

---

## 🏗️ 架构设计

### 项目结构

```
workflow.go/
├── cmd/
│   ├── workflow/        # 主命令入口
│   └── install/         # 安装命令
├── internal/
│   ├── cli/             # CLI 命令定义
│   ├── commands/        # 命令实现
│   ├── lib/             # 核心业务逻辑
│   │   ├── git/         # Git 操作
│   │   ├── github/      # GitHub API
│   │   ├── jira/        # Jira API
│   │   ├── llm/         # LLM 集成
│   │   ├── http/         # HTTP 客户端
│   │   ├── config/      # 配置管理
│   │   └── ...
│   └── utils/           # 工具函数
├── pkg/                 # 公共包（可选）
├── scripts/             # 安装脚本
├── docs/                # 文档
├── go.mod
├── go.sum
└── Makefile
```

### 模块依赖关系

```
CLI 层 (cmd/)
    ↓
命令层 (internal/commands/)
    ↓
业务逻辑层 (internal/lib/)
    ↓
工具层 (internal/utils/)
```

---

## 📊 依赖库对比表

### 核心依赖

| 功能模块 | Rust 依赖 | Go 替代方案 | 优先级 |
|---------|----------|------------|--------|
| CLI 框架 | `clap` | `cobra` | 高 |
| 配置管理 | `toml` + `serde` | `viper` + `pelletier/go-toml` | 高 |
| Git 操作 | `git2` | `go-git` | 高 |
| HTTP 客户端 | `reqwest` | `net/http` + `resty` | 高 |
| JSON 处理 | `serde_json` | `encoding/json` | 高 |
| 交互式输入 | `inquire` / `dialoguer` | `survey` | 中 |
| 进度条 | `indicatif` | `cheggaaa/pb` | 中 |
| 表格显示 | `tabled` | `olekukonko/tablewriter` | 中 |
| 日志 | `tracing` | `log` + `logrus` | 中 |
| 模板引擎 | `handlebars` | `text/template` 或 `aymerick/raymond` | 中 |
| 剪贴板 | `clipboard` | `atotto/clipboard` | 低 |
| 浏览器 | `open` | `pkg/browser` | 低 |
| 文件压缩 | `zip` + `tar` + `flate2` | `archive/zip` + `archive/tar` | 高 |
| 校验和 | `sha2` | `crypto/sha256` | 高 |
| GitHub SDK | 自定义 | `google/go-github` | 中 |
| Jira SDK | 自定义 | `andygrunwald/go-jira` | 中 |
| LLM SDK | 自定义 | `sashabaranov/go-openai` | 中 |

---

## 🚀 迁移计划

### 阶段一：基础设施（2-3 周）

**目标**：搭建项目基础架构和核心模块

1. **项目初始化**
   - 创建 Go 项目结构
   - 配置 `go.mod` 和依赖管理
   - 设置构建脚本和 Makefile

2. **CLI 框架集成**
   - 集成 `cobra` 框架
   - 实现命令结构（主命令和子命令）
   - 实现 Shell 补全生成

3. **配置管理**
   - 实现 TOML 配置文件读写
   - 实现配置验证和迁移
   - 实现配置导入/导出

4. **HTTP 客户端**
   - 实现统一的 HTTP 客户端封装
   - 实现重试机制
   - 实现认证支持

5. **错误处理**
   - 定义错误类型和错误码
   - 实现错误包装和用户友好消息

### 阶段二：核心功能（4-6 周）

**目标**：实现 Git 操作和基础命令

1. **Git 操作模块**
   - 集成 `go-git` 库
   - 实现分支管理（创建、删除、切换、重命名、同步）
   - 实现 Commit 操作（amend、reword、squash）
   - 实现 Stash 管理
   - 实现 Tag 管理

2. **基础命令**
   - 实现 `check` 命令（环境检查）
   - 实现 `setup` 命令（配置初始化）
   - 实现 `config` 命令（配置管理）
   - 实现 `completion` 命令（Shell 补全）

3. **用户交互**
   - 集成交互式对话框库
   - 实现进度条和 Spinner
   - 实现表格显示

### 阶段三：API 集成（3-4 周）

**目标**：实现 GitHub、Jira、LLM 集成

1. **GitHub API 集成**
   - 实现 GitHub API 客户端
   - 实现 PR 操作（创建、合并、关闭、查询、评论、批准）
   - 实现 GitHub 账号管理

2. **Jira API 集成**
   - 实现 Jira API 客户端
   - 实现 Ticket 查询、评论、附件下载
   - 实现日志操作（下载、搜索、查找）

3. **LLM 集成**
   - 实现统一的 LLM 接口
   - 实现 OpenAI、DeepSeek、Proxy 提供者
   - 实现 PR 标题生成、总结等功能

### 阶段四：高级功能（2-3 周）

**目标**：实现剩余功能和系统集成

1. **代理管理**
   - 实现系统代理检测和开关
   - 实现代理配置生成

2. **文件操作**
   - 实现 ZIP/TAR 解压
   - 实现校验和验证
   - 实现日志下载

3. **系统集成**
   - 实现剪贴板操作
   - 实现浏览器打开
   - 实现 Shell 检测和重载

4. **生命周期管理**
   - 实现 `install` 命令
   - 实现 `update` 命令
   - 实现 `uninstall` 命令

### 阶段五：测试和优化（2-3 周）

**目标**：完善测试、优化性能、完善文档

1. **测试**
   - 单元测试（覆盖率 > 80%）
   - 集成测试
   - E2E 测试

2. **性能优化**
   - 启动速度优化
   - 二进制体积优化
   - 内存占用优化

3. **文档**
   - API 文档
   - 用户手册
   - 迁移指南

---

## ⚠️ 技术挑战和解决方案

### 挑战 1：Git 操作兼容性

**问题**：`go-git` 是纯 Go 实现，某些底层操作可能不如 `git2` 灵活。

**解决方案**：
- 优先使用 `go-git` 实现大部分功能
- 对于复杂操作，使用 `os/exec` 执行 git 命令作为备选
- 实现 Git 命令包装层，统一接口

### 挑战 2：模板引擎兼容性

**问题**：Go 标准库的 `text/template` 语法与 Handlebars 不同。

**解决方案**：
- 评估现有模板数量和使用情况
- 如果模板较少，考虑迁移到 Go 模板语法
- 如果模板较多，使用 `aymerick/raymond` 保持兼容性

### 挑战 3：错误处理模式

**问题**：Go 的错误处理模式与 Rust 的 `Result<T, E>` 不同。

**解决方案**：
- 使用 `pkg/errors` 进行错误包装
- 定义统一的错误接口和错误码
- 实现错误链和上下文信息

### 挑战 4：并发模型差异

**问题**：Go 的 goroutine 模型与 Rust 的 async/await 不同。

**解决方案**：
- 利用 Go 的 goroutine 简化并发代码
- 使用 channel 进行并发通信
- 使用 `sync.WaitGroup` 等待并发任务完成

### 挑战 5：类型安全

**问题**：Go 的类型系统不如 Rust 严格。

**解决方案**：
- 充分利用 Go 的接口和类型断言
- 使用代码生成工具（如 `stringer`）生成类型安全的代码
- 编写充分的单元测试

---

## 📈 性能目标

### 启动速度
- **目标**：< 200ms（Rust 版本 < 100ms）
- **优化策略**：
  - 减少初始化代码
  - 延迟加载非关键模块
  - 优化依赖导入

### 二进制体积
- **目标**：< 30MB（Rust 版本 5-15MB）
- **优化策略**：
  - 使用 `-ldflags="-s -w"` 减小体积
  - 移除调试信息
  - 静态链接依赖

### 内存占用
- **目标**：< 50MB（运行时）
- **优化策略**：
  - 及时释放资源
  - 使用对象池
  - 优化数据结构

---

## 🔧 开发工具和流程

### 开发工具
- **Go 版本**：Go 1.21+
- **代码格式化**：`gofmt`、`goimports`
- **Lint 工具**：`golangci-lint`
- **测试工具**：`go test`、`testify`
- **构建工具**：`Makefile`、`goreleaser`（发布）

### 开发流程
1. **代码规范**：遵循 Go 官方代码规范
2. **测试驱动**：先写测试，再实现功能
3. **代码审查**：所有代码必须经过审查
4. **持续集成**：GitHub Actions 自动构建和测试

---

## 📝 迁移检查清单

### 功能完整性
- [ ] CLI 命令系统（所有命令）
- [ ] Git 操作（分支、Commit、Stash、Tag）
- [ ] GitHub API 集成
- [ ] Jira API 集成
- [ ] LLM 集成
- [ ] 配置管理
- [ ] 文件操作（ZIP/TAR 解压）
- [ ] 用户交互（对话框、进度条）
- [ ] 系统集成（剪贴板、浏览器）
- [ ] Shell 补全
- [ ] 代理管理
- [ ] 日志和追踪

### 非功能性需求
- [ ] 跨平台支持（macOS、Linux、Windows）
- [ ] 启动速度 < 200ms
- [ ] 二进制体积 < 30MB
- [ ] 测试覆盖率 > 80%
- [ ] 文档完整性
- [ ] 错误处理完善
- [ ] 用户体验一致

### 发布准备
- [ ] 构建脚本
- [ ] 安装脚本（macOS/Linux/Windows）
- [ ] 发布流程
- [ ] 版本管理
- [ ] 更新机制

---

## 📚 参考资料

### Go 生态资源
- [Cobra 官方文档](https://github.com/spf13/cobra)
- [Viper 官方文档](https://github.com/spf13/viper)
- [go-git 官方文档](https://github.com/go-git/go-git)
- [Go 官方文档](https://go.dev/doc/)

### 项目参考
- [现有 Rust 实现](../workflow.rs/)
- [技术选型对比](./python_rust_go_comparison.md)
- [Python CLI 框架分析](./python_cli_frameworks_analysis.md)

---

## 🎯 总结

### 迁移优势
1. **开发效率**：Go 的快速编译和简洁语法
2. **并发能力**：goroutine 简化并发编程
3. **标准库**：丰富的标准库减少外部依赖
4. **团队熟悉度**：如果团队更熟悉 Go

### 迁移挑战
1. **性能**：启动速度和二进制体积可能不如 Rust
2. **类型安全**：Go 的类型系统不如 Rust 严格
3. **生态差异**：某些库的功能和 API 可能不同

### 建议
- **如果团队熟悉 Go**：迁移是可行的，可以充分利用 Go 的优势
- **如果追求极致性能**：建议继续使用 Rust
- **如果快速开发优先**：Go 是更好的选择

---

**最后更新**: 2025-12-28


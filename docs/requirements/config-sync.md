# 跨平台配置同步功能需求

> 实现 Workflow CLI 配置在 Linux/Windows/macOS 之间的安全同步功能，支持多种同步后端（Git、云存储、Gist），确保敏感信息（API Token、密钥）的安全性。

---

## 📋 目录

- [概述](#-概述)
- [需求分析](#-需求分析)
- [技术方案](#-技术方案)
- [设计方案](#-设计方案)
- [实施计划](#-实施计划)
- [安全考虑](#-安全考虑)
- [相关文档](#-相关文档)

---

## 📋 概述

### 当前状态

- **状态**: ⏳ 待实施
- **实现度**: 0%
- **优先级**: 中
- **分类**: 配置管理 / 跨平台功能

### 目标

实现 Workflow CLI 配置文件的跨平台同步功能，允许用户在多个设备（Linux、Windows、macOS）之间安全地同步配置，包括：

1. **安全同步**：确保包含敏感信息（API Token、密钥）的配置文件在同步过程中得到加密保护
2. **多后端支持**：支持多种同步后端（私有 Git 仓库、云存储、GitHub Gist）
3. **自动化**：提供简单的 CLI 命令实现配置的上传和下载
4. **冲突处理**：处理多设备修改配置时的冲突情况

### 已完成

- ✅ 配置系统已实现（`internal/config`）
- ✅ 配置文件位置已标准化（遵循 XDG 规范）
- ✅ 敏感信息识别和过滤机制已实现

### 待实现

- ⏳ 配置加密/解密功能
- ⏳ 同步后端抽象层
- ⏳ Git 后端实现
- ⏳ 云存储后端实现
- ⏳ Gist 后端实现
- ⏳ 密钥管理集成（系统密钥链）
- ⏳ 冲突检测和合并策略
- ⏳ CLI 命令实现（`workflow config sync`）

---

## 需求分析

### 功能需求

#### 1. 配置加密

**需求描述**：配置文件包含敏感信息（GitHub API Token、Jira API Token、LLM API Keys），必须加密后才能同步。

**具体要求**：
- 支持 AES-256 加密算法
- 加密密钥存储在系统密钥链中（macOS Keychain、Windows Credential Manager、Linux Secret Service）
- 支持全文件加密和字段级加密两种模式
- 提供加密/解密命令

**命令示例**：

```bash
# 加密配置文件
workflow config sync encrypt

# 解密配置文件
workflow config sync decrypt
```

#### 2. 多后端同步

**需求描述**：支持多种同步后端，用户可以选择最适合的方案。

**支持的后端**：
- **Git 后端**：私有 Git 仓库（GitHub/GitLab/Gitea）
- **云存储后端**：OneDrive、Google Drive、Dropbox、坚果云等
- **Gist 后端**：GitHub Gist（私有）

**命令示例**：

```bash
# 初始化同步（选择后端）
workflow config sync init

# 上传配置到远程
workflow config sync push

# 从远程下载配置
workflow config sync pull

# 查看同步状态
workflow config sync status
```

#### 3. 冲突处理

**需求描述**：当多个设备修改配置时，需要检测和处理冲突。

**处理策略**：
- 检测远程配置和本地配置的差异
- 提供交互式合并界面
- 支持自动合并（基于时间戳或用户选择）
- 保留冲突备份

**命令示例**：

```bash
# 拉取配置（自动处理冲突）
workflow config sync pull --auto-merge

# 手动解决冲突
workflow config sync pull --interactive
```

#### 4. 密钥管理

**需求描述**：加密密钥需要安全存储，使用系统密钥链。

**具体要求**：
- macOS：使用 Keychain
- Windows：使用 Credential Manager
- Linux：使用 Secret Service（如 GNOME Keyring、KWallet）

**实现建议**：
- 使用 `github.com/99designs/keyring` 库
- 密钥标识：`workflow-cli:config-encryption-key`
- 首次使用时生成随机密钥并存储

### 非功能需求

#### 1. 安全性

- **加密强度**：使用 AES-256-GCM 加密算法
- **密钥管理**：密钥不存储在配置文件中，使用系统密钥链
- **传输安全**：所有网络传输使用 HTTPS/TLS
- **访问控制**：配置文件权限设置为 600（仅所有者可读）

#### 2. 可用性

- **简单易用**：提供直观的 CLI 命令
- **错误处理**：清晰的错误提示和解决建议
- **文档完善**：提供详细的使用文档和示例

#### 3. 性能

- **加密性能**：加密/解密操作应在 1 秒内完成
- **同步性能**：同步操作应在 5 秒内完成（网络正常时）

#### 4. 兼容性

- **跨平台**：支持 Linux、Windows、macOS
- **向后兼容**：不影响现有配置系统
- **可选功能**：同步功能为可选，不影响核心功能

---

## 技术方案

### 方案对比

#### 方案 A：加密 + 私有 Git 仓库（推荐）

**优点**：
- 版本控制，可追溯变更历史
- 私有仓库，访问可控
- 支持自动化（CI/CD）
- 跨平台支持好

**缺点**：
- 需要 Git 操作（pull/push）
- 需要管理 Git 凭据
- 需要额外工具（git-crypt/git-secret）或自行实现加密

**适用场景**：需要版本控制，可接受 Git 操作

#### 方案 B：AES 加密 + 云存储（推荐）

**优点**：
- 使用简单，自动同步
- 跨平台客户端支持好
- 无需额外工具

**缺点**：
- 需要加密配置内容
- 依赖第三方服务
- 同步可能有延迟

**适用场景**：追求简单，已有云存储

#### 方案 C：混合方案

**优点**：
- 灵活，可扩展
- 支持多种后端
- 用户可选择最适合的方案

**缺点**：
- 实现复杂度较高
- 需要维护多个后端

**适用场景**：需要灵活性和可扩展性

### 推荐方案

**采用方案 C（混合方案）**，原因：
1. 提供最大灵活性，用户可根据需求选择
2. 可逐步实现，先实现一个后端，再扩展其他后端
3. 便于未来扩展新的同步后端

---

## 设计方案

### 架构设计

```
┌─────────────────────────────────────────┐
│         CLI Commands Layer              │
│  workflow config sync [push/pull/init]  │
└─────────────────┬───────────────────────┘
                  │
┌─────────────────▼───────────────────────┐
│      Sync Manager (核心协调层)          │
│  - 加密/解密管理                         │
│  - 冲突检测和处理                        │
│  - 后端路由                              │
└─────────────────┬───────────────────────┘
                  │
      ┌───────────┼───────────┐
      │           │           │
┌─────▼─────┐ ┌───▼────┐ ┌───▼────┐
│ Git Backend│ │ Cloud │ │ Gist   │
│            │ │Backend │ │Backend │
└────────────┘ └────────┘ └────────┘
      │           │           │
┌─────▼───────────────────────▼─────┐
│      Keyring (密钥管理)             │
│  - macOS Keychain                  │
│  - Windows Credential Manager      │
│  - Linux Secret Service            │
└────────────────────────────────────┘
```

### 模块设计

#### 1. 加密模块 (`internal/config/sync/encrypt`)

**职责**：
- 配置文件加密/解密
- 密钥生成和管理
- 加密算法实现（AES-256-GCM）

**接口**：

```go
type Encryptor interface {
    Encrypt(data []byte) ([]byte, error)
    Decrypt(encrypted []byte) ([]byte, error)
    GetKey() ([]byte, error)
}
```

#### 2. 同步后端接口 (`internal/config/sync/backend`)

**职责**：
- 定义同步后端接口
- 提供统一的同步操作接口

**接口**：

```go
type Backend interface {
    Name() string
    Push(data []byte) error
    Pull() ([]byte, error)
    Exists() (bool, error)
    Init() error
}
```

#### 3. Git 后端 (`internal/config/sync/backend/git`)

**职责**：
- 实现 Git 仓库同步
- 处理 Git 操作（clone、push、pull）
- 管理 Git 凭据

**实现要点**：
- 使用 `gogit` 或 `go-git` 库
- 支持 SSH 和 HTTPS 两种方式
- 自动处理 Git 凭据

#### 4. 云存储后端 (`internal/config/sync/backend/cloud`)

**职责**：
- 实现云存储同步
- 支持多种云存储服务

**实现要点**：
- 检测云存储同步目录
- 使用文件系统监控或轮询检测变更
- 支持 OneDrive、Google Drive、Dropbox 等

#### 5. Gist 后端 (`internal/config/sync/backend/gist`)

**职责**：
- 实现 GitHub Gist 同步
- 使用 GitHub API 上传/下载配置

**实现要点**：
- 使用 GitHub API v3/v4
- 支持私有 Gist
- 使用 GitHub Token 认证

#### 6. 密钥管理 (`internal/config/sync/keyring`)

**职责**：
- 封装系统密钥链操作
- 提供跨平台密钥存储接口

**实现要点**：
- 使用 `github.com/99designs/keyring`
- 密钥标识：`workflow-cli:config-encryption-key`
- 首次使用时生成随机密钥

### 数据流设计

#### 上传流程（Push）

```
1. 读取本地配置文件
2. 使用密钥加密配置
3. 根据配置的后端类型，调用对应的 Push 方法
4. 上传加密后的配置到远程
5. 显示成功消息
```

#### 下载流程（Pull）

```
1. 根据配置的后端类型，调用对应的 Pull 方法
2. 从远程下载加密后的配置
3. 检测本地配置是否存在
4. 如果存在，比较差异（冲突检测）
5. 如果有冲突，提供合并选项
6. 使用密钥解密配置
7. 保存到本地配置文件
8. 显示成功消息
```

### 配置文件扩展

在现有配置中添加同步相关配置：

```toml
[sync]
enabled = true
backend = "git"  # git, cloud, gist
auto_sync = false

[sync.git]
repository = "https://github.com/user/config-repo.git"
branch = "main"
path = "workflow-config.enc"

[sync.cloud]
provider = "onedrive"  # onedrive, googledrive, dropbox
path = ".workflow/config.enc"

[sync.gist]
gist_id = "abc123..."
```

---

## 实施计划

### 阶段 1：基础功能（核心加密和密钥管理）

**时间估算**：5-7 天

- [ ] **任务 1.1**：实现密钥管理模块
  - 集成 `github.com/99designs/keyring`
  - 实现跨平台密钥存储（macOS Keychain、Windows Credential Manager、Linux Secret Service）
  - 实现密钥生成和获取接口
- [ ] **任务 1.2**：实现加密模块
  - 实现 AES-256-GCM 加密算法
  - 实现加密/解密接口
  - 编写单元测试
- [ ] **任务 1.3**：实现 CLI 加密命令
  - `workflow config sync encrypt` - 加密配置文件
  - `workflow config sync decrypt` - 解密配置文件
  - 集成到现有配置系统

**依赖关系**：
- 需要 `github.com/99designs/keyring` 依赖
- 需要 `golang.org/x/crypto` 依赖

### 阶段 2：同步后端抽象层和 Git 后端

**时间估算**：7-10 天

- [ ] **任务 2.1**：实现同步后端接口
  - 定义 `Backend` 接口
  - 实现同步管理器（Sync Manager）
  - 实现后端注册机制
- [ ] **任务 2.2**：实现 Git 后端
  - 集成 `go-git` 库
  - 实现 Git 仓库操作（clone、push、pull）
  - 实现 Git 凭据管理
  - 支持 SSH 和 HTTPS
- [ ] **任务 2.3**：实现 CLI Git 同步命令
  - `workflow config sync init --backend git` - 初始化 Git 同步
  - `workflow config sync push` - 推送配置
  - `workflow config sync pull` - 拉取配置
  - `workflow config sync status` - 查看同步状态

**依赖关系**：
- 需要 `github.com/go-git/go-git` 依赖
- 依赖阶段 1 的加密模块

### 阶段 3：云存储后端

**时间估算**：5-7 天

- [ ] **任务 3.1**：实现云存储检测
  - 检测常见的云存储同步目录
  - 支持 OneDrive、Google Drive、Dropbox、坚果云
- [ ] **任务 3.2**：实现云存储后端
  - 实现文件系统监控或轮询
  - 实现配置文件的自动同步检测
- [ ] **任务 3.3**：实现 CLI 云存储同步命令
  - `workflow config sync init --backend cloud` - 初始化云存储同步
  - 集成到现有同步命令

**依赖关系**：
- 依赖阶段 1 和阶段 2

### 阶段 4：Gist 后端和冲突处理

**时间估算**：5-7 天

- [ ] **任务 4.1**：实现 Gist 后端
  - 集成 GitHub API
  - 实现 Gist 创建、更新、读取
  - 支持私有 Gist
- [ ] **任务 4.2**：实现冲突检测和处理
  - 实现配置差异检测
  - 实现交互式合并界面
  - 实现自动合并策略
  - 实现冲突备份
- [ ] **任务 4.3**：实现 CLI Gist 同步命令
  - `workflow config sync init --backend gist` - 初始化 Gist 同步
  - `workflow config sync pull --auto-merge` - 自动合并
  - `workflow config sync pull --interactive` - 交互式合并

**依赖关系**：
- 需要 GitHub API 客户端库
- 依赖阶段 1、2、3

### 阶段 5：测试和文档

**时间估算**：3-5 天

- [ ] **任务 5.1**：编写单元测试
  - 加密模块测试
  - 各后端测试
  - 冲突处理测试
- [ ] **任务 5.2**：编写集成测试
  - 端到端同步测试
  - 跨平台测试
- [ ] **任务 5.3**：编写文档
  - 使用文档
  - API 文档
  - 故障排除指南

**依赖关系**：
- 依赖所有前序阶段

---

## 安全考虑

### 加密安全

1. **加密算法**：使用 AES-256-GCM，提供认证加密
2. **密钥管理**：密钥存储在系统密钥链，不存储在配置文件中
3. **密钥生成**：使用加密安全的随机数生成器（`crypto/rand`）

### 传输安全

1. **HTTPS/TLS**：所有网络传输使用 HTTPS
2. **Git SSH**：支持 SSH 方式，避免传输密钥

### 存储安全

1. **文件权限**：配置文件权限设置为 600（仅所有者可读）
2. **密钥隔离**：加密密钥与配置文件分离存储

### 访问控制

1. **后端认证**：各后端需要相应的认证（Git Token、GitHub Token 等）
2. **权限检查**：同步前检查用户权限

---

## 📊 任务统计

| 状态 | 数量 | 说明 |
|-----|------|------|
| ✅ 已完成 | 0 个 | 尚未开始实施 |
| 🚧 进行中 | 0 个 | 尚未开始实施 |
| ⏳ 待实施 | 20 个 | 所有任务待实施 |
| **总计** | **20** | - |

---

## 🔍 故障排除

### 问题 1：密钥获取失败

**症状**：执行加密/解密操作时提示密钥获取失败

**解决方案**：

1. 检查系统密钥链是否可用
2. 尝试重新生成密钥：`workflow config sync encrypt --regenerate-key`
3. 检查密钥链权限设置

### 问题 2：Git 同步失败

**症状**：执行 `workflow config sync push` 时提示 Git 操作失败

**解决方案**：

1. 检查 Git 仓库 URL 是否正确
2. 检查 Git 凭据是否配置
3. 检查网络连接
4. 检查仓库权限

### 问题 3：冲突处理

**症状**：拉取配置时提示配置冲突

**解决方案**：

1. 使用 `--interactive` 选项手动解决冲突
2. 使用 `--auto-merge` 选项自动合并（基于时间戳）
3. 查看冲突备份文件

### 问题 4：云存储同步目录未找到

**症状**：初始化云存储同步时提示找不到同步目录

**解决方案**：

1. 确认云存储客户端已安装并运行
2. 确认同步目录已设置
3. 手动指定同步目录路径

---

## 📚 相关文档

- [配置架构文档](../architecture/config.md) - 配置系统架构说明
- [安全规范](../development/references/security.md) - 安全开发规范
- [配置参考文档](../development/references/configuration.md) - 配置使用参考

---

## ✅ 检查清单

实施本需求时，请确保：

- [ ] 所有敏感信息（API Token、密钥）必须加密
- [ ] 加密密钥存储在系统密钥链，不存储在配置文件中
- [ ] 所有网络传输使用 HTTPS/TLS
- [ ] 配置文件权限设置为 600（仅所有者可读）
- [ ] 提供清晰的错误提示和解决建议
- [ ] 编写完整的单元测试和集成测试
- [ ] 编写详细的使用文档
- [ ] 支持跨平台（Linux、Windows、macOS）
- [ ] 向后兼容，不影响现有配置系统
- [ ] 同步功能为可选，不影响核心功能

---

**最后更新**: 2025-01-28

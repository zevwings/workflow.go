# 配置管理规范

> 本文档定义了 Workflow CLI 项目的配置管理规范和最佳实践，所有贡献者都应遵循这些规范。

---

## 📋 目录

- [概述](#-概述)
- [配置验证规则](#-配置验证规则)
- [配置迁移规则](#-配置迁移规则)
- [配置默认值管理规则](#-配置默认值管理规则)
- [相关文档](#-相关文档)

---

## 📋 概述

本文档定义了配置管理规范，包括配置验证、配置迁移和默认值管理。

### 核心原则

- **验证配置**：所有配置加载时必须验证配置的有效性
- **迁移机制**：配置格式变化时必须提供迁移机制
- **默认值**：所有配置项都应提供合理的默认值

### 使用场景

- 添加新配置项时参考
- 修改配置格式时使用
- 配置验证时参考

---

## 配置验证规则

所有配置加载时必须验证配置的有效性，确保配置的正确性和安全性。

### 配置验证时机

配置验证应在以下时机进行：

1. **配置加载时**：使用 `workflow config validate` 命令验证配置
2. **配置更新时**：更新配置后自动验证
3. **程序启动时**：可选，在关键配置缺失时提示用户

### 配置验证内容

配置验证应检查以下内容：

- **格式验证**：配置文件格式是否正确（TOML、JSON、YAML）
- **必需字段**：检查必需字段是否存在
- **字段类型**：验证字段类型是否正确
- **值有效性**：验证字段值的有效性
  - URL 格式（必须以 `http://` 或 `https://` 开头）
  - 邮箱格式（必须包含 `@`）
  - 路径格式（路径是否存在、是否可访问）
  - 枚举值（是否在允许的枚举值范围内）

### 配置验证实现

使用配置验证函数进行配置验证：

```go
import "github.com/your-org/workflow/internal/commands/config"

// 验证配置
err := config.Validate(nil, false, false)
if err != nil {
	return err
}

// 自动修复配置错误
err = config.Validate(nil, true, false)
if err != nil {
	return err
}

// 严格模式（警告视为错误）
err = config.Validate(nil, false, true)
if err != nil {
	return err
}
```

### 配置错误消息

配置验证失败时，应提供清晰的错误消息：

```go
// ✅ 好的错误消息
type ValidationError struct {
	Field         string
	Message       string
	Fixable       bool
	FixSuggestion string
}

err := ValidationError{
	Field:   "jira.email",
	Message: fmt.Sprintf("Invalid email format: '%s'", email),
	Fixable: false,
}

// ✅ 提供修复建议
err = ValidationError{
	Field:   "jira.service_address",
	Message: fmt.Sprintf("Invalid URL format: '%s' (must start with http:// or https://)", serviceAddress),
	Fixable: true,
	FixSuggestion: fmt.Sprintf("Updated 'jira.service_address' from '%s' to 'https://%s'",
		serviceAddress, strings.TrimPrefix(strings.TrimPrefix(serviceAddress, "http://"), "https://")),
}
```

### 配置验证失败处理

配置验证失败时的处理流程：

1. **显示错误信息**：列出所有验证错误和警告
2. **提供修复建议**：对于可修复的错误，提供修复建议
3. **自动修复**：使用 `--fix` 选项自动修复可修复的错误
4. **退出码**：验证失败时返回非零退出码（用于 CI/CD）

---

## 配置迁移规则

当配置格式发生变化时，必须提供配置迁移机制，确保用户配置能够平滑升级。

### 迁移版本管理

**重要**：迁移版本号**独立于**软件版本号！

- **软件版本**（如 `1.4.8`）：表示软件本身的版本，在 `go.mod` 中定义
- **迁移版本**（如 `v1.0.0`）：表示**配置格式的版本**，只有当配置格式发生变化时才需要迁移

**迁移版本命名规范**：
- 使用语义化版本（Semantic Versioning）
- **Major**（v1.0.0 → v2.0.0）：重大配置格式变化，不向后兼容
- **Minor**（v1.0.0 → v1.1.0）：新增配置项或格式变化，向后兼容
- **Patch**（v1.0.0 → v1.0.1）：通常不使用，因为配置格式变化通常需要 minor 或 major

### 添加新迁移版本的步骤

当需要添加新的迁移版本时，按以下步骤操作：

1. **创建迁移实现文件**：
   ```go
   // internal/commands/migrate/v1_1_0.go
   // v1.1.0 迁移实现

   package migrate

   func MigrateV1_1_0(dryRun bool, cleanup bool) error {
       // 1. 检测需要迁移的内容
       // 2. 执行迁移逻辑
       // 3. 可选：清理旧文件
       return nil
   }
   ```

2. **在 migrate.go 中注册新版本**：
   ```go
   // internal/commands/migrate/migrate.go
   var Migrations = []string{
       "v1.0.0",
       "v1.1.0", // 添加新版本
   }
   ```

4. **创建迁移文档**：
   - 使用迁移文档模板：`docs/templates/migration/migration.template`
   - 创建迁移文档：`docs/migration/{旧版本}-to-{新版本}.md`
   - 更新迁移文档索引：`docs/migration/README.md`

5. **创建迁移脚本**（如需要）：
   - Shell 脚本：`scripts/migrate/{旧版本}-to-{新版本}.sh`
   - PowerShell 脚本：`scripts/migrate/{旧版本}-to-{新版本}.ps1`

### 迁移实现要求

迁移实现应遵循以下要求：

- **幂等性**：迁移可以多次执行而不产生副作用
- **可回滚**：迁移前备份原始配置，支持回滚
- **预览模式**：支持 `--dry-run` 预览迁移结果
- **清理选项**：支持 `--cleanup` 清理旧配置文件
- **错误处理**：迁移失败时提供清晰的错误信息
- **日志记录**：记录迁移历史，避免重复迁移

### 迁移历史管理

使用 `MigrationHistory` 管理迁移历史：

```go
import "github.com/your-org/workflow/internal/commands/migrate"

// 检查是否已迁移
hasMigrated, err := migrate.HasMigrated("v1.1.0")
if err != nil {
	return err
}
if hasMigrated {
	return nil // 已迁移，跳过
}

// 执行迁移
err = migrate.MigrateV1_1_0(dryRun, cleanup)
if err != nil {
	return err
}

// 记录迁移历史
err = migrate.Record("v1.1.0")
if err != nil {
	return err
}
```

---

## 配置默认值管理规则

所有配置项都应提供合理的默认值，确保程序在配置文件缺失或字段缺失时仍能正常运行。

### 默认值定义方式

使用结构体字段标签定义默认值：

```go
type Settings struct {
	// Jira 配置
	Jira JiraSettings `toml:"jira,omitempty"`

	// 日志配置
	Log LogSettings `toml:"log,omitempty"`
}

// NewSettings 创建默认设置
func NewSettings() *Settings {
	return &Settings{
		Jira: NewJiraSettings(),
		Log:  NewLogSettings(),
	}
}
```

### 默认值变更的影响评估

修改默认值时，必须评估以下影响：

1. **向后兼容性**：默认值变更是否会影响现有用户
2. **用户体验**：新默认值是否提供更好的用户体验
3. **安全性**：新默认值是否引入安全风险
4. **性能影响**：新默认值是否影响程序性能
5. **文档更新**：必须更新相关文档说明默认值变更

**默认值变更流程**：

1. **评估影响**：评估默认值变更的影响范围
2. **更新代码**：更新默认值实现
3. **更新文档**：更新配置文档和 CHANGELOG.md
4. **测试验证**：验证新默认值的行为
5. **发布说明**：在发布说明中说明默认值变更

### 配置加载时的默认值处理

配置加载时，如果配置文件不存在或字段缺失，应使用默认值：

```go
// Load 从配置文件加载设置
// 如果配置文件不存在或字段缺失，使用默认值
func Load() (*Settings, error) {
	configPath, err := GetConfigPath()
	if err != nil {
		return NewSettings(), nil // 路径获取失败，返回默认值
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return NewSettings(), nil // 文件不存在，返回默认值
	}

	content, err := os.ReadFile(configPath)
	if err != nil {
		return NewSettings(), nil // 读取失败，返回默认值
	}

	var settings Settings
	if err := toml.Unmarshal(content, &settings); err != nil {
		return NewSettings(), nil // 解析失败，返回默认值
	}

	// 合并默认值
	return mergeWithDefaults(&settings), nil
}
```

---

## 🔍 故障排除

### 问题 1：配置验证失败

**症状**：配置验证失败但错误信息不清晰

**解决方案**：

1. 检查配置验证错误消息是否清晰
2. 提供修复建议
3. 支持自动修复（如适用）

### 问题 2：配置迁移失败

**症状**：配置迁移失败

**解决方案**：

1. 检查迁移实现是否幂等
2. 提供回滚机制
3. 记录迁移历史避免重复迁移

---

## 📚 相关文档

### 架构文档

- [配置验证命令架构文档](../../architecture/config.md#4-配置验证命令-validaters) - 详细的配置验证实现说明
- [迁移系统架构文档](../../../architecture/migrate.md) - 详细的迁移系统说明
- [Settings 模块架构文档](../../architecture/settings.md#3-defaults默认值模块) - 默认值模块说明

### 迁移文档

- [迁移文档索引](../../../migration/README.md) - 迁移文档列表和编写规范

### 代码实现

- `internal/commands/config/validate.go`（如适用） - 配置验证实现
- `internal/commands/migrate/` - 迁移实现代码
- `internal/lib/config/manager.go` - ConfigManager 实现
- `internal/lib/config/defaults.go`（如适用） - 默认值函数

---

## ✅ 检查清单

使用本规范时，请确保：

- [ ] 配置加载时已验证配置有效性
- [ ] 配置格式变化时已提供迁移机制
- [ ] 所有配置项已提供合理的默认值
- [ ] 默认值变更已评估影响并更新文档

---

**最后更新**: 2025-12-23


# 测试规范

> 本文档是 Workflow CLI 项目测试规范的索引和总览，包含所有测试相关的规范和最佳实践。

本目录用于存放 Workflow CLI 项目的测试规范文档，定义测试组织、编写、执行等各个环节的规范和最佳实践。

## 📋 说明

测试规范文档用于指导开发者编写高质量、可维护的测试代码，确保测试覆盖全面、执行高效、结果可靠。

## 📝 文档规范

### 模板使用

编写测试代码时，应参考以下模板：

**测试模板位置**：`docs/templates/testing/`

该目录提供了以下测试模板：

- **标准测试模板**（`test-case.template`）：标准测试用例的模板，包含完整的 AAA 结构
- **被忽略测试模板**（`ignored-test.template`）：被忽略测试的完整模板，包含5个必需的文档部分
- **Mock测试模板**（`mock-test.template`）：使用 Mock 的测试模板

**使用场景**：
- 编写新测试时，参考对应的模板结构
- 确保测试文档注释符合规范
- 保持测试代码风格一致

### 文档命名

测试规范文档使用描述性名称命名，例如：
- `organization.md` - 测试组织规范
- `writing.md` - 测试编写规范
- `commands.md` - 测试命令参考

### 文档位置

- 测试规范文档存放在 `docs/testing/` 目录下
- 参考文档存放在 `docs/testing/references/` 目录下

---

## 📋 目录

- [概述](#-概述)
- [核心规范](#-核心规范)
- [参考文档](#-参考文档)
- [测试模板](#-测试模板)
- [快速开始](#-快速开始)
- [相关文档](#-相关文档)

---

## 📋 概述

本文档目录包含所有测试规范文档，分为以下类别：

- **核心规范**：日常测试必须遵循的规范（测试组织、测试编写、测试命令）
- **参考文档**：详细的参考指南（工具使用、Mock测试、覆盖率、性能测试等）
- **测试模板**：常用的测试模板和代码片段

### 快速导航

| 类别 | 文档 | 说明 |
|------|------|------|
| **核心规范** | [测试组织规范](./organization.md) | 测试类型、文件结构、命名约定、共享工具 |
| | [测试编写规范](./writing.md) | AAA模式、命名规范、独立性、断言最佳实践 |
| | [测试命令参考](./commands.md) | 常用测试命令、调试命令、Makefile命令 |
| **参考文档** | [测试工具指南](./references/tools.md) | testify、go-cmp、httptest |
| | [测试环境工具指南](./references/environments.md) | 测试环境隔离、临时目录管理 |
| | [测试辅助工具指南](./references/helpers.md) | 测试辅助函数、测试数据生成 |
| | [Mock测试指南](./references/mock-server.md) | HTTP Mock、接口Mock、测试替身 |
| | [被忽略测试规范](./references/ignored-tests.md) | 文档格式、测试类型模板、最佳实践 |
| | [覆盖率测试指南](./references/coverage.md) | 覆盖率工具、报告生成、提升技巧 |
| | [性能测试指南](./references/performance.md) | 基准测试、性能要求、优化建议 |
| | [单元测试指南](./references/unit-tests.md) | 编写规范、组织方式、最佳实践 |
| | [集成测试指南](./references/integration-tests.md) | 环境配置、数据隔离、清理机制 |
| | [跨平台测试方案](./cross-platform.md) | 跨平台测试策略、CI/CD 测试 |
| **测试模板** | [测试文档模板](./documentation-template.md) | 标准化测试文档模板（推荐） |
| | [标准测试模板](../templates/testing/test-case.template) | 标准测试用例模板 |
| | [被忽略测试模板](../templates/testing/ignored-test.template) | 被忽略测试的完整模板 |
| | [Mock测试模板](../templates/testing/mock-test.template) | Mock 测试模板 |

---

## 核心规范

核心规范是日常测试必须遵循的规范，建议优先阅读：

### [测试组织规范](./organization.md)

定义测试组织结构、命名约定和共享工具使用规范。

**关键内容**：
- 测试类型（单元测试、集成测试、表驱动测试）
- 测试文件组织结构（`*_test.go` 文件）
- 测试文件命名约定
- 共享测试工具（testutils 包）
- 测试数据管理（testdata 目录）

**快速参考**：
```
internal/
├── lib/
│   ├── config/
│   │   ├── manager.go
│   │   └── manager_test.go      # 单元测试
│   └── http/
│       ├── client.go
│       └── client_test.go
testdata/                          # 测试数据
├── fixtures/
│   ├── sample_github_pr.json
│   └── sample_jira_response.json
└── integration/                   # 集成测试数据
```

### [测试编写规范](./writing.md)

定义测试编写规范和最佳实践。

**关键内容**：
- 测试结构（AAA 模式：Arrange-Act-Assert）
- 测试命名规范
- 测试独立性原则
- 测试覆盖原则
- 断言最佳实践
- 错误处理测试
- 边界条件测试
- 表驱动测试

**快速参考**：
```go
func TestParseTicketID(t *testing.T) {
    // Arrange: 准备测试数据
    input := "PROJ-123"
    expected := "PROJ-123"

    // Act: 执行被测试的功能
    result := ParseTicketID(input)

    // Assert: 验证结果
    assert.Equal(t, expected, result)
}
```

### [测试命令参考](./commands.md)

提供常用测试命令的快速参考。

**关键内容**：
- 基本测试命令
- 测试类型命令（单元、集成、基准测试）
- Makefile 测试命令
- 测试调试命令
- 常用命令速查

**快速参考**：
```bash
# 运行所有测试
go test ./...

# 运行特定包的测试
go test ./internal/lib/config

# 运行测试并显示覆盖率
go test -cover ./...

# 生成覆盖率报告
make test-coverage
```

---

## 参考文档

参考文档提供详细的参考指南，按需查阅：

### [测试工具指南](./references/tools.md)

介绍常用测试工具的使用方法。

**包含工具**：
- `testify` - 断言和Mock框架
- `go-cmp` - 深度比较工具
- `httptest` - HTTP测试工具
- 测试环境工具（临时目录、环境变量隔离）
- 测试辅助工具（测试数据生成、CLI命令测试）

### [Mock测试指南](./references/mock-server.md)

详细说明 Mock 测试的使用方法。

**关键内容**：
- HTTP Mock 基本使用（httptest）
- 接口Mock（testify/mock）
- 测试替身（Test Doubles）
- 验证Mock调用
- Mock最佳实践

### [测试辅助工具指南](./references/helpers.md)

详细介绍测试辅助工具的使用方法，包括测试辅助函数、测试数据生成和路径获取函数。

**关键内容**：
- 测试辅助函数基本使用和API参考
- 测试数据生成方法
- **路径获取函数**：`TestHomeDir()`, `TestConfigDir()`, `TestDataDir()`, `TestCacheDir()`
- 测试环境隔离和跨平台兼容性
- 与源代码行为一致性
- 最佳实践和使用示例

### [被忽略测试规范](./references/ignored-tests.md)

定义被忽略测试的文档规范。

**关键内容**：
- 统一文档格式（5个必需部分）
- 6种测试类型模板
- 文档编写最佳实践
- 文档维护清单

### [覆盖率测试指南](./references/coverage.md)

介绍测试覆盖率的检查和提升方法。

**关键内容**：
- 覆盖率工具安装和使用
- 生成覆盖率报告
- 覆盖率分析方法
- 覆盖率提升技巧

### [性能测试指南](./references/performance.md)

介绍性能测试和基准测试的方法。

**关键内容**：
- 性能测试要求
- 基准测试（`go test -bench`）
- 性能测试报告
- 性能优化建议

### [单元测试指南](./references/unit-tests.md)

介绍单元测试的编写规范、组织方式和最佳实践。

**关键内容**：
- 单元测试组织方式（文件位置、包名）
- 编写规范（AAA 模式、命名规范）
- 测试模式（表驱动测试、错误处理、边界条件）
- 常见场景和示例
- 与集成测试的区别

### [集成测试指南](./references/integration-tests.md)

介绍集成测试的环境配置和最佳实践。

**关键内容**：
- 集成测试环境配置
- 数据隔离规则
- 清理机制
- 临时文件管理

### [跨平台测试方案](./cross-platform.md)

定义跨平台测试的策略、方法和最佳实践。

**关键内容**：
- 支持的平台和目标架构
- 本地测试方法
- CI/CD 跨平台测试流程
- 平台特定测试处理
- 常见问题和解决方案
- 跨平台测试最佳实践

---

## 测试模板

项目提供了常用的测试模板，帮助快速创建标准化的测试代码：

### [测试文档模板](./documentation-template.md) ⭐ 推荐

**标准化测试文档模板**，包含：
- 标准测试模板（基础版和完整版）
- 被忽略测试模板
- 表驱动测试模板
- 集成测试模板
- 错误处理测试模板
- 实际使用示例

**使用场景**：所有新编写的测试都应该参考此模板编写文档注释。

### [标准测试模板](../templates/testing/test-case.template)

标准测试用例的模板，包含完整的 AAA 结构。

**使用场景**：编写新的单元测试时参考此模板。

### [被忽略测试模板](../templates/testing/ignored-test.template)

被忽略测试的完整模板，包含5个必需的文档部分。

**使用场景**：编写需要被忽略的测试（如集成测试、平台特定测试）时参考此模板。

### [Mock测试模板](../templates/testing/mock-test.template)

使用 Mock 的测试模板。

**使用场景**：编写需要 Mock 外部依赖的测试时参考此模板。

## 🎯 使用指南

### 创建新测试

1. 根据测试类型选择合适的模板：
   - 标准单元测试 → 使用 `test-case.template`
   - 需要 Mock 的测试 → 使用 `mock-test.template`
   - 被忽略的测试 → 使用 `ignored-test.template`

2. 参考模板结构编写测试代码和文档注释

3. 确保测试符合[测试编写规范](./writing.md)的要求

### 参考模板位置

所有测试模板位于 `docs/templates/testing/` 目录下。

---

## 快速开始

### 新手开发者学习路径

1. **阅读核心规范**（~1小时）
   - [测试组织规范](./organization.md) - 了解测试结构
   - [测试编写规范](./writing.md) - 学习编写规范
   - [测试命令参考](./commands.md) - 掌握常用命令

2. **实践编写测试**（~2小时）
   - 参考[测试文档模板](./documentation-template.md)编写测试文档
   - 使用[标准测试模板](../templates/testing/test-case.template)编写测试代码
   - 参考[测试工具指南](./references/tools.md)
   - 运行测试并查看结果

3. **深入学习**（按需）
   - 需要隔离的测试环境 → 阅读 [测试环境工具指南](./references/environments.md)
   - 需要 CLI 命令测试辅助 → 阅读 [测试辅助工具指南](./references/helpers.md)
   - 需要 Mock 外部 API → 阅读 [Mock测试指南](./references/mock-server.md)
   - 需要提升覆盖率 → 阅读 [覆盖率测试指南](./references/coverage.md)

### 快速检查清单

开始编写测试前，请确保：

- [ ] 已阅读[测试组织规范](./organization.md)和[测试编写规范](./writing.md)
- [ ] 了解[测试命令](./commands.md)的基本用法
- [ ] 测试文件放在正确的位置（与源码同目录的 `*_test.go` 文件）
- [ ] 测试命名遵循规范（`Test` 前缀，描述性名称）
- [ ] 测试使用 AAA 模式（Arrange-Act-Assert）
- [ ] 测试之间相互独立，不共享状态

### 常用命令速查

```bash
# 运行所有测试
go test ./...

# 运行特定包的测试
go test ./internal/lib/config

# 运行特定测试函数
go test -run TestParseTicketID ./internal/lib/config

# 运行测试并显示详细输出
go test -v ./...

# 运行测试并显示覆盖率
go test -cover ./...

# 生成覆盖率报告
make test-coverage

# 运行基准测试
go test -bench=. ./...

# 运行被忽略的测试
go test -run TestIgnored ./...
```

---

## ❓ 常见问题

### Q: 为什么测试需要使用隔离环境？

A: 测试隔离确保：
- **测试之间不会相互影响**：每个测试运行在独立的临时目录中
- **测试不会污染真实系统路径**：测试不会修改用户的真实配置文件
- **测试可以在并行环境中安全运行**：支持并行测试执行，提高测试效率
- **测试结果可重复**：测试环境一致，结果可重复

详细说明请参考 [测试环境工具指南](./references/environments.md)。

### Q: 什么时候应该使用 `TestHomeDir()` 而不是 `os.UserHomeDir()`？

A:
- **使用 `TestHomeDir()`**：当测试需要支持测试环境隔离时
  - 测试需要访问主目录下的配置文件
  - 测试需要支持环境变量隔离
  - 测试需要跨平台兼容性

- **使用 `os.UserHomeDir()`**：当测试需要访问真实的系统路径时
  - 需要从真实的全局配置复制
  - 需要验证真实系统路径的行为
  - 需要访问用户的实际配置文件

详细说明请参考 [测试辅助工具指南 - 路径获取函数](./references/helpers.md#3-路径获取函数)。

### Q: 如何确保测试在不同平台上都能运行？

A:
- **使用统一的路径获取函数**：`TestHomeDir()`, `TestConfigDir()` 等函数自动处理平台差异
- **使用构建标签标记平台特定代码**：明确标记平台特定的测试和代码
- **在 CI/CD 中运行跨平台测试**：确保所有平台上的测试都能通过
- **参考平台差异分析文档**：了解平台差异和注意事项
- **参考跨平台测试方案**：了解完整的跨平台测试策略和方法

相关文档：
- [跨平台测试方案](./cross-platform.md) - 完整的跨平台测试指南

### Q: 测试环境会自动清理吗？

A: 是的，测试环境使用 `t.Cleanup()` 自动清理：

- **临时目录**：使用 `t.TempDir()` 创建的临时目录会在测试结束后自动清理
- **环境变量**：使用 `t.Setenv()` 设置的环境变量会在测试结束后自动恢复
- **文件资源**：使用 `defer` 或 `t.Cleanup()` 确保资源清理

**无需手动清理**：
```go
func TestExample(t *testing.T) {
    tempDir := t.TempDir() // 自动清理
    t.Setenv("HOME", tempDir) // 自动恢复

    // 测试代码
    // 测试结束后自动清理，无需手动清理
}
```

### Q: 临时目录应该使用哪个函数？

A:
- **临时目录**：使用 `t.TempDir()` 或 `os.MkdirTemp()`
  - `t.TempDir()` 是推荐方式，自动清理
  - `os.MkdirTemp()` 需要手动清理

- **标准目录**：使用统一的路径获取函数
  - `TestHomeDir()` - 主目录
  - `TestConfigDir()` - 配置目录
  - `TestDataDir()` - 数据目录
  - `TestCacheDir()` - 缓存目录

详细说明请参考 [测试辅助工具指南 - 路径获取函数](./references/helpers.md#3-路径获取函数)。

---

## 📚 相关文档

### 测试模板

- [标准测试模板](../templates/testing/test-case.template)
- [被忽略测试模板](../templates/testing/ignored-test.template)
- [Mock测试模板](../templates/testing/mock-test.template)

### 开发规范

- [开发规范总览](../development/README.md) - 开发规范索引
- [代码风格规范](../development/code-style.md) - 代码格式化和 Lint
- [错误处理规范](../development/error-handling.md) - 错误处理最佳实践

### 测试审查

- [测试用例检查指南](../development/references/review-test-case.md) - AI 助手测试审查指南
- [测试覆盖检查机制](../development/references/test-coverage-check.md) - 测试覆盖检查流程

### 其他指南

- [文档编写指南](../guidelines/document.md) - 文档编写规范和模板

---

## ✅ 测试质量标准

> **测试覆盖率目标**：详见 [测试组织规范 - 测试覆盖率](./organization.md#-测试覆盖率)

一个高质量的测试应该满足：

1. **清晰性** - 测试名称和代码清晰表达测试意图
2. **独立性** - 测试之间相互独立，不共享状态
3. **完整性** - 覆盖成功路径、错误路径和边界条件
4. **可维护性** - 代码简洁，易于理解和修改
5. **快速性** - 单元测试 < 100ms，集成测试 < 1s
6. **稳定性** - 测试结果可重复，不受外部因素影响

---

**最后更新**: 2025-01-28

---

## 📝 变更历史

### 2025-01-28
- **重写文档**：从 Rust 风格完全重写为 Go 风格
- **更新测试工具**：从 Rust 工具（cargo, pretty_assertions, rstest）更新为 Go 工具（go test, testify, go-cmp）
- **更新测试组织**：从 Rust 的 `tests/` 目录结构更新为 Go 的 `*_test.go` 文件结构
- **更新测试命令**：从 `cargo test` 更新为 `go test`
- **新增路径获取函数文档**：在 `helpers.md` 中添加路径获取函数章节，说明 `TestHomeDir()` 等函数的使用
- **新增测试基础设施最佳实践**：在 `writing.md` 中添加测试基础设施最佳实践章节
- **新增常见问题章节**：在 `README.md` 中添加常见问题（FAQ）章节，解答测试环境、路径获取等常见问题

# 构建标签说明

## 概述

`testutils` 包使用构建标签 `//go:build test`，确保只在测试时编译，不会被打包到 release 中。

## 工作原理

### 构建标签

所有 `testutils` 包的文件都包含以下构建标签：

```go
//go:build test

package testutils
```

这意味着：
- ✅ **使用 `-tags=test` 时**：testutils 包会被编译，可以在测试中使用
- ❌ **不使用标签时**：testutils 包不会被编译，如果生产代码导入会失败

### 验证

```bash
# ✅ 正常构建（不包含 testutils）
go build ./cmd/workflow

# ✅ 使用标签构建 testutils
go build -tags=test ./internal/testutils

# ❌ 不使用标签构建 testutils（会失败）
go build ./internal/testutils
# 输出：package github.com/zevwings/workflow/internal/testutils: build constraints exclude all Go files
```

## 使用方法

### 运行测试

```bash
# 使用 Makefile（已自动包含 -tags=test）
make test
make test-coverage

# 手动运行测试
go test -tags=test ./...

# 运行特定包的测试
go test -tags=test ./internal/config
```

### 在测试代码中使用

```go
package config_test

import (
    "testing"
    "github.com/zevwings/workflow/internal/testutils"
)

func TestExample(t *testing.T) {
    // 使用 testutils
    homeDir := testutils.TestHomeDir(t)
    // ...
}
```

### 注意事项

1. **测试时必须使用 `-tags=test`**：否则 testutils 包不会被编译
2. **生产代码不能导入 testutils**：如果导入，编译会失败（这是好的，提醒开发者）
3. **Makefile 已更新**：`make test` 和 `make test-coverage` 已自动包含 `-tags=test`

## 为什么使用构建标签？

- ✅ **避免打包到 release**：testutils 是测试专用工具，不应该出现在生产代码中
- ✅ **编译时检查**：如果生产代码错误导入了 testutils，编译会失败，及时发现问题
- ✅ **清晰的职责分离**：明确区分测试代码和生产代码
- ✅ **减少二进制大小**：testutils 不会被包含在最终二进制文件中

## 相关文档

- [testutils README](./README.md) - 完整使用文档
- [Go 构建标签文档](https://pkg.go.dev/cmd/go#hdr-Build_constraints)


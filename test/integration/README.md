# test/integration 目录

本目录用于存放**跨包级别**或**端到端（E2E）**的集成测试。

## 设计原则

- **包级集成测试**（仅依赖单个包，对外部服务做真实调用）
  - 继续放在各自包目录中，例如：
    - `internal/jira/jira_client_integration_test.go`
  - 使用 `//go:build integration` 构建标签进行区分。

- **跨包 / E2E 集成测试**
  - 放在 `test/integration/` 下，例如：
    - `test/integration/workflow_integration_test.go`
  - 同样使用 `//go:build integration` 构建标签。

## 运行方式

```bash
# 仅运行跨包 / E2E 集成测试
go test -tags=integration ./test/integration/...

# 配合包级集成测试，一起运行所有集成测试
go test -tags=integration ./...
```

## 编写约定

- 测试文件命名：`*_integration_test.go`
- 包名统一使用：`package integration`
- 尽量只依赖公共 API，不直接操作各包的私有实现



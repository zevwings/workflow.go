# test/integration 目录

本目录用于存放**所有集成测试**，包括包级集成测试和跨包/端到端（E2E）集成测试。

## 设计原则

- **所有集成测试**（包括包级和跨包/E2E）
  - 统一放在 `test/integration/` 目录下
  - 使用 `//go:build integration` 构建标签进行区分
  - 包名统一使用 `package integration`
  - 通过导入包的方式使用各包的公共 API

- **包级集成测试**（仅依赖单个包，对外部服务做真实调用）
  - 例如：`test/integration/jira_client_integration_test.go`
  - 测试单个包的功能，但使用真实的外部服务（如 Jira API）

- **跨包 / E2E 集成测试**
  - 例如：`test/integration/workflow_integration_test.go`
  - 测试多个包的交互和完整工作流程

## 运行方式

```bash
# 运行所有集成测试（包括包级和跨包/E2E）
go test -tags=integration ./test/integration/...

# 运行特定集成测试
go test -tags=integration ./test/integration/... -run TestJiraClient
```

## 编写约定

- 测试文件命名：`*_integration_test.go`
- 包名统一使用：`package integration`
- 尽量只依赖公共 API，不直接操作各包的私有实现



# 快速启动指南

## 前置要求

- Go 1.21 或更高版本
- Git（用于版本控制）

## 安装步骤

### 1. 初始化项目

```bash
# 下载依赖
go mod download
go mod tidy
```

### 2. 构建项目

```bash
# 使用 Makefile
make build

# 或直接使用 go build
go build -o bin/workflow ./cmd/workflow
```

### 3. 运行基础命令

```bash
# 查看帮助
./bin/workflow --help

# 查看版本
./bin/workflow version

# 初始化配置
./bin/workflow setup

# 查看配置
./bin/workflow config show

# 验证配置
./bin/workflow config validate

# 检查环境
./bin/workflow check
```

## 开发模式

### 运行测试

```bash
make test
```

### 格式化代码

```bash
make fmt
```

### 代码检查

```bash
# 需要先安装 golangci-lint
# go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
make lint
```

### 直接运行（不构建）

```bash
make run
# 或
go run ./cmd/workflow --help
```

## 项目结构说明

```
workflow.go/
├── cmd/workflow/          # 主入口
│   └── main.go
├── internal/
│   ├── cli/               # CLI 根命令
│   │   └── root.go
│   ├── commands/          # 命令实现
│   │   ├── setup.go       # 配置初始化
│   │   ├── version.go     # 版本信息
│   │   ├── config.go      # 配置管理
│   │   └── check.go       # 环境检查
│   ├── lib/               # 核心业务逻辑
│   │   ├── config/        # 配置管理
│   │   │   └── manager.go
│   │   └── http/          # HTTP 客户端
│   │       └── client.go
│   └── utils/             # 工具函数
│       ├── output.go      # 输出格式化
│       └── table.go       # 表格显示
├── go.mod
├── go.sum
├── Makefile
└── README.md
```

## 下一步

1. 运行 `workflow setup` 初始化配置
2. 运行 `workflow check` 检查环境
3. 查看 `go_migration_plan.md` 了解完整的迁移计划
4. 开始实现其他命令模块

## 常见问题

### Q: 构建失败，提示找不到依赖？

A: 运行 `go mod download` 和 `go mod tidy` 下载依赖。

### Q: 配置文件在哪里？

A: 配置文件默认在 `~/.workflow/config.toml`。

### Q: 如何添加新命令？

A: 在 `internal/commands/` 目录下创建新文件，然后在 `internal/cli/root.go` 中注册命令。

## 贡献

欢迎提交 Issue 和 Pull Request！


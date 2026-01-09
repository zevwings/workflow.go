.PHONY: build install clean test run help

# 变量定义
BINARY_NAME=workflow
VERSION?=0.1.0
BUILD_DATE=$(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
GIT_COMMIT=$(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
LDFLAGS=-ldflags "-X github.com/zevwings/workflow/cmd/workflow.version=$(VERSION) -X github.com/zevwings/workflow/cmd/workflow.buildDate=$(BUILD_DATE) -X github.com/zevwings/workflow/cmd/workflow.gitCommit=$(GIT_COMMIT) -s -w"

# 默认目标
.DEFAULT_GOAL := help

# 构建二进制文件（不包含示例代码）
build:
	@echo "构建 $(BINARY_NAME)..."
	@go build $(LDFLAGS) -o bin/$(BINARY_NAME) ./cmd/workflow
	@echo "构建完成: bin/$(BINARY_NAME)"

# 安装到系统
install: build
	@echo "安装 $(BINARY_NAME)..."
	@sudo cp bin/$(BINARY_NAME) /usr/local/bin/
	@echo "安装完成"

# 运行程序（不包含示例代码）
run:
	@go run $(LDFLAGS) ./cmd/workflow

# 运行示例程序
run-example:
	@go run -tags=example $(LDFLAGS) ./cmd/example

# 清理构建文件
clean:
	@echo "清理构建文件..."
	@rm -rf bin/
	@go clean
	@echo "清理完成"

# 运行测试
test:
	@echo "运行测试..."
	@go test -tags=test -v ./...

# 运行测试并生成覆盖率报告
test-coverage:
	@echo "运行测试并生成覆盖率报告..."
	@go test -tags=test -v -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html
	@echo "覆盖率报告已生成: coverage.html"

# 运行跨包 / E2E 集成测试（test/integration）
test-integration:
	@echo "运行跨包 / E2E 集成测试..."
	@go test -tags=integration -v ./test/integration/...

# 格式化代码
fmt:
	@echo "格式化代码..."
	@go fmt ./...
	@gofmt -s -w .

# 代码检查
lint:
	@echo "代码检查..."
	@golangci-lint run

# 下载依赖
deps:
	@echo "下载依赖..."
	@go mod download
	@go mod tidy

# 更新依赖
update-deps:
	@echo "更新依赖..."
	@go get -u ./...
	@go mod tidy

# 生成 Shell 补全脚本
completion:
	@echo "生成 Shell 补全脚本..."
	@mkdir -p completions
	@go run ./cmd/workflow completion bash > completions/workflow.bash
	@go run ./cmd/workflow completion zsh > completions/workflow.zsh
	@go run ./cmd/workflow completion fish > completions/workflow.fish
	@go run ./cmd/workflow completion powershell > completions/workflow.ps1
	@echo "补全脚本已生成到 completions/ 目录"

# 显示帮助信息
help:
	@echo "可用的 Make 目标:"
	@echo "  build          - 构建二进制文件（不包含示例）"
	@echo "  build-example  - 构建示例程序（独立的 workflow-example）"
	@echo "  install        - 安装到系统 (/usr/local/bin)"
	@echo "  run            - 运行程序（不包含示例）"
	@echo "  run-example    - 运行示例程序"
	@echo "  clean          - 清理构建文件"
	@echo "  test           - 运行测试"
	@echo "  test-coverage  - 运行测试并生成覆盖率报告"
	@echo "  fmt            - 格式化代码"
	@echo "  lint           - 代码检查"
	@echo "  deps           - 下载依赖"
	@echo "  update-deps    - 更新依赖"
	@echo "  completion     - 生成 Shell 补全脚本"
	@echo "  help           - 显示此帮助信息"


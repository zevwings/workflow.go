# make/Makefile.build.mk
# 构建和安装模块

# 变量定义
BINARY_NAME = workflow
VERSION ?= 0.1.0
BUILD_DATE = $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
GIT_COMMIT = $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
LDFLAGS = -ldflags "-X github.com/zevwings/workflow/cmd/workflow.version=$(VERSION) -X github.com/zevwings/workflow/cmd/workflow.buildDate=$(BUILD_DATE) -X github.com/zevwings/workflow/cmd/workflow.gitCommit=$(GIT_COMMIT) -s -w"

# Help 信息
define HELP_BUILD
	@echo "构建相关："
	@echo "  make build          - 构建二进制文件（不包含示例）"
	@echo "  make run            - 运行程序（不包含示例）"
	@echo "  make run-example    - 运行示例程序"
	@echo "  make clean          - 清理构建文件"
	@echo ""
	@echo "安装相关："
	@echo "  make install        - 安装到系统 (/usr/local/bin)"
	@echo ""
endef

# ============================================
# 构建相关目标
# ============================================

# 构建二进制文件（不包含示例代码）
build:
	@echo "构建 $(BINARY_NAME)..."
	@go build $(LDFLAGS) -o bin/$(BINARY_NAME) ./cmd/workflow
	@echo "构建完成: bin/$(BINARY_NAME)"

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

# ============================================
# 安装和部署相关目标
# ============================================

# 安装到系统
install: build
	@echo "安装 $(BINARY_NAME)..."
	@sudo cp bin/$(BINARY_NAME) /usr/local/bin/
	@echo "安装完成"

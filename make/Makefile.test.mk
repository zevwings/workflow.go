# make/Makefile.test.mk
# 测试模块

# Help 信息
define HELP_TEST
	@echo "测试相关："
	@echo "  make test                         - 运行测试"
	@echo "  make test-integration              - 运行跨包 / E2E 集成测试"
	@echo ""
	@echo "覆盖率相关："
	@echo "  make test-coverage                - 运行测试并生成覆盖率报告（标准 HTML）"
	@echo "  make test-coverage-ui             - 运行测试并生成覆盖率报告（标准 HTML，别名）"
	@echo "  make test-coverage-interactive    - 使用 gocovsh 交互式查看覆盖率（终端 UI，类似 cargo-tarpaulin，推荐）"
	@echo "  make test-coverage-treemap        - 生成覆盖率树状图可视化"
	@echo "  make open-coverage                - 打开标准覆盖率 HTML 报告"
	@echo "  make open-coverage-ui             - 打开覆盖率 HTML 报告"
	@echo ""
endef

# ============================================
# 测试运行
# ============================================

# 运行测试
test:
	@echo "运行测试..."
	@go test -tags=test -v ./...

# 运行跨包 / E2E 集成测试（test/integration）
test-integration:
	@echo "运行跨包 / E2E 集成测试..."
	@go test -tags=integration -v ./test/integration/...

# ============================================
# 覆盖率报告
# ============================================

# 运行测试并生成覆盖率报告（标准 HTML）
test-coverage:
	@echo "运行测试并生成覆盖率报告..."
	@mkdir -p coverage
	@go test -tags=test -v -coverprofile=coverage/coverage.out ./...
	@go tool cover -html=coverage/coverage.out -o coverage/coverage.html
	@echo "覆盖率报告已生成: coverage/coverage.html"

# 运行测试并生成美观的覆盖率报告（使用标准工具，兼容性更好）
test-coverage-ui:
	@echo "运行测试并生成覆盖率报告..."
	@mkdir -p coverage
	@go test -tags=test -v -coverprofile=coverage/coverage.out ./...
	@go tool cover -html=coverage/coverage.out -o coverage/coverage-ui.html
	@echo "覆盖率报告已生成: coverage/coverage-ui.html"
	@echo "使用 'make open-coverage-ui' 打开报告"
	@echo ""
	@echo "提示: 如需更美观的报告，可以使用:"
	@echo "  - make test-coverage-interactive  (终端交互式 UI，类似 cargo-tarpaulin)"
	@echo "  - make test-coverage-treemap      (树状图可视化)"

# 使用 gocovsh 交互式查看覆盖率（终端 UI，类似 cargo-tarpaulin，推荐）
test-coverage-interactive:
	@echo "运行测试并使用 gocovsh 查看覆盖率..."
	@if ! command -v gocovsh > /dev/null; then \
		echo "安装 gocovsh..."; \
		if ! go install github.com/orlangure/gocovsh@latest 2>/dev/null; then \
			echo "错误: 无法安装 gocovsh，请检查网络连接或 Go 环境"; \
			echo "提示: 可以手动安装: go install github.com/orlangure/gocovsh@latest"; \
			exit 1; \
		fi; \
	fi
	@mkdir -p coverage
	@go test -tags=test -v -coverprofile=coverage/coverage.out ./...
	@echo ""
	@echo "启动 gocovsh 交互式界面..."
	@echo "提示: 使用 j/k 键导航，Enter 选择，Esc 退出，q 退出"
	@gocovsh --profile=coverage/coverage.out

# 生成覆盖率树状图可视化（使用 go-cover-treemap）
test-coverage-treemap:
	@echo "运行测试并生成覆盖率树状图..."
	@if ! command -v go-cover-treemap > /dev/null; then \
		echo "安装 go-cover-treemap..."; \
		go install github.com/nikolaydubina/go-cover-treemap@latest; \
	fi
	@mkdir -p coverage
	@go test -tags=test -v -coverprofile=coverage/coverage.out ./...
	@go-cover-treemap -coverprofile=coverage/coverage.out > coverage/coverage-treemap.svg
	@echo "覆盖率树状图已生成: coverage/coverage-treemap.svg"
	@if command -v open > /dev/null; then \
		open coverage/coverage-treemap.svg; \
	elif command -v xdg-open > /dev/null; then \
		xdg-open coverage/coverage-treemap.svg; \
	elif command -v start > /dev/null; then \
		start coverage/coverage-treemap.svg; \
	fi

# 打开覆盖率 HTML 报告（标准）
open-coverage:
	@if [ ! -f coverage/coverage.html ]; then \
		echo "错误: 覆盖率报告不存在，请先运行 'make test-coverage'"; \
		exit 1; \
	fi
	@echo "打开覆盖率报告..."
	@if command -v open > /dev/null; then \
		open coverage/coverage.html; \
	elif command -v xdg-open > /dev/null; then \
		xdg-open coverage/coverage.html; \
	elif command -v start > /dev/null; then \
		start coverage/coverage.html; \
	else \
		echo "无法自动打开，请手动打开: coverage/coverage.html"; \
	fi

# 打开美观的覆盖率 HTML 报告
open-coverage-ui:
	@if [ ! -f coverage/coverage-ui.html ]; then \
		echo "错误: 美观的覆盖率报告不存在，请先运行 'make test-coverage-ui'"; \
		exit 1; \
	fi
	@echo "打开美观的覆盖率报告..."
	@if command -v open > /dev/null; then \
		open coverage/coverage-ui.html; \
	elif command -v xdg-open > /dev/null; then \
		xdg-open coverage/coverage-ui.html; \
	elif command -v start > /dev/null; then \
		start coverage/coverage-ui.html; \
	else \
		echo "无法自动打开，请手动打开: coverage/coverage-ui.html"; \
	fi

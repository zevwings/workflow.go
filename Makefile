# 主 Makefile
# 包含所有功能模块（按依赖顺序）

# 声明所有伪目标（统一管理）
.PHONY: help build install clean run run-example test test-integration test-coverage test-coverage-ui test-coverage-interactive test-coverage-treemap open-coverage open-coverage-ui fmt lint download-deps update-deps completion

# 设置默认目标
.DEFAULT_GOAL := help

# 包含功能模块（按依赖顺序）
include make/Makefile.build.mk       # 1. 构建和安装（命令: build, install, run, run-example, clean | 变量: BINARY_NAME, VERSION, BUILD_DATE, GIT_COMMIT, LDFLAGS）
include make/Makefile.lint.mk        # 2. 代码检查（命令: fmt, lint）
include make/Makefile.test.mk        # 3. 测试（命令: test, test-integration, test-coverage, test-coverage-ui, test-coverage-interactive, test-coverage-treemap, open-coverage, open-coverage-ui）
include make/Makefile.tools.mk       # 4. 工具（命令: completion）
include make/Makefile.deps.mk        # 5. 依赖管理（命令: download-deps, update-deps）

# 集成所有模块的 help 信息
help:
	@echo "可用的 Make 目标："
	@echo ""
	$(HELP_BUILD)
	$(HELP_LINT)
	$(HELP_TEST)
	$(HELP_TOOLS)
	$(HELP_DEPS)

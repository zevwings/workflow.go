# make/Makefile.deps.mk
# 依赖管理模块

# Help 信息
define HELP_DEPS
	@echo "依赖管理相关："
	@echo "  make download-deps   - 下载依赖"
	@echo "  make update-deps     - 更新依赖"
	@echo ""
endef

# ============================================
# 依赖管理
# ============================================

# 下载依赖
download-deps:
	@echo "下载依赖..."
	@go mod download
	@go mod tidy
	@echo "✓ 依赖下载完成"

# 更新依赖
update-deps:
	@echo "更新依赖..."
	@go get -u ./...
	@go mod tidy
	@echo "✓ 依赖更新完成"

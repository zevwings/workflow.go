# make/Makefile.tools.mk
# 工具模块

# Help 信息
define HELP_TOOLS
	@echo "工具相关："
	@echo "  make completion     - 生成 Shell 补全脚本"
	@echo ""
endef

# ============================================
# Shell 补全
# ============================================

# 生成 Shell 补全脚本
completion:
	@echo "生成 Shell 补全脚本..."
	@mkdir -p completions
	@go run ./cmd/workflow completion bash > completions/workflow.bash
	@go run ./cmd/workflow completion zsh > completions/workflow.zsh
	@go run ./cmd/workflow completion fish > completions/workflow.fish
	@go run ./cmd/workflow completion powershell > completions/workflow.ps1
	@echo "补全脚本已生成到 completions/ 目录"

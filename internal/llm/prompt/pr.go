package prompt

import (
	"fmt"

	"github.com/zevwings/workflow/internal/config"
	"github.com/zevwings/workflow/internal/llm"
)

// RewordPRSystemPrompt PR Reword 的 system prompt
//
// 用于根据当前 PR 标题和 PR diff 生成简洁的 PR 标题和描述。
// 从嵌入的模板文件中加载。
var RewordPRSystemPrompt = MustLoadTemplate("pr-reword.md")

// GenerateSummarizePRSystemPrompt 根据语言生成 PR 总结的 system prompt
//
// 参数:
//   - cfg: 配置管理器（用于获取语言设置）
//
// 返回:
//   - string: 根据语言定制的 system prompt
//
// 说明:
//
//	语言选择优先级：配置文件 > 默认值（"en"）
//	如果配置文件中的语言代码不在支持列表中，将使用英文作为默认语言。
func GenerateSummarizePRSystemPrompt(cfg *config.GlobalManager) string {
	// 获取 JSON 响应示例（动态内容）
	summarizeResponseExample := "{\n" +
		"      \\\"summary\\\": \\\"# Add User Authentication\\\\\\\\n\\\\\\\\n## Overview\\\\\\\\nThis PR adds user authentication functionality to the application.\\\\\\\\n\\\\\\\\n## Requirements Analysis\\\\\\\\n\\\\\\\\n### Business Requirements\\\\\\\\nDevelopers need a secure way to authenticate users...\\\\\\\\n\\\\\\\\n### Functional Requirements\\\\\\\\nThe system accepts user credentials and returns authentication tokens...\\\\\\\\n\\\\\\\\n## Key Changes\\\\\\\\n- Added login endpoint\\\\\\\\n- Implemented JWT token generation\\\\\\\\n\\\\\\\\n## Files Changed\\\\\\\\n- src/auth/login.ts: Added login handler\\\\\\\\n- src/auth/jwt.ts: Added token generation\\\\\\\\n\\\\\\\\n## Technical Details\\\\\\\\nImplemented JWT-based authentication:\\\\\\\\n\\\\\\\\n[code block: typescript]\\\\\\\\nfunction generateToken(user: User): string {\\\\\\\\n  return jwt.sign({ userId: user.id }, secret);\\\\\\\\n}\\\\n[code block end]\\\\\\\\n\\\\\\\\n## Testing\\\\\\\\nAdded unit tests for authentication flow.\\\\\\\\n\\\\\\\\n## Usage Instructions\\\\\\\\nRun npm run test to execute tests.\\\",\\n" +
		"      \\\"filename\\\": \\\"add-user-authentication\\\"\\n" +
		"    }"

	// 从嵌入的模板文件中加载基础 prompt，然后拼接动态内容
	basePrompt := MustLoadTemplate("pr-summary.md")
	fullPrompt := fmt.Sprintf(basePrompt, summarizeResponseExample)

	// 使用 LLM 模块的语言增强功能
	return llm.GetLanguageRequirement(fullPrompt, cfg)
}

package pr

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/zevwings/workflow/internal/llm"
	"github.com/zevwings/workflow/internal/llm/prompt"
)

// GeneratePRContent 同时生成分支名、PR 标题、描述和 scope（通过一个 LLM 请求）
//
// 根据 commit 标题和 git diff 生成符合规范的分支名、PR 标题、描述和 scope。
// 分支名和 PR 标题都会自动翻译为英文（如果输入包含非英文内容）。
//
// 参数:
//   - commitTitle: commit 标题或描述
//   - existsBranches: 已存在的分支列表（可选）
//   - gitDiff: Git 工作区和暂存区的修改内容（可选，用于生成描述和提取 scope）
//   - client: LLM 客户端实例
//
// 返回:
//   - *llm.PullRequestContent: PR 内容，包含分支名、PR 标题、描述和 scope
//   - error: 如果 LLM API 调用失败或响应格式不正确，返回相应的错误信息
func GeneratePRContent(commitTitle string, existsBranches []string, gitDiff string, client *llm.LLMClient) (*llm.PullRequestContent, error) {
	// 构建请求参数
	userPrompt := buildCreateUserPrompt(commitTitle, existsBranches, gitDiff)
	systemPrompt := prompt.GenerateBranchSystemPrompt

	params := &llm.LLMRequestParams{
		SystemPrompt: systemPrompt,
		UserPrompt:   userPrompt,
		MaxTokens:    nil, // 不限制，确保有足够空间返回完整的 JSON（包括 description）
		Temperature:  0.5,
	}

	// 调用 LLM API
	response, err := client.Call(params)
	if err != nil {
		return nil, fmt.Errorf("调用 LLM API 生成分支名失败 (commit title: '%s'): %w", commitTitle, err)
	}

	// 解析响应
	content, err := parseCreateResponse(response, commitTitle)
	if err != nil {
		return nil, fmt.Errorf("解析 LLM 响应失败 (commit title: '%s'): %w", commitTitle, err)
	}

	return content, nil
}

// buildCreateUserPrompt 生成同时生成分支名和 PR 标题的 user prompt
func buildCreateUserPrompt(commitTitle string, existsBranches []string, gitDiff string) string {
	// 提取分支列表，如果没有或为空则使用空数组
	baseBranchNames := existsBranches
	if len(baseBranchNames) == 0 {
		baseBranchNames = []string{}
	}

	// 组装 prompt 内容
	// 明确说明优先级：commit title 是主要输入，git diff 仅用于验证
	parts := []string{
		fmt.Sprintf("Commit title (PRIMARY INPUT): %s", commitTitle),
		"",
		"Instructions:",
		"- Generate PR title primarily based on the commit title above",
		"- Use git changes below only to verify and refine, not to replace the commit title's intent",
		"- Focus on the business intent expressed in the commit title, not implementation details",
	}

	if len(baseBranchNames) > 0 {
		parts = append(parts, "")
		parts = append(parts, fmt.Sprintf("Existing base branch names: %s", strings.Join(baseBranchNames, ", ")))
	}

	if gitDiff != "" && strings.TrimSpace(gitDiff) != "" {
		// 限制 git diff 长度，避免超过 LLM token 限制
		// create 主要用于生成标题和描述，不需要完整 diff
		const maxDiffLength = 10000 // create 只需要了解主要变更
		diffTrimmed := TruncateDiff(gitDiff, maxDiffLength)
		parts = append(parts, "")
		parts = append(parts, "Git changes (for verification only):")
		parts = append(parts, diffTrimmed)
	}

	return strings.Join(parts, "\n")
}

// parseCreateResponse 解析 LLM 返回的 JSON 响应，提取分支名、PR 标题、描述和 scope
//
// 从 LLM 的 JSON 响应中提取 `branch_name`、`pr_title`、`description` 和 `scope` 字段。
// 支持处理包含 markdown 代码块的响应格式。
func parseCreateResponse(response, commitTitle string) (*llm.PullRequestContent, error) {
	// 使用公共方法提取并修复 JSON（修复转义问题）
	jsonStr := ExtractAndFixJSON(response)

	// 解析 JSON
	var jsonData map[string]interface{}
	if err := json.Unmarshal([]byte(jsonStr), &jsonData); err != nil {
		return nil, fmt.Errorf("解析 LLM 响应为 JSON 失败。原始响应: %s: %w", jsonStr, err)
	}

	// 提取 branch_name
	branchNameRaw, ok := jsonData["branch_name"]
	if !ok {
		return nil, fmt.Errorf("LLM 响应中缺少 'branch_name' 字段")
	}
	branchName, ok := branchNameRaw.(string)
	if !ok {
		return nil, fmt.Errorf("LLM 响应中 'branch_name' 字段类型错误")
	}

	// 提取 pr_title
	prTitleRaw, ok := jsonData["pr_title"]
	if !ok {
		return nil, fmt.Errorf("LLM 响应中缺少 'pr_title' 字段")
	}
	prTitle, ok := prTitleRaw.(string)
	if !ok {
		return nil, fmt.Errorf("LLM 响应中 'pr_title' 字段类型错误")
	}

	// description 是可选的
	var description *string
	if descRaw, ok := jsonData["description"]; ok {
		if descStr, ok := descRaw.(string); ok {
			descTrimmed := strings.TrimSpace(descStr)
			if descTrimmed != "" {
				description = &descTrimmed
			}
		}
	}

	// scope 是可选的
	var scope *string
	if scopeRaw, ok := jsonData["scope"]; ok {
		if scopeStr, ok := scopeRaw.(string); ok {
			scopeTrimmed := strings.TrimSpace(scopeStr)
			if scopeTrimmed != "" {
				scope = &scopeTrimmed
			}
		}
	}

	// 清理分支名，确保只保留 ASCII 字符
	cleanedBranchName := sanitizeBranchName(strings.TrimSpace(branchName))

	return &llm.PullRequestContent{
		BranchName:  cleanedBranchName,
		PRTitle:     strings.TrimSpace(prTitle),
		Description: description,
		Scope:       scope,
	}, nil
}

// sanitizeBranchName 清理分支名，确保只保留 ASCII 字符
func sanitizeBranchName(name string) string {
	var result strings.Builder
	for _, r := range name {
		// 只保留字母、数字、连字符和下划线
		if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') || r == '-' || r == '_' {
			result.WriteRune(r)
		}
	}
	return result.String()
}

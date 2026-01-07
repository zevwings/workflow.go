package pr

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/zevwings/workflow/internal/llm"
	"github.com/zevwings/workflow/internal/llm/prompt"
)

// RewordPR 基于当前 PR 标题和 PR diff 生成更新的 PR 标题和描述
//
// 根据当前 PR 标题和 PR diff 内容生成更新的标题和完整的描述，用于更新现有 PR。
// 与 `GeneratePRContent` 流程保持一致：标题主要基于当前标题，PR diff 用于验证和细化。
// 与 `SummarizePR()` 不同，这个方法生成的是适合作为 PR 元数据的标题和描述列表，
// 而不是详细的总结文档。
//
// 参数:
//   - prDiff: PR 的 diff 内容（用于验证和细化标题）
//   - currentTitle: 当前 PR 标题（主要输入，如果包含 markdown 格式如 `#` 会保留）
//   - client: LLM 客户端实例
//
// 返回:
//   - *llm.PullRequestReword: PR Reword 结果，包含标题和描述
//   - error: 如果 LLM API 调用失败或响应格式不正确，返回相应的错误信息
func RewordPR(prDiff string, currentTitle *string, client *llm.LLMClient) (*llm.PullRequestReword, error) {
	// 构建请求参数
	userPrompt := buildRewordUserPrompt(prDiff, currentTitle)
	systemPrompt := prompt.RewordPRSystemPrompt

	params := &llm.LLMRequestParams{
		SystemPrompt: systemPrompt,
		UserPrompt:   userPrompt,
		MaxTokens:    nil, // 增加 token 限制以支持更完整的描述
		Temperature:  0.5,
	}

	// 调用 LLM API
	response, err := client.Call(params)
	if err != nil {
		titleStr := "nil"
		if currentTitle != nil {
			titleStr = *currentTitle
		}
		return nil, fmt.Errorf("调用 LLM API 重写 PR 失败 (current title: '%s'): %w", titleStr, err)
	}

	// 解析响应
	reword, err := parseRewordResponse(response, currentTitle)
	if err != nil {
		titleStr := "nil"
		if currentTitle != nil {
			titleStr = *currentTitle
		}
		return nil, fmt.Errorf("解析 LLM 响应失败 (current title: '%s'): %w", titleStr, err)
	}

	return reword, nil
}

// buildRewordUserPrompt 生成 PR reword 的 user prompt
//
// 与 create 流程保持一致：当前标题作为主要输入，PR diff 用于验证和细化。
func buildRewordUserPrompt(prDiff string, currentTitle *string) string {
	parts := []string{}

	// 如果有当前标题，将其作为主要输入（与 create 流程一致）
	if currentTitle != nil && *currentTitle != "" {
		parts = append(parts, fmt.Sprintf("Current PR title (PRIMARY INPUT): %s", *currentTitle))
		parts = append(parts, "")
		parts = append(parts, "Instructions:")
		parts = append(parts, "- Generate PR title primarily based on the current PR title above")
		parts = append(parts, "- Use PR diff below only to verify and refine, not to replace the current title's intent")
		parts = append(parts, "- Focus on the business intent expressed in the current title, not implementation details")
		parts = append(parts, "")
	} else {
		// 如果没有当前标题，回退到基于 PR diff 生成
		parts = append(parts, "Instructions:")
		parts = append(parts, "- Generate a PR title and description based on the PR diff below")
		parts = append(parts, "")
	}

	if prDiff != "" && strings.TrimSpace(prDiff) != "" {
		// 限制 diff 长度，避免超过 LLM token 限制
		// reword 只需要了解主要变更，不需要完整 diff
		const maxDiffLength = 12000 // reword 需要比 summary 少一些上下文
		diffTrimmed := TruncateDiff(prDiff, maxDiffLength)
		parts = append(parts, "PR Diff (for verification only):")
		parts = append(parts, diffTrimmed)
	}

	return strings.Join(parts, "\n")
}

// parseRewordResponse 解析 LLM 返回的 JSON 响应，提取 PR 标题和描述
//
// 从 LLM 的 JSON 响应中提取 `pr_title` 和 `description` 字段。
// 支持处理包含 markdown 代码块的响应格式。
func parseRewordResponse(response string, currentTitle *string) (*llm.PullRequestReword, error) {
	// 使用公共方法提取并修复 JSON（修复转义问题）
	jsonStr := ExtractAndFixJSON(response)

	// 解析 JSON
	var jsonData map[string]interface{}
	if err := json.Unmarshal([]byte(jsonStr), &jsonData); err != nil {
		return nil, fmt.Errorf("解析 LLM 响应为 JSON 失败。原始响应: %s: %w", jsonStr, err)
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

	return &llm.PullRequestReword{
		PRTitle:     strings.TrimSpace(prTitle),
		Description: description,
	}, nil
}

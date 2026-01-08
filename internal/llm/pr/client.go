package pr

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/zevwings/workflow/internal/llm/client"
	"github.com/zevwings/workflow/internal/llm/prompt"
	"github.com/zevwings/workflow/internal/llm/utils"
)

// PullRequestLLMClient PR LLM 客户端
//
// 封装所有 PR 相关的 LLM 操作，包括生成 PR 内容、总结 PR、重写 PR 和总结文件变更。
// 提供统一的接口和配置管理。
type PullRequestLLMClient struct {
	llmClient *client.LLMClient
	lang      *client.SupportedLanguage
}

// NewPullRequestLLMClient 创建新的 PR LLM 客户端
//
// 参数:
//   - llmClient: LLM 客户端实例（不能为 nil）
//   - lang: 语言配置（如果为 nil，使用默认英文配置）
//
// 返回:
//   - *PullRequestLLMClient: PR LLM 客户端实例
func NewPullRequestLLMClient(llmClient *client.LLMClient, lang *client.SupportedLanguage) *PullRequestLLMClient {
	if llmClient == nil {
		panic("pr.NewPullRequestLLMClient: llmClient 不能为 nil")
	}
	return &PullRequestLLMClient{
		llmClient: llmClient,
		lang:      lang,
	}
}

// GenerateContent 生成 PR 内容（分支名、标题、描述和 scope）
//
// 根据 commit 标题和 git diff 生成符合规范的分支名、PR 标题、描述和 scope。
// 分支名和 PR 标题都会自动翻译为英文（如果输入包含非英文内容）。
//
// 参数:
//   - commitTitle: commit 标题或描述
//   - existsBranches: 已存在的分支列表（可选）
//   - gitDiff: Git 工作区和暂存区的修改内容（可选，用于生成描述和提取 scope）
//
// 返回:
//   - *PullRequestContent: PR 内容，包含分支名、PR 标题、描述和 scope
//   - error: 如果 LLM API 调用失败或响应格式不正确，返回相应的错误信息
func (c *PullRequestLLMClient) GenerateContent(commitTitle string, existsBranches []string, gitDiff string) (*PullRequestContent, error) {
	return GeneratePRContent(commitTitle, existsBranches, gitDiff, c.llmClient)
}

// Summarize 生成 PR 总结文档和文件名
//
// 根据 PR 的 diff 内容生成总结文档和合适的文件名。
//
// 参数:
//   - prTitle: PR 标题
//   - prDiff: PR 的 diff 内容
//
// 返回:
//   - *PullRequestSummary: PR 总结结果，包含总结文档和文件名
//   - error: 如果 LLM API 调用失败或响应格式不正确，返回相应的错误信息
func (c *PullRequestLLMClient) Summarize(prTitle, prDiff string) (*PullRequestSummary, error) {
	return SummarizePR(prTitle, prDiff, c.lang, c.llmClient)
}

// Reword 重写 PR 标题和描述
//
// 根据当前 PR 标题和 PR diff 内容生成更新的标题和完整的描述，用于更新现有 PR。
// 与 `GenerateContent` 流程保持一致：标题主要基于当前标题，PR diff 用于验证和细化。
// 与 `Summarize()` 不同，这个方法生成的是适合作为 PR 元数据的标题和描述列表，
// 而不是详细的总结文档。
//
// 参数:
//   - prDiff: PR 的 diff 内容（用于验证和细化标题）
//   - currentTitle: 当前 PR 标题（主要输入，如果包含 markdown 格式如 `#` 会保留）
//
// 返回:
//   - *PullRequestReword: PR Reword 结果，包含标题和描述
//   - error: 如果 LLM API 调用失败或响应格式不正确，返回相应的错误信息
func (c *PullRequestLLMClient) Reword(prDiff string, currentTitle *string) (*PullRequestReword, error) {
	return RewordPR(prDiff, currentTitle, c.llmClient)
}

// SummarizeFileChange 总结单个文件变更
//
// 根据文件的 diff 内容生成该文件的修改总结。
//
// 参数:
//   - filePath: 文件路径
//   - fileDiff: 文件的 diff 内容
//
// 返回:
//   - string: 文件的修改总结（纯文本）
//   - error: 如果 LLM API 调用失败，返回相应的错误信息
func (c *PullRequestLLMClient) SummarizeFileChange(filePath, fileDiff string) (string, error) {
	return SummarizeFileChange(filePath, fileDiff, c.lang, c.llmClient)
}

// ============================================================================
// GeneratePRContent 相关函数
// ============================================================================

// GeneratePRContent 同时生成分支名、PR 标题、描述和 scope（通过一个 LLM 请求）
//
// 根据 commit 标题和 git diff 生成符合规范的分支名、PR 标题、描述和 scope。
// 分支名和 PR 标题都会自动翻译为英文（如果输入包含非英文内容）。
//
// 参数:
//   - commitTitle: commit 标题或描述
//   - existsBranches: 已存在的分支列表（可选）
//   - gitDiff: Git 工作区和暂存区的修改内容（可选，用于生成描述和提取 scope）
//   - llmClient: LLM 客户端实例
//
// 返回:
//   - *PullRequestContent: PR 内容，包含分支名、PR 标题、描述和 scope
//   - error: 如果 LLM API 调用失败或响应格式不正确，返回相应的错误信息
func GeneratePRContent(commitTitle string, existsBranches []string, gitDiff string, llmClient *client.LLMClient) (*PullRequestContent, error) {
	// 构建请求参数
	userPrompt := buildCreateUserPrompt(commitTitle, existsBranches, gitDiff)
	systemPrompt := prompt.GenerateBranchSystemPrompt

	params := &client.LLMRequestParams{
		SystemPrompt: systemPrompt,
		UserPrompt:   userPrompt,
		MaxTokens:    nil, // 不限制，确保有足够空间返回完整的 JSON（包括 description）
		Temperature:  0.5,
	}

	// 调用 LLM API
	response, err := llmClient.Call(params)
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
		parts = append(parts, "")
		parts = append(parts, "Git changes (for verification only):")
		parts = append(parts, gitDiff)
	}

	return strings.Join(parts, "\n")
}

// parseCreateResponse 解析 LLM 返回的 JSON 响应，提取分支名、PR 标题、描述和 scope
//
// 从 LLM 的 JSON 响应中提取 `branch_name`、`pr_title`、`description` 和 `scope` 字段。
// 支持处理包含 markdown 代码块的响应格式。
func parseCreateResponse(response, commitTitle string) (*PullRequestContent, error) {
	// 使用公共方法提取并修复 JSON（修复转义问题）
	jsonStr := utils.ExtractAndFixJSON(response)

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
	cleanedBranchName := utils.SanitizeBranchName(strings.TrimSpace(branchName))

	return &PullRequestContent{
		BranchName:  cleanedBranchName,
		PRTitle:     strings.TrimSpace(prTitle),
		Description: description,
		Scope:       scope,
	}, nil
}

// ============================================================================
// SummarizePR 相关函数
// ============================================================================

// SummarizePR 生成 PR 总结文档和文件名
//
// 根据 PR 的 diff 内容生成总结文档和合适的文件名。
//
// 参数:
//   - prTitle: PR 标题
//   - prDiff: PR 的 diff 内容
//   - lang: 语言配置（如果为 nil，使用默认英文配置）
//   - llmClient: LLM 客户端实例
//
// 返回:
//   - *PullRequestSummary: PR 总结结果，包含总结文档和文件名
//   - error: 如果 LLM API 调用失败或响应格式不正确，返回相应的错误信息
func SummarizePR(prTitle, prDiff string, lang *client.SupportedLanguage, llmClient *client.LLMClient) (*PullRequestSummary, error) {
	// 构建请求参数
	userPrompt := buildSummaryUserPrompt(prTitle, prDiff)
	// 根据语言生成 system prompt
	systemPrompt := prompt.GenerateSummarizePRSystemPrompt(lang)

	params := &client.LLMRequestParams{
		SystemPrompt: systemPrompt,
		UserPrompt:   userPrompt,
		MaxTokens:    nil, // 增加 token 数量，确保有足够空间返回完整的总结文档
		Temperature:  0.3, // 降低温度，使输出更稳定
	}

	// 调用 LLM API
	response, err := llmClient.Call(params)
	if err != nil {
		return nil, fmt.Errorf("调用 LLM API 总结 PR 失败 (PR title: '%s'): %w", prTitle, err)
	}

	// 解析响应
	summary, err := parseSummaryResponse(response, prTitle)
	if err != nil {
		return nil, fmt.Errorf("解析 LLM 响应失败 (PR title: '%s'): %w", prTitle, err)
	}

	return summary, nil
}

// buildSummaryUserPrompt 生成 PR 总结的 user prompt
func buildSummaryUserPrompt(prTitle, prDiff string) string {
	parts := []string{fmt.Sprintf("PR Title: %s", prTitle)}

	if prDiff != "" && strings.TrimSpace(prDiff) != "" {
		parts = append(parts, fmt.Sprintf("PR Diff:\n%s", prDiff))
	}

	return strings.Join(parts, "\n\n")
}

// parseSummaryResponse 解析 LLM 返回的 JSON 响应，提取总结文档和文件名
//
// 从 LLM 的 JSON 响应中提取 `summary` 和 `filename` 字段。
// 支持处理包含 markdown 代码块的响应格式。
func parseSummaryResponse(response, prTitle string) (*PullRequestSummary, error) {
	// 使用公共方法提取并修复 JSON（修复转义问题）
	jsonStr := utils.ExtractAndFixJSON(response)

	// 解析 JSON
	var jsonData map[string]interface{}
	if err := json.Unmarshal([]byte(jsonStr), &jsonData); err != nil {
		return nil, fmt.Errorf("解析 LLM 响应为 JSON 失败。原始响应: %s: %w", jsonStr, err)
	}

	// 提取 summary
	summaryRaw, ok := jsonData["summary"]
	if !ok {
		return nil, fmt.Errorf("LLM 响应中缺少 'summary' 字段")
	}
	summary, ok := summaryRaw.(string)
	if !ok {
		return nil, fmt.Errorf("LLM 响应中 'summary' 字段类型错误")
	}

	// 提取 filename
	filenameRaw, ok := jsonData["filename"]
	if !ok {
		return nil, fmt.Errorf("LLM 响应中缺少 'filename' 字段")
	}
	filename, ok := filenameRaw.(string)
	if !ok {
		return nil, fmt.Errorf("LLM 响应中 'filename' 字段类型错误")
	}

	// 清理文件名，确保只包含有效的文件名字符
	cleanedFilename := utils.CleanFilename(filename)

	if cleanedFilename == "" {
		return nil, fmt.Errorf("清理后的文件名为空")
	}

	return &PullRequestSummary{
		Summary:  strings.TrimSpace(summary),
		Filename: cleanedFilename,
	}, nil
}

// ============================================================================
// RewordPR 相关函数
// ============================================================================

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
//   - llmClient: LLM 客户端实例
//
// 返回:
//   - *PullRequestReword: PR Reword 结果，包含标题和描述
//   - error: 如果 LLM API 调用失败或响应格式不正确，返回相应的错误信息
func RewordPR(prDiff string, currentTitle *string, llmClient *client.LLMClient) (*PullRequestReword, error) {
	// 构建请求参数
	userPrompt := buildRewordUserPrompt(prDiff, currentTitle)
	systemPrompt := prompt.RewordPRSystemPrompt

	params := &client.LLMRequestParams{
		SystemPrompt: systemPrompt,
		UserPrompt:   userPrompt,
		MaxTokens:    nil, // 增加 token 限制以支持更完整的描述
		Temperature:  0.5,
	}

	// 调用 LLM API
	response, err := llmClient.Call(params)
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
		parts = append(parts, "PR Diff (for verification only):")
		parts = append(parts, prDiff)
	}

	return strings.Join(parts, "\n")
}

// parseRewordResponse 解析 LLM 返回的 JSON 响应，提取 PR 标题和描述
//
// 从 LLM 的 JSON 响应中提取 `pr_title` 和 `description` 字段。
// 支持处理包含 markdown 代码块的响应格式。
func parseRewordResponse(response string, currentTitle *string) (*PullRequestReword, error) {
	// 使用公共方法提取并修复 JSON（修复转义问题）
	jsonStr := utils.ExtractAndFixJSON(response)

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

	return &PullRequestReword{
		PRTitle:     strings.TrimSpace(prTitle),
		Description: description,
	}, nil
}

// ============================================================================
// SummarizeFileChange 相关函数
// ============================================================================

// SummarizeFileChange 生成单个文件的修改总结
//
// 根据文件的 diff 内容生成该文件的修改总结。
//
// 参数:
//   - filePath: 文件路径
//   - fileDiff: 文件的 diff 内容
//   - lang: 语言配置（如果为 nil，使用默认英文配置）
//   - llmClient: LLM 客户端实例
//
// 返回:
//   - string: 文件的修改总结（纯文本）
//   - error: 如果 LLM API 调用失败，返回相应的错误信息
func SummarizeFileChange(filePath, fileDiff string, lang *client.SupportedLanguage, llmClient *client.LLMClient) (string, error) {
	// 构建请求参数
	userPrompt := buildFileSummaryUserPrompt(filePath, fileDiff)
	// 根据语言生成 system prompt
	systemPrompt := prompt.GenerateSummarizeFileChangeSystemPrompt(lang)

	params := &client.LLMRequestParams{
		SystemPrompt: systemPrompt,
		UserPrompt:   userPrompt,
		MaxTokens:    nil, // 单个文件的总结应该比较简短
		Temperature:  0.3,
	}

	// 调用 LLM API
	response, err := llmClient.Call(params)
	if err != nil {
		return "", fmt.Errorf("调用 LLM API 总结文件修改失败 (file path: '%s'): %w", filePath, err)
	}

	// 清理响应（移除可能的 markdown 代码块包装）
	summary := cleanFileChangeSummaryResponse(response)

	return summary, nil
}

// buildFileSummaryUserPrompt 生成单个文件修改总结的 user prompt
func buildFileSummaryUserPrompt(filePath, fileDiff string) string {
	return fmt.Sprintf("File path: %s\n\nFile diff:\n%s", filePath, fileDiff)
}

// cleanFileChangeSummaryResponse 清理文件修改总结响应
//
// 移除可能的 markdown 代码块包装，返回纯文本。
func cleanFileChangeSummaryResponse(response string) string {
	return utils.ExtractJSONFromMarkdown(response)
}

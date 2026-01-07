package pr

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/zevwings/workflow/internal/config"
	"github.com/zevwings/workflow/internal/llm"
	"github.com/zevwings/workflow/internal/llm/prompt"
)

// SummarizePR 生成 PR 总结文档和文件名
//
// 根据 PR 的 diff 内容生成总结文档和合适的文件名。
//
// 参数:
//   - prTitle: PR 标题
//   - prDiff: PR 的 diff 内容
//   - cfgMgr: 配置管理器（用于语言设置）
//   - client: LLM 客户端实例
//
// 返回:
//   - *llm.PullRequestSummary: PR 总结结果，包含总结文档和文件名
//   - error: 如果 LLM API 调用失败或响应格式不正确，返回相应的错误信息
func SummarizePR(prTitle, prDiff string, cfgMgr *config.GlobalManager, client *llm.LLMClient) (*llm.PullRequestSummary, error) {
	// 构建请求参数
	userPrompt := buildSummaryUserPrompt(prTitle, prDiff)
	// 根据语言生成 system prompt（语言选择逻辑在 prompt 生成函数内部处理）
	systemPrompt := prompt.GenerateSummarizePRSystemPrompt(cfgMgr)

	params := &llm.LLMRequestParams{
		SystemPrompt: systemPrompt,
		UserPrompt:   userPrompt,
		MaxTokens:    nil, // 增加 token 数量，确保有足够空间返回完整的总结文档
		Temperature:  0.3, // 降低温度，使输出更稳定
	}

	// 调用 LLM API
	response, err := client.Call(params)
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
		// 限制 diff 长度，避免请求过大
		// 对于总结，我们需要更多的 diff 内容，但也要避免超过 token 限制
		const maxDiffLength = 15000 // 增加长度，因为总结需要更多上下文
		diffTrimmed := TruncateDiff(prDiff, maxDiffLength)
		parts = append(parts, fmt.Sprintf("PR Diff:\n%s", diffTrimmed))
	}

	return strings.Join(parts, "\n\n")
}

// parseSummaryResponse 解析 LLM 返回的 JSON 响应，提取总结文档和文件名
//
// 从 LLM 的 JSON 响应中提取 `summary` 和 `filename` 字段。
// 支持处理包含 markdown 代码块的响应格式。
func parseSummaryResponse(response, prTitle string) (*llm.PullRequestSummary, error) {
	// 使用公共方法提取并修复 JSON（修复转义问题）
	jsonStr := ExtractAndFixJSON(response)

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
	cleanedFilename := cleanFilename(filename)

	if cleanedFilename == "" {
		return nil, fmt.Errorf("清理后的文件名为空")
	}

	return &llm.PullRequestSummary{
		Summary:  strings.TrimSpace(summary),
		Filename: cleanedFilename,
	}, nil
}

// cleanFilename 清理文件名，确保只包含有效的文件名字符
func cleanFilename(filename string) string {
	// 转小写
	cleaned := strings.ToLower(strings.TrimSpace(filename))
	// 替换空格为连字符
	cleaned = strings.ReplaceAll(cleaned, " ", "-")

	// 只保留字母、数字、连字符和下划线
	var result strings.Builder
	for _, r := range cleaned {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') || r == '-' || r == '_' {
			result.WriteRune(r)
		}
	}

	cleaned = result.String()

	// 移除 .md 扩展名（如果存在），因为我们会自动添加
	cleaned = strings.TrimSuffix(cleaned, ".md")

	return cleaned
}

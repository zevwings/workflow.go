package llm

// LLMRequestParams LLM 请求参数
//
// 包含调用 LLM API 所需的所有参数。
type LLMRequestParams struct {
	// SystemPrompt 系统提示词
	SystemPrompt string `json:"system_prompt"`
	// UserPrompt 用户提示词
	UserPrompt string `json:"user_prompt"`
	// MaxTokens 最大 token 数（nil 表示不限制，使用模型默认最大值）
	MaxTokens *int `json:"max_tokens,omitempty"`
	// Temperature 温度参数（控制输出的随机性，范围 0.0-2.0）
	Temperature float32 `json:"temperature"`
	// Model 模型名称（可选，如果为空则从配置获取）
	Model string `json:"model,omitempty"`
}

// DefaultLLMRequestParams 返回默认的 LLM 请求参数
//
// 返回:
//   - *LLMRequestParams: 使用默认值的请求参数
func DefaultLLMRequestParams() *LLMRequestParams {
	temperature := float32(0.5)
	return &LLMRequestParams{
		SystemPrompt: "",
		UserPrompt:   "",
		MaxTokens:    nil,
		Temperature:  temperature,
		Model:        "",
	}
}

// ChatCompletionResponse OpenAI Chat Completions API 响应
//
// 完整的 OpenAI 标准响应格式，支持所有标准字段和扩展字段。
type ChatCompletionResponse struct {
	// ID 响应唯一标识符
	ID string `json:"id"`
	// Object 对象类型，固定为 "chat.completion"
	Object string `json:"object"`
	// Created 创建时间戳（Unix 时间戳）
	Created int64 `json:"created"`
	// Model 使用的模型名称
	Model string `json:"model"`
	// SystemFingerprint 系统指纹（可选）
	SystemFingerprint *string `json:"system_fingerprint,omitempty"`
	// Choices 选择列表
	Choices []ChatCompletionChoice `json:"choices"`
	// Usage Token 使用统计
	Usage Usage `json:"usage"`
}

// ChatCompletionChoice Chat Completion 选择项
type ChatCompletionChoice struct {
	// Index 选择索引
	Index int `json:"index"`
	// Message 消息对象
	Message ChatMessage `json:"message"`
	// Logprobs 对数概率（可选）
	Logprobs interface{} `json:"logprobs,omitempty"`
	// FinishReason 完成原因
	FinishReason string `json:"finish_reason"`
}

// ChatMessage Chat 消息
type ChatMessage struct {
	// Role 消息角色
	Role string `json:"role"`
	// Content 消息内容（可能为 null）
	Content *string `json:"content,omitempty"`
}

// Usage Token 使用统计
type Usage struct {
	// PromptTokens 提示词 token 数
	PromptTokens int `json:"prompt_tokens"`
	// CompletionTokens 完成 token 数
	CompletionTokens int `json:"completion_tokens"`
	// TotalTokens 总 token 数
	TotalTokens int `json:"total_tokens"`
}

// PullRequestContent PR 内容，包含分支名、PR 标题、描述和 scope
//
// 由 LLM 生成的分支名、PR 标题、描述和 scope，用于创建 Pull Request。
type PullRequestContent struct {
	// BranchName 分支名称（小写，使用连字符分隔）
	BranchName string
	// PRTitle PR 标题（简洁，不超过 8 个单词）
	PRTitle string
	// Description PR 描述（基于 Git 修改内容生成，可选）
	Description *string
	// Scope Commit scope（从 git diff 提取，用于 Conventional Commits 格式，可选）
	//
	// Scope 表示变更涉及的模块或功能区域，例如 "api", "auth", "jira" 等。
	// 如果无法确定 scope，此字段为 nil。
	Scope *string
}

// PullRequestReword PR Reword 结果，包含标题和描述
//
// 由 LLM 基于当前 PR 标题和 PR diff 生成的标题和完整描述，用于更新现有 PR。
type PullRequestReword struct {
	// PRTitle PR 标题（简洁，不超过 8 个单词，主要基于当前标题，如果当前标题包含 markdown 格式如 `#` 会保留）
	PRTitle string
	// Description PR 描述（基于 PR diff 生成的完整描述列表，包含所有重要变更，可选）
	Description *string
}

// PullRequestSummary PR 总结结果，包含总结文档和文件名
//
// 由 LLM 生成的 PR 总结文档和对应的文件名。
type PullRequestSummary struct {
	// Summary PR 总结文档（Markdown 格式）
	Summary string
	// Filename 文件名（不含路径和扩展名）
	Filename string
}

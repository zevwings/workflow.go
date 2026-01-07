package pr

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

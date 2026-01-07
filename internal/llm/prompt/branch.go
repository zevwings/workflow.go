package prompt

// GenerateBranchSystemPrompt 生成分支名的 system prompt
//
// 用于根据 commit 标题和 git 变更生成分支名、PR 标题和描述。
// 从嵌入的模板文件中加载。
var GenerateBranchSystemPrompt = MustLoadTemplate("branch.md")

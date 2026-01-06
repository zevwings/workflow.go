# 测试用例审查指南

> 本文档专为 AI 助手设计，提供测试用例检查的核心原则和实用指导，帮助 AI 高效地进行测试覆盖分析和测试用例生成。

## 📋 目录

- [核心原则](#-核心原则)
- [检查目标](#-检查目标)
- [测试边界和范围](#-测试边界和范围)
- [检查流程](#-检查流程)
- [项目测试现状](#-项目测试现状)
- [测试示例](#-测试示例)
- [快速检查脚本](#-快速检查脚本)
- [检查报告模板](#-检查报告模板)

---

## 🎯 核心原则

**测试边界**：测试我们自己的业务逻辑，不测试外部依赖和第三方库。

**检查重点**：
- ✅ 业务逻辑、数据转换、状态管理、错误处理
- ✅ CLI 参数解析、用户交互、输出格式化
- ✅ API 封装、重试机制、响应处理
- ❌ 外部工具功能、第三方库实现、远程 API 业务逻辑

---

## 🎯 检查目标

测试用例检查的主要目标：

1. **确保测试覆盖完整**：所有业务逻辑、CLI 命令、错误处理都有对应的测试用例
2. **确保测试边界正确**：测试我们自己的代码，不测试外部依赖和第三方库
3. **确保测试质量**：测试用例合理、有效，使用合适的测试工具和方法
4. **识别缺失测试**：发现未测试的功能和边界情况
5. **优化测试结构**：确保测试组织清晰，易于维护

### 检查范围

- **单元测试**：`#[cfg(test)]` 模块中的测试
- **集成测试**：`test/` 目录中的测试文件或 `*_test.go` 文件
- **文档测试**：文档中的代码示例（doctest）
- **测试工具使用**：是否使用推荐的测试工具（rstest、pretty_assertions、mockito、insta 等）

### 检查原则

1. **测试边界原则**：测试我们自己的业务逻辑，不测试外部依赖和第三方库
2. **全面性原则**：检查所有模块的测试覆盖情况
3. **质量原则**：确保测试用例合理、有效，使用合适的测试工具
4. **可操作性原则**：检查结果应提供明确的改进建议

---

## 🎯 测试边界和范围

> **核心原则**：我们应该测试自己的业务逻辑，而不是测试外部依赖和第三方库的实现。

### ✅ 需要测试的内容

#### 1. 业务逻辑层

测试我们自己实现的业务规则和处理逻辑：

- ✅ **业务规则实现**
  - 分支前缀处理（如 `feature/` + 分支名）
  - 合并策略选择（如根据配置选择 `--no-ff` 或 `--ff-only`）
  - 数据验证规则（如分支名称验证、配置验证）
  - 业务流程控制（如 PR 创建流程、日志工作流）

- ✅ **数据转换和处理**
  - JSON/TOML 数据解析后的转换
  - 数据格式化（如日期格式、标题生成）
  - 数据聚合和计算（如统计、汇总）

- ✅ **状态管理**
  - 状态转换逻辑（如 Git 仓库状态、工作树状态）
  - 状态解析（如 `WorktreeStatus` 结构体的正确性）
  - 状态验证（如检查是否有未提交的更改）

- ✅ **数据结构的正确性**
  - 自定义数据结构（如 `CommitInfo`、`WorktreeStatus`、`BranchInfo`）
  - 数据结构的序列化/反序列化
  - 数据结构的默认值和验证

#### 2. CLI 命令层测试

测试命令行接口的正确性和用户体验：

- ✅ **命令参数解析**
  - 参数验证逻辑（必需参数、可选参数、默认值）
  - 参数类型转换和格式验证（如 Jira ID 格式、PR 标题长度）
  - 参数组合的有效性检查（互斥参数、依赖参数）
  - 错误参数的处理和友好提示

- ✅ **命令执行流程**
  - 命令的主要执行路径和业务逻辑
  - 前置条件检查（Git 仓库、配置文件、网络连接）
  - 命令间的依赖关系处理和执行顺序
  - 命令执行的状态管理和错误恢复

- ✅ **用户交互测试**
  - Dialog 组件的配置和验证逻辑
  - 用户输入的处理和验证（输入框、选择框、多选框）
  - 交互流程的正确性（确认、取消、重试）
  - 用户体验的一致性（提示信息、错误处理）

- ✅ **输出格式化**
  - 多种输出格式的正确性（Table、JSON、YAML、Markdown）
  - 输出内容的一致性和完整性
  - 错误消息和警告信息的格式化
  - 国际化和本地化支持

**CLI 测试示例**：
```go
// ✅ 测试命令参数解析（我们的业务逻辑）
func TestPRCreateArgsValidation(t *testing.T) {
	args := &PRCreateArgs{
		JiraTicket: "PROJ-123",
		Title:      "Test PR",
		Description: "",
		DryRun:     false,
	}

	// 测试我们的参数验证逻辑
	err := ValidatePRCreateArgs(args)
	assert.NoError(t, err)

	// 测试无效的 Jira ID
	invalidArgs := &PRCreateArgs{
		JiraTicket: "invalid-id",
		Title:      args.Title,
	}
	err = ValidatePRCreateArgs(invalidArgs)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Jira ID 格式无效")
}

// ✅ 测试用户交互逻辑（我们的业务逻辑）
func TestBranchSelectionDialogConfig(t *testing.T) {
	branches := []string{"main", "develop", "feature/test"}
	dialog := CreateBranchSelectionDialog(branches)

	// 测试我们的 Dialog 配置逻辑
	assert.Equal(t, "选择目标分支:", dialog.Prompt())
	assert.Equal(t, 3, len(dialog.Options()))
	assert.Equal(t, 0, dialog.DefaultIndex())
	assert.True(t, dialog.EnableFilter()) // 启用模糊匹配
}

// ✅ 测试输出格式化（我们的业务逻辑）
func TestPRListOutputFormats(t *testing.T) {
	prs := CreateMockPRList()

	// 测试表格格式
	tableOutput := FormatPRListAsTable(prs)
	assert.Contains(t, tableOutput, "ID")
	assert.Contains(t, tableOutput, "Title")
	assert.Contains(t, tableOutput, "Status")

	// 测试 JSON 格式
	jsonOutput := FormatPRListAsJSON(prs)
	var parsed []map[string]interface{}
	err := json.Unmarshal([]byte(jsonOutput), &parsed)
	assert.NoError(t, err)
	assert.Equal(t, len(prs), len(parsed))
}

// ✅ 测试命令执行前置条件（我们的业务逻辑）
func TestPRCreatePreconditions(t *testing.T) {
	// 测试 Git 仓库检查
	tempDir := t.TempDir()
	os.Chdir(tempDir)

	err := PrCreateCommand.ValidateGitRepo()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "不是 Git 仓库")

	// 测试配置文件检查
	err = PrCreateCommand.ValidateGitHubConfig()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "GitHub 配置")
}
```

#### 3. 错误处理逻辑

测试我们如何处理错误，而不是测试错误是否会发生：

- ✅ **异常情况处理**
  - 检查错误情况是否被正确捕获
  - 检查错误恢复机制是否正确
  - 检查失败后的清理逻辑

- ✅ **错误消息和上下文**
  - 检查错误消息是否清晰准确
  - 检查错误上下文是否完整（使用 `fmt.Errorf` 和 `%w`）
  - 检查错误类型是否正确传递

- ✅ **错误传播**
  - 检查错误是否正确向上传播
  - 检查错误转换是否正确（如将底层错误转换为业务错误）

#### 3. 边界条件

测试各种边界情况和特殊输入：

- ✅ **输入边界**
  - 空输入（空字符串、空数组、空配置）
  - 最大/最小值（长度限制、数值范围）
  - 特殊字符（Unicode、换行符、特殊符号）
  - 无效输入（格式错误、类型错误）

- ✅ **输出边界**
  - 空结果处理
  - 超大结果处理
  - 格式化输出的正确性

- ✅ **并发和异步场景**
  - 并发函数的正确性测试（使用 goroutine 和 channel）
  - 并发执行的安全性测试（数据竞争、死锁检测）
  - 并发限制和资源管理（如并发执行器的限制）
  - 超时和取消机制测试（任务超时、用户取消）
  - 异步错误处理和传播
  - 并发场景下的状态一致性

#### 4. 集成逻辑

测试我们如何封装和使用外部依赖：

- ✅ **API 调用封装**
  - 参数构造和传递
  - 请求配置（headers、auth、timeout）
  - 响应处理和数据提取

- ✅ **数据解析和转换**
  - API 响应的解析
  - 数据映射到内部数据结构
  - 错误响应的处理

- ✅ **重试和容错机制**
  - 重试逻辑的正确性（指数退避、最大重试次数）
  - 可重试错误的判断
  - 重试失败的处理

### ❌ 不需要测试的内容

#### 1. 外部依赖和第三方库

**核心原则**：不要测试外部工具和库的实现，它们已经有自己的测试。

#### 判断依据

使用以下问题判断是否需要测试：

1. **这段代码是谁写的？**
   - ❌ 外部库/工具的作者 → 不需要测试
   - ✅ 我们自己的团队 → 需要测试

2. **这段代码在哪里维护？**
   - ❌ 在外部仓库（如 crates.io、GitHub、系统工具） → 不需要测试
   - ✅ 在我们的项目中 → 需要测试

3. **测试的目的是什么？**
   - ❌ 验证外部库是否按文档工作 → 不需要测试（信任外部库）
   - ✅ 验证我们的代码逻辑是否正确 → 需要测试

#### 不需要测试的典型场景

- ❌ **外部命令行工具的功能**
  - 例如：Git、Docker、npm 等命令的基本功能
  - 不要测试：命令本身是否正确执行
  - 应该测试：我们如何构建命令参数、如何解析命令输出

- ❌ **第三方库的 API 实现**
  - 例如：HTTP 客户端、数据库驱动、序列化库
  - 不要测试：库的内部实现和协议处理
  - 应该测试：我们如何配置和使用这些库

- ❌ **远程 API 服务的业务逻辑**
  - 例如：GitHub API、Jira API、云服务 API
  - 不要测试：API 是否返回正确的业务数据
  - 应该测试：我们如何调用 API、如何处理 API 响应

- ❌ **标准库和系统调用**
  - 例如：文件系统、进程管理、网络操作
  - 不要测试：标准库的正确性
  - 应该测试：我们的文件处理逻辑、错误处理

- ❌ **语言和框架的核心功能**
  - 例如：Go 标准库、语言特性、编译器行为
  - 不要测试：语言本身的功能
  - 应该测试：我们使用这些功能的业务逻辑

#### 2. 测试策略

对于外部依赖，采用以下策略：

- ✅ **使用 Mock 和 Stub 隔离测试**
  - 使用 `httptest` Mock HTTP API
  - 使用测试工具模拟 Git 仓库状态
  - 使用 Stub 模拟外部依赖的返回值

- ✅ **测试我们的代码如何使用外部依赖**
  - 测试我们传递给外部依赖的参数是否正确
  - 测试我们如何处理外部依赖的返回值
  - 测试我们如何处理外部依赖的错误

- ✅ **测试边界和异常情况**
  - 测试外部依赖返回错误时的处理
  - 测试外部依赖返回异常数据时的处理
  - 测试外部依赖超时或不可用时的处理

### 📚 具体示例

#### 示例 0: 项目实际结构对照

**实际项目模块结构**：

##### Core 业务模块 (`internal/lib/`)
- ✅ **Base 模块** (`lib/`): HTTP 客户端、LLM 客户端、Settings、Dialog、Logger 等基础组件
- ✅ **Git 模块** (`lib/git/`): 分支管理、提交管理、仓库操作、Stash 管理
- ✅ **Jira 模块** (`lib/jira/`): API 集成、日志管理、附件处理、状态管理
- ✅ **PR 模块** (`lib/pr/`): GitHub 平台集成、LLM 生成、Body 解析
- ✅ **Branch 模块** (`lib/branch/`): 分支命名、LLM 生成、同步管理
- ✅ **Commit 模块** (`lib/commit/`): 提交修改、重写、压缩
- ✅ **Template 模块** (`lib/template/`): 模板配置、引擎、变量管理
- ✅ **Proxy 模块** (`lib/proxy/`): 代理配置生成和管理
- ✅ **Repo 模块** (`lib/repo/`): 仓库配置管理
- ✅ **Rollback 模块** (`lib/rollback/`): 操作回滚

##### CLI 命令层 (`internal/commands/`)
- ✅ **配置管理**: `config/`, `github/`, `check/`, `proxy/`, `llm/`
- ✅ **业务功能**: `pr/`, `jira/`, `branch/`, `commit/`, `stash/`, `log/`
- ✅ **工具管理**: `lifecycle/`, `migrate/`, `repo/`, `alias/`, `tag/`

**实际测试工具配置** (`go.mod`):
```go
require (
    github.com/stretchr/testify v1.8.4    // 清晰的断言输出
    github.com/google/go-cmp v0.6.0      // 深度比较工具
    github.com/gorilla/mux v1.8.1        // HTTP API Mock 测试
    github.com/spf13/cobra v1.8.0        // CLI 命令测试
)
```

**实际测试目录结构**:
```
internal/
├── lib/
│   ├── config/
│   │   └── manager_test.go      # 配置管理测试
│   └── http/
│       └── client_test.go        # HTTP 客户端测试
test/
├── cli/            # CLI 命令测试（所有 commands/ 对应测试）
│   ├── basic_cli_test.go        # 基础 CLI 测试
│   ├── integration_cli_test.go  # CLI 集成测试
│   └── [各命令测试文件]        # PR、Branch、Config 等
├── integration/    # 集成测试
└── fixtures/       # 测试数据文件
```

**测试覆盖现状**:
- 🟢 **已完整覆盖**: Base 模块（LLM、Settings、Dialog）、CLI 参数解析、PR 模块、Jira 模块
- 🟢 **CLI 测试工具**: 已添加 testify、go-cmp 和完整的测试辅助工具
- 🟡 **部分覆盖**: HTTP 模块、Completion 模块、Proxy 模块、CLI 集成测试（基础框架已建立）
- 🔴 **缺失覆盖**: Git 模块、Template 模块、Branch 模块、Commit 模块、Stash 模块

#### 示例 1: Git 模块

**✅ 应该测试的**：

```go
// ✅ 测试分支前缀处理逻辑（我们的业务逻辑）
func TestFormatBranchNameWithPrefix(t *testing.T) {
	result := FormatBranchName("feature", "login")
	assert.Equal(t, "feature/login", result)
}

// ✅ 测试合并策略选择逻辑（我们的业务逻辑）
func TestMergeStrategySelection(t *testing.T) {
	strategy := DetermineMergeStrategy(true, false)
	assert.Equal(t, MergeStrategyNoFastForward, strategy)
}

// ✅ 测试分支名称验证逻辑（我们的业务逻辑）
func TestValidateBranchName(t *testing.T) {
	err := ValidateBranchName("feature/login")
	assert.NoError(t, err)

	err = ValidateBranchName("invalid//name")
	assert.Error(t, err)

	err = ValidateBranchName("")
	assert.Error(t, err)
}

// ✅ 测试 Git 命令执行失败时的错误处理（我们的错误处理）
func TestBranchCreateErrorHandling(t *testing.T) {
	// 使用 Mock 模拟 Git 命令失败
	result, err := GitBranch.Create("invalid/name")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "分支名称无效")
	assert.Nil(t, result)
}

// ✅ 测试 CommitInfo 数据结构解析（我们的数据处理）
func TestParseCommitInfo(t *testing.T) {
	output := "abc123\nJohn Doe\n2024-01-01\nInitial commit"
	info, err := CommitInfoFromOutput(output)
	assert.NoError(t, err)
    assert.Equal(t, "abc123", info.Hash)
    assert.Equal(t, "John Doe", info.Author)
}
```

**❌ 不应该测试的**：

```go
// ❌ 不要测试 Git 命令本身是否正确（这是 Git 的责任）
func TestGitBranchCommandCreatesBranch(t *testing.T) {
	// 这是在测试 Git 本身，而不是我们的代码
	cmd := exec.Command("git", "branch", "test")
	cmd.Run()

	output, _ := exec.Command("git", "branch", "--list", "test").Output()
	assert.Contains(t, string(output), "test")
}

// ❌ 不要测试 Git 参数的功能（这是 Git 的责任）
func TestGitMergeFFOnlyParameter(t *testing.T) {
	// 这是在测试 Git 的 --ff-only 参数，而不是我们的代码
	cmd := exec.Command("git", "merge", "--ff-only", "feature")
	cmd.Run()
}

// ❌ 不要测试 Git 的底层实现（这是 Git 的责任）
func TestGitInternalMergeAlgorithm(t *testing.T) {
	// 这是在测试 Git 的合并算法，而不是我们的代码
}
```

#### 示例 2: HTTP 模块

**✅ 应该测试的**：

```go
// ✅ 测试请求配置构建逻辑（我们的业务逻辑）
func TestBuildRequestWithAuth(t *testing.T) {
	client := NewHTTPClient()
	request := client.Request("https://api.example.com").
		WithAuth("token", "abc123").
		Build()

	assert.Contains(t, request.Header.Get("Authorization"), "Bearer abc123")
}

// ✅ 测试重试逻辑（我们的业务逻辑）
func TestRetryOnNetworkError(t *testing.T) {
	// 使用 httptest 模拟服务器
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	client := NewHTTPClient()
	_, err := client.Get(server.URL + "/api")
	assert.Error(t, err)
	// 验证重试逻辑
}

// ✅ 测试响应数据解析（我们的业务逻辑）
func TestParseAPIResponse(t *testing.T) {
	json := `{"id": 123, "name": "test"}`
	var data APIResponse
	err := json.Unmarshal([]byte(json), &data)
	assert.NoError(t, err)
	assert.Equal(t, 123, data.ID)
	assert.Equal(t, "test", data.Name)
}
```

**❌ 不应该测试的**：

```go
// ❌ 不要测试 net/http 是否正确发送 HTTP 请求（这是标准库的责任）
func TestHTTPSendsHTTPRequest(t *testing.T) {
	// 这是在测试标准库，而不是我们的代码
	resp, err := http.Get("https://httpbin.org/get")
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

// ❌ 不要测试 HTTP 协议的正确性（这是标准协议）
func TestHTTPProtocol(t *testing.T) {
	// 这是在测试 HTTP 协议，而不是我们的代码
}
```

#### 示例 3: Jira 模块

**✅ 应该测试的**：

```go
// ✅ 测试 Jira API 请求构建（我们的业务逻辑）
func TestBuildJiraSearchRequest(t *testing.T) {
	query := BuildJiraQuery("PROJECT-123")
	assert.Equal(t, "project = PROJECT AND key = PROJECT-123", query)
}

// ✅ 测试 Jira 响应数据转换（我们的业务逻辑）
func TestConvertJiraIssueToInternalFormat(t *testing.T) {
	jiraIssue := MockJiraIssue()
	issue := IssueFromJiraResponse(jiraIssue)
	assert.Equal(t, "PROJECT-123", issue.Key)
	assert.Equal(t, "Test Issue", issue.Summary)
}

// ✅ 测试日志格式化（我们的业务逻辑）
func TestFormatWorklog(t *testing.T) {
	worklog := Worklog{
		TimeSpent: 3600,
		Comment:   "Fixed bug",
	}
	formatted := FormatWorklog(worklog)
	assert.Equal(t, "1h - Fixed bug", formatted)
}
```

**❌ 不应该测试的**：

```go
// ❌ 不要测试 Jira API 本身的功能（这是 Jira 的责任）
func TestJiraAPIReturnsCorrectIssue(t *testing.T) {
	// 这是在测试 Jira API，而不是我们的代码
	issue, err := jiraClient.GetIssue("PROJECT-123")
	assert.NoError(t, err)
	assert.Equal(t, "Expected Summary", issue.Fields.Summary)
}

// ❌ 不要测试 Jira 的业务逻辑（这是 Jira 的责任）
func TestJiraCalculatesTimeTracking(t *testing.T) {
	// 这是在测试 Jira 的时间跟踪逻辑，而不是我们的代码
}
```

#### 示例 4: 并发和异步测试

**✅ 应该测试的**：

```go
// ✅ 测试并发函数的正确性（我们的业务逻辑）
func TestConcurrentHTTPRequests(t *testing.T) {
	client := NewHTTPClient()
	urls := []string{"url1", "url2", "url3"}

	// 测试我们的并发请求逻辑
	results := client.FetchAll(urls)
	assert.Equal(t, 3, len(results))
	for _, r := range results {
		assert.NoError(t, r.Error)
	}
}

// ✅ 测试并发执行器的限制（我们的业务逻辑）
func TestConcurrentExecutorLimits(t *testing.T) {
    executor := NewConcurrentExecutor(2) // 最大2个并发
    tasks := createTestTasks(5) // 5个任务

    // 测试我们的并发控制逻辑
    startTime := time.Now()
    results, err := executor.ExecuteAll(tasks)
    duration := time.Since(startTime)

    // 验证结果和并发限制
    assert.NoError(t, err)
    assert.Equal(t, 5, len(results))
    assert.GreaterOrEqual(t, duration, 500*time.Millisecond) // 至少需要3轮执行
}

// ✅ 测试超时和取消机制（我们的业务逻辑）
func TestTaskTimeoutHandling(t *testing.T) {
    executor := NewConcurrentExecutor(1)
    timeoutTask := createLongRunningTask(10 * time.Second)

    // 测试我们的超时处理逻辑
    ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
    defer cancel()

    _, err := executor.ExecuteWithContext(ctx, timeoutTask)
    assert.Error(t, err) // 应该超时
    assert.True(t, errors.Is(err, context.DeadlineExceeded))
}

// ✅ 测试并发安全性（我们的业务逻辑）
func TestConcurrentStateConsistency(t *testing.T) {
    var mu sync.Mutex
    sharedState := make([]int, 0)

    var wg sync.WaitGroup
    for i := 0; i < 10; i++ {
        wg.Add(1)
        go func(idx int) {
            defer wg.Done()
            mu.Lock()
            defer mu.Unlock()
            sharedState = append(sharedState, idx)
        }(i)
    }

    // 等待所有任务完成
    wg.Wait()

    // 验证状态一致性
    mu.Lock()
    defer mu.Unlock()
    assert.Equal(t, 10, len(sharedState))
}

// ✅ 测试错误处理和重试逻辑（我们的业务逻辑）
func TestAsyncErrorPropagation(t *testing.T) {
    client := NewHTTPClient()

    // 测试我们的错误处理和重试逻辑
    _, err := client.FetchWithRetry("invalid-url", 3)
    assert.Error(t, err)

    assert.Contains(t, err.Error(), "网络请求失败")
    assert.Contains(t, err.Error(), "重试 3 次后仍然失败")
}
```

**❌ 不应该测试的**：

```go
// ❌ 不要测试 Go 运行时的正确性（这是 Go 运行时的责任）
func TestGoroutineBehavior(t *testing.T) {
    // 这是在测试 Go 运行时本身，而不是我们的代码
    go func() {
        time.Sleep(100 * time.Millisecond)
    }()
}

// ❌ 不要测试标准库的并发原语（这是标准库的责任）
func TestMutexLocking(t *testing.T) {
    // 这是在测试 Mutex 的实现，而不是我们的代码
    var mu sync.Mutex
    mu.Lock()
    defer mu.Unlock()
    assert.Equal(t, 0, 0)
}
```

### 🎯 测试边界总结

| 测试类型 | 应该测试 ✅ | 不应该测试 ❌ |
|---------|------------|--------------|
| **业务逻辑** | 我们的业务规则、数据转换、状态管理 | 外部工具的业务逻辑 |
| **错误处理** | 我们如何处理错误、错误消息、错误恢复 | 外部工具是否会产生错误 |
| **数据结构** | 我们的数据结构、序列化、验证 | 标准库的数据结构 |
| **API 集成** | 我们如何调用 API、处理响应、错误处理 | API 本身的实现和正确性 |
| **Git 操作** | 我们的 Git 封装、参数构建、结果解析 | Git 命令本身的功能 |
| **HTTP 请求** | 我们的请求配置、重试逻辑、响应处理 | HTTP 客户端库的实现 |
| **CLI 命令** | 我们的参数解析、执行流程、用户交互 | Cobra 库本身的参数解析功能 |
| **并发异步** | 我们的并发控制、goroutine 逻辑、错误处理 | Go 运行时和标准库并发原语 |

**关键原则**：**测试我们自己写的代码，信任外部依赖已经过充分测试。**

---

## 🔄 检查流程

### 步骤 1：项目结构扫描

#### 1.1 收集源代码信息
```bash
# 列出所有 lib 模块
echo "=== Core 业务模块 (internal/lib/) ==="
find internal/lib -name "*.go" -not -name "*_test.go" | sort

# 列出所有 commands 模块
echo "=== CLI 命令模块 (internal/commands/) ==="
find internal/commands -name "*.go" -not -name "*_test.go" | sort

# 统计模块数量
echo "Core 模块数量: $(find internal/lib -name "*.go" -not -name "*_test.go" | wc -l)"
echo "Commands 模块数量: $(find internal/commands -name "*.go" -not -name "*_test.go" | wc -l)"
```

#### 1.2 收集测试文件信息
```bash
# 列出所有测试文件
echo "=== 测试文件 (*_test.go) ==="
find . -name "*_test.go" | sort

# 统计测试文件数量
echo "测试文件数量: $(find . -name "*_test.go" | wc -l)"

# 检查测试目录结构
echo "=== 测试目录结构 ==="
tree test/ 2>/dev/null || find test -type d | sort
```

### 步骤 2：覆盖情况检查

#### 2.1 模块覆盖对比
```bash
# 创建模块覆盖检查脚本
cat > check-_coverage.sh << 'EOF'
#!/bin/bash
echo "=== 模块覆盖情况检查 ==="

echo "🟢 已覆盖的模块:"
for lib_file in $(find internal/lib -name "*.go" -not -name "*_test.go"); do
    module_name=$(basename $(dirname $lib_file))
    test_file="${lib_file%%.go}_test.go"
    if [[ -f "$test_file" ]] && [[ $(grep -c "func Test" "$test_file" 2>/dev/null || echo 0) -gt 0 ]]; then
        echo "  ✅ $module_name ($(basename $lib_file))"
    fi
done

echo ""
echo "🟡 部分覆盖的模块:"
for lib_file in $(find internal/lib -name "*.go" -not -name "*_test.go"); do
    module_name=$(basename $(dirname $lib_file))
    test_file="${lib_file%%.go}_test.go"
    if [[ -f "$test_file" ]] && [[ $(grep -c "func Test" "$test_file" 2>/dev/null || echo 0) -eq 0 ]]; then
        echo "  ⚠️  $module_name (测试文件存在但为空)"
    fi
done

echo ""
echo "🔴 缺失覆盖的模块:"
for lib_file in $(find internal/lib -name "*.go" -not -name "*_test.go"); do
    module_name=$(basename $(dirname $lib_file))
    test_file="${lib_file%%.go}_test.go"
    if [[ ! -f "$test_file" ]]; then
        module_name=$(basename $(dirname $lib_file))
        echo "  ❌ $module_name (无测试文件)"
    fi
done
EOF

chmod +x check-_coverage.sh
./check-_coverage.sh
```

#### 2.2 功能覆盖检查
```bash
# 检查公共函数覆盖情况
echo "=== 公共函数覆盖检查 ==="
for module in internal/lib/*/*.go; do
    if [[ "$module" == *"_test.go" ]]; then
        continue
    fi
    module_name=$(basename $(dirname $module))
    echo "检查模块: $module_name"

    # 提取公共函数
    pub_functions=$(grep -n "func [A-Z]" "$module" 2>/dev/null | head -5)
    if [[ -n "$pub_functions" ]]; then
        echo "  公共函数:"
        echo "$pub_functions" | sed 's/^/    /'

        # 检查对应测试
        test_file="${module%%.go}_test.go"
        if [[ -f "$test_file" ]]; then
            test_count=$(grep -c "func Test" "$test_file" 2>/dev/null || echo 0)
            echo "  测试用例数量: $test_count"
        else
            echo "  ❌ 无测试文件"
        fi
    fi
    echo ""
done
```

### 步骤 3：测试质量评估

#### 3.1 测试工具使用检查
```bash
echo "=== 测试工具使用情况 ==="

# 检查表驱动测试使用
table_driven_count=$(grep -r "t.Run\|range.*testCases" . -name "*_test.go" 2>/dev/null | wc -l)
echo "📊 表驱动测试: $table_driven_count 个文件"

# 检查 testify 使用
testify_count=$(grep -r "github.com/stretchr/testify" . -name "*_test.go" 2>/dev/null | wc -l)
echo "📊 testify 使用: $testify_count 个文件"

# 检查 go-cmp 使用
go_cmp_count=$(grep -r "github.com/google/go-cmp" . -name "*_test.go" 2>/dev/null | wc -l)
echo "📊 go-cmp 深度比较: $go_cmp_count 个文件"

# 检查 httptest Mock 测试
httptest_count=$(grep -r "httptest" . -name "*_test.go" 2>/dev/null | wc -l)
echo "📊 httptest Mock 测试: $httptest_count 个文件"

# 检查并发测试
goroutine_count=$(grep -r "go func\|goroutine" . -name "*_test.go" 2>/dev/null | wc -l)
echo "📊 goroutine 并发测试: $goroutine_count 个文件"
```

#### 3.2 测试结构检查
```bash
echo "=== 测试结构和质量检查 ==="

# 检查测试命名规范
echo "🔍 测试命名规范检查:"
non_standard_tests=$(grep -r "func test" . -name "*_test.go" | grep -v "func Test" | wc -l)
if [[ $non_standard_tests -eq 0 ]]; then
    echo "  ✅ 所有测试都遵循 Test 命名规范"
else
    echo "  ⚠️  发现 $non_standard_tests 个不规范的测试命名"
fi

# 检查测试文档注释
documented_tests=$(grep -r "// Test" . -name "*_test.go" 2>/dev/null | wc -l)
total_tests=$(grep -r "func Test" . -name "*_test.go" 2>/dev/null | wc -l)
echo "📝 测试文档覆盖: $documented_tests/$total_tests"

# 检查错误处理测试
error_tests=$(grep -r "assert\.Error\|assert\.NoError\|require\.Error" . -name "*_test.go" 2>/dev/null | wc -l)
echo "🚨 错误处理测试: $error_tests 个"

# 检查边界条件测试
boundary_tests=$(grep -r "empty\|nil\|zero\|max\|min\|boundary" . -name "*_test.go" 2>/dev/null | wc -l)
echo "🎯 边界条件测试: $boundary_tests 个"
```

### 步骤 4：缺失测试识别

#### 4.1 生成缺失测试报告
```bash
cat > generate-_missing-_tests-_report.sh << 'EOF'
#!/bin/bash
echo "# 缺失测试分析报告"
echo ""
echo "## 📊 统计概览"
echo ""

# 统计总体情况
total_lib_modules=$(find internal/lib -name "*.go" -not -name "*_test.go" | wc -l)
total_test_files=$(find . -name "*_test.go" | wc -l)
total_tests=$(grep -r "func Test" . -name "*_test.go" 2>/dev/null | wc -l)

echo "- **Core 模块总数**: $total_lib_modules"
echo "- **测试文件总数**: $total_test_files"
echo "- **测试用例总数**: $total_tests"
echo ""

# 计算覆盖率
covered_modules=0
for lib_file in $(find internal/lib -name "*.go" -not -name "*_test.go"); do
    test_file="${lib_file%%.go}_test.go"
    if [[ -f "$test_file" ]] && [[ $(grep -c "func Test" "$test_file" 2>/dev/null || echo 0) -gt 0 ]]; then
        covered_modules=$((covered_modules + 1))
    fi
done

coverage_percent=$(echo "scale=1; $covered_modules * 100 / $total_lib_modules" | bc -l 2>/dev/null || echo "0")
echo "- **模块覆盖率**: $covered_modules/$total_lib_modules ($coverage_percent%)"
echo ""

echo "## 🔴 完全缺失测试的模块"
echo ""
for lib_file in $(find internal/lib -name "*.go" -not -name "*_test.go"); do
    test_file="${lib_file%%.go}_test.go"
    if [[ ! -f "$test_file" ]]; then
        module_name=$(basename $(dirname $lib_file))
        echo "- ❌ **$module_name** (\`$(basename $lib_file)\`)"
        # 尝试识别主要功能
        pub_functions=$(grep "func [A-Z]" "$lib_file" 2>/dev/null | head -3 | sed 's/.*func \([^(]*\).*/  - \1()/')
        if [[ -n "$pub_functions" ]]; then
            echo "  - 主要功能:"
            echo "$pub_functions"
        fi
        echo ""
    fi
done

echo "## 🟡 测试文件为空的模块"
echo ""
for lib_file in $(find internal/lib -name "*.go" -not -name "*_test.go"); do
    test_file="${lib_file%%.go}_test.go"
    if [[ -f "$test_file" ]] && [[ $(grep -c "func Test" "$test_file" 2>/dev/null || echo 0) -eq 0 ]]; then
        module_name=$(basename $(dirname $lib_file))
        echo "- ⚠️  **$module_name** (测试文件存在但无实际测试)"
        echo ""
    fi
done

echo "## 🟢 已完整覆盖的模块"
echo ""
for lib_file in $(find internal/lib -name "*.go" -not -name "*_test.go"); do
    test_file="${lib_file%%.go}_test.go"
    if [[ -f "$test_file" ]] && [[ $(grep -c "func Test" "$test_file" 2>/dev/null || echo 0) -gt 0 ]]; then
        test_count=$(grep -c "func Test" "$test_file" 2>/dev/null || echo 0)
        module_name=$(basename $(dirname $lib_file))
        echo "- ✅ **$module_name** ($test_count 个测试)"
    fi
done
EOF

chmod +x generate-_missing-_tests-_report.sh
./generate-_missing-_tests-_report.sh
```

### 步骤 5：生成检查报告

#### 5.1 创建完整报告
```bash
cat > generate-_full-_report.sh << 'EOF'
#!/bin/bash
REPORT_FILE="report/TEST_COVERAGE_REPORT_$(date +%Y%m%d_%H%M%S).md"
mkdir -p report

echo "# 测试用例检查报告" > $REPORT_FILE
echo "" >> $REPORT_FILE
echo "**生成时间**: $(date '+%Y-%m-%d %H:%M:%S')" >> $REPORT_FILE
echo "" >> $REPORT_FILE

# 执行所有检查并写入报告
echo "## 📈 覆盖情况总结" >> $REPORT_FILE
echo "" >> $REPORT_FILE

# 这里可以调用之前的检查脚本并将结果追加到报告中
./check-_coverage.sh >> $REPORT_FILE 2>&1
echo "" >> $REPORT_FILE

./generate-_missing-_tests-_report.sh >> $REPORT_FILE 2>&1

echo "## 🛠️ 改进建议" >> $REPORT_FILE
echo "" >> $REPORT_FILE
echo "### 高优先级" >> $REPORT_FILE
echo "- [ ] 补充 Git 模块核心功能测试" >> $REPORT_FILE
echo "- [ ] 添加 Template 模块测试" >> $REPORT_FILE
echo "- [ ] 实现 Branch 模块测试" >> $REPORT_FILE
echo "" >> $REPORT_FILE
echo "### 中优先级" >> $REPORT_FILE
echo "- [ ] 完善 HTTP 模块重试逻辑测试" >> $REPORT_FILE
echo "- [ ] 增强 Proxy 模块配置验证测试" >> $REPORT_FILE
echo "" >> $REPORT_FILE

echo "报告已生成: $REPORT_FILE"
EOF

chmod +x generate-_full-_report.sh
./generate-_full-_report.sh
```

---

## 📊 项目测试现状

### 测试工具配置
```toml
[dev-dependencies]
pretty_assertions = "1.4"    # 清晰的断言输出
rstest = "0.18"             # 参数化测试
mockito = "1.2"             # HTTP Mock 测试
insta = "1.38"              # 快照测试
tempfile = "3.8"            # 临时文件管理
```

### 测试目录结构
```
tests/
├── base/           # Base 模块测试
├── cli/            # CLI 命令测试
├── git/            # Git 模块测试（目前为空）
├── jira/           # Jira 模块测试
├── pr/             # PR 模块测试
├── common/         # 共享测试工具
└── fixtures/       # 测试数据文件
```

### 当前覆盖状态
- 🟢 **已完整覆盖**：Base（LLM、Settings、Dialog）、PR、Jira、CLI 参数解析
- 🟡 **部分覆盖**：HTTP、Completion、Proxy
- 🔴 **缺失覆盖**：Git、Template、Branch、Commit、Stash

---
## 🎯 测试示例

### ✅ 好的测试示例

```go
// ✅ 测试业务逻辑：分支名称格式化
func TestFormatBranchName(t *testing.T) {
	assert.Equal(t, "feature/login", FormatBranchName("feature", "login"))
	assert.Equal(t, "test", FormatBranchName("", "test"))
}

// ✅ 测试错误处理：参数验证
func TestValidateJiraID(t *testing.T) {
	err := ValidateJiraID("PROJ-123")
	assert.NoError(t, err)

	err = ValidateJiraID("invalid")
	assert.Error(t, err)
}

// ✅ 使用 Mock 测试：HTTP API 调用
func TestGitHubAPICall(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"name": "repo"}`))
	}))
	defer server.Close()

	result, err := githubClient.GetRepo("owner/repo")
	assert.NoError(t, err)
	assert.NotNil(t, result)
}

// ✅ 表驱动测试：多种输入验证
func TestBranchNameValidation(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		{"valid", "feature/test", true},
		{"invalid", "invalid//name", false},
		{"empty", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, IsValidBranchName(tt.input))
		})
	}
}
```

### ❌ 避免的测试

```go
// ❌ 不要测试外部工具
func TestGitCommand(t *testing.T) {
	cmd := exec.Command("git", "status")
	cmd.Run()
}

// ❌ 不要测试第三方库
func TestHTTPClient(t *testing.T) {
	resp, err := http.Get("https://api.github.com")
	assert.NoError(t, err)
	defer resp.Body.Close()
}
```

---

## 📋 快速检查脚本

### 检查测试覆盖
```bash
# 检查缺失的测试文件
for module in internal/lib/*/*.go; do
    if [[ "$module" == *"_test.go" ]]; then
        continue
    fi
    test_file="${module%%.go}_test.go"
    if [[ ! -f "$test_file" ]]; then
        echo "❌ 缺失测试: $module"
    fi
done

# 检查空的测试文件
find . -name "*_test.go" -exec sh -c 'if [ $(grep -c "func Test" "$1") -eq 0 ]; then echo "⚠️  空测试文件: $1"; fi' _ {} \;
```

### 检查测试工具使用
```bash
# 检查是否使用推荐的测试工具
grep -r "github.com/stretchr/testify" . -name "*_test.go" || echo "❌ 未使用 testify"
grep -r "github.com/google/go-cmp" . -name "*_test.go" || echo "❌ 未使用 go-cmp"
grep -r "net/http/httptest" . -name "*_test.go" || echo "❌ 未使用 httptest"
grep -r "t.Run\|range.*testCases" . -name "*_test.go" || echo "❌ 未使用表驱动测试"
```

---

## 📊 检查报告模板

```markdown
# 测试用例检查报告

## 📈 覆盖情况总结
- **总测试文件数**: X 个
- **已覆盖模块**: X/Y (Z%)
- **测试用例总数**: ~X 个

## 🟢 已完整覆盖
- Base 模块 (LLM、Settings、Dialog)
- PR 模块 (GitHub 集成、LLM 生成)
- Jira 模块 (API 集成、日志管理)

## 🟡 部分覆盖
- HTTP 模块 (缺少重试逻辑测试)
- Proxy 模块 (缺少配置验证测试)

## 🔴 缺失覆盖
- Git 模块 (测试文件为空)
- Template 模块 (无测试文件)
- Branch 模块 (无测试文件)

## 🛠️ 改进建议
1. **优先级 1**: 补充 Git 模块核心功能测试
2. **优先级 2**: 添加 Template 和 Branch 模块测试
3. **优先级 3**: 完善 HTTP 和 Proxy 模块测试

## 📋 行动计划
- [ ] 创建 Git 模块测试框架
- [ ] 实现分支管理功能测试
- [ ] 添加提交管理功能测试
- [ ] 补充错误处理和边界条件测试
```

---

## 📚 参考文档

- [测试规范指南](../../testing/README.md) - 详细的测试组织和最佳实践
- [开发规范索引](../../development/README.md) - 开发规范总览
- [提交前检查指南](../workflows/pre-commit.md) - 测试检查的简要说明

---

## 📋 检查清单

### 测试边界检查清单

- [ ] 已明确测试边界：测试我们自己的业务逻辑，不测试外部依赖
- [ ] 已识别需要测试的内容（业务逻辑、CLI、错误处理等）
- [ ] 已识别不需要测试的内容（外部工具、第三方库等）

### 测试覆盖检查清单

- [ ] 所有 Lib 层模块都有对应的测试文件
- [ ] 所有 Commands 层模块都有对应的测试文件
- [ ] 新增的公共函数是否有单元测试？
- [ ] 新增的 CLI 命令是否有集成测试？
- [ ] 错误处理路径是否有测试？
- [ ] 边界情况是否有测试？

### 测试质量检查清单

- [ ] 是否使用了推荐的测试工具（rstest、pretty_assertions、mockito、insta 等）？
- [ ] 测试命名是否遵循 `test_` 规范？
- [ ] 测试是否有文档注释？
- [ ] 测试结构是否清晰，易于维护？

### 测试工具使用检查清单

- [ ] 是否使用 `rstest` 进行参数化测试？
- [ ] 是否使用 `pretty_assertions` 提供清晰的断言输出？
- [ ] 是否使用 `mockito` 进行 HTTP API Mock 测试？
- [ ] 是否使用 `insta` 进行快照测试（如适用）？
- [ ] 是否使用 goroutine 和 channel 进行并发测试（如适用）？

### 缺失测试识别检查清单

- [ ] 已识别完全缺失测试的模块
- [ ] 已识别测试文件为空的模块
- [ ] 已识别部分覆盖的模块
- [ ] 已生成缺失测试报告

---

**最后更新**: 2025-01-27

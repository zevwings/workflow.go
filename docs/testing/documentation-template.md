# 测试文档模板

> 本文档提供标准化的测试文档模板，确保所有测试文档格式统一、内容完整。

---

## 📋 目录

- [标准测试模板](#标准测试模板)
- [被忽略测试模板](#被忽略测试模板)
- [参数化测试模板](#参数化测试模板)
- [集成测试模板](#集成测试模板)
- [错误处理测试模板](#错误处理测试模板)
- [使用指南](#使用指南)

---

## 📝 标准测试模板

### 基础模板

```go
// 测试 {功能描述}
//
// ## 测试目的
// 验证 {功能} 能够 {预期行为}。
//
// ## 测试场景
// 1. {步骤1}
// 2. {步骤2}
// 3. {步骤3}
//
// ## 预期结果
// - {结果1}
// - {结果2}
```

### 完整模板（包含可选部分）

```go
// 测试 {功能描述}
//
// ## 测试目的
// 验证 {功能} 能够 {预期行为}。
//
// ## 测试场景
// 1. {步骤1}
// 2. {步骤2}
// 3. {步骤3}
//
// ## 预期结果
// - {结果1}
// - {结果2}
//
// ## 技术细节（可选）
// - 使用 {工具/库} 进行测试
// - 测试隔离方式：{隔离机制}
// - Mock 服务器配置：{配置说明}
//
// ## 注意事项（可选）
// - {特殊要求或限制}
// - {平台相关说明}
```

### 实际示例

```go
// 测试 GET 请求成功响应
//
// ## 测试目的
// 验证 HTTP 客户端能够正确发送 GET 请求并处理成功响应。
//
// ## 测试场景
// 1. 配置 Mock 服务器返回 200 状态码和 JSON 响应
// 2. 发送 GET 请求
// 3. 验证响应状态码和内容
//
// ## 预期结果
// - 响应状态码为 200
// - 响应标记为成功
// - Mock 服务器收到预期的请求
//
// ## 技术细节
// - 使用 MockServer 模拟 HTTP 服务器
// - 测试隔离方式：每个测试使用独立的 MockServer
func TestGetRequestWithSuccessResponseReturnsResponse(t *testing.T) {
    // Arrange: 准备 Mock 服务器和响应
    mockServer := setupMockServer(t)
    defer mockServer.Close()
    // ...

    // Act: 发送 GET 请求
    response, err := client.Get(url, config)
    if err != nil {
        t.Fatalf("unexpected error: %v", err)
    }

    // Assert: 验证响应状态和内容
    if response.StatusCode != 200 {
        t.Errorf("expected status 200, got %d", response.StatusCode)
    }
    if !response.IsSuccess() {
        t.Error("expected response to be successful")
    }
}
```

---

## 🚫 被忽略测试模板

### 基础模板

```go
// 测试 {功能描述}
//
// ## 测试目的
// 验证 {功能} 能够 {预期行为}。
//
// ## 为什么被忽略
// - {原因1}
// - {原因2}
// - {原因3}
//
// ## 如何手动运行
// ```bash
// go test -v -run TestName -tags=integration
// ```
//
// ## 测试场景
// 1. {步骤1}
// 2. {步骤2}
//
// ## 预期结果
// - {结果1}
```

### 实际示例

```go
// 测试在非Git仓库中执行PR命令
//
// ## 测试目的
// 验证当不在Git仓库中时，PR命令能够正确检测并返回清晰的错误消息。
//
// ## 为什么被忽略
// - **可能初始化客户端**: 即使在非Git仓库中，可能仍尝试初始化Jira/GitHub客户端
// - **长时间阻塞**: 客户端初始化可能导致测试阻塞
// - **平台限制**: Windows上已通过构建标签跳过
// - **CI时间考虑**: 避免在CI中长时间阻塞
//
// ## 如何手动运行
// ```bash
// go test -v -run TestPRWithoutGitRepo -tags=integration
// ```
// 注意：此测试在非Windows平台上运行
//
// ## 测试场景
// 1. 创建临时目录（非Git仓库）
// 2. 在该目录中执行`pr create --dry-run`命令
// 3. 命令应该失败
// 4. 验证错误消息提示"没有Git仓库"
//
// ## 预期结果
// - 命令执行失败（非零退出码）
// - stderr包含错误消息
// - 错误消息清晰说明原因（不在Git仓库中）
// +build !windows

func TestPRWithoutGitRepo(t *testing.T) {
    if testing.Short() {
        t.Skip("skipping integration test in short mode")
    }
    // ...
}
```

---

## 🔄 参数化测试模板

### 基础模板

```go
// 测试 {功能描述}（参数化测试）
//
// ## 测试目的
// 使用参数化测试验证 {功能} 能够处理多种输入情况。
//
// ## 测试场景
// 测试 {输入类型1}、{输入类型2}、{输入类型3} 等各种输入
//
// ## 预期结果
// - 所有测试用例都能正确处理
// - 输出符合预期
func Test{功能}WithVariousInputs(t *testing.T) {
    tests := []struct {
        name     string
        input    {类型}
        expected {类型}
    }{
        {
            name:     "{用例1描述}",
            input:    {输入1},
            expected: {预期1},
        },
        {
            name:     "{用例2描述}",
            input:    {输入2},
            expected: {预期2},
        },
        {
            name:     "{用例3描述}",
            input:    {输入3},
            expected: {预期3},
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // Arrange: 准备测试数据

            // Act: 执行被测试的功能
            result := functionUnderTest(tt.input)

            // Assert: 验证结果
            if result != tt.expected {
                t.Errorf("expected %v, got %v", tt.expected, result)
            }
        })
    }
}
```

### 实际示例

```go
// 测试创建 PR 请求序列化（参数化测试）
//
// ## 测试目的
// 使用参数化测试验证 CreatePullRequestRequest 能够正确序列化为 JSON。
//
// ## 测试场景
// 测试正常输入、空字符串、特殊字符、多行文本等各种输入
//
// ## 预期结果
// - JSON 字段存在且值正确
// - 序列化/反序列化一致
func TestCreatePRRequestSerializationWithVariousInputsSerializesCorrectly(t *testing.T) {
    tests := []struct {
        name  string
        title string
        body  string
        head  string
        base  string
    }{
        {
            name:  "正常输入",
            title: "Test PR",
            body:  "Test body",
            head:  "feature/test",
            base:  "main",
        },
        {
            name:  "空字符串",
            title: "",
            body:  "",
            head:  "",
            base:  "",
        },
        {
            name:  "特殊字符和多行文本",
            title: "Long Title with Special Chars !@#",
            body:  "Long Body\nwith\nmultiple\nlines",
            head:  "feature/long-branch-name",
            base:  "develop",
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // Arrange: 准备 CreatePullRequestRequest 实例
            request := CreatePullRequestRequest{
                Title: tt.title,
                Body:  tt.body,
                Head:  tt.head,
                Base:  tt.base,
            }

            // Act: 序列化为 JSON
            jsonBytes, err := json.Marshal(request)
            if err != nil {
                t.Fatalf("failed to marshal: %v", err)
            }

            var jsonValue map[string]interface{}
            if err := json.Unmarshal(jsonBytes, &jsonValue); err != nil {
                t.Fatalf("failed to unmarshal: %v", err)
            }

            // Assert: 验证 JSON 字段存在且值正确
            if title, ok := jsonValue["title"].(string); !ok || title != tt.title {
                t.Errorf("expected title %q, got %v", tt.title, jsonValue["title"])
            }
            if body, ok := jsonValue["body"].(string); !ok || body != tt.body {
                t.Errorf("expected body %q, got %v", tt.body, jsonValue["body"])
            }
        })
    }
}
```

---

## 🔗 集成测试模板

### 基础模板

```go
// 测试 {功能} 集成场景
//
// ## 测试目的
// 验证 {功能} 在完整工作流程中的行为。
//
// ## 测试场景
// 1. {初始化步骤}
// 2. {执行步骤1}
// 3. {执行步骤2}
// 4. {验证步骤}
//
// ## 预期结果
// - {结果1}
// - {结果2}
//
// ## 技术细节
// - 使用 {测试环境} 进行集成测试
// - 测试隔离方式：{隔离机制}
func Test{功能}Integration(t *testing.T) {
    if testing.Short() {
        t.Skip("skipping integration test in short mode")
    }

    // Arrange: 准备集成测试环境
    env, err := NewCLITestEnv()
    if err != nil {
        t.Fatalf("failed to create test env: %v", err)
    }
    defer env.Cleanup()

    if err := env.InitGitRepo(); err != nil {
        t.Fatalf("failed to init git repo: %v", err)
    }

    // Act: 执行完整工作流程
    // ...

    // Assert: 验证最终状态
    // ...
}
```

### 实际示例

```go
// 测试 Git 提交工作流程集成
//
// ## 测试目的
// 验证 Git 提交功能在完整工作流程中的行为，包括文件修改、暂存和提交。
//
// ## 测试场景
// 1. 创建 Git 测试环境
// 2. 创建并修改文件
// 3. 暂存文件
// 4. 创建提交
// 5. 验证提交成功
//
// ## 预期结果
// - 文件修改成功
// - 暂存成功
// - 提交成功
// - 工作树状态为干净
//
// ## 技术细节
// - 使用 GitTestEnv 进行集成测试
// - 测试隔离方式：每个测试使用独立的临时目录
func TestGitCommitWorkflowIntegration(t *testing.T) {
    if testing.Short() {
        t.Skip("skipping integration test in short mode")
    }

    // Arrange: 准备 Git 测试环境
    env, err := NewGitTestEnv()
    if err != nil {
        t.Fatalf("failed to create git test env: %v", err)
    }
    defer env.Cleanup()

    // Act: 执行完整工作流程
    testFile := filepath.Join(env.Path(), "test.txt")
    if err := os.WriteFile(testFile, []byte("test content"), 0644); err != nil {
        t.Fatalf("failed to write test file: %v", err)
    }

    // 暂存文件
    if err := GitCommit.StageFiles([]string{testFile}); err != nil {
        t.Fatalf("failed to stage files: %v", err)
    }

    // 创建提交
    if err := GitCommit.Create("test: add test file"); err != nil {
        t.Fatalf("failed to create commit: %v", err)
    }

    // Assert: 验证提交成功
    status, err := GitCommit.GetWorktreeStatus()
    if err != nil {
        t.Fatalf("failed to get worktree status: %v", err)
    }
    if status.StagedCount != 0 {
        t.Errorf("expected staged count 0, got %d", status.StagedCount)
    }
    if status.ModifiedCount != 0 {
        t.Errorf("expected modified count 0, got %d", status.ModifiedCount)
    }
}
```

---

## ⚠️ 错误处理测试模板

### 基础模板

```go
// 测试 {功能} 错误处理
//
// ## 测试目的
// 验证 {功能} 在 {错误条件} 时能够正确处理错误。
//
// ## 测试场景
// 1. {准备错误条件}
// 2. 调用 {功能}
// 3. 验证错误处理
//
// ## 预期结果
// - 返回错误（error != nil）
// - 错误消息清晰
// - 错误类型正确
func Test{功能}With{错误条件}ReturnsError(t *testing.T) {
    // Arrange: 准备错误条件
    // ...

    // Act: 调用功能
    result, err := functionUnderTest(input)

    // Assert: 验证错误处理
    if err == nil {
        t.Fatal("expected error, got nil")
    }
    if result != nil {
        t.Errorf("expected nil result, got %v", result)
    }
    if !strings.Contains(err.Error(), "expected error message") {
        t.Errorf("error message should contain 'expected error message', got: %v", err)
    }
}
```

### 实际示例

```go
// 测试 GET 请求网络错误处理
//
// ## 测试目的
// 验证 HTTP 客户端在网络错误时能够正确处理错误。
//
// ## 测试场景
// 1. 配置 Mock 服务器返回网络错误
// 2. 发送 GET 请求
// 3. 验证错误处理
//
// ## 预期结果
// - 返回错误（error != nil）
// - 错误消息包含网络错误信息
// - 错误类型为网络错误
func TestGetRequestWithNetworkErrorHandlesGracefully(t *testing.T) {
    // Arrange: 准备 Mock 服务器返回网络错误
    mockServer := setupMockServer(t)
    defer mockServer.Close()
    url := fmt.Sprintf("%s/test", mockServer.BaseURL)

    // 配置服务器返回连接错误
    mockServer.Mock("GET", "/test", func(w http.ResponseWriter, r *http.Request) {
        w.WriteHeader(http.StatusInternalServerError)
    })

    // Act: 发送 GET 请求
    client := http.DefaultClient
    config := NewRequestConfig()
    _, err := client.Get(url, config)

    // Assert: 验证错误处理
    if err == nil {
        t.Fatal("expected error, got nil")
    }
    errorMsg := err.Error()
    if !strings.Contains(errorMsg, "network") && !strings.Contains(errorMsg, "connection") {
        t.Errorf("error message should contain 'network' or 'connection', got: %v", errorMsg)
    }
}
```

---

## 📖 使用指南

### 何时使用哪个模板

| 测试类型 | 使用模板 | 说明 |
|---------|---------|------|
| 标准功能测试 | 标准测试模板 | 大多数测试使用此模板 |
| 被忽略的测试 | 被忽略测试模板 | 需要 `t.Skip()` 或构建标签的测试 |
| 多输入测试 | 参数化测试模板 | 使用表格驱动测试的测试 |
| 端到端测试 | 集成测试模板 | 测试完整工作流程 |
| 错误场景测试 | 错误处理测试模板 | 测试错误处理逻辑 |

### 文档编写原则

1. **完整性**: 包含所有必需部分（测试目的、场景、预期结果）
2. **清晰性**: 使用简洁明了的语言
3. **一致性**: 遵循模板格式，保持统一
4. **准确性**: 确保文档与实际测试代码一致

### 可选部分使用指南

- **技术细节**: 当测试使用特殊工具或技术时添加
- **注意事项**: 当测试有特殊要求或限制时添加
- **平台相关说明**: 当测试行为因平台而异时添加

### 文档更新

- 修改测试代码时，同步更新文档
- 代码审查时检查文档完整性
- 定期审查文档格式一致性

---

## 📚 相关文档

- [测试组织规范](./organization.md)
- [测试编写规范](./writing.md)
- [测试命令参考](./commands.md)

---

**文档版本**: 1.0
**最后更新**: 2024年（当前）


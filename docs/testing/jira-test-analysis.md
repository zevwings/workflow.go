# Jira 模块测试用例分析

> 本文档分析 `internal/jira` 模块的测试策略和测试用例设计。

---

## 📋 目录

- [模块结构分析](#-模块结构分析)
- [测试策略](#-测试策略)
- [测试用例设计](#-测试用例设计)
- [Mock 策略](#-mock-策略)
- [测试文件组织](#-测试文件组织)
- [实施优先级](#-实施优先级)

---

## 📋 模块结构分析

### 模块组件

1. **helpers.go** - 辅助函数（纯函数，无外部依赖）
   - `ValidateTicketKey()` - 验证 Ticket Key 格式
   - `NormalizeTicketKey()` - 规范化 Ticket Key
   - `ExtractProjectKey()` - 提取项目 Key
   - `ExtractTicketNumber()` - 提取 Ticket 编号

2. **client.go** - 底层客户端封装
   - `NewClient()` - 创建客户端
   - `WithContext()` - 使用自定义 context
   - `GetJiraClient()` - 获取底层客户端
   - `GetContext()` - 获取 context

3. **jira_client.go** - 高级封装客户端
   - `NewJiraClient()` - 创建高级客户端
   - 各种业务方法（GetTicketInfo, GetAttachments, AddComment 等）

4. **api/** - API 模块
   - `issue.go` - Issue API
   - `project.go` - Project API
   - `user.go` - User API

---

## 🎯 测试策略

### 测试分层

1. **单元测试（Unit Tests）**
   - 测试纯函数（helpers.go）
   - 测试客户端创建和配置
   - 使用 Mock 测试 API 调用

2. **集成测试（Integration Tests）**
   - 测试与真实 Jira API 的交互（可选，需要 API 密钥）
   - 使用构建标签 `//go:build integration`

### Mock 策略

由于 Jira 是外部 API，需要使用 Mock：

1. **HTTP Mock（推荐）**
   - 使用 `httptest.NewServer` 创建 Mock Jira 服务器
   - 模拟 Jira REST API 响应

2. **接口 Mock（可选）**
   - 如果重构为接口，可以使用 `testify/mock`
   - 当前实现直接使用 `go-jira` SDK，较难 Mock

---

## 📝 测试用例设计

### 1. helpers_test.go - 辅助函数测试

#### 1.1 ValidateTicketKey 测试

```go
func TestValidateTicketKey(t *testing.T) {
    tests := []struct {
        name    string
        input   string
        wantErr bool
    }{
        {
            name:    "valid ticket key",
            input:    "PROJ-123",
            wantErr:  false,
        },
        {
            name:    "valid ticket key lowercase",
            input:    "proj-123",
            wantErr:  false,
        },
        {
            name:    "empty string",
            input:    "",
            wantErr:  true,
        },
        {
            name:    "invalid format - no dash",
            input:    "PROJ123",
            wantErr:  true,
        },
        {
            name:    "invalid format - multiple dashes",
            input:    "PROJ-123-456",
            wantErr:  true,
        },
        {
            name:    "missing project key",
            input:    "-123",
            wantErr:  true,
        },
        {
            name:    "missing ticket number",
            input:    "PROJ-",
            wantErr:  true,
        },
        {
            name:    "only dash",
            input:    "-",
            wantErr:  true,
        },
        {
            name:    "minimum valid length",
            input:    "A-1",
            wantErr:  false,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            err := ValidateTicketKey(tt.input)
            if tt.wantErr {
                assert.Error(t, err)
            } else {
                assert.NoError(t, err)
            }
        })
    }
}
```

#### 1.2 NormalizeTicketKey 测试

```go
func TestNormalizeTicketKey(t *testing.T) {
    tests := []struct {
        name     string
        input    string
        expected string
    }{
        {
            name:     "uppercase",
            input:    "PROJ-123",
            expected: "PROJ-123",
        },
        {
            name:     "lowercase",
            input:    "proj-123",
            expected: "PROJ-123",
        },
        {
            name:     "mixed case",
            input:    "Proj-123",
            expected: "PROJ-123",
        },
        {
            name:     "with spaces",
            input:    "  proj-123  ",
            expected: "PROJ-123",
        },
        {
            name:     "with leading spaces",
            input:    "  PROJ-123",
            expected: "PROJ-123",
        },
        {
            name:     "with trailing spaces",
            input:    "PROJ-123  ",
            expected: "PROJ-123",
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            result := NormalizeTicketKey(tt.input)
            assert.Equal(t, tt.expected, result)
        })
    }
}
```

#### 1.3 ExtractProjectKey 测试

```go
func TestExtractProjectKey(t *testing.T) {
    tests := []struct {
        name     string
        input    string
        expected string
    }{
        {
            name:     "valid ticket key",
            input:    "PROJ-123",
            expected: "PROJ",
        },
        {
            name:     "lowercase",
            input:    "proj-123",
            expected: "proj",
        },
        {
            name:     "long project key",
            input:    "VERY-LONG-PROJECT-123",
            expected: "VERY",
        },
        {
            name:     "single character",
            input:    "A-1",
            expected: "A",
        },
        {
            name:     "invalid format",
            input:    "invalid",
            expected: "invalid",
        },
        {
            name:     "empty string",
            input:    "",
            expected: "",
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            result := ExtractProjectKey(tt.input)
            assert.Equal(t, tt.expected, result)
        })
    }
}
```

#### 1.4 ExtractTicketNumber 测试

```go
func TestExtractTicketNumber(t *testing.T) {
    tests := []struct {
        name     string
        input    string
        expected string
    }{
        {
            name:     "valid ticket key",
            input:    "PROJ-123",
            expected: "123",
        },
        {
            name:     "large number",
            input:    "PROJ-999999",
            expected: "999999",
        },
        {
            name:     "single digit",
            input:    "PROJ-1",
            expected: "1",
        },
        {
            name:     "invalid format",
            input:    "invalid",
            expected: "",
        },
        {
            name:     "empty string",
            input:    "",
            expected: "",
        },
        {
            name:     "no dash",
            input:    "PROJ123",
            expected: "",
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            result := ExtractTicketNumber(tt.input)
            assert.Equal(t, tt.expected, result)
        })
    }
}
```

---

### 2. client_test.go - 客户端测试

#### 2.1 NewClient 测试

```go
func TestNewClient(t *testing.T) {
    tests := []struct {
        name    string
        config  *Config
        wantErr bool
    }{
        {
            name: "valid config",
            config: &Config{
                URL:      "https://test.atlassian.net",
                Username: "test@example.com",
                Token:    "test-token",
            },
            wantErr: false,
        },
        {
            name:    "nil config",
            config:  nil,
            wantErr: true,
        },
        {
            name: "empty URL",
            config: &Config{
                URL:      "",
                Username: "test@example.com",
                Token:    "test-token",
            },
            wantErr: true,
        },
        {
            name: "empty username",
            config: &Config{
                URL:      "https://test.atlassian.net",
                Username: "",
                Token:    "test-token",
            },
            wantErr: true,
        },
        {
            name: "empty token",
            config: &Config{
                URL:      "https://test.atlassian.net",
                Username: "test@example.com",
                Token:    "",
            },
            wantErr: true,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            client, err := NewClient(tt.config)
            if tt.wantErr {
                assert.Error(t, err)
                assert.Nil(t, client)
            } else {
                // 注意：这里需要 Mock Jira 客户端创建
                // 或者使用集成测试
                // 由于 go-jira 会尝试连接，单元测试可能需要跳过
                // 或者使用接口 Mock
            }
        })
    }
}
```

#### 2.2 WithContext 测试

```go
func TestClient_WithContext(t *testing.T) {
    config := &Config{
        URL:      "https://test.atlassian.net",
        Username: "test@example.com",
        Token:    "test-token",
    }

    client, err := NewClient(config)
    if err != nil {
        t.Skip("需要有效的 Jira 配置或 Mock")
    }

    ctx := context.WithTimeout(context.Background(), 5*time.Second)
    newClient := client.WithContext(ctx)

    assert.NotEqual(t, client, newClient) // 应该是新实例
    assert.Equal(t, newClient.GetContext(), ctx)
    assert.Equal(t, client.GetJiraClient(), newClient.GetJiraClient()) // 共享底层客户端
}
```

#### 2.3 GetJiraClient 和 GetContext 测试

```go
func TestClient_GetJiraClient(t *testing.T) {
    // 需要 Mock 或集成测试
}

func TestClient_GetContext(t *testing.T) {
    config := &Config{
        URL:      "https://test.atlassian.net",
        Username: "test@example.com",
        Token:    "test-token",
    }

    client, err := NewClient(config)
    if err != nil {
        t.Skip("需要有效的 Jira 配置或 Mock")
    }

    ctx := client.GetContext()
    assert.NotNil(t, ctx)
}
```

---

### 3. jira_client_test.go - 高级客户端测试

#### 3.1 NewJiraClient 测试

```go
func TestNewJiraClient(t *testing.T) {
    config := &Config{
        URL:      "https://test.atlassian.net",
        Username: "test@example.com",
        Token:    "test-token",
    }

    client, err := NewJiraClient(config)
    if err != nil {
        t.Skip("需要有效的 Jira 配置或 Mock")
    }

    assert.NotNil(t, client)
    assert.NotNil(t, client.GetClient())
    assert.NotNil(t, client.GetIssueAPI())
    assert.NotNil(t, client.GetProjectAPI())
    assert.NotNil(t, client.GetUserAPI())
}
```

#### 3.2 GetTicketInfo 测试

```go
func TestJiraClient_GetTicketInfo(t *testing.T) {
    // 使用 Mock HTTP 服务器
    server := setupMockJiraServer(t)
    defer server.Close()

    // 创建客户端（需要修改以支持自定义 URL）
    // 或者使用接口 Mock

    tests := []struct {
        name    string
        ticket  string
        wantErr bool
    }{
        {
            name:    "valid ticket",
            ticket:  "PROJ-123",
            wantErr: false,
        },
        {
            name:    "invalid format",
            ticket:  "invalid",
            wantErr: true,
        },
        {
            name:    "empty string",
            ticket:  "",
            wantErr: true,
        },
        {
            name:    "lowercase ticket",
            ticket:  "proj-123",
            wantErr: false, // 应该被规范化
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // 测试代码
        })
    }
}
```

#### 3.3 MoveTicket 测试

```go
func TestJiraClient_MoveTicket(t *testing.T) {
    // 测试状态转换逻辑
    tests := []struct {
        name         string
        ticket       string
        status       string
        transitions  []cloud.Transition
        wantErr      bool
        expectedErr  string
    }{
        {
            name:   "valid transition",
            ticket: "PROJ-123",
            status: "In Progress",
            transitions: []cloud.Transition{
                {ID: "11", Name: "In Progress"},
                {ID: "21", Name: "Done"},
            },
            wantErr: false,
        },
        {
            name:   "status not found",
            ticket: "PROJ-123",
            status: "Invalid Status",
            transitions: []cloud.Transition{
                {ID: "11", Name: "In Progress"},
            },
            wantErr:     true,
            expectedErr: "未找到状态转换",
        },
        {
            name:    "invalid ticket format",
            ticket:  "invalid",
            status:  "In Progress",
            wantErr: true,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // 使用 Mock 测试
        })
    }
}
```

---

### 4. api/issue_test.go - Issue API 测试

#### 4.1 GetIssue 测试

```go
func TestIssueAPI_GetIssue(t *testing.T) {
    // 使用 Mock HTTP 服务器模拟 Jira API
    server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        assert.Equal(t, "/rest/api/3/issue/PROJ-123", r.URL.Path)
        assert.Equal(t, http.MethodGet, r.Method)

        w.WriteHeader(http.StatusOK)
        w.Write([]byte(`{
            "id": "10000",
            "key": "PROJ-123",
            "fields": {
                "summary": "Test Issue"
            }
        }`))
    }))
    defer server.Close()

    // 创建 Mock 客户端（需要适配 go-jira SDK）
    // ...
}
```

#### 4.2 GetIssueAttachments 测试

```go
func TestIssueAPI_GetIssueAttachments(t *testing.T) {
    tests := []struct {
        name        string
        ticket      string
        attachments []*cloud.Attachment
        wantErr     bool
    }{
        {
            name:   "with attachments",
            ticket: "PROJ-123",
            attachments: []*cloud.Attachment{
                {ID: "10000", Filename: "test.txt"},
            },
            wantErr: false,
        },
        {
            name:        "no attachments",
            ticket:      "PROJ-123",
            attachments: nil,
            wantErr:     false,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // Mock 测试
        })
    }
}
```

---

### 5. api/project_test.go - Project API 测试

#### 5.1 GetProject 测试

```go
func TestProjectAPI_GetProject(t *testing.T) {
    // Mock 测试
}
```

#### 5.2 GetProjectStatuses 测试

```go
func TestProjectAPI_GetProjectStatuses(t *testing.T) {
    // 注意：当前实现返回空列表
    // 测试应该验证返回空列表的行为
}
```

---

### 6. api/user_test.go - User API 测试

#### 6.1 GetCurrentUser 测试

```go
func TestUserAPI_GetCurrentUser(t *testing.T) {
    // Mock 测试
}
```

#### 6.2 FindUsers 测试

```go
func TestUserAPI_FindUsers(t *testing.T) {
    tests := []struct {
        name    string
        query   string
        users   []*cloud.User
        wantErr bool
    }{
        {
            name:  "find users",
            query: "john",
            users: []*cloud.User{
                {AccountID: "123", DisplayName: "John Doe"},
            },
            wantErr: false,
        },
        {
            name:    "empty query",
            query:   "",
            users:   nil,
            wantErr: false, // 可能返回空列表或错误
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // Mock 测试
        })
    }
}
```

---

## 🔧 Mock 策略

### 方案 1：HTTP Mock 服务器（推荐）

创建 Mock Jira REST API 服务器：

```go
func setupMockJiraServer(t *testing.T) *httptest.Server {
    t.Helper()

    server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // 根据路径和方法返回不同的响应
        switch {
        case r.URL.Path == "/rest/api/3/issue/PROJ-123" && r.Method == http.MethodGet:
            w.WriteHeader(http.StatusOK)
            w.Write([]byte(`{
                "id": "10000",
                "key": "PROJ-123",
                "fields": {
                    "summary": "Test Issue",
                    "status": {"name": "To Do"}
                }
            }`))
        case r.URL.Path == "/rest/api/3/issue/PROJ-123/transitions" && r.Method == http.MethodGet:
            w.WriteHeader(http.StatusOK)
            w.Write([]byte(`{
                "transitions": [
                    {"id": "11", "name": "In Progress"},
                    {"id": "21", "name": "Done"}
                ]
            }`))
        default:
            w.WriteHeader(http.StatusNotFound)
        }
    }))

    return server
}
```

### 方案 2：接口抽象（重构后）

如果重构为接口，可以使用 `testify/mock`：

```go
type JiraClientInterface interface {
    GetIssue(ticket string) (*cloud.Issue, error)
    // ...
}

type MockJiraClient struct {
    mock.Mock
}

func (m *MockJiraClient) GetIssue(ticket string) (*cloud.Issue, error) {
    args := m.Called(ticket)
    return args.Get(0).(*cloud.Issue), args.Error(1)
}
```

---

## 📁 测试文件组织

```
internal/jira/
├── helpers.go
├── helpers_test.go          # 辅助函数测试
├── client.go
├── client_test.go            # 客户端测试
├── jira_client.go
├── jira_client_test.go      # 高级客户端测试
├── types.go
├── api/
│   ├── issue.go
│   ├── issue_test.go        # Issue API 测试
│   ├── project.go
│   ├── project_test.go      # Project API 测试
│   ├── user.go
│   └── user_test.go         # User API 测试
└── testdata/                # 测试数据（可选）
    └── fixtures/
        └── jira_responses.json
```

---

## 🎯 实施优先级

### 阶段 1：基础测试（高优先级）

1. ✅ **helpers_test.go** - 纯函数测试，无外部依赖
   - 易于实现
   - 高覆盖率
   - 快速执行

2. ✅ **client_test.go** - 客户端创建和配置测试
   - 测试配置验证
   - 测试 Context 处理

### 阶段 2：API 测试（中优先级）

3. ⚠️ **api/issue_test.go** - Issue API 测试
   - 需要 Mock Jira API
   - 测试主要业务逻辑

4. ⚠️ **api/project_test.go** - Project API 测试
   - 相对简单

5. ⚠️ **api/user_test.go** - User API 测试
   - 相对简单

### 阶段 3：集成测试（低优先级）

6. 🔄 **jira_client_test.go** - 高级客户端测试
   - 依赖 API 测试完成
   - 可能需要集成测试环境

### 阶段 4：集成测试（可选）

7. 🔄 集成测试（需要真实 Jira 环境）
   - 使用构建标签 `//go:build integration`
   - 需要 API 密钥配置

---

## 📝 注意事项

### 1. go-jira SDK 的 Mock 挑战

`go-jira` SDK 直接创建 HTTP 客户端，较难 Mock。建议：

- **方案 A**：使用 `httptest.NewServer` 创建 Mock 服务器，修改客户端 URL
- **方案 B**：重构为接口，使用 `testify/mock`
- **方案 C**：使用集成测试（需要真实环境）

### 2. 测试数据

建议创建 `testdata/fixtures/` 目录存放 Jira API 响应示例：

```json
{
  "issue": {
    "id": "10000",
    "key": "PROJ-123",
    "fields": {
      "summary": "Test Issue"
    }
  },
  "transitions": [
    {"id": "11", "name": "In Progress"}
  ]
}
```

### 3. 错误处理测试

确保测试各种错误情况：
- 网络错误
- API 错误（404, 401, 500 等）
- 无效输入
- 空响应

### 4. 测试覆盖率目标

- **helpers.go**: 100%（纯函数）
- **client.go**: > 80%
- **jira_client.go**: > 70%
- **api/**: > 70%

---

## 🔗 相关文档

- [测试编写规范](../writing.md)
- [单元测试指南](./references/unit-tests.md)
- [Mock 测试指南](./references/mock-server.md)
- [测试组织规范](./organization.md)

---

**最后更新**: 2025-01-28


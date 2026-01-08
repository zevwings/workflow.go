package github

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"github.com/google/go-github/v57/github"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/oauth2"
)

// ==================== 测试辅助函数 ====================

// setupMockGitHubServer 创建 Mock GitHub API 服务器
func setupMockGitHubServer(t *testing.T, handler http.HandlerFunc) *httptest.Server {
	t.Helper()

	if handler == nil {
		handler = defaultMockGitHubHandler
	}

	server := httptest.NewServer(handler)
	return server
}

// defaultMockGitHubHandler 默认的 Mock GitHub API 处理函数
func defaultMockGitHubHandler(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path

	switch {
	case path == "/repos/owner/repo" && r.Method == http.MethodGet:
		// 获取仓库信息
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(&github.Repository{
			DefaultBranch: github.String("main"),
		})

	case path == "/repos/owner/repo/pulls" && r.Method == http.MethodPost:
		// 创建 PR
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(&github.PullRequest{
			Number:  github.Int(123),
			HTMLURL: github.String("https://github.com/owner/repo/pull/123"),
		})

	case path == "/repos/owner/repo/pulls/123" && r.Method == http.MethodGet:
		// 获取 PR 信息
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		mergeable := true
		json.NewEncoder(w).Encode(&github.PullRequest{
			Number:    github.Int(123),
			State:     github.String("open"),
			Merged:    github.Bool(false),
			Mergeable: &mergeable,
			UpdatedAt: &github.Timestamp{Time: time.Now()},
		})

	case path == "/repos/owner/repo/pulls/123" && r.Method == http.MethodPatch:
		// 更新 PR
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(&github.PullRequest{
			Number: github.Int(123),
		})

	case path == "/repos/owner/repo/pulls/123/merge" && r.Method == http.MethodPut:
		// 合并 PR
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(&github.PullRequestMergeResult{
			Merged: github.Bool(true),
		})

	case path == "/repos/owner/repo/git/refs/heads/feature-branch" && r.Method == http.MethodDelete:
		// 删除分支
		w.WriteHeader(http.StatusNoContent)

	case path == "/repos/owner/repo/issues/123/comments" && r.Method == http.MethodPost:
		// 添加评论
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(&github.IssueComment{
			ID:   github.Int64(1),
			Body: github.String("Test comment"),
		})

	case path == "/repos/owner/repo/pulls/123/reviews" && r.Method == http.MethodPost:
		// 批准 PR
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(&github.PullRequestReview{
			ID:    github.Int64(1),
			State: github.String("APPROVED"),
		})

	case path == "/repos/owner/repo/pulls" && r.Method == http.MethodGet:
		// 列出 PRs
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode([]*github.PullRequest{
			{
				Number:    github.Int(123),
				Title:     github.String("Test PR 1"),
				State:     github.String("open"),
				HTMLURL:   github.String("https://github.com/owner/repo/pull/123"),
				CreatedAt: &github.Timestamp{Time: time.Now()},
				UpdatedAt: &github.Timestamp{Time: time.Now()},
				User: &github.User{
					Login: github.String("testuser"),
				},
			},
		})

	default:
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(&github.ErrorResponse{
			Message: "Not Found",
		})
	}
}

// createTestGitHubClient 创建用于测试的 GitHub 客户端
func createTestGitHubClient(t *testing.T, serverURL string) *GitHub {
	t.Helper()

	// 解析 Mock 服务器 URL
	mockURL, err := url.Parse(serverURL)
	require.NoError(t, err)

	// 创建 OAuth token source
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: "test-token"},
	)

	// 创建自定义 Transport，将请求转发到 Mock 服务器
	mockTransport := &mockTransport{
		mockURL:       mockURL,
		baseTransport: oauth2.NewClient(context.Background(), ts).Transport,
	}

	// 创建 HTTP 客户端
	httpClient := &http.Client{
		Transport: mockTransport,
	}

	// 创建 GitHub 客户端
	client := github.NewClient(httpClient)

	return &GitHub{
		client: client,
		ctx:    context.Background(),
		owner:  "owner",
		repo:   "repo",
	}
}

// mockTransport 自定义 Transport，将请求转发到 Mock 服务器
type mockTransport struct {
	mockURL       *url.URL
	baseTransport http.RoundTripper
}

func (m *mockTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	// 创建新请求，避免修改原始请求
	newReq := req.Clone(req.Context())

	// 修改请求 URL，指向 Mock 服务器，但保留原始路径和查询参数
	newReq.URL.Scheme = m.mockURL.Scheme
	newReq.URL.Host = m.mockURL.Host

	return m.baseTransport.RoundTrip(newReq)
}

// ==================== NewGitHub 测试 ====================

func TestNewGitHub(t *testing.T) {
	tests := []struct {
		name    string
		token   string
		owner   string
		repo    string
		wantErr bool
	}{
		{
			name:    "有效配置",
			token:   "test-token",
			owner:   "owner",
			repo:    "repo",
			wantErr: false,
		},
		{
			name:    "空 token",
			token:   "",
			owner:   "owner",
			repo:    "repo",
			wantErr: true,
		},
		{
			name:    "空 owner",
			token:   "test-token",
			owner:   "",
			repo:    "repo",
			wantErr: true,
		},
		{
			name:    "空 repo",
			token:   "test-token",
			owner:   "owner",
			repo:    "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gh, err := NewGitHub(tt.token, tt.owner, tt.repo)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, gh)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, gh)
				assert.Equal(t, tt.owner, gh.owner)
				assert.Equal(t, tt.repo, gh.repo)
			}
		})
	}
}

// ==================== GetPlatformName 测试 ====================

func TestGitHub_GetPlatformName(t *testing.T) {
	gh, err := NewGitHub("test-token", "owner", "repo")
	require.NoError(t, err)

	assert.Equal(t, "github", gh.GetPlatformName())
}

// ==================== GetDefaultBranch 测试 ====================

func TestGitHub_GetDefaultBranch(t *testing.T) {
	server := setupMockGitHubServer(t, nil)
	defer server.Close()

	gh := createTestGitHubClient(t, server.URL)

	branch, err := gh.GetDefaultBranch(context.Background())
	assert.NoError(t, err)
	assert.Equal(t, "main", branch)
}

func TestGitHub_GetDefaultBranch_Error(t *testing.T) {
	server := setupMockGitHubServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/repos/owner/repo" {
			w.WriteHeader(http.StatusNotFound)
		}
	})
	defer server.Close()

	gh := createTestGitHubClient(t, server.URL)

	_, err := gh.GetDefaultBranch(context.Background())
	assert.Error(t, err)
}

func TestGitHub_GetDefaultBranch_NoDefaultBranch(t *testing.T) {
	server := setupMockGitHubServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/repos/owner/repo" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(&github.Repository{
				DefaultBranch: nil,
			})
		}
	})
	defer server.Close()

	gh := createTestGitHubClient(t, server.URL)

	_, err := gh.GetDefaultBranch(context.Background())
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "no default branch")
}

func TestGitHub_GetDefaultBranch_EmptyDefaultBranch(t *testing.T) {
	server := setupMockGitHubServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/repos/owner/repo" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			emptyBranch := ""
			json.NewEncoder(w).Encode(&github.Repository{
				DefaultBranch: &emptyBranch,
			})
		}
	})
	defer server.Close()

	gh := createTestGitHubClient(t, server.URL)

	_, err := gh.GetDefaultBranch(context.Background())
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "no default branch")
}

// ==================== CreatePullRequest 测试 ====================

func TestGitHub_CreatePullRequest(t *testing.T) {
	server := setupMockGitHubServer(t, nil)
	defer server.Close()

	gh := createTestGitHubClient(t, server.URL)

	ctx := context.Background()
	prURL, err := gh.CreatePullRequest(ctx, "Test PR", "Test body", "feature-branch", nil)

	assert.NoError(t, err)
	assert.Equal(t, "https://github.com/owner/repo/pull/123", prURL)
}

func TestGitHub_CreatePullRequest_WithTargetBranch(t *testing.T) {
	server := setupMockGitHubServer(t, nil)
	defer server.Close()

	gh := createTestGitHubClient(t, server.URL)

	ctx := context.Background()
	targetBranch := "develop"
	prURL, err := gh.CreatePullRequest(ctx, "Test PR", "Test body", "feature-branch", &targetBranch)

	assert.NoError(t, err)
	assert.Equal(t, "https://github.com/owner/repo/pull/123", prURL)
}

func TestGitHub_CreatePullRequest_Error(t *testing.T) {
	server := setupMockGitHubServer(t, func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		switch {
		case path == "/repos/owner/repo" && r.Method == http.MethodGet:
			// 获取默认分支成功
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(&github.Repository{
				DefaultBranch: github.String("main"),
			})
		case path == "/repos/owner/repo/pulls" && r.Method == http.MethodPost:
			// 创建 PR 失败
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(&github.ErrorResponse{
				Message: "Validation Failed",
			})
		}
	})
	defer server.Close()

	gh := createTestGitHubClient(t, server.URL)

	ctx := context.Background()
	_, err := gh.CreatePullRequest(ctx, "Test PR", "Test body", "feature-branch", nil)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to create PR")
}

func TestGitHub_CreatePullRequest_GetDefaultBranchFails(t *testing.T) {
	server := setupMockGitHubServer(t, func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		switch {
		case path == "/repos/owner/repo" && r.Method == http.MethodGet:
			// 获取默认分支失败，应该回退到 "main"
			w.WriteHeader(http.StatusNotFound)
		case path == "/repos/owner/repo/pulls" && r.Method == http.MethodPost:
			// 创建 PR 成功（使用默认的 "main"）
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusCreated)
			json.NewEncoder(w).Encode(&github.PullRequest{
				Number:  github.Int(123),
				HTMLURL: github.String("https://github.com/owner/repo/pull/123"),
			})
		}
	})
	defer server.Close()

	gh := createTestGitHubClient(t, server.URL)

	ctx := context.Background()
	prURL, err := gh.CreatePullRequest(ctx, "Test PR", "Test body", "feature-branch", nil)

	// 即使获取默认分支失败，也应该使用默认值 "main" 成功创建 PR
	assert.NoError(t, err)
	assert.Equal(t, "https://github.com/owner/repo/pull/123", prURL)
}

// ==================== GetPullRequestStatus 测试 ====================

func TestGitHub_GetPullRequestStatus(t *testing.T) {
	server := setupMockGitHubServer(t, nil)
	defer server.Close()

	gh := createTestGitHubClient(t, server.URL)

	ctx := context.Background()
	status, err := gh.GetPullRequestStatus(ctx, "123")

	assert.NoError(t, err)
	assert.NotNil(t, status)
	assert.Equal(t, "open", status.State)
	assert.False(t, status.Merged)
	assert.NotNil(t, status.Mergeable)
	assert.True(t, *status.Mergeable)
}

func TestGitHub_GetPullRequestStatus_WithURL(t *testing.T) {
	server := setupMockGitHubServer(t, nil)
	defer server.Close()

	gh := createTestGitHubClient(t, server.URL)

	ctx := context.Background()
	status, err := gh.GetPullRequestStatus(ctx, "https://github.com/owner/repo/pull/123")

	assert.NoError(t, err)
	assert.NotNil(t, status)
}

func TestGitHub_GetPullRequestStatus_InvalidPRID(t *testing.T) {
	server := setupMockGitHubServer(t, nil)
	defer server.Close()

	gh := createTestGitHubClient(t, server.URL)

	ctx := context.Background()
	_, err := gh.GetPullRequestStatus(ctx, "invalid")

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid PR ID")
}

func TestGitHub_ClosePullRequest_InvalidPRID(t *testing.T) {
	server := setupMockGitHubServer(t, nil)
	defer server.Close()

	gh := createTestGitHubClient(t, server.URL)

	ctx := context.Background()
	err := gh.ClosePullRequest(ctx, "invalid")

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid PR ID")
}

func TestGitHub_UpdatePullRequest_InvalidPRID(t *testing.T) {
	server := setupMockGitHubServer(t, nil)
	defer server.Close()

	gh := createTestGitHubClient(t, server.URL)

	ctx := context.Background()
	newTitle := "Updated Title"
	err := gh.UpdatePullRequest(ctx, "invalid", &newTitle, nil, nil)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid PR ID")
}

func TestGitHub_AddComment_InvalidPRID(t *testing.T) {
	server := setupMockGitHubServer(t, nil)
	defer server.Close()

	gh := createTestGitHubClient(t, server.URL)

	ctx := context.Background()
	err := gh.AddComment(ctx, "invalid", "Test comment")

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid PR ID")
}

func TestGitHub_ApprovePullRequest_InvalidPRID(t *testing.T) {
	server := setupMockGitHubServer(t, nil)
	defer server.Close()

	gh := createTestGitHubClient(t, server.URL)

	ctx := context.Background()
	err := gh.ApprovePullRequest(ctx, "invalid")

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid PR ID")
}

// ==================== ListPullRequests 测试 ====================

func TestGitHub_ListPullRequests(t *testing.T) {
	server := setupMockGitHubServer(t, nil)
	defer server.Close()

	gh := createTestGitHubClient(t, server.URL)

	ctx := context.Background()
	prs, err := gh.ListPullRequests(ctx, "open", 10)

	assert.NoError(t, err)
	assert.Len(t, prs, 1)
	assert.Equal(t, 123, prs[0].Number)
	assert.Equal(t, "Test PR 1", prs[0].Title)
	assert.Equal(t, "open", prs[0].State)
}

func TestGitHub_ListPullRequests_InvalidState(t *testing.T) {
	server := setupMockGitHubServer(t, nil)
	defer server.Close()

	gh := createTestGitHubClient(t, server.URL)

	ctx := context.Background()
	prs, err := gh.ListPullRequests(ctx, "invalid", 10)

	// 无效状态应该回退到 "open"
	assert.NoError(t, err)
	assert.NotNil(t, prs)
}

func TestGitHub_ListPullRequests_EmptyState(t *testing.T) {
	server := setupMockGitHubServer(t, nil)
	defer server.Close()

	gh := createTestGitHubClient(t, server.URL)

	ctx := context.Background()
	prs, err := gh.ListPullRequests(ctx, "", 10)

	// 空状态应该使用默认值 "open"
	assert.NoError(t, err)
	assert.NotNil(t, prs)
}

func TestGitHub_ListPullRequests_ClosedState(t *testing.T) {
	server := setupMockGitHubServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/repos/owner/repo/pulls" && r.Method == http.MethodGet {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode([]*github.PullRequest{})
		}
	})
	defer server.Close()

	gh := createTestGitHubClient(t, server.URL)

	ctx := context.Background()
	prs, err := gh.ListPullRequests(ctx, "closed", 10)

	assert.NoError(t, err)
	assert.NotNil(t, prs)
	assert.Len(t, prs, 0)
}

func TestGitHub_ListPullRequests_AllState(t *testing.T) {
	server := setupMockGitHubServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/repos/owner/repo/pulls" && r.Method == http.MethodGet {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode([]*github.PullRequest{})
		}
	})
	defer server.Close()

	gh := createTestGitHubClient(t, server.URL)

	ctx := context.Background()
	prs, err := gh.ListPullRequests(ctx, "all", 10)

	assert.NoError(t, err)
	assert.NotNil(t, prs)
}

// ==================== MergePullRequest 测试 ====================

func TestGitHub_MergePullRequest(t *testing.T) {
	server := setupMockGitHubServer(t, nil)
	defer server.Close()

	gh := createTestGitHubClient(t, server.URL)

	ctx := context.Background()
	err := gh.MergePullRequest(ctx, "123", "squash", false)

	assert.NoError(t, err)
}

func TestGitHub_MergePullRequest_WithDeleteBranch(t *testing.T) {
	server := setupMockGitHubServer(t, func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		switch {
		case path == "/repos/owner/repo/pulls/123" && r.Method == http.MethodGet:
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(&github.PullRequest{
				Number: github.Int(123),
				Head: &github.PullRequestBranch{
					Ref: github.String("feature-branch"),
				},
			})
		case path == "/repos/owner/repo/pulls/123/merge" && r.Method == http.MethodPut:
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(&github.PullRequestMergeResult{
				Merged: github.Bool(true),
			})
		case path == "/repos/owner/repo/git/refs/heads/feature-branch" && r.Method == http.MethodDelete:
			w.WriteHeader(http.StatusNoContent)
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	})
	defer server.Close()

	gh := createTestGitHubClient(t, server.URL)

	ctx := context.Background()
	err := gh.MergePullRequest(ctx, "123", "merge", true)

	assert.NoError(t, err)
}

func TestGitHub_MergePullRequest_InvalidPRID(t *testing.T) {
	server := setupMockGitHubServer(t, nil)
	defer server.Close()

	gh := createTestGitHubClient(t, server.URL)

	ctx := context.Background()
	err := gh.MergePullRequest(ctx, "invalid", "squash", false)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid PR ID")
}

func TestGitHub_MergePullRequest_InvalidMethod(t *testing.T) {
	server := setupMockGitHubServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/repos/owner/repo/pulls/123/merge" && r.Method == http.MethodPut {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(&github.PullRequestMergeResult{
				Merged: github.Bool(true),
			})
		}
	})
	defer server.Close()

	gh := createTestGitHubClient(t, server.URL)

	ctx := context.Background()
	// 使用无效的合并方法，应该回退到 "squash"
	err := gh.MergePullRequest(ctx, "123", "invalid-method", false)

	assert.NoError(t, err)
}

func TestGitHub_MergePullRequest_MergeError(t *testing.T) {
	server := setupMockGitHubServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/repos/owner/repo/pulls/123/merge" && r.Method == http.MethodPut {
			w.WriteHeader(http.StatusMethodNotAllowed)
			json.NewEncoder(w).Encode(&github.ErrorResponse{
				Message: "Pull Request is not mergeable",
			})
		}
	})
	defer server.Close()

	gh := createTestGitHubClient(t, server.URL)

	ctx := context.Background()
	err := gh.MergePullRequest(ctx, "123", "squash", false)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to merge PR")
}

func TestGitHub_MergePullRequest_GetPRInfoError(t *testing.T) {
	server := setupMockGitHubServer(t, func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		switch {
		case path == "/repos/owner/repo/pulls/123/merge" && r.Method == http.MethodPut:
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(&github.PullRequestMergeResult{
				Merged: github.Bool(true),
			})
		case path == "/repos/owner/repo/pulls/123" && r.Method == http.MethodGet:
			w.WriteHeader(http.StatusNotFound)
		}
	})
	defer server.Close()

	gh := createTestGitHubClient(t, server.URL)

	ctx := context.Background()
	err := gh.MergePullRequest(ctx, "123", "squash", true)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to get PR info")
}

func TestGitHub_MergePullRequest_DeleteBranchError(t *testing.T) {
	server := setupMockGitHubServer(t, func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		switch {
		case path == "/repos/owner/repo/pulls/123" && r.Method == http.MethodGet:
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(&github.PullRequest{
				Number: github.Int(123),
				Head: &github.PullRequestBranch{
					Ref: github.String("feature-branch"),
				},
			})
		case path == "/repos/owner/repo/pulls/123/merge" && r.Method == http.MethodPut:
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(&github.PullRequestMergeResult{
				Merged: github.Bool(true),
			})
		case path == "/repos/owner/repo/git/refs/heads/feature-branch" && r.Method == http.MethodDelete:
			w.WriteHeader(http.StatusNotFound)
		}
	})
	defer server.Close()

	gh := createTestGitHubClient(t, server.URL)

	ctx := context.Background()
	err := gh.MergePullRequest(ctx, "123", "squash", true)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to delete branch")
}

// ==================== ClosePullRequest 测试 ====================

func TestGitHub_ClosePullRequest(t *testing.T) {
	server := setupMockGitHubServer(t, nil)
	defer server.Close()

	gh := createTestGitHubClient(t, server.URL)

	ctx := context.Background()
	err := gh.ClosePullRequest(ctx, "123")

	assert.NoError(t, err)
}

// ==================== UpdatePullRequest 测试 ====================

func TestGitHub_UpdatePullRequest(t *testing.T) {
	server := setupMockGitHubServer(t, nil)
	defer server.Close()

	gh := createTestGitHubClient(t, server.URL)

	ctx := context.Background()
	newTitle := "Updated Title"
	err := gh.UpdatePullRequest(ctx, "123", &newTitle, nil, nil)

	assert.NoError(t, err)
}

func TestGitHub_UpdatePullRequest_WithState(t *testing.T) {
	server := setupMockGitHubServer(t, nil)
	defer server.Close()

	gh := createTestGitHubClient(t, server.URL)

	ctx := context.Background()
	state := "closed"
	err := gh.UpdatePullRequest(ctx, "123", nil, nil, &state)

	assert.NoError(t, err)
}

// ==================== AddComment 测试 ====================

func TestGitHub_AddComment(t *testing.T) {
	server := setupMockGitHubServer(t, nil)
	defer server.Close()

	gh := createTestGitHubClient(t, server.URL)

	ctx := context.Background()
	err := gh.AddComment(ctx, "123", "Test comment")

	assert.NoError(t, err)
}

// ==================== ApprovePullRequest 测试 ====================

func TestGitHub_ApprovePullRequest(t *testing.T) {
	server := setupMockGitHubServer(t, nil)
	defer server.Close()

	gh := createTestGitHubClient(t, server.URL)

	ctx := context.Background()
	err := gh.ApprovePullRequest(ctx, "123")

	assert.NoError(t, err)
}

func TestGitHub_ClosePullRequest_Error(t *testing.T) {
	server := setupMockGitHubServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/repos/owner/repo/pulls/123" && r.Method == http.MethodPatch {
			w.WriteHeader(http.StatusNotFound)
		}
	})
	defer server.Close()

	gh := createTestGitHubClient(t, server.URL)

	ctx := context.Background()
	err := gh.ClosePullRequest(ctx, "123")

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to close PR")
}

func TestGitHub_GetPullRequestStatus_Error(t *testing.T) {
	server := setupMockGitHubServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/repos/owner/repo/pulls/123" && r.Method == http.MethodGet {
			w.WriteHeader(http.StatusNotFound)
		}
	})
	defer server.Close()

	gh := createTestGitHubClient(t, server.URL)

	ctx := context.Background()
	_, err := gh.GetPullRequestStatus(ctx, "123")

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to get PR")
}

func TestGitHub_UpdatePullRequest_Error(t *testing.T) {
	server := setupMockGitHubServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/repos/owner/repo/pulls/123" && r.Method == http.MethodPatch {
			w.WriteHeader(http.StatusNotFound)
		}
	})
	defer server.Close()

	gh := createTestGitHubClient(t, server.URL)

	ctx := context.Background()
	newTitle := "Updated Title"
	err := gh.UpdatePullRequest(ctx, "123", &newTitle, nil, nil)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to update PR")
}

func TestGitHub_AddComment_Error(t *testing.T) {
	server := setupMockGitHubServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/repos/owner/repo/issues/123/comments" && r.Method == http.MethodPost {
			w.WriteHeader(http.StatusNotFound)
		}
	})
	defer server.Close()

	gh := createTestGitHubClient(t, server.URL)

	ctx := context.Background()
	err := gh.AddComment(ctx, "123", "Test comment")

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to add comment")
}

func TestGitHub_ApprovePullRequest_Error(t *testing.T) {
	server := setupMockGitHubServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/repos/owner/repo/pulls/123/reviews" && r.Method == http.MethodPost {
			w.WriteHeader(http.StatusForbidden)
		}
	})
	defer server.Close()

	gh := createTestGitHubClient(t, server.URL)

	ctx := context.Background()
	err := gh.ApprovePullRequest(ctx, "123")

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to approve PR")
}

func TestGitHub_ListPullRequests_Error(t *testing.T) {
	server := setupMockGitHubServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/repos/owner/repo/pulls" && r.Method == http.MethodGet {
			w.WriteHeader(http.StatusInternalServerError)
		}
	})
	defer server.Close()

	gh := createTestGitHubClient(t, server.URL)

	ctx := context.Background()
	_, err := gh.ListPullRequests(ctx, "open", 10)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to list PRs")
}


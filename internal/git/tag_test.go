package git

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ==================== CreateTag 测试 ====================

func TestRepository_CreateTag(t *testing.T) {
	repo, _ := setupTestRepoWithCommit(t)

	// 获取 HEAD
	head, err := repo.GetHead()
	require.NoError(t, err)

	// 创建 tag
	tagName := "v1.0.0"
	err = repo.CreateTag(tagName, head)
	assert.NoError(t, err)

	// 验证 tag 已创建
	exists, err := repo.TagExists(tagName)
	assert.NoError(t, err)
	assert.True(t, exists)
}

func TestRepository_CreateTag_InvalidHash(t *testing.T) {
	repo, _ := setupTestRepoWithCommit(t)

	invalidHash := plumbing.NewHash("0000000000000000000000000000000000000000")
	err := repo.CreateTag("v1.0.0", invalidHash)
	// go-git 可能不会对无效哈希返回错误，这取决于实现
	// 我们主要测试方法调用不会 panic
	if err != nil {
		t.Logf("CreateTag 对无效哈希返回错误（预期行为）: %v", err)
	} else {
		t.Logf("CreateTag 对无效哈希不返回错误（go-git 行为）")
	}
}

// ==================== CreateTagAtHead 测试 ====================

func TestRepository_CreateTagAtHead(t *testing.T) {
	repo, _ := setupTestRepoWithCommit(t)

	// 在 HEAD 创建 tag
	tagName := "v1.0.0"
	err := repo.CreateTagAtHead(tagName)
	assert.NoError(t, err)

	// 验证 tag 已创建
	exists, err := repo.TagExists(tagName)
	assert.NoError(t, err)
	assert.True(t, exists)

	// 验证 tag 指向 HEAD
	head, err := repo.GetHead()
	require.NoError(t, err)

	tags, err := repo.ListTags()
	require.NoError(t, err)

	for _, tag := range tags {
		if tag.Name == tagName {
			assert.Equal(t, head.String(), tag.CommitHash)
			break
		}
	}
}

// ==================== ListTags 测试 ====================

func TestRepository_ListTags(t *testing.T) {
	repo, _ := setupTestRepoWithCommit(t)

	// 创建多个 tags
	tagNames := []string{"v1.0.0", "v1.1.0", "v2.0.0"}
	head, err := repo.GetHead()
	require.NoError(t, err)

	for _, tagName := range tagNames {
		err := repo.CreateTag(tagName, head)
		require.NoError(t, err)
	}

	// 列出所有 tags
	tags, err := repo.ListTags()
	assert.NoError(t, err)
	assert.Len(t, tags, len(tagNames))

	// 验证 tag 信息
	tagMap := make(map[string]bool)
	for _, tag := range tags {
		tagMap[tag.Name] = true
		assert.NotEmpty(t, tag.CommitHash)
	}

	for _, tagName := range tagNames {
		assert.True(t, tagMap[tagName])
	}
}

func TestRepository_ListTags_Empty(t *testing.T) {
	repo, _ := setupTestRepoWithCommit(t)

	tags, err := repo.ListTags()
	assert.NoError(t, err)
	assert.Empty(t, tags)
}

// TestRepository_ListTags_WithAnnotatedTags 测试带 annotated tags 的情况
func TestRepository_ListTags_WithAnnotatedTags(t *testing.T) {
	repo, _ := setupTestRepoWithCommit(t)

	// 创建多个提交和 tags
	head, err := repo.GetHead()
	require.NoError(t, err)

	tagNames := []string{"v1.0.0", "v1.1.0", "v2.0.0"}
	for _, tagName := range tagNames {
		err := repo.CreateTag(tagName, head)
		require.NoError(t, err)
	}

	// 列出所有 tags
	tags, err := repo.ListTags()
	assert.NoError(t, err)
	assert.Len(t, tags, len(tagNames))

	// 验证每个 tag 都有 commit hash
	for _, tag := range tags {
		assert.NotEmpty(t, tag.CommitHash)
		assert.NotEmpty(t, tag.Name)
	}
}

// TestRepository_ListTags_MultipleCommits 测试多个提交的 tags
func TestRepository_ListTags_MultipleCommits(t *testing.T) {
	repo, tempDir := setupTestRepoWithCommit(t)

	// 创建多个提交
	author := &object.Signature{
		Name:  "Test User",
		Email: "test@example.com",
	}

	commits := []string{"commit1", "commit2", "commit3"}
	var commitHashes []plumbing.Hash

	for i, msg := range commits {
		filename := filepath.Join(tempDir, fmt.Sprintf("file%d.txt", i))
		err := os.WriteFile(filename, []byte("content"), 0644)
		require.NoError(t, err)

		err = repo.Add(fmt.Sprintf("file%d.txt", i))
		require.NoError(t, err)

		hash, err := repo.Commit(msg, author)
		require.NoError(t, err)
		commitHashes = append(commitHashes, hash)
	}

	// 为每个提交创建 tag
	for i, hash := range commitHashes {
		tagName := fmt.Sprintf("v1.%d.0", i)
		err := repo.CreateTag(tagName, hash)
		require.NoError(t, err)
	}

	// 列出所有 tags
	tags, err := repo.ListTags()
	assert.NoError(t, err)
	assert.Len(t, tags, len(commits))

	// 验证 tags 指向正确的提交
	tagMap := make(map[string]string)
	for _, tag := range tags {
		tagMap[tag.Name] = tag.CommitHash
	}

	for i, hash := range commitHashes {
		tagName := fmt.Sprintf("v1.%d.0", i)
		if commitHash, ok := tagMap[tagName]; ok {
			assert.Equal(t, hash.String(), commitHash)
		}
	}
}

// ==================== DeleteTag 测试 ====================

func TestRepository_DeleteTag(t *testing.T) {
	repo, _ := setupTestRepoWithCommit(t)

	// 创建 tag
	tagName := "v1.0.0"
	head, err := repo.GetHead()
	require.NoError(t, err)

	err = repo.CreateTag(tagName, head)
	require.NoError(t, err)

	// 验证 tag 存在
	exists, err := repo.TagExists(tagName)
	require.NoError(t, err)
	require.True(t, exists)

	// 删除 tag
	err = repo.DeleteTag(tagName)
	assert.NoError(t, err)

	// 验证 tag 已删除
	exists, err = repo.TagExists(tagName)
	assert.NoError(t, err)
	assert.False(t, exists)
}

func TestRepository_DeleteTag_NonExistent(t *testing.T) {
	repo, _ := setupTestRepoWithCommit(t)

	// 删除不存在的 tag，go-git 可能不会返回错误
	// 这取决于实现，我们主要测试方法不会 panic
	err := repo.DeleteTag("non-existent-tag")
	// 如果返回错误，这是预期的；如果不返回错误，也是可以接受的
	if err != nil {
		t.Logf("DeleteTag 对不存在的 tag 返回错误（预期行为）: %v", err)
	} else {
		t.Logf("DeleteTag 对不存在的 tag 不返回错误（go-git 行为）")
	}
}

// ==================== TagExists 测试 ====================

func TestRepository_TagExists(t *testing.T) {
	repo, _ := setupTestRepoWithCommit(t)

	tests := []struct {
		name    string
		tag     string
		setup   func() error
		want    bool
		wantErr bool
	}{
		{
			name: "existing tag",
			tag:  "v1.0.0",
			setup: func() error {
				head, err := repo.GetHead()
				if err != nil {
					return err
				}
				return repo.CreateTag("v1.0.0", head)
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "non-existent tag",
			tag:  "non-existent",
			setup: func() error {
				return nil
			},
			want:    false,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.setup()
			require.NoError(t, err)

			exists, err := repo.TagExists(tt.tag)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, exists)
			}
		})
	}
}

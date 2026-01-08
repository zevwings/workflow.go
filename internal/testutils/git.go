//go:build test

package testutils

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	gitv5 "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/stretchr/testify/require"
	"github.com/zevwings/workflow/internal/git"
)

// GitTestRepo 封装了测试 Git 仓库，提供便捷的访问方法
type GitTestRepo struct {
	repo    *git.Repository
	path    string
	tempDir string
	isBare  bool
	cleanup func()
}

// Repository 返回底层的 git.Repository 对象
func (r *GitTestRepo) Repository() *git.Repository {
	return r.repo
}

// Path 返回仓库路径
func (r *GitTestRepo) Path() string {
	return r.path
}

// TempDir 返回临时目录路径（仅对普通仓库有效）
func (r *GitTestRepo) TempDir() string {
	return r.tempDir
}

// IsBare 返回是否为 bare 仓库
func (r *GitTestRepo) IsBare() bool {
	return r.isBare
}

// Close 手动清理资源（通常不需要调用，Build 时会自动注册清理）
func (r *GitTestRepo) Close() {
	if r.cleanup != nil {
		r.cleanup()
	}
}

// GitTestRepoBuilder 用于构建 Git 测试仓库的构建器
type GitTestRepoBuilder struct {
	// 基础配置
	defaultBranch string
	bare          bool
	tempDir       string // 如果为空，将使用 t.TempDir()

	// Git 用户配置
	userName  string
	userEmail string

	// 文件操作（按顺序执行）
	files []fileOperation

	// 提交操作
	commits []commitOperation

	// 分支操作
	branches []branchOperation

	// 远程操作
	remotes []remoteOperation

	// Tag 操作
	tags []tagOperation
}

type fileOperation struct {
	path    string
	content []byte
	mode    os.FileMode
}

type commitOperation struct {
	message string
	author  *object.Signature
	files   []string // 如果为空，提交所有更改
}

type branchOperation struct {
	name   string
	create bool // true=创建并切换，false=仅创建
}

type remoteOperation struct {
	name string
	url  string
}

type tagOperation struct {
	name    string
	message string
}

// NewGitTestRepo 创建新的 Git 测试仓库构建器
func NewGitTestRepo() *GitTestRepoBuilder {
	return &GitTestRepoBuilder{
		defaultBranch: "main",
		userName:      "Test User",
		userEmail:     "test@example.com",
		files:         make([]fileOperation, 0),
		commits:       make([]commitOperation, 0),
		branches:      make([]branchOperation, 0),
		remotes:       make([]remoteOperation, 0),
		tags:          make([]tagOperation, 0),
	}
}

// WithDefaultBranch 设置默认分支
func (b *GitTestRepoBuilder) WithDefaultBranch(branch string) *GitTestRepoBuilder {
	b.defaultBranch = branch
	return b
}

// WithBare 设置为 bare 仓库（用于模拟远程仓库）
func (b *GitTestRepoBuilder) WithBare(bare bool) *GitTestRepoBuilder {
	b.bare = bare
	return b
}

// WithTempDir 设置临时目录（如果为空，将使用 t.TempDir()）
func (b *GitTestRepoBuilder) WithTempDir(dir string) *GitTestRepoBuilder {
	b.tempDir = dir
	return b
}

// WithUser 设置 Git 用户信息
func (b *GitTestRepoBuilder) WithUser(name, email string) *GitTestRepoBuilder {
	b.userName = name
	b.userEmail = email
	return b
}

// WithFile 添加文件到仓库
func (b *GitTestRepoBuilder) WithFile(path string, content []byte) *GitTestRepoBuilder {
	b.files = append(b.files, fileOperation{
		path:    path,
		content: content,
		mode:    0644,
	})
	return b
}

// WithFileString 添加文件到仓库（字符串内容）
func (b *GitTestRepoBuilder) WithFileString(path, content string) *GitTestRepoBuilder {
	return b.WithFile(path, []byte(content))
}

// WithFileMode 设置最后一个文件的权限模式
func (b *GitTestRepoBuilder) WithFileMode(mode os.FileMode) *GitTestRepoBuilder {
	if len(b.files) > 0 {
		b.files[len(b.files)-1].mode = mode
	}
	return b
}

// WithCommit 添加提交操作
func (b *GitTestRepoBuilder) WithCommit(message string) *GitTestRepoBuilder {
	return b.WithCommitWithAuthor(message, nil)
}

// WithCommitWithAuthor 添加提交操作（指定作者）
func (b *GitTestRepoBuilder) WithCommitWithAuthor(message string, author *object.Signature) *GitTestRepoBuilder {
	b.commits = append(b.commits, commitOperation{
		message: message,
		author:  author,
		files:   nil, // 提交所有更改
	})
	return b
}

// WithCommitFiles 添加提交操作（仅提交指定文件）
func (b *GitTestRepoBuilder) WithCommitFiles(message string, files ...string) *GitTestRepoBuilder {
	b.commits = append(b.commits, commitOperation{
		message: message,
		author:  nil,
		files:   files,
	})
	return b
}

// WithBranch 创建分支
func (b *GitTestRepoBuilder) WithBranch(name string) *GitTestRepoBuilder {
	b.branches = append(b.branches, branchOperation{
		name:   name,
		create: false,
	})
	return b
}

// WithBranchAndCheckout 创建并切换到分支
func (b *GitTestRepoBuilder) WithBranchAndCheckout(name string) *GitTestRepoBuilder {
	b.branches = append(b.branches, branchOperation{
		name:   name,
		create: true,
	})
	return b
}

// WithRemote 添加远程仓库
func (b *GitTestRepoBuilder) WithRemote(name, url string) *GitTestRepoBuilder {
	b.remotes = append(b.remotes, remoteOperation{
		name: name,
		url:  url,
	})
	return b
}

// WithTag 添加 Tag
func (b *GitTestRepoBuilder) WithTag(name string) *GitTestRepoBuilder {
	return b.WithTagMessage(name, "")
}

// WithTagMessage 添加 Tag（带消息）
func (b *GitTestRepoBuilder) WithTagMessage(name, message string) *GitTestRepoBuilder {
	b.tags = append(b.tags, tagOperation{
		name:    name,
		message: message,
	})
	return b
}

// Build 构建 Git 测试仓库
// 仓库会在测试结束时自动清理（使用 t.Cleanup）
func (b *GitTestRepoBuilder) Build(t *testing.T) *GitTestRepo {
	t.Helper()

	// 确定临时目录
	tempDir := b.tempDir
	if tempDir == "" {
		tempDir = t.TempDir()
	}

	var repo *git.Repository
	var err error
	var repoPath string

	if b.bare {
		// 创建 bare 仓库
		opts := &gitv5.PlainInitOptions{
			InitOptions: gitv5.InitOptions{
				DefaultBranch: plumbing.ReferenceName(fmt.Sprintf("refs/heads/%s", b.defaultBranch)),
			},
			Bare: true,
		}

		gitRepo, err := gitv5.PlainInitWithOptions(tempDir, opts)
		require.NoError(t, err, "创建 bare 仓库失败")

		// 配置 Git 用户
		config, err := gitRepo.Config()
		require.NoError(t, err)
		config.User.Name = b.userName
		config.User.Email = b.userEmail
		err = gitRepo.SetConfig(config)
		require.NoError(t, err)

		// bare 仓库不能直接使用我们的 Repository 包装
		// 但我们可以返回路径供测试使用
		repoPath = tempDir
	} else {
		// 创建普通仓库
		repo, err = git.Init(tempDir, b.defaultBranch)
		require.NoError(t, err, "初始化 Git 仓库失败")

		// 配置 Git 用户
		config, err := repo.Repo().Config()
		require.NoError(t, err)
		config.User.Name = b.userName
		config.User.Email = b.userEmail
		err = repo.Repo().SetConfig(config)
		require.NoError(t, err)

		repoPath = repo.Path()
	}

	// 创建测试仓库对象
	testRepo := &GitTestRepo{
		path:    repoPath,
		tempDir: tempDir,
		isBare:  b.bare,
	}

	if !b.bare {
		testRepo.repo = repo

		// 执行文件操作
		for _, fileOp := range b.files {
			filePath := filepath.Join(tempDir, fileOp.path)
			err := os.MkdirAll(filepath.Dir(filePath), 0755)
			require.NoError(t, err, "创建文件目录失败")

			err = os.WriteFile(filePath, fileOp.content, fileOp.mode)
			require.NoError(t, err, "写入文件失败")
		}

		// 执行提交操作
		for _, commitOp := range b.commits {
			// 添加文件到暂存区
			if len(commitOp.files) > 0 {
				// 只添加指定文件
				for _, file := range commitOp.files {
					err := repo.Add(file)
					require.NoError(t, err, "添加文件到暂存区失败")
				}
			} else {
				// 添加所有更改
				err := repo.AddAll()
				require.NoError(t, err, "添加所有文件到暂存区失败")
			}

			// 准备作者信息
			author := commitOp.author
			if author == nil {
				author = &object.Signature{
					Name:  b.userName,
					Email: b.userEmail,
					When:  time.Now(),
				}
			}

			// 提交
			_, err := repo.Commit(commitOp.message, author)
			require.NoError(t, err, "提交失败")
		}

		// 执行分支操作
		for _, branchOp := range b.branches {
			if branchOp.create {
				// 创建并切换
				err := repo.CreateAndCheckoutBranch(branchOp.name)
				require.NoError(t, err, "创建并切换分支失败")
			} else {
				// 仅创建
				err := repo.CreateBranch(branchOp.name)
				require.NoError(t, err, "创建分支失败")
			}
		}

		// 执行远程操作
		for _, remoteOp := range b.remotes {
			err := repo.AddRemote(remoteOp.name, remoteOp.url)
			require.NoError(t, err, "添加远程仓库失败")
		}

		// 执行 Tag 操作（需要先有提交）
		for _, tagOp := range b.tags {
			if tagOp.message != "" {
				// 创建 annotated tag（需要直接使用 go-git API）
				head, err := repo.GetHead()
				require.NoError(t, err, "获取 HEAD 失败，无法创建 tag")

				gitRepo := repo.Repo()
				tagObj := &object.Tag{
					Name:    tagOp.name,
					Message: tagOp.message,
					Tagger: object.Signature{
						Name:  b.userName,
						Email: b.userEmail,
						When:  time.Now(),
					},
					Target: head,
				}

				// 创建 tag 对象
				obj := gitRepo.Storer.NewEncodedObject()
				err = tagObj.Encode(obj)
				require.NoError(t, err, "编码 tag 对象失败")

				hash, err := gitRepo.Storer.SetEncodedObject(obj)
				require.NoError(t, err, "保存 tag 对象失败")

				// 创建引用
				ref := plumbing.NewHashReference(
					plumbing.ReferenceName(fmt.Sprintf("refs/tags/%s", tagOp.name)),
					hash,
				)
				err = gitRepo.Storer.SetReference(ref)
				require.NoError(t, err, "创建 tag 引用失败")
			} else {
				// 创建 lightweight tag（使用我们的 Repository API）
				err := repo.CreateTagAtHead(tagOp.name)
				require.NoError(t, err, "创建 tag 失败")
			}
		}
	}

	// 注册清理函数（bare 仓库由 t.TempDir() 自动清理，普通仓库也是）
	testRepo.cleanup = func() {
		// 实际上，由于使用了 t.TempDir()，清理是自动的
		// 这里保留接口，以便将来可能需要手动清理
	}

	// 使用 t.Cleanup 确保资源清理
	t.Cleanup(func() {
		if testRepo.cleanup != nil {
			testRepo.cleanup()
		}
	})

	return testRepo
}

// NewBareGitTestRepo 创建新的 bare Git 测试仓库构建器（用于模拟远程仓库）
func NewBareGitTestRepo() *GitTestRepoBuilder {
	return NewGitTestRepo().WithBare(true)
}

// SetupTestRepo 在临时目录中创建并初始化一个 Git 仓库
// 这是对 GitTestRepoBuilder 的便捷包装，提供与内部测试函数相同的接口
//
// 返回:
//   - repo: 初始化的 Repository 实例
//   - tempDir: 临时目录路径
//
// 示例:
//
//	repo, tempDir := testutils.SetupTestRepo(t)
func SetupTestRepo(t *testing.T) (*git.Repository, string) {
	t.Helper()
	testRepo := NewGitTestRepo().
		WithDefaultBranch("main").
		Build(t)
	return testRepo.Repository(), testRepo.Path()
}

// SetupTestRepoWithCommit 创建带有一个初始提交的测试仓库
// 这是对 GitTestRepoBuilder 的便捷包装，提供与内部测试函数相同的接口
//
// 返回:
//   - repo: 初始化的 Repository 实例
//   - tempDir: 临时目录路径
//
// 示例:
//
//	repo, tempDir := testutils.SetupTestRepoWithCommit(t)
func SetupTestRepoWithCommit(t *testing.T) (*git.Repository, string) {
	t.Helper()
	testRepo := NewGitTestRepo().
		WithDefaultBranch("main").
		WithFileString("test.txt", "test content").
		WithCommit("Initial commit").
		Build(t)
	return testRepo.Repository(), testRepo.Path()
}

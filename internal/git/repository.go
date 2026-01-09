package git

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/zevwings/workflow/internal/logging"
)

// Repository Git 仓库封装
type Repository struct {
	repo     *git.Repository
	worktree *git.Worktree
	path     string
}

// Open 打开指定路径的 Git 仓库
func Open(path string) (*Repository, error) {
	logger := logging.GetLogger()
	logger.Debugf("Opening Git repository: %s", path)

	absPath, err := filepath.Abs(path)
	if err != nil {
		logger.WithError(err).WithField("path", path).Error("Failed to resolve absolute path")
		return nil, fmt.Errorf("failed to get absolute path: %w", err)
	}

	repo, err := git.PlainOpen(absPath)
	if err != nil {
		logger.WithError(err).WithField("path", absPath).Error("Failed to open Git repository")
		return nil, fmt.Errorf("failed to open repository: %w", err)
	}

	worktree, err := repo.Worktree()
	if err != nil {
		logger.WithError(err).WithField("path", absPath).Error("Failed to get worktree")
		return nil, fmt.Errorf("failed to get worktree: %w", err)
	}

	logger.WithField("path", absPath).Debug("Git repository opened successfully")

	return &Repository{
		repo:     repo,
		worktree: worktree,
		path:     absPath,
	}, nil
}

// OpenCurrent 打开当前目录的 Git 仓库
func OpenCurrent() (*Repository, error) {
	wd, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("failed to get current directory: %w", err)
	}
	return Open(wd)
}

// Init 初始化一个新的 Git 仓库
func Init(path string, initialBranch string) (*Repository, error) {
	logger := logging.GetLogger()
	logger.WithFields(logging.Fields{
		"path":           path,
		"initial_branch": initialBranch,
	}).Info("Initializing new Git repository")

	absPath, err := filepath.Abs(path)
	if err != nil {
		return nil, fmt.Errorf("failed to get absolute path: %w", err)
	}

	opts := &git.PlainInitOptions{
		InitOptions: git.InitOptions{
			DefaultBranch: plumbing.ReferenceName(fmt.Sprintf("refs/heads/%s", initialBranch)),
		},
	}

	repo, err := git.PlainInitWithOptions(absPath, opts)
	if err != nil {
		logger.WithError(err).WithFields(logging.Fields{
			"path":           absPath,
			"initial_branch": initialBranch,
		}).Error("Failed to initialize Git repository")
		return nil, fmt.Errorf("failed to initialize repository: %w", err)
	}

	worktree, err := repo.Worktree()
	if err != nil {
		return nil, fmt.Errorf("failed to get worktree: %w", err)
	}

	logger.WithField("path", absPath).Info("Git repository initialized successfully")

	return &Repository{
		repo:     repo,
		worktree: worktree,
		path:     absPath,
	}, nil
}

// IsGitRepo 检查指定路径是否是 Git 仓库
func IsGitRepo(path string) bool {
	_, err := Open(path)
	return err == nil
}

// Path 返回仓库路径
func (r *Repository) Path() string {
	return r.path
}

// Repo 返回底层的 git.Repository 对象（用于高级操作）
func (r *Repository) Repo() *git.Repository {
	return r.repo
}

// Worktree 返回底层的 git.Worktree 对象（用于高级操作）
func (r *Repository) Worktree() *git.Worktree {
	return r.worktree
}

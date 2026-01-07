package git

import (
	"fmt"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/plumbing/storer"
)

// Status 获取工作区状态
func (r *Repository) Status() (*StatusInfo, error) {
	status, err := r.worktree.Status()
	if err != nil {
		return nil, fmt.Errorf("failed to get status: %w", err)
	}

	info := &StatusInfo{
		ModifiedFiles:  []string{},
		StagedFiles:    []string{},
		UntrackedFiles: []string{},
	}

	for file, fileStatus := range status {
		if fileStatus.Staging == git.Untracked || fileStatus.Worktree == git.Untracked {
			info.UntrackedFiles = append(info.UntrackedFiles, file)
		} else {
			if fileStatus.Staging != git.Unmodified {
				info.StagedFiles = append(info.StagedFiles, file)
			}
			if fileStatus.Worktree != git.Unmodified {
				info.ModifiedFiles = append(info.ModifiedFiles, file)
			}
		}
	}

	return info, nil
}

// HasChanges 检查是否有未提交的更改
func (r *Repository) HasChanges() (bool, error) {
	status, err := r.Status()
	if err != nil {
		return false, err
	}

	return len(status.ModifiedFiles) > 0 ||
		len(status.StagedFiles) > 0 ||
		len(status.UntrackedFiles) > 0, nil
}

// Add 添加文件到暂存区
func (r *Repository) Add(path string) error {
	_, err := r.worktree.Add(path)
	if err != nil {
		return fmt.Errorf("failed to add file %s: %w", path, err)
	}
	return nil
}

// AddAll 添加所有文件到暂存区
func (r *Repository) AddAll() error {
	_, err := r.worktree.Add(".")
	if err != nil {
		return fmt.Errorf("failed to add all files: %w", err)
	}
	return nil
}

// Commit 提交更改
func (r *Repository) Commit(message string, author *object.Signature) (plumbing.Hash, error) {
	if author == nil {
		// 尝试从配置获取作者信息
		config, err := r.repo.Config()
		if err == nil {
			author = &object.Signature{
				Name:  config.User.Name,
				Email: config.User.Email,
				When:  time.Now(),
			}
		} else {
			author = &object.Signature{
				Name:  "Unknown",
				Email: "unknown@example.com",
				When:  time.Now(),
			}
		}
	}

	hash, err := r.worktree.Commit(message, &git.CommitOptions{
		Author: author,
	})
	if err != nil {
		return plumbing.ZeroHash, fmt.Errorf("failed to commit: %w", err)
	}

	return hash, nil
}

// GetHead 获取 HEAD 的提交哈希
func (r *Repository) GetHead() (plumbing.Hash, error) {
	ref, err := r.repo.Head()
	if err != nil {
		return plumbing.ZeroHash, fmt.Errorf("failed to get HEAD: %w", err)
	}
	return ref.Hash(), nil
}

// GetCommit 获取指定哈希的提交信息
func (r *Repository) GetCommit(hash plumbing.Hash) (*CommitInfo, error) {
	commit, err := r.repo.CommitObject(hash)
	if err != nil {
		return nil, fmt.Errorf("failed to get commit: %w", err)
	}

	return &CommitInfo{
		Hash:    hash.String(),
		Message: commit.Message,
		Author:  fmt.Sprintf("%s <%s>", commit.Author.Name, commit.Author.Email),
		Date:    commit.Author.When.Format(time.RFC3339),
	}, nil
}

// GetLastCommit 获取最后一次提交信息
func (r *Repository) GetLastCommit() (*CommitInfo, error) {
	hash, err := r.GetHead()
	if err != nil {
		return nil, err
	}
	return r.GetCommit(hash)
}

// Log 获取提交历史
func (r *Repository) Log(from plumbing.Hash, limit int) ([]CommitInfo, error) {
	commitIter, err := r.repo.Log(&git.LogOptions{
		From:  from,
		Order: git.LogOrderCommitterTime,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get log: %w", err)
	}

	commits := []CommitInfo{}
	count := 0
	err = commitIter.ForEach(func(commit *object.Commit) error {
		if limit > 0 && count >= limit {
			return storer.ErrStop
		}

		commits = append(commits, CommitInfo{
			Hash:    commit.Hash.String(),
			Message: commit.Message,
			Author:  fmt.Sprintf("%s <%s>", commit.Author.Name, commit.Author.Email),
			Date:    commit.Author.When.Format(time.RFC3339),
		})
		count++
		return nil
	})
	if err != nil && err != storer.ErrStop {
		return nil, fmt.Errorf("failed to iterate commits: %w", err)
	}

	return commits, nil
}

// ResolveRevision 解析引用为提交哈希
func (r *Repository) ResolveRevision(rev string) (plumbing.Hash, error) {
	hash, err := r.repo.ResolveRevision(plumbing.Revision(rev))
	if err != nil {
		return plumbing.ZeroHash, fmt.Errorf("failed to resolve revision %s: %w", rev, err)
	}
	return *hash, nil
}


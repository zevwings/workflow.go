package git

import (
	"fmt"
	"strings"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
)

// CurrentBranch 获取当前分支名
func (r *Repository) CurrentBranch() (string, error) {
	ref, err := r.repo.Head()
	if err != nil {
		return "", fmt.Errorf("failed to get HEAD: %w", err)
	}

	if !ref.Name().IsBranch() {
		return "", fmt.Errorf("HEAD is not pointing to a branch (detached HEAD)")
	}

	return ref.Name().Short(), nil
}

// CreateBranch 创建新分支
func (r *Repository) CreateBranch(name string) error {
	headRef, err := r.repo.Head()
	if err != nil {
		return fmt.Errorf("failed to get HEAD: %w", err)
	}

	branchRef := plumbing.NewHashReference(
		plumbing.ReferenceName(fmt.Sprintf("refs/heads/%s", name)),
		headRef.Hash(),
	)

	err = r.repo.Storer.SetReference(branchRef)
	if err != nil {
		return fmt.Errorf("failed to create branch: %w", err)
	}

	return nil
}

// CheckoutBranch 切换到指定分支
func (r *Repository) CheckoutBranch(name string) error {
	branchRef := plumbing.ReferenceName(fmt.Sprintf("refs/heads/%s", name))
	err := r.worktree.Checkout(&git.CheckoutOptions{
		Branch: branchRef,
		Force:  false,
	})
	if err != nil {
		return fmt.Errorf("failed to checkout branch %s: %w", name, err)
	}
	return nil
}

// CreateAndCheckoutBranch 创建并切换到新分支
func (r *Repository) CreateAndCheckoutBranch(name string) error {
	headRef, err := r.repo.Head()
	if err != nil {
		return fmt.Errorf("failed to get HEAD: %w", err)
	}

	branchRef := plumbing.ReferenceName(fmt.Sprintf("refs/heads/%s", name))
	err = r.worktree.Checkout(&git.CheckoutOptions{
		Branch: branchRef,
		Create: true,
		Hash:   headRef.Hash(),
	})
	if err != nil {
		return fmt.Errorf("failed to create and checkout branch %s: %w", name, err)
	}

	return nil
}

// DeleteBranch 删除本地分支
func (r *Repository) DeleteBranch(name string) error {
	branchRef := plumbing.ReferenceName(fmt.Sprintf("refs/heads/%s", name))

	// 检查是否是当前分支
	currentBranch, err := r.CurrentBranch()
	if err == nil && currentBranch == name {
		return fmt.Errorf("cannot delete current branch: %s", name)
	}

	err = r.repo.Storer.RemoveReference(branchRef)
	if err != nil {
		return fmt.Errorf("failed to delete branch %s: %w", name, err)
	}

	return nil
}

// ListBranches 列出所有本地分支
func (r *Repository) ListBranches() ([]BranchInfo, error) {
	branches := []BranchInfo{}
	currentBranch, _ := r.CurrentBranch()

	iter, err := r.repo.Branches()
	if err != nil {
		return nil, fmt.Errorf("failed to list branches: %w", err)
	}

	err = iter.ForEach(func(ref *plumbing.Reference) error {
		if ref.Name().IsBranch() {
			branches = append(branches, BranchInfo{
				Name:   ref.Name().Short(),
				IsHead: ref.Name().Short() == currentBranch,
			})
		}
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("failed to iterate branches: %w", err)
	}

	return branches, nil
}

// BranchExists 检查分支是否存在
func (r *Repository) BranchExists(name string) (bool, error) {
	branchRef := plumbing.ReferenceName(fmt.Sprintf("refs/heads/%s", name))
	_, err := r.repo.Storer.Reference(branchRef)
	if err == nil {
		return true, nil
	}
	if err == plumbing.ErrReferenceNotFound {
		return false, nil
	}
	return false, fmt.Errorf("failed to check branch existence: %w", err)
}

// GetDefaultBranch 获取默认分支名
func (r *Repository) GetDefaultBranch() (string, error) {
	// 尝试从远程获取默认分支
	remotes, err := r.repo.Remotes()
	if err == nil && len(remotes) > 0 {
		remote := remotes[0]
		refs, err := remote.List(&git.ListOptions{})
		if err == nil {
			for _, ref := range refs {
				if ref.Name() == plumbing.HEAD {
					// 解析符号引用
					target := ref.Target()
					if target.IsBranch() {
						return target.Short(), nil
					}
				}
			}
		}
	}

	// 回退到本地 HEAD
	headRef, err := r.repo.Head()
	if err != nil {
		return "", fmt.Errorf("failed to get default branch: %w", err)
	}

	if headRef.Name().IsBranch() {
		return headRef.Name().Short(), nil
	}

	// 尝试常见分支名
	commonBranches := []string{"main", "master", "develop", "dev"}
	for _, branch := range commonBranches {
		exists, _ := r.BranchExists(branch)
		if exists {
			return branch, nil
		}
	}

	return "", fmt.Errorf("failed to determine default branch")
}

// ListRemoteBranches 列出所有远程分支
func (r *Repository) ListRemoteBranches(remoteName string) ([]string, error) {
	remote, err := r.repo.Remote(remoteName)
	if err != nil {
		return nil, fmt.Errorf("failed to get remote %s: %w", remoteName, err)
	}

	refs, err := remote.List(&git.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to list remote refs: %w", err)
	}

	branches := []string{}
	for _, ref := range refs {
		if ref.Name().IsBranch() {
			// 移除远程前缀，例如 "refs/remotes/origin/main" -> "main"
			branchName := strings.TrimPrefix(ref.Name().String(), fmt.Sprintf("refs/remotes/%s/", remoteName))
			if branchName != "" {
				branches = append(branches, branchName)
			}
		}
	}

	return branches, nil
}


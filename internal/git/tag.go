package git

import (
	"fmt"

	"github.com/go-git/go-git/v5/plumbing"
)

// CreateTag 创建 tag
func (r *Repository) CreateTag(name string, commitHash plumbing.Hash) error {
	_, err := r.repo.CreateTag(name, commitHash, nil)
	if err != nil {
		return fmt.Errorf("failed to create tag %s: %w", name, err)
	}
	return nil
}

// CreateTagAtHead 在当前 HEAD 创建 tag
func (r *Repository) CreateTagAtHead(name string) error {
	hash, err := r.GetHead()
	if err != nil {
		return err
	}
	return r.CreateTag(name, hash)
}

// ListTags 列出所有 tag
func (r *Repository) ListTags() ([]TagInfo, error) {
	tags := []TagInfo{}

	iter, err := r.repo.Tags()
	if err != nil {
		return nil, fmt.Errorf("failed to list tags: %w", err)
	}

	err = iter.ForEach(func(ref *plumbing.Reference) error {
		tagInfo := TagInfo{
			Name: ref.Name().Short(),
		}

		// 获取 tag 指向的 commit hash
		if ref.Hash() != plumbing.ZeroHash {
			tagInfo.CommitHash = ref.Hash().String()
		} else {
			// 如果是 annotated tag，需要解析
			tagObj, err := r.repo.TagObject(ref.Hash())
			if err == nil {
				tagInfo.CommitHash = tagObj.Target.String()
			}
		}

		tags = append(tags, tagInfo)
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("failed to iterate tags: %w", err)
	}

	return tags, nil
}

// DeleteTag 删除本地 tag
func (r *Repository) DeleteTag(name string) error {
	tagRef := plumbing.ReferenceName(fmt.Sprintf("refs/tags/%s", name))
	err := r.repo.Storer.RemoveReference(tagRef)
	if err != nil {
		return fmt.Errorf("failed to delete tag %s: %w", name, err)
	}
	return nil
}

// TagExists 检查 tag 是否存在
func (r *Repository) TagExists(name string) (bool, error) {
	tagRef := plumbing.ReferenceName(fmt.Sprintf("refs/tags/%s", name))
	_, err := r.repo.Storer.Reference(tagRef)
	if err == nil {
		return true, nil
	}
	if err == plumbing.ErrReferenceNotFound {
		return false, nil
	}
	return false, fmt.Errorf("failed to check tag existence: %w", err)
}


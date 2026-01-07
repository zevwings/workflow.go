package git

import (
	"github.com/go-git/go-git/v5/plumbing"
)

// BranchInfo 分支信息
type BranchInfo struct {
	Name   string
	IsHead bool
}

// CommitInfo 提交信息
type CommitInfo struct {
	Hash    string
	Message string
	Author  string
	Date    string
}

// TagInfo Tag 信息
type TagInfo struct {
	Name       string
	CommitHash string
}

// RemoteInfo 远程仓库信息
type RemoteInfo struct {
	Name string
	URL  string
}

// StatusInfo 工作区状态信息
type StatusInfo struct {
	ModifiedFiles  []string
	StagedFiles    []string
	UntrackedFiles []string
}

// ReferenceName 引用名称类型
type ReferenceName = plumbing.ReferenceName

// Hash 哈希类型
type Hash = plumbing.Hash

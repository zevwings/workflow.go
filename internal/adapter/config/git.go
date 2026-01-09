package config

import (
	"github.com/zevwings/workflow/internal/config"
	"github.com/zevwings/workflow/internal/git"
)

// gitAdapterWrapper 包装 git 适配器，实现 config.GitRepository 接口
//
// 由于接口定义在 config 包中，git 包无法直接实现接口。
// 这里创建一个包装器，将 git 适配器的方法调用转发到 config 接口。
type gitAdapterWrapper struct {
	impl *git.GitAdapter
}

func (w *gitAdapterWrapper) GetRepoPath() string {
	return w.impl.GetRepoPath()
}

func (w *gitAdapterWrapper) IsGitRepo(path string) bool {
	return w.impl.IsGitRepo(path)
}

func (w *gitAdapterWrapper) Open(path string) (config.GitRepo, error) {
	repo, err := w.impl.Open(path)
	if err != nil {
		return nil, err
	}
	// GitRepoAdapter 实现了 GetRemoteURL 方法，满足 config.GitRepo 接口
	return &gitRepoWrapper{impl: repo}, nil
}

// gitRepoWrapper 包装 git repo 适配器，实现 config.GitRepo 接口
type gitRepoWrapper struct {
	impl *git.GitRepoAdapter
}

func (w *gitRepoWrapper) GetRemoteURL(name string) (string, error) {
	return w.impl.GetRemoteURL(name)
}

// NewRepoManagerWithDefaultGit 创建仓库配置管理器（便捷函数）
//
// 此函数使用默认的 git 模块实现，内部调用 config.GlobalRepoManager。
// 将 Git 适配器的创建移到 adapter 层，避免 config 包依赖 git 包。
//
// 参数:
//   - repoPath: 仓库根目录路径（如果为空，使用当前目录）
//
// 返回:
//   - *config.RepoManager: 仓库配置管理器实例
//   - error: 如果创建失败，返回错误
func NewRepoManagerWithDefaultGit(repoPath string) (*config.RepoManager, error) {
	// 使用默认的 git 适配器，通过包装器实现接口
	gitAdapter := &gitAdapterWrapper{impl: git.NewGitAdapter(repoPath)}
	return config.GlobalRepoManager(gitAdapter)
}

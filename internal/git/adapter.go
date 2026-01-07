package git

// GitAdapter Git 适配器
//
// 这是 git 模块对 config 模块接口的适配器实现。
// 用于解耦 config 模块对 git 模块的直接依赖。
//
// 注意：接口定义在 config 包中，git 包无法直接引用。
// 这里提供实现，调用方（commands 层）负责创建包装器实现接口。
type GitAdapter struct {
	impl     *gitAdapterImpl
	repoPath string
}

// NewGitAdapter 创建新的 Git 适配器
//
// 参数:
//   - repoPath: 仓库根目录路径（如果为空，使用当前目录）
//
// 返回一个 GitAdapter 实例，调用方需要创建包装器来实现 config.GitRepository 接口。
func NewGitAdapter(repoPath string) *GitAdapter {
	return &GitAdapter{
		impl:     &gitAdapterImpl{},
		repoPath: repoPath,
	}
}

// GetRepoPath 获取仓库路径
func (g *GitAdapter) GetRepoPath() string {
	return g.repoPath
}

// IsGitRepo 检查指定路径是否是 Git 仓库
func (g *GitAdapter) IsGitRepo(path string) bool {
	return g.impl.IsGitRepo(path)
}

// Open 打开指定路径的 Git 仓库
//
// 返回一个 GitRepoAdapter 实例，调用方需要创建包装器来实现 config.GitRepo 接口。
func (g *GitAdapter) Open(path string) (*GitRepoAdapter, error) {
	return g.impl.Open(path)
}

// gitAdapterImpl 内部实现
type gitAdapterImpl struct{}

// IsGitRepo 检查指定路径是否是 Git 仓库
func (g *gitAdapterImpl) IsGitRepo(path string) bool {
	return IsGitRepo(path)
}

// Open 打开指定路径的 Git 仓库
func (g *gitAdapterImpl) Open(path string) (*GitRepoAdapter, error) {
	repo, err := Open(path)
	if err != nil {
		return nil, err
	}
	return &GitRepoAdapter{repo: repo}, nil
}

// GitRepoAdapter Git 仓库适配器
//
// 实现了与 config.GitRepo 相同的方法签名。
type GitRepoAdapter struct {
	repo *Repository
}

// GetRemoteURL 获取指定名称的 remote URL
func (g *GitRepoAdapter) GetRemoteURL(name string) (string, error) {
	return g.repo.GetRemoteURL(name)
}

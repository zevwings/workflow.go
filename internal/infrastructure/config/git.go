package config

import (
	"github.com/zevwings/workflow/internal/config"
	"github.com/zevwings/workflow/internal/git"
)

// gitAdapterWrapper wraps git adapter to implement config.GitRepository interface
//
// Since the interface is defined in the config package, the git package cannot directly implement the interface.
// Here we create a wrapper that forwards git adapter method calls to the config interface.
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
	// GitRepoAdapter implements GetRemoteURL method, satisfying config.GitRepo interface
	return &gitRepoWrapper{impl: repo}, nil
}

// gitRepoWrapper wraps git repo adapter to implement config.GitRepo interface
type gitRepoWrapper struct {
	impl *git.GitRepoAdapter
}

func (w *gitRepoWrapper) GetRemoteURL(name string) (string, error) {
	return w.impl.GetRemoteURL(name)
}

// NewRepoManagerWithDefaultGit creates repository configuration manager (convenience function)
//
// This function uses the default git module implementation, internally calling config.GlobalRepoManager.
// Moves Git adapter creation to infrastructure layer to avoid config package depending on git package.
//
// Parameters:
//   - repoPath: Repository root directory path (if empty, use current directory)
//
// Returns:
//   - *config.RepoManager: Repository configuration manager instance
//   - error: Returns error if creation fails
func NewRepoManagerWithDefaultGit(repoPath string) (*config.RepoManager, error) {
	// Use default git adapter, implement interface through wrapper
	gitAdapter := &gitAdapterWrapper{impl: git.NewGitAdapter(repoPath)}
	return config.GlobalRepoManager(gitAdapter)
}

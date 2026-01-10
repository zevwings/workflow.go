package repo

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/zevwings/workflow/internal/git"
	infrastructureconfig "github.com/zevwings/workflow/internal/infrastructure/config"
	"github.com/zevwings/workflow/internal/prompt"
)

// NewShowCmd creates the repo show command
func NewShowCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "show",
		Short: "Display current repository configuration",
		Long:  `Display current repository configuration including branch settings and templates.`,
		RunE:  runShow,
	}

	return cmd
}

func runShow(cmd *cobra.Command, args []string) error {
	msg := prompt.GetMessage()

	msg.Info("Repository Configuration")
	msg.Break()

	// 1. Get repository name
	gitRepo, err := git.OpenCurrent()
	if err != nil {
		return fmt.Errorf("不在 Git 仓库中: %w", err)
	}

	repoName, err := getRepoName(gitRepo)
	if err != nil {
		return fmt.Errorf("获取仓库名称失败: %w", err)
	}

	msg.Info("Repository: %s", repoName)
	msg.Break()

	// 2. Show branch configuration
	msg.Info("Branch Configuration")
	msg.Break('-', 40)

	manager, err := infrastructureconfig.NewRepoManagerWithDefaultGit("")
	if err != nil {
		return fmt.Errorf("初始化配置管理器失败: %w", err)
	}

	// Load from personal preference config
	prefix := manager.GetBranchPrefix()

	if prefix != "" {
		msg.Info("Prefix: %s (personal preference)", prefix)
	} else {
		msg.Info("Prefix: (not set)")
		msg.Info("Run 'workflow repo setup' to configure branch prefix")
	}

	// 3. Show template configuration
	msg.Break()
	msg.Info("Template Configuration")
	msg.Break('-', 40)

	// Load template config
	if err := manager.Load(); err != nil {
		msg.Debug("Failed to load template config, using defaults")
	}

	templateConfig := manager.GetTemplateConfig()

	// Commit template
	msg.Info("Commit Template:")
	useScope := false
	if scope, ok := templateConfig.Commit["use_scope"].(bool); ok {
		useScope = scope
	}
	msg.Info("  Use scope: %v", useScope)

	commitTemplate := "default"
	if tmpl, ok := templateConfig.Commit["default"].(string); ok && tmpl != "" {
		commitTemplate = tmpl
	}
	msg.Info("  Template: %s", commitTemplate)

	// Branch template
	msg.Break()
	msg.Info("Branch Template:")
	branchTemplate := "default"
	if tmpl, ok := templateConfig.Branch["default"].(string); ok && tmpl != "" {
		branchTemplate = tmpl
	}
	msg.Info("  Default: %s", branchTemplate)

	// Pull request template
	msg.Break()
	msg.Info("Pull Request Template:")
	prTemplate := "default"
	if tmpl, ok := templateConfig.PullRequests["default"].(string); ok && tmpl != "" {
		prTemplate = tmpl
	}
	msg.Info("  Template: %s", prTemplate)

	return nil
}

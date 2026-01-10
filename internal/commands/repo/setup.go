package repo

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/zevwings/workflow/internal/config"
	"github.com/zevwings/workflow/internal/git"
	infrastructureconfig "github.com/zevwings/workflow/internal/infrastructure/config"
	"github.com/zevwings/workflow/internal/prompt"
	"github.com/zevwings/workflow/internal/prompt/form"
)

// NewSetupCmd creates the repo setup command
func NewSetupCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "setup",
		Short: "Interactively initialize repository-level configuration",
		Long: `Interactively initialize repository-level configuration.

This command helps you configure:
- Branch prefix (personal preference)
- Commit templates (project standard)
- Branch templates (project standard)
- Pull request templates (project standard)
- Auto-accept change type in PR creation (personal preference)

Project standard configurations are saved to .workflow/config.toml (can be committed to Git).
Personal preference configurations are saved to your global config directory (not committed).`,
		RunE: runSetup,
	}

	return cmd
}

// Ensure ensures repository configuration exists
//
// This function should be called at the beginning of branch/commit/pr operations.
// Checks if repo setup has been completed.
//
// If configuration doesn't exist, it will:
// 1. Check if in interactive environment
// 2. Prompt user to run setup
// 3. Run setup automatically if user confirms
//
// Returns error only if setup is required and fails.
func Ensure() error {
	// 1. Check if in interactive environment
	if !isTerminal() {
		return nil // Non-interactive environment, skip check
	}

	// 2. Check if configuration exists
	exists, err := checkConfigExists()
	if err != nil {
		return fmt.Errorf("检查配置失败: %w", err)
	}
	if exists {
		return nil // Configuration exists, no need to setup
	}

	// 3. Configuration doesn't exist or is incomplete
	msg := prompt.GetMessage()
	msg.Break()
	msg.Warning("Repository configuration not found or incomplete.")
	msg.Info("Project-level configuration helps:")
	msg.Print("  - Share branch prefix and commit template settings with your team")
	msg.Print("  - Automatically configure commit message format")
	msg.Print("  - Manage ignored branches")
	msg.Break()

	// 4. Ask user if they want to run setup
	shouldSetup, err := prompt.Confirm().
		Prompt("Run 'workflow repo setup' to configure this repository?").
		Default(true).
		Run()
	if err != nil {
		return fmt.Errorf("获取用户确认失败: %w", err)
	}

	if shouldSetup {
		// 5. Run setup
		msg.Break()
		msg.Info("Running repository setup...")
		msg.Break()

		if err := Run(); err != nil {
			return fmt.Errorf("运行仓库设置失败: %w", err)
		}

		msg.Break()
		msg.Success("Repository configuration completed!")
		msg.Break()
	} else {
		msg.Info("Skipping repository setup. You can run 'workflow repo setup' later.")
	}

	return nil
}

// Run runs repository setup
//
// This method can be called:
// 1. Directly by users: `workflow repo setup`
// 2. By other commands: `repo.Run()`
func Run() error {
	msg := prompt.GetMessage()

	// 1. Check if in Git repository
	gitRepo, err := git.OpenCurrent()
	if err != nil {
		return fmt.Errorf("不在 Git 仓库中，请在 Git 仓库中运行此命令: %w", err)
	}

	// Get repository name
	repoName, err := getRepoName(gitRepo)
	if err != nil {
		return fmt.Errorf("获取仓库名称失败: %w", err)
	}

	msg.Info("Repository: %s", repoName)
	msg.Break()

	// 2. Load existing configuration (if exists)
	manager, err := infrastructureconfig.NewRepoManagerWithDefaultGit("")
	if err != nil {
		return fmt.Errorf("初始化配置管理器失败: %w", err)
	}

	// Load existing config
	if err := manager.Load(); err != nil {
		// If load fails, continue with empty config
		msg.Debug("Failed to load existing config, starting with empty config")
	}

	// 3. Collect configuration information
	config, err := collectConfig(manager)
	if err != nil {
		return fmt.Errorf("收集配置失败: %w", err)
	}

	// 4. Save configuration
	if err := saveConfig(manager, config); err != nil {
		return fmt.Errorf("保存配置失败: %w", err)
	}

	msg.Break()
	msg.Success("Repository configuration saved successfully!")
	msg.Debug("Project template configuration: %s", manager.GetPublicConfigPath())
	msg.Debug("Personal preference configuration: %s", manager.GetPrivateConfigPath())
	msg.Success("You can commit the project template configuration to Git to share with your team.")

	return nil
}

func runSetup(cmd *cobra.Command, args []string) error {
	return Run()
}

// collectConfig collects configuration interactively
func collectConfig(manager *config.RepoManager) (*repoConfig, error) {
	msg := prompt.GetMessage()

	// Prepare existing values
	currentPrefix := manager.GetBranchPrefix()
	currentAutoAccept := manager.GetAutoAcceptChangeType()

	// Get existing template config
	templateConfig := manager.GetTemplateConfig()
	currentUseScope := false
	if useScope, ok := templateConfig.Commit["use_scope"].(bool); ok {
		currentUseScope = useScope
	}

	// Group 1: Personal Preference Configuration
	msg.Info("Personal Preference Configuration")
	msg.Break('-', 40)
	msg.Debug("These settings are personal preferences and will be saved to your global config.")
	msg.Break()

	personalForm := prompt.Form().
		SetTitle("Personal Preference Configuration").
		AddInput(form.InputFormField{
			Key:          "branch_prefix",
			Prompt:       getBranchPrefixPrompt(currentPrefix),
			DefaultValue: currentPrefix,
			Validator:    nil,
		}).
		AddConfirm(form.ConfirmFormField{
			Key:          "auto_accept_change_type",
			Prompt:       "Auto-accept auto-selected change type in PR creation? (skip confirmation prompt)",
			DefaultValue: currentAutoAccept,
		})

	personalResult, err := personalForm.Run()
	if err != nil {
		return nil, fmt.Errorf("收集个人偏好配置失败: %w", err)
	}

	// Group 2: Project Template Configuration
	msg.Break()
	msg.Info("Project Template Configuration")
	msg.Break('-', 40)
	msg.Debug("These settings are project standards and will be saved to .workflow/config.toml (can be committed to Git).")
	msg.Break()

	projectForm := prompt.Form().
		SetTitle("Project Template Configuration").
		AddConfirm(form.ConfirmFormField{
			Key:          "use_scope",
			Prompt:       "Use scope for commit messages?",
			DefaultValue: currentUseScope,
		}).
		AddConfirm(form.ConfirmFormField{
			Key:          "configure_commit_template",
			Prompt:       "Configure commit templates?",
			DefaultValue: false,
		}).
		AddInput(form.InputFormField{
			Key:          "custom_commit_template",
			Prompt:       "Enter custom commit template:",
			DefaultValue: "",
			Validator:    nil,
			Condition: func(r *form.FormResult) bool {
				return r.GetBool("configure_commit_template")
			},
		}).
		AddConfirm(form.ConfirmFormField{
			Key:          "configure_branch_template",
			Prompt:       "Configure branch templates?",
			DefaultValue: false,
		}).
		AddInput(form.InputFormField{
			Key:          "custom_branch_template",
			Prompt:       "Enter custom default branch template:",
			DefaultValue: "",
			Validator:    nil,
			Condition: func(r *form.FormResult) bool {
				return r.GetBool("configure_branch_template")
			},
		}).
		AddConfirm(form.ConfirmFormField{
			Key:          "configure_pr_template",
			Prompt:       "Configure pull request templates?",
			DefaultValue: false,
		}).
		AddInput(form.InputFormField{
			Key:          "custom_pr_template",
			Prompt:       "Enter custom pull request template:",
			DefaultValue: "",
			Validator:    nil,
			Condition: func(r *form.FormResult) bool {
				return r.GetBool("configure_pr_template")
			},
		})

	projectResult, err := projectForm.Run()
	if err != nil {
		return nil, fmt.Errorf("收集项目模板配置失败: %w", err)
	}

	// Build configuration from results
	cfg := &repoConfig{
		// Personal preference
		branchPrefix:         strings.TrimSpace(personalResult.GetString("branch_prefix")),
		autoAcceptChangeType: personalResult.GetBool("auto_accept_change_type"),
		// Project template
		useScope:                projectResult.GetBool("use_scope"),
		configureCommitTemplate: projectResult.GetBool("configure_commit_template"),
		customCommitTemplate:    strings.TrimSpace(projectResult.GetString("custom_commit_template")),
		configureBranchTemplate: projectResult.GetBool("configure_branch_template"),
		customBranchTemplate:    strings.TrimSpace(projectResult.GetString("custom_branch_template")),
		configurePRTemplate:     projectResult.GetBool("configure_pr_template"),
		customPRTemplate:        strings.TrimSpace(projectResult.GetString("custom_pr_template")),
	}

	return cfg, nil
}

// repoConfig represents the collected configuration
type repoConfig struct {
	// Personal preference
	branchPrefix         string
	autoAcceptChangeType bool
	// Project template
	useScope                bool
	configureCommitTemplate bool
	customCommitTemplate    string
	configureBranchTemplate bool
	customBranchTemplate    string
	configurePRTemplate     bool
	customPRTemplate        string
}

// saveConfig saves the configuration to files
func saveConfig(manager *config.RepoManager, cfg *repoConfig) error {
	// Save public configuration (project template)
	if err := savePublicConfig(manager, cfg); err != nil {
		return fmt.Errorf("保存项目公共配置失败: %w", err)
	}

	// Save private configuration (personal preference)
	if err := savePrivateConfig(manager, cfg); err != nil {
		return fmt.Errorf("保存个人偏好配置失败: %w", err)
	}

	return nil
}

// savePublicConfig saves public configuration to .workflow/config.toml
func savePublicConfig(manager *config.RepoManager, cfg *repoConfig) error {
	// Update template config
	templateConfig := manager.GetTemplateConfig()

	// Commit template
	if cfg.useScope {
		templateConfig.Commit["use_scope"] = true
	} else {
		templateConfig.Commit["use_scope"] = false
	}

	if cfg.configureCommitTemplate && cfg.customCommitTemplate != "" {
		templateConfig.Commit["default"] = cfg.customCommitTemplate
	}

	// Branch template
	if cfg.configureBranchTemplate && cfg.customBranchTemplate != "" {
		templateConfig.Branch["default"] = cfg.customBranchTemplate
	}

	// PR template
	if cfg.configurePRTemplate && cfg.customPRTemplate != "" {
		templateConfig.PullRequests["default"] = cfg.customPRTemplate
	}

	// Save
	return manager.Save()
}

// savePrivateConfig saves private configuration to global config directory
func savePrivateConfig(manager *config.RepoManager, cfg *repoConfig) error {
	// Load existing private config
	privateConfig := manager.LoadPrivateConfig()
	if privateConfig == nil {
		privateConfig = &config.PrivateRepoConfig{
			Repositories: make(map[string]config.PrivateRepoSection),
		}
	}

	// Get repo ID
	repoID := manager.GetRepoID()

	// Get or create repo section
	repoSection, exists := privateConfig.Repositories[repoID]
	if !exists {
		repoSection = config.PrivateRepoSection{}
	}

	// Update branch prefix
	if cfg.branchPrefix != "" {
		if repoSection.Branch == nil {
			repoSection.Branch = &config.BranchConfig{}
		}
		repoSection.Branch.Prefix = &cfg.branchPrefix
	}

	// Update auto-accept change type
	repoSection.AutoAcceptChangeType = &cfg.autoAcceptChangeType

	// Save back
	privateConfig.Repositories[repoID] = repoSection

	// Save to file
	return manager.SavePrivateConfig(privateConfig)
}

// Helper functions

func getBranchPrefixPrompt(currentPrefix string) string {
	if currentPrefix != "" {
		return fmt.Sprintf("Enter branch prefix (press Enter to keep current: %s):", currentPrefix)
	}
	return "Enter branch prefix (optional, press Enter to skip, e.g., 'feature', 'fix'):"
}

func getRepoName(gitRepo *git.Repository) (string, error) {
	// Try to get remote URL
	url, err := gitRepo.GetRemoteURL("origin")
	if err != nil {
		// If no remote, use directory name
		wd, err := os.Getwd()
		if err != nil {
			return "unknown", nil
		}
		return filepath.Base(wd), nil
	}

	// Extract repo name from URL
	repoName, err := git.ExtractRepoName(url)
	if err != nil {
		return "unknown", nil
	}
	return repoName, nil
}

func isTerminal() bool {
	fileInfo, _ := os.Stdin.Stat()
	return (fileInfo.Mode() & os.ModeCharDevice) != 0
}

func checkConfigExists() (bool, error) {
	manager, err := infrastructureconfig.NewRepoManagerWithDefaultGit("")
	if err != nil {
		return false, err
	}

	// Check if private config exists and is configured
	privateConfig := manager.LoadPrivateConfig()
	if privateConfig == nil {
		return false, nil
	}

	repoID := manager.GetRepoID()
	repoSection, exists := privateConfig.Repositories[repoID]
	if !exists {
		return false, nil
	}

	// Check if any configuration is set
	hasConfig := repoSection.Branch != nil || repoSection.AutoAcceptChangeType != nil
	return hasConfig, nil
}

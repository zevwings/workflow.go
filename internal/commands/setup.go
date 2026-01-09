package commands

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/zevwings/workflow/internal/config"
	"github.com/zevwings/workflow/internal/prompt"
	"github.com/zevwings/workflow/internal/prompt/form"
)

// NewSetupCmd creates the setup command
func NewSetupCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "setup",
		Short: "Initialize or update configuration (interactive)",
		Long: `Interactive initialization or update of Workflow CLI configuration.
This command will guide you through the configuration process, including GitHub, Jira, and other services.`,
		RunE: runSetup,
	}

	return cmd
}

func runSetup(cmd *cobra.Command, args []string) error {
	msg := prompt.NewMessage(false)

	msg.Info("Starting Workflow CLI initialization...")

	// Step 1: Load GlobalConfig
	manager, err := config.Global()
	if err != nil {
		return fmt.Errorf("failed to create config manager: %w", err)
	}

	configExists := false
	if err := manager.Load(); err == nil {
		configExists = true
	} else {
		// Config file doesn't exist, Load() has set Config to zero value
		configExists = false
	}

	cfg := manager.Config

	// Step 1: GitHub configuration
	if err := handleGitHubConfig(msg, cfg, configExists); err != nil {
		return err
	}

	// Step 2: Jira configuration
	if err := handleJiraConfig(msg, cfg, configExists); err != nil {
		return err
	}

	// Step 3: LLM configuration
	if err := handleLLMConfig(msg, cfg, configExists); err != nil {
		return err
	}

	// Step 4: Log configuration
	if err := handleLogConfig(msg, cfg); err != nil {
		return err
	}

	// Save configuration
	manager.Config = cfg
	if err := manager.Save(); err != nil {
		return fmt.Errorf("failed to save configuration: %w", err)
	}

	msg.Success("Configuration saved to: %s", manager.GetConfigPath())
	msg.Info("You can use 'workflow config show' to view the configuration")

	return nil
}

// handleGitHubConfig handles GitHub configuration
func handleGitHubConfig(msg *prompt.Message, cfg *config.GlobalConfig, configExists bool) error {
	hasGitHub := len(cfg.GitHub.Accounts) > 0

	// Always print the info message about detected accounts
	msg.Break()
	if hasGitHub {
		msg.Info("The following GitHub accounts were detected:")
	} else {
		msg.Info("No GitHub accounts were detected.")
	}
	msg.Break()

	if hasGitHub && configExists {
		// Build account selection options
		options := []string{}
		accountMap := make(map[int]string) // Map option index to account name
		currentIndex := 0
		optionIndex := 0

		// First add current account (if exists)
		if cfg.GitHub.Current != "" {
			options = append(options, fmt.Sprintf("Keep current account (%s)", cfg.GitHub.Current))
			accountMap[optionIndex] = cfg.GitHub.Current
			currentIndex = optionIndex
			optionIndex++
		}

		// Add other accounts
		for _, account := range cfg.GitHub.Accounts {
			if account.Name != cfg.GitHub.Current {
				options = append(options, fmt.Sprintf("Use existing account (%s)", account.Name))
				accountMap[optionIndex] = account.Name
				optionIndex++
			}
		}

		// Add new account option
		options = append(options, "Add new account")
		addNewAccountIndex := optionIndex

		choice, err := prompt.AskSelect(prompt.SelectField{
			Message:      "GitHub account management",
			Options:      options,
			DefaultIndex: currentIndex,
			ResultTitle:  "GitHub account management",
		})
		if err != nil {
			return err
		}

		// If "Add new account" is selected
		if choice == addNewAccountIndex {
			// Continue with the add new account flow below
		} else if accountName, ok := accountMap[choice]; ok {
			// Select to use existing account (including keep current)
			cfg.GitHub.Current = accountName
			return nil
		}
	}

	// Collect GitHub account information
	result, err := prompt.Form().
		SetTitle("GitHub Configuration").
		AddInput(form.InputFormField{
			Key:          "name",
			Prompt:       "Please enter your account name (required)",
			DefaultValue: "",
			Validator:    nil,
			ResultTitle:  "Your account name",
		}).
		AddInput(form.InputFormField{
			Key:          "email",
			Prompt:       "Please enter your email (required)",
			DefaultValue: "",
			Validator:    nil,
			ResultTitle:  "Your email",
		}).
		AddPassword(form.PasswordFormField{
			Key:          "api_token",
			Prompt:       "Please enter your GitHub Personal Access Token (required)",
			DefaultValue: "",
			Validator:    nil,
			ResultTitle:  "Your GitHub Personal Access Token",
		}).
		Run()
	if err != nil {
		return fmt.Errorf("failed to configure GitHub: %w", err)
	}

	accountName := result.GetString("name")
	if accountName == "" {
		accountName = "default"
	}
	email := result.GetString("email")
	apiToken := result.GetString("api_token")

	if apiToken != "" {
		// Check if account already exists
		accountExists := false
		for idx := range cfg.GitHub.Accounts {
			if cfg.GitHub.Accounts[idx].Name == accountName {
				cfg.GitHub.Accounts[idx].APIToken = apiToken
				if email != "" {
					cfg.GitHub.Accounts[idx].Email = email
				}
				accountExists = true
				break
			}
		}

		if !accountExists {
			cfg.GitHub.Accounts = append(cfg.GitHub.Accounts, config.GitHubAccount{
				Name:     accountName,
				Email:    email,
				APIToken: apiToken,
			})
		}

		if cfg.GitHub.Current == "" {
			cfg.GitHub.Current = accountName
		}
	}

	return nil
}

// handleJiraConfig handles Jira configuration
func handleJiraConfig(msg *prompt.Message, cfg *config.GlobalConfig, configExists bool) error {
	hasJira := cfg.Jira.ServiceAddress != "" || cfg.Jira.APIToken != ""

	// Always print the info message about detected configuration
	msg.Break()
	if hasJira {
		msg.Info("Jira configuration detected.")
	} else {
		msg.Info("No Jira configuration detected.")
	}
	msg.Break()

	if hasJira && configExists {
		keepJira, err := prompt.AskConfirm(prompt.ConfirmField{
			Message:     "Existing Jira configuration detected. Do you want to keep the current values?",
			DefaultYes:  true,
			ResultTitle: "Keep Jira configuration",
		})
		if err != nil {
			return err
		}
		if keepJira {
			return nil
		}
		// If user chooses not to keep, continue with update flow
	}

	// If doesn't exist or choose to update
	serviceAddressPrompt := "Please enter your Jira service address (required)"
	if hasJira && cfg.Jira.ServiceAddress != "" {
		serviceAddressPrompt = "Please enter your Jira service address (press Enter to keep)"
	} else {
		serviceAddressPrompt = "Please enter your Jira service address (required)"
	}

	emailPrompt := "Please enter your Jira email (required)"
	if hasJira && cfg.Jira.Email != "" {
		emailPrompt = "Please enter your Jira email (press Enter to keep)"
	} else {
		emailPrompt = "Please enter your Jira email (required)"
	}

	tokenPrompt := "Please enter your Jira API token (required)"
	var tokenDefaultValue string
	if hasJira && cfg.Jira.APIToken != "" {
		tokenPrompt = "Please enter your Jira API token (press Enter to keep)"
		tokenDefaultValue = cfg.Jira.APIToken
	} else {
		tokenPrompt = "Please enter your Jira API token (required)"
		tokenDefaultValue = ""
	}

	// Collect Jira configuration information
	result, err := prompt.Form().
		SetTitle("Jira Configuration").
		AddInput(form.InputFormField{
			Key:          "service_address",
			Prompt:       serviceAddressPrompt,
			DefaultValue: cfg.Jira.ServiceAddress,
			Validator:    nil,
			ResultTitle:  "Your Jira service address",
		}).
		AddInput(form.InputFormField{
			Key:          "email",
			Prompt:       emailPrompt,
			DefaultValue: cfg.Jira.Email,
			Validator:    nil,
			ResultTitle:  "Your Jira email",
		}).
		AddPassword(form.PasswordFormField{
			Key:          "api_token",
			Prompt:       tokenPrompt,
			DefaultValue: tokenDefaultValue,
			Validator:    nil,
			ResultTitle:  "Your Jira API token",
		}).
		Run()
	if err != nil {
		return fmt.Errorf("failed to configure Jira: %w", err)
	}

	serviceAddress := result.GetString("service_address")
	email := result.GetString("email")
	apiToken := result.GetString("api_token")

	// Use existing values if user pressed Enter (empty input)
	if serviceAddress != "" {
		cfg.Jira.ServiceAddress = serviceAddress
	}
	if email != "" {
		cfg.Jira.Email = email
	}
	if apiToken != "" {
		cfg.Jira.APIToken = apiToken
	}

	return nil
}

// handleLLMConfig handles LLM configuration
func handleLLMConfig(msg *prompt.Message, cfg *config.GlobalConfig, configExists bool) error {
	hasLLM := cfg.LLM.Provider != ""

	// Always print the info message about detected configuration
	msg.Break()
	if hasLLM {
		msg.Info("LLM configuration detected (Provider: %s).", cfg.LLM.Provider)
	} else {
		msg.Info("No LLM configuration detected.")
	}
	msg.Break()

	if hasLLM && configExists {
		keepLLM, err := prompt.AskConfirm(prompt.ConfirmField{
			Message:     fmt.Sprintf("Existing LLM configuration detected (Provider: %s). Do you want to keep the current values?", cfg.LLM.Provider),
			DefaultYes:  true,
			ResultTitle: "Keep LLM configuration",
		})
		if err != nil {
			return err
		}
		if keepLLM {
			return nil
		}
		// If user chooses not to keep, skip the "Do you want to configure LLM?" question
		// and go directly to the configuration flow
	} else {
		// If doesn't exist, ask if user wants to configure LLM
		configureLLM, err := prompt.AskConfirm(prompt.ConfirmField{
			Message:     "Do you want to configure LLM?",
			DefaultYes:  false,
			ResultTitle: "Configure LLM",
		})
		if err != nil {
			return err
		}

		if !configureLLM {
			return nil
		}
	}

	// Select LLM provider type
	providerOptions := []string{"openai", "deepseek", "proxy"}
	providerPrompt := "Please select your LLM provider (required)"
	if hasLLM {
		providerPrompt = fmt.Sprintf("Please select your LLM provider [current: %s]", cfg.LLM.Provider)
	}

	defaultProviderIndex := 0
	if hasLLM {
		for i, provider := range providerOptions {
			if provider == cfg.LLM.Provider {
				defaultProviderIndex = i
				break
			}
		}
	}

	providerIndex, err := prompt.AskSelect(prompt.SelectField{
		Message:      providerPrompt,
		Options:      providerOptions,
		DefaultIndex: defaultProviderIndex,
		ResultTitle:  "Your LLM provider",
	})
	if err != nil {
		return err
	}

	var result *prompt.FormResult

	switch providerIndex {
	case 0: // OpenAI
		apiKeyPrompt := "Please enter your OpenAI API key (required)"
		var apiKeyDefaultValue string
		if cfg.LLM.OpenAI.APIKey != "" {
			apiKeyPrompt = "Please enter your OpenAI API key (press Enter to keep)"
			apiKeyDefaultValue = cfg.LLM.OpenAI.APIKey
		} else {
			apiKeyPrompt = "Please enter your OpenAI API key (required)"
			apiKeyDefaultValue = ""
		}

		modelPrompt := "Please enter your OpenAI model (required)"
		var modelDefaultValue string
		if cfg.LLM.OpenAI.Model != "" {
			modelPrompt = "Please enter your OpenAI model (press Enter to keep)"
			modelDefaultValue = cfg.LLM.OpenAI.Model
		} else {
			modelPrompt = "Please enter your OpenAI model (required)"
			modelDefaultValue = ""
		}

		result, err = prompt.Form().
			SetTitle("OpenAI Configuration").
			AddPassword(form.PasswordFormField{
				Key:          "api_key",
				Prompt:       apiKeyPrompt,
				DefaultValue: apiKeyDefaultValue,
				Validator:    prompt.ValidateRequired(),
				ResultTitle:  "Your OpenAI API key",
			}).
			AddInput(form.InputFormField{
				Key:          "model",
				Prompt:       modelPrompt,
				DefaultValue: modelDefaultValue,
				Validator:    prompt.ValidateRequired(),
				ResultTitle:  "Your OpenAI model",
			}).
			Run()
		if err != nil {
			return fmt.Errorf("failed to configure OpenAI: %w", err)
		}
		cfg.LLM.Provider = "openai"
		apiKey := result.GetString("api_key")
		if apiKey != "" {
			cfg.LLM.OpenAI.APIKey = apiKey
		}
		model := result.GetString("model")
		if model != "" {
			cfg.LLM.OpenAI.Model = model
		}

	case 1: // DeepSeek
		apiKeyPrompt := "Please enter your DeepSeek API key (required)"
		var apiKeyDefaultValue string
		if cfg.LLM.DeepSeek.APIKey != "" {
			apiKeyPrompt = "Please enter your DeepSeek API key (press Enter to keep)"
			apiKeyDefaultValue = cfg.LLM.DeepSeek.APIKey
		} else {
			apiKeyPrompt = "Please enter your DeepSeek API key (required)"
			apiKeyDefaultValue = ""
		}

		modelPrompt := "Please enter your DeepSeek model (required)"
		var modelDefaultValue string
		if cfg.LLM.DeepSeek.Model != "" {
			modelPrompt = "Please enter your DeepSeek model (press Enter to keep)"
			modelDefaultValue = cfg.LLM.DeepSeek.Model
		} else {
			modelPrompt = "Please enter your DeepSeek model (required)"
			modelDefaultValue = ""
		}

		result, err = prompt.Form().
			SetTitle("DeepSeek Configuration").
			AddPassword(form.PasswordFormField{
				Key:          "api_key",
				Prompt:       apiKeyPrompt,
				DefaultValue: apiKeyDefaultValue,
				Validator:    prompt.ValidateRequired(),
				ResultTitle:  "Your DeepSeek API key",
			}).
			AddInput(form.InputFormField{
				Key:          "model",
				Prompt:       modelPrompt,
				DefaultValue: modelDefaultValue,
				Validator:    prompt.ValidateRequired(),
				ResultTitle:  "Your DeepSeek model",
			}).
			Run()
		if err != nil {
			return fmt.Errorf("failed to configure DeepSeek: %w", err)
		}
		cfg.LLM.Provider = "deepseek"
		apiKey := result.GetString("api_key")
		if apiKey != "" {
			cfg.LLM.DeepSeek.APIKey = apiKey
		}
		model := result.GetString("model")
		if model != "" {
			cfg.LLM.DeepSeek.Model = model
		}

	case 2: // Custom Provider (Proxy)
		urlPrompt := "Please enter your LLM proxy URL (required)"
		if cfg.LLM.Proxy.URL != "" {
			urlPrompt = "Please enter your LLM proxy URL (press Enter to keep)"
		} else {
			urlPrompt = "Please enter your LLM proxy URL (required)"
		}

		keyPrompt := "Please enter your LLM proxy key (required)"
		var keyDefaultValue string
		if cfg.LLM.Proxy.APIKey != "" {
			keyPrompt = "Please enter your LLM proxy key (press Enter to keep)"
			keyDefaultValue = cfg.LLM.Proxy.APIKey
		} else {
			keyPrompt = "Please enter your LLM proxy key (required)"
			keyDefaultValue = ""
		}

		modelPrompt := "Please enter your LLM model (required)"
		if cfg.LLM.Proxy.Model != "" {
			modelPrompt = "Please enter your LLM model (press Enter to keep)"
		} else {
			modelPrompt = "Please enter your LLM model (required)"
		}

		result, err = prompt.Form().
			SetTitle("Custom Provider (Proxy) Configuration").
			AddInput(form.InputFormField{
				Key:          "url",
				Prompt:       urlPrompt,
				DefaultValue: cfg.LLM.Proxy.URL,
				Validator:    prompt.ValidateRequired(),
				ResultTitle:  "Your LLM proxy URL",
			}).
			AddPassword(form.PasswordFormField{
				Key:          "api_key",
				Prompt:       keyPrompt,
				DefaultValue: keyDefaultValue,
				Validator:    nil,
				ResultTitle:  "Your LLM proxy key",
			}).
			AddInput(form.InputFormField{
				Key:          "model",
				Prompt:       modelPrompt,
				DefaultValue: cfg.LLM.Proxy.Model,
				Validator:    prompt.ValidateRequired(),
				ResultTitle:  "Your LLM model",
			}).
			Run()
		if err != nil {
			return fmt.Errorf("failed to configure custom provider (proxy): %w", err)
		}
		cfg.LLM.Provider = "proxy"
		url := result.GetString("url")
		if url != "" {
			cfg.LLM.Proxy.URL = url
		} else if cfg.LLM.Proxy.URL == "" {
			return fmt.Errorf("LLM proxy API URL is required")
		}
		apiKey := result.GetString("api_key")
		if apiKey != "" {
			cfg.LLM.Proxy.APIKey = apiKey
		} else if cfg.LLM.Proxy.APIKey == "" {
			return fmt.Errorf("LLM proxy API key is required")
		}
		model := result.GetString("model")
		if model != "" {
			cfg.LLM.Proxy.Model = model
		} else if cfg.LLM.Proxy.Model == "" {
			return fmt.Errorf("Model is required for proxy provider")
		}
	}

	return nil
}

// handleLogConfig handles Log configuration
func handleLogConfig(msg *prompt.Message, cfg *config.GlobalConfig) error {
	// Always print the info message about detected configuration
	msg.Break()
	if cfg.Log.Level != "" {
		msg.Info("Log configuration detected (Level: %s).", cfg.Log.Level)
	} else {
		msg.Info("No log configuration detected.")
	}
	msg.Break()

	// If log level is already configured, default to yes (enabled)
	defaultEnableLogging := cfg.Log.Level != ""
	enableLogging, err := prompt.AskConfirm(prompt.ConfirmField{
		Message:     "Do you want to enable logging?",
		DefaultYes:  defaultEnableLogging,
		ResultTitle: "Enable logging",
	})
	if err != nil {
		return err
	}

	if !enableLogging {
		// If user chooses not to enable logging, clear the existing value
		cfg.Log.Level = ""
		return nil
	}

	// Select log level
	levelOptions := []string{"error", "warn", "info", "debug"}
	defaultIndex := 2 // info is the default

	// If configuration exists, find the corresponding index
	if cfg.Log.Level != "" {
		for i, level := range levelOptions {
			if level == strings.ToLower(cfg.Log.Level) {
				defaultIndex = i
				break
			}
		}
	}

	levelPrompt := "Please select your log level (required)"
	if cfg.Log.Level != "" {
		levelPrompt = fmt.Sprintf("Please select your log level [current: %s]", cfg.Log.Level)
	}

	levelIndex, err := prompt.AskSelect(prompt.SelectField{
		Message:      levelPrompt,
		Options:      levelOptions,
		DefaultIndex: defaultIndex,
		ResultTitle:  "Your log level",
	})
	if err != nil {
		return err
	}

	cfg.Log.Level = levelOptions[levelIndex]

	return nil
}

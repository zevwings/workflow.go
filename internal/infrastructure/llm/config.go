package llm

import (
	"fmt"

	"github.com/zevwings/workflow/internal/config"
	"github.com/zevwings/workflow/internal/llm"
)

// ============================================================================
// Type Definitions
// ============================================================================

// llmConfigProvider implements the llm.LLMConfigProvider interface
//
// Wraps config.LLMConfig to provide interface implementation.
// Encapsulates LLM configuration conversion logic into the infrastructure layer to avoid other packages directly depending on config package's concrete types.
type llmConfigProvider struct {
	llmConfig *config.LLMConfig
}

// ============================================================================
// Interface Implementation
// ============================================================================

// GetProviderConfig gets provider configuration
//
// Gets the current provider's configuration information (APIKey, Model, URL) from config.LLMConfig,
// and converts it to llm.ProviderConfig format for return.
//
// Returns:
//   - *llm.ProviderConfig: LLM provider configuration (APIKey, Model, URL)
//   - error: Returns error if configuration is invalid or retrieval fails
func (p *llmConfigProvider) GetProviderConfig() (*llm.ProviderConfig, error) {
	apiKey, model, url, err := p.llmConfig.CurrentProvider()
	if err != nil {
		return nil, fmt.Errorf("failed to get LLM provider configuration: %w", err)
	}

	return &llm.ProviderConfig{
		APIKey: apiKey,
		Model:  model,
		URL:    url,
	}, nil
}

// GetLanguage gets language configuration
//
// Gets the current language configuration information from config.LLMConfig,
// and converts it to llm.SupportedLanguage format for return.
// If language is not configured, returns default English configuration.
//
// Returns:
//   - *llm.SupportedLanguage: Language configuration information (Code, Name, NativeName, InstructionTemplate)
//   - error: Returns error if language configuration is invalid or retrieval fails
func (p *llmConfigProvider) GetLanguage() (*llm.SupportedLanguage, error) {
	lang, err := p.llmConfig.CurrentLanguage()
	if err != nil {
		return nil, fmt.Errorf("failed to get language configuration: %w", err)
	}

	return &llm.SupportedLanguage{
		Code:                lang.Code,
		Name:                lang.Name,
		NativeName:          lang.NativeName,
		InstructionTemplate: lang.InstructionTemplate,
	}, nil
}

// ============================================================================
// Private Helper Methods
// ============================================================================

// getLLMConfig gets LLM configuration from global config manager
//
// Internal helper method to get global config manager and return LLM configuration.
// Will panic if config manager initialization fails.
//
// Returns:
//   - *config.LLMConfig: LLM configuration instance
//
// Note: Can directly use manager.LLMConfig to access, this method is kept for backward compatibility.
func getLLMConfig() *config.LLMConfig {
	manager, err := config.Global()
	if err != nil {
		panic(fmt.Errorf("infrastructure/llm.getLLMConfig: failed to initialize global config manager: %w", err))
	}
	// Ensure configuration is loaded
	if err := manager.Load(); err != nil {
		// If loading fails, use default configuration
	}
	// Directly return convenience field
	return manager.LLMConfig
}

// ============================================================================
// Public Constructors
// ============================================================================

// NewLLMConfigProvider creates LLM configuration provider
//
// Gets LLM configuration from global config manager and creates llm.LLMConfigProvider interface implementation.
// This function converts config.LLMConfig to llm.LLMConfigProvider interface implementation,
// moving configuration adapter creation to infrastructure layer to avoid other packages directly depending on config package's concrete types.
//
// Returns:
//   - llm.LLMConfigProvider: LLM configuration provider interface implementation
//
// Note:
//   - Function will panic if global config manager initialization fails
//
// Usage example:
//
//	provider := infrastructurellm.NewLLMConfigProvider()
//	prClient := llm.NewPullRequestLLMClient(provider)
func NewLLMConfigProvider() llm.LLMConfigProvider {
	llmConfig := getLLMConfig()
	return &llmConfigProvider{llmConfig: llmConfig}
}

// NewBranchLLMClient creates branch LLM client
//
// Creates and returns branch LLM client instance from global configuration.
// Internally automatically creates configuration provider and LLM client, simplifying client creation process.
//
// Returns:
//   - *llm.BranchLLMClient: Branch LLM client instance
//
// Note:
//   - Function will panic if configuration is invalid or creation fails
//
// Usage example:
//
//	branchClient := infrastructurellm.NewBranchLLMClient()
//	translated, err := branchClient.TranslateToEnglish("Hello")
func NewBranchLLMClient() *llm.BranchLLMClient {
	provider := NewLLMConfigProvider()
	return llm.NewBranchLLMClient(provider)
}

// NewPullRequestLLMClient creates PR LLM client
//
// Creates and returns PR LLM client instance from global configuration.
// Internally automatically creates configuration provider and LLM client, simplifying client creation process.
//
// Returns:
//   - *llm.PullRequestLLMClient: PR LLM client instance
//
// Note:
//   - Function will panic if configuration is invalid or creation fails
//
// Usage example:
//
//	prClient := infrastructurellm.NewPullRequestLLMClient()
//	content, err := prClient.GenerateContent("fix: bug", nil, "")
//	summary, err := prClient.Summarize("PR Title", "PR Diff")
func NewPullRequestLLMClient() *llm.PullRequestLLMClient {
	provider := NewLLMConfigProvider()
	return llm.NewPullRequestLLMClient(provider)
}

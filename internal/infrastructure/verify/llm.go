package verify

import (
	"fmt"
	"strings"

	"github.com/zevwings/workflow/internal/config"
	"github.com/zevwings/workflow/internal/infrastructure/llm"
	llmclient "github.com/zevwings/workflow/internal/llm/client"
	"github.com/zevwings/workflow/internal/prompt"
	"github.com/zevwings/workflow/internal/util"
)

// VerifyLLMConfig verifies LLM configuration
//
// Verifies LLM configuration and connection by sending a simple "hello" test request.
// Verification includes:
//   - Configuration completeness check
//   - API connection test
//   - API Key validity verification
//   - LLM response verification
func VerifyLLMConfig(llmConfig *config.LLMConfig) {
	if llmConfig.Provider == "" {
		return
	}

	msg := prompt.GetMessage()
	msg.Info("LLM Configuration")

	// Prepare configuration information for table display
	var model, key, language string
	language = llmConfig.Language
	if language == "" {
		language = "en"
	}

	switch llmConfig.Provider {
	case "openai":
		model = llmConfig.OpenAI.Model
		if model == "" {
			model = "gpt-3.5-turbo"
		}
		key = util.MaskSensitiveValue(llmConfig.OpenAI.APIKey)
	case "deepseek":
		model = llmConfig.DeepSeek.Model
		if model == "" {
			model = "deepseek-chat"
		}
		key = util.MaskSensitiveValue(llmConfig.DeepSeek.APIKey)
	case "proxy":
		model = llmConfig.Proxy.Model
		// For proxy type, don't display URL in table to avoid table being too wide
		// URL information can be displayed in detailed output
		key = util.MaskSensitiveValue(llmConfig.Proxy.APIKey)
	default:
		model = "-"
		key = "-"
	}

	// 1. Configuration completeness verification
	apiKey, _, _, err := llmConfig.CurrentProvider()
	var status string
	var systemPrompt, userPrompt, testResponse string

	if err != nil {
		status = fmt.Sprintf("✗ Configuration error: %v", err)
	} else if apiKey == "" {
		status = "✗ API Key not configured"
	} else {
		// 2. Create LLM client and send test request
		systemPrompt = "You are a helpful assistant."
		userPrompt = "Say hello"

		spinner := prompt.NewSpinner("Verifying LLM connection...")
		err = spinner.Do(func() error {
			// Create configuration provider
			provider := llm.NewLLMConfigProvider()

			// Get ProviderConfig
			providerConfig, err := provider.GetProviderConfig()
			if err != nil {
				return fmt.Errorf("failed to get LLM provider configuration: %w", err)
			}

			// Create LLM client
			llmClient := llmclient.Global(providerConfig)

			// Send simple test request
			params := &llmclient.LLMRequestParams{
				SystemPrompt: systemPrompt,
				UserPrompt:   userPrompt,
				Temperature:  0.7,
				// MaxTokens:    intPtr(10), // Limit token count to save costs
			}

			testResponse, err = llmClient.Call(params)
			if err != nil {
				return err
			}

			// Verify response is not empty
			if testResponse == "" {
				return fmt.Errorf("LLM returned empty response")
			}

			return nil
		})
		// Note: spinner.Do() internally calls Stop() automatically via defer, no need to call manually

		// 3. Set verification result status
		if err != nil {
			// Provide different prompts based on error type
			errorMsg := err.Error()
			if strings.Contains(errorMsg, "API key") || strings.Contains(errorMsg, "API Key") {
				status = "✗ API Key invalid or not configured"
			} else if strings.Contains(errorMsg, "timeout") {
				status = "✗ Connection timeout"
			} else if strings.Contains(errorMsg, "network") {
				status = "✗ Network connection failed"
			} else {
				status = fmt.Sprintf("✗ Verification failed: %v", err)
			}
		} else {
			status = "✓ Verification successful"
		}
	}

	// 4. Display verification results in vertical table format
	table := prompt.NewTable([]string{"", ""})
	table.AddRow([]string{"Provider", llmConfig.Provider})
	table.AddRow([]string{"Model", model})
	table.AddRow([]string{"Key", key})
	table.AddRow([]string{"Output Language", language})
	table.AddRow([]string{"Status", status})
	table.Render()

	// 5. If verification is successful, output success message with test details
	if testResponse != "" {
		msg.Info("  System prompt: %s", systemPrompt)
		msg.Info("  User prompt: %s", userPrompt)
		msg.Info("  Response: %s", testResponse)
		msg.Success("LLM verified successfully!")
	}

	msg.Break()
}

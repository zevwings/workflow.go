package verify

import (
	"context"
	"time"

	"github.com/zevwings/workflow/internal/config"
	"github.com/zevwings/workflow/internal/http"
	"github.com/zevwings/workflow/internal/prompt"
)

// VerifyEnvironment verifies environment checks (configuration file and network connection)
// Internally creates a table and renders results, consistent with other verification functions
func VerifyEnvironment() {
	msg := prompt.GetMessage()
	msg.Info("Environment Configuration")

	table := prompt.NewTable([]string{"Check Item", "Status", "Description"})

	// Verify configuration file
	manager, err := config.Global()
	var configFileOK bool
	if err != nil {
		table.AddRow([]string{"Configuration File", "✗", "Configuration file does not exist or is invalid"})
		configFileOK = false
	} else if err := manager.Load(); err != nil {
		table.AddRow([]string{"Configuration File", "✗", "Configuration file does not exist or is invalid"})
		configFileOK = false
	} else {
		table.AddRow([]string{"Configuration File", "✓", "Configuration file exists and is valid"})
		configFileOK = true
	}

	// Verify network connection
	client := http.Global()
	spinner := prompt.NewSpinner("Verifying network connection...")
	var respStatusCode int
	var networkOK bool
	err = spinner.Do(func() error {
		// Check GitHub connection (using timeout context)
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		// Use resty client to send request
		restyClient := client.GetRestyClient()
		resp, err := restyClient.R().SetContext(ctx).Get("https://api.github.com")
		if err != nil {
			return err
		}

		// Check response status code
		respStatusCode = resp.StatusCode()
		return nil
	})
	spinner.Stop()

	if err != nil {
		table.AddRow([]string{"Network Connection", "✗", "Network connection failed"})
		networkOK = false
	} else if respStatusCode == 200 {
		table.AddRow([]string{"Network Connection", "✓", "Network connection is normal"})
		networkOK = true
	} else {
		table.AddRow([]string{"Network Connection", "✗", "Network connection failed"})
		networkOK = false
	}

	table.Render()

	// If all checks passed, output success message
	if configFileOK && networkOK {
		msg.Success("Environment verified successfully! Configuration File: ✓ Network Connection: ✓")
	}
	msg.Break()
}

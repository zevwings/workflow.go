package verify

import (
	"github.com/zevwings/workflow/internal/config"
	"github.com/zevwings/workflow/internal/jira"
	"github.com/zevwings/workflow/internal/prompt"
	"github.com/zevwings/workflow/internal/util"
)

// VerifyJiraConfig verifies Jira configuration
func VerifyJiraConfig(jiraConfig *config.JiraConfig) {
	if jiraConfig.Email == "" || jiraConfig.APIToken == "" || jiraConfig.ServiceAddress == "" {
		return
	}

	msg := prompt.GetMessage()
	msg.Info("Jira Configuration")
	table := prompt.NewTable([]string{"Email", "Service Address", "API Token"})

	// Verify Jira authentication
	jiraConfigForAuth := &jira.Config{
		ServiceAddress: jiraConfig.ServiceAddress,
		Email:          jiraConfig.Email,
		APIToken:       jiraConfig.APIToken,
	}

	var jiraResult *jira.AuthResult
	var err error

	spinner := prompt.NewSpinner("Verifying Jira configuration...")
	err = spinner.Do(func() error {
		jiraResult, err = jira.ValidateAuth(jiraConfigForAuth)
		return err
	})

	spinner.Stop()

	if err != nil {
		msg.Warning("Jira verification error: %v", err)
	}

	table.AddRow([]string{
		jiraConfig.Email,
		jiraConfig.ServiceAddress,
		util.MaskSensitiveValue(jiraConfig.APIToken),
	})
	table.Render()

	if jiraResult != nil && jiraResult.Valid {
		if accountID, ok := jiraResult.Details["account_id"].(string); ok && accountID != "" {
			msg.Success("Jira verified successfully! Email: %s (Account ID: %s)", jiraConfig.Email, accountID)
		} else {
			msg.Success("Jira verified successfully! Email: %s", jiraConfig.Email)
		}
	} else if jiraResult != nil {
		msg.Error("Jira verification failed: %s", jiraResult.Message)
	}
	msg.Break()
}

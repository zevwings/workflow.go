package verify

import (
	"fmt"

	"github.com/zevwings/workflow/internal/config"
	"github.com/zevwings/workflow/internal/pr/github"
	"github.com/zevwings/workflow/internal/prompt"
	"github.com/zevwings/workflow/internal/util"
)

// VerifyGitHubConfig verifies GitHub configuration
func VerifyGitHubConfig(githubConfig *config.GitHubConfig) bool {
	if len(githubConfig.Accounts) == 0 {
		return true
	}

	msg := prompt.GetMessage()
	msg.Info("GitHub Configuration")
	table := prompt.NewTable([]string{"Name", "Email", "API Token", "Status", "Verification"})

	allValid := true
	for _, account := range githubConfig.Accounts {
		status := ""
		if account.Name == githubConfig.Current {
			status = "Current"
		}

		// Verify go-github
		var githubResult *github.AuthResult
		var githubErr error

		// Check configuration completeness
		if account.APIToken == "" {
			githubResult = &github.AuthResult{
				Valid:   false,
				Message: "GitHub API Token not configured",
				Details: make(map[string]interface{}),
			}
		} else {
			spinner := prompt.NewSpinner(fmt.Sprintf("Verifying go-github for %s...", account.Name))
			githubErr = spinner.Do(func() error {
				githubResult, githubErr = github.ValidateAuth(account.APIToken)
				return githubErr
			})
			spinner.Stop()

			// If email exists in configuration but not in API response, use email from configuration
			if githubResult != nil && account.Email != "" && githubResult.Details["email"] == nil {
				if githubResult.Details == nil {
					githubResult.Details = make(map[string]interface{})
				}
				githubResult.Details["email"] = account.Email
			}
		}

		// Determine verification result
		githubValid := githubResult != nil && githubResult.Valid && githubErr == nil

		// Display verification result (compact format: github or empty)
		var verification string
		if githubValid {
			verification = "github"
		} else {
			verification = ""
			allValid = false
		}

		// If github verification fails, set allValid = false
		if !githubValid {
			allValid = false
		}

		table.AddRow([]string{
			account.Name,
			account.Email,
			util.MaskSensitiveValue(account.APIToken),
			status,
			verification,
		})
	}

	table.Render()
	if allValid {
		msg.Success("All %d GitHub account(s) verified successfully!", len(githubConfig.Accounts))
	} else {
		msg.Warning("Some GitHub account(s) verification failed. Please check the configuration.")
	}
	msg.Break()

	return allValid
}

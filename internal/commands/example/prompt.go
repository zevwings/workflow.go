//go:build example

package example

import (
	"fmt"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"github.com/zevwings/workflow/internal/prompt"
)

// NewDemoCmd creates a command to demonstrate all Prompt component features
func NewDemoCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "demo-prompt",
		Short: "Demonstrate interactive features of all Prompt components",
		Long: `Demonstrate features of all Prompt components:
- Message: Message output (Info, Success, Warning, Error)
- Confirm: Yes/No confirmation
- Input: Text input
- Password: Password input
- Select: Single selection
- MultiSelect: Multiple selection
- Spinner: Loading indicator
- Table: Table display

This demo will sequentially demonstrate the basic usage of all components.`,
		RunE: runDemo,
	}

	return cmd
}

func runDemo(cmd *cobra.Command, args []string) error {
	msg := prompt.GetMessage()

	msg.Info("Welcome to Prompt component demonstration")
	msg.Break()
	msg.Info("This demo will sequentially demonstrate the following components:")
	msg.Print("  1. Message - Message output")
	msg.Print("  2. Confirm - Confirmation dialog")
	msg.Print("  3. Input - Text input")
	msg.Print("  4. Password - Password input")
	msg.Print("  5. Select - Single selection")
	msg.Print("  6. MultiSelect - Multiple selection")
	msg.Print("  7. Spinner - Loading indicator")
	msg.Print("  8. Table - Table display")
	msg.Break()

	// 1. Message demonstration
	msg.Info("=== 1. Message Output ===")
	msg.Success("Success message example")
	msg.Warning("Warning message example")
	msg.Error("Error message example")
	msg.Break()

	// 2. Confirm demonstration
	msg.Info("=== 2. Confirm Dialog ===")
	confirm, err := prompt.AskConfirm(prompt.ConfirmField{
		Message:    "Do you want to continue the demonstration?",
		DefaultYes: true,
	})
	if err != nil {
		return fmt.Errorf("confirmation failed: %w", err)
	}
	if confirm {
		msg.Success("You selected: Yes")
	} else {
		msg.Warning("You selected: No")
	}
	msg.Break()

	// 3. Input demonstration
	msg.Info("=== 3. Text Input ===")
	name, err := prompt.Input().
		Prompt("Please enter your name").
		DefaultValue("John").
		Validate(prompt.ValidateRequired()).
		Run()
	if err != nil {
		return fmt.Errorf("input failed: %w", err)
	}
	msg.Success("Your name is: %s", name)
	msg.Break()

	// 4. Password demonstration
	msg.Info("=== 4. Password Input ===")
	password, err := prompt.Password().
		Prompt("Please enter password (at least 6 characters)").
		Validate(func(s string) error {
			if len(s) < 6 {
				return fmt.Errorf("length must be at least 6 characters")
			}
			return nil
		}).
		Run()
	if err != nil {
		return fmt.Errorf("input failed: %w", err)
	}
	maskedPassword := maskPassword(password)
	msg.Success("Password entered (length: %d characters, display: %s)", len(password), maskedPassword)
	msg.Break()

	// 5. Select demonstration
	msg.Info("=== 5. Single Selection ===")
	msg.Print("Tip: Use ↑/↓ arrow keys to navigate, Enter to confirm")
	options := []string{"Option A", "Option B", "Option C", "Option D"}
	selectedIndex, err := prompt.Select().
		Prompt("Please select an option").
		Options(options).
		Default(0).
		Run()
	if err != nil {
		return fmt.Errorf("selection failed: %w", err)
	}
	msg.Success("You selected: %s (index: %d)", options[selectedIndex], selectedIndex)
	msg.Break()

	// 6. MultiSelect demonstration
	msg.Info("=== 6. Multiple Selection ===")
	msg.Print("Tip: Use ↑/↓ arrow keys to navigate, Space to toggle selection, Enter to confirm")
	features := []string{"Feature A", "Feature B", "Feature C", "Feature D"}
	selectedIndices, err := prompt.MultiSelect().
		Prompt("Please select features to enable (multiple selection allowed)").
		Options(features).
		Default([]int{0}).
		Run()
	if err != nil {
		return fmt.Errorf("selection failed: %w", err)
	}
	if len(selectedIndices) == 0 {
		msg.Warning("You did not select any features")
	} else {
		var selectedNames []string
		for _, idx := range selectedIndices {
			selectedNames = append(selectedNames, features[idx])
		}
		msg.Success("You selected: %s (indices: %v)", strings.Join(selectedNames, ", "), selectedIndices)
	}
	msg.Break()

	// 7. Spinner demonstration
	msg.Info("=== 7. Loading Indicator ===")
	spinner := prompt.NewSpinner("Processing...")
	spinner.Start()
	time.Sleep(2 * time.Second)
	spinner.WithSuccess("Processing completed")
	msg.Break()

	// 8. Table demonstration
	msg.Info("=== 8. Table Display ===")
	table := prompt.NewTable([]string{"Component", "Status", "Description"})
	table.AddRow([]string{"Message", "✓", "Message output function is normal"})
	table.AddRow([]string{"Confirm", "✓", "Confirmation dialog function is normal"})
	table.AddRow([]string{"Input", "✓", "Text input function is normal"})
	table.AddRow([]string{"Password", "✓", "Password input function is normal"})
	table.AddRow([]string{"Select", "✓", "Single selection function is normal"})
	table.AddRow([]string{"MultiSelect", "✓", "Multiple selection function is normal"})
	table.AddRow([]string{"Spinner", "✓", "Loading indicator function is normal"})
	table.AddRow([]string{"Table", "✓", "Table display function is normal"})
	table.Render()
	msg.Break()

	msg.Success("Demonstration completed! All components are functioning normally.")

	return nil
}

// maskPassword masks password display (only shows first 2 characters)
func maskPassword(password string) string {
	if len(password) <= 2 {
		return strings.Repeat("*", len(password))
	}
	return password[:2] + strings.Repeat("*", len(password)-2)
}

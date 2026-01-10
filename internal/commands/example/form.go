//go:build example

package example

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/zevwings/workflow/internal/prompt"
	"github.com/zevwings/workflow/internal/prompt/form"
)

// NewDemoFormCmd creates a command to demonstrate Form module features
func NewDemoFormCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "demo-form",
		Short: "Demonstrate interactive form features of Form module",
		Long: `Demonstrate Form module features:
- Basic flow: Multi-layer sequential execution (confirm -> input -> select -> input), supports title
- Complex flow: Combination of nested forms and conditional fields, supports nested form title

This demo showcases the core usage scenarios of the Form module, including title functionality.`,
		RunE: runDemoForm,
	}

	return cmd
}

func runDemoForm(cmd *cobra.Command, args []string) error {
	msg := prompt.GetMessage()

	msg.Info("Welcome to Form module demonstration")
	msg.Break()
	msg.Info("This demo will demonstrate the following scenarios:")
	msg.Print("  1. Basic flow - Multi-layer sequential execution (confirm -> input -> select -> input)")
	msg.Print("  2. Complex flow - Nested forms and conditional field combination")
	msg.Break()

	// 1. Basic flow - Multi-layer sequential execution
	msg.Info("=== 1. Basic Flow - Multi-layer Sequential Execution ===")
	result1, err := prompt.Form().
		SetTitle("User Registration Form").
		AddConfirm(form.ConfirmFormField{
			Key:          "agree",
			Prompt:       "Do you agree?",
			DefaultValue: false,
		}).
		AddInput(form.InputFormField{
			Key:          "name",
			Prompt:       "Please enter your name",
			DefaultValue: "",
			Validator:    nil,
		}).
		AddInput(form.InputFormField{
			Key:          "email",
			Prompt:       "Please enter your email",
			DefaultValue: "",
			Validator:    prompt.ValidateEmail(),
		}).
		AddSelect(form.SelectFormField{
			Key:          "role",
			Prompt:       "Please select a role",
			Options:      []string{"Admin", "User"},
			DefaultIndex: 0,
		}).
		AddInput(form.InputFormField{
			Key:          "department",
			Prompt:       "Please enter department",
			DefaultValue: "",
			Validator:    nil,
		}).
		AddInput(form.InputFormField{
			Key:          "remark",
			Prompt:       "Please enter remarks",
			DefaultValue: "",
			Validator:    nil,
		}).
		Run()
	if err != nil {
		return fmt.Errorf("form execution failed: %w", err)
	}
	roleIndex := result1.GetInt("role")
	roles := []string{"Admin", "User"}
	msg.Success("Form results:")
	msg.Print("  agree: %v", result1.GetBool("agree"))
	msg.Print("  name: %s", result1.GetString("name"))
	msg.Print("  email: %s", result1.GetString("email"))
	msg.Print("  role: %s (index: %d)", roles[roleIndex], roleIndex)
	msg.Print("  department: %s", result1.GetString("department"))
	msg.Print("  remark: %s", result1.GetString("remark"))
	msg.Break()

	// 2. Complex flow - Nested forms and conditional field combination
	msg.Info("=== 2. Complex Flow - Nested Forms and Conditional Field Combination ===")
	innerUserForm := prompt.Form().
		SetTitle("User Information").
		AddInput(form.InputFormField{
			Key:          "name",
			Prompt:       "Name",
			DefaultValue: "",
			Validator:    prompt.ValidateRequired(),
		}).
		AddInput(form.InputFormField{
			Key:          "email",
			Prompt:       "Email",
			DefaultValue: "",
			Validator:    prompt.ValidateEmail(),
		}).
		AddSelect(form.SelectFormField{
			Key:          "role",
			Prompt:       "Role",
			Options:      []string{"Admin", "User"},
			DefaultIndex: 0,
		})

	innerAddressForm := prompt.Form().
		SetTitle("Address Information").
		AddInput(form.InputFormField{
			Key:          "city",
			Prompt:       "City",
			DefaultValue: "",
			Validator:    nil,
		}).
		AddInput(form.InputFormField{
			Key:          "street",
			Prompt:       "Street",
			DefaultValue: "",
			Validator:    nil,
		})

	result2, err := prompt.Form().
		SetTitle("Create User").
		AddConfirm(form.ConfirmFormField{
			Key:          "createUser",
			Prompt:       "Do you want to create a user?",
			DefaultValue: false,
		}).
		AddForm(form.NestedFormField{
			Key:        "user",
			Prompt:     "User Information",
			NestedForm: innerUserForm,
			Condition: func(r *prompt.FormResult) bool {
				return r.GetBool("createUser")
			},
		}).
		AddConfirm(form.ConfirmFormField{
			Key:          "hasAddress",
			Prompt:       "Do you have an address?",
			DefaultValue: false,
		}).
		AddForm(form.NestedFormField{
			Key:        "address",
			Prompt:     "Address Information",
			NestedForm: innerAddressForm,
			Condition: func(r *prompt.FormResult) bool {
				return r.GetBool("hasAddress")
			},
		}).
		AddMultiSelect(form.MultiSelectFormField{
			Key:             "tags",
			Prompt:          "Tags",
			Options:         []string{"VIP", "Premium", "Standard"},
			DefaultSelected: []int{},
		}).
		Run()
	if err != nil {
		return fmt.Errorf("form execution failed: %w", err)
	}
	msg.Success("Form results:")
	msg.Print("  createUser: %v", result2.GetBool("createUser"))
	if result2.GetBool("createUser") {
		userResult := result2.GetForm("user")
		if userResult != nil {
			roleIndex2 := userResult.GetInt("role")
			roles2 := []string{"Admin", "User"}
			msg.Print("  user.name: %s", userResult.GetString("name"))
			msg.Print("  user.email: %s", userResult.GetString("email"))
			msg.Print("  user.role: %s (index: %d)", roles2[roleIndex2], roleIndex2)
		}
	} else {
		msg.Print("  user: (not filled, because createUser is false)")
	}
	msg.Print("  hasAddress: %v", result2.GetBool("hasAddress"))
	if result2.GetBool("hasAddress") {
		addressResult := result2.GetForm("address")
		if addressResult != nil {
			msg.Print("  address.city: %s", addressResult.GetString("city"))
			msg.Print("  address.street: %s", addressResult.GetString("street"))
		}
	} else {
		msg.Print("  address: (not filled, because hasAddress is false)")
	}
	tags := result2.GetIntSlice("tags")
	tagNames := []string{"VIP", "Premium", "Standard"}
	var selectedTags []string
	for _, idx := range tags {
		if idx >= 0 && idx < len(tagNames) {
			selectedTags = append(selectedTags, tagNames[idx])
		}
	}
	msg.Print("  tags: %v", selectedTags)
	msg.Break()

	msg.Success("Form module demonstration completed!")

	return nil
}

//go:build example

package example

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/zevwings/workflow/internal/prompt"
	"github.com/zevwings/workflow/internal/prompt/form"
)

// NewDemoFormCmd 创建一个演示 Form 模块功能的命令
func NewDemoFormCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "demo-form",
		Short: "演示 Form 模块的交互式表单功能",
		Long: `演示 Form 模块的功能：
- 基本流程：多层顺序执行（confirm -> input -> select -> input），支持 title
- 复杂流程：嵌套表单与条件字段的组合，支持嵌套表单的 title

这个 demo 展示了 Form 模块的核心使用场景，包括 title 功能。`,
		RunE: runDemoForm,
	}

	return cmd
}

func runDemoForm(cmd *cobra.Command, args []string) error {
	msg := prompt.NewMessage(false)

	msg.Info("欢迎使用 Form 模块演示")
	msg.Break()
	msg.Info("本演示将展示以下场景：")
	msg.Println("  1. 基本流程 - 多层顺序执行（confirm -> input -> select -> input）")
	msg.Println("  2. 复杂流程 - 嵌套表单与条件字段组合")
	msg.Break()

	// 1. 基本流程 - 多层顺序执行
	msg.Info("=== 1. 基本流程 - 多层顺序执行 ===")
	result1, err := prompt.Form().
		SetTitle("用户注册表单").
		AddConfirm(form.ConfirmFormField{
			Key:          "agree",
			Prompt:       "是否同意？",
			DefaultValue: false,
		}).
		AddInput(form.InputFormField{
			Key:         "name",
			Prompt:      "请输入姓名",
			DefaultValue: "",
			Validator:   nil,
		}).
		AddInput(form.InputFormField{
			Key:         "email",
			Prompt:      "请输入邮箱",
			DefaultValue: "",
			Validator:   prompt.ValidateEmail(),
		}).
		AddSelect(form.SelectFormField{
			Key:          "role",
			Prompt:       "请选择角色",
			Options:      []string{"Admin", "User"},
			DefaultIndex: 0,
		}).
		AddInput(form.InputFormField{
			Key:         "department",
			Prompt:      "请输入部门",
			DefaultValue: "",
			Validator:   nil,
		}).
		AddInput(form.InputFormField{
			Key:         "remark",
			Prompt:      "请输入备注",
			DefaultValue: "",
			Validator:   nil,
		}).
		Run()
	if err != nil {
		return fmt.Errorf("表单执行失败: %w", err)
	}
	roleIndex := result1.GetInt("role")
	roles := []string{"Admin", "User"}
	msg.Success("表单结果：")
	msg.Println("  agree: %v", result1.GetBool("agree"))
	msg.Println("  name: %s", result1.GetString("name"))
	msg.Println("  email: %s", result1.GetString("email"))
	msg.Println("  role: %s (索引: %d)", roles[roleIndex], roleIndex)
	msg.Println("  department: %s", result1.GetString("department"))
	msg.Println("  remark: %s", result1.GetString("remark"))
	msg.Break()

	// 2. 复杂流程 - 嵌套表单与条件字段组合
	msg.Info("=== 2. 复杂流程 - 嵌套表单与条件字段组合 ===")
	innerUserForm := prompt.Form().
		SetTitle("用户信息").
		AddInput(form.InputFormField{
			Key:         "name",
			Prompt:      "姓名",
			DefaultValue: "",
			Validator:   prompt.ValidateRequired(),
		}).
		AddInput(form.InputFormField{
			Key:         "email",
			Prompt:      "邮箱",
			DefaultValue: "",
			Validator:   prompt.ValidateEmail(),
		}).
		AddSelect(form.SelectFormField{
			Key:          "role",
			Prompt:       "角色",
			Options:      []string{"Admin", "User"},
			DefaultIndex: 0,
		})

	innerAddressForm := prompt.Form().
		SetTitle("地址信息").
		AddInput(form.InputFormField{
			Key:         "city",
			Prompt:      "城市",
			DefaultValue: "",
			Validator:   nil,
		}).
		AddInput(form.InputFormField{
			Key:         "street",
			Prompt:      "街道",
			DefaultValue: "",
			Validator:   nil,
		})

	result2, err := prompt.Form().
		SetTitle("创建用户").
		AddConfirm(form.ConfirmFormField{
			Key:          "createUser",
			Prompt:       "是否创建用户？",
			DefaultValue: false,
		}).
		AddForm(form.NestedFormField{
			Key:        "user",
			Prompt:     "用户信息",
			NestedForm: innerUserForm,
			Condition: func(r *prompt.FormResult) bool {
				return r.GetBool("createUser")
			},
		}).
		AddConfirm(form.ConfirmFormField{
			Key:          "hasAddress",
			Prompt:       "是否有地址？",
			DefaultValue: false,
		}).
		AddForm(form.NestedFormField{
			Key:        "address",
			Prompt:     "地址信息",
			NestedForm: innerAddressForm,
			Condition: func(r *prompt.FormResult) bool {
				return r.GetBool("hasAddress")
			},
		}).
		AddMultiSelect(form.MultiSelectFormField{
			Key:             "tags",
			Prompt:          "标签",
			Options:         []string{"VIP", "Premium", "Standard"},
			DefaultSelected: []int{},
		}).
		Run()
	if err != nil {
		return fmt.Errorf("表单执行失败: %w", err)
	}
	msg.Success("表单结果：")
	msg.Println("  createUser: %v", result2.GetBool("createUser"))
	if result2.GetBool("createUser") {
		userResult := result2.GetForm("user")
		if userResult != nil {
			roleIndex2 := userResult.GetInt("role")
			roles2 := []string{"Admin", "User"}
			msg.Println("  user.name: %s", userResult.GetString("name"))
			msg.Println("  user.email: %s", userResult.GetString("email"))
			msg.Println("  user.role: %s (索引: %d)", roles2[roleIndex2], roleIndex2)
		}
	} else {
		msg.Println("  user: (未填写，因为 createUser 为 false)")
	}
	msg.Println("  hasAddress: %v", result2.GetBool("hasAddress"))
	if result2.GetBool("hasAddress") {
		addressResult := result2.GetForm("address")
		if addressResult != nil {
			msg.Println("  address.city: %s", addressResult.GetString("city"))
			msg.Println("  address.street: %s", addressResult.GetString("street"))
		}
	} else {
		msg.Println("  address: (未填写，因为 hasAddress 为 false)")
	}
	tags := result2.GetIntSlice("tags")
	tagNames := []string{"VIP", "Premium", "Standard"}
	var selectedTags []string
	for _, idx := range tags {
		if idx >= 0 && idx < len(tagNames) {
			selectedTags = append(selectedTags, tagNames[idx])
		}
	}
	msg.Println("  tags: %v", selectedTags)
	msg.Break()

	msg.Success("Form 模块演示完成！")

	return nil
}

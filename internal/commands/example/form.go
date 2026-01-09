//go:build example

package example

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/zevwings/workflow/internal/prompt"
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
	msg.Println("")
	msg.Info("本演示将展示以下场景：")
	msg.Println("  1. 基本流程 - 多层顺序执行（confirm -> input -> select -> input）")
	msg.Println("  2. 复杂流程 - 嵌套表单与条件字段组合")
	msg.Println("")

	// 1. 基本流程 - 多层顺序执行
	msg.Info("=== 1. 基本流程 - 多层顺序执行 ===")
	result1, err := prompt.Form().
		SetTitle("用户注册表单").
		AddConfirm("agree", "是否同意？", false).
		AddInput("name", "请输入姓名", "", nil).
		AddInput("email", "请输入邮箱", "", prompt.ValidateEmail()).
		AddSelect("role", "请选择角色", []string{"Admin", "User"}, 0).
		AddInput("department", "请输入部门", "", nil).
		AddInput("remark", "请输入备注", "", nil).
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
	msg.Println("")

	// 2. 复杂流程 - 嵌套表单与条件字段组合
	msg.Info("=== 2. 复杂流程 - 嵌套表单与条件字段组合 ===")
	innerUserForm := prompt.Form().
		SetTitle("用户信息").
		AddInput("name", "姓名", "", prompt.ValidateRequired()).
		AddInput("email", "邮箱", "", prompt.ValidateEmail()).
		AddSelect("role", "角色", []string{"Admin", "User"}, 0)

	innerAddressForm := prompt.Form().
		SetTitle("地址信息").
		AddInput("city", "城市", "", nil).
		AddInput("street", "街道", "", nil)

	result2, err := prompt.Form().
		SetTitle("创建用户").
		AddConfirm("createUser", "是否创建用户？", false).
		AddForm("user", "用户信息", innerUserForm).
		Condition(func(r *prompt.FormResult) bool {
			return r.GetBool("createUser")
		}).
		AddConfirm("hasAddress", "是否有地址？", false).
		AddForm("address", "地址信息", innerAddressForm).
		Condition(func(r *prompt.FormResult) bool {
			return r.GetBool("hasAddress")
		}).
		AddMultiSelect("tags", "标签", []string{"VIP", "Premium", "Standard"}, []int{}).
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
	msg.Println("")

	msg.Success("Form 模块演示完成！")

	return nil
}

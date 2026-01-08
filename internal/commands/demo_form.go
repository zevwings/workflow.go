package commands

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
- 多层顺序执行：confirm -> input -> select -> input 等组合
- 表单嵌套：form 嵌套 form，支持多层嵌套
- 条件执行：基于前面字段的值条件性地执行后续字段
- 数据收集：自动收集所有字段的值

这个 demo 会依次展示 Form 模块的各种使用场景。`,
		RunE: runDemoForm,
	}

	return cmd
}

func runDemoForm(cmd *cobra.Command, args []string) error {
	msg := prompt.NewMessage(false)

	msg.Info("欢迎使用 Form 模块演示")
	msg.Println("")
	msg.Info("本演示将依次展示以下场景：")
	msg.Println("  1. 基本使用 - confirm 后多个输入")
	msg.Println("  2. select 后多个输入")
	msg.Println("  3. 复杂组合 - 多层顺序执行")
	msg.Println("  4. 单个嵌套表单")
	msg.Println("  5. 多个嵌套表单（并列）")
	msg.Println("  6. 条件字段")
	msg.Println("  7. 复杂表单（组合所有特性）")
	msg.Println("")

	// 1. 基本使用 - confirm 后多个输入
	msg.Info("=== 1. 基本使用 - confirm 后多个输入 ===")
	result1, err := prompt.Form().
		AddConfirm("createUser", "是否创建用户？", false).
		AddInput("name", "请输入姓名", "", nil).
		AddInput("email", "请输入邮箱", "", prompt.ValidateEmail()).
		AddInput("phone", "请输入电话", "", nil).
		Run()
	if err != nil {
		return fmt.Errorf("表单执行失败: %w", err)
	}
	msg.Success("表单结果：")
	msg.Println("  createUser: %v", result1.GetBool("createUser"))
	msg.Println("  name: %s", result1.GetString("name"))
	msg.Println("  email: %s", result1.GetString("email"))
	msg.Println("  phone: %s", result1.GetString("phone"))
	msg.Println("")

	// 2. select 后多个输入
	msg.Info("=== 2. select 后多个输入 ===")
	result2, err := prompt.Form().
		AddSelect("userType", "请选择用户类型", []string{"个人", "企业"}, 0).
		AddInput("name", "请输入姓名/公司名", "", nil).
		AddInput("email", "请输入邮箱", "", prompt.ValidateEmail()).
		AddInput("address", "请输入地址", "", nil).
		Run()
	if err != nil {
		return fmt.Errorf("表单执行失败: %w", err)
	}
	userTypeIndex := result2.GetInt("userType")
	userTypes := []string{"个人", "企业"}
	msg.Success("表单结果：")
	msg.Println("  userType: %s (索引: %d)", userTypes[userTypeIndex], userTypeIndex)
	msg.Println("  name: %s", result2.GetString("name"))
	msg.Println("  email: %s", result2.GetString("email"))
	msg.Println("  address: %s", result2.GetString("address"))
	msg.Println("")

	// 3. 复杂组合 - 多层顺序执行
	msg.Info("=== 3. 复杂组合 - 多层顺序执行 ===")
	result3, err := prompt.Form().
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
	roleIndex := result3.GetInt("role")
	roles := []string{"Admin", "User"}
	msg.Success("表单结果：")
	msg.Println("  agree: %v", result3.GetBool("agree"))
	msg.Println("  name: %s", result3.GetString("name"))
	msg.Println("  email: %s", result3.GetString("email"))
	msg.Println("  role: %s (索引: %d)", roles[roleIndex], roleIndex)
	msg.Println("  department: %s", result3.GetString("department"))
	msg.Println("  remark: %s", result3.GetString("remark"))
	msg.Println("")

	// 4. 单个嵌套表单
	msg.Info("=== 4. 单个嵌套表单 ===")
	addressForm := prompt.Form().
		AddInput("city", "请输入城市", "", nil).
		AddInput("street", "请输入街道", "", nil).
		AddInput("zip", "请输入邮编", "", nil)

	result4, err := prompt.Form().
		AddInput("name", "请输入姓名", "", nil).
		AddForm("address", "地址信息", addressForm).
		AddInput("phone", "请输入电话", "", nil).
		Run()
	if err != nil {
		return fmt.Errorf("表单执行失败: %w", err)
	}
	addressResult := result4.GetForm("address")
	msg.Success("表单结果：")
	msg.Println("  name: %s", result4.GetString("name"))
	if addressResult != nil {
		msg.Println("  address.city: %s", addressResult.GetString("city"))
		msg.Println("  address.street: %s", addressResult.GetString("street"))
		msg.Println("  address.zip: %s", addressResult.GetString("zip"))
	}
	msg.Println("  phone: %s", result4.GetString("phone"))
	msg.Println("")

	// 5. 多个嵌套表单（并列）
	msg.Info("=== 5. 多个嵌套表单（并列） ===")
	userForm := prompt.Form().
		AddInput("name", "姓名", "", nil).
		AddInput("email", "邮箱", "", nil)

	addressForm2 := prompt.Form().
		AddInput("city", "城市", "", nil).
		AddInput("street", "街道", "", nil)

	contactForm := prompt.Form().
		AddInput("phone", "电话", "", nil).
		AddInput("wechat", "微信", "", nil)

	result5, err := prompt.Form().
		AddForm("user", "用户信息", userForm).
		AddForm("address", "地址信息", addressForm2).
		AddForm("contact", "联系方式", contactForm).
		Run()
	if err != nil {
		return fmt.Errorf("表单执行失败: %w", err)
	}
	userResult := result5.GetForm("user")
	addressResult2 := result5.GetForm("address")
	contactResult := result5.GetForm("contact")
	msg.Success("表单结果：")
	if userResult != nil {
		msg.Println("  user.name: %s", userResult.GetString("name"))
		msg.Println("  user.email: %s", userResult.GetString("email"))
	}
	if addressResult2 != nil {
		msg.Println("  address.city: %s", addressResult2.GetString("city"))
		msg.Println("  address.street: %s", addressResult2.GetString("street"))
	}
	if contactResult != nil {
		msg.Println("  contact.phone: %s", contactResult.GetString("phone"))
		msg.Println("  contact.wechat: %s", contactResult.GetString("wechat"))
	}
	msg.Println("")

	// 6. 条件字段
	msg.Info("=== 6. 条件字段 ===")
	result6, err := prompt.Form().
		AddConfirm("hasEmail", "是否有邮箱？", false).
		AddInput("email", "请输入邮箱", "", prompt.ValidateEmail()).
			Condition(func(r *prompt.FormResult) bool {
				return r.GetBool("hasEmail")
			}).
		AddConfirm("hasPhone", "是否有电话？", false).
		AddInput("phone", "请输入电话", "", nil).
			Condition(func(r *prompt.FormResult) bool {
				return r.GetBool("hasPhone")
			}).
		Run()
	if err != nil {
		return fmt.Errorf("表单执行失败: %w", err)
	}
	msg.Success("表单结果：")
	msg.Println("  hasEmail: %v", result6.GetBool("hasEmail"))
	if result6.GetBool("hasEmail") {
		msg.Println("  email: %s", result6.GetString("email"))
	} else {
		msg.Println("  email: (未填写，因为 hasEmail 为 false)")
	}
	msg.Println("  hasPhone: %v", result6.GetBool("hasPhone"))
	if result6.GetBool("hasPhone") {
		msg.Println("  phone: %s", result6.GetString("phone"))
	} else {
		msg.Println("  phone: (未填写，因为 hasPhone 为 false)")
	}
	msg.Println("")

	// 7. 复杂表单（组合所有特性）
	msg.Info("=== 7. 复杂表单（组合所有特性） ===")
	innerUserForm := prompt.Form().
		AddInput("name", "姓名", "", prompt.ValidateRequired()).
		AddInput("email", "邮箱", "", prompt.ValidateEmail()).
		AddSelect("role", "角色", []string{"Admin", "User"}, 0)

	innerAddressForm := prompt.Form().
		AddInput("city", "城市", "", nil).
		AddInput("street", "街道", "", nil)

	result7, err := prompt.Form().
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
	msg.Println("  createUser: %v", result7.GetBool("createUser"))
	if result7.GetBool("createUser") {
		userResult7 := result7.GetForm("user")
		if userResult7 != nil {
			roleIndex7 := userResult7.GetInt("role")
			roles7 := []string{"Admin", "User"}
			msg.Println("  user.name: %s", userResult7.GetString("name"))
			msg.Println("  user.email: %s", userResult7.GetString("email"))
			msg.Println("  user.role: %s (索引: %d)", roles7[roleIndex7], roleIndex7)
		}
	} else {
		msg.Println("  user: (未填写，因为 createUser 为 false)")
	}
	msg.Println("  hasAddress: %v", result7.GetBool("hasAddress"))
	if result7.GetBool("hasAddress") {
		addressResult7 := result7.GetForm("address")
		if addressResult7 != nil {
			msg.Println("  address.city: %s", addressResult7.GetString("city"))
			msg.Println("  address.street: %s", addressResult7.GetString("street"))
		}
	} else {
		msg.Println("  address: (未填写，因为 hasAddress 为 false)")
	}
	tags := result7.GetIntSlice("tags")
	tagNames := []string{"VIP", "Premium", "Standard"}
	var selectedTags []string
	for _, idx := range tags {
		if idx >= 0 && idx < len(tagNames) {
			selectedTags = append(selectedTags, tagNames[idx])
		}
	}
	msg.Println("  tags: %v", selectedTags)
	msg.Println("")

	msg.Success("Form 模块演示完成！所有功能正常。")

	return nil
}


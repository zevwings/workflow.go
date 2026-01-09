package form

import (
	"github.com/zevwings/workflow/internal/prompt/common"
)

// InputProvider 输入提供者接口（用于避免循环依赖）
type InputProvider interface {
	AskInput(field InputField) (string, error)
	AskPassword(field PasswordField) (string, error)
}

// InputField 输入字段配置（用于接口）
type InputField struct {
	Message      string
	DefaultValue string
	Validator    interface{}
	ResultTitle  string
	Config       *common.PromptConfig
}

// PasswordField 密码字段配置（用于接口）
type PasswordField struct {
	Message      string
	DefaultValue string
	Validator    interface{}
	ResultTitle  string
	Config       *common.PromptConfig
}

package form

// FormConfig Form 模块使用的配置
// 通过这个配置，form 包可以访问 prompt 包的格式化函数
type FormConfig struct {
	FormatPrompt func(message string) string
	FormatAnswer func(value string) string
	FormatError  func(message string) string
	FormatHint   func(message string) string
	// AskInputFunc 用于调用 prompt.AskInput 的函数（避免循环依赖）
	AskInputFunc func(message string, defaultValue string, validator interface{}) (string, error)
	// AskPasswordFunc 用于调用 prompt.AskPassword 的函数（避免循环依赖）
	AskPasswordFunc func(message string, validator interface{}) (string, error)
}

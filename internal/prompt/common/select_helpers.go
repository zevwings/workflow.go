package common

import (
	"sort"

	"github.com/zevwings/workflow/internal/prompt/io"
)

// SelectSetup 选择功能的通用设置
// 包含格式化提示消息、创建管理器等公共逻辑
type SelectSetup struct {
	PromptMsg         string
	PromptMsgWithPrefix string
	RawModeMgr       *io.RawModeManager
	Parser           *io.EscapeSequenceParser
	Renderer         *io.InteractiveRenderer
}

// SetupInteractiveSelect 设置交互式选择的通用组件
// 返回设置好的 SelectSetup 结构
func SetupInteractiveSelect(cfg BasePromptConfig) SelectSetup {
	// 格式化提示消息，未输入时使用 "? " 前缀
	promptMsg := cfg.Config.FormatPrompt(cfg.Message)
	promptMsgWithPrefix := FormatPromptWithPrefix(promptMsg, cfg.Config)

	// 创建原始模式管理器、转义序列解析器和渲染器
	rawModeMgr := io.NewRawModeManager(cfg.Terminal)
	parser := io.NewEscapeSequenceParser(cfg.Terminal)
	renderer := io.NewInteractiveRenderer(cfg.Terminal)

	return SelectSetup{
		PromptMsg:         promptMsg,
		PromptMsgWithPrefix: promptMsgWithPrefix,
		RawModeMgr:       rawModeMgr,
		Parser:           parser,
		Renderer:         renderer,
	}
}

// SelectFallbackOptions 选择 fallback 的选项配置
type SelectFallbackOptions struct {
	// FormatOptionLine 格式化单个选项行的函数
	// 参数: index - 选项索引, option - 选项文本, isDefault - 是否为默认选中
	// 返回: 格式化后的行文本
	FormatOptionLine func(index int, option string, isDefault bool) string

	// GetDefaultIndex 获取默认索引的函数
	GetDefaultIndex func() int

	// ParseInput 解析用户输入的函数
	// 参数: input - 用户输入的字符串
	// 返回: 选中的索引, 是否有效
	ParseInput func(input string) (int, bool)

	// FormatSelectedOption 格式化选中选项的函数（用于显示结果）
	FormatSelectedOption func(index int) string

	// InputPrompt 输入提示文本（如 "请选择 (1-3): "）
	InputPrompt string

	// ResultPrefix 结果显示前缀（如 "已选择: "）
	ResultPrefix string
}

// ExecuteSelectFallback 执行选择 fallback 的通用框架
// 处理通用的 fallback 流程：格式化提示、显示选项列表、读取输入、解析输入、显示结果
//
// 参数:
//   - terminal: 终端接口
//   - message: 原始提示消息
//   - config: 提示配置
//   - options: 选项列表
//   - fallbackOptions: fallback 选项配置
//
// 返回:
//   - selectedIndex: 选中的索引
//   - error: 错误（仅在读取输入失败时返回，其他情况返回默认值）
func ExecuteSelectFallback(
	terminal io.TerminalIO,
	message string,
	config PromptConfig,
	options []string,
	fallbackOptions SelectFallbackOptions,
) (int, error) {
	// 格式化提示消息，未输入时使用 "? " 前缀
	promptMsg := config.FormatPrompt(message)
	promptMsgWithPrefix := FormatPromptWithPrefix(promptMsg, config)

	// 显示提示消息
	terminal.Println(promptMsgWithPrefix)

	// 显示选项列表
	defaultIndex := fallbackOptions.GetDefaultIndex()
	for i, option := range options {
		isDefault := (i == defaultIndex)
		line := fallbackOptions.FormatOptionLine(i, option, isDefault)
		terminal.Print(line)
	}

	// 显示输入提示
	if fallbackOptions.InputPrompt != "" {
		terminal.Print(fallbackOptions.InputPrompt)
	}

	// 读取输入
	inputLine, err := terminal.ReadLine()
	if err != nil {
		// 如果读取失败，返回默认值（不返回错误，因为这是正常的 fallback 行为）
		return defaultIndex, nil
	}

	// 解析输入
	selectedIndex, valid := fallbackOptions.ParseInput(inputLine)
	if !valid {
		// 如果解析失败，返回默认值
		return defaultIndex, nil
	}

	// 显示选择结果
	selectedText := fallbackOptions.FormatSelectedOption(selectedIndex)
	resultText := fallbackOptions.ResultPrefix + selectedText
	terminal.Println(resultText)

	return selectedIndex, nil
}

// MultiSelectFallbackOptions 多选 fallback 的选项配置
type MultiSelectFallbackOptions struct {
	// FormatOptionLine 格式化单个选项行的函数
	// 参数: index - 选项索引, option - 选项文本, isSelected - 是否已选中
	// 返回: 格式化后的行文本
	FormatOptionLine func(index int, option string, isSelected bool) string

	// GetDefaultSelected 获取默认选中项的函数
	GetDefaultSelected func() map[int]bool

	// ParseInput 解析用户输入的函数
	// 参数: input - 用户输入的字符串
	// 返回: 选中的索引列表
	ParseInput func(input string) []int

	// FormatSelectedOptions 格式化选中选项的函数（用于显示结果）
	FormatSelectedOptions func(selectedIndices []int) string

	// Instructions 额外的说明文本（可选）
	Instructions string

	// InputPrompt 输入提示文本（如 "请选择 (例如: 1,3,5): "）
	InputPrompt string

	// ResultPrefix 结果显示前缀（如 "已选择: "）
	ResultPrefix string

	// EmptyResultText 无选择时的显示文本（如 "未选择任何选项"）
	EmptyResultText string
}

// ExecuteMultiSelectFallback 执行多选 fallback 的通用框架
// 处理通用的 fallback 流程：格式化提示、显示选项列表、读取输入、解析输入、显示结果
//
// 参数:
//   - terminal: 终端接口
//   - message: 原始提示消息
//   - config: 提示配置
//   - options: 选项列表
//   - fallbackOptions: fallback 选项配置
//
// 返回:
//   - selectedIndices: 选中的索引列表
//   - error: 错误（仅在读取输入失败时返回，其他情况返回默认值）
func ExecuteMultiSelectFallback(
	terminal io.TerminalIO,
	message string,
	config PromptConfig,
	options []string,
	fallbackOptions MultiSelectFallbackOptions,
) ([]int, error) {
	// 格式化提示消息，未输入时使用 "? " 前缀
	promptMsg := config.FormatPrompt(message)
	promptMsgWithPrefix := FormatPromptWithPrefix(promptMsg, config)

	// 显示提示消息
	terminal.Println(promptMsgWithPrefix)

	// 显示说明文本（如果有）
	if fallbackOptions.Instructions != "" {
		terminal.Println(fallbackOptions.Instructions)
		terminal.Println("")
	}

	// 显示选项列表
	defaultSelected := fallbackOptions.GetDefaultSelected()
	for i, option := range options {
		isSelected := defaultSelected[i]
		line := fallbackOptions.FormatOptionLine(i, option, isSelected)
		terminal.Print(line)
	}

	// 显示输入提示
	if fallbackOptions.InputPrompt != "" {
		terminal.Print(fallbackOptions.InputPrompt)
	}

	// 读取输入
	input, err := terminal.ReadLine()
	if err != nil {
		// 如果读取失败，返回默认值（不返回错误，因为这是正常的 fallback 行为）
		return mapToSlice(defaultSelected), nil
	}

	// 解析输入
	selectedSlice := fallbackOptions.ParseInput(input)

	// 显示选择结果
	if len(selectedSlice) > 0 {
		selectedText := fallbackOptions.FormatSelectedOptions(selectedSlice)
		resultText := fallbackOptions.ResultPrefix + selectedText
		terminal.Println(resultText)
	} else if fallbackOptions.EmptyResultText != "" {
		terminal.Println(fallbackOptions.EmptyResultText)
	}

	return selectedSlice, nil
}

// mapToSlice 将 map[int]bool 转换为排序后的 []int
// 这是一个通用辅助函数，用于多选 fallback
func mapToSlice(selected map[int]bool) []int {
	if len(selected) == 0 {
		return []int{} // 返回空切片而不是 nil
	}
	indices := make([]int, 0, len(selected))
	for idx := range selected {
		indices = append(indices, idx)
	}
	// 使用标准库的快速排序算法（O(n log n)）
	sort.Ints(indices)
	return indices
}

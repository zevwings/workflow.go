package selectpkg

import (
	"fmt"

	"github.com/zevwings/workflow/internal/prompt/common"
)


// SelectConfig 选择功能配置
type SelectConfig struct {
	common.BasePromptConfig
	// Options 选项列表
	Options []string
	// DefaultIndex 默认选中的索引
	DefaultIndex int
}

// Select 选择选项（使用配置结构体）
func Select(cfg SelectConfig) (int, error) {
	if len(cfg.Options) == 0 {
		return -1, fmt.Errorf("选项列表不能为空")
	}

	handler := NewSelectHandler(cfg.Options, cfg.DefaultIndex, cfg.Config)

	// 确保默认索引有效
	currentIndex := handler.ValidateAndAdjustDefaultIndex()

	// 设置交互式选择的通用组件
	setup := common.SetupInteractiveSelect(cfg.BasePromptConfig)
	promptMsg := setup.PromptMsg
	promptMsgWithPrefix := setup.PromptMsgWithPrefix
	rawModeMgr := setup.RawModeMgr
	parser := setup.Parser
	renderer := setup.Renderer

	// 使用原始模式管理器执行交互逻辑
	var result int

	err := rawModeMgr.WithRawModeAndFallback(
		func() error {
			// 使用通用渲染函数渲染选项列表
			formatLine := func(index int, currentIdx int) (string, bool) {
				return handler.FormatOptionLine(index, currentIdx)
			}
			// 使用闭包捕获 currentIndex 的引用，支持动态更新
			getCurrentIndex := func() int {
				return currentIndex
			}
			renderSelect := common.RenderOptions(
				cfg.Terminal,
				renderer,
				len(cfg.Options),
				getCurrentIndex,
				formatLine,
				"使用 ↑/↓ 导航，回车确认",
				cfg.Config,
			)

			// 使用渲染器渲染提示和初始界面（未输入时使用 "? " 前缀）
			if err := renderer.RenderWithPrompt(promptMsgWithPrefix, renderSelect); err != nil {
				return err
			}

			// 使用通用输入处理函数
			return common.HandleInteractiveInput(
				parser,
				cfg.Terminal,
				&currentIndex,
				func(idx int, dir string) (int, bool) {
					return handler.ProcessArrowKey(idx, dir)
				},
				func() (bool, error) {
					selectedText := handler.FormatSelectedOption(currentIndex)
					if err := common.FormatResultWithTitle(cfg.Terminal, promptMsg, selectedText, nil, false, cfg.Message, cfg.Config.FormatResultTitle, cfg.Config.FormatAnswerPrefix); err != nil {
						return false, err
					}
					result = currentIndex
					return true, nil
				},
				nil, // select 不需要处理空格键
				func() {
					// 使用渲染器的 ReRender 方法确保正确清除之前的内容
					renderer.ReRender(renderSelect)
				},
			)
		},
		func() error {
			// Fallback: 如果无法设置原始模式，使用简单编号选择
			selectedIndex, err := selectFallback(cfg)
			result = selectedIndex
			return err
		},
	)

	if err != nil {
		// 如果 err 是取消错误，返回 -1 和错误
		return -1, err
	}

	return result, nil
}

// selectFallback 回退方案：如果无法设置原始模式，使用简单编号选择
func selectFallback(cfg SelectConfig) (int, error) {
	handler := NewSelectHandler(cfg.Options, cfg.DefaultIndex, cfg.Config)
	defaultIndex := handler.ValidateAndAdjustDefaultIndex()

	// 使用通用的 fallback 框架
	return common.ExecuteSelectFallback(
		cfg.Terminal,
		cfg.Message,
		cfg.Config,
		cfg.Options,
		common.SelectFallbackOptions{
			FormatOptionLine: func(index int, option string, isDefault bool) string {
				marker := " "
				if isDefault {
					marker = "*"
				}
				return fmt.Sprintf("  %s %d. %s\n", marker, index+1, option)
			},
			GetDefaultIndex: func() int {
				return defaultIndex
			},
			ParseInput: func(input string) (int, bool) {
				var num int
				_, err := fmt.Sscanf(input, "%d", &num)
				if err != nil {
					return 0, false
				}
				// 验证输入范围并转换为索引（用户输入是 1-based）
				if num < 1 || num > len(cfg.Options) {
					return 0, false
				}
				return num - 1, true
			},
			FormatSelectedOption: func(index int) string {
				return handler.FormatSelectedOption(index)
			},
			InputPrompt:   fmt.Sprintf("请选择 (1-%d): ", len(cfg.Options)),
			ResultPrefix:  "已选择: ",
		},
	)
}

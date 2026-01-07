package selectpkg

import (
	"fmt"
	"os"

	"golang.org/x/term"
)

// Config 选择功能配置
type Config struct {
	// 格式化函数
	FormatPrompt func(message string) string
	FormatAnswer func(value string) string
	FormatHint   func(message string) string
}

// Select 选择选项
func Select(message string, options []string, defaultIndex int, config Config) (int, error) {
	if len(options) == 0 {
		return -1, fmt.Errorf("选项列表不能为空")
	}

	// 确保默认索引有效
	if defaultIndex < 0 || defaultIndex >= len(options) {
		defaultIndex = 0
	}

	// 格式化提示消息
	promptMsg := config.FormatPrompt(message)

	// 获取文件描述符
	fd := int(os.Stdin.Fd())

	// 保存当前终端状态
	oldState, err := term.MakeRaw(fd)
	if err != nil {
		// 如果无法设置原始模式，回退到简单选择
		return selectFallback(message, options, defaultIndex, config)
	}
	defer func() {
		// 恢复光标显示
		fmt.Print("\033[?25h")
		term.Restore(fd, oldState)
	}()

	// 隐藏光标
	fmt.Print("\033[?25l")

	currentIndex := defaultIndex

	// 先输出提示消息，确保换行并重置颜色
	fmt.Print(promptMsg)
	fmt.Print("\033[0m") // 重置所有格式（包括颜色）
	fmt.Print("\r\n")    // 回车+换行，确保光标在新行的行首
	fmt.Print("\r\n")

	// 保存光标位置（在选项列表之前，确保在新行的行首）
	fmt.Print("\033[s")

	// 渲染选项列表
	renderSelect := func(isFirst bool) {
		if !isFirst {
			// 恢复之前保存的光标位置
			fmt.Print("\033[u")
			// 确保回到行首（防止 ANSI 代码导致的位置偏移）
			fmt.Print("\r")
			// 清除从光标位置到屏幕底部的所有内容
			fmt.Print("\033[J")
		}

		// 渲染选项（确保每行从行首开始）
		for i, option := range options {
			// 确保从行首开始
			fmt.Print("\r")
			if i == currentIndex {
				// 当前选中的选项：高亮显示
				selectedText := config.FormatAnswer("> " + option)
				fmt.Print(selectedText)
			} else {
				// 其他选项：普通显示
				fmt.Print("  " + option)
			}
			// 清除到行尾，然后换行
			fmt.Print("\033[K")
			fmt.Println()
		}

		// 在选项列表之后显示提示信息
		hintMsg := config.FormatHint("使用 ↑/↓ 导航，回车确认")
		fmt.Print("\r")
		fmt.Print(hintMsg)
		fmt.Print("\r\n")

		fmt.Print("\033[?25l") // 确保光标隐藏
	}

	// 初始渲染
	renderSelect(true)

	// 读取输入
	var buf [1]byte

	for {
		n, err := os.Stdin.Read(buf[:])
		if err != nil || n == 0 {
			term.Restore(fd, oldState)
			return defaultIndex, fmt.Errorf("读取输入失败: %w", err)
		}

		char := buf[0]

		// 处理转义序列（箭头键）
		if char == '\x1b' {
			// 读取转义序列的后续字符（最多读取2个字符）
			seq := make([]byte, 0, 3)
			seq = append(seq, char)

			// 读取下一个字符
			n2, err2 := os.Stdin.Read(buf[:])
			if err2 != nil || n2 == 0 {
				continue
			}
			seq = append(seq, buf[0])

			// 检查第二个字符
			if buf[0] == '[' || buf[0] == 'O' {
				// 读取第三个字符
				n3, err3 := os.Stdin.Read(buf[:])
				if err3 != nil || n3 == 0 {
					continue
				}
				seq = append(seq, buf[0])

				// 检查是否是箭头键
				seqStr := string(seq[1:])
				if seqStr == "[A" || seqStr == "OA" {
					// 上箭头
					if currentIndex > 0 {
						currentIndex--
						renderSelect(false)
					}
					continue
				}
				if seqStr == "[B" || seqStr == "OB" {
					// 下箭头
					if currentIndex < len(options)-1 {
						currentIndex++
						renderSelect(false)
					}
					continue
				}
			}
			// 不是有效的箭头键序列，忽略
			continue
		}

		// 处理回车键（确认选择）
		if char == '\r' || char == '\n' {
			// 恢复光标位置并清除内容
			fmt.Print("\033[u")
			fmt.Print("\033[J")
			// 显示结果
			fmt.Print(promptMsg)
			fmt.Print(" ")
			selectedText := config.FormatAnswer(options[currentIndex])
			fmt.Println(selectedText)
			// 重置所有 ANSI 格式，确保后续输出格式正确
			fmt.Print("\033[0m")
			term.Restore(fd, oldState)
			return currentIndex, nil
		}

		// 处理 Ctrl+C
		if char == 3 { // Ctrl+C
			// 恢复光标位置并清除内容
			fmt.Print("\033[u")
			fmt.Print("\033[J")
			term.Restore(fd, oldState)
			fmt.Println()
			return -1, fmt.Errorf("用户取消输入")
		}

		// 其他字符：静默忽略
	}
}

// selectFallback 回退方案：如果无法设置原始模式，使用简单编号选择
func selectFallback(message string, options []string, defaultIndex int, config Config) (int, error) {
	// 格式化提示消息
	promptMsg := config.FormatPrompt(message)

	// 显示选项列表
	fmt.Println(promptMsg)
	for i, option := range options {
		marker := " "
		if i == defaultIndex {
			marker = "*"
		}
		fmt.Printf("  %s %d. %s\n", marker, i+1, option)
	}

	// 提示输入
	fmt.Print("请选择 (1-", len(options), "): ")

	// 读取输入
	var input int
	_, err := fmt.Scanf("%d", &input)
	if err != nil {
		// 如果读取失败，返回默认值
		return defaultIndex, nil
	}

	// 验证输入范围
	if input < 1 || input > len(options) {
		return defaultIndex, nil
	}

	// 显示选择结果
	selectedText := config.FormatAnswer(options[input-1])
	fmt.Println("已选择:", selectedText)

	return input - 1, nil
}

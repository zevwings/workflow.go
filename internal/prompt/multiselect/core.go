package multiselect

import (
	"fmt"
	"os"
	"strings"

	"golang.org/x/term"
)

// Config 多选功能配置
type Config struct {
	// 格式化函数
	FormatPrompt func(message string) string
	FormatAnswer func(value string) string
	FormatHint   func(message string) string
}

// MultiSelect 多选功能
func MultiSelect(message string, options []string, defaultSelected []int, config Config) ([]int, error) {
	if len(options) == 0 {
		return nil, fmt.Errorf("选项列表不能为空")
	}

	// 验证并清理默认选中项
	selected := make(map[int]bool)
	for _, idx := range defaultSelected {
		if idx >= 0 && idx < len(options) {
			selected[idx] = true
		}
	}

	// 格式化提示消息
	promptMsg := config.FormatPrompt(message)

	// 获取文件描述符
	fd := int(os.Stdin.Fd())

	// 保存当前终端状态
	oldState, err := term.MakeRaw(fd)
	if err != nil {
		// 如果无法设置原始模式，回退到简单多选
		return multiselectFallback(message, options, defaultSelected, config)
	}
	defer func() {
		// 恢复光标显示
		fmt.Print("\033[?25h")
		term.Restore(fd, oldState)
	}()

	// 隐藏光标
	fmt.Print("\033[?25l")

	currentIndex := 0
	if len(defaultSelected) > 0 && defaultSelected[0] >= 0 && defaultSelected[0] < len(options) {
		currentIndex = defaultSelected[0]
	}

	// 先输出提示消息，确保换行并重置颜色
	fmt.Print(promptMsg)
	fmt.Print("\033[0m") // 重置所有格式（包括颜色）
	fmt.Print("\r\n")    // 回车+换行，确保光标在新行的行首
	fmt.Print("\r\n")

	// 保存光标位置（在选项列表之前，确保在新行的行首）
	fmt.Print("\033[s")

	// 渲染选项列表
	renderMultiSelect := func(isFirst bool) {
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
			prefix := "  "
			if i == currentIndex {
				prefix = "> "
			}

			marker := "[ ]"
			if selected[i] {
				marker = "[x]"
			}

			if i == currentIndex {
				// 当前光标所在行：高亮显示
				line := fmt.Sprintf("%s%s %s", prefix, marker, option)
				highlightedLine := config.FormatAnswer(line)
				fmt.Print(highlightedLine)
			} else {
				// 其他行：普通显示
				line := fmt.Sprintf("%s%s %s", prefix, marker, option)
				fmt.Print(line)
			}
			// 清除到行尾，然后换行
			fmt.Print("\033[K")
			fmt.Println()
		}

		// 在选项列表之后显示提示信息
		hintMsg := config.FormatHint("使用 ↑/↓ 导航，空格键切换选择，回车确认")
		fmt.Print("\r")
		fmt.Print(hintMsg)
		fmt.Print("\r\n")

		fmt.Print("\033[?25l") // 确保光标隐藏
	}

	// 初始渲染
	renderMultiSelect(true)

	// 读取输入
	var buf [1]byte

	for {
		n, err := os.Stdin.Read(buf[:])
		if err != nil || n == 0 {
			term.Restore(fd, oldState)
			return mapToSlice(selected), fmt.Errorf("读取输入失败: %w", err)
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
						renderMultiSelect(false)
					}
					continue
				}
				if seqStr == "[B" || seqStr == "OB" {
					// 下箭头
					if currentIndex < len(options)-1 {
						currentIndex++
						renderMultiSelect(false)
					}
					continue
				}
			}
			// 不是有效的箭头键序列，忽略
			continue
		}

		// 处理空格键（切换选择状态）
		if char == ' ' {
			// 切换当前选项的选中状态
			if selected[currentIndex] {
				delete(selected, currentIndex)
			} else {
				selected[currentIndex] = true
			}
			renderMultiSelect(false)
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

			selectedIndices := mapToSlice(selected)
			if len(selectedIndices) == 0 {
				fmt.Println("(未选择)")
			} else {
				var selectedOptions []string
				for _, idx := range selectedIndices {
					selectedOptions = append(selectedOptions, options[idx])
				}
				selectedText := config.FormatAnswer(strings.Join(selectedOptions, ", "))
				fmt.Println(selectedText)
			}

			// 重置所有 ANSI 格式，确保后续输出格式正确
			fmt.Print("\033[0m")

			term.Restore(fd, oldState)
			return selectedIndices, nil
		}

		// 处理 Ctrl+C
		if char == 3 { // Ctrl+C
			// 恢复光标位置并清除内容
			fmt.Print("\033[u")
			fmt.Print("\033[J")
			term.Restore(fd, oldState)
			fmt.Println()
			return nil, fmt.Errorf("用户取消输入")
		}

		// 其他字符：静默忽略
	}
}

// mapToSlice 将 map[int]bool 转换为排序后的 []int
func mapToSlice(selected map[int]bool) []int {
	var indices []int
	for idx := range selected {
		indices = append(indices, idx)
	}
	// 简单排序（冒泡排序）
	for i := 0; i < len(indices); i++ {
		for j := i + 1; j < len(indices); j++ {
			if indices[i] > indices[j] {
				indices[i], indices[j] = indices[j], indices[i]
			}
		}
	}
	return indices
}

// multiselectFallback 回退方案：如果无法设置原始模式，使用简单多选
func multiselectFallback(message string, options []string, defaultSelected []int, config Config) ([]int, error) {
	// 格式化提示消息
	promptMsg := config.FormatPrompt(message)

	// 显示选项列表
	fmt.Println(promptMsg)
	fmt.Println("请输入选项编号（多个选项用逗号分隔，如：1,3,5）")
	fmt.Println()

	selectedMap := make(map[int]bool)
	for _, idx := range defaultSelected {
		if idx >= 0 && idx < len(options) {
			selectedMap[idx] = true
		}
	}

	for i, option := range options {
		marker := "[ ]"
		if selectedMap[i] {
			marker = "[x]"
		}
		fmt.Printf("  %s %d. %s\n", marker, i+1, option)
	}

	// 提示输入
	fmt.Print("请选择 (例如: 1,3,5): ")

	// 读取输入
	var input string
	_, err := fmt.Scanln(&input)
	if err != nil {
		// 如果读取失败，返回默认值
		return defaultSelected, nil
	}

	// 解析输入
	input = strings.TrimSpace(input)
	if input == "" {
		return defaultSelected, nil
	}

	// 解析逗号分隔的数字
	parts := strings.Split(input, ",")
	selectedIndices := make(map[int]bool)

	for _, part := range parts {
		part = strings.TrimSpace(part)
		var num int
		_, err := fmt.Sscanf(part, "%d", &num)
		if err != nil {
			continue
		}
		if num >= 1 && num <= len(options) {
			selectedIndices[num-1] = true
		}
	}

	// 显示选择结果
	selectedSlice := mapToSlice(selectedIndices)
	if len(selectedSlice) > 0 {
		var selectedOptions []string
		for _, idx := range selectedSlice {
			selectedOptions = append(selectedOptions, options[idx])
		}
		selectedText := config.FormatAnswer(strings.Join(selectedOptions, ", "))
		fmt.Println("已选择:", selectedText)
	} else {
		fmt.Println("未选择任何选项")
	}

	return selectedSlice, nil
}

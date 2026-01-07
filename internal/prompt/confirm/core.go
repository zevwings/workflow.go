package confirm

import (
	"fmt"
	"os"
	"strings"

	"golang.org/x/term"
)

// Config 确认功能配置
type Config struct {
	// 格式化函数
	FormatPrompt func(message string) string
	FormatAnswer func(value string) string
}

// Confirm 内部实现函数
func Confirm(message string, defaultYes bool, config Config) (bool, error) {
	// 格式化提示消息
	promptMsg := config.FormatPrompt(message)

	// 根据 defaultYes 构建不同的提示（显示 [Y/n] 或 [y/N]）
	var promptText string
	if defaultYes {
		promptText = fmt.Sprintf("%s [Y/n] ", promptMsg)
	} else {
		promptText = fmt.Sprintf("%s [y/N] ", promptMsg)
	}

	// 输出提示
	fmt.Print(promptText)

	// 获取文件描述符
	fd := int(os.Stdin.Fd())

	// 保存当前终端状态
	oldState, err := term.MakeRaw(fd)
	if err != nil {
		// 如果无法设置原始模式，回退到普通输入
		return confirmFallback(message, defaultYes, config)
	}
	defer func() {
		// 恢复光标显示
		fmt.Print("\033[?25h")
		term.Restore(fd, oldState)
	}()

	// 隐藏光标
	fmt.Print("\033[?25l")

	// 读取单个字符，检测到 y/n 立即返回
	var buf [1]byte
	for {
		n, err := os.Stdin.Read(buf[:])
		if err != nil || n == 0 {
			term.Restore(fd, oldState)
			return defaultYes, fmt.Errorf("读取输入失败: %w", err)
		}

		char := buf[0]

		// 处理回车键（使用默认值，在同一行显示 yes 或 no）
		if char == '\r' || char == '\n' {
			// 清除 [Y/n] 部分，显示格式化的 yes 或 no
			fmt.Print("\r")      // 回到行首
			fmt.Print("\033[K")  // 清除到行尾
			fmt.Print(promptMsg) // 显示提示消息（带颜色）
			fmt.Print(" ")       // 空格分隔
			var answerText string
			if defaultYes {
				answerText = config.FormatAnswer("yes")
			} else {
				answerText = config.FormatAnswer("no")
			}
			fmt.Println(answerText) // 显示格式化的答案并换行
			// 重置所有 ANSI 格式，确保后续输出格式正确
			fmt.Print("\033[0m")
			term.Restore(fd, oldState)
			return defaultYes, nil
		}

		// 处理 Ctrl+C
		if char == 3 { // Ctrl+C
			term.Restore(fd, oldState)
			fmt.Println()
			return false, fmt.Errorf("用户取消输入")
		}

		// 转换为小写进行比较
		charLower := strings.ToLower(string(char))[0]

		// 处理 yes 输入：y 或 Y（在同一行显示 yes 后确认）
		if charLower == 'y' {
			// 清除 [Y/n] 部分，显示格式化的 yes
			fmt.Print("\r")      // 回到行首
			fmt.Print("\033[K")  // 清除到行尾
			fmt.Print(promptMsg) // 显示提示消息（带颜色）
			fmt.Print(" ")       // 空格分隔
			answerText := config.FormatAnswer("yes")
			fmt.Println(answerText) // 显示格式化的答案并换行
			// 重置所有 ANSI 格式，确保后续输出格式正确
			fmt.Print("\033[0m")
			term.Restore(fd, oldState)
			return true, nil
		}

		// 处理 no 输入：n 或 N（在同一行显示 no 后确认）
		if charLower == 'n' {
			// 清除 [Y/n] 部分，显示格式化的 no
			fmt.Print("\r")      // 回到行首
			fmt.Print("\033[K")  // 清除到行尾
			fmt.Print(promptMsg) // 显示提示消息（带颜色）
			fmt.Print(" ")       // 空格分隔
			answerText := config.FormatAnswer("no")
			fmt.Println(answerText) // 显示格式化的答案并换行
			// 重置所有 ANSI 格式，确保后续输出格式正确
			fmt.Print("\033[0m")
			term.Restore(fd, oldState)
			return false, nil
		}

		// 其他字符：静默忽略，继续等待有效输入
		// 不做任何处理，直接继续循环
	}
}

// confirmFallback 回退方案：如果无法设置原始模式，使用普通输入
func confirmFallback(message string, defaultYes bool, config Config) (bool, error) {
	// 格式化提示消息
	promptMsg := config.FormatPrompt(message)

	// 根据 defaultYes 构建不同的提示（显示 [Y/n] 或 [y/N]）
	var promptText string
	if defaultYes {
		promptText = fmt.Sprintf("%s [Y/n] ", promptMsg)
	} else {
		promptText = fmt.Sprintf("%s [y/N] ", promptMsg)
	}

	// 输出提示
	fmt.Print(promptText)

	// 读取一行输入
	var input string
	_, err := fmt.Scanln(&input)
	if err != nil {
		// 如果读取失败（可能是空输入），返回默认值并显示格式化结果
		fmt.Print("\r")      // 回到行首
		fmt.Print("\033[K")  // 清除到行尾
		fmt.Print(promptMsg) // 显示提示消息（带颜色）
		fmt.Print(" ")       // 空格分隔
		var answerText string
		if defaultYes {
			answerText = config.FormatAnswer("yes")
		} else {
			answerText = config.FormatAnswer("no")
		}
		fmt.Println(answerText) // 显示格式化的答案并换行
		fmt.Print("\033[0m")    // 重置所有 ANSI 格式
		return defaultYes, nil
	}

	// 清理输入
	input = strings.TrimSpace(strings.ToLower(input))

	// 处理空输入（使用默认值）
	if input == "" {
		// 显示格式化结果
		fmt.Print("\r")      // 回到行首
		fmt.Print("\033[K")  // 清除到行尾
		fmt.Print(promptMsg) // 显示提示消息（带颜色）
		fmt.Print(" ")       // 空格分隔
		var answerText string
		if defaultYes {
			answerText = config.FormatAnswer("yes")
		} else {
			answerText = config.FormatAnswer("no")
		}
		fmt.Println(answerText) // 显示格式化的答案并换行
		fmt.Print("\033[0m")    // 重置所有 ANSI 格式
		return defaultYes, nil
	}

	// 处理 yes 输入
	if input == "y" || input == "yes" {
		// 显示格式化结果
		fmt.Print("\r")      // 回到行首
		fmt.Print("\033[K")  // 清除到行尾
		fmt.Print(promptMsg) // 显示提示消息（带颜色）
		fmt.Print(" ")       // 空格分隔
		answerText := config.FormatAnswer("yes")
		fmt.Println(answerText) // 显示格式化的答案并换行
		fmt.Print("\033[0m")    // 重置所有 ANSI 格式
		return true, nil
	}

	// 处理 no 输入
	if input == "n" || input == "no" {
		// 显示格式化结果
		fmt.Print("\r")      // 回到行首
		fmt.Print("\033[K")  // 清除到行尾
		fmt.Print(promptMsg) // 显示提示消息（带颜色）
		fmt.Print(" ")       // 空格分隔
		answerText := config.FormatAnswer("no")
		fmt.Println(answerText) // 显示格式化的答案并换行
		fmt.Print("\033[0m")    // 重置所有 ANSI 格式
		return false, nil
	}

	// 非法输入，返回默认值并显示格式化结果
	fmt.Print("\r")      // 回到行首
	fmt.Print("\033[K")  // 清除到行尾
	fmt.Print(promptMsg) // 显示提示消息（带颜色）
	fmt.Print(" ")       // 空格分隔
	var answerText string
	if defaultYes {
		answerText = config.FormatAnswer("yes")
	} else {
		answerText = config.FormatAnswer("no")
	}
	fmt.Println(answerText) // 显示格式化的答案并换行
	fmt.Print("\033[0m")    // 重置所有 ANSI 格式
	return defaultYes, nil
}

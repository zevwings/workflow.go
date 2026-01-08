package prompt

import (
	"fmt"
	"io"
	"os"
	"sync"
	"time"

	"github.com/charmbracelet/lipgloss"
)

// Spinner 加载指示器
type Spinner struct {
	message      string
	spinner      []string
	interval     time.Duration
	style        lipgloss.Style  // 兼容旧版本，如果设置了 spinnerStyle 和 messageStyle 则优先使用
	spinnerStyle *lipgloss.Style // spinner 字符的样式（使用指针以便检查是否设置）
	messageStyle *lipgloss.Style // 消息文本的样式（使用指针以便检查是否设置）
	stopChan     chan struct{}
	stopped      bool
	mu           sync.Mutex
	writer       io.Writer
	currentFrame int
	cursorHidden bool // 标记光标是否已隐藏
}

// SpinnerOption 配置选项
type SpinnerOption func(*Spinner)

// WithSpinner 设置自定义 spinner 字符序列
func WithSpinner(frames []string) SpinnerOption {
	return func(s *Spinner) {
		if len(frames) > 0 {
			s.spinner = frames
		}
	}
}

// WithInterval 设置更新间隔
func WithInterval(d time.Duration) SpinnerOption {
	return func(s *Spinner) {
		if d > 0 {
			s.interval = d
		}
	}
}

// WithStyle 设置自定义样式（同时应用于 spinner 和文本）
func WithStyle(style lipgloss.Style) SpinnerOption {
	return func(s *Spinner) {
		s.style = style
	}
}

// WithSpinnerStyle 设置 spinner 字符的样式
func WithSpinnerStyle(style lipgloss.Style) SpinnerOption {
	return func(s *Spinner) {
		s.spinnerStyle = &style
	}
}

// WithMessageStyle 设置消息文本的样式
func WithMessageStyle(style lipgloss.Style) SpinnerOption {
	return func(s *Spinner) {
		s.messageStyle = &style
	}
}

// WithWriter 设置输出流
func WithWriter(w io.Writer) SpinnerOption {
	return func(s *Spinner) {
		s.writer = w
	}
}

// NewSpinner 创建新的加载指示器
func NewSpinner(message string, opts ...SpinnerOption) *Spinner {
	theme := GetTheme()

	// 默认 spinner 动画帧
	defaultSpinner := []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}

	// 默认样式（使用 InfoStyle）
	defaultStyle := theme.InfoStyle

	s := &Spinner{
		message:      message,
		spinner:      defaultSpinner,
		interval:     100 * time.Millisecond,
		style:        defaultStyle,
		spinnerStyle: nil, // 默认未设置
		messageStyle: nil, // 默认未设置
		stopChan:     make(chan struct{}),
		writer:       os.Stdout,
	}

	// 应用选项
	for _, opt := range opts {
		opt(s)
	}

	return s
}

// Start 启动加载指示器
func (s *Spinner) Start() {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.stopped {
		// 如果已经停止，重新创建 stopChan
		s.stopChan = make(chan struct{})
		s.stopped = false
		s.currentFrame = 0
	}

	// 隐藏光标
	s.hideCursor()

	go s.run()
}

// Stop 停止加载指示器
func (s *Spinner) Stop() {
	s.mu.Lock()
	defer s.mu.Unlock()

	if !s.stopped {
		close(s.stopChan)
		s.stopped = true
		// 清除当前行
		s.clearLine()
		// 恢复光标
		s.showCursor()
	}
}

// UpdateMessage 更新消息文本
func (s *Spinner) UpdateMessage(message string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.message = message
}

// run 运行加载动画（在 goroutine 中执行）
func (s *Spinner) run() {
	ticker := time.NewTicker(s.interval)
	defer ticker.Stop()

	for {
		select {
		case <-s.stopChan:
			return
		case <-ticker.C:
			s.render()
		}
	}
}

// render 渲染当前帧
func (s *Spinner) render() {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.stopped {
		return
	}

	// 清除当前行
	s.clearLine()

	// 获取当前帧
	frame := s.spinner[s.currentFrame%len(s.spinner)]
	s.currentFrame++

	// 构建输出文本
	theme := GetTheme()
	var text string

	if s.message != "" {
		// 分别渲染 spinner 和文本
		var spinnerText string
		var messageText string

		if theme.EnableColor {
			// 检查是否设置了独立的样式
			if s.spinnerStyle != nil || s.messageStyle != nil {
				// 使用独立的样式
				if s.spinnerStyle != nil {
					spinnerText = s.spinnerStyle.Render(frame)
				} else {
					// 如果没有设置 spinnerStyle，使用默认样式
					spinnerText = s.style.Render(frame)
				}

				if s.messageStyle != nil {
					messageText = s.messageStyle.Render(s.message)
				} else {
					// 如果没有设置 messageStyle，使用默认样式
					messageText = s.style.Render(s.message)
				}
			} else {
				// 使用统一样式
				spinnerText = s.style.Render(frame)
				messageText = s.style.Render(s.message)
			}
		} else {
			// 颜色未启用，直接使用原始文本
			spinnerText = frame
			messageText = s.message
		}

		text = fmt.Sprintf("%s %s", spinnerText, messageText)
	} else {
		// 只有 spinner，没有消息
		if theme.EnableColor {
			if s.spinnerStyle != nil {
				text = s.spinnerStyle.Render(frame)
			} else {
				text = s.style.Render(frame)
			}
		} else {
			text = frame
		}
	}

	// 输出（不换行，使用 \r 回到行首）
	fmt.Fprint(s.writer, text+"\r")
}

// clearLine 清除当前行
func (s *Spinner) clearLine() {
	// 使用 ANSI 转义序列清除当前行
	fmt.Fprint(s.writer, "\033[2K\r")
}

// hideCursor 隐藏光标
func (s *Spinner) hideCursor() {
	if !s.cursorHidden {
		// ANSI 转义序列：隐藏光标
		fmt.Fprint(s.writer, "\033[?25l")
		s.cursorHidden = true
	}
}

// showCursor 显示光标
func (s *Spinner) showCursor() {
	if s.cursorHidden {
		// ANSI 转义序列：显示光标
		fmt.Fprint(s.writer, "\033[?25h")
		s.cursorHidden = false
	}
}

// WithSuccess 停止并显示成功消息
func (s *Spinner) WithSuccess(message string) {
	s.Stop()
	theme := GetTheme()
	formatted := formatMessage("✓", message, theme.SuccessStyle)
	fmt.Fprintln(s.writer, formatted)
}

// WithError 停止并显示错误消息
func (s *Spinner) WithError(message string) {
	s.Stop()
	theme := GetTheme()
	formatted := formatMessage("✗", message, theme.ErrorStyle)
	fmt.Fprintln(s.writer, formatted)
}

// WithInfo 停止并显示信息消息
func (s *Spinner) WithInfo(message string) {
	s.Stop()
	theme := GetTheme()
	formatted := formatMessage("ℹ", message, theme.InfoStyle)
	fmt.Fprintln(s.writer, formatted)
}

// Do 执行一个函数并显示加载状态
func (s *Spinner) Do(fn func() error) error {
	s.Start()
	defer s.Stop()

	err := fn()
	if err != nil {
		return err
	}
	return nil
}

//go:build test

package testutils

import (
	"bytes"
	"context"
	"os"
	"testing"
	"time"

	"github.com/creack/pty"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

// TerminalTester 终端测试器接口
// 这个接口抽象了终端测试功能，方便以后更换底层实现（creack/pty, termtest, tap 等）
type TerminalTester interface {
	// TTY 返回终端文件描述符，用于替换 os.Stdin
	TTY() *os.File

	// WriteInput 写入模拟输入
	WriteInput(t *testing.T, data string) error

	// ReadOutput 读取捕获的输出
	ReadOutput() string

	// Close 清理资源
	Close() error
}

// PTYTerminalTester 基于 creack/pty 的终端测试器实现
type PTYTerminalTester struct {
	pty       *os.File
	tty       *os.File
	oldStdin  *os.File
	oldStdout *os.File
	stdoutR   *os.File
	stdoutW   *os.File
	stdoutBuf *bytes.Buffer
	done      chan struct{}
}

// NewTerminalTester 创建新的终端测试器
// 使用 creack/pty 作为底层实现
func NewTerminalTester(t *testing.T) TerminalTester {
	// 创建伪终端
	ptyFile, ttyFile, err := pty.Open()
	require.NoError(t, err, "创建伪终端失败")

	// 保存原始的 stdin/stdout
	oldStdin := os.Stdin
	oldStdout := os.Stdout

	// 替换 stdin
	os.Stdin = ttyFile

	// 创建管道来捕获 stdout
	stdoutR, stdoutW, err := os.Pipe()
	require.NoError(t, err, "创建 stdout 管道失败")

	// 替换 stdout
	os.Stdout = stdoutW

	// 创建输出缓冲区
	stdoutBuf := &bytes.Buffer{}
	done := make(chan struct{})

	// 在 goroutine 中读取 stdout
	go func() {
		defer close(done)
		buf := make([]byte, 4096)
		for {
			n, err := stdoutR.Read(buf)
			if n > 0 {
				stdoutBuf.Write(buf[:n])
			}
			if err != nil {
				break
			}
		}
	}()

	return &PTYTerminalTester{
		pty:       ptyFile,
		tty:       ttyFile,
		oldStdin:  oldStdin,
		oldStdout: oldStdout,
		stdoutR:   stdoutR,
		stdoutW:   stdoutW,
		stdoutBuf: stdoutBuf,
		done:      done,
	}
}

// TTY 返回终端文件描述符
func (pt *PTYTerminalTester) TTY() *os.File {
	return pt.tty
}

// WriteInput 写入模拟输入
// 注意：为了确保输入能被正确读取，建议在 prompt 函数调用之后才调用此方法
// 或者在 goroutine 中添加小延迟，确保终端模式设置完成
func (pt *PTYTerminalTester) WriteInput(t *testing.T, data string) error {
	// 确保 PTY 准备就绪
	// 给终端一些时间来设置原始模式（如果正在设置的话）
	// 这是一个最小延迟，确保写入操作在正确的时机执行
	_, err := pt.pty.WriteString(data)
	require.NoError(t, err, "写入输入失败")
	return err
}

// ReadOutput 读取捕获的输出
func (pt *PTYTerminalTester) ReadOutput() string {
	return pt.stdoutBuf.String()
}

// Close 清理资源
func (pt *PTYTerminalTester) Close() error {
	var errs []error

	// 关闭 stdout 写入端，让读取端知道结束
	if pt.stdoutW != nil {
		if err := pt.stdoutW.Close(); err != nil {
			errs = append(errs, err)
		}
	}

	// 等待读取完成
	if pt.done != nil {
		<-pt.done
	}

	// 关闭 stdout 读取端
	if pt.stdoutR != nil {
		if err := pt.stdoutR.Close(); err != nil {
			errs = append(errs, err)
		}
	}

	if pt.pty != nil {
		if err := pt.pty.Close(); err != nil {
			errs = append(errs, err)
		}
	}

	if pt.tty != nil {
		if err := pt.tty.Close(); err != nil {
			errs = append(errs, err)
		}
	}

	// 恢复原始的 stdin/stdout
	os.Stdin = pt.oldStdin
	os.Stdout = pt.oldStdout

	if len(errs) > 0 {
		return errs[0]
	}
	return nil
}

// WithTerminal 在终端测试环境中执行函数
// 这是一个便捷函数，自动处理资源清理
func WithTerminal(t *testing.T, fn func(tester TerminalTester)) {
	tester := NewTerminalTester(t)
	defer func() {
		require.NoError(t, tester.Close(), "清理终端测试器失败")
	}()

	fn(tester)
}

// PromptTestBuilder 提示测试构建器
// 用于简化终端交互测试的编写
type PromptTestBuilder struct {
	inputs      []string          // 输入序列
	delays      []time.Duration   // 每次输入后的延迟（可选）
	startDelay  time.Duration     // 开始输入前的延迟
	timeout     time.Duration     // 测试超时时间（可选，默认 5 秒）
	tester      TerminalTester    // 终端测试器（可选，如果不提供会自动创建）
	outputCheck func(string) bool // 输出验证函数（可选）
}

// NewPromptTest 创建新的提示测试构建器
func NewPromptTest() *PromptTestBuilder {
	return &PromptTestBuilder{
		inputs:      make([]string, 0),
		delays:      make([]time.Duration, 0),
		startDelay:  50 * time.Millisecond, // 默认延迟，确保 prompt 函数准备好
		timeout:     5 * time.Second,       // 默认超时 5 秒
		outputCheck: nil,
	}
}

// WithInput 添加单个输入（自动添加换行符）
// 例如：WithInput("y") 会发送 "y\n"
// 对于空字符串，会发送 "\n"（回车）
func (b *PromptTestBuilder) WithInput(input string) *PromptTestBuilder {
	// 对于空字符串，直接添加换行符
	if input == "" {
		b.inputs = append(b.inputs, "\n")
		return b
	}
	// 确保输入以换行符结尾
	if input[len(input)-1] != '\n' {
		input += "\n"
	}
	b.inputs = append(b.inputs, input)
	return b
}

// WithInputs 添加多个输入
func (b *PromptTestBuilder) WithInputs(inputs ...string) *PromptTestBuilder {
	for _, input := range inputs {
		b.WithInput(input)
	}
	return b
}

// WithRawInput 添加原始输入（不自动添加换行符）
// 用于转义序列等特殊输入
func (b *PromptTestBuilder) WithRawInput(input string) *PromptTestBuilder {
	b.inputs = append(b.inputs, input)
	return b
}

// WithStartDelay 设置开始输入前的延迟
// 默认 50ms，确保 prompt 函数设置好终端模式
func (b *PromptTestBuilder) WithStartDelay(delay time.Duration) *PromptTestBuilder {
	b.startDelay = delay
	return b
}

// WithDelay 为上次添加的输入设置延迟
// 用于在多个输入之间添加延迟
func (b *PromptTestBuilder) WithDelay(delay time.Duration) *PromptTestBuilder {
	b.delays = append(b.delays, delay)
	return b
}

// WithInputAndDelay 添加输入并设置延迟（便捷方法）
func (b *PromptTestBuilder) WithInputAndDelay(input string, delay time.Duration) *PromptTestBuilder {
	b.WithInput(input)
	b.WithDelay(delay)
	return b
}

// WithOutputCheck 设置输出验证函数
func (b *PromptTestBuilder) WithOutputCheck(check func(string) bool) *PromptTestBuilder {
	b.outputCheck = check
	return b
}

// WithTimeout 设置测试超时时间
// 默认 5 秒，如果测试在超时时间内未完成，会返回超时错误
func (b *PromptTestBuilder) WithTimeout(timeout time.Duration) *PromptTestBuilder {
	b.timeout = timeout
	return b
}

// Run 执行提示测试
// fn: 执行 prompt 调用的函数，返回结果和错误
// 返回: prompt 函数的结果和错误
func (b *PromptTestBuilder) Run(t *testing.T, fn func() (interface{}, error)) (interface{}, error) {
	return b.RunWithTester(t, nil, fn)
}

// RunWithTester 使用指定的终端测试器执行提示测试
// 如果不提供 tester，会自动创建并在测试结束时清理
func (b *PromptTestBuilder) RunWithTester(t *testing.T, tester TerminalTester, fn func() (interface{}, error)) (interface{}, error) {
	// 如果没有提供 tester，创建一个新的
	createdTester := tester == nil
	if createdTester {
		tester = NewTerminalTester(t)
		t.Cleanup(func() {
			if err := tester.Close(); err != nil {
				t.Logf("清理终端测试器失败: %v", err)
			}
		})
	}

	// 保存原来的 tester 以便输出验证
	b.tester = tester

	// 创建带超时的 context
	ctx, cancel := context.WithTimeout(context.Background(), b.timeout)
	defer cancel()

	// 创建结果和错误通道
	type testResult struct {
		value interface{}
		err   error
	}
	resultChan := make(chan testResult, 1)

	// 在 goroutine 中写入输入
	go func() {
		// 等待 prompt 函数设置终端模式
		if b.startDelay > 0 {
			time.Sleep(b.startDelay)
		}

		// 依次写入所有输入
		for i, input := range b.inputs {
			// 检查是否已超时
			select {
			case <-ctx.Done():
				return
			default:
			}

			// 在 goroutine 中使用 t.Logf 而不是 require，避免测试框架问题
			if err := tester.WriteInput(t, input); err != nil {
				t.Logf("写入输入失败: %v", err)
				return
			}

			// 如果设置了延迟，等待
			if i < len(b.delays) && b.delays[i] > 0 {
				select {
				case <-ctx.Done():
					return
				case <-time.After(b.delays[i]):
				}
			} else if i < len(b.inputs)-1 {
				// 如果下一个输入是转义序列，添加默认延迟
				nextInput := b.inputs[i+1]
				if len(nextInput) > 0 && nextInput[0] == 0x1b {
					select {
					case <-ctx.Done():
						return
					case <-time.After(10 * time.Millisecond):
					}
				}
			}
		}
	}()

	// 在另一个 goroutine 中执行 prompt 函数
	go func() {
		value, err := fn()
		select {
		case resultChan <- testResult{value: value, err: err}:
		case <-ctx.Done():
			// 如果已超时，忽略结果
		}
	}()

	// 等待结果或超时
	select {
	case res := <-resultChan:
		// 如果设置了输出验证，执行验证
		if b.outputCheck != nil && b.tester != nil {
			output := b.tester.ReadOutput()
			if !b.outputCheck(output) {
				t.Errorf("输出验证失败: %s", output)
			}
		}
		return res.value, res.err
	case <-ctx.Done():
		// 超时，返回超时错误
		t.Errorf("测试超时: %v", b.timeout)
		return nil, context.DeadlineExceeded
	}
}

// RunPromptTest 便捷函数，用于快速执行简单的提示测试
// 示例：
//
//	result, err := testutils.RunPromptTest(t, "y", func() (interface{}, error) {
//	    return prompt.AskConfirm("是否继续？", false)
//	})
func RunPromptTest(t *testing.T, input string, fn func() (interface{}, error)) (interface{}, error) {
	return NewPromptTest().WithInput(input).Run(t, fn)
}

// BasePromptTestSuite 是所有 prompt 测试套件的基类
// 提供通用的 runPromptTest 方法，消除各个测试套件中的重复代码
//
// 使用方法：
//
//	type InputTestSuite struct {
//	    testutils.BasePromptTestSuite
//	}
//
//	func (s *InputTestSuite) TestSomething() {
//	    result, err := s.RunPromptTest(
//	        func(pt *PromptTestBuilder) *PromptTestBuilder {
//	            return pt.WithInput("test")
//	        },
//	        func() (interface{}, error) {
//	            return prompt.AskInput("请输入", "")
//	        },
//	    )
//	    s.NoError(err)
//	    s.Equal("test", result)
//	}
type BasePromptTestSuite struct {
	suite.Suite
}

// RunPromptTest 执行 prompt 测试的辅助方法
// 这是所有 prompt 测试套件共享的方法，消除了重复代码
//
// 参数：
//   - builder: 可选，用于配置 PromptTestBuilder 的函数
//   - fn: 执行 prompt 调用的函数，返回结果和错误
//
// 返回：
//   - 结果和错误
func (s *BasePromptTestSuite) RunPromptTest(
	builder func(*PromptTestBuilder) *PromptTestBuilder,
	fn func() (interface{}, error),
) (interface{}, error) {
	pt := NewPromptTest()
	if builder != nil {
		pt = builder(pt)
	}
	return pt.Run(s.T(), fn)
}

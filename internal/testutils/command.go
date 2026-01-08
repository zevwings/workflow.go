//go:build test

package testutils

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

// CommandResult 命令执行结果
type CommandResult struct {
	Stdout string
	Stderr string
	Err    error
}

// ExecuteCommand 执行 CLI 命令
//
// 执行命令并返回标准输出。
//
// 参数:
//   - t: 测试对象
//   - command: 命令名称
//   - args: 命令参数
//
// 返回:
//   - string: 标准输出
//   - error: 执行错误
func ExecuteCommand(t *testing.T, command string, args ...string) (string, error) {
	cmd := exec.Command(command, args...)
	output, err := cmd.Output()
	return string(output), err
}

// ExecuteCommandWithEnv 执行 CLI 命令（带环境变量）
//
// 执行命令并返回标准输出，使用指定的环境变量。
//
// 参数:
//   - t: 测试对象
//   - env: 环境变量映射
//   - command: 命令名称
//   - args: 命令参数
//
// 返回:
//   - string: 标准输出
//   - error: 执行错误
func ExecuteCommandWithEnv(t *testing.T, env map[string]string, command string, args ...string) (string, error) {
	cmd := exec.Command(command, args...)

	// 设置环境变量
	cmd.Env = os.Environ()
	for k, v := range env {
		cmd.Env = append(cmd.Env, k+"="+v)
	}

	output, err := cmd.Output()
	return string(output), err
}

// ExecuteCommandWithDir 执行 CLI 命令（带工作目录）
//
// 执行命令并返回标准输出，使用指定的工作目录。
//
// 参数:
//   - t: 测试对象
//   - dir: 工作目录
//   - command: 命令名称
//   - args: 命令参数
//
// 返回:
//   - string: 标准输出
//   - error: 执行错误
func ExecuteCommandWithDir(t *testing.T, dir string, command string, args ...string) (string, error) {
	// 确保目录存在
	absDir, err := filepath.Abs(dir)
	require.NoError(t, err, "获取绝对路径失败: %s", dir)
	require.DirExists(t, absDir, "工作目录不存在: %s", absDir)

	cmd := exec.Command(command, args...)
	cmd.Dir = absDir

	output, err := cmd.Output()
	return string(output), err
}

// ExecuteCommandCapture 执行 CLI 命令并捕获输出
//
// 执行命令并捕获标准输出、标准错误和错误信息。
//
// 参数:
//   - t: 测试对象
//   - command: 命令名称
//   - args: 命令参数
//
// 返回:
//   - *CommandResult: 命令执行结果
func ExecuteCommandCapture(t *testing.T, command string, args ...string) *CommandResult {
	cmd := exec.Command(command, args...)

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()

	return &CommandResult{
		Stdout: stdout.String(),
		Stderr: stderr.String(),
		Err:    err,
	}
}

// ExecuteCommandCaptureWithEnv 执行 CLI 命令并捕获输出（带环境变量）
//
// 执行命令并捕获标准输出、标准错误和错误信息，使用指定的环境变量。
//
// 参数:
//   - t: 测试对象
//   - env: 环境变量映射
//   - command: 命令名称
//   - args: 命令参数
//
// 返回:
//   - *CommandResult: 命令执行结果
func ExecuteCommandCaptureWithEnv(t *testing.T, env map[string]string, command string, args ...string) *CommandResult {
	cmd := exec.Command(command, args...)

	// 设置环境变量
	cmd.Env = os.Environ()
	for k, v := range env {
		cmd.Env = append(cmd.Env, k+"="+v)
	}

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()

	return &CommandResult{
		Stdout: stdout.String(),
		Stderr: stderr.String(),
		Err:    err,
	}
}

// ExecuteCommandCaptureWithDir 执行 CLI 命令并捕获输出（带工作目录）
//
// 执行命令并捕获标准输出、标准错误和错误信息，使用指定的工作目录。
//
// 参数:
//   - t: 测试对象
//   - dir: 工作目录
//   - command: 命令名称
//   - args: 命令参数
//
// 返回:
//   - *CommandResult: 命令执行结果
func ExecuteCommandCaptureWithDir(t *testing.T, dir string, command string, args ...string) *CommandResult {
	// 确保目录存在
	absDir, err := filepath.Abs(dir)
	require.NoError(t, err, "获取绝对路径失败: %s", dir)
	require.DirExists(t, absDir, "工作目录不存在: %s", absDir)

	cmd := exec.Command(command, args...)
	cmd.Dir = absDir

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err = cmd.Run()

	return &CommandResult{
		Stdout: stdout.String(),
		Stderr: stderr.String(),
		Err:    err,
	}
}

package commands

import (
	"context"
	"os/exec"
	"time"

	"github.com/spf13/cobra"
	"github.com/zevwings/workflow/internal/config"
	"github.com/zevwings/workflow/internal/http"
	"github.com/zevwings/workflow/internal/prompt"
)

// NewCheckCmd 创建 check 命令
func NewCheckCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "check",
		Short: "运行环境检查",
		Long: `运行环境检查，包括：
- Git 状态检查
- 网络连接检查
- 配置文件检查`,
		RunE: runCheck,
	}

	return cmd
}

func runCheck(cmd *cobra.Command, args []string) error {
	out := prompt.NewMessage(false)
	out.Info("开始环境检查...")
	out.Println("")

	table := prompt.NewTable([]string{"检查项", "状态", "说明"})

	// 1. 检查 Git
	gitOK := checkGit(out)
	if gitOK {
		table.AddRow([]string{"Git", "✓", "Git 已安装"})
	} else {
		table.AddRow([]string{"Git", "✗", "Git 未安装或不在 PATH 中"})
	}

	// 2. 检查 Git 仓库
	repoOK := checkGitRepo(out)
	if repoOK {
		table.AddRow([]string{"Git 仓库", "✓", "当前目录是 Git 仓库"})
	} else {
		table.AddRow([]string{"Git 仓库", "✗", "当前目录不是 Git 仓库"})
	}

	// 3. 检查配置文件
	configOK := checkConfig(out)
	if configOK {
		table.AddRow([]string{"配置文件", "✓", "配置文件存在且有效"})
	} else {
		table.AddRow([]string{"配置文件", "✗", "配置文件不存在或无效"})
	}

	// 4. 检查网络连接
	networkOK := checkNetwork(out)
	if networkOK {
		table.AddRow([]string{"网络连接", "✓", "网络连接正常"})
	} else {
		table.AddRow([]string{"网络连接", "✗", "网络连接失败"})
	}

	out.Println("")
	table.Render()
	out.Println("")

	if gitOK && repoOK && configOK && networkOK {
		out.Success("所有检查通过")
		return nil
	}

	out.Warning("部分检查未通过，请根据上述信息进行修复")
	return nil
}

func checkGit(out *prompt.Message) bool {
	cmd := exec.Command("git", "version")
	if err := cmd.Run(); err != nil {
		return false
	}
	return true
}

func checkGitRepo(out *prompt.Message) bool {
	cmd := exec.Command("git", "rev-parse", "--git-dir")
	if err := cmd.Run(); err != nil {
		return false
	}
	return true
}

func checkConfig(out *prompt.Message) bool {
	manager, err := config.Global()
	if err != nil {
		return false
	}

	if err := manager.Load(); err != nil {
		return false
	}

	return true
}

func checkNetwork(out *prompt.Message) bool {
	client := http.Global()

	// 检查 GitHub 连接（使用超时上下文）
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 使用 resty 客户端发送请求
	restyClient := client.GetRestyClient()
	resp, err := restyClient.R().SetContext(ctx).Get("https://api.github.com")
	if err != nil {
		return false
	}

	// 检查响应状态码
	return resp.StatusCode() == 200
}

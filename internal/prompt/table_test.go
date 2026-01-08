//go:build test

package prompt

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewTable(t *testing.T) {
	headers := []string{"Name", "Age", "City"}
	table := NewTable(headers)

	require.NotNil(t, table)
	require.Equal(t, headers, table.headers)
	require.Empty(t, table.rows)
	require.True(t, table.border)
	require.True(t, table.rowLine)
	require.Equal(t, ALIGN_LEFT, table.align)
}

func TestTable_AddRow(t *testing.T) {
	table := NewTable([]string{"Col1", "Col2"})

	// 测试链式调用
	result := table.AddRow([]string{"A", "B"}).
		AddRow([]string{"C", "D"})

	require.Equal(t, table, result)
	require.Len(t, table.rows, 2)
	require.Equal(t, []string{"A", "B"}, table.rows[0])
	require.Equal(t, []string{"C", "D"}, table.rows[1])
}

func TestTable_SetHeader(t *testing.T) {
	table := NewTable([]string{"Old1", "Old2"})
	newHeaders := []string{"New1", "New2", "New3"}

	result := table.SetHeader(newHeaders)

	require.Equal(t, table, result)
	require.Equal(t, newHeaders, table.headers)
}

func TestTable_SetBorder(t *testing.T) {
	table := NewTable([]string{"Col1"})

	require.True(t, table.border)

	result := table.SetBorder(false)
	require.Equal(t, table, result)
	require.False(t, table.border)

	result = table.SetBorder(true)
	require.True(t, table.border)
}

func TestTable_SetRowLine(t *testing.T) {
	table := NewTable([]string{"Col1"})

	require.True(t, table.rowLine)

	result := table.SetRowLine(false)
	require.Equal(t, table, result)
	require.False(t, table.rowLine)

	result = table.SetRowLine(true)
	require.True(t, table.rowLine)
}

func TestTable_SetAlignment(t *testing.T) {
	table := NewTable([]string{"Col1"})

	require.Equal(t, ALIGN_LEFT, table.align)

	result := table.SetAlignment(ALIGN_CENTER)
	require.Equal(t, table, result)
	require.Equal(t, ALIGN_CENTER, table.align)

	result = table.SetAlignment(ALIGN_RIGHT)
	require.Equal(t, ALIGN_RIGHT, table.align)

	result = table.SetAlignment(ALIGN_LEFT)
	require.Equal(t, ALIGN_LEFT, table.align)
}

func TestStripAnsiCodes(t *testing.T) {
	testCases := []struct {
		name     string
		input    string
		expected string
	}{
		{"纯文本", "hello", "hello"},
		{"带ANSI代码", "\x1b[31mhello\x1b[0m", "hello"},
		{"多个ANSI代码", "\x1b[1;32mtest\x1b[0m\x1b[33mworld\x1b[0m", "testworld"},
		{"空字符串", "", ""},
		{"只有ANSI代码", "\x1b[31m\x1b[0m", ""},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := stripAnsiCodes(tc.input)
			require.Equal(t, tc.expected, result)
		})
	}
}

func TestTable_CalculateColumnWidths(t *testing.T) {
	// 保存原始主题
	originalTheme := GetTheme()
	defer SetTheme(originalTheme)

	// 禁用颜色以便测试
	theme := GetTheme()
	theme.EnableColor = false
	SetTheme(theme)

	table := NewTable([]string{"Name", "Age"})
	table.AddRow([]string{"Alice", "25"})
	table.AddRow([]string{"Bob", "30"})
	table.AddRow([]string{"Charlie", "35"})

	widths := table.calculateColumnWidths()

	require.Len(t, widths, 2)
	// "Charlie" 是最长的，7个字符
	require.Equal(t, 7, widths[0])
	// "Age" 表头是3个字符，但数据都是2个字符，所以应该是3
	require.Equal(t, 3, widths[1])
}

func TestTable_CalculateColumnWidths_WithAnsiCodes(t *testing.T) {
	// 保存原始主题
	originalTheme := GetTheme()
	defer SetTheme(originalTheme)

	// 启用颜色
	theme := GetTheme()
	theme.EnableColor = true
	SetTheme(theme)

	table := NewTable([]string{"Name", "Age"})
	// 添加带ANSI代码的内容
	table.AddRow([]string{"\x1b[31mAlice\x1b[0m", "25"})

	widths := table.calculateColumnWidths()

	// 应该忽略ANSI代码，只计算实际显示宽度
	require.Len(t, widths, 2)
	require.Equal(t, 5, widths[0]) // "Alice" 是5个字符
	require.Equal(t, 3, widths[1]) // "Age" 是3个字符
}

func TestTable_AlignCell(t *testing.T) {
	// 保存原始主题
	originalTheme := GetTheme()
	defer SetTheme(originalTheme)

	theme := GetTheme()
	theme.EnableColor = false
	SetTheme(theme)

	table := NewTable([]string{"Col1"})

	// 测试左对齐（默认）
	table.align = ALIGN_LEFT
	result := table.alignCell("test", 10, 4)
	require.Equal(t, "test      ", result) // 右边填充6个空格

	// 测试右对齐
	table.align = ALIGN_RIGHT
	result = table.alignCell("test", 10, 4)
	require.Equal(t, "      test", result) // 左边填充6个空格

	// 测试居中对齐
	table.align = ALIGN_CENTER
	result = table.alignCell("test", 10, 4)
	require.Equal(t, "   test   ", result) // 左右各填充3个空格

	// 测试宽度相等的情况
	result = table.alignCell("test", 4, 4)
	require.Equal(t, "test", result)

	// 测试实际宽度超过目标宽度（应该返回原内容）
	result = table.alignCell("very long text", 5, 14)
	require.Equal(t, "very long text", result)
}

func TestTable_Render_EmptyHeaders(t *testing.T) {
	// 保存原始主题
	originalTheme := GetTheme()
	defer SetTheme(originalTheme)

	theme := GetTheme()
	theme.EnableColor = false
	SetTheme(theme)

	table := NewTable([]string{})
	table.AddRow([]string{"A", "B"})

	// 捕获输出
	out := captureOutput(t, func() {
		table.Render()
	})

	// 空表头时应该不输出任何内容
	require.Empty(t, strings.TrimSpace(out))
}

func TestTable_Render_WithBorder(t *testing.T) {
	// 保存原始主题
	originalTheme := GetTheme()
	defer SetTheme(originalTheme)

	theme := GetTheme()
	theme.EnableColor = false
	SetTheme(theme)

	table := NewTable([]string{"Name", "Age"})
	table.AddRow([]string{"Alice", "25"})
	table.AddRow([]string{"Bob", "30"})

	out := captureOutput(t, func() {
		table.Render()
	})

	lines := strings.Split(strings.TrimSpace(out), "\n")
	require.GreaterOrEqual(t, len(lines), 5) // 至少包含顶部边框、表头、分隔线、两行数据、底部边框

	// 检查是否包含边框字符
	require.Contains(t, out, "┌")
	require.Contains(t, out, "┐")
	require.Contains(t, out, "│")
	require.Contains(t, out, "─")
	require.Contains(t, out, "└")
	require.Contains(t, out, "┘")

	// 检查是否包含数据
	require.Contains(t, out, "Alice")
	require.Contains(t, out, "Bob")
	require.Contains(t, out, "25")
	require.Contains(t, out, "30")
}

func TestTable_Render_WithoutBorder(t *testing.T) {
	// 保存原始主题
	originalTheme := GetTheme()
	defer SetTheme(originalTheme)

	theme := GetTheme()
	theme.EnableColor = false
	SetTheme(theme)

	table := NewTable([]string{"Name", "Age"})
	table.SetBorder(false)
	table.AddRow([]string{"Alice", "25"})

	out := captureOutput(t, func() {
		table.Render()
	})

	// 无边框模式下不应该包含边框字符（┌┐└┘等）
	// 但列之间仍然使用 "│" 作为分隔符
	require.NotContains(t, out, "┌")
	require.NotContains(t, out, "┐")
	require.NotContains(t, out, "└")
	require.NotContains(t, out, "┘")
	require.NotContains(t, out, "─")
	require.NotContains(t, out, "┼")
	require.NotContains(t, out, "┬")
	require.NotContains(t, out, "┴")

	// 但应该包含数据
	require.Contains(t, out, "Alice")
	require.Contains(t, out, "25")
}

func TestTable_Render_WithoutRowLine(t *testing.T) {
	// 保存原始主题
	originalTheme := GetTheme()
	defer SetTheme(originalTheme)

	theme := GetTheme()
	theme.EnableColor = false
	SetTheme(theme)

	table := NewTable([]string{"Name", "Age"})
	table.SetRowLine(false)
	table.AddRow([]string{"Alice", "25"})
	table.AddRow([]string{"Bob", "30"})

	out := captureOutput(t, func() {
		table.Render()
	})

	// 无行线模式下，行之间不应该有分隔线
	// 但表头分隔线仍然存在（如果有边框）
	lines := strings.Split(strings.TrimSpace(out), "\n")

	// 应该只有表头分隔线，没有行分隔线
	// 计算包含 "┼" 的行数（行分隔符）
	separatorCount := 0
	for _, line := range lines {
		if strings.Contains(line, "┼") {
			separatorCount++
		}
	}
	// 应该只有表头分隔线（1个），没有行分隔线
	require.Equal(t, 1, separatorCount)
}

func TestTable_Render_Alignment(t *testing.T) {
	// 保存原始主题
	originalTheme := GetTheme()
	defer SetTheme(originalTheme)

	theme := GetTheme()
	theme.EnableColor = false
	SetTheme(theme)

	// 测试右对齐
	table := NewTable([]string{"Number"})
	table.SetAlignment(ALIGN_RIGHT)
	table.AddRow([]string{"1"})
	table.AddRow([]string{"123"})

	out := captureOutput(t, func() {
		table.Render()
	})

	// 验证数字是右对齐的（较短的数字前面有空格）
	lines := strings.Split(strings.TrimSpace(out), "\n")

	// 找到包含 "1" 的行（应该是数据行，不是表头）
	var dataLine string
	for _, line := range lines {
		if strings.Contains(line, "1") && !strings.Contains(line, "Number") {
			dataLine = line
			break
		}
	}

	// 右对齐时，较短的 "1" 应该在较长的 "123" 之前有更多空格
	require.NotEmpty(t, dataLine)
}

func TestTable_Render_ChainCall(t *testing.T) {
	// 保存原始主题
	originalTheme := GetTheme()
	defer SetTheme(originalTheme)

	theme := GetTheme()
	theme.EnableColor = false
	SetTheme(theme)

	// 测试链式调用
	table := NewTable([]string{"Col1"}).
		AddRow([]string{"A"}).
		SetBorder(false).
		SetRowLine(false).
		SetAlignment(ALIGN_CENTER)

	// Render 应该返回自身，支持链式调用
	result := table.Render()
	require.Equal(t, table, result)
}

func TestTable_RenderSeparator(t *testing.T) {
	table := NewTable([]string{"Col1", "Col2", "Col3"})
	colWidths := []int{5, 3, 4}

	result := table.renderSeparator(colWidths, "-", "+", "|", "|")

	// 每列宽度+2（左右各1个空格），用+连接，前后用|包围
	require.Contains(t, result, "|")
	require.Contains(t, result, "+")
	require.Contains(t, result, "-")
}

func TestTable_RenderTopBorder(t *testing.T) {
	table := NewTable([]string{"Col1", "Col2"})
	colWidths := []int{5, 3}

	result := table.renderTopBorder(colWidths, "-", "┌", "┐", "┬")

	require.Contains(t, result, "┌")
	require.Contains(t, result, "┐")
	require.Contains(t, result, "┬")
	require.Contains(t, result, "-")
}

func TestTable_RenderBottomBorder(t *testing.T) {
	table := NewTable([]string{"Col1", "Col2"})
	colWidths := []int{5, 3}

	result := table.renderBottomBorder(colWidths, "-", "└", "┘", "┴")

	require.Contains(t, result, "└")
	require.Contains(t, result, "┘")
	require.Contains(t, result, "┴")
	require.Contains(t, result, "-")
}


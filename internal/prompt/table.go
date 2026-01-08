package prompt

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/mattn/go-runewidth"
)

// Table 表格工具
type Table struct {
	headers []string
	rows    [][]string
	border  bool
	rowLine bool
	align   Alignment
	theme   Theme
}

// Alignment 对齐方式
type Alignment int

const (
	ALIGN_LEFT Alignment = iota
	ALIGN_CENTER
	ALIGN_RIGHT
)

// NewTable 创建新的表格
func NewTable(headers []string) *Table {
	return &Table{
		headers: headers,
		rows:    make([][]string, 0),
		border:  true,
		rowLine: true,
		align:   ALIGN_LEFT,
		theme:   GetTheme(),
	}
}

// AddRow 添加行（支持链式调用）
func (t *Table) AddRow(row []string) *Table {
	t.rows = append(t.rows, row)
	return t
}

// Render 渲染表格（支持链式调用）
func (t *Table) Render() *Table {
	if len(t.headers) == 0 {
		return t
	}

	// 计算每列的最大宽度
	colWidths := t.calculateColumnWidths()

	// 使用主题中的边框样式
	borderStyle := t.theme.BorderStyle
	if !t.theme.EnableColor {
		borderStyle = lipgloss.NewStyle()
	}

	// 边框字符
	vertical := "│"
	horizontal := "─"
	cross := "┼"
	topLeft := "┌"
	topRight := "┐"
	bottomLeft := "└"
	bottomRight := "┘"
	topCross := "┬"
	bottomCross := "┴"
	leftCross := "├"
	rightCross := "┤"

	if !t.border {
		vertical = " "
		horizontal = " "
		cross = " "
		topLeft = " "
		topRight = " "
		bottomLeft = " "
		bottomRight = " "
		topCross = " "
		bottomCross = " "
		leftCross = " "
		rightCross = " "
	}

	// 渲染表头
	headerRow := t.renderRow(t.headers, colWidths, true, borderStyle)
	headerSeparator := t.renderSeparator(colWidths, horizontal, cross, leftCross, rightCross)

	// 组装完整表格
	var lines []string

	// 顶部边框
	if t.border {
		topBorder := t.renderTopBorder(colWidths, horizontal, topLeft, topRight, topCross)
		lines = append(lines, borderStyle.Render(topBorder))
	}

	// 表头
	var headerLine string
	if t.border {
		headerLine = borderStyle.Render(vertical) + " " + headerRow + " " + borderStyle.Render(vertical)
	} else {
		headerLine = " " + headerRow + " "
	}
	lines = append(lines, headerLine)

	// 表头分隔线
	if t.border {
		lines = append(lines, borderStyle.Render(headerSeparator))
	}

	// 渲染数据行
	for i, row := range t.rows {
		// 数据行
		dataRow := t.renderRow(row, colWidths, false, borderStyle)
		var rowLine string
		if t.border {
			rowLine = borderStyle.Render(vertical) + " " + dataRow + " " + borderStyle.Render(vertical)
		} else {
			rowLine = " " + dataRow + " "
		}
		lines = append(lines, rowLine)

		// 行分隔线（最后一行后不添加）
		if t.rowLine && i < len(t.rows)-1 {
			if t.border {
				rowSeparator := t.renderSeparator(colWidths, horizontal, cross, leftCross, rightCross)
				lines = append(lines, borderStyle.Render(rowSeparator))
			}
			// 无边框模式下，不添加行分隔线（保持简洁）
		}
	}

	// 底部边框
	if t.border {
		bottomBorder := t.renderBottomBorder(colWidths, horizontal, bottomLeft, bottomRight, bottomCross)
		lines = append(lines, borderStyle.Render(bottomBorder))
	}

	// 输出表格
	fmt.Println(strings.Join(lines, "\n"))
	return t
}

// stripAnsiCodes 去除 ANSI 转义码，返回纯文本（用于计算显示宽度）
func stripAnsiCodes(s string) string {
	ansiRegex := regexp.MustCompile(`\x1b\[[0-9;]*[a-zA-Z]`)
	return ansiRegex.ReplaceAllString(s, "")
}

// calculateColumnWidths 计算每列的最大宽度
func (t *Table) calculateColumnWidths() []int {
	colWidths := make([]int, len(t.headers))

	// 计算表头宽度（使用 runewidth 计算实际显示宽度）
	for i, header := range t.headers {
		// 去除可能的 ANSI 代码后计算宽度
		cleanHeader := stripAnsiCodes(header)
		width := runewidth.StringWidth(cleanHeader)
		if width > colWidths[i] {
			colWidths[i] = width
		}
	}

	// 计算数据行宽度
	for _, row := range t.rows {
		for i := 0; i < len(row) && i < len(colWidths); i++ {
			// 去除可能的 ANSI 代码后计算宽度
			cleanCell := stripAnsiCodes(row[i])
			width := runewidth.StringWidth(cleanCell)
			if width > colWidths[i] {
				colWidths[i] = width
			}
		}
	}

	// 确保最小宽度为 1
	for i := range colWidths {
		if colWidths[i] < 1 {
			colWidths[i] = 1
		}
	}

	return colWidths
}

// renderRow 渲染一行数据
func (t *Table) renderRow(row []string, colWidths []int, isHeader bool, borderStyle lipgloss.Style) string {
	var cells []string

	for i := 0; i < len(colWidths); i++ {
		var cell string
		if i < len(row) {
			cell = row[i]
		}

		// 先去除 ANSI 代码计算实际宽度，然后对齐
		cleanCell := stripAnsiCodes(cell)
		actualWidth := runewidth.StringWidth(cleanCell)

		// 对齐处理（基于实际显示宽度）
		cell = t.alignCell(cell, colWidths[i], actualWidth)

		// 表头样式（在对齐后应用，这样不会影响宽度计算）
		if isHeader && t.theme.EnableColor {
			cell = t.theme.TitleStyle.Render(cell)
		}

		cells = append(cells, cell)
	}

	// 根据是否有边框选择分隔符
	var separator string
	if !t.border {
		// 无边框模式：使用简单的分隔符
		separator = " │ "
	} else if t.theme.EnableColor {
		// 有边框且启用颜色：使用带样式的分隔符
		separator = " " + borderStyle.Render("│") + " "
	} else {
		// 有边框但未启用颜色：使用普通分隔符
		separator = " │ "
	}
	return strings.Join(cells, separator)
}

// alignCell 对齐单元格内容
// cell: 原始单元格内容（可能包含 ANSI 代码）
// targetWidth: 目标宽度（列宽）
// actualWidth: 实际显示宽度（已去除 ANSI 代码）
func (t *Table) alignCell(cell string, targetWidth int, actualWidth int) string {
	if actualWidth >= targetWidth {
		// 如果实际宽度超过目标宽度，需要截断
		// 这里简单处理：返回原内容（实际应该截断，但为了简化先这样）
		return cell
	}

	padding := targetWidth - actualWidth
	switch t.align {
	case ALIGN_CENTER:
		leftPad := padding / 2
		rightPad := padding - leftPad
		return strings.Repeat(" ", leftPad) + cell + strings.Repeat(" ", rightPad)
	case ALIGN_RIGHT:
		return strings.Repeat(" ", padding) + cell
	default: // ALIGN_LEFT
		return cell + strings.Repeat(" ", padding)
	}
}

// renderSeparator 渲染分隔线
func (t *Table) renderSeparator(colWidths []int, horizontal, cross, left, right string) string {
	var parts []string
	for _, width := range colWidths {
		parts = append(parts, strings.Repeat(horizontal, width+2))
	}
	return left + strings.Join(parts, cross) + right
}

// renderTopBorder 渲染顶部边框
func (t *Table) renderTopBorder(colWidths []int, horizontal, left, right, cross string) string {
	var parts []string
	for _, width := range colWidths {
		parts = append(parts, strings.Repeat(horizontal, width+2))
	}
	return left + strings.Join(parts, cross) + right
}

// renderBottomBorder 渲染底部边框
func (t *Table) renderBottomBorder(colWidths []int, horizontal, left, right, cross string) string {
	var parts []string
	for _, width := range colWidths {
		parts = append(parts, strings.Repeat(horizontal, width+2))
	}
	return left + strings.Join(parts, cross) + right
}

// SetHeader 设置表头（支持链式调用）
func (t *Table) SetHeader(headers []string) *Table {
	t.headers = headers
	return t
}

// SetBorder 设置边框（支持链式调用）
func (t *Table) SetBorder(border bool) *Table {
	t.border = border
	return t
}

// SetRowLine 设置行线（支持链式调用）
func (t *Table) SetRowLine(rowLine bool) *Table {
	t.rowLine = rowLine
	return t
}

// SetAlignment 设置对齐方式（支持链式调用）
func (t *Table) SetAlignment(align Alignment) *Table {
	t.align = align
	return t
}

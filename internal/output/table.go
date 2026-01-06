package output

import (
	"os"

	"github.com/olekukonko/tablewriter"
)

// Table 表格工具
type Table struct {
	table *tablewriter.Table
}

// NewTable 创建新的表格
func NewTable(headers []string) *Table {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader(headers)
	table.SetBorder(true)
	table.SetRowLine(true)
	table.SetAutoWrapText(false)
	table.SetAlignment(tablewriter.ALIGN_LEFT)

	return &Table{table: table}
}

// AddRow 添加行
func (t *Table) AddRow(row []string) {
	t.table.Append(row)
}

// Render 渲染表格
func (t *Table) Render() {
	t.table.Render()
}

// SetHeader 设置表头
func (t *Table) SetHeader(headers []string) {
	t.table.SetHeader(headers)
}

// SetBorder 设置边框
func (t *Table) SetBorder(border bool) {
	t.table.SetBorders(tablewriter.Border{Left: border, Top: border, Right: border, Bottom: border})
}

// SetRowLine 设置行线
func (t *Table) SetRowLine(rowLine bool) {
	t.table.SetRowLine(rowLine)
}

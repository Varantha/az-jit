package list

import (
	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/lipgloss"
	azure "github.com/varantha/az-jit/internal/azure/resources"
)

var (
	activeRowStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("2"))
)

func RolesToTable(roles []azure.EligibleRoleAssignment) table.Model {
	columnSizes := map[string]int{
		"Role":       0,
		"Scope Type": 0,
		"Scope":      0,
	}

	rows := []table.Row{}

	for _, role := range roles {
		rows = append(rows, table.Row{role.RoleName, role.ResourceScopeType, role.ResourceScopeName})
		if len(role.RoleName) > columnSizes["Role"] {
			columnSizes["Role"] = len(role.RoleName)
		}
		if len(role.ResourceScopeType) > columnSizes["Scope Type"] {
			columnSizes["Scope Type"] = len(role.ResourceScopeType)
		}
		if len(role.ResourceScopeName) > columnSizes["Scope"] {
			columnSizes["Scope"] = len(role.ResourceScopeName)
		}
	}

	columns := []table.Column{}
	for name, size := range columnSizes {
		columns = append(columns, table.Column{
			Title: name,
			Width: size + 2,
		})
	}

	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithHeight(7),
	)

	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		BorderBottom(true).
		Bold(false)
	s.Selected = s.Selected.
		Foreground(lipgloss.Color("229")).
		Background(lipgloss.Color("57")).
		Bold(false)
	t.SetStyles(s)
	return t
}

func styleRow(row table.Row, st lipgloss.Style) table.Row {
	out := make(table.Row, len(row))
	for i, cell := range row {
		out[i] = st.Bold(false).Render(cell)
	}
	return out
}

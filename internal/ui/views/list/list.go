package list

import (
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	azure "github.com/varantha/az-jit/internal/azure/resources"
	"github.com/varantha/az-jit/internal/clients"
)

var baseStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.NormalBorder()).
	BorderForeground(lipgloss.Color("240"))

type Model struct {
	Client        *clients.OmniClient
	EligibleRoles []azure.EligibleRoleAssignment
	Table         table.Model
}

func InitialModel(client *clients.OmniClient) Model {
	return Model{
		Client: client,
	}
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var cmd tea.Cmd
	m.Table, cmd = m.Table.Update(msg)
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			return m, tea.Batch(
				tea.Printf("Let's activate: %s!", m.Table.SelectedRow()[0]),
			)
		}
	}

	return m, cmd
}

func (m Model) View() string {
	return baseStyle.Render(m.Table.View()) + "\n"
}

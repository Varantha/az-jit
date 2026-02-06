package ui

import (
	"context"

	tea "github.com/charmbracelet/bubbletea"
	azure "github.com/varantha/az-jit/internal/azure/resources"
	"github.com/varantha/az-jit/internal/clients"
	"github.com/varantha/az-jit/internal/ui/views/list"
)

type Screen int

const (
	screenList Screen = iota
	screenActivate
	screenLoading
)

type model struct {
	Client        *clients.OmniClient
	SelectedRoles []azure.EligibleRoleAssignment
	ListModel     list.Model

	Screen Screen
}

func InitialModel(client *clients.OmniClient) model {
	return model{
		Client:        client,
		Screen:        screenList,
		SelectedRoles: []azure.EligibleRoleAssignment{},
		ListModel:     list.InitialModel(nil),
	}
}

func (m model) Init() tea.Cmd {
	return func() tea.Msg {
		ctx := context.Background()
		eligibleRoles, err := azure.GetEligibleRoles(ctx, m.Client)
		if err != nil {
			return err
		}
		return eligibleRoles
	}
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case error:
		return m, nil
	case []azure.EligibleRoleAssignment:
		m.ListModel.Table = list.RolesToTable(msg)
		return m, nil
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		}
	}
	if m.Screen == screenList {
		var cmd tea.Cmd
		m.ListModel, cmd = m.ListModel.Update(msg)
		return m, cmd
	}
	return m, nil
}

func (m model) View() string {
	switch m.Screen {
	case screenList:
		return m.ListModel.View()
	}
	return ""
}

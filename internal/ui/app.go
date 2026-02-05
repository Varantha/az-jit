package ui

import (
	"context"

	tea "github.com/charmbracelet/bubbletea"
	azure "github.com/varantha/az-jit/internal/azure/resources"
	"github.com/varantha/az-jit/internal/clients"
)

type screen int

const (
	screenList screen = iota
	screenActivate
	screenLoading
)

type model struct {
	client        *clients.OmniClient
	eligibleRoles *[]azure.EligibleRoleAssignment
	screen        screen
}

func InitialModel(client *clients.OmniClient) model {
	return model{
		client: client,
		screen: screenList,
	}
}

func (m model) Init() tea.Cmd {
	return func() tea.Msg {
		ctx := context.Background()
		eligibleRoles, err := azure.GetEligibleRoles(ctx, m.client)
		if err != nil {
			return err
		}
		return eligibleRoles
	}
}

func (m model) Update(msg tea.Msg) (model, tea.Cmd) {
	switch msg := msg.(type) {
	case error:
		return m, nil
	case *[]azure.EligibleRoleAssignment:
		m.eligibleRoles = msg
		m.screen = screenList
		return m, nil
	default:
		return m, nil
	}
}

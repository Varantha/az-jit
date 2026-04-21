package ui

import (
	"context"

	tea "charm.land/bubbletea/v2"
	"github.com/varantha/az-jit/internal/auth"
	"github.com/varantha/az-jit/internal/clients"
)

// identityLoadedMsg delivers the authed user once auth.GetIdentity returns.
type identityLoadedMsg struct {
	User auth.AuthedUser
}

// loadErrorMsg surfaces any error from an async load (identity, roles, …)
// back into Update so the UI can switch to the error screen.
type loadErrorMsg struct {
	Err error
}

// loadIdentityCmd wraps auth.GetIdentity in a tea.Cmd so it runs off the
// event loop. On success emits identityLoadedMsg, on failure loadErrorMsg.
func loadIdentityCmd(client *clients.OmniClient) tea.Cmd {
	return func() tea.Msg {
		user, err := auth.GetIdentity(context.Background(), client.CliCredential)
		if err != nil {
			return loadErrorMsg{Err: err}
		}
		return identityLoadedMsg{User: *user}
	}
}

package ui

import (
	tea "charm.land/bubbletea/v2"
	"github.com/varantha/az-jit/internal/auth"
	"github.com/varantha/az-jit/internal/clients"
)

type screen int

const (
	// screenLoading is the default on startup: it shows while the initial
	// identity / role loads are in flight. Once identityLoadedMsg arrives
	// Update swaps to screenRoleList.
	screenLoading screen = iota
	screenError
	screenRoleList
	screenConfig
)

// contentWidth is the fixed render width used by both screens. The mockups
// were designed around ~96 columns which fits comfortably in a standard
// 80-col terminal with some overflow — crank this up to 120+ once you wire
// in tea.WindowSizeMsg-driven sizing.
const contentWidth = 96

type model struct {
	// Client is the Azure + Graph SDK façade. It is accepted here so that the
	// Init / Update hooks can fire API calls without having to reach through
	// package globals. No calls are made yet — wire them up in Init() and
	// extra message handlers in Update().
	Client *clients.OmniClient

	screen screen
	width  int
	height int

	// Shared static state used by both views.
	authedUser auth.AuthedUser

	// loadErr holds the most recent fatal load error. When non-nil the
	// screen is switched to screenError and the user can only quit.
	loadErr error

	// Role-list screen state. `roles` holds every eligible / active role
	// across all three domains; the visible slice is filtered by activeTab.
	roles     []Role
	activeTab Domain
	cursorIdx int
	selected  map[string]bool

	// Config screen state.
	activations     []Activation
	defaultReason   string
	defaultTicket   string
	defaultDuration int // index into [30m, 1h, 2h, 4h, 8h]; 3 = "4h"
	configCursor    int
}

// InitialModel builds the starting state for the TUI. It is pure: no I/O,
// no errors. All real loads (identity, roles, …) are kicked off from Init()
// as tea.Cmds; their results land in Update as typed messages.
//
// The screen starts on screenLoading; once identityLoadedMsg arrives Update
// switches to screenRoleList. Mock fixtures still seed the role/activation
// slices so the list/config screens render the moment auth succeeds — swap
// those out as the real loaders get wired up.
func InitialModel(client *clients.OmniClient) model {
	reason, ticket, _ := mockDefaults()

	return model{
		Client:          client,
		screen:          screenLoading,
		width:           contentWidth,
		height:          30,
		roles:           mockAllRoles(),
		activeTab:       DomainEntra,
		cursorIdx:       0,
		selected:        selectedIDs(),
		activations:     mockActivations(),
		defaultReason:   reason,
		defaultTicket:   ticket,
		defaultDuration: 3, // "4h"
		configCursor:    0,
	}
}

// Init kicks off the initial identity load. Extend with tea.Batch(...) as
// the role loaders are added:
//
//	return tea.Batch(
//	    loadIdentityCmd(m.Client),
//	    loadEntraRolesCmd(m.Client),
//	    loadAzureRolesCmd(m.Client),
//	    loadGroupRolesCmd(m.Client),
//	)
//
// Each loader emits its own *LoadedMsg on success and a shared
// loadErrorMsg on failure (see commands.go).
func (m model) Init() tea.Cmd {
	return loadIdentityCmd(m.Client)
}

// Update is intentionally minimal: enough to demo nav between the two screens
// and exercise the cursor. All the real interactions (filter mode, edit mode,
// default adjusting, tab switching, selection toggle, activation POST, etc.)
// are left as TODOs for you to wire up.
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

	case identityLoadedMsg:
		m.authedUser = msg.User
		m.screen = screenRoleList

	case loadErrorMsg:
		m.loadErr = msg.Err
		m.screen = screenError

	case tea.KeyPressMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit

		case "enter":
			if m.screen == screenRoleList {
				m.screen = screenConfig
			}

		case "esc":
			if m.screen == screenConfig {
				m.screen = screenRoleList
			}

		case "up", "k":
			if m.screen == screenRoleList && m.cursorIdx > 0 {
				m.cursorIdx--
			}
			if m.screen == screenConfig && m.configCursor > 0 {
				m.configCursor--
			}

		case "down", "j":
			if m.screen == screenRoleList && m.cursorIdx < len(m.visibleRoles())-1 {
				m.cursorIdx++
			}
			if m.screen == screenConfig && m.configCursor < len(m.activations)-1 {
				m.configCursor++
			}

		case "left", "h":
			if m.screen == screenRoleList {
				m.activeTab = prevTab(m.activeTab)
				m.cursorIdx = 0
			}

		case "right", "l", "tab":
			if m.screen == screenRoleList {
				m.activeTab = nextTab(m.activeTab)
				m.cursorIdx = 0
			}

		case "space":
			if m.screen == screenRoleList {
				visible := m.visibleRoles()
				if m.cursorIdx < len(visible) {
					r := visible[m.cursorIdx]
					if !r.Active {
						if m.selected == nil {
							m.selected = map[string]bool{}
						}
						m.selected[r.ID] = !m.selected[r.ID]
					}
				}
			}
		}
	}
	return m, nil
}

// tabOrder is the left-to-right order of the role-list tabs. Left/right
// arrows cycle through this slice.
var tabOrder = []Domain{DomainEntra, DomainAzure, DomainGroups}

func nextTab(d Domain) Domain {
	for i, t := range tabOrder {
		if t == d {
			return tabOrder[(i+1)%len(tabOrder)]
		}
	}
	return tabOrder[0]
}

func prevTab(d Domain) Domain {
	for i, t := range tabOrder {
		if t == d {
			return tabOrder[(i-1+len(tabOrder))%len(tabOrder)]
		}
	}
	return tabOrder[0]
}

// visibleRoles returns the roles filtered to the currently active tab.
func (m model) visibleRoles() []Role {
	out := make([]Role, 0, len(m.roles))
	for _, r := range m.roles {
		if r.Domain == m.activeTab {
			out = append(out, r)
		}
	}
	return out
}

// tabCounts returns the number of roles per domain, for the badge pills.
func (m model) tabCounts() map[Domain]int {
	counts := map[Domain]int{}
	for _, r := range m.roles {
		counts[r.Domain]++
	}
	return counts
}

func (m model) View() tea.View {
	var content string
	switch m.screen {
	case screenLoading:
		content = m.renderLoading()
	case screenError:
		content = m.renderError()
	case screenRoleList:
		content = m.renderRoleList()
	case screenConfig:
		content = m.renderConfig()
	}
	v := tea.NewView(content)
	v.AltScreen = true
	return v
}

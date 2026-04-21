package ui

import (
	"strings"

	"charm.land/lipgloss/v2"
)

// renderLoading is the pre-auth splash shown while loadIdentityCmd is in
// flight. Kept intentionally plain — no header, no tabs — because authedUser
// isn't populated yet.
func (m model) renderLoading() string {
	msg := sScreenTitle.Render("az-jit") + "\n\n" +
		sSec.Render("Loading identity…")

	box := lipgloss.NewStyle().
		Width(contentWidth).
		Height(10).
		Align(lipgloss.Center, lipgloss.Center).
		Render(msg)

	var sb strings.Builder
	sb.WriteString(box)
	sb.WriteString("\n")
	sb.WriteString(divider(contentWidth))
	sb.WriteString("\n")
	sb.WriteString(" " + keybinds([]keybind{{"q", "quit"}}))
	return sb.String()
}

// renderError is shown whenever a loadErrorMsg lands. It keeps the wrapper
// layout (divider + footer) consistent with the other screens so terminal
// state stays tidy on exit.
func (m model) renderError() string {
	title := sStatusWarn.Bold(true).Render("Something went wrong")
	detail := sText.Render(m.loadErr.Error())
	hint := sDimmer.Render("Is the Azure CLI logged in? Try `az login` then retry.")

	body := title + "\n\n" + detail + "\n\n" + hint

	box := lipgloss.NewStyle().
		Width(contentWidth).
		Padding(2, 4).
		Render(body)

	var sb strings.Builder
	sb.WriteString(box)
	sb.WriteString("\n")
	sb.WriteString(divider(contentWidth))
	sb.WriteString("\n")
	sb.WriteString(" " + keybinds([]keybind{{"q", "quit"}}))
	return sb.String()
}

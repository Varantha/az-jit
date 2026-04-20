package ui

import (
	"fmt"
	"strings"
	"time"

	"charm.land/lipgloss/v2"
)

// Column widths for the config table. They sum (with the gutter, cursor and
// badge) to contentWidth:  2 + 2 + 8 + 28 + 24 + 14 + 14 + 4 = 96.
const (
	cfGutter   = 2
	cfCursor   = 2
	cfBadge    = 8 // "[ENTRA] " / "[AZURE] " — padded
	cfRoleName = 28
	cfReason   = 24
	cfTicket   = 14
	cfDuration = 14
	cfStatus   = 4
)

// Duration scale used by the defaults slider. Left/right arrow would step
// through this; we render index defaultDuration.
var durationScale = []time.Duration{
	30 * time.Minute,
	1 * time.Hour,
	2 * time.Hour,
	4 * time.Hour,
	8 * time.Hour,
}

func (m model) renderConfig() string {
	var sb strings.Builder
	sb.WriteString(m.renderUserHeader())
	sb.WriteString("\n")
	sb.WriteString(divider(contentWidth))
	sb.WriteString("\n")
	sb.WriteString(m.renderConfigTitle())
	sb.WriteString("\n\n")
	sb.WriteString(m.renderDefaultsHeader())
	sb.WriteString("\n")
	sb.WriteString(m.renderDefaultsRow())
	sb.WriteString("\n")
	sb.WriteString(divider(contentWidth))
	sb.WriteString("\n")
	sb.WriteString(m.renderConfigColumnHeaders())
	sb.WriteString("\n")
	sb.WriteString(divider(contentWidth))
	sb.WriteString("\n")
	for i, a := range m.activations {
		sb.WriteString(m.renderConfigRow(a, i == m.configCursor))
		sb.WriteString("\n")
	}
	sb.WriteString(divider(contentWidth))
	sb.WriteString("\n")
	sb.WriteString(m.renderConfigFooter())
	return sb.String()
}

// "CONFIGURE ACTIVATION  [5 roles]" on the left, hint on the right.
func (m model) renderConfigTitle() string {
	title := sScreenTitle.Render("CONFIGURE ACTIVATION")
	countBadge := sBadgePurple.Render(fmt.Sprintf(" %d roles ", len(m.activations)))
	left := title + "  " + countBadge
	right := sDim.Render("fill once, override per-row")
	return layoutLR(contentWidth, left, right)
}

func (m model) renderDefaultsHeader() string {
	return sDefaultsHeader.Render("◆ DEFAULTS")
}

// One-line defaults row: reason stretches across the ROLE+REASON columns
// (since there's no per-row role in the defaults), ticket in its column,
// duration as value + slider + ←→ hint in its column.
func (m model) renderDefaultsRow() string {
	// Reason stretches cfRoleName + cfReason.
	reasonStretch := cfRoleName + cfReason
	reason := padRight(sText.Render(ellipsize(m.defaultReason, reasonStretch)), reasonStretch)

	ticket := padRight(sText.Render(ellipsize(m.defaultTicket, cfTicket)), cfTicket)

	// Duration widget: "4h ├──●──┤ ←→"
	dur := sText.Render(formatDuration(durationScale[m.defaultDuration]))
	widget := dur + " " + renderSlider(len(durationScale), m.defaultDuration) + " " + sTer.Render("←→")
	duration := padRight(widget, cfDuration)

	return pad(cfGutter+cfCursor+cfBadge) + reason + ticket + duration
}

func (m model) renderConfigColumnHeaders() string {
	return pad(cfGutter+cfCursor+cfBadge) +
		sColHeader.Render(padRight("ROLE", cfRoleName)) +
		sColHeader.Render(padRight("REASON", cfReason)) +
		sColHeader.Render(padRight("TICKET", cfTicket)) +
		sColHeader.Render(padRight("DUR", cfDuration))
}

// renderConfigRow renders one selected role with its config values.
// It shows: inherited values in dim, overrides in bright + pencil prefix,
// clamped durations in amber, and a status ✓/⚠ at the end.
func (m model) renderConfigRow(a Activation, isCursor bool) string {
	r := a.Role

	// Cursor glyph.
	cursor := " "
	if isCursor {
		cursor = "▸"
	}

	// Domain badge.
	var badge string
	switch r.Domain {
	case DomainEntra:
		badge = sBadgeEntra.Render("ENTRA")
	case DomainAzure:
		badge = sBadgeAzure.Render("AZURE")
	case DomainGroups:
		badge = sBadgeAzure.Render("GROUP") // pick a teal ramp if you add one
	}

	// Role name cell.
	nameCell := padRight(r.Name, cfRoleName)

	// Reason cell — override vs inherited vs missing.
	reasonCell := renderReasonCell(a, m.defaultReason)

	// Ticket cell — override vs inherited vs missing vs n/a.
	ticketCell := renderTicketCell(a, m.defaultTicket)

	// Duration cell — shows clamped value if role max is below default.
	durCell := renderDurationCell(a, durationScale[m.defaultDuration])

	// Status indicator.
	status := padLeft(configRowStatus(a, m.defaultReason, m.defaultTicket), cfStatus)

	// Compose the row.
	var row string
	if isCursor {
		row = pad(cfGutter) +
			padRight(sKeyBind.Render(cursor), cfCursor) +
			padRight(badge, cfBadge) +
			sText.Bold(true).Render(nameCell) +
			reasonCell +
			ticketCell +
			durCell +
			status
		return sRowCursor.Width(contentWidth).Render(row)
	}

	return pad(cfGutter) +
		pad(cfCursor) +
		padRight(badge, cfBadge) +
		sText.Render(nameCell) +
		reasonCell +
		ticketCell +
		durCell +
		status
}

// renderReasonCell formats the reason column:
//   - override (act.Reason non-empty and != default): "✎ <val>" in primary text
//   - inherited: default value in dim
func renderReasonCell(a Activation, defaultReason string) string {
	if a.Reason != "" && strings.TrimSpace(a.Reason) != defaultReason {
		content := sOverride.Render("✎ ") + sText.Render(ellipsize(a.Reason, cfReason-3))
		return padRight(content, cfReason)
	}
	return sSec.Render(padRight(ellipsize(defaultReason, cfReason), cfReason))
}

// renderTicketCell handles the four cases (inherited, overridden, missing, n/a).
func renderTicketCell(a Activation, defaultTicket string) string {
	r := a.Role
	if !r.RequireTicket {
		return sDimmer.Render(padRight("n/a", cfTicket))
	}
	// In this mock, Activation.Ticket = " " means "user cleared ticket" → missing required.
	if strings.TrimSpace(a.Ticket) == "" && a.Ticket != "" {
		return sStatusWarn.Render(padRight("⚠ needed", cfTicket))
	}
	if a.Ticket != "" && a.Ticket != defaultTicket {
		return sText.Render(padRight("✎ "+a.Ticket, cfTicket))
	}
	return sSec.Render(padRight(defaultTicket, cfTicket))
}

// renderDurationCell shows the applied duration. If the role's MaxDuration
// is shorter than the default, the value is clamped and shown in amber.
func renderDurationCell(a Activation, defaultDuration time.Duration) string {
	r := a.Role
	applied := defaultDuration
	if a.Duration > 0 {
		applied = a.Duration
	}
	clamped := false
	if r.MaxDuration > 0 && applied > r.MaxDuration {
		applied = r.MaxDuration
		clamped = true
	}

	text := formatDuration(applied)
	if clamped {
		return lipgloss.NewStyle().Foreground(cAmberT).Render(padRight(text+" max", cfDuration))
	}
	return sSec.Render(padRight(text, cfDuration))
}

// configRowStatus returns a styled "✓" or "⚠".
func configRowStatus(a Activation, defaultReason, defaultTicket string) string {
	r := a.Role

	reasonOk := true
	if r.RequireReason {
		reasonOk = a.Reason != "" || defaultReason != ""
	}

	ticketOk := true
	if r.RequireTicket {
		switch {
		case strings.TrimSpace(a.Ticket) == "" && a.Ticket != "":
			ticketOk = false // explicitly cleared
		case a.Ticket == "" && defaultTicket == "":
			ticketOk = false
		}
	}

	if reasonOk && ticketOk {
		return sStatusOK.Render("✓")
	}
	return sStatusWarn.Render("⚠")
}

// Footer: readiness badge on the left, keybinds trailing.
func (m model) renderConfigFooter() string {
	ready, total := readinessCounts(m.activations, m.defaultReason, m.defaultTicket)

	var badge string
	if ready < total {
		badge = sBadgeAmber.Render(fmt.Sprintf(" ⚠ %d / %d ready ", ready, total))
	} else {
		badge = sBadgePurple.Bold(true).Render(fmt.Sprintf(" %d / %d ready ", ready, total))
	}

	activateStyle := sKeyBind
	activateDesc := sKeyBindDesc
	if ready < total {
		activateStyle = sKeyBindDim
		activateDesc = sKeyBindDim
	}

	binds := keybindsCustom([]customKeybind{
		{sKeyBind.Render("↑↓"), sKeyBindDesc.Render("move")},
		{sKeyBind.Render("tab"), sKeyBindDesc.Render("defaults")},
		{sKeyBind.Render("e"), sKeyBindDesc.Render("edit")},
		{sKeyBind.Render("r"), sKeyBindDesc.Render("reset")},
		{sKeyBind.Render("x"), sKeyBindDesc.Render("remove")},
		{activateStyle.Render("↵"), activateDesc.Render("activate")},
		{sKeyBind.Render("esc"), sKeyBindDesc.Render("back")},
		{sKeyBind.Render("q"), sKeyBindDesc.Render("quit")},
	})

	return " " + badge + "  " + binds
}

func readinessCounts(acts []Activation, defaultReason, defaultTicket string) (ready, total int) {
	total = len(acts)
	for _, a := range acts {
		s := configRowStatus(a, defaultReason, defaultTicket)
		// Plain-text check: the rendered string for OK starts with "✓".
		if strings.Contains(s, "✓") {
			ready++
		}
	}
	return
}

// renderSlider draws a short horizontal slider with the dot at `idx` of `n`
// stops. The rendered width (including brackets) is 8 characters so the full
// defaults duration widget ("4h ├──●──┤ ←→") fits in cfDuration=14 columns.
func renderSlider(n, idx int) string {
	const inner = 6 // track width between the brackets
	if n <= 1 {
		return sDim.Render("├" + strings.Repeat("─", inner) + "┤")
	}
	pos := idx * (inner - 1) / (n - 1)
	var out strings.Builder
	out.WriteString(sDim.Render("├"))
	for i := 0; i < inner; i++ {
		if i == pos {
			out.WriteString(sOverride.Render("●"))
		} else {
			out.WriteString(sDim.Render("─"))
		}
	}
	out.WriteString(sDim.Render("┤"))
	return out.String()
}

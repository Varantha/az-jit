package ui

import (
	"fmt"
	"strings"
	"time"

	"charm.land/lipgloss/v2"
)

// Fixed column widths for the role-list table. Name / scope-type / scope
// widths are derived per-render from the visible rows (see listColumnWidths)
// so long scope types can't bleed into the next column.
const (
	lrGutter = 2
	lrCursor = 2 // "▸ " or "  "
	lrCheck  = 2 // "■ " / "□ " / "● "
	lrTime   = 14 // right-aligned time-left, or blank

	lrColGap    = 2  // minimum spacing between dynamic columns
	lrNameFloor = 16 // don't let ROLE get narrower than this
)

// listCols holds the render widths for the three dynamic columns of the
// role-list table. They sum to contentWidth - (gutter + cursor + check + time).
type listCols struct {
	name, typ, scope int
}

// listColumnWidths sizes scope-type and scope to their longest value (so
// nothing gets truncated unnecessarily) and hands ROLE the remaining space.
// If content is pathologically wide, ROLE is floored and its cell gets
// ellipsized at render time.
func listColumnWidths(visible []Role) listCols {
	const fixed = lrGutter + lrCursor + lrCheck + lrTime
	flex := contentWidth - fixed

	typ := lipgloss.Width("SCOPE TYPE")
	scope := lipgloss.Width("SCOPE")
	for _, r := range visible {
		if w := lipgloss.Width(r.ScopeType); w > typ {
			typ = w
		}
		if w := lipgloss.Width(r.Scope); w > scope {
			scope = w
		}
	}
	typ += lrColGap
	scope += lrColGap

	name := flex - typ - scope
	if name < lrNameFloor {
		name = lrNameFloor
	}
	return listCols{name: name, typ: typ, scope: scope}
}

func (m model) renderRoleList() string {
	visible := m.visibleRoles()
	cols := listColumnWidths(visible)

	var sb strings.Builder
	sb.WriteString(m.renderUserHeader())
	sb.WriteString("\n")
	sb.WriteString(divider(contentWidth))
	sb.WriteString("\n")
	sb.WriteString(m.renderTabs())
	sb.WriteString("\n")
	sb.WriteString(divider(contentWidth))
	sb.WriteString("\n")
	sb.WriteString(renderListColumnHeaders(cols))
	sb.WriteString("\n")
	sb.WriteString(divider(contentWidth))
	sb.WriteString("\n")
	for i, r := range visible {
		sb.WriteString(m.renderListRow(r, i == m.cursorIdx, cols))
		sb.WriteString("\n")
	}
	sb.WriteString(divider(contentWidth))
	sb.WriteString("\n")
	sb.WriteString(m.renderListFooter())
	return sb.String()
}

// User / tenant header, single line.
func (m model) renderUserHeader() string {
	connectedDot := sConnected.Render("●")

	sep := sDim.Render("·")
	name := sText.Bold(true).Render(m.authedUser.Name)
	email := sTer.Render(m.authedUser.Email)
	tenant := sTenantID.Render("tenant " + m.authedUser.TenantID)
	left := name + "  " + sep + " " + email + "  " + sep + " " + tenant
	return layoutLR(contentWidth, left, connectedDot)
}

// Tab bar with count pills and a right-aligned filter hint.
func (m model) renderTabs() string {
	counts := m.tabCounts()

	tab := func(label string, d Domain) string {
		count := fmt.Sprintf(" %d ", counts[d])
		if m.activeTab == d {
			return sTabActive.Render(label) + sBadgePurple.Render(count)
		}
		return sTabInactive.Render(label) + sBadgeNeutral.Render(count)
	}

	left := tab("Entra Roles", DomainEntra) +
		"  " + tab("Azure Roles", DomainAzure) +
		"  " + tab("Groups", DomainGroups)
	right := sDimmer.Render("/ to filter")
	return layoutLR(contentWidth, left, right)
}

func renderListColumnHeaders(cols listCols) string {
	// Same column layout as rows, offset by gutter + cursor + checkbox.
	return pad(lrGutter) + pad(lrCursor) + pad(lrCheck) +
		sColHeader.Render(padRight("ROLE", cols.name)) +
		sColHeader.Render(padRight("SCOPE TYPE", cols.typ)) +
		sColHeader.Render(padRight("SCOPE", cols.scope))
}

// renderListRow renders one eligible (or currently active) role row.
func (m model) renderListRow(r Role, isCursor bool, cols listCols) string {
	// Cursor and selection glyphs.
	cursor := " "
	if isCursor {
		cursor = "▸"
	}

	var check string
	switch {
	case r.Active:
		check = "●"
	case m.selected[r.ID]:
		check = "■"
	default:
		check = "□"
	}

	// Cells — ellipsize to the dynamic column widths so long values don't
	// bleed across the column boundary.
	nameCell := padRight(ellipsize(r.Name, cols.name), cols.name)
	typeCell := padRight(ellipsize(r.ScopeType, cols.typ), cols.typ)
	scopeCell := padRight(ellipsize(r.Scope, cols.scope), cols.scope)

	timeCell := ""
	if r.Active {
		timeCell = formatDuration(r.TimeLeft) + " left"
	}
	timeCell = padLeft(timeCell, lrTime)

	// Compose with per-cell colours.
	var row string
	switch {
	case r.Active:
		row = pad(lrGutter) +
			padRight(sActiveDot.Render(cursor), lrCursor) +
			padRight(sActiveDot.Render(check), lrCheck) +
			sActiveRole.Render(nameCell) +
			sActiveScope.Render(typeCell) +
			sActiveScope.Render(scopeCell) +
			sActiveLeft.Render(timeCell)
		return sRowActive.Width(contentWidth).Render(row)

	case isCursor:
		row = pad(lrGutter) +
			padRight(sKeyBind.Render(cursor), lrCursor) +
			padRight(sKeyBind.Render(check), lrCheck) +
			sText.Bold(true).Render(nameCell) +
			sText.Render(typeCell) +
			sText.Render(scopeCell) +
			sSec.Render(timeCell)
		return sRowCursor.Width(contentWidth).Render(row)

	default:
		var checkStyle lipgloss.Style
		if m.selected[r.ID] {
			checkStyle = sOverride // purple filled
		} else {
			checkStyle = sDim
		}
		return pad(lrGutter) +
			pad(lrCursor) + // no cursor glyph
			padRight(checkStyle.Render(check), lrCheck) +
			sText.Render(nameCell) +
			sSec.Render(typeCell) +
			sSec.Render(scopeCell) +
			sSec.Render(timeCell)
	}
}

// Footer: selection count badge on the left, keybinds trailing right.
func (m model) renderListFooter() string {
	selectedCount := 0
	for _, v := range m.selected {
		if v {
			selectedCount++
		}
	}

	badge := ""
	if selectedCount > 0 {
		badge = sBadgePurple.Bold(true).Render(fmt.Sprintf(" %d selected ", selectedCount))
	}

	binds := keybinds([]keybind{
		{"↑↓", "move"},
		{"←→", "tab"},
		{"space", "select"},
		{"↵", "activate"},
		{"/", "filter"},
		{"?", "help"},
		{"q", "quit"},
	})

	if badge == "" {
		return " " + binds
	}
	return " " + badge + "  " + binds
}

// formatDuration renders short human durations like "1h 42m" / "30m".
func formatDuration(d time.Duration) string {
	if d <= 0 {
		return ""
	}
	h := int(d / time.Hour)
	m := int((d % time.Hour) / time.Minute)
	switch {
	case h > 0 && m > 0:
		return fmt.Sprintf("%dh %dm", h, m)
	case h > 0:
		return fmt.Sprintf("%dh", h)
	default:
		return fmt.Sprintf("%dm", m)
	}
}

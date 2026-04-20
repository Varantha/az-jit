package ui

import (
	"strings"

	"charm.land/lipgloss/v2"
)

// pad returns a string of n spaces. Used for fixed-width gutters.
func pad(n int) string {
	if n <= 0 {
		return ""
	}
	return strings.Repeat(" ", n)
}

// padRight pads s on the right with spaces so its visible width is w.
// If s is already wider, returns s unchanged.
func padRight(s string, w int) string {
	actual := lipgloss.Width(s)
	if actual >= w {
		return s
	}
	return s + strings.Repeat(" ", w-actual)
}

// padLeft pads s on the left with spaces so its visible width is w.
func padLeft(s string, w int) string {
	actual := lipgloss.Width(s)
	if actual >= w {
		return s
	}
	return strings.Repeat(" ", w-actual) + s
}

// ellipsize truncates s to visible width w, appending "…" if truncated.
// It counts runes rather than bytes so unicode content doesn't get mangled.
func ellipsize(s string, w int) string {
	if lipgloss.Width(s) <= w {
		return s
	}
	if w <= 1 {
		return "…"
	}
	runes := []rune(s)
	for len(runes) > 0 && lipgloss.Width(string(runes)+"…") > w {
		runes = runes[:len(runes)-1]
	}
	return string(runes) + "…"
}

// layoutLR composes a single-line row of fixed width, with `left` flush left,
// `right` flush right, and spaces filling the middle.
func layoutLR(width int, left, right string) string {
	lw := lipgloss.Width(left)
	rw := lipgloss.Width(right)
	gap := width - lw - rw
	if gap < 1 {
		gap = 1
	}
	return left + strings.Repeat(" ", gap) + right
}

// keybind is a single "key + description" footer item.
type keybind struct {
	Key  string
	Desc string
}

func keybinds(kbs []keybind) string {
	parts := make([]string, 0, len(kbs))
	for _, kb := range kbs {
		parts = append(parts, sKeyBind.Render(kb.Key)+" "+sKeyBindDesc.Render(kb.Desc))
	}
	return strings.Join(parts, "  ")
}

// customKeybind allows pre-rendered key/desc strings (so individual bindings
// can be dimmed or recoloured, e.g. ↵ activate when not ready).
type customKeybind struct {
	Key  string
	Desc string
}

func keybindsCustom(kbs []customKeybind) string {
	parts := make([]string, 0, len(kbs))
	for _, kb := range kbs {
		parts = append(parts, kb.Key+" "+kb.Desc)
	}
	return strings.Join(parts, "  ")
}

package ui

import "charm.land/lipgloss/v2"

// Color tokens — match the SVG mockup palette.
var (
	cPurple      = lipgloss.Color("#a78bfa")
	cPurpleLight = lipgloss.Color("#c4b5fd")
	cPurpleBg    = lipgloss.Color("#3b2d5c")
	cCursorBg    = lipgloss.Color("#2d1f5f")
	cActiveTabBg = lipgloss.Color("#1a1532")

	cAzureLight = lipgloss.Color("#93c5fd")
	cAzureBg    = lipgloss.Color("#1e3a5f")

	cAmber   = lipgloss.Color("#d97706")
	cAmberT  = lipgloss.Color("#c48a30")
	cAmberBg = lipgloss.Color("#141008")
	cAmberS  = lipgloss.Color("#7a6030")
	cAmberR  = lipgloss.Color("#2a1a08")

	cGreen = lipgloss.Color("#86efac")

	cText    = lipgloss.Color("#dde0e8")
	cSec     = lipgloss.Color("#7a8290")
	cTer     = lipgloss.Color("#5d6373")
	cDim     = lipgloss.Color("#4a5065")
	cDimmer  = lipgloss.Color("#3d4357")
	cDivider = lipgloss.Color("#262a38")
	cTabBg   = lipgloss.Color("#1e2230")
)

// Text styles.
var (
	sText      = lipgloss.NewStyle().Foreground(cText)
	sSec       = lipgloss.NewStyle().Foreground(cSec)
	sTer       = lipgloss.NewStyle().Foreground(cTer)
	sDim       = lipgloss.NewStyle().Foreground(cDim)
	sDimmer    = lipgloss.NewStyle().Foreground(cDimmer)
	sTenantID  = lipgloss.NewStyle().Foreground(cTer)
	sConnected = lipgloss.NewStyle().Foreground(cGreen)

	sScreenTitle = lipgloss.NewStyle().Foreground(cPurpleLight).Bold(true)
	sColHeader   = lipgloss.NewStyle().Foreground(cTer)

	sKeyBind     = lipgloss.NewStyle().Foreground(cPurple)
	sKeyBindDesc = lipgloss.NewStyle().Foreground(cSec)
	sKeyBindDim  = lipgloss.NewStyle().Foreground(cDimmer)

	sStatusOK   = lipgloss.NewStyle().Foreground(cGreen)
	sStatusWarn = lipgloss.NewStyle().Foreground(cAmber)

	sOverride       = lipgloss.NewStyle().Foreground(cPurple)
	sDefaultsHeader = lipgloss.NewStyle().Foreground(cPurple).Bold(true)

	// Tabs.
	sTabActive   = lipgloss.NewStyle().Foreground(cPurpleLight).Bold(true).Background(cActiveTabBg).Padding(0, 2)
	sTabInactive = lipgloss.NewStyle().Foreground(cSec).Padding(0, 2)

	// Pill-shaped count / status badges.
	sBadgePurple  = lipgloss.NewStyle().Foreground(cPurpleLight).Background(cPurpleBg).Padding(0, 1)
	sBadgeNeutral = lipgloss.NewStyle().Foreground(cSec).Background(cTabBg).Padding(0, 1)
	sBadgeAmber   = lipgloss.NewStyle().Foreground(cAmber).Background(cAmberR).Bold(true).Padding(0, 1)

	// Domain badges.
	sBadgeEntra = lipgloss.NewStyle().Foreground(cPurpleLight).Background(cPurpleBg).Bold(true).Padding(0, 1)
	sBadgeAzure = lipgloss.NewStyle().Foreground(cAzureLight).Background(cAzureBg).Bold(true).Padding(0, 1)

	// Row backgrounds — apply Foreground explicitly inside because background
	// colours don't compose through nested styles in lipgloss.
	sRowCursor = lipgloss.NewStyle().Background(cCursorBg).Foreground(cText)
	sRowActive = lipgloss.NewStyle().Background(cAmberBg).Foreground(cAmberT)

	// Active-role inline styles (inside sRowActive).
	sActiveDot   = lipgloss.NewStyle().Foreground(cAmber).Background(cAmberBg)
	sActiveRole  = lipgloss.NewStyle().Foreground(cAmberT).Background(cAmberBg)
	sActiveScope = lipgloss.NewStyle().Foreground(cAmberS).Background(cAmberBg)
	sActiveLeft  = lipgloss.NewStyle().Foreground(cAmberS).Background(cAmberBg)
)

// Divider draws a full-width horizontal rule in the muted divider colour.
// Lipgloss handles terminal width; we pass it in explicitly so it matches
// the rendered content block.
func divider(width int) string {
	return lipgloss.NewStyle().Foreground(cDivider).Render(repeat("─", width))
}

func repeat(s string, n int) string {
	if n <= 0 {
		return ""
	}
	out := make([]byte, 0, len(s)*n)
	for i := 0; i < n; i++ {
		out = append(out, s...)
	}
	return string(out)
}

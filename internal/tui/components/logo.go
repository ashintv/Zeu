package components

import (
	"github.com/charmbracelet/lipgloss"
)

var (
	// Claude-inspired palette: warm coral + cool slate
	coral     = lipgloss.NewStyle().Foreground(lipgloss.Color("#CC785C")) // warm terracotta
	sand      = lipgloss.NewStyle().Foreground(lipgloss.Color("#D4A27A")) // muted amber
	slate     = lipgloss.NewStyle().Foreground(lipgloss.Color("#8B9EB7")) // cool slate blue
	brandName = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#CC785C")).
			Bold(true)
	tagline = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#8B9EB7")).
		Italic(true)
)

func Logo() string {
	// Minimal geometric mark — a stylized "Z" built from blocks
	// Coral top bar, sand diagonal, slate bottom bar
	mark := lipgloss.JoinVertical(
		lipgloss.Left,
		coral.Render("  ████████"),
		sand.Render("     ███  "),
		sand.Render("    ███   "),
		slate.Render("   ███    "),
		slate.Render("  ████████"),
	)

	

	logo := lipgloss.JoinHorizontal(lipgloss.Top, mark)


	return lipgloss.JoinVertical(lipgloss.Left,
		"",
		logo,
		"",
	)
}
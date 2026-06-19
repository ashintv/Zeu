package components

import (
	"charm.land/lipgloss/v2"
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

func Logo(frame int) string {
	// Colors to cycle: coral, sand, sand, slate, slate
	colors := []string{
		"#CC785C", // coral
		"#D4A27A", // sand
		"#D4A27A", // sand
		"#8B9EB7", // slate
		"#8B9EB7", // slate
	}

	getColor := func(lineIdx int) lipgloss.Style {
		idx := (lineIdx - frame) % len(colors)
		if idx < 0 {
			idx += len(colors)
		}
		c := colors[idx]
		return lipgloss.NewStyle().Foreground(lipgloss.Color(c))
	}

	// Minimal geometric mark — a stylized "Z" built from blocks
	mark := lipgloss.JoinVertical(
		lipgloss.Left,
		getColor(0).Render("  ████████"),
		getColor(1).Render("     ███  "),
		getColor(2).Render("    ███   "),
		getColor(3).Render("   ███    "),
		getColor(4).Render("  ████████"),
	)

	logo := lipgloss.JoinHorizontal(lipgloss.Top, mark)

	return logo
}

func LogoSmall() string {
	return brandName.Render(" Zeu")
}

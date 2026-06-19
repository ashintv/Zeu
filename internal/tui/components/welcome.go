package components

import (
	

	"charm.land/lipgloss/v2"
)

var (
	textS = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#8B9EB7")).
		Italic(true)

	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#CC785C"))

	textStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#C4CCDA"))

	linkStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#D4A27A"))

	separatorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#3D4A57"))

	subtleStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#6B7F96"))
)

func Welcome(width int, logoFrame int) string {
	logo := Logo(logoFrame)

	content := lipgloss.JoinVertical(
		lipgloss.Left,
		brandName.Render("Zeu ai"),
		tagline.Render("orchestrate everything"),

		"",
		tagline.Render("⭐ Star the project:"),
		linkStyle.Render("github.com/ashintv/zeu"),
	)

	body := lipgloss.JoinHorizontal(
		lipgloss.Center,
		logo,
		"        ", // 8 spaces padding
		content,
	)

	hints := subtleStyle.Render("  ctrl+c quit  •  enter send")

	// line := separatorStyle.Render(
	// 	strings.Repeat("─", max(0, width)),
	// )

	return body + "\n\n" + hints + "\n\n"  
}

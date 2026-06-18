package components

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

var (
	textS = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#8B9EB7")).
		Italic(true)

	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#00D7FF"))

	textStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#CBD5E1"))

	linkStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#10B981"))

	separatorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#334155"))

	subtleStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#94A3B8"))
)

func Welcome(width int) string {
	logo := Logo()

	// Right side: name + tagline
	// label := lipgloss.JoinVertical(
	//

	// )

	content := lipgloss.NewStyle().
		MarginTop(1).
		Render(
			lipgloss.JoinVertical(
				lipgloss.Left,
				brandName.Render("Zeu ai"),
				tagline.Render("orchestrate everything"),

				"",
				tagline.Render("⭐ Star the project:"),
				linkStyle.Render("github.com/ashintv/zeu"),

			),
		)

	body := lipgloss.JoinHorizontal(
		lipgloss.Top,
		logo,
		"    ",
		"    ",
		content,
	)

	line := separatorStyle.Render(
		strings.Repeat("─", max(0, width)),
	)

	return body + "\n\n" + line
}

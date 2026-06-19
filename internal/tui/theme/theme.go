package theme

import (
	"time"

	"charm.land/lipgloss/v2"
)

const (
	colorCoral = "#CC785C"
	colorSand  = "#D4A27A"
	colorSlate = "#8B9EB7"

	colorCoralLight = "#E0967A"
	colorCoralDark  = "#A35A3E"
	colorSandDim    = "#B8956E"
	colorSlateDim   = "#6B7F96"
	colorSlateDark  = "#4E6175"
	colorGutter     = "#3D4A57"
	colorCream      = "#F0E6DC"
	colorAsh        = "#C4CCDA"
)

type Theme struct {
	User           lipgloss.Style
	Agent          lipgloss.Style
	ToolCall       lipgloss.Style
	ToolResult     lipgloss.Style
	Error          lipgloss.Style
	System         lipgloss.Style
	TimestampStyle lipgloss.Style
	Border         lipgloss.Style
	Prompt         lipgloss.Style
	Input          lipgloss.Style
}

var Default = LogoTheme()

func LogoTheme() Theme {
	return Theme{
		User: lipgloss.NewStyle().
			Foreground(lipgloss.Color(colorCoral)),
		Agent: lipgloss.NewStyle().
			Foreground(lipgloss.Color(colorSlate)),
		ToolCall: lipgloss.NewStyle().
			Foreground(lipgloss.Color(colorSand)),
		ToolResult: lipgloss.NewStyle().
			Foreground(lipgloss.Color(colorSandDim)),
		Error: lipgloss.NewStyle().
			Foreground(lipgloss.Color(colorCoralDark)).
			Bold(true),
		System: lipgloss.NewStyle().
			Foreground(lipgloss.Color(colorSlateDim)),

		TimestampStyle: lipgloss.NewStyle().Foreground(lipgloss.Color(colorSlateDark)),
		Border:         lipgloss.NewStyle().Foreground(lipgloss.Color(colorGutter)),

		Prompt: lipgloss.NewStyle().Foreground(lipgloss.Color(colorCream)).Bold(true),
		Input:  lipgloss.NewStyle().Foreground(lipgloss.Color(colorAsh)),
	}
}

func (t Theme) RoleStyle(role string) lipgloss.Style {
	switch role {
	case "user":
		return t.User
	case "agent", "assistant", "ai", "model":
		return t.Agent
	case "tool_call", "tool_use", "function_call":
		return t.ToolCall
	case "tool_result", "tool_response", "function_response":
		return t.ToolResult
	case "error":
		return t.Error
	default:
		return t.System
	}
}

func roleTag(role string) string {
	switch role {
	case "user":
		return " USER "
	case "agent", "assistant", "ai", "model":
		return " ZEU  "
	case "tool_call", "tool_use", "function_call":
		return " TOOL "
	case "tool_result", "tool_response", "function_response":
		return " RSLT "
	case "error":
		return " ERR   "
	default:
		return " SYS  "
	}
}

func (t Theme) RoleLabel(role string) string {
	bg := t.RoleStyle(role).GetForeground()
	badge := lipgloss.NewStyle().
		Background(bg).
		Foreground(lipgloss.Color("#1A1A1A")).
		Bold(true)
	return badge.Render(roleTag(role))
}

func (t Theme) Timestamp() string {
	return t.TimestampStyle.Render(time.Now().Format("15:04:05"))
}

func Now() string {
	return time.Now().Format("15:04:05")
}

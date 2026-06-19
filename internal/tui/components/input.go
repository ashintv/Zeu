package components

import (
	"strings"

	"charm.land/bubbles/v2/key"
	"charm.land/bubbles/v2/textarea"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

type Input struct {
	textarea textarea.Model
}

func NewInput() Input {
	ta := textarea.New()
	ta.Focus()
	ta.Placeholder = "Type your query here..."
	ta.ShowLineNumbers = false
	ta.Prompt = ""
	ta.SetHeight(1)

	ta.KeyMap.InsertNewline = key.NewBinding(
		key.WithKeys("shift+enter", "alt+enter"),
	)

	styles := textarea.DefaultDarkStyles()
	styles.Focused.CursorLine = lipgloss.NewStyle()
	styles.Blurred.CursorLine = lipgloss.NewStyle()
	ta.SetStyles(styles)

	return Input{
		textarea: ta,
	}
}

func predictHeight(ta textarea.Model, msg tea.Msg) int {
	val := ta.Value()
	w := ta.Width()
	if w <= 0 {
		w = 40
	}

	keyMsg, ok := msg.(tea.KeyMsg)
	if !ok {
		return -1
	}

	if keyMsg.String() == "shift+enter" || keyMsg.String() == "alt+enter" {
		val += "\n"
	} else if len(keyMsg.String()) == 1 {
		val += keyMsg.String()
	} else if keyMsg.String() == "backspace" {
		if len(val) > 0 {
			val = val[:len(val)-1]
		}
	}

	rawLines := strings.Split(val, "\n")
	height := 0
	for _, line := range rawLines {
		lineLen := len(line)
		if lineLen == 0 {
			height++
		} else {
			height += (lineLen + w - 1) / w
		}
	}

	if height > 5 {
		height = 5
	}
	if height < 1 {
		height = 1
	}
	return height
}

func (i Input) Update(msg tea.Msg) (Input, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		i.textarea.SetWidth(msg.Width - 4)
	}

	if predHeight := predictHeight(i.textarea, msg); predHeight > 0 {
		i.textarea.SetHeight(predHeight)
	}

	i.textarea, cmd = i.textarea.Update(msg)

	val := i.textarea.Value()
	w := i.textarea.Width()
	if w <= 0 {
		w = 40
	}
	rawLines := strings.Split(val, "\n")
	height := 0
	for _, line := range rawLines {
		lineLen := len(line)
		if lineLen == 0 {
			height++
		} else {
			height += (lineLen + w - 1) / w
		}
	}

	if height > 5 {
		height = 5
	}
	if height < 1 {
		height = 1
	}
	i.textarea.SetHeight(height)

	return i, cmd
}

func (i Input) Value() string {
	return i.textarea.Value()
}

func (i *Input) Reset() {
	i.textarea.Reset()
}

func (i Input) View() string {
	topBottomBorder := lipgloss.Border{
		Top:    "─",
		Bottom: "─",
	}

	borderStyle := lipgloss.NewStyle().
		Border(topBottomBorder, true, false, true, false).
		BorderForeground(lipgloss.Color("#8B9EB7")).
		Blink(true).
		Padding(0, 1)

	return borderStyle.Render(i.textarea.View())
}

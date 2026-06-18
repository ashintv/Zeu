package tui

import (
	"github.com/ashintv/Zeu/internal/tui/components"
	tea "github.com/charmbracelet/bubbletea"
)

type model struct {
	Width int
}

func initialModel() model {
	return model{}
}

func GetInitModel() func() model {
	return initialModel
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
		
	case tea.WindowSizeMsg:
		m.Width = msg.Width
	case tea.KeyMsg:
		switch msg.String() {

		case "ctrl+c", "q":
			return m, tea.Quit

		}
	}

	return m, nil
}

func (m model) View() string {

	return components.Welcome(m.Width)
}

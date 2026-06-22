// Package tui provides terminal user interface components
package tui

import (
	"strings"
	"time"

	"charm.land/bubbles/v2/spinner"
	"charm.land/bubbles/v2/textarea"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"

	"github.com/ashintv/Zeu/internal/harness"

	"github.com/ashintv/Zeu/internal/tui/components"
	"github.com/ashintv/Zeu/internal/tui/handler"
	"github.com/ashintv/Zeu/internal/tui/theme"
	"github.com/ashintv/Zeu/internal/types"
)

type logoTickMsg time.Time

func logoTick() tea.Cmd {
	return tea.Tick(150*time.Millisecond, func(t time.Time) tea.Msg {
		return logoTickMsg(t)
	})
}

type model struct {
	Width         int
	Input         components.Input
	msgTimes      []string
	streamCh      <-chan types.Coversation
	streamingConv *types.Coversation
	chatLoader    spinner.Model
	toolLoader    spinner.Model
	chatLoading   bool
	toolLoading   bool
	logoFrame     int
	startTime     time.Time
	chatHandler   *handler.ChatHandler
	agent         *harness.Agent
}

var loadingStyle = lipgloss.NewStyle().
	Foreground(lipgloss.Color("#8B9EB7"))

func GetNewCli(agent *harness.Agent) model {
	chatSpinner := spinner.New()
	chatSpinner.Spinner = spinner.Jump

	toolSpinner := spinner.New()
	toolSpinner.Spinner = spinner.Meter

	chatHandler := handler.NewChatHandler(agent)

	return model{
		Input:       components.NewInput(),
		msgTimes:    make([]string, 0),
		chatLoader:  chatSpinner,
		toolLoader:  toolSpinner,
		startTime:   time.Now(),
		chatHandler: chatHandler,
	}
}

func (m model) Init() tea.Cmd {
	return tea.Batch(
		m.chatLoader.Tick,
		m.toolLoader.Tick,
		textarea.Blink,
		logoTick(),
	)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {

	case logoTickMsg:
		if time.Since(m.startTime) < 5*time.Second || m.chatLoading || m.toolLoading {
			m.logoFrame++
		} else {
			m.logoFrame = 0
		}
		cmds = append(cmds, logoTick())

	case tea.WindowSizeMsg:
		m.Width = msg.Width

	case tea.KeyMsg:
		switch msg.String() {

		case "ctrl+c":
			return m, tea.Quit

		case "q":
			if strings.TrimSpace(m.Input.Value()) == "" {
				return m, tea.Quit
			}

		case "enter":
			if m.chatLoading {
				return m, nil
			}

			value := strings.TrimSpace(m.Input.Value())
			if value == "" {
				return m, nil
			}

			m.msgTimes = append(m.msgTimes, theme.Now())

			m.Input.Reset()
			m.chatLoading = true
			m.streamingConv = nil

			ch := m.chatHandler.Invoke(value)
			m.streamCh = ch
			m.msgTimes = append(m.msgTimes, theme.Now())

			return m, m.chatHandler.WaitForToken(ch)
		}

	case handler.StreamTokenMsg:
		conv := types.Coversation(msg)
		m.streamingConv = &conv
		// Agent loop automatically updates state, we just track current streaming conversation
		return m, m.chatHandler.WaitForToken(m.streamCh)

	case handler.StreamDoneMsg:
		m.chatLoading = false
		m.streamCh = nil
		m.streamingConv = nil
	}

	var cmd tea.Cmd

	m.chatLoader, cmd = m.chatLoader.Update(msg)
	cmds = append(cmds, cmd)

	m.toolLoader, cmd = m.toolLoader.Update(msg)
	cmds = append(cmds, cmd)

	m.Input, cmd = m.Input.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m model) RenderChat() string {
	th := theme.Default
	var rendered []string

	messages := m.chatHandler.GetMessages()

	// If we're streaming, append the streaming conversation to display
	if m.chatLoading && m.streamingConv != nil {
		messages = append(messages, *m.streamingConv)
	}

	for i, msg := range messages {
		content := msg.Content

		if msg.Role == "user" {
			textLength := 0
			lines := strings.Split(content, "\n")
			for _, line := range lines {
				if len(line) > textLength {
					textLength = len(line)
				}
			}

			lineLength := textLength + 2

			width := m.Width
			if width <= 0 {
				width = 40
			}
			if lineLength > width {
				lineLength = width
			}
			if lineLength < 3 {
				lineLength = 3
			}

			if i > 0 {
				rendered = append(rendered, "")
			}
			sepLine := th.Border.Render(strings.Repeat("─", lineLength))
			rendered = append(rendered, sepLine)

			textColorStyle := lipgloss.NewStyle().Foreground(th.User.GetForeground())
			contentLines := strings.Split(content, "\n")
			var prefixedLines []string
			for idx, cLine := range contentLines {
				styledLine := textColorStyle.Render(cLine)
				if idx == 0 {
					prefixedLines = append(prefixedLines, textColorStyle.Render(">")+" "+styledLine)
				} else {
					prefixedLines = append(prefixedLines, "  "+styledLine)
				}
			}
			rendered = append(rendered, strings.Join(prefixedLines, "\n"))

			rendered = append(rendered, "")
		} else if msg.Role == "agent" || msg.Role == "assistant" || msg.Role == "ai" || msg.Role == "model" {
			styledContent := th.RoleStyle(msg.Role).Render(content)
			rendered = append(rendered, styledContent)
		} else {
			styledContent := th.RoleStyle(msg.Role).Render(content)
			label := th.RoleLabel(msg.Role)
			rendered = append(rendered, label+"  "+styledContent)

			rendered = append(rendered, "")
		}
	}

	return lipgloss.JoinVertical(lipgloss.Left, rendered...)
}

func (m model) View() tea.View {
	var loading string

	if m.chatLoading {
		loading = loadingStyle.Render(
			m.chatLoader.View()+" Working...",
		) + "\n"
	}

	v := tea.NewView(
		lipgloss.JoinVertical(
			lipgloss.Left,
			components.Welcome(m.Width, m.logoFrame),
			m.RenderChat(),
			loading,
			m.Input.View(),
		),
	)
	v.AltScreen = true
	return v
}

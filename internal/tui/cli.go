package tui

import (
	"strings"
	"time"

	"charm.land/bubbles/v2/spinner"
	"charm.land/bubbles/v2/textarea"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"

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
	Width       int
	Input       components.Input
	Messages    []types.Coversation
	msgTimes    []string
	streamCh    <-chan string
	streamBuf   string
	chatLoader  spinner.Model
	toolLoader  spinner.Model
	chatLoading bool
	toolLoading bool
	logoFrame   int
	startTime   time.Time
	chatHandler *handler.ChatHandler
}

var loadingStyle = lipgloss.NewStyle().
	Foreground(lipgloss.Color("#8B9EB7"))

func initialModel() model {
	chatSpinner := spinner.New()
	chatSpinner.Spinner = spinner.Jump

	toolSpinner := spinner.New()
	toolSpinner.Spinner = spinner.Meter

	cfg := handler.DefaultConfig()
	chatHandler := handler.NewChatHandler(cfg)

	return model{
		Input:       components.NewInput(),
		Messages:    make([]types.Coversation, 0),
		msgTimes:    make([]string, 0),
		chatLoader:  chatSpinner,
		toolLoader:  toolSpinner,
		startTime:   time.Now(),
		chatHandler: chatHandler,
	}
}

func GetInitModel() func() model {
	return initialModel
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

			m.Messages = append(m.Messages, types.Coversation{
				Role:    "user",
				Content: value,
			})
			m.msgTimes = append(m.msgTimes, theme.Now())

			m.Input.Reset()
			m.chatLoading = true
			m.streamBuf = ""

			ch := m.chatHandler.FakeAI()
			m.streamCh = ch

			m.Messages = append(m.Messages, types.Coversation{
				Role:    "agent",
				Content: "",
			})
			m.msgTimes = append(m.msgTimes, theme.Now())

			return m, m.chatHandler.WaitForToken(ch)
		}

	case handler.StreamTokenMsg:
		m.streamBuf += string(msg)
		if len(m.Messages) > 0 {
			m.Messages[len(m.Messages)-1].Content = m.streamBuf
		}
		return m, m.chatHandler.WaitForToken(m.streamCh)

	case handler.StreamDoneMsg:
		m.chatLoading = false
		m.streamCh = nil
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

	for i, msg := range m.Messages {
		if msg.Role == "user" {
			textLength := 0
			lines := strings.Split(msg.Content, "\n")
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
			contentLines := strings.Split(msg.Content, "\n")
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
		} else if msg.Role == "agent" || msg.Role == "assistant" || msg.Role == "ai" || msg.Role == "model" {
			content := th.RoleStyle(msg.Role).Render(msg.Content)
			rendered = append(rendered, content)
		} else {
			content := th.RoleStyle(msg.Role).Render(msg.Content)
			label := th.RoleLabel(msg.Role)
			rendered = append(rendered, label+"  "+content)

			rendered = append(rendered, "")
		}
	}

	return lipgloss.JoinVertical(lipgloss.Left, rendered...)
}

func (m model) View() tea.View {
	var loading string

	if m.chatLoading {
		loading = loadingStyle.Render(
			m.chatLoader.View()+" Thinking...",
		) + "\n"
	}

	v := tea.NewView(
		lipgloss.JoinVertical(
			lipgloss.Left,
			"",
			"",
			components.Welcome(m.Width, m.logoFrame),
			m.RenderChat(),
			loading,
			m.Input.View(),
		),
	)
	v.AltScreen = true
	return v
}

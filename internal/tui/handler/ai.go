package handler

import (
	"strings"
	"time"

	tea "charm.land/bubbletea/v2"
)

type StreamTokenMsg string
type StreamDoneMsg struct{}

type ChatHandler struct {
	Config ModelConfig
}

func NewChatHandler(cfg ModelConfig) *ChatHandler {
	return &ChatHandler{
		Config: cfg,
	}
}

func (h *ChatHandler) FakeAI() <-chan string {
	ch := make(chan string)
	go func() {
		defer close(ch)
		response := "Hello from Zeu! 👋 I'm your AI assistant, ready to help you with anything you need. How can I assist you today?"
		words := strings.Fields(response)
		for i, word := range words {
			time.Sleep(60 * time.Millisecond)
			if i > 0 {
				ch <- " "
			}
			ch <- word
		}
	}()
	return ch
}

func (h *ChatHandler) WaitForToken(ch <-chan string) tea.Cmd {
	return func() tea.Msg {
		token, ok := <-ch
		if !ok {
			return StreamDoneMsg{}
		}
		return StreamTokenMsg(token)
	}
}

package handler

import (
	"context"

	tea "charm.land/bubbletea/v2"
	"github.com/ashintv/Zeu/internal/ai"
	"github.com/ashintv/Zeu/internal/harness"
	"github.com/ashintv/Zeu/internal/types"
)

type StreamTokenMsg types.Coversation
type StreamDoneMsg struct{}

type ChatHandler struct {
	agent *harness.Agent
}

func NewChatHandler(agent *harness.Agent) *ChatHandler {
	return &ChatHandler{
		agent: agent,
	}
}

func (h *ChatHandler) UpdateChatModel(ai *ai.AI) {
	h.agent.Ai = ai
}

func (h *ChatHandler) Invoke(prompt string) <-chan types.Coversation {
	stream := h.agent.Invoke(context.Background(), prompt)
	return stream
}

func (h *ChatHandler) WaitForToken(ch <-chan types.Coversation) tea.Cmd {
	return func() tea.Msg {
		conv, ok := <-ch
		if !ok {
			return StreamDoneMsg{}
		}
		return StreamTokenMsg(conv)
	}
}

func (h *ChatHandler) GetHistory() []types.Coversation {
	return h.agent.GetState()
}

func (h *ChatHandler) GetMessages() []types.Coversation {
	return h.agent.GetState()
}

func (h *ChatHandler) SetMessages(messages []types.Coversation) {
	h.agent.SetState(messages)
}

func (h *ChatHandler) AddMessage(msg types.Coversation) {
	messages := h.agent.GetState()
	messages = append(messages, msg)
	h.agent.SetState(messages)
}

func (h *ChatHandler) GetAgent() *harness.Agent {
	return h.agent
}

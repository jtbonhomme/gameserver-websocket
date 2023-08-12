package main

import (
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const (
	focusMethod  int = 0
	focusPayload int = 1
)

var (
	focusedStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#FAFAFA"))
	noStyle      = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
)

type RPC struct {
	Method  string
	Payload string
}

type model struct {
	method  textinput.Model
	payload textinput.Model
	focus   int
	err     error
	RPC     RPC
	res     chan RPC
}

func initialModel() model {
	tiMethod := textinput.New()
	tiMethod.Placeholder = "method"
	tiMethod.Focus()
	tiMethod.CharLimit = 156
	tiMethod.Width = 20
	tiMethod.PromptStyle = focusedStyle
	tiMethod.TextStyle = focusedStyle

	tiPayload := textinput.New()
	tiPayload.Placeholder = `{"data": "payload"}`
	tiPayload.Blur()
	tiPayload.CharLimit = 156
	tiPayload.Width = 20
	tiPayload.PromptStyle = noStyle
	tiPayload.TextStyle = noStyle

	return model{
		method:  tiMethod,
		payload: tiPayload,
		focus:   0,
		err:     nil,
	}
}

func (m model) Init() tea.Cmd {
	return textinput.Blink
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		// Exit
		case tea.KeyCtrlC, tea.KeyEsc:
			os.Exit(1)
		// Enter
		case tea.KeyEnter:
			if m.focus == focusMethod {
				m.focus = focusPayload
				m.method.Blur()
				m.method.PromptStyle = noStyle
				m.method.TextStyle = noStyle
				m.payload.PromptStyle = focusedStyle
				m.payload.TextStyle = focusedStyle
				return m, m.payload.Focus()
			}

			if m.focus == focusPayload {
				m.payload.Blur()
				m.payload.PromptStyle = noStyle
				m.payload.TextStyle = noStyle
				return m, tea.Quit
			}
			m.res <- m.RPC
			return m, tea.Quit
		default:
			switch m.focus {
			case focusMethod:
				m.RPC.Method += msg.String()
			case focusPayload:
				m.RPC.Payload += msg.String()
			}
		}
	}

	// Handle character input and blinking
	cmd := m.updateInputs(msg)

	return m, cmd
}

func (m *model) updateInputs(msg tea.Msg) tea.Cmd {
	cmds := make([]tea.Cmd, 2)

	// Only text inputs with Focus() set will respond, so it's safe to simply
	// update all of them here without any further logic.
	m.method, cmds[focusMethod] = m.method.Update(msg)
	m.payload, cmds[focusPayload] = m.payload.Update(msg)

	return tea.Batch(cmds...)
}

func (m model) View() string {
	return fmt.Sprintf(
		"Execute Remote Procedure Calls from REPL:\n\n%s\n%s\n\n%s",
		m.method.View(),
		m.payload.View(),
		"(esc or ctrl+c to quit)",
	) + "\n\n"
}

package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"

	centrifuge "github.com/centrifugal/centrifuge-go"
	"github.com/rs/zerolog"
)

func newClient(log *zerolog.Logger) *centrifuge.Client {
	wsURL := "ws://localhost:8000/connection/websocket"
	c := centrifuge.NewJsonClient(wsURL, centrifuge.Config{
		Name:    "centrifuge-go",
		Data:    []byte(`{"name":"totoro"}`),
		Version: "0.0.1",
	})

	c.OnConnecting(func(_ centrifuge.ConnectingEvent) {
		log.Info().Msg("Connecting")
	})
	c.OnConnected(func(_ centrifuge.ConnectedEvent) {
		log.Info().Msg("Connected")
	})
	c.OnDisconnected(func(_ centrifuge.DisconnectedEvent) {
		log.Info().Msg("Disconnected")
	})
	c.OnError(func(e centrifuge.ErrorEvent) {
		log.Info().Msgf("error: %s", e.Error.Error())
	})
	c.OnMessage(func(e centrifuge.MessageEvent) {
		log.Info().Msgf("Message received from server %s", string(e.Data))

		// When issue blocking requests from inside event handler we must use
		// a goroutine. Otherwise, connection read loop will be blocked.
		go func() {
			result, err := c.RPC(context.Background(), "method", []byte(`{"action":"eat"}`))
			if err != nil {
				log.Info().Msgf("%s", err.Error())
				return
			}
			log.Printf("RPC result 2: %s", string(result.Data))
		}()
	})
	return c
}

func main() {
	var err error

	// Init logger
	zerolog.SetGlobalLevel(zerolog.DebugLevel)
	output := zerolog.ConsoleWriter{
		Out:           os.Stderr,
		TimeFormat:    time.RFC3339,
		FormatMessage: func(i interface{}) string { return fmt.Sprintf("[client] %s", i) },
	}
	log := zerolog.New(output).With().Timestamp().Logger()
	log.Info().Msg("start client")
	c := newClient(&log)
	defer c.Close()

	err = c.Connect()
	if err != nil {
		log.Panic().Msgf("connect error: %s", err.Error())
	}

	p := tea.NewProgram(initialModel())
	msg, err := p.Run()
	if err != nil {
		log.Panic().Err(err)
	}
	rpc := msg.(model).RPC

	result, err := c.RPC(context.Background(), rpc.Method, []byte(rpc.Payload))
	if err != nil {
		log.Error().Msgf("error executing RPC: %s", err.Error())
	}
	log.Printf("RPC result: %s", string(result.Data))

	log.Info().Msg("exit")
}

const (
	focusMethod  int = 0
	focusPayload int = 1
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

	tiPayload := textinput.New()
	tiPayload.Placeholder = `{"data": "payload"}`
	tiPayload.Blur()
	tiPayload.CharLimit = 156
	tiPayload.Width = 20

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
				return m, m.payload.Focus()
			}

			if m.focus == focusPayload {
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

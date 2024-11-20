package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/rs/zerolog/log"
)

// Messages
type (
	msgCountDown struct{}
	msgTick      struct{}
	msgError     error
	msgJoke      string
)

// model is the bubbletea model that contains the render
// and update logic
type model struct {
	width, height int
	joke          string
	err           error
	countdown     int
}

// Init sets up the initial command
func (m *model) Init() tea.Cmd {
	return tea.Batch(tick(), fetchJoke())
}

// Update handles messages and updates the model
func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	log.Debug().Type("msg", msg).Msg("new message")

	switch msg := msg.(type) {
	case msgTick:
		m.countdown = 5 // Reset countdown to 5 seconds
		return m, tea.Batch(tick(), fetchJoke(), countdownTick())

	case msgCountDown:
		if m.countdown > 0 {
			m.countdown--
			return m, countdownTick()
		}
		return m, nil

	case msgJoke:
		m.joke = string(msg)
		m.err = nil
		return m, nil

	case msgError:
		// Update the error in the model
		m.err = msg
		return m, nil

	case tea.KeyMsg:
		if msg.String() == "q" || msg.String() == "esc" || msg.String() == "ctrl-c" {
			return m, tea.Quit
		}
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	}

	return m, nil
}

// View renders the UI
func (m *model) View() string {
	const (
		HeaderHeight = 3
		FooterHeight = 1
	)

	header := StyleHeader.Width(m.width).Render("Random Dad Joke")

	footerRender := StyleFooter.Width(m.width).Render

	footer := footerRender(fmt.Sprintf("Next joke in %d seconds | Press q or ESC to quit", m.countdown))
	if m.countdown == 0 {
		footer = footerRender("Next joke loading... | Press q or ESC to quit")
	}

	contentRender := lipgloss.NewStyle().
		Width(m.width).
		Height(m.height-HeaderHeight-FooterHeight). // accommodate header and footer
		Align(lipgloss.Center, lipgloss.Center).
		Render

	content := contentRender(m.joke)
	if m.err != nil {
		content = contentRender(m.err.Error())
	}

	return lipgloss.JoinVertical(lipgloss.Top, header, content, footer)
}

// tick creates a message after 5 seconds
func tick() tea.Cmd {
	return tea.Tick(5*time.Second, func(time.Time) tea.Msg {
		return msgTick{}
	})
}

type Joke struct {
	Joke string `json:"joke"`
}

func fetchJoke() tea.Cmd {
	return func() tea.Msg {
		log.Debug().Msg("fetching new joke")
		req, err := http.NewRequest("GET", "https://icanhazdadjoke.com/", nil)
		if err != nil {
			return msgError(err)
		}

		req.Header.Set("Accept", "application/json")

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			return msgError(err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			log.Error().Int("status", resp.StatusCode).Msg("got bad status code while fetching joke")
			return msgError(fmt.Errorf("bad status: %s", resp.Status))
		}

		data := Joke{}

		err = json.NewDecoder(resp.Body).Decode(&data)
		if err != nil {
			log.Err(err).Msg("error decoding json")
			return msgError(err)
		}

		log.Debug().Str("joke", data.Joke).Msg("got data joke")
		return msgJoke(data.Joke)
	}
}

func countdownTick() tea.Cmd {
	return tea.Tick(1*time.Second, func(time.Time) tea.Msg {
		return msgCountDown{}
	})
}

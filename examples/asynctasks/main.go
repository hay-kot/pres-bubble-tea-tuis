package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	f, err := os.OpenFile("./devlog.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		panic(err)
	}

	log.Logger = log.Output(zerolog.ConsoleWriter{
		Out: f,
	})

	log.Info().Msg("starting joke generator")
	err = run()
	if err != nil {
		log.Fatal().Err(err).Msg("failed to run pageviews demo")
	}

	log.Info().Msg("closing joke generator")
}

func run() error {
	// Initialize the Bubble Tea program
	p := tea.NewProgram(&model{}, tea.WithAltScreen())
	_, err := p.Run()
	return err
}

type model struct {
	width, height int
	joke          string
	err           error
	countdown     int
}

type (
	countdownTickMsg struct{}
	tickMsg          struct{}
	errorMsg         error
	jokeMsg          string
)

// Init sets up the initial command
func (m *model) Init() tea.Cmd {
	return tea.Batch(tick(), fetchJoke())
}

// Update handles messages and updates the model
func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	log.Debug().Type("msg", msg).Msg("new message")

	switch msg := msg.(type) {
	case tickMsg:
		m.countdown = 5 // Reset countdown to 5 seconds
		return m, tea.Batch(tick(), fetchJoke(), countdownTick())

	case countdownTickMsg:
		if m.countdown > 0 {
			m.countdown--
			return m, countdownTick()
		}
		return m, nil

	case jokeMsg:
		m.joke = string(msg)
		m.err = nil
		return m, nil

	case errorMsg:
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

	header := lipgloss.NewStyle().
		Align(lipgloss.Center).
		Width(m.width).
		Background(lipgloss.Color("#14532d")).
		Foreground(lipgloss.Color("#dcfce7")).
		PaddingTop(1).
		PaddingBottom(1).
		Render("Random Dad Joke")

	footerRender := lipgloss.NewStyle().
		Align(lipgloss.Center).
		Background(lipgloss.Color("#14532d")).
		Foreground(lipgloss.Color("#dcfce7")).
		Width(m.width).
		Render

	footer := footerRender(fmt.Sprintf("Next joke in %d seconds | Press q or ESC to quit", m.countdown))
	if m.countdown == 0 {
		footer = footerRender("Next joke loading... | Press q or ESC to quit")
	}

	contentRender := lipgloss.NewStyle().
		Width(m.width).
		// accommodate header and footer
		Height(m.height-HeaderHeight-FooterHeight).
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
		return tickMsg{}
	})
}

// fetchJoke fetches a random dad joke from the API
func fetchJoke() tea.Cmd {
	return func() tea.Msg {
		req, err := http.NewRequest("GET", "https://icanhazdadjoke.com/", nil)
		if err != nil {
			return errorMsg(err)
		}

		// Set the Accept header to request JSON response
		req.Header.Set("Accept", "application/json")

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			return errorMsg(err)
		}
		defer resp.Body.Close()

		// Check for successful status code
		if resp.StatusCode != http.StatusOK {
			return errorMsg(fmt.Errorf("bad status: %s", resp.Status))
		}

		// Parse the JSON response
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return errorMsg(err)
		}

		var data struct {
			Joke string `json:"joke"`
		}
		if err := json.Unmarshal(body, &data); err != nil {
			return errorMsg(err)
		}
		return jokeMsg(data.Joke)
	}
}

func countdownTick() tea.Cmd {
	return tea.Tick(1*time.Second, func(time.Time) tea.Msg {
		return countdownTickMsg{}
	})
}

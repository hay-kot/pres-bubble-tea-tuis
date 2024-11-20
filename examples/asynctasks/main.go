package main

import (
	"os"

	tea "github.com/charmbracelet/bubbletea"
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

package main

import "github.com/charmbracelet/lipgloss"

var StyleHeader = lipgloss.NewStyle().
	Align(lipgloss.Center).
	Background(lipgloss.Color("#14532d")).
	Foreground(lipgloss.Color("#dcfce7")).
	PaddingTop(1).
	PaddingBottom(1)

var StyleFooter = lipgloss.NewStyle().
	Align(lipgloss.Center).
	Background(lipgloss.Color("#14532d")).
	Foreground(lipgloss.Color("#dcfce7"))

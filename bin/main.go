package main

// A standalone Bubble Tea client for the Star Wars ASCII animation.

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/taigrr/ssh-wars/asciimation"
)

func main() {
	renderer := lipgloss.DefaultRenderer()
	model := asciimation.NewDefaultModel(renderer)

	program := tea.NewProgram(model)
	if _, err := program.Run(); err != nil {
		fmt.Printf("Error: %v", err)
		os.Exit(1)
	}
}

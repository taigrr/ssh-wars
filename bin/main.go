package main

// A standalone Bubble Tea client for the Star Wars ASCII animation.

import (
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/progress"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/taigrr/ssh-wars/asciimation"
)

func main() {
	renderer := lipgloss.DefaultRenderer()
	m := asciimation.New(renderer)
	m.Progress = asciimation.ModelProg{
		Progress: progress.New(progress.WithSolidFill("#174ea6")),
		MaxWidth: 65,
		Padding:  2,
	}
	m.Help = asciimation.NewHelpModel(renderer)
	m.Speed = 15

	p := tea.NewProgram(m)
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error: %v", err)
		os.Exit(1)
	}
}

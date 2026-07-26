package main

// A standalone Bubble Tea client for the Star Wars ASCII animation.

import (
	"fmt"
	"os"

	tea "charm.land/bubbletea/v2"
	"github.com/taigrr/ssh-wars/asciimation"
)

func main() {
	model := asciimation.NewDefaultModel()

	program := tea.NewProgram(model)
	if _, err := program.Run(); err != nil {
		fmt.Printf("Error: %v", err)
		os.Exit(1)
	}
}

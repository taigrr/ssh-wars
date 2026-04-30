package main

// An example Bubble Tea server. This will put an ssh session into alt screen
// and continually print up to date terminal information.

import (
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/progress"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/muesli/termenv"
	"github.com/taigrr/ssh-wars/asciimation"
)

func main() {
	m := asciimation.New()
	profile := termenv.ColorProfile()
	m.Progress = asciimation.ModelProg{Progress: progress.New(
		progress.WithSolidFill("#174ea6"),
		progress.WithColorProfile(profile),
	),
		MaxWidth: 65,
		Padding:  2,
	}
	m.Help = asciimation.NewHelpModel()
	m = m.UpdateRenderer(profile)
	m.Speed = 15

	p := tea.NewProgram(m)
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error: %v", err)
		os.Exit(1)
	}
}

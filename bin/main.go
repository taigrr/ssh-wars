package main

// An example Bubble Tea server. This will put an ssh session into alt screen
// and continually print up to date terminal information.

import (
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/progress"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/taigrr/ssh-wars/asciimation"
)

func main() {
	m := asciimation.New()
	df := lipgloss.NewDoeFoot()
	m.Progress = asciimation.ModelProg{Progress: progress.New(progress.WithSolidFill("#174ea6"),
		progress.WithColorProfile(df.Profile)),
		MaxWidth: 65,
		Padding:  2}
	m.Help = asciimation.NewHelpModel()
	m = m.UpdateDoeFoot(df)
	m.Speed = 15

	p := tea.NewProgram(m)
	if err := p.Start(); err != nil {
		fmt.Printf("Error: %v", err)
		os.Exit(1)
	}

}
